// examples/basic/main.go
package main

import (
	"fmt"
	"log"
	"time"

	"norm/builder"
	"norm/model"
	"norm/types"
)

// User 实体，自定义标签
type User struct {
	_		  struct{}  `cypher:"label:Customer,VIP"`
	ID		  int64	  `cypher:"id;omitempty"`
	Username  string	  `cypher:"username;unique"`
	Email	  string	  `cypher:"email;required"`
	Active	  bool	  `cypher:"active"`
	CreatedAt time.Time `cypher:"created_at"`
}

// Company 实体，使用默认标签 "Company"
type Company struct {
	ID		 int64  `cypher:"id;omitempty"`
	Name	 string `cypher:"name;required"`
	Industry string `cypher:"industry"`
}

func main() {
	fmt.Println("=== Cypher ORM Basic Examples with New Tags ===")

	// 1. 创建实体注册表
	registry := model.NewRegistry()

	// 2. 注册实体
	if err := registry.Register(User{}); err != nil {
		log.Fatalf("Failed to register User entity: %v", err)
	}
	if err := registry.Register(Company{}); err != nil {
		log.Fatalf("Failed to register Company entity: %v", err)
	}
	fmt.Println("✅ Entities registered successfully")

	// 检查注册的元数据
	userMeta, _ := registry.Get("User")
	fmt.Printf("User labels: %v\n", userMeta.Labels) // 应该输出 [Customer, VIP]
	companyMeta, _ := registry.Get("Company")
	fmt.Printf("Company labels: %v\n", companyMeta.Labels) // 应该输出 [Company]

	// --- 示例 1: 使用自定义标签创建用户 ---
	fmt.Println("\n--- Example 1: Create User with Custom Labels ---")
	user := User{
		Username:  "test_user",
		Email:	   "test@example.com",
		Active:	   true,
		CreatedAt: time.Now(),
	}
	qb1 := builder.NewQueryBuilder(registry)
	result1, err := qb1.CreateEntity(user).Return("u").Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printQueryResult("Create User", result1)

	// --- 示例 2: 使用默认标签查询公司 ---
	fmt.Println("\n--- Example 2: Query Company with Default Label ---")
	qb2 := builder.NewQueryBuilder(registry)
	result2, err := qb2.
		Match("(c:Company)").
		Where("c.name = $name").
		Return("c").
		SetParameter("name", "Norm Corp").
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printQueryResult("Query Company", result2)

	fmt.Println("\n✅ All basic examples executed!")
}

func printQueryResult(title string, result types.QueryResult) {
	fmt.Printf("\n--- %s ---\n", title)
	fmt.Println("Cypher Query:")
	fmt.Println(result.Query)
	fmt.Println("Parameters:")
	for k, v := range result.Parameters {
		fmt.Printf("  %s: %v\n", k, v)
	}
	fmt.Printf("Valid: %t\n", result.Valid)
	if !result.Valid && len(result.Errors) > 0 {
		fmt.Println("Errors:")
		for _, e := range result.Errors {
			fmt.Printf("  - %s\n", e.Message)
		}
	}
}