package models

import "time"

type StopCondition struct {
	AutomaticMode bool
	Type          string
	Value         int
	TimeLimit     float64
}

type VisualizationConfig struct {
	Enabled             bool
	UpdateInterval      time.Duration
	ShowRealtimeMetrics bool
	ProgressBarWidth    int
}

type RandomConfig struct {
	Seed         int64
	Distribution string
}

type LoggingConfig struct {
	Level        string
	LogToFile    bool
	LogFilePath  string
	OutputFormat string
}

type SimulationConfig struct {
	SimulationTime float64
	ArrivalRate    float64
	ServiceRate    float64
	MaxQueueSize   int
	MaxCustomers   int
	StopCondition  StopCondition
	Visualization  VisualizationConfig
	Random         RandomConfig
	Logging        LoggingConfig
}
