package main

import (
	"bytes"
	"image/png"
	"io/ioutil"
	"net/http"
	"path"
	"strconv"

	"github.com/bgzzz/go-elevate/mercator"
	log "github.com/sirupsen/logrus"
)

const (
	PNG_URL_PREFIX = `http://s3.amazonaws.com/elevation-tiles-prod/terrarium/`
	ZOOM_LEVEL     = 15
	FILE_RES       = ".png"
)

// GetHeight sends HeightItem while receive png and calculate height
func GetHeight(coord *Coord, res chan<- HeightItem) {
	var height = 0.0
	var err error
	var vErr *ValidationError

	//calculate coords in format
	m := mercator.NewMercatorWithLatLong(coord.Lat, coord.Lon, ZOOM_LEVEL)

	// sending request

	rsp, err := http.Get(PNG_URL_PREFIX +
		path.Join(strconv.Itoa(int(m.Zoom)), strconv.Itoa(int(m.Tile.X)),
			strconv.Itoa(int(m.Tile.Y))+FILE_RES))

	if err != nil {
		log.Error(err.Error())
	} else {
		log.Debug(rsp.Request.URL.String())

		defer rsp.Body.Close()

		body_byte, err := ioutil.ReadAll(rsp.Body)
		if err != nil {
			panic(err)
		}

		//calculating height
		height, err = calculateHeight(body_byte, m.PixelOnTile)

		if err != nil {
			log.Error(err.Error())
		} else {
			log.WithFields(log.Fields{"lon": coord.Lon,
				"lat":    coord.Lat,
				"zoom":   m.Zoom,
				"height": height}).Debug("height calculated")
		}
	}

	//construct the error if needed
	if err != nil {
		vErr = &ValidationError{
			Code:        HEIGHTS_CALCULATION_FAILED_ERROR,
			Description: err.Error(),
		}

	}

	//send message back
	res <- HeightItem{
		Point:  *coord,
		Height: height,
		Error:  vErr,
	}

}

// calculateHeight return height based on png response body
func calculateHeight(body []byte, pixel mercator.MercatorCoord) (float64, error) {

	read := bytes.NewReader(body)
	img, err := png.Decode(read)
	if err != nil {
		log.Error(err.Error())
		return 0.0, err
	}

	color := img.At(int(pixel.X), int(pixel.Y))
	r, g, b, _ := color.RGBA()

	return float64(uint8(r))*256.0 + float64(uint8(g)) + float64(uint8(b))/256.0 - 32768.0, nil

}
