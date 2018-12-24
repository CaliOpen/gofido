package store

import (
	"github.com/tstranex/u2f"
)

// Define an interface for all store operations
type StoreInterface interface {
	NewChallenge(string) (u2f.Challenge, error)
	GetChallenge(string, string) (u2f.Challenge, error)
	NewRegistration(string, u2f.Challenge, u2f.RegisterResponse) (*u2f.Registration, error)
	GetRegistrations(string) ([]u2f.Registration, error)
	InsertKeyChallenge(string, []byte, u2f.Challenge) error
	GetKeyChallenges(string, []byte) []KeyChallenge
	GetKeyCounter(string, []byte) (KeyCounter, error)
	UpdateCounter(string, []byte, uint32) error
}
