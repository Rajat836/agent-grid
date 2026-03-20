CREATE TABLE IF NOT EXISTS service_team_roles (
    id BIGSERIAL PRIMARY KEY,
    service_id BIGINT NOT NULL,
    team_id BIGINT NOT NULL,
    role VARCHAR(64) NOT NULL,        -- e.g. 'owner', 'contributor', 'on-call'
    assigned_at TIMESTAMP DEFAULT NOW(),
    revoked_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_service_team_roles_service_id ON service_team_roles(service_id);
CREATE INDEX idx_service_team_roles_team_id ON service_team_roles(team_id);
CREATE INDEX idx_service_team_roles_role ON service_team_roles(role);
