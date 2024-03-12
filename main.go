package main

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"
)

const dir = "data/"

var (
	Misc = Teal
	Info = Green
	Fata = Red
)

var (
	Red   = Color("\033[1;31m%s\033[0m")
	Green = Color("\033[1;32m%s\033[0m")
	Teal  = Color("\033[1;36m%s\033[0m")
)

func Color(colorString string) func(...interface{}) string {
	sprint := func(args ...interface{}) string {
		return fmt.Sprintf(colorString,
			fmt.Sprint(args...))
	}
	return sprint
}

func main() {
	fmt.Println(Misc("---- starting ----"))

	client := &http.Client{
		Transport: &http.Transport{},
	}
	names, err := loadList()
	if err != nil {
		fmt.Println(Fata(err))
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

	fmt.Println(Misc("---- done ----"))

}

func loadList() ([]string, error) {
	fileName := os.Args[1]
	if fileName == "" {
		return nil, fmt.Errorf("no file name provided")
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(file), "\r\n"), nil
}

func getCsv(client *http.Client, name string) {
	values := url.Values{
		"emittentName": {name},
	}
	domain := fmt.Sprintf("https://portal.mvp.bafin.de/database/DealingsInfo/sucheForm.do?meldepflichtigerName=&zeitraum=0&d-4000784-e=1&emittentButton=Suche+Emittent&%s&zeitraumVon=&emittentIsin=&6578706f7274=1&zeitraumBis=", values.Encode())
	fmt.Println(Misc("getting:"), name)
	res, err := client.Get(domain)
	if err != nil {
		fmt.Println(Fata(err))
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
		fmt.Println(Fata(err))
		return
	}
	csv = csv[:n]
	err = saveFile(dir, name, csv)
	if err != nil {
		fmt.Println(Fata(err))
	}
}

func saveFile(dir, name string, content []byte) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	// overwrite file
	filePath := path.Join(dir, name+".csv")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(Fata(err))
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		fmt.Println(Fata(err))
		return err
	}

	cwd, err := os.Executable()
	if err != nil {
		fmt.Println(Fata(err))
		return err
	}

	fmt.Println(Info("saved:"), filepath.Join(cwd, dir, name+".csv"))
	return nil
}

func mergeCsv() {
	fmt.Println(Misc("merging csv files..."))
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println(Fata(err))
	}

	var contentList []string
	header, err := readFile(dir + files[0].Name())
	if err != nil {
		fmt.Println(Fata(err))
	}
	header = header[:strings.Index(string(header), "\n")]
	contentList = append(contentList, string(header))

	for _, file := range files {
		if file.IsDir() || file.Name() == "#merged.csv" {
			continue
		}
		content, err := readFile(dir + file.Name())
		if err != nil {
			fmt.Println(Fata(err))
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

	err = saveFile(dir, "#merged", []byte(strings.Join(contentList, "\n")))
	if err != nil {
		fmt.Println(Fata(err))
	}

	fmt.Println(Info("merging files done"))
}

func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}