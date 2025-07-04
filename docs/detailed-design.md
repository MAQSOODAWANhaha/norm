# Cypher ORM 简化详细设计文档

## 目录
1. [核心类型和结构](#核心类型和结构)
2. [查询构建器系统](#查询构建器系统)
3. [模型管理系统](#模型管理系统)
4. [类型转换系统](#类型转换系统)
5. [验证系统](#验证系统)
6. [解析系统](#解析系统)
7. [实现指南](#实现指南)
8. [代码示例](#代码示例)
9. [测试指南](#测试指南)

## 核心类型和结构
// ... (内容保持不变)

## 查询构建器系统
// ... (内容保持不变)

## 简化实体解析系统

### 核心设计原则

我们移除了复杂的实体注册表和预注册机制，改为直接从结构体实例进行实时解析。这大大简化了使用方式：

**新的设计优势:**
- 无需预注册实体类型
- 直接从结构体标签提取信息
- 运行时反射解析
- 更简洁的 API

### 标签格式

**支持的标签格式:**

```go
type User struct {
    _        struct{} `cypher:"label:User,VIP"`     // 指定节点标签
    ID       int64    `cypher:"id,omitempty"`       // 属性映射，空值忽略
    Username string   `cypher:"username"`           // 简单属性映射
    Email    string   `cypher:"email"`              
    Active   bool     `cypher:"active"`             
}
```

**标签选项说明:**
- `label:Label1,Label2` - 指定多个节点标签
- `omitempty` - 零值时忽略该字段
- 第一部分为属性名，省略时使用字段名小写

### 核心实现

#### `builder/entity.go`

```go
// EntityInfo 包含从实体中解析出的信息
type EntityInfo struct {
    Labels     []string                   // 节点标签
    Properties map[string]interface{}     // 属性键值对
}

// ParseEntity 直接从结构体实例解析标签和属性信息
func ParseEntity(entity interface{}) (*EntityInfo, error) {
    // 实时解析，无需预注册
    // 支持指针和值类型
    // 自动类型转换
    // 处理 omitempty 选项
}
```

#### 使用示例

```go
// 1. 定义实体
type User struct {
    _        struct{} `cypher:"label:User,Person"`
    Username string   `cypher:"username"`
    Email    string   `cypher:"email"`
    Active   bool     `cypher:"active"`
}

// 2. 直接使用，无需注册
user := User{Username: "john", Email: "john@example.com", Active: true}

qb := builder.NewQueryBuilder()  // 无需传入注册表
result, _ := qb.
    CreateEntity(user).
    Return("u").
    Build()

// 生成: CREATE (:User:Person{username: $username_1, email: $email_2, active: $active_3})
```

## 类型转换系统
// ... (内容保持不变)

// ... (其余部分保持不变)