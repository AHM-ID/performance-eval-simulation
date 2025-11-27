package simulation

import (
	"container/heap"
	"math"
	"math/rand"
	"time"
	"des/models"
)

type eventHeap []*models.Event

func (h eventHeap) Len() int           { return len(h) }
func (h eventHeap) Less(i, j int) bool { return h[i].Timestamp < h[j].Timestamp }
func (h eventHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *eventHeap) Push(x interface{}) {
	*h = append(*h, x.(*models.Event))
}

func (h *eventHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

type EventList struct {
	events eventHeap
}

func NewEventList() *EventList {
	h := make(eventHeap, 0)
	heap.Init(&h)
	return &EventList{events: h}
}

func (el *EventList) Push(event *models.Event) {
	heap.Push(&el.events, event)
}

func (el *EventList) Pop() *models.Event {
	if el.IsEmpty() {
		return nil
	}
	return heap.Pop(&el.events).(*models.Event)
}

func (el *EventList) Peek() *models.Event {
	if el.IsEmpty() {
		return nil
	}
	return el.events[0]
}

func (el *EventList) IsEmpty() bool {
	return len(el.events) == 0
}

func (el *EventList) Size() int {
	return len(el.events)
}

type EventManager struct {
	eventList *EventList
	rng       *rand.Rand
	config    *models.SimulationConfig
}

func NewEventManager(config *models.SimulationConfig) *EventManager {
	seed := config.Random.Seed
	if seed == -1 {
		seed = time.Now().UnixNano()
	}

	return &EventManager{
		eventList: NewEventList(),
		rng:       rand.New(rand.NewSource(seed)),
		config:    config,
	}
}

func (em *EventManager) ScheduleEvent(eventType models.EventType, timestamp float64, customer *models.Customer) {
	event := &models.Event{
		Type:      eventType,
		Timestamp: timestamp,
		Customer:  customer,
	}
	em.eventList.Push(event)
}

func (em *EventManager) GetNextEvent() *models.Event {
	return em.eventList.Pop()
}

func (em *EventManager) PeekNextEvent() *models.Event {
	return em.eventList.Peek()
}

func (em *EventManager) GenerateExponential(rate float64) float64 {
	if rate <= 0 {
		return 1.0
	}
	u := em.rng.Float64()
	for u == 0.0 || u == 1.0 {
		u = em.rng.Float64()
	}
	return -math.Log(1.0-u) / rate
}

func (em *EventManager) GenerateUniform(min, max float64) float64 {
	return min + em.rng.Float64()*(max-min)
}

func (em *EventManager) GenerateConstant(value float64) float64 {
	return value
}

func (em *EventManager) GetInterarrivalTime() float64 {
	switch em.config.Random.Distribution {
	case "uniform":
		return em.GenerateUniform(0.5/em.config.ArrivalRate, 1.5/em.config.ArrivalRate)
	case "constant":
		return 1.0 / em.config.ArrivalRate
	default:
		return em.GenerateExponential(em.config.ArrivalRate)
	}
}

func (em *EventManager) GetServiceTime() float64 {
	switch em.config.Random.Distribution {
	case "uniform":
		return em.GenerateUniform(0.5/em.config.ServiceRate, 1.5/em.config.ServiceRate)
	case "constant":
		return 1.0 / em.config.ServiceRate
	default:
		return em.GenerateExponential(em.config.ServiceRate)
	}
}

func (em *EventManager) HasEvents() bool {
	return !em.eventList.IsEmpty()
}

func (em *EventManager) ClearEvents() {
	em.eventList = NewEventList()
}

func (em *EventManager) GetEventCount() int {
	return em.eventList.Size()
}