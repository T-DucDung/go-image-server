package main

import (
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

func main() {
	// setup router
	router := mux.NewRouter()

	// GET / : say hello
	router.HandleFunc("/", Logger(handleHomeGet, "home")).Methods("GET")

	// GET /raw/<filePath>.{jpg,png,svg,pdf} : output image as it is on the disk if it exists in the requested format
	router.HandleFunc(
		"/raw/{filePath:[a-zA-Z0-9/_\\-\\.]+}.{extension:(?:jpg|png|svg|pdf)}",
		Logger(handleRaw, "raw"),
	).Methods("GET")

	// GET /<width>w/<filePath>.{jpg,png} : resize image to given width and reencode it to the desired output format
	router.HandleFunc(
		"/{width:[0-9]+}w/{filePath:[a-zA-Z0-9/_\\-\\.]+}.{extension:(?:jpg|png)}",
		Logger(handleFixedWidth, "fixedWidth"),
	).Methods("GET")

	// GET /<height>p/<filePath>.{jpg,png} : resize image to given height and reencode it to the desired output format
	router.HandleFunc(
		"/{height:[0-9]+}p/{filePath:[a-zA-Z0-9/_\\-\\.]+}.{extension:(?:jpg|png)}",
		Logger(handleFixedHeight, "fixedHeight"),
	).Methods("GET")

	// add 404 handler
	router.NotFoundHandler = Logger(handleNotFound, "notFound")

	// start server
	log.Fatal(router, http.ListenAndServe(Bind+":"+strconv.Itoa(Port), router))
}
