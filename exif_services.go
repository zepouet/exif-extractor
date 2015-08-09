package main

import (
	//"fmt"
	"github.com/rwcarlsen/goexif/exif"
	"github.com/zepouet/exif-extractor/api"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// TODO
func avian_carrier(res chan api.ExifInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("3 :: wg.done")

	for e := range res {
		log.Println("3 :: telegraf the data : ", e)
	}
}

func converter(inChan chan string, outChan chan api.ExifInfo, wg *sync.WaitGroup) {
	defer wg.Done()
	defer log.Println("2 :: wg.done")

	for file := range inChan {

		log.Println("2 :: convert file to exif : ", file)
		f, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}

		exif, err := exif.Decode(f)
		if err != nil {
			log.Fatal(err)
		}

		info := api.ExifInfo{FileName: file}
		info.Decode(exif)

		outChan <- info
	}
}

func main() {

	t0 := time.Now()

	// halt and catch fire
	runtime.GOMAXPROCS(runtime.NumCPU())

	// worker wait group
	var wgConvert sync.WaitGroup
	var wgTelegraf sync.WaitGroup
	workers := 12
	wgConvert.Add(workers)
	wgTelegraf.Add(workers)

	// Convert FileInfo into exif struct
	filesChannel := make(chan string, workers*2)
	exifChannel := make(chan api.ExifInfo, workers)

	// number of photos
	var n api.AtomicInt

	// start in background the converter worker (from file to exif)
	for i := 0; i < workers; i++ {
		go converter(filesChannel, exifChannel, &wgConvert)
	}

	// start in background the avian carrier for dropping exif to telegraf
	for i := 0; i < workers; i++ {
		go avian_carrier(exifChannel, &wgTelegraf)
	}

	// Callback method for each resource (file or directory)
	callback := func(path string, f os.FileInfo, err error) error {

		// TODO : if file is an image (known format)
		if err == nil && !f.IsDir() &&
			(strings.HasSuffix(f.Name(), ".jpg") || strings.HasSuffix(f.Name(), ".nef") || strings.HasSuffix(f.Name(), ".png")) {
			log.Printf("1 :: Walking in the trees : %s with %d bytes\n", path, f.Size())
			n.Add(1)
			filesChannel <- path
		}
		return nil
	}

	// walk into the tree in a dark forest
	for _, p := range os.Args[1:] {
		filepath.Walk(p, callback)
	}

	// we can close the channel because filepath.Walk is blocking
	// do not run 'filepath.Walk' into a goroutine
	// else you could not know when it ends to close the channel
	log.Println("Close Files Channel")
	close(filesChannel)

	// signals the end of last worker
	log.Println("Waiting for last worker converter..")
	wgConvert.Wait()
	log.Println("Last worker converter is dead")

	// all the converter have done their job to push exif into outchan.
	// so we can close outchan
	close(exifChannel)

	// signals the end of last worker
	log.Println("Waiting for last worker emetter..")
	wgTelegraf.Wait()
	log.Println("Last worker emetter is dead")

	// success
	log.Printf("Exit %v\n", time.Since(t0))

	//os.Exit(0)

}
