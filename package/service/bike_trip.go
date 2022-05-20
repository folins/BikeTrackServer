package service

import (
	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/repository"
)

type BikeTripService struct {
	repo repository.BikeTrip
}

func NewBikeTripService(repo repository.BikeTrip) *BikeTripService {
	return &BikeTripService{repo: repo}
}

func (s *BikeTripService) Create(userId int, trips biketrackserver.BikeTrip) (int, error) {
	return s.repo.Create(userId, trips)
}

func (s *BikeTripService) GetAll(userId int) ([]biketrackserver.BikeTrip, error) {
	return s.repo.GetAll(userId)
}

func (s *BikeTripService) GetById(userId, tripId int) (biketrackserver.BikeTrip, error) {
	return s.repo.GetById(userId, tripId)
}

func (s *BikeTripService) Delete(userId, tripId int) error {
	return s.repo.Delete(userId, tripId)
}
