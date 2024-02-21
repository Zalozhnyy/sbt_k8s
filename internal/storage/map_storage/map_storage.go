package mapstorage

import (
	"sync"
	"time"

	"github.com/Zalozhnyy/sbt_k8s/internal/http-server/handlers/objects"
	"github.com/Zalozhnyy/sbt_k8s/internal/storage"
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
