package service

import "shorty/internal/repository"

type StatService struct {
	repo *repository.StatRepository
}

func NewStatService(repo *repository.StatRepository) *StatService {
	return &StatService{repo: repo}
}
