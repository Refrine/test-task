include .env
export 


env-up:
	docker compose up postgres

env-down:
	docker compose down postgres