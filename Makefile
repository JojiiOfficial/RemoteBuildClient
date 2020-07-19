build:
	go mod download
	go build -o rbuild

default: build

upgrade:
	go mod download
	go get -u -v
	go mod tidy
	go mod verify

test:
	go test

man: build
	./rbuild --help-man | man -l -

install:
	@if ! test -f rbuild;then echo 'run "make build" first'; exit 1; fi

ifneq ($(shell id -u), 0)
	@echo "You must be root to perform this action."
	@exit 1
endif
	@mkdir -p /usr/local/share/man/man8
	cp rbuild /usr/bin/rbuild
	/usr/bin/rbuild --help-man > rbuild.1
	install -Dm644 rbuild.1 /usr/share/man/man8/rbuild.8
	@rm rbuild.1
	@echo Installed successfully!

uninstall:
ifneq ($(shell id -u), 0)
	@echo "You must be root to perform this action."
	@exit 1
endif
	rm /usr/bin/rbuild
	rm -f /usr/share/man/man8/rbuild.8
	@echo Uninstalled successfully!

clean:
	rm -f rbuild.1
	rm -f rbuild
	rm -f main
