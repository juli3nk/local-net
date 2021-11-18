package dns

import (
	"fmt"
	"strings"

	"github.com/miekg/dns"
)

func ServerRespond(domain string, dnsServer string) error {
	if dnsServer == "" {
		return fmt.Errorf("dns server is an empty string")
	}
	if !strings.HasSuffix(dnsServer, ":53") {
		dnsServer = fmt.Sprintf("%s:53", dnsServer)
	}

	dnsRequest := new(dns.Msg)
	if dnsRequest == nil {
		return fmt.Errorf("can not new dnsRequest")
	}
	dnsClient := new(dns.Client)
	if dnsClient == nil {
		return fmt.Errorf("can not new dnsClient")
	}

	dnsRequest.SetQuestion(domain + ".", dns.TypeA)

	dnsRequest.SetEdns0(4096, true)

	_, _, err := dnsClient.Exchange(dnsRequest, dnsServer)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
