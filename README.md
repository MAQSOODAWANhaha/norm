# Norm: 一个轻量级、流式的 Go Cypher ORM

Norm 是一个为 Go 设计的轻量级、流式且功能强大的 Cypher 查询构建器，旨在简化与 Neo4j 和其他兼容 Cypher 的数据库的交互。它提供了一个类型安全且直观的 API，用于在不编写原始查询字符串的情况下构建从简单到复杂的 Cypher 查询。

## ✨ 特性

- **流式查询构建器**: 通过链式调用 `Match()`、`Where()`、`Return()` 等方法，逐步构建查询。
- **结构体到 Cypher 的映射**: 使用结构体标签（struct tag）自动将 Go 结构体解析为 Cypher 节点和关系模式。
- **完整的图模式构建**: 支持通过 `PatternBuilder` 显式定义节点、关系、方向和变长路径。
- **高级数据操作**:
    - 使用 `MERGE` 并通过 `OnCreate` 和 `OnMatch` 实现复杂的“存在则更新，否则创建”逻辑。
    - 使用 `UNWIND` 展开列表数据。
    - 使用 `REMOVE` 移除节点属性或标签。
- **子查询与集合操作**: 支持 `CALL { ... }` 嵌入子查询，以及 `UNION` 和 `UNION ALL` 合并结果集。
- **丰富的表达式支持**: 提供广泛的函数库，用于聚合、字符串操作、数学计算、列表处理等。
- **复杂的条件查询**: 使用丰富的谓词函数轻松创建嵌套的 `AND`/`OR` 条件。
- **零依赖**: 使用纯 Go 编写，无任何外部依赖。

## 🚀 快速开始

### 安装

```sh
go get github.com/your-username/norm
```

### 快速示例

使用 `cypher` 标签定义你的实体结构体，然后构建查询。

**场景**: 查找一个用户和他创建的第一篇文章。

```go
package main

import (
	"fmt"
	"github.com/your-username/norm/builder"
	"github.com/your-username/norm/types"
)

// User 代表用户节点
type User struct {
	_    struct{} `cypher:"label:User"`
	Name string   `cypher:"name"`
}

// Post 代表文章节点
type Post struct {
	_     struct{} `cypher:"label:Post"`
	Title string   `cypher:"title"`
}

func main() {
	qb := builder.NewQueryBuilder()
	user := &User{}
	post := &Post{}

	// 构建查询
	pattern := builder.NewPatternBuilder().
		StartNode(builder.Node("u", "User")).
		Relationship(builder.Outgoing("WROTE").Variable("r")).
		EndNode(builder.Node("p", "Post")).
		Build()

	result, err := qb.
		MatchPattern(pattern).
		Where(builder.Eq("u.name", "Alice")).
		Return("u.name", "p.title").
		OrderBy("r.createdAt").
		Limit(1).
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
MATCH (u:User)-[r:WROTE]->(p:Post)
WHERE (u.name = $p1)
RETURN u.name, p.title
ORDER BY r.createdAt
LIMIT 1

参数:
map[p1:Alice]
```

## 核心功能与高级用法

### 1. 关系模式构建 (`PatternBuilder`)

当简单的 `Match(&User{})` 不足以描述复杂的图关系时，你可以使用 `PatternBuilder` 来精确定义模式。

**示例**: 查找一位用户写过的、并且被标记为 "Go" 的所有文章。

```go
pattern := builder.NewPatternBuilder().
    StartNode(builder.Node("u", "User")).
    Relationship(builder.Outgoing("WROTE")).
    EndNode(builder.Node("p", "Post")).
    Relationship(builder.Incoming("HAS_TAG")).
    EndNode(builder.Node("t", "Tag")).
    Build()

result, _ := builder.NewQueryBuilder().
    MatchPattern(pattern).
    Where(
        builder.Eq("u.name", "Bob"),
        builder.Eq("t.name", "Go"),
    ).
    Return("p.title").
    Build()

// 生成的 Cypher:
// MATCH (u:User)-[:WROTE]->(p:Post)<-[:HAS_TAG]-(t:Tag)
// WHERE (u.name = $p1 AND t.name = $p2)
// RETURN p.title
```

### 2. 高级 `MERGE` 用法 (`OnCreate` / `OnMatch`)

`MERGE` 用于确保图中不存在重复数据。你可以使用 `OnCreate` 和 `OnMatch` 来指定当节点是新建的或已存在时应执行的附加操作。

**示例**: 如果用户 "Charlie" 不存在，则创建他并记录创建时间；如果他已存在，则更新他的最后访问时间。

