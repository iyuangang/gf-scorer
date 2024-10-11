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
    var err error
    if s.config.Database.Type == "postgres" {
        _, err = s.db.Exec(`
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
    } else { // SQLite
        _, err = s.db.Exec(`
            CREATE TABLE IF NOT EXISTS key_id_scores (
                keyid TEXT PRIMARY KEY,
                rl_score INTEGER,
                il_score INTEGER,
                dl_score INTEGER,
                ml_score INTEGER,
                letters_count INTEGER,
                score REAL,
                filename TEXT,
                create_time DATETIME
            )
        `)
    }
    if err != nil {
        log.Fatalf("Failed to create key_id_scores table: %v", err)
    }

    if s.config.Database.Type == "postgres" {
        _, err = s.db.Exec(`
            CREATE TABLE IF NOT EXISTS gpg_ed25519_keys (
                fingerprint VARCHAR(255) PRIMARY KEY,
                public_key TEXT,
                private_key TEXT,
                rl_score INT,
                il_score INT,
                dl_score INT,
                ml_score INT,
                score INT,
                letters_count INT
            )
        `)
    } else { // SQLite
        _, err = s.db.Exec(`
            CREATE TABLE IF NOT EXISTS gpg_ed25519_keys (
                fingerprint TEXT PRIMARY KEY,
                public_key TEXT,
                private_key TEXT,
                rl_score INTEGER,
                il_score INTEGER,
                dl_score INTEGER,
                ml_score INTEGER,
                score INTEGER,
                letters_count INTEGER
            )
        `)
    }
    if err != nil {
        log.Fatalf("Failed to create gpg_ed25519_keys table: %v", err)
    }
}
