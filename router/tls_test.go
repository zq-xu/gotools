package router

import (
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

// go test -v \
// config.go \
// tls.go tls_test.go \
// -test.run TestGeneratePEM -count=1
func TestGeneratePEM(t *testing.T) {
	Convey("TestGeneratePEM", t, func() {
		Convey("print pem", func() {
			certPem, keyPem, err := GeneratePEM(CertOptions{
				Hosts:        []string{"localhost"},
				Organization: "Organization",
				IsCA:         false,
			})
			So(err, ShouldBeNil)
			fmt.Println(fmt.Sprintf("certPem is: %s", string(certPem)))
			fmt.Println(fmt.Sprintf("keyPem is: %s", string(keyPem)))
		})
	})
}
