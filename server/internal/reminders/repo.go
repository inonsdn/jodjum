package reminders

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ONE_TIME_REMINDER_TYPE = 1
	DAILY_REMINDER_TYPE    = 2
	MONTHLY_REMINDER_TYPE  = 3
	YEARLY_REMINDER_TYPE   = 4
)

type RemindersRepo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *RemindersRepo {
	return &RemindersRepo{
		db: db,
	}
}

type Reminder struct {
	Id              uuid.UUID `json:"id"`
	UserId          uuid.UUID `json:"user_id"`
	RemindTimestamp time.Time `json:"remind_timestamp"`
	ReminderType    int       `json:"reminder_type"`
	IsActive        bool      `json:"is_active"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
}

func (r *RemindersRepo) GetRemindersOfUser(ctx context.Context, userId uuid.UUID) ([]Reminder, error) {
	const statement = `
		SELECT
			id,
			user_id,
			remind_timestamp,
			reminder_type,
			is_active,
			name,
			description
		FROM reminders
		WHERE user_id = $1
		ORDER BY remind_timestamp ASC
	`

	rows, err := r.db.Query(ctx, statement, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reminderList := make([]Reminder, 0)

	for rows.Next() {
		var reminder Reminder

		err := rows.Scan(
			&reminder.Id,
			&reminder.UserId,
			&reminder.RemindTimestamp,
			&reminder.ReminderType,
			&reminder.IsActive,
			&reminder.Name,
			&reminder.Description,
		)
		if err != nil {
			return nil, err
		}

		reminderList = append(reminderList, reminder)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return reminderList, nil
}

func (r *RemindersRepo) GetReminderById(ctx context.Context, userId uuid.UUID, reminderId uuid.UUID) (Reminder, error) {
	const statement = `
		SELECT
			id,
			user_id,
			remind_timestamp,
			reminder_type,
			is_active,
			name,
			description
		FROM reminders
		WHERE id = $1 AND user_id = $2
	`

	var reminder Reminder
	err := r.db.QueryRow(ctx, statement, reminderId, userId).Scan(
		&reminder.Id,
		&reminder.UserId,
		&reminder.RemindTimestamp,
		&reminder.ReminderType,
		&reminder.IsActive,
		&reminder.Name,
		&reminder.Description,
	)
	if err != nil {
		return Reminder{}, err
	}

	return reminder, nil
}

func (r *RemindersRepo) CreateReminder(ctx context.Context, userId uuid.UUID, remindTimestamp time.Time, reminderType int, isActive bool, name string, description string) (Reminder, error) {
	const statement = `
		INSERT INTO reminders (user_id, remind_timestamp, reminder_type, is_active, name, description)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, user_id, remind_timestamp, reminder_type, is_active, name, description
	`

	var reminder Reminder
	err := r.db.QueryRow(ctx, statement, userId, remindTimestamp, reminderType, isActive, name, description).Scan(
		&reminder.Id,
		&reminder.UserId,
		&reminder.RemindTimestamp,
		&reminder.ReminderType,
		&reminder.IsActive,
		&reminder.Name,
		&reminder.Description,
	)
	if err != nil {
		return Reminder{}, err
	}

	return reminder, nil
}

func (r *RemindersRepo) UpdateReminder(ctx context.Context, userId uuid.UUID, reminderId uuid.UUID, remindTimestamp time.Time, reminderType int, isActive bool, name string, description string) (Reminder, error) {
	const statement = `
		UPDATE reminders
		SET remind_timestamp = $1, reminder_type = $2, is_active = $3, name = $4, description = $5
		WHERE id = $6 AND user_id = $7
		RETURNING id, user_id, remind_timestamp, reminder_type, is_active, name, description
	`

	var reminder Reminder
	err := r.db.QueryRow(ctx, statement, remindTimestamp, reminderType, isActive, name, description, reminderId, userId).Scan(
		&reminder.Id,
		&reminder.UserId,
		&reminder.RemindTimestamp,
		&reminder.ReminderType,
		&reminder.IsActive,
		&reminder.Name,
		&reminder.Description,
	)
	if err != nil {
		return Reminder{}, err
	}

	return reminder, nil
}

func (r *RemindersRepo) GetNextReminder(ctx context.Context) (Reminder, error) {
	const statement = `
		SELECT
			id,
			user_id,
			remind_timestamp,
			reminder_type,
			is_active,
			name,
			description
		FROM reminders
		WHERE is_active = TRUE
		ORDER BY remind_timestamp ASC
		LIMIT 1
	`

	var reminder Reminder
	err := r.db.QueryRow(ctx, statement).Scan(
		&reminder.Id,
		&reminder.UserId,
		&reminder.RemindTimestamp,
		&reminder.ReminderType,
		&reminder.IsActive,
		&reminder.Name,
		&reminder.Description,
	)
	if err != nil {
		return Reminder{}, err
	}

	return reminder, nil
}
