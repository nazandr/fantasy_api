package matches

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"github.com/nazandr/fantasy_api/internal/app/server"
	"github.com/nazandr/fantasy_api/internal/app/store"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	tz = time.FixedZone("UTC+3", +3*60*60)
)

type MatchServer struct {
	LastMatch int64
	Store     *store.Store
	Server    *server.APIServer
}

type match []struct {
	MatchID     int64  `json:"match_id"`
	StartTime   int    `json:"start_time"`
	RadiantName string `json:"radiant_name"`
	DireName    string `json:"dire_name"`
	Leagueid    int    `json:"leagueid"`
	LeagueName  string `json:"league_name"`
	SeriesId    int64  `json:"series_id"`
	RadiantWin  bool   `json:"radiant_win"`
}

func NewMatchServer(s *server.APIServer) *MatchServer {
	return &MatchServer{
		Store:  s.Store,
		Server: s,
	}
}

func (m *MatchServer) Ticker() {
	if err := m.Serve(); err != nil {
		m.Server.Logger.Info(err)
	}
	m.Server.Logger.Info("served")
	ticker := time.NewTicker(1 * time.Hour)
	for range ticker.C {
		if err := m.Serve(); err != nil {
			m.Server.Logger.Info(err)
		}
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

	for i := 0; i < len(matches); i++ {
		t := parseTime(matches[i].StartTime).In(tz)
		// m.Server.Logger.Info(m.Store.PlayerCards().IsTeam(matches[i].RadiantName))
		// m.Server.Logger.Info(matches[i].RadiantName)
		// m.Server.Logger.Info(m.Store.PlayerCards().IsTeam(matches[i].DireName))
		// m.Server.Logger.Info(matches[i].DireName)
		if !m.Store.PlayerCards().IsTeam(matches[i].RadiantName) && !m.Store.PlayerCards().IsTeam(matches[i].DireName) {
			continue
		}

		series, err := m.Store.Series().FindSeries(matches[i].SeriesId)
		if err == mongo.ErrNoDocuments {
			ser := models.NewSeries()
			ser.Teams = append(ser.Teams, matches[i].RadiantName)
			ser.Teams = append(ser.Teams, matches[i].DireName)
			ser.Date = t
			ser.SeriesId = matches[i].SeriesId
			newMatch := models.NewMatch()
			newMatch.Teams = append(newMatch.Teams, matches[i].RadiantName)
			newMatch.Teams = append(newMatch.Teams, matches[i].DireName)
			newMatch.CalcPoints(matches[i].MatchID)
			ser.Matches = append(ser.Matches, *newMatch)

			if err := m.Store.Series().Create(ser); err != nil {
				return err
			}
			continue
		} else if err != nil {
			return err
		}

		st := false
		for _, v := range series.Matches {
			if v.MatchID == matches[i].MatchID {
				st = true
			}
		}
		if !st {
			newMatch := models.NewMatch()
			newMatch.Teams = append(newMatch.Teams, matches[i].RadiantName)
			newMatch.Teams = append(newMatch.Teams, matches[i].DireName)
			newMatch.CalcPoints(matches[i].MatchID)
			series.Matches = append(series.Matches, *newMatch)
			if err := m.Store.Series().UpdateSeries(series); err != nil {
				return err
			}
		}

	}

	return nil
}

func parseTime(unix int) time.Time {
	i, _ := strconv.ParseInt(strconv.Itoa(unix), 10, 64)
	return time.Unix(i, 0)
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
