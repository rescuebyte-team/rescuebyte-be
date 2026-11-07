package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type commentRepository struct {
	dbGetter database.DBGetter
}

func NewCommentRepository(dbGetter database.DBGetter) contracts.CommentRepository {
	return &commentRepository{dbGetter: dbGetter}
}
