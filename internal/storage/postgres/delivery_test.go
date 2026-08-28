package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// Unit Test — validates delivery insertion outcomes. [Delivery Idempotency]
func TestDeliveryStoreTryCreate(t *testing.T) {
	tests := []struct {
		name    string
		result  *sqlmock.Rows
		err     error
		wantNew bool
		wantErr bool
	}{
		{name: "[New Delivery] returned ID marks new", result: sqlmock.NewRows([]string{"delivery_id"}).AddRow("delivery-1"), wantNew: true},
		{name: "[Duplicate Delivery] no returned row marks existing", result: sqlmock.NewRows([]string{"delivery_id"}), wantNew: false},
		{name: "[Database Error] query error propagates", err: errors.New("database unavailable"), wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, mock, err := sqlmock.New()
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			expectation := mock.ExpectQuery(regexp.QuoteMeta(`
	    INSERT INTO webhook_deliveries (delivery_id)
	    VALUES ($1)
	    ON CONFLICT (delivery_id) DO NOTHING
	    RETURNING delivery_id
	`)).WithArgs("delivery-1")
			if tt.err != nil {
				expectation.WillReturnError(tt.err)
			} else {
				expectation.WillReturnRows(tt.result)
			}

			gotNew, gotErr := NewDeliveryStore(db).TryCreate(context.Background(), "delivery-1")
			if gotNew != tt.wantNew || (gotErr != nil) != tt.wantErr {
				t.Fatalf("TryCreate() = (%t, %v), want (%t, error=%t)", gotNew, gotErr, tt.wantNew, tt.wantErr)
			}
			if err := mock.ExpectationsWereMet(); err != nil {
				t.Fatal(err)
			}
		})
	}
}
