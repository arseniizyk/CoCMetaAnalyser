package coc

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func (c *CoC) GetPlayerInfo(playerTag string) (*Player, error) {
	if string(playerTag[0]) != "#" {
		playerTag = "#" + playerTag
	}

	playerTag = url.PathEscape(playerTag)
	url := fmt.Sprintf(c.url+"/players/%s", playerTag)

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

	var player Player
	err = json.NewDecoder(resp.Body).Decode(&player)
	if err != nil {
		return nil, err
	}

	return &player, nil
}
