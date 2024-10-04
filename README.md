# caching-proxy

A CLI tool that starts a caching proxy server, it will forward requests to the actual server and cache the responses. If the same request is made again, it will return the cached response instead of forwarding the request to the server.

For detailed project instructions, please visit: [Caching Proxy](https://roadmap.sh/projects/caching-server).

Features

- Request Forwarding: Forwards incoming requests to the specified origin server.
- Response Caching: Caches responses for faster retrieval on subsequent requests.
- Cache Management: Automatically clear cached data after 1hour.
- Configurable Server: Set port and origin server from the command line.

## Installation

Clone the repository and build the tool.

```bash
git clone https://github.com/username/caching-proxy.git
cd caching-proxy
go build -o caching-proxy
```

## Usage

```bash
caching-proxy --port 3000 --origin http://dummyjson.com
```
