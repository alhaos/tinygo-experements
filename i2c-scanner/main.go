package main

import (
	"fmt"
	"machine"
	"time"
)

func main() {

	time.Sleep(time.Second * 2)

	i2c := machine.I2C0
	err := i2c.Configure(machine.I2CConfig{
		Frequency: 400 * machine.KHz,
		SDA:       machine.GP4,
		SCL:       machine.GP5,
	})

	if err != nil {
		deadLoopErr(err)
	}

	performFullScan(i2c)
}

func performFullScan(i2c *machine.I2C) {

	print("Full scan\n")
	found := 0

	for addr := uint8(0x08); addr < uint8(0x77); addr++ {
		if checkAddress(i2c, addr) {
			printDeviceInfo(addr)
			found++
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func checkAddress(i2c *machine.I2C, addr uint8) bool {
	err := i2c.WriteRegister(addr, 0x00, []byte{})
	return err == nil
}

func printDeviceInfo(addr uint8) {

	fmt.Printf("0x%02X", addr)

	// Определяем тип устройства по адресу
	switch addr {
	case 0x1D, 0x1E:
		fmt.Printf(" - Акселерометр/магнитометр (FXOS8700, LSM303, ADXL345)")
	case 0x20, 0x21, 0x22, 0x23, 0x24, 0x25, 0x26, 0x27:
		fmt.Printf(" - I/O расширитель (PCF8574, PCF8574A, MCP23008)")
	case 0x3C, 0x3D:
		fmt.Printf(" - OLED дисплей (SSD1306)")
	case 0x40:
		fmt.Printf(" - Датчик температуры/влажности (BMP180, PCA9685)")
	case 0x48, 0x49, 0x4A, 0x4B:
		fmt.Printf(" - АЦП (PCF8591, ADS1115)")
	case 0x50, 0x51, 0x52, 0x53, 0x54, 0x55, 0x56, 0x57:
		fmt.Printf(" - EEPROM (24Cxx)")
	case 0x68:
		fmt.Printf(" - RTC или гироскоп (DS1307, MPU6050/9250)")
	case 0x76, 0x77:
		fmt.Printf(" - Датчик давления/влажности (BME280, BMP280)")
	default:
		fmt.Printf(" - Неизвестное устройство")
	}

	fmt.Printf(" (десятичное: %d)\r\n", addr)

}

func deadLoopErr(err error) {
	for {
		print(err.Error())
		time.Sleep(time.Second)
	}
}
