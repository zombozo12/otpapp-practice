package redis

import (
	"database/sql"
)

type (
	RedisPhone struct {
		Number    string       `json:"number"`
		Code      string       `json:"code"`
		ExpiredAt sql.NullTime `json:"expired_at"`
	}
)
