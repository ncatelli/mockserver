package logging

import (
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

// Middleware is a latency middleware that causes a response to sleep
// for a predetermined amount of time simulating transit/processing
// latency
type Middleware struct {
	Latency int // A static latency in milliseconds
	Min     int // A minumum latency for a random range
	Max     int // a maximum latency for a random range
}

// Init takes a configuration mapping for either static latency or a latency
// range.
func (latency *Middleware) Init(conf map[string]string) error {
	if v, prs := conf["latency"]; prs == true {
		l, e := strconv.Atoi(v)
		if e != nil {
			return e
		}

		latency.Latency = l
	}

	if v, prs := conf["min"]; prs == true {
		min, e := strconv.Atoi(v)
		if e != nil {
			return e
		}

		latency.Min = min
	}

	if v, prs := conf["max"]; prs == true {
		max, e := strconv.Atoi(v)
		if e != nil {
			return e
		}

		latency.Max = max
	}

	return nil
}

// Middleware implements the Middleware interface and injects latency into a
// response based on the configurations defined on the handler.
func (latency *Middleware) Middleware(next http.Handler) http.Handler {
	var duration int

	if latency.Latency > 0 {
		duration = latency.Latency
	} else if (latency.Max >= latency.Min) && latency.Min > 0 {
		diff := latency.Max - latency.Min
		r := rand.Intn(diff)
		duration = r + diff
	}

	if duration > 0 {
		time.Sleep(time.Duration(duration) * time.Millisecond)
	}

	return next
}
