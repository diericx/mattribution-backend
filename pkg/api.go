package api

import (
	"database/sql"
	"fmt"
)

type track struct {
	ID       string
	userID   string
	fpHash   string // fingerprint hash
	pageURL  string // optional (website specific)
	path     string // optional ()
	referrer string
	extra    string // (optional) extra json
}

type trackPG struct {
	db *sql.DB
}

type referrerCount struct {
	referrer string
	count    int
}

func NewTrackPG(ip string) (trackPG, error) {
	connStr := "user=Zac dbname="
	db, err := sql.Open("mattribution", connStr)
	if err != nil {
		return trackPG{}, err
	}
	return trackPG{
		db: db,
	}, nil
}

func (t trackPG) Store(track track) error {
	query := fmt.Sprintf(`
	INSERT INTO tracks(userID, fpHash, pageURL, path, referrer, data)
	VALUES(%s, %s, %s, %s, %s, %s)
	`, track.userID, track.fpHash, track.pageURL, track.path, track.referrer, track.extra)

	rows, err := t.db.Query(query)
	if err != nil {
		return err
	}
	rows.Close()

	return nil
}

// func (t trackPG) GetTopExternalRefererrers() []referrerCount {
// 	t.db.Query(``)
// }
