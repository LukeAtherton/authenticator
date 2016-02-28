// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"encoding/xml"
	"time"

	. "github.com/lukeatherton/identity"
)

type Credentials struct {
	XMLName               xml.Name  `json:"-" xml:"credentials" bson:"-"`
	Id                    ID        `json:"id" xml:"id" bson:"id,omitempty"`
	Email                 string    `json:"email" xml:"email" bson:"email"`
	Salt                  []byte    `json:"salt" xml:"salt" bson:"salt"`
	Key                   []byte    `json:"key" xml:"key" bson:"key"`
	IsEmailVerified       bool      `json:"isEmailVerified" xml:"isEmailVerified" bson:"isEmailVerified"`
	EmailVerificationCode string    `json:"emailVerificationCode" xml:"emailVerificationCode" bson:"emailVerificationCode"`
	CreatedDate           time.Time `json:"createdDate" xml:"createdDate"  bson:"createdDate"`
	LastModifiedDate      time.Time `json:"lastModifiedDate" xml:"lastModifiedDate"  bson:"lastModifiedDate"`
	ConfirmedDate         time.Time `json:"confirmedDate" xml:"confirmedDate"  bson:"confirmedDate"`
}

type AuthResponse struct {
	XMLName xml.Name `json:"-" xml:"auth_response" bson:"-"`
	Id      ID       `json:"id" xml:"id" bson:"id"`
	Email   string   `json:"email" xml:"email" bson:"email"`
	Token   string   `json:"token" xml:"token" bson:"token"`
}

type UserRegistrationView struct {
	XMLName  xml.Name `json:"-" xml:"user_registration"`
	Email    string   `json:"email" xml:"email"`
	Password string   `json:"password" xml:"password"`
}

type LoginView struct {
	XMLName  xml.Name `json:"-" xml:"login"`
	Email    string   `json:"email" xml:"email"`
	Password string   `json:"password" xml:"password"`
}

type ChangePasswordView struct {
	XMLName     xml.Name `json:"-" xml:"password_change_request"`
	OldPassword string   `json:"oldPassword" xml:"oldPassword"`
	NewPassword string   `json:"newPassword" xml:"newPassword"`
}

type ErrorResponse struct {
	XMLName xml.Name `json:"-" xml:"error_response" bson:"-"`
	Errors  []*Error `json:"errors" xml:"errors" bson:"errors"`
}
