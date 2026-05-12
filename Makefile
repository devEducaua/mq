
TARGET = smpq
TARGETDIR = bin

SRCS = $(wildcard *go)

all: $(TARGET)

$(TARGET): $(SRCS)
	mkdir -p $(TARGETDIR)
	go build -o $(TARGETDIR)/$@ .

clean:
	rm -r $(TARGETDIR)
