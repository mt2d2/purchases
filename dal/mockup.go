package dal

import (
	"time"

	"github.com/jmoiron/sqlx"
)

import _ "github.com/mattn/go-sqlite3" // We need this for our mockup.

// Mockup database for testing
const MockupDB = `
PRAGMA foreign_keys=OFF;
BEGIN TRANSACTION;
CREATE TABLE "purchase" (
        id    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
        name  TEXT NOT NULL,
        cost  REAL NOT NULL,
        time_bought   TIMESTAMP NOT NULL
);
INSERT INTO "purchase" VALUES(1,'Test1',20.0,'2015-10-14 22:12:33');
INSERT INTO "purchase" VALUES(2,'Test2',30.0,'2015-10-14 22:12:28');
COMMIT;
`

func mockupDB() (*sqlx.DB, error) {
	// :memory: databases aren't shared amongst connections
	// https://groups.google.com/forum/#!topic/golang-nuts/AYZl1lNxCfA
	db, err := sqlx.Open("sqlite3", "file:dummy.db?mode=memory&cache=shared")
	if err != nil {
		return nil, err
	}
	if _, err := db.Exec(MockupDB); err != nil {
		return nil, err
	}

	return db, nil
}

func mockup1() *Purchase {
	return &Purchase{
		uint64(1), "Test1", 20.0, time.Date(2015, 10, 14, 22, 12, 33, 0, time.UTC),
	}
}

// GetMockup2 returns the second mockup in the db
func mockup2() *Purchase {
	return &Purchase{
		uint64(2), "Test2", 30.0, time.Date(2015, 10, 14, 22, 12, 28, 0, time.UTC),
	}
}
