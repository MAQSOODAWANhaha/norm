// examples/advanced/main.go
package main

import (
	"fmt"
	"log"
	"time"

	"norm/builder"
	"norm/types"
)

// User 实体
type User struct {
	ID        int64     `cypher:"id,omitempty"`
	Username  string    `cypher:"username,required,unique"`
	Email     string    `cypher:"email,required,unique"`
	Active    bool      `cypher:"active"`
	CreatedAt time.Time `cypher:"created_at"`
}

// Post 实体
type Post struct {
	ID        int64     `cypher:"id,omitempty"`
	Title     string    `cypher:"title,required"`
	Content   string    `cypher:"content,required"`
	Published bool      `cypher:"published"`
	CreatedAt time.Time `cypher:"created_at"`
}

func main() {
	fmt.Println("=== Cypher ORM Advanced Examples ===")

	// 无需注册实体！直接使用
	fmt.Println("✅ Using simplified entity parsing")

	// --- 示例 1: 使用 WITH 和别名 ---
	fmt.Println("\n--- Example 1: Using WITH and Aliases ---")
	qb1 := builder.NewQueryBuilder()
	res1, err := qb1.
		Match("(p:Post)").
		Where("p.published = true").
		With(builder.As("count(p)", "publishedCount")).
		Return("publishedCount").
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printQueryResult("Count published posts with WITH", res1)

	// --- 示例 2: 在 RETURN 中使用别名 ---
	fmt.Println("\n--- Example 2: Using Aliases in RETURN ---")
	qb2 := builder.NewQueryBuilder()
	res2, err := qb2.
		Match("(u:User)").
		Where("u.active = true").
		Return(builder.As("u.username", "active_user")).
		Limit(5).
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printQueryResult("Find active users with alias", res2)

	// --- 示例 3: 结合 OptionalMatch ---
	fmt.Println("\n--- Example 3: Using OptionalMatch ---")
	qb3 := builder.NewQueryBuilder()
	res3, err := qb3.
		Match("(u:User {username: $username})").
		OptionalMatch("(u)-[:WROTE]->(p:Post)").
		Return("u.username", builder.As("count(p)", "post_count")).
		SetParameter("username", "jane_doe").
		Build()
	if err != nil {
		log.Fatalf("Build failed: %v", err)
	}
	printQueryResult("User post count with OptionalMatch", res3)

	fmt.Println("\n✅ All advanced examples executed!")
}

func printQueryResult(title string, result types.QueryResult) {
	fmt.Printf("--- %s ---\n", title)
	fmt.Println("Cypher Query:")
	fmt.Println(result.Query)
	fmt.Println("Parameters:")
	for k, v := range result.Parameters {
		fmt.Printf("  %s: %v\n", k, v)
	}
	fmt.Printf("Valid: %t\n", result.Valid)
	if !result.Valid {
		fmt.Println("Validation Errors:")
		for _, e := range result.Errors {
			fmt.Printf("  - Type: %s, Message: %s\n", e.Type, e.Message)
		}
	}
	fmt.Println()
}