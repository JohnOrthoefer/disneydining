all: diningUpdate diningWeb install

diningUpdate:
	$(MAKE) -C cmd/diningUpdate all

diningWeb:
	$(MAKE) -C cmd/diningWeb all

install:
	-mkdir bin
	cp cmd/diningUpdate/diningUpdate bin/
	cp cmd/diningWeb/diningWeb bin/

clean:
	$(MAKE) -C cmd/diningUpdate clean
	$(MAKE) -C cmd/diningWeb clean

