cd /d %~dp0
cd ..
go build -ldflags="-s -w" -o GoIPMailer.exe
rem ..\upx-3.95-win64\upx.exe -9 GoIPMailer.exe
pause