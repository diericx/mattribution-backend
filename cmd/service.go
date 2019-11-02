package main

import (
	"fmt"

	"github.com/diericx/tracker/backend/pkg/conversionrule"
	"github.com/diericx/tracker/backend/pkg/track"
)

// Service holds the high level business logic
type Service struct {
	tracksRepo track.Repository
	crRepo     conversionrule.Repository
}

// NewService creates a new Service object
func NewService(tracksRepo track.Repository, crRepo conversionrule.Repository) Service {
	return Service{
		tracksRepo: tracksRepo,
		crRepo:     crRepo,
	}
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Tracks
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

// New will create and store a new Track object
func (s Service) NewTrack(ownerID int, t track.Track) (id int, err error) {
	valid := t.IsValid()
	if !valid {
		return 0, fmt.Errorf("Invalid track: %v", t)
	}

	return s.tracksRepo.Store(ownerID, t)
}

// GetAll will query for tracks for a specific user
func (s Service) GetTracksByAttributeAndValue(userID int, attribute string, value string) ([]track.Track, error) {
	return s.tracksRepo.FindByAttributeAndValue(userID, attribute, value)
}

// GetAll will query for tracks for a specific user
func (s Service) GetAllTracks(userID int) ([]track.Track, error) {
	return s.tracksRepo.FindByOwnerID(userID)
}

// GetDailyCounts will aggregate tracks into a daily count
func (s Service) GetDailyTrackCounts(userID int) ([]track.DailyCount, error) {
	return s.tracksRepo.FindDailyCounts(userID)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Conversion Rules
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
