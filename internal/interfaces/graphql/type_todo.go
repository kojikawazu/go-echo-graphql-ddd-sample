package interfaces_graphql

import "github.com/graphql-go/graphql"

// Todo型
var todoType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Todo",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"completed":   &graphql.Field{Type: graphql.Boolean},
		"userId":      &graphql.Field{Type: graphql.String},
	},
})
