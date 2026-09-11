package workers

import (
	redisClient "github.com/redis/go-redis/v9"

	"log/slog"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
	redis_queue "github.com/max-fletcher/golang_web_server_boilerplate/internal/queue/redis"
)

// func New() does the following:
//  1. Create dispatcher and bind events with corresponding handlers/functions/actions to it(via a map)
//  2. Create new queue.Worker(internal/queue/worker.go)
//  3. Bind to the queue.Worker:
//     a) *redis.Client
//     b) stream
//     c) consumerGroup
//     d) consumer
//     e) dispatcher's Dispatch as handler(very important)
//     f) logger
//     and return it(basically returns a redis worker). This worker(queue.Worker) can be ran(contains a Run method) to spawn
//     a running worker(using redis client) that will listen to events published.
//
// One important thing to keep in mind is that even though we bound the Dispatch method of dispatcher to handler, the dispatch method's
// event-handler bindings isn't gone; it lives somewhere in memory, and when we are calling queue.worker.handler, the bindings are used
// to map out which method/handler to execute
func New(
	redisClient *redisClient.Client,
	// cacheClient cache.Cache, // Bind and dependency here that you may want to pass down to other workers below
	logger *slog.Logger,
) *redis_queue.Worker {
	dispatcher := queue.NewDispatcher()

	userWorker := users.NewWorker(
	// cacheClient // To bind any depencencies to this worker(if needed), send it here
	)

	dispatcher.Register(
		string(events.QueueEventUserCreated),
		userWorker.HandleCreated,
	)

	dispatcher.Register(
		string(events.QueueEventUserUpdated),
		userWorker.HandleUpdated,
	)

	dispatcher.Register(
		string(events.QueueEventUserDeleted),
		userWorker.HandleDeleted,
	)

	postWorker := posts.NewWorker(
	// cacheClient // To bind any depencencies to this worker(if needed), send it here
	)

	dispatcher.Register(
		string(events.QueueEventPostCreated),
		postWorker.HandleCreated,
	)

	dispatcher.Register(
		string(events.QueueEventPostUpdated),
		postWorker.HandleUpdated,
	)

	dispatcher.Register(
		string(events.QueueEventPostDeleted),
		postWorker.HandleDeleted,
	)

	return redis_queue.NewWorker(
		redisClient,
		events.DefaultStream,
		"app-workers",
		"worker-1",
		dispatcher.Dispatch,
		logger,
	)
}
