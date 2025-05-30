-- =============================================================================
-- FINANCIAL & PROCUREMENT QUERIES
-- =============================================================================

-- name: CreateFinancialProcurement :one
INSERT INTO financial_procurement (
    brief_id, total_available_budget, budget_source, budget_approval_workflow, procurement_process,
    purchasing_policies, payment_terms_constraints, financial_approval_levels, budget_cycle_timing,
    cost_justification_requirements, roi_calculation_method, payback_period_expectations,
    financing_options, contract_terms_requirements, legal_review_process, insurance_requirements
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, total_available_budget, budget_source, budget_approval_workflow, procurement_process,
    purchasing_policies, payment_terms_constraints, financial_approval_levels, budget_cycle_timing,
    cost_justification_requirements, roi_calculation_method, payback_period_expectations,
    financing_options, contract_terms_requirements, legal_review_process, insurance_requirements,
    created_at, updated_at;

-- name: GetFinancialProcurementByBriefID :one
SELECT id, brief_id, total_available_budget, budget_source, budget_approval_workflow, procurement_process,
    purchasing_policies, payment_terms_constraints, financial_approval_levels, budget_cycle_timing,
    cost_justification_requirements, roi_calculation_method, payback_period_expectations,
    financing_options, contract_terms_requirements, legal_review_process, insurance_requirements,
    created_at, updated_at
FROM financial_procurement
WHERE brief_id = $1;

-- name: UpdateFinancialProcurement :one
UPDATE financial_procurement
SET total_available_budget = $2,
    budget_source = $3,
    budget_approval_workflow = $4,
    procurement_process = $5,
    purchasing_policies = $6,
    payment_terms_constraints = $7,
    financial_approval_levels = $8,
    budget_cycle_timing = $9,
    cost_justification_requirements = $10,
    roi_calculation_method = $11,
    payback_period_expectations = $12,
    financing_options = $13,
    contract_terms_requirements = $14,
    legal_review_process = $15,
    insurance_requirements = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, total_available_budget, budget_source, budget_approval_workflow, procurement_process,
    purchasing_policies, payment_terms_constraints, financial_approval_levels, budget_cycle_timing,
    cost_justification_requirements, roi_calculation_method, payback_period_expectations,
    financing_options, contract_terms_requirements, legal_review_process, insurance_requirements,
    created_at, updated_at;

-- name: DeleteFinancialProcurement :exec
DELETE FROM financial_procurement
WHERE brief_id = $1;