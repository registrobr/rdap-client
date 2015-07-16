package protocol

type Conformance struct {
	Levels []string `json:"rdapConformance,omitempty"`
}

type ConformanceSetter interface {
	SetConformance([]string)
}

func (l *Conformance) SetConformance(levels []string) {
	l.Levels = levels
}
