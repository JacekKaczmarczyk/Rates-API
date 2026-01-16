package main

import (
	"net/http"

	"github.com/JacekKaczmarczyk/nbpApi/providers"
	"github.com/gin-gonic/gin"
)

var defaultCodes = []string{"USD", "EUR"}
var providersMap = make(map[string]providers.Provider)

func init() {
	providersMap["nbp"] = providers.NewNbpProvider()
}

func getCurrencies(c *gin.Context) {
	codes, date := parseQueryParams(c)
	providerName := c.DefaultQuery("provider", "nbp")

	provider, exists := providersMap[providerName]
	if !exists {
		c.JSON(http.StatusBadRequest, providers.ErrorResponse{
			Message: "unknown provider: " + providerName,
			Details: "supported providers: nbp",
		})
		return
	}

	response, err, statusCode := provider.GetCurrencies(codes, date)
	if err != nil {
		c.JSON(statusCode, providers.ErrorResponse{
			Message: "failed to fetch currencies",
			Details: err.Error(),
		})
		return
	}

	c.JSON(statusCode, response)
}

func parseQueryParams(c *gin.Context) ([]string, string) {
	codes, ok := c.GetQueryArray("codes")
	if !ok {
		codes = defaultCodes
	}

	date := c.DefaultQuery("date", "")
	return codes, date
}

func main() {
	router := gin.Default()
	router.GET("/currencies", getCurrencies)
	router.Run("localhost:8000")
}
