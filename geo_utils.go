package main

import (
	"image"
	"image/png"
	"io"
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

// SendMercatorReq send request to get png file form the cloud
// request URL is made via mercator coordinates
func SendMercatorReq(m *mercator.Mercator) (image.Image, error) {
	rsp, err := http.Get(PNG_URL_PREFIX +
		path.Join(strconv.Itoa(int(m.Zoom)), strconv.Itoa(int(m.Tile.X)),
			strconv.Itoa(int(m.Tile.Y))+FILE_RES))

	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	log.Debug(rsp.Request.URL.String())

	var r io.Reader
	r = rsp.Body
	defer rsp.Body.Close()

	img, err := png.Decode(r)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}

	return img, nil
}

// GetHeight return HeightItem over channel
// HeightItem is calculated via evaluation of requested png file
func GetHeight(coord *Coord, res chan<- HeightItem) {
	var err error
	var vErr *ValidationError

	//calculate coords in format
	m := mercator.NewMercatorWithLatLong(coord.Lat, coord.Lon, ZOOM_LEVEL)

	// sending request to cloud to get png file encoding height
	img, err := SendMercatorReq(m)

	//calculate height
	var height = 0.0
	if err == nil {
		height, err = calculateHeight(img, m.PixelOnTile)
	}

	//check for errors
	//construct the error if needed
	if err != nil {
		vErr = &ValidationError{
			Code:        HEIGHTS_CALCULATION_FAILED_ERROR,
			Description: err.Error(),
		}
	} else {
		log.WithFields(log.Fields{"lon": coord.Lon,
			"lat":    coord.Lat,
			"zoom":   m.Zoom,
			"height": height}).Debug("height calculated")
	}

	//send message back
	res <- HeightItem{
		Point:  *coord,
		Height: height,
		Error:  vErr,
	}

}

// calculateHeight return height based on png response body
func calculateHeight(img image.Image, pixel mercator.MercatorCoord) (float64, error) {

	color := img.At(int(pixel.X), int(pixel.Y))
	r, g, b, _ := color.RGBA()

	return float64(uint8(r))*256.0 + float64(uint8(g)) + float64(uint8(b))/256.0 - 32768.0, nil

}
