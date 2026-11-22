package services

type Catalog struct{ services map[string]Service }

func NewCatalog() *Catalog { return &Catalog{services: make(map[string]Service)} }

func (c *Catalog) Register(service Service) { c.services[service.Name] = service }

func (c *Catalog) Resolve(name string) (Service, bool) {
	service, ok := c.services[name]
	return service, ok
}
