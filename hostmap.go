package hostmap

import (
	"fmt"
	"net/http"
)

type HostMap struct {
	Map map[string]string
	LoadMethod func() (map[string]string, error)
	AddMethod func(v1, v2 string) error
	DelMethod func(v1 string) error
	ModifyMethod func(v1, v2 string) error
}

func New(load func() (map[string]string, error), 
		add func(v1, v2 string) error, 
		del func(v1 string) error, 
		modify func(v1, v2 string) error) *HostMap {

	return &HostMap{
		LoadMethod: load,
		AddMethod: add,
		DelMethod: del,
		ModifyMethod: modify,
	}
}

func (hostMap *HostMap) Load() error {
	m, err := hostMap.LoadMethod()
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
		return hostMap.AddMethod(v1, v2)
	} else {
		return fmt.Errorf("HostMap: Key %s already exists.", v1)
	}
}

func (hostMap *HostMap) Del(v1 string) error {
	_, exists := hostMap.Map[v1]
	if exists {
		delete(hostMap.Map, v1)
		return hostMap.DelMethod(v1)
	} else {
		return fmt.Errorf("HostMap: Key %s does not exist.", v1)
	}
}

func (hostMap *HostMap) Modify(v1, v2 string) error {
	_, exists := hostMap.Map[v1]
	if exists {
		hostMap.Map[v1] = v2
		return hostMap.ModifyMethod(v1, v2)
	} else {
		return fmt.Errorf("HostMap: Key %s does not exist.", v1)
	}
}