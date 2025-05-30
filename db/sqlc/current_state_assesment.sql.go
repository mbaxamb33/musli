// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: current_state_assesment.sql

package db

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createCurrentStateAssessment = `-- name: CreateCurrentStateAssessment :one

INSERT INTO current_state_assessment (
    brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)
RETURNING id, brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues,
    created_at, updated_at
`

type CreateCurrentStateAssessmentParams struct {
	BriefID                     uuid.NullUUID  `json:"brief_id"`
	CurrentSolutionProvider     sql.NullString `json:"current_solution_provider"`
	CurrentSolutionSatisfaction sql.NullInt32  `json:"current_solution_satisfaction"`
	SpecificPainPoints          sql.NullString `json:"specific_pain_points"`
	WorkaroundSolutions         sql.NullString `json:"workaround_solutions"`
	CostOfStatusQuo             sql.NullString `json:"cost_of_status_quo"`
	SwitchingBarriers           sql.NullString `json:"switching_barriers"`
	ContractEndDates            sql.NullTime   `json:"contract_end_dates"`
	RenewalTiming               sql.NullString `json:"renewal_timing"`
	VendorRelationshipHealth    sql.NullString `json:"vendor_relationship_health"`
	SupportSatisfaction         sql.NullString `json:"support_satisfaction"`
	FunctionalityGaps           sql.NullString `json:"functionality_gaps"`
	PerformanceIssues           sql.NullString `json:"performance_issues"`
	ScalabilityConstraints      sql.NullString `json:"scalability_constraints"`
	IntegrationChallenges       sql.NullString `json:"integration_challenges"`
	UserAdoptionIssues          sql.NullString `json:"user_adoption_issues"`
}

