-- +goose Up
ALTER TABLE quote
RENAME TO quote_old;

CREATE TABLE quote (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    number INTEGER,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    submitter TEXT NOT NULL,
    channel TEXT,
    added TIMESTAMP NOT NULL
);

INSERT INTO quote (
    number, content, author, submitter, added
)
SELECT * FROM quote_old;

DROP TABLE quote_old;

-- +goose Down
ALTER TABLE quote
RENAME TO quote_old;

CREATE TABLE quote (
    number INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    submitter TEXT NOT NULL,
    added TIMESTAMP NOT NULL
);

INSERT INTO quote (
    number, content, author, submitter, added
)
SELECT 
    number, content, author, submitter, added
FROM quote_old;

DROP TABLE quote_old;