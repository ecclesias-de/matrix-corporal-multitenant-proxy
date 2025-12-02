package container

import (
	"devture-matrix-corporal/corporal/httphelp"
	"matrix-corporal-multitenant-proxy/api"
	"matrix-corporal-multitenant-proxy/api/handler"
	"matrix-corporal-multitenant-proxy/client"
	"matrix-corporal-multitenant-proxy/config"
	"matrix-corporal-multitenant-proxy/instance"
	"matrix-corporal-multitenant-proxy/reconcile"
	"matrix-corporal-multitenant-proxy/store"
	"time"
)

type Service interface {
	Start() error
	Stop() error
}

func Bootstrap(config config.Config) []Service {
	instanceRepository := instance.NewInMemoryInstanceRepository()

	services := make([]Service, 0)

	for _, instanceConfig := range config.Instances {
		instance := &instance.Instance{
			Store:  store.BootstrapPersistentStore(config.StoragePath, instanceConfig.HomeServerName),
			Config: instanceConfig,
			Client: client.NewClient(instanceConfig),
		}
		instanceRepository.AddInstance(instanceConfig.HomeServerName, instance)

		reconciler := reconcile.NewReconciler(instance, time.Millisecond*time.Duration(config.ReconcileRetryInterval))
		services = append(services, reconciler)
	}

	policyHandler := handler.NewPolicyHandler(instanceRepository)

	server := api.NewServer(config.ListenAddress, []httphelp.HandlerRegistrator{policyHandler})
	services = append(services, server)

	return services
}
