package middleware

import (
	"fmt"
	"net/http"
)

type RouteMux struct {
	*http.ServeMux
	routes []string
}

func NewRouteMux() *RouteMux {
	return &RouteMux{
		ServeMux: http.NewServeMux(),
		routes:   make([]string, 0),
	}
}

func (m *RouteMux) Handle(pattern string, handler http.Handler) {
	m.routes = append(m.routes, pattern)
	m.ServeMux.Handle(pattern, handler)
}

func (m *RouteMux) HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request)) {
	m.routes = append(m.routes, pattern)
	m.ServeMux.HandleFunc(pattern, handler)
}

func (m *RouteMux) PrintRoutes(ip string, port string) {
	fmt.Printf("\nServer Information:\n")
	fmt.Printf("==================\n")
	fmt.Printf("Local IP: %s\n", ip)
	fmt.Printf("Port: %s\n", port)
	fmt.Printf("==================\n")
	fmt.Printf("\nAvailable Routes:\n")
	fmt.Printf("==================\n")
	for _, route := range m.routes {
		fmt.Printf("• http://%s%s%s\n", ip, port, route)
	}
	fmt.Printf("==================\n")
}
