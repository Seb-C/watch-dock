package pkg

import (
	"testing"
)

func TestGetHello(t *testing.T) {
	t.Run("Should return hello world", func(t *testing.T) {
		if got, expect := GetHello(), "Hello World!"; got != expect {
			t.Fatalf("Expected '%+v', got '%+v'", expect, got)
		}
	})
}
