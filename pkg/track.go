package track

import (
	"database/sql"
	"errors"
	"fmt"
)

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Entities
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

type Track struct {
	ID       int    `json:"id"`
	UserID   string `json:"userId"`
	FpHash   string `json:"fpHash"`  // fingerprint hash
	PageURL  string `json:"pageURL"` // optional (website specific)
	Path     string `json:"path"`    // optional ()
	Referrer string `json:"referrer"`
	Extra    string `json:"extra"` // (optional) extra json
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Domain
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

// Repository repository interface to model how we interact with our repo (storage)
type Repository interface {
	// Find(id int) (track, error)
	Store(t Track) (id int, err error)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Use case
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

type Service struct {
	r Repository
}

func NewService(r Repository) Service {
	return Service{
		r: r,
	}
}

// IsValid checks to see if a track object is valid
func (s Service) IsValid(t Track) bool {
	if t.UserID == "" && t.FpHash != "" {
		return false
	}

	return true
}

func (s Service) New(t Track) (id int, err error) {
	valid := s.IsValid(t)
	if !valid {
		return 0, errors.New("Invalid track")
	}

	return s.r.Store(t)
}

// =~=~=~=~=~=~=~=~=~=~=~=~=~=~
// Framework & Driver
// =~=~=~=~=~=~=~=~=~=~=~=~=~=~

// PostgresRepo repo implemented in postgres
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
func (p PostgresRepo) Store(t Track) (id int, err error) {
	sqlStatement :=
		`INSERT INTO public.tracks (user_id, fp_hash, page_url, path, referrer, extra)
	VALUES($1, $2, $3, $4, $5, $6)
	RETURNING id`

	id = 0
	//  Default json
	if t.Extra == "" {
		t.Extra = "{}"
	}

	err = p.db.QueryRow(sqlStatement, t.UserID, t.FpHash, t.PageURL, t.Path, t.Referrer, t.Extra).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}
