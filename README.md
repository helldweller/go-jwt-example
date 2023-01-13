# REST API for REACT app

![code checks](https://github.com/helldweller/go-jwt-example/actions/workflows/audit.yml/badge.svg)

## Develop

### In minikube (linux)

    minikube start
    eval $(minikube docker-env)
    kubectl config use-context minikube
    kubectl create ns go-api
    kubectl config set-context --current --namespace=go-api

    # DB
    helm -n go-api install \
        --set global.postgresql.auth.postgresPassword=<POSTGRES_ADMIN_PASSWORD> \
        --set global.postgresql.auth.username=go-api \
        --set global.postgresql.auth.password=<POSTGRES_PASSWORD for user go-api> \
        --set global.postgresql.auth.database=go-api \
        postgres bitnami/postgresql

    # APP
    cp skaffold/secrets-example.yaml skaffold/secrets.yaml
    # and edit secrets
    kubectl apply -n go-api -f skaffold/secrets.yaml
    skaffold dev

### Swagger

    echo 'export PATH=$PATH:$HOME/go/bin' >> $HOME/.profile && bash
    go install github.com/swaggo/swag/cmd/swag@latest
    swag init --dir ./src/internal -g app/app.go --parseInternal --parseDependency --output ./src/internal/docs

http://127.0.0.1:8080/swagger/index.html

## To Do

* promhttp not collecting metrics for all routes (github.com/zsais/go-gin-prometheus?)
* gin access log to json, use LOG_LEVEL, GIN_MODE=release | gin.SetMode(gin.ReleaseMode)
* console command for db migrate, try use golang-migrate/migrate and atlas
* gorm healthcheck (readiness/liveness/startup probes)
* package/main/internal/controllers/userController.go:63 ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)
* password validation
* Error: internal/middleware/requireAuth.go:34:2: this value of err is never used (SA4006)
* internal/middleware/requireAuth.go:36:16: error strings should not be capitalized (ST1005)
* {"error":"Failed to create user"} when user exists package/main/internal/controllers/userController.go:75 ERROR: duplicate key value violates unique constraint "users_email_key" (SQLSTATE 23505)
* swagger example param only req values
* logout token blacklist in redis
