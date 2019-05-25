package main

import (
	"fmt"
	mosaic "mosaic_web/mosaic"
)

func main() {
	mosaic.TILESDB = mosaic.TilesDB()
	fmt.Println(mosaic.TILESDB)
	fmt.Println("Starting mosaic generater ...")
	// mosaic.NoConcurrency("./cat_l.jpg", "./tiles", "./mosaic-cat-l-50-60-pro.jpg", 50)
	mosaic.FanOutWithChannel("./cat_l.jpg", "./tiles", "./mosaic-cat-l-50-60-fan.jpg", 50)
}
