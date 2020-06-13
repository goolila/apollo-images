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
	"time"
)

var (
	wg sync.WaitGroup
	// flags
	outputDir string
	missionId int
	sleep     int
	onlyHR    bool
)

// check if given string is an image and should be downloaded
func okToDownload(s string, hr bool) bool {
	if strings.Contains(s, "jpg") == false {
		return false
	}
	if strings.Contains(s, "..") == true {
		return false
	}
	if hr && strings.Contains(s, "HR") == false {
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

func worker(queue chan string, id int, mID int) {
	for link := range queue {
		imageLink := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%d/", mID) + link
		dest := path.Join(missionDir(mID), link)
		fmt.Printf("[worker %d] downloading  %s to %s \n", id, imageLink, dest)
		time.Sleep(time.Duration(sleep) * time.Millisecond)
		if err := DownloadFile(imageLink, dest); err != nil {
			log.Fatalf("error happened on downloading %s. %+v", imageLink, err)
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
		fmt.Printf("🚀 It was so much fun: https://en.wikipedia.org/wiki/Apollo_%d\n", mission)
	}
}

func ensureDir(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		log.Printf("outputDir %s does not exist. it will be created.\n", dir)
		err = os.MkdirAll(dir, 755)
		if err != nil {
			return err
		}
	}
	return nil
}

func missionDir(id int) string {
	return path.Join(outputDir, strconv.Itoa(id))
}

func main() {
	flag.StringVar(&outputDir, "output", "/tmp/apollo-images", "output outputDir to save photos on")
	flag.IntVar(&missionId, "mission", 11, "number of apollo mission")
	flag.BoolVar(&onlyHR, "hr", true, "download only high res photos")
	flag.IntVar(&sleep, "sleep", 250, "ms to sleep before queueing new url")
	flag.Parse()

	validateMission(missionId)

	mDir := missionDir(missionId)
	if err := ensureDir(mDir); err != nil {
		log.Fatalf("output dir %s does not exists and could not be created: %s", mDir, err)
	}

	// queue of jobs
	q := make(chan string)
	// init workers
	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(q, i, missionId)
	}

	// init colly
	c := colly.NewCollector()
	// root url to be visited
	url := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%[1]d/images%[1]d.html", missionId)

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		imageLink := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%d/", missionId) + link
		if okToDownload(imageLink, onlyHR) {
			// send link to job queue
			time.Sleep(time.Duration(sleep) * time.Millisecond)
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
