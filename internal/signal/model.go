package signal

type Signal struct {
	ID   uint64 `json:"id"`
	Freq uint64 `json:"frequency"`
	Msg  string `json:"message"`
}

type SignalUpdateRequest struct {
	ID  uint64 `json:"id"`
	Msg string `json:"message"`
}
