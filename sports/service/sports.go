package service

import (
	"git.neds.sh/matty/entain/sports/db"
	"git.neds.sh/matty/entain/sports/proto/sports"
	"golang.org/x/net/context"
)

type Sports interface {
	// ListSports will return a collection of races.
	ListSports(ctx context.Context, in *sports.ListSportsRequest) (*sports.ListSportsResponse, error)
	// GetSportsEventDetails will returns details of a race
	GetSportsEventDetails(ctx context.Context, in *sports.GetSportsDetailsRequest) (*sports.GetSportsDetailsResponse, error)
}

// sportsService implements the Sports interface.
type sportsService struct {
	sportsRepo db.SportsRepo
}

// NewSportsService instantiates and returns a new sportsService.
func NewSportsService(sportsRepo db.SportsRepo) Sports {
	return &sportsService{sportsRepo}
}

func (s *sportsService) ListSports(ctx context.Context, in *sports.ListSportsRequest) (*sports.ListSportsResponse, error) {
	sportsEvents, err := s.sportsRepo.List(in.Filter)
	if err != nil {
		return nil, err
	}

	return &sports.ListSportsResponse{SportsEvents: sportsEvents}, nil
}

func (s *sportsService) GetSportsEventDetails(ctx context.Context, in *sports.GetSportsDetailsRequest) (*sports.GetSportsDetailsResponse, error) {
	eventDetails, err := s.sportsRepo.GetSportsEventDetails(in.Id)
	if err != nil {
		return nil, err
	}
	return &sports.GetSportsDetailsResponse{SportsEvents: eventDetails}, nil
}
