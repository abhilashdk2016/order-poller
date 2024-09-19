package storage

import (
	"context"
	"fmt"

	"github.com/abhilashdk2016/order-poller/models"
	"github.com/jackc/pgx/v4"
)

type Repository struct {
	DB *pgx.Conn
}

func (r *Repository) GetOutbox() error {
	outbox_items := []models.Outbox{}
	tx, _ := r.DB.BeginTx(context.Background(), pgx.TxOptions{})
	defer tx.Rollback(context.Background())

	rows, err := tx.Query(context.Background(), "SELECT id, payload, is_processed FROM outbox WHERE is_processed = false")
	if err != nil {
		fmt.Println("tx.Query select unprocessed outbox items error: ", err)
		return nil
	}
	defer rows.Close()

	if rows.Err() != nil {
		fmt.Println("row.Err()", err)
		return nil
	}

	for rows.Next() {
		outbox := models.Outbox{}
		err := rows.Scan(&outbox.ID, &outbox.Payload, &outbox.IsProcessed)
		if err != nil {
			return fmt.Errorf("unable to scan outbox row: %w", err)
		}
		outbox_items = append(outbox_items, outbox)
	}

	fmt.Println(outbox_items)

	_, execErr := tx.Exec(context.Background(), "Update outbox SET is_processed = true WHERE is_processed = false")
	if execErr != nil {
		fmt.Println(execErr.Error())
		return fmt.Errorf("tx.Exec updating is_processed = true error: %w", err)
	}
	if err := tx.Commit(context.Background()); err != nil {
		fmt.Println(err.Error())
		return fmt.Errorf("tx.Commit %w", err)
	}
	fmt.Println("outbox table updated")
	return nil
}
