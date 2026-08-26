package postgres

import (
	"context"
	"database/sql"
)

type DeliveryStore struct {
	DB *sql.DB
}

// constructor
func NewDeliveryStore(db *sql.DB) *DeliveryStore {
	return &DeliveryStore{
		DB: db,
	}
}

func (s *DeliveryStore) TryCreate(
	ctx context.Context,
	deliveryID string,
) (bool, error) {

	var insertedID string

	err := s.DB.QueryRowContext(ctx, `
        INSERT INTO webhook_deliveries (delivery_id)
        VALUES ($1)
        ON CONFLICT (delivery_id) DO NOTHING
        RETURNING delivery_id
    `, deliveryID).Scan(&insertedID)

	if err == sql.ErrNoRows {
		// Delivery already exists.
		return false, nil
	}

	if err != nil {
		return false, err
	}
	return true, nil
}
