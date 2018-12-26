package main

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
)

//errors
const (
	MarshallingError = iota
	SubrequestsFailed
	HeightCalculationFailed
	TimerExpired
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
	Code        int    `json:"items"`
	Description string `json:"description"`
}

// GetOneHeightHandler handles the functionality on HTTP
// request for ground height of one particular point
// Response:
//        200: heighResponse
//		  400: validationError
func GetOneHeightHandler(c echo.Context) error {
	coord := new(Coord)
	if err := c.Bind(coord); err != nil {
		return ResponseError(c, MarshallingError,
			err.Error())
	}

	return ResponseHeights(c, []*Coord{coord})
}

// GetMultipleHeightsHandler handles the functionality on HTTP
// request for ground height of multiple particular point
// Response:
//        200: heighResponse
//		  400: validationError
func GetMultipleHeightsHandler(c echo.Context) error {
	req := new(HeightRequest)
	if err := c.Bind(req); err != nil {
		return ResponseError(c, MarshallingError,
			err.Error())
	}

	return ResponseHeights(c, req.Coords)
}

// GetRsp return response based on the requested coordinates
func GetRsp(coords []*Coord) HeightResponse {
	var rsp HeightResponse

	nReqHeights := len(coords)
	msgr := make(chan HeightItem)

	timer := time.NewTimer(time.Second * time.Duration(TIMER_DURATION_SEC))
	// sending tasks for request and calculation
	for i := 0; i < nReqHeights; i++ {
		go GetHeight(coords[i], msgr)
	}

	//receiving the results
	for nReqHeights > 0 {
		select {
		case item := <-msgr:
			{
				if item.Error != nil {
					rsp.Error = &ValidationError{
						Code:        SubrequestsFailed,
						Description: item.Error.Description,
					}
				}

				rsp.Items = append(rsp.Items, item)
				nReqHeights--
			}
		case <-timer.C:
			rsp.Error = &ValidationError{
				Code:        TimerExpired,
				Description: string(TIMER_DURATION_SEC) + " sec expired",
			}
		}
	}

	return rsp
}

// ResponseError sends error structure as HTTP response
func ResponseError(c echo.Context,
	errCode int, description string) error {
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
