package main

import (
	"strconv"
	"os"
	"log"
)

var Port int
var Bind string
var ImageDir string
var MaxAge int
var JpgQuality int

func init() {

	var err error
	Port, err = strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		Port = 80
	}

	Bind, bindIsset := os.LookupEnv("BIND")
	if !bindIsset {
		// listen on ipv4 loopback only
		Bind = "127.0.0.1"
	}

	ImageDir = os.Getenv("IMAGE_DIR")
	if _, err := os.Stat(ImageDir); err != nil {
		log.Fatalf("IMAGE_DIR=%v does not exist", ImageDir)
		os.Exit(1)
	}

	MaxAge, err = strconv.Atoi(os.Getenv("MAX_AGE"))
	if err != nil {
		MaxAge = 0
	}

	JpgQuality, err = strconv.Atoi(os.Getenv("JPG_QUALITY"))
	if err != nil {
		JpgQuality = 90
	}

	log.Print("configuration:")
	log.Printf("  Port=%v", Port)
	log.Printf("  Bind=%v", Bind)
	log.Printf("  ImageDir=%v", ImageDir)
	log.Printf("  MaxAge=%v", MaxAge)
	log.Printf("  JpgQuality=%v", JpgQuality)
}
