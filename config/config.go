package config

type PrintFrequency int

const (
	PRINT_MONTHLY = iota
	PRINT_YEARLY
	PRINT_MILESTONE
)

// config represents initial command line configuration
type Config struct {
	jsonPath       string
	printFrequency PrintFrequency
}
