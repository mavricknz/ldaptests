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

func TestLocalBindWithTimeout(t *testing.T) {
	fmt.Printf("TestLocalBindWithTimeout: starting...\n")
	fmt.Printf("Expecting a %v error\n", ldap.ErrorNetwork)
	l := ldap.NewLDAPConnection(server, port)
	l.NetworkConnectTimeout = 5 * time.Second
	l.ReadTimeout = 5 * time.Microsecond
	err := l.Connect()
	if err != nil {
		t.Errorf(err.Error() + "\n")
		return
	}
	defer l.Close()
	err = l.Bind(binddn, passwd)
	if err != nil {
		if err.ResultCode == ldap.ErrorNetwork {
			fmt.Println(err.Error())
		} else {
			fmt.Println(err.Error())
		}
	} else {
		t.Errorf("Should have timed out, it's possible not.\n")
		return
	}
	fmt.Printf("TestLocalBindWithTimeout: finished...\n")
}
