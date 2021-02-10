# godmin

RestApi admin panel

### Default environments

	PORT=8080
	LOG_LEVEL=debug

    # database
    DATABASE_HOST=localhost
    DATABASE_PORT=5432
    DATABASE_NAME=godmin_db_dev
    DATABASE_USER=godmin
    DATABASE_PASSWORD=password
    DATABASE_SSL_MODE=disable
    DATABASE_MAX_OPEN_CONNS=1000
    DATABASE_MAX_IDLE_CONNS=15
    DATABASE_CONN_MAX_IDLE_TIME=30s
    DATABASE_CONN_MAX_LIFE_TIME=0

    # jwt
    JWT_ACCESS_SECRET=secret;)
    JWT_REFRESH_SECRET=secret;)

    #redis
	REDIS_URL=localhost:6379

### TODO

- Tests
- CRUD
  - models
  - filters
  - pagination (without counter)