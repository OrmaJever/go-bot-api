cd /var/bot
[ ! -f "main" ] && go mod download && go build -o main ./cmd

./main