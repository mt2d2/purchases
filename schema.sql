CREATE TABLE "purchase" (
        id    INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT UNIQUE,
        name  TEXT NOT NULL,
        cost  REAL NOT NULL,
        time_bought   TIMESTAMP NOT NULL
);

CREATE INDEX time_bough_index ON purchase (time_bought);
