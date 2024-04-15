package utils_test

import (
	"github.com/ricardojonathanromero/lambda-golang-example/get-all-documents-lambda/internal/utils"
	"testing"
)

func TestToString(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		mapInput := map[string]string{
			"message": "hello world",
		}

		result := utils.ToString(mapInput)
		if len(result) == 0 {
			t.Errorf("no len: %d", len(result))
			t.FailNow()
		}
	})

	t.Run("error", func(t *testing.T) {
		value := make(chan int)
		result := utils.ToString(value)
		if len(result) != 0 {
			t.Errorf("no len: %d", len(result))
			t.FailNow()
		}
	})
}