```go
user := &User{Name: "Charlie"}

result, _ := builder.NewQueryBuilder().
    Merge(user).As("u").
    OnCreate(map[string]interface{}{
        "u.createdAt": builder.Timestamp(), // 使用 Cypher 的 timestamp() 函数
    }).
    OnMatch(map[string]interface{}{
        "u.lastSeen": builder.Timestamp(),
    }).
    Return("u").
    Build()

// 生成的 Cypher:
// MERGE (u:User {name: $p1})
// ON CREATE SET u.createdAt = timestamp()
// ON MATCH SET u.lastSeen = timestamp()
// RETURN u
```

### 3. 子查询 (`CALL`)

`CALL { ... }` 允许你在一个查询内部执行一个独立的子查询，这对于聚合或复杂的逻辑非常有用。

**示例**: 查找所有文章及其作者数量。

```go
post := &Post{}
subQuery := builder.NewQueryBuilder().
    With("p"). // 从外部查询接收 'p'
    MatchPattern(
        builder.NewPatternBuilder().
            StartNode(builder.Node("p")).
            Relationship(builder.Incoming("WROTE")).
            EndNode(builder.Node("u", "User")).
            Build(),
    ).
    Return(builder.Count("u").BuildAs("authorCount"))

result, _ := builder.NewQueryBuilder().
    Match(post).As("p").
    Call(subQuery).
    Return("p.title", "authorCount").
    Build()

// 生成的 Cypher:
// MATCH (p:Post)
// CALL {
// WITH p
// MATCH (p)<-[:WROTE]-(u:User)
// RETURN count(u) AS authorCount
// }
// RETURN p.title, authorCount
```

### 4. 集合操作 (`UNION`)

使用 `UNION` 或 `UNION ALL` 来合并来自两个或多个查询的结果。

**示例**: 查找所有标记为 "Go" 或 "Database" 的文章标题。

```go
// 查询 'Go' 标签的文章
query1, _ := builder.NewQueryBuilder().
    MatchPattern(
        builder.NewPatternBuilder().
            StartNode(builder.Node("p", "Post")).
            Relationship(builder.Incoming("HAS_TAG")).
            EndNode(builder.Node("t", "Tag", builder.Eq("name", "Go"))).
            Build(),
    ).
    Return("p.title AS title").
    Build()

// 查询 'Database' 标签的文章
query2, _ := builder.NewQueryBuilder().
    MatchPattern(
        builder.NewPatternBuilder().
            StartNode(builder.Node("p", "Post")).
            Relationship(builder.Incoming("HAS_TAG")).
            EndNode(builder.Node("t", "Tag", builder.Eq("name", "Database"))).
            Build(),
    ).
    Return("p.title AS title").
    Build()

// 合并结果
// 注意：在实际使用中，你需要一个方法来组合这些查询
// 此处仅为演示目的
finalQuery := query1.Query + "\nUNION\n" + query2.Query

// 生成的 Cypher:
// MATCH (p:Post)<-[:HAS_TAG]-(t:Tag {name: 'Go'})
// RETURN p.title AS title
// UNION
// MATCH (p:Post)<-[:HAS_TAG]-(t:Tag {name: 'Database'})
// RETURN p.title AS title
```

## 📖 查询构建器 API

`QueryBuilder` 提供了一个流式接口来构建 Cypher 查询。

| 方法 | 描述 |
|---|---|
| `Match(entity)` | 开始一个 `MATCH` 子句。 |
| `OptionalMatch(entity)` | 开始一个 `OPTIONAL MATCH` 子句。 |
| `Create(entity)` | 开始一个 `CREATE` 子句。 |
| `Merge(entity)` | 开始一个 `MERGE` 子句。 |
| `MatchPattern(pattern)` | 使用 `PatternBuilder` 开始一个 `MATCH` 子句。 |
| `As(alias)` | 为前一个模式设置别名。 |
| `Where(conditions...)` | 添加 `WHERE` 条件。 |
| `Set(properties)` | 添加 `SET` 子句以更新属性。 |
| `OnCreate(properties)` | 在 `MERGE` 创建新节点时执行 `SET`。 |
| `OnMatch(properties)` | 在 `MERGE` 匹配到现有节点时执行 `SET`。 |
| `Remove(items...)` | 添加 `REMOVE` 子句以移除属性或标签。 |
| `Delete(variables...)` | 添加 `DELETE` 子句。 |
| `DetachDelete(variables...)` | 添加 `DETACH DELETE` 子句。 |
| `Return(expressions...)` | 指定返回值。 |
| `With(expressions...)` | 将变量传递给下一个查询部分。 |
| `Unwind(list, alias)` | 展开列表为行。 |
| `Call(subQuery)` | 执行一个子查询。 |
| `Union()` / `UnionAll()` | 合并查询结果。 |
| `OrderBy(fields...)` | 对结果进行排序。 |
| `Skip(count)` | 跳过指定数量的结果。 |
| `Limit(count)` | 限制结果的数量。 |
| `Build()` | 构建最终的查询和参数。 |

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

该项目基于 MIT 许可证。详情请参阅 `LICENSE` 文件。