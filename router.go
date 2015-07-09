package webengine

import (
	"fmt"
	"strings"
	"regexp"
	"reflect"
	"net/http"
)

type Node struct {
	children map[string]*Node
	name string
	regex *regexp.Regexp
	value interface{}
}

func NewNode() *Node {
	return &Node{ children: make(map[string]*Node) }
}

func (n *Node) Add( key string, value interface{} ){
	n.add( splitKey(key), value )
}

func (n *Node) Find( key string ) interface{} {
	return n.find( splitKey(key) )
}

func (n *Node) find( elements []string ) interface{} {
	el, els := elements[0], elements[1:]

	child, ok := n.children[el]
	if !ok {
		for _, node := range n.children {
			if node.regex != nil && node.regex.MatchString( el ) {
				child = node
				break
			}
		}
	}

	if len(els) == 0 {
		return child
	}
	return child.find( els )
}

func (n *Node) add( elements []string, value interface{} ){
	parent, children := elements[0], elements[1:]

	child, ok := n.children[parent]
	if !ok {
		child = NewNode()
		n.children[parent] = child
	}

	if lenParent := len(parent); string(parent[0]) == "{" && string(parent[lenParent-1]) == "}" {
		cnt := strings.Split(parent[1:lenParent-1], ":")
		child.name = cnt[0]
		var regex string

		if len(cnt) > 1 {
			if cnt[1] == "*" {
				regex = "(.*)"
			} else {
				regex = "("+cnt[1]+")"
			}
		} else {
			regex = "([^/]+)"
		}

		child.regex = regexp.MustCompile(regex)
	}

	if len(children) == 0 {
		child.value = value
		return
	}

	child.add( children, value )
}

func splitKey( key string ) []string {
	keys := strings.Split(key, "/")
	kl := len(keys)
	if keys[kl-1] == "" {
		keys[kl-1] = "/"
	}
	return keys
}

type RouterHandler func( http.ResponseWriter, *http.Request )

type Router struct {
	routes *Node
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	node := router.routes.Find( r.Method + path )
	if reflect.ValueOf(node).IsNil() || node.(*Node).value == nil {
		//route.Handler(w, r)
		fmt.Fprintf(w, "404 Route not found, %q", path)
		return
	}

	route := node.(*Node).value.(Route)
	fmt.Fprintf(w, "Route found, %+v", route )

	ctrl := Registry.NewInstance(route.Controller)
	fmt.Println("Route found, %+v", route, ctrl )
	res  := CallMethod( ctrl, route.Method )
	fmt.Println( res )
}

type Route struct {
	Controller string
	Method string
	//Handler RouterHandler
}

type RouterConfig struct {
	Routes []struct{
		Name string `yaml: name`
		Method string `yaml: method`
		Path string `yaml: path`
		Controller string `yaml: controller`
	}
}

func NewRouter() *Router{
	router := Router{NewNode()}
	routerConfig := RouterConfig{}
	LoadConfig("conf/routes.yaml", &routerConfig)

	for _, route := range routerConfig.Routes {
		action := strings.Split(route.Controller, ".")
		controller := strings.ToUpper( string(action[0][0]) ) + string(action[0][1:])
		method := strings.ToUpper( string(action[1][0]) ) + string(action[1][1:])
		router.routes.Add( strings.ToUpper( route.Method )+route.Path, Route{ controller, method } )
	}

	return &router
}

func CallMethod(i interface{}, methodName string) interface{} {
	var ptr reflect.Value
	var value reflect.Value
	var finalMethod reflect.Value

	value = reflect.ValueOf(i)

	// if we start with a pointer, we need to get value pointed to
	// if we start with a value, we need to get a pointer to that value
	if value.Type().Kind() == reflect.Ptr {
		ptr = value
		value = ptr.Elem()
	} else {
		ptr = reflect.New(reflect.TypeOf(i))
		temp := ptr.Elem()
		temp.Set(value)
	}

	// check for method on value
	method := value.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}
	// check for method on pointer
	method = ptr.MethodByName(methodName)
	if method.IsValid() {
		finalMethod = method
	}

	if (finalMethod.IsValid()) {
		res := finalMethod.Call([]reflect.Value{})
		// check for a return value
		if len(res) == 0 {
			return ""
		}

		return res[0].Interface()
	}

	// return or panic, method not found of either type
	return ""
}
