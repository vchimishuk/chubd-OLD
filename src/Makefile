include $(GOROOT)/src/Make.inc

all: chubd

chubd: server.$(O) protocol.$(O) audio.$(O) mp3.$(O) ogg.$(O) config.$(O) playlist.$(O) player.$(O)
	$(GC) main.go
	$(LD) -o chubd main.$(O)

config.$(O): config/config.go
	$(GC) -o config.$(O) config/config.go

server.$(O): server.go
	$(GC) server.go

player.$(O): player.go playlist.$(O)
	$(GC) player.go

protocol.$(O): server.$(O) vfs.$(O) protocol.go player.$(O)
	$(GC) protocol.go

vfs.$(O): vfs/vfs.go vfs/track.go vfs/directory.go vfs/entry.go vfs/path.go audio.$(O) config.$(O)
	$(GC) -o vfs.$(O) vfs/vfs.go vfs/track.go vfs/directory.go vfs/entry.go vfs/path.go

playlist.$(O): playlist/playlist.go
	$(GC) -o playlist.$(O) playlist/playlist.go

audio.$(O): audio/tagreader.go audio/tag.go
	$(GC) -o audio.$(O) audio/tagreader.go audio/tag.go

mp3.$(O): mp3/tagreader.go audio.$(O)
	$(GC) -o mp3.$(O) mp3/tagreader.go

ogg.$(O): ogg/tagreader.go audio.$(O)
	$(GC) -o ogg.$(O) ogg/tagreader.go

clean:
	rm -f *.$(O) chubd

format:
	find . -type f -name '*.go' -exec gofmt -w {} \;

run: all
	./chubd
