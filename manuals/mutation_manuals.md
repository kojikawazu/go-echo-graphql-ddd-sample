# Mutationマニュアル

## URL

以下URLでアクセすること。
メソッドはPOST。

```txt
[オリジン]/graphql
```

## Todo生成

```graphql
mutation ($description: String!, $completed: Boolean!, $userId: String!) {
  createTodo(description: $description, completed: $completed, userId: $userId) {
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
mutation ($id: String!, $description: String!, $completed: Boolean!, $userId: String!) {
  updateTodo(id: $id, description: $description, completed: $completed, userId: $userId) {
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