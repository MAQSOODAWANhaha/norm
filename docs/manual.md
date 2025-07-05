我来帮你梳理Neo4j Cypher语句的完整语法手册。首先让我获取官方文档的最新内容。
让我获取更详细的Cypher语法和子句信息：
现在让我获取写入和管理子句的信息：
让我获取更多关于聚合、排序和其他高级功能的信息：
让我再获取一些关于函数和高级特性的信息：
基于搜索结果，我将为您整理一个完整的Neo4j Cypher 5语法功能手册。

# Neo4j Cypher 5 完整语法功能手册

## 1. 概述

Cypher是Neo4j的声明式图查询语言，类似于SQL但专为图数据库设计。本手册涵盖Cypher 5的所有语法特性 [ref:1,2]。

## 2. 基础语法规则

### 2.1 语法特性
- **大小写敏感性**：Cypher关键字不区分大小写，但变量名区分大小写 [ref:18]
- **模式匹配**：使用ASCII艺术风格的模式来描述图结构 [ref:3]
- **子句链接**：查询由多个子句链接组成，类似SQL结构 [ref:8]

### 2.2 注释
```cypher
// 单行注释
/* 多行注释 */
```

## 3. 读取子句 (Reading Clauses)

### 3.1 MATCH子句
**用途**：指定Neo4j在数据库中搜索的模式，是获取数据的主要方式 [ref:12,14]

**基本语法**：
```cypher
MATCH (node:Label)
MATCH (node:Label {property: value})
MATCH (node1)-[:RELATIONSHIP]->(node2)
MATCH (node1)-[:RELATIONSHIP*1..3]->(node2)  // 可变长度路径
```

**高级用法**：
```cypher
// 多个模式匹配
MATCH (n:Person), (m:Movie)
MATCH (n:Person)-[:ACTED_IN]->(m:Movie)

// 多标签匹配
MATCH (n:Person:Actor)

// 使用变量长度关系
MATCH (n:Person)-[:FRIEND*2..4]->(friend)
```

### 3.2 OPTIONAL MATCH子句
**用途**：类似SQL的LEFT JOIN，即使模式不匹配也会返回结果 [ref:11]

**语法**：
```cypher
OPTIONAL MATCH (person:Person)-[:ACTED_IN]->(movie:Movie)
```

**注意事项**：
- WHERE子句是模式描述的一部分，在匹配过程中考虑，而非匹配后过滤 [ref:11]

### 3.3 WHERE子句
**用途**：为MATCH和OPTIONAL MATCH添加模式约束 [ref:13]

**语法**：
```cypher
MATCH (n:Person)
WHERE n.age > 30 AND n.name STARTS WITH 'A'

// 在关系中使用
MATCH (n:Person)-[r:ACTED_IN]->(m:Movie)
WHERE r.year > 2000
```

**支持的操作符**：
- 比较：`=`, `<>`, `<`, `<=`, `>`, `>=`
- 逻辑：`AND`, `OR`, `NOT`
- 字符串：`STARTS WITH`, `ENDS WITH`, `CONTAINS`
- 正则表达式：`=~`
- 列表：`IN`
- 存在性：`IS NULL`, `IS NOT NULL`

### 3.4 UNWIND子句
**用途**：将列表展开为行

**语法**：
```cypher
UNWIND [1, 2, 3] AS x
RETURN x
```

## 4. 写入子句 (Writing Clauses)

### 4.1 CREATE子句
**用途**：创建节点和关系 [ref:9]

**创建节点**：
```cypher
CREATE (n:Person {name: 'John', age: 30})
CREATE (n:Person:Actor {name: 'Tom'})  // 多标签
```

**创建关系**：
```cypher
CREATE (a:Person {name: 'Alice'})-[:KNOWS {since: 2020}]->(b:Person {name: 'Bob'})
```

**创建路径**：
```cypher
CREATE (a:Person {name: 'Alice'})-[:KNOWS]->(b:Person {name: 'Bob'})-[:WORKS_AT]->(c:Company {name: 'Neo4j'})
```

### 4.2 DELETE子句
**用途**：删除节点和关系 [ref:20]

**语法**：
```cypher
// 删除关系
MATCH (n:Person)-[r:KNOWS]->(m:Person)
DELETE r

// 删除节点（必须先删除相关关系）
MATCH (n:Person {name: 'John'})
DETACH DELETE n  // 自动删除相关关系

// 条件删除
MATCH (n:Person)
WHERE n.age < 18
DELETE n
```

### 4.3 SET子句
**用途**：设置属性和标签

**设置属性**：
```cypher
MATCH (n:Person {name: 'John'})
SET n.age = 31, n.city = 'New York'

// 使用表达式
SET n.age = n.age + 1

// 设置所有属性
SET n = {name: 'John', age: 31, city: 'New York'}

// 追加属性
SET n += {phone: '123-456-7890'}
```

