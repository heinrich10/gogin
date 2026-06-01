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
		expectedPage   int
	}{
		{
			name:           "Default values",
			queryParams:    map[string]string{},
			expectedLimit:  10,
			expectedOffset: 0,
			expectedPage:   1,
		},
		{
			name: "Custom limit and page",
			queryParams: map[string]string{
				"limit": "20",
				"page":  "2",
			},
			expectedLimit:  20,
			expectedOffset: 20,
			expectedPage:   2,
		},
		{
			name: "Invalid limit and page",
			queryParams: map[string]string{
				"limit": "abc",
				"page":  "xyz",
			},
			expectedLimit:  10,
			expectedOffset: 0,
			expectedPage:   1,
		},
		{
			name: "Limit exceeds max",
			queryParams: map[string]string{
				"limit": "200",
			},
			expectedLimit:  100,
			expectedOffset: 0,
			expectedPage:   1,
		},
		{
			name: "Negative values",
			queryParams: map[string]string{
				"limit": "-10",
				"page":  "-1",
			},
			expectedLimit:  10,
			expectedOffset: 0,
			expectedPage:   1,
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

			limit, offset, page := Paginate(c)
			assert.Equal(t, tt.expectedLimit, limit)
			assert.Equal(t, tt.expectedOffset, offset)
			assert.Equal(t, tt.expectedPage, page)
		})
	}
}
