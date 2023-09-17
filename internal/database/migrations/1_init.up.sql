CREATE TABLE "directories"
(
    "dir_id"       SERIAL PRIMARY KEY,
    "path"         TEXT NOT NULL UNIQUE,
    "last_scanned" TIMESTAMPTZ
);

CREATE TABLE "covers"
(
    "cover_id"      SERIAL PRIMARY KEY,
    "dir_id"        INTEGER NOT NULL,
    "relative_path" TEXT    NOT NULL,
    "filename"      TEXT    NOT NULL,
    "format"        TEXT    NOT NULL,
    "width_px"      INTEGER NOT NULL,
    "height_px"     INTEGER NOT NULL,
    "size_byte"     BIGINT  NOT NULL,
    "hash_sha_256"  TEXT    NOT NULL,
    FOREIGN KEY ("dir_id") REFERENCES "directories" ("dir_id")
);

CREATE TABLE "tracks"
(
    "track_id"       SERIAL PRIMARY KEY,
    "dir_id"         INTEGER NOT NULL,
    "cover_id"       INTEGER,
    "relative_path"  TEXT    NOT NULL,
    "filename"       TEXT    NOT NULL,
    "duration_ms"    BIGINT  NOT NULL,
    "size_byte"      BIGINT  NOT NULL,
    "audio_codec"    TEXT    NOT NULL,
    "bitrate_kbps"   INTEGER NOT NULL,
    "sample_rate_hz" INTEGER NOT NULL,
    "channels"       INTEGER NOT NULL,
    "hash_sha_256"   TEXT    NOT NULL,
    FOREIGN KEY ("dir_id") REFERENCES "directories" ("dir_id"),
    FOREIGN KEY ("cover_id") REFERENCES "covers" ("cover_id")
);
