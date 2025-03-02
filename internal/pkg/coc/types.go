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
}

type LeagueInfo struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	// IconUrls struct {
	// 	Small  string `json:"small"`
	// 	Tiny   string `json:"tiny"`
	// 	Medium string `json:"medium"`
	// } `json:"iconUrls"`
}

type LeagueSeasons struct {
	Items []struct {
		ID string `json:"id"`
	} `json:"items"`
}

type LeagueSeasonRanking struct {
	Items []struct {
		// dont need every line instead of tag, because GetPlayerItems is working by tag
		Tag string `json:"tag"`
		// Name        string `json:"name"`
		// ExpLevel    int    `json:"expLevel"`
		// Trophies    int    `json:"trophies"`
		// AttackWins  int    `json:"attackWins"`
		// DefenseWins int    `json:"defenseWins"`
		// Rank        int    `json:"rank"`
		// 	Clan        struct {
		// 		Tag       string `json:"tag"`
		// 		Name      string `json:"name"`
		// 		BadgeUrls struct {
		// 			Small  string `json:"small"`
		// 			Large  string `json:"large"`
		// 			Medium string `json:"medium"`
		// 		} `json:"badgeUrls"`
		// 	} `json:"clan,omitempty"`
		// } `json:"items"`
		// Paging struct {
		// 	Cursors struct {
		// 		After string `json:"after"`
		// 	} `json:"cursors"`
		// } `json:"paging"`
	}
}

func New(client *http.Client, api string) *CoC {
	return &CoC{
		url:    "https://api.clashofclans.com/v1",
		client: client,
		api:    api,
	}
}
