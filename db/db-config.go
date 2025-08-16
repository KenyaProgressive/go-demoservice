package db

const InsertionQueryCustomers string = "INSERT INTO customers(customer_id, email, name, phone) VALUES ($1, $2, $3, $4)"
const InsertionQueryOrderInfo string = "INSERT INTO order_info(order_uid, customer_id, track_number, entry, locale, delivery_service, shardkey, sm_id, oof_shard, internal_signature, date_created)" +
	"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)"

const InsertionQueryPayments string = "INSERT INTO payments(transaction, request_id, provider, bank, amount, currency, payment_dt, delivery_cost, custom_fee, goods_total)" +
	"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)"

const InsertionQueryDeliveries string = "INSERT INTO deliveries(order_uid, name, phone, zip, city, address, region, email)" +
	"VALUES ($1, $2, $3, $4, $5, $6, $7, $8)"

const InsertionQueryOrderItems string = "INSERT INTO order_items(chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status, order_uid)" +
	"VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)"
