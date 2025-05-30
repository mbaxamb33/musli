-- 000006_add_brief_system.up.sql
-- Migration Up: Add brief system to existing schema
-- Implements hierarchical sales intelligence with 150 fields across 9 categories

-- Master Brief Management
-- Top-level container for sales intelligence, links to users/companies/contacts
CREATE TABLE master_briefs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    cognito_sub VARCHAR NOT NULL REFERENCES users(cognito_sub) ON DELETE CASCADE, -- Owner of the brief
    company_id INTEGER REFERENCES companies(company_id) ON DELETE SET NULL, -- Link to existing company
    contact_id INTEGER REFERENCES contacts(contact_id) ON DELETE SET NULL, -- Link to existing contact
    company_reference VARCHAR(255) NOT NULL, -- Text reference to company
    contact_reference VARCHAR(255) NOT NULL, -- Text reference to contact
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Brief Types and Management
-- Enum types for categorizing briefs by type and processing stage
CREATE TYPE brief_type AS ENUM ('master', 'stage_specific', 'regular');
CREATE TYPE brief_tag AS ENUM ('initial', 'specific', 'updated', 'final');

-- Individual briefs under a master brief
-- Each brief can be stage-specific or regular, contains text content and attachments
CREATE TABLE briefs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_brief_id UUID REFERENCES master_briefs(id) ON DELETE CASCADE, -- Parent master brief
    brief_type brief_type NOT NULL, -- Type of brief (master/stage_specific/regular)
    brief_tag brief_tag NOT NULL, -- Processing stage tag
    title VARCHAR(255), -- Brief title
    text_content TEXT, -- Main text content of the brief
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Media Attachments (extends existing datasources concept)
-- Links briefs to existing datasource files (images, documents, voice memos)
CREATE TABLE brief_attachments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    datasource_id INTEGER REFERENCES datasources(datasource_id) ON DELETE CASCADE, -- Existing datasource
    attachment_type VARCHAR(50), -- 'image', 'document', 'voice_memo'
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Company Intelligence (Fields 1-15)
-- Basic company information, financials, market position, and business context
CREATE TABLE company_intelligence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 1: Legal business name
    company_name VARCHAR(255),
    -- Field 2: Business model and core activities
    company_overview TEXT,
    -- Field 3: Primary industry classification
    industry_sector VARCHAR(100),
    -- Field 4: Annual revenue in USD
    company_revenue DECIMAL(15,2),
    -- Field 5: Total number of employees
    employee_count INTEGER,
    -- Field 6: Locations where company operates
    geographic_footprint TEXT,
    -- Field 7: Ownership structure and subsidiaries
    parent_company VARCHAR(255),
    -- Field 8: Competitive ranking in industry
    market_position VARCHAR(100),
    -- Field 9: Recent press releases or news coverage
    recent_news_events TEXT,
    -- Field 10: Credit rating and financial stability
    financial_health VARCHAR(50),
    -- Field 11: Revenue and employee growth trends
    growth_trajectory TEXT,
    -- Field 12: External forces affecting business
    market_pressures TEXT,
    -- Field 13: Compliance requirements and changes
    regulatory_environment TEXT,
    -- Field 14: Recent or planned M&A activity
    merger_acquisition_activity TEXT,
    -- Field 15: Key competitors and market dynamics
    competitive_landscape TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Strategic Context (Fields 16-30)
