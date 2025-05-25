package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CacheEntry struct {
	Key       string        `json:"key"`
	Data      any           `json:"data"`
	Timestamp time.Time     `json:"timestamp"`
	TTL       time.Duration `json:"ttl"`
}

type Cache struct {
	entries map[string]*CacheEntry
	dir     string
}

func LoadCache() (*Cache, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("error finding home directory: %w", err)
	}

	cacheHome := os.Getenv("XDG_CACHE_HOME")
	if cacheHome == "" {
		cacheHome = filepath.Join(home, ".cache")
	}
	cacheDir := filepath.Join(cacheHome, "gh-ci")

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	c := &Cache{
		entries: make(map[string]*CacheEntry),
		dir:     cacheDir,
	}

	cacheFile := filepath.Join(cacheDir, "cache.json")
	if data, err := os.ReadFile(cacheFile); err == nil {
		json.Unmarshal(data, &c.entries)
	}

	return c, nil
}

func (c *Cache) Get(key string) (any, bool) {
	hashedKey := c.hashKey(key)
	entry, exists := c.entries[hashedKey]

	if !exists {
		return nil, false
	}

	if time.Since(entry.Timestamp) > entry.TTL {
		delete(c.entries, hashedKey)
		c.save()
		return nil, false
	}

	return entry.Data, true
}

func (c *Cache) Set(key string, data any, ttl time.Duration) error {
	hashedKey := c.hashKey(key)
	c.entries[hashedKey] = &CacheEntry{
		Key:       key,
		Data:      data,
		Timestamp: time.Now(),
		TTL:       ttl,
	}

	return c.save()
}

func (c *Cache) Delete(key string) {
	hashedKey := c.hashKey(key)
	delete(c.entries, hashedKey)
	c.save()
}

func (c *Cache) Clear() error {
	c.entries = make(map[string]*CacheEntry)
	return c.save()
}

func (c *Cache) save() error {
	cacheFile := filepath.Join(c.dir, "cache.json")
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0644)
}

func (c *Cache) hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

func (c *Cache) GetFileCache(key string) (string, bool) {
	hashedKey := c.hashKey(key)
	filePath := filepath.Join(c.dir, hashedKey+".zip")

	if entry, exists := c.entries[hashedKey]; exists {
		if time.Since(entry.Timestamp) > entry.TTL {
			delete(c.entries, hashedKey)
			os.Remove(filePath)
			c.save()
			return "", false
		}

		if _, err := os.Stat(filePath); err == nil {
			return filePath, true
		}
	}

	return "", false
}

func (c *Cache) SetFileCache(key string, data []byte, ttl time.Duration) (string, error) {
	hashedKey := c.hashKey(key)
	filePath := filepath.Join(c.dir, hashedKey+".zip")

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	c.entries[hashedKey] = &CacheEntry{
		Key:       key,
		Data:      filePath,
		Timestamp: time.Now(),
		TTL:       ttl,
	}

	if err := c.save(); err != nil {
		return "", err
	}

	return filePath, nil
}
