package main

import (
	"context"
	"log"
	"net"
	"time"

	"github.com/miekg/dns"
)

func main() {

	target := "microsoft.com"
	server := "8.8.8.8"

  // Using Go's net.Resolver
  rv := &net.Resolver{
    PreferGo: true,
    Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
      d := net.Dialer{
        Timeout: time.Second * 10,
      }
      return d.DialContext(ctx, network, server+":53")
    },
  }
  start := time.Now()
  ip, err := rv.LookupHost(context.Background(), target)
  log.Printf("Dialer took %v", time.Since(start))
  if err != nil {
    log.Printf("Dialer ERROR: %+v", err)
  } else {
    log.Printf("DNS Resolved by Dialer: %v", ip)
  }
  log.Println()

  // Use the DNS library
	c := dns.Client{}
	c.Net = ""
	m := dns.Msg{}
  // IPv4 addresses
	m.SetQuestion(target+".", dns.TypeA)
	r, t, err := c.Exchange(&m, server+":53")
	if err != nil {
		log.Printf("IPv4 ERROR: %+v", err)
	} else {
  	log.Printf("IPv4 took %v", t)
    //log.Printf("Reply: %+v", r)
  	if len(r.Answer) == 0 {
  		log.Printf("No IPv4 results")
  	}
  	for _, ans := range r.Answer {
  		Arecord := ans.(*dns.A)
  		log.Printf("%s", Arecord.A)
  	}
  }
  // IPv6 addresses
  m = dns.Msg{}
  m.SetQuestion(target+".", dns.TypeAAAA)
  r, t, err = c.Exchange(&m, server+":53")
  if err != nil {
    log.Printf("IPv6 ERROR: %+v", err)
  } else {
    log.Printf("IPv6 took %v", t)
    //log.Printf("Reply: %+v", r)
    if len(r.Answer) == 0 {
      log.Printf("No IPv6 results")
    }
    for _, ans := range r.Answer {
      Arecord := ans.(*dns.AAAA)
      log.Printf("%s", Arecord.AAAA)
    }
  }
}
