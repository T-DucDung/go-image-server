package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/handlers"

	"github.com/gorilla/mux"
)

func HomeHandler(rw http.ResponseWriter, r *http.Request) {
	http.ServeFile(rw, r, "index.html")
}

func main() {
	// read configuration and exit on error
	ReadConfig()

	// setup router
	router := mux.NewRouter()

	// GET / : say hello
	router.HandleFunc("/", HomeHandler).Methods("GET")

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

	router.HandleFunc("/upload", uploadFile).Methods("POST")

	// Where ORIGIN_ALLOWED is like `scheme://dns[:port]`, or `*` (insecure)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED")})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	// add 404 handler
	router.NotFoundHandler = Logger(handleNotFound, "notFound")

	// start server
	log.Fatal(router, http.ListenAndServe(Bind+":"+strconv.Itoa(Port), handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	// Parse our multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.
	r.ParseMultipartForm(10 << 20)
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("img", "upload-*.png")
	if err != nil {
		fmt.Println(err)
	}
	defer tempFile.Close()

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	// return that we have successfully uploaded our file!
	fmt.Fprintf(w, "Successfully Uploaded File\n")
}
