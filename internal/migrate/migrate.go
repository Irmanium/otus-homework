package migrate

import (
	_ "github.com/lib/pq" // for goose
	"github.com/pressly/goose/v3"
)

func Up(dbString, migrationDir string) {
	db, err := goose.OpenDBWithDriver("postgres", dbString)
	if err != nil {
		panic(err)
	}

	defer func() {
		if err = db.Close(); err != nil {
			panic(err)
		}
	}()

	if err = goose.Up(db, migrationDir); err != nil {
		panic(err)
	}
}
