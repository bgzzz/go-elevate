package main

import (
	"net/http"

	"github.com/labstack/echo"
)

//errors
const (
	MarshallingError = iota
	SubrequestsFailed
	HeightCalculationFailed
)

//statuses string
const (
	ERROR_STATUS   = "ERROR"
	SUCCESS_STATUS = "SUCCESS"
)

//Coord holds Longitude and Latitude values
//Is used in as a parameter in requests
type Coord struct {

	// Longitude coordinate.
	Lon float64 `json:"lon" query:"lon"`

	// Latitude coordinate.
	Lat float64 `json:"lat" query:"lat"`
}

//HeightItem a struct that holds information about point
//coordinate and ground height
type HeightItem struct {
	Point  Coord   `json:"point"`
	Height float64 `json:"height"`

	//Item specififc error refers to subreqeusts
	Error *ValidationError `json:"error,omitempty"`
}

//HeightResponse is used on response of the API and holds list of
//HeightItem and optionally error descriptor
type HeightResponse struct {
	Items []HeightItem `json:"items,omitempty"`

	//Error is request specific error refers to the request itself
	Error *ValidationError `json:"error,omitempty"`
}

//HeightRequest is used on request of the API and holds list of
//coords that needs calculation of height
type HeightRequest struct {
	Coords []*Coord `json:"coords"`
}

// A ValidationError is a response to represent error
//
// trololo:response validationError
type ValidationError struct {
	Code        int    `json:"items"`
	Description string `json:"description"`
}

// GetOneHeightHandler handles the functionality on HTTP
// request for ground heigth of one particular point
// Response:
//        200: heighResponse
//		  400: validationError
func GetOneHeightHandler(c echo.Context) error {
	var rsp HeightResponse

	coord := new(Coord)
	if err := c.Bind(coord); err != nil {
		return c.JSON(http.StatusBadRequest, HeightResponse{
			Error: &ValidationError{
				Code:        MarshallingError,
				Description: err.Error(),
			},
		})
	}

	rsp = GetRsp([]*Coord{coord})

	status := http.StatusOK
	if rsp.Error != nil {
		status = http.StatusBadRequest
	}

	return c.JSON(status, rsp)
}

// GetMultipleHeightsHandler handles the functionality on HTTP
// request for ground heigth of multiple particular point
// Response:
//        200: heighResponse
//		  400: validationError
func GetMultipleHeightsHandler(c echo.Context) error {
	var rsp HeightResponse

	req := new(HeightRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, HeightResponse{
			Error: &ValidationError{
				Code:        MarshallingError,
				Description: err.Error(),
			},
		})
	}

	rsp = GetRsp(req.Coords)

	status := http.StatusOK
	if rsp.Error != nil {
		status = http.StatusBadRequest
	}

	return c.JSON(status, rsp)
}

func GetRsp(coords []*Coord) HeightResponse {
	var rsp HeightResponse

	nReqHeights := len(coords)
	msgr := make(chan HeightItem)

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
		}
	}

	return rsp
}
