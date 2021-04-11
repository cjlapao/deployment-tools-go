package service

import (
	"github.com/cjlapao/common-go/log"
	"github.com/cjlapao/common-go/version"
	"github.com/cjlapao/deployment-tools-go/executioncontext"
	"github.com/cjlapao/deployment-tools-go/servicebuscli"
)

type ServiceProvider struct {
	Context    *executioncontext.Context
	Version    *version.Version
	Logger     *log.Logger
	ServiceBus *servicebuscli.ServiceBusCli
}

var globalProvider *ServiceProvider

func CreateProvider() *ServiceProvider {
	if globalProvider != nil {
		return globalProvider
	}

	globalProvider = &ServiceProvider{}
	globalProvider.Context = executioncontext.Get()
	globalProvider.Logger = log.Get()
	globalProvider.Version = version.Get()

	return globalProvider
}

func GetServiceProvider() *ServiceProvider {
	return globalProvider
}
