# Bafin Downloader

downloads the latest financial reports based of ISIN from
the [Bafin website](https://portal.mvp.bafin.de/database/DealingsInfo/sucheForm.do) based on the emittent name and
returns the data as single and merged csv files.

## Usage

1. Download the executable
2. open the inside the folder with the executable a terminal (Windows `shift + right click` -> `open command window / terminal here`)
3. run the executable with the following command: `./BafinDownloader <input_file>`, eg. `./BafinDownloader input.txt`

> The input file should contain a list of Company ISIN, each on a new line.

 ```
DE0012359698
DE0005859123
DE0005859456
```

## Build locally

1. Clone the repository
2. Run `go mod tidy` or `make setup`
3. Run `go build .` or `make build`
