all: diningUpdate diningWeb diningNotify install

diningUpdate:
	$(MAKE) -C cmd/diningUpdate all

diningWeb:
	$(MAKE) -C cmd/diningWeb all

diningNotify:
	$(MAKE) -C cmd/diningNotify all

install:
	-mkdir bin
	cp cmd/diningUpdate/diningUpdate bin/
	cp cmd/diningWeb/diningWeb bin/
	cp cmd/diningNotify/diningNotify bin/

clean:
	$(MAKE) -C cmd/diningUpdate clean
	$(MAKE) -C cmd/diningWeb clean
	$(MAKE) -C cmd/diningNotify clean
	rm bin/*

