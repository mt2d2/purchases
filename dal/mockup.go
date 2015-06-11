package dal

import "github.com/jmoiron/sqlx"

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
INSERT INTO "purchase" VALUES(1,'Test1',20.0,0);
INSERT INTO "purchase" VALUES(2,'Test2',30.0,0);
COMMIT;
`

// GetMockupDB gets a shared, cached sqlite in memory database with the mockup data for testing.
func GetMockupDB() (*sqlx.DB, error) {
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
