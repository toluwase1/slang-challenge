package server

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"math"
	"net/http"
	"os"
	"slang/activity"
	"slang/api"
	"strings"
	"time"
)

//const (
//UNAUTHORIZED = 401 Unauthorized
//400 Bad request
//429 Too many requests (en el caso de ser rate-limited)
//)

func  PostActivitiesToEndPoint(body []byte) (resp uint, err error) {
	var bearer = os.Getenv("AUTHORIZATION_HEADER_TOKEN")

	req, err := http.NewRequest("POST", os.Getenv("SEND_BODY"), strings.NewReader(string(body)))

	req.Header.Add("Authorization", bearer)

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		log.Println("error occurred while retrieving activities from the server.\n[ERROR] -", err)
	}


	if response.StatusCode == 204 {
		return http.StatusNoContent, nil
	} else if response.StatusCode == 401 {
		return  http.StatusUnauthorized, nil
	} else {
		return http.StatusInternalServerError, nil
	}
}

func GetActivities() gin.HandlerFunc {
	return func(c *gin.Context) {
		data, err := api.FindActivitiesFromApi()
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
		m, err := json.Marshal(&sessions)
		PostActivitiesToEndPoint(m)
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
