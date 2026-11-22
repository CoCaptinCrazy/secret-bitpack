package scanning

type Scanner interface {
	Scan([]byte) (Report, error)
}
