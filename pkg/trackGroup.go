package trackgroup

// TrackGroup is a high level object that collects individual Track objects
type TrackGroup struct {
	ID   int
	Name string
}

// Repository repository interface to model how we interact with our repo (storage)
type Repository interface {
	// Find(id int) (track, error)
	Store(ownerID int, t TrackGroup) (id int, err error)
	FindByOwnerID(ownerID int) ([]TrackGroup, error)
}
