-- Migration Down: Drop the new schema structure

-- First drop all indexes
DROP INDEX IF EXISTS idx_proposition_drafts_sales_process_id;
DROP INDEX IF EXISTS idx_customer_needs_sales_process_id;
DROP INDEX IF EXISTS idx_analysis_inputs_analysis_id;
DROP INDEX IF EXISTS idx_analyses_sales_process_id;
DROP INDEX IF EXISTS idx_meetings_contact_id;
DROP INDEX IF EXISTS idx_meetings_sales_process_id;
DROP INDEX IF EXISTS idx_tasks_sales_process_id;
DROP INDEX IF EXISTS idx_sales_processes_contact_id;
DROP INDEX IF EXISTS idx_sales_processes_user_id;
DROP INDEX IF EXISTS idx_contact_news_contact_id;
DROP INDEX IF EXISTS idx_company_news_company_id;
DROP INDEX IF EXISTS idx_paragraphs_datasource_id;
DROP INDEX IF EXISTS idx_projects_user_id;
DROP INDEX IF EXISTS idx_contacts_company_id;
DROP INDEX IF EXISTS idx_companies_user_id;

-- Drop tables in reverse order of creation to avoid foreign key constraints
DROP TABLE IF EXISTS project_datasources;
DROP TABLE IF EXISTS proposition_drafts;
DROP TABLE IF EXISTS needs_datasource_matches;
DROP TABLE IF EXISTS customer_needs;
DROP TABLE IF EXISTS analysis_contact_info;
DROP TABLE IF EXISTS analysis_inputs;
DROP TABLE IF EXISTS analyses;
DROP TABLE IF EXISTS meetings;
DROP TABLE IF EXISTS tasks;
DROP TABLE IF EXISTS sales_process_projects;
DROP TABLE IF EXISTS sales_processes;
DROP TABLE IF EXISTS contact_news;
DROP TABLE IF EXISTS company_news;
DROP TABLE IF EXISTS paragraphs;
DROP TABLE IF EXISTS datasources;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS users;

-- Drop custom types
DROP TYPE IF EXISTS input_type;
DROP TYPE IF EXISTS task_status;
DROP TYPE IF EXISTS datasource_type;