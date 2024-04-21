package allocate

import (
	"fmt"
	"net"
)

func IsPortFree(port int) bool {
	address := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		// Port is not available
		return false
	}
	defer func(listener net.Listener) {
		err := listener.Close()
		if err != nil {
			fmt.Printf("listener cannot be closed: %v", err)
		}
	}(listener)
	// Port is available
	return true
}
