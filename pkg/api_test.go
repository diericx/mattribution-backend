package api_test

import (
	"testing"

	api "github.com/diericx/tracker/backend/pkg"
)

func TestNewTrack(t *testing.T) {
	_, err := api.NewTrackPG("localhost")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
}
