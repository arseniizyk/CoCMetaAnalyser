package coc

import (
	"fmt"
	"log"
	"time"
)

type PlayerInRanking struct {
	Tag string `json:"tag"`
}

func (c *CoC) GetMetaItems(limit int, season string) (map[string]map[string]int, error) {
	start := time.Now()
	items := make(map[string]map[string]int)
	unavailable := 0

	ranking, err := c.GetLeagueSeasonRanking(limit, season)
	if err != nil {
		return nil, err
	}

	// 20ms(50rps) to avoid rate limits from coc api
	ticker := time.NewTicker(20 * time.Millisecond)

	for _, p := range ranking.Players {
		<-ticker.C

		c.wg.Add(1)

		go func(player PlayerInRanking) {
			defer c.wg.Done()
			log.Printf("Working with: %s", player.Tag)
			equipments, err := c.GetPlayerInfo(player.Tag)
			if err != nil {
				log.Println("player with tag", player.Tag, "was banned or deleted")

				c.mu.Lock()
				unavailable++
				c.mu.Unlock()

				return
			}
			for _, hero := range equipments.Heroes {
				for _, equip := range hero.Equipment {
					c.mu.Lock()
					if _, exists := items[hero.Name]; !exists {
						items[hero.Name] = make(map[string]int)
					}
					items[hero.Name][equip.Name]++
					c.mu.Unlock()
				}
			}

		}(p)
	}

	c.wg.Wait()

	log.Println(unavailable)
	log.Println(time.Since(start))
	return items, nil
}

func (c *CoC) GetMetaItemPairs(limit int, season string) (map[string]map[string]int, error) {
	start := time.Now()

	pairs := make(map[string]map[string]int)
	unavailablePlayers := 0

	ranking, err := c.GetLeagueSeasonRanking(limit, season)
	if err != nil {
		return nil, err
	}

	// 20ms(50rps) to avoid rate limits from coc api
	ticker := time.NewTicker(20 * time.Millisecond)
	defer ticker.Stop()

	for _, p := range ranking.Players {
		<-ticker.C

		c.wg.Add(1)

		go func(player PlayerInRanking) {
			defer c.wg.Done()
			log.Printf("Working with: %s\n", player.Tag)

			playerInfo, err := c.GetPlayerInfo(player.Tag)
			if err != nil {
				log.Printf("Player with tag: %s unavailable", player.Tag)
				c.mu.Lock()
				unavailablePlayers++
				c.mu.Unlock()
				return
			}

			for _, hero := range playerInfo.Heroes {
				if len(hero.Equipment) == 0 {
					continue
				}
				pair := fmt.Sprintf("%s + %s", hero.Equipment[0].Name, hero.Equipment[1].Name)

				c.mu.Lock()
				if _, exists := pairs[hero.Name]; !exists {
					pairs[hero.Name] = make(map[string]int)
				}
				pairs[hero.Name][pair]++
				c.mu.Unlock()
			}
		}(p)
	}

	c.wg.Wait()
	log.Println(unavailablePlayers)
	log.Println(time.Since(start))
	return pairs, nil
}
