package metrics

import (
	"runtime"
	"runtime/debug"
	"runtime/pprof"
	"time"

	kit "github.com/go-kit/kit/metrics"
)

const (
	MetricDebugGCLast       = "debug.gc.last"
	MetricDebugGCNum        = "debug.gc.num"
	MetricDebugGCPause      = "debug.gc.pause"
	MetricDebugGCPauseTotal = "debug.gc.pause_total"
	MetricDebugGCReadStats  = "debug.gc.read_stats"

	MetricRuntimeReadMemStats = "runtime.read_mem_stats"

	MetricMemStatsAlloc      = "runtime.mem_stats.alloc"
	MetricMemStatsTotalAlloc = "runtime.mem_stats.total_alloc"
	MetricMemStatsSys        = "runtime.mem_stats.sys"
	MetricMemStatsLookups    = "runtime.mem_stats.lookups"
	MetricMemStatsMallocs    = "runtime.mem_stats.mallocs"
	MetricMemStatsFrees      = "runtime.mem_stats.frees"

	MetricMemStatsHeapAlloc    = "runtime.mem_stats.heap_alloc"
	MetricMemStatsHeapSys      = "runtime.mem_stats.heap_sys"
	MetricMemStatsHeapIdle     = "runtime.mem_stats.heap_idle"
	MetricMemStatsHeapInuse    = "runtime.mem_stats.heap_inuse"
	MetricMemStatsHeapReleased = "runtime.mem_stats.heap_released"
	MetricMemStatsHeapObjects  = "runtime.mem_stats.heap_objects"

	MetricMemStatsStackInuse  = "runtime.mem_stats.stack_inuse"
	MetricMemStatsStackSys    = "runtime.mem_stats.stack_sys"
	MetricMemStatsMSpanInuse  = "runtime.mem_stats.m_span_inuse"
	MetricMemStatsMSpanSys    = "runtime.mem_stats.m_span_sys"
	MetricMemStatsMCacheInuse = "runtime.mem_stats.m_cache_inuse"
	MetricMemStatsMCacheSys   = "runtime.mem_stats.m_cache_sys"
	MetricMemStatsBuckHashSys = "runtime.mem_stats.buck_hash_sys"
	MetricMemStatsGCSys       = "runtime.mem_stats.gc_sys"
	MetricMemStatsOtherSys    = "runtime.mem_stats.other_sys"

	MetricMemStatsNextGC        = "runtime.mem_stats.next_gc"
	MetricMemStatsLastGC        = "runtime.mem_stats.last_gc"
	MetricMemStatsPauseTotalNs  = "runtime.mem_stats.pause_total_ns"
	MetricMemStatsPauseNs       = "runtime.mem_stats.pause_ns"
	MetricMemStatsNumGC         = "runtime.mem_stats.num_gc"
	MetricMemStatsGCCPUFraction = "runtime.mem_stats.gc_cpu_fraction"
	MetricMemStatsEnableGC      = "runtime.mem_stats.enabled_gc"
	MetricMemStatsDebugGC       = "runtime.mem_stats.debug_gc"

	MetricRuntimeNumCgoCall   = "runtime.num_cgo_call"
	MetricRuntimeNumGoroutine = "runtime.num_goroutine"
	MetricRuntimeNumThread    = "runtime.num_thread"
)

