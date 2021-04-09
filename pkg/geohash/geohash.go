package geohash

import (
	"bytes"
)

const (
	MAX_LAT=90
	MIN_LAT=-90
	MAX_LONG=180
	MIN_LONG=-180
	
)

var BASE32_VALS = []byte("0123456789bcdefghjkmnpqrstuvwxyz")


type geohash interface {
	encode(lat, lon int)
}

type Geom struct {
	lat float64
	long float64
}

// Encode a latitude and longitude into a geohash
// Set precision
func Encode(lat,lng float64,_precision ...int) string {
	var precision int
	if len(_precision) == 0 {
		precision = 12
	}
	precision = _precision[0]


	var mid float64 = 0
	bits := 0
	bits_total := 0
	hash_count := 0
	lat_bnd := []float64{-90, 90}
	lng_bnd := []float64{-180, 180}
	var chars bytes.Buffer
	for chars.Len() < precision {
		if (bits_total % 2 == 0) {
			mid = (lng_bnd[0] + lng_bnd[1])/2
			if (lng > mid) {
			hash_count = (hash_count << 1) + 1;
        	lng_bnd[0] = mid;
			} else {
			hash_count = (hash_count << 1) + 0;
        	lng_bnd[1] = mid;
			}
		} else {
			mid = (lat_bnd[0] + lat_bnd[1])/2
			if (lat > mid) {
			hash_count = (hash_count << 1) + 1;
        	lat_bnd[0] = mid;
			} else {
			hash_count = (hash_count << 1) + 0;
        	lat_bnd[1] = mid;
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
	lat_bnd := []float64{-90, 90}
	lng_bnd := []float64{-180, 180}
	var hash_val byte
	bits_total := 0
	var mid float64 = 0

	for i := 0; i < len(hashstring); i++ {
		curr := []byte{hashstring[i]}
		hash_val = byte(bytes.Index(BASE32_VALS,curr))
		for bits := 4; bits >= 0; bits-- {
			bit := (int(hash_val) >> bits) & 1;
			
			if (bits_total % 2 == 0) {
				mid = (lng_bnd[0] + lng_bnd[1]) / 2;
				if (bit == 1) {
				lng_bnd[0] = mid;
				} else {
				lng_bnd[1] = mid;
				}
			} else {
				mid = (lat_bnd[0] + lat_bnd[1])/2
				if (bit == 1) {
				lat_bnd[0] = mid;
				} else {
				lat_bnd[1] = mid;
				}
		}
		bits_total++
	}
}

return []float64{lat_bnd[0],lng_bnd[0],lat_bnd[1],lng_bnd[1]}
}
//Decode a hash string into corresponding latitude and longitude
func Decode(hashstring string) (position []float64)  {
	bbox := Decode_Bbox(hashstring)
	lat := (bbox[0] + bbox[2]) / 2;
  	lon := (bbox[1] + bbox[3]) / 2;
	
	position = []float64{lat,lon}
	return
}

func FromPolygon(){}