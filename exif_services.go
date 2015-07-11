package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
)

type exif struct {
	FileName     string
	Path         string
	CameraModel  string
	Focal        string
	Aperture     string
	ISO          string
	TimeShooting string
}

// TODO
func avian_carrier(res chan exif) {
	for s := range res {
		fmt.Println("avian_carrier : ", s.FileName)
	}
}

func main() {

	// halt and catch fire
	runtime.GOMAXPROCS(runtime.NumCPU())

	// worker wait group
	var wg sync.WaitGroup
	workers := 16
	wg.Add(workers)

	// Convert FileInfo into exif struct
	inChan := make(chan os.FileInfo, workers*2)
	outChan := make(chan exif, workers)

	converter := func(files chan os.FileInfo) {
		defer wg.Done()
		for file := range files {
			fmt.Println("converter : ", file.Name())
			outChan <- exif{FileName: file.Name()}
		}
	}

	// prepare the converter worker (from file to exif)
	for i := 0; i < workers; i++ {
		go converter(inChan)
	}

	// prepare the avian carrier for dropping packet to telegraf
	go avian_carrier(outChan)

	// Callback method for each resource (file or directory)
	callback := func(path string, f os.FileInfo, err error) error {
		//fmt.Printf("%s with %d bytes\n", path, f.Size())
		// TODO : if file is an image (known format)
		if err == nil {
			inChan <- f
		}
		return nil
	}

	// walk into the tree in a dark forest
	for _, p := range os.Args[1:] {
		filepath.Walk(p, callback)
	}

	// we can close the channel because filepath.Walk is blocking
	// do not run 'callback' method it into a goroutine
	// else you could not know when it ends to close the channel
	close(inChan)

	// signals the end of last worker
	wg.Wait()

	// all the converter have done their job to push exif into outchan.
	// so we can close outchan
	close(outChan)

	// success
	os.Exit(0)
}
