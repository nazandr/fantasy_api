package playercards

import (
	"math/rand"
	"time"

	"github.com/nazandr/fantasy_api/internal/app/store"
)

type Pack struct {
	Cards [5]store.PlayerCard
}

func NewPack() *Pack {
	return &Pack{}
}

func (p *Pack) OpenCommonPack(s *store.Store) error {
	it := []int{0, 1, 2, 3}
	w := []float32{0.7, 0.23, 0.05, 0.2}
	c := false

	players, err := s.PlayerCards().GetAll()
	if err != nil {
		return err
	}
	for i := 0; i < 5; i++ {
		p.Cards[i] = players[rand.Intn(len(players))]
		p.Cards[i].Rarity = randomWithProbabilitis(it, w)
		if p.Cards[i].Rarity > 0 {
			c = true
		}
	}
	if !c {
		p.Cards[rand.Intn(3)].Rarity = 1
	}
	return nil
}

func randomWithProbabilitis(items []int, weights []float32) int {
	rand.Seed(time.Now().UnixNano())
	var (
		sumWeight []float32
		sum       float32
	)
	for _, v := range weights {
		sum += v
		sumWeight = append(sumWeight, sum)
	}

	ri := rand.Float32()

	for i, v := range sumWeight {
		if v >= ri {
			return items[i]
		}
	}
	return 5
}
