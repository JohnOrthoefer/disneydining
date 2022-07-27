all: diningUpdate diningWeb

diningUpdate:
	$(MAKE) -C cmd/diningUpdate all

diningWeb:
	$(MAKE) -C cmd/diningWeb all

clean:
	$(MAKE) -C cmd/diningUpdate clean
	$(MAKE) -C cmd/diningWeb clean

