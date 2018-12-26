package mercator

import (
	"fmt"
	"reflect"
	"testing"
)

func TestMercator(t *testing.T) {
	lat, long := 41.850, -87.650
	zoom := 15

	// benchmark is taken from here
	// https://developers.google.com/maps/documentation/javascript/examples/map-coordinates
	benchmarkMercator := &Mercator{
		Latitude:  lat,
		Longitude: long,
		Zoom:      uint(zoom),
		//this one a little bit different in the tail of the float
		World: struct {
			X float64
			Y float64
		}{
			X: 65.67111111111112,
			Y: 95.17492654697409,
		},
		Pixel: MercatorCoord{
			X: 2151910,
			Y: 3118691,
		},
		Tile: MercatorCoord{
			X: 8405,
			Y: 12182,
		},
		PixelOnTile: MercatorCoord{
			X: 230,
			Y: 99,
		},
	}

	m := NewMercatorWithLatLong(lat, long, uint(zoom))

	t.Log(fmt.Sprintf("benchmark: %+v", benchmarkMercator))
	t.Log(fmt.Sprintf("calculated: %+v", m))

	if !reflect.DeepEqual(benchmarkMercator, m) {
		t.Fail()
	}

}