-- Corporate strategy, initiatives, goals, and transformation plans
CREATE TABLE strategic_context (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 16: Overall corporate strategy and direction
    business_strategy TEXT,
    -- Field 17: Key projects driving transformation
    strategic_initiatives TEXT,
    -- Field 18: Current quarter's top 3-5 priorities
    quarterly_priorities TEXT,
    -- Field 19: Year-end targets and objectives
    annual_goals TEXT,
    -- Field 20: Digital or operational transformation plans
    transformation_agenda TEXT,
    -- Field 21: Current state of digital adoption
    digital_maturity VARCHAR(50),
    -- Field 22: Areas of R&D investment
    innovation_focus TEXT,
    -- Field 23: Process inefficiencies and bottlenecks
    operational_challenges TEXT,
    -- Field 24: Mandates to reduce expenses
    cost_reduction_pressures TEXT,
    -- Field 25: Growth expectations and timelines
    revenue_growth_targets TEXT,
    -- Field 26: Productivity improvement requirements
    efficiency_mandates TEXT,
    -- Field 27: New regulations requiring action
    compliance_drivers TEXT,
    -- Field 28: Key risks being addressed
    risk_management_priorities TEXT,
    -- Field 29: ESG commitments and targets
    sustainability_goals TEXT,
    -- Field 30: Planned technology investments
    technology_roadmap TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Buying Committee Intelligence (Fields 31-50)
-- Decision makers, influencers, blockers, and committee dynamics
CREATE TABLE buying_committee (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 31: Person with budget authority
    economic_buyer_name VARCHAR(255),
    -- Field 32: Job title and level
    economic_buyer_title VARCHAR(255),
    -- Field 33: Level of decision-making power (1-10)
    economic_buyer_influence INTEGER CHECK (economic_buyer_influence BETWEEN 1 AND 10),
    -- Field 34: Personal and professional drivers
    economic_buyer_motivations TEXT,
    -- Field 35: Person evaluating technical requirements
    technical_buyer_name VARCHAR(255),
    -- Field 36: Key technical evaluation criteria
    technical_buyer_concerns TEXT,
    -- Field 37: End users involved in evaluation
    user_buyer_representatives TEXT,
    -- Field 38: Internal advocate for your solution
    coach_champion_name VARCHAR(255),
    -- Field 39: Political capital and reach (1-10)
    coach_influence_level INTEGER CHECK (coach_influence_level BETWEEN 1 AND 10),
    -- Field 40: People opposing the purchase
    blocker_identification TEXT,
    -- Field 41: Specific objections and resistance points
    blocker_concerns TEXT,
    -- Field 42: How group makes decisions together
    committee_dynamics TEXT,
    -- Field 43: Steps from evaluation to signature
    decision_making_process TEXT,
    -- Field 44: Level of agreement needed
    consensus_requirements TEXT,
    -- Field 45: Each person's comfort with change
    individual_risk_tolerance TEXT,
    -- Field 46: How this decision affects careers
    career_motivations TEXT,
    -- Field 47: How individuals measure success
    personal_success_metrics TEXT,
    -- Field 48: Who influences whom internally
    relationship_mapping TEXT,
    -- Field 49: Preferred meeting styles and frequency
    communication_preferences TEXT,
    -- Field 50: Informal power structures
    influence_networks TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Current State Assessment (Fields 51-65)
-- Existing solutions, pain points, satisfaction levels, and switching barriers
CREATE TABLE current_state_assessment (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 51: Existing vendor or internal solution
    current_solution_provider VARCHAR(255),
    -- Field 52: Satisfaction level (1-10)
    current_solution_satisfaction INTEGER CHECK (current_solution_satisfaction BETWEEN 1 AND 10),
    -- Field 53: Exact problems with current state
    specific_pain_points TEXT,
    -- Field 54: Manual processes compensating for gaps
    workaround_solutions TEXT,
    -- Field 55: Financial impact of not changing
    cost_of_status_quo DECIMAL(15,2),
    -- Field 56: Obstacles to changing vendors
    switching_barriers TEXT,
    -- Field 57: When current contracts expire
    contract_end_dates DATE,
    -- Field 58: Decision points for renewals
    renewal_timing TEXT,
    -- Field 59: Quality of current vendor relationship
    vendor_relationship_health VARCHAR(50),
    -- Field 60: Experience with current support
    support_satisfaction TEXT,
    -- Field 61: Missing features in current solution
    functionality_gaps TEXT,
    -- Field 62: Speed, reliability, or capacity problems
    performance_issues TEXT,
    -- Field 63: Limits preventing growth
    scalability_constraints TEXT,
    -- Field 64: Problems connecting systems
    integration_challenges TEXT,
    -- Field 65: End user resistance or confusion
    user_adoption_issues TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Competitive Intelligence (Fields 66-80)
-- Competitor analysis, evaluation criteria, and decision timeline
CREATE TABLE competitive_intelligence (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 66: Other vendors being considered
    competitors_in_evaluation TEXT,
    -- Field 67: Any favoritism toward specific vendors
    preferred_vendor_bias TEXT,
    -- Field 68: Past relationships and experiences
    previous_vendor_history TEXT,
    -- Field 69: Competitors' advantages in this deal
    competitive_strengths TEXT,
    -- Field 70: Areas where competitors fall short
    competitive_weaknesses TEXT,
    -- Field 71: Budget range and price sensitivity
    pricing_expectations TEXT,
    -- Field 72: How solutions compare feature-by-feature
    feature_comparison_matrix JSONB,
    -- Field 73: Factors that will determine winner
    vendor_selection_criteria TEXT,
    -- Field 74: Relative importance of each criterion
    criteria_weighting JSONB,
    -- Field 75: Steps in vendor assessment
    evaluation_process TEXT,
    -- Field 76: Customer references needed
    reference_requirements TEXT,
    -- Field 77: Demonstration requirements
    proof_of_concept_needs TEXT,
    -- Field 78: Test implementation parameters
    pilot_program_scope TEXT,
    -- Field 79: Executive presentation requirements
    final_presentation_format TEXT,
    -- Field 80: Key dates and final decision deadline
    decision_timeline DATE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Financial & Procurement (Fields 81-95)
-- Budget details, approval processes, and procurement requirements
CREATE TABLE financial_procurement (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 81: Complete budget allocation
    total_available_budget DECIMAL(15,2),
    -- Field 82: Department or cost center funding
    budget_source VARCHAR(255),
    -- Field 83: Steps to approve spending
    budget_approval_workflow TEXT,
    -- Field 84: Purchasing department requirements
    procurement_process TEXT,
    -- Field 85: Company procurement rules
    purchasing_policies TEXT,
    -- Field 86: Required payment schedules
    payment_terms_constraints TEXT,
    -- Field 87: Who can approve what amounts
    financial_approval_levels TEXT,
    -- Field 88: When budgets reset or refresh
    budget_cycle_timing TEXT,
    -- Field 89: ROI documentation needed
    cost_justification_requirements TEXT,
    -- Field 90: How they measure return on investment
    roi_calculation_method TEXT,
    -- Field 91: Time to break even
    payback_period_expectations TEXT,
    -- Field 92: Leasing or payment plan preferences
    financing_options TEXT,
    -- Field 93: Legal and commercial terms
    contract_terms_requirements TEXT,
    -- Field 94: Contract approval workflow
    legal_review_process TEXT,
    -- Field 95: Coverage and liability needs
    insurance_requirements TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Project Requirements (Fields 96-110)
-- Implementation scope, timeline, resources, and project management
CREATE TABLE project_requirements (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 96: Breadth and depth of implementation
    project_scope TEXT,
    -- Field 97: How success will be measured
    success_criteria TEXT,
    -- Field 98: Project schedule and milestones
    implementation_timeline TEXT,
    -- Field 99: Internal team assignments
    resource_allocation TEXT,
    -- Field 100: Roles and responsibilities
    project_team_structure TEXT,
    -- Field 101: How to handle organizational change
    change_management_approach TEXT,
    -- Field 102: User education and certification needs
    training_requirements TEXT,
    -- Field 103: Phased vs. big bang implementation
    rollout_strategy TEXT,
    -- Field 104: Test deployment scope and goals
    pilot_phase_design TEXT,
    -- Field 105: Contingencies for potential issues
    risk_mitigation_plan TEXT,
    -- Field 106: Stakeholder updates and messaging
    communication_plan TEXT,
    -- Field 107: How to involve affected parties
    stakeholder_engagement TEXT,
    -- Field 108: KPIs to track post-implementation
    performance_metrics TEXT,
    -- Field 109: Decision-making and oversight model
    governance_structure TEXT,
    -- Field 110: How to handle problems and conflicts
    escalation_procedures TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Technical & Integration (Fields 111-125)
-- IT architecture, security, compliance, and integration requirements
CREATE TABLE technical_integration (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 111: Current IT infrastructure and platforms
    technical_architecture TEXT,
    -- Field 112: Data protection and access controls
    security_requirements TEXT,
    -- Field 113: Industry regulations and certifications
    compliance_standards TEXT,
    -- Field 114: Systems that must connect
    integration_points TEXT,
    -- Field 115: Information to transfer from old systems
    data_migration_scope TEXT,
    -- Field 116: Configuration and development requirements
    customization_needs TEXT,
    -- Field 117: Expected growth and capacity needs
    scalability_requirements TEXT,
    -- Field 118: Speed and reliability expectations
    performance_benchmarks TEXT,
    -- Field 119: Backup and continuity requirements
    disaster_recovery_needs TEXT,
    -- Field 120: Data protection and recovery procedures
    backup_requirements TEXT,
    -- Field 121: User permissions and authentication
    access_control_requirements TEXT,
    -- Field 122: Activity logging and compliance tracking
    audit_trail_needs TEXT,
    -- Field 123: Analytics and dashboard requirements
    reporting_capabilities TEXT,
    -- Field 124: Third-party integrations and data exchange
    api_requirements TEXT,
    -- Field 125: Smartphone and tablet functionality
    mobile_access_needs TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Behavioral & Psychological Insights (Fields 126-140)
-- Decision-making patterns, communication styles, and organizational behavior
CREATE TABLE behavioral_insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 126: How individuals and groups decide
    decision_making_style TEXT,
    -- Field 127: Comfort with uncertainty and change (1-10)
    risk_aversion_level INTEGER CHECK (risk_aversion_level BETWEEN 1 AND 10),
    -- Field 128: How organization handles new initiatives
    change_adoption_patterns TEXT,
    -- Field 129: Willingness to try new approaches
    innovation_appetite TEXT,
    -- Field 130: How agreement is reached
    consensus_building_approach TEXT,
    -- Field 131: How disagreements are handled
    conflict_resolution_style TEXT,
    -- Field 132: Formal vs. informal information flow
    communication_patterns TEXT,
    -- Field 133: What creates credibility and confidence
    trust_building_factors TEXT,
    -- Field 134: Credentials and proof points needed
    credibility_requirements TEXT,
    -- Field 135: Transactional vs. partnership approach
    relationship_preferences TEXT,
    -- Field 136: Productive meeting styles and structures
    meeting_effectiveness TEXT,
    -- Field 137: Communication speed and reliability
    follow_up_responsiveness TEXT,
    -- Field 138: Level of detail and formality expected
    documentation_preferences TEXT,
    -- Field 139: Executive vs. technical focus
    presentation_style_preferences TEXT,
    -- Field 140: Collaborative vs. competitive style
    negotiation_approach TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Sales Process Tracking (Fields 141-150)
-- Sales funnel management, deal progression, and probability assessment
CREATE TABLE sales_process_tracking (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    brief_id UUID REFERENCES briefs(id) ON DELETE CASCADE, -- Parent brief
    -- Field 141: How opportunity was generated
    lead_source VARCHAR(255),
    -- Field 142: Current position in sales funnel
    opportunity_stage VARCHAR(100),
    -- Field 143: Likelihood of closing (0-100%)
    probability_percentage INTEGER CHECK (probability_percentage BETWEEN 0 AND 100),
    -- Field 144: Deal size multiplied by probability
    weighted_value DECIMAL(15,2),
    -- Field 145: Immediate steps to advance deal
    next_action_required TEXT,
    -- Field 146: Critical events and decision points
    key_milestones TEXT,
    -- Field 147: Speed of progression through stages
    sales_velocity DECIMAL(10,2),
    -- Field 148: Current energy and urgency level
    deal_momentum VARCHAR(50),
    -- Field 149: Standing relative to other vendors
    competitive_position VARCHAR(100),
    -- Field 150: Strengths and weaknesses affecting outcome
    win_probability_factors TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Ground Truth - Master aggregated values
-- Stores the "best" value for each field across all briefs under a master
-- Uses confidence scoring and source tracking for data quality
CREATE TABLE ground_truth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    master_brief_id UUID REFERENCES master_briefs(id) ON DELETE CASCADE, -- Parent master brief
    field_name VARCHAR(100), -- Name of the field (e.g., "company_name", "budget")
    field_value TEXT, -- Aggregated/best value for this field
    confidence_score DECIMAL(3,2) CHECK (confidence_score BETWEEN 0 AND 1), -- Quality score 0-1
    source_brief_ids UUID[], -- Array of brief IDs that contributed to this value
    last_updated TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- When this ground truth was last calculated
);

-- Link briefs to sales processes
-- Junction table connecting the brief system to existing sales workflow
CREATE TABLE sales_process_briefs (
    sales_process_id INTEGER REFERENCES sales_processes(sales_process_id) ON DELETE CASCADE, -- Existing sales process
    master_brief_id UUID REFERENCES master_briefs(id) ON DELETE CASCADE, -- Master brief
    PRIMARY KEY (sales_process_id, master_brief_id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Performance Indexes
-- Strategic indexes for common query patterns and foreign key lookups
CREATE INDEX idx_master_briefs_cognito_sub ON master_briefs(cognito_sub); -- User ownership queries
CREATE INDEX idx_master_briefs_company_id ON master_briefs(company_id); -- Company-based queries
CREATE INDEX idx_master_briefs_contact_id ON master_briefs(contact_id); -- Contact-based queries
CREATE INDEX idx_briefs_master_brief_id ON briefs(master_brief_id); -- Brief hierarchy navigation
CREATE INDEX idx_brief_attachments_brief_id ON brief_attachments(brief_id); -- Attachment lookups
CREATE INDEX idx_company_intelligence_brief_id ON company_intelligence(brief_id); -- Intelligence queries
CREATE INDEX idx_strategic_context_brief_id ON strategic_context(brief_id);
CREATE INDEX idx_buying_committee_brief_id ON buying_committee(brief_id);
CREATE INDEX idx_current_state_assessment_brief_id ON current_state_assessment(brief_id);
CREATE INDEX idx_competitive_intelligence_brief_id ON competitive_intelligence(brief_id);
CREATE INDEX idx_financial_procurement_brief_id ON financial_procurement(brief_id);
CREATE INDEX idx_project_requirements_brief_id ON project_requirements(brief_id);
CREATE INDEX idx_technical_integration_brief_id ON technical_integration(brief_id);
CREATE INDEX idx_behavioral_insights_brief_id ON behavioral_insights(brief_id);
CREATE INDEX idx_sales_process_tracking_brief_id ON sales_process_tracking(brief_id);
CREATE INDEX idx_ground_truth_master_brief_id ON ground_truth(master_brief_id); -- Ground truth aggregation
CREATE INDEX idx_ground_truth_field_name ON ground_truth(field_name); -- Field-based queries
CREATE INDEX idx_sales_process_briefs_sales_process_id ON sales_process_briefs(sales_process_id); -- Sales integration