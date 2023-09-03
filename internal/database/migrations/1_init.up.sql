CREATE TABLE "files"(
    "file_id" SERIAL PRIMARY KEY,
    "path" TEXT NOT NULL,
    "size" BIGINT NOT NULL,
    "format" TEXT NOT NULL,
    "date_added" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "hash" TEXT NOT NULL
);