**设置标签**：
```cypher
MATCH (n:Person {name: 'John'})
SET n:Actor:Director
```

### 4.4 REMOVE子句
**用途**：移除属性和标签 [ref:23]

**语法**：
```cypher
// 移除属性
MATCH (n:Person {name: 'John'})
REMOVE n.age

// 移除标签
MATCH (n:Person {name: 'John'})
REMOVE n:Actor

// 移除多个
MATCH (n:Person {name: 'John'})
REMOVE n.age, n:Actor
```

### 4.5 MERGE子句
**用途**：匹配现有模式或创建新数据 [ref:19]

**基本语法**：
```cypher
MERGE (n:Person {name: 'John'})
```

**使用ON CREATE和ON MATCH**：
```cypher
MERGE (n:Person {name: 'John'})
ON CREATE SET n.created = timestamp()
ON MATCH SET n.lastSeen = timestamp()
```

**合并关系**：
```cypher
MATCH (a:Person {name: 'Alice'}), (b:Person {name: 'Bob'})
MERGE (a)-[:KNOWS]->(b)
```

## 5. 返回和投影子句

### 5.1 RETURN子句
**用途**：定义查询的输出

**基本语法**：
```cypher
RETURN n
RETURN n.name, n.age
RETURN n.name AS name, n.age AS age
```

**DISTINCT**：
```cypher
RETURN DISTINCT n.name
```

### 5.2 WITH子句
**用途**：将查询部分链接在一起，将一部分的结果作为下一部分的起点 [ref:7]

**语法**：
```cypher
MATCH (person:Person)
WITH person, person.age + 10 AS futureAge
WHERE futureAge > 35
RETURN person.name, futureAge
```

### 5.3 ORDER BY子句
**用途**：对结果进行排序 [ref:24]

**语法**：
```cypher
RETURN n.name, n.age
ORDER BY n.age DESC, n.name ASC

// 使用表达式排序
ORDER BY n.age + n.experience DESC

// 在WITH中使用
WITH n, n.age + 10 AS futureAge
ORDER BY futureAge
```

### 5.4 LIMIT和SKIP子句
**用途**：限制结果数量和跳过记录

**语法**：
```cypher
RETURN n.name
ORDER BY n.age
SKIP 10
LIMIT 20

// 使用变量
RETURN n.name
LIMIT $limit
SKIP $offset
```

## 6. 聚合函数

### 6.1 基本聚合函数 [ref:25]
```cypher
// 计数
RETURN count(*)
RETURN count(n)
RETURN count(DISTINCT n.name)

// 数值聚合
RETURN sum(n.age)
RETURN avg(n.age)
RETURN min(n.age)
RETURN max(n.age)

// 集合聚合
RETURN collect(n.name)
RETURN collect(DISTINCT n.name)
```

### 6.2 聚合规则
- 聚合函数创建分组
- 在ORDER BY中使用聚合必须包含在RETURN中 [ref:25]
- DISTINCT修饰符可用于所有聚合函数 [ref:31]

## 7. 函数库

### 7.1 字符串函数 [ref:34]
```cypher
// 长度和大小
RETURN size('hello')
RETURN length('hello')

// 大小写转换
RETURN toLower('Hello')
RETURN toUpper('Hello')

// 字符串操作
RETURN substring('hello', 1, 3)
RETURN replace('hello world', 'world', 'Neo4j')
RETURN split('a,b,c', ',')
RETURN trim('  hello  ')

// 字符串测试
RETURN startsWith('hello', 'he')
RETURN endsWith('hello', 'lo')
RETURN contains('hello', 'ell')

// 类型转换
RETURN toString(123)
RETURN toStringOrNull(null)
```

### 7.2 列表函数 [ref:33]
```cypher
// 列表操作
RETURN head([1, 2, 3])
RETURN tail([1, 2, 3])
RETURN last([1, 2, 3])
RETURN size([1, 2, 3])

// 列表处理
RETURN reverse([1, 2, 3])
RETURN sort([3, 1, 2])
RETURN range(1, 10)
RETURN range(1, 10, 2)

// 列表谓词
RETURN any(x IN [1, 2, 3] WHERE x > 2)
RETURN all(x IN [1, 2, 3] WHERE x > 0)
RETURN none(x IN [1, 2, 3] WHERE x > 5)
RETURN single(x IN [1, 2, 3] WHERE x = 2)

// 列表过滤和转换
RETURN filter(x IN [1, 2, 3, 4, 5] WHERE x > 3)
RETURN extract(x IN [1, 2, 3] | x * 2)
RETURN reduce(total = 0, x IN [1, 2, 3] | total + x)
```

