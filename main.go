package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"github.com/pathakanu/loadbalancer/server"
)



// simpleServer holds the information about a backend server.
type simpleServer struct {
	addr  string
	proxy *httputil.ReverseProxy
}

// newSimpleServer creates a new backend server instance.
func newSimpleServer(addr string) *simpleServer {
	serverURL, err := url.Parse(addr)
	if err != nil {
		panic(err)
	}

	return &simpleServer{
		addr:  addr,
		proxy: httputil.NewSingleHostReverseProxy(serverURL),
	}
}

func (s *simpleServer) Address() string { return s.addr }

// IsAlive checks if the backend server is alive. For this example, we'll always assume it is.
func (s *simpleServer) IsAlive() bool { return true }

// Serve forwards the request to the backend server.
func (s *simpleServer) Serve(w http.ResponseWriter, r *http.Request) {
	s.proxy.ServeHTTP(w, r)
}

// LoadBalancer holds the state for our load balancer.
type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []server.Server
	mu              sync.Mutex
}

// NewLoadBalancer creates a new LoadBalancer.
func NewLoadBalancer(port string, servers []server.Server) *LoadBalancer {
	return &LoadBalancer{
		port:    port,
		servers: servers,
	}
}

// getNextAvailableServer finds the next available server to send a request to, using a round-robin algorithm.
func (lb *LoadBalancer) getNextAvailableServer() server.Server {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	return server
}

// serveProxy is the handler for all incoming requests.
func (lb *LoadBalancer) serveProxy(w http.ResponseWriter, r *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding request to address %q\n", targetServer.Address())
	targetServer.Serve(w, r)
}

func main() {
	// Create our list of backend servers.
	servers := []server.Server{
		newSimpleServer("https://www.facebook.com"),
		newSimpleServer("https://www.bing.com")}

	// Create the load balancer.
	lb := NewLoadBalancer("8000", servers)

	handleRedirect := func(w http.ResponseWriter, r *http.Request) {
		lb.serveProxy(w, r)
	}

	http.HandleFunc("/", handleRedirect)

	// // Start a backend server on port 8081.
	// go func() {
	// 	http.ListenAndServe(":8081", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "Hello from Backend Server 1")
	// 	}))
	// }()

	// // Start another backend server on port 8082.
	// go func() {
	// 	http.ListenAndServe(":8082", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintln(w, "Hello from Backend Server 2")
	// 	}))
	// }()

	// // Set up the handler for the load balancer.
	// http.HandleFunc("/", lb.serveProxy)

	fmt.Printf("Load Balancer started at :%s\n", lb.port)
	// Start the load balancer server.
	http.ListenAndServe(":"+lb.port, nil)
}
