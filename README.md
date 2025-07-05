# Norm: 一个轻量级、流式的 Go Cypher ORM

Norm 是一个为 Go 设计的轻量级、流式且功能强大的 Cypher 查询构建器，旨在简化与 Neo4j 和其他兼容 Cypher 的数据库的交互。它提供了一个类型安全且直观的 API，用于在不编写原始查询字符串的情况下构建复杂的 Cypher 查询。

## ✨ 特性

- **流式查询构建器**: 通过链式调用 `Match()`、`Where()`、`Return()` 等方法，逐步构建查询。
- **结构体到 Cypher 的映射**: 使用结构体标签（struct tag）自动将 Go 结构体解析为 Cypher 节点模式。
- **丰富的表达式支持**: 提供广泛的函数库，用于聚合、字符串操作、数学计算、列表处理等。
- **复杂的条件查询**: 使用丰富的谓词函数轻松创建嵌套的 `AND`/`OR` 条件。
- **类型安全设计**: 利用 Go 的类型系统在编译时捕获错误。
- **查询验证**: 提供基础的验证功能，在将查询发送到数据库之前捕获常见的语法错误。
- **零依赖**: 使用纯 Go 编写，无任何外部依赖。

## 🚀 快速开始

### 安装

```sh
go get github.com/your-username/norm
```

### 快速示例

使用 `cypher` 标签定义你的实体结构体：

```go
package main

import (
	"fmt"
	"github.com/your-username/norm/builder"
	"github.com/your-username/norm/types"
)

// User 代表图中的一个节点
type User struct {
	_        struct{} `cypher:"label:User,Person"`
	Username string   `cypher:"username"`
	Email    string   `cypher:"email"`
	Active   bool     `cypher:"active"`
	Age      int      `cypher:"age"`
}

func main() {
	qb := builder.NewQueryBuilder()
	user := &User{}

	// 构建一个查询，查找年龄大于 25 岁的活跃用户
	result, err := qb.
		Match(user).As("u").
		Where(
			builder.Gt("u.age", 25),
			builder.Eq("u.active", true),
		).
		Return("u.username", "u.email").
		OrderBy("u.username DESC").
		Limit(10).
		Build()

	if err != nil {
		panic(err)
	}

	// 打印生成的查询和参数
	fmt.Println("生成的查询:")
	fmt.Println(result.Query)
	fmt.Println("\n参数:")
	fmt.Printf("%v\n", result.Parameters)
}
```

### 输出

```
生成的查询:
MATCH (u:User:Person)
WHERE (u.age > $p1 AND u.active = $p2)
RETURN u.username, u.email
ORDER BY u.username DESC
LIMIT 10

参数:
map[p1:25 p2:true]
```

## 📖 文档

### 查询构建器 API

`QueryBuilder` 提供了一个流式接口来构建 Cypher 查询。

| 方法            | 描述                               | 示例                                          |
|-----------------|------------------------------------|-----------------------------------------------|
| `Match()`       | 开始一个 `MATCH` 子句。            | `qb.Match(&User{}).As("u")`                   |
| `OptionalMatch()`| 开始一个 `OPTIONAL MATCH` 子句。   | `qb.OptionalMatch(&Department{}).As("d")`     |
| `Create()`      | 开始一个 `CREATE` 子句。           | `qb.Create(&User{...}).As("u")`               |
| `Merge()`       | 开始一个 `MERGE` 子句。            | `qb.Merge(&User{...}).As("u")`                |
| `Where()`       | 添加带条件的 `WHERE` 子句。        | `qb.Where(builder.Gt("u.age", 18))`           |
| `Set()`         | 添加 `SET` 子句以更新属性。        | `qb.Set("u.active", false)`                   |
| `Remove()`      | 添加 `REMOVE` 子句。               | `qb.Remove("u.property")`                     |
| `Delete()`      | 添加 `DELETE` 子句。               | `qb.Delete("u")`                              |
| `DetachDelete()`| 添加 `DETACH DELETE` 子句。        | `qb.DetachDelete("u")`                        |
| `Return()`      | 指定返回值。                       | `qb.Return("u.name", "u.email")`              |
| `With()`        | 将变量传递给下一个查询部分。       | `qb.With("u")`                                |
| `OrderBy()`     | 对结果进行排序。                   | `qb.OrderBy("u.name DESC")`                   |
| `Skip()`        | 跳过指定数量的结果。               | `qb.Skip(10)`                                 |
| `Limit()`       | 限制结果的数量。                   | `qb.Limit(20)`                                |
| `Build()`       | 构建最终的查询和参数。             | `result, err := qb.Build()`                   |

### 结构体标签 DSL

使用 `cypher` 结构体标签来控制结构体如何映射到 Cypher 节点。

- `label`: 指定节点标签。支持多标签，用逗号分隔。如果省略，则使用结构体名称作为默认标签。
  - `cypher:"label:User,Person"` (多标签示例)
  - `type MyNode struct { _ struct{} }` (将自动生成 `MyNode` 标签)
- `property_name`: 覆盖默认的属性名称（默认为小写的字段名）。
  - `cypher:"username"`
- `omitempty`: 如果字段为零值（例如 `0`, `""`, `false`），则在查询中排除该字段。
  - `cypher:"age,omitempty"`
- `-`: 始终忽略该字段。
  - `cypher:"-"`

### 表达式与条件

Norm 在 `builder` 包中提供了一套丰富的函数来创建复杂的表达式和条件。

#### 条件函数

- `Eq()`, `Neq()`, `Gt()`, `Gte()`, `Lt()`, `Lte()`
- `Contains()`, `StartsWith()`, `EndsWith()`
- `In()`, `IsNull()`, `IsNotNull()`
- `And()`, `Or()`, `Not()` 用于逻辑分组。

**示例:**

```go
qb.Where(
    builder.And(
        builder.Gt("u.age", 18),
        builder.Or(
            builder.Eq("u.department", "Sales"),
            builder.Contains("u.email", "@example.com")
        )
    )
)
```

#### 函数表达式

Norm 支持广泛的 Cypher 函数：

- **聚合**: `Count()`, `Sum()`, `Avg()`, `Min()`, `Max()`, `Collect()`
- **字符串**: `Upper()`, `Lower()`, `Substring()`, `Replace()`
- **数学**: `Abs()`, `Round()`, `Sqrt()`, `Sin()`, `Cos()`
- **列表**: `Size()`, `Labels()`, `Keys()`, `Range()`
- **路径**: `ShortestPath()`, `Nodes()`, `Relationships()`

**示例:**

```go
qb.Return(
    builder.Count("u").BuildAs("total_users"),
    builder.Avg("u.salary").BuildAs("avg_salary")
)
```

## 🏗️ 架构

Norm 采用清晰且模块化的架构设计：

- **`builder/`**: 包含流式查询构建器、表达式辅助函数和实体解析逻辑。
- **`types/`**: 定义核心数据结构，如 `QueryResult` 和 `Condition`。
- **`validator/`**: 为生成的 Cypher 查询提供基础的语法验证。
- **`docs/`**: 包含详细的设计和架构文档。

其核心原理是将一系列 Go 方法调用转换为结构化的 Cypher 子句列表，然后将其编译为带有参数化值的最终查询字符串。

## 🤝 贡献

欢迎参与贡献！如果您发现错误、有功能建议或任何问题，请随时提交 Pull Request 或创建 Issue。

## 📜 许可证

该项目基于 MIT 许可证。详情请参阅 [LICENSE](LICENSE) 文件。