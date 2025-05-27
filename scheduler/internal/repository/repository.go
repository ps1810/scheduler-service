package repository

import "scheduler/internal/db/sqlite"

type Repository struct {
	Job JobRepository
}

func NewRepository() *Repository {
	db := sqlite.GetSqliteDB()
	return &Repository{
		Job: NewJobRepository(db),
	}
}
