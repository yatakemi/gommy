package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	pb "gopkg.in/cheggaaa/pb.v2"
)

// Config for toml
type Config struct {
	Datetime DatetimeConfig
	Data     DataConfig
	Header   []HeaderConfig
}

// DatetimeConfig for toml
type DatetimeConfig struct {
	DatetimeFormat string           `toml:"datetimeFormat"`
	Start          time.Time        `toml:"start"`
	End            time.Time        `toml:"end"`
	Column         int              `toml:"column"`
	Sampling       SamplingDatetime `toml:"sampling"`
}

// SamplingDatetime for toml
type SamplingDatetime struct {
	Num  int    `toml:"num"`
	Unit string `toml:"unit"`
}

// DataConfig for toml
type DataConfig struct {
	Min       float64        `toml:"min"`
	Max       float64        `toml:"max"`
	Pointtype string         `toml:"pointtype"`
	Abnormal  []AbnormalData `toml:"abnormal"`
	Tag       []TagData      `toml:"tag"`
	Datetime  []DatetimeData `toml:"datetime"`
}

// AbnormalData for toml
type AbnormalData struct {
	Min        float64                `toml:"min"`
	Max        float64                `toml:"max"`
	Pointtype  string                 `toml:"pointtype"`
	Column     int                    `toml:"column"`
	Start      time.Time              `toml:"start"`
	End        time.Time              `toml:"end"`
	Transition TransitionAbnormalData `toml:"transition"`
}

// TransitionAbnormalData for toml
type TransitionAbnormalData struct {
	Num  int    `toml:"num"`
	Unit string `toml:"unit"`
}

// TagData for toml
type TagData struct {
	Column int      `toml:"column"`
	Rate   int      `toml:"rate"` // TODO only rate=100
	Value  []string `toml:"value"`
}

// DatetimeData for toml
type DatetimeData struct {
	Column         int    `toml:"column"`
	DatetimeFormat string `toml:"datetimeFormat"`
	Add            int    `toml:"add"` // [ms]
}

// HeaderConfig for toml
type HeaderConfig struct {
	Row []string `toml:"row"`
}

// failOnError function
func failOnError(err error) {
	if err != nil {
		log.Fatal("Error:", err)
	}
}

// Exists function is checking a file exists
func Exists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

