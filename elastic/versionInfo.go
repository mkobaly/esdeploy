package elastic

import (
	"time"
)

type VersionInfo struct {
	ID         string    `json:"id"`
	Folder     string    `json:"folder"`
	File       string    `json:"file"`
	Machine    string    `json:"machine"`
	DateRunUtc time.Time `json:"dateRunUtc"`
}
