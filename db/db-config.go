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

const SelectQueryInfoWithoutItems string = "SELECT " +
	"o.order_uid, " +
	"o.track_number, " +
	"o.entry, " +
	"d.name AS delivery_name, " +
	"d.phone AS delivery_phone, " +
	"d.zip AS delivery_zip, " +
	"d.city AS delivery_city, " +
	"d.address AS delivery_address, " +
	"d.region AS delivery_region, " +
	"d.email AS delivery_email, " +
	"p.transaction AS payment_transaction, " +
	"p.request_id AS payment_request_id, " +
	"p.currency AS payment_currency, " +
	"p.provider AS payment_provider, " +
	"p.amount AS payment_amount, " +
	"p.payment_dt AS payment_dt, " +
	"p.bank AS payment_bank, " +
	"p.delivery_cost AS payment_delivery_cost, " +
	"p.goods_total AS payment_goods_total, " +
	"p.custom_fee AS payment_custom_fee, " +
	"o.locale, " +
	"o.internal_signature, " +
	"o.customer_id, " +
	"o.delivery_service, " +
	"o.shardkey, " +
	"o.sm_id, " +
	"o.date_created, " +
	"o.oof_shard " +
	"FROM order_info o " +
	"JOIN deliveries d ON o.order_uid = d.order_uid " +
	"JOIN payments p ON o.order_uid = p.transaction " +
	"WHERE o.order_uid = $1;"

const SelectQueryItems string = "SELECT " +
	"chrt_id, " +
	"track_number, " +
	"price, " +
	"rid, " +
	"name, " +
	"sale, " +
	"size, " +
	"total_price, " +
	"nm_id, " +
	"brand, " +
	"status " +
	"FROM order_items " +
	"WHERE order_uid = $1;"
