package rabbitmq

import (
	"log"
	"strings"
	"sync"
	"time"

	"github.com/streadway/amqp"
	"golang.org/x/time/rate"
)

type WorkerPool struct {
	QueueName   string
	RateLimiter *rate.Limiter
	Workers     int
	MaxWorkers  int
}

func NewWorkerPool(queueName string, rateLimit *rate.Limiter, workers int, maxWorkers int) *WorkerPool {
	return &WorkerPool{
		QueueName:   queueName,
		RateLimiter: rateLimit,
		Workers:     workers,
		MaxWorkers:  maxWorkers,
	}
}

func Scheduler(pools []*WorkerPool, ch *amqp.Channel) {
	var wg sync.WaitGroup
	ch.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	for _, pool := range pools {
		wg.Add(1)
		go func(pool *WorkerPool) {
			defer wg.Done()
			for i := 0; i < pool.Workers; i++ {
				go Worker(pool.QueueName, pool.RateLimiter, ch)
			}
		}(pool)
	}

	// Dynamic priority adjustment loop
	go func() {
		for {
			time.Sleep(7 * time.Second) // Adjust priorities every 10 seconds

			for _, pool := range pools {
				q, err := ch.QueueInspect(pool.QueueName)
				if err != nil {
					log.Printf("Failed to inspect queue %s: %v", pool.QueueName, err)
					continue
				}

				if q.Messages == 0 && pool.Workers > 1 {
					log.Printf("Queue %s is idle, reallocating workers", pool.QueueName)
					// Reallocate workers to another pool (this is a simplified example)
					reallocateWorkers(pools, ch)
				} else {
					log.Println("Queue is not idle, scaling workers...")
					scaleWorkers(pool, ch)
				}
			}
			println("Total messages: ", getTotalMessages(getQueueStats(pools, ch)))
			println()
		}
	}()

	wg.Wait()
}

func scaleWorkers(pool *WorkerPool, ch *amqp.Channel) {
	q, err := ch.QueueInspect(pool.QueueName)
	if err != nil {
		log.Printf("Failed to inspect queue %s: %v", pool.QueueName, err)
		return
	}

	messagesInQueue := q.Messages
	// required_workers = messagesInQueue / 5
	newWorkers := 1 + (messagesInQueue / 5)

	if newWorkers > pool.MaxWorkers {
		newWorkers = pool.MaxWorkers
	}

	pool.Workers = newWorkers

	log.Printf("Scaled workers for queue %s to %d", pool.QueueName, pool.Workers)
}

func getQueueStats(pools []*WorkerPool, ch *amqp.Channel) map[string]int {
	// Gather queue statistics
	queueStats := make(map[string]int)
	for _, pool := range pools {
		q, err := ch.QueueInspect(pool.QueueName)
		if err != nil {
			log.Printf("Failed to inspect queue %s: %v", pool.QueueName, err)
			continue
		}
		queueStats[pool.QueueName] = q.Messages
	}

	return queueStats
}

func getTotalMessages(queueStats map[string]int) int {
	totalMessages := 0
	for _, messages := range queueStats {
		totalMessages += messages
	}
	return totalMessages
}

func reallocateWorkers(pools []*WorkerPool, ch *amqp.Channel) {

	queueStats := getQueueStats(pools, ch)

	totalMessages := getTotalMessages(queueStats)
	println("Total messages: ", totalMessages)

	println("Reallocating workers...")
	// Adjust worker allocation based on queue backlog and priority
	for _, pool := range pools {
		if totalMessages == 0 {
			pool.Workers = 1 // At least one worker per pool
		} else {
			// Adjust based on the proportion of messages in the queue
			// and the priority level of the pool
			priorityWeight := getPriorityWeight(pool.QueueName)
			pool.Workers = int(float64(pool.Workers) * (float64(queueStats[pool.QueueName]) / float64(totalMessages)) * priorityWeight)

			// Ensure at least one worker is always allocated to avoid starvation
			if pool.Workers < 1 {
				pool.Workers = 1
			}

			// Ensure a maximum of 10 workers per pool
			if pool.Workers > pool.MaxWorkers {
				pool.Workers = pool.MaxWorkers
			}

		}

		log.Printf("Reallocated %d workers to queue %s", pool.Workers, pool.QueueName)

	}

	println("Workers reallocated")
	println()
}

func getPriorityWeight(queueName string) float64 {
	switch {
	case strings.Contains(queueName, "high"):
		return 1.5
	case strings.Contains(queueName, "medium"):
		return 1.0
	case strings.Contains(queueName, "low"):
		return 0.5
	default:
		return 1.0
	}
}
