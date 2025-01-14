package ssid

import (
	"fmt"
	"sync"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

type SSID struct {
	ID string
}

func (s SSID) String() string {
	return s.ID
}

type SequenceGenerator struct {
	DB               *gorm.DB
	Config           PrefixConfig    // Map of prefix and corresponding config (length, etc.)
	createdSequences map[string]bool // Cache for created sequences
	mu               sync.Mutex      // Mutex for thread-safe operations
}

type PrefixConfig map[string]struct {
	Length  int
	SeqName string
}

func NewSequenceGenerator(db *gorm.DB, config PrefixConfig) (*SequenceGenerator, error) {
	return &SequenceGenerator{
		DB:               db,
		Config:           config,
		createdSequences: make(map[string]bool),
	}, nil
}

func (sg *SequenceGenerator) GenerateID(prefix string) (SSID, error) {

	// Retrieve the configuration for the given prefix
	config, ok := sg.Config[prefix]
	{
		if !ok {
			return SSID{}, fmt.Errorf("no configuration found for prefix %s", prefix)
		}

		// Ensure the sequence exists, creating it only once
		if err := sg.ensureSequenceExists(config.SeqName); err != nil {
			return SSID{}, err
		}
	}

	// Get the next sequence value
	var nextVal int64
	{
		err := sg.DB.Raw("SELECT nextval($1)", config.SeqName).Scan(&nextVal).Error
		if err != nil {
			return SSID{}, fmt.Errorf("failed to generate next ID: %w", err)
		}
	}

	// Pad the sequence value with leading zeros and format the ID
	id := fmt.Sprintf("%s%0*d", prefix, config.Length, nextVal)

	return SSID{ID: id}, nil
}

func (sg *SequenceGenerator) ensureSequenceExists(seqName string) error {

	// Use a mutex to ensure thread-safe access
	sg.mu.Lock()
	defer sg.mu.Unlock()

	// Check if the sequence has already been created
	if sg.createdSequences[seqName] {
		return nil
	}

	// Create the sequence if it doesn't exist
	query := "CREATE SEQUENCE IF NOT EXISTS " + pq.QuoteIdentifier(seqName) + " INCREMENT BY 1 START WITH 1 MINVALUE 1 NO CYCLE"
	{
		if err := sg.DB.Exec(query).Error; err != nil {
			return fmt.Errorf("failed to create sequence %s: %w", seqName, err)
		}
	}

	// Mark the sequence as created
	sg.createdSequences[seqName] = true

	return nil
}
