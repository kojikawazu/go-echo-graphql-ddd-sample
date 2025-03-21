package interfaces_graphql

import "github.com/graphql-go/graphql"

// ユーザー型
var userType = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":       &graphql.Field{Type: graphql.String},
		"username": &graphql.Field{Type: graphql.String},
		"email":    &graphql.Field{Type: graphql.String},
	},
})
