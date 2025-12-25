package database

import (
	"bytes"
	"database/sql"
	"dmensions/internal/model"
	"encoding/binary"
	"os"
	"path/filepath"
)

type Storage struct {
	db *sql.DB
}

func NewStorage(db *sql.DB) *Storage {
	return &Storage{db: db}
}

// --- Initialise a database --- //

func InitDB() (*sql.DB, error) {
	home, _ := os.UserHomeDir()
	appDir := filepath.Join(home, ".dmensions")

	if err := os.MkdirAll(appDir, 0755); err != nil {
		return nil, err
	}

	dbPath := filepath.Join(appDir, "data.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		return nil, err
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
		entity_id INTEGER,
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

// --- Methods for Storage --- //

// Saves a word into database
func (s *Storage) SaveWord(word string, vector []float32) error {
	blob, err := VectorToBlob(vector)
	if err != nil {
		return err
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	res, err := tx.Exec("INSERT OR IGNORE INTO entities (content) VALUES (?)", word)
	if err != nil {
		tx.Rollback()
		return err
	}

	var entityID int64
	rows, _ := res.RowsAffected()
	if rows == 0 {
		err = tx.QueryRow("SELECT id FROM entities WHERE content = ?", word).Scan(&entityID)
	} else {
		entityID, _ = res.LastInsertId()
	}

	_, err = tx.Exec("INSERT INTO embeddings (entity_id, vector) VALUES (?, ?)", entityID, blob)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// Get all vectors from database
func (s *Storage) GetAllVectors() ([]model.WordData, error) {
	rows, err := s.db.Query(`
		SELECT e.id, e.content, em.vector 
		FROM entities e 
		JOIN embeddings em ON e.id = em.entity_id`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []model.WordData
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
