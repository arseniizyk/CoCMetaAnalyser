package coc

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

func (c *CoC) GetPlayerInfo(playerTag string) (*Player, error) {
	if playerTag == "" {
		return nil, fmt.Errorf("player tag is empty")
	}

	if !strings.HasPrefix(playerTag, "#") {
		playerTag = "#" + playerTag
	}

	playerTag = url.PathEscape(playerTag)
	endpoint := fmt.Sprintf("%s/players/%s", c.url, playerTag)

	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", "Bearer "+c.api)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	var player Player
	if err := parseResponse(resp, &player); err != nil {
		return nil, err
	}

	return &player, nil
}
