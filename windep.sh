#!/bin/bash
echo "Start to Vano deploy"
GOOS=windows GOARCH=amd64 go build -o "license.exe" cmd/apiserver/main.go
FILE=/home/rura/mnt/Vano/LicenseServer/license.exe
if [ -f "$FILE" ]; then
    echo "Mounted the server drive"
else
    echo "Mounting the server drive"
    sudo mount -t cifs -o username=Ivan,password=162747 \\\\192.168.115.120\\JanFant /home/rura/mnt/Vano
fi
sudo cp *.exe /home/rura/mnt/Vano/LicenseServer/
