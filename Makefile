PREFIX ?= /usr/local
BINDIR = $(PREFIX)/bin
MANDIR = $(PREFIX)/share/man/man1

TARGET = mq
TARGETDIR = bin
MANPAGE = mq.1

SRCS = $(wildcard *go)

all: $(TARGET)

$(TARGET): $(SRCS)
	mkdir -p $(TARGETDIR)
	go build -o $(TARGETDIR)/$@ .

install:
	cp $(TARGETDIR)/$(TARGET) $(BINDIR)/
	cp ./$(MANPAGE) $(MANDIR)/
	gzip -f $(MANDIR)/$(MANPAGE)

uninstall:
	rm -f $(BINDIR)/$(TARGET)
	rm -f $(MANDIR)$(MANPAGE).gz

clean:
	rm -r $(TARGETDIR)
