include .env
export $(shell sed 's/=.*//' .env)


run:
	@go run cmd/*.go http
dev:
	go run cmd/*.go dev
tel:
	go run cmd/*.go telegram

seed:
	@go run cmd/*.go seed
