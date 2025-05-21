package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/HdrHistogram/hdrhistogram-go"
	pb "github.com/uppercasee/drls/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gopkg.in/yaml.v3"
)

type LatencyStats struct {
	Min   time.Duration `json:"min"`
	Max   time.Duration `json:"max"`
	Avg   time.Duration `json:"avg"`
	P50   time.Duration `json:"p50"`
	P95   time.Duration `json:"p95"`
	P99   time.Duration `json:"p99"`
	Count int64         `json:"count"`
}

type ClientStats struct {
	ClientID     string       `json:"client_id"`
	Allowed      int          `json:"allowed"`
	Denied       int          `json:"denied"`
	Errors       int          `json:"errors"`
	LatencyStats LatencyStats `json:"latency_stats"`
}

type BenchmarkResults struct {
	Clients []ClientStats `json:"clients"`
}

type LoadConfig struct {
	Clients      int    `yaml:"clients"`
	RPSPerClient int    `yaml:"rps_per_client"`
	DurationSec  int    `yaml:"duration_seconds"`
	GRPCTarget   string `yaml:"grpc_target"`
	Pattern      string `yaml:"pattern"`

	Burst struct {
		SizeMultiplier int `yaml:"size_multiplier"`
		IntervalMs     int `yaml:"interval_ms"`
	} `yaml:"burst"`

	Ramp struct {
		Step       int `yaml:"step"`
		IntervalMs int `yaml:"interval_ms"`
	} `yaml:"ramp"`
}

func main() {
	// load config
	var cfg LoadConfig
	data, err := os.ReadFile("benchmark/config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse config: %v", err)
	}

	grpcAddr := cfg.GRPCTarget
	numClients := cfg.Clients
	rpsPerClient := cfg.RPSPerClient
	testDuration := time.Duration(cfg.DurationSec) * time.Second

	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewRateLimiterServiceClient(conn)
	var wg sync.WaitGroup

	start := time.Now()
	end := start.Add(testDuration)

	clientStatsMap := make(map[string]*ClientStats)
	histograms := make(map[string]*hdrhistogram.Histogram)

	for i := range numClients {
		clientID := fmt.Sprintf("client-%d", i)
		clientStatsMap[clientID] = &ClientStats{ClientID: clientID}
		histograms[clientID] = hdrhistogram.New(1, 10_000_000_000, 3)
	}

	var mu sync.Mutex

	for i := range numClients {
		wg.Add(1)
		clientID := fmt.Sprintf("client-%d", i)
		go func(clientID string) {
			defer wg.Done()

			switch cfg.Pattern {
			case "burst":
				burstSize := rpsPerClient * cfg.Burst.SizeMultiplier
				burstInterval := time.Duration(cfg.Burst.IntervalMs) * time.Millisecond
				ticker := time.NewTicker(burstInterval)
				defer ticker.Stop()

				for now := range ticker.C {
					if now.After(end) {
						return
					}
					for range burstSize {
						go sendRequest(client, clientID, clientStatsMap, histograms, &mu)
					}
				}

			case "ramp":
				rampStep := cfg.Ramp.Step
				rampInterval := time.Duration(cfg.Ramp.IntervalMs) * time.Millisecond
				currentRPS := 1
				maxRPS := rpsPerClient

				for now := time.Now(); now.Before(end); now = time.Now() {
					for range currentRPS {
						go sendRequest(client, clientID, clientStatsMap, histograms, &mu)
					}
					if currentRPS < maxRPS {
						currentRPS += rampStep
						if currentRPS > maxRPS {
							currentRPS = maxRPS
						}
					}
					time.Sleep(rampInterval)
				}

			default: // "steady"
				ticker := time.NewTicker(time.Second / time.Duration(rpsPerClient))
				defer ticker.Stop()

				for now := range ticker.C {
					if now.After(end) {
						return
					}
					sendRequest(client, clientID, clientStatsMap, histograms, &mu)
				}
			}
		}(clientID)
	}

	wg.Wait()

	for clientID, stats := range clientStatsMap {
		h := histograms[clientID]
		if h.TotalCount() > 0 {
			stats.LatencyStats = LatencyStats{
				Min:   time.Duration(h.Min()) / time.Millisecond,
				Max:   time.Duration(h.Max()) / time.Millisecond,
				Avg:   time.Duration(h.Mean()) / time.Millisecond,
				P50:   time.Duration(h.ValueAtQuantile(50)) / time.Millisecond,
				P95:   time.Duration(h.ValueAtQuantile(95)) / time.Millisecond,
				P99:   time.Duration(h.ValueAtQuantile(99)) / time.Millisecond,
				Count: h.TotalCount(),
			}
		}
	}

	results := BenchmarkResults{Clients: []ClientStats{}}
	for _, stats := range clientStatsMap {
		results.Clients = append(results.Clients, *stats)
	}

	outputPath := fmt.Sprintf("benchmark/results/results-3-%s.json", cfg.Pattern)
	os.MkdirAll("benchmark/results", os.ModePerm)
	file, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create result file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(results); err != nil {
		log.Fatalf("Failed to write results to JSON: %v", err)
	}

	fmt.Println("🔍 Benchmark results written to:", outputPath)
}

func sendRequest(client pb.RateLimiterServiceClient, clientID string, clientStatsMap map[string]*ClientStats, histograms map[string]*hdrhistogram.Histogram, mu *sync.Mutex) {
	start := time.Now()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	res, err := client.Check(ctx, &pb.CheckRequest{ClientId: clientID})
	duration := time.Since(start)

	mu.Lock()
	defer mu.Unlock()

	// Update stats
	stats := clientStatsMap[clientID]
	if err != nil {
		stats.Errors++
	} else if res.Allowed {
		stats.Allowed++
	} else {
		stats.Denied++
	}

	// Record latency in histogram in microseconds (hdrhistogram expects int64)
	histograms[clientID].RecordValue(int64(duration))
}
