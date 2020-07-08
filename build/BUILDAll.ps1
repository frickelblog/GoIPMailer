cd ..
echo "============================"
echo "   GOIPMailer BuildScript"
echo "============================"

echo "Create Static Ressources (Release)"
go-bindata -verbose -o bindata.go ./assets/...

echo "Build Windows/amd64..."
$env:GOOS="windows"; $env:GOARCH="amd64"; go build -ldflags="-s -w -X main.version=1.0.0" -o ./build/GoIPMailer_amd64.exe
echo "Build Windows/x86..."
$env:GOOS="windows"; $env:GOARCH="386"; go build -ldflags="-s -w -X main.version=1.0.0" -o ./build/GoIPMailer_x86.exe
echo "Build Linux/amd64..."
$env:GOOS="linux"; $env:GOARCH="amd64"; go build -ldflags="-s -w -X main.version=1.0.0" -o ./build/GoIPMailer_amd64.bin
echo "Build Linux/x86..."
$env:GOOS="linux"; $env:GOARCH="386"; go build -ldflags="-s -w -X main.version=1.0.0" -o ./build/GoIPMailer_x86.bin
echo "Build Linux/arm..."
$env:GOOS="linux"; $env:GOARCH="arm"; go build -ldflags="-s -w -X main.version=1.0.0" -o ./build/GoIPMailer.arm

echo "UPX - Compress Static Binaries"
..\upx-3.96-win64\upx.exe -9 ./build/GoIPMailer_amd64.exe
..\upx-3.96-win64\upx.exe -9 ./build/GoIPMailer_x86.exe
..\upx-3.96-win64\upx.exe -9 ./build/GoIPMailer_amd64.bin
..\upx-3.96-win64\upx.exe -9 ./build/GoIPMailer_x86.bin
..\upx-3.96-win64\upx.exe -9 ./build/GoIPMailer.arm

echo "Create Static Ressources (Debug)"
go-bindata -verbose -debug -o bindata.go ./assets/...
pause