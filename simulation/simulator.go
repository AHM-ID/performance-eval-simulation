package simulation

import (
	"des/models"
	"fmt"
	"math"
	"time"
)

type DiscreteEventSimulator struct {
	state      *models.SystemState
	config     *models.SimulationConfig
	events     *EventManager
	stats      *EnhancedStatisticsCollector
	visualizer *TerminalVisualizer
	customerID int
	eventLog   []*models.EventLogEntry
}

func (sim *DiscreteEventSimulator) GetState() *models.SystemState {
	return sim.state
}

func (sim *DiscreteEventSimulator) GetEvents() *EventManager {
	return sim.events
}

func (sim *DiscreteEventSimulator) GetStats() *EnhancedStatisticsCollector {
	return sim.stats
}

func (sim *DiscreteEventSimulator) GetVisualizer() *TerminalVisualizer {
	return sim.visualizer
}

func (sim *DiscreteEventSimulator) GetEventLog() []*models.EventLogEntry {
	return sim.eventLog
}

func (sim *DiscreteEventSimulator) ProcessEvent(event *models.Event) {
	timeDiff := event.Timestamp - sim.state.LastEventTime
	if timeDiff > 0 {
		sim.stats.UpdatePreEvent(sim.state, timeDiff)
	}

	sim.processEvent(event)
	sim.logEvent(event)

	sim.state.Clock = event.Timestamp
	sim.state.LastEventTime = sim.state.Clock
	sim.state.EventsProcessed++
}

func (sim *DiscreteEventSimulator) ShouldStop() bool {
	return sim.state.Clock >= sim.config.SimulationTime
}

func NewSimulator(config *models.SimulationConfig) *DiscreteEventSimulator {
	sim := &DiscreteEventSimulator{
		config:     config,
		events:     NewEventManager(config),
		visualizer: NewTerminalVisualizer(),
		customerID: 1,
		eventLog:   make([]*models.EventLogEntry, 0),
	}
	sim.initializeState()
	return sim
}

func (sim *DiscreteEventSimulator) Initialize() {
	sim.initializeState()
	firstArrivalTime := sim.events.GetInterarrivalTime()
	sim.events.ScheduleEvent(models.EventArrival, firstArrivalTime, nil)
}

func (sim *DiscreteEventSimulator) initializeState() {
	sim.state = &models.SystemState{
		Clock:             0,
		ServerBusy:        false,
		Queue:             make([]*models.Customer, 0),
		NextArrivalTime:   0,
		NextDepartureTime: math.Inf(1),
		LastEventTime:     0,
		CustomersServed:   0,
		TotalCustomers:    0,
		TotalDelay:        0,
		AreaUnderQ:        0,
		AreaUnderB:        0,
		EventsProcessed:   0,
		RejectedCustomers: 0,
	}
	sim.stats = NewStatisticsCollector(sim.config)
	sim.customerID = 1
	sim.eventLog = make([]*models.EventLogEntry, 0)
}

func (sim *DiscreteEventSimulator) Run() *models.SimulationResults {
	if sim.config.Visualization.Enabled {
		sim.visualizer.ClearScreen()
		sim.visualizer.DisplayHeader(sim.config)
	}
	startTime := time.Now()

	for {
		nextEvent := sim.events.PeekNextEvent()
		if nextEvent == nil {
			break
		}

		if nextEvent.Timestamp > sim.config.SimulationTime {
			break
		}

		timeDiff := nextEvent.Timestamp - sim.state.LastEventTime
		if timeDiff > 0 {
			sim.stats.UpdatePreEvent(sim.state, timeDiff)
		}

		event := sim.events.GetNextEvent()
		sim.state.Clock = event.Timestamp
		sim.processEvent(event)

		sim.state.LastEventTime = sim.state.Clock
		sim.state.EventsProcessed++

		if sim.config.Visualization.Enabled {
			sim.visualizer.ClearScreen()
			sim.visualizer.DisplayHeader(sim.config)
			peekNext := sim.events.PeekNextEvent()
			sim.visualizer.DisplayState(sim.state, peekNext, sim.config)
			time.Sleep(sim.config.Visualization.UpdateInterval)
		}
	}

	if sim.state.Clock < sim.config.SimulationTime {
		finalTimeDiff := sim.config.SimulationTime - sim.state.LastEventTime
		if finalTimeDiff > 0 {
			sim.stats.UpdatePreEvent(sim.state, finalTimeDiff)
		}
		sim.state.Clock = sim.config.SimulationTime
	}

	metrics := sim.stats.CalculateFinalMetrics(sim.state)

	return &models.SimulationResults{
		Config:   sim.config,
		Metrics:  metrics,
		State:    sim.state,
		Runtime:  time.Since(startTime),
		EventLog: sim.eventLog,
	}
}

