package mercator

import (
	"math"
)

const (
	TILE_SIZE = 256
)

type MercatorCoord struct {
	X int64
	Y int64
}

type MercatorWorld struct {
	X float64
	Y float64
}

type Mercator struct {
	Latitude    float64
	Longitude   float64
	Zoom        uint
	World       MercatorWorld
	Pixel       MercatorCoord
	Tile        MercatorCoord
	PixelOnTile MercatorCoord
}

func (m *Mercator) Deg2num() (world MercatorWorld, pixel MercatorCoord,
	tile MercatorCoord, pixelOnTile MercatorCoord) {

	var siny = math.Sin(m.Latitude * math.Pi / 180)

	// taken from here
	// https://developers.google.com/maps/documentation/javascript/examples/map-coordinates
	// Truncating to 0.9999 effectively limits latitude to 89.189. This is
	// about a third of a tile past the edge of the world tile.
	siny = math.Min(math.Max(siny, -0.9999), 0.9999)

	world = MercatorWorld{
		X: TILE_SIZE * (0.5 + m.Longitude/360),
		Y: TILE_SIZE * (0.5 - math.Log((1+siny)/(1-siny))/(4*math.Pi)),
	}

	scale := 1 << m.Zoom

	pixel = MercatorCoord{
		X: int64(math.Floor(world.X * float64(scale))),
		Y: int64(math.Floor(world.Y * float64(scale))),
	}

	tile = MercatorCoord{
		X: int64(math.Floor(world.X * float64(scale) / TILE_SIZE)),
		Y: int64(math.Floor(world.Y * float64(scale) / TILE_SIZE)),
	}

	pixelOnTile = MercatorCoord{
		X: int64(math.Mod(world.X*float64(scale), float64(TILE_SIZE))),
		Y: int64(math.Mod(world.Y*float64(scale), float64(TILE_SIZE))),
	}
	return
}

func NewMercatorWithLatLong(lat float64, long float64, z uint) (m *Mercator) {
	m = new(Mercator)

	m.Latitude = lat
	m.Longitude = long
	m.Zoom = z
	m.World, m.Pixel, m.Tile, m.PixelOnTile = m.Deg2num()
	return
}
