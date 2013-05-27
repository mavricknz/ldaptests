package ldaptests

import (
	"fmt"
	"github.com/mavricknz/ldap"
	//"runtime/pprof"
	"testing"
	"time"
)

func TestLocalBind(t *testing.T) {
	fmt.Printf("TestLocalBind: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	l.Debug = true
	l.NetworkConnectTimeout = 5 * time.Second
	err := l.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Close()
	err = l.Bind(binddn, passwd)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("TestLocalBind: finished...\n")
}

func TestLocalBindHammer(t *testing.T) {
	fmt.Printf("TestLocalBindHammer: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	l.NetworkConnectTimeout = 5 * time.Second
	// l.Debug = true
	err := l.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Close()
	for i := 0; i < 100; i++ {
		err = l.Bind(binddn, passwd)
		if err != nil {
			t.Error(err)
			return
		}
	}
	fmt.Printf("TestLocalBindHammer: finished...\n")
}

// Really just a test of setting the timeout.
func TestLocalBindWithTimeout(t *testing.T) {
	fmt.Printf("TestLocalBindWithTimeout: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	l.NetworkConnectTimeout = 5 * time.Second
	l.ReadTimeout = 5 * time.Second
	err := l.Connect()
	if err != nil {
		t.Error(err)
		return
	}
	defer l.Close()
	err = l.Bind(binddn, passwd)
	if err != nil {
		t.Error("Timed out in with a bind timeout of 5 seconds!")
		return
	}
	fmt.Printf("TestLocalBindWithTimeout: finished...\n")
}
