-- +goose Up
CREATE TABLE quote (
    number INTEGER PRIMARY KEY AUTOINCREMENT,
    content TEXT NOT NULL,
    author TEXT NOT NULL,
    submitter TEXT NOT NULL,
    added TIMESTAMP NOT NULL
);

CREATE TABLE message (
    id TEXT PRIMARY KEY,
    content TEXT NOT NULL,
    username TEXT NOT NULL,
    channel TEXT NOT NULL,
    created TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE message;
DROP TABLE quote;