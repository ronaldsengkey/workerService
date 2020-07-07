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

const (
	BILLFAZZ_URL_SANDBOX = "https://secure.billfazz.com/sandbox/api/v1"
	BILLFAZZ_KEY_SANDBOX = "dXNlcm5hbWU6NGI1NTIyOGRmZGYyOTgyNDBhZDA5MjU4NWMyM2NiZWQ3ZGRjMTgyMjFjMTUxZDUzNDU="
	ApiUrl = "http://localhost:8089"
)