// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	. "github.com/lukeatherton/domain-events"
	"github.com/satori/go.uuid"
)

func GetIdParam(param string, c *gin.Context) uuid.UUID {
	id, err := uuid.FromString(param)
	if err != nil {
		errorResponse := ErrorResponse{
			Errors: []*Error{NewError(ErrCodeValueRequired, fmt.Sprintf("valid id required", param))},
		}
		c.JSON(http.StatusBadRequest, errorResponse)
		c.Abort()
		return uuid.Nil
	}

	return id
}

//Authenticates a User and returns an Auth token for use in future requests
func Authenticate(c *gin.Context) {
	auth := c.MustGet("auth").(Authenticator)
	repo := c.MustGet("repo").(Repo)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var view *LoginView
	c.Bind(&view)

	email := strings.ToLower(view.Email)
	password := view.Password

	if email != "" && password != "" {
		token, auth_err := auth.Authenticate(email, password, "")
		if auth_err != nil {
			//TODO: Fix this, reveals user details to client
			c.JSON(http.StatusUnauthorized, auth_err.Error())
			return
		}

		userId, _ := repo.FindEmail(email)

		response := AuthResponse{
			Id:    userId,
			Email: email,
			Token: token,
		}

		c.JSON(http.StatusOK, response)
		return
	}

	errorResponse := ErrorResponse{
		Errors: []*Error{NewError(ErrCodeValueRequired, "email and password are required")},
	}

	c.JSON(http.StatusUnauthorized, errorResponse)
	return
}

func CheckEmail(c *gin.Context) {
	repo := c.MustGet("repo").(Repo)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the query string arguments, if any
	qs := c.Request.URL.Query()
	email := strings.ToLower(qs.Get("email"))

	if email != "" {

		user, err := repo.FindEmail(email)

		if err != nil && err.Error() != "not found" {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		if user != uuid.Nil {
			c.JSON(http.StatusFound, email)
			return
		}

		c.JSON(http.StatusOK, email)
		return
	}

	c.JSON(http.StatusBadRequest, NewError(ErrCodeValueRequired, fmt.Sprintf("email is a required field")))
	return
}

func VerifyEmail(c *gin.Context) {
	repo := c.MustGet("repo").(Repo)
	publisher := c.MustGet("publisher").(Publisher)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	// Get the query string arguments, if any
	qs := c.Request.URL.Query()
	email := strings.ToLower(qs.Get("email"))
	code := qs.Get("code")

	if code != "" && email != "" {

		userId, err := repo.FindEmail(email)

		if err != nil && err.Error() != "not found" {
			c.JSON(http.StatusBadRequest, NewError(ErrCodeValueRequired, fmt.Sprintf("email verification failed")))
			return
		}

		if userId != uuid.Nil {
			user, err := repo.GetCredentials(userId)

			if err != nil && err.Error() != "not found" {
				c.JSON(http.StatusBadRequest, NewError(ErrCodeValueRequired, fmt.Sprintf("email verification failed")))
				return
			}

			if user.EmailVerificationCode == code {
				user.IsEmailVerified = true
				repo.SaveCredentials(userId, user)

				if c.Request.Header.Get("CID") == "" {
					go publisher.PublishMessage(NewEmailVerifiedEvent(userId, email, uuid.Nil))
				}

				c.JSON(http.StatusOK, "email verified")
				return
			}
		}

		c.JSON(http.StatusBadRequest, NewError(ErrCodeValueRequired, fmt.Sprintf("email verification failed")))
		return
	}

	c.JSON(http.StatusBadRequest, NewError(ErrCodeValueRequired, fmt.Sprintf("email and code are required fields")))
	return
}

func RegisterUser(c *gin.Context) {
	auth := c.MustGet("auth").(Authenticator)
	repo := c.MustGet("repo").(Repo)
	publisher := c.MustGet("publisher").(Publisher)

	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	var view *UserRegistrationView
	c.Bind(&view)

	credentials, validate_err := DecodeRegistrationDetails(view)
	if validate_err != nil {
		c.JSON(http.StatusBadRequest, validate_err)
		return
	}

	credentials.Email = strings.ToLower(credentials.Email)

	duplicateUserId, _ := repo.FindEmail(credentials.Email)

	if duplicateUserId != uuid.Nil {
		c.JSON(http.StatusBadRequest, NewError(ErrCodeAlreadyExists, "user email exists"))
		return
	}

	credentials.Id = uuid.NewV4()
	credentials.IsEmailVerified = false
	credentials.EmailVerificationCode = uuid.NewV4().String()

	err := repo.SaveCredentials(credentials.Id, credentials)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	token, auth_err := auth.Authenticate(view.Email, view.Password, "")
	if auth_err != nil {
		//TODO: Fix this, reveals user details to client
		c.JSON(http.StatusUnauthorized, auth_err.Error())
		return
	}

	go publisher.PublishMessage(NewUserRegisteredEvent(credentials.Id, credentials.Email, uuid.Nil))
	go publisher.PublishMessage(NewEmailVerificationPendingEvent(credentials.Id, credentials.Email, credentials.EmailVerificationCode, uuid.Nil))

	response := AuthResponse{
		Id:    credentials.Id,
		Email: credentials.Email,
		Token: token,
	}

	c.JSON(http.StatusCreated, response)
	return
}

// Parse the request body, load into an Registration structure.
func DecodeRegistrationDetails(reg *UserRegistrationView) (*Credentials, *Error) {

	if reg.Password == "" {
		return nil, NewError(ErrCodeValueRequired, fmt.Sprintf("password is a required field"))
	}

	if reg.Email == "" {
		return nil, NewError(ErrCodeValueRequired, fmt.Sprintf("email is a required field"))
	}

	passwordKey := DeriveKey(reg.Password)

	credentials := &Credentials{Email: reg.Email, Key: passwordKey.Key, Salt: passwordKey.Salt, IsEmailVerified: false}

	return credentials, nil
}

func ChangePassword(c *gin.Context) {
	auth := c.MustGet("auth").(Authenticator)
	repo := c.MustGet("repo").(Repo)
	id := GetIdParam(c.MustGet("userId").(string), c)
	email := c.MustGet("email").(string)

	var view *ChangePasswordView
	c.Bind(&view)

	isValid, validation_err := validatePasswordInfo(view)
	if !isValid {
		c.JSON(http.StatusBadRequest, validation_err.Error())
		return
	}

	_, auth_err := auth.Authenticate(email, view.OldPassword, "")
	if auth_err != nil {
		c.JSON(http.StatusUnauthorized, "")
		return
	}

	credentials, err := repo.GetCredentials(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if credentials != nil {
		passwordKey := DeriveKey(view.NewPassword)

		credentials.Salt = passwordKey.Salt
		credentials.Key = passwordKey.Key

		err := repo.SaveCredentials(id, credentials)

		if err != nil {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}

		token, auth_err := auth.Authenticate(email, view.NewPassword, "")
		if auth_err != nil {
			c.JSON(http.StatusUnauthorized, "")
			return
		}

		response := AuthResponse{
			Id:    id,
			Email: email,
			Token: token,
		}

		c.JSON(http.StatusCreated, response)
		return
	}

	c.JSON(http.StatusInternalServerError, "password change failed")
	return
}

func validatePasswordInfo(form *ChangePasswordView) (bool, *Error) {

	if form.OldPassword == "" {
		return false, NewError(ErrCodeValueRequired, fmt.Sprintf("old password is a required field"))
	}

	if form.NewPassword == "" {
		return false, NewError(ErrCodeValueRequired, fmt.Sprintf("new password is a required field"))
	}

	return true, nil

}
