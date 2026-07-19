package things

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ThingsRepo struct {
	db *pgxpool.Pool
}

func NewRepo(db *pgxpool.Pool) *ThingsRepo {
	return &ThingsRepo{
		db: db,
	}
}

type Things struct {
	Id          uuid.UUID `json:"id"`
	UserId      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	ExpiredAt   time.Time `json:"expires_at"`
	CreatedAt   time.Time `json:"created_at"`
}

func (r *ThingsRepo) GetThingsOfUser(ctx context.Context, userId uuid.UUID) ([]Things, error) {
	const statement = `
		SELECT
			id,
			user_id,
			name,
			description,
			quantity,
			expires_at,
			created_at
		FROM things
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(ctx, statement, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	thingsList := make([]Things, 0)

	for rows.Next() {
		var things Things

		err := rows.Scan(
			&things.Id,
			&things.UserId,
			&things.Name,
			&things.Description,
			&things.Quantity,
			&things.ExpiredAt,
			&things.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		thingsList = append(thingsList, things)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return thingsList, nil
}

func (r *ThingsRepo) GetThingsById(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID) (Things, error) {
	const statement = `
		SELECT
			id,
			user_id,
			name,
			description,
			quantity,
			expires_at,
			created_at
		FROM things
		WHERE id = $1 AND user_id = $2
	`

	var things Things
	err := r.db.QueryRow(ctx, statement, thingsId, userId).Scan(
		&things.Id,
		&things.UserId,
		&things.Name,
		&things.Description,
		&things.Quantity,
		&things.ExpiredAt,
		&things.CreatedAt,
	)
	if err != nil {
		return Things{}, err
	}

	return things, nil
}

func (r *ThingsRepo) CreateThings(ctx context.Context, userId uuid.UUID, name string, description string, quantity int, expiredAt time.Time) (Things, error) {
	const statement = `
		INSERT INTO things (user_id, name, description, quantity, expires_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, user_id, name, description, quantity, expires_at, created_at
	`

	var things Things
	err := r.db.QueryRow(ctx, statement, userId, name, description, quantity, expiredAt).Scan(
		&things.Id,
		&things.UserId,
		&things.Name,
		&things.Description,
		&things.Quantity,
		&things.ExpiredAt,
		&things.CreatedAt,
	)
	if err != nil {
		return Things{}, err
	}

	return things, nil
}

func (r *ThingsRepo) DeleteThings(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID) error {
	const statement = `
		DELETE FROM things
		WHERE id = $1 AND user_id = $2
	`

	tag, err := r.db.Exec(ctx, statement, thingsId, userId)
	if err != nil {
		return err
	}

	// No row deleted means it does not exist (or belongs to another user).
	if tag.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}

	return nil
}

func (r *ThingsRepo) UpdateThings(ctx context.Context, userId uuid.UUID, thingsId uuid.UUID, name string, description string, quantity int, expiredAt time.Time) (Things, error) {
	const statement = `
		UPDATE things
		SET name = $1, description = $2, quantity = $3, expires_at = $4
		WHERE id = $5 AND user_id = $6
		RETURNING id, user_id, name, description, quantity, expires_at, created_at
	`

	var things Things
	err := r.db.QueryRow(ctx, statement, name, description, quantity, expiredAt, thingsId, userId).Scan(
		&things.Id,
		&things.UserId,
		&things.Name,
		&things.Description,
		&things.Quantity,
		&things.ExpiredAt,
		&things.CreatedAt,
	)
	if err != nil {
		return Things{}, err
	}

	return things, nil
}
