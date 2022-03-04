package activity

type Activity struct {
	ID int64 `json:"id"`
	UserID string `json:"user_id"`
	AnsweredAt string `json:"answered_at"`
	FirstSeenAt string `json:"first_seen_at"`
}

type Activities struct {
	Activities []Activity `json:"activities"`
}

type UserSessions struct {
	EndedAt string	`json:"ended_at"`
	StartedAt string `json:"started_at"`
	ActivityID []int	`json:"activity_id"`
	Duration float64 `json:"duration_seconds"`
}