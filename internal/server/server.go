package server

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/cthulhu-platform/auth/internal/handlers"
	"github.com/cthulhu-platform/auth/internal/service"
	"github.com/cthulhu-platform/common/pkg/messages"
	"github.com/wagslane/go-rabbitmq"
)

const (
	AuthExchange = "auth"
)

const (
	RoutingKeyGenerateTokens     = "auth.generate_tokens"
	RoutingKeyValidateToken     = "auth.validate_token"
	RoutingKeyValidateUserID    = "auth.validate_user"
)

type RMQServer struct {
	Conn    *rabbitmq.Conn
	Service service.AuthenticationService
}

// NewRMQServer creates a new RabbitMQ server instance
func NewRMQServer(conn *rabbitmq.Conn, s service.AuthenticationService) *RMQServer {
	return &RMQServer{
		Conn:    conn,
		Service: s,
	}
}

// Start sets up and starts all consumers
func (s *RMQServer) Start() {
	// Create diagnose consumer
	diagnoseConsumer, err := rabbitmq.NewConsumer(
		s.Conn,
		"auth.diagnose",
		rabbitmq.WithConsumerOptionsRoutingKey(messages.TopicDiagnoseServicesAll),
		rabbitmq.WithConsumerOptionsExchangeName(messages.DiagnoseExchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsConsumerName("auth_diagnose"),
	)
	if err != nil {
		log.Fatalf("Failed to create diagnose consumer: %v", err)
	}
	defer diagnoseConsumer.Close()

	// Create GenerateTokens consumer
	generateTokensConsumer, err := rabbitmq.NewConsumer(
		s.Conn,
		"auth.generate_tokens",
		rabbitmq.WithConsumerOptionsRoutingKey(RoutingKeyGenerateTokens),
		rabbitmq.WithConsumerOptionsExchangeName(AuthExchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsConsumerName("auth_generate_tokens"),
	)
	if err != nil {
		log.Fatalf("Failed to create GenerateTokens consumer: %v", err)
	}
	defer generateTokensConsumer.Close()

	// Create ValidateAccessToken consumer
	validateTokenConsumer, err := rabbitmq.NewConsumer(
		s.Conn,
		"auth.validate_token",
		rabbitmq.WithConsumerOptionsRoutingKey(RoutingKeyValidateToken),
		rabbitmq.WithConsumerOptionsExchangeName(AuthExchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsConsumerName("auth_validate_token"),
	)
	if err != nil {
		log.Fatalf("Failed to create ValidateAccessToken consumer: %v", err)
	}
	defer validateTokenConsumer.Close()

	// Create ValidateUserID consumer
	validateUserIDConsumer, err := rabbitmq.NewConsumer(
		s.Conn,
		"auth.validate_user",
		rabbitmq.WithConsumerOptionsRoutingKey(RoutingKeyValidateUserID),
		rabbitmq.WithConsumerOptionsExchangeName(AuthExchange),
		rabbitmq.WithConsumerOptionsExchangeDeclare,
		rabbitmq.WithConsumerOptionsExchangeKind("topic"),
		rabbitmq.WithConsumerOptionsExchangeDurable,
		rabbitmq.WithConsumerOptionsQueueDurable,
		rabbitmq.WithConsumerOptionsConsumerName("auth_validate_user"),
	)
	if err != nil {
		log.Fatalf("Failed to create ValidateUserID consumer: %v", err)
	}
	defer validateUserIDConsumer.Close()

	log.Printf("Auth service started and listening for messages:")
	log.Printf("  - Diagnose: exchange=%s, routing_key=%s", messages.DiagnoseExchange, messages.TopicDiagnoseServicesAll)
	log.Printf("  - GenerateTokens: exchange=%s, routing_key=%s", AuthExchange, RoutingKeyGenerateTokens)
	log.Printf("  - ValidateAccessToken: exchange=%s, routing_key=%s", AuthExchange, RoutingKeyValidateToken)
	log.Printf("  - ValidateUserID: exchange=%s, routing_key=%s", AuthExchange, RoutingKeyValidateUserID)

	// Setup graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigs
		log.Printf("Received signal: %v, stopping consumers...", sig)
		diagnoseConsumer.Close()
		generateTokensConsumer.Close()
		validateTokenConsumer.Close()
		validateUserIDConsumer.Close()
	}()

	// Start all consumers in goroutines
	go func() {
		log.Println("Starting diagnose consumer...")
		if err := diagnoseConsumer.Run(handlers.HandleDiagnoseMessage(s.Service)); err != nil {
			log.Fatalf("Diagnose consumer error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting GenerateTokens consumer...")
		if err := generateTokensConsumer.Run(handlers.HandleGenerateTokens(s.Service)); err != nil {
			log.Fatalf("GenerateTokens consumer error: %v", err)
		}
	}()

	go func() {
		log.Println("Starting ValidateAccessToken consumer...")
		if err := validateTokenConsumer.Run(handlers.HandleValidateAccessToken(s.Service)); err != nil {
			log.Fatalf("ValidateAccessToken consumer error: %v", err)
		}
	}()

	// Block main thread - wait for messages (Run is blocking)
	log.Println("Starting ValidateUserID consumer and waiting for messages...")
	if err := validateUserIDConsumer.Run(handlers.HandleValidateUserID(s.Service)); err != nil {
		log.Fatalf("ValidateUserID consumer error: %v", err)
	}

	log.Println("Shutting down gracefully...")
}

