package postgres

import (
	"context"
	"database/sql"
)

type  NewDelivery struct{
	DB *sql.DB
}

// constructor 
func DeliveryStore (db *sql.DB) *NewDelivery{
	return &NewDelivery{
		DB : db,
	}
}
func (s *NewDelivery) tryCreate(ctx context.Context , deliveryId string )(bool, error){
	var returnVal string 
	err :=s.DB.QueryRowContext(ctx , `
	INSERT INTO webhook_deliveries (delivery_id)
		VALUES ($1)
		ON CONFLICT (delivery_id) DO NOTHING
		RETURNING delivery_id`,deliveryId).Scan(&returnVal);

		if err == sql.ErrNoRows{
			// if no rows return 
			return false , nil
		}
		if err !=nil {
			return false , err
		}
		return true , nil 
}