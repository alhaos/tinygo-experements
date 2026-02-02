package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/adxl345"
	"tinygo.org/x/drivers/ssd1306"
)

const varbose = false

type Position struct {
	x int32
	y int32
}

func main() {

	time.Sleep(time.Second * 2)

	info("start delay success")

	i2c, err := initI2c(machine.GP4, machine.GP5)
	if err != nil {
		deadLoopPrint(err)
	}

	display, err := initOLED(i2c)
	if err != nil {
		deadLoopPrint(err)
	}

	accel := adxl345.New(i2c)

	accel.Configure()

	currentPosition := Position{
		x: 128 / 2,
		y: 64 / 2,
	}

	var nextPosition Position

	light := color.RGBA{255, 255, 255, 255}
	dark := color.RGBA{0, 0, 0, 255}

	for {

		x, y, _, err := accel.ReadAcceleration()
		if err != nil {
			fmt.Printf("Read acceleration error: %s", err.Error())
		}

		if x > 100 {
			nextPosition.x = currentPosition.x + 1
		}

		if x < -100 {
			nextPosition.x = currentPosition.x - 1
		}

		if y > 100 {
			nextPosition.y = currentPosition.y - 1
		}

		if y < -100 {
			nextPosition.y = currentPosition.y + 1
		}

		display.SetPixel(
			int16(currentPosition.x),
			int16(currentPosition.y),
			dark,
		)

		display.SetPixel(
			int16(nextPosition.x),
			int16(nextPosition.y),
			light,
		)

		currentPosition = nextPosition

		display.Display()

		time.Sleep(time.Millisecond * 50)
	}
}

func DrawChar(display *ssd1306.Device, x, y int16, ch rune) {
	data, ok := font5x8[ch]
	if !ok {
		data = font5x8['?']
	}

	var light = color.RGBA{255, 255, 255, 255}
	var dark = color.RGBA{0, 0, 0, 255}

	// fmt.Printf("found char [%c]\n", ch)

	for colNum, col := range data {
		xd := x + int16(colNum)
		info(fmt.Sprintf("%d: %08b %d\n", colNum, col, xd))
		for rowNum := range 8 {
			if (col>>rowNum)&1 == 1 {
				display.SetPixel(xd, y+int16(rowNum), light)
			} else {
				display.SetPixel(xd, y+int16(rowNum), dark)
			}
		}
	}
}

func DrawText(display *ssd1306.Device, x, y int16, text string) {
	for i, r := range []rune(text) {
		info(fmt.Sprintf("rune: [%c], i: %d\n", r, i))
		DrawChar(display, x+int16(i*6), y, r)
	}
}

func initI2c(sdaPin machine.Pin, sclPin machine.Pin) (*machine.I2C, error) {

	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GP4,
		SCL:       machine.GP5,
	})

	if err != nil {
		return nil, err
	}

	return i2c, nil
}

func initOLED(i2c *machine.I2C) (*ssd1306.Device, error) {

	display := ssd1306.NewI2C(i2c)

	display.Configure(ssd1306.Config{
		Width:    128,
		Height:   64,
		Address:  0x3C, // Ваш адрес
		VccState: ssd1306.SWITCHCAPVCC,
	})

	display.ClearDisplay()

	return display, nil
}

func deadLoopPrint(err error) {
	for {
		fmt.Println(err.Error())
		time.Sleep(time.Second)
	}
}

func info(msg string) {
	if varbose {
		fmt.Println(msg)
	}
}
