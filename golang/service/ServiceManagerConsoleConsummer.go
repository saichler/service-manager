package service

func (sm *ServiceManager) Prompt() string {
	return sm.node.NetworkID().String()
}

func (sm *ServiceManager) InputReceived(line string) string {
	if line == "shutdown" {
		sm.shutdownMtx.Broadcast()
	}
	return ""
}

func (sm *ServiceManager) SupportedCommands() map[string]string {
	m := make(map[string]string)
	m["shutdown"] = "Shutdown the service manager."
	return m
}
