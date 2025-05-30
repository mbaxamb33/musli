// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: financial_procurement.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createFinancialProcurement = `-- name: CreateFinancialProcurement :one

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
    created_at, updated_at
`

type CreateFinancialProcurementParams struct {
	BriefID                       uuid.NullUUID  `json:"brief_id"`
	TotalAvailableBudget          sql.NullString `json:"total_available_budget"`
	BudgetSource                  sql.NullString `json:"budget_source"`
	BudgetApprovalWorkflow        sql.NullString `json:"budget_approval_workflow"`
	ProcurementProcess            sql.NullString `json:"procurement_process"`
	PurchasingPolicies            sql.NullString `json:"purchasing_policies"`
	PaymentTermsConstraints       sql.NullString `json:"payment_terms_constraints"`
	FinancialApprovalLevels       sql.NullString `json:"financial_approval_levels"`
	BudgetCycleTiming             sql.NullString `json:"budget_cycle_timing"`
	CostJustificationRequirements sql.NullString `json:"cost_justification_requirements"`
	RoiCalculationMethod          sql.NullString `json:"roi_calculation_method"`
	PaybackPeriodExpectations     sql.NullString `json:"payback_period_expectations"`
	FinancingOptions              sql.NullString `json:"financing_options"`
	ContractTermsRequirements     sql.NullString `json:"contract_terms_requirements"`
	LegalReviewProcess            sql.NullString `json:"legal_review_process"`
	InsuranceRequirements         sql.NullString `json:"insurance_requirements"`
}

// =============================================================================
// FINANCIAL & PROCUREMENT QUERIES
// =============================================================================
func (q *Queries) CreateFinancialProcurement(ctx context.Context, arg CreateFinancialProcurementParams) (FinancialProcurement, error) {
	row := q.db.QueryRowContext(ctx, createFinancialProcurement,
		arg.BriefID,
		arg.TotalAvailableBudget,
		arg.BudgetSource,
		arg.BudgetApprovalWorkflow,
		arg.ProcurementProcess,
		arg.PurchasingPolicies,
		arg.PaymentTermsConstraints,
		arg.FinancialApprovalLevels,
		arg.BudgetCycleTiming,
		arg.CostJustificationRequirements,
		arg.RoiCalculationMethod,
		arg.PaybackPeriodExpectations,
		arg.FinancingOptions,
		arg.ContractTermsRequirements,
		arg.LegalReviewProcess,
		arg.InsuranceRequirements,
	)
	var i FinancialProcurement
	err := row.Scan(
		&i.ID,
		&i.BriefID,
		&i.TotalAvailableBudget,
		&i.BudgetSource,
		&i.BudgetApprovalWorkflow,
		&i.ProcurementProcess,
		&i.PurchasingPolicies,
		&i.PaymentTermsConstraints,
		&i.FinancialApprovalLevels,
		&i.BudgetCycleTiming,
		&i.CostJustificationRequirements,
		&i.RoiCalculationMethod,
		&i.PaybackPeriodExpectations,
		&i.FinancingOptions,
		&i.ContractTermsRequirements,
		&i.LegalReviewProcess,
		&i.InsuranceRequirements,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteFinancialProcurement = `-- name: DeleteFinancialProcurement :exec
DELETE FROM financial_procurement
WHERE brief_id = $1
`

func (q *Queries) DeleteFinancialProcurement(ctx context.Context, briefID uuid.NullUUID) error {
	_, err := q.db.ExecContext(ctx, deleteFinancialProcurement, briefID)
	return err
}

const getFinancialProcurementByBriefID = `-- name: GetFinancialProcurementByBriefID :one
SELECT id, brief_id, total_available_budget, budget_source, budget_approval_workflow, procurement_process,
    purchasing_policies, payment_terms_constraints, financial_approval_levels, budget_cycle_timing,
    cost_justification_requirements, roi_calculation_method, payback_period_expectations,
    financing_options, contract_terms_requirements, legal_review_process, insurance_requirements,
    created_at, updated_at
