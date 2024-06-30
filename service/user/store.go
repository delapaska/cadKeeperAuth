package user

import (
	"database/sql"
	"fmt"

	"github.com/delapaska/cadKeeperAuth/models"
)

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) GetUserByEmail(email string) (*models.User, error) {

	rows, err := s.db.Query(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))

	if err != nil {
		return nil, err
	}

	u := new(models.User)

	for rows.Next() {

		u, err = scanRowIntoUser(rows)

		if err != nil {
			return nil, err
		}

	}

	if u.ID == 0 {

		return nil, fmt.Errorf("user not found")
	}

	return u, nil
}

func scanRowIntoUser(rows *sql.Rows) (*models.User, error) {
	user := new(models.User)

	err := rows.Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.Password,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *Store) GetUserById(id int) (*models.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", id)
	rows, err := s.db.Query(query)

	if err != nil {
		return nil, err
	}

	u := new(models.User)
	for rows.Next() {
		u, err = scanRowIntoUser(rows)
		if err != nil {
			return nil, err
		}

	}
	if u.ID == 0 {
		return nil, fmt.Errorf("user not found")
	}
	return u, nil
}

func (s *Store) CreateUser(user models.User) error {

	query := fmt.Sprintf("INSERT INTO users (email, username, password) VALUES('%s', '%s', '%s')", user.Email, user.Username, user.Password)

	_, err := s.db.Exec(query)

	if err != nil {
		return err
	}
	return nil
}
