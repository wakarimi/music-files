package models

type Track struct {
	TrackId      int    `db:"track_id"`
	DirId        int    `db:"dir_id"`
	CoverId      *int   `db:"cover_id"`
	RelativePath string `db:"relative_path"`
	Filename     string `db:"filename"`
	DurationMs   int64  `db:"duration_ms"`
	Size         int64  `db:"size"`
	AudioCodec   string `db:"audio_codec"`
	BitrateKbps  int    `db:"bitrate_kbps"`
	SampleRateHz int    `db:"sample_rate_hz"`
	Channels     int    `db:"channels"`
	HashSha256   string `db:"hash_sha_256"`
}
