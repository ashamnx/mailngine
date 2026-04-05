package template

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/mailngine/mailngine/internal/db/sqlcdb"
)

// ErrTemplateNotFound is returned when a template cannot be found.
var ErrTemplateNotFound = errors.New("template not found")

// ErrTemplateNameConflict is returned when a template name already exists for the org.
var ErrTemplateNameConflict = errors.New("template name already exists")

// Service handles template operations.
type Service struct {
	db      *pgxpool.Pool
	queries *sqlcdb.Queries
	logger  zerolog.Logger
}

// NewService creates a new template Service.
func NewService(db *pgxpool.Pool, logger zerolog.Logger) *Service {
	return &Service{
		db:      db,
		queries: sqlcdb.New(db),
		logger:  logger,
	}
}

// Create creates a new template for the given organization.
func (s *Service) Create(ctx context.Context, orgID uuid.UUID, name, subject, htmlBody, textBody string, variables []string) (*sqlcdb.Template, error) {
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("marshal variables: %w", err)
	}

	tmpl, err := s.queries.CreateTemplate(ctx, sqlcdb.CreateTemplateParams{
		OrgID:     orgID,
		Name:      name,
		Subject:   subject,
		HtmlBody:  htmlBody,
		TextBody:  pgtype.Text{String: textBody, Valid: textBody != ""},
		Variables: variablesJSON,
	})
	if err != nil {
		if isUniqueViolation(err) {
			return nil, ErrTemplateNameConflict
		}
		return nil, fmt.Errorf("create template: %w", err)
	}

	return &tmpl, nil
}

// Get retrieves a template by ID and org.
func (s *Service) Get(ctx context.Context, orgID, templateID uuid.UUID) (*sqlcdb.Template, error) {
	tmpl, err := s.queries.GetTemplate(ctx, sqlcdb.GetTemplateParams{
		ID:    templateID,
		OrgID: orgID,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		return nil, fmt.Errorf("get template: %w", err)
	}

	return &tmpl, nil
}

// List returns all templates for an organization.
func (s *Service) List(ctx context.Context, orgID uuid.UUID) ([]sqlcdb.Template, error) {
	return s.queries.ListTemplatesByOrg(ctx, orgID)
}

// Update updates an existing template.
func (s *Service) Update(ctx context.Context, orgID, templateID uuid.UUID, name, subject, htmlBody, textBody string, variables []string) (*sqlcdb.Template, error) {
	variablesJSON, err := json.Marshal(variables)
	if err != nil {
		return nil, fmt.Errorf("marshal variables: %w", err)
	}

	tmpl, err := s.queries.UpdateTemplate(ctx, sqlcdb.UpdateTemplateParams{
		ID:        templateID,
		OrgID:     orgID,
		Name:      name,
		Subject:   subject,
		HtmlBody:  htmlBody,
		TextBody:  pgtype.Text{String: textBody, Valid: textBody != ""},
		Variables: variablesJSON,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrTemplateNotFound
		}
		if isUniqueViolation(err) {
			return nil, ErrTemplateNameConflict
		}
		return nil, fmt.Errorf("update template: %w", err)
	}

	return &tmpl, nil
}

// Delete removes a template.
func (s *Service) Delete(ctx context.Context, orgID, templateID uuid.UUID) error {
	return s.queries.DeleteTemplate(ctx, sqlcdb.DeleteTemplateParams{
		ID:    templateID,
		OrgID: orgID,
	})
}

// Render loads a template and replaces {{variable}} placeholders with values from the data map.
// It returns the rendered subject, HTML body, and text body.
func (s *Service) Render(ctx context.Context, orgID, templateID uuid.UUID, data map[string]string) (subject, htmlBody, text string, err error) {
	tmpl, err := s.Get(ctx, orgID, templateID)
	if err != nil {
		return "", "", "", err
	}

	subject = replacePlaceholders(tmpl.Subject, data, false)
	htmlBody = replacePlaceholders(tmpl.HtmlBody, data, true)
	text = replacePlaceholders(tmpl.TextBody.String, data, false)

	return subject, htmlBody, text, nil
}

// replacePlaceholders replaces {{key}} placeholders in the input string with values from data.
// When escapeHTML is true, values are HTML-escaped to prevent XSS in HTML email bodies.
func replacePlaceholders(input string, data map[string]string, escapeHTML bool) string {
	result := input
	for key, value := range data {
		placeholder := "{{" + key + "}}"
		v := value
		if escapeHTML {
			v = html.EscapeString(v)
		}
		result = strings.ReplaceAll(result, placeholder, v)
	}
	return result
}

// isUniqueViolation checks if the error is a PostgreSQL unique constraint violation (23505).
func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	// pgx wraps PostgreSQL errors as *pgconn.PgError.
	var pgErr interface{ SQLState() string }
	if errors.As(err, &pgErr) {
		return pgErr.SQLState() == "23505"
	}
	return false
}
