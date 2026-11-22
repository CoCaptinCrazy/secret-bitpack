package grpc

import "fmt"

type Service struct{ Registry *Registry }

func (s Service) Invoke(name string, request []byte) ([]byte, error) {
	handler, ok := s.Registry.Lookup(name)
	if !ok {
		return nil, fmt.Errorf("unknown service %q", name)
	}
	return handler(request)
}
