build:
	go build -o bin/qc .

.PHONY: install
install: build
	if [ -f ~/.local/bin/qc ]; then rm ~/.local/bin/qc; fi
	cp ./bin/qc ~/.local/bin/qc
