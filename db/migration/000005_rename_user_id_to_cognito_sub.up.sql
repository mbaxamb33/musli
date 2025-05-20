-- Migration to rename user_id columns to cognito_sub for better clarity and consistency

-- Step 1: Drop existing foreign key constraints
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_user_id_fkey;
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_user_id_fkey;
ALTER TABLE sales_processes DROP CONSTRAINT IF EXISTS sales_processes_user_id_fkey;

-- Step 2: Rename user_id columns to cognito_sub
ALTER TABLE companies RENAME COLUMN user_id TO cognito_sub;
ALTER TABLE projects RENAME COLUMN user_id TO cognito_sub;
ALTER TABLE sales_processes RENAME COLUMN user_id TO cognito_sub;

-- Step 3: Re-add the foreign key constraints with new column name
ALTER TABLE companies ADD CONSTRAINT companies_cognito_sub_fkey 
    FOREIGN KEY (cognito_sub) REFERENCES users (cognito_sub) ON DELETE CASCADE;
    
ALTER TABLE projects ADD CONSTRAINT projects_cognito_sub_fkey 
    FOREIGN KEY (cognito_sub) REFERENCES users (cognito_sub) ON DELETE CASCADE;
    
ALTER TABLE sales_processes ADD CONSTRAINT sales_processes_cognito_sub_fkey 
    FOREIGN KEY (cognito_sub) REFERENCES users (cognito_sub) ON DELETE CASCADE;

-- Step 4: Update indexes to match new column names
DROP INDEX IF EXISTS idx_companies_user_id;
DROP INDEX IF EXISTS idx_projects_user_id;
DROP INDEX IF EXISTS idx_sales_processes_user_id;

CREATE INDEX idx_companies_cognito_sub ON companies(cognito_sub);
CREATE INDEX idx_projects_cognito_sub ON projects(cognito_sub);
CREATE INDEX idx_sales_processes_cognito_sub ON sales_processes(cognito_sub);