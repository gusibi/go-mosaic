package mosaic

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"time"
)

// Handler function for fan-out with individual channels
func FanOutWithChannel(imgPath, tiles, outFile string, tileSize int) error {
	t0 := time.Now()
	// get the content from the POSTed form
	file, err := os.Open(imgPath)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	//
	// decode and get original image
	original, _, _ := image.Decode(file)
	// 调整图片大小，目标宽为tileSize 的50 倍
	newWidth := tileSize * 60
	original = Resize(original, newWidth)
	bounds := original.Bounds()
	db := CloneTilesDB()

	c1 := CutWithChannel(original, &db, tileSize, bounds.Min.X, bounds.Min.Y, bounds.Max.X/2, bounds.Max.Y/2)
	c2 := CutWithChannel(original, &db, tileSize, bounds.Max.X/2, bounds.Min.Y, bounds.Max.X, bounds.Max.Y/2)
	c3 := CutWithChannel(original, &db, tileSize, bounds.Min.X, bounds.Max.Y/2, bounds.Max.X/2, bounds.Max.Y)
	c4 := CutWithChannel(original, &db, tileSize, bounds.Max.X/2, bounds.Max.Y/2, bounds.Max.X, bounds.Max.Y)

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
	fmt.Println(images)
	return nil
}

// fan-out with individual channels //////////////////////////////////////////

func CutWithChannel(original image.Image, db *map[string][3]float64, tileSize, x1, y1, x2, y2 int) <-chan string {
	c := make(chan string)
	go func() {
		newimage := image.NewNRGBA(image.Rect(x1, y1, x2, y2))
		sp := image.Point{0, 0}

		for y := y1; y < y2; y = y + tileSize {
			for x := x1; x < x2; x = x + tileSize {
				// get the RGBA value
				r, g, b, _ := original.At(x, y).RGBA()
				color := [3]float64{float64(r), float64(g), float64(b)}
				nearest := Nearest(color, db)
				file, err := os.Open(nearest)
				if err == nil {
					img, _, err := image.Decode(file)
					if err == nil {
						// t := imaging.Resize(img, tileSize, tileSize, imaging.Lanczos)
						t := Resize(img, tileSize)
						tile := t.SubImage(t.Bounds())
						tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
						draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
					} else {
						fmt.Println("error in decoding nearest", err, nearest)
					}
				} else {
					fmt.Println("error opening file when creating mosaic:", nearest)
				}
				file.Close()
			}
		}
		buf2 := new(bytes.Buffer)
		jpeg.Encode(buf2, newimage, nil)
		c <- base64.StdEncoding.EncodeToString(buf2.Bytes())

	}()
	return c
}
