include $(GOROOT)/src/Make.inc

all: chubd

chubd: server.$(O) protocol.$(O) audio.$(O) mp3.$(O) ogg.$(O) config.$(O) playlist.$(O) player.$(O) utils.$(O)
	$(GC) main.go
	$(LD) -o chubd main.$(O)

config.$(O): config/config.go
	$(GC) -o config.$(O) config/config.go

server.$(O): server.go
	$(GC) server.go

player.$(O): player/player.go player/playingroutine.go playlist.$(O) audio.$(O) ogg.$(O) alsa.$(O)
	$(GC) -o player.$(O) player/player.go player/playingroutine.go

protocol.$(O): server.$(O) vfs.$(O) protocol.go player.$(O)
	$(GC) protocol.go

vfs.$(O): vfs/vfs.go vfs/track.go vfs/directory.go vfs/entry.go vfs/path.go audio.$(O) config.$(O)
	$(GC) -o vfs.$(O) vfs/vfs.go vfs/track.go vfs/directory.go vfs/entry.go vfs/path.go

playlist.$(O): playlist/playlist.go
	$(GC) -o playlist.$(O) playlist/playlist.go

audio.$(O): audio/decoder.go audio/output.go audio/tagreader.go audio/tag.go
	$(GC) -o audio.$(O) audio/decoder.go audio/output.go audio/tagreader.go audio/tag.go

mp3.$(O): mp3/tagreader.go audio.$(O)
	$(GC) -o mp3.$(O) mp3/tagreader.go

ogg.$(O): ogg/ogg.go ogg/tagreader.go ogg/decoder.go audio.$(O) utils.$(O)
	$(GC) -o ogg.$(O) ogg/ogg.go ogg/tagreader.go ogg/decoder.go

alsa.$(O): alsa/alsa.go audio.$(O)
	$(GC) -o alsa.$(O) alsa/alsa.go

utils.$(O): utils/utils.go 
	$(GC) -o utils.$(O) utils/utils.go

clean:
	rm -f *.$(O) chubd

format:
	find . -type f -name '*.go' -exec gofmt -w {} \;

run: all
	./chubd

kill:
	killall chubd
