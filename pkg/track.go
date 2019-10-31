package track

import (
	"database/sql"
	"fmt"
	"time"
)

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Domain
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

// Track holds tracking data for a specific event
type Track struct {
	ID              int       `json:"id"`
	UserID          string    `json:"userId"`
	FpHash          string    `json:"fpHash"`   // fingerprint hash
	PageURL         string    `json:"pageURL"`  // optional (website specific)
	PagePath        string    `json:"pagePath"` // optional ()
	PageTitle       string    `json:"pageTitle"`
	PageReferrer    string    `json:"pageReferrer"`
	Event           string    `json:"event"`
	CampaignSource  string    `json:"campaignSource"`
	CampaignMedium  string    `json:"campaignMedium"`
	CampaignName    string    `json:"campaignName"`
	CampaignContent string    `json:"campaignContent"`
	SentAt          time.Time `json:"sentAt"`
	IP              string
	Extra           string `json:"extra"` // (optional) extra json
}

// Repository repository interface to model how we interact with our repo (storage)
type Repository interface {
	// Find(id int) (track, error)
	Store(ownerID int, t Track) (id int, err error)
	FindByOwnerID(ownerID int) ([]Track, error)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Use case
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

// Service holds the high level business logic
type Service struct {
	r Repository
}

// NewService creates a new Service object
func NewService(r Repository) Service {
	return Service{
		r: r,
	}
}

// IsValid checks to see if a track object is valid
func (s Service) IsValid(t Track) bool {
	if t.UserID == "" && t.FpHash == "" {
		return false
	}

	return true
}

// New will create and store a new Track object
func (s Service) New(ownerID int, t Track) (id int, err error) {
	valid := s.IsValid(t)
	if !valid {
		return 0, fmt.Errorf("Invalid track: %v", t)
	}

	return s.r.Store(ownerID, t)
}

// GetTracksForUser will query for tracks for a specific user
func (s Service) GetTracksForUser(userID int) ([]Track, error) {
	return s.r.FindByOwnerID(userID)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Framework & Driver
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

// PostgresRepo is a repo implemented in postgres
type PostgresRepo struct {
	db *sql.DB
}

// NewPostgresRepo inits a new postgres repo object
func NewPostgresRepo(host string, port int, user string, password string, dbName string) (PostgresRepo, error) {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbName)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return PostgresRepo{}, err
	}

	err = db.Ping()
	if err != nil {
		return PostgresRepo{}, err
	}

	return PostgresRepo{
		db: db,
	}, nil
}

// GetDB returns the db object
func (p PostgresRepo) GetDB() *sql.DB {
	return p.db
}

// Store stores a new track object
func (p PostgresRepo) Store(ownerID int, t Track) (id int, err error) {
	sqlStatement :=
		`INSERT INTO public.tracks (owner_id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, received_at, extra)
	VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
	RETURNING id`

	id = 0
	//  Default json
	if t.Extra == "" {
		t.Extra = "{}"
	}

	err = p.db.QueryRow(sqlStatement, ownerID, t.UserID, t.FpHash, t.PageURL, t.PagePath, t.PageReferrer, t.PageTitle, t.Event, t.CampaignSource, t.CampaignMedium, t.CampaignName, t.CampaignContent, t.SentAt.Format(time.RFC3339), time.Now().Format(time.RFC3339), t.Extra).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// FindByOwnerID finds all track objects by owner id
func (p PostgresRepo) FindByOwnerID(ownerID int) ([]Track, error) {
	sqlStatement :=
		`SELECT id, user_id, fp_hash, page_url, page_path, page_referrer, page_title, event, campaign_source, campaign_medium, campaign_name, campaign_content, sent_at, extra from public.tracks
		WHERE owner_id = $1`

	rows, err := p.db.Query(sqlStatement, ownerID)
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	defer rows.Close()

	tracks := []Track{}
	for rows.Next() {
		var t Track
		err = rows.Scan(&t.ID, &t.UserID, &t.FpHash, &t.PageURL, &t.PagePath, &t.PageReferrer, &t.PageTitle, &t.Event, &t.CampaignSource, &t.CampaignMedium, &t.CampaignName, &t.CampaignContent, &t.SentAt, &t.Extra)
		if err != nil {
			// handle this error
			return nil, err
		}
		tracks = append(tracks, t)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return tracks, nil
}
