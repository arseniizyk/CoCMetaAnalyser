package coc

import (
	"fmt"
	"log"
	"sync"
	"time"
)

type PlayerInRanking struct {
	Tag string `json:"tag"`
}

type MetaResultSet struct {
	ItemsTotal  map[string]map[string]int
	PairsTotal  map[string]map[string]int
	ItemsTop100 map[string]map[string]int
	PairsTop100 map[string]map[string]int
}

const (
	maxConcurrentRequests = 25
	requestInterval       = 20 * time.Millisecond // ~50 req/sec rate limit
)

// GetMetaFull fetches top ranking players and collects both single item and pair stats,
// automatically separating Top 100 players from the total ranking in a single pass.
func (c *CoC) GetMetaFull(limit int, season string) (*MetaResultSet, error) {
	start := time.Now()
	res := &MetaResultSet{
		ItemsTotal:  make(map[string]map[string]int),
		PairsTotal:  make(map[string]map[string]int),
		ItemsTop100: make(map[string]map[string]int),
		PairsTop100: make(map[string]map[string]int),
	}
	unavailable := 0

	ranking, err := c.GetLeagueSeasonRanking(limit, season)
	if err != nil {
		return nil, err
	}

	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, maxConcurrentRequests)

	for i, p := range ranking.Players {
		rank := i + 1
		<-ticker.C
		sem <- struct{}{}
		wg.Add(1)

		go func(player PlayerInRanking, playerRank int) {
			defer wg.Done()
			defer func() { <-sem }()

			playerInfo, err := c.GetPlayerInfo(player.Tag)
			if err != nil {
				log.Printf("Player with tag: %s unavailable or banned: %v", player.Tag, err)
				mu.Lock()
				unavailable++
				mu.Unlock()
				return
			}

			mu.Lock()
			defer mu.Unlock()

			for _, hero := range playerInfo.Heroes {
				// Record single items
				for _, equip := range hero.Equipment {
					// Add to total
					if _, exists := res.ItemsTotal[hero.Name]; !exists {
						res.ItemsTotal[hero.Name] = make(map[string]int)
					}
					res.ItemsTotal[hero.Name][equip.Name]++

					// Add to top 100 if player is in top 100
					if playerRank <= 100 {
						if _, exists := res.ItemsTop100[hero.Name]; !exists {
							res.ItemsTop100[hero.Name] = make(map[string]int)
						}
						res.ItemsTop100[hero.Name][equip.Name]++
					}
				}

				// Record equipment pairs (requires at least 2 items)
				if len(hero.Equipment) >= 2 {
					eq1 := hero.Equipment[0].Name
					eq2 := hero.Equipment[1].Name
					// Lexicographical sort to ensure order-independent pairing
					if eq1 > eq2 {
						eq1, eq2 = eq2, eq1
					}
					pair := fmt.Sprintf("%s + %s", eq1, eq2)

					// Add to total
					if _, exists := res.PairsTotal[hero.Name]; !exists {
						res.PairsTotal[hero.Name] = make(map[string]int)
					}
					res.PairsTotal[hero.Name][pair]++

					// Add to top 100 if player is in top 100
					if playerRank <= 100 {
						if _, exists := res.PairsTop100[hero.Name]; !exists {
							res.PairsTop100[hero.Name] = make(map[string]int)
						}
						res.PairsTop100[hero.Name][pair]++
					}
				}
			}
		}(p, rank)
	}

	wg.Wait()

	log.Printf("Finished in %v. Unavailable players: %d", time.Since(start), unavailable)
	return res, nil
}

func (c *CoC) GetMetaAndPairs(limit int, season string) (map[string]map[string]int, map[string]map[string]int, error) {
	res, err := c.GetMetaFull(limit, season)
	if err != nil {
		return nil, nil, err
	}
	return res.ItemsTotal, res.PairsTotal, nil
}

func (c *CoC) GetMetaItems(limit int, season string) (map[string]map[string]int, error) {
	items, _, err := c.GetMetaAndPairs(limit, season)
	return items, err
}

func (c *CoC) GetMetaItemPairs(limit int, season string) (map[string]map[string]int, error) {
	_, pairs, err := c.GetMetaAndPairs(limit, season)
	return pairs, err
}
