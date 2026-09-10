package workers

import (
	"github.com/redis/go-redis/v9"

	"log/slog"

	"github.com/max-fletcher/golang_web_server_boilerplate/internal/events"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/posts"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/modules/users"
	"github.com/max-fletcher/golang_web_server_boilerplate/internal/queue"
)

func New(
	redisClient *redis.Client,
	// cacheClient cache.Cache, // Bind and dependency here that you may want to pass down to other workers below
	logger *slog.Logger,
) *queue.Worker {
	dispatcher := queue.NewDispatcher()

	userWorker := users.NewWorker(
	// cacheClient // To bind any depencencies to this worker(if needed), send it here
	)

	postWorker := posts.NewWorker(
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

	return queue.NewWorker(
		redisClient,
		events.DefaultStream,
		"app-workers",
		"worker-1",
		dispatcher.Dispatch,
		logger,
	)
}
