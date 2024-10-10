package scorer

import (
	"crypto"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"
)

type KeyInfo struct {
	Fingerprint   string
	PublicKey     string
	PrivateKey    string
	Score         int
	LettersCount  int
}

func (s *Scorer) GenerateKeys(numKeys int, numWorkers int) error {

	var wg sync.WaitGroup
	keysPerWorker := numKeys / numWorkers
	keyInfoChan := make(chan KeyInfo, numKeys)

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func(workerId int) {
			defer wg.Done()
			for j := 0; j < keysPerWorker; j++ {
				if keyInfo, err := s.generateKeyPair(); err == nil {
					keyInfoChan <- keyInfo
				}
			}
		}(i)
	}

	go func() {
		wg.Wait()
		close(keyInfoChan)
	}()

	return s.processKeyInfo(keyInfoChan)
}

func (s *Scorer) generateKeyPair() (KeyInfo, error) {

	entity, err := openpgp.NewEntity("Tim Yuan", "Comment", "yuangang@me.com", &packet.Config{
		DefaultHash:   crypto.SHA256,
		Time:          time.Now,
		Algorithm:     packet.PubKeyAlgoEdDSA,
	})
	if err != nil {
		return KeyInfo{}, fmt.Errorf("failed to create entity: %w", err)
	}

	fingerprint := fmt.Sprintf("%x", entity.PrimaryKey.Fingerprint)
	scores := calculateScores(fingerprint[len(fingerprint)-16:])
	totalScore := float64(scores.RLScore + scores.ILScore + scores.DLScore + scores.MLScore)
	
	if totalScore <= 200 && scores.LettersCount >= 5 {
		return KeyInfo{}, fmt.Errorf("key does not meet criteria")
	}

	pubKeyBuf := new(strings.Builder)
	privKeyBuf := new(strings.Builder)

	pubKeyArmor, err := armor.Encode(pubKeyBuf, openpgp.PublicKeyType, nil)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("failed to encode public key: %w", err)
	}
	entity.Serialize(pubKeyArmor)
	pubKeyArmor.Close()

	privKeyArmor, err := armor.Encode(privKeyBuf, openpgp.PrivateKeyType, nil)
	if err != nil {
		return KeyInfo{}, fmt.Errorf("failed to encode private key: %w", err)
	}
	entity.SerializePrivate(privKeyArmor, nil)
	privKeyArmor.Close()

	return KeyInfo{
		Fingerprint:  fingerprint,
		PublicKey:    pubKeyBuf.String(),
		PrivateKey:   privKeyBuf.String(),
		Score:        int(totalScore),
		LettersCount: scores.LettersCount,
	}, nil
}

func (s *Scorer) processKeyInfo(keyInfoChan <-chan KeyInfo) error {
	batch := make([]KeyInfo, 0, s.config.Processing.BatchSize)
	for keyInfo := range keyInfoChan {
		batch = append(batch, keyInfo)
		if len(batch) >= s.config.Processing.BatchSize {
			if err := s.insertKeyBatch(batch); err != nil {
				return err
			}
			batch = batch[:0]
		}
	}
	if len(batch) > 0 {
		return s.insertKeyBatch(batch)
	}
	return nil
}

func (s *Scorer) insertKeyBatch(batch []KeyInfo) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(`
		INSERT INTO gpg_keys (fingerprint, public_key, private_key, score, letters_count)
		VALUES ($1, $2, $3, $4, $5)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, keyInfo := range batch {
		_, err := stmt.Exec(
			keyInfo.Fingerprint,
			keyInfo.PublicKey,
			keyInfo.PrivateKey,
			keyInfo.Score,
			keyInfo.LettersCount,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
