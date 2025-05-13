package main

import (
	"github.com/kavindu-udara/fileshear.git/fileserver"
	"github.com/kavindu-udara/fileshear.git/internal"
)

func main(){

	// check LAN connection
	if !internal.CheckNetwork() {
		println("No LAN connected!")
		return
	}

	// get the ip addres
	ip, err := internal.GetIp()
	if(err != err){
		println("Error when try to get the ip")
	}
	println("Local ip address : " + ip)

	// start the API
	fileserver.API(ip)
}
