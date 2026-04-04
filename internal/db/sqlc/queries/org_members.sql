-- name: AddOrgMember :one
INSERT INTO org_members (org_id, user_id, role, invited_by)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetOrgMember :one
SELECT * FROM org_members WHERE org_id = $1 AND user_id = $2;

-- name: ListOrgMembers :many
SELECT om.*, u.email, u.name, u.avatar_url
FROM org_members om
JOIN users u ON u.id = om.user_id
WHERE om.org_id = $1
ORDER BY om.joined_at;

-- name: UpdateOrgMemberRole :exec
UPDATE org_members SET role = $3 WHERE org_id = $1 AND user_id = $2;

-- name: RemoveOrgMember :exec
DELETE FROM org_members WHERE org_id = $1 AND user_id = $2;

-- name: ListUserOrgs :many
SELECT o.*, om.role
FROM organizations o
JOIN org_members om ON om.org_id = o.id
WHERE om.user_id = $1
ORDER BY o.name;
