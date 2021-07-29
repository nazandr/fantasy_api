package store

import (
	"time"

	"github.com/nazandr/fantasy_api/internal/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var TZ = time.FixedZone("UTC+3", +3*60*60)

type SeriesCollection struct {
	Store      *Store
	Collection *mongo.Collection
}

func (c *SeriesCollection) Create(series *models.Series) error {
	_, err := c.Collection.InsertOne(c.Store.context, series)

	if err != nil {
		return err
	}

	return nil
}

func (c *SeriesCollection) UpdateSeries(series models.Series) error {
	_, err := c.Collection.ReplaceOne(c.Store.context, bson.M{"_id": series.ID}, series)

	if err != nil {
		return err
	}

	return nil
}

func (c *SeriesCollection) FindSeries(seriesId int64) (models.Series, error) {
	r := c.Collection.FindOne(c.Store.context, bson.M{"series_id": seriesId})
	series := models.NewSeries()
	if err := r.Decode(series); err != nil {
		return *series, err
	}

	return *series, nil
}

func (c *SeriesCollection) GetAll() ([]models.Series, error) {
	var series []models.Series
	cursor, err := c.Collection.Find(c.Store.context, bson.D{{}})
	if err != nil {
		return nil, err
	}

	for cursor.Next(c.Store.context) {
		var oneSeries models.Series

		if err := cursor.Decode(&oneSeries); err != nil {
			return nil, err
		}
		series = append(series, oneSeries)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	cursor.Close(c.Store.context)

	return series, nil
}

func (c *SeriesCollection) FindByDate(data time.Time) ([]models.Series, error) {
	gd := data.Truncate(24 * time.Hour).In(TZ)
	t := time.Date(data.Year(), data.Month(), data.Day(), 23, 59, 59, data.Nanosecond(), data.Location())
	cursor, err := c.Collection.Find(c.Store.context, bson.M{"date": bson.M{"$gte": gd, "$lte": t}})
	var series []models.Series

	for cursor.Next(c.Store.context) {
		var oneSeries models.Series

		if err := cursor.Decode(&oneSeries); err != nil {
			return nil, err
		}
		series = append(series, oneSeries)
	}

	if cursor.Err() != nil {
		return nil, err
	}

	cursor.Close(c.Store.context)

	return series, nil
}
