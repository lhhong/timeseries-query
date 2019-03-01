package loader

import (
	"bufio"
	"encoding/csv"
	"github.com/lhhong/timeseries-query/pkg/config"
	"github.com/lhhong/timeseries-query/pkg/sectionindex"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
)

// LoadData Loads data from commands
func LoadData(cmd *cobra.Command, conf *config.AppConfig, repo *repository.Repository) {

	indexOnly, _ := cmd.Flags().GetBool("index-only")
	if indexOnly {
		group, _ := cmd.Flags().GetString("groupname")
		loadIndex(group, conf.Env, repo)
		return
	}
	swift, _ := cmd.Flags().GetBool("swift-data")
	if swift {
		LoadSwift(cmd, conf, repo)
		return
	}

	group, _ := cmd.Flags().GetString("groupname")
	data, _ := cmd.Flags().GetString("datafile")

	seriesCol, _ := cmd.Flags().GetInt("series")
	dateCol, _ := cmd.Flags().GetInt("date")
	valCol, _ := cmd.Flags().GetInt("value")

	readCsvAndSave(repo, group, data, seriesCol, dateCol, valCol)
	log.Println("Starting to load index")

	loadIndex(group, conf.Env, repo)
}

func readCsvAndSave(repo *repository.Repository, group string, data string, seriesCol int, dateCol int, valCol int) {
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
	series := make(map[string]bool)
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
					Groupname: group,
					Series:    line[seriesCol],
					Smooth:    0,
					Seq:       seq,
					Value:     value,
				})
				// For storing series info
				series[line[seriesCol]] = true
			}
		}
		bulkSave(values, repo)
	}
SaveAndExit:
	bulkSave(values, repo)
	saveSeries(series, group, repo)
}

func loadIndex(group string, env string, repo *repository.Repository) {
	ind := sectionindex.InitDefaultIndex()
	CalcAndSaveIndexDetails(repo, ind, env, group)
}

func saveSeries(series map[string]bool, group string, repo *repository.Repository) {
	var seriesInfo []*repository.SeriesInfo
	for s := range series {
		seriesInfo = append(seriesInfo, &repository.SeriesInfo{
			Groupname: group,
			Series:    s,
			Nsmooth:   0,
			Type:      "date",
		})
	}
	err := repo.BulkSaveSeriesInfoUnsafe(seriesInfo)
	if err != nil {
		log.Println(err)
	}
}

func bulkSave(vals []repository.RawData, repo *repository.Repository) {
	if len(vals) > 0 {
		err := repo.BulkSaveRawDataUnsafe(vals)
		if err != nil {
			log.Println(err)
		}
	}
}
