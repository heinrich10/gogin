package util

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestPaginate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		queryParams    map[string]string
		expectedLimit  int
		expectedOffset int
	}{
		{
			name:           "Default values",
			queryParams:    map[string]string{},
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name: "Custom limit and page",
			queryParams: map[string]string{
				"limit": "20",
				"page":  "2",
			},
			expectedLimit:  20,
			expectedOffset: 20,
		},
		{
			name: "Invalid limit and page",
			queryParams: map[string]string{
				"limit": "abc",
				"page":  "xyz",
			},
			expectedLimit:  10,
			expectedOffset: 0,
		},
		{
			name: "Limit exceeds max",
			queryParams: map[string]string{
				"limit": "200",
			},
			expectedLimit:  100,
			expectedOffset: 0,
		},
		{
			name: "Negative values",
			queryParams: map[string]string{
				"limit": "-10",
				"page":  "-1",
			},
			expectedLimit:  10,
			expectedOffset: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			req := httptest.NewRequest("GET", "/", nil)
			q := req.URL.Query()
			for k, v := range tt.queryParams {
				q.Add(k, v)
			}
			req.URL.RawQuery = q.Encode()
			c.Request = req

			limit, offset := Paginate(c)
			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedOffset, offset)
		})
	}
}
