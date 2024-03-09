CREATE TABLE directories
(
    id            SERIAL PRIMARY KEY,
    parent_dir_id INTEGER   NULL,
    relative_path TEXT      NOT NULL,
    last_scanned  TIMESTAMP NULL,
    FOREIGN KEY (parent_dir_id) REFERENCES directories (id),
    UNIQUE (relative_path, parent_dir_id)
);
CREATE INDEX idx_directories_parent ON directories (parent_dir_id);

CREATE TYPE extension_type AS ENUM ('unknown', 'audio', 'video', 'image', 'document');

CREATE TABLE extensions
(
    id   SERIAL PRIMARY KEY,
    name TEXT UNIQUE    NOT NULL,
    type extension_type NOT NULL
);
CREATE INDEX idx_extensions_type ON extensions (type);

CREATE TABLE files
(
    id           SERIAL PRIMARY KEY,
    dir_id       INTEGER   NOT NULL,
    name         TEXT      NOT NULL,
    extension_id INTEGER   NOT NULL,
    size_byte    BIGINT    NOT NULL,
    last_update  TIMESTAMP NOT NULL,
    FOREIGN KEY (dir_id) REFERENCES directories (id),
    UNIQUE (dir_id, name)
);
CREATE INDEX idx_files_dir ON files (dir_id);