-- Migration to add multiple datasource relationships to companies and contacts

-- Step 1: Create junction table for companies and datasources
CREATE TABLE company_datasources (
    company_id INTEGER NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    datasource_id INTEGER NOT NULL REFERENCES datasources(datasource_id) ON DELETE CASCADE,
    PRIMARY KEY (company_id, datasource_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 2: Create junction table for contacts and datasources
CREATE TABLE contact_datasources (
    contact_id INTEGER NOT NULL REFERENCES contacts(contact_id) ON DELETE CASCADE,
    datasource_id INTEGER NOT NULL REFERENCES datasources(datasource_id) ON DELETE CASCADE,
    PRIMARY KEY (contact_id, datasource_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Step 3: Create indexes for better performance
CREATE INDEX idx_company_datasources_company_id ON company_datasources(company_id);
CREATE INDEX idx_company_datasources_datasource_id ON company_datasources(datasource_id);
CREATE INDEX idx_contact_datasources_contact_id ON contact_datasources(contact_id);
CREATE INDEX idx_contact_datasources_datasource_id ON contact_datasources(datasource_id);

-- Step 4: Migration down (rollback) script
-- This should be added to the down migration file
/*
DROP INDEX IF EXISTS idx_contact_datasources_datasource_id;
DROP INDEX IF EXISTS idx_contact_datasources_contact_id;
DROP INDEX IF EXISTS idx_company_datasources_datasource_id;
DROP INDEX IF EXISTS idx_company_datasources_company_id;

DROP TABLE IF EXISTS contact_datasources;
DROP TABLE IF EXISTS company_datasources;
*/