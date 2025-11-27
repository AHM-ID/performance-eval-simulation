package simulation

import (
	"des/logging"
	"des/models"
	"fmt"
	"strings"
	"time"
)

type TerminalVisualizer struct {
	logger *logging.Logger
}

func NewTerminalVisualizer() *TerminalVisualizer {
	return &TerminalVisualizer{}
}

func (tv *TerminalVisualizer) SetLogger(logger *logging.Logger) {
	tv.logger = logger
}

func (tv *TerminalVisualizer) ClearScreen() {
	fmt.Print("\033[2J")
	fmt.Print("\033[H")
}

func (tv *TerminalVisualizer) DisplayHeader(config *models.SimulationConfig) {
	header := fmt.Sprintf("%s\nDISCRETE EVENT SIMULATION - SINGLE SERVER QUEUEING SYSTEM\n%s\nConfiguration: Arrival Rate=%.2f, Service Rate=%.2f, Max Queue=%d\n%s\n",
		strings.Repeat("=", 80),
		strings.Repeat("=", 80),
		config.ArrivalRate, config.ServiceRate, config.MaxQueueSize,
		strings.Repeat("-", 80))

	if tv.logger != nil {
		tv.logger.LogTerminal(header)
	} else {
		fmt.Print(header)
	}
}

func (tv *TerminalVisualizer) DisplayState(state *models.SystemState, nextEvent *models.Event, config *models.SimulationConfig) {
	tv.ClearScreen()
	tv.DisplayHeader(config)

	stateStr := fmt.Sprintf("SIMULATION TIME: %8.2f     EVENTS PROCESSED: %6d\n%s\n",
		state.Clock, state.EventsProcessed, strings.Repeat("-", 80))

	serverStatus := "IDLE"
	if state.ServerBusy {
		serverStatus = "BUSY"
	}
	stateStr += fmt.Sprintf("SERVER STATUS: %-6s    CUSTOMERS SERVED: %6d\n", serverStatus, state.CustomersServed)
	stateStr += fmt.Sprintf("QUEUE LENGTH: %3d/%3d    REJECTED CUSTOMERS: %4d\n",
		len(state.Queue), config.MaxQueueSize, state.RejectedCustomers)

	if nextEvent != nil {
		eventType := "ARRIVAL"
		if nextEvent.Type == models.EventDeparture {
			eventType = "DEPARTURE"
		} else if nextEvent.Type == models.EventTermination {
			eventType = "TERMINATION"
		}
		stateStr += fmt.Sprintf("NEXT EVENT: %-12s at TIME: %8.2f\n", eventType, nextEvent.Timestamp)
	}

	queueDisplay := make([]string, len(state.Queue))
	for i := range queueDisplay {
		queueDisplay[i] = "C"
	}
	stateStr += fmt.Sprintf("QUEUE: [%s]\n", strings.Join(queueDisplay, ""))

	progress := (state.Clock / config.SimulationTime) * float64(config.Visualization.ProgressBarWidth)
	if progress > float64(config.Visualization.ProgressBarWidth) {
		progress = float64(config.Visualization.ProgressBarWidth)
	}
	if progress < 0 {
		progress = 0
	}
	progressPercent := (state.Clock / config.SimulationTime) * 100
	if progressPercent < 0 {
		progressPercent = 0
	}
	if progressPercent > 100 {
		progressPercent = 100
	}
	stateStr += fmt.Sprintf("PROGRESS: [%s%s] %6.1f%%\n",
		strings.Repeat("#", int(progress)),
		strings.Repeat("-", config.Visualization.ProgressBarWidth-int(progress)),
		progressPercent)

	avgWait := 0.0
	if state.CustomersServed > 0 {
		avgWait = state.TotalDelay / float64(state.CustomersServed)
	}

	queueLength := float64(len(state.Queue))

	serverUtil := 0.0
	if state.Clock > 0 {
		serverUtil = (state.AreaUnderB / state.Clock) * 100
		if serverUtil < 0 {
			serverUtil = 0
		}
		if serverUtil > 100 {
			serverUtil = 100
		}
	}

	throughput := 0.0
	if state.Clock > 0 {
		throughput = float64(state.CustomersServed) / state.Clock
		if throughput < 0 {
			throughput = 0
		}
	}

	stateStr += fmt.Sprintf("REAL-TIME METRICS:\n")
	stateStr += fmt.Sprintf("  Avg Wait Time: %6.2f    Queue Length: %6.2f\n", avgWait, queueLength)
	stateStr += fmt.Sprintf("  Server Util:   %6.2f%%   Throughput:   %6.2f cust/time\n", serverUtil, throughput)
	stateStr += fmt.Sprintf("%s\n", strings.Repeat("-", 80))

	if tv.logger != nil {
		tv.logger.LogTerminal(stateStr)
	} else {
		fmt.Print(stateStr)
	}
}

