package auth

type Permissions map[string]bool

func (p Permissions) Allows(permission string) bool {
	return p["*"] || p[permission]
}
