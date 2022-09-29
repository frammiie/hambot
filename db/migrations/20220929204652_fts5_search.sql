-- +goose Up
CREATE VIRTUAL TABLE quote_fts USING fts5(
    id UNINDEXED,
    number UNINDEXED,
    content,
    author UNINDEXED,
    submitter UNINDEXED,
    channel UNINDEXED,
    added UNINDEXED,
    tokenize = "unicode61 remove_diacritics 0",
    content_rowid='id',
    content='quote',
);

CREATE VIRTUAL TABLE message_fts USING fts5(
    id UNINDEXED,
    content,
    username UNINDEXED,
    channel UNINDEXED,
    created UNINDEXED,
    tokenize = "unicode61 remove_diacritics 0",
    content='message',
);

-- +goose StatementBegin
CREATE TRIGGER quote_fts_ai AFTER INSERT ON quote BEGIN
    INSERT INTO quote_fts(
        rowid, number, content, author, submitter, channel, added
    )
    VALUES (
        new.id, new.number, new.content, new.author, new.submitter,
        new.channel, new.content
    );
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER quote_fts_ad AFTER DELETE ON quote BEGIN
    INSERT INTO quote_fts(
        quote_fts, rowid, number, content, author, submitter, channel, added
    )
    VALUES(
        'delete', new.id, new.number, new.content, new.author, new.submitter,
        new.channel, new.content
    );
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER quote_fts_au AFTER UPDATE ON quote BEGIN
    INSERT INTO quote_fts(
        quote_fts, rowid, number, content, author, submitter, channel, added
    )
    VALUES(
        'delete', old.id, old.number, old.content, old.author, old.submitter,
        old.channel, old.content
    );

    INSERT INTO quote_fts(
        rowid, number, content, author, submitter, channel, added
    )
    VALUES (
        new.id, new.number, new.content, new.author, new.submitter,
        new.channel, new.content
    );
END;
-- +goose StatementEnd

---

-- +goose StatementBegin
CREATE TRIGGER message_fts_ai AFTER INSERT ON message BEGIN
    INSERT INTO message_fts(
        rowid, id, content, username, channel, created
    )
    VALUES (
        new.rowid, new.id, new.content, new.username, new.channel, new.created
    );
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER message_fts_ad AFTER DELETE ON message BEGIN
    INSERT INTO message_fts(
        message_fts, rowid, id, content, username, channel, created
    )
    VALUES(
        'delete', old.rowid, old.id, old.content, old.username,
        old.channel, old.created
    );
END;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER message_fts_au AFTER UPDATE ON message BEGIN
        INSERT INTO message_fts(
        message_fts, rowid, id, content, username, channel, created
    )
    VALUES(
        'delete', old.rowid, old.id, old.content, old.username,
        old.channel, old.created
    );

    INSERT INTO message_fts(
        rowid, id, content, username, channel, created
    )
    VALUES (
        new.rowid, new.id, new.content, new.username, new.channel, new.created
    );
END;
-- +goose StatementEnd

INSERT INTO quote_fts (channel, content)
SELECT channel, content FROM quote;

INSERT INTO message_fts (channel, content)
SELECT channel, content FROM message;

-- +goose Down
DROP TABLE quote_fts;
DROP TABLE message_fts;

DROP TRIGGER quote_fts_ai;
DROP TRIGGER quote_fts_ad;
DROP TRIGGER quote_fts_au;

DROP TRIGGER message_fts_ai;
DROP TRIGGER message_fts_ad;
DROP TRIGGER message_fts_au;
