CREATE TABLE "directories"(
    "dir_id" SERIAL PRIMARY KEY,
    "path" TEXT NOT NULL,
    "date_added" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "last_scanned" TIMESTAMPTZ
);

CREATE TABLE "music_files"(
    "music_file_id" SERIAL PRIMARY KEY,
    "path" TEXT NOT NULL,
    "size" BIGINT NOT NULL,
    "format" TEXT NOT NULL,
    "date_added" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "hash" TEXT NOT NULL
);

CREATE TABLE "covers"(
    "cover_id" SERIAL PRIMARY KEY,
    "path" TEXT NOT NULL,
    "size" BIGINT NOT NULL,
    "format" TEXT NOT NULL,
    "date_added" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "hash" TEXT NOT NULL
);
