package getx

import (
	"net/http"
	"testing"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"

	"github.com/zq-xu/gotools/routerx/auth"
)

func contextWithRoles(roles string) *gin.Context {
	gin.SetMode(gin.TestMode)

	c, _ := gin.CreateTestContext(nil)
	c.Set("JWT_PAYLOAD", jwt.MapClaims{
		auth.AuthAccountIDToken:       "1",
		auth.AuthAccountNameToken:     "Alan",
		auth.AuthAccountUsernameToken: "alan",
		auth.AuthUserRolesToken:       roles,
		auth.AuthAccountStatusToken:   "1",
	})

	return c
}

func TestValidateAuthInRolesAllowsMatchedRole(t *testing.T) {
	ei := ValidateAuthInRoles("admin", "operator")(contextWithRoles("viewer,operator"))
	if ei != nil {
		t.Fatalf("role operator should pass, got %v", ei.ErrorMessage())
	}
}

func TestValidateAuthInRolesDeniesUnmatchedRole(t *testing.T) {
	ei := ValidateAuthInRoles("admin")(contextWithRoles("viewer"))
	if ei == nil {
		t.Fatal("role viewer must not pass a check that requires admin")
	}

	if ei.StatusCode() != http.StatusForbidden {
		t.Fatalf("status = %d, want %d", ei.StatusCode(), http.StatusForbidden)
	}
}

func TestValidateAuthInRolesDeniesEmptyRoles(t *testing.T) {
	ei := ValidateAuthInRoles("admin")(contextWithRoles(""))
	if ei == nil {
		t.Fatal("empty roles must not pass")
	}
}

func TestValidateAuthInRolesDeniesWithoutClaims(t *testing.T) {
	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(nil)

	ei := ValidateAuthInRoles("admin")(c)
	if ei == nil {
		t.Fatal("missing claims must not pass")
	}

	if ei.StatusCode() != http.StatusUnauthorized {
		t.Fatalf("status = %d, want %d", ei.StatusCode(), http.StatusUnauthorized)
	}
}
