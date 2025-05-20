-- 000003_add_cognito_sub.up.sql
ALTER TABLE users ADD COLUMN cognito_sub VARCHAR UNIQUE NOT NULL;
