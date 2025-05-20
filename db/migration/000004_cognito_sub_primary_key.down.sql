-- 000004_cognito_sub_primary_key.down.sql
-- Step 1: Remove foreign key constraints pointing to cognito_sub
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_user_id_fkey;
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_user_id_fkey;
ALTER TABLE sales_processes DROP CONSTRAINT IF EXISTS sales_processes_user_id_fkey;
ALTER TABLE contacts DROP CONSTRAINT IF EXISTS contacts_company_id_fkey;

-- Step 2: Add back user_id as a SERIAL column to users
ALTER TABLE users DROP CONSTRAINT users_pkey;
ALTER TABLE users ADD COLUMN user_id SERIAL PRIMARY KEY;

-- Step 3: Revert tables to use user_id
ALTER TABLE companies ADD COLUMN IF NOT EXISTS user_id_old INT;
ALTER TABLE projects ADD COLUMN IF NOT EXISTS user_id_old INT;
ALTER TABLE sales_processes ADD COLUMN IF NOT EXISTS user_id_old INT;

UPDATE companies SET user_id_old = (SELECT u.user_id FROM users u WHERE u.cognito_sub = companies.user_id);
UPDATE projects SET user_id_old = (SELECT u.user_id FROM users u WHERE u.cognito_sub = projects.user_id);
UPDATE sales_processes SET user_id_old = (SELECT u.user_id FROM users u WHERE u.cognito_sub = sales_processes.user_id);

-- Step 4: Drop cognito_sub foreign keys and restore user_id foreign keys
ALTER TABLE companies DROP COLUMN user_id;
ALTER TABLE companies RENAME COLUMN user_id_old TO user_id;
ALTER TABLE companies ADD CONSTRAINT companies_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

ALTER TABLE projects DROP COLUMN user_id;
ALTER TABLE projects RENAME COLUMN user_id_old TO user_id;
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

ALTER TABLE sales_processes DROP COLUMN user_id;
ALTER TABLE sales_processes RENAME COLUMN user_id_old TO user_id;
ALTER TABLE sales_processes ADD CONSTRAINT sales_processes_user_id_fkey FOREIGN KEY (user_id) REFERENCES users (user_id) ON DELETE CASCADE;

-- Reconnect contacts to companies
ALTER TABLE contacts ADD CONSTRAINT contacts_company_id_fkey FOREIGN KEY (company_id) REFERENCES companies (company_id) ON DELETE CASCADE;