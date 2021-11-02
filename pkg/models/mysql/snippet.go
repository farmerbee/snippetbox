package mysql

import (
	"blog/pkg/models"
	"database/sql"
	"errors"
)

type SnippetModel struct {
	DB *sql.DB
}

func (s *SnippetModel) Insert(title, content, expire string) (int, error) {
	statement := `INSERT INTO (snippets title, content, created, expire) VALUES
				(?, ?, UTC_TIMESTAMP, DATE_ADD(UTC_TIMESTAMP, INTERVAL ? DAY))`

	res, err := s.DB.Exec(statement, title, content, expire)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (s *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `SELECT id,title, content, created, expire FROM snippets WHERE expire > UTC_TIMESTAMP() AND id = ?`

	snippet := &models.Snippet{}
	err := s.DB.QueryRow(query, id).Scan(&snippet.Id, &snippet.Title, &snippet.Content, &snippet.Created, &snippet.Expire)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows){
			return nil, models.ErrNoRecord
		}
		return nil, err
	}


	return snippet, nil
}

// Lasted get the top ten newest records
func (s *SnippetModel) Lasted() ([]*models.Snippet, error) {
	query := `SELECT * FROM snippets WHERE expire > UTC_TIMESTAMP ORDER BY created DESC LIMIT 10`

	rows, err := s.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var snippets []*models.Snippet
	for rows.Next() {
		s := &models.Snippet{}

		err = rows.Scan(&s.Id, &s.Title, &s.Content, &s.Created, &s.Expire)
		if err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
