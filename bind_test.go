package ldaptests

import (
	"fmt"
	"ldap"
	"testing"
	"time"
)

func TestLocalBind(t *testing.T) {
	fmt.Printf("TestLocalBind: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	l.NetworkConnectTimeout = 5 * time.Second
	err := l.Connect()
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	defer l.Close()
	err = l.Bind(binddn, passwd)
	if err != nil {
		t.Errorf(err.Error())
		return
	}
	fmt.Printf("TestLocalBind: finished...\n")
}

// Really just a test of setting the timeout.
func TestLocalBindWithTimeout(t *testing.T) {
	fmt.Printf("TestLocalBindWithTimeout: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	l.NetworkConnectTimeout = 5 * time.Second
	l.ReadTimeout = 5 * time.Second
	err := l.Connect()
	if err != nil {
		t.Errorf(err.Error() + "\n")
		return
	}
	defer l.Close()
	err = l.Bind(binddn, passwd)
	if err != nil {
		t.Errorf("Timed out in with a bind timeout of 5 seconds!\n")
		return
	}
	fmt.Printf("TestLocalBindWithTimeout: finished...\n")
}
