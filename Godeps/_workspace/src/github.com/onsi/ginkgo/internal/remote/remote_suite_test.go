package remote_test

import (
	. "github.com/lukeatherton/authenticator/Godeps/_workspace/src/github.com/onsi/ginkgo"
	. "github.com/lukeatherton/authenticator/Godeps/_workspace/src/github.com/onsi/gomega"
	"testing"
)

func TestRemote(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Remote Spec Forwarding Suite")
}
