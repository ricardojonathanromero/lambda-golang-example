package services_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"testing"
)

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Suite Service")
}
