package frm

import (
	"bufio"
	"os"

)

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
