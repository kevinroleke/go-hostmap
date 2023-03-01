package hostmap

import (
	"fmt"
	"net/http"
	"database/sql"
)

type HostMap struct {
	Map map[string]string
	DB *sql.DB
	LoadMethod func(db *sql.DB) (map[string]string, error)
	AddMethod func(db *sql.DB, v1, v2 string) error
	DelMethod func(db *sql.DB, v1 string) error
	ModifyMethod func(db *sql.DB, v1, v2 string) error
}

func New(db *sql.DB, 
		load func() (map[string]string, error), 
		add func(v1, v2 string) error, 
		del func(v1 string) error, 
		modify func(v1, v2 string) error) *HostMap {

	return &HostMap{
		DB: db,
		LoadMethod: load,
		AddMethod: add,
		DelMethod: del,
		ModifyMethod: modify,
	}
}

func (hostMap *HostMap) Load() error {
	m, err := hostMap.LoadMethod(hostMap.DB)
	if err != nil {
		return err
	}

	hostMap.Map = m
	return err
}

func (hostMap *HostMap) Get(r *http.Request) (string, error) {
	val, ok := hostMap.Map[r.Host]
	if !ok {
		return "", fmt.Errorf("HostMap: Host %s does not exist in the map.", r.Host)
	} else {
		return val, nil
	}
}

func (hostMap *HostMap) Add(v1, v2 string) error {
	_, exists := hostMap.Map[v1]
	if !exists {
		hostMap.Map[v1] = v2
		return hostMap.AddMethod(hostMap.DB, v1, v2)
	} else {
		return fmt.Errorf("HostMap: Key %s already exists.", v1)
	}
}

func (hostMap *HostMap) Del(v1 string) error {
	_, exists := hostMap.Map[v1]
	if exists {
		delete(hostMap.Map, v1)
		return hostMap.DelMethod(hostMap.DB, v1)
	} else {
		return fmt.Errorf("HostMap: Key %s does not exist.", v1)
	}
}

func (hostMap *HostMap) Modify(v1, v2 string) error {
	_, exists := hostMap.Map[v1]
	if exists {
		hostMap.Map[v1] = v2
		return hostMap.ModifyMethod(hostMap.DB, v1, v2)
	} else {
		return fmt.Errorf("HostMap: Key %s does not exist.", v1)
	}
}