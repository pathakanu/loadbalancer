# Go Load Balancer

This is a simple HTTP load balancer written in Go. It distributes incoming requests across a set of backend servers using a round-robin algorithm.

## How to Run

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/pathakanu/loadbalancer.git
    cd loadbalancer
    ```

2.  **Run the load balancer:**
    ```bash
    go run main.go
    ```

    The load balancer will start on port 8000.

3.  **Send requests to the load balancer:**
    You can use `curl` or any other HTTP client to send requests to the load balancer.

    ```bash
    curl http://localhost:8000
    ```

## Backend Servers

The backend servers are defined in the `main.go` file. By default, the load balancer is configured to use the following backend servers:

*   `https://www.facebook.com`
*   `https://www.bing.com`

You can modify the `servers` slice in the `main` function to add or remove backend servers.

## Load Balancer Logic

The load balancer uses a simple round-robin algorithm to select the next backend server to forward a request to. It maintains a counter that is incremented for each request. The server is selected by taking the counter modulo the number of available servers.

The load balancer also checks if a server is "alive" before forwarding a request to it. In this example, the `IsAlive()` function always returns `true`, but it can be extended to implement a more sophisticated health check mechanism.
