// examples/advanced/main.go
package main

import (
    "fmt"
    "log"
    "time"
    "norm/builder"
    "norm/model"
)

// User 用户实体
type User struct {
    ID        int64     `cypher:"id,omitempty"`
    Username  string    `cypher:"username,required,unique"`
    Email     string    `cypher:"email,required,unique"`
    Password  string    `cypher:"password,required"`
    Avatar    string    `cypher:"avatar"`
    Bio       string    `cypher:"bio"`
    Active    bool      `cypher:"active"`
    CreatedAt time.Time `cypher:"created_at"`
    UpdatedAt time.Time `cypher:"updated_at"`
    
    // 关系
    Posts     []Post    `relationship:"AUTHORED,outgoing"`
    Follows   []User    `relationship:"FOLLOWS,outgoing"`
    Followers []User    `relationship:"FOLLOWS,incoming"`
    Likes     []Post    `relationship:"LIKES,outgoing"`
}

// Post 文章实体
type Post struct {
    ID        int64     `cypher:"id,omitempty"`
    Title     string    `cypher:"title,required"`
    Content   string    `cypher:"content,required"`
    Slug      string    `cypher:"slug,unique"`
    Published bool      `cypher:"published"`
    Views     int       `cypher:"views"`
    CreatedAt time.Time `cypher:"created_at"`
    UpdatedAt time.Time `cypher:"updated_at"`
    
    // 关系
    Author   User        `relationship:"AUTHORED,incoming"`
    Tags     []Tag       `relationship:"TAGGED,outgoing"`
    Comments []Comment   `relationship:"HAS_COMMENT,outgoing"`
    Likes    []User      `relationship:"LIKES,incoming"`
}

// Tag 标签实体
type Tag struct {
    ID          int64  `cypher:"id,omitempty"`
    Name        string `cypher:"name,required,unique"`
    Description string `cypher:"description"`
    Color       string `cypher:"color"`
    
    // 关系
    Posts []Post `relationship:"TAGGED,incoming"`
}

// Comment 评论实体
type Comment struct {
    ID        int64     `cypher:"id,omitempty"`
    Content   string    `cypher:"content,required"`
    CreatedAt time.Time `cypher:"created_at"`
    UpdatedAt time.Time `cypher:"updated_at"`
    
    // 关系
    Author User `relationship:"COMMENTED,incoming"`
    Post   Post `relationship:"HAS_COMMENT,incoming"`
}

