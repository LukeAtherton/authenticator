// Copyright (c) Luke Atherton 2015

package authenticator

import (
	"encoding/xml"
	"time"
)

type DomainEvent interface {
	GetMessageType() (messageType string)
	GetHeader() (header *MessageHeader)
}

type MessageHeader struct {
	CorrelationId string    `json:"c_id" xml:"c_id"`
	TriggeredById string    `json:"t_id" xml:"t_id"`
	MessageType   string    `json:"message_type" xml:"message_type"`
	TimeStamp     time.Time `json:"timestamp" xml:"timestamp"`
}

func (h *MessageHeader) GetHeader() (header *MessageHeader) {
	return h
}

func (h *MessageHeader) GetMessageType() (messageType string) {
	return h.MessageType
}

//=====================================================================================

type UserRegistered struct {
	XMLName        xml.Name `json:"-" xml:"user_registered"`
	*MessageHeader `json:"header" xml:"header"`
	Id             string `json:"id" xml:"id"`
	Email          string `json:"email" xml:"email"`
}

func NewUserRegisteredEvent(id ID, email string, triggeredBy string) *UserRegistered {
	header := &MessageHeader{CorrelationId: NewSequentialUUID().String(), TriggeredById: triggeredBy, MessageType: "User.Registered", TimeStamp: time.Now().UTC()}
	return &UserRegistered{MessageHeader: header, Id: id.String(), Email: email}
}

//=====================================================================================

type EmailVerificationPending struct {
	XMLName               xml.Name `json:"-" xml:"email_verification_pending"`
	*MessageHeader        `json:"header" xml:"header"`
	Id                    string `json:"id" xml:"id"`
	Email                 string `json:"email" xml:"email"`
	EmailVerificationCode string `json:"email_verification_code" xml:"email_verification_code"`
}

func NewEmailVerificationPendingEvent(id ID, email string, emailVerificationCode string, triggeredBy string) *EmailVerificationPending {
	header := &MessageHeader{CorrelationId: NewSequentialUUID().String(), TriggeredById: triggeredBy, MessageType: "Email.Verification.Pending", TimeStamp: time.Now().UTC()}
	return &EmailVerificationPending{MessageHeader: header, Id: id.String(), Email: email, EmailVerificationCode: emailVerificationCode}
}

//=====================================================================================

type EmailVerified struct {
	XMLName        xml.Name `json:"-" xml:"email_verified"`
	*MessageHeader `json:"header" xml:"header"`
	Id             string `json:"id" xml:"id"`
	Email          string `json:"email" xml:"email"`
}

func NewEmailVerifiedEvent(id ID, email string, triggeredBy string) *EmailVerified {
	header := &MessageHeader{CorrelationId: NewSequentialUUID().String(), TriggeredById: triggeredBy, MessageType: "Email.Verified", TimeStamp: time.Now().UTC()}
	return &EmailVerified{MessageHeader: header, Id: id.String(), Email: email}
}
