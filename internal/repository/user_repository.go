package repository

import "shorty/pkg/db"

type UserRepository struct {
	Database *db.DB
}

func NewUserRepository(db *db.DB) *UserRepository {
	return &UserRepository{Database: db}
}

func (u *UserRepository) Create() {}

func (u *UserRepository) GetAll() {}

func (u *UserRepository) GetByID() {}

func (u *UserRepository) GetByEmail() {}

func (u *UserRepository) Update() {}

func (u *UserRepository) Delete() {}
