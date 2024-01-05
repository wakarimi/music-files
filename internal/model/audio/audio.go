package audio

import "time"

type Audio struct {
	ID                int       `db:"id"`
	DirID             int       `db:"dir_id"`
	Filename          string    `db:"filename"`
	Extension         string    `db:"extension"`
	SizeByte          int64     `db:"size_byte"`
	DurationMs        int64     `db:"duration_ms"`
	BitrateKbps       int       `db:"bitrate_kbps"`
	SampleRateHz      int       `db:"sample_rate_hz"`
	ChannelsN         int       `db:"channels_n"`
	SHA256            string    `db:"sha_256"`
	LastContentUpdate time.Time `db:"last_content_update"`
}
