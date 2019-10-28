package service_manager

import (
	"errors"
	. "github.com/saichler/utils/golang"
	"plugin"
)

func (sm *ServiceManager) AddService(service IService, commands IServiceCommands, handlers IServiceMessageHandlers) {
	container, ok := sm.containers.Get(service.Topic()).(*ServiceContainer)
	if !ok {
		container = NewServiceContainer(service.Topic())
		sm.containers.Put(service.Topic(), container)
	}
	container.AddService(service, handlers, sm)
	for _, cmd := range commands.Commands() {
		sm.console.RegisterCommand(cmd)
	}
}

func (sm *ServiceManager) LoadService(filename string) error {
	servicePlugin, e := plugin.Open(filename)
	if e != nil {
		return Error("Unable to load serivce plugin:", e)
	}
	service, e := getServiceFromPlugin(servicePlugin)
	if e != nil {
		return Error(e.Error())
	}
	commands, e := getCommandsFromPlugin(servicePlugin)
	if e != nil {
		return Error(e.Error())
	}
	handlers, e := getHandlersFromPlugin(servicePlugin)
	if e != nil {
		return Error(e.Error())
	}

	sm.AddService(service, commands, handlers)

	return nil
}

func getServiceFromPlugin(plugin *plugin.Plugin) (IService, error) {
	svr, e := plugin.Lookup("Service")
	if e != nil {
		return nil, e
	}

	ptr, ok := svr.(*IService)
	if !ok {
		msg := "Service is not of type IService, please check that it implements IService and that Service is a pointer."
		return nil, errors.New(msg)
	}
	return *ptr, nil
}

func getCommandsFromPlugin(plugin *plugin.Plugin) (IServiceCommands, error) {
	svr, e := plugin.Lookup("Commands")
	if e != nil {
		return nil, e
	}

	ptr, ok := svr.(*IServiceCommands)
	if !ok {
		msg := "Commands is not of type IServiceCommands, please check that it implements IServiceCommands and that Commands is a pointer."
		return nil, errors.New(msg)
	}
	return *ptr, nil
}

func getHandlersFromPlugin(plugin *plugin.Plugin) (IServiceMessageHandlers, error) {
	svr, e := plugin.Lookup("Handlers")
	if e != nil {
		return nil, e
	}

	ptr, ok := svr.(*IServiceMessageHandlers)
	if !ok {
		msg := "Commands is not of type IServiceMessageHandlers, please check that it implements IServiceMessageHandlers and that Handlers is a pointer."
		return nil, errors.New(msg)
	}
	return *ptr, nil
}
