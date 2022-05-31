package frm

import (
	"bufio"
	"log"
	"os"
)

func OpenFile(filename string, mode os.FileMode) *os.File {
	f, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}

	return f
}

func ReadData(f *os.File) ([]string, error) {
	scanner := bufio.NewScanner(f)
	var data []string

	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		Logger.Fatal().Err(err)
		return []string{}, err
	}

	return data, nil
}

func WriteData(f *os.File, data string) error {
	if _, err := f.WriteString((data + "\n")); err != nil {
		Logger.Fatal().Err(err)
		return err
	}
	return nil
}
