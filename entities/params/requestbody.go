package params

// RequestBody represents type for request body structures to simplify validation
type RequestBody interface {
	TrimSpaces()
}
