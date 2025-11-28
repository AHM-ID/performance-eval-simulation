package simulation

import (
	"des/models"
	"math"
	"sort"
)

type EnhancedStatisticsCollector struct {
	metrics        *models.ComprehensiveMetrics
	customerStats  []*models.CustomerStats
	waitTimes      []float64
	systemTimes    []float64
	config         *models.SimulationConfig
	maxQueueLength int
	maxWaitTime    float64
}

func NewStatisticsCollector(config *models.SimulationConfig) *EnhancedStatisticsCollector {
	return &EnhancedStatisticsCollector{
		metrics: &models.ComprehensiveMetrics{
			WaitTimePercentiles:   make(map[string]float64),
			SystemTimePercentiles: make(map[string]float64),
		},
		customerStats:  []*models.CustomerStats{},
		waitTimes:      []float64{},
		systemTimes:    []float64{},
		config:         config,
		maxQueueLength: 0,
		maxWaitTime:    0,
	}
}

func (sc *EnhancedStatisticsCollector) UpdatePreEvent(state *models.SystemState, timeDiff float64) {
	if timeDiff <= 0 {
		return
	}
	currentQueueLength := len(state.Queue)
	state.AreaUnderQ += float64(currentQueueLength) * timeDiff

	if currentQueueLength > sc.maxQueueLength {
		sc.maxQueueLength = currentQueueLength
	}

	if state.ServerBusy {
		state.AreaUnderB += timeDiff
	}
}

func (sc *EnhancedStatisticsCollector) RecordCustomerCompletion(customer *models.Customer) {
	waitTime := math.Max(0, customer.ServiceStart-customer.ArrivalTime)
	systemTime := math.Max(0, customer.ExitTime-customer.ArrivalTime)
	stats := &models.CustomerStats{
		WaitTime:   waitTime,
		SystemTime: systemTime,
	}
	sc.customerStats = append(sc.customerStats, stats)
	sc.waitTimes = append(sc.waitTimes, waitTime)
	sc.systemTimes = append(sc.systemTimes, systemTime)
	if waitTime > sc.maxWaitTime {
		sc.maxWaitTime = waitTime
	}
}

func (sc *EnhancedStatisticsCollector) CalculateFinalMetrics(state *models.SystemState) *models.ComprehensiveMetrics {
	sc.metrics.TotalCustomers = state.TotalCustomers
	sc.metrics.RejectedCustomers = state.RejectedCustomers
	sc.metrics.MaxQueueLength = sc.maxQueueLength
	sc.metrics.MaxWaitTime = sc.maxWaitTime

	if state.Clock > 0 {
		sc.metrics.AverageQueueLength = state.AreaUnderQ / state.Clock
		sc.metrics.ServerUtilization = state.AreaUnderB / state.Clock
		sc.metrics.ServerBusyTime = state.AreaUnderB
		sc.metrics.ServerIdleTime = state.Clock - state.AreaUnderB
		sc.metrics.Throughput = float64(state.CustomersServed) / state.Clock
	} else {
		sc.metrics.AverageQueueLength = 0
		sc.metrics.ServerUtilization = 0
		sc.metrics.ServerBusyTime = 0
		sc.metrics.ServerIdleTime = 0
		sc.metrics.Throughput = 0
	}

	if state.CustomersServed > 0 {
		sc.metrics.AverageWaitTime = state.TotalDelay / float64(state.CustomersServed)
	} else {
		sc.metrics.AverageWaitTime = 0
	}

	totalSystemTime := 0.0
	for _, stats := range sc.customerStats {
		totalSystemTime += stats.SystemTime
	}
	if len(sc.customerStats) > 0 {
		sc.metrics.AverageSystemTime = totalSystemTime / float64(len(sc.customerStats))
	} else {
		sc.metrics.AverageSystemTime = 0
	}

	sc.metrics.AverageInSystem = sc.metrics.AverageQueueLength + sc.metrics.ServerUtilization
	if state.Clock > 0 && sc.config.MaxQueueSize > 0 {
		sc.metrics.QueueProbability = state.AreaUnderQ / state.Clock / float64(sc.config.MaxQueueSize)
	} else {
		sc.metrics.QueueProbability = 0
	}

	if state.TotalCustomers > 0 {
		sc.metrics.BlockingProbability = float64(state.RejectedCustomers) / float64(state.TotalCustomers)
	} else {
		sc.metrics.BlockingProbability = 0
	}

	sc.calculateVariances()
	sc.calculatePercentiles()
	sc.calculateConfidenceIntervals()

	return sc.metrics
}

func (sc *EnhancedStatisticsCollector) calculateVariances() {
	if len(sc.waitTimes) < 2 {
		sc.metrics.WaitTimeVariance = 0
		sc.metrics.SystemTimeVariance = 0
		return
	}

	waitSumSq := 0.0
	for _, wt := range sc.waitTimes {
		waitSumSq += (wt - sc.metrics.AverageWaitTime) * (wt - sc.metrics.AverageWaitTime)
	}
	sc.metrics.WaitTimeVariance = waitSumSq / float64(len(sc.waitTimes))

	systemSumSq := 0.0
	for _, st := range sc.systemTimes {
		systemSumSq += (st - sc.metrics.AverageSystemTime) * (st - sc.metrics.AverageSystemTime)
	}
	sc.metrics.SystemTimeVariance = systemSumSq / float64(len(sc.systemTimes))
}

func (sc *EnhancedStatisticsCollector) calculatePercentiles() {
	if len(sc.waitTimes) == 0 {
		return
	}
	waitTimes := append([]float64{}, sc.waitTimes...)
	sort.Float64s(waitTimes)
	sc.metrics.WaitTimePercentiles = map[string]float64{
		"50th": sc.calculatePercentile(waitTimes, 0.5),
		"75th": sc.calculatePercentile(waitTimes, 0.75),
		"90th": sc.calculatePercentile(waitTimes, 0.9),
		"95th": sc.calculatePercentile(waitTimes, 0.95),
	}
}

func (sc *EnhancedStatisticsCollector) calculatePercentile(data []float64, percentile float64) float64 {
	if len(data) == 0 {
		return 0
	}
	pos := percentile * float64(len(data)-1)
	lower := int(math.Floor(pos))
	upper := int(math.Ceil(pos))
	if lower == upper {
		return data[lower]
	}
	weight := pos - float64(lower)
	return data[lower]*(1-weight) + data[upper]*weight
}

func (sc *EnhancedStatisticsCollector) calculateConfidenceIntervals() {
	if len(sc.waitTimes) < 2 {
		sc.metrics.WaitTimeConfidence = [2]float64{0, 0}
		sc.metrics.SystemTimeConfidence = [2]float64{0, 0}
		return
	}

	waitStdErr := math.Sqrt(sc.metrics.WaitTimeVariance/float64(len(sc.waitTimes))) * 1.96
	sc.metrics.WaitTimeConfidence = [2]float64{
		sc.metrics.AverageWaitTime - waitStdErr,
		sc.metrics.AverageWaitTime + waitStdErr,
	}

	if len(sc.systemTimes) >= 2 {
		systemStdErr := math.Sqrt(sc.metrics.SystemTimeVariance/float64(len(sc.systemTimes))) * 1.96
		sc.metrics.SystemTimeConfidence = [2]float64{
			sc.metrics.AverageSystemTime - systemStdErr,
			sc.metrics.AverageSystemTime + systemStdErr,
		}
	}
}
