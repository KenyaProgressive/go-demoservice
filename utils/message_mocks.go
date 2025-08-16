package utils

const (
	logTimeFormat string = "02.01.2006 15:04:05"
)

type Message struct {
	OrderUId          string   `json:"order_uid"`
	TrackNumber       string   `json:"track_number"`
	Entry             string   `json:"entry"`
	Delivery          Delivery `json:"delivery"`
	Payment           Payment  `json:"payment"`
	Items             []Items  `json:"items"`
	Locale            string   `json:"locale"`
	InternalSignature string   `json:"internal_signature"`
	CustomerId        string   `json:"customer_id"`
	DeliveryService   string   `json:"delivery_service"`
	ShardKey          string   `json:"shardkey"`
	SmId              uint     `json:"sm_id"`
	DateCreated       string   `json:"date_created"`
	OofShard          string   `json:"oof_shard"`
}

type Delivery struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	ZipCode string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"address"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}

type Payment struct {
	Transaction  string `json:"transaction"`
	RequestId    string `json:"request_id"`
	Currency     string `json:"currency"`
	Provider     string `json:"provider"`
	Amount       uint   `json:"amount"`
	PaymentDt    uint   `json:"payment_dt"`
	Bank         string `json:"bank"`
	DeliveryCost uint   `json:"delivery_cost"`
	GoodsTotal   uint   `json:"goods_total"`
	CustomFee    uint   `json:"custom_fee"`
}

type Items struct {
	ChrtID      uint   `json:"chrt_id"`
	TrackNumber string `json:"track_number"`
	Price       uint   `json:"price"`
	Rid         string `json:"rid"`
	Name        string `json:"name"`
	Sale        uint   `json:"sale"`
	Size        string `json:"size"`
	TotalPrice  uint   `json:"total_price"`
	NmID        uint   `json:"nm_id"`
	Brand       string `json:"brand"`
	Status      int    `json:"status"`
}

var TestMessageZipCodes [10]string = [10]string{
	"10001",
	"10115",
	"110001",
	"75001",
	"150-0001",
	"1012",
	"00175",
	"30301",
	"90210",
	"60601"}

var TestMessageBankNames [10]string = [10]string{
	"Alpha",
	"Deutsche Bank",
	"Bank of America",
	"UBS Group",
	"Ironwood Savings & Loan",
	"Blue Horizon Bank",
	"BNP Paribas",
	"Royal Bank of Canada",
	"Standard Chartered Bank",
	"ANZ Bank"}
