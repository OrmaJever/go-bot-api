package telegram

type PollAnswer struct {
	PollId    string  `json:"poll_id"`
	VoterChat *Chat   `json:"voter_chat"`
	User      *User   `json:"user"`
	OptionIds []int64 `json:"option_ids"`
}
