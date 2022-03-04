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

	var bearer = "Bearer " + AUTHORIZATION_HEADER_TOKEN

	req, err := http.NewRequest("GET", USER_ACTIVITIES, nil)

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error occurred while retrieving activities from the server.\n[ERROR] -", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error occurred while closing")
		}
	}(response.Body)
	body, err := ioutil.ReadAll(response.Body)

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

