// Copyright (c) Luke Atherton 2015

package authenticator

import (
	. "github.com/lukeatherton/domain-events"
	. "github.com/lukeatherton/identity"
)

const (
	SERVICE_NAME = "authenticator"
)

type DeliveryOpened struct {
	*MessageHeader  `json:"header" bson:"header"`
	Id              ID     `json:"id" bson:"id"`
	ItemDescription string `json:"item_description" bson:"item_description"`
	PickupAddress   string `json:"pickup_address" bson:"pickup_address"`
	DeliveryAddress string `json:"delivery_address" bson:"delivery_address"`
}

func NewDeliveryOpenedEvent(id ID, itemDescription string, pickupAddress string, deliveryAddress string, senderId ID) DeliveryOpened {
	return DeliveryOpened{
		MessageHeader:   BuildHeader("Delivery.Opened", &EventSource{Service: SERVICE_NAME, UserId: senderId}),
		Id:              id,
		ItemDescription: itemDescription,
		PickupAddress:   pickupAddress,
		DeliveryAddress: deliveryAddress,
	}
}

//=====================================================================================

type UserRegistered struct {
	*MessageHeader `json:"header" xml:"header"`
	Id             ID     `json:"id" xml:"id"`
	Email          string `json:"email" xml:"email"`
}

func NewUserRegisteredEvent(id ID, email string, senderId ID) UserRegistered {
	return UserRegistered{
		MessageHeader: BuildHeader("User.Registered", &EventSource{Service: SERVICE_NAME, UserId: senderId}),
		Id:            id,
		Email:         email,
	}
}

//=====================================================================================

type EmailVerificationPending struct {
	*MessageHeader        `json:"header" xml:"header"`
	Id                    ID     `json:"id" xml:"id"`
	Email                 string `json:"email" xml:"email"`
	EmailVerificationCode string `json:"email_verification_code" xml:"email_verification_code"`
}

func NewEmailVerificationPendingEvent(id ID, email string, emailVerificationCode string, senderId ID) EmailVerificationPending {
	return EmailVerificationPending{
		MessageHeader: BuildHeader("Email.Verification.Pending", &EventSource{Service: SERVICE_NAME, UserId: senderId}),
		Id:            id,
		Email:         email,
		EmailVerificationCode: emailVerificationCode,
	}
}

//=====================================================================================

type EmailVerified struct {
	*MessageHeader `json:"header" xml:"header"`
	Id             ID     `json:"id" xml:"id"`
	Email          string `json:"email" xml:"email"`
}

func NewEmailVerifiedEvent(id ID, email string, senderId ID) EmailVerified {
	return EmailVerified{
		MessageHeader: BuildHeader("Email.Verified", &EventSource{Service: SERVICE_NAME, UserId: senderId}),
		Id:            id,
		Email:         email,
	}
}
