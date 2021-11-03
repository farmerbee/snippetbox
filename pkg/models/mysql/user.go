package mysql

import "database/sql"

type UserModel struct {
	DB *sql.DB
}

func (u *UserModel) Insert(name, email, password string) error {
	return nil
}

func (u *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

func (u *UserModel) Get(id int) (*UserModel, error) {
	return nil, nil
}
