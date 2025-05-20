-- Migration to revert the renaming of cognito_sub back to user_id

-- Step 1: Drop existing foreign key constraints
ALTER TABLE companies DROP CONSTRAINT IF EXISTS companies_cognito_sub_fkey;
ALTER TABLE projects DROP CONSTRAINT IF EXISTS projects_cognito_sub_fkey;
ALTER TABLE sales_processes DROP CONSTRAINT IF EXISTS sales_processes_cognito_sub_fkey;

-- Step 2: Rename cognito_sub columns back to user_id
ALTER TABLE companies RENAME COLUMN cognito_sub TO user_id;
ALTER TABLE projects RENAME COLUMN cognito_sub TO user_id;
ALTER TABLE sales_processes RENAME COLUMN cognito_sub TO user_id;

-- Step 3: Re-add the foreign key constraints with original column name
ALTER TABLE companies ADD CONSTRAINT companies_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users (cognito_sub) ON DELETE CASCADE;
    
ALTER TABLE projects ADD CONSTRAINT projects_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users (cognito_sub) ON DELETE CASCADE;
    
ALTER TABLE sales_processes ADD CONSTRAINT sales_processes_user_id_fkey 
    FOREIGN KEY (user_id) REFERENCES users (cognito_sub) ON DELETE CASCADE;

-- Step 4: Update indexes to match original column names
DROP INDEX IF EXISTS idx_companies_cognito_sub;
DROP INDEX IF EXISTS idx_projects_cognito_sub;
DROP INDEX IF EXISTS idx_sales_processes_cognito_sub;

CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_sales_processes_user_id ON sales_processes(user_id);