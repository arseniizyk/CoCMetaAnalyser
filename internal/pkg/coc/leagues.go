package coc

import (
	"fmt"
	"net/http"
)

func (c *CoC) GetLeaguesInfo() (*LeagueInfo, error) {
	leagueID := 29000022 // Legendary League ID (season info only from legendary league)
	endpoint := fmt.Sprintf("%s/leagues/%d", c.url, leagueID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var league LeagueInfo
	if err := parseResponse(resp, &league); err != nil {
		return nil, err
	}

	return &league, nil
}

func (c *CoC) GetLeagueSeasons() (*LeagueSeasons, error) {
	leagueID := 29000022
	endpoint := fmt.Sprintf("%s/leagues/%d/seasons", c.url, leagueID)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var seasons LeagueSeasons
	if err := parseResponse(resp, &seasons); err != nil {
		return nil, err
	}

	return &seasons, nil
}

func (c *CoC) GetLeagueSeasonRanking(limit int, season string) (*LeagueSeasonRanking, error) {
	if limit < 100 {
		limit = 100
	}
	if limit > 25000 {
		limit = 25000
	}

	leagueID := 29000022
	endpoint := fmt.Sprintf("%s/leagues/%d/seasons/%s?limit=%d", c.url, leagueID, season, limit)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var ranking LeagueSeasonRanking
	if err := parseResponse(resp, &ranking); err != nil {
		return nil, err
	}

	return &ranking, nil
}
