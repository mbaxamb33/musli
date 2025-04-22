-- Revert the initial database migration for project-resource-company matching platform
-- Drop tables in reverse order of creation to respect foreign key constraints

-- Drop indexes first
DROP INDEX IF EXISTS idx_processed_company_data_company_id;
DROP INDEX IF EXISTS idx_web_scrape_data_company_id;
DROP INDEX IF EXISTS idx_company_contacts_contact_id;
DROP INDEX IF EXISTS idx_company_contacts_company_id;
DROP INDEX IF EXISTS idx_project_companies_company_id;
DROP INDEX IF EXISTS idx_project_companies_project_id;
DROP INDEX IF EXISTS idx_project_resources_resource_id;
DROP INDEX IF EXISTS idx_project_resources_project_id;
DROP INDEX IF EXISTS idx_resources_category_id;
DROP INDEX IF EXISTS idx_projects_user_id;

-- Data Collection Tables
DROP TABLE IF EXISTS processed_company_data;
DROP TABLE IF EXISTS web_scrape_data;
DROP TABLE IF EXISTS data_sources;

-- Supportive Tables
DROP TABLE IF EXISTS matching_scores_detail;
DROP TABLE IF EXISTS approach_strategies;
DROP TABLE IF EXISTS matching_criteria;

-- Relationship Tables
DROP TABLE IF EXISTS company_contacts;
DROP TABLE IF EXISTS project_companies;
DROP TABLE IF EXISTS project_processing;
DROP TABLE IF EXISTS project_resources;

-- Core Tables
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS resources;
DROP TABLE IF EXISTS resource_categories;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;
