package handler

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"html/template"
	"image"
	"image/jpeg"
	mosaic "mosaic_web/mosaic"
	"net/http"
	"strconv"
	"time"
)

// Handler function for fan-out with individual channels
func FanOutWithChannelHandlerFunc(w http.ResponseWriter, r *http.Request) {
	t0 := time.Now()
	// get the content from the POSTed form
	r.ParseMultipartForm(10485760) // max body in memory is 10MB
	file, _, _ := r.FormFile("image")
	defer file.Close()
	tileSize, _ := strconv.Atoi(r.FormValue("tile_size"))
	//
	//   // decode and get original image
	original, _, _ := image.Decode(file)
	bounds := original.Bounds()
	db := mosaic.CloneTilesDB()

	c1 := mosaic.CutWithChannel(original, &db, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	c2 := mosaic.CutWithChannel(original, &db, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	c3 := mosaic.CutWithChannel(original, &db, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	c4 := mosaic.CutWithChannel(original, &db, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)

	buf1 := new(bytes.Buffer)
	jpeg.Encode(buf1, original, nil)
	originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())
	t1 := time.Now()
	images := map[string]string{
		"original": originalStr,
		"part1":    <-c1,
		"part2":    <-c2,
		"part3":    <-c3,
		"part4":    <-c4,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}

	t, _ := template.ParseFiles("results_parts.html")
	t.Execute(w, images)
}
