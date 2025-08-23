package repository

import (
	. "neema.co.za/rest/utils/database"
)

type Repository struct {
	*Database
}

func NewRepository(database *Database) *Repository {
	return &Repository{database}
}
