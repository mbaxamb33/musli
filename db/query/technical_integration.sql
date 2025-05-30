-- =============================================================================
-- TECHNICAL & INTEGRATION QUERIES
-- =============================================================================

-- name: CreateTechnicalIntegration :one
INSERT INTO technical_integration (
    brief_id, technical_architecture, security_requirements, compliance_standards, integration_points,
    data_migration_scope, customization_needs, scalability_requirements, performance_benchmarks,
    disaster_recovery_needs, backup_requirements, access_control_requirements, audit_trail_needs,
    reporting_capabilities, api_requirements, mobile_access_needs
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, technical_architecture, security_requirements, compliance_standards, integration_points,
    data_migration_scope, customization_needs, scalability_requirements, performance_benchmarks,
    disaster_recovery_needs, backup_requirements, access_control_requirements, audit_trail_needs,
    reporting_capabilities, api_requirements, mobile_access_needs, created_at, updated_at;

-- name: GetTechnicalIntegrationByBriefID :one
SELECT id, brief_id, technical_architecture, security_requirements, compliance_standards, integration_points,
    data_migration_scope, customization_needs, scalability_requirements, performance_benchmarks,
    disaster_recovery_needs, backup_requirements, access_control_requirements, audit_trail_needs,
    reporting_capabilities, api_requirements, mobile_access_needs, created_at, updated_at
FROM technical_integration
WHERE brief_id = $1;

-- name: UpdateTechnicalIntegration :one
UPDATE technical_integration
SET technical_architecture = $2,
    security_requirements = $3,
    compliance_standards = $4,
    integration_points = $5,
    data_migration_scope = $6,
    customization_needs = $7,
    scalability_requirements = $8,
    performance_benchmarks = $9,
    disaster_recovery_needs = $10,
    backup_requirements = $11,
    access_control_requirements = $12,
    audit_trail_needs = $13,
    reporting_capabilities = $14,
    api_requirements = $15,
    mobile_access_needs = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, technical_architecture, security_requirements, compliance_standards, integration_points,
    data_migration_scope, customization_needs, scalability_requirements, performance_benchmarks,
    disaster_recovery_needs, backup_requirements, access_control_requirements, audit_trail_needs,
    reporting_capabilities, api_requirements, mobile_access_needs, created_at, updated_at;

-- name: DeleteTechnicalIntegration :exec
DELETE FROM technical_integration
WHERE brief_id = $1;