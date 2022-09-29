package lib

import (
	"encoding/json"
	"os"
	"time"
)

type LevelStat struct {
	TimeBestMicroseconds int64 `json:"_timeBestMicroseconds"` // public long _timeBestMicroseconds
}

type LevelStats struct {
	Values []LevelStat `json:"values"`
}

type PlayerSaveData struct {
	LevelStats LevelStats `json:"levelStats"`
}

func GetPlainSave(path string) (plainSave []byte, err error) {
	xoredSaveContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	saveContent := make([]byte, len(xoredSaveContent))

	xorKey := GetXorKey()
	for i := 0; i < len(xoredSaveContent); i++ {
		saveContent[i] = xoredSaveContent[i] ^ xorKey[i%len(xorKey)]
	}

	return saveContent, nil
}

func retrievePSD(path string) (psd *PlayerSaveData, err error) {
	saveContent, err := GetPlainSave(path)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(saveContent, &psd)
	if err != nil {
		return nil, err
	}

	return psd, nil
}

func RetrievePBs(path string) (pbTimes []time.Duration, err error) {
	psd, err := retrievePSD(path)
	if err != nil {
		return nil, err
	}

	pbTimes = []time.Duration{}

	for i := range GetLevels() {
		bestTime := (time.Duration(psd.LevelStats.Values[i].TimeBestMicroseconds) * time.Microsecond)
		bestTime = bestTime.Truncate(time.Microsecond * 1000)
		pbTimes = append(pbTimes, bestTime)
	}

	return pbTimes, nil
}
