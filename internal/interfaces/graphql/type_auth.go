package interfaces_graphql

import "github.com/graphql-go/graphql"

// LoginPayload型
var loginPayload = graphql.NewObject(graphql.ObjectConfig{
	Name: "LoginPayload",
	Fields: graphql.Fields{
		"token": &graphql.Field{Type: graphql.String},
	},
})
