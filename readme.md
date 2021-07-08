# Power-Meter

Power-Meter is a OpenHardwareMonitor client that calculates approximate electricity cost and consumption. This program is not for exact calculation. If you need exact consumption metrics, consider buying a kill-a-watt device.

## Installation

Firstly, install [openhardwaremonitor](https://openhardwaremonitor.org/downloads/) and make sure remote web server on port `8085` enabled. If you change the port, do not forget to change in global variables too.

![openhardwaremonitor](./images/tutorial.png)

You need to set your device sensor paths and other variables for proper calculation.

```go
    var (
        cpuPowerPath = "AMD Ryzen 7 3700X, Powers, CPU Package"
        gpuPowerPath = "NVIDIA NVIDIA GeForce RTX 2070 SUPER, Powers, GPU Power"
        usbCount     = 3
        hddCount     = 3
        fanCount     = 5
        pricePerkWh  = 0.85
        currency     = "TL"
        offset       = 0.0
        port         = "8085"
    )
```

```bash
go install .
```

## Usage

```bash
power-meter
```
## Roadmap
- [ ] Get variable from cmd args.
- [ ] Dynamic parsing according to path.


## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)