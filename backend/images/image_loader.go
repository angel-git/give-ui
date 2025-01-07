package images

type ImageLoader interface {
	LoadImage(source string, sessionId string, imageHash string) (string, error)
}
