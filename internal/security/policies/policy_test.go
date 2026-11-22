package policies

import "testing"

func TestPolicyAllows(t *testing.T) {
	allowAll := Policy{Name: "allow-all", Action: "*", Effect: "allow"}
	denySpecific := Policy{Name: "deny", Action: "delete", Effect: "deny"}
	if !allowAll.Allows("ignite") {
		t.Fatal("allow-all policy rejected an action")
	}
	if denySpecific.Allows("delete") {
		t.Fatal("deny policy allowed an action")
	}
}
