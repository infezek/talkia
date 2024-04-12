package util_url

import (
	"fmt"
	"os"
)

func New(url string) string {
	if os.Getenv("CHAT_ENVIRONMENT") == "development" {
		return url
	}
	if string(url[0]) != "/" {
		url = "/" + url
	}
	prefix := os.Getenv("CHAT_PREFIX_URL")
	if prefix == "" {
		return url
	}
	if string(prefix[0]) != "/" {
		prefix = "/" + prefix
	}
	if string(prefix[len(prefix)-1]) == "/" {
		prefix = prefix[:len(prefix)-1]
	}
	return fmt.Sprintf("%s%s", prefix, url)
}
