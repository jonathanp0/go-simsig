package gateway

//message received on MVT_ALL_TOC Channel
type TrainMovementMessage struct {
	TrainLocation *TrainLocation `json:"train_location"`
	TrainDelay    *TrainDelay    `json:"train_delay"`
}

type TrainLocation struct {
	Headcode string `json:"headcode"`
	UID      string `json:"uid"`
	Action   string `json:"action"`
	Location string `json:"location"`
	Platform string `json:"platform"`
	Time     uint64 `json:"time"`
	AspPass  uint8  `json:"aspPass"`
	AspAppr  uint8  `json:"aspAppr"`
}

type TrainDelay struct {
	Headcode string `json:"headcode"`
	UID      string `json:"uid"`
	Delay    int
}

//message received on SimSig channel
type SimSigMessage struct {
	Clock *ClockMsg `json:"clock_msg"`
}

type ClockMsg struct {
	AreaID string `json:"area_id"`
	Clock  uint
}
