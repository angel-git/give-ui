package images

import (
	"spt-give-ui/backend/api"
)

type ServerImageLoader struct{}

func (e *ServerImageLoader) LoadImage(url string, sessionId string, imageHash string) (string, error) {
	return api.LoadImage(url, sessionId, imageHash)
}