// =============================================================================
// CURRENT STATE ASSESSMENT QUERIES
// =============================================================================
func (q *Queries) CreateCurrentStateAssessment(ctx context.Context, arg CreateCurrentStateAssessmentParams) (CurrentStateAssessment, error) {
	row := q.db.QueryRowContext(ctx, createCurrentStateAssessment,
		arg.BriefID,
		arg.CurrentSolutionProvider,
		arg.CurrentSolutionSatisfaction,
		arg.SpecificPainPoints,
		arg.WorkaroundSolutions,
		arg.CostOfStatusQuo,
		arg.SwitchingBarriers,
		arg.ContractEndDates,
		arg.RenewalTiming,
		arg.VendorRelationshipHealth,
		arg.SupportSatisfaction,
		arg.FunctionalityGaps,
		arg.PerformanceIssues,
		arg.ScalabilityConstraints,
		arg.IntegrationChallenges,
		arg.UserAdoptionIssues,
	)
	var i CurrentStateAssessment
	err := row.Scan(
		&i.ID,
		&i.BriefID,
		&i.CurrentSolutionProvider,
		&i.CurrentSolutionSatisfaction,
		&i.SpecificPainPoints,
		&i.WorkaroundSolutions,
		&i.CostOfStatusQuo,
		&i.SwitchingBarriers,
		&i.ContractEndDates,
		&i.RenewalTiming,
		&i.VendorRelationshipHealth,
		&i.SupportSatisfaction,
		&i.FunctionalityGaps,
		&i.PerformanceIssues,
		&i.ScalabilityConstraints,
		&i.IntegrationChallenges,
		&i.UserAdoptionIssues,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteCurrentStateAssessment = `-- name: DeleteCurrentStateAssessment :exec
DELETE FROM current_state_assessment
WHERE brief_id = $1
`

func (q *Queries) DeleteCurrentStateAssessment(ctx context.Context, briefID uuid.NullUUID) error {
	_, err := q.db.ExecContext(ctx, deleteCurrentStateAssessment, briefID)
	return err
}

const getCurrentStateAssessmentByBriefID = `-- name: GetCurrentStateAssessmentByBriefID :one
SELECT id, brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues,
    created_at, updated_at
FROM current_state_assessment
WHERE brief_id = $1
`

func (q *Queries) GetCurrentStateAssessmentByBriefID(ctx context.Context, briefID uuid.NullUUID) (CurrentStateAssessment, error) {
	row := q.db.QueryRowContext(ctx, getCurrentStateAssessmentByBriefID, briefID)
	var i CurrentStateAssessment
	err := row.Scan(
		&i.ID,
		&i.BriefID,
		&i.CurrentSolutionProvider,
		&i.CurrentSolutionSatisfaction,
		&i.SpecificPainPoints,
		&i.WorkaroundSolutions,
		&i.CostOfStatusQuo,
		&i.SwitchingBarriers,
		&i.ContractEndDates,
		&i.RenewalTiming,
		&i.VendorRelationshipHealth,
		&i.SupportSatisfaction,
		&i.FunctionalityGaps,
		&i.PerformanceIssues,
		&i.ScalabilityConstraints,
		&i.IntegrationChallenges,
		&i.UserAdoptionIssues,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const updateCurrentStateAssessment = `-- name: UpdateCurrentStateAssessment :one
UPDATE current_state_assessment
SET current_solution_provider = $2,
    current_solution_satisfaction = $3,
    specific_pain_points = $4,
    workaround_solutions = $5,
    cost_of_status_quo = $6,
    switching_barriers = $7,
    contract_end_dates = $8,
    renewal_timing = $9,
    vendor_relationship_health = $10,
    support_satisfaction = $11,
    functionality_gaps = $12,
    performance_issues = $13,
    scalability_constraints = $14,
    integration_challenges = $15,
    user_adoption_issues = $16,
    updated_at = CURRENT_TIMESTAMP
WHERE brief_id = $1
RETURNING id, brief_id, current_solution_provider, current_solution_satisfaction, specific_pain_points,
    workaround_solutions, cost_of_status_quo, switching_barriers, contract_end_dates,
    renewal_timing, vendor_relationship_health, support_satisfaction, functionality_gaps,
    performance_issues, scalability_constraints, integration_challenges, user_adoption_issues,
    created_at, updated_at
`

type UpdateCurrentStateAssessmentParams struct {
	BriefID                     uuid.NullUUID  `json:"brief_id"`
	CurrentSolutionProvider     sql.NullString `json:"current_solution_provider"`
	CurrentSolutionSatisfaction sql.NullInt32  `json:"current_solution_satisfaction"`
	SpecificPainPoints          sql.NullString `json:"specific_pain_points"`
	WorkaroundSolutions         sql.NullString `json:"workaround_solutions"`
	CostOfStatusQuo             sql.NullString `json:"cost_of_status_quo"`
	SwitchingBarriers           sql.NullString `json:"switching_barriers"`
	ContractEndDates            sql.NullTime   `json:"contract_end_dates"`
	RenewalTiming               sql.NullString `json:"renewal_timing"`
	VendorRelationshipHealth    sql.NullString `json:"vendor_relationship_health"`
	SupportSatisfaction         sql.NullString `json:"support_satisfaction"`
	FunctionalityGaps           sql.NullString `json:"functionality_gaps"`
	PerformanceIssues           sql.NullString `json:"performance_issues"`
	ScalabilityConstraints      sql.NullString `json:"scalability_constraints"`
	IntegrationChallenges       sql.NullString `json:"integration_challenges"`
	UserAdoptionIssues          sql.NullString `json:"user_adoption_issues"`
}

func (q *Queries) UpdateCurrentStateAssessment(ctx context.Context, arg UpdateCurrentStateAssessmentParams) (CurrentStateAssessment, error) {
	row := q.db.QueryRowContext(ctx, updateCurrentStateAssessment,
		arg.BriefID,
		arg.CurrentSolutionProvider,
		arg.CurrentSolutionSatisfaction,
		arg.SpecificPainPoints,
		arg.WorkaroundSolutions,
		arg.CostOfStatusQuo,
		arg.SwitchingBarriers,
		arg.ContractEndDates,
		arg.RenewalTiming,
		arg.VendorRelationshipHealth,
		arg.SupportSatisfaction,
		arg.FunctionalityGaps,
		arg.PerformanceIssues,
		arg.ScalabilityConstraints,
		arg.IntegrationChallenges,
		arg.UserAdoptionIssues,
	)
	var i CurrentStateAssessment
	err := row.Scan(
		&i.ID,
		&i.BriefID,
		&i.CurrentSolutionProvider,
		&i.CurrentSolutionSatisfaction,
		&i.SpecificPainPoints,
		&i.WorkaroundSolutions,
		&i.CostOfStatusQuo,
		&i.SwitchingBarriers,
		&i.ContractEndDates,
		&i.RenewalTiming,
		&i.VendorRelationshipHealth,
		&i.SupportSatisfaction,
		&i.FunctionalityGaps,
		&i.PerformanceIssues,
		&i.ScalabilityConstraints,
		&i.IntegrationChallenges,
		&i.UserAdoptionIssues,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
