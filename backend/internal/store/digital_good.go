package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
)

func scanDigitalGood(row interface{ Scan(dest ...any) error }) (*DigitalGood, error) {
	var good DigitalGood
	var description sql.NullString
	var category sql.NullString
	var gameID sql.NullString
	if err := row.Scan(&good.ID, &good.Name, &description, &category, &gameID, &good.CreatedAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}
		return nil, err
	}
	if description.Valid {
		good.Description = &description.String
	}
	if category.Valid {
		good.Category = &category.String
	}
	if gameID.Valid {
		id, err := uuid.Parse(gameID.String)
		if err != nil {
			return nil, err
		}
		good.GameID = &id
	}
	return &good, nil
}

const digitalGoodColumns = `id, name, description, category, game_id, created_at`

func (s *Store) ListDigitalGoods(ctx context.Context, gameID *uuid.UUID) ([]DigitalGood, error) {
	var (
		rows *sql.Rows
		err  error
	)
	if gameID != nil {
		rows, err = s.db.QueryContext(ctx, `
			SELECT `+digitalGoodColumns+`
			FROM digital_goods
			WHERE game_id = $1
			ORDER BY name ASC
		`, *gameID)
	} else {
		rows, err = s.db.QueryContext(ctx, `
			SELECT `+digitalGoodColumns+`
			FROM digital_goods
			ORDER BY name ASC
		`)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []DigitalGood
	for rows.Next() {
		var good DigitalGood
		var description sql.NullString
		var category sql.NullString
		var rowGameID sql.NullString
		if err := rows.Scan(&good.ID, &good.Name, &description, &category, &rowGameID, &good.CreatedAt); err != nil {
			return nil, err
		}
		if description.Valid {
			good.Description = &description.String
		}
		if category.Valid {
			good.Category = &category.String
		}
		if rowGameID.Valid {
			id, err := uuid.Parse(rowGameID.String)
			if err != nil {
				return nil, err
			}
			good.GameID = &id
		}
		goods = append(goods, good)
	}
	return goods, rows.Err()
}

func (s *Store) GetDigitalGoodByID(ctx context.Context, id uuid.UUID) (*DigitalGood, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT `+digitalGoodColumns+`
		FROM digital_goods
		WHERE id = $1
	`, id)
	return scanDigitalGood(row)
}

func (s *Store) CreateDigitalGood(ctx context.Context, name string, description *string, gameID *uuid.UUID) (*DigitalGood, error) {
	var gameIDValue any
	if gameID != nil {
		gameIDValue = *gameID
	}

	row := s.db.QueryRowContext(ctx, `
		INSERT INTO digital_goods (name, description, game_id)
		VALUES ($1, $2, $3)
		RETURNING `+digitalGoodColumns+`
	`, name, description, gameIDValue)
	return scanDigitalGood(row)
}
