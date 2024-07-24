package helper

import (
	"strings"
)

func IsImage(filename string) bool {
	extensions := []string{".jpg", ".jpeg", ".png", ".gif", ".bmp", ".svg"}
	ext := strings.ToLower(filename[len(filename)-4:])
	for _, e := range extensions {
		if ext == e {
			return true
		}
	}
	return false
}
