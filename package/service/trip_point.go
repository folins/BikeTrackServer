package service

import (
	"github.com/folins/biketrackserver"
	"github.com/folins/biketrackserver/package/repository"
)

type TripPointService struct {
	repo repository.TripPoint
	tripRepo repository.BikeTrip
}

func NewTripPointService(repo repository.TripPoint, tripRepo repository.BikeTrip) *TripPointService {
	return &TripPointService{
		repo: repo,
		tripRepo: tripRepo,
	}
}

func (s *TripPointService) Create(userId, tripId int, point biketrackserver.TripPoint) (int, error) {
	_, err := s.tripRepo.GetById(userId, tripId)
	if err != nil {
		// Trip is not exist
		return 0, err
	}
	
	return	s.repo.Create(tripId, point)
}

func (s *TripPointService) GetAll(userId, tripId int) ([]biketrackserver.TripPoint, error) {
	_, err := s.tripRepo.GetById(userId, tripId)
	if err != nil {
		// Trip is not exist
		return nil, err
	}

	return s.repo.GetAll(tripId)
}

func (s *TripPointService) GetById(userId, tripId, pointId int) (biketrackserver.TripPoint, error) {
	_, err := s.tripRepo.GetById(userId, tripId)
	if err != nil {
		// Trip is not exist
		return biketrackserver.TripPoint{}, err
	}

	return s.repo.GetById(tripId, pointId)
}