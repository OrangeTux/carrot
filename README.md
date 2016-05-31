# Carrot
Carrot is a small Go tool to read data from the P1 port of my smart meter and
write it to InfluxDB.

# Usage
```
Usage of ./carrot:
  -baudrate int
        specify the baudrate (default 115200)
  -device string
        specify the device (default "/dev/ttyUSB0")
  -influx-db string
        specify InfluxDB database (default "rommel")
  -influx-host string
        specify host of InfluxDB (default "http://mon.mijnbaopt.nl:8086")
  -influx-password string
        specify InfluxDB password
  -influx-user string
        specify InfluxDB user
```

# License
MIT
