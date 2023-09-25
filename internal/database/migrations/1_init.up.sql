CREATE TABLE directories
(
    dir_id        SERIAL PRIMARY KEY,
    name          TEXT NOT NULL,
    parent_dir_id INTEGER NULL,
    last_scanned  TIMESTAMP NULL,
    FOREIGN KEY (parent_dir_id) REFERENCES directories (dir_id),
    UNIQUE (name, parent_dir_id)
);

CREATE TABLE covers
(
    cover_id            SERIAL PRIMARY KEY,
    dir_id              INTEGER    NOT NULL,
    filename            TEXT       NOT NULL,
    extension           VARCHAR(5) NOT NULL,
    size_byte           BIGINT     NOT NULL,
    width_px            INTEGER    NOT NULL,
    height_px           INTEGER    NOT NULL,
    sha_256             CHAR(64)   NOT NULL,
    last_content_update TIMESTAMP  NOT NULL,
    FOREIGN KEY (dir_id) REFERENCES directories (dir_id),
    UNIQUE (dir_id, filename)
);

CREATE TABLE tracks
(
    track_id            SERIAL PRIMARY KEY,
    dir_id              INTEGER    NOT NULL,
    filename            TEXT       NOT NULL,
    extension           VARCHAR(5) NOT NULL,
    size_byte           BIGINT     NOT NULL,
    duration_ms         BIGINT     NOT NULL,
    bitrate_kbps        INTEGER    NOT NULL,
    sample_rate_hz      INTEGER    NOT NULL,
    channels_n          INTEGER    NOT NULL,
    sha_256             CHAR(64)   NOT NULL,
    last_content_update TIMESTAMP  NOT NULL,
    FOREIGN KEY (dir_id) REFERENCES directories (dir_id),
    UNIQUE (dir_id, filename)
);

CREATE UNIQUE INDEX idx_covers_sha_256 ON covers (sha_256);
CREATE UNIQUE INDEX idx_tracks_sha_256 ON tracks (sha_256);
CREATE INDEX idx_directories_parent_dir_id ON directories (parent_dir_id);
