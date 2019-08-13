module github.com/Liquid-Labs/lc-authorizations-model

require (
	firebase.google.com/go v3.9.0+incompatible
	github.com/Liquid-Labs/go-rest v1.0.0-prototype.4
	github.com/Liquid-Labs/lc-authentication-api v0.0.0-20190812225013-10df0f6f9995
	github.com/Liquid-Labs/lc-entities-model v1.0.0-alpha.0
	github.com/Liquid-Labs/lc-rdb-service v1.0.0-alpha.1
	github.com/Liquid-Labs/terror v1.0.0-alpha.0
	github.com/go-pg/pg v8.0.5+incompatible
)

replace github.com/Liquid-Labs/lc-authentication-api => ../lc-authentication-api

replace github.com/Liquid-Labs/go-rest => ../go-rest

replace github.com/Liquid-Labs/lc-entities-model => ../lc-entities-model
