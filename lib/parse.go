package lib

import (
	"encoding/json"
	"os"
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

func RetrievePSD(path string) (psd *PlayerSaveData, err error) {
	xoredSaveContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	saveContent := make([]byte, len(xoredSaveContent))

	xorKey := getXorKey()
	for i := 0; i < len(xoredSaveContent); i++ {
		saveContent[i] = xoredSaveContent[i] ^ xorKey[i%len(xorKey)]
	}

	err = json.Unmarshal(saveContent, &psd)
	if err != nil {
		return nil, err
	}

	return psd, nil
}
