package memory

type NoCacheError struct {
	error
}

func (n NoCacheError) Error() string {
	return "Cache not found"
}
