package mosaic

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"time"

	"github.com/disintegration/imaging"
)

func NoConcurrency(imgPath, tiles, outFile string, tileSize int) error {
	t0 := time.Now()
	// get the content from the POSTed form
	// r.ParseMultipartForm(10485760) // max body in memory is 10MB
	file, err := os.Open(imgPath)
	defer file.Close()
	if err != nil {
		log.Println(err)
		return err
	}
	// tile size
	// decode and get original image
	original, _, _ := image.Decode(file)
	// 调整图片大小，目标宽为tileSize 的50 倍
	newWidth := tileSize * 60
	original = ResizePro(original, newWidth)
	bounds := original.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	fmt.Println(width, height)
	// create a new image for the mosaic
	newimage := image.NewNRGBA(image.Rect(bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y))
	// build up the tiles database
	db := CloneTilesDB()
	// source point for each tile, which starts with 0, 0 of each tile
	sp := image.Point{0, 0}
	for y := bounds.Min.Y; y < bounds.Max.Y; y = y + tileSize {
		for x := bounds.Min.X; x < bounds.Max.X; x = x + tileSize {
			// use the top left most pixel color in each rectangle for the average color
			r, g, b, _ := original.At(x, y).RGBA()
			color := [3]float64{float64(r), float64(g), float64(b)}
			// get the closest tile from the tiles DB
			nearest := Nearest(color, &db)
			file, err := os.Open(nearest)
			if err == nil {
				img, _, err := image.Decode(file)
				if err == nil {
					// resize the tile to the correct size and the image
					t := imaging.Resize(img, tileSize, tileSize, imaging.Lanczos)
					tile := t.SubImage(t.Bounds())
					tileBounds := image.Rect(x, y, x+tileSize, y+tileSize)
					// draw the tile into the mosaic
					draw.Draw(newimage, tileBounds, tile, sp, draw.Src)
				} else {
					log.Println("error in decoding nearest", err, nearest)
					return err
				}
			} else {
				fmt.Println("error opening file when creating mosaic:", nearest)
			}
			file.Close()
		}
	}

	// buf1 := new(bytes.Buffer)
	// jpeg.Encode(buf1, original, nil)
	// originalStr := base64.StdEncoding.EncodeToString(buf1.Bytes())

	outfile, err := os.Create(outFile)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	var opt jpeg.Options
	opt.Quality = 100
	jpeg.Encode(outfile, newimage, &opt)
	// mosaic := base64.StdEncoding.EncodeToString(buf2.Bytes())
	t1 := time.Now()
	images := map[string]string{
		// "original": originalStr,
		"mosaic":   outFile,
		"duration": fmt.Sprintf("%v ", t1.Sub(t0)),
	}
	fmt.Printf("succeed images: %+v", images)
	// t, _ := template.ParseFiles("results.html")
	// t.Execute(w, images)
	return nil
}
