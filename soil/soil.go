package soil

import (
	"encoding/binary"
	"errors"
	"fmt"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/physic"
)

const (
	touchBase          byte = 0x0f
	touchChannelOffset byte = 0x10
	statusBase         byte = 0x00
	statusTemp         byte = 0x04
)

type Opts struct {
	Addr uint16
}

var DefaultOpts = Opts{
	Addr: 0x36,
}

type Dev struct {
	c    conn.Conn
	opts Opts
}

type SensorValues struct {
	Capacitance uint16
	Temperature physic.Temperature
}

func NewI2C(b i2c.Bus, opts *Opts) (*Dev, error) {
	switch opts.Addr {
	case 0x36, 0x37, 0x38, 0x39:
	default:
		return nil, errors.New("soil: given address not supported by device")
	}
	dev := &Dev{
		c:    &i2c.Dev{Bus: b, Addr: opts.Addr},
		opts: *opts,
	}
	return dev, nil
}

func (d *Dev) Sense(values *SensorValues) error {
	cap, err := d.senseCapacitance()
	if err != nil {
		return err
	}

	temp, err := d.senseTemperature()
	if err != nil {
		return err
	}

	values.Capacitance = cap
	values.Temperature = temp

	return nil
}

func (d *Dev) senseCapacitance() (uint16, error) {
	// r/w Tx doesn't work well, need to wait 5 milliseconds between write and read
	if err := d.c.Tx([]byte{touchBase, touchChannelOffset}, []byte{}); err != nil {
		return 0, err
	}

	time.Sleep(5 * time.Millisecond)

	read := make([]byte, 2)
	if err := d.c.Tx([]byte{}, read); err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Millisecond)

	cap := binary.BigEndian.Uint16(read)
	if cap > 4095 {
		return 0, errors.New(fmt.Sprintf("bad sample: %d", cap))
	}

	return cap, nil
}

func (d *Dev) senseTemperature() (physic.Temperature, error) {
	// r/w Tx doesn't work well, need to wait 5 milliseconds between write and read
	if err := d.c.Tx([]byte{statusBase, statusTemp}, []byte{}); err != nil {
		return 0, err
	}

	time.Sleep(5 * time.Millisecond)

	read := make([]byte, 4)
	if err := d.c.Tx([]byte{}, read); err != nil {
		return 0, err
	}

	time.Sleep(1 * time.Millisecond)

	read[0] = read[0] & 0x3F
	tempRaw := binary.BigEndian.Uint32(read)
	temp := physic.ZeroCelsius + physic.Temperature(0.00001525878*float64(tempRaw))*physic.Celsius

	return temp, nil
}

func (d *Dev) Halt() error {
	return nil
}
