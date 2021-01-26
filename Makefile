.PHONY: all

all:
	echo "linux and windows targets available"

.PHONY: linux
linux:
	go build -o mandelbrot .

.PHONY: windows
windows:
	CC=x86_64-w64-mingw32-gcc.exe go build -ldflags '-w -extldflags "-static"' .
