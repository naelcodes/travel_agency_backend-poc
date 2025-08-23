package managers

import (
	"context"
	"fmt"
	"reflect"

	"neema.co.za/rest/utils/logger"
)

type Dependency func(context.Context) (any, error)

type DependencyManager struct {
	dependencies map[string]Dependency
}

func NewDependencyManager() *DependencyManager {
	return &DependencyManager{
		dependencies: make(map[string]Dependency),
	}
}

func (d *DependencyManager) Get(key string) Dependency {
	return d.dependencies[key]
}

func (d *DependencyManager) Add(moduleExports ...any) {
	for i := 0; i < len(moduleExports); i++ {
		ModuleReflectedValue := reflect.ValueOf(moduleExports[i])
		moduleReflectedType := ModuleReflectedValue.Type()
		for j := 0; j < moduleReflectedType.NumMethod(); j++ {
			method := moduleReflectedType.Method(j)
			logger.Info(fmt.Sprintf(" Exported function Name: %v", method.Name))
			d.dependencies[method.Name] = createFunction(ModuleReflectedValue.MethodByName(method.Name))

		}
	}
}

func (d *DependencyManager) GetAll() map[string]Dependency {
	return d.dependencies
}

func createFunction(reflectedMethod reflect.Value) Dependency {
	return func(context context.Context) (any, error) {
		rResults := reflectedMethod.Call([]reflect.Value{reflect.ValueOf(context)})
		if rResults[1].Interface() != nil {
			return rResults[0].Interface(), rResults[1].Interface().(error)
		}
		return rResults[0].Interface(), nil

	}
}
