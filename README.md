# SSID (Specific Sequence ID) Generator

# Documentation:

## Overview
The `SequenceGenerator` is a Go package designed to generate unique, sequential IDs with configurable prefixes and lengths. It uses PostgreSQL sequences for ID generation, ensuring thread safety and efficiency.

### Key Features
- Thread-safe sequence creation and management.
- Configurable prefixes and ID lengths.
- Safe handling of dynamic SQL queries.
- Minimal database overhead with caching of created sequences.

---

## Installation

### Prerequisites
- Go 1.18 or later.
- PostgreSQL database.
- GORM as the ORM library.
- `github.com/lib/pq` for PostgreSQL support.

### Install the Package
```bash
go get github.com/your-repo/ssid
```

---

## Usage

### Define Configuration
Define the prefix and sequence configuration for the IDs you want to generate:

```go
import (
	"github.com/your-repo/ssid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Initialize GORM DB connection
	dsn := "host=localhost user=postgres password=yourpassword dbname=yourdb port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Define sequence configurations
	config := ssid.PrefixConfig{
		"PL": {Length: 8, SeqName: "pl_sequence"},
		"REF": {Length: 10, SeqName: "ref_sequence"},
	}

	// Create a new SequenceGenerator instance
	generator, err := ssid.NewSequenceGenerator(db, config)
	if err != nil {
		panic(err)
	}

	// Generate an ID
	id, err := generator.GenerateID("PL")
	if err != nil {
		panic(err)
	}

	fmt.Println("Generated ID:", id.String())
}
```

---

## API Reference

### `NewSequenceGenerator`
```go
func NewSequenceGenerator(db *gorm.DB, config map[string]PrefixConfig) (*SequenceGenerator, error)
```
**Description**:
Creates a new instance of `SequenceGenerator` with the provided database connection and configuration.

**Parameters**:
- `db`: A GORM database connection.
- `config`: A map of prefixes to their respective configurations.

**Returns**:
- `*SequenceGenerator`: A pointer to the initialized generator.
- `error`: Error if initialization fails.

---

### `GenerateID`
```go
func (sg *SequenceGenerator) GenerateID(prefix string) (SSID, error)
```
**Description**:
Generates a new ID based on the given prefix.

**Parameters**:
- `prefix`: The prefix for the ID (e.g., "PL").

**Returns**:
- `SSID`: The generated ID.
- `error`: Error if generation fails.

---

### `SSID.String`
```go
func (s SSID) String() string
```
**Description**:
Returns the string representation of the generated ID.

---

## Examples

### Generate IDs for Multiple Prefixes
```go
id1, err := generator.GenerateID("PL")
if err != nil {
	log.Fatal(err)
}
fmt.Println("Generated ID 1:", id1.String())

id2, err := generator.GenerateID("REF")
if err != nil {
	log.Fatal(err)
}
fmt.Println("Generated ID 2:", id2.String())
```

---

## Best Practices
- Always validate the `PrefixConfig` during initialization to ensure correctness.
- Use `pq.QuoteIdentifier` to safely handle dynamic sequence names.
- Avoid hardcoding prefixes or sequence names; define them in a centralized configuration.

---

## Error Handling
Common errors include:
- Missing configuration for a prefix.
- Database connection issues.
- Failure to create or access sequences.

Handle errors gracefully to ensure smooth operation.

---

## License
This package is licensed under the MIT License. See LICENSE for details.

---

## Contributing
Contributions are welcome! Please open an issue or submit a pull request on GitHub.

