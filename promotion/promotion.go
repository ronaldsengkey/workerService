package promotion

import(
	"log"
	"time"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	// "bytes"
	// "strconv"
	"github.com/gomodule/redigo/redis"
	"workerservice/model"
)

const apiurl = model.ApiTransactionUrl
const notifUrl = model.ApiNotificationUrl

type Response = model.PromotionsResponse
type Data = model.Promotions

func GeneratePromotionList(){
	log.Printf("start GeneratePromotionList")
	response := getPromotionData()
	log.Println(response)
	if response.ResponseCode == "200" {
		currentTime := time.Now()
		now := currentTime.Format("2006-01-02")
		for _, data := range response.Data {
			// log.Println(data.Id + "---")
			if data.Blast_date != "" {
				dates := strings.Split(data.Blast_date, ",")
				for _, date := range dates {
					// log.Println(now + ": "  + date)
					if now == date {
						addToNotifService(data)
						// log.Println("true")
						break
					}
				}
			}
			if data.Blast_times != "" {
				// log.Println(now + ": " + data.Startdate)
				if now == data.Startdate {
					addToNotifService(data)
					// log.Println("true")
				}
			}
		}
	} else {
		log.Println("error GeneratePromotionList: ", response)
	}
}

func getPromotionData() Response{
	url := apiurl + "/transaction/cronjob/getPromotion"
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("GET", url, nil)
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

func storeDataToRedis(data Data){
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()
	// conn.Do("del", "notifService")
	// conn.Do("rpush", "notifService", "1")
	// conn.Do("rpush", "notifService", "10")
	// value, err := redis.Strings(conn.Do("LRANGE", "notifService", 0, -1))
	// if err != nil {
	//   log.Println(err)
	// }
	// log.Println("value: ", value)
	// pop, err := redis.Strings(conn.Do("LPOP", "notifService"))
	// if err != nil {
	//   log.Println(err)
	// }
	// log.Println("pop: ", pop)
	neo, err := json.Marshal(&data)
    if err != nil {
        log.Println(err)
    }
	log.Println(string(neo))
	conn.Do("rpush", "notifService", string(neo))
}

func addToNotifService(data Data){
	param, err := json.Marshal(&data)
    if err != nil {
        log.Println(err)
    }
	log.Println(string(param))
	url := notifUrl + "/notification/cronjob/saveNotifPromo"
	timeout := time.Duration(5 * time.Second)
	client := http.Client{
		Timeout: timeout,
	}
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("workerService", "workerService")
	req.Header.Add("param", string(param))
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
	log.Println("res: ", res.ResponseCode)
	// return res
}