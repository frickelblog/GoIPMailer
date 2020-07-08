cd /d %~dp0
cd ..
go-bindata -verbose -o bindata.go ./assets/...
go build -ldflags="-s -w" -o GoIPMailer.exe
..\upx-3.96-win64\upx.exe -9 GoIPMailer.exe
go-bindata -debug -verbose -o bindata.go ./assets/...
pause