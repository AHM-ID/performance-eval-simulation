package main

import (
	"bufio"
	"des/config"
	"des/logging"
	"des/models"
	"des/simulation"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	configFile := flag.String("config", "config.yml", "Path to configuration file")
	runMode := flag.String("mode", "automatic", "Run mode: automatic or manual")
	flag.Parse()

	cfg, err := config.LoadConfig(*configFile)
	if err != nil {
		fmt.Printf("Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	var automaticValue bool
	switch strings.ToLower(*runMode) {
	case "automatic", "a":
		automaticValue = true
	case "manual", "m":
		automaticValue = false
	default:
		fmt.Printf("Invalid mode: %s\nValid options: automatic (a), manual (m)\n", *runMode)
		os.Exit(1)
	}

	logger := logging.NewLogger(
		cfg.Logging.Level,
		cfg.Logging.LogToFile,
		cfg.Logging.LogFilePath,
		cfg.Logging.OutputFormat,
	)

	initializeSimulation(logger, cfg)

	simulator := simulation.NewSimulator(cfg)
	simulator.GetVisualizer().SetLogger(logger)
	simulator.Initialize()

	cfg.StopCondition.AutomaticMode = automaticValue

	if cfg.StopCondition.AutomaticMode {
		runAutomaticSimulation(simulator, logger, cfg)
	} else {
		runManualSimulation(simulator, logger, cfg)
	}

	logger.LogInfo("Simulation completed successfully")
}

func initializeSimulation(logger *logging.Logger, cfg *models.SimulationConfig) {
	logger.LogInfo("=== DISCRETE EVENT SIMULATION INITIALIZATION ===")
	logger.LogInfo(fmt.Sprintf("Automatic Mode: %v", cfg.StopCondition.AutomaticMode))
	logger.LogInfo(fmt.Sprintf("Simulation Time: %.2f", cfg.SimulationTime))
	logger.LogInfo(fmt.Sprintf("Arrival Rate: %.2f", cfg.ArrivalRate))
	logger.LogInfo(fmt.Sprintf("Service Rate: %.2f", cfg.ServiceRate))
	logger.LogInfo(fmt.Sprintf("Max Queue Size: %d", cfg.MaxQueueSize))
	logger.LogInfo(fmt.Sprintf("Max Customers: %d", cfg.MaxCustomers))
	logger.LogInfo(fmt.Sprintf("Stop Condition: %s", cfg.StopCondition.Type))
}

func runAutomaticSimulation(simulator *simulation.DiscreteEventSimulator, logger *logging.Logger, cfg *models.SimulationConfig) {
	logger.LogInfo("Starting simulation in AUTOMATIC mode")
	if simulator.GetState().Clock == 0 && simulator.GetEvents().PeekNextEvent() == nil {
		simulator.Initialize()
	}
	results := simulator.Run()
	simulator.GetVisualizer().DisplayResults(results, cfg)
	simulator.GetVisualizer().DisplayExecutionInfo(results.Runtime, results.State.EventsProcessed)
}

func runManualSimulation(simulator *simulation.DiscreteEventSimulator, logger *logging.Logger, cfg *models.SimulationConfig) {
	logger.LogInfo("Starting simulation in MANUAL mode")
	logger.LogInfo("Press ENTER key to advance to next event")

	simulator.Initialize()
	scanner := bufio.NewScanner(os.Stdin)

	for {
		if simulator.ShouldStop() {
			break
		}

		nextEvent := simulator.GetEvents().PeekNextEvent()
		if nextEvent == nil {
			break
		}

		if nextEvent.Timestamp > cfg.SimulationTime {
			break
		}

		simulator.GetVisualizer().ClearScreen()
		simulator.GetVisualizer().DisplayHeader(cfg)
		simulator.GetVisualizer().DisplayState(simulator.GetState(), nextEvent, cfg)

		fmt.Print("\n\tPress SPACE to continue (or 'q' to quit): ")
		scanner.Scan()
		input := strings.TrimSpace(scanner.Text())

		if input == "q" || input == "Q" {
			logger.LogInfo("Manual simulation stopped by user")
			break
		}

		event := simulator.GetEvents().GetNextEvent()
		if event == nil {
			break
		}

		timeDiff := event.Timestamp - simulator.GetState().LastEventTime
		if timeDiff > 0 {
			simulator.GetStats().UpdatePreEvent(simulator.GetState(), timeDiff)
		}

		simulator.GetState().Clock = event.Timestamp
		simulator.ProcessEvent(event)

		if simulator.ShouldStop() {
			break
		}
	}

	if simulator.GetState().Clock < cfg.SimulationTime {
		finalTimeDiff := cfg.SimulationTime - simulator.GetState().LastEventTime
		if finalTimeDiff > 0 {
			simulator.GetStats().UpdatePreEvent(simulator.GetState(), finalTimeDiff)
		}
		simulator.GetState().Clock = cfg.SimulationTime
	}

	metrics := simulator.GetStats().CalculateFinalMetrics(simulator.GetState())
	results := &models.SimulationResults{
		Config:   cfg,
		Metrics:  metrics,
		State:    simulator.GetState(),
		Runtime:  time.Duration(0),
		EventLog: simulator.GetEventLog(),
	}

	simulator.GetVisualizer().DisplayResults(results, cfg)
}
