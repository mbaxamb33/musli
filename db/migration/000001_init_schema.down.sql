-- Migration Down: Drop Project Management Platform Database Schema

-- Drop indexes first
DROP INDEX IF EXISTS idx_project_companies_score;
DROP INDEX IF EXISTS idx_datasource_type;
DROP INDEX IF EXISTS idx_contact_company;
DROP INDEX IF EXISTS idx_project_user;

-- Drop tables in reverse order of creation to avoid foreign key constraints issues
DROP TABLE IF EXISTS matching_score_criteria;
DROP TABLE IF EXISTS contact_news;
DROP TABLE IF EXISTS company_news;
DROP TABLE IF EXISTS paragraphs;
DROP TABLE IF EXISTS datasources;
DROP TABLE IF EXISTS project_files;
DROP TABLE IF EXISTS project_contacts;
DROP TABLE IF EXISTS project_companies;
DROP TABLE IF EXISTS contacts;
DROP TABLE IF EXISTS companies;
DROP TABLE IF EXISTS projects;
DROP TABLE IF EXISTS users;