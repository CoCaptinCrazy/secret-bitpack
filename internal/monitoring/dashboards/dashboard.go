package dashboards

type Dashboard struct {
	Name    string
	Widgets []Widget
}

func (d *Dashboard) Add(widget Widget) { d.Widgets = append(d.Widgets, widget) }
