package ldaptests

import (
	"fmt"
	"github.com/mavricknz/ldap"
	"testing"
)

func TestLocalSearch(t *testing.T) {
	fmt.Printf("TestLocalSearch: starting...\n")
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

	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[0],
		attributes,
	)

	sr, err := l.Search(search_request)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("TestLocalSearch: %s -> num of entries = %d\n", search_request.Filter, len(sr.Entries))
}

func TestLocalSearchWithPaging(t *testing.T) {
	fmt.Printf("TestLocalSearchWithPaging: starting...\n")
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

	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[0],
		attributes,
	)
	sr, err := l.SearchWithPaging(search_request, 3)

	if err != nil {
		t.Error(err)
		return
	}

	fmt.Printf("TestLocalSearchWithPaging: %s -> num of entries = %d\n", search_request.Filter, len(sr.Entries))
}

func testLocalMultiGoroutineSearch(
	t *testing.T, l *ldap.LDAPConnection,
	results chan *ldap.SearchResult, i int) {
	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[i],
		attributes,
	)
	sr, err := l.Search(search_request)

	if err != nil {
		t.Error(err)
		results <- nil
		return
	}

	results <- sr
}

func TestLocalMultiGoroutineSearch(t *testing.T) {
	fmt.Printf("TestLocalMultiGoroutineSearch: starting...\n")
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

	results := make([]chan *ldap.SearchResult, len(filters))
	for i := range filters {
		results[i] = make(chan *ldap.SearchResult)
		go testLocalMultiGoroutineSearch(t, l, results[i], i)
	}
	for i := range filters {
		sr := <-results[i]
		if sr == nil {
			t.Errorf("Did not receive results from goroutine for %q", filters[i])
		} else {
			fmt.Printf("TestLocalMultiGoroutineSearch(%d): %s -> num of entries = %d\n", i, filters[i], len(sr.Entries))
		}
	}
}

func TestLocalCompare(t *testing.T) {
	fmt.Printf("TestLocalCompare: starting...\n")
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

	fmt.Printf("Comparing: %s : sn=Boy which is True\n", addDNs[0])
	compareReq := ldap.NewCompareRequest(addDNs[0], "sn", "Boy")
	result, cerr := l.Compare(compareReq)
	if cerr != nil {
		t.Error(err)
		return
	}
	if result != true {
		t.Error("Compare Result should have been true")
		return
	}
	fmt.Printf("Compare Result : %v\n", result)

	fmt.Printf("Comparing: %s : sn=BoyIsThisWrong which is False\n", addDNs[0])
	compareReq = ldap.NewCompareRequest(addDNs[0], "sn", "BoyIsThisWrong")
	result, cerr = l.Compare(compareReq)
	if cerr != nil {
		t.Error(cerr)
		return
	}
	if result == true {
		t.Error("Compare Result should have been false")
		return
	}
	fmt.Printf("Compare Result : %v\n", result)

	fmt.Printf("Deleting: %s\n", addDNs[0])
	delRequest := ldap.NewDeleteRequest(addDNs[0])
	err = l.Delete(delRequest)
	if err != nil {
		t.Errorf("Delete : %s : %s\n", addDNs[0], err)
		return
	}
}

func TestLocalControlMatchedValuesRequest(t *testing.T) {
	fmt.Printf("LocalControlMatchedValuesRequest: starting...\n")
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

	fmt.Printf("Modify: %s = {aaa, bbb, ccc}\n", "description")
	mod := ldap.NewMod(ldap.ModAdd, "description", []string{"aaa", "bbb", "ccc", "aabb"})
	modreq := ldap.NewModifyRequest(addDNs[0])
	modreq.AddMod(mod)
	err = l.Modify(modreq)
	if err != nil {
		t.Errorf("Modify: %s : %s\n", addDNs[0], err)
		return
	}

	control := ldap.NewControlMatchedValuesRequest(true, "(description=aaa)")
	fmt.Println(control.String())
	fmt.Printf("Search: (objectclass=*), (description=aaa) via MatchedValuesRequest\n")
	search_request := ldap.NewSimpleSearchRequest(
		addDNs[0],
		ldap.ScopeBaseObject,
		"(objectclass=*)",
		[]string{"description"},
	)
	search_request.AddControl(control)
	//l.Debug = true
	sr, err := l.Search(search_request)
	if err != nil {
		t.Errorf("Search: %s : %s\n", addDNs[0], err)
		return
	}
	//l.Debug = false
	fmt.Println("Search Result:")
	fmt.Print(sr)

	control = ldap.NewControlMatchedValuesRequest(true, "(description=a*)")
	fmt.Println(control.String())
	fmt.Printf("Search: (objectclass=*), (description=a*) via MatchedValuesRequest\n")
	search_request = ldap.NewSimpleSearchRequest(
		addDNs[0],
		ldap.ScopeBaseObject,
		"(objectclass=*)",
		[]string{"description"},
	)
	search_request.AddControl(control)
	//l.Debug = true
	sr, err = l.Search(search_request)
	if err != nil {
		t.Errorf("Search: %s : %s\n", addDNs[0], err)
		return
	}
	//l.Debug = false
	fmt.Println("Search Result:")
	fmt.Print(sr)

	fmt.Printf("Deleting: %s\n", addDNs[0])
	delRequest := ldap.NewDeleteRequest(addDNs[0])
	err = l.Delete(delRequest)

	if err != nil {
		t.Errorf("Delete : %s : %s\n", addDNs[0], err)
		return
	}
}

