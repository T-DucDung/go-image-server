package main

import (
	"net/http"
	"time"
	"os"
)

func SetLastModified(w http.ResponseWriter, modtime time.Time) {
	if !isZeroTime(modtime) {
		w.Header().Set("Last-Modified", modtime.UTC().Format(http.TimeFormat))
	}
}

// return true if a valid IF-Modified-Since header is sent and it is newer than the given lastModified
func CheckIfModifiedSince(r *http.Request, lastModified time.Time) (isNotModified bool) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return false
	}
	ims := r.Header.Get("If-Modified-Since")
	if ims == "" || isZeroTime(lastModified) {
		return false
	}
	t, err := http.ParseTime(ims)
	if err != nil {
		return false
	}
	// The Date-Modified header truncates sub-second precision, so
	// use mtime < t+1s instead of mtime <= t to check for unmodified.
	return lastModified.Before(t.Add(1 * time.Second))
}

var unixEpochTime = time.Unix(0, 0)

func isZeroTime(t time.Time) bool {
	return t.IsZero() || t.Equal(unixEpochTime)
}

func GetLastModified(info os.FileInfo) time.Time {
	return info.ModTime()
}
