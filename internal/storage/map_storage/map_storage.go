package mapstorage

import (
	"context"
	"log/slog"
	"sync"
	"time"

	"github.com/Zalozhnyy/sbt_k8s/internal/http-server/handlers/objects"
	"github.com/Zalozhnyy/sbt_k8s/internal/storage"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	jsonsWithoutTime = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sbt_k8s_number_of_json_without_expire_time",
		Help: "Total number of messages without expire time",
	})
	totalNumberOfSavedJsons = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "sbt_k8s_number_of_jsons",
		Help: "Total number of messages",
	})
)

type MapStorage struct {
	m    *sync.RWMutex
	data map[string]objects.JsonDto
}

func New() *MapStorage {
	return &MapStorage{
		m:    &sync.RWMutex{},
		data: make(map[string]objects.JsonDto),
	}
}

func (m *MapStorage) Save(id string, obj objects.JsonDto) error {
	if obj.ExpiredTime.IsZero() {
		jsonsWithoutTime.Inc()
	}
	m.m.Lock()
	m.data[id] = obj
	m.m.Unlock()

	return nil
}

func (m *MapStorage) Get(id string) (objects.JsonDto, error) {
	m.m.RLock()
	js, ok := m.data[id]
	m.m.RUnlock()
	if !ok {
		return objects.JsonDto{}, storage.ErrDoNotExists
	}

	if !js.ExpiredTime.IsZero() && time.Now().Before(js.ExpiredTime) {
		return objects.JsonDto{}, storage.ErrExpired
	}

	return js, nil
}

func cleanMap(storage *MapStorage) {

	storage.m.Lock()
	defer storage.m.Unlock()

	for key, js := range storage.data {
		if !js.ExpiredTime.IsZero() && time.Now().After(js.ExpiredTime) {
			delete(storage.data, key)
		}
	}

}

func MapCleaner(ctx context.Context, log *slog.Logger, storage *MapStorage) {
	ticker := time.NewTicker(time.Minute * 10)
	defer ticker.Stop()
LOOP:
	for {
		select {
		case <-ctx.Done():
			log.Info("map cleaner stopped")
			break LOOP
		case <-ticker.C:
			log.Info("start clean map storage")
			cleanMap(storage)
			log.Info("cleaning finish")
		}
	}
}

func GetLenOfMapStorage(ctx context.Context, log *slog.Logger, storage *MapStorage) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

LOOP:
	for {
		select {
		case <-ctx.Done():
			break LOOP
		case <-ticker.C:
			totalNumberOfSavedJsons.Set(float64(len(storage.data)))
		}
	}

}
