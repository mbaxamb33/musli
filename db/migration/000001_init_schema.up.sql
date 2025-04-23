-- Migration Up: Create Project Management Platform Database Schema

-- Users Table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects Table
CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(user_id),
    project_name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Companies Table
CREATE TABLE companies (
    company_id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    website VARCHAR(255),
    industry VARCHAR(100),
    description TEXT,
    headquarters_location VARCHAR(255),
    founded_year INTEGER,
    is_public BOOLEAN,
    ticker_symbol VARCHAR(10),
    scrape_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contacts Table
CREATE TABLE contacts (
    contact_id SERIAL PRIMARY KEY,
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    linkedin_profile VARCHAR(255),
    job_title VARCHAR(255),
    company_id INTEGER REFERENCES companies(company_id),
    location VARCHAR(255),
    bio TEXT,
    scrape_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Project-Company Association Table
CREATE TABLE project_companies (
    project_id INTEGER REFERENCES projects(project_id),
    company_id INTEGER REFERENCES companies(company_id),
    association_notes TEXT,
    matching_score DECIMAL(5,2),
    approach_strategy TEXT,
    PRIMARY KEY (project_id, company_id)
);

-- Project-Contact Association Table
CREATE TABLE project_contacts (
    project_id INTEGER REFERENCES projects(project_id),
    contact_id INTEGER REFERENCES contacts(contact_id),
    association_notes TEXT,
    PRIMARY KEY (project_id, contact_id)
);

-- File/Document Storage Table
CREATE TABLE project_files (
    file_id SERIAL PRIMARY KEY,
    project_id INTEGER REFERENCES projects(project_id),
    file_name VARCHAR(255) NOT NULL,
    file_type VARCHAR(50) NOT NULL,
    file_path VARCHAR(500) NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    file_size INTEGER
);

-- Datasource Table
CREATE TABLE datasources (
    datasource_id SERIAL PRIMARY KEY,
    source_type VARCHAR(50) NOT NULL,
    source_id INTEGER,
    extraction_timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Paragraphs Table
CREATE TABLE paragraphs (
    paragraph_id SERIAL PRIMARY KEY,
    datasource_id INTEGER REFERENCES datasources(datasource_id),
    content TEXT NOT NULL,
    main_idea TEXT,
    classification VARCHAR(100),
    confidence_score DECIMAL(5,2)
);

-- Company News Table
CREATE TABLE company_news (
    news_id SERIAL PRIMARY KEY,
    company_id INTEGER REFERENCES companies(company_id),
    title VARCHAR(255) NOT NULL,
    publication_date DATE,
    source VARCHAR(255),
    url VARCHAR(500),
    summary TEXT,
    sentiment VARCHAR(50),
    datasource_id INTEGER REFERENCES datasources(datasource_id)
);

-- Contact News/Mentions Table
CREATE TABLE contact_news (
    mention_id SERIAL PRIMARY KEY,
    contact_id INTEGER REFERENCES contacts(contact_id),
    title VARCHAR(255) NOT NULL,
    publication_date DATE,
    source VARCHAR(255),
    url VARCHAR(500),
    summary TEXT,
    datasource_id INTEGER REFERENCES datasources(datasource_id)
);

-- Matching Score Criteria Table
CREATE TABLE matching_score_criteria (
    criteria_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    weight DECIMAL(5,2) NOT NULL
);

-- Indexes for performance optimization
CREATE INDEX idx_project_user ON projects(user_id);
CREATE INDEX idx_contact_company ON contacts(company_id);
CREATE INDEX idx_datasource_type ON datasources(source_type);
CREATE INDEX idx_project_companies_score ON project_companies(matching_score);