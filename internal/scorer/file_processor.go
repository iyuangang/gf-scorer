package scorer

import (
	"gf-scorer/internal/config"
	"gf-scorer/internal/metrics"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

func (s *Scorer) ProcessInput(path string, cfg config.ProcessingConfig) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}

	if info.IsDir() {
		return s.processDirectory(path, cfg)
	}
	return s.processFile(path, cfg.BatchSize)
}

func (s *Scorer) processDirectory(dirPath string, cfg config.ProcessingConfig) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	errChan := make(chan error, len(files))
	semaphore := make(chan struct{}, cfg.MaxConcurrentFiles)

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".log") {
			wg.Add(1)
			go func(f os.DirEntry) {
				semaphore <- struct{}{} // Acquire semaphore
				defer func() {
					<-semaphore // Release semaphore
					wg.Done()
				}()
				if err := s.processFile(filepath.Join(dirPath, f.Name()), cfg.BatchSize); err != nil {
					errChan <- err
				}
			}(file)
		}
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	for err := range errChan {
		if err != nil {
			return err // Return the first error encountered
		}
	}

	return nil
}

func (s *Scorer) processFile(filePath string, batchSize int) error {
	startTime := time.Now()
	defer func() {
		metrics.RecordFileProcessingTime(time.Since(startTime))
		metrics.IncrementProcessedFiles()
	}()

	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	lines := strings.Split(string(content), "\n")
	lineCount := 0
	batch := make([]ScoreRecord, 0, batchSize)
	filename := filepath.Base(filePath)

	for _, line := range lines {
		if line = strings.TrimSpace(line); line != "" {
			scores := calculateScores(line)
			totalScore := scores.RLScore + scores.ILScore + scores.DLScore + scores.MLScore

			record := ScoreRecord{
				KeyID:        line,
				RLScore:      scores.RLScore,
				ILScore:      scores.ILScore,
				DLScore:      scores.DLScore,
				MLScore:      scores.MLScore,
				LettersCount: scores.LettersCount,
				TotalScore:   totalScore,
				Filename:     filename,
				CreateTime:   time.Now(),
			}

			batch = append(batch, record)
			lineCount++

			if len(batch) >= batchSize {
				if err := s.processBatch(batch); err != nil {
					return err
				}
				batch = batch[:0] // Clear the batch
			}
		}
	}

	// Process any remaining records
	if len(batch) > 0 {
		if err := s.processBatch(batch); err != nil {
			return err
		}
	}

	metrics.IncrementProcessedLines(lineCount)
	return nil
}
