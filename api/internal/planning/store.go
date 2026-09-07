package planning

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound is returned when a plan lookup finds no plan for that date.
var ErrNotFound = errors.New("plan not found")

// Store persists daily plans in Postgres.
type Store struct {
	db *sql.DB
}

// NewStore returns a Store backed by db.
func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func emptyPlan(date string) Plan {
	return Plan{Date: date, TaskIDs: []string{}}
}

// GetByDate returns the plan for date, or a zero-value unstarted Plan if
// none has been created yet.
func (s *Store) GetByDate(ctx context.Context, date string) (Plan, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT plan_date, intention, started_at, review_notes, reviewed_at, created_at, updated_at
		FROM daily_plans WHERE plan_date = $1
	`, date)

	var p Plan
	var planDate sql.NullTime
	err := row.Scan(&planDate, &p.Intention, &p.StartedAt, &p.ReviewNotes, &p.ReviewedAt, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return emptyPlan(date), nil
	}
	if err != nil {
		return Plan{}, fmt.Errorf("get plan: %w", err)
	}

	p.Date = planDate.Time.Format(dateLayout)
	p.Started = p.StartedAt != nil
	p.Reviewed = p.ReviewedAt != nil

	taskIDs, err := s.taskIDsForDate(ctx, date)
	if err != nil {
		return Plan{}, err
	}
	p.TaskIDs = taskIDs

	return p, nil
}

func (s *Store) taskIDsForDate(ctx context.Context, date string) ([]string, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT dpt.task_id
		FROM daily_plan_tasks dpt
		JOIN daily_plans dp ON dp.id = dpt.plan_id
		WHERE dp.plan_date = $1
		ORDER BY dpt.position
	`, date)
	if err != nil {
		return nil, fmt.Errorf("list plan tasks: %w", err)
	}
	defer rows.Close()

	ids := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("list plan tasks: %w", err)
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Start creates or updates the plan for date: sets its intention and task
// list, and stamps started_at if it hasn't already been started.
func (s *Store) Start(ctx context.Context, date, intention string, taskIDs []string) (Plan, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return Plan{}, fmt.Errorf("start plan: %w", err)
	}
	defer tx.Rollback()

	var planID string
	err = tx.QueryRowContext(ctx, `
		INSERT INTO daily_plans (plan_date, intention, started_at)
		VALUES ($1, $2, now())
		ON CONFLICT (plan_date) DO UPDATE
		SET intention = EXCLUDED.intention,
		    started_at = COALESCE(daily_plans.started_at, EXCLUDED.started_at),
		    updated_at = now()
		RETURNING id
	`, date, intention).Scan(&planID)
	if err != nil {
		return Plan{}, fmt.Errorf("start plan: %w", err)
	}

	if _, err := tx.ExecContext(ctx, `DELETE FROM daily_plan_tasks WHERE plan_id = $1`, planID); err != nil {
		return Plan{}, fmt.Errorf("start plan: %w", err)
	}

	for i, taskID := range taskIDs {
		if _, err := tx.ExecContext(ctx, `
			INSERT INTO daily_plan_tasks (plan_id, task_id, position) VALUES ($1, $2, $3)
		`, planID, taskID, i); err != nil {
			return Plan{}, fmt.Errorf("start plan: attach task %s: %w", taskID, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return Plan{}, fmt.Errorf("start plan: %w", err)
	}

	return s.GetByDate(ctx, date)
}

// Review sets the review notes and reviewed_at timestamp on an
// already-started plan for date. It returns ErrNotFound if no plan for that
// date has been started yet.
func (s *Store) Review(ctx context.Context, date, notes string) (Plan, error) {
	res, err := s.db.ExecContext(ctx, `
		UPDATE daily_plans
		SET review_notes = $1, reviewed_at = now(), updated_at = now()
		WHERE plan_date = $2
	`, notes, date)
	if err != nil {
		return Plan{}, fmt.Errorf("review plan: %w", err)
	}
	n, err := res.RowsAffected()
	if err != nil {
		return Plan{}, fmt.Errorf("review plan: %w", err)
	}
	if n == 0 {
		return Plan{}, ErrNotFound
	}

	return s.GetByDate(ctx, date)
}
