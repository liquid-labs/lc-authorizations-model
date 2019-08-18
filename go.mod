module github.com/Liquid-Labs/lc-authorizations-model

require (
	github.com/Liquid-Labs/lc-authentication-api v0.0.0-20190812225013-10df0f6f9995
	github.com/Liquid-Labs/lc-containers-model v0.0.0-20190818011142-594fc04ecdf1
	github.com/Liquid-Labs/lc-entities-model v1.0.0-alpha.0
	github.com/Liquid-Labs/lc-rdb-service v1.0.0-alpha.1
	github.com/Liquid-Labs/lc-users-model v1.0.0-alpha.0
	github.com/Liquid-Labs/strkit v0.0.0-20190818184832-9e3e35dcfc9c
	github.com/Liquid-Labs/terror v1.0.0-alpha.1
	github.com/go-pg/pg v8.0.5+incompatible
	github.com/stretchr/testify v1.4.0
)

replace github.com/Liquid-Labs/lc-authentication-api => ../lc-authentication-api

replace github.com/Liquid-Labs/go-rest => ../go-rest

replace github.com/Liquid-Labs/lc-entities-model => ../lc-entities-model

replace github.com/Liquid-Labs/lc-rdb-service => ../lc-rdb-service

replace github.com/Liquid-Labs/lc-users-model => ../lc-users-model

replace github.com/Liquid-Labs/terror => ../terror

replace github.com/Liquid-Labs/lc-containers-model => ../lc-containers-model
