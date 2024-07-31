BIN := soulkiller
all: $(BIN) install 

$(BIN): main.go
	go build -o soulkiller main.go

install: $(BIN)
	sudo install -m 0755 -o root -g root soulkiller /sbin/soulkiller
	sudo systemctl restart soulkiller

