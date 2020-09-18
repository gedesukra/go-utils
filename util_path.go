package goutils

import (
	// "fmt"

	"net/http"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

func GetRootPath() string {
	_, dirname, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(dirname)
	rootPath := filepath.Dir(currentDir)
	return rootPath
}

func ShiftPath(p string) (head, tail string) {
	p = path.Clean("/" + p)
	i := strings.Index(p[1:], "/") + 1
	if i <= 0 {
		return p[1:], "/"
	}
	return p[1:i], p[i:]
}

func GetCompleteUrl(r *http.Request) string {
	var head string
	head, r.URL.Path = ShiftPath(r.URL.Path)
	completeUrl := "/" + head + r.URL.Path

	if completeUrl == "//" {
		completeUrl = "/"
	}

	return completeUrl
}
