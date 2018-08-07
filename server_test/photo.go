package main

import (
	"strings"
	"os"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"
)

type Photo struct {
	Created int64
	Lon float64
	Lat float64
	relativePath string
	Name string
}

type photoFactory struct{
	string
}

func(f photoFactory) accept(names string) bool {
	if strings.HasSuffix(names, "jpg") {
		return true
	}

	return false
}

func (f photoFactory) CreatePhoto(names string) *Photo {
	r, _ := os.Open(names)
	exif.RegisterParsers(mknote.All...)
	x, _ := exif.Decode(r)
	timeTaken, _ := x.DateTime()
	timeTaken2 := timeTaken.Unix()
	lat, long, _ := x.LatLong()
	relativePath := "./" + names

	components := strings.Split(names, "/")

	return &Photo{
		Created: timeTaken2, Lon: long, Lat: lat, relativePath: relativePath,
		Name:components[len(components)-1],
	} //lat, long, relativePath}
}
