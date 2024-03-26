CREATE SCHEMA IF NOT EXISTS tenant;

CREATE TABLE IF NOT EXISTS tenant.tenant (
    id VARCHAR(255) PRIMARY KEY,
    "name" VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS tenant.tenant_user_access (
    tenant_id VARCHAR(255) NOT NULL,
    user_account_id VARCHAR(255) NOT NULL,
    access_level VARCHAR(255) NOT NULL,
    PRIMARY KEY (tenant_id, user_account_id),
    FOREIGN KEY (tenant_id) REFERENCES tenant.tenant(id),
    FOREIGN KEY (user_account_id) REFERENCES user_account.user_account(id)
);