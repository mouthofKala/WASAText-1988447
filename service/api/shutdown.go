package api

// Close should close everything opened in the lifecycle of the `_router`; for example, background goroutines.
func (rt *_router) Close() error {
	return nil
}

// do NOT use for app lifecycle resources, that is done in main
