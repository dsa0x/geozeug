package geohash

import (
	"bytes"
	"fmt"
	"math"
)

const (
	MAX_LAT  = 90
	MIN_LAT  = -90
	MAX_LONG = 180
	MIN_LONG = -180
)

var BASE32_VALS = []byte("0123456789bcdefghjkmnpqrstuvwxyz")

type geohash interface {
	encode(lat, lon int)
}

type Geom struct {
	lat  float64
	long float64
}

// Encode a latitude and longitude into a geohash
// Set precision
func Encode(lat, lng float64, _precision ...int) string {
	var precision int
	if len(_precision) <= 0 {
		precision = 12
	} else {
		precision = _precision[0]
	}

	var mid float64 = 0
	bits := 0
	bits_total := 0
	hash_count := 0
	lat_bnd := []float64{MIN_LAT, MAX_LAT}
	lng_bnd := []float64{MIN_LONG, MAX_LONG}
	var chars bytes.Buffer
	for chars.Len() < precision {
		if bits_total%2 == 0 {
			mid = (lng_bnd[0] + lng_bnd[1]) / 2
			if lng > mid {
				hash_count = (hash_count << 1) + 1
				lng_bnd[0] = mid
			} else {
				hash_count = (hash_count << 1) + 0
				lng_bnd[1] = mid

			}
		} else {
			mid = (lat_bnd[0] + lat_bnd[1]) / 2
			if lat > mid {
				hash_count = (hash_count << 1) + 1
				lat_bnd[0] = mid
			} else {
				hash_count = (hash_count << 1) + 0
				lat_bnd[1] = mid
			}
		}
		bits_total++
		if bits++; bits == 5 {
			chars.WriteByte(BASE32_VALS[hash_count])
			bits = 0
			hash_count = 0

		}
	}
	return chars.String()
}

//Decode_Bbox decodes a hashstring into its corresponding bounding box
func Decode_Bbox(hashstring string) []float64 {
	lat_bnd := []float64{MIN_LAT, MAX_LAT}
	lng_bnd := []float64{MIN_LONG, MAX_LONG}
	var hash_val byte
	bits_total := 0
	var mid float64 = 0

	for i := 0; i < len(hashstring); i++ {
		curr := []byte{hashstring[i]}
		hash_val = byte(bytes.Index(BASE32_VALS, curr))
		for bits := 4; bits >= 0; bits-- {
			bit := (int(hash_val) >> bits) & 1

			if bits_total%2 == 0 {
				mid = (lng_bnd[0] + lng_bnd[1]) / 2
				if bit == 1 {
					lng_bnd[0] = mid
				} else {
					lng_bnd[1] = mid
				}
			} else {
				mid = (lat_bnd[0] + lat_bnd[1]) / 2
				if bit == 1 {
					lat_bnd[0] = mid
				} else {
					lat_bnd[1] = mid
				}
			}
			bits_total++
		}
	}

	return []float64{lat_bnd[0], lng_bnd[0], lat_bnd[1], lng_bnd[1]}
}

//Decode a hash string into corresponding latitude and longitude
func Decode(hashstring string) (position, errs []float64) {
	bbox := Decode_Bbox(hashstring)
	lat := (bbox[0] + bbox[2]) / 2
	lon := (bbox[1] + bbox[3]) / 2
	latErr := bbox[2] - lat
	lonErr := bbox[3] - lon
	position = []float64{lat, lon}
	errs = []float64{latErr, lonErr}
	return
}

//Neighbor finds the neighbor of a geohash in a certain direction.
// direction example: {1,0} for "north", {0,1} for "east"
func Neighbor(hashstring string, direction []float64) string {
	position, errs := Decode(hashstring)
	nLat := position[0] + direction[0]*errs[0]*2
	nLon := position[0] + direction[1]*errs[0]*2
	return Encode(validateLat(nLat), validateLon(nLon), 12)
}

// DirtoArr helper method to convert string direction like "north" to float slice
func DirtoArr(direction string) []float64 {
	switch direction {
	case "north":
		return []float64{1, 0}
	case "east":
		return []float64{0, 1}
	case "west":
		return []float64{-1, 0}
	case "south":
		return []float64{-1, 0}
	default:
		return []float64{1, 0}
	}
}

//Neighbors return hashstring of all neighbors
func Neigbors(hashstring string) []string {
	hashLen := len(hashstring)

	position, errs := Decode(hashstring)
	lat := position[0]
	lon := position[1]
	latErr := errs[0] * 2
	lonErr := errs[1] * 2

	hashList := []string{}
	for i := -1; i < 2; i++ {
		for j := -1; j < 2; j++ {
			if i == 0 && j == 0 {
				continue
			}
			nLat := validateLat(lat + float64(i)*latErr)
			nLon := validateLon(lon + float64(j)*lonErr)
			nHash := Encode(nLat, nLon, hashLen)
			hashList = append(hashList, nHash)

		}
	}

	return hashList
}

func validateLon(lon float64) float64 {
	if lon > MAX_LONG {
		return math.Mod((MIN_LONG)+lon, MAX_LONG)
	} else if lon < MIN_LONG {
		return math.Mod(MAX_LONG+lon, MAX_LONG)
	}
	return lon
}
func validateLat(lat float64) float64 {
	if lat > MAX_LAT {
		return MAX_LAT
	}
	if lat < MIN_LAT {
		return MIN_LAT
	}
	return lat
}

//TODO
//Bbox return all hashstrings within the bounding box
func bbox(lat, lon []float64, precision int) []string {
	if precision <= 0 {
		precision = 12
	}
	hashSouthWest := Encode(lat[0], lon[0], precision)
	hashNorthEast := Encode(lat[1], lon[1], precision)
	_, errs := Decode(hashSouthWest)

	boxSouthWest := Decode_Bbox(hashSouthWest)
	boxNorthEast := Decode_Bbox(hashNorthEast)
	fmt.Println(errs[0])
	latStep := math.Round((boxNorthEast[0] - boxSouthWest[0]) / errs[0] * 2)
	lonStep := math.Round((boxNorthEast[1] - boxSouthWest[1]) / errs[1] * 2)
	hashList := []string{}
	fmt.Println(lat, latStep)
	for lat := 0; float64(lat) <= latStep; lat++ {
		for lon := 0; float64(lon) <= lonStep; lon++ {
			pos := []float64{float64(lat), float64(lon)}
			hashList = append(hashList, Neighbor(hashSouthWest, pos))
		}
	}

	return hashList
}

func FromPolygon() {}
