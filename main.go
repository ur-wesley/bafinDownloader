package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/logrusorgru/aurora/v4"
)

const dir = "data/"

func main() {
	fmt.Println(aurora.Blue("---- starting ----"))

	client := &http.Client{
		Transport: &http.Transport{},
	}
	names, err := LoadList()
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
		return
	}

	var wg sync.WaitGroup
	for _, name := range names {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			getCsv(client, name)
		}(name)
	}
	wg.Wait()

	mergeCsv()

	fmt.Println(aurora.Blue("---- done ----"))
}

func getCsv(client *http.Client, name string) {
	values := url.Values{
		"emittentName": {name},
	}
	domain := fmt.Sprintf("https://portal.mvp.bafin.de/database/DealingsInfo/sucheForm.do?meldepflichtigerName=&zeitraum=0&d-4000784-e=1&emittentButton=Suche+Emittent&%s&zeitraumVon=&emittentIsin=&6578706f7274=1&zeitraumBis=", values.Encode())
	fmt.Println(aurora.Green("getting:"), name)
	res, err := client.Get(domain)
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
		return
	}
	if res.StatusCode != 200 {
		fmt.Println("error: ", res.StatusCode)
		return
	}
	defer res.Body.Close()

	csv := make([]byte, 1048576)
	n, err := res.Body.Read(csv)
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
		return
	}
	csv = csv[:n]
	err = SaveFile(dir, name, csv)
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
	}
}

func mergeCsv() {
	fmt.Println(aurora.Green("merging csv files..."))
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
	}

	var contentList []string
	header, err := ReadFile(dir + files[0].Name())
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
	}
	header = header[:strings.Index(string(header), "\n")]
	contentList = append(contentList, string(header))

	for _, file := range files {
		if file.IsDir() || file.Name() == "#merged.csv" {
			continue
		}
		content, err := ReadFile(dir + file.Name())
		if err != nil {
			fmt.Println(aurora.White(err).BgRed())
		}
		lines := strings.Split(string(content), "\n")
		for i := 0; i < len(lines); i++ {
			if lines[i] == "" {
				lines = append(lines[:i], lines[i+1:]...)
				i--
			}
		}
		lines = lines[1:]
		contentList = append(contentList, strings.Join(lines, "\n"))
	}

	err = SaveFile(dir, "#merged", []byte(strings.Join(contentList, "\n")))
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
	}

	fmt.Println(aurora.Green("merging files done"))
}