// utils/query_helpers.go
package helper

import (
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/pratikgagare03/feedback/logger"
	"github.com/pratikgagare03/feedback/models"
)

// ParseQueryInt parses a query parameter as an integer and applies a default value if the parameter is missing or invalid.
func ParseQueryInt(c *gin.Context, key string, defaultValue int, limit int) (int, error) {
	valueStr := c.Query(key)
	if valueStr == "" {
		return defaultValue, nil
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return 0, fmt.Errorf("invalid %s: %v", key, err)
	}
	if value <= 0 || value > limit {
		logger.Logs.Error().Msgf("out of range %s: %d", key, value)
		logger.Logs.Info().Msgf("setting %s to default value %d", key, defaultValue)
		return defaultValue, nil
	}
	return value, nil
}

func GetParamsModel(c *gin.Context) (models.QueryParams, error) {
	params := models.QueryParams{
		DateFrom:   c.Query("dateFrom"),
		DateTo:     c.Query("dateTo"),
		Query:      c.Query("query"),
		InQuestion: c.Query("inQuestion"),
		InAnswer:   c.Query("inAnswer"),
	}
	var parseErr error
	if params.DateFrom != "" {
		params.DateFrom, params.DateTo, parseErr = GetParsedDateRange(params.DateFrom, params.DateTo)
		if parseErr != nil {
			return params, parseErr
		}
	}
	return params, nil
}
