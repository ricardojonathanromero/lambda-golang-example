package repository_test

import (
	"github.com/jarcoal/httpmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"testing"
)

var _ = BeforeSuite(func() {
	// set http mock handler for dummy tests
	httpmock.ActivateNonDefault(http.DefaultClient)
})

var _ = AfterSuite(func() {
	httpmock.DeactivateAndReset()
})

func TestRepository(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Repository Suite")
}
