package health

import "fmt"

type StatusError struct{ Code int }

func (e *StatusError) Error() string { return fmt.Sprintf("health check returned HTTP %d", e.Code) }
