package balance

import(
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"github.com/gomodule/redigo/redis"
	"workerservice/model"
)

const apiurl = model.ApiUrl

type Response = model.Response

func getCustomer() Response{
	log.Printf("start getCustomer")
	url := apiurl + "/wallet/cronjob/getCustomer"
	log.Printf(url)
	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("clientKey", "clientKey")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	// log.Println(res.ResponseMessage)
	log.Printf("end getCustomer")
	return res
} 

func getSaldo(id string) Response{
	log.Printf("start getSaldo")
	url := apiurl + "/wallet/cronjob/getSaldo/" + id
	timeout := time.Duration(60 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("clientKey", "clientKey")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	// log.Println(res.ResponseMessage)
	log.Printf("end getSaldo")
	return res
}

func GenerateSaldo() {
	log.Printf("start generateSaldo")
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	response := getCustomer()
	log.Println(response)
	if response.ResponseCode == "200" {
		for _, data := range response.Data {
			res := getSaldo(data.Id)
			log.Println(res);
			if res.ResponseCode == "200" {
				_, err = conn.Do("HMSET", "customerSaldo:" + data.Id, "id", data.Id, "nominal", res.Data[0].Nominal, "periode", res.Data[0].Periode)
				if err != nil {
					log.Println(err)
				}
				// log.Println("HMSET", "customerSaldo:" + data.Id, "id", data.Id, "nominal", res.Data[0].Nominal, "periode", res.Data[0].Periode)
			} else {
				log.Println("error getSaldo: ", res)
			}
		}
	} else {
		log.Println("error getCustomer: ", response)
	}
	log.Printf("end generateSaldo")
}

func GenerateCustomerSaldo() {
	log.Println("start GenerateCustomerSaldo")
	res := saveCustomerSaldo()
	log.Println(res)
	log.Println("end GenerateCustomerSaldo")
}

func saveCustomerSaldo() Response{
	log.Printf("start saveCustomerSaldo")
	url := apiurl + "/wallet/cronjob/saveCustomerSaldo"
	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("serviceCode", "eyJjb2RlIjoid29ya2VyU2VydmljZVVsdGlwYXkiLCJjcmVhdGVkX2F0IjoiMTYxMjc1NjY4MTkwOSJ9")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	log.Println(res.ResponseMessage)
	log.Printf("end saveCustomerSaldo")
	return res
}

func saveSaldoToDB(id, nominal, periode string) Response{
	log.Printf("start saveSaldoToDB")
	url := apiurl + "/wallet/cronjob/saveSaldo"
	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("clientKey", "clientKey")
	req.Header.Add("id", id)
	req.Header.Add("nominal", nominal)
	req.Header.Add("periode", periode)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	log.Println(res.ResponseMessage)
	log.Printf("end saveSaldoToDB")
	return res
}

func SaveSaldo() {
	log.Printf("start saveSaldo")
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	response := getCustomer()
	if response.ResponseCode == "200" {
		for _, data := range response.Data {
			id, err := redis.String(conn.Do("HGET", "customerSaldo:" + data.Id, "id"))
			if err != nil {
				log.Println(err)
				continue
			}
			nominal, err := redis.String(conn.Do("HGET", "customerSaldo:" + data.Id, "nominal"))
			if err != nil {
				log.Println(err)
				continue 
			}
			periode, err := redis.String(conn.Do("HGET", "customerSaldo:" + data.Id, "periode"))
			if err != nil {
				log.Println(err)
				continue 
			}
			// log.Println(id, ": " , nominal, ": ", periode)
			saveSaldoToDB(id, nominal, periode)
			// log.Println(res)
			_, err = conn.Do("DEL", "customerSaldo:" + data.Id)
			if err != nil {
				log.Println(err)
			}
		}
	} else {
		log.Println("error getCustomer: ", response)
	}
	log.Printf("end saveSaldo")
}

func CheckLastDay() bool{
	log.Printf("start CheckLastDay")
	currentTime := time.Now().Format("2006-01-02")
	lastDay := getFirstDay().Format("2006-01-02")
	log.Println("current: ", currentTime)
	log.Println("lastDay: ", lastDay)
	if currentTime == lastDay {
		log.Println("CheckLastDay: true")
		return true
	} else {
		log.Println("CheckLastDay: false")
		return false
	}
}

func getFirstDay() time.Time{
	log.Printf("start getFirstDay")
	now := time.Now()
    currentYear, currentMonth, _ := now.Date()
    currentLocation := now.Location()

    firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
    // lastOfMonth := firstOfMonth.AddDate(0, 1, -1)

    // fmt.Println(firstOfMonth)
	// fmt.Println(lastOfMonth)
	log.Printf("end getFirstDay")
	return firstOfMonth
	// return lastOfMonth
}