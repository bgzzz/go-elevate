package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/apeyroux/gosm"
	log "github.com/sirupsen/logrus"
)

//for now
const (
	PNG_URL_PREFIX = `http://s3.amazonaws.com/elevation-tiles-prod/terrarium/`
	ZOOM_LEVEL     = 15
	FILE_RES       = ".png"
)

//GetHeight sends HeightItem while receive png and calculate height
func GetHeight(coord *Coord, res chan<- HeightItem) {
	var height = 0.0
	var err error
	var vErr *ValidationError

	//calculate coords in format
	tile := gosm.NewTileWithLatLong(coord.Lat, coord.Lon, ZOOM_LEVEL)

	// sending request

	rsp, err := http.Get(PNG_URL_PREFIX +
		path.Join(strconv.Itoa(tile.Z), strconv.Itoa(tile.X),
			strconv.Itoa(tile.Y)+FILE_RES))

	log.Debug(rsp.Request.URL.String())

	if err != nil {
		log.Error(err.Error())
	} else {
		defer rsp.Body.Close()

		body_byte, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			panic(err)
		}

		//calculating height
		height, err = calculateHeight(body_byte)

		if err != nil {
			log.Error(err.Error())
		} else {
			log.WithFields(log.Fields{"lon": coord.Lon,
				"lat":    coord.Lat,
				"zoom":   tile.Z,
				"height": height}).Debug("height calculated")
		}
	}

	//construct the error if needed
	if err != nil {
		vErr = &ValidationError{
			Code:        HeightCalculationFailed,
			Description: err.Error(),
		}
	} else {

		//send message back
		res <- HeightItem{
			Point:  *coord,
			Height: height,
			Error:  vErr,
		}
	}

}

//calculateHeight return height based on png responce body
func calculateHeight(body []byte) (float64, error) {

	read := bytes.NewReader(body)
	img, err := png.Decode(read)
	if err != nil {
		log.Error(err.Error())
		return 0.0, err
	}

	// var totalR, totalG, totalB float64 = 0.0, 0.0, 0.0
	// for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
	// 	for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {

	// 		color := img.At(x, y)

	// 		r, g, b, _ := color.RGBA()
	// 		totalR += float64(uint8(r))
	// 		totalG += float64(uint8(g))
	// 		totalB += float64(uint8(b))
	// 	}
	// }

	// take the middle of the square
	color := img.At(128, 128)
	r, g, b, _ := color.RGBA()

	// x, y := float64(img.Bounds().Max.X), float64(img.Bounds().Max.Y)
	// return ((totalR/(x*y))*256.0 + (totalG / (x * y)) + (totalB/(x*y))/256.0) - 32768.0, nil

	return float64(uint8(r))*256.0 + float64(uint8(g)) + float64(uint8(b))/256.0 - 32768.0, nil

}
