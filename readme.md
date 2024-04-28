# Install dependency
go get -d ./...

# Create Makefile
Ini contoh Makefile

run: build
    @cd ./bin/
    @main.exe <- sesuaikan kalo di windows pake .exe, kalo linux ./main kalo gak salah

build:
    @go build -o ./bin/main.exe ./cmd/api/main.go <- yang disini juga kalo dilinux jadi ./bin/main

test:
    @go test ./...


# Run project

Linux:
```
make run
```

Windows:
```
nmake run
```
