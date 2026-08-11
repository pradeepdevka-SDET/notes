package store

import (
	"database/sql"
	"notesapi/models"
)

type Store struct {
	db *sql.DB // lowercase =private, nobody outside this package touches the raw DB
}

// New builds a Store.
func New(db *sql.DB) *Store {
	return &Store{db: db}
}
func (s *Store) GetNotes() ([]models.Note, error) {
	rows, err := s.db.Query("SELECT id, title FROM notes ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	notes := []models.Note{}
	for rows.Next() {
		var n models.Note
		if err := rows.Scan(&n.ID, &n.Title); err != nil {
			return nil, err
		}
		notes = append(notes, n)
	}
	if err := rows.Err(); err != nil {
		return notes, err
	}
	return notes, nil
}

func (s *Store) GetNote(id int) (models.Note, error) {
	var n models.Note
	err := s.db.QueryRow("SELECT id, title FROM notes WHERE id =$1", id).Scan(&n.ID, &n.Title)
	return n, err
}

func (s *Store) CreateNote(title string) (models.Note, error) {
	n := models.Note{Title: title}
	err := s.db.QueryRow("INSERT INTO notes(title) VALUES ($1) RETURNING id", title).Scan(&n.ID)
	return n, err
}
func (s *Store) UpdateNote(id int, title string) (int64, error) {
	result, err := s.db.Exec("UPDATE notes SET title=$1 WHERE id=$2", title, id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

func (s *Store) DeleteNote(id int) (int64, error) {
	result, err := s.db.Exec("DELETE FROM notes WHERE id=$1", id)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
