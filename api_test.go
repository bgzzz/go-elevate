package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
)

var (
	multipleCoordsReq = `{"coords":[{"lon":138.72905,"lat":35.360638},{"lat":27.986065,"lon":86.922623}]}`
	multipleCoordsRsp = `{"items":[{"point":{"lon":138.72905,"lat":35.360638},"height":3685.69921875},{"point":{"lon":86.922623,"lat":27.986065},"height":8368.4140625}]}`

	multipleCoordsRspMap map[string]HeightItem
	singleCoordRsp       = `{"items":[{"point":{"lon":138.72905,"lat":35.360638},"height":3685.69921875}]}`
)

func getExpectedRsp(js string, t *testing.T) (exp HeightResponse) {
	multipleCoordsRspMap = map[string]HeightItem{}
	err := json.Unmarshal([]byte(js), &exp)
	if err != nil {
		t.Fatalf("error while read %s", err.Error())
	}

	//this is done for simplicity
	for i := range exp.Items {
		multipleCoordsRspMap[fmt.Sprintf("%v_%v",
			exp.Items[i].Point.Lat, exp.Items[i].Point.Lon)] = exp.Items[i]
	}
	return
}

func getRxedRsp(b *bytes.Buffer, t *testing.T) (rxed HeightResponse) {
	body, err := ioutil.ReadAll(b)
	if err != nil {
		t.Fatalf("error while read %s", err.Error())
	}
	err = json.Unmarshal(body, &rxed)
	if err != nil {
		t.Fatalf("error while read %s", err.Error())
	}
	return
}

func prepareContext(req *http.Request, t *testing.T) (c echo.Context,
	rec *httptest.ResponseRecorder) {

	rec = httptest.NewRecorder()

	e := echo.New()
	c = e.NewContext(req, rec)

	return
}

func handlerAssert(c echo.Context,
	rec *httptest.ResponseRecorder,
	expected HeightResponse,
	hand func(c echo.Context) error,
	t *testing.T) {
	// Assertions
	if assert.NoError(t, hand(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)

		rxedRsp := getRxedRsp(rec.Body, t)
		t.Logf("expected: %+v", expected)
		t.Logf("rxed: %+v", rxedRsp)

		for i := range rxedRsp.Items {
			if rxedRsp.Items[i].Error != nil {
				t.Fatalf("err: %+v", rxedRsp.Items[i].Error)
			}

			item, ok := multipleCoordsRspMap[fmt.Sprintf("%v_%v",
				rxedRsp.Items[i].Point.Lat, rxedRsp.Items[i].Point.Lon)]

			if !ok {
				t.Fatalf("there is no such point %+v", rxedRsp.Items[i].Point)
			}

			if !reflect.DeepEqual(rxedRsp.Items[i], item) {
				t.Logf("expected: %+v", item)
				t.Logf("rxed: %+v", rxedRsp.Items[i])
				t.Fatalf("Item is not equal")
			}

		}

	}
}

func TestGetReq(t *testing.T) {
	q := make(url.Values)
	q.Set(QUERY_COORDS_KEY, `[{"lon":138.72905, "lat":35.360638}]`)
	req := httptest.NewRequest(http.MethodGet, `/heights?`+q.Encode(), nil)

	c, rec := prepareContext(req, t)

	expected := getExpectedRsp(singleCoordRsp, t)

	handlerAssert(c, rec, expected, GetHeightsHandler, t)

}

func TestPostMultipleReq(t *testing.T) {
	// Setup
	req := httptest.NewRequest(http.MethodPost, "/heights",
		strings.NewReader(multipleCoordsReq))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	c, rec := prepareContext(req, t)

	expected := getExpectedRsp(multipleCoordsRsp, t)

	handlerAssert(c, rec, expected, PostMultipleHeightsHandler, t)

}
