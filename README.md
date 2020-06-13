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
$ ./apollo-images -mission 12 -sleep 500
🚀 It was so much fun: https://en.wikipedia.org/wiki/Apollo_11Visiting https://www.hq.nasa.gov/alsj/a11/images11.html
[worker 1] downloading  https://www.hq.nasa.gov/alsj/a11/a11LunOrb5076HR.jpg to /tmp/apollo-images/11/a11LunOrb5076HR.jpg
[worker 2] downloading  https://www.hq.nasa.gov/alsj/a11/a11pan1040226combHR.jpg to /tmp/apollo-images/11/a11pan1040226combHR.jpg
[worker 3] downloading  https://www.hq.nasa.gov/alsj/a11/a11pan1040226lftHR.jpg to /tmp/apollo-images/11/a11pan1040226lftHR.jpg
[worker 0] downloading  https://www.hq.nasa.gov/alsj/a11/a11pan1040226rghtHR.jpg to /tmp/apollo-images/11/a11pan1040226rghtHR.jpg
[worker 3] downloading  https://www.hq.nasa.gov/alsj/a11/a11pan1093226HR.jpg to /tmp/apollo-images/11/a11pan1093226HR.jpg
```