FROM financial_procurement
WHERE brief_id = $1
`

func (q *Queries) GetFinancialProcurementByBriefID(ctx context.Context, briefID uuid.NullUUID) (FinancialProcurement, error) {
	row := q.db.QueryRowContext(ctx, getFinancialProcurementByBriefID, briefID)
	var i FinancialProcurement
	err := row.Scan(
		&i.ID,
		&i.BriefID,
		&i.TotalAvailableBudget,
		&i.BudgetSource,
		&i.BudgetApprovalWorkflow,
		&i.ProcurementProcess,
		&i.PurchasingPolicies,
		&i.PaymentTermsConstraints,
		&i.FinancialApprovalLevels,
		&i.BudgetCycleTiming,
		&i.CostJustificationRequirements,
		&i.RoiCalculationMethod,
		&i.PaybackPeriodExpectations,
		&i.FinancingOptions,
		&i.ContractTermsRequirements,
		&i.LegalReviewProcess,
		&i.InsuranceRequirements,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateFinancialProcurement = `-- name: UpdateFinancialProcurement :one
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
    created_at, updated_at
`

type UpdateFinancialProcurementParams struct {
	BriefID                       uuid.NullUUID  `json:"brief_id"`
	TotalAvailableBudget          sql.NullString `json:"total_available_budget"`
	BudgetSource                  sql.NullString `json:"budget_source"`
	BudgetApprovalWorkflow        sql.NullString `json:"budget_approval_workflow"`
	ProcurementProcess            sql.NullString `json:"procurement_process"`
	PurchasingPolicies            sql.NullString `json:"purchasing_policies"`
	PaymentTermsConstraints       sql.NullString `json:"payment_terms_constraints"`
	FinancialApprovalLevels       sql.NullString `json:"financial_approval_levels"`
	BudgetCycleTiming             sql.NullString `json:"budget_cycle_timing"`
	CostJustificationRequirements sql.NullString `json:"cost_justification_requirements"`
	RoiCalculationMethod          sql.NullString `json:"roi_calculation_method"`
	PaybackPeriodExpectations     sql.NullString `json:"payback_period_expectations"`
	FinancingOptions              sql.NullString `json:"financing_options"`
	ContractTermsRequirements     sql.NullString `json:"contract_terms_requirements"`
	LegalReviewProcess            sql.NullString `json:"legal_review_process"`
	InsuranceRequirements         sql.NullString `json:"insurance_requirements"`
}

func (q *Queries) UpdateFinancialProcurement(ctx context.Context, arg UpdateFinancialProcurementParams) (FinancialProcurement, error) {
	row := q.db.QueryRowContext(ctx, updateFinancialProcurement,
		arg.BriefID,
		arg.TotalAvailableBudget,
		arg.BudgetSource,
		arg.BudgetApprovalWorkflow,
		arg.ProcurementProcess,
		arg.PurchasingPolicies,
		arg.PaymentTermsConstraints,
		arg.FinancialApprovalLevels,
		arg.BudgetCycleTiming,
		arg.CostJustificationRequirements,
		arg.RoiCalculationMethod,
		arg.PaybackPeriodExpectations,
		arg.FinancingOptions,
		arg.ContractTermsRequirements,
		arg.LegalReviewProcess,
		arg.InsuranceRequirements,
	)
	var i FinancialProcurement
	err := row.Scan(
		&i.ID,
		&i.BriefID,
		&i.TotalAvailableBudget,
		&i.BudgetSource,
		&i.BudgetApprovalWorkflow,
		&i.ProcurementProcess,
		&i.PurchasingPolicies,
		&i.PaymentTermsConstraints,
		&i.FinancialApprovalLevels,
		&i.BudgetCycleTiming,
		&i.CostJustificationRequirements,
		&i.RoiCalculationMethod,
		&i.PaybackPeriodExpectations,
		&i.FinancingOptions,
		&i.ContractTermsRequirements,
		&i.LegalReviewProcess,
		&i.InsuranceRequirements,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
