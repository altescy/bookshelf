NAME    := bookshelf
PWD     := $(shell pwd)
GOCMD   := go
GOBUILD := $(GOCMD) build
GOTEST  := $(GOCMD) test
SOURCE  := $(PWD)
TARGET  := $(PWD)/bin/$(NAME)

$(TARGET):
	$(GOBUILD)  -o $(TARGET) $(SOURCE)

.PHONY: run
run: $(TARGET)
	$(TARGET)

.PHONY: test
test:
	$(GOTEST) $(PWD)/...


.PHONY: clean
clean:
	rm -rf $(PWD)/bin

all: clean $(TARGET)
