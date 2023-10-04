package telegram

type PollOption struct {
	Text       string `json:"text"`
	VoterCount int32  `json:"voter_count"`
}
