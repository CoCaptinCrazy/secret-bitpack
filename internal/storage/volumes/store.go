package volumes

type Store struct{ volumes map[string]Volume }

func NewStore() *Store { return &Store{volumes: make(map[string]Volume)} }

func (s *Store) Put(volume Volume) { s.volumes[volume.Name] = volume }

func (s *Store) Get(name string) (Volume, bool) {
	volume, ok := s.volumes[name]
	return volume, ok
}
