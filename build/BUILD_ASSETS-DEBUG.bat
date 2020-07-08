cd /d %~dp0
cd ..
go-bindata -debug -verbose -o bindata.go ./assets/...
pause