type counter struct {
	EntryCount          int
	ReferenceCount      int
	AbandonAtEntryCount int
}

func (c *counter) ProcessDiscreteResult(sr *ldap.DiscreteSearchResult, connInfo *ldap.ConnectionInfo) (stopProcessing bool, err error) {
	switch sr.SearchResultType {
	case ldap.SearchResultEntry:
		fmt.Printf("result entry: %s\n", sr.Entry.DN)
		c.EntryCount++
		if c.AbandonAtEntryCount != 0 {
			if c.EntryCount == c.AbandonAtEntryCount {
				fmt.Printf("Abandoning at request: %d\n", connInfo.MessageID)
				err = connInfo.Conn.Abandon(connInfo.MessageID)
				// While we are abandoning the results its not an error in this case.
				return true, nil
			}
		}
	case ldap.SearchResultDone:
		fmt.Println("results done")
	case ldap.SearchResultReference:
		fmt.Println("result referral")
		c.ReferenceCount++
	}
	return false, nil
}

func TestLocalSearchWithHandler(t *testing.T) {
	fmt.Printf("TestLocalSearchWithHandler: starting...\n")

	l := ldap.NewLDAPConnection(server, port)
	err := l.Connect()

	// l.Debug = true
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
	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[0],
		attributes,
	)

	l.Debug = false

	// Blocking
	fmt.Println("Blocking version...")
	resultCounter := new(counter)
	err = l.SearchWithHandler(search_request, resultCounter, nil)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("TestLocalSearchWithHandler: %s entries = %d, Referrals = %d\n",
		search_request.Filter, resultCounter.EntryCount, resultCounter.ReferenceCount)

	// Non-Blocking
	fmt.Println("Non-Blocking version...")
	resultChan := make(chan error)
	resultCounter = new(counter)
	go l.SearchWithHandler(search_request, resultCounter, resultChan)
	fmt.Println("do stuff ...")
	err = <-resultChan
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("TestLocalSearchWithHandler - go routine: %s entries = %d, Referrals = %d\n",
		search_request.Filter, resultCounter.EntryCount, resultCounter.ReferenceCount)

	// TODO blocking + abandon non-trival version.

	// Non-Blocking + Abandoning
	fmt.Println("Non-Blocking + Abandon version...")
	resultChan = make(chan error)
	resultCounter = new(counter)
	resultCounter.AbandonAtEntryCount = 4
	go l.SearchWithHandler(search_request, resultCounter, resultChan)
	err = <-resultChan
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("TestLocalSearchWithHandler - go routine: %s entries = %d, Referrals = %d\n",
		search_request.Filter, resultCounter.EntryCount, resultCounter.ReferenceCount)
}

func TestLocalSearchPagingWithHandler(t *testing.T) {
	fmt.Printf("TestLocalSearchPagingWithHandler: starting...\n")

	l := ldap.NewLDAPConnection(server, port)
	err := l.Connect()

	// l.Debug = true
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
	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[0],
		attributes,
	)

	l.Debug = false
	pagingControl := ldap.NewControlPaging(2)
	search_request.Controls = append(search_request.Controls, pagingControl)

	for {
		sr := new(ldap.SearchResult)
		err = l.SearchWithHandler(search_request, sr, nil)
		if err != nil {
			t.Error(err)
			return
		}
		_, pagingResponsePacket := ldap.FindControl(sr.Controls, ldap.ControlTypePaging)
		if pagingResponsePacket == nil {
			t.Errorf("Expected Paging Control.")
		}
		pagingControl.Cookie = pagingResponsePacket.(*ldap.ControlPaging).Cookie
		ldap.ReplaceControl(search_request.Controls, pagingControl)
		fmt.Printf("TestLocalSearchPagingWithHandler: %s entries = %d, Referrals = %d\n",
			search_request.Filter, len(sr.Entries), len(sr.Referrals))
		if len(pagingControl.Cookie) == 0 {
			return
		}
	}
}

