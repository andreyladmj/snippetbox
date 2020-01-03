package mysql

import (
	"andreyladmj/snippetbox/pkg/models"
	"database/sql"
)

type SnippetModel struct {
	DB *sql.DB
}

// prepared statements exists on db connection. So if connection is closed or
// in use the statement will be re-prepared on another connection. Mysql has 16382 max prep statements

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	stmt := `INSERT INTO snippets (title, content, created, expires)
		VALUES(?, ?, UTC_TIMESTAMP(), DATE_ADD(UTC_TIMESTAMP(), INTERVAL ? DAY))`

	result, err := m.DB.Exec(stmt, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() AND id = ?`

	row := m.DB.QueryRow(stmt, id)
	s := &models.Snippet{}

	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

	if err == sql.ErrNoRows {
		return nil, models.ErrNoRecord // because we don't want relate to sql or datastore-specific errors
	}

	if err != nil {
		return nil, err
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	stmt := `SELECT id, title, content, created, expires
		FROM snippets
		WHERE expires > UTC_TIMESTAMP() 
		ORDER BY created DESC
		LIMIT 10`

	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}
		err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)

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

func (m *SnippetModel) ExampleTransaction() error {
	tx, err := m.DB.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO ...")
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("UPDATE ...")
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	return err
}
