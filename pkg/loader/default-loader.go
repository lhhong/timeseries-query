package loader

import (
	"bufio"
	"encoding/csv"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
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

	batchSize := 50000
	values := make([]repository.RawData, 0, batchSize)
	for {
		values = values[:0]
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
				values = append(values, repository.RawData{
					group, line[seriesCol], 0, seq, value})
			}
		}
		bulkSave(values, repo)
	}
SaveAndExit:
	bulkSave(values, repo)
}

func bulkSave(vals []repository.RawData, repo *repository.Repository) {
	if len(vals) > 0 {
		err := repo.BulkSaveRawDataUnsafe(vals)
		if err != nil {
			log.Println(err)
		}
	}
}
