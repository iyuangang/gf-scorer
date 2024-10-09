package scorer

import (
	"database/sql"
	"log"
)

type Scorer struct {
    db *sql.DB
}

func New(db *sql.DB) *Scorer {
    s := &Scorer{db: db}
    s.ensureTableExists()
    return s
}

func (s *Scorer) ensureTableExists() {
    _, err := s.db.Exec(`
        CREATE TABLE IF NOT EXISTS key_id_scores (
            keyid VARCHAR(255) PRIMARY KEY,
            rl_score INT,
            il_score INT,
            dl_score INT,
            ml_score INT,
            letters_count INT,
            score FLOAT,
            filename VARCHAR(255),
            create_time TIMESTAMP WITH TIME ZONE
        )
    `)
    if err != nil {
        log.Fatalf("Failed to create table: %v", err)
    }
}
