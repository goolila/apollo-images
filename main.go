package main

import (
	"flag"
	"fmt"
	"github.com/gocolly/colly"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

var outputDir = ""
var wg sync.WaitGroup


// check if given string is an image and should be downloaded
func okToDownload(s string, hr bool) bool {
	if strings.Contains(s, "jpg") == false {
		return false
	}
	if strings.Contains(s, "..") == true {
		return false
	}
	if hr && strings.Contains(s, "HR") == false{
		return false
	}
	return true
}


func DownloadFile(url string, filepath string) error {
    resp, err := http.Get(url)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    // Create the file
    out, err := os.Create(filepath)
    if err != nil {
        return err
    }
    defer out.Close()

    // Write the body to file
    _, err = io.Copy(out, resp.Body)
    return err
}

func worker(queue chan string, worknumber int, mID string) {
	for link := range queue {
		imageLink := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%s/", mID) + link
		dest := path.Join(outputDir, link)
		log.Printf("[worker %d] downloading  %s to %s \n", worknumber, imageLink, dest)
		if err := DownloadFile(imageLink, dest); err != nil {
			panic(err)
		}
		wg.Done()
	}
}

func validateMission(mission int) {
    switch m := mission; {
	case m < 11:
		log.Fatalf("Apollo 11 was the spaceflight that first landed humans on the Moon")
	case m > 17:
		log.Fatalf("Apollo 17 (December 7–19, 1972) was the final mission of NASA's Apollo program.\n")
	default:
		log.Printf("It was so much fun: https://en.wikipedia.org/wiki/Apollo_%d", mission)
    }
}

func main() {
	var output = flag.String("output", "/tmp", "output outputDir to save photos on")
	var mission = flag.Int("mission", 11, "number of apollo mission" )
	var onlyHR = flag.Bool("hr", true, "download only high res photos")
	flag.Parse()

	validateMission(*mission)

	outputDir = *output
	missionID := strconv.Itoa(*mission)

	if _, err := os.Stat(outputDir); os.IsNotExist(err) {
		log.Fatalf("outputDir %s does not exist", outputDir)
	}

	// queue of jobs
	q := make(chan string)
	// init workers
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(q, i, missionID)
	}

	// init colly
	c := colly.NewCollector()
	// root url to be visited
	url := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%[1]s/images%[1]s.html", missionID)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		imageLink := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%s/", missionID) + link
		if okToDownload(imageLink, *onlyHR) {
			// send link to job queue
			wg.Add(1)
			q <- link
		}
	})
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	if err := c.Visit(url); err != nil {
		log.Fatalf("error happened on visting root url %s %v+\n", url, err)
	}
	wg.Wait()
}