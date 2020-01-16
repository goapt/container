package container

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

type key struct {
	name string
	rt   reflect.Type
}

// Container
type Container struct {
	instance sync.Map
}

// NewContainer
func New() *Container {
	return &Container{}
}

// Register will bind an abstraction to container
func (p *Container) Register(resolver interface{}) {
	rt := reflect.TypeOf(resolver)
	if rt == nil {
		panic("the resolver cant't an untyped nil")
	}

	if rt.Kind() != reflect.Func {
		panic("the resolver must be a function")
	}

	for i := 0; i < rt.NumOut(); i++ {
		p.instance.Store(p.key(rt.Out(i)), resolver)
	}
}

// arguments will return resolved arguments of the given function.
func (p *Container) arguments(function interface{}) []reflect.Value {
	rt := reflect.TypeOf(function)
	argumentsCount := rt.NumIn()
	arguments := make([]reflect.Value, argumentsCount)

	for i := 0; i < argumentsCount; i++ {
		abstraction := rt.In(i)
		var instance interface{}

		if concrete, ok := p.instance.Load(p.key(abstraction)); ok {
			instance = p.invoke(concrete)
		} else {
			panic("no concrete found for the abstraction: " + abstraction.String())
		}

		arguments[i] = reflect.ValueOf(instance)
	}

	return arguments
}

func (p *Container) key(rt reflect.Type) key {
	return key{
		name: "",
		rt:   rt,
	}
}

// invoke will call the given function and return its returned value.
// It only works for functions that return a single value.
func (p *Container) invoke(function interface{}) interface{} {
	return reflect.ValueOf(function).Call(p.arguments(function))[0].Interface()
}

// Make will resolve the dependency and return a appropriate concrete of the given abstraction.
// It can take an abstraction (interface reference) and fill it with the related implementation.
// It also can takes a function (receiver) with one or more arguments of the abstractions (interfaces) that need to be
// resolved, Container will invoke the receiver function and pass the related implementations.
func (p *Container) Make(receiver interface{}) []reflect.Value {
	rt := reflect.TypeOf(receiver)
	if rt == nil {
		panic("cannot detect type of the receiver, make sure your are passing reference of the object")
	}

	if rt.Kind() == reflect.Ptr {
		abstraction := reflect.TypeOf(receiver).Elem()

		if concrete, ok := p.instance.Load(p.key(abstraction)); ok {
			instance := p.invoke(concrete)
			reflect.ValueOf(receiver).Elem().Set(reflect.ValueOf(instance))
			return nil
		}
	}

	if reflect.TypeOf(receiver).Kind() == reflect.Func {
		ret := reflect.ValueOf(receiver).Call(p.arguments(receiver))
		return ret
	}

	panic("the receiver must be either a reference or a callback")
}

// String is return all contaniner instance
func (p *Container) String() string {
	lines := make([]string, 0)
	lines = append(lines, "container:")
	p.instance.Range(func(k, value interface{}) bool {
		ks := k.(key)
		line := fmt.Sprintf(`  {name:"%s",rt:%#v}: %#v`, ks.name, ks.rt.String(), value)
		lines = append(lines, line)

		return true
	})
	return strings.Join(lines, "\n")
}
