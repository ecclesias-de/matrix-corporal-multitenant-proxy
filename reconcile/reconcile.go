package reconcile

import (
	"devture-matrix-corporal/corporal/policy"
	"matrix-corporal-multitenant-proxy/client"
	"matrix-corporal-multitenant-proxy/instance"
	"matrix-corporal-multitenant-proxy/store"
	"time"

	"github.com/sirupsen/logrus"
)

type Reconciler struct {
	updateChanel  chan *store.UpdateEvent
	client        client.Client
	retryInterval time.Duration
	ticker        *time.Ticker
	stop          chan any
}

func NewReconciler(instance *instance.Instance, retryInterval time.Duration) *Reconciler {
	return &Reconciler{
		updateChanel:  instance.Store.GetUpdateChannel(),
		client:        *instance.Client,
		retryInterval: retryInterval,
		stop:          make(chan any),
	}
}

func (me *Reconciler) Start() error {
	me.ticker = time.NewTicker(me.retryInterval)
	me.ticker.Stop()

	go me.listen()

	return nil
}

func (me *Reconciler) Stop() error {
	me.ticker.Stop()
	me.stop <- true

	return nil
}

func (me *Reconciler) listen() {
	var policy *policy.Policy

	for {
		select {
		case <-me.stop:
			return
		case <-me.ticker.C:
		case updateEvent, active := <-me.updateChanel:
			logrus.Debugf(`Reconciling outgoing policy`)
			if !active {
				return
			}

			policy = BuildOutgoingPolicy(updateEvent)
			me.ticker.Reset(me.retryInterval)
		}

		err := me.client.UpdatePolicy(policy)
		if err == nil {
			me.ticker.Stop()
		} else {
			logrus.Debugf(`Failed putting outgoing policy: %s`, err)
		}
	}
}
