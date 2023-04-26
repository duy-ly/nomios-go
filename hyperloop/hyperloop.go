package hyperloop

type Hyperloop struct {
	running bool
}

func (h *Hyperloop) Start() {

}

func (h *Hyperloop) Stop() {

}

func (h *Hyperloop) IsRunning() bool {
	return h.running
}

func (h *Hyperloop) OnError() {

}
