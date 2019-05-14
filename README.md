# timeseries-query

## Set-up

* Install Go. https://golang.org/
* Clone the repository to `$GOPATH/src/github.com/lhhong/timeseries-query`
  + `npm install -g bower`
  + `cd` into project, `bower install`
* Install MariaDB
  + Current installation properties: 
    - user: `root` | password: *none*
	- user: `dbuser` | password: `user_password`
* Install Redis
  + Installation on Windows requires WSL https://docs.microsoft.com/en-us/windows/wsl/install-win10
  + Current installation: 
    - Ubuntu subsystem installed, user: `ubuntu` | password: `ubuntu`
	- Redis installed and running on that subsystem

## To index data

Go to the main directory of the application using powershell or cmd. Make sure to have the folder `index` present.
The following operations saves the data to database, index the data, then output binary file of the index to the `index` folder

### Stocks data set

```go run cmd/loader/loader.go -f <file location of stocks data csv> -n stocks -s 0 -d 1 -v 5```

We also accept any time series data in the form of csv file

```go run cmd/loader/loader.go -f <file location of data csv> -n <name of dataset> -s <column number of series label in csv> -d <column number of date in csv> -v <column number of value in csv>```

### ECG dataset

```go run cmd/loader/loader.go --ecg-data -n ECG --dir <directory containing all ecg series to index>```

### Indexing data already in the database

We also have an option to index data already on the database

```go run cmd/loader/loader.go --index-only -n <name of dataset [stocks, ECG, ...]>```

## Running the Server

Go to the main directory of the application using powershell or cmd. 

Default configurations can be found in `conf/default.toml`.
Use `seriesGroup` config to decide which dataset to load onto the system.
You may create your own custom configuration file and it will overwrite differences with `default.toml`.

```go run cmd/server/server.go [-c <custom config file>]```

The query page will be available at http://localhost:8080 or the httpServer port defined in configuration.

## Building Executable

Instead of using `go run`, we can build exe files using `go build cmd/loader/loader.go` or `go build cmd/server/server.go`.
Run the exe in powershell with the corresponding arguments.
Make sure that the `index` and `conf` folders are in your current directory when running.

