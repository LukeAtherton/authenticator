package ghttp_test

import (
	. "github.com/lukeatherton/authenticator/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/lukeatherton/authenticator/Godeps/_workspace/src/github.com/onsi/gomega"
	"testing"
)

func TestGHTTP(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GHTTP Suite")
}
