package coc

import (
	"encoding/json"
	"fmt"
	"net/http"
)

const (
	badResponse = "Status code: %v, message: %v"
)

func (c *CoC) GetLeaguesInfo() (*LeagueInfo, error) {
	leagueID := 29000022 // Legendary League ID(CoC can get season info only from legendary league)
	url := fmt.Sprintf(c.url+"/leagues/%d", leagueID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var badResp BadResponse
		err = json.NewDecoder(resp.Body).Decode(&badResp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(badResponse, resp.StatusCode, badResp)
	}

	var league LeagueInfo
	err = json.NewDecoder(resp.Body).Decode(&league)
	if err != nil {
		return nil, err
	}

	return &league, nil
}

func (c *CoC) GetLeagueSeasons() (*LeagueSeasons, error) {
	leagueID := 29000022
	url := fmt.Sprintf(c.url+"/leagues/%d/seasons", leagueID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var badResp BadResponse
		err = json.NewDecoder(resp.Body).Decode(&badResp)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf(badResponse, resp.StatusCode, badResp)
	}

	var seasons LeagueSeasons
	err = json.NewDecoder(resp.Body).Decode(&seasons)
	if err != nil {
		return nil, err
	}

	return &seasons, nil
}

func (c *CoC) GetLeagueSeasonRanking(limit int) (*LeagueSeasonRanking, error) {
	// leagueID hardcoded, cuz league seasons information is available only for Legendary League
	// but seasonID could be different
	if limit > 25000 {
		limit = 25000
	}
	leagueID := 29000022
	seasonID := "2025-02"
	url := fmt.Sprintf(c.url+"/leagues/%d/seasons/%s", leagueID, seasonID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		var badResp BadResponse
		err := json.NewDecoder(resp.Body).Decode(&badResp)
		if err != nil {
			return nil, err
		}

		return nil, fmt.Errorf(badResponse, resp.StatusCode, badResp)
	}

	var ranking LeagueSeasonRanking
	err = json.NewDecoder(resp.Body).Decode(&ranking)
	if err != nil {
		return nil, err
	}

	return &ranking, nil
}
