cd /d %~dp0
cd ..
go-bindata -verbose -o bindata.go ./assets/...
pause