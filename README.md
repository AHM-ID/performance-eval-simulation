# Discrete Event Simulation Engine

This repository contains a modular, high-performance discrete event simulation (DES) framework written in Go. It is designed to be extensible, configurable, and suitable for research, academic use, and performance-driven simulation environments.

The project implements a classic single-server queueing model with optional real-time visualization, detailed logging, configurable randomness, and pluggable simulation components.

---

## Table of Contents

1. Overview
2. Core Architecture
3. Event Processing Pipeline
4. Configuration System
5. Simulation Flow
6. Logging System
7. Visualization Engine
8. Statistical Metrics
9.  How to Run

---

# 1. Overview

This project implements a discrete event simulator built around a priority-queue event list. The engine supports multiple distribution models, customizable stopping conditions, pluggable visualization, and file-based logging with timestamped output.

The simulation models a queueing system with the following characteristics:

* A single server
* Random arrival and service processes
* Configurable queue capacity
* Event-driven execution
* Real-time console-based visualization
* Statistical analysis of performance metrics

---

# 2. Core Architecture

The architecture follows a modular pattern composed of the following key components:

### Event Manager

Handles:

* Event scheduling
* Managing the priority queue
* Generating interarrival and service times
* Configurable random distribution (exponential, uniform, constant)

### Simulator

Responsible for:

* Executing events in chronological order
* Managing queue and server state
* Updating time, metrics, and stopping conditions

### Statistics Collector

Collects continuous and discrete metrics:

* Waiting times
* System times
* Queue lengths
* Server utilization
* Throughput
* Time-dependent averages

### Logging Engine

Provides configurable file/text logging:

* JSON or text output mode

### Visualization Engine

Renders real-time simulation state:

* Server status
* Queue state
* Next event
* Metrics
* Progress bar

---

# 3. Event Processing Pipeline

Events are stored inside a binary heap prioritizing the earliest timestamp. Each event contains:

* Event type (arrival or departure)
* Timestamp
* Customer reference

Processing loop:

1. Retrieve next event from the event list
2. Advance simulation time
3. Apply the event logic (arrival or departure)
4. Update statistics
5. Render state (if enabled)
6. Check stop conditions

The simulation continues until one of the following occurs:

* Time limit reached
* Maximum customer count reached
* Event list is exhausted

---

# 4. Configuration System

Configuration is loaded from a YAML file and mapped into a structured `SimulationConfig` object.

### Supported Fields

* Simulation parameters (arrival rate, service rate, queue size)
* Random settings and seed
* Visualization options
* Logging settings

---

# 5. Simulation Flow

### Initialization

* Load configuration
* Create event manager
* Initialize statistics and visualization modules
  
### Event Loop

* Process next event
* Handle arrival or departure
* Update queue and server
* Schedule future events
* Log state changes
* Render visualization

### Termination

* Triggered when stop condition is met
* Final metrics are computed
* Results printed to terminal and logs

---

# 6. Logging System

The logging engine supports:

* Text or JSON formats
* Automatic timestamped file naming
* Creation of missing directories
* Configurable log levels

---

# 7. Visualization Engine

The terminal visualization displays:

* Server status
* Queue state
* Next event information
* Real-time metrics
* Progress bar

Example output:

```
SERVER STATUS: BUSY      CUSTOMERS SERVED:    105
QUEUE LENGTH:  19/20     REJECTED:            82
NEXT EVENT: DEPARTURE    at TIME:   100.28
QUEUE: [CCCCCCCCCCCCCCCCC]
PROGRESS: [############################################]  98%
```

The rendering is efficiently managed to avoid flickering.

---

# 8. Statistical Metrics

The simulation collects the following performance measurements:

* Average wait time
* Average time in system
* Server utilization
* Throughput
* Queue probability
* Rejection probability
* Maximum queue length
* Variance of wait and system times
* Percentiles (50th, 75th, 90th, 95th)
* Confidence intervals

---

# 9. How to Run

### Build

```
go build -o sim
```

### Run (Automatic Mode)

```
./sim -config config.yml -mode=automatic
```

### Run (Manual Mode)

```
./sim -config config.yml -mode=manual
```
