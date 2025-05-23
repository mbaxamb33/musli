// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: analysis_contact_info.sql

package db

import (
	"context"
	"database/sql"
)

const createAnalysisContactInfo = `-- name: CreateAnalysisContactInfo :one
INSERT INTO analysis_contact_info (
    analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information, updated_at
`

type CreateAnalysisContactInfoParams struct {
	AnalysisID          int32          `json:"analysis_id"`
	Problems            sql.NullString `json:"problems"`
	Needs               sql.NullString `json:"needs"`
	Urgency             sql.NullString `json:"urgency"`
	Priorities          sql.NullString `json:"priorities"`
	DecisionProcess     sql.NullString `json:"decision_process"`
	Budget              sql.NullString `json:"budget"`
	Resources           sql.NullString `json:"resources"`
	RelevantInformation sql.NullString `json:"relevant_information"`
}

func (q *Queries) CreateAnalysisContactInfo(ctx context.Context, arg CreateAnalysisContactInfoParams) (AnalysisContactInfo, error) {
	row := q.db.QueryRowContext(ctx, createAnalysisContactInfo,
		arg.AnalysisID,
		arg.Problems,
		arg.Needs,
		arg.Urgency,
		arg.Priorities,
		arg.DecisionProcess,
		arg.Budget,
		arg.Resources,
		arg.RelevantInformation,
	)
	var i AnalysisContactInfo
	err := row.Scan(
		&i.AnalysisID,
		&i.Problems,
		&i.Needs,
		&i.Urgency,
		&i.Priorities,
		&i.DecisionProcess,
		&i.Budget,
		&i.Resources,
		&i.RelevantInformation,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteAnalysisContactInfo = `-- name: DeleteAnalysisContactInfo :exec
DELETE FROM analysis_contact_info
WHERE analysis_id = $1
`

func (q *Queries) DeleteAnalysisContactInfo(ctx context.Context, analysisID int32) error {
	_, err := q.db.ExecContext(ctx, deleteAnalysisContactInfo, analysisID)
	return err
}

const getAnalysisContactInfo = `-- name: GetAnalysisContactInfo :one
SELECT analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information, updated_at
FROM analysis_contact_info
WHERE analysis_id = $1
`

func (q *Queries) GetAnalysisContactInfo(ctx context.Context, analysisID int32) (AnalysisContactInfo, error) {
	row := q.db.QueryRowContext(ctx, getAnalysisContactInfo, analysisID)
	var i AnalysisContactInfo
	err := row.Scan(
		&i.AnalysisID,
		&i.Problems,
		&i.Needs,
		&i.Urgency,
		&i.Priorities,
		&i.DecisionProcess,
		&i.Budget,
		&i.Resources,
		&i.RelevantInformation,
		&i.UpdatedAt,
	)
	return i, err
}

const updateAnalysisContactInfo = `-- name: UpdateAnalysisContactInfo :one
UPDATE analysis_contact_info
SET problems = $2,
    needs = $3,
    urgency = $4,
    priorities = $5,
    decision_process = $6,
    budget = $7,
    resources = $8,
    relevant_information = $9,
    updated_at = CURRENT_TIMESTAMP
WHERE analysis_id = $1
RETURNING analysis_id, problems, needs, urgency, priorities, decision_process, budget, resources, relevant_information, updated_at
`

type UpdateAnalysisContactInfoParams struct {
	AnalysisID          int32          `json:"analysis_id"`
	Problems            sql.NullString `json:"problems"`
	Needs               sql.NullString `json:"needs"`
	Urgency             sql.NullString `json:"urgency"`
	Priorities          sql.NullString `json:"priorities"`
	DecisionProcess     sql.NullString `json:"decision_process"`
	Budget              sql.NullString `json:"budget"`
	Resources           sql.NullString `json:"resources"`
	RelevantInformation sql.NullString `json:"relevant_information"`
}

func (q *Queries) UpdateAnalysisContactInfo(ctx context.Context, arg UpdateAnalysisContactInfoParams) (AnalysisContactInfo, error) {
	row := q.db.QueryRowContext(ctx, updateAnalysisContactInfo,
		arg.AnalysisID,
		arg.Problems,
		arg.Needs,
		arg.Urgency,
		arg.Priorities,
		arg.DecisionProcess,
		arg.Budget,
		arg.Resources,
		arg.RelevantInformation,
	)
	var i AnalysisContactInfo
	err := row.Scan(
		&i.AnalysisID,
		&i.Problems,
		&i.Needs,
		&i.Urgency,
		&i.Priorities,
		&i.DecisionProcess,
		&i.Budget,
		&i.Resources,
		&i.RelevantInformation,
		&i.UpdatedAt,
	)
	return i, err
}
