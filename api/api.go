package api

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"slang/activity"
)

const (
	USER_ACTIVITIES = "https://api.slangapp.com/challenges/v1/activities/"
	AUTHORIZATION_HEADER_TOKEN = "MzU6bXpycGdMcVZWL2hhbVdybW8rRjlCVlRTSUdvQ1IycGVKSDJlTnNHcE5ZWT0="

)
func FindActivitiesFromApi() (activities *[]activity.Activity, err error) {
	// Create a Bearer string by appending string access token
	var bearer = "Bearer " + AUTHORIZATION_HEADER_TOKEN
	// Create a new request using http
	req, err := http.NewRequest("GET", USER_ACTIVITIES, nil)

	// add authorization header to the req
	req.Header.Add("Authorization", bearer)
	//req.Header.Add("Accept", "application/json")
	// Send req using http Client
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error occurred while retrieving activities from the server.\n[ERROR] -", err)
	}

	//response, err := http.Get(USER_ACTIVITIES)
	//if err != nil {
	//	log.Println("error occurred while retrieving activities from the server")
	//	return nil, err
	//}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error occurred while closing")
		}
	}(response.Body)
	body, err := ioutil.ReadAll(response.Body)
	//fmt.Println(body)
	if err != nil {
		log.Println("failed to read response body from Api")
		return nil, err
	}
	act := &activity.Activities{}
	err = json.Unmarshal(body, act)
	if err != nil {
		log.Println("failed to unmarshal response body from api")
		return nil, err
	}
	fmt.Println(&act.Activities)

	return &act.Activities, err
}

func GetJson(url string, target interface{}) error {
	response, err := http.Get(url)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	return json.NewDecoder(response.Body).Decode(target)
}

