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

func XorSave(xorSave []byte) (plainSave []byte) {
	saveContent := make([]byte, len(xorSave))

	xorKey := GetXorKey()
	for i := 0; i < len(xorSave); i++ {
		saveContent[i] = xorSave[i] ^ xorKey[i%len(xorKey)]
	}

	return saveContent
}

func RetrievePBsFromSave(save []byte) (pbTimes []time.Duration, err error) {
	saveContent := XorSave(save)

	var psd PlayerSaveData

	err = json.Unmarshal(saveContent, &psd)
	if err != nil {
		return nil, err
	}

	return psd.Durations(), nil
}

func RetrievePBsFromDisk(path string) (pbTimes []time.Duration, err error) {
	xoredSaveContent, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return RetrievePBsFromSave(xoredSaveContent)
}

func (p *PlayerSaveData) Durations() []time.Duration {
	pbTimes := []time.Duration{}

	for i := range GetLevels() {
		bestTime := (time.Duration(p.LevelStats.Values[i].TimeBestMicroseconds) * time.Microsecond)
		bestTime = bestTime.Truncate(time.Microsecond * 1000)
		pbTimes = append(pbTimes, bestTime)
	}

	return pbTimes
}
