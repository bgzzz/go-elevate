package main

import (
	"fmt"
	"testing"
)

// Benchmark data is taken from here
// https://www.ngs.noaa.gov/cgi-bin/ds_mark.prl?PidBox=LW5521
func TestCalculateHeight(t *testing.T) {
	var benchmarkHeight = 1
	var benchmarkLat, benchmarkLong = 41.74774375, -71.31694444
	// sending tasks for request and calculation
	msgr := make(chan HeightItem)

	go GetHeight(&Coord{Lat: benchmarkLat,
		Lon: benchmarkLong,
	}, msgr)

	item := <-msgr

	t.Log(fmt.Sprintf("benchmark %+v", benchmarkHeight))
	t.Log(fmt.Sprintf("resp coords %+v", item))

	//TODO check if it is network problems then skip the test

	if item.Error != nil {
		t.Log(fmt.Sprintf("%+v", item.Error))
		t.Fail()
	}

	// 1m max error
	// error might be due to non accurate tail of the lat and long
	if int(item.Height) != benchmarkHeight {
		t.Log("Height is not equal to benchmark")
		t.Fail()
	}
}
