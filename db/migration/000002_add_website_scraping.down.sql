-- Revert website scraping capabilities

-- Drop indexes first
DROP INDEX IF EXISTS idx_website_pages_path;
DROP INDEX IF EXISTS idx_website_pages_parent;
DROP INDEX IF EXISTS idx_website_pages_website;
DROP INDEX IF EXISTS idx_company_websites_company;

-- Remove constraint on datasources
ALTER TABLE datasources DROP CONSTRAINT IF EXISTS check_source_type;

-- Drop tables in reverse order of creation to avoid foreign key constraint issues
DROP TABLE IF EXISTS website_pages;
DROP TABLE IF EXISTS company_websites;
