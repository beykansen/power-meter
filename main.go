package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"

	"github.com/beykansen/power-meter/pkg"
)

type ContextData struct {
	cpuPowerPath string
	gpuPowerPath string
	usbCount     int
	hddCount     int
	fanCount     int
	pricePerKwh  float64
	currency     string
	offset       float64
	port         string
}

func main() {
	cpuPowerPath := flag.String("cpu-power-path", "AMD Ryzen 7 3700X_Powers_CPU Package", "cpu power metric path on open hardware monitor gui seperated with under score")
	gpuPowerPath := flag.String("gpu-power-path", "NVIDIA NVIDIA GeForce RTX 2070 SUPER_Powers_GPU Power", "gpu power metric path on open hardware monitor gui seperated with under score")
	usbCount := flag.Int("usb-count", 4, "usb count")
	hddCount := flag.Int("hdd-count", 3, "hdd count")
	fanCount := flag.Int("fan-count", 5, "fan count")
	pricePerKwh := flag.Float64("price", 0.85, "price per kWh")
	currency := flag.String("currency", "TL", "currency")
	offset := flag.Float64("offset", 77.0, "positive or negative offset watt for ram, motherboard and other peripherals that don't too much fluctuate power consumption according to your usage.")
	port := flag.String("port", "8085", "open hardware monitor web service port")
	flag.Parse()

	recordings := pkg.NewRecordings()
	go func() {
		time.AfterFunc(1*time.Hour, func() { recordings.Reset() })
	}()

	contextData := &ContextData{
		cpuPowerPath: *cpuPowerPath,
		gpuPowerPath: *gpuPowerPath,
		usbCount:     *usbCount,
		hddCount:     *hddCount,
		fanCount:     *fanCount,
		pricePerKwh:  *pricePerKwh,
		currency:     *currency,
		offset:       *offset,
		port:         *port,
	}
	drawUi(recordings, contextData)
}

func drawUi(recordings *pkg.Recordings, contextData *ContextData) {
	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize ui: %v", err)
	}
	defer ui.Close()
	metrics := getMetrics(recordings, contextData)

	liveIndicator := widgets.NewParagraph()
	liveIndicator.Title = "POWER-METER"
	liveIndicator.Text = "LIVE! PRESS q TO QUIT"
	liveIndicator.SetRect(0, 0, 70, 3)
	liveIndicator.TextStyle.Fg = ui.ColorWhite
	liveIndicator.BorderStyle.Fg = ui.ColorCyan
	updateParagraph := func(count int) {
		if count%2 == 0 {
			liveIndicator.TextStyle.Fg = ui.ColorRed
		} else {
			liveIndicator.TextStyle.Fg = ui.ColorWhite
		}
	}

	resultText := widgets.NewParagraph()
	resultText.Title = "SUMMARY"
	resultText.Text = metrics.text
	resultText.SetRect(0, 3, 70, 10)
	resultText.TextStyle.Fg = ui.ColorWhite
	resultText.BorderStyle.Fg = ui.ColorCyan

	currentConsumption := widgets.NewPlot()
	currentConsumption.Title = "CURRENT CONSUMPTION (W)"
	currentConsumption.Data = make([][]float64, 1)
	currentConsumption.Data[0] = []float64{metrics.currentWatt}
	currentConsumption.SetRect(0, 10, 70, 20)
	currentConsumption.AxesColor = ui.ColorWhite
	currentConsumption.Marker = widgets.MarkerDot
	currentConsumption.BorderStyle.Fg = ui.ColorCyan

	costs := widgets.NewPlot()
	costs.Title = fmt.Sprintf("COSTS (%s)", contextData.currency)
	costs.Data = make([][]float64, 2)
	costs.Data[0] = []float64{metrics.dailyElectricCost}
	costs.Data[1] = []float64{metrics.monthlyElectricCost}
	costs.SetRect(0, 40, 70, 20)
	costs.AxesColor = ui.ColorWhite
	costs.Marker = widgets.MarkerBraille
	costs.LineColors[0] = ui.ColorCyan
	costs.LineColors[1] = ui.ColorYellow

	draw := func(count int) {
		metrics := getMetrics(recordings, contextData)
		resultText.Text = metrics.text
		currentConsumption.Data[0] = append(currentConsumption.Data[0], metrics.currentWatt)
		costs.Data[0] = append(costs.Data[0], metrics.dailyElectricCost)
		costs.Data[1] = append(costs.Data[1], metrics.monthlyElectricCost)
		ui.Render(liveIndicator, resultText, currentConsumption, costs)
	}

	tickerCount := 1
	draw(tickerCount)
	tickerCount++
	uiEvents := ui.PollEvents()
	ticker := time.NewTicker(time.Second).C
	for {
		select {
		case e := <-uiEvents:
			switch e.ID {
			case "q", "<C-c>":
				return
			}
		case <-ticker:
			updateParagraph(tickerCount)
			draw(tickerCount)
			tickerCount++
		}
	}
}

type getRecordingResult struct {
	text                string
	currentWatt         float64
	dailyElectricCost   float64
	monthlyElectricCost float64
}

func getMetrics(recordings *pkg.Recordings, contextData *ContextData) *getRecordingResult {
	result := &getRecordingResult{}
	totalWatt, err := pkg.Calculate(contextData.cpuPowerPath, contextData.gpuPowerPath, contextData.usbCount, contextData.hddCount, contextData.fanCount, contextData.port)
	if err != nil {
		log.Fatal(err)
	}
	totalWatt += contextData.offset
	totalKw := totalWatt / 1000

	recordings.Add(totalKw)
	result.currentWatt = totalWatt

	currentElectricCostText := fmt.Sprintf("Your pc consumes %.3f kW %.3f W right now.", totalKw, totalWatt)

	result.dailyElectricCost = recordings.GetMedian() * 24 * contextData.pricePerKwh
	dailyElectricCostText := fmt.Sprintf("Today's approximate electricity cost %.3f %s", result.dailyElectricCost, contextData.currency)

	result.monthlyElectricCost = result.dailyElectricCost * 30
	monthlyElectricCostText := fmt.Sprintf("This Month's approximate electricity cost %.3f %s", result.monthlyElectricCost, contextData.currency)

	result.text = fmt.Sprintf("%s %s %s", currentElectricCostText, dailyElectricCostText, monthlyElectricCostText)
	return result
}
