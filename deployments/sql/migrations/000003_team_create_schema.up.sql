CREATE SCHEMA IF NOT EXISTS team;

CREATE TABLE IF NOT EXISTS team.team (
    id VARCHAR(255) PRIMARY KEY,
    tenant_id VARCHAR(255) NOT NULL,
    "name" VARCHAR(255) NOT NULL,
    FOREIGN KEY (tenant_id) REFERENCES tenant.tenant(id),
);
