package otp

import (
	"account-service/internal/domain/secret"
	"account-service/internal/service/account"
	"account-service/pkg/sms"
)

// Configuration is an alias for a function that will take in a pointer to a Service and modify it
type Configuration func(s *Service) error

// Service is an implementation of the Service
type Service struct {
	smsClient        *sms.Client
	secretRepository secret.Repository
	accountService   *account.Service
}

// New takes a variable amount of Configuration functions and returns a new Service
// Each Configuration will be called in the order they are passed in
func New(configs ...Configuration) (s *Service, err error) {
	// Create the service
	s = &Service{}

	// Apply all Configurations passed in
	for _, cfg := range configs {
		// Pass the service into the configuration function
		if err = cfg(s); err != nil {
			return
		}
	}
	return
}

// WithSecretRepository applies a given secret repository to the Service
func WithSecretRepository(secretRepository secret.Repository) Configuration {
	// return a function that matches the Configuration alias,
	// You need to return this so that the parent function can take in all the needed parameters
	return func(s *Service) error {
		s.secretRepository = secretRepository
		return nil
	}
}

// WithAccountService applies a given account service to the Service
func WithAccountService(accountService *account.Service) Configuration {
	// Create the account service, if we needed parameters, such as connection strings they could be inputted here
	return func(s *Service) error {
		s.accountService = accountService
		return nil
	}
}

// WithSMSClient applies a given sms client to the Service
func WithSMSClient(smsClient *sms.Client) Configuration {
	// Create the sms client, if we needed parameters, such as connection strings they could be inputted here
	return func(s *Service) error {
		s.smsClient = smsClient
		return nil
	}
}
