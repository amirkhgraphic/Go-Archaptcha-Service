package services

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"strings"
	"sync"
	"time"
)

var (
	ErrChallengeEmpty   = errors.New("challenge_id is required")
	ErrChallengeInvalid = errors.New("challenge_id is invalid or expired")
	ErrChallengeNetwork = errors.New("temporary arcaptcha network issue")
)

// ArcaptchaService is an in-memory fake arcaptcha validator used for local testing.
// It can mint one-time challenges and validate them with optional network/error simulation.
type ArcaptchaService struct {
	mu          sync.Mutex
	challenges  map[string]challengeInfo
	validForSec int64
}

type challengeInfo struct {
	createdAt time.Time
}

func NewArcaptchaService() *ArcaptchaService {
	return &ArcaptchaService{
		challenges:  make(map[string]challengeInfo),
		validForSec: 10 * 60, // 10 minutes
	}
}

// Arcaptcha is a shared singleton used across handlers.
var Arcaptcha = NewArcaptchaService()

// GenerateChallenge returns a new challenge token ready to be validated later.
func (s *ArcaptchaService) GenerateChallenge() string {
	token := randomToken()
	s.mu.Lock()
	s.challenges[token] = challengeInfo{createdAt: time.Now()}
	s.mu.Unlock()
	return token
}

// RegisterChallenge allows seeding a predictable token for tests.
func (s *ArcaptchaService) RegisterChallenge(token string) {
	if strings.TrimSpace(token) == "" {
		return
	}
	s.mu.Lock()
	s.challenges[token] = challengeInfo{createdAt: time.Now()}
	s.mu.Unlock()
}

// PeekChallenge validates a token without consuming it (used by the fake verify endpoint).
func (s *ArcaptchaService) PeekChallenge(challengeID string) error {
	return s.validate(challengeID, false)
}

// ValidateChallenge validates and consumes a token (used by protected endpoints).
func (s *ArcaptchaService) ValidateChallenge(challengeID string) error {
	return s.validate(challengeID, true)
}

func (s *ArcaptchaService) validate(challengeID string, consume bool) error {
	if strings.TrimSpace(challengeID) == "" {
		return ErrChallengeEmpty
	}

	// Simple hook to mimic upstream hiccups.
	if strings.HasSuffix(challengeID, "-neterr") {
		return ErrChallengeNetwork
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	info, ok := s.challenges[challengeID]
	if !ok {
		return ErrChallengeInvalid
	}

	// Expire old challenges to avoid unbounded growth.
	if s.validForSec > 0 && time.Since(info.createdAt).Seconds() > float64(s.validForSec) {
		delete(s.challenges, challengeID)
		return ErrChallengeInvalid
	}

	if consume {
		delete(s.challenges, challengeID)
	}
	return nil
}

// ActiveChallenges exposes the number of available tokens (handy for debugging/tests).
func (s *ArcaptchaService) ActiveChallenges() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.challenges)
}

func randomToken() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return "arcaptcha_" + hex.EncodeToString(buf)
}
