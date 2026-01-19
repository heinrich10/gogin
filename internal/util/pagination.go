package util

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

func Paginate(c *gin.Context) (limit, offset int) {
	const (
		defaultLimit = 10
		defaultPage  = 1
		maxLimit     = 100
	)

	var page int

	if s := strings.TrimSpace(c.Query("limit")); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			if v <= 0 {
				v = defaultLimit
			}
			if v > maxLimit {
				v = maxLimit
			}
			limit = int(v)
		}
	}

	if s := strings.TrimSpace(c.Query("page")); s != "" {
		if v, err := strconv.ParseInt(s, 10, 64); err == nil {
			if v <= 0 {
				v = defaultPage
			}
			page = int(v)
		}
	}

	offset = page * limit
	if offset <= 0 {
		offset = 0
	}
	return limit, offset
}
