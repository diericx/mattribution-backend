package conversionrule

import (
	"database/sql"
	"fmt"
)

// ConversionRule is a high level object that collects individual Track objects
type ConversionRule struct {
	ID        int
	OwnerID   int
	Attribute string
	Value     string
}

// Repository repository interface to model how we interact with our repo (storage)
type Repository interface {
	// Find(id int) (track, error)
	Store(ownerID int, cr ConversionRule) (id int, err error)
	FindByOwnerID(ownerID int) ([]ConversionRule, error)
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
func (p PostgresRepo) Store(ownerID int, cr ConversionRule) (id int, err error) {
	sqlStatement :=
		`INSERT INTO public.conversion_rules (owner_id, attribute, value)
	VALUES($1, $2)
	RETURNING id`

	err = p.db.QueryRow(sqlStatement, ownerID, cr.OwnerID, cr.Attribute, cr.Value).Scan(&id)
	if err != nil {
		return id, err
	}

	return id, nil
}

// FindByOwnerID finds all track objects by owner id
func (p PostgresRepo) FindByOwnerID(ownerID int) ([]ConversionRule, error) {
	sqlStatement :=
		`SELECT id, attribute, value from public.tracks
		WHERE owner_id = $1`

	rows, err := p.db.Query(sqlStatement, ownerID)
	if err != nil {
		// handle this error better than this
		return nil, err
	}
	defer rows.Close()

	conversionRules := []ConversionRule{}
	for rows.Next() {
		var cr ConversionRule
		err = rows.Scan(&cr.ID, &cr.Attribute, &cr.Value)
		if err != nil {
			// handle this error
			return nil, err
		}
		conversionRules = append(conversionRules, cr)
	}
	// get any error encountered during iteration
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	return conversionRules, nil
}
