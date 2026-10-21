package coc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CoC struct {
	url    string
	api    string
	client *http.Client
}

type BadResponse struct {
	Reason  string `json:"reason"`
	Message string `json:"message"`
	Type    string `json:"type"`
}

type Player struct {
	Tag    string `json:"tag"`
	Heroes []struct {
		Name      string `json:"name"`
		Equipment []struct {
			Name string `json:"name"`
		} `json:"equipment,omitempty"`
		Village string `json:"village"`
	} `json:"heroes"`
}

type LeagueInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type LeagueSeasons struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

type LeagueSeasonRanking struct {
	Players []struct {
		Tag string `json:"tag"`
	} `json:"items"`
}

func New(client *http.Client, api string) *CoC {
	return &CoC{
		url:    "https://api.clashofclans.com/v1",
		client: client,
		api:    api,
	}
}

func parseResponse(resp *http.Response, target any) error {
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("reading response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var badResp BadResponse
		if err := json.Unmarshal(bodyBytes, &badResp); err == nil && badResp.Message != "" {
			return fmt.Errorf("status code: %d, reason: %s, message: %s", resp.StatusCode, badResp.Reason, badResp.Message)
		}
		return fmt.Errorf("status code: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	if err := json.Unmarshal(bodyBytes, target); err != nil {
		return fmt.Errorf("decoding response JSON: %w", err)
	}

	return nil
}
