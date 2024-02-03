package methodHttp

const (
	updateFlag      = "update"
	updateEmpty     = 0
	updateRunning   = 1
	updateCompleted = 2
	updateFailed    = 3
)

var (
	err      error
	jsonText []byte
)
