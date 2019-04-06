package loader

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"

	"github.com/lhhong/timeseries-query/pkg/config"
	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
)

func LoadECG(cmd *cobra.Command, conf *config.AppConfig, repo *repository.Repository) {

	dir, _ := cmd.Flags().GetString("dir")
	group, _ := cmd.Flags().GetString("groupname")

	postfix := ".csv"
	re := regexp.MustCompile("(.*)" + regexp.QuoteMeta(postfix))
	var files []string
	var paths []string
	start := time.Now()

	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, info.Name())
		paths = append(paths, path)
		return nil
	})

	series := make(map[string]bool)
	for i, f := range files {
		if match := re.FindStringSubmatch(f); len(match) > 0 {
			extractedSeries := readAndSaveECGSeries(repo, group, paths[i], match[1])
			for s := range extractedSeries {
				series[s] = true
			}
		}
	}
	saveSeries(series, group, repo)
	elapsed := time.Since(start)

	log.Println("Starting to load index")

	loadIndex(group, conf.Env, repo)
	log.Printf("Storing query to database took %s", elapsed)
}

func readAndSaveECGSeries(repo *repository.Repository, group string, path string, name string) map[string]bool {
	log.Println("Saving", name)
	csvFile, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()

	reader := csv.NewReader(bufio.NewReader(csvFile))
	_, err = reader.Read()
	if err == io.EOF {
		return nil
	}
	_, err = reader.Read()
	if err == io.EOF {
		return nil
	}

	batchSize := 10000
	count := 0
	ind := 0
	series := make(map[string]bool)
	var values []repository.RawData
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
			panic(err)
		} else {
			x, err := strconv.ParseInt(line[0], 10, 64)
			if err != nil {
				log.Println(err)
				continue
			}

			seriesPartition := name + "_" + fmt.Sprintf("%03d", ind)
			if exists := series[seriesPartition]; !exists {
				series[seriesPartition] = true
			}
			v, err := strconv.ParseFloat(line[1], 64)
			if err != nil {
				log.Println(err)
				continue
			}
			values = append(values, repository.RawData{
				Groupname: group,
				Series:    seriesPartition,
				Seq:       x,
				Ind:       count,
				Value:     v,
			})
		}
		count++
		if count >= batchSize {
			bulkSave(values, repo)
			values = values[:0]
			count = 0
			ind++
		}
	}
	bulkSave(values, repo)
	return series
}
