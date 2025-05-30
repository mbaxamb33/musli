-- =============================================================================
-- COMPANY INTELLIGENCE QUERIES
-- =============================================================================

-- name: CreateCompanyIntelligence :one
INSERT INTO company_intelligence (
    brief_id, company_name, company_overview, industry_sector, company_revenue, employee_count,
    geographic_footprint, parent_company, market_position, recent_news_events, financial_health,
    growth_trajectory, market_pressures, regulatory_environment, merger_acquisition_activity,
    competitive_landscape
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, company_name, company_overview, industry_sector, company_revenue, employee_count,
    geographic_footprint, parent_company, market_position, recent_news_events, financial_health,
    growth_trajectory, market_pressures, regulatory_environment, merger_acquisition_activity,
    competitive_landscape, created_at, updated_at;

-- name: GetCompanyIntelligenceByBriefID :one
SELECT id, brief_id, company_name, company_overview, industry_sector, company_revenue, employee_count,
    geographic_footprint, parent_company, market_position, recent_news_events, financial_health,
    growth_trajectory, market_pressures, regulatory_environment, merger_acquisition_activity,
    competitive_landscape, created_at, updated_at
FROM company_intelligence
WHERE brief_id = $1;

-- name: UpdateCompanyIntelligence :one
UPDATE company_intelligence
SET company_name = $2,
    company_overview = $3,
    industry_sector = $4,
    company_revenue = $5,
    employee_count = $6,
    geographic_footprint = $7,
    parent_company = $8,
    market_position = $9,
    recent_news_events = $10,
    financial_health = $11,
    growth_trajectory = $12,
    market_pressures = $13,
    regulatory_environment = $14,
    merger_acquisition_activity = $15,
    competitive_landscape = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, company_name, company_overview, industry_sector, company_revenue, employee_count,
    geographic_footprint, parent_company, market_position, recent_news_events, financial_health,
    growth_trajectory, market_pressures, regulatory_environment, merger_acquisition_activity,
    competitive_landscape, created_at, updated_at;

-- name: DeleteCompanyIntelligence :exec
DELETE FROM company_intelligence
WHERE brief_id = $1;