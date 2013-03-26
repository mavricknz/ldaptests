package ldaptests

import (
	"fmt"
	"ldap"
	"testing"
	"time"
)

func TestLocalConnect(t *testing.T) {
	fmt.Printf("TestLocalConnect: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	err := l.Connect()
	if err != nil {
		t.Errorf(err.Error() + "\n")
		return
	}
	defer l.Close()
	fmt.Printf("TestLocalConnect: finished...\n")
}

func TestLocalConnectTimeout(t *testing.T) {
	fmt.Printf("TestLocalConnectTimeout: starting...\n")
	fmt.Printf("Expecting a %v error\n", ldap.ErrorNetwork)
	l := ldap.NewLDAPConnection(server, port)
	l.NetworkConnectTimeout = 1 * time.Microsecond
	err := l.Connect()
	if err != nil {
		if err.ResultCode == ldap.ErrorNetwork {
			fmt.Println(err.Error())
		} else {
			t.Errorf(err.Error() + "\n")
			return
		}
	}
	fmt.Printf("TestLocalConnectTimeout: finished...\n")
}
