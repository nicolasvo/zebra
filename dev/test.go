package main

import (
	"encoding/csv"
	"log"
	"os"
)

func contains(a string, list [][]string) int {
	for i, b := range list {
		if b[0] == a {
			return i
		}
	}
	return 0
}

func setUserLanguage(user string, language string, dataFile string) {
	f, err := os.OpenFile(dataFile, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}

	r := csv.NewReader(f)
	data, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}

	res := contains(m.Sender.ID, data)
	if res != 0 {
		data[res][1] = language
	} else {
		data = append(data, []string{user, language})
	}
	f.Truncate(0)
	f.Seek(0, 0)
	w := csv.NewWriter(f)
	w.WriteAll(data)
}

func main() {
}
