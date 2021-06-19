package matches

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/nazandr/fantasy_api/internal/app/store"
)

type MatchServer struct {
	Tournaments tournaments
	LastMatch   int64
	Store       *store.Store
}

type tournaments []struct {
	Leagueid   int    `json:"leagueid"`
	LeagueName string `json:"league_name"`
}

type match []struct {
	MatchID     int64  `json:"match_id"`
	StartTime   int    `json:"start_time"`
	RadiantName string `json:"radiant_name"`
	DireName    string `json:"dire_name"`
	Leagueid    int    `json:"leagueid"`
	LeagueName  string `json:"league_name"`
	RadiantWin  bool   `json:"radiant_win"`
}

func NewMatchServer(s *store.Store) *MatchServer {
	return &MatchServer{
		Tournaments: RefreshTournaments(),
		Store:       s,
	}
}

func (m *MatchServer) Serve() error {
	res, err := request("https://api.opendota.com/api/proMatches", "GET")
	if err != nil {
		return err
	}
	matches := match{}
	if err := json.NewDecoder(res.Body).Decode(&matches); err != nil {
		return err
	}

	for i := 0; i < 15; i++ {
		if matches[i].MatchID == m.LastMatch {
			break
		}
		for _, v := range m.Tournaments {
			if v.Leagueid == matches[i].Leagueid {
				newMatch := models.NewMatch()
				newMatch.CalcPoints(matches[i].MatchID)
				if err := m.Store.Matches().Create(newMatch); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func RefreshTournaments() tournaments {
	file, _ := ioutil.ReadFile("configs/tournaments.json")
	t := tournaments{}
	json.Unmarshal([]byte(file), &t)
	return t
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
