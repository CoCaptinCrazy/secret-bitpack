package execution

type Environment map[string]string

func (e Environment) Clone() Environment {
	clone := make(Environment, len(e))
	for key, value := range e {
		clone[key] = value
	}
	return clone
}
