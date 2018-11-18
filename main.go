package main

import "log"
import i2c "github.com/d2r2/go-i2c"

func main() {
	// Create new connection to I2C bus on 2 line with address 0x27
	i2c, err := i2c.NewI2C(0x29, 1)
	if err != nil {
		log.Fatal(err)
	}
	// Free I2C connection on exit
	defer i2c.Close()
	// Here goes code specific for sending and reading data
	// to and from device connected via I2C bus, like:
	_, err = i2c.WriteBytes([]byte{0x1, 0xF3})
	if err != nil {
		log.Fatal(err)
	}
}
