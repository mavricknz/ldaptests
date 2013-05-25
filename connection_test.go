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
		t.Error(err)
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
		fmt.Print(err) // not an error
		return
	}
	defer l.Close()
	fmt.Print("TestLocalConnectTimeout: finished...\n")
}
