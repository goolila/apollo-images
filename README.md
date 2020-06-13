## apollo image downloader

CLI to scrape NASA website to download images of Apollo missions.

## Build

Make sure you have go installed (tested with go version `1.14`).

```
git clone https://gitlab.com/goolila/apollo-images
cd apollo-images
go get -v -t -d ./...
go build
```

## Usage
```
Usage of ./apollo-images:
  -hr
        download only high res photos (default true)
  -mission int
        number of apollo mission (default 11)
  -output string
        output outputDir to save photos on (default "/tmp/apollo-images")
  -sleep int
        ms to sleep before queueing new url (default 250)
```

Example:
```
$ ./apollo-images -mission 12 -sleep 750
🚀 It was so much fun: https://en.wikipedia.org/wiki/Apollo_12
Visiting https://www.hq.nasa.gov/alsj/a12/images12.html
[worker 1] downloading  https://www.hq.nasa.gov/alsj/a12/ap12-S69-57076HR.jpg to /tmp/apollo-images/12/ap12-S69-57076HR.jpg
[worker 0] downloading  https://www.hq.nasa.gov/alsj/a12/a12-lam7HR.jpg to /tmp/apollo-images/12/a12-lam7HR.jpg
[worker 2] downloading  https://www.hq.nasa.gov/alsj/a12/a12-lse76gHRsite4actual.jpg to /tmp/apollo-images/12/a12-lse76gHRsite4actual.jpg
[worker 3] downloading  https://www.hq.nasa.gov/alsj/a12/a12-lse730HR.jpg to /tmp/apollo-images/12/a12-lse730HR.jpg
```
