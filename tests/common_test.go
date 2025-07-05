// tests/common_test.go
package tests

// User is a common struct used for testing across the package.
type User struct {
	_        struct{} `cypher:"label:User,Person"`
	Username string   `cypher:"username"`
	Email    string   `cypher:"email"`
	Age      int      `cypher:"age,omitempty"`
	Active   bool     `cypher:"active"`
}
