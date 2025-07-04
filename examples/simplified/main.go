// examples/simplified/main.go
package main

import (
	"fmt"
	"log"
	"time"

	"norm/builder"
	"norm/types"
)

// User 用户实体 - 使用新的简化标签格式
type User struct {
	_        struct{} `cypher:"label:User,Person"`     // 指定多个标签
	ID       int64    `cypher:"id,omitempty"`          // 属性映射，空值时忽略
	Username string   `cypher:"username"`              // 简单属性映射
	Email    string   `cypher:"email"`                 
	Age      int      `cypher:"age,omitempty"`         
	Active   bool     `cypher:"active"`                
	Salary   int      `cypher:"salary,omitempty"`      
	CreatedAt time.Time `cypher:"created_at,omitempty"` 
}

// Company 公司实体 - 使用默认标签
type Company struct {
	ID       int64  `cypher:"id,omitempty"`
	Name     string `cypher:"name"`
	Industry string `cypher:"industry,omitempty"`
	Founded  int    `cypher:"founded,omitempty"`
}

// Product 产品实体 - 演示更多标签选项
type Product struct {
	_           struct{} `cypher:"label:Product,Item"`  // 多标签
	ID          int64    `cypher:"id,omitempty"`
	Name        string   `cypher:"name"`
	Description string   `cypher:"description,omitempty"`
	Price       float64  `cypher:"price"`
	InStock     bool     `cypher:"in_stock"`
	Tags        []string `cypher:"tags,omitempty"`      // 列表属性
}

func main() {
	fmt.Println("=== Simplified Cypher ORM Examples ===\n")

	// 无需注册实体！直接使用

	// ================================
	// 1. 基本实体创建
	// ================================
	fmt.Println("--- 1. Basic Entity Creation ---")
	
	user := User{
		Username:  "john_doe",
		Email:     "john@example.com",
		Age:       30,
		Active:    true,
		Salary:    75000,
		CreatedAt: time.Now(),
	}

	qb1 := builder.NewQueryBuilder()
	result1, err := qb1.
		CreateEntity(user).
		Return("u").
		Build()
	
	if err != nil {
		log.Fatal(err)
	}
	printQuery("Create User", result1)

	// ================================
	// 2. 实体匹配查询
	// ================================
	fmt.Println("--- 2. Entity Matching ---")
	
	searchUser := User{
		Username: "john_doe",
		Active:   true,
	}

	qb2 := builder.NewQueryBuilder()
	result2, _ := qb2.
		MatchEntity(searchUser).
		Return("u.username", "u.email", "u.age").
		Build()
	
	printQuery("Match Active User", result2)

	// ================================
	// 3. 公司实体示例
	// ================================
	fmt.Println("--- 3. Company Entity Example ---")
	
	company := Company{
		Name:     "Tech Corp",
		Industry: "Technology",
		Founded:  2010,
	}

	qb3 := builder.NewQueryBuilder()
	result3, _ := qb3.
		MergeEntity(company).
		Return("c").
		Build()
	
	printQuery("Merge Company", result3)

	// ================================
	// 4. 产品实体与复杂属性
	// ================================
	fmt.Println("--- 4. Product with Complex Properties ---")
	
	product := Product{
		Name:        "Laptop Pro",
		Description: "High-performance laptop for professionals",
		Price:       1999.99,
		InStock:     true,
		Tags:        []string{"electronics", "computers", "professional"},
	}

	qb4 := builder.NewQueryBuilder()
	result4, _ := qb4.
		CreateEntity(product).
		Return("p").
		Build()
	
	printQuery("Create Product", result4)

	// ================================
	// 5. 条件查询组合
	// ================================
	fmt.Println("--- 5. Complex Query with Conditions ---")
	
	qb5 := builder.NewQueryBuilder()
	result5, _ := qb5.
		Match("(u:User)").
		Where(builder.AndConditions(
			builder.Gt("u.age", 25),
			builder.Lt("u.salary", 100000),
			builder.Eq("u.active", true),
		)).
		Return(
			"u.username",
			"u.age", 
			"u.salary",
			builder.Round("u.salary / 12").BuildAs("monthly_salary"),
		).
		OrderBy("u.salary DESC").
		Limit(10).
		Build()
	
	printQuery("Complex User Query", result5)

	// ================================
	// 6. 实体解析演示
	// ================================
	fmt.Println("--- 6. Entity Parsing Demo ---")
	
	// 演示直接解析实体信息
	entityInfo, err := builder.ParseEntity(user)
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Printf("User Entity Info:\n")
	fmt.Printf("  Labels: %v\n", entityInfo.Labels)
	fmt.Printf("  Properties: %d items\n", len(entityInfo.Properties))
	for k, v := range entityInfo.Properties {
		fmt.Printf("    %s: %v (%T)\n", k, v, v)
	}

	// ================================
	// 7. 多实体关系查询
	// ================================
	fmt.Println("\n--- 7. Multi-Entity Relationship Query ---")
	
	qb7 := builder.NewQueryBuilder()
	result7, _ := qb7.
		Match("(u:User)").
		Match("(c:Company)").
		Create("(u)-[:WORKS_AT {start_date: $start_date}]->(c)").
		SetParameter("start_date", "2024-01-01").
		Return("u.username", "c.name").
		Build()
	
	printQuery("Create Relationship", result7)

	// ================================
	// 8. 聚合查询
	// ================================
	fmt.Println("--- 8. Aggregation Query ---")
	
	qb8 := builder.NewQueryBuilder()
	result8, _ := qb8.
		Match("(u:User)").
		Where(builder.IsNotNull("u.salary")).
		Return(
			builder.Count("u").BuildAs("total_users"),
			builder.Avg("u.salary").BuildAs("avg_salary"),
			builder.Min("u.age").BuildAs("youngest"),
			builder.Max("u.age").BuildAs("oldest"),
		).
		Build()
	
	printQuery("User Statistics", result8)

	fmt.Println("\n✅ All simplified examples completed!")
	fmt.Println("\n🎯 Key Benefits:")
	fmt.Println("  - No entity registration required")
	fmt.Println("  - Direct struct-to-Cypher conversion")
	fmt.Println("  - Automatic label extraction")
	fmt.Println("  - Type-safe property mapping")
	fmt.Println("  - Flexible tag-based configuration")
}

func printQuery(title string, result types.QueryResult) {
	fmt.Printf("--- %s ---\n", title)
	fmt.Println("Cypher Query:")
	fmt.Println(result.Query)
	if len(result.Parameters) > 0 {
		fmt.Println("Parameters:")
		for k, v := range result.Parameters {
			fmt.Printf("  %s: %v (%T)\n", k, v, v)
		}
	}
	fmt.Printf("Valid: %t\n\n", result.Valid)
}