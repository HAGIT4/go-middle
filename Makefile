generate-swagger:
	swag init --o ./docs -d ./internal/server/api -g router.go --parseDependency