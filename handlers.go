package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var contentTypes = map[string]string{
	"png": "image/png",
	"jpg": "image/jpeg",
	"svg": "image/svg+xml",
	"pdf": "application/pdf",
}

func writeImageHeaders(w http.ResponseWriter, contentType string, lastModified time.Time) {
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", fmt.Sprintf("public, max-age=%d", MaxAge))
	SetLastModified(w, lastModified)
}

func handleRaw(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	filePath := vars["filePath"]
	extension := vars["extension"]

	inpFile := ImageDir + filePath + "." + extension

	fileInfo, err := os.Stat(inpFile)
	if err != nil {
		log.Printf("cannot find file: %v", inpFile)
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// check if StatusNotModified can be answered
	lastModified := GetLastModified(fileInfo)
	if CheckIfModifiedSince(r, lastModified) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	fi, err := os.Open(inpFile)
	if err != nil {
		log.Printf("cannot open file: %v", inpFile)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	defer fi.Close()

	writeImageHeaders(w, contentTypes[extension], lastModified)

	// send file
	_, err = io.Copy(w, fi)

	if err != nil {
		log.Printf("handleRaw: copy error: %v", err)
		return
	}

	return
}

func handleFixedWidth(w http.ResponseWriter, r *http.Request) {
	// read request variables
	vars := mux.Vars(r)
	filePath := vars["filePath"]
	extension := vars["extension"]
	width, _ := strconv.Atoi(vars["width"])

	handleResize(w, r, width, 0, filePath, extension)
}

func handleFixedHeight(w http.ResponseWriter, r *http.Request) {
	// read request variables
	vars := mux.Vars(r)
	filePath := vars["filePath"]
	extension := vars["extension"]
	height, _ := strconv.Atoi(vars["height"])

	handleResize(w, r, 0, height, filePath, extension)
}

func handleResize(w http.ResponseWriter, r *http.Request, width, height int, filePath, extension string) {
	inpFilePathWoExt := ImageDir + filePath

	fileName, ending, fileInfo, err := getFileNameAndStat(inpFilePathWoExt)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	// check if StatusNotModified can be answered
	lastModified := GetLastModified(fileInfo)
	if CheckIfModifiedSince(r, lastModified) {
		w.WriteHeader(http.StatusNotModified)
		return
	}

	// read image
	img, err := ImageRead(fileName, ending, width, height)
	if err != nil {
		log.Printf("cannot read file: %v", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	// output image
	writeImageHeaders(w, contentTypes[extension], lastModified)
	switch extension {
	case "jpg":
		err = ImageWriteJpg(w, img)
	case "png":
		err = ImageWritePng(w, img)
	default:
		log.Printf("unsupported extension: %v", extension)
	}

	if err != nil {
		log.Printf("handleFixedHeight: error while writing output: %v", err)
		return
	}
}

var fileEndings = []string{"png", "jpg"}

func getFileNameAndStat(inpFilePathWoExt string) (fileName, ending string, fileInfo os.FileInfo, err error) {
	for _, ending = range fileEndings {
		fileName = inpFilePathWoExt + "." + ending
		if fileInfo, err = os.Stat(fileName); err == nil {
			return
		}
	}

	return "", "", nil,
		errors.New(fmt.Sprintf("cannot find inpFilePathWoExt=%v", inpFilePathWoExt))
}

func handleNotFound(w http.ResponseWriter, r *http.Request) {
	http.Error(w, fmt.Sprintf("Route not found: %v", r.RequestURI), http.StatusNotFound)
}
