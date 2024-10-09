package metrics

import (
	"expvar"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ProcessedFiles = expvar.NewInt("processed_files")
	ProcessedLines = expvar.NewInt("processed_lines")

	FileProcessingTime = promauto.NewHistogram(prometheus.HistogramOpts{
		Name: "file_processing_time_seconds",
		Help: "Time taken to process a file",
	})

	ScoresCalculated = promauto.NewCounter(prometheus.CounterOpts{
		Name: "scores_calculated_total",
		Help: "The total number of scores calculated",
	})
)

func RecordFileProcessingTime(duration time.Duration) {
	FileProcessingTime.Observe(duration.Seconds())
}

func IncrementProcessedFiles() {
	ProcessedFiles.Add(1)
}

func IncrementProcessedLines(count int) {
	ProcessedLines.Add(int64(count))
}

func IncrementScoresCalculated() {
	ScoresCalculated.Inc()
}