func (sim *DiscreteEventSimulator) processEvent(event *models.Event) {
	switch event.Type {
	case models.EventArrival:
		sim.processArrival()
	case models.EventDeparture:
		sim.processDeparture(event.Customer)
	}
}

func (sim *DiscreteEventSimulator) processArrival() {
	serviceTime := sim.events.GetServiceTime()
	customer := &models.Customer{
		ID:          sim.customerID,
		ArrivalTime: sim.state.Clock,
		ServiceTime: serviceTime,
		Status:      models.CustomerWaiting,
	}
	sim.customerID++
	sim.state.TotalCustomers++

	logMessage := fmt.Sprintf("Customer %d arrived at time %.2f, service time=%.2f",
		customer.ID, sim.state.Clock, serviceTime)
	if sim.visualizer.logger != nil {
		sim.visualizer.logger.LogInfo(logMessage)
	}

	if !sim.state.ServerBusy {
		sim.state.ServerBusy = true
		customer.ServiceStart = sim.state.Clock
		customer.Status = models.CustomerInService
		departureTime := sim.state.Clock + customer.ServiceTime
		sim.events.ScheduleEvent(models.EventDeparture, departureTime, customer)
		sim.state.CustomersServed++
	} else if len(sim.state.Queue) < sim.config.MaxQueueSize {
		sim.state.Queue = append(sim.state.Queue, customer)
	} else {
		customer.Status = models.CustomerRejected
		sim.state.RejectedCustomers++
	}

	nextArrivalTime := sim.state.Clock + sim.events.GetInterarrivalTime()
	if nextArrivalTime <= sim.config.SimulationTime {
		sim.events.ScheduleEvent(models.EventArrival, nextArrivalTime, nil)
	}
}

func (sim *DiscreteEventSimulator) processDeparture(customer *models.Customer) {
	if customer != nil {
		customer.ExitTime = sim.state.Clock
		sim.stats.RecordCustomerCompletion(customer)
		logMessage := fmt.Sprintf("Customer %d departed at time %.2f",
			customer.ID, sim.state.Clock)
		if sim.visualizer.logger != nil {
			sim.visualizer.logger.LogInfo(logMessage)
		}
	}

	if len(sim.state.Queue) > 0 {
		nextCustomer := sim.state.Queue[0]
		sim.state.Queue = sim.state.Queue[1:]

		delay := sim.state.Clock - nextCustomer.ArrivalTime
		if delay < 0 {
			delay = 0
		}
		sim.state.TotalDelay += delay

		sim.state.ServerBusy = true
		nextCustomer.ServiceStart = sim.state.Clock
		nextCustomer.Status = models.CustomerInService
		departureTime := sim.state.Clock + nextCustomer.ServiceTime
		sim.events.ScheduleEvent(models.EventDeparture, departureTime, nextCustomer)
		sim.state.CustomersServed++
	} else {
		sim.state.ServerBusy = false
	}
}

func (sim *DiscreteEventSimulator) logEvent(event *models.Event) {
	eventType := ""
	action := ""

	switch event.Type {
	case models.EventArrival:
		eventType = "ARRIVAL"
		action = "Customer arrived"
	case models.EventDeparture:
		eventType = "DEPARTURE"
		action = "Customer departed"
	}

	logEntry := &models.EventLogEntry{
		Time:       event.Timestamp,
		EventType:  eventType,
		QueueSize:  len(sim.state.Queue),
		ServerBusy: sim.state.ServerBusy,
		Action:     action,
	}

	if event.Customer != nil {
		logEntry.CustomerID = event.Customer.ID
	}

	sim.eventLog = append(sim.eventLog, logEntry)
}
