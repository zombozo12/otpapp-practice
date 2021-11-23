package redis

import "time"

type (
	RedisPhone struct {
		Number    string    `json:"number"`
		Code      string    `json:"code"`
		ExpiredAt time.Time `json:"expired_at"`
	}
)
