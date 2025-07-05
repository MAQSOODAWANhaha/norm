# Cypher ORM Expression Builder Features

## 概述

基于官方 Neo4j Cypher 文档，我们为 Cypher ORM 包添加了全面的表达式支持，涵盖了所有主要的 Cypher 函数和操作符。

## 🎯 已实现功能

### 1. 基本比较操作符
- `=`, `<>`, `<`, `<=`, `>`, `>=`
- `CONTAINS`, `STARTS WITH`, `ENDS WITH`
- `IN`, `IS NULL`, `IS NOT NULL`
- 正则表达式匹配 (`=~`)

```go
// 示例
builder.Eq("age", 25)              // u.age = $age_1 (假设别名为 u)
builder.Gt("salary", 50000)        // u.salary > $salary_2
builder.Contains("name", "john")   // u.name CONTAINS $name_3
builder.In("status", 1, 2, 3)      // u.status IN $status_4
```

### 2. 聚合函数 (Aggregating Functions)
- `count()`, `count(DISTINCT)`
- `sum()`, `avg()`, `min()`, `max()`
- `collect()`, `collect(DISTINCT)`

```go
// 示例
builder.Count("u").BuildAs("total_users")
builder.Avg("u.salary").BuildAs("avg_salary")
builder.CollectDistinct("u.department").BuildAs("departments")
```

### 3. 字符串函数 (String Functions)
- `lower()`, `upper()`, `trim()`, `ltrim()`, `rtrim()`
- `replace()`, `substring()`, `split()`
- `toString()`, `left()`, `right()`, `reverse()`

```go
// 示例
builder.Upper("u.name").BuildAs("upper_name")
builder.Replace("u.email", "'@gmail.com'", "'@company.com'")
builder.Substring("u.description", "0", "100")
```

### 4. 数学函数 (Mathematical Functions)
- 基本函数: `abs()`, `ceil()`, `floor()`, `round()`, `sign()`
- 对数函数: `exp()`, `log()`, `log10()`, `sqrt()`
- 三角函数: `sin()`, `cos()`, `tan()`
- 随机数: `rand()`

```go
// 示例
builder.Round("u.salary / 12").BuildAs("monthly_salary")
builder.Abs("u.score - 100").BuildAs("score_diff")
builder.Sqrt("u.area").BuildAs("side_length")
```

### 5. 列表函数 (List Functions)
- `size()`, `head()`, `last()`, `tail()`
- `range()`, `keys()`, `labels()`, `type()`

```go
// 示例
builder.Size("u.skills").BuildAs("skill_count")
builder.Labels("u").BuildAs("user_labels")
builder.Range("1", "10").BuildAs("numbers")
```

### 6. 谓词函数 (Predicate Functions)
- `exists()`, `isEmpty()`
- `all()`, `any()`, `none()`, `single()`

```go
// 示例
builder.Exists("(u)-[:KNOWS]->()")
builder.Any("x", "u.skills", "x = 'programming'")
builder.IsEmpty("u.errors")
```

### 7. 标量函数 (Scalar Functions)
- `coalesce()`, `elementId()`, `id()`, `properties()`
- `startNode()`, `endNode()`

```go
// 示例
builder.Coalesce("u.nickname", "u.username", "'Anonymous'")
builder.ElementId("u").BuildAs("user_id")
builder.Properties("u").BuildAs("user_props")
```

### 8. 时间函数 (Temporal Functions)
- `date()`, `datetime()`, `time()`
- `localtime()`, `localdatetime()`, `duration()`

```go
// 示例
builder.Date().BuildAs("today")
builder.DateTime("'2024-01-01T00:00:00Z'")
builder.Duration("'P1Y2M3DT4H5M6S'")
```

### 9. 路径函数 (Path Functions)
- `length()`, `nodes()`, `relationships()`
- `shortestPath()`, `allShortestPaths()`

