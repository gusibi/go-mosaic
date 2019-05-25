package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	handler "mosaic_web/handler"
	mosaic "mosaic_web/mosaic"
	"net/http"
	"runtime"
)

func uploadHandlerFunc(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles("upload.html")
	t.Execute(w, len(mosaic.TILESDB))
}

func fetchHandlerFunc(w http.ResponseWriter, r *http.Request) {
	files, _ := ioutil.ReadDir("tiles")
	t, _ := template.ParseFiles("fetch.html")
	t.Execute(w, len(files))
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	fmt.Println("Starting mosaic server ...")
	mux := http.NewServeMux()
	files := http.FileServer(http.Dir("public"))
	mux.Handle("/static/", http.StripPrefix("/static/", files))

	mux.HandleFunc("/", uploadHandlerFunc)
	mux.HandleFunc("/fetch", fetchHandlerFunc)
	mux.HandleFunc("/reload", handler.ReloadTilesDBHandlerFunc)
	mux.HandleFunc("/fetch_tiles", handler.FetchTilesHandlerFunc)
	mux.HandleFunc("/mosaic_no_concurrency", handler.NoConcurrencyHandlerFunc)
	mux.HandleFunc("/mosaic_fanout_channel", handler.FanOutWithChannelHandlerFunc)
	mux.HandleFunc("/mosaic_fanout_fanin", handler.FanOutFanInHandlerFunc)

	server := &http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}
	mosaic.TILESDB = mosaic.TilesDB()
	fmt.Println("Mosaic server started.", mosaic.TILESDB)
	server.ListenAndServe()

}
