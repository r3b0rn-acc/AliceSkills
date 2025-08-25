package skills

type Event struct {
	Version string   `json:"version"`
	Session struct{} `json:"session"`
	Request struct {
		OriginalUtterance string `json:"original_utterance"`
	} `json:"request"`
}

type Response struct {
	Version string   `json:"version"`
	Session struct{} `json:"session"`
	Result  struct {
		Text       string `json:"text"`
		EndSession bool   `json:"end_session"`
	} `json:"response"`
}
