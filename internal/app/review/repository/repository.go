package repository

import (
	"github.com/rescuebyte/rescuebyte-be/domain/contracts"
	"github.com/rescuebyte/rescuebyte-be/internal/infra/database"
)

type reviewRepository struct {
	dbGetter database.DBGetter
}

func NewReviewRepository(dbGetter database.DBGetter) contracts.ReviewRepository {
	return &reviewRepository{dbGetter: dbGetter}
}
