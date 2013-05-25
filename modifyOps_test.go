package ldaptests

import (
	"fmt"
	"ldap"
	"testing"
)

var addDNs []string = []string{"cn=Jon Boy,ou=People,dc=example,dc=com"}
var addAttrs []ldap.EntryAttribute = []ldap.EntryAttribute{
	ldap.EntryAttribute{
		Name: "objectclass",
		Values: []string{
			"person", "inetOrgPerson", "organizationalPerson", "top",
		},
	},
	ldap.EntryAttribute{
		Name: "cn",
		Values: []string{
			"Jon Boy",
		},
	},
	ldap.EntryAttribute{
		Name: "givenName",
		Values: []string{
			"Jon",
		},
	},
	ldap.EntryAttribute{
		Name: "sn",
		Values: []string{
			"Boy",
		},
	},
}

func TestLocalAddAndDelete(t *testing.T) {
	fmt.Printf("TestLocalAddAndDelete: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
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

	addReq := ldap.NewAddRequest(addDNs[0])
	for _, attr := range addAttrs {
		addReq.AddAttribute(&attr)
	}
	fmt.Printf("Adding: %s\n", addDNs[0])
	err = l.Add(addReq)
	if err != nil {
		t.Errorf("Add : %s : %s\n", addDNs[0], err)
		return
	}
	fmt.Printf("Deleting: %s\n", addDNs[0])
	delRequest := ldap.NewDeleteRequest(addDNs[0])
	err = l.Delete(delRequest)
	if err != nil {
		t.Errorf("Delete : %s : %s\n", addDNs[0], err)
		return
	}
}

func TestLocalControlPermissiveModifyRequest(t *testing.T) {
	fmt.Printf("ControlPermissiveModifyRequest: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
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

	addReq := ldap.NewAddRequest(addDNs[0])
	for _, attr := range addAttrs {
		addReq.AddAttribute(&attr)
	}
	fmt.Printf("Adding: %s\n", addDNs[0])
	err = l.Add(addReq)
	if err != nil {
		t.Errorf("Add : %s : %s\n", addDNs[0], err)
		return
	}

	modreq := ldap.NewModifyRequest(addDNs[0])
	mod := ldap.NewMod(ldap.ModAdd, "description", []string{"aaa"})
	modreq.AddMod(mod)
	fmt.Println(modreq)
	err = l.Modify(modreq)
	if err != nil {
		t.Errorf("Modify : %s : %s\n", addDNs[0], err)
		return
	}

	mod = ldap.NewMod(ldap.ModAdd, "description", []string{"aaa", "bbb", "ccc"})
	modreq = ldap.NewModifyRequest(addDNs[0])
	modreq.AddMod(mod)
	control := ldap.NewControlString(ldap.ControlTypePermissiveModifyRequest, true, "")
	fmt.Println(control.String())
	modreq.AddControl(control)
	fmt.Println(modreq)
	err = l.Modify(modreq)
	if err != nil {
		t.Errorf("Modify (Permissive): %s : %s\n", addDNs[0], err)
		return
	}

	mod = ldap.NewMod(ldap.ModAdd, "description", []string{"aaa", "bbb", "ccc", "ddd"})
	modreq = ldap.NewModifyRequest(addDNs[0])
	modreq.AddMod(mod)
	control = ldap.NewControlPermissiveModifyRequest(false)
	fmt.Println(control.String())
	modreq.AddControl(control)
	fmt.Println(modreq)
	err = l.Modify(modreq)
	if err != nil {
		t.Errorf("Modify (Permissive): %s : %s\n", addDNs[0], err)
		return
	}

	fmt.Printf("Deleting: %s\n", addDNs[0])
	delRequest := ldap.NewDeleteRequest(addDNs[0])
	err = l.Delete(delRequest)

	if err != nil {
		t.Errorf("Delete : %s : %s\n", addDNs[0], err)
		return
	}
}
