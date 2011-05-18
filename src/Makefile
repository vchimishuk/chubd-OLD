include $(GOROOT)/src/Make.inc

all: chubd

chubd: server.$(O) protocol.$(O)
	$(GC) main.go
	$(LD) -o chubd main.$(O)

server.$(O): server.go
	$(GC) server.go

protocol.$(O): server.$(O) filesystem.$(O) protocol.go
	$(GC) protocol.go

filesystem.$(O): filesystem.go track.go directory.go entry.go
	$(GC) filesystem.go track.go directory.go entry.go

clean:
	rm -f *.$(O) chubd

format:
	gofmt -w *.go
