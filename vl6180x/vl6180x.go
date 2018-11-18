package vl6180x

import (
	"bytes"
	"encoding/binary"
	"time"

	i2c "github.com/d2r2/go-i2c"

	log "github.com/sirupsen/logrus"
)

const (
	// The fixed I2C addres
	VL6180X_DEFAULT_I2C_ADDR = 0x29

	///! Device model identification number
	VL6180X_REG_IDENTIFICATION_MODEL_ID = 0x000
	///! Interrupt configuration
	VL6180X_REG_SYSTEM_INTERRUPT_CONFIG = 0x014
	///! Interrupt clear bits
	VL6180X_REG_SYSTEM_INTERRUPT_CLEAR = 0x015
	///! Fresh out of reset bit
	VL6180X_REG_SYSTEM_FRESH_OUT_OF_RESET = 0x016
	///! Trigger Ranging
	VL6180X_REG_SYSRANGE_START = 0x018
	///! Trigger Lux Reading
	VL6180X_REG_SYSALS_START = 0x038
	///! Lux reading gain
	VL6180X_REG_SYSALS_ANALOGUE_GAIN = 0x03F
	///! Integration period for ALS mode, high byte
	VL6180X_REG_SYSALS_INTEGRATION_PERIOD_HI = 0x040
	///! Integration period for ALS mode, low byte
	VL6180X_REG_SYSALS_INTEGRATION_PERIOD_LO = 0x041
	///! Specific error codes
	VL6180X_REG_RESULT_RANGE_STATUS = 0x04d
	///! Interrupt status
	VL6180X_REG_RESULT_INTERRUPT_STATUS_GPIO = 0x04f
	///! Light reading value
	VL6180X_REG_RESULT_ALS_VAL = 0x050
	///! Ranging reading value
	VL6180X_REG_RESULT_RANGE_VAL = 0x062

	VL6180X_ALS_GAIN_1    = 0x06 ///< 1x gain
	VL6180X_ALS_GAIN_1_25 = 0x05 ///< 1.25x gain
	VL6180X_ALS_GAIN_1_67 = 0x04 ///< 1.67x gain
	VL6180X_ALS_GAIN_2_5  = 0x03 ///< 2.5x gain
	VL6180X_ALS_GAIN_5    = 0x02 ///< 5x gain
	VL6180X_ALS_GAIN_10   = 0x01 ///< 1=0x gain
	VL6180X_ALS_GAIN_20   = 0x00 ///< 2=0x gain
	VL6180X_ALS_GAIN_40   = 0x07 ///< 4=0x gain

	VL6180X_ERROR_NONE        = 0  ///< Success!
	VL6180X_ERROR_SYSERR_1    = 1  ///< System error
	VL6180X_ERROR_SYSERR_5    = 5  ///< Sysem error
	VL6180X_ERROR_ECEFAIL     = 6  ///< Early convergence estimate fail
	VL6180X_ERROR_NOCONVERGE  = 7  ///< No target detected
	VL6180X_ERROR_RANGEIGNORE = 8  ///< Ignore threshold check failed
	VL6180X_ERROR_SNR         = 11 ///< Ambient conditions too high
	VL6180X_ERROR_RAWUFLOW    = 12 ///< Raw range algo underflow
	VL6180X_ERROR_RAWOFLOW    = 13 ///< Raw range algo overflow
	VL6180X_ERROR_RANGEUFLOW  = 14 ///< Raw range algo underflow
	VL6180X_ERROR_RANGEOFLOW  = 15 ///< Raw range algo overflow
)

type Vl6180x struct {
	i2cAddr   uint8
	i2cDevice *i2c.I2C
	ioTimeout time.Duration
}

func NewVl6180x(i2cAddr uint8) *Vl6180x {
	v := &Vl6180x{}
	v.i2cAddr = i2cAddr

	i2cDevice, err := i2c.NewI2C(0x29, 0)
	if err != nil {
		log.Fatal(err)
	}
	v.i2cDevice = i2cDevice

	// Check model id
	modelId, err := v.i2cDevice.ReadRegU8(VL6180X_REG_IDENTIFICATION_MODEL_ID)
	if err != nil {
		log.Fatal(err)
	}
	if modelId != 0xB4 {
		log.Error("The model id is %x not 0xB4", modelId)
		return nil
	}

	return v
}

func (v *Vl6180x) Config(i2c *i2c.I2C) {

}

func (v *Vl6180x) WriteByte(reg uint16, value byte) error {
	return v.Write(reg, []byte{value})
}

func (v *Vl6180x) Write(reg uint16, value []byte) error {
	bytesBuffer := new(bytes.Buffer)
	binary.Write(bytesBuffer, binary.LittleEndian, reg)
	binary.Write(bytesBuffer, binary.LittleEndian, value)

	_, err := v.i2cDevice.WriteBytes(bytesBuffer.Bytes())
	if err != nil {
		return err
	}
	log.Debugf("Write %x to 0x%0X", value, reg)

	return nil
}

