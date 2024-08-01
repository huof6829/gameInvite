# Go parameters
GOCMD=GO111MODULE=on go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test

.PHONY: build stop show log api mysql redis

build: 
	go build backend.go

runf:
	go run backend.go -f=etc/backend-api.yaml

rund: 
	go build backend.go 
	nohup /home/mart/gameInvite/backend -f=/home/mart/gameInvite/etc/backend-api.yaml > /home/mart/gameInvite/run.log 2>&1 &

stop: 
	pkill -f backend

show:
	netstat -tunlp|grep -e backend

log:
	tail -f run.log

api:
	goctl api go --api backend.api --dir . --style go_zero

mysql:
	mysql -uadmin -pYHt6_8YeHiG_C?ixZP

redis:
	redis-cli -h 127.0.0.1 -p 6379 

redisdelid:
	redis-cli keys "cache:userBind:id:*" | xargs redis-cli del

# goctl model mysql datasource -style=go_zero -c -url="admin:YHt6_8YeHiG_C?ixZP@tcp(127.0.0.1:3306)/savvy_gameing" --home=/home/mart/.goctl/backend -t=sys_invite  -dir=internal/model/sys_invite 