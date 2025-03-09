package service

import (
	"log"
	"sync"
	"time"

	"github.com/arseniizyk/CoCMetaAnalyser/internal/pkg/coc"
)

type Service interface {
	GetLeaguesInfo() (*coc.LeagueInfo, error)
	GetLeagueSeasons() (*coc.LeagueSeasons, error)
	GetLeagueSeasonRanking(limit int, season string) (*coc.LeagueSeasonRanking, error)
	GetPlayerInfo(playerTag string) (*coc.Player, error)
	GetItemsMeta(season string, players int) (map[string]map[string]int, error)
	GetItemPairsMeta(seasons string, players int) (map[string]map[string]int, error)
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
	start := time.Now()
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
	log.Println(time.Since(start))
	return result, nil
}

func (s *service) GetItemPairsMeta(seasons string, players int) (map[string]map[string]int, error) {
	start := time.Now()
	wg := sync.WaitGroup{}
	result := make(map[string]map[string]int)
	var mu sync.Mutex
	bannedOrDeleted := 0

	rankings, err := s.GetLeagueSeasonRanking(players, seasons)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()

	for i, player := range rankings.Items {
		<-ticker.C

		wg.Add(1)

		go func(p struct {
			Tag string `json:"tag"`
		}, i int) {
			defer wg.Done()
			log.Println(i, "| working with ", player.Tag)

			equipments, err := s.GetPlayerInfo(player.Tag)
			if err != nil {
				log.Println("player with tag", player.Tag, "was banned or deleted")
				mu.Lock()
				bannedOrDeleted++
				mu.Unlock()
				return
			}

			for _, hero := range equipments.Heroes {
				// to skip heroes from builder base
				if hero.Name == "Battle Copter" || hero.Name == "Battle Machine" {
					continue
				}
				pair := ""
				for _, equip := range hero.Equipment {
					if pair != "" {
						pair += " + "
					}
					pair += equip.Name
				}

				mu.Lock()
				if _, exists := result[hero.Name]; !exists {
					result[hero.Name] = make(map[string]int)
				}
				result[hero.Name][pair]++
				mu.Unlock()
			}
		}(player, i)
	}

	wg.Wait()
	log.Println(bannedOrDeleted)
	log.Println(time.Since(start))
	return result, nil
}
