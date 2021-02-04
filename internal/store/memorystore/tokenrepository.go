package memorystore

import (
	"godmin/internal/dto"
	"strconv"
	"time"
)

type TokenRepository struct {
	store *Store
}

func (r *TokenRepository) Create(userId uint64, t *dto.Token) error {
	at := time.Unix(t.AccessTokenExpires, 0) //converting Unix to UTC(to Time object)
	rt := time.Unix(t.RefreshTokenExpires, 0)
	now := time.Now()

	errAccess := r.store.client.Set(t.AccessUuid, strconv.Itoa(int(userId)), at.Sub(now)).Err()
	if errAccess != nil {
		return errAccess
	}

	errRefresh := r.store.client.Set(t.RefreshUuid, strconv.Itoa(int(userId)), rt.Sub(now)).Err()
	if errRefresh != nil {
		return errRefresh
	}

	return nil
}

func (r *TokenRepository) Find(accessUuid string) (uint64, error) {
	userIdRaw, err := r.store.client.Get(accessUuid).Result()
	if err != nil {
		return 0, err
	}

	userId, _ := strconv.ParseUint(userIdRaw, 10, 64)

	return userId, nil
}

func (r *TokenRepository) Delete(accessUuid string) (int64, error) {
	deleted, err := r.store.client.Del(accessUuid).Result()
	if err != nil {
		return 0, err
	}

	return deleted, nil
}
