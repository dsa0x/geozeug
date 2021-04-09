package main

import (
	"fmt"

	"github.com/dsa0x/geozeug/pkg/geohash"
)


func main()  {
	fmt.Println(geohash.Encode(52,34,6))
	fmt.Println(geohash.Decode("uc0pvd"))
}