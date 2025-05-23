package cache

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
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

// LoadCache charge le cache depuis le fichier ou crée un nouveau cache
func LoadCache() (*Cache, error) {
	var cacheDir string

	// Utiliser XDG_CONFIG_HOME si défini, sinon ~/.config
	if xdgConfig := os.Getenv("XDG_CONFIG_HOME"); xdgConfig != "" {
		cacheDir = filepath.Join(xdgConfig, "gh-ci", ".cache")
	} else {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, err
		}
		cacheDir = filepath.Join(homeDir, ".config", "gh-ci", ".cache")
	}

	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return nil, err
	}

	c := &Cache{
		entries: make(map[string]*CacheEntry),
		dir:     cacheDir,
	}

	// Charger le cache existant
	cacheFile := filepath.Join(cacheDir, "cache.json")
	if data, err := os.ReadFile(cacheFile); err == nil {
		json.Unmarshal(data, &c.entries)
	}

	return c, nil
}

// Get récupère une entrée du cache
func (c *Cache) Get(key string) (any, bool) {
	hashedKey := c.hashKey(key)
	entry, exists := c.entries[hashedKey]

	if !exists {
		return nil, false
	}

	// Vérifier si l'entrée a expiré
	if time.Since(entry.Timestamp) > entry.TTL {
		delete(c.entries, hashedKey)
		c.save()
		return nil, false
	}

	return entry.Data, true
}

// Set ajoute une entrée au cache
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

// Delete supprime une entrée du cache
func (c *Cache) Delete(key string) {
	hashedKey := c.hashKey(key)
	delete(c.entries, hashedKey)
	c.save()
}

// Clear vide tout le cache
func (c *Cache) Clear() error {
	c.entries = make(map[string]*CacheEntry)
	return c.save()
}

// save sauvegarde le cache sur disque
func (c *Cache) save() error {
	cacheFile := filepath.Join(c.dir, "cache.json")
	data, err := json.MarshalIndent(c.entries, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(cacheFile, data, 0644)
}

// hashKey crée un hash SHA256 de la clé pour éviter les problèmes de caractères spéciaux
func (c *Cache) hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// GetFileCache récupère le chemin d'un fichier en cache
func (c *Cache) GetFileCache(key string) (string, bool) {
	hashedKey := c.hashKey(key)
	filePath := filepath.Join(c.dir, hashedKey+".zip")

	if entry, exists := c.entries[hashedKey]; exists {
		// Vérifier si l'entrée a expiré
		if time.Since(entry.Timestamp) > entry.TTL {
			delete(c.entries, hashedKey)
			os.Remove(filePath)
			c.save()
			return "", false
		}

		// Vérifier si le fichier existe toujours
		if _, err := os.Stat(filePath); err == nil {
			return filePath, true
		}
	}

	return "", false
}

// SetFileCache sauvegarde un fichier en cache
func (c *Cache) SetFileCache(key string, data []byte, ttl time.Duration) (string, error) {
	hashedKey := c.hashKey(key)
	filePath := filepath.Join(c.dir, hashedKey+".zip")

	// Sauvegarder le fichier
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", err
	}

	// Ajouter l'entrée au cache
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
