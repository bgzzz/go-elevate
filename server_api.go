package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/labstack/echo"
)

//errors
const (
	MARSHALLING_ERROR                = "Marshalling failed"
	SUBREQUESTS_FAILED_ERROR         = "Subrequests failed"
	HEIGHTS_CALCULATION_FAILED_ERROR = "Heights calculation failed"
	TIME_EXPIRED_ERROR               = "Timer expired error"
	INTERNAL_ERROR                   = "Internal server error"
)

//query keys
const (
	QUERY_COORDS_KEY = "coords"
)

// limit the request in time
const (
	TIMER_DURATION_SEC = 5
)

// Coord holds Longitude and Latitude values
// It is used as a parameter in requests
type Coord struct {

	// Longitude coordinate.
	Lon float64 `json:"lon" query:"lon"`

	// Latitude coordinate.
	Lat float64 `json:"lat" query:"lat"`
}

// HeightItem a struct that represents information about
// coordinates of the point and ground height
type HeightItem struct {
	Point  Coord   `json:"point"`
	Height float64 `json:"height"`

	//Item specific error refers to sub requests
	Error *ValidationError `json:"error,omitempty"`
}

// HeightResponse is used on response of the API and holds list of
// HeightItem and optionally error descriptor
type HeightResponse struct {
	Items []HeightItem `json:"items,omitempty"`

	//Error is request specific error refers to the request itself
	Error *ValidationError `json:"error,omitempty"`
}

// HeightRequest is used on request of the API and holds list of
// coords that needs calculation of the height
type HeightRequest struct {
	Coords []*Coord `json:"coords"`
}

// A ValidationError is a response to represent error
type ValidationError struct {
	Code        string `json:"code"`
	Description string `json:"description"`
}

// GetHeightsHandler handles the functionality on HTTP
// request for ground height of multiple particular points
// Response:
//        200: heighResponse
//		  400: validationError
func GetHeightsHandler(c echo.Context) error {

	coordsQuery := c.Request().URL.Query().Get(QUERY_COORDS_KEY)

	coords := []*Coord{}
	err := json.Unmarshal([]byte(coordsQuery), &coords)
	if err != nil {
		return ResponseError(c, MARSHALLING_ERROR,
			err.Error())
	}

	return ResponseHeights(c, coords)
}

// PostMultipleHeightsHandler handles the functionality on HTTP
// request for ground height of multiple particular point
// Response:
//        200: heighResponse
//		  400: validationError
func PostMultipleHeightsHandler(c echo.Context) error {
	req := new(HeightRequest)
	if err := c.Bind(req); err != nil {
		return ResponseError(c, MARSHALLING_ERROR,
			err.Error())
	}

	return ResponseHeights(c, req.Coords)
}

// GetRsp return response based on the requested coordinates
func GetRsp(coords []*Coord) HeightResponse {
	var rsp HeightResponse

	msgr := make(chan HeightItem)

	timer := time.NewTimer(time.Second * time.Duration(TIMER_DURATION_SEC))

	// this is done to avoid making requests with same lat lon
	coordsHeightMap := map[Coord]HeightItem{}
	for i := range coords {

		coordsHeightMap[*coords[i]] = HeightItem{}
	}

	nReqHeights := len(coordsHeightMap)

	// sending tasks for request and calculation
	for k, _ := range coordsHeightMap {
		go GetHeight(k, msgr)
	}

	//receiving the results
	for nReqHeights > 0 {
		select {
		case item := <-msgr:
			{
				if item.Error != nil {
					rsp.Error = &ValidationError{
						Code:        SUBREQUESTS_FAILED_ERROR,
						Description: item.Error.Description,
					}
				}

				coordsHeightMap[item.Point] = item
				nReqHeights--
			}
		case <-timer.C:
			rsp.Error = &ValidationError{
				Code:        TIME_EXPIRED_ERROR,
				Description: string(TIMER_DURATION_SEC) + " sec expired",
			}
			return rsp
		}
	}

	// making array of HeighItems in the same order as request
	for i := range coords {
		v := coordsHeightMap[*coords[i]]
		rsp.Items = append(rsp.Items, v)
	}

	return rsp
}

// ResponseError sends error structure as HTTP response
func ResponseError(c echo.Context,
	errCode string, description string) error {
	return c.JSON(http.StatusBadRequest, HeightResponse{
		Error: &ValidationError{
			Code:        errCode,
			Description: description,
		},
	})
}

// ResponseHeights sends heights response structure as HTTP response
func ResponseHeights(c echo.Context, coords []*Coord) error {
	rsp := GetRsp(coords)

	status := http.StatusOK
	if rsp.Error != nil {
		status = http.StatusBadRequest
	}

	return c.JSON(status, rsp)
}
