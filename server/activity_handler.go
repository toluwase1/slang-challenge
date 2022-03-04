package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"slang/activity"
	"slang/api"
	"sort"
	"time"
)

func GetActivities() gin.HandlerFunc {
	return func(c *gin.Context) {
		var activities []activity.Activity
		fmt.Println("I got here")
		data, err := api.FindActivitiesFromApi()
		fmt.Println(data)
		if err != nil {
			log.Println(err)
		}
		//var movies = s.Cache.Get("movies")
		if data != nil {
			//data, err := api.FindActivitiesFromApi()
			//if err != nil {
			//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			//	return
			//}
			result := *data
			sort.Slice(result, func(i, j int) bool {
				return result[i].AnsweredAt > result[j].FirstSeenAt
			})
			sort.Slice(result, func(i, j int) bool {
				layout := "2006-01-02T15:04:05.000Z"
				//str1 := result[i].AnsweredAt
				t1, _ := time.Parse(layout, result[i].AnsweredAt)
				t2, _ := time.Parse(layout, result[i].FirstSeenAt)
				t3, _ := time.Parse(layout, result[j].AnsweredAt)
				t4, _ := time.Parse(layout, result[j].FirstSeenAt)


				return t1.Unix() - t2.Unix() > t3.Unix()-t4.Unix()
			})

			activities = result
			log.Println("activity log")
		}
		c.JSON(http.StatusOK, activities)
	}
}
