package cmd

import (
	"errors"
	"flag"
	"fmt"
	"log"

	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"

	"github.com/asssaf/stemma-soil-go/soil"
)

type ReadCommand struct {
	fs   *flag.FlagSet
	addr int
}

func NewReadCommand() *ReadCommand {
	c := &ReadCommand{
		fs: flag.NewFlagSet("read", flag.ExitOnError),
	}

	c.fs.IntVar(&c.addr, "address", 0, "Device address (0x36-0x39)")

	// c.fs.Usage = func() {
	// 	fmt.Fprintf(flag.CommandLine.Output(), "Usage: automationhat %s <1|2|3> <on|off>\n", c.fs.Name())
	// }

	return c
}

func (c *ReadCommand) Name() string {
	return c.fs.Name()
}

func (c *ReadCommand) Init(args []string) error {
	if err := c.fs.Parse(args); err != nil {
		return err
	}

	flag.Usage = c.fs.Usage

	return nil
}

func (c *ReadCommand) Execute() error {
	// Make sure periph is initialized.
	if _, err := host.Init(); err != nil {
		log.Fatal(err)
	}

	i2cPort, err := i2creg.Open("/dev/i2c-1")
	if err != nil {
		log.Fatal(err)
	}

	opts := soil.DefaultOpts
	if c.addr != 0 {
		if c.addr < 0x36 || c.addr > 0x39 {
			return errors.New(fmt.Sprintf("given address not supported by device: %x", c.addr))
		}
		opts.Addr = uint16(c.addr)
	}

	dev, err := soil.NewI2C(i2cPort, &opts)
	if err != nil {
		log.Fatal(err)
	}
	defer dev.Halt()

	values := soil.SensorValues{}
	if err := dev.Sense(&values); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Capacitance: %d\n", values.Capacitance)
	fmt.Printf("Temperature: %s\n", values.Temperature)

	return nil
}
