package loader

import (
	"bufio"
	"encoding/csv"
	"github.com/lhhong/timeseries-query/backend/pkg/repository"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"
	"strconv"
)

func LoadData(cmd *cobra.Command) {

	var err error

	group, _ := cmd.Flags().GetString("groupname")
	data, _ := cmd.Flags().GetString("datafile")

	seriesCol, _ := cmd.Flags().GetInt("series")
	dateCol, _ := cmd.Flags().GetInt("date")
	valCol, _ := cmd.Flags().GetInt("value")

	csvFile, err := os.Open(data)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	_, err = reader.Read()
	if err == io.EOF {
		return
	}
	batchSize := 10000
	valueArgs := make([]interface{}, 0, batchSize*4)
	for {
		for i := 0; i < batchSize; i++ {
			valueArgs = valueArgs[:0]
			line, err := reader.Read()
			if err == io.EOF {
				goto SaveAndExit
			} else if err != nil {
				log.Fatal(err)
				panic(err)
			} else {
				value, err := strconv.ParseFloat(line[valCol], 64)
				if err != nil {
					log.Println(err)
					continue
				}
				valueArgs = append(valueArgs, group)
				valueArgs = append(valueArgs, line[seriesCol])
				valueArgs = append(valueArgs, line[dateCol])
				valueArgs = append(valueArgs, value)
			}
			bulkSave(valueArgs)
		}
	}
SaveAndExit:
	bulkSave(valueArgs)
}

func bulkSave(vals []interface{}) {
	if len(vals) > 0 {
		err := repository.Repo.BulkSaveRawData(vals)
		if err != nil {
			log.Println(err)
		}
	}
}
