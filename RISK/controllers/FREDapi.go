package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type freddata struct {
	RealtimeStart    string    `json:"realtime_start"`
	RealtimeEnd      string    `json:"realtime_end"`
	ObservationStart string    `json:"observation_start"`
	ObservationEnd   string    `json:"observation_end"`
	Units            string    `json:"units"`
	OutputType       int       `json:"output_type"`
	FileType         string    `json:"file_type"`
	OrderBy          string    `json:"order_by"`
	SortOrder        string    `json:"sort_order"`
	Count            int       `json:"count"`
	Offset           int       `json:"offset"`
	Limit            int       `json:"limit"`
	Observations     []obstype `json:"observations"`
}

type obstype struct {
	RealtimeStart string `json:"realtime_start"`
	RealtimeEnd   string `json:"realtime_end"`
	Date          string `json:"date"`
	Value         string `json:"value"`
}

type fredstruct struct {
	Date  string
	Value float64
}

func fredapi(seriesid string) (finaldata map[int]fredstruct) {
	apistring := "https://api.stlouisfed.org/fred/series/observations?series_id=" + seriesid + "&api_key=12a363ff0f6f0f9a26d3ba9a9cbcc88b&file_type=json"
	fred := &freddata{}
	response, err := http.Get(apistring)
	finaldata = make(map[int]fredstruct)
	if err != nil {
		fmt.Printf("main.go: http.Get Request failed with error %s\n", err)
	} else {
		data, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("main.go: ioutil.ReadAll failed with error %s\n", err)
		}
		err2 := json.Unmarshal([]byte(data), &fred)
		if err2 != nil {
			fmt.Printf("main.go: json.Unmarshal failed with error %s\n", err2)
		}
		var obsdata []obstype
		obsdata = fred.Observations
		for k, v := range obsdata {
			val, _ := strconv.ParseFloat(v.Value, 64)
			finaldata[k] = fredstruct{
				Date:  v.Date,
				Value: val,
			}
		}
		return finaldata
	}
	return finaldata
}

func fredapisorted(seriesid string) (sorted map[int]fredstruct) {
	unsorted := fredapi(seriesid)
	len := len(unsorted)
	sorted = make(map[int]fredstruct, len)
	for k := 0; k < len; k++ {
		sorted[k] = fredstruct{
			Date:  unsorted[k].Date,
			Value: unsorted[k].Value,
		}
		//fmt.Println(unsorted[k].Date, ": ", unsorted[k].Value)
	}
	return sorted
}

/*func main() {
	x := fredapisorted("A191RL1Q225SBEA")
	fmt.Println(x)
}
*/
