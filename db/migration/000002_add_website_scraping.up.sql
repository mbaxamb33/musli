-- Add website scraping capabilities

-- Table to store company websites
CREATE TABLE company_websites (
    website_id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES companies(company_id) NOT NULL,
    base_url VARCHAR(500) NOT NULL,
    site_title VARCHAR(255),
    last_scraped_at TIMESTAMP,
    scrape_frequency_days INTEGER DEFAULT 30,
    is_active BOOLEAN DEFAULT TRUE,
    datasource_id INTEGER REFERENCES datasources(datasource_id),
    UNIQUE(company_id, base_url)
);

-- Table to store individual pages from a website
CREATE TABLE website_pages (
    page_id SERIAL PRIMARY KEY,
    website_id INTEGER REFERENCES company_websites(website_id) NOT NULL,
    url VARCHAR(1000) NOT NULL,
    path VARCHAR(500) NOT NULL,
    title VARCHAR(500),
    parent_page_id INTEGER REFERENCES website_pages(page_id),
    depth INTEGER NOT NULL,
    last_extracted_at TIMESTAMP,
    extract_status VARCHAR(50),
    datasource_id INTEGER REFERENCES datasources(datasource_id),
    UNIQUE(website_id, url)
);

-- Add a new source type option for datasources
ALTER TABLE datasources ADD CONSTRAINT check_source_type 
    CHECK (source_type IN ('website', 'page', 'api', 'file', 'manual', 'other'));

-- Create indexes for performance
CREATE INDEX idx_company_websites_company ON company_websites(company_id);
CREATE INDEX idx_website_pages_website ON website_pages(website_id);
CREATE INDEX idx_website_pages_parent ON website_pages(parent_page_id);
CREATE INDEX idx_website_pages_path ON website_pages(path);