package cosmos

import (
	"testing"
)

const (
	testKey = "C2y6yDjf5/R+ob0N8A7Cgv30VRDJIWEHLM+4QDU5DE2nQ9nDuVTqobD4b8mGGyPMbIZnqyMsEcaGQy67XIw/Jw=="
)

func TestParseKey(t *testing.T) {
	_, err := ParseKey(testKey)
	if err != nil {
		t.Errorf("failed to parse key: %v", err)
	}
}
