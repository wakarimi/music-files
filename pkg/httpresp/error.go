package httpresp

type Error struct {
	Message          string  `json:"message"`
	LocalizedMessage *string `json:"localizedMessage,omitempty"`
	Error            string  `json:"Error"`
}
