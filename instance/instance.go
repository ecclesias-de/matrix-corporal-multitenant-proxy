package instance

import (
	"matrix-corporal-multitenant-proxy/client"
	"matrix-corporal-multitenant-proxy/config"
	"matrix-corporal-multitenant-proxy/store"
)

type Instance struct {
	Store  store.Store
	Config config.Instance
	Client *client.Client
}

type InstanceRepository interface {
	GetInstance(homeServerName string) *Instance
	AddInstance(homeServerName string, instance *Instance)
}

type InMemoryInstanceRepository struct {
	instances map[string]*Instance
}

func (me *InMemoryInstanceRepository) GetInstance(homeServerName string) *Instance {
	return me.instances[homeServerName]
}

func (me *InMemoryInstanceRepository) AddInstance(homeServerName string, instance *Instance) {
	me.instances[homeServerName] = instance
}

func NewInMemoryInstanceRepository() *InMemoryInstanceRepository {
	return &InMemoryInstanceRepository{instances: make(map[string]*Instance)}
}
