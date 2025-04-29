-- Migration Down: Drop the junction tables and indexes

-- Drop indexes
DROP INDEX IF EXISTS idx_contact_datasources_datasource_id;
DROP INDEX IF EXISTS idx_contact_datasources_contact_id;
DROP INDEX IF EXISTS idx_company_datasources_datasource_id;
DROP INDEX IF EXISTS idx_company_datasources_company_id;

-- Drop junction tables
DROP TABLE IF EXISTS contact_datasources;
DROP TABLE IF EXISTS company_datasources;