### 7.3 数学函数 [ref:39]
```cypher
// 基本数学
RETURN abs(-5)
RETURN ceil(3.7)
RETURN floor(3.7)
RETURN round(3.7)
RETURN sign(-5)

// 三角函数
RETURN sin(3.14159)
RETURN cos(3.14159)
RETURN tan(3.14159)
RETURN asin(0.5)
RETURN acos(0.5)
RETURN atan(1)

// 指数和对数
RETURN exp(2)
RETURN log(10)
RETURN log10(100)
RETURN sqrt(16)
RETURN pow(2, 3)

// 随机数
RETURN rand()
RETURN floor(rand() * 100)
```

### 7.4 节点和关系函数
```cypher
// 节点函数
RETURN id(n)
RETURN labels(n)
RETURN keys(n)
RETURN properties(n)

// 关系函数
RETURN id(r)
RETURN type(r)
RETURN startNode(r)
RETURN endNode(r)
RETURN keys(r)
RETURN properties(r)

// 路径函数
RETURN length(path)
RETURN nodes(path)
RETURN relationships(path)
```

### 7.5 日期时间函数
```cypher
// 当前时间
RETURN datetime()
RETURN date()
RETURN time()
RETURN localtime()
RETURN localdatetime()

// 时间戳
RETURN timestamp()

// 日期解析
RETURN date('2023-01-01')
RETURN datetime('2023-01-01T12:00:00')

// 日期操作
RETURN date() + duration('P1D')
RETURN datetime() - duration('PT1H')
```

## 8. 高级特性

### 8.1 条件表达式
```cypher
// CASE表达式
RETURN 
  CASE n.age 
    WHEN 18 THEN 'adult'
    WHEN null THEN 'unknown'
    ELSE 'minor'
  END

// 简单条件
RETURN 
  CASE 
    WHEN n.age >= 18 THEN 'adult'
    ELSE 'minor'
  END
```

### 8.2 路径模式
```cypher
// 变长路径
MATCH (a)-[:KNOWS*1..3]->(b)

// 最短路径
MATCH p = shortestPath((a)-[:KNOWS*]-(b))

// 所有最短路径
MATCH p = allShortestPaths((a)-[:KNOWS*]-(b))
```

### 8.3 子查询
```cypher
// 存在子查询
MATCH (p:Person)
WHERE EXISTS {
  (p)-[:ACTED_IN]->(:Movie)
}

// 计数子查询
MATCH (p:Person)
RETURN p.name, COUNT {
  (p)-[:ACTED_IN]->(:Movie)
} AS movieCount
```

### 8.4 列表推导式
```cypher
// 列表推导
RETURN [x IN range(1, 10) WHERE x % 2 = 0 | x * 2]

// 路径推导
MATCH p = (a)-[:KNOWS*]-(b)
RETURN [n IN nodes(p) | n.name]
```

## 9. 管理子句

### 9.1 数据库管理
```cypher
// 显示数据库
SHOW DATABASES

// 使用数据库
USE database_name

// 创建数据库
CREATE DATABASE database_name

// 删除数据库
DROP DATABASE database_name
```

### 9.2 约束和索引
```cypher
// 创建唯一约束
CREATE CONSTRAINT person_name_unique FOR (p:Person) REQUIRE p.name IS UNIQUE

// 创建索引
CREATE INDEX person_age FOR (p:Person) ON (p.age)

// 显示约束和索引
SHOW CONSTRAINTS
SHOW INDEXES
```

### 9.3 函数和过程
```cypher
// 显示函数
SHOW FUNCTIONS

// 显示过程
SHOW PROCEDURES

// 调用过程
CALL db.labels()
CALL apoc.help('text')
```

## 10. 性能优化

### 10.1 查询计划
```cypher
// 查看执行计划
EXPLAIN MATCH (n:Person) RETURN n

// 查看实际执行计划
PROFILE MATCH (n:Person) RETURN n
```

### 10.2 索引提示
```cypher
// 使用索引提示
MATCH (n:Person)
USING INDEX n:Person(name)
WHERE n.name = 'John'
```

### 10.3 查询优化建议
1. 使用参数化查询
2. 创建适当的索引
3. 避免笛卡尔积
4. 使用LIMIT限制结果
5. 优化WHERE条件的顺序

## 11. 事务和错误处理

### 11.1 事务管理
```cypher
// 显式事务（在驱动程序中）
BEGIN
  CREATE (n:Person {name: 'John'})
  CREATE (m:Person {name: 'Jane'})
  CREATE (n)-[:KNOWS]->(m)
COMMIT

// 回滚事务
ROLLBACK
```

### 11.2 错误处理
```cypher
// 使用FOREACH进行条件操作
FOREACH (ignoreMe IN CASE WHEN condition THEN [1] ELSE [] END |
  CREATE (n:Person {name: 'John'})
)
```

这个功能手册涵盖了Neo4j Cypher 5的所有主要语法特性，每个部分都包含了详细的语法说明和实用示例。通过这个手册，您可以全面掌握Cypher查询语言的使用方法。
