build:
	docker-compose build

stop:
	docker-compose down

clean:
	docker-compose down -v

run:
	docker-compose up -d

rerun: | stop run

start: | build run

restart: | stop start

fresh: | clean start
