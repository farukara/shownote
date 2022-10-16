.PHONY: all
all: sn

sn: main.go
	go build -o $@ $^
	cp $@ /usr/local/bin/
	mkdir -pv ~/.config/shownote
	cp ./config.yaml ~/.config/shownote/
	rm sn
