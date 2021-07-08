package pkg

import (
	"fmt"
	"strconv"
	"strings"
)

func Calculate(cpuPowerPath string, gpuPowerPath string, usbCount int, hddCount int, fanCount int, port string) (float64, error) {
	metrics, err := GetMetrics(port)
	if err != nil {
		return 0.0, fmt.Errorf("open hardware monitor cannot be found on https://localhost:8085. please be sure open hardware monitor is running and webserver enabled on 8085 port")
	}

	rawCpuPowerText, err := parseMetrics(strings.Split(cpuPowerPath, ","), metrics)
	if err != nil {
		return 0.0, err
	}
	rawGpuPowerText, err := parseMetrics(strings.Split(gpuPowerPath, ","), metrics)
	if err != nil {
		return 0.0, err
	}
	cpuPower, err := parseAndClearPower(rawCpuPowerText)
	if err != nil {
		return 0.0, err
	}
	gpuPower, err := parseAndClearPower(rawGpuPowerText)
	if err != nil {
		return 0.0, err
	}

	//approximate watt usages
	fanWatt := 2.3
	hddWatt := 5.50
	usbWatt := 2.5
	totalFanPower := float64(fanCount) * fanWatt
	totalHddWatt := float64(hddCount) * hddWatt
	totalUsbWatt := float64(usbCount) * usbWatt

	return cpuPower + gpuPower + totalFanPower + totalHddWatt + totalUsbWatt, nil
}

func parseAndClearPower(power string) (float64, error) {
	power = strings.TrimSpace(power)
	power = strings.ReplaceAll(power, ",", ".")
	power = strings.ReplaceAll(power, "W", "")
	return strconv.ParseFloat(strings.TrimSpace(power), 64)
}

func parseMetrics(pathVar []string, metrics *ResponseDto) (string, error) {
	for _, children := range metrics.Children[0].Children {
		if children.Text == strings.TrimSpace(pathVar[0]) {
			for _, children2 := range children.Children {
				if children2.Text == strings.TrimSpace(pathVar[1]) {
					for _, children3 := range children2.Children {
						if children3.Text == strings.TrimSpace(pathVar[2]) {
							return children3.Value, nil
						}
					}
				}
			}
		}
	}
	return "", fmt.Errorf("power not found")
}
