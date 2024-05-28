package rabbitmq

import (
	"bytes"
	"encoding/gob"
)

func DecodeAllUserMetricsMessage(b []byte) ([]UserMetricsMessage, error) {
	var res []UserMetricsMessage
	if err := gob.NewDecoder(bytes.NewReader(b)).Decode(&res); err != nil {
		return []UserMetricsMessage{}, err
	}

	return res, nil
}
