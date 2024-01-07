package response

type Error struct {
	Message string `json:"message"`
	Reason  string `json:"reason"`
}
