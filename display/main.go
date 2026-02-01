package main

import (
	"fmt"
	"image/color"
	"machine"
	"time"

	"tinygo.org/x/drivers/ssd1306"
)

func main() {
	time.Sleep(time.Second * 2)
	fmt.Print("start delay success")

	i2c, err := initI2c(machine.GP4, machine.GP5)
	if err != nil {
		deadLoopErr(err)
	}

	display, err := initOLED(i2c)
	if err != nil {
		deadLoopErr(err)
	}

	x := int16(0)
	moveRight := true
	white := color.RGBA{255, 255, 255, 255}
	black := color.RGBA{0, 0, 0, 255}

	// Очищаем весь экран в начале
	display.FillRectangle(0, 0, 128, 64, black)
	display.Display()

	for {
		// Стираем только старую позицию квадрата
		display.FillRectangle(x, 0, 16, 16, black)

		// Обновляем позицию
		if moveRight {
			x++
		} else {
			x--
		}

		if x <= 0 {
			moveRight = true
			x = 0
		}

		if x >= 112 {
			moveRight = false
			x = 112
		}

		// Рисуем квадрат в новой позиции
		display.FillRectangle(x, 0, 16, 16, white)

		// Обновляем дисплей
		display.Display()

		time.Sleep(time.Millisecond * 20) // ~50 FPS
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

func deadLoopErr(err error) {
	for {
		print(err.Error())
		time.Sleep(time.Second)
	}
}
