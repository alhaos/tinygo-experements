package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {

	time.Sleep(time.Second * 3)

	i2c, err := initI2c(machine.GP4, machine.GP5)
	if err != nil {
		deadLoopPrint(err)
	}

	data := make([]byte, 1)

	// Читаем из регистра 0xD0 (Chip ID)
	err = i2c.ReadRegister(uint8(0x76), 0xD0, data)
	if err != nil {
		fmt.Printf("Ошибка: %s\n", err.Error())
		return
	}

	id := data[0]
	fmt.Printf("Chip ID: 0x%X\n", id)

	switch id {
	case 0x58:
		fmt.Println("Это BMP280 (только давление и температура)")
	case 0x60:
		fmt.Println("Это BME280 (давление, температура и влажность)")
	default:
		fmt.Println("Неизвестное устройство")
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

func deadLoopPrint(err error) {
	for {
		fmt.Println(err.Error())
		time.Sleep(time.Second)
	}
}
