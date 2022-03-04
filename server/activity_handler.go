package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"slang/activity"
	"slang/api"
	"time"
)

func GetActivities() gin.HandlerFunc {
	return func(c *gin.Context) {
		//var activities []activity.Activity
		fmt.Println("I got here")
		data, err := api.FindActivitiesFromApi()
		fmt.Println(data)
		if err != nil {
			log.Println(err)
		}
		sessions := activity.Sessions{
			GotSessions: map[string][]activity.UserSessions{},
		}
		for _, act := range *data {
			actualSession := NewUserSession(act)
			userSessions, ok := sessions.GotSessions[act.UserID]
			if !ok {
				sessions.GotSessions[act.UserID] = append(sessions.GotSessions[act.UserID], actualSession)
			} else {

				updated := false
				var inputIndex int
				for index, userSession := range userSessions {
					if isAfterTime(userSession.StartedAt, actualSession.StartedAt) {
						if isIn5MinutesSession(userSession.StartedAt, actualSession.EndedAt) {
							userSessions[index] = changeLastInSession(userSession, actualSession)
							updated = true
						}
					} else {
						if isIn5MinutesSession(userSession.StartedAt, actualSession.EndedAt) {
							userSessions[index] = changeFirstInSession(userSession, actualSession)
							updated = true
						} else {
							inputIndex = index
						}
					}
				}
				if !updated {
					if inputIndex == 0 {
						slice := []activity.UserSessions{}
						slice = append(slice, actualSession)
						sessions.GotSessions[act.UserID] = append(slice, sessions.GotSessions[act.UserID]...)
					}else if inputIndex != 0 {
						lastElements := sessions.GotSessions[act.UserID][inputIndex:]
						sessions.GotSessions[act.UserID] = append(sessions.GotSessions[act.UserID][:inputIndex], actualSession)
						sessions.GotSessions[act.UserID] = append(sessions.GotSessions[act.UserID], lastElements...)

					} else {
						sessions.GotSessions[act.UserID] = append(sessions.GotSessions[act.UserID], actualSession)
					}
				}
			}
		}
		//var movies = s.Cache.Get("movies")
		//if data != nil {
		//	//data, err := api.FindActivitiesFromApi()
		//	//if err != nil {
		//	//	c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		//	//	return
		//	//}
		//	result := *data
		//	sort.Slice(result, func(i, j int) bool {
		//		return result[i].AnsweredAt > result[j].FirstSeenAt
		//	})
		//
		//	sort.Slice(result, func(i, j int) bool {
		//		//str1 := result[i].AnsweredAt
		//		t1, _ := time.Parse(time.RFC3339, result[i].AnsweredAt)
		//		t2, _ := time.Parse(time.RFC3339, result[i].FirstSeenAt)
		//		t3, _ := time.Parse(time.RFC3339, result[j].AnsweredAt)
		//		t4, _ := time.Parse(time.RFC3339, result[j].FirstSeenAt)
		//
		//
		//		return t1.Unix() - t2.Unix() > t3.Unix()-t4.Unix()
		//	})
		//
		//	activities = result
		//	log.Println("activity log")
		//}
		c.JSON(http.StatusOK, gin.H{"user_activities": sessions.GotSessions})
	}
}


func TimeDifference(start string, end string) float64 {
	startT, _ := time.Parse(time.RFC3339, start)
	endT , _ := time.Parse(time.RFC3339, end)
	abc := math.Abs(float64(startT.Unix() - endT.Unix()))
	return abc
}

func NewUserSession(user activity.Activity) activity.UserSessions {
	act := activity.UserSessions{
		EndedAt:   user.AnsweredAt,
		StartedAt: user.FirstSeenAt,
		Duration:  TimeDifference(user.FirstSeenAt, user.AnsweredAt),
	}
	act.ActivityID = append(act.ActivityID, user.ID)
	return act
}

func isAfterTime(timestamp1 string, timestamp2 string) bool {
	timeDifferenceInSeconds := TimeDifference(timestamp1, timestamp2)
	return timeDifferenceInSeconds > 0
}

func isIn5MinutesSession(timestamp1 string, timestamp2 string) bool {
	timeDifferenceInSeconds := TimeDifference(timestamp1, timestamp2)
	return timeDifferenceInSeconds <= 300
}

func changeUserActivity(sessionActivity activity.UserSessions, actualSessionActivity activity.UserSessions) activity.UserSessions {
	updateSession := sessionActivity
	updateSession.ActivityID = append(updateSession.ActivityID, actualSessionActivity.ActivityID...)
	updateSession.Duration = TimeDifference(sessionActivity.StartedAt, sessionActivity.EndedAt)
	return updateSession
}

func changeFirstInSession(sessionActivity activity.UserSessions, actualSessionActivity activity.UserSessions) activity.UserSessions {
	updateSession := changeUserActivity(sessionActivity, actualSessionActivity)
	updateSession.StartedAt = actualSessionActivity.StartedAt
	return updateSession
}

func changeLastInSession(sessionActivity activity.UserSessions, actualSessionActivity activity.UserSessions) activity.UserSessions {
	updateSession := changeUserActivity(sessionActivity, actualSessionActivity)
	updateSession.EndedAt = actualSessionActivity.EndedAt
	return updateSession
}
