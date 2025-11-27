package models

import "time"

// ComprehensiveMetrics holds ALL performance metrics
type ComprehensiveMetrics struct {
	AverageWaitTime        float64
	AverageQueueLength     float64
	ServerUtilization      float64
	AverageSystemTime      float64
	AverageInSystem        float64
	Throughput             float64
	QueueProbability       float64
	RejectedCustomers      int
	TotalCustomers         int
	BlockingProbability    float64
	WaitTimeVariance       float64
	SystemTimeVariance     float64
	MaxQueueLength         int
	MaxWaitTime            float64
	ServerIdleTime         float64
	ServerBusyTime         float64
	WaitTimeConfidence     [2]float64
	SystemTimeConfidence   [2]float64
	WaitTimePercentiles    map[string]float64
	SystemTimePercentiles  map[string]float64
}

// SimulationResults holds complete simulation results
type SimulationResults struct {
	Config      *SimulationConfig
	Metrics     *ComprehensiveMetrics
	State       *SystemState
	Runtime     time.Duration
	EventLog    []*EventLogEntry
}

// SystemState represents the current state of the simulation
type SystemState struct {
	Clock              float64
	ServerBusy         bool
	Queue              []*Customer
	NextArrivalTime    float64
	NextDepartureTime  float64
	CustomersServed    int
	TotalCustomers     int
	TotalDelay         float64
	AreaUnderQ         float64
	AreaUnderB         float64
	LastEventTime      float64
	EventsProcessed    int
	RejectedCustomers  int
}

// Event represents a discrete event in the simulation
type Event struct {
	Type      EventType
	Timestamp float64
	Customer  *Customer
	Priority  int
}

type EventType int

const (
	EventArrival EventType = iota
	EventDeparture
	EventTermination
)
