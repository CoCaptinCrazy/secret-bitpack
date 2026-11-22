package sandbox

type Isolation interface {
	Apply(Profile) error
}
