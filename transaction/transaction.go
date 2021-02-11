package transaction

import(
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	// "strings"
	// "bytes"
	// "strconv"
	// "github.com/gomodule/redigo/redis"
	"workerservice/model"
)

const apiurl = model.ApiTransactionUrl
const notifUrl = model.ApiNotificationUrl

type Response = model.PromotionsResponse
type Data = model.Promotions

func GenerateInOut(){
	log.Printf("start GenerateInOut")
	response := postData()
	log.Println(response)
	log.Printf("end GenerateInOut")
}

func postData() Response{
	url := apiurl + "/transaction/summaryInOutTransaction"
	timeout := time.Duration(120 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("POST", url, nil)
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
	// log.Println(string(body))
	if err != nil {
		log.Println(err)
	}
	var res Response
	json.Unmarshal(body, &res)
	return res
}