package database

import (
	"bytes"
	"database/sql"
	"dmensions/internal/ai"
	"dmensions/internal/model"
	"encoding/binary"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	_ "modernc.org/sqlite"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// --- Init Methods --- //

func InitDB() (*sql.DB, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	appDir := filepath.Join(home, ".dmensions")
	dbPath := filepath.Join(appDir, "data.db")

	dbExists, err := Exists(dbPath)

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL;"); err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
	}

	if !dbExists {
		if err := PopulateDB(NewStorage(db)); err != nil {
			return nil, err
		}
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS entities (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		content TEXT UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS embeddings (
		entity_id INTEGER UNIQUE,
		vector BLOB,
		FOREIGN KEY(entity_id) REFERENCES entities(id)
	);`
	_, err := db.Exec(query)
	return err
}

// --- Logic Help --- //

func VectorToBlob(vector []float32) ([]byte, error) {
	buf := new(bytes.Buffer)
	err := binary.Write(buf, binary.LittleEndian, vector)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BlobToVector(blob []byte) ([]float32, error) {
	vec := make([]float32, len(blob)/4)
	err := binary.Read(bytes.NewReader(blob), binary.LittleEndian, &vec)
	return vec, err
}

func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if errors.Is(err, fs.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// --- Methods for Storage --- //

func (s *Storage) SaveWord(word string) error {
	// Pre-process
	word = strings.ToLower(strings.TrimSpace(word))

	// Get Embedding
	vec, err := ai.GetEmbedding(word)
	if err != nil {
		return fmt.Errorf("ai error: %w", err)
	}

	blob, err := VectorToBlob(vec)
	if err != nil {
		return err
	}

	// Start Transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Insert Word (OR IGNORE if it exists)
	_, err = tx.Exec(
		`INSERT OR IGNORE INTO entities (content) VALUES (?)`,
		word,
	)
	if err != nil {
		return err
	}

	// Get the ID (needed for the foreign key in embeddings table)
	var entityID int64
	err = tx.QueryRow(
		`SELECT id FROM entities WHERE content = ?`,
		word,
	).Scan(&entityID)
	if err != nil {
		return err
	}

	// Upsert the embedding
	_, err = tx.Exec(`
		INSERT INTO embeddings (entity_id, vector)
		VALUES (?, ?)
		ON CONFLICT(entity_id)
		DO UPDATE SET vector = excluded.vector
	`, entityID, blob)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (s *Storage) GetAllVectors() ([]model.WordData, error) {
	rows, err := s.db.Query(`SELECT e.id, e.content, em.vector FROM entities e JOIN embeddings em ON e.id = em.entity_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entities := make([]model.WordData, 0, 100)

	for rows.Next() {
		var id int64
		var word string
		var blob []byte

		if err := rows.Scan(&id, &word, &blob); err != nil {
			return nil, err
		}

		vector, err := BlobToVector(blob)
		if err != nil {
			return nil, err
		}

		entities = append(entities, model.WordData{
			ID:     id,
			Word:   word,
			Vector: vector,
		})
	}
	return entities, nil
}

// --- Methods for Setup --- //

func PopulateDB(s *Storage) error {
	words := []string{
		// colors
		// "red", "orange", "yellow", "green", "cyan", "blue", "purple",
		// numbers
		// "zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten",
		// "0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "10",
		// animals
		"anteater", "ape", "armadillo", "baboon", "bat", "bear", "beetle", "bongo", "camel", "centipede", "chameleon", "cheetah", "cockatoo", "crane", "crocodile", "deer", "duck", "eagle", "elephant", "flamingo", "fox", "giraffe", "hamster", "hawk", "hermit crab",
		"hippo", "hippopotamus", "horse", "hummingbird", "hyena", "iguana", "impala", "jaguar", "kangaroo", "kingfisher", "kite", "kiwi", "koala", "komodo dragon", "leopard", "lion", "lizard", "mole", "monkey", "newt", "opossum", "orangutan", "ostrich", "owl", "panda", "panther",
		// "parrot", "peakock", "penguin", "pigeon", "platypus", "puffin", "rabbit", "rattlesnake", "red panda", "reindeer", "rhinoceros", "scorpion", "seal", "snake", "sparrow", "squirrel", "swan", "tiger", "turkey", "turtle", "vulture", "walrus", "wolf", "woodpecker", "yak", "zebra",
	}
	for _, word := range words {
		if err := s.SaveWord(word); err != nil {
			return err
		}
	}
	return nil
}

func ImportDB(s *sql.DB, path string) {
	// Puts an SQLite into the ./dmensions
}

func ExportDB(s *sql.DB, path string) {
	// Exports current SQLite to the path
}
