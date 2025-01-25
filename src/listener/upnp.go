package listener

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/huin/goupnp/dcps/internetgateway2"
)

type RouterClient interface {
	AddPortMapping(
		NewRemoteHost string,
		NewExternalPort uint16,
		NewProtocol string,
		NewInternalPort uint16,
		NewInternalClient string,
		NewEnabled bool,
		NewPortMappingDescription string,
		NewLeaseDuration uint32,
	) (err error)

	GetExternalIPAddress() (
		NewExternalIPAddress string,
		err error,
	)
}

func getClient[T RouterClient](clientChan chan RouterClient, errChan chan error, f func() (clients []T, errors []error, err error)) {
	clients, _, err := f()
	errChan <- err
	if len(clients) == 1 {
		clientChan <- clients[0]
		return
	}
	clientChan <- nil
}

func PickRouterClient() (RouterClient, error) {
	clientChan := make(chan RouterClient, 3)
	errChan := make(chan error, 3)

	go getClient(clientChan, errChan, internetgateway2.NewWANIPConnection2Clients)
	go getClient(clientChan, errChan, internetgateway2.NewWANIPConnection1Clients)
	go getClient(clientChan, errChan, internetgateway2.NewWANPPPConnection1Clients)

	for i := 0; i < 3; i++ {
		err := <-errChan
		if err != nil {
			return nil, err
		}
	}
	for i := 0; i < 3; i++ {
		client := <-clientChan
		if client != nil {
			return client, nil
		}
	}
	return nil, fmt.Errorf("no suitable UPnP clients found")
}

func ForwardNAT(listenPort uint16) error {
	client, err := PickRouterClient()
	if err != nil {
		return err
	}

	externalIP, err := client.GetExternalIPAddress()
	if err != nil {
		return err
	}
	lanIP, err := getLANIP()
	if err != nil {
		return err
	}
	log.Print("Forwarding ", lanIP, ":", listenPort, " to ", externalIP, ":", listenPort, ".\n")

	return client.AddPortMapping(
		"",
		// External port number to expose to Internet:
		listenPort,
		// Forward TCP (this could be "UDP" if we wanted that instead).
		"UDP",
		// Internal port number on the LAN to forward to.
		// Some routers might not support this being different to the external
		// port number.
		listenPort,
		// Internal address on the LAN we want to forward to.
		lanIP,
		// Enabled:
		true,
		// Informational description for the client requesting the port forwarding.
		"Akademi",
		// How long should the port forward last for in seconds.
		// If you want to keep it open for longer and potentially across router
		// resets, you might want to periodically request before this elapses.
		3600,
	)
}

// Get the local machines's IP address.
func getLANIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("lan ip not found")
}

// The UPnPWorker goroutine is responsible for making sure
// the listen port of the application is accessible to the
// outside world.
func UPnPWorker(errChan chan error, listenPort uint16) {
	for {
		err := ForwardNAT(listenPort)
		if err != nil {
			log.Print(err)
			if err.Error() == "no suitable UPnP clients found" {
				log.Print("Proceeding without UPnP.")
				return
			}
		}
		time.Sleep(time.Minute * 30)
	}
}
