package my_app_test

import (
	"fmt"
	"testing"

	"github.com/gophteam/jenkins_demo/common"
)

// TestSnakeToCaseTD !
func TestSnakeToCaseTD(t *testing.T) {
	var tests = []struct {
		p1   string
		want string
	}{
		{"LastName", "last_name"},
		{"SSSNumber", "sss_number"},
		{"StudentIDNumber", "student_id_number"},
	}

	for _, tt := range tests {
		testName := fmt.Sprintf("expect '%s' converted to '%s'", tt.p1, tt.want)
		t.Run(testName, func(t *testing.T) {
			got := common.ToSnakeCase(tt.p1)
			if got != tt.want {
				t.Errorf("want %s, got %s", tt.want, got)
			}
		})
	}
}
