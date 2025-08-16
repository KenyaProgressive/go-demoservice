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

	db.SetMaxOpenConns(50)

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
		utils.DbLogger.Errorf("Transaction-Push wasn't started wtih error: %s", err)
		return err
	}

	defer func() {
		if err != nil { // Transaction-Push Error
			errRb := tx.Rollback()

			if errRb != nil { // Rollback Error
				utils.DbLogger.Errorf("Rollback wasn't completed with error: %s", errRb)
			} else { // Success Rollback
				utils.DbLogger.Warnf("Transaction-Push was failed. Doing rollback")
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

func GetInfoAndCreateMessage(db *sql.DB, target_uuid string, msg *utils.Message) error {

	withoutItemsTx, errWithoutItemsTx := db.Begin()
	if errWithoutItemsTx != nil {
		utils.DbLogger.Errorf("Info getting failed with error: %s", errWithoutItemsTx)
		return errWithoutItemsTx
	}

	defer func() {
		if errWithoutItemsTx != nil { // Transaction-Select Error
			_ = withoutItemsTx.Rollback()
		}
	}()

	resultWithoutItems := withoutItemsTx.QueryRow(SelectQueryInfoWithoutItems, target_uuid)

	txItems, errTxItems := db.Begin()
	if errTxItems != nil {
		utils.DbLogger.Errorf("Info getting failed with error: %s", errTxItems)
		return errTxItems
	}

	defer func() {
		if errTxItems != nil { // Transaction-Select Error
			_ = txItems.Rollback()
		}
	}()

	resultItems, errItems := txItems.Query(SelectQueryItems, target_uuid)
	if errItems != nil {
		utils.DbLogger.Errorf("Query was failed with error: %s", errItems)
		return errItems
	}

	//utils.DbLogger.Debug(resultItems.Columns())

	errCompletingMessage := PushingToMessageStructure(resultWithoutItems, resultItems, msg)
	if errCompletingMessage != nil {
		utils.DbLogger.Errorf("Error in pushing data to structure: %s", errCompletingMessage)
		return errCompletingMessage
	}

	withoutItemsTx.Commit()
	txItems.Commit()

	return nil
}

func PushingToMessageStructure(withoutItemsRow *sql.Row, itemsRows *sql.Rows, msg *utils.Message) error {
	errScanWithoutItems := withoutItemsRow.Scan(
		&msg.OrderUId,
		&msg.TrackNumber,
		&msg.Entry,
		&msg.Delivery.Name,
		&msg.Delivery.Phone,
		&msg.Delivery.ZipCode,
		&msg.Delivery.City,
		&msg.Delivery.Address,
		&msg.Delivery.Region,
		&msg.Delivery.Email,
		&msg.Payment.Transaction,
		&msg.Payment.RequestId,
		&msg.Payment.Currency,
		&msg.Payment.Provider,
		&msg.Payment.Amount,
		&msg.Payment.PaymentDt,
		&msg.Payment.Bank,
		&msg.Payment.DeliveryCost,
		&msg.Payment.GoodsTotal,
		&msg.Payment.CustomFee,
		&msg.Locale,
		&msg.InternalSignature,
		&msg.CustomerId,
		&msg.DeliveryService,
		&msg.ShardKey,
		&msg.SmId,
		&msg.DateCreated,
		&msg.OofShard)

	if errScanWithoutItems != nil {
		return errScanWithoutItems
	}

	var itemInfo utils.Items

	for itemsRows.Next() {
		errScanItems := itemsRows.Scan(
			&itemInfo.ChrtID,
			&itemInfo.TrackNumber,
			&itemInfo.Price,
			&itemInfo.Rid,
			&itemInfo.Name,
			&itemInfo.Sale,
			&itemInfo.Size,
			&itemInfo.TotalPrice,
			&itemInfo.NmID,
			&itemInfo.Brand,
			&itemInfo.Status,
		)
		if errScanItems != nil {
			return errScanItems
		}
		msg.Items = append(msg.Items, itemInfo)
	}

	return nil
}
