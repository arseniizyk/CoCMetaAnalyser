package service

import (
	"github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/coc"
)

type Service interface {
	GetLeaguesInfo() (*coc.LeagueInfo, error)
	GetLeagueSeasons() (*coc.LeagueSeasons, error)
	GetLeagueSeasonRanking(limit int, season string) (*coc.LeagueSeasonRanking, error)
	GetPlayerInfo(playerTag string) (*coc.Player, error)
	GetMetaItems(limit int, season string) (map[string]map[string]int, error)
	GetMetaItemPairs(limit int, seasons string) (map[string]map[string]int, error)
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

func (s *service) GetLeagueSeasonRanking(limit int, season string) (*coc.LeagueSeasonRanking, error) {
	return s.cocClient.GetLeagueSeasonRanking(limit, season)
}

func (s *service) GetPlayerInfo(playerTag string) (*coc.Player, error) {
	return s.cocClient.GetPlayerInfo(playerTag)
}

func (s *service) GetMetaItems(limit int, season string) (map[string]map[string]int, error) {
	return s.cocClient.GetMetaItems(limit, season)
}

func (s *service) GetMetaItemPairs(limit int, season string) (map[string]map[string]int, error) {
	return s.cocClient.GetMetaItemPairs(limit, season)
}
