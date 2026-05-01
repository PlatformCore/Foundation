package metrics

import (
	"sort"
	"sync"
)

// SimpleRegistry is the clean, lightweight compatibility surface merged into the
// heavy gometrics core. It preserves the old observability/metrics/basic API,
// but lives in the core package so new services can import gometrics directly.
type SimpleRegistry struct {
	mu         sync.RWMutex
	counters   map[string]int64
	gauges     map[string]float64
	histograms map[string][]float64
}

func NewSimpleRegistry() *SimpleRegistry {
	return &SimpleRegistry{counters: map[string]int64{}, gauges: map[string]float64{}, histograms: map[string][]float64{}}
}

func (r *SimpleRegistry) Inc(name string) { r.Add(name, 1) }
func (r *SimpleRegistry) Add(name string, delta int64) {
	if name == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters[name] += delta
}
func (r *SimpleRegistry) Set(name string, value float64) {
	if name == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.gauges[name] = value
}
func (r *SimpleRegistry) Observe(name string, value float64) {
	if name == "" {
		return
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.histograms[name] = append(r.histograms[name], value)
}

type SimpleSnapshot struct {
	Counters   map[string]int64
	Gauges     map[string]float64
	Histograms map[string]SimpleHistogramSnapshot
}

type SimpleHistogramSnapshot struct {
	Count int
	Sum   float64
	Avg   float64
	Min   float64
	Max   float64
}

func (r *SimpleRegistry) Snapshot() SimpleSnapshot {
	r.mu.RLock()
	defer r.mu.RUnlock()
	counters := make(map[string]int64, len(r.counters))
	for k, v := range r.counters {
		counters[k] = v
	}
	gauges := make(map[string]float64, len(r.gauges))
	for k, v := range r.gauges {
		gauges[k] = v
	}
	hists := make(map[string]SimpleHistogramSnapshot, len(r.histograms))
	for k, points := range r.histograms {
		hists[k] = summarizeSimple(points)
	}
	return SimpleSnapshot{Counters: counters, Gauges: gauges, Histograms: hists}
}

func (r *SimpleRegistry) Counter(name string) int64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.counters[name]
}
func (r *SimpleRegistry) Gauge(name string) float64 {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.gauges[name]
}
func (r *SimpleRegistry) MetricNames() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	set := map[string]struct{}{}
	for k := range r.counters {
		set[k] = struct{}{}
	}
	for k := range r.gauges {
		set[k] = struct{}{}
	}
	for k := range r.histograms {
		set[k] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func (r *SimpleRegistry) Reset() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.counters = map[string]int64{}
	r.gauges = map[string]float64{}
	r.histograms = map[string][]float64{}
}
func (s SimpleSnapshot) MetricNames() []string {
	set := map[string]struct{}{}
	for k := range s.Counters {
		set[k] = struct{}{}
	}
	for k := range s.Gauges {
		set[k] = struct{}{}
	}
	for k := range s.Histograms {
		set[k] = struct{}{}
	}
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
func summarizeSimple(values []float64) SimpleHistogramSnapshot {
	if len(values) == 0 {
		return SimpleHistogramSnapshot{}
	}
	min, max := values[0], values[0]
	sum := 0.0
	for _, v := range values {
		sum += v
		if v < min {
			min = v
		}
		if v > max {
			max = v
		}
	}
	return SimpleHistogramSnapshot{Count: len(values), Sum: sum, Avg: sum / float64(len(values)), Min: min, Max: max}
}
