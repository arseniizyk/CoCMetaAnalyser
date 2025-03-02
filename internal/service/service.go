package service

import "github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/coc"

type Service interface {
	GetLeaguesInfo() (*coc.LeagueInfo, error)
	GetLeagueSeasons() (*coc.LeagueSeasons, error)
	GetLeagueSeasonRanking(limit int) (*coc.LeagueSeasonRanking, error)
}

type service struct {
	cocClient *coc.CoC
}

func New(cocClient *coc.CoC) Service {
	return &service{cocClient: cocClient}
}

func (s *service) GetLeaguesInfo() (*coc.LeagueInfo, error) {
	return s.cocClient.GetLeaguesInfo()
}

func (s *service) GetLeagueSeasons() (*coc.LeagueSeasons, error) {
	return s.cocClient.GetLeagueSeasons()
}

func (s *service) GetLeagueSeasonRanking(limit int) (*coc.LeagueSeasonRanking, error) {
	return s.cocClient.GetLeagueSeasonRanking(limit)
}
