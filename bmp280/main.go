package main

import (
	"errors"
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/bmp280"
	"tinygo.org/x/drivers/ssd1306"
)

const varbose = false

func main() {

	time.Sleep(time.Second * 2)

	info("start delay success")

	i2c, err := initI2c(machine.GP4, machine.GP5)
	if err != nil {
		deadLoopPrint(err)
	}

	sensor := bmp280.New(i2c)
	sensor.Address = 0x76 // Явно указываем адрес

	// 1. Сначала проверяем соединение
	if !sensor.Connected() {
		deadLoopPrint(errors.New("bmp280 sensor not connected"))
	}

	// 2. Только потом конфигурируем
	sensor.Configure(
		bmp280.STANDBY_1000MS,
		bmp280.FILTER_16X,
		bmp280.SAMPLING_16X,
		bmp280.SAMPLING_16X,
		bmp280.MODE_NORMAL,
	)

	display, err := initOLED(i2c)
	if err != nil {
		deadLoopPrint(err)
	}

	for {

		temperature, err := sensor.ReadTemperature()
		if err != nil {
			deadLoopPrint(err)
		}
		pressure, err := sensor.ReadPressure()
		if err != nil {
			deadLoopPrint(err)
		}

		temperatureMsg := fmt.Sprintf("%-20s", fmt.Sprintf("t: %d гр.", temperature/1000))
		pressureMsg := fmt.Sprintf("%-20s", fmt.Sprintf("p: %d ммРс", int32(float32(pressure)*0.75006)))

		DrawText(display, 10, 10, temperatureMsg, color.RGBA{255, 255, 255, 255})
		DrawText(display, 10, 20, pressureMsg, color.RGBA{255, 255, 255, 255})
		display.Display()

		time.Sleep(time.Second)
	}
}

func DrawChar(display *ssd1306.Device, x, y int16, ch rune, c color.RGBA) {
	data, ok := font5x8[ch]
	if !ok {
		data = font5x8['?']
	}

	// fmt.Printf("found char [%c]\n", ch)

	for colNum, col := range data {
		xd := x + int16(colNum)
		info(fmt.Sprintf("%d: %08b %d\n", colNum, col, xd))
		for rowNum := range 8 {
			if (col>>rowNum)&1 == 1 {
				display.SetPixel(xd, y+int16(rowNum), c)
			} else {
				display.SetPixel(xd, y+int16(rowNum), color.RGBA{0, 0, 0, 255})
			}
		}
	}

	fmt.Println()
}

func DrawText(display *ssd1306.Device, x, y int16, text string, c color.RGBA) {
	for i, r := range []rune(text) {
		info(fmt.Sprintf("rune: [%c], i: %d\n", r, i))
		DrawChar(display, x+int16(i*6), y, r, c)
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
