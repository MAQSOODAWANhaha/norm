// examples/expressions/main.go
package main

import (
	"fmt"
	"norm/builder"
	"norm/types"
)

// User 用户实体
type User struct {
	ID       int64  `cypher:"id,omitempty"`
	Username string `cypher:"username,required,unique"`
	Email    string `cypher:"email,required"`
	Age      int    `cypher:"age"`
	Active   bool   `cypher:"active"`
	Salary   int    `cypher:"salary"`
}

// Post 文章实体
type Post struct {
	ID        int64  `cypher:"id,omitempty"`
	Title     string `cypher:"title,required"`
	Content   string `cypher:"content"`
	Views     int    `cypher:"views"`
	Published bool   `cypher:"published"`
}

func main() {
	fmt.Println("=== Cypher ORM Expression Examples ===\n")

	// 无需注册实体！

	// ================================
	// 1. 基本比较操作符示例
	// ================================
	fmt.Println("--- 1. Basic Comparison Operators ---")
	
	qb1 := builder.NewQueryBuilder()
	result1, _ := qb1.
		Match("(u:User)").
		Where(builder.AndConditions(
			builder.Gt("u.age", 25),
			builder.Lt("u.salary", 100000),
			builder.Eq("u.active", true),
		)).
		Return("u.username", "u.age", "u.salary").
		Build()
	printQuery("Active users aged > 25 with salary < 100000", result1)

	// ================================
	// 2. 字符串操作示例
	// ================================
	fmt.Println("\n--- 2. String Operations ---")
	
	qb2 := builder.NewQueryBuilder()
	result2, _ := qb2.
		Match("(u:User)").
		Where(builder.OrConditions(
			builder.Like("u.username", "john"),
			builder.NewExpression().Property("u.email").EndsWith("@gmail.com").Build(),
		)).
		Return(
			builder.Upper("u.username").BuildAs("upper_username"),
			builder.Lower("u.email").BuildAs("lower_email"),
		).
		Build()
	printQuery("Users with 'john' in username OR gmail addresses", result2)

	// ================================
	// 3. 聚合函数示例
	// ================================
	fmt.Println("\n--- 3. Aggregation Functions ---")
	
	qb3 := builder.NewQueryBuilder()
	result3, _ := qb3.
		Match("(u:User)").
		Where(builder.IsNotNull("u.salary")).
		Return(
			builder.Count("u").BuildAs("total_users"),
			builder.Avg("u.salary").BuildAs("avg_salary"),
			builder.Min("u.age").BuildAs("min_age"),
			builder.Max("u.age").BuildAs("max_age"),
			builder.Sum("u.salary").BuildAs("total_salary"),
		).
		Build()
	printQuery("User statistics", result3)

	// ================================
	// 4. 列表和范围操作示例
	// ================================
	fmt.Println("\n--- 4. List and Range Operations ---")
	
	qb4 := builder.NewQueryBuilder()
	result4, _ := qb4.
		Match("(u:User)").
		Where(builder.AndConditions(
			builder.InList("u.age", 25, 30, 35, 40),
			builder.Between("u.salary", 50000, 150000),
		)).
		Return("u.username", "u.age", "u.salary").
		OrderBy("u.age").
		Build()
	printQuery("Users with specific ages and salary range", result4)

	// ================================
	// 5. 数学函数示例
	// ================================
	fmt.Println("\n--- 5. Mathematical Functions ---")
	
	qb5 := builder.NewQueryBuilder()
	result5, _ := qb5.
		Match("(u:User)").
		Where(builder.IsNotNull("u.salary")).
		Return(
			"u.username",
			"u.salary",
			builder.Round("u.salary / 12").BuildAs("monthly_salary"),
			builder.Abs("u.salary - 75000").BuildAs("salary_diff_from_75k"),
		).
		Build()
	printQuery("Salary calculations", result5)

	// ================================
	// 6. CASE 表达式示例
	// ================================
	fmt.Println("\n--- 6. CASE Expressions ---")
	
	salaryCategory := builder.NewCase().
		When("u.salary >= 100000", "'High'").
		When("u.salary >= 50000", "'Medium'").
		Else("'Low'").
		End()
	
	qb6 := builder.NewQueryBuilder()
	result6, _ := qb6.
		Match("(u:User)").
		Where(builder.IsNotNull("u.salary")).
		Return(
			"u.username",
			"u.salary",
			salaryCategory.BuildAs("salary_category"),
		).
		OrderBy("u.salary DESC").
		Build()
	printQuery("User salary categories", result6)

	// ================================
	// 7. 时间函数示例
	// ================================
	fmt.Println("\n--- 7. Temporal Functions ---")
	
	qb7 := builder.NewQueryBuilder()
	result7, _ := qb7.
		Match("(u:User)").
		Return(
			"u.username",
			builder.Date().BuildAs("current_date"),
			builder.DateTime().BuildAs("current_datetime"),
		).
		Limit(3).
		Build()
	printQuery("Current date and time", result7)

	// ================================
	// 8. 复杂表达式组合示例
	// ================================
	fmt.Println("\n--- 8. Complex Expression Combinations ---")
	
	// 创建复杂的WHERE条件
	complexCondition := builder.AndConditions(
		builder.OrConditions(
			builder.Gt("u.age", 30),
			builder.Like("u.username", "admin"),
		),
		builder.Eq("u.active", true),
		builder.Ne("u.email", "''"),
	)
	
	qb8 := builder.NewQueryBuilder()
	result8, _ := qb8.
		Match("(u:User)").
		Where(complexCondition).
		With(
			"u",
			builder.NewCase().
				When("u.age >= 50", "'Senior'").
				When("u.age >= 30", "'Mid-level'").
				Else("'Junior'").
				End().BuildAs("experience_level"),
		).
		Return(
			"u.username",
			"u.age",
			"experience_level",
			builder.Upper("u.email").BuildAs("email_upper"),
		).
		OrderBy("u.age DESC").
		Build()
	printQuery("Complex user categorization", result8)

	// ================================
	// 9. 存在性和谓词函数示例
	// ================================
	fmt.Println("\n--- 9. Existence and Predicate Functions ---")
	
	qb9 := builder.NewQueryBuilder()
	result9, _ := qb9.
		Match("(u:User)").
		OptionalMatch("(u)-[:AUTHORED]->(p:Post)").
		Where(builder.Exists("(u)-[:AUTHORED]->()").Text).
		Return(
			"u.username",
			builder.Count("p").BuildAs("post_count"),
			builder.Collect("p.title").BuildAs("post_titles"),
		).
		Build()
	printQuery("Users who have authored posts", result9)

	// ================================
	// 10. 路径和节点函数示例
	// ================================
	fmt.Println("\n--- 10. Path and Node Functions ---")
	
	qb10 := builder.NewQueryBuilder()
	result10, _ := qb10.
		Match("(u:User)-[r]->(p:Post)").
		Return(
			"u.username",
			"p.title",
			builder.Type("r").BuildAs("relationship_type"),
			builder.Labels("u").BuildAs("user_labels"),
			builder.Keys("u").BuildAs("user_properties"),
		).
		Limit(5).
		Build()
	printQuery("Relationship and node information", result10)

	fmt.Println("\n✅ All expression examples completed!")
}

func printQuery(title string, result types.QueryResult) {
	fmt.Printf("--- %s ---\n", title)
	fmt.Println("Cypher Query:")
	fmt.Println(result.Query)
	if len(result.Parameters) > 0 {
		fmt.Println("Parameters:")
		for k, v := range result.Parameters {
			fmt.Printf("  %s: %v\n", k, v)
		}
	}
	fmt.Printf("Valid: %t\n\n", result.Valid)
}