var (
	debugMetricGCLast       kit.Gauge
	debugMetricGCNum        kit.Gauge
	debugMetricGCPause      kit.Histogram
	debugMetricGCPauseTotal kit.Gauge
	debugMetricGCReadStats  Timer

	runtimeMetricReadMemStats Timer

	runtimeMetricMemStatsAlloc      kit.Gauge
	runtimeMetricMemStatsTotalAlloc kit.Gauge
	runtimeMetricMemStatsSys        kit.Gauge
	runtimeMetricMemStatsLookups    kit.Gauge
	runtimeMetricMemStatsMallocs    kit.Gauge
	runtimeMetricMemStatsFrees      kit.Gauge

	runtimeMetricMemStatsHeapAlloc    kit.Gauge
	runtimeMetricMemStatsHeapSys      kit.Gauge
	runtimeMetricMemStatsHeapIdle     kit.Gauge
	runtimeMetricMemStatsHeapInuse    kit.Gauge
	runtimeMetricMemStatsHeapReleased kit.Gauge
	runtimeMetricMemStatsHeapObjects  kit.Gauge

	runtimeMetricMemStatsStackInuse  kit.Gauge
	runtimeMetricMemStatsStackSys    kit.Gauge
	runtimeMetricMemStatsMSpanInuse  kit.Gauge
	runtimeMetricMemStatsMSpanSys    kit.Gauge
	runtimeMetricMemStatsMCacheInuse kit.Gauge
	runtimeMetricMemStatsMCacheSys   kit.Gauge
	runtimeMetricMemStatsBuckHashSys kit.Gauge
	runtimeMetricMemStatsGCSys       kit.Gauge
	runtimeMetricMemStatsOtherSys    kit.Gauge

	runtimeMetricMemStatsNextGC        kit.Gauge
	runtimeMetricMemStatsLastGC        kit.Gauge
	runtimeMetricMemStatsPauseTotalNs  kit.Gauge
	runtimeMetricMemStatsPauseNs       kit.Histogram
	runtimeMetricMemStatsNumGC         kit.Gauge
	runtimeMetricMemStatsGCCPUFraction kit.Gauge
	runtimeMetricMemStatsEnableGC      kit.Gauge
	runtimeMetricMemStatsDebugGC       kit.Gauge

	runtimeMetricNumCgoCall   kit.Gauge
	runtimeMetricNumGoroutine kit.Gauge
	runtimeMetricNumThread    kit.Gauge

	debugGCStats    debug.GCStats
	runtimeMemStats runtime.MemStats

	runtimeNumCgoCalls         int64
	runtimeFrees               uint64
	runtimeLookups             uint64
	runtimeMallocs             uint64
	runtimeNumGC               uint32
	runtimeThreadCreateProfile = pprof.Lookup("threadcreate")
)

func init() {
	debugGCStats.Pause = make([]time.Duration, 11)
}

func (c *Component) MetricsCapture() {
	CaptureDebugMetrics()
	CaptureRuntimeMetrics()
}

func (c *Component) MetricsRegister(m *Component) {
	RegisterDebugMetrics(m)
	RegisterRuntimeMetrics(m)
}

func CaptureDebugMetrics() {
	gcLast := debugGCStats.LastGC

	t := time.Now()
	debug.ReadGCStats(&debugGCStats)
	debugMetricGCReadStats.ObserveDurationByTime(t)

	debugMetricGCLast.Set(float64(debugGCStats.LastGC.UnixNano()))
	debugMetricGCNum.Set(float64(debugGCStats.NumGC))
	debugMetricGCPauseTotal.Set(float64(debugGCStats.PauseTotal))

	if gcLast != debugGCStats.LastGC && len(debugGCStats.Pause) > 0 {
		debugMetricGCPause.Observe(float64(debugGCStats.Pause[0]))
	}
}

