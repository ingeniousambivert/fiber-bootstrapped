package core

import (
	"fmt"
)

type Services map[string]*Service
type App struct {
	Services Services
}

func InitApp() *App {
	return &App{
		Services: make(Services),
	}
}

func (a *App) InitServices(services Services) {
	a.Services = services
}

func (a *App) Service(name string) (*Service, error) {
	service, ok := a.Services[name]
	if !ok {
		return nil, fmt.Errorf("service %s not found", name)
	}
	return service, nil
}
