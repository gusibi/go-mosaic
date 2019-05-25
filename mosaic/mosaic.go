package mosaic

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"io/ioutil"
	"math"
	"mime"
	"os"

	"github.com/disintegration/imaging"
)

// Guess image format from gif/jpeg/png/webp
func guessImageFormat(r io.Reader) (format string, err error) {
	_, format, err = image.DecodeConfig(r)
	return
}

// Guess image mime types from gif/jpeg/png/webp
func guessImageMimeTypes(r io.Reader) string {
	format, _ := guessImageFormat(r)
	if format == "" {
		return ""
	}
	return mime.TypeByExtension("." + format)
}

// TILESDB global tiles database
var TILESDB map[string][3]float64

// CloneTilesDB clone db
func CloneTilesDB() map[string][3]float64 {
	db := make(map[string][3]float64)
	for k, v := range TILESDB {
		db[k] = v
	}
	return db
}

// TilesDB populate a tiles database in memory
func TilesDB() map[string][3]float64 {
	fmt.Println("Start populating tiles db ...")
	db := make(map[string][3]float64)
	files, _ := ioutil.ReadDir("tiles")
	for _, f := range files {
		name := "tiles/" + f.Name()
		file, err := os.Open(name)
		if err == nil {
			// t := guessImageMimeTypes(file)
			// fmt.Println("image type", t, file)
			img, _, err := image.Decode(file)
			if err == nil {
				db[name] = averageColor(img)
			} else {
				fmt.Println("error in populating tiles db:", err, name)
			}
		} else {
			fmt.Println("cannot open file", name, "when populating tiles db:", err)
		}
		file.Close()
	}
	fmt.Println("Finished populating tiles db.")
	return db
}

// NearestCache nearest file cache
var NearestCache = make(map[[3]float64]string)

// Nearest find the nearest matching image
func Nearest(target [3]float64, db *map[string][3]float64) string {
	var filename string
	smallest := 1000000.0
	// fmt.Println(len(*db))
	// filename, ok := NearestCache[target]
	ok := false
	if ok {
		// fmt.Println("ok", filename)
		return filename
	} else {
		for k, v := range *db {
			dist := distance(target, v)
			if dist < smallest {
				filename, smallest = k, dist
			}
		}
		// delete(*db, filename)
		NearestCache[target] = filename
		return filename
	}
}

// find the Eucleadian distance between 2 points
func distance(p1 [3]float64, p2 [3]float64) float64 {
	return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

// find the square
func sq(n float64) float64 {
	return n * n
}

// ResizePro 使用imaging 调整图片大小
func ResizePro(img image.Image, newWidth int) *image.NRGBA {
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Println(width, height)
	// create a new image for the mosaic
	newHeight := float64(height) * (float64(newWidth) / float64(width))
	return imaging.Resize(img, newWidth, int(newHeight), imaging.Lanczos)
}

// Resize an image by its ratio e.g. ratio 2 means reduce the size by 1/2, 10 means reduce the size by 1/10
// 调整图片大小
func Resize(in image.Image, newWidth int) *image.NRGBA {
	bounds := in.Bounds()
	width := bounds.Max.X - bounds.Min.X
	ratio := float64(width) / float64(newWidth)
	minX, maxX := float64(bounds.Min.X), float64(bounds.Max.X)
	minY, maxY := float64(bounds.Min.Y), float64(bounds.Max.Y)
	out := image.NewNRGBA(image.Rect(int(minX/ratio), int(minY/ratio), int(maxX/ratio), int(maxY/ratio)))
	for y, j := int(minY), int(minY); y < int(maxY); y, j = y+int(ratio), j+1 {
		for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+int(ratio), i+1 {
			r, g, b, a := in.At(x, y).RGBA()
			out.SetNRGBA(i, j, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
		}
	}
	return out
}

// find the average color of the picture
func averageColor(img image.Image) [3]float64 {
	bounds := img.Bounds()
	r, g, b := 0.0, 0.0, 0.0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r1, g1, b1, _ := img.At(x, y).RGBA()
			r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
		}
	}
	totalPixels := float64(bounds.Max.X * bounds.Max.Y)
	return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}
