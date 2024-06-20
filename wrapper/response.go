package wrapper

import (
	"io"
	"time"
)

type Response struct {
	CacheTime time.Duration     `json:"cacheTime"`
	ServedBy  string            `json:"servedBy"`
	Additions map[string]string `json:"additions"`
	Data      *io.Reader        `json:"data"`
}
