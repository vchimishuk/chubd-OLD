include $(GOROOT)/src/Make.inc

all: chubd

chubd: server.$(O) protocol.$(O) audio.$(O) mp3.$(O) config.$(O)
	$(GC) main.go
	$(LD) -o chubd main.$(O)

config.$(O): config/config.go
	$(GC) -o config.$(O) config/config.go

server.$(O): server.go
	$(GC) server.go

protocol.$(O): server.$(O) vfs.$(O) protocol.go
	$(GC) protocol.go

vfs.$(O): vfs/vfs.go vfs/track.go vfs/directory.go vfs/entry.go audio.$(O) config.$(O)
	$(GC) -o vfs.$(O) vfs/vfs.go vfs/track.go vfs/directory.go vfs/entry.go

audio.$(O): audio/tagreader.go audio/tag.go
	$(GC) -o audio.$(O) audio/tagreader.go audio/tag.go

mp3.$(O): mp3/tagreader.go audio.$(O)
	$(GC) -o mp3.$(O) mp3/tagreader.go

clean:
	rm -f *.$(O) chubd

format:
	find . -type f -name '*.go' -exec gofmt -w {} \;

run: all
	./chubd
