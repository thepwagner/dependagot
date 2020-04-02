package sandbox

type Sandboxen map[string]*Sandbox

func (s Sandboxen) NewSandbox() (string, *Sandbox) {
	sandboxID := "test"
	box := &Sandbox{}
	s[sandboxID] = box
	return sandboxID, box
}
