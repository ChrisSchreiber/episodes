.PHONY: build clean

build: clean episodes

episodes:
	go build ./ 

clean:
	rm -f episodes