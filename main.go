package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gosuri/uilive"

	"github.com/beykansen/power-meter/pkg"
)

// Change these lines according to your setup.
var (
	cpuPowerPath = "AMD Ryzen 7 3700X, Powers, CPU Package"
	gpuPowerPath = "NVIDIA NVIDIA GeForce RTX 2070 SUPER, Powers, GPU Power"
	usbCount     = 3
	hddCount     = 3
	fanCount     = 5
	pricePerkWh  = 0.85
	currency     = "TL"
	offset       = 77.0
	port         = "8085"
)

func main() {
	recordings := pkg.NewRecordings()
	go func() {
		time.AfterFunc(1*time.Hour, func() { recordings.Reset() })
	}()
	writer := uilive.New()
	writer.Start()
	defer writer.Stop()

	for {
		_, _ = fmt.Fprintf(writer, "LIVE UPDATE: %s", getRecording(recordings))
		time.Sleep(1 * time.Second)
	}
}

func getRecording(recordings *pkg.Recordings) string {
	totalWatt, err := pkg.Calculate(cpuPowerPath, gpuPowerPath, usbCount, hddCount, fanCount, port)
	if err != nil {
		log.Fatal(err)
	}
	totalWatt += offset
	totalKw := totalWatt / 1000

	recordings.Add(totalKw)

	currentElectricCostText := fmt.Sprintf("Your pc consumes %.3f kW %.3f W right now.", totalKw, totalWatt)
	dailyElectricCost := recordings.GetMedian() * 24 * pricePerkWh
	dailyElectricCostText := fmt.Sprintf("Today's approximate electricity cost %.3f %s", dailyElectricCost, currency)
	monthlyElectricCost := dailyElectricCost * 30
	monthlyElectricCostText := fmt.Sprintf("This Month's approximate electricity cost %.3f %s", monthlyElectricCost, currency)

	return fmt.Sprintf("%s %s %s", currentElectricCostText, dailyElectricCostText, monthlyElectricCostText)
}