func CaptureRuntimeMetrics() {
	t := time.Now()
	runtime.ReadMemStats(&runtimeMemStats)
	runtimeMetricReadMemStats.ObserveDurationByTime(t)

	runtimeMetricMemStatsAlloc.Set(float64(runtimeMemStats.Alloc))
	runtimeMetricMemStatsTotalAlloc.Set(float64(runtimeMemStats.TotalAlloc))
	runtimeMetricMemStatsSys.Set(float64(runtimeMemStats.Sys))
	runtimeMetricMemStatsLookups.Set(float64(runtimeMemStats.Lookups - runtimeLookups))
	runtimeMetricMemStatsMallocs.Set(float64(runtimeMemStats.Mallocs - runtimeMallocs))
	runtimeMetricMemStatsFrees.Set(float64(runtimeMemStats.Frees - runtimeFrees))

	runtimeMetricMemStatsHeapAlloc.Set(float64(runtimeMemStats.HeapAlloc))
	runtimeMetricMemStatsHeapSys.Set(float64(runtimeMemStats.HeapSys))
	runtimeMetricMemStatsHeapIdle.Set(float64(runtimeMemStats.HeapIdle))
	runtimeMetricMemStatsHeapInuse.Set(float64(runtimeMemStats.HeapInuse))
	runtimeMetricMemStatsHeapReleased.Set(float64(runtimeMemStats.HeapReleased))
	runtimeMetricMemStatsHeapObjects.Set(float64(runtimeMemStats.HeapObjects))

	runtimeMetricMemStatsStackInuse.Set(float64(runtimeMemStats.StackInuse))
	runtimeMetricMemStatsStackSys.Set(float64(runtimeMemStats.StackSys))
	runtimeMetricMemStatsMSpanInuse.Set(float64(runtimeMemStats.MSpanInuse))
	runtimeMetricMemStatsMSpanSys.Set(float64(runtimeMemStats.MSpanSys))
	runtimeMetricMemStatsMCacheInuse.Set(float64(runtimeMemStats.MCacheInuse))
	runtimeMetricMemStatsMCacheSys.Set(float64(runtimeMemStats.MCacheSys))
	runtimeMetricMemStatsBuckHashSys.Set(float64(runtimeMemStats.BuckHashSys))
	runtimeMetricMemStatsGCSys.Set(float64(runtimeMemStats.GCSys))
	runtimeMetricMemStatsOtherSys.Set(float64(runtimeMemStats.OtherSys))

	runtimeMetricMemStatsNextGC.Set(float64(runtimeMemStats.NextGC))
	runtimeMetricMemStatsLastGC.Set(float64(runtimeMemStats.LastGC))
	runtimeMetricMemStatsPauseTotalNs.Set(float64(runtimeMemStats.PauseTotalNs))

	i := runtimeNumGC % uint32(len(runtimeMemStats.PauseNs))
	ii := runtimeMemStats.NumGC % uint32(len(runtimeMemStats.PauseNs))
	if runtimeMemStats.NumGC-runtimeNumGC >= uint32(len(runtimeMemStats.PauseNs)) {
		for i = 0; i < uint32(len(runtimeMemStats.PauseNs)); i++ {
			runtimeMetricMemStatsPauseNs.Observe(float64(runtimeMemStats.PauseNs[i]))
		}
	} else {
		if i > ii {
			for ; i < uint32(len(runtimeMemStats.PauseNs)); i++ {
				runtimeMetricMemStatsPauseNs.Observe(float64(runtimeMemStats.PauseNs[i]))
			}
			i = 0
		}
		for ; i < ii; i++ {
			runtimeMetricMemStatsPauseNs.Observe(float64(runtimeMemStats.PauseNs[i]))
		}
	}

	runtimeMetricMemStatsNumGC.Set(float64(runtimeMemStats.NumGC - runtimeNumGC))
	runtimeMetricMemStatsGCCPUFraction.Set(gcCPUFraction(&runtimeMemStats))

	if runtimeMemStats.EnableGC {
		runtimeMetricMemStatsEnableGC.Set(1)
	} else {
		runtimeMetricMemStatsEnableGC.Set(0)
	}

	if runtimeMemStats.DebugGC {
		runtimeMetricMemStatsDebugGC.Set(1)
	} else {
		runtimeMetricMemStatsDebugGC.Set(0)
	}

	currentNumCgoCalls := getNumCgoCall()
	runtimeMetricNumCgoCall.Set(float64(currentNumCgoCalls - runtimeNumCgoCalls))

	runtimeMetricNumGoroutine.Set(float64(runtime.NumGoroutine()))
	runtimeMetricNumThread.Set(float64(runtimeThreadCreateProfile.Count()))

	runtimeNumCgoCalls = currentNumCgoCalls
	runtimeFrees = runtimeMemStats.Frees
	runtimeLookups = runtimeMemStats.Lookups
	runtimeMallocs = runtimeMemStats.Mallocs
	runtimeNumGC = runtimeMemStats.NumGC
}

