package middleware

type StopSign interface {
	// set the stop sign, if the stop sign has been signed, return false
	Sign() bool
	//get the info whether the stop sign has been sent
	Signed() bool
	//reset the stop sign
	Reset()
	//code represents the code of stop sign processor
	Deal(code string)
	//get the stop sign processed count from one stop sign processor
	DealCount(code string) uint32
	//get the total count of stop sign processed by all stop sign processors
	DealTotal() uint32
	//get the summary info
	Summary() string
}
