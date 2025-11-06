package tasmota

import "encoding/json"

// Version represents the library version.
const Version = "0.1.0"

// UserAgent is the User-Agent header sent with requests.
var UserAgent = "tasmota-go/" + Version

// unmarshalJSON is a helper function to unmarshal JSON with proper error handling.
func unmarshalJSON(data []byte, v interface{}) error {
	if err := json.Unmarshal(data, v); err != nil {
		return NewError(ErrorTypeParse, "failed to parse JSON response", err)
	}
	return nil
}
