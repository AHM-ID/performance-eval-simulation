package models

// Customer represents a customer in the system
type Customer struct {
	ID           int
	ArrivalTime  float64
	ServiceTime  float64
	ServiceStart float64
	ExitTime     float64
	Status       CustomerStatus
}

type CustomerStatus int

const (
	CustomerWaiting CustomerStatus = iota
	CustomerInService
	CustomerCompleted
	CustomerRejected
)

// CustomerStats holds statistics for a customer
type CustomerStats struct {
	WaitTime   float64
	SystemTime float64
}

// EventLogEntry for detailed logging
type EventLogEntry struct {
	Time       float64
	EventType  string
	CustomerID int
	QueueSize  int
	ServerBusy bool
	Action     string
}
