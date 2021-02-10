package model

type Customer struct {
	Id string `json:"id"`
	Nominal string `json:"nominal"`
	Periode string `json:"periode"`
}

type Response struct {
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Data []Customer
}

type BillfazzProduct struct {
	Code string `json:"code"`
	Name string `json:"name"`
	Description string `json:"description"`
	Type string `json:"type"`
	OperatorCode string `json:"operatorCode"`
	AdminPrice int `json:"adminProce"`
	SellPrice int `json:"sellPrice"`
	Active bool `json:"active"`
	Problem bool `json:"problem"`
}

type BillfazzResponse struct {
	Data []BillfazzProduct
}

type Promotions struct {
	Id string `json:"id"`
	Startdate string `json:"startdate"`
	Enddate string `json:"enddate"`
	Promoname string `json:"promoname"`
	Attach_code string `json:"attach_code"`
	Start_date string `json:"start_date"`
	End_date string `json:"end_date"`
	Amount string `json:"amount"`
	Content string `json:"content"`
	Blast_date string `json:"blast_date"`
	Blast_times string `json:"blast_times"`
	Blast_category string `json:"blast_category"`
}

type PromotionsResponse struct {
	ResponseCode string `json:"responseCode"`
	ResponseMessage string `json:"responseMessage"`
	Data []Promotions
}

const (
	// BILLFAZZ_URL_SANDBOX = "https://secure.billfazz.com/sandbox/api/v1"
	// BILLFAZZ_KEY_SANDBOX = "dXNlcm5hbWU6NGI1NTIyOGRmZGYyOTgyNDBhZDA5MjU4NWMyM2NiZWQ3ZGRjMTgyMjFjMTUxZDUzNDU="
	// ApiUrl = "https://sandbox.api.ultipay.id:8443"
	// ApiTransactionUrl = "https://sandbox.api.ultipay.id:8443"
	// ApiNotificationUrl = "https://sandbox.api.ultipay.id:8443"
	BILLFAZZ_URL_SANDBOX = "https://secure.billfazz.com/sandbox/api/v1"
	BILLFAZZ_KEY_SANDBOX = "dXNlcm5hbWU6NGI1NTIyOGRmZGYyOTgyNDBhZDA5MjU4NWMyM2NiZWQ3ZGRjMTgyMjFjMTUxZDUzNDU="
	ApiUrl = "http://localhost:8443"
	ApiTransactionUrl = "http://localhost:8443"
	ApiNotificationUrl = "http://localhost:8443"
)