// Question function
func Question(q string) bool {
	result := true
	fmt.Print(color.GreenString(q))

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		i := scanner.Text()

		if i == "Y" || i == "y" || i == "" {
			break
		} else if i == "N" || i == "n" {
			result = false
			break
		} else {
			fmt.Println(color.RedString("Please answer Y or N"))
			fmt.Print(color.GreenString(q))
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return result
}

// getDatetimeDuration function
func getDatetimeDuration(unitStr string) time.Duration {
	var timeDuration time.Duration
	unitStr = strings.ToLower(unitStr)
	// TODO use switch
	if unitStr == "min" {
		timeDuration = time.Minute
	} else if unitStr == "sec" {
		timeDuration = time.Second
	} else if unitStr == "hour" {
		timeDuration = time.Hour
	} else { // default
		timeDuration = time.Minute
	}
	return timeDuration
}

// randFloatn generate random float number
func randFloatn(max float64) float64 {
	return rand.Float64() * max
}

// getNormalData generate normal data
func getNormalData(config Config) string {
	switch strings.ToLower(config.Data.Pointtype) {
	case "int":
		return strconv.Itoa(rand.Intn(int(config.Data.Max-config.Data.Min)) + int(config.Data.Min))
	default:
		// case of float
		return fmt.Sprint(randFloatn(config.Data.Max-config.Data.Min) + config.Data.Min)
	}
}

// getAbnormalData generate abnormal data
func getAbnormalData(abnormalData AbnormalData, max float64, min float64) string {
	switch strings.ToLower(abnormalData.Pointtype) {
	case "int":
		return strconv.Itoa(rand.Intn(int(max-min)) + int(min))
	default:
		// case of float
		return fmt.Sprint(randFloatn(max-min) + min)
	}
}

// getDatetimeData
func getDatetimeData(datetimeData DatetimeData, current time.Time) string {
	t := current.Add(time.Duration(datetimeData.Add) * time.Millisecond)
	return t.Format(datetimeData.DatetimeFormat)
}

// combination
func combine(s1 [][]string, s2 []string) [][]string {
	var result [][]string
	for _, v1 := range s1 {
		for _, v2 := range s2 {
			cpy := make([]string, len(v1))
			copy(cpy, v1)
			result = append(result, append(cpy, v2))
		}
	}
	return result
}

// getProbabilityFilteredArray
func getProbabilityFilteredArray(rate int, data []string) []string {
	// not allow rate=0
	if rate <= 0 {
		rate = 100
	}

	// filter
	var result []string
	for _, v := range data {
		if rate < rand.Intn(100) {
			continue
		}
		result = append(result, v)
	}
	return result
}

// indexOf int array
func indexOf(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 // not found.
}

// indexOf AbnormalData array TODO use interface
func indexOfColumn(columnNumber int, abnormalData []AbnormalData) int {
	for i, v := range abnormalData {
		if columnNumber == v.Column {
			return i
		}
	}
	return -1 // not found.
}

// indexOf DatetimeData array TODO use interface
func indexOfDatetimeDataColumn(columnNumber int, datetimeData []DatetimeData) int {
	for i, v := range datetimeData {
		if columnNumber == v.Column {
			return i
		}
	}
	return -1 // not found.
}

// getCombineTagList
func getCombineTagList(config Config) ([][]string, []int) {
	var tagDataColumnIndexList []int // init slice for tag data column index
	var combineList [][]string
	for _, v := range config.Data.Tag {
		targetList := getProbabilityFilteredArray(v.Rate, v.Value)
		tagDataColumnIndexList = append(tagDataColumnIndexList, v.Column-1)
		if len(combineList) == 0 {
			sliceSize := len(targetList)
			for i := 0; i < sliceSize; i++ {
				end := i + 1
				if sliceSize < end {
					end = sliceSize
				}
				combineList = append(combineList, targetList[i:end])
			}
			continue
		}
		combineList = combine(combineList, targetList)
	}
	return combineList, tagDataColumnIndexList
}

// TODO use interface
func isBetween(start time.Time, end time.Time, target time.Time) bool {
	if (start.Before(target) || start.Equal(target)) && (end.After(target) || end.Equal(target)) {
		return true
	}
	return false
}

func getProgressCount(end time.Time, start time.Time, unitStr string) int {
	duration := end.Sub(start)
	var progressCount int
	unitStr = strings.ToLower(unitStr)
	// TODO use switch
	if unitStr == "min" {
		progressCount = int(duration.Minutes())
	} else if unitStr == "sec" {
		progressCount = int(duration.Seconds())
	} else if unitStr == "hour" {
		progressCount = int(duration.Hours())
	} else { // default
		progressCount = int(duration.Minutes())
	}
	return progressCount
}

// genarete dummy data
func generator(filename string, config Config) {
	log.Printf("datetime range: %+v to %+v\n", config.Datetime.Start, config.Datetime.End)

	// TODO
	//書き込みファイル準備
	file, err := os.Create(filename)
	// file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666) // ファイルがあれば追記
	failOnError(err)
	defer file.Close()

	// writer := csv.NewWriter(transform.NewWriter(file2, japanese.ShiftJIS.NewEncoder()))
	// writer := csv.NewWriter(transform.NewWriter(file2, japanese.EUCJP.NewEncoder()))
	writer := csv.NewWriter(file) // utf8
	writer.UseCRLF = true         // デフォルトはLFのみ

	// write header
	for _, v := range config.Header {
		writer.Write(v.Row)
	}
	writer.Flush()

	// set configuration for dummy data
	timeDuration := getDatetimeDuration(config.Datetime.Sampling.Unit) // set time duration unit for datetime
	current := config.Datetime.Start                                   // set start datetime
	layout := config.Datetime.DatetimeFormat                           // set datetime format
	columnSize := len(config.Header[0].Row)                            // set column size for data
	rand.Seed(time.Now().UnixNano())                                   // set random seed

	progressCount := getProgressCount(config.Datetime.End, current, config.Datetime.Sampling.Unit)
	progressBar := pb.StartNew(progressCount)
	// write data
	for current.Before(config.Datetime.End) { // range loop
		// generate tag data combination list
		combineList, tagDataColumnIndexList := getCombineTagList(config)

		for _, combineListValue := range combineList {
			dataRow := make([]string, columnSize) // init slice for raw data
			for i := 0; i < columnSize; i++ {
				if tagDataColumnNumberIndex := indexOf(i, tagDataColumnIndexList); tagDataColumnNumberIndex > -1 {
					// set tag data
					dataRow[i] = combineListValue[tagDataColumnNumberIndex]
				} else if config.Datetime.Column-1 == i {
					// set a base datetime data
					dataRow[i] = current.Format(layout)
				} else if abnormalDataColumnNumberIndex := indexOfColumn(i+1, config.Data.Abnormal); abnormalDataColumnNumberIndex > -1 && isBetween(config.Data.Abnormal[abnormalDataColumnNumberIndex].Start, config.Data.Abnormal[abnormalDataColumnNumberIndex].End, current) {
					// set abnormal data
					abnormalData := config.Data.Abnormal[abnormalDataColumnNumberIndex]

					transitionStart := abnormalData.Start
					transitionEnd := abnormalData.Start.Add(time.Duration(abnormalData.Transition.Num) * getDatetimeDuration(abnormalData.Transition.Unit))
					allCount := float64(getProgressCount(transitionEnd, transitionStart, config.Datetime.Sampling.Unit))
					transitionCount := 1 + float64(getProgressCount(current, abnormalData.Start, config.Datetime.Sampling.Unit))

					unitMax := (abnormalData.Max - config.Data.Max) / allCount
					unitMin := (abnormalData.Min - config.Data.Min) / allCount

					var max, min float64
					if transitionCount < allCount {
						max = unitMax * transitionCount
						min = unitMin * transitionCount
					} else {
						max = abnormalData.Max
						min = abnormalData.Min
					}
					dataRow[i] = getAbnormalData(abnormalData, max, min)
				} else if datetimeDataColumnNumberIndex := indexOfDatetimeDataColumn(i+1, config.Data.Datetime); datetimeDataColumnNumberIndex > -1 {
					// set datetime data
					dataRow[i] = getDatetimeData(config.Data.Datetime[datetimeDataColumnNumberIndex], current)
				} else {
					// set normal data
					dataRow[i] = getNormalData(config)
				}
			}
			writer.Write(dataRow)
		}
		current = current.Add(
			time.Duration(config.Datetime.Sampling.Num) * timeDuration,
		) // set current datetime
		progressBar.Increment()
	}
	writer.Flush()
	progressBar.Finish()
}

func main() {
	app := cli.NewApp()
	app.Name = "DummyGenerator"
	app.Usage = "This app create to the dummy data files."
	app.Version = "0.0.1"

	// flags
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Value: "./config.toml",
			Usage: "a config file path",
		},
		cli.StringFlag{
			Name:  "output, o",
			Value: "./dummyData.csv",
			Usage: "a creating dummy data file path",
		},
	}

	// action
	app.Action = func(c *cli.Context) error {
		// parameter check
		param := make(map[string]string)
		param["config"] = c.String("config")
		param["output"] = c.String("output")

		log.Printf("\"config\": %#v\n", param["config"])
		log.Printf("\"output\": %#v\n", param["output"])

		if !Exists(param["config"]) {
			log.Fatalf(color.RedString("%s not find"), param["config"])
		}
		if Exists(param["output"]) {
			if Question(fmt.Sprintf("%s already exists. if you wish to overwrite the file, press enter.[Y/n]", param["output"])) {
				log.Printf(color.GreenString("Overwrite the %s"), param["output"])
			} else {
				log.Fatalf(color.RedString("Rename or delete the %s"), param["output"])
			}
		}

		// config parser
		var config Config
		_, err := toml.DecodeFile(param["config"], &config)
		failOnError(err)

		// fmt.Printf("max datetime:%d\n", config.Data.Max)
		// fmt.Printf("min datetime:%d\n", config.Data.Min)
		// for k, v := range config.Header {
		// 	fmt.Printf("header row%#v %#v\n", k, v.Row)
		// }

		// create dummy data file
		generator(param["output"], config)

		return nil
	}

	app.Run(os.Args)
}
