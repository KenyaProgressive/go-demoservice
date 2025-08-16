package db

import (
	"database/sql"
	"encoding/json"
	"go-demoservice/utils"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func MakeDbConnection() (*sql.DB, error) {
	db, err := sql.Open("pgx", utils.ConnectString)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		utils.DbLogger.Error(err)
		return nil, err
	}
	if err := createDbAndTables(db); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(350)

	return db, nil
}

func createDbAndTables(db *sql.DB) error {
	// if _, err := db.Exec("CREATE DATABASE test_db;"); err != nil {
	// 	return err
	// }

	startQuerysBytes, err := os.ReadFile("db/creation_query.sql")
	if err != nil {
		return err
	}

	startQuerys := string(startQuerysBytes)

	if _, err := db.Exec(startQuerys); err != nil {
		return err
	}

	return nil

}

func PrepareMessagesAndPushToDb(db *sql.DB, ValueToPush []byte) error {

	var msg utils.Message
	if err := json.Unmarshal(ValueToPush, &msg); err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		utils.DbLogger.Errorf("Transaction wasn't started wtih error: %s", err)
		return err
	}

	defer func() {
		if err != nil { // Transaction Error
			errRb := tx.Rollback()

			if errRb != nil { // Rollback Error
				utils.DbLogger.Errorf("Rollback wasn't completed with error: %s", errRb)
			} else { // Success Rollback
				utils.DbLogger.Warnf("Transaction was failed. Doing rollback")
			}
		}
	}()

	if _, err := tx.Exec(InsertionQueryCustomers, msg.CustomerId, msg.Delivery.Email, msg.Delivery.Name,
		msg.Delivery.Phone); err != nil {
		return err
	}

	if _, err := tx.Exec(InsertionQueryOrderInfo, msg.OrderUId, msg.CustomerId, msg.TrackNumber,
		msg.Entry, msg.Locale, msg.DeliveryService, msg.ShardKey,
		msg.SmId, msg.OofShard, msg.InternalSignature, msg.DateCreated); err != nil {
		return err
	}

	if _, err := tx.Exec(InsertionQueryDeliveries, msg.OrderUId, msg.Delivery.Name, msg.Delivery.Phone,
		msg.Delivery.ZipCode, msg.Delivery.City, msg.Delivery.Address, msg.Delivery.Region,
		msg.Delivery.Email); err != nil {
		return err
	}

	if _, err := tx.Exec(InsertionQueryPayments, msg.Payment.Transaction, msg.Payment.RequestId, msg.Payment.Provider,
		msg.Payment.Bank, msg.Payment.Amount, msg.Payment.Currency, msg.Payment.PaymentDt,
		msg.Payment.DeliveryCost, msg.Payment.CustomFee, msg.Payment.GoodsTotal); err != nil {
		return err
	}

	if _, err := tx.Exec(InsertionQueryOrderItems, msg.Items[0].ChrtID, msg.Items[0].TrackNumber, msg.Items[0].Price,
		msg.Items[0].Rid, msg.Items[0].Name, msg.Items[0].Sale, msg.Items[0].Size,
		msg.Items[0].TotalPrice, msg.Items[0].NmID, msg.Items[0].Brand, msg.Items[0].Status,
		msg.OrderUId); err != nil {
		return err
	}

	return tx.Commit()
}
