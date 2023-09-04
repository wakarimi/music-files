CREATE TABLE "directories"
(
    "dir_id"       SERIAL PRIMARY KEY,
    "path"         TEXT        NOT NULL UNIQUE,
    "date_added"   TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_scanned" TIMESTAMPTZ
);

CREATE TABLE "music_files"
(
    "music_file_id" SERIAL PRIMARY KEY,
    "dir_id"        INTEGER     NOT NULL,
    "path"          TEXT        NOT NULL,
    "size"          BIGINT      NOT NULL,
    "format"        TEXT        NOT NULL,
    "date_added"    TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("dir_id") REFERENCES "directories" ("dir_id")
);

CREATE TABLE "cover_files"
(
    "cover_file_id"   SERIAL PRIMARY KEY,
    "dir_id"     INTEGER     NOT NULL,
    "path"       TEXT        NOT NULL,
    "size"       BIGINT      NOT NULL,
    "format"     TEXT        NOT NULL,
    "date_added" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY ("dir_id") REFERENCES "directories" ("dir_id")
);
