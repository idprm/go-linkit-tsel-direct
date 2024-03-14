package repository

import (
	"context"
	"encoding/json"
	"time"

	"github.com/idprm/go-linkit-tsel/internal/domain/entity"
	"github.com/redis/go-redis/v9"
)

type VerifyRepository struct {
	rdb *redis.Client
}

type IVerifyRepository interface {
	Set(*entity.Verify) error
	Get(string) (*entity.Verify, error)
}

func NewVerifyRepository(rdb *redis.Client) *VerifyRepository {
	return &VerifyRepository{
		rdb: rdb,
	}
}

func (r *VerifyRepository) Set(t *entity.Verify) error {
	verify, _ := json.Marshal(t)
	err := r.rdb.Set(context.TODO(), t.GetToken(), string(verify), 60*time.Minute).Err()
	if err != nil {
		return err
	}
	return nil
}

func (r *VerifyRepository) Get(token string) (*entity.Verify, error) {
	val, err := r.rdb.Get(context.TODO(), token).Result()
	if err != nil {
		return nil, err
	}
	var v *entity.Verify
	json.Unmarshal([]byte(val), &v)
	return v, nil
}
