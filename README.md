# Door

## Build

```
./scripts/build.sh
```

## Install

```
scp door doorpi:
scp door.service doorpi:/etc/systemd/system/door.service
# ssh doorpi
# replace username/password in service file
sudo cp door.service /etc/systemd/system/door.service
sudo chown root:root /etc/systemd/system/door.service
sudo systemctl start door
sudo systemctl enable door # start on boot
```

## Raspberry Pi Config

1. Port forward 80/443 to the device
1. [Disable wifi sleep mode](https://www.heelpbook.net/2021/raspberry-pi-4-preventing-wifi-module-to-go-to-sleep-mode/)
