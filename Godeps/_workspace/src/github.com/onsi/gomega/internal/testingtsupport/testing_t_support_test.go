package testingtsupport_test

import (
	. "github.com/lukeatherton/authenticator/Godeps/_workspace/src/github.com/onsi/gomega"
	"testing"
)

func TestTestingT(t *testing.T) {
	RegisterTestingT(t)
	Ω(true).Should(BeTrue())
}
