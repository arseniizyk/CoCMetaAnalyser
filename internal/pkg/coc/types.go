package coc

import (
	"net/http"
	"sync"
)

type CoC struct {
	url    string
	api    string
	client *http.Client
	wg     sync.WaitGroup
	mu     sync.Mutex
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
		mu:     sync.Mutex{},
		wg:     sync.WaitGroup{},
	}
}