func (v *Vl6180x) loadSettings() {
	log.Info("Loading vl6180x settings")

	//v.WriteByte(0x0207, 0x01)
	// private settings from page 24 of app note
	v.WriteByte(0x0207, 0x01)
	v.WriteByte(0x0208, 0x01)
	v.WriteByte(0x0096, 0x00)
	v.WriteByte(0x0097, 0xfd)
	v.WriteByte(0x00e3, 0x00)
	v.WriteByte(0x00e4, 0x04)
	v.WriteByte(0x00e5, 0x02)
	v.WriteByte(0x00e6, 0x01)
	v.WriteByte(0x00e7, 0x03)
	v.WriteByte(0x00f5, 0x02)
	v.WriteByte(0x00d9, 0x05)
	v.WriteByte(0x00db, 0xce)
	v.WriteByte(0x00dc, 0x03)
	v.WriteByte(0x00dd, 0xf8)
	v.WriteByte(0x009f, 0x00)
	v.WriteByte(0x00a3, 0x3c)
	v.WriteByte(0x00b7, 0x00)
	v.WriteByte(0x00bb, 0x3c)
	v.WriteByte(0x00b2, 0x09)
	v.WriteByte(0x00ca, 0x09)
	v.WriteByte(0x0198, 0x01)
	v.WriteByte(0x01b0, 0x17)
	v.WriteByte(0x01ad, 0x00)
	v.WriteByte(0x00ff, 0x05)
	v.WriteByte(0x0100, 0x05)
	v.WriteByte(0x0199, 0x05)
	v.WriteByte(0x01a6, 0x1b)
	v.WriteByte(0x01ac, 0x3e)
	v.WriteByte(0x01a7, 0x1f)
	v.WriteByte(0x0030, 0x00)

	// Recommended : Public registers - See data sheet for more detail
	v.WriteByte(0x0011, 0x10) // Enables polling for 'New Sample ready'
	// when measurement completes
	v.WriteByte(0x010a, 0x30) // Set the averaging sample period
	// (compromise between lower noise and
	// increased execution time)
	v.WriteByte(0x003f, 0x46) // Sets the light and dark gain (upper
	// nibble). Dark gain should not be
	// changed.
	v.WriteByte(0x0031, 0xFF) // sets the # of range measurements after
	// which auto calibration of system is
	// performed
	v.WriteByte(0x0040, 0x63) // Set ALS integration time to 100ms
	v.WriteByte(0x002e, 0x01) // perform a single temperature calibration
	// of the ranging sensor

	// Optional: Public registers - See data sheet for more detail
	v.WriteByte(0x001b, 0x09) // Set default ranging inter-measurement
	// period to 100ms
	v.WriteByte(0x003e, 0x31) // Set default ALS inter-measurement period
	// to 500ms
	v.WriteByte(0x0014, 0x24) // Configures interrupt on 'New Sample
	// Ready threshold event'
}

func (v *Vl6180x) ReadByte(reg uint8) byte {
	data, err := v.i2cDevice.ReadRegU8(reg)
	if err != nil {
		log.Error("Can't read byte from %x", reg)
		log.Fatal(err)
		return 0x00
	}

	return data
}

func (v *Vl6180x) ReadBytes(reg uint8, n int) []byte {
	buf, c, err := v.i2cDevice.ReadRegBytes(reg, n)
	if err != nil {
		log.Error("Can't read reg %x", reg)
		return []byte{}
	}
}

func (v *Vl6180x) ReadRangeStatus() uint8 {
	return (v.ReadByte(VL6180X_REG_RESULT_RANGE_STATUS) >> 4)
}

func (v *Vl6180x) ReadLux(gain uint8) {
	var reg byte = 0

	reg = v.ReadByte(VL6180X_REG_SYSTEM_INTERRUPT_CONFIG)
	reg = reg & ^byte(0x38)
	reg = reg | (0x04 << 3) // IRQ on ALS ready

	v.WriteByte(VL6180X_REG_SYSTEM_INTERRUPT_CONFIG, reg)

	// 100 ms integration period
	v.WriteByte(VL6180X_REG_SYSALS_INTEGRATION_PERIOD_HI, 0)
	v.WriteByte(VL6180X_REG_SYSALS_INTEGRATION_PERIOD_LO, 100)

	// analog gain
	if gain > VL6180X_ALS_GAIN_40 {
		gain = VL6180X_ALS_GAIN_40
	}
	v.WriteByte(VL6180X_REG_SYSALS_ANALOGUE_GAIN, 0x40|gain)

	// start ALS
	v.WriteByte(VL6180X_REG_SYSALS_START, 0x01)

	// Poll until "New Sample Ready threshold event" is set
	for ((v.ReadByte(VL6180X_REG_RESULT_INTERRUPT_STATUS_GPIO) >> 3) & 0x07) != 4 {
		log.Debug("Poll until 'New Sample Ready threshold event' is set")
	}

}

func (v *Vl6180x) ReadRange() uint8 {
	for (v.ReadByte(VL6180X_REG_RESULT_RANGE_STATUS) & 0x01) == 0 {
		log.Debug("Device is not ready to read the range")
	}

	v.WriteByte(VL6180X_REG_SYSRANGE_START, 0x01)

	// Poll until bit 2 is set
	for (v.ReadByte(VL6180X_REG_RESULT_INTERRUPT_STATUS_GPIO) & 0x04) == 0 {
		log.Debug("Poll until the interrupt status is 0x04")
	}

	// Read range in mm
	result := v.ReadByte(VL6180X_REG_RESULT_RANGE_VAL)

	// Clear interrupt
	v.WriteByte(VL6180X_REG_SYSTEM_INTERRUPT_CLEAR, 0x07)

	return result
}

func (v *Vl6180x) Close() {
	v.i2cDevice.Close()
}

func (v *Vl6180x) setTimeout(timeout time.Duration) {
	v.ioTimeout = timeout
}

func (v *Vl6180x) Init(i2c *i2c.I2C) error {
	v.setTimeout(time.Millisecond * 1000)

	return nil
}