func TestLocalConnAndSearch(t *testing.T) {
	fmt.Printf("TestLocalConnAndSearch: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	err := l.Connect()

	// l.Debug = true
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
	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[0],
		attributes,
	)
	// ber.Debug = true
	sr, err := l.Search(search_request)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("TestLocalSearch: %s -> num of entries = %d\n", search_request.Filter, len(sr.Entries))
}

func TestLocalOrderedSearch(t *testing.T) {
	fmt.Printf("TestLocalOrderedSearch: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	err := l.Connect()

	// l.Debug = true
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
	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		filters[3],
		attributes,
	)

	serverSideSortAttrRuleOrder := ldap.ServerSideSortAttrRuleOrder{
		AttributeName: "cn",
		OrderingRule:  "",
		ReverseOrder:  false,
	}
	sortKeyList := make([]ldap.ServerSideSortAttrRuleOrder, 0, 1)
	sortKeyList = append(sortKeyList, serverSideSortAttrRuleOrder)
	sortControl := ldap.NewControlServerSideSortRequest(sortKeyList, true)
	fmt.Println(sortControl.String())
	search_request.AddControl(sortControl)
	l.Debug = false
	sr, err := l.Search(search_request)
	if err != nil {
		t.Error(err)
		return
	}
	_, sssResponse := ldap.FindControl(sr.Controls, ldap.ControlTypeServerSideSortResponse)
	if sssResponse != nil {
		fmt.Println(sssResponse.String())
	}
	fmt.Printf("TestLocalSearch: %s -> num of entries = %d\n", search_request.Filter, len(sr.Entries))
}

func TestLocalVlvSearch(t *testing.T) {
	fmt.Printf("TestLocalVlvSearch: starting...\n")
	l := ldap.NewLDAPConnection(server, port)
	err := l.Connect()

	// l.Debug = true
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
	search_request := ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		"(cn=*)",
		attributes,
	)
	vlvControl := new(ldap.ControlVlvRequest)
	vlvControl.BeforeCount = 0
	vlvControl.AfterCount = 3

	offset := new(ldap.VlvOffSet)
	offset.Offset = 1
	offset.ContentCount = 3

	vlvControl.ByOffset = offset

	//pack, _ := vlvControl.Encode()
	//fmt.Println(hex.Dump(pack.Bytes()))

	search_request.AddControl(vlvControl)

	serverSideSortAttrRuleOrder := ldap.ServerSideSortAttrRuleOrder{
		AttributeName: "cn",
		OrderingRule:  "",
		ReverseOrder:  false,
	}
	sortKeyList := make([]ldap.ServerSideSortAttrRuleOrder, 0, 1)
	sortKeyList = append(sortKeyList, serverSideSortAttrRuleOrder)
	sortControl := ldap.NewControlServerSideSortRequest(sortKeyList, true)
	search_request.AddControl(sortControl)

	l.Debug = false
	sr, err := l.Search(search_request)
	if err != nil {
		t.Error(err)
		return
	}
	_, vlvResp := ldap.FindControl(sr.Controls, ldap.ControlTypeVlvResponse)
	if vlvResp != nil {
		fmt.Println(vlvResp.String())
	}
	for _, entry := range sr.Entries {
		fmt.Println(entry.GetAttributeValues("cn")[0])
	}
	fmt.Printf("TestLocalVlvSearch (byOffSet): %s -> num of entries = %d\n", search_request.Filter, len(sr.Entries))

	search_request = ldap.NewSimpleSearchRequest(
		base_dn,
		ldap.ScopeWholeSubtree,
		"(cn=*)",
		attributes,
	)

	vlvControl = new(ldap.ControlVlvRequest)
	vlvControl.BeforeCount = 0
	vlvControl.AfterCount = 3
	vlvControl.GreaterThanOrEqual = "Aaren Amar"

	//pack, _ := vlvControl.Encode()
	//fmt.Println(hex.Dump(pack.Bytes()))

	search_request.AddControl(vlvControl)
	search_request.AddControl(sortControl)

	sr, err = l.Search(search_request)
	if err != nil {
		t.Error(err)
		return
	}
	_, vlvResp = ldap.FindControl(sr.Controls, ldap.ControlTypeVlvResponse)
	if vlvResp != nil {
		fmt.Println(vlvResp.String())
	}
	for _, entry := range sr.Entries {
		fmt.Println(entry.GetAttributeValues("cn")[0])
	}
	fmt.Printf("TestLocalVlvSearch (value): %s -> num of entries = %d\n", search_request.Filter, len(sr.Entries))
	fmt.Printf("TestLocalVlvSearch: Finished.\n")
}
