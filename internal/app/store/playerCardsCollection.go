package store

import (
	"math/rand"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type PlayerCardsCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

type Pack struct {
	Cards []models.PlayerCard
}

func NewPack() *Pack {
	return &Pack{
		Cards: make([]models.PlayerCard, 5),
	}
}

func (c *PlayerCardsCollection) GetAll() ([]models.PlayerCard, error) {
	var playersList []models.PlayerCard
	cursor, err := c.Collection.Find(c.Store.context, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(c.Store.context) {
		var player models.PlayerCard

		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}
		playersList = append(playersList, player)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	cursor.Close(c.Store.context)

	return playersList, nil
}

func (c *PlayerCardsCollection) OpenCommonPack() (*Pack, error) {
	it := []int{0, 1, 2, 3}
	w := []float32{0.7, 0.23, 0.05, 0.2}
	rareCard := false

	players, err := c.GetAll()
	if err != nil {
		return nil, err
	}

	p := NewPack()

	for i := 0; i < 5; i++ {
		p.Cards[i] = players[rand.Intn(len(players))]
		p.Cards[i].Rarity = models.RandomWithProbabilitis(it, w)
		p.Cards[i].Id = primitive.NewObjectID()
		if p.Cards[i].Rarity > 0 {
			rareCard = true
			buffs := models.SetBuffs(p.Cards[i].Rarity)
			p.Cards[i].Buffs = buffs
		}
	}
	if !rareCard {
		n := rand.Intn(4)
		p.Cards[n].Rarity = 1
		buffs := models.SetBuffs(p.Cards[n].Rarity)
		p.Cards[n].Buffs = buffs
	}

	// i := reflect.ValueOf(fp).Elem().FieldByName(n.NameOfFild).Float()
	// fmt.Println(i * float64(n.Multiplaier))

	return p, nil
}
