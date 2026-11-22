package volumes

type Volume struct {
	Name      string
	Driver    string
	MountPath string
	ReadOnly  bool
}
