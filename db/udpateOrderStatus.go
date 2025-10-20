package db

import (
	"AdminPanelAPI/apperrors"
	"context"
	"fmt"

	ekc_db "github.com/Hugokoks/kratomclub-go-common/db"
)

func UpdateOrderStatus(ctx context.Context, orderID string, newStatus string) error {

	query := `update orders set status = $1 WHERE number $2;`

	cmdTag, err := ekc_db.Pool.Exec(ctx, query, newStatus, orderID)
	if err != nil {
		return err
	}
	if cmdTag.RowsAffected() == 0 {

		return fmt.Errorf("objednávka s ID %s nebyla nalezena: %w", orderID, apperrors.ErrOrdersNotFound)
	}
	return nil
}
