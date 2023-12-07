package handlers

import "keeper/storage"

type Handler struct {
	userStorage         *storage.UserStorage
	refreshTokenStorage *storage.RefreshTokenStorage
}

func NewHandler(us *storage.UserStorage, rts *storage.RefreshTokenStorage) *Handler {
	return &Handler{
		userStorage:         us,
		refreshTokenStorage: rts,
	}
}
