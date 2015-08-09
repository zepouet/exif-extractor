package api

import (
	"log"
	"github.com/rwcarlsen/goexif/exif"
)

type ExifInfo struct {
	FileName     string
	Path         string
	CameraModel  string
	Focal        int64
	Aperture     int64
	ISO          string
	TimeShooting string
}

func (info *ExifInfo) Decode(x *exif.Exif) {

	// normally, don't ignore errors!
	camModel, _ := x.Get(exif.Model)
	info.CameraModel, _ = camModel.StringVal()

	// retrieve first (only) rat. value
	focal, _ := x.Get(exif.FocalLength)
	numer, _, _ := focal.Rat2(0)
	//fmt.Printf("\nFocal : %v/%v", numer, denom)
	info.Focal = numer

	// retrieve first (only) rat. value
	aperture, _ := x.Get(exif.FNumber)
	numer, _, _ = aperture.Rat2(0)
	//fmt.Printf("\nAperture : %v/%v", numer, denom)
	info.Aperture = numer

	iso, _ := x.Get(exif.ISOSpeedRatings)
	//fmt.Printf("\n%v", iso)
	info.ISO = iso.String()

}

func (e ExifInfo) ToString() {
	log.Println(e.FileName + " :: " + e.ISO)
}




