package util

import (
	"encoding/json"
	"io"
	"spt-give-ui/backend/http"
)

func GetJson(url string, sessionId string, target interface{}) error {
	r, err := http.DoGet(url, sessionId)
	if err != nil {
		return err
	}
	defer r.Body.Close()
	return json.NewDecoder(r.Body).Decode(target)
}

func GetRawBytes(url string, sessionId string) ([]byte, error) {
	r, err := http.DoGet(url, sessionId)
	if err != nil {
		return nil, err
	}
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func ParseByteResponse(profiles []byte, target interface{}) error {
	return json.Unmarshal(profiles, target)
}
