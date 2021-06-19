package models

import (
	"encoding/json"
	"net/http"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Match struct {
	ID      primitive.ObjectID `bson:"_id" json:"_id"`
	Teams   []string           `bson:"teams" json:"teams"`
	MatchID int64              `bson:"match_id"`
	Points  []Points
}
type Points struct {
	PlayerID      int     `bson:"account_id"`
	Total         float32 `bson:"name"`
	Kills         float32 `bson:"kills"`
	Deaths        float32 `bson:"deaths"`
	Assists       float32 `bson:"assists"`
	LastHits      float32 `bson:"last_hits"`
	Gpm           float32 `bson:"gold_per_min"`
	TowerKills    int     `bson:"tower_kills"`
	RoshanKils    int     `bson:"roshan_kills"`
	Participation float32 `bson:"teamfight_participation"`
	Observers     float32 `bson:"observer_uses"`
	CampStacked   float32 `bson:"camps_stacked"`
	Runs          float32 `bson:"rune_pickups"`
	FirsBlood     int     `bson:"firstblood_claimed"`
	Stuns         float32 `bson:"stuns"`
}

type player struct {
	SteamId       int     `json:"account_id"`
	Name          string  `json:"name"`
	Kills         int     `json:"kills"`
	Deaths        int     `json:"deaths"`
	Assists       int     `json:"assists"`
	LastHits      int     `json:"last_hits"`
	Gpm           int     `json:"gold_per_min"`
	TowerKills    int     `json:"tower_Kills"`
	RoshanKils    int     `json:"roshan_Kills"`
	Participation float32 `json:"teamfight_participation"`
	Observers     int     `json:"observer_uses"`
	CampStacked   int     `json:"camps_stacked"`
	Runs          int     `json:"rune_pickups"`
	FirsBlood     int     `json:"firstblood_claimed"`
	Stuns         float32 `json:"stuns"`
}

type match struct {
	MatchId int      `json:"match_id"`
	Players []player `json:"players"`
}

func NewMatch() *Match {
	return &Match{
		ID:      primitive.NewObjectID(),
		Teams:   []string{},
		MatchID: 0,
		Points:  make([]Points, 10),
	}
}

func (m *Match) CalcPoints(matchID int64) {
	res, err := request("https://api.opendota.com/api/matches/"+string(matchID), "GET")
	if err != nil {
		return
	}

	match := match{}
	if err := json.NewDecoder(res.Body).Decode(&match); err != nil {
		return
	}

	for i, v := range match.Players {
		m.Points[i].PlayerID = v.SteamId
		m.Points[i].Kills = (float32(v.Kills) * 0.3)
		m.Points[i].Deaths = float32(v.Deaths) * -0.3
		m.Points[i].Assists = float32(v.Assists) * 0.15
		m.Points[i].LastHits = float32(v.LastHits) * 0.003
		m.Points[i].Gpm = float32(v.Gpm) * 0.002
		m.Points[i].TowerKills = v.TowerKills * 1
		m.Points[i].RoshanKils = v.RoshanKils * 1
		m.Points[i].Participation = float32(v.Participation) * 3
		m.Points[i].Observers = float32(v.Observers) * 0.5
		m.Points[i].CampStacked = float32(v.CampStacked) * 0.5
		m.Points[i].Runs = float32(v.Runs) * 0.25
		m.Points[i].FirsBlood = v.FirsBlood * 4
		m.Points[i].Stuns = float32(v.Stuns) * 0.05
		m.Points[i].Total = (m.Points[i].Kills + m.Points[i].Deaths + m.Points[i].Assists +
			m.Points[i].LastHits + m.Points[i].Gpm + float32(m.Points[i].TowerKills+
			m.Points[i].RoshanKils) + m.Points[i].Participation + m.Points[i].Observers +
			m.Points[i].CampStacked + m.Points[i].Runs + float32(m.Points[i].FirsBlood) +
			m.Points[i].Stuns)
	}
}

func request(url string, method string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	c := &http.Client{}
	req.Header.Add("Content-Type", "application/json")
	res, err := c.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
