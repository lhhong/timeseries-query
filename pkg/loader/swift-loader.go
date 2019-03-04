package loader

import (
	"log"
	"os"
	"path/filepath"
	"regexp"

	fits "github.com/astrogo/fitsio"
	"github.com/lhhong/timeseries-query/pkg/config"
	"github.com/lhhong/timeseries-query/pkg/repository"
	"github.com/spf13/cobra"
)

func LoadSwift(cmd *cobra.Command, conf *config.AppConfig, repo *repository.Repository) {

	dir, _ := cmd.Flags().GetString("dir")
	group, _ := cmd.Flags().GetString("groupname")

	prefix := "BAT_58m_snapshot_SWIFT_"
	postfix := ".lc"
	re := regexp.MustCompile(regexp.QuoteMeta(prefix) + "(.*)" + regexp.QuoteMeta(postfix))

	var files []string
	var paths []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		files = append(files, info.Name())
		paths = append(paths, path)
		return nil
	})

	series := make(map[string]bool)
	for i, f := range files {
		if match := re.FindStringSubmatch(f); len(match) > 0 {
			readAndSaveSwiftSeries(repo, group, paths[i], match[1])
			series[match[1]] = true
		}
	}
	saveSeries(series, group, repo)
	log.Println("Starting to load index")

	loadIndex(group, conf.Env, repo)
}

func readAndSaveSwiftSeries(repo *repository.Repository, group string, path string, name string) {

	var values []repository.RawData

	log.Printf("Group: %s, Series: %s\n", group, name)
	r, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer r.Close()
	f, err := fits.Open(r)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// get the second HDU
	table := f.HDU(1).(*fits.Table)

	nrows := table.NumRows()

	// using a map
	yy := make(map[string]interface{})

	rows, err := table.Read(0, nrows)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	index := 0
	for rows.Next() {
		err = rows.Scan(&yy)
		if err != nil {
			panic(err)
		}
		values = append(values, repository.RawData{
			Groupname: group,
			Series:    name,
			Seq:       int64(yy["TIME"].(float64)),
			Ind:       index,
			Value:     float64(yy["RATE"].([9]float32)[8]),
		})
		index++
	}

	err = rows.Err()
	if err != nil {
		panic(err)
	}

	bulkSave(values, repo)
}
