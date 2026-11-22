package pull

type Manager struct{ Workers []Worker }

func (m *Manager) Add(worker Worker) { m.Workers = append(m.Workers, worker) }
