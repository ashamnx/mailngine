// Package team implements team and organization management for Hello Mail.
package team

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	sqlcdb "github.com/hellomail/hellomail/internal/db/sqlcdb"
)

// Valid roles for organization members.
var validRoles = map[string]bool{
	"viewer": true,
	"member": true,
	"admin":  true,
	"owner":  true,
}

// Sentinel errors returned by the team service.
var (
	ErrUserNotRegistered = errors.New("user not registered")
	ErrAlreadyMember     = errors.New("user is already a member of this organization")
	ErrMemberNotFound    = errors.New("member not found")
	ErrInvalidRole       = errors.New("invalid role: must be one of viewer, member, admin, owner")
	ErrLastOwner         = errors.New("cannot remove or change the role of the last owner")
	ErrCannotRemoveSelf  = errors.New("cannot remove yourself; use the leave endpoint instead")
	ErrForbidden         = errors.New("insufficient permissions")
	ErrOrgNotFound       = errors.New("organization not found")
	ErrNameRequired      = errors.New("name is required")
)

// Service provides team and organization management operations.
type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	logger  zerolog.Logger
}

// NewService creates a new team Service with the given dependencies.
func NewService(db *pgxpool.Pool, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		logger:  logger.With().Str("component", "team_service").Logger(),
	}
}

// ListMembers returns all members of the given organization.
func (s *Service) ListMembers(ctx context.Context, orgID uuid.UUID) ([]sqlcdb.ListOrgMembersRow, error) {
	members, err := s.queries.ListOrgMembers(ctx, orgID)
	if err != nil {
		return nil, fmt.Errorf("list org members: %w", err)
	}
	return members, nil
}

// InviteMember adds an existing user to the organization with the given role.
// The inviter must have admin or owner role.
func (s *Service) InviteMember(ctx context.Context, orgID, inviterID uuid.UUID, email, role string) error {
	role = strings.ToLower(strings.TrimSpace(role))
	if !validRoles[role] {
		return ErrInvalidRole
	}

	// Look up the user by email.
	user, err := s.queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrUserNotRegistered
		}
		return fmt.Errorf("lookup user by email: %w", err)
	}

	// Check if user is already a member.
	_, err = s.queries.GetOrgMember(ctx, sqlcdb.GetOrgMemberParams{
		OrgID:  orgID,
		UserID: user.ID,
	})
	if err == nil {
		return ErrAlreadyMember
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("check existing membership: %w", err)
	}

	// Add the user as a member.
	_, err = s.queries.AddOrgMember(ctx, sqlcdb.AddOrgMemberParams{
		OrgID:     orgID,
		UserID:    user.ID,
		Role:      role,
		InvitedBy: pgtype.UUID{Bytes: inviterID, Valid: true},
	})
	if err != nil {
		return fmt.Errorf("add org member: %w", err)
	}

	return nil
}

// UpdateRole changes the role of a member in the organization.
// It validates the role and prevents changing the last owner's role.
func (s *Service) UpdateRole(ctx context.Context, orgID, memberUserID uuid.UUID, role string) error {
	role = strings.ToLower(strings.TrimSpace(role))
	if !validRoles[role] {
		return ErrInvalidRole
	}

	// Verify the target member exists.
	member, err := s.queries.GetOrgMember(ctx, sqlcdb.GetOrgMemberParams{
		OrgID:  orgID,
		UserID: memberUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMemberNotFound
		}
		return fmt.Errorf("get org member: %w", err)
	}

	// If the member is currently an owner and the new role is not owner,
	// ensure they are not the last owner.
	if member.Role == "owner" && role != "owner" {
		if err := s.ensureNotLastOwner(ctx, orgID); err != nil {
			return err
		}
	}

	err = s.queries.UpdateOrgMemberRole(ctx, sqlcdb.UpdateOrgMemberRoleParams{
		OrgID:  orgID,
		UserID: memberUserID,
		Role:   role,
	})
	if err != nil {
		return fmt.Errorf("update org member role: %w", err)
	}

	return nil
}

// RemoveMember removes a member from the organization.
// It prevents removing the last owner and prevents self-removal.
func (s *Service) RemoveMember(ctx context.Context, orgID, memberUserID uuid.UUID) error {
	// Verify the target member exists.
	member, err := s.queries.GetOrgMember(ctx, sqlcdb.GetOrgMemberParams{
		OrgID:  orgID,
		UserID: memberUserID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrMemberNotFound
		}
		return fmt.Errorf("get org member: %w", err)
	}

	// Prevent removing the last owner.
	if member.Role == "owner" {
		if err := s.ensureNotLastOwner(ctx, orgID); err != nil {
			return err
		}
	}

	err = s.queries.RemoveOrgMember(ctx, sqlcdb.RemoveOrgMemberParams{
		OrgID:  orgID,
		UserID: memberUserID,
	})
	if err != nil {
		return fmt.Errorf("remove org member: %w", err)
	}

	return nil
}

// GetOrg retrieves an organization by ID.
func (s *Service) GetOrg(ctx context.Context, orgID uuid.UUID) (*sqlcdb.Organization, error) {
	org, err := s.queries.GetOrganization(ctx, orgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("get organization: %w", err)
	}
	return &org, nil
}

// UpdateOrg updates the organization name.
func (s *Service) UpdateOrg(ctx context.Context, orgID uuid.UUID, name string) (*sqlcdb.Organization, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, ErrNameRequired
	}

	org, err := s.queries.UpdateOrganizationName(ctx, sqlcdb.UpdateOrganizationNameParams{
		ID:   orgID,
		Name: name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrOrgNotFound
		}
		return nil, fmt.Errorf("update organization name: %w", err)
	}
	return &org, nil
}

// ensureNotLastOwner verifies that the org has more than one owner.
// Returns ErrLastOwner if only one owner remains.
func (s *Service) ensureNotLastOwner(ctx context.Context, orgID uuid.UUID) error {
	members, err := s.queries.ListOrgMembers(ctx, orgID)
	if err != nil {
		return fmt.Errorf("list org members for owner check: %w", err)
	}

	ownerCount := 0
	for _, m := range members {
		if m.Role == "owner" {
			ownerCount++
		}
	}

	if ownerCount <= 1 {
		return ErrLastOwner
	}
	return nil
}
