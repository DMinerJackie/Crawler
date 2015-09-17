package main

import (
	"encoding/csv"
	"fmt"
	"os"
	//"net/url"
)

func ExportToCSV(startHost string, visited map[string]bool) {
	csvfile, err := os.Create(startHost + ".csv")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer csvfile.Close()

	writer := csv.NewWriter(csvfile)
	for columns, _ := range visited {
		temp := []string{columns}
		err = writer.Write(temp)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
	}
	writer.Flush()
}
