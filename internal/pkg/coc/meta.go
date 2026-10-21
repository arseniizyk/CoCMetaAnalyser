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

const (
	maxConcurrentRequests = 25
	requestInterval       = 20 * time.Millisecond // ~50 req/sec rate limit
)

// GetMetaAndPairs fetches top ranking players and collects both single item stats and pair stats in a single pass.
func (c *CoC) GetMetaAndPairs(limit int, season string) (map[string]map[string]int, map[string]map[string]int, error) {
	start := time.Now()
	items := make(map[string]map[string]int)
	pairs := make(map[string]map[string]int)
	unavailable := 0

	ranking, err := c.GetLeagueSeasonRanking(limit, season)
	if err != nil {
		return nil, nil, err
	}

	ticker := time.NewTicker(requestInterval)
	defer ticker.Stop()

	var wg sync.WaitGroup
	var mu sync.Mutex
	sem := make(chan struct{}, maxConcurrentRequests)

	for _, p := range ranking.Players {
		<-ticker.C
		sem <- struct{}{}
		wg.Add(1)

		go func(player PlayerInRanking) {
			defer wg.Done()
			defer func() { <-sem }()

			log.Printf("Working with: %s", player.Tag)
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
				// Single items
				for _, equip := range hero.Equipment {
					if _, exists := items[hero.Name]; !exists {
						items[hero.Name] = make(map[string]int)
					}
					items[hero.Name][equip.Name]++
				}

				// Equipment pairs (requires at least 2 items)
				if len(hero.Equipment) >= 2 {
					eq1 := hero.Equipment[0].Name
					eq2 := hero.Equipment[1].Name
					// Sort lexicographically to ensure order-independent pairing
					if eq1 > eq2 {
						eq1, eq2 = eq2, eq1
					}
					pair := fmt.Sprintf("%s + %s", eq1, eq2)

					if _, exists := pairs[hero.Name]; !exists {
						pairs[hero.Name] = make(map[string]int)
					}
					pairs[hero.Name][pair]++
				}
			}
		}(p)
	}

	wg.Wait()

	log.Printf("Finished in %v. Unavailable players: %d", time.Since(start), unavailable)
	return items, pairs, nil
}

func (c *CoC) GetMetaItems(limit int, season string) (map[string]map[string]int, error) {
	items, _, err := c.GetMetaAndPairs(limit, season)
	return items, err
}

func (c *CoC) GetMetaItemPairs(limit int, season string) (map[string]map[string]int, error) {
	_, pairs, err := c.GetMetaAndPairs(limit, season)
	return pairs, err
}