```go
// 示例
builder.Length("p").BuildAs("path_length")
builder.Nodes("p").BuildAs("path_nodes")
builder.ShortestPath("(a)-[*]-(b)")
```

### 10. 高级表达式构建器

#### CASE 表达式
```go
caseExpr := builder.NewCase().
    When("u.age >= 65", "'Senior'").
    When("u.age >= 18", "'Adult'").
    Else("'Minor'").
    End().BuildAs("age_category")
```

#### 表达式构建器 (已废弃)
**注意**: `ExpressionBuilder` 是一个较旧的、已废弃的特性。推荐使用新的基于函数的方法（例如 `builder.Eq()`）来构建条件。

```go
// 旧版 ExpressionBuilder
expr := builder.NewExpression().
    Property("u.age").
    GreaterThan(18).
    And("u.active = true").
    Build()
```

#### 逻辑操作符组合
```go
condition := builder.And(
    builder.Gt("u.age", 18),
    builder.Eq("u.active", true),
    builder.Or(
        builder.Contains("u.email", "@company.com"),
        builder.Eq("u.department", "IT"),
    ),
)
```

## 🚀 使用示例

### 复杂查询示例
```go
qb := builder.NewQueryBuilder()
user := &User{} // 假设 User struct 已定义

result, _ := qb.
    Match(user).As("u").
    Where(
        builder.Gt("age", 25),
        builder.Lt("age", 65),
        builder.Contains("email", "@company.com"),
        builder.In("department", "IT", "Engineering", "Data"),
    ).
    With(
        "u",
        builder.NewCase().
            When("u.salary >= 100000", "'Senior'").
            When("u.salary >= 50000", "'Mid'").
            Else("'Junior'").
            End().BuildAs("level"),
    ).
    Return(
        builder.Upper("u.name").BuildAs("name"),
        "u.department",
        "level",
        builder.Round("u.salary / 12").BuildAs("monthly_salary"),
    ).
    OrderBy("level", "u.department").
    Build()
```

### 聚合分析示例
```go
result, _ := qb.
    Match("(u:User)-[:WORKS_AT]->(c:Company)").
    Where(builder.IsNotNull("u.salary")).
    Return(
        "c.name",
        builder.Count("u").BuildAs("employee_count"),
        builder.Avg("u.salary").BuildAs("avg_salary"),
        builder.Min("u.salary").BuildAs("min_salary"),
        builder.Max("u.salary").BuildAs("max_salary"),
        builder.Sum("u.salary").BuildAs("total_payroll"),
    ).
    OrderBy("avg_salary DESC").
    Build()
```

## 📊 支持的操作符和函数统计

- **比较操作符**: 12+ 种
- **聚合函数**: 8+ 种
- **字符串函数**: 12+ 种
- **数学函数**: 15+ 种
- **列表函数**: 8+ 种
- **谓词函数**: 6+ 种
- **标量函数**: 6+ 种
- **时间函数**: 6+ 种
- **路径函数**: 5+ 种

## ✅ 特性

1. **完整的 Cypher 函数支持**: 覆盖官方文档中的所有主要函数
2. **类型安全**: 利用 Go 的类型系统确保正确性
3. **流畅接口**: 支持链式调用和组合
4. **别名支持**: 所有表达式都支持 AS 别名
5. **参数化查询**: 自动处理参数绑定
6. **表达式组合**: 支持复杂的逻辑表达式组合
7. **易用性**: 提供便利函数简化常用操作

## 🎯 覆盖的 Neo4j 功能

✅ 所有基本比较操作符  
✅ 所有聚合函数  
✅ 所有字符串操作函数  
✅ 所有数学计算函数  
✅ 所有列表处理函数  
✅ 所有谓词判断函数  
✅ 所有标量函数  
✅ 所有时间处理函数  
✅ 所有路径分析函数  
✅ CASE 条件表达式  
✅ 复杂表达式组合  

这个实现为 Cypher ORM 提供了企业级的表达式构建能力，完全符合 Neo4j Cypher 语言规范。