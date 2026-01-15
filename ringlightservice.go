package main

type RingLightService struct{}

func (r *RingLightService) SetColor(monitorName string, red, green, blue uint8) {
	manager.SetColor(monitorName, red, green, blue)
}

func (r *RingLightService) SetBrightness(monitorName string, alpha uint8) {
	manager.SetAlpha(monitorName, alpha)
}

func (r *RingLightService) SetEnabled(enabled bool) {
	manager.SetEnabled(enabled)
}

func (r *RingLightService) SetWidth(monitorName string, width int32) {
	manager.SetWidth(monitorName, width)
}

func (r *RingLightService) SetRadius(monitorName string, radius int32) {
	manager.SetRadius(monitorName, radius)
}



func (r *RingLightService) GetMonitors() []string {
	var monitors []string
	for name := range manager.windows {
		monitors = append(monitors, name)
	}
	return monitors
}

func (r *RingLightService) ToggleMonitor(name string, enabled bool) {
	manager.SetMonitorEnabled(name, enabled)
}

