package handlers

import (
	"encoding/json"
	"log"

	"github.com/cthulhu-platform/auth/internal/service"
	"github.com/cthulhu-platform/common/pkg/messages"
	"github.com/wagslane/go-rabbitmq"
)

// HandleGenerateTokens processes token generation requests from RabbitMQ
func HandleGenerateTokens(s service.AuthenticationService) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("Handler received GenerateTokens request - RoutingKey: %s, Body: %s", d.RoutingKey, string(d.Body))

		// Unmarshal the request
		var req struct {
			UserID   string `json:"user_id"`
			Email    string `json:"email"`
			Provider string `json:"provider"`
		}
		if err := json.Unmarshal(d.Body, &req); err != nil {
			log.Printf("Failed to unmarshal GenerateTokens request: %v", err)
			return rabbitmq.NackRequeue
		}

		// Call service method
		tokenPair, err := s.GenerateTokens(req.UserID, req.Email, req.Provider)
		if err != nil {
			log.Printf("Failed to generate tokens: %v", err)
			return rabbitmq.NackRequeue
		}

		// Marshal response
		response, err := json.Marshal(tokenPair)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return rabbitmq.NackRequeue
		}

		log.Printf("Generated tokens successfully")
		// TODO: Publish response to reply queue
		_ = response
		return rabbitmq.Ack
	}
}

// HandleValidateAccessToken processes token validation requests from RabbitMQ
func HandleValidateAccessToken(s service.AuthenticationService) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("Handler received ValidateAccessToken request - RoutingKey: %s, Body: %s", d.RoutingKey, string(d.Body))

		// Unmarshal the request
		var req struct {
			Token string `json:"token"`
		}
		if err := json.Unmarshal(d.Body, &req); err != nil {
			log.Printf("Failed to unmarshal ValidateAccessToken request: %v", err)
			return rabbitmq.NackRequeue
		}

		// Call service method
		claims, err := s.ValidateAccessToken(req.Token)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			return rabbitmq.NackRequeue
		}

		// Marshal response
		response, err := json.Marshal(claims)
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return rabbitmq.NackRequeue
		}

		log.Printf("Token validated successfully")
		// TODO: Publish response to reply queue
		_ = response
		return rabbitmq.Ack
	}
}

// HandleValidateUserID processes user ID validation requests from RabbitMQ
func HandleValidateUserID(s service.AuthenticationService) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("Handler received ValidateUserID request - RoutingKey: %s, Body: %s", d.RoutingKey, string(d.Body))

		// Unmarshal the request
		var req struct {
			UserID string `json:"user_id"`
		}
		if err := json.Unmarshal(d.Body, &req); err != nil {
			log.Printf("Failed to unmarshal ValidateUserID request: %v", err)
			return rabbitmq.NackRequeue
		}

		// Call service method
		valid, err := s.ValidateUserID(req.UserID)
		if err != nil {
			log.Printf("Failed to validate user ID: %v", err)
			return rabbitmq.NackRequeue
		}

		// Marshal response
		response, err := json.Marshal(struct {
			Valid bool `json:"valid"`
		}{Valid: valid})
		if err != nil {
			log.Printf("Failed to marshal response: %v", err)
			return rabbitmq.NackRequeue
		}

		log.Printf("User ID validation completed: %v", valid)
		// TODO: Publish response to reply queue
		_ = response
		return rabbitmq.Ack
	}
}

// HandleDiagnoseMessage processes diagnose messages from RabbitMQ
func HandleDiagnoseMessage(s service.AuthenticationService) rabbitmq.Handler {
	return func(d rabbitmq.Delivery) rabbitmq.Action {
		log.Printf("Handler received message - RoutingKey: %s, Body: %s", d.RoutingKey, string(d.Body))

		// Unmarshal the message
		var msg messages.DiagnoseMessage
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			log.Printf("Failed to unmarshal diagnose message: %v", err)
			return rabbitmq.NackRequeue
		}

		// Print the diagnosis check message
		log.Printf("Diagnosis check from routing key: %s", d.RoutingKey)
		log.Printf("Received message - TransactionID: %s, Operation: %s, Message: %s",
			msg.TransactionID, msg.Operation, msg.Message)

		// Process the message through the service
		if err := s.HandleDiagnoseMessage(msg.TransactionID, msg.Operation, msg.Message); err != nil {
			log.Printf("Failed to handle diagnose message: %v", err)
			return rabbitmq.NackRequeue
		}

		return rabbitmq.Ack
	}
}

