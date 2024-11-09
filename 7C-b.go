//go:build !solution

package httpgauge

import (
	"io"
	"net/http"
	"sort"
	"strconv"
	"sync"

	"github.com/go-chi/chi/v5"
)

type Gauge struct {
	cd map[string]int
}

func New() *Gauge {
	return &Gauge{
		cd: make(map[string]int),
	}
}

func (g *Gauge) Snapshot() map[string]int {
	return g.cd
}

type Cds struct {
	way   string
	total int
}

func (g *Gauge) ServeHTTP(writer http.ResponseWriter, req *http.Request) {
	cds := make([]Cds, 0, len(g.cd))
	for path, cnt := range g.cd {
		cds = append(cds, Cds{path, cnt})
	}
	sort.Slice(cds, func(i, j int) bool {
		return cds[j].way > cds[i].way
	})
	for _, cdn := range cds {
		itoa := strconv.Itoa(cdn.total) + "\n"
		io.WriteString(writer, cdn.way+" "+itoa)
	}
}

func (g *Gauge) Wrap(next http.Handler) http.Handler {
	var error any
	mtx := sync.Mutex{}
	return http.HandlerFunc(func(wr http.ResponseWriter, rq *http.Request) {
		func() {
			defer func() {
				error = recover()
			}()
			next.ServeHTTP(wr, rq)
		}()

		mtx.Lock()

		pat := chi.RouteContext(rq.Context()).RoutePattern()
		g.cd[pat] += 1
		mtx.Unlock()
		if error != nil {
			panic(error)
		}
	})
}
