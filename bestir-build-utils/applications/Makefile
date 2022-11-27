# it no worky
set-connection-string:
	DATABASE_URL="postgresql://prometheus:AmV_tvERFYLSvxfhsGiZig@unique-quokka-4306.6wr.cockroachlabs.cloud:26257/defaultdb?sslmode=verify-full"

# needs escape character I think
get-connection-string:
	echo $DATABASE_URL

local-cockroach-exec:
	cockroach sql --insecure --host=localhost:26257

local-cockroach-deps:
	echo "no not that kinda roach"

deps:
	go get github.com/badoux/checkmail
	go get github.com/jinzhu/gorm
	go get golang.org/x/crypto/bcrypt
	go get github.com/dgrijalva/jwt-go
	go get github.com/gorilla/mux
	go get github.com/jinzhu/gorm/dialects/mysql" //If using mysql 
	go get github.com/jinzhu/gorm/dialects/postgres //If using postgres
	go get github.com/joho/godotenv
	go get gopkg.in/go-playground/assert.v1

mod:
	go mod tidy 
	go mod vendor 

local-build-up:
	echo "dockerize this bitch s'il vouz plait"

local-build-down:
	echo "go to sleep"

local-test:
	echo "congrats you're HIV positive!"

# we're gonna wanna parameterize this bad boi later
local-db-up:
	 docker run -d --name ms -p 3306:3306 -e MYSQL_ROOT_PASSWORD=password mysql
# docker exec -d ms bash mysql-u root -ppassword # i don't think I actually need this command since I'm not mannually logging in

local-db-migrate:
	echo "is this even necessary for GORM? or is this only for schema first ORMs like sql-c?"
	 

# okay so it looks like we've managed to login successfully w/ the script
# now we need to find the right connection string 

# this need a check to make sure that mysql is actuall up 
local-db-down:
	docker stop ms 
	docker rm ms


# idk but the mod target was breaking
# local-env-up: mod local-db-up local-db-migrate local-build-up local-test 
local-env-up: local-db-up local-db-migrate local-build-up local-test 
	

local-env-down: local-build-down local-db-down 

local-stop:
	docker kill permissionmaking

local-start:
	docker run --name permissionmaking -d \
		--network=bridge --rm \
		-e MYSQL_USER=usr \
		-e MYSQL_PASSWORD=permission \
		-e MYSQL_DATABASE=permission \
		-e MYSQL_ALLOW_EMPTY_PASSWORD=yes \
		-p 3306:3306 mysql

	DD_DISABLE=true \
		API_ADDRESS=0.0.0.0:8080 \
		permission_DB_USER=usr \
		permission_DB_PASSWORD=permission \
		permission_DB_HOST=127.0.0.1 \
		permission_DB_NAME=permission \
		permission_DB_PORT=3306 \
		ENABLE_MIGRATE=true \
		go run ./cmd/bestir-permissionmaking-service/

local-restart: local-stop local-start