-- Migration Up: Implement the new schema structure

-- Users table
CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username VARCHAR(100) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Company table
CREATE TABLE companies (
    company_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    company_name VARCHAR(200) NOT NULL,
    industry VARCHAR(100),
    website VARCHAR(255),
    address TEXT,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contacts table
CREATE TABLE contacts (
    contact_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    position VARCHAR(100),
    email VARCHAR(255),
    phone VARCHAR(50),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects table
CREATE TABLE projects (
    project_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    project_name VARCHAR(200) NOT NULL,
    main_idea TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Datasource types enum
CREATE TYPE datasource_type AS ENUM ('mp3', 'website', 'word_document', 'pdf', 'excel', 'powerpoint', 'plain_text');

-- Datasources table
CREATE TABLE datasources (
    datasource_id SERIAL PRIMARY KEY,
    source_type datasource_type NOT NULL,
    link VARCHAR(255),
    file_data BYTEA,
    file_name VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Paragraphs table
CREATE TABLE paragraphs (
    paragraph_id SERIAL PRIMARY KEY,
    datasource_id INTEGER NOT NULL REFERENCES datasources(datasource_id) ON DELETE CASCADE,
    title VARCHAR(255),
    main_idea TEXT,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Company news table
CREATE TABLE company_news (
    company_news_id SERIAL PRIMARY KEY,
    company_id INTEGER NOT NULL REFERENCES companies(company_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    datasource_id INTEGER REFERENCES datasources(datasource_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Contact news table
CREATE TABLE contact_news (
    contact_news_id SERIAL PRIMARY KEY,
    contact_id INTEGER NOT NULL REFERENCES contacts(contact_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    datasource_id INTEGER REFERENCES datasources(datasource_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sales process table
CREATE TABLE sales_processes (
    sales_process_id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    contact_id INTEGER NOT NULL REFERENCES contacts(contact_id),
    overall_matching_score DECIMAL(5,2),
    status VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Projects in sales process (junction table)
CREATE TABLE sales_process_projects (
    sales_process_id INTEGER NOT NULL REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE,
    project_id INTEGER NOT NULL REFERENCES projects(project_id) ON DELETE CASCADE,
    PRIMARY KEY (sales_process_id, project_id)
);

-- Task status enum
CREATE TYPE task_status AS ENUM ('not_started', 'in_progress', 'completed', 'stopped');

-- Tasks table
CREATE TABLE tasks (
    task_id SERIAL PRIMARY KEY,
    sales_process_id INTEGER NOT NULL REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status task_status NOT NULL DEFAULT 'not_started',
    due_date DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Meetings table
CREATE TABLE meetings (
    meeting_id SERIAL PRIMARY KEY,
    sales_process_id INTEGER NOT NULL REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE,
    contact_id INTEGER NOT NULL REFERENCES contacts(contact_id),
    task_id INTEGER REFERENCES tasks(task_id),
    meeting_time TIMESTAMP NOT NULL,
    meeting_place VARCHAR(255),
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Analysis table
CREATE TABLE analyses (
    analysis_id SERIAL PRIMARY KEY,
    sales_process_id INTEGER NOT NULL REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE,
    version INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (sales_process_id, version)
);

-- Input types enum
CREATE TYPE input_type AS ENUM ('personal_input', 'meeting_input', 'other');

-- Analysis inputs table
CREATE TABLE analysis_inputs (
    input_id SERIAL PRIMARY KEY,
    analysis_id INTEGER NOT NULL REFERENCES analyses(analysis_id) ON DELETE CASCADE,
    input_type input_type NOT NULL,
    datasource_id INTEGER REFERENCES datasources(datasource_id),
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Analysis contact information table
CREATE TABLE analysis_contact_info (
    analysis_id INTEGER PRIMARY KEY REFERENCES analyses(analysis_id) ON DELETE CASCADE,
    problems TEXT,
    needs TEXT,
    urgency TEXT,
    priorities TEXT,
    decision_process TEXT,
    budget TEXT,
    resources TEXT,
    relevant_information TEXT,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Customer needs table
CREATE TABLE customer_needs (
    need_id SERIAL PRIMARY KEY,
    sales_process_id INTEGER NOT NULL REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE,
    need_description TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Customer needs to project datasources matching table
CREATE TABLE needs_datasource_matches (
    match_id SERIAL PRIMARY KEY,
    need_id INTEGER NOT NULL REFERENCES customer_needs(need_id) ON DELETE CASCADE,
    datasource_id INTEGER NOT NULL REFERENCES datasources(datasource_id) ON DELETE CASCADE,
    match_score DECIMAL(5,2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (need_id, datasource_id)
);

-- Proposition drafts table
CREATE TABLE proposition_drafts (
    draft_id SERIAL PRIMARY KEY,
    sales_process_id INTEGER NOT NULL REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    content TEXT,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Project datasources junction table
CREATE TABLE project_datasources (
    project_id INTEGER NOT NULL REFERENCES projects(project_id) ON DELETE CASCADE,
    datasource_id INTEGER NOT NULL REFERENCES datasources(datasource_id) ON DELETE CASCADE,
    PRIMARY KEY (project_id, datasource_id)
);

-- Create indexes for better performance
CREATE INDEX idx_companies_user_id ON companies(user_id);
CREATE INDEX idx_contacts_company_id ON contacts(company_id);
CREATE INDEX idx_projects_user_id ON projects(user_id);
CREATE INDEX idx_paragraphs_datasource_id ON paragraphs(datasource_id);
CREATE INDEX idx_company_news_company_id ON company_news(company_id);
CREATE INDEX idx_contact_news_contact_id ON contact_news(contact_id);
CREATE INDEX idx_sales_processes_user_id ON sales_processes(user_id);
CREATE INDEX idx_sales_processes_contact_id ON sales_processes(contact_id);
CREATE INDEX idx_tasks_sales_process_id ON tasks(sales_process_id);
CREATE INDEX idx_meetings_sales_process_id ON meetings(sales_process_id);
CREATE INDEX idx_meetings_contact_id ON meetings(contact_id);
CREATE INDEX idx_analyses_sales_process_id ON analyses(sales_process_id);
CREATE INDEX idx_analysis_inputs_analysis_id ON analysis_inputs(analysis_id);
CREATE INDEX idx_customer_needs_sales_process_id ON customer_needs(sales_process_id);
CREATE INDEX idx_proposition_drafts_sales_process_id ON proposition_drafts(sales_process_id);