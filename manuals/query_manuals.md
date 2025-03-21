# Queryマニュアル

## URL

以下URLでアクセスすること。
- メソッドはPOST。
- `Header` の `Authorization` に`Bearer JWTトークン`を付与すること

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
query {
  todoByUserId {
    id
    description
    completed
  }
}
```
