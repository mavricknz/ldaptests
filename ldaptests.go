package ldaptests

var binddn string = "cn=directory manager"
var passwd string = "qwerty"
var server string = "localhost"
var port uint16 = 1389
var base_dn string = "dc=example,dc=com"

var filters []string = []string{
	"(sn=Abb*)",
	"(uniqueMember=*)",
	"(|(uniqueMember=*)(sn=Abbie))",
	"(&(objectclass=person)(cn=ab*))",
	`(&(objectclass=person)(cn=\41\42*))`, // same as above
	"(&(objectclass=person)(cn=ko*))",
	"(&(|(sn=an*)(sn=ba*))(!(sn=bar*)))",
	"(&(ou:dn:=people)(sn=aa*))",
}

var attributes []string = []string{
	"cn",
	"description",
}
