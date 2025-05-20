-- 000004_cognito_sub_primary_key.up.sql
-- Step 1: Drop foreign key constraints referencing users.user_id
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_user_id_fkey;
ALTER TABLE contacts DROP CONSTRAINT IF EXISTS contacts_company_id_fkey;
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_user_id_fkey;
ALTER TABLE sales_processes DROP CONSTRAINT IF EXISTS sales_processes_user_id_fkey;

-- Step 2: Add cognito_sub reference in tables before dropping user_id
ALTER TABLE companies ADD COLUMN IF NOT EXISTS user_cognito_sub VARCHAR;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS user_cognito_sub VARCHAR;
ALTER TABLE sales_processes ADD COLUMN IF NOT EXISTS user_cognito_sub VARCHAR;

-- Step 3: Update tables to use cognito_sub before dropping user_id
UPDATE companies SET user_cognito_sub = (SELECT u.cognito_sub FROM users u WHERE u.user_id = companies.user_id);
UPDATE projects SET user_cognito_sub = (SELECT u.cognito_sub FROM users u WHERE u.user_id = projects.user_id);
UPDATE sales_processes SET user_cognito_sub = (SELECT u.cognito_sub FROM users u WHERE u.user_id = sales_processes.user_id);

-- Step 4: Make cognito_sub the primary key
ALTER TABLE users DROP CONSTRAINT IF EXISTS users_pkey;
ALTER TABLE users ADD PRIMARY KEY (cognito_sub);
ALTER TABLE users DROP COLUMN IF EXISTS user_id;

-- Step 5: Rename and enforce foreign key constraints
ALTER TABLE companies DROP COLUMN IF EXISTS user_id;
ALTER TABLE companies RENAME COLUMN user_cognito_sub TO user_id;
ALTER TABLE companies ADD CONSTRAINT companies_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (cognito_sub) ON DELETE CASCADE;

ALTER TABLE projects DROP COLUMN IF EXISTS user_id;
ALTER TABLE projects RENAME COLUMN user_cognito_sub TO user_id;
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (cognito_sub) ON DELETE CASCADE;

ALTER TABLE sales_processes DROP COLUMN IF EXISTS user_id;
ALTER TABLE sales_processes RENAME COLUMN user_cognito_sub TO user_id;
ALTER TABLE sales_processes ADD CONSTRAINT sales_processes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (cognito_sub) ON DELETE CASCADE;

-- Reconnect contacts to companies after the changes
ALTER TABLE contacts ADD CONSTRAINT contacts_company_id_fkey FOREIGN KEY (company_id) REFERENCES companies (company_id) ON DELETE CASCADE;