## apollo image downloader

This package scrape NASA website to download images of Apollo missions.

## Build
```
git clone https://gitlab.com/goolila/apollo-images
cd apollo-images
GO111MODULE=on go build
```


## Usage
```
Usage of ./apollo-images:
  -hr
    	download only high res photos (default true)
  -mission int
    	number of apollo mission (default 11)
  -output string
    	output dir to save photos on (default "/tmp")
```

Example:
```
mkdir -p /tmp/apollo11
./apollo-images -mission 11 -output /tmp/apollo11
```
