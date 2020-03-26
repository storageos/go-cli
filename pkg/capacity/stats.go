package capacity

// Stats struct for CapacityStats
type Stats struct {
	// Total bytes in the filesystem
	Total uint64 `json:"total,omitempty"`
	// Free bytes in the filesystem available to root user
	Free uint64 `json:"free,omitempty"`
	// Byte value available to an unprivileged user
	Available uint64 `json:"available,omitempty"`
}
