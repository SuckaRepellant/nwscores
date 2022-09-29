package lib

import (
	"context"
	"strings"
	"time"

	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
)

func RetrieveWRsFromGoogle(apiKey string, sheet string, ranges string) (wrTimes []time.Duration, err error) {
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return nil, err
	}

	wrTimes = []time.Duration{}

	for _, rng := range strings.Split(ranges, ";") {
		resp, err := srv.Spreadsheets.Values.Get(sheet, rng).Do()
		if err != nil {
			return nil, err
		}
		for _, row := range resp.Values {
			replacedTime := strings.Replace(row[0].(string), ":", "m", 1) + "s"
			wrTime, _ := time.ParseDuration(replacedTime)
			wrTimes = append(wrTimes, wrTime)
		}
	}
	return wrTimes, nil
}
