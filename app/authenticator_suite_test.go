// Copyright (c) Luke Atherton 2015

package authenticator_test

import (
	"fmt"
	"testing"

	. "github.com/lukeatherton/authenticator/app"
	"github.com/modocache/gory"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestAuthServer(t *testing.T) {
	defineFactories()
	RegisterFailHandler(Fail)
	RunSpecs(t, "Auth Server Suite")
}

func defineFactories() {

	gory.Define("userRegistration", UserRegistrationView{},
		func(factory gory.Factory) {
			factory["Email"] = gory.Sequence(
				func(n int) interface{} {
					return fmt.Sprintf("latherton%d@example.com", n)
				})
			factory["Password"] = "secret"
		})

	gory.Define("userRegistrationDuplicate", UserRegistrationView{},
		func(factory gory.Factory) {
			factory["Email"] = "latherton@example.com"
			factory["Password"] = "secret"
		})

	gory.Define("userRegistrationDupEmail", UserRegistrationView{},
		func(factory gory.Factory) {
			factory["Email"] = gory.Sequence(
				func(n int) interface{} {
					return fmt.Sprintf("latherton@example.com")
				})
			factory["Password"] = "secretpassword"
		})

	gory.Define("userRegistrationMissingEmail", UserRegistrationView{},
		func(factory gory.Factory) {
			factory["Password"] = "secretpassword"
		})

	gory.Define("loginValid", LoginView{},
		func(factory gory.Factory) {
			factory["Email"] = "latherton@example.com"
			factory["Password"] = "secret"
		})

	gory.Define("passwordChangeRequest", ChangePasswordView{},
		func(factory gory.Factory) {
			factory["OldPassword"] = "secret"
			factory["NewPassword"] = "new secret"
		})

	gory.Define("passwordChangeRequestWrongPassword", ChangePasswordView{},
		func(factory gory.Factory) {
			factory["OldPassword"] = "wrong"
			factory["NewPassword"] = "new secret"
		})

	gory.Define("passwordChangeRequestMissingInfo", ChangePasswordView{},
		func(factory gory.Factory) {
			factory["OldPassword"] = "secret"
		})

}
