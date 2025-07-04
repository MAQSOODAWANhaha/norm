// examples/basic/main.go
package main

import (
    "fmt"
    "log"
    "time"
    "norm/builder"
    "norm/model"
)

// Person 用户实体
type Person struct {
    ID        int64     `cypher:"id,omitempty"`
    Name      string    `cypher:"name,required"`
    Age       int       `cypher:"age"`
    Email     string    `cypher:"email,unique"`
    Active    bool      `cypher:"active"`
    CreatedAt time.Time `cypher:"created_at"`
    
    // 关系定义
    Friends []Person `relationship:"FRIEND,outgoing"`
    WorksAt Company  `relationship:"WORKS_AT,outgoing"`
}

// Company 公司实体
type Company struct {
    ID       int64  `cypher:"id,omitempty"`
    Name     string `cypher:"name,required"`
    Industry string `cypher:"industry"`
    
    // 关系
    Employees []Person `relationship:"WORKS_AT,incoming"`
}

// Post 文章实体
type Post struct {
    ID        int64     `cypher:"id,omitempty"`
    Title     string    `cypher:"title,required"`
    Content   string    `cypher:"content"`
    Published bool      `cypher:"published"`
    CreatedAt time.Time `cypher:"created_at"`
    
    // 关系
    Author Person `relationship:"AUTHORED,incoming"`
}

func main() {
    fmt.Println("=== Cypher ORM 示例 ===\n")
    
    // 创建实体注册表
    registry := model.NewEntityRegistry()
    
    // 注册实体
    err := registry.Register(Person{})
    if err != nil {
        log.Fatal("注册 Person 实体失败:", err)
    }
    
    err = registry.Register(Company{})
    if err != nil {
        log.Fatal("注册 Company 实体失败:", err)
    }
    
    err = registry.Register(Post{})
    if err != nil {
        log.Fatal("注册 Post 实体失败:", err)
    }
    
    fmt.Println("✅ 实体注册成功")
    
    // 示例 1：基本节点查询
    fmt.Println("\n--- 示例 1：基本节点查询 ---")
    qb := builder.NewQueryBuilder(registry)
    
    result, err := qb.
        Match("(p:Person {name: $name})").
        Return("p").
        SetParameter("name", "张三").
        Build()
    
    if err != nil {
        log.Fatal("构建查询失败:", err)
    }
    
    fmt.Println("Cypher 查询:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
    
    // 示例 2：使用实体创建查询
    fmt.Println("\n--- 示例 2：使用实体创建查询 ---")
    person := Person{
        Name:      "李四",
        Age:       30,
        Email:     "lisi@example.com",
        Active:    true,
        CreatedAt: time.Now(),
    }
    
    result, err = builder.NewQueryBuilder(registry).
        CreateEntity(person).
        Return("p").
        Build()
    
    if err != nil {
        log.Fatal("构建创建查询失败:", err)
    }
    
    fmt.Println("Cypher 查询:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
    
    // 示例 3：复杂关系查询
    fmt.Println("\n--- 示例 3：复杂关系查询 ---")
    result, err = builder.NewQueryBuilder(registry).
        Match("(p:Person)-[:FRIEND]->(f:Person)").
        Where("p.age > $minAge").
        With("p, collect(f) as friends").
        Match("(p)-[:WORKS_AT]->(c:Company)").
        Return("p.name as person_name, size(friends) as friend_count, c.name as company_name").
        OrderBy("friend_count DESC").
        Limit(10).
        SetParameter("minAge", 25).
        Build()
    
    if err != nil {
        log.Fatal("构建复杂查询失败:", err)
    }
    
    fmt.Println("Cypher 查询:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
    
    // 示例 4：公司和员工关系
    fmt.Println("\n--- 示例 4：公司和员工关系 ---")
    company := Company{
        Name:     "科技有限公司",
        Industry: "软件开发",
    }
    
    employee := Person{
        Name:   "王五",
        Age:    28,
        Email:  "wangwu@company.com",
        Active: true,
    }
    
    // 先创建公司和员工
    result, err = builder.NewQueryBuilder(registry).
        CreateEntity(company).
        CreateEntity(employee).
        Build()
    
    if err != nil {
        log.Fatal("构建创建查询失败:", err)
    }
    
    fmt.Println("创建公司和员工:")
    fmt.Println(result.Query)
    
    // 创建工作关系
    result, err = builder.NewQueryBuilder(registry).
        Match("(p:Person {name: $person_name})").
        Match("(c:Company {name: $company_name})").
        Create("(p)-[:WORKS_AT {start_date: $start_date}]->(c)").
        SetParameter("person_name", "王五").
        SetParameter("company_name", "科技有限公司").
        SetParameter("start_date", time.Now().Format("2006-01-02")).
        Build()
    
    if err != nil {
        log.Fatal("构建关系查询失败:", err)
    }
    
    fmt.Println("\n创建工作关系:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
    
    // 示例 5：使用节点构建器
    fmt.Println("\n--- 示例 5：使用节点构建器 ---")
    nodeBuilder := builder.NewNodeBuilder()
    nodePattern := nodeBuilder.
        Variable("user").
        Labels("Person", "User").
        Property("name", "$userName").
        Property("age", "$userAge").
        Build()
    
    fmt.Println("节点模式:", nodePattern)
    
    // 示例 6：使用关系构建器
    fmt.Println("\n--- 示例 6：使用关系构建器 ---")
    relBuilder := builder.NewRelationshipBuilder()
    relPattern := relBuilder.
        Variable("r").
        Type("KNOWS").
        Direction(builder.DirectionOutgoing).
        Property("since", "$since").
        Length(1, 3).
        Build()
    
    fmt.Println("关系模式:", relPattern)
    
    // 示例 7：组合查询
    fmt.Println("\n--- 示例 7：组合查询 ---")
    result, err = builder.NewQueryBuilder(registry).
        Match("(p1:Person {name: $name1})").
        Match("(p2:Person {name: $name2})").
        Create("(p1)-[:FRIEND {created_at: $created_at}]->(p2)").
        Create("(p2)-[:FRIEND {created_at: $created_at}]->(p1)").
        SetParameter("name1", "张三").
        SetParameter("name2", "李四").
        SetParameter("created_at", time.Now().Format(time.RFC3339)).
        Build()
    
    if err != nil {
        log.Fatal("构建组合查询失败:", err)
    }
    
    fmt.Println("Cypher 查询:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
    
    // 示例 8：聚合查询
    fmt.Println("\n--- 示例 8：聚合查询 ---")
    result, err = builder.NewQueryBuilder(registry).
        Match("(p:Person)-[:WORKS_AT]->(c:Company)").
        With("c, count(p) as employee_count").
        Where("employee_count > $min_employees").
        Return("c.name as company, c.industry, employee_count").
        OrderByDesc("employee_count").
        SetParameter("min_employees", 5).
        Build()
    
    if err != nil {
        log.Fatal("构建聚合查询失败:", err)
    }
    
    fmt.Println("Cypher 查询:")
    fmt.Println(result.Query)
    fmt.Println("参数:")
    for k, v := range result.Parameters {
        fmt.Printf("  %s: %v\n", k, v)
    }
    
    fmt.Println("\n✅ 所有示例执行完成！")
}