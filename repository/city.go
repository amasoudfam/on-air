package repository

import (
	"errors"
	"on-air/models"
	"on-air/server/services"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type City struct {
	APIMockClient *services.APIMockClient
	DB            *gorm.DB
	SyncPeriod    time.Duration
}

func (c *City) SyncCities() {
	ticker := time.NewTicker(c.SyncPeriod)
	done := make(chan bool)

	go func() {
		cities, err := c.APIMockClient.GetCities()
		if err != nil {
			logrus.Error("city_repository_sync_cities:", err)
		} else {
			err := c.StoreCities(cities)
			if err != nil {
				logrus.Error("city_repository_sync_cities:", err)
			}
		}

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				cities, err := c.APIMockClient.GetCities()
				if err != nil {
					logrus.Error("city_repository_sync_cities:", err)
				} else {
					err := c.StoreCities(cities)
					if err != nil {
						logrus.Error("city_repository_sync_cities:", err)
					}
				}
			}
		}
	}()
}

func (c *City) StoreCities(cities []string) error {
	for _, cityName := range cities {
		city := models.City{Name: cityName, CountryID: 1}
		err := c.DB.FirstOrCreate(&city, models.City{Name: cityName}).Error
		if err != nil {
			logrus.Error("city_repository_store_cities:", err)
			return errors.New("failed to store cities")
		}
	}

	return nil
}

func FindCityByName(db *gorm.DB, Name string) (*models.City, error) {
	var city models.City
	err := db.Where("name = ?", Name).First(&city).Error
	if err != nil {
		return nil, err
	}

	return &city, nil
}
