package lib

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

type RemoteWRs struct {
	WhenScraped time.Time
	WRs         []time.Duration
}

func RetrieveWRsFromWeb(url string) (wrTimes *RemoteWRs, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	wholeBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	resp.Body.Close()

	var latest RemoteWRs
	err = json.Unmarshal(wholeBody, &latest)
	if err != nil {
		return nil, err
	}

	return &latest, nil
}

func RetrieveWRsFromGoogle(apiKey string, sheet string, ranges string) (wrTimes *RemoteWRs, err error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	times := []time.Duration{}

	for _, rng := range strings.Split(ranges, ";") {
		resp, err := srv.Spreadsheets.Values.Get(sheet, rng).Do()
		if err != nil {
			return nil, err
		}
		for _, row := range resp.Values {
			replacedTime := strings.Replace(row[0].(string), ":", "m", 1) + "s"
			wrTime, _ := time.ParseDuration(replacedTime)
			times = append(times, wrTime)
		}
	}
	return &RemoteWRs{WhenScraped: time.Now(), WRs: times}, nil
}
