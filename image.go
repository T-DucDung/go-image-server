package main

import (
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"errors"
	"os"
	"github.com/disintegration/imaging"
)

func ImageRead(filePath, ending string, width, height int) (img image.Image, err error) {

	switch (ending) {
	case "png":
		return ImageReadPng(filePath, width, height)
	case "jpg":
		return ImageReadJpg(filePath, width, height)
	}

	return nil, errors.New("unsupported ending")
}

func ImageReadJpg(filePath string, width, height int) (img image.Image, err error) {
	inpImgF, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inpImgF.Close()

	img, err = jpeg.Decode(inpImgF)
	if err != nil {
		return
	}

	img = imageResize(img, width, height)
	return
}

func ImageReadPng(filePath string, width, height int) (img image.Image, err error) {
	inpImgF, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer inpImgF.Close()

	img, err = png.Decode(inpImgF)
	if err != nil {
		return
	}

	img = imageResize(img, width, height)
	return
}

func imageResize(inpImg image.Image, width, height int) (oupImg image.Image) {
	return imaging.Resize(inpImg, width, height, imaging.Box)
}

func ImageWriteJpg(oupWriter io.Writer, image image.Image) (err error) {
	return jpeg.Encode(oupWriter, image, &jpeg.Options{Quality: JpgQuality})
}

func ImageWritePng(oupWriter io.Writer, image image.Image) (err error) {
	return png.Encode(oupWriter, image)
}
