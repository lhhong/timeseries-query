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
	"time"
)

// LoadData Loads data from commands
func LoadData(cmd *cobra.Command, repo *repository.Repository) {

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
	valueArgs := make([]interface{}, 0, batchSize*repository.NRawDataEle)
	for {
		valueArgs = valueArgs[:0]
		for i := 0; i < batchSize; i++ {
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
				t, err := time.Parse("2006-01-02", line[dateCol])
				if err != nil {
					log.Println(err)
					continue
				}
				seq := t.Unix()
				valueArgs = append(valueArgs, group)
				valueArgs = append(valueArgs, line[seriesCol])
				valueArgs = append(valueArgs, 0) // Smooth iteration
				valueArgs = append(valueArgs, seq)
				valueArgs = append(valueArgs, value)
			}
		}
		bulkSave(valueArgs, repo)
	}
SaveAndExit:
	bulkSave(valueArgs, repo)
}

func bulkSave(vals []interface{}, repo *repository.Repository) {
	if len(vals) > 0 {
		err := repo.BulkSaveRawData(vals)
		if err != nil {
			log.Println(err)
		}
	}
}