func RegisterDebugMetrics(r *Component) {
	debugMetricGCLast = r.NewGauge(MetricDebugGCLast)
	debugMetricGCNum = r.NewGauge(MetricDebugGCNum)
	debugMetricGCPause = r.NewHistogram(MetricDebugGCPause)
	debugMetricGCPauseTotal = r.NewGauge(MetricDebugGCPauseTotal)
	debugMetricGCReadStats = r.NewTimer(MetricDebugGCReadStats)
}

func RegisterRuntimeMetrics(r *Component) {
	runtimeMetricReadMemStats = r.NewTimer(MetricRuntimeReadMemStats)

	runtimeMetricMemStatsAlloc = r.NewGauge(MetricMemStatsAlloc)
	runtimeMetricMemStatsTotalAlloc = r.NewGauge(MetricMemStatsTotalAlloc)
	runtimeMetricMemStatsSys = r.NewGauge(MetricMemStatsSys)
	runtimeMetricMemStatsLookups = r.NewGauge(MetricMemStatsLookups)
	runtimeMetricMemStatsMallocs = r.NewGauge(MetricMemStatsMallocs)
	runtimeMetricMemStatsFrees = r.NewGauge(MetricMemStatsFrees)

	runtimeMetricMemStatsHeapAlloc = r.NewGauge(MetricMemStatsHeapAlloc)
	runtimeMetricMemStatsHeapSys = r.NewGauge(MetricMemStatsHeapSys)
	runtimeMetricMemStatsHeapIdle = r.NewGauge(MetricMemStatsHeapIdle)
	runtimeMetricMemStatsHeapInuse = r.NewGauge(MetricMemStatsHeapInuse)
	runtimeMetricMemStatsHeapReleased = r.NewGauge(MetricMemStatsHeapReleased)
	runtimeMetricMemStatsHeapObjects = r.NewGauge(MetricMemStatsHeapObjects)

	runtimeMetricMemStatsStackInuse = r.NewGauge(MetricMemStatsStackInuse)
	runtimeMetricMemStatsStackSys = r.NewGauge(MetricMemStatsStackSys)
	runtimeMetricMemStatsMSpanInuse = r.NewGauge(MetricMemStatsMSpanInuse)
	runtimeMetricMemStatsMSpanSys = r.NewGauge(MetricMemStatsMSpanSys)
	runtimeMetricMemStatsMCacheInuse = r.NewGauge(MetricMemStatsMCacheInuse)
	runtimeMetricMemStatsMCacheSys = r.NewGauge(MetricMemStatsMCacheSys)
	runtimeMetricMemStatsBuckHashSys = r.NewGauge(MetricMemStatsBuckHashSys)
	runtimeMetricMemStatsGCSys = r.NewGauge(MetricMemStatsGCSys)
	runtimeMetricMemStatsOtherSys = r.NewGauge(MetricMemStatsOtherSys)

	runtimeMetricMemStatsNextGC = r.NewGauge(MetricMemStatsNextGC)
	runtimeMetricMemStatsLastGC = r.NewGauge(MetricMemStatsLastGC)
	runtimeMetricMemStatsPauseTotalNs = r.NewGauge(MetricMemStatsPauseTotalNs)
	runtimeMetricMemStatsPauseNs = r.NewHistogram(MetricMemStatsPauseNs)
	runtimeMetricMemStatsNumGC = r.NewGauge(MetricMemStatsNumGC)
	runtimeMetricMemStatsGCCPUFraction = r.NewGauge(MetricMemStatsGCCPUFraction)
	runtimeMetricMemStatsEnableGC = r.NewGauge(MetricMemStatsEnableGC)
	runtimeMetricMemStatsDebugGC = r.NewGauge(MetricMemStatsDebugGC)

	runtimeMetricNumCgoCall = r.NewGauge(MetricRuntimeNumCgoCall)
	runtimeMetricNumGoroutine = r.NewGauge(MetricRuntimeNumGoroutine)
	runtimeMetricNumThread = r.NewGauge(MetricRuntimeNumThread)
}
