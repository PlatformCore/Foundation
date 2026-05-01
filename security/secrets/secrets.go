package secrets

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"io"
	"sync"
	"time"
)

var ErrSecretNotFound = errors.New("secret not found")

type Secret struct {
	Name      string
	Version   string
	Value     []byte
	CreatedAt time.Time
	ExpiresAt time.Time
}

type Store interface {
	Get(ctx context.Context, name string) (Secret, error)
	Put(ctx context.Context, secret Secret) error
	Delete(ctx context.Context, name string) error
}

type MemoryStore struct {
	mu sync.RWMutex
	m  map[string]Secret
}

func NewMemoryStore() *MemoryStore { return &MemoryStore{m: map[string]Secret{}} }

func (s *MemoryStore) Get(ctx context.Context, name string) (Secret, error) {
	select {
	case <-ctx.Done():
		return Secret{}, ctx.Err()
	default:
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	sec, ok := s.m[name]
	if !ok {
		return Secret{}, ErrSecretNotFound
	}
	if !sec.ExpiresAt.IsZero() && time.Now().After(sec.ExpiresAt) {
		return Secret{}, ErrSecretNotFound
	}
	return sec, nil
}

func (s *MemoryStore) Put(ctx context.Context, sec Secret) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if sec.CreatedAt.IsZero() {
		sec.CreatedAt = time.Now()
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[sec.Name] = sec
	return nil
}

func (s *MemoryStore) Delete(ctx context.Context, name string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, name)
	return nil
}

type Envelope struct {
	KeyID string
	Data  string
}

func Encrypt(keyID string, key, plaintext []byte) (Envelope, error) {
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return Envelope{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return Envelope{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return Envelope{}, err
	}
	out := append(nonce, gcm.Seal(nil, nonce, plaintext, []byte(keyID))...)
	return Envelope{KeyID: keyID, Data: base64.RawStdEncoding.EncodeToString(out)}, nil
}

func Decrypt(key []byte, env Envelope) ([]byte, error) {
	raw, err := base64.RawStdEncoding.DecodeString(env.Data)
	if err != nil {
		return nil, err
	}
	k := sha256.Sum256(key)
	block, err := aes.NewCipher(k[:])
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(raw) < gcm.NonceSize() {
		return nil, errors.New("invalid envelope")
	}
	nonce, data := raw[:gcm.NonceSize()], raw[gcm.NonceSize():]
	return gcm.Open(nil, nonce, data, []byte(env.KeyID))
}
