package main

import (
	"fmt"
	"github.com/logrusorgru/aurora/v4"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func LoadList() ([]string, error) {
	if len(os.Args) < 2 {
		return nil, fmt.Errorf("no file name provided")
	}
	fileName := os.Args[1]
	if fileName == "" {
		return nil, fmt.Errorf("no file name provided")
	}
	file, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(file), "\r\n")
	for i := 0; i < len(list); i++ {
		if list[i] == "" {
			list = append(list[:i], list[i+1:]...)
			i--
		}
	}
	return list, nil
}

func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func SaveFile(dir, name string, content []byte) error {
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		err := os.Mkdir(dir, 0755)
		if err != nil {
			return err
		}
	}
	filePath := path.Join(dir, name+".csv")
	file, err := os.Create(filePath)
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
		return err
	}
	defer file.Close()
	_, err = file.Write(content)
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Println(aurora.White(err).BgRed())
		return err
	}

	fmt.Println(aurora.Black("saved:").BgGreen(), aurora.White(filepath.Join(cwd, dir, name+".csv")))
	return nil
}