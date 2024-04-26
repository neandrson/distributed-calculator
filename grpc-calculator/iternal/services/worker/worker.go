package worker

import (
	"culc/iternal/config"
	"culc/iternal/model"
	"fmt"
	"math"
	"strconv"
	"sync"
	"time"
)

type Worker struct {
	ID     int
	Status string
}

type Client struct {
	ID              string
	ActiveWorker    int
	Workers         []*Worker
	workerAvailable chan struct{}
	cfg             config.ConfigTime
}

type CalculatorServer struct {
	clients    map[string]*Client
	workerPool chan *Worker
	mu         sync.Mutex
}

func NewWorker(id int) *Worker {
	return &Worker{
		ID:     id,
		Status: "active",
	}
}

func NewClient(clientID string, workerCount int) *Client {
	conf := config.LoadConfigTime("./config/time.yaml")
	client := &Client{
		ID:              clientID,
		ActiveWorker:    workerCount,
		Workers:         make([]*Worker, workerCount),
		workerAvailable: make(chan struct{}),
		cfg:             *conf,
	}

	for i := 0; i < workerCount; i++ {
		client.Workers[i] = NewWorker(i)
	}

	return client
}

func NewCalculatorServer(initialWorkerCount int) *CalculatorServer {

	return &CalculatorServer{
		clients:    make(map[string]*Client),
		workerPool: make(chan *Worker, initialWorkerCount),
	}
}

func (s *CalculatorServer) AddClient(clientID string, workerCount int) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.clients[clientID]; !ok {
		client := NewClient(clientID, workerCount)
		s.clients[clientID] = client
	}
}
func (s *Client) CountEx(num1 string, num2 string, oper string) (float64, error) {

	num1Float, err := strconv.ParseFloat(num1, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert num1 to float: %v", err)
	}
	num2Float, err := strconv.ParseFloat(num2, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to convert num2 to float: %v", err)
	}
	s.SetWorkerStatusBusy()
	defer s.SetWorkerStatusActive()
	res, err := s.CalcExpression(num1Float, oper, num2Float)
	if err != nil {
		s.SetWorkerStatusActive()
		return 0, fmt.Errorf("ошибка при расчете: %e", err)
	}

	return res, nil
}
func (s *CalculatorServer) GetClientWithMaxFreeWorkers() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.clients) == 0 {
		s.AddClient("1", 8)
	}
	maxFreeWorkerCount := 0
	var clientIDWithMaxFreeWorkers string

	for clientID, client := range s.clients {
		freeWorkerCount := 0
		for _, worker := range client.Workers {
			if worker.Status == "active" {
				freeWorkerCount++
			}
		}
		if freeWorkerCount > maxFreeWorkerCount {
			maxFreeWorkerCount = freeWorkerCount
			clientIDWithMaxFreeWorkers = clientID
		}
	}

	return clientIDWithMaxFreeWorkers
}
func (s *CalculatorServer) GetClientByID(clientID string) *Client {
	s.mu.Lock()
	defer s.mu.Unlock()

	if client, ok := s.clients[clientID]; ok {
		return client
	}

	return nil
}
func (s *Client) CalcExpression(Num1 float64, operator string, Num2 float64) (float64, error) {
	res := 0.0
	var Error error
	switch operator {
	case "/":

		if Num2 == 0.0 {
			Error = fmt.Errorf("встечено деление на ноль")
		}
		fmt.Println(s.cfg.Division)
		res = Num1 / Num2
		time.Sleep(s.cfg.Division)
	case "*":

		res = Num1 * Num2
		time.Sleep(s.cfg.MultiP)
	case "-":

		res = Num1 - Num2
		time.Sleep(s.cfg.Minus)
	case "+":

		res = Num1 + Num2
		time.Sleep(s.cfg.Plus)
	case "^":

		res = math.Pow(Num1, Num2)
		time.Sleep(s.cfg.Exponent)
	default:
		Error = fmt.Errorf("оператор не найден")
	}
	if Error != nil {
		return 0, Error
	}
	return res, nil
}
func (c *Client) SetWorkerStatusBusy() {

	for _, worker := range c.Workers {

		if worker.Status == "active" {
			worker.Status = "busy"
			return
		}
	}

	<-c.workerAvailable

	for _, worker := range c.Workers {
		if worker.Status == "active" {
			worker.Status = "busy"
			return
		}
	}
}
func (s *CalculatorServer) GetFreeWorkerCounts() map[string]model.StatusClient {
	s.mu.Lock()
	defer s.mu.Unlock()

	freeWorkerCounts := make(map[string]model.StatusClient)

	for clientID, client := range s.clients {

		freeCount := 0
		AllWorkers := 0
		for _, worker := range client.Workers {
			AllWorkers++
			if worker.Status == "active" {
				freeCount++
			}
		}

		freeWorkerCounts[clientID] = model.StatusClient{ActiveWorkers: freeCount, AllWorkers: AllWorkers}
	}

	return freeWorkerCounts
}

func (c *Client) SetWorkerStatusActive() {

	for _, worker := range c.Workers {
		if worker.Status == "busy" {
			worker.Status = "active"
			select {
			case c.workerAvailable <- struct{}{}:

			case <-time.After(2 * time.Second):

				return
			}

		}
	}
}
func (s *CalculatorServer) UpdateTime(cfg config.ConfigTime) {

	for _, client := range s.clients {
		client.cfg = cfg
	}
}
func (c *Client) GetConfig() config.ConfigTime {
	return c.cfg
}
