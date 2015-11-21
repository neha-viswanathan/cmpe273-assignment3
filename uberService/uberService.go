/*
* CMPE 273 Assignment 3 - Trip planner using Uber API
* Neha Viswanathan
* 010029097
*/

package uberService

//import statements
import (
	"fmt"
  "net/http"
  "encoding/json"
  "io/ioutil"
  "bytes"
)

//struct declarations - start
//List of price estimates struct
type PriceList struct {
	Price []PriceEstimate `json:"priceEst"`
}

// Price estimate struct for Uber
type PriceEstimate struct {
	ProductId string `json:"product_id"`
	CurrencyCode string `json:"currency_code"`
	DisplayName string `json:"display_name"`
	Estimate string `json:"estimate"`
	LowEstimate int `json:"low_estimate"`
	HighEstimate int `json:"high_estimate"`
	SurgeMultiplier float64 `json:"surge_multiplier"`
	Duration int `json:"duration"`
	Distance float64 `json:"distance"`
}

//Output struct for Uber
type UberOut struct {
	Distance float64
  Cost int
	Duration int
}

//ETA struct for Uber
type UberETA struct {
  Request_id string `json:"request_id"`
  Status string	`json:"status"`
  Vehicle string `json:"vehicle"`
  Driver string	`json:"driver"`
  Location string	`json:"location"`
  ETA int	`json:"eta"`
  SurgeMultiplier float64 `json:"surge_multiplier"`
}
//struct declarations - end

//This function is to determine price rate for Uber
func GetUberRate(startLati, startLongi, endLati, endLongi string) UberOut {
	client := &http.Client{}
	reqURL := fmt.Sprintf("https://sandbox-api.uber.com/v1/estimates/price?start_latitude=%s&start_longitude=%s&end_latitude=%s&end_longitude=%s&server_token=81WOPreAmQKFjpTEoVN6tdQXD98XZGlttN5fB1Ia", startLati, startLongi, endLati, endLongi)
	fmt.Println("URL created is :: " + reqURL)

	req, err := http.NewRequest("GET", reqURL , nil)
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while sending request to Uber :: ", err);
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while trying to read response :: ", err);
	}

	var pl PriceList
	err = json.Unmarshal(body, &pl)
	if err != nil {
		fmt.Println("Error while trying to unmarshall response :: ", err);
	}

	var uOut UberOut
	uOut.Distance = pl.Price[0].Distance
  uOut.Cost = pl.Price[0].LowEstimate
	uOut.Duration = pl.Price[0].Duration

  return uOut
}

//This function is to determine estimated arrival time (ETA) for Uber
func GetUberETA(startLati, startLongi, endLati, endLongi string) int {
	var jsonStr = []byte(`{"start_latitude":"` + startLati + `","start_longitude":"` + startLongi + `","end_latitude":"` + endLati + `","end_longitude":"` + endLongi + `","product_id":"a1111c8c-c720-46c3-8534-2fcdd730040d"}`)
	reqURL := "https://sandbox-api.uber.com/v1/requests"

	req, err := http.NewRequest("POST", reqURL, bytes.NewBuffer(jsonStr))
	req.Header.Set("Authorization", "Bearer eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJzY29wZXMiOlsicmVxdWVzdCJdLCJzdWIiOiJjZjg5YjZlMi04NTBkLTRhZjktYmU3MC05OTJmMDcyOTU0Y2YiLCJpc3MiOiJ1YmVyLXVzMSIsImp0aSI6ImQ4OGRiYTJhLTE0OWUtNDc5MC04ODJkLTMxYjIzZGM2NWZlYiIsImV4cCI6MTQ1MDY4Mjc3NiwiaWF0IjoxNDQ4MDkwNzc1LCJ1YWN0IjoiVzA3NGZTRDZNZmMzYVRLN1dORVNIeGoyRmJiOU9SIiwibmJmIjoxNDQ4MDkwNjg1LCJhdWQiOiJUd0FuODhkUGxkS0NtcEpnU29Da0RUQllqUDBsbHNkQyJ9.OksqiPE908841Cy0hvTiXFvkBU3SoeSVaY79AO5UjkbKCMXkMvMkevX6Z_EAtPicc0FM6ZgjrInhIHQbxpg9KsgTAa0i13MsDgm6jsU6pm9_6Gj6e5b59WIy1EcRf4A9Z6IfLiudCd6XKuqm-HwxFST2ttNmg9UNeU2V6wC-xkdAZosriq_qvHxviCFwF2HAkuYiN6m3MNR3hPGl9ikTZ1dPClvZPFSzi80QYAkVdGyiVQYjPmYrvDnQdmTyD4qE4xWSFmlvjkxVvDxRR0_6_4Y8YjnsTFus1_2YWOyhmGT7QqDNCXIEAsniDfxXYTRyosSb8wjgE9nvrGgHPQQb4g")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while sending request to Uber :: ", err);
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error while reading response :: ", err);
	}

	var ue UberETA
	err = json.Unmarshal(body, &ue)
	if err != nil {
		fmt.Println("Error while unmarshalling response :: ", err);
	}

	eta := ue.ETA
	return eta
}
