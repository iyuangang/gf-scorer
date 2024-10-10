package scorer

import (
	"time"
)

type ScoreRecord struct {
	KeyID        string
	RLScore      int
	ILScore      int
	DLScore      int
	MLScore      int
	LettersCount int
	TotalScore   int
	Filename     string
	CreateTime   time.Time
}

func (s *Scorer) processBatch(records []ScoreRecord) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO key_id_scores (keyid, rl_score, il_score, dl_score, ml_score, letters_count, score, filename, create_time)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		_, err := stmt.Exec(
			record.KeyID,
			record.RLScore,
			record.ILScore,
			record.DLScore,
			record.MLScore,
			record.LettersCount,
			record.TotalScore,
			record.Filename,
			record.CreateTime,
		)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	//metrics.IncrementScoresCalculated(len(records))
	return nil
}
