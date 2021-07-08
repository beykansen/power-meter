# Power-Meter

Power-Meter is an OpenHardwareMonitor client that calculates approximate electricity cost and consumption. This program is not for exact calculation. If you need exact consumption metrics, consider buying a kill-a-watt device.

## Installation

Firstly, install [openhardwaremonitor](https://openhardwaremonitor.org/downloads/) and make sure remote web server on port `8085` enabled. If you change the port, do not forget to change in global variables too. OpenHardwareMonitor should be running background.

![openhardwaremonitor](./images/tutorial.png)

You need to set your device sensor paths and other variables for proper calculation. You can split children with comma. This means that ``Cpu Package`` under ``Power`` under ``AMD Ryzen 7 3700X`` on OpenHardwareMonitor GUI. 

You can add positive or negative offset watt for ram, motherboard and other peripherals that don't too much fluctuate power consumption according to your usage.
- Around 3 W of power for every 8 GB of DDR3 or DDR4 memory can be allocated.
- Around 30 W for regular motherboard, 65 W for high-end motherboard can be allocated.
- According to these; I've 32 GB DDR4 ram with high-end motherboard; so I've set (3 W * 4) for ram and 65 W for motherboard.
```go
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
```

```bash
go install .
```

## Usage

```bash
power-meter
```
## Roadmap & Known Issues
- [ ] Get variable from cmd args.
- [ ] Dynamic parsing according to path.
- [ ] Live update doesn't work properly with Windows CMD and Powershell. Please, if you are running on Windows, consider using Git Bash right now.

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## License
[MIT](https://choosealicense.com/licenses/mit/)