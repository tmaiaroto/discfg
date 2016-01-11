package config

// discfg status codes (not unlike HTTP status codes, but different numbers)
const (
	// EcodeKeyNotFound      = 100
	StatusContinue           = 100
	StatusSwitchingProtocols = 101

	StatusOK       = 200
	StatusCreated  = 201
	StatusAccepted = 202
)

var statusText = map[int]string{
	// EcodeKeyNotFound:      "Key not found",
	StatusContinue:           "Continue",
	StatusSwitchingProtocols: "Switching Protocols",

	StatusOK:       "OK",
	StatusCreated:  "Created",
	StatusAccepted: "Accepted",
}

// StatusText returns a text for the discfg status code.
// It returns the empty string if the code is unknown.
func StatusText(code int) string {
	return statusText[code]
}
