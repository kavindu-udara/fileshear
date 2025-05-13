package internal

import (
	"errors"
	"net"
)

// check is user connected to the LAN network
func CheckNetwork () bool {
	// get all network interfaces 
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _,iface := range interfaces {
		// skip loopback and interfaces that are down
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}

		// get interface addresses
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
            if ipnet, ok := addr.(*net.IPNet); ok {
                // Skip loopback and invalid addresses
                if ipnet.IP.IsLoopback() || ipnet.IP.To4() == nil {
                    continue
                }
                // If we find a valid IPv4 address, consider it connected
                return true
            }
        }
	}
	return false
}

// check local ip address and return it
func GetIp() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		// check if the address is an IP network
		if ipnet, ok := addr.(*net.IPNet); ok {
			// skip loopback addresses
			if ipnet.IP.IsLoopback() {
				continue
			}
			// get IPv4 addresses
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", errors.New("No IP address found !")
}
