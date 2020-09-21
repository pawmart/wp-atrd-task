package util

import "github.com/gin-gonic/gin"

//PrepareResponse serialize struct depending on `Accept` header and add to response
func PrepareResponse(responseCode int, obj interface{}, c *gin.Context) {
	if c.Request.Header.Get("Accept") == "application/xml" {
		c.XML(responseCode, obj)
	} else {
		c.JSON(responseCode, obj)
	}
}
