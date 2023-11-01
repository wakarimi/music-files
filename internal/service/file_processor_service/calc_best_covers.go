package file_processor_service

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"music-files/internal/errors"
	"sort"
)

func (s *Service) CalcBestCovers(tx *sqlx.Tx, audioFileIds []int) (bestCoverIds []int, err error) {
	coverIds := make([]int, 0)
	for _, audioFileId := range audioFileIds {
		cover, err := s.GetCoverForAudioFile(tx, audioFileId)
		if _, ok := err.(errors.NotFound); ok {
			continue
		} else if err != nil {
			log.Error().Err(err).Msg("Failed to get cover")
			return make([]int, 0), err
		}
		coverIds = append(coverIds, cover.CoverId)
	}

	bestCoverIds = uniqueSortedByFrequency(coverIds)
	return bestCoverIds, err
}

type numberFrequency struct {
	Number    int
	Frequency int
}

func uniqueSortedByFrequency(arr []int) []int {
	frequencyMap := make(map[int]int)
	for _, num := range arr {
		frequencyMap[num]++
	}

	frequencies := make([]numberFrequency, 0, len(frequencyMap))
	for num, freq := range frequencyMap {
		frequencies = append(frequencies, numberFrequency{Number: num, Frequency: freq})
	}

	sort.Slice(frequencies, func(i, j int) bool {
		return frequencies[i].Frequency > frequencies[j].Frequency
	})

	result := make([]int, len(frequencies))
	for i, freq := range frequencies {
		result[i] = freq.Number
	}

	return result
}
