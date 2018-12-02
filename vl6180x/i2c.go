package vl6180x

// Package i2c provides low level control over the linux i2c bus.
//
// Before usage you should load the i2c-dev kernel module
//
//      sudo modprobe i2c-dev
//
// Each i2c bus can address 127 independent i2c devices, and most
// linux systems contain several buses.

import (
	"encoding/hex"
	"fmt"
	"os"
	"syscall"

	log "github.com/sirupsen/logrus"
)

const (
	I2C_SLAVE = 0x0703
)

// I2C represents a connection to I2C-device.
type I2C struct {
	addr uint8
	bus  int
	rc   *os.File
}

// NewI2C opens a connection for I2C-device.
// SMBus (System Management Bus) protocol over I2C
// supported as well: you should preliminary specify
// register address to read from, either write register
// together with the data in case of write operations.
func NewI2C(addr uint8, bus int) (*I2C, error) {
	f, err := os.OpenFile(fmt.Sprintf("/dev/i2c-%d", bus), os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}
	if err := ioctl(f.Fd(), I2C_SLAVE, uintptr(addr)); err != nil {
		return nil, err
	}
	v := &I2C{rc: f, bus: bus, addr: addr}
	return v, nil
}

// GetBus return bus line, where I2C-device is allocated.
func (v *I2C) GetBus() int {
	return v.bus
}

// GetBus return device occupied address in the bus.
func (v *I2C) GetAddr() uint8 {
	return v.addr
}

func (v *I2C) write(buf []byte) (int, error) {
	return v.rc.Write(buf)
}

// Write sends bytes to the remote I2C-device. The interpretation of
// the message is implementation-dependant.
func (v *I2C) WriteBytes(buf []byte) (int, error) {
	log.Debugf("Write %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return v.write(buf)
}

func (v *I2C) read(buf []byte) (int, error) {
	return v.rc.Read(buf)
}

// ReadBytes reads bytes from I2C-device.
// Number of bytes read correspond to buf parameter length.
func (v *I2C) ReadBytes(buf []byte) (int, error) {
	n, err := v.read(buf)
	if err != nil {
		return n, err
	}
	log.Debugf("Read %d hex bytes: [%+v]", len(buf), hex.EncodeToString(buf))
	return n, nil
}

// Close I2C-connection.
func (v *I2C) Close() error {
	return v.rc.Close()
}

func (v *I2C) uint16LE(reg uint16) []byte {
	regBytes := make([]byte, 2)
	regBytes[0] = byte((reg >> 8) & 0x00FF)
	regBytes[1] = byte(reg & 0x00FF)

	return regBytes
}

// ReadRegBytes read count of n byte's sequence from I2C-device
// starting from reg address.
// SMBus (System Management Bus) protocol over I2C.
func (v *I2C) ReadRegBytes(reg uint16, n int) ([]byte, int, error) {
	log.Debugf("Read %d bytes starting from reg 0x%0d...", n, reg)

	_, err := v.WriteBytes(v.uint16LE(reg))
	if err != nil {
		return nil, 0, err
	}
	buf := make([]byte, n)
	c, err := v.ReadBytes(buf)
	if err != nil {
		return nil, 0, err
	}
	return buf, c, nil

}

func (v *I2C) WriteRegU8(reg uint16, value uint8) error {
	buf := v.uint16LE(reg)
	_, err := v.WriteBytes(append(buf, value))
	if err != nil {
		return err
	}
	log.Debugf("Write U8 %d to reg 0x%0X", value, reg)
	return nil
}

func (v *I2C) WriteRegU16(reg uint16, value uint16) error {
	buf := v.uint16LE(reg)
	_, err := v.WriteBytes(append(buf, v.uint16LE(value)...))
	if err != nil {
		return err
	}
	log.Debugf("Write U16 %d to reg 0x%0X", value, reg)
	return nil
}

func ioctl(fd, cmd, arg uintptr) error {
	_, _, err := syscall.Syscall6(syscall.SYS_IOCTL, fd, cmd, arg, 0, 0, 0)
	if err != 0 {
		return err
	}
	return nil
}
