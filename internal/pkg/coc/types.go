package coc

import "net/http"

type CoC struct {
	url    string
	client *http.Client
	api    string
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
	Items []struct {
		Tag string `json:"tag"`
	}
}

func New(client *http.Client, api string) *CoC {
	return &CoC{
		url:    "https://api.clashofclans.com/v1",
		client: client,
		api:    api,
	}
}
