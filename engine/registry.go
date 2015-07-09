package webengine

import (
	"reflect"
	"fmt"
)

var Registry RegistryObj

type RegistryObj struct {
	data map[string]reflect.Type
}

func (r *RegistryObj) Register( name string, val interface{} ){
	fmt.Println("Registering " + name)
	r.data[name] = reflect.TypeOf(val)
}

func (r *RegistryObj) NewInstance( name string ) interface{}{
	val, ok := r.data[name]
	if !ok {
		panic("Key doesn't exist "+name)
	}

	v := reflect.New(val).Elem()
	return v.Interface()
}

func init(){
	Registry = RegistryObj{ make(map[string]reflect.Type) }
}
