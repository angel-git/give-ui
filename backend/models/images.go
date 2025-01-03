package models

type CacheImageResponse struct {
	Error       *string `json:"error"`
	ImageBase64 *string `json:"imageBase64"`
}
