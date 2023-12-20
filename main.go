package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	MAX_REQUEST            = 1
	ADD_ONE_REQUEST_SECOND = 10
)

type ClientDetails struct {
	count       int
	lastUpdated time.Time
}

var clientList map[string]ClientDetails = map[string]ClientDetails{}

/*
{
	"1.2":"1",
	"2.3":"2",
	"3.4":"3"
	{
		"count":"10",
		"last":"12345678"
	}
}
*/

func main() {

	r := gin.Default()
	r.GET("/unlimited", func(c *gin.Context) {
		c.JSON(http.StatusOK, map[string]string{
			"message": "Unlimited! Let's Go!",
		})
	})

	r.GET("/limited", func(c *gin.Context) {
		if rl := RateLimit("1"); rl {
			c.JSON(http.StatusTooManyRequests, map[string]string{
				"message": "To Many Request",
			})
			return
		}

		c.JSON(http.StatusOK, map[string]string{
			"message": "Limited, don't over use me!",
		})
	})

	r.Run()

}

func RateLimit(ip string) bool {
	data, ok := clientList[ip]

	if !ok {
		clientList[ip] = ClientDetails{count: MAX_REQUEST - 1, lastUpdated: time.Now()}
		return false
	} else {
		timeDiff := time.Now().Second() - data.lastUpdated.Second()
		if timeDiff >= ADD_ONE_REQUEST_SECOND || data.count > 0 {
			data.count = data.count + timeDiff
			if data.count > MAX_REQUEST {
				data.count = MAX_REQUEST
			}
			data.count--
			return false
		} else {
			return true
		}
	}

}
