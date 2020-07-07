package billfazz

import(
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"bytes"
	"strconv"
	"workerservice/model"
)

const (
	BILLFAZZ_URL_SANDBOX = model.BILLFAZZ_URL_SANDBOX
	BILLFAZZ_KEY_SANDBOX = model.BILLFAZZ_KEY_SANDBOX
	apiUrl = model.ApiUrl
)

type BillfazzResponse = model.BillfazzResponse
type Response = model.Response
type BillfazzProduct = model.BillfazzProduct

func BillfazzCronjob() {
	log.Println("start billfazzCronjob")
	res := getBillfazzProduct()
	// log.Println(len(res.Data))
	if len(res.Data) > 0 {
		resDel := delBillfazzProduct()
		// log.Println(resDel)
		if resDel.ResponseCode == "200" {
			for _, data := range res.Data {
				saveBillfazzProduct(data)
			}
		} else {
			log.Println("Error Delete Data: ", resDel)
		}
	} else {
		log.Println("Error Get Data: ", res)
	}
	log.Println("end billfazzCronjob")
}

func getBillfazzProduct() BillfazzResponse{
	url := BILLFAZZ_URL_SANDBOX + "/products/client"
	timeout := time.Duration(10 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", "Basic " + BILLFAZZ_KEY_SANDBOX)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	var res BillfazzResponse
	json.Unmarshal(body, &res)
	if len(res.Data) < 1 {
		// log.Println(string(body))
	}
	return res
}

func delBillfazzProduct() Response{
	url := apiUrl + "/wallet/cronjob/billfazz"
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("workerService", "workerService")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Println(string(body))
	var res Response
	json.Unmarshal(body, &res)
	// log.Println(res.Data[0].Code)
	return res
}

func saveBillfazzProduct(data BillfazzProduct) Response{
	stringUrl := apiUrl + "/wallet/cronjob/billfazz"
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	reqBody, _ := json.Marshal(map[string]string{
		"code": data.Code,
		"name": data.Name,
		"description": data.Description,
		"type": data.Type,
		"operatorCode": data.OperatorCode,
		"adminPrice": strconv.Itoa(data.AdminPrice),
		"sellPrice": strconv.Itoa(data.SellPrice),
		"active": strconv.FormatBool(data.Active),
		"problem": strconv.FormatBool(data.Problem),
	})
	req, err := http.NewRequest("POST", stringUrl, bytes.NewBuffer(reqBody))
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("workerService", "workerService")
	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer resp.Body.Close()
	resBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
	}
	// log.Println(string(body))
	var res Response
	json.Unmarshal(resBody, &res)
	// log.Println(res.Data[0].Code)
	return res
}