# Cypher ORM 详细设计文档

## 目录
1. [核心类型和结构](#1-核心类型和结构)
2. [查询构建器系统](#2-查询构建器系统)
3. [实体解析系统](#3-实体解析系统)
4. [表达式和条件系统](#4-表达式和条件系统)
5. [验证系统](#5-验证系统)
6. [完整代码示例](#6-完整代码示例)
7. [实现指南和未来工作](#7-实现指南和未来工作)

---

## 1. 核心类型和结构

核心类型定义在 `types/` 目录下，是构建和表示 Cypher 查询的基础。

### `types/core.go`

- **`QueryResult`**: 这是 `Build()` 方法的返回结果，包含了最终生成的 Cypher 查询字符串和参数映射。
  ```go
  type QueryResult struct {
      Query      string                 `json:"query"`
      Parameters map[string]interface{} `json:"parameters"`
      Valid      bool                   `json:"valid"`
      Errors     []ValidationError      `json:"errors"`
  }
  ```

- **`Clause` 和 `ClauseType`**: 查询构建器将一个查询分解为多个子句（`Clause`），例如 `MATCH`、`WHERE`、`RETURN` 等。`ClauseType` 是这些子句类型的枚举。
  ```go
  type ClauseType string
  const (
      MatchClause    ClauseType = "MATCH"
      OptionalClause ClauseType = "OPTIONAL MATCH"
      CreateClause   ClauseType = "CREATE"
      MergeClause    ClauseType = "MERGE"
      WhereClause    ClauseType = "WHERE"
      SetClause      ClauseType = "SET"
      RemoveClause   ClauseType = "REMOVE"
      DeleteClause   ClauseType = "DELETE"
      DetachDeleteClause ClauseType = "DETACH DELETE"
      ReturnClause   ClauseType = "RETURN"
      WithClause     ClauseType = "WITH"
      UnwindClause   ClauseType = "UNWIND"
      OrderByClause  ClauseType = "ORDER BY"
      SkipClause     ClauseType = "SKIP"
      LimitClause    ClauseType = "LIMIT"
      UnionClause    ClauseType = "UNION"
      UnionAllClause ClauseType = "UNION ALL"
      CallClause     ClauseType = "CALL"
      ForEachClause  ClauseType = "FOREACH"
      // ... 其他子句
  )

  type Clause struct {
      Type    ClauseType
      Content string
  }
  ```

### `types/predicate.go`

- **`Condition`**: 这是一个接口，代表 `WHERE` 子句中的一个条件。它可以是一个简单的谓词，也可以是一个逻辑组合。
- **`Predicate`**: 代表一个基本的 "属性-操作符-值" 表达式，例如 `u.name = 'Alice'`。
  ```go
  type Predicate struct {
      Property string
      Operator Operator
      Value    interface{}
      Not      bool
  }
  ```
- **`LogicalGroup`**: 代表由 `AND` 或 `OR` 连接的一组 `Condition`，可以实现复杂的逻辑嵌套。
  ```go
  type LogicalGroup struct {
      Operator   Operator
      Conditions []Condition
  }
  ```

---

## 2. 查询构建器系统

查询构建器是 ORM 的核心，负责通过流畅的链式调用来构建 Cypher 查询。

### `builder/query.go`

- **`QueryBuilder` 接口**: 定义了所有可用的查询构建方法，例如 `Match()`, `Where()`, `Return()` 等。
- **`cypherQueryBuilder` 结构体**: `QueryBuilder` 接口的实现。它内部维护一个 `clauses` 切片，每次调用一个方法（如 `Match`），就会向该切片中添加一个新的 `Clause`。

#### 工作机制

1. **初始化**: `builder.NewQueryBuilder()` 创建一个 `cypherQueryBuilder` 实例。
2. **链式调用**: 每个构建方法（如 `Match`, `Where`）都会修改 `cypherQueryBuilder` 的内部状态（主要是 `clauses` 和 `parameters`），然后返回其自身的指针，从而实现链式调用。
3. **实体处理**: 当 `Match`, `Create` 等方法接收一个实体（struct）时，它会将实体和子句类型暂存到 `pendingEntity` 和 `pendingClause` 中。
4. **别名处理**: `As(alias)` 方法会将别名与暂存的实体关联起来，然后调用 `finalizePendingClause()` 来构建完整的模式字符串（例如 `(u:User:Person)`）并添加到 `clauses` 中。
5. **构建**: `Build()` 方法会遍历 `clauses` 切片，将所有子句按顺序拼接成最终的 Cypher 查询字符串，并返回 `QueryResult`。

#### 示例：
```go
qb := builder.NewQueryBuilder()
// 1. Match(&User{}) -> pendingEntity = &User{}, pendingClause = "MATCH"
// 2. As("u") -> currentAlias = "u", finalizePendingClause() -> clauses = [MATCH (u:User:Person)]
// 3. Where(...) -> clauses = [..., WHERE (u.age > $age_1)]
// 4. Build() -> "MATCH (u:User:Person)\nWHERE (u.age > $age_1)"
```

---

## 3. 实体解析系统

实体解析系统负责将 Go 的 struct 转换为 Cypher 的节点模式。

### `builder/entity.go`

- **`ParseEntity(entity interface{})`**: 这是实体解析的核心函数。它接收一个 struct 实例，通过反射来解析其标签和字段。

#### 标签格式

ORM 支持通过 struct tag 来自定义节点标签和属性。

```go
type User struct {
    // 使用 "label" 标签指定一个或多个节点标签，用逗号分隔。
    // 如果省略此标签，将自动使用结构体名称作为默认标签。
    _        struct{} `cypher:"label:User,Person"`
    
    // 字段名默认为属性名的小写形式
    Username string   `cypher:"username"`
    
    // 使用 "omitempty" 忽略零值字段
    Age      int      `cypher:"age,omitempty"`
    
    // 使用 "-" 标签完全忽略该字段
    Password string   `cypher:"-"`
}

// 示例：没有明确指定标签的结构体将自动获得 "Product" 标签
type Product struct {
    _          struct{}
    ProductID  string `cypher:"product_id"`
    Name       string `cypher:"name"`
}
```

#### 解析流程

1. **`ParseEntity`** 接收一个实体。
2. **`parseLabels`** 被调用来解析标签：
   - 它会查找一个名为 `_` 的匿名空结构体字段。
   - 如果该字段存在并且有 `cypher:"label:..."` 标签，它会提取这些标签，并支持逗号分隔的多个标签。
   - 如果没有找到 `cypher:"label:..."` 标签，或者标签中没有指定有效标签，它将使用 struct 的类型名作为默认标签。
3. **`ParseEntity`** 遍历所有字段：
   - 它会读取 `cypher` 标签来确定属性名。如果标签为空，则使用字段名的小写形式。
   - 如果标签包含 `omitempty`，则在字段为零值时忽略该属性。
   - 最终，所有非零值或未被忽略的字段都会被添加到 `EntityInfo.Properties` 映射中。

---

## 4. 表达式和条件系统

本系统提供了丰富的函数来构建复杂的查询表达式和 `WHERE` 条件。

### `builder/expression.go`

该文件包含大量的辅助函数，用于生成 Cypher 函数和操作符的字符串表示。

- **聚合函数**: `Count()`, `Sum()`, `Avg()` 等。
- **字符串函数**: `Upper()`, `Lower()`, `Substring()` 等。
- **数学函数**: `Abs()`, `Round()`, `Sqrt()` 等。
- **列表/标量/路径函数**: `Size()`, `Labels()`, `ShortestPath()` 等。

所有这些函数都返回一个 `Expression` 结构体，该结构体可以携带一个别名。

```go
// builder.Avg("u.salary") 返回 Expression{Text: "avg(u.salary)"}
// .BuildAs("avg_salary") -> Expression{Text: "avg(u.salary)", Alias: "avg_salary"}
// 在 Build() 中被格式化为 "avg(u.salary) AS avg_salary"
```

### `builder/expression.go` (条件函数)

- **`Eq()`, `Gt()`, `Contains()`** 等函数用于创建 `types.Predicate`。
- **`And()`, `Or()`** 函数用于创建 `types.LogicalGroup`，可以将多个 `Condition` 组合起来。

```go
// 使用示例
builder.Where(
    // 创建一个 Predicate: u.active = true
    builder.Eq("u.active", true),
    // 创建一个 LogicalGroup: (u.age > 18 OR u.department = 'IT')
    builder.Or(
        builder.Gt("u.age", 18),
        builder.Eq("u.department", "IT"),
    ),
)
```

`buildConditionString` 方法会递归地处理这些 `Condition`，生成最终的 `WHERE` 子句字符串，并自动处理参数化。

---

## 5. 验证系统

验证系统确保生成的查询在语法上是基本正确的。

### `validator/query.go`

- **`NewQueryValidator()`**: 创建一个验证器实例。
- **`Validate(query string)`**: 执行验证。

#### 当前的验证能力

- **空查询检查**: 确保查询不为空。
- **括号匹配**: 检查 `()`, `[]`, `{}` 是否正确配对。
- **关键字检查**: 确保查询中至少包含一个有效的 Cypher 关键字（如 `MATCH`, `CREATE`）。

这是一个基础的、非侵入式的验证，不能完全保证 Cypher 的语义正确性，但可以捕捉到许多常见的语法错误。

---

## 6. 完整代码示例

下面是一个结合了多个功能的复杂查询示例，以展示 ORM 的能力。

```go
package main

import (
	"fmt"
	"norm/builder"
	"norm/types"
)

// 定义实体
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

	// 构建查询
	result, err := qb.
		// 匹配 User 节点，别名为 u
		Match(user).As("u").
		// 添加复杂的 WHERE 条件
		Where(
			builder.Gt("u.age", 25),
			builder.In("u.department", "Sales", "Marketing"),
			builder.Or(
				builder.Contains("u.email", "@example.com"),
				builder.Eq("u.active", true),
			),
		).
		// 使用 WITH 子句传递和计算中间结果
		With(
			"u.username AS name",
			builder.NewCase().
				When("u.age >= 60", "'Senior'").
				When("u.age >= 30", "'Mid'").
				Else("'Junior'").
				End().BuildAs("level"),
		).
		// 返回最终结果
		Return("name", "level").
		OrderBy("level", "name DESC").
		Limit(10).
		Build()

	if err != nil {
		panic(err)
	}

	// 打印生成的查询和参数
	fmt.Println("Generated Query:")
	fmt.Println(result.Query)
	fmt.Println("\nParameters:")
	fmt.Printf("%v\n", result.Parameters)
}
```

---

## 7. 实现指南和未来工作

### 实现指南

- **添加新函数**: 要添加对新 Cypher 函数的支持，只需在 `builder/expression.go` 中添加一个新的辅助函数，使其返回一个 `Expression` 即可。
- **添加新子句**: 要支持新的 Cypher 子句，需要在 `types/core.go` 的 `ClauseType` 中添加一个常量，然后在 `builder/query.go` 中为 `QueryBuilder` 接口和 `cypherQueryBuilder` 实现添加一个新方法。

### 未来工作

- **关系构建**: `builder/relationship.go` 尚未实现。未来需要添加对关系模式的完整支持，例如 `(u)-[:KNOWS]->(f)`。
- **更强的验证**: `validator/query.go` 中的结构和参数验证尚未实现。需要一个更强大的验证器来检查子句的顺序是否正确，以及参数是否被正确使用。
- **事务支持**: 当前 ORM 不处理数据库连接和事务。未来可以考虑添加一个可选的执行层来管理这些。
- **CALL { subquery }**: 支持子查询将大大增强 ORM 的能力，允许更复杂的查询组合。

