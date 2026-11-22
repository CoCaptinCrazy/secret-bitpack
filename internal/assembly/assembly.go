package assembly

import "errors"

type Plan struct {
	Name  string
	Steps []string
}

func (p Plan) Validate() error {
	if p.Name == "" {
		return errors.New("assembly name is required")
	}
	return nil
}