func main() {
    fmt.Println("=== Cypher ORM 高级示例 ===\n")
    
    // 创建实体注册表
    registry := model.NewEntityRegistry()
    
    // 注册所有实体
    entities := []interface{}{
        User{}, Post{}, Tag{}, Comment{},
    }
    
    for _, entity := range entities {
        err := registry.Register(entity)
        if err != nil {
            log.Fatal("注册实体失败:", err)
        }
    }
    
    fmt.Println("✅ 所有实体注册成功")
    
    // 高级示例 1：创建完整的用户和文章数据
    fmt.Println("\n--- 高级示例 1：创建完整数据 ---")
    
    // 创建用户
    user := User{
        Username:  "johndoe",
        Email:     "john@example.com",
        Password:  "hashed_password",
        Avatar:    "/avatars/john.jpg",
        Bio:       "软件开发工程师，热爱技术分享",
        Active:    true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    
    qb := builder.NewQueryBuilder(registry)
    result, err := qb.CreateEntity(user).Return("u").Build()
    if err != nil {
        log.Fatal("创建用户查询失败:", err)
    }
    
    fmt.Println("创建用户:")
    fmt.Println(result.Query)
    
    // 创建标签
    tags := []Tag{
        {Name: "Go语言", Description: "Go编程语言相关", Color: "#00ADD8"},
        {Name: "数据库", Description: "数据库技术", Color: "#336791"},
        {Name: "微服务", Description: "微服务架构", Color: "#FF6B6B"},
    }
    
    for _, tag := range tags {
        result, err = builder.NewQueryBuilder(registry).
            CreateEntity(tag).
            Build()
        if err != nil {
            log.Fatal("创建标签查询失败:", err)
        }
        fmt.Printf("创建标签 %s\n", tag.Name)
    }
    
    // 高级示例 2：复杂的社交网络查询
    fmt.Println("\n--- 高级示例 2：社交网络查询 ---")
    
    // 查找用户的朋友的朋友（二度连接）
    result, err = builder.NewQueryBuilder(registry).
        Match("(user:User {username: $username})").
        Match("(user)-[:FOLLOWS]->(friend)-[:FOLLOWS]->(fof)").
        Where("fof <> user AND NOT (user)-[:FOLLOWS]->(fof)").
        Return("DISTINCT fof.username as suggested_user, fof.bio").
        OrderBy("suggested_user").
        Limit(10).
        SetParameter("username", "johndoe").
        Build()
    
    if err != nil {
        log.Fatal("构建社交查询失败:", err)
    }
    
    fmt.Println("推荐关注用户:")
    fmt.Println(result.Query)
    
    // 高级示例 3：内容推荐算法
    fmt.Println("\n--- 高级示例 3：内容推荐算法 ---")
    
    // 基于用户兴趣标签推荐文章
    result, err = builder.NewQueryBuilder(registry).
        Match("(user:User {username: $username})").
        Match("(user)-[:LIKES]->(liked_post)-[:TAGGED]->(tag)").
        Match("(tag)<-[:TAGGED]-(recommended_post)").
        Where("recommended_post <> liked_post AND NOT (user)-[:LIKES]->(recommended_post)").
        With("recommended_post, count(tag) as matching_tags").
        Where("matching_tags >= $min_matching_tags").
        Match("(author)-[:AUTHORED]->(recommended_post)").
        Return("recommended_post.title, recommended_post.slug, author.username, matching_tags").
        OrderByDesc("matching_tags", "recommended_post.views").
        Limit(5).
        SetParameter("username", "johndoe").
        SetParameter("min_matching_tags", 2).
        Build()
    
    if err != nil {
        log.Fatal("构建推荐查询失败:", err)
    }
    
    fmt.Println("文章推荐:")
    fmt.Println(result.Query)
    
    // 高级示例 4：统计分析查询
    fmt.Println("\n--- 高级示例 4：统计分析查询 ---")
    
    // 用户活跃度统计
    result, err = builder.NewQueryBuilder(registry).
        Match("(u:User)").
        OptionalMatch("(u)-[:AUTHORED]->(posts)").
        OptionalMatch("(u)-[:LIKES]->(liked)").
        OptionalMatch("(u)-[:FOLLOWS]->(following)").
        OptionalMatch("(followers)-[:FOLLOWS]->(u)").
        With("u, count(posts) as post_count, count(liked) as like_count, count(following) as following_count, count(followers) as follower_count").
        Return("u.username, post_count, like_count, following_count, follower_count, (post_count * 3 + like_count + following_count + follower_count) as activity_score").
        OrderByDesc("activity_score").
        Limit(20).
        Build()
    
    if err != nil {
        log.Fatal("构建统计查询失败:", err)
    }
    
    fmt.Println("用户活跃度统计:")
    fmt.Println(result.Query)
    
    // 高级示例 5：时间序列分析
    fmt.Println("\n--- 高级示例 5：时间序列分析 ---")
    
    // 分析最近7天的文章发布趋势
    result, err = builder.NewQueryBuilder(registry).
        Match("(p:Post)").
        Where("p.created_at >= $start_date AND p.published = true").
        With("date(p.created_at) as publish_date, count(p) as daily_posts").
        Return("publish_date, daily_posts").
        OrderBy("publish_date").
        SetParameter("start_date", time.Now().AddDate(0, 0, -7).Format("2006-01-02")).
        Build()
    
    if err != nil {
        log.Fatal("构建时间分析查询失败:", err)
    }
    
    fmt.Println("文章发布趋势:")
    fmt.Println(result.Query)
    
    // 高级示例 6：路径查询
    fmt.Println("\n--- 高级示例 6：路径查询 ---")
    
    // 查找两个用户之间的最短关注路径
    result, err = builder.NewQueryBuilder(registry).
        Match("p = shortestPath((user1:User {username: $user1})-[:FOLLOWS*1..6]->(user2:User {username: $user2}))").
        Return("length(p) as path_length, [n in nodes(p) | n.username] as path").
        SetParameter("user1", "alice").
        SetParameter("user2", "bob").
        Build()
    
    if err != nil {
        log.Fatal("构建路径查询失败:", err)
    }
    
    fmt.Println("用户关系路径:")
    fmt.Println(result.Query)
    
    // 高级示例 7：图算法应用
    fmt.Println("\n--- 高级示例 7：图算法应用 ---")
    
    // 使用 PageRank 算法找出影响力用户
    result, err = builder.NewQueryBuilder(registry).
        Call("gds.pageRank.stream", "user_network").
        Return("nodeId, score").
        OrderByDesc("score").
        Limit(10).
        Build()
    
    if err != nil {
        log.Fatal("构建图算法查询失败:", err)
    }
    
    fmt.Println("影响力用户排名:")
    fmt.Println(result.Query)
    
    // 高级示例 8：数据更新和维护
    fmt.Println("\n--- 高级示例 8：数据更新和维护 ---")
    
    // 批量更新用户最后活跃时间
    result, err = builder.NewQueryBuilder(registry).
        Match("(u:User)").
        Where("u.active = true").
        Set("u.last_active = $current_time", "u.updated_at = $current_time").
        Return("count(u) as updated_users").
        SetParameter("current_time", time.Now().Format(time.RFC3339)).
        Build()
    
    if err != nil {
        log.Fatal("构建更新查询失败:", err)
    }
    
    fmt.Println("批量更新用户:")
    fmt.Println(result.Query)
    
    // 高级示例 9：条件查询
    fmt.Println("\n--- 高级示例 9：条件查询 ---")
    
    // 使用 CASE 语句进行条件分类
    result, err = builder.NewQueryBuilder(registry).
        Match("(u:User)").
        OptionalMatch("(u)-[:AUTHORED]->(posts)").
        With("u, count(posts) as post_count").
        Return(`u.username, 
               post_count,
               CASE 
                 WHEN post_count >= 50 THEN 'Prolific Writer'
                 WHEN post_count >= 10 THEN 'Active Writer'
                 WHEN post_count >= 1 THEN 'Casual Writer'
                 ELSE 'Reader'
               END as user_type`).
        OrderByDesc("post_count").
        Build()
    
    if err != nil {
        log.Fatal("构建条件查询失败:", err)
    }
    
    fmt.Println("用户分类:")
    fmt.Println(result.Query)
    
    // 高级示例 10：子查询
    fmt.Println("\n--- 高级示例 10：子查询 ---")
    
    // 查找有相同兴趣标签的用户群组
    result, err = builder.NewQueryBuilder(registry).
        Match("(u:User)-[:LIKES]->(p:Post)-[:TAGGED]->(t:Tag)").
        With("t, collect(DISTINCT u) as users").
        Where("size(users) >= $min_users").
        Unwind("users", "user").
        With("t, user, size(users) as group_size").
        Return("t.name as tag, user.username, group_size").
        OrderBy("tag", "user.username").
        SetParameter("min_users", 3).
        Build()
    
    if err != nil {
        log.Fatal("构建子查询失败:", err)
    }
    
    fmt.Println("兴趣群组:")
    fmt.Println(result.Query)
    
    fmt.Println("\n✅ 所有高级示例执行完成！")
    fmt.Println("\n这些示例展示了 Cypher ORM 的强大功能：")
    fmt.Println("- 复杂的社交网络查询")
    fmt.Println("- 推荐算法实现")
    fmt.Println("- 统计分析和聚合")
    fmt.Println("- 时间序列分析")
    fmt.Println("- 图路径查询")
    fmt.Println("- 图算法应用")
    fmt.Println("- 数据维护操作")
    fmt.Println("- 条件查询和分类")
    fmt.Println("- 子查询和集合操作")
}