# Mutationマニュアル

## URL

以下URLでアクセスすること。
- メソッドはPOST。
- `Header` の `c` に`Bearer JWTトークン`を付与すること

```txt
[オリジン]/graphql
```

## Todo生成

```graphql
mutation ($description: String!, $completed: Boolean!) {
  createTodo(description: $description, completed: $completed) {
    id
    description
    completed
  }
}
```

- graphql variables

```json
{
    "description": "",
    "completed": false,
    "userId": ""
}
```

## Todo更新

```graphql
mutation ($id: String!, $description: String!, $completed: Boolean!) {
  updateTodo(id: $id, description: $description, completed: $completed) {
    id
    description
    completed
    userId
  }
}
```

- graphql variables

```json
{
    "id": "",
    "description": "",
    "completed": false,
    "userId": ""
}
```

## Todo削除

```graphql
mutation ($id: String!) {
  deleteTodo(id: $id)
}
```

- graphql variables

```json
{
    "id": ""
}
```

## ログイン

- `Header` の `Authorization` に`Bearer JWTトークン`を付与は不要。

```graphql
mutation ($email: String!, $password: String!) {
  login(email: $email, password: $password)
}
```

- graphql variables

```json
{
    "email": "",
    "password": ""
}
```