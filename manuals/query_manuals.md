# Queryマニュアル

## URL

以下URLでアクセすること。
メソッドはPOST。

```txt
[オリジン]/graphql
```

## ユーザー全取得

- query

```graphql
query {
  users {
    id
    username
    email
  }
}
```

## Todo全取得

- query

```graphql
query {
  todos {
    id
    description
    completed
  }
}
```

## IDによる取得

- query

```graphql
query ($id: String!) {
  todo(id: $id) {
    id
    description
    completed
  }
}
```

- graphql variables

```json
{
  "id": ""
}
```

## ユーザーIDによる取得

```graphql
query ($userId: String!) {
  todoByUserId(userId: $userId) {
    id
    description
    completed
  }
}
```

- graphql variables

```json
{
  "userId": ""
}
```