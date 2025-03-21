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

// DeleteTodoPayload型
var deleteTodoPayload = graphql.NewObject(graphql.ObjectConfig{
	Name: "DeleteTodoPayload",
	Fields: graphql.Fields{
		"success": &graphql.Field{Type: graphql.Boolean},
		"message": &graphql.Field{Type: graphql.String},
	},
})
