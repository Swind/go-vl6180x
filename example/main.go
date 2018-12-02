package main

import (
	"fmt"

	log "github.com/sirupsen/logrus"
	vl6180x "github.com/swind/go-vl6180x/vl6180x"
)

func main() {
	device, err := vl6180x.NewVl6180x(0x29, 1)
	if err != nil {
		log.Fatal(err)
	}
	device.LoadSettings()
	fmt.Println("scaling 1:", device.ReadRange())

	device.SetScaling(2)
	fmt.Println("scaling 2:", device.ReadRange())

	device.SetScaling(3)
	fmt.Println("scaling 3:", device.ReadRange())
}
