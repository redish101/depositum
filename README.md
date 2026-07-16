# depositum

**D**istributed **E**lectronic **P**reservation and **O**rganization **S**ystem with **I**ntegrated **T**racking and **U**nified **M**anagement

一个档案/图书管理系统。

## API

depositum 提供 HTTP Restful API 与 GraphQL。

### Restful

depositum 提供了比较规范的 Restful API，启动在`/api/v1`，例如：

```http
POST /api/v1/objects
GET /api/v1/objects/1
PATCH /api/v1/objects/1
DELETE /api/v1/objects/1
```

并计划（comming s∞n）提供 OpenAPI 文档。

### GraphQL

depositum 在提供 Restful API 的同时，在`/graphql`提供 GraphQL 及其演练场。文档的话，GraphQL 演练场本身功能就已足够，不再单独提供
