package config

import (
	"des/models"
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type YAMLConfig struct {
	Simulation struct {
		SimulationTime float64 `yaml:"simulation_time"`
		ArrivalRate    float64 `yaml:"arrival_rate"`
		ServiceRate    float64 `yaml:"service_rate"`
		MaxQueueSize   int     `yaml:"max_queue_size"`
		MaxCustomers   int     `yaml:"max_customers"`
		StopCondition  struct {
			AutomaticMode bool    `yaml:"automatic_mode"`
			Type          string  `yaml:"type"`
			Value         float64 `yaml:"value"`
			TimeLimit     float64 `yaml:"time_limit"`
		} `yaml:"stop_condition"`
		Visualization struct {
			Enabled             bool `yaml:"enabled"`
			UpdateIntervalMs    int  `yaml:"update_interval_ms"`
			ShowRealtimeMetrics bool `yaml:"show_realtime_metrics"`
			ProgressBarWidth    int  `yaml:"progress_bar_width"`
		} `yaml:"visualization"`
		Random struct {
			Seed         int64  `yaml:"seed"`
			Distribution string `yaml:"distribution"`
		} `yaml:"random"`
		Logging struct {
			Level        string `yaml:"level"`
			LogToFile    bool   `yaml:"log_to_file"`
			LogFilePath  string `yaml:"log_file_path"`
			OutputFormat string `yaml:"output_format"`
		} `yaml:"logging"`
	} `yaml:"simulation"`
}

func LoadConfig(filename string) (*models.SimulationConfig, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var yamlConfig YAMLConfig
	err = yaml.Unmarshal(data, &yamlConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %v", err)
	}

	return convertToModel(&yamlConfig), nil
}

func convertToModel(yamlConfig *YAMLConfig) *models.SimulationConfig {
	return &models.SimulationConfig{
		SimulationTime: yamlConfig.Simulation.SimulationTime,
		ArrivalRate:    yamlConfig.Simulation.ArrivalRate,
		ServiceRate:    yamlConfig.Simulation.ServiceRate,
		MaxQueueSize:   yamlConfig.Simulation.MaxQueueSize,
		MaxCustomers:   yamlConfig.Simulation.MaxCustomers,
		StopCondition: models.StopCondition{
			AutomaticMode: yamlConfig.Simulation.StopCondition.AutomaticMode,
			Type:          yamlConfig.Simulation.StopCondition.Type,
			Value:         int(yamlConfig.Simulation.StopCondition.Value),
			TimeLimit:     yamlConfig.Simulation.StopCondition.TimeLimit,
		},
		Visualization: models.VisualizationConfig{
			Enabled:             yamlConfig.Simulation.Visualization.Enabled,
			UpdateInterval:      time.Duration(yamlConfig.Simulation.Visualization.UpdateIntervalMs) * time.Millisecond,
			ShowRealtimeMetrics: yamlConfig.Simulation.Visualization.ShowRealtimeMetrics,
			ProgressBarWidth:    yamlConfig.Simulation.Visualization.ProgressBarWidth,
		},
		Random: models.RandomConfig{
			Seed:         yamlConfig.Simulation.Random.Seed,
			Distribution: yamlConfig.Simulation.Random.Distribution,
		},
		Logging: models.LoggingConfig{
			Level:        yamlConfig.Simulation.Logging.Level,
			LogToFile:    yamlConfig.Simulation.Logging.LogToFile,
			LogFilePath:  yamlConfig.Simulation.Logging.LogFilePath,
			OutputFormat: yamlConfig.Simulation.Logging.OutputFormat,
		},
	}
}
