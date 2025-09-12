package domain

type ServiceDefinition struct {
	ServiceName  string
	Image        string
	PortForwards map[uint16]uint16 // map[hostport]containerport
	Environment  map[string]string
}
