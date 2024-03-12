# Bafin Downloader

downloads the latest financial reports from
the [Bafin website](https://portal.mvp.bafin.de/database/DealingsInfo/sucheForm.do) based on the emittent name and
returns the data as single and merged csv files.

## Usage

1. Download the executable
2. open the folder with the executable (Windows `shift + right click` -> `open command window / terminal here`)
3. run the executable with the following command: `./BafinDownloader <input_file>`, eg. `./BafinDownloader input.txt`

> The input file should contain a list of Company names, each on a new line.

 ```
Company 1
Company 2
Company 3
```

## Build locally

1. Clone the repository
2. Run `go mod tidy` or `make setup`
3. Run `go build .` or `make build`