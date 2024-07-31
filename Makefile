soulkiller: main.go
	go build -o soulkiller main.go

install: soulkiller
	sudo install -m 0755 -o root -g root soulkiller /sbin/soulkiller
	sudo systemctl restart soulkiller

