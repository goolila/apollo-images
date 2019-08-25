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
	"strconv"
	"strings"
)

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

func main() {
	var outputDir = flag.String("output", "/tmp", "output dir to save photos on")
	var missionID = flag.Int("mission", 11, "number of apollo mission" )
	var onlyHR = flag.Bool("hr", true, "download only high res photos")

	flag.Parse()

	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		log.Fatalf("dir %s does not exist", *outputDir)
	}

	missionIDstr := strconv.Itoa(*missionID)
	url := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%[1]s/images%[1]s.html", missionIDstr)

	c := colly.NewCollector()
	nDone := 0



	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		imageLink := fmt.Sprintf("https://www.hq.nasa.gov/alsj/a%s/", missionIDstr) + link
		if okToDownload(imageLink, *onlyHR) {
				dest := path.Join(*outputDir, link)
				log.Printf("downloading  %s to %s \n", imageLink, dest)
				if err := DownloadFile(imageLink, dest); err != nil {
					panic(err)
    			}

		}
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})
	fmt.Printf("downloaded %d to %s", nDone, *outputDir)

	c.Visit(url)
}