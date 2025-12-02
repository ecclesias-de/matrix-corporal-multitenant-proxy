package store

import (
	"devture-matrix-corporal/corporal/policy"
	"encoding/json"
	"errors"
	"io"
	"matrix-corporal-multitenant-proxy/mtp_policy"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
)

type PersistentStore struct {
	InMemoryStore
	path string
}

type persistedPolicy struct {
	BasePolicy       *policy.Policy                        `json:"basePolicy"`
	ExchangePolicies map[string]*mtp_policy.ExchangePolicy `json:"exchangePolicies"`
}

func BootstrapPersistentStore(storagePath string, tenant string) *PersistentStore {
	path := filepath.Join(storagePath, tenant+".json")

	inMemoryStore, err := load(path)
	if err != nil {
		logrus.Warn(err)
	}
	if inMemoryStore == nil {
		inMemoryStore = NewInMemoryStore()
	}

	store := &PersistentStore{
		InMemoryStore: *inMemoryStore,
		path:          path,
	}

	store.start()

	return store
}

func (me *PersistentStore) start() {
	updateChan := me.GetUpdateChannel()

	go func() {
		for {
			event, active := <-updateChan
			if !active {
				break
			}

			if err := me.persist(event); err != nil {
				logrus.Warn(err)
			}
		}
	}()
}

func (me *PersistentStore) persist(event *UpdateEvent) error {
	policyJson, err := json.Marshal(persistedPolicy{
		BasePolicy:       event.BasePolicies,
		ExchangePolicies: me.exchangePolicies,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(me.path, policyJson, 0600); err != nil {
		return err
	}

	return nil
}

func load(path string) (*InMemoryStore, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	policyJson, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	var policy persistedPolicy
	if err := json.Unmarshal(policyJson, &policy); err != nil {
		return nil, err
	}

	if policy.ExchangePolicies == nil {
		policy.ExchangePolicies = make(map[string]*mtp_policy.ExchangePolicy)
	}

	return &InMemoryStore{
		basePolicy:       policy.BasePolicy,
		exchangePolicies: policy.ExchangePolicies,
		channels:         make([]chan *UpdateEvent, 0),
	}, nil
}
