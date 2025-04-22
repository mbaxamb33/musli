-- Initial database migration for project-resource-company matching platform

-- Core Tables

-- Users table
CREATE TABLE IF NOT EXISTS users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(100) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects table
CREATE TABLE IF NOT EXISTS projects (
    project_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    start_date DATE,
    end_date DATE,
    status VARCHAR(20) DEFAULT 'active',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
);

-- Resource Categories table (created before Resources table due to foreign key)
CREATE TABLE IF NOT EXISTS resource_categories (
    category_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Resources table
CREATE TABLE IF NOT EXISTS resources (
    resource_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    category_id INTEGER,
    unit VARCHAR(20),
    cost_per_unit DECIMAL(10, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (category_id) REFERENCES resource_categories(category_id) ON DELETE SET NULL
);

-- Companies table
CREATE TABLE IF NOT EXISTS companies (
    company_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    industry VARCHAR(50),
    size VARCHAR(20),
    location VARCHAR(100),
    website VARCHAR(255),
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contacts table
CREATE TABLE IF NOT EXISTS contacts (
    contact_id SERIAL PRIMARY KEY,
    first_name VARCHAR(50) NOT NULL,
    last_name VARCHAR(50) NOT NULL,
    title VARCHAR(100),
    email VARCHAR(100),
    phone VARCHAR(20),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Approach Strategies table (created before Project Companies table)
CREATE TABLE IF NOT EXISTS approach_strategies (
    strategy_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    recommended_score_min DECIMAL(5, 2),
    recommended_score_max DECIMAL(5, 2),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Relationship Tables

-- Project Resources table
CREATE TABLE IF NOT EXISTS project_resources (
    project_resource_id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    resource_id INTEGER NOT NULL,
    quantity DECIMAL(10, 2) NOT NULL DEFAULT 1,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    FOREIGN KEY (resource_id) REFERENCES resources(resource_id) ON DELETE CASCADE,
    UNIQUE(project_id, resource_id)
);

-- Project Processing table
CREATE TABLE IF NOT EXISTS project_processing (
    processing_id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    processing_type VARCHAR(50) NOT NULL,
    description TEXT,
    status VARCHAR(20) DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE
);

-- Project Companies table
CREATE TABLE IF NOT EXISTS project_companies (
    project_company_id SERIAL PRIMARY KEY,
    project_id INTEGER NOT NULL,
    company_id INTEGER NOT NULL,
    matching_score DECIMAL(5, 2),
    approach_strategy_id INTEGER,
    status VARCHAR(20) DEFAULT 'potential',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_id) REFERENCES projects(project_id) ON DELETE CASCADE,
    FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE CASCADE,
    FOREIGN KEY (approach_strategy_id) REFERENCES approach_strategies(strategy_id) ON DELETE SET NULL,
    UNIQUE(project_id, company_id)
);

-- Company Contacts table
CREATE TABLE IF NOT EXISTS company_contacts (
    company_contact_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    contact_id INTEGER NOT NULL,
    is_primary BOOLEAN DEFAULT false,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE CASCADE,
    FOREIGN KEY (contact_id) REFERENCES contacts(contact_id) ON DELETE CASCADE,
    UNIQUE(company_id, contact_id)
);

-- Supportive Tables

-- Matching Criteria table
CREATE TABLE IF NOT EXISTS matching_criteria (
    criteria_id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    description TEXT,
    weight DECIMAL(3, 2) DEFAULT 1.00,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Matching Scores Detail table
CREATE TABLE IF NOT EXISTS matching_scores_detail (
    score_detail_id SERIAL PRIMARY KEY,
    project_company_id INTEGER NOT NULL,
    criteria_id INTEGER NOT NULL,
    score DECIMAL(5, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (project_company_id) REFERENCES project_companies(project_company_id) ON DELETE CASCADE,
    FOREIGN KEY (criteria_id) REFERENCES matching_criteria(criteria_id) ON DELETE CASCADE
);

-- Data Collection Tables

-- Web Scrape Data table
CREATE TABLE IF NOT EXISTS web_scrape_data (
    scrape_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    source_url VARCHAR(255) NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    content TEXT,
    scrape_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    is_processed BOOLEAN DEFAULT false,
    FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE CASCADE
);

-- Processed Company Data table
CREATE TABLE IF NOT EXISTS processed_company_data (
    data_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL,
    data_type VARCHAR(50) NOT NULL,
    data_key VARCHAR(100) NOT NULL,
    data_value TEXT,
    confidence_score DECIMAL(5, 2),
    source_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (company_id) REFERENCES companies(company_id) ON DELETE CASCADE,
    FOREIGN KEY (source_id) REFERENCES web_scrape_data(scrape_id) ON DELETE SET NULL
);

-- Data Sources table
CREATE TABLE IF NOT EXISTS data_sources (
    source_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    url_pattern VARCHAR(255),
    api_endpoint VARCHAR(255),
    api_key VARCHAR(255),
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_resources_category_id ON resources(category_id);
CREATE INDEX idx_project_resources_project_id ON project_resources(project_id);
CREATE INDEX idx_project_resources_resource_id ON project_resources(resource_id);
CREATE INDEX idx_project_companies_project_id ON project_companies(project_id);
CREATE INDEX idx_project_companies_company_id ON project_companies(company_id);
CREATE INDEX idx_company_contacts_company_id ON company_contacts(company_id);
CREATE INDEX idx_company_contacts_contact_id ON company_contacts(contact_id);
CREATE INDEX idx_web_scrape_data_company_id ON web_scrape_data(company_id);
CREATE INDEX idx_processed_company_data_company_id ON processed_company_data(company_id);