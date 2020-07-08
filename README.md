# GoIPMailer


**Buildscripte**

Im Ordner `build` befinden sich diverse Buildscripte für Windows (cmd/powershell):
 * `BUILD.bat` - führt einen normalen build für das aktuelle Betriebssystem aus
 * `BUILD_ASSETS-DEBUG.bat` - erzeugt die `bindata.go` Datei für die Asset-Verwaltung (z.B. eingebettetes RDP-Template) als links
 * `BUILD_ASSETS-RELEASE.bat` - erzeugt die `bindata.go` Datei für die Asset-Verwaltung (z.B. eingebettetes RDP-Template) als Eingebettete Binär-Strings
 * `BUILDUPX.bat` - wie `Build.bat` mit anschluißender UPX-Komprimierung
 * `BUILDAll.ps1` - Powershell-Script welches Binaries für Windows/32, Win/64, Linux/32, Linux/64, Linux/ARM kompiliert und UPX-Komprimiert
 

**Konfiguration**  
Die Konfiguration erfolgt in der Datei `config.json`, z.B:
```
{
    "SMTPHost":"mail.smtpserver.de",
    "SMTPPort":25,
    "SMTPFROM":"goipmailer@smtpserver.de",
    "SMTPTO":"my.emailadress@smtpserver.de",
    "SMTPUser":"",
    "SMTPPass":"",
    "IPFilter":"192.168.1",
    "RDPUser":"\\LokalBenutzer",
    "RDPAttachFile":true,
    "RDPAttachBodyLine":true,
    "SSHUser":"",
    "SSHAttachBodyLine":false,
    "UserFile":"README.md",
    "UserAttachFile":false,
    "UserBodyLine":"Benutzerdefinierter Text (Kann auch <b>HTML</b> enthalten)</br>",
    "UserAttachBodyLine":true
}
```


----
### Installation als Dienst unter Windows:
```
sc create GoIPMailer binpath="P:\Go\gowiki\gowiki_amd64.exe" start=delayed-auto DisplayName="GoIPMailer"
sc start GoIPMailer
sc query GoIPMailer
```

Deinstallation:
```
sc stop GoIPMailer
sc delete GoIPMailer
```

----


### Installation Dienst unter Linux:

#### Systemd

Vorraussetzung: `systemctl enable systemd-networkd.service systemd-networkd-wait-online.service`

Systemd-Script `/lib/systemd/system/GoIPMailer.service`:
```
[Unit]
Description=GoIPMailer
Documentation=https://github.com/frickelblog/goipmailer
After=systemd-networkd-wait-online.service
Wants=systemd-networkd-wait-online.service

[Service]
Type=simple
ExecStart=/opt/frickelblog/GoIPMailer/GoIPMailer_amd64.bin
StandardOutput=null
Restart=on-failure
RestartSec=10s

[Install]
WantedBy=multi-user.target
Alias=GoIPMailer.service
```

Systemd Script enablen: `systemctl enable GoIPMailer.service`


----