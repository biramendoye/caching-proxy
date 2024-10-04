package server

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func (srv *Server) routes() http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("/", srv.HandleRequests)

	return srv.logRequests()(router)
}

func (srv *Server) HandleRequests(w http.ResponseWriter, r *http.Request) {
	url := fmt.Sprintf("%s%s", srv.config.forwardedURL, r.RequestURI)
	v, ok := srv.cache.Get(url)
	if ok {
		w.Header().Add("X-Cache", "HIT")
		w.Write(v.([]byte))
		return
	}

	srv.logger.Info("forward request", "from", r.RequestURI, "to", url)
	proxyReq, err := http.NewRequest(r.Method, url, r.Body)
	if err != nil {
		srv.logger.Error(err.Error())
		return
	}

	proxyReq.Header = make(http.Header)
	for h, val := range r.Header {
		proxyReq.Header[h] = val
	}
	copyHeader(proxyReq.Header, r.Header)

	client := &http.Client{}
	resp, err := client.Do(proxyReq)
	if err != nil {
		srv.logger.Error(err.Error())
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		srv.logger.Error(err.Error())
		return
	}

	srv.cache.Set(url, body, time.Hour)

	copyHeader(w.Header(), resp.Header)
	w.Header().Add("X-Cache", "MISS")
	w.Write(body)

}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
