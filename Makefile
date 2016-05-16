SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=deapexer

.DEFAULT_GOAL: $(BINARY)

$(BINARY):	$(SOURCES)
	go build -o ${BINARY} deapexer.go

.PHONY: install
install:
	mkdir -p /etc/deapexer
	cp ${BINARY} /usr/local/bin/${BINARY}
	cp config.json /etc/deapexer/config.json

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
