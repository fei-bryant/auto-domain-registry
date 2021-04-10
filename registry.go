package main

import (
	"context"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"os"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)

}

const (
	DefaultReadLimit = 1 << 20
	WHOIS            = "whois.verisign-grs.com"
)

var (
	errBodyNil     = errors.New("Search Body is nil")
	errAvailableIP = errors.New("Available ip failed")
)

func DomainRegister(domains []string) error {
	return register(domains)
}

func register(domains []string) error {
	ips := getAvailableIPs()
	domainsLength := len(domains)
	ipLength := len(ips)

	if domainsLength > ipLength {
		//TODO
	}
	if err := Search(ips, domains[:ipLength]); err != nil {
		return err
	}

	return nil

}

func Search(ips, domains []string) error {

	for i := 0; i < len(ips); i++ {
		go searchAndRegister(ips[i], domains[i])
	}
	return nil
}

func searchAndRegister(ip, domain string) {
	for {
		err := domainUsed(ip, domain)
		if err != nil && err == errBodyNil {
			log.WithFields(log.Fields{
				"domain": domain,
				"ip":     ip,
			}).Info("域名未注册")
			registerIP(domain)
			continue
		} else if err != nil && err != errBodyNil {
			log.WithFields(log.Fields{
				"domain": domain,
				"ip":     ip,
			}).Error(err.Error())
			continue
		} else {
			log.WithFields(log.Fields{
				"domain": domain,
				"ip":     ip,
			}).Info("域名已经注册")
		}

	}

}

func registerIP(domain string) {

}

func domainUsed(ip, domain string) error {
	d := net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: 0,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", "whois.verisign-grs.com:43")
	if err != nil {
		log.Error(err.Error())
		return err
	}
	defer conn.Close()

	if _, err := conn.Write([]byte(domain + "\r\n")); err != nil {
		log.Error(err.Error())
		return err
	}

	body, err := ioutil.ReadAll(io.LimitReader(conn, DefaultReadLimit))
	if err != nil {
		log.Println(err.Error())
		return err
	}
	if body == nil {
		return errBodyNil
	}
	log.Debug(string(body))
	return nil
}

func getLocalNetInterfaceIP() ([]string, error) {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}

	for _, addr := range addrs {
		if addr.String() == "127.0.0.1" {
			break
		}
		ips = append(ips, strings.Split(addr.String(), "/")[0])
	}
	return ips, nil
}

func getAvailableIPs() (res []string) {
	ips, err := getLocalNetInterfaceIP()
	if err != nil {
		log.Error(err.Error())
		return res
	}

	for _, ip := range ips {
		ip, err := ipAvailable(ip)
		if err != nil {
			log.WithFields(log.Fields{
				"ip": ip,
			}).Info("IP地址不可用")
			continue
		}
		res = append(res, ip)
	}

	return res

}

func ipAvailable(ip string) (string, error) {
	if available(ip) {
		return ip, nil
	}
	return ip, errAvailableIP

}

func available(ip string) bool {
	d := net.Dialer{
		LocalAddr: &net.TCPAddr{
			IP:   net.ParseIP(ip),
			Port: 0,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()
	conn, err := d.DialContext(ctx, "tcp", WHOIS+":43")
	if err != nil {
		return false
	}
	defer conn.Close()
	return true
}
