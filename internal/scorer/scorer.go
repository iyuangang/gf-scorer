package scorer

import (
	"database/sql"
	"gf-scorer/internal/config"
	"log"
)

type Scorer struct {
    db *sql.DB
    config *config.Config
}

func New(db *sql.DB, cfg *config.Config) *Scorer {
    s := &Scorer{db: db, config: cfg}
    s.ensureTablesExist()
    return s
}

func (s *Scorer) ensureTablesExist() {
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
        log.Fatalf("Failed to create key_id_scores table: %v", err)
    }

    _, err = s.db.Exec(`
        CREATE TABLE IF NOT EXISTS gpg_keys (
            fingerprint VARCHAR(255) PRIMARY KEY,
            public_key TEXT,
            private_key TEXT,
            score INT,
            letters_count INT
        )
    `)
    if err != nil {
        log.Fatalf("Failed to create gpg_keys table: %v", err)
    }
}
