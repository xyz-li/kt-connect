package common

import (
	"fmt"
	"github.com/miekg/dns"
	"github.com/rs/zerolog/log"
	"strconv"
)

// SetupDnsServer start dns server on specified port
func SetupDnsServer(dnsHandler dns.Handler, port int, net string) {
	srv := &dns.Server{
		Addr: ":" + strconv.Itoa(port),
		Net: net,
		Handler: dnsHandler,
	}
	err := srv.ListenAndServe()
	if err != nil {
		log.Error().Err(err).Msgf("Failed to start dns server")
		panic(err.Error())
	}
}

// NsLookup query domain record, dnsServerAddr use '<ip>:<port>' format
func NsLookup(domain string, qtype uint16, net, dnsServerAddr string) (*dns.Msg, error) {
	c := new(dns.Client)
	c.Net = net
	msg := new(dns.Msg)
	msg.RecursionDesired = true
	msg.SetQuestion(domain, qtype)
	res, _, err := c.Exchange(msg, dnsServerAddr)
	if err != nil {
		return nil, err
	}
	if res.Rcode == dns.RcodeNameError {
		return nil, DomainNotExistError{domain}
	} else if res.Rcode != dns.RcodeSuccess {
		return nil, fmt.Errorf("response code %d", res.Rcode)
	}
	return res, nil
}
