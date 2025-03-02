package service

import (
	"log"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/coc"
)

type Service interface {
	GetLeaguesInfo() (*coc.LeagueInfo, error)
	GetLeagueSeasons() (*coc.LeagueSeasons, error)
	GetLeagueSeasonRanking(limit int, season string) (*coc.LeagueSeasonRanking, error)
	GetPlayerInfo(playerTag string) (*coc.Player, error)
	GetItemsMeta(season string, players int) (map[string]map[string]int, error)
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

func (s *service) GetItemsMeta(season string, players int) (map[string]map[string]int, error) {
	result := make(map[string]map[string]int)
	bannedOrDeleted := 0

	rankings, err := s.GetLeagueSeasonRanking(players, season)
	if err != nil {
		return nil, err
	}
	count := 0
	for _, player := range rankings.Items {
		count++
		log.Println(count, "| working with ", player.Tag)
		equipments, err := s.GetPlayerInfo(player.Tag)
		if err != nil {
			log.Println("player with tag", player.Tag, "was banned or deleted")
			bannedOrDeleted++
			continue
		}
		for _, hero := range equipments.Heroes {
			for _, equip := range hero.Equipment {
				if _, exists := result[hero.Name]; !exists {
					result[hero.Name] = make(map[string]int)
				}
				result[hero.Name][equip.Name]++
			}
		}
	}

	log.Println(bannedOrDeleted)
	return result, nil
}
