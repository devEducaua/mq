PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin

TARGET = mq
TARGETDIR = bin

SRCS = $(wildcard *go)

all: $(TARGET)

$(TARGET): $(SRCS)
	mkdir -p $(TARGETDIR)
	go build -o $(TARGETDIR)/$@ .

install:
	cp $(TARGETDIR)/$(TARGET) $(BINDIR)/

uninstall:
	rm $(BINDIR)/$(TARGET)

clean:
	rm -r $(TARGETDIR)