func (tv *TerminalVisualizer) DisplayResults(results *models.SimulationResults, config *models.SimulationConfig) {
	metrics := results.Metrics
	state := results.State

	resultsStr := fmt.Sprintf("\n%s\nSIMULATION RESULTS\n%s\n",
		strings.Repeat("=", 80),
		strings.Repeat("=", 80))

	resultsStr += fmt.Sprintf("PERFORMANCE METRICS:\n")
	resultsStr += fmt.Sprintf("  Average Wait Time in Queue:   %12.4f time units\n", metrics.AverageWaitTime)
	resultsStr += fmt.Sprintf("  Average Time in System:       %12.4f time units\n", metrics.AverageSystemTime)
	resultsStr += fmt.Sprintf("  Average Queue Length:         %12.4f customers\n", metrics.AverageQueueLength)
	resultsStr += fmt.Sprintf("  Average Customers in System:  %12.4f customers\n", metrics.AverageInSystem)
	resultsStr += fmt.Sprintf("  Server Utilization:           %12.4f %%\n", metrics.ServerUtilization*100)
	resultsStr += fmt.Sprintf("  System Throughput:            %12.4f customers/time unit\n", metrics.Throughput)
	resultsStr += fmt.Sprintf("  Queue Probability:            %12.4f\n", metrics.QueueProbability)
	resultsStr += fmt.Sprintf("  Rejected Customers:           %12d\n", metrics.RejectedCustomers)

	resultsStr += fmt.Sprintf("  Blocking Probability:         %12.4f %%\n", metrics.BlockingProbability*100)
	resultsStr += fmt.Sprintf("  Max Queue Length:             %12d\n", metrics.MaxQueueLength)
	resultsStr += fmt.Sprintf("  Max Wait Time:                %12.4f\n", metrics.MaxWaitTime)
	resultsStr += fmt.Sprintf("  Wait Time Variance:           %12.4f\n", metrics.WaitTimeVariance)
	resultsStr += fmt.Sprintf("  System Time Variance:         %12.4f\n", metrics.SystemTimeVariance)

	if len(metrics.WaitTimePercentiles) > 0 {
		resultsStr += fmt.Sprintf("\nWAIT TIME PERCENTILES:\n")
		resultsStr += fmt.Sprintf("  50th (Median):              %12.4f\n", metrics.WaitTimePercentiles["50th"])
		resultsStr += fmt.Sprintf("  75th:                       %12.4f\n", metrics.WaitTimePercentiles["75th"])
		resultsStr += fmt.Sprintf("  90th:                       %12.4f\n", metrics.WaitTimePercentiles["90th"])
		resultsStr += fmt.Sprintf("  95th:                       %12.4f\n", metrics.WaitTimePercentiles["95th"])
	}

	resultsStr += fmt.Sprintf("\n95%% CONFIDENCE INTERVALS:\n")
	resultsStr += fmt.Sprintf("  Wait Time:        [%8.4f, %8.4f]\n",
		metrics.WaitTimeConfidence[0], metrics.WaitTimeConfidence[1])
	resultsStr += fmt.Sprintf("  System Time:      [%8.4f, %8.4f]\n",
		metrics.SystemTimeConfidence[0], metrics.SystemTimeConfidence[1])

	resultsStr += fmt.Sprintf("\nSIMULATION SUMMARY:\n")
	resultsStr += fmt.Sprintf("  Total Simulation Time:        %12.2f time units\n", state.Clock)
	resultsStr += fmt.Sprintf("  Total Customers Arrived:      %12d\n", state.TotalCustomers)
	resultsStr += fmt.Sprintf("  Total Customers Served:       %12d\n", state.CustomersServed)
	resultsStr += fmt.Sprintf("  Total Events Processed:       %12d\n", state.EventsProcessed)
	resultsStr += fmt.Sprintf("  Area under Q(t):              %12.4f\n", state.AreaUnderQ)
	resultsStr += fmt.Sprintf("  Area under B(t):              %12.4f\n", state.AreaUnderB)
	resultsStr += fmt.Sprintf("  Real Execution Time:          %12v\n", results.Runtime)

	resultsStr += fmt.Sprintf("%s\n", strings.Repeat("=", 80))

	if tv.logger != nil {
		tv.logger.LogTerminal(resultsStr)
	} else {
		fmt.Print(resultsStr)
	}
}

func (tv *TerminalVisualizer) DisplayExecutionInfo(executionTime time.Duration, eventsProcessed int) {
	infoStr := fmt.Sprintf("\nEXECUTION INFORMATION:\n")
	infoStr += fmt.Sprintf("  Real-time execution: %v\n", executionTime)
	if executionTime > 0 {
		infoStr += fmt.Sprintf("  Performance:         %.2f events/second\n",
			float64(eventsProcessed)/executionTime.Seconds())
	}
	infoStr += fmt.Sprintf("%s\n", strings.Repeat("=", 80))

	if tv.logger != nil {
		tv.logger.LogTerminal(infoStr)
	} else {
		fmt.Print(infoStr)
	}
}
