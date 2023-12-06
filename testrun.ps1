# test file
go build -o target/main.exe main.go routers.go
if ($LASTEXITCODE -ne 0) {
    exit $LASTEXITCODE
}
target/main.exe
