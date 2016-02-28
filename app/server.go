// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/lukeatherton/domain-events"
)

func InitApiServices(publisher Publisher, repo Repo, auth Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("publisher", publisher)
		c.Set("repo", repo)
		c.Set("auth", auth)
		c.Next()
	}
}

func NewRouter(publisher Publisher, repo Repo, auth Authenticator) (router *gin.Engine) {
	r := gin.Default()

	gin.SetMode(gin.TestMode)

	r.Use(InitApiServices(publisher, repo, auth))

	r.GET("/status", func(c *gin.Context) {
		c.String(200, "OK")
	})

	api := r.Group("/api")
	{
		api.OPTIONS("/auth", SendOptions("POST", false))
		api.POST("/auth", AllowOrigin("*"), Authenticate)

		api.OPTIONS("/emails", SendOptions("GET", false))
		api.GET("/emails", AllowOrigin("*"), CheckEmail)

		api.OPTIONS("/registrations", SendOptions("POST", false))
		api.POST("/registrations", AllowOrigin("*"), RegisterUser)

		api.OPTIONS("/verification", SendOptions("GET", true))
		api.GET("/verification", AllowOrigin("*"), VerifyEmail)

		api.OPTIONS("/credentials/updaterequests", SendOptions("POST", true))
		api.POST("/credentials/updaterequests", AllowOrigin("*"), Authorization(auth), ChangePassword)
	}

	return r
}

func SendOptions(methods string, isAuthRequired bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Add("Access-Control-Allow-Methods", methods)

		allowHeaders := "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Accept"

		if isAuthRequired {
			allowHeaders += ", Authorization"
		}

		c.Writer.Header().Add("Access-Control-Allow-Headers", allowHeaders)
		c.Writer.Header().Add("Access-Control-Allow-Credentials", "true")
		c.Next()
	}
}

func AllowOrigin(origins string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// allow cross domain AJAX requests
		c.Writer.Header().Set("Access-Control-Allow-Origin", origins)

		c.Next()
	}
}

func Authorization(auth Authenticator) gin.HandlerFunc {
	return func(c *gin.Context) {

		authorizationHeader := c.Request.Header["Authorization"]

		if len(authorizationHeader) < 1 {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"user\"")
			c.JSON(http.StatusUnauthorized, http.StatusText(401))
			return
		}

		authorizationArray := strings.SplitN(authorizationHeader[0], " ", 2)

		if len(authorizationArray) != 2 || authorizationArray[0] != "Bearer" {
			c.Writer.Header().Set("WWW-Authenticate", "Bearer realm=\"user\"")
			c.JSON(http.StatusUnauthorized, http.StatusText(401))
			return
		}

		if !auth.ValidateToken(authorizationArray[1]) {
			c.JSON(http.StatusUnauthorized, "authorization failed")
			return
		}

		// allow cross domain AJAX requests
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

		id, _ := auth.GetTokenClaim(authorizationArray[1], "id")
		c.Set("userId", id)

		email, _ := auth.GetTokenClaim(authorizationArray[1], "email")
		c.Set("email", email)

		c.Next()
	}
}
