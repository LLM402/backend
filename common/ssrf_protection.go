package common

import (
	"fmt"
	"net"
	"net/url"
	"strconv"
	"strings"
)


type SSRFProtection struct {
	AllowPrivateIp         bool
	DomainFilterMode       bool     
	DomainList             []string 
	IpFilterMode           bool     
	IpList                 []string 
	AllowedPorts           []int    
	ApplyIPFilterForDomain bool     
}


var DefaultSSRFProtection = &SSRFProtection{
	AllowPrivateIp:   false,
	DomainFilterMode: true,
	DomainList:       []string{},
	IpFilterMode:     true,
	IpList:           []string{},
	AllowedPorts:     []int{},
}


func isPrivateIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	
	private := []net.IPNet{
		{IP: net.IPv4(10, 0, 0, 0), Mask: net.CIDRMask(8, 32)},     
		{IP: net.IPv4(172, 16, 0, 0), Mask: net.CIDRMask(12, 32)},  
		{IP: net.IPv4(192, 168, 0, 0), Mask: net.CIDRMask(16, 32)}, 
		{IP: net.IPv4(127, 0, 0, 0), Mask: net.CIDRMask(8, 32)},    
		{IP: net.IPv4(169, 254, 0, 0), Mask: net.CIDRMask(16, 32)}, 
		{IP: net.IPv4(224, 0, 0, 0), Mask: net.CIDRMask(4, 32)},    
		{IP: net.IPv4(240, 0, 0, 0), Mask: net.CIDRMask(4, 32)},    
	}

	for _, privateNet := range private {
		if privateNet.Contains(ip) {
			return true
		}
	}

	
	if ip.To4() == nil {
		
		if ip.Equal(net.IPv6loopback) {
			return true
		}
		
		if strings.HasPrefix(ip.String(), "fe80:") {
			return true
		}
		
		if strings.HasPrefix(ip.String(), "fc") || strings.HasPrefix(ip.String(), "fd") {
			return true
		}
	}

	return false
}



func parsePortRanges(portConfigs []string) ([]int, error) {
	var ports []int

	for _, config := range portConfigs {
		config = strings.TrimSpace(config)
		if config == "" {
			continue
		}

		if strings.Contains(config, "-") {
			
			parts := strings.Split(config, "-")
			if len(parts) != 2 {
				return nil, fmt.Errorf("invalid port range format: %s", config)
			}

			startPort, err := strconv.Atoi(strings.TrimSpace(parts[0]))
			if err != nil {
				return nil, fmt.Errorf("invalid start port in range %s: %v", config, err)
			}

			endPort, err := strconv.Atoi(strings.TrimSpace(parts[1]))
			if err != nil {
				return nil, fmt.Errorf("invalid end port in range %s: %v", config, err)
			}

			if startPort > endPort {
				return nil, fmt.Errorf("invalid port range %s: start port cannot be greater than end port", config)
			}

			if startPort < 1 || startPort > 65535 || endPort < 1 || endPort > 65535 {
				return nil, fmt.Errorf("port range %s contains invalid port numbers (must be 1-65535)", config)
			}

			
			for port := startPort; port <= endPort; port++ {
				ports = append(ports, port)
			}
		} else {
			
			port, err := strconv.Atoi(config)
			if err != nil {
				return nil, fmt.Errorf("invalid port number: %s", config)
			}

			if port < 1 || port > 65535 {
				return nil, fmt.Errorf("invalid port number %d (must be 1-65535)", port)
			}

			ports = append(ports, port)
		}
	}

	return ports, nil
}


func (p *SSRFProtection) isAllowedPort(port int) bool {
	if len(p.AllowedPorts) == 0 {
		return true 
	}

	for _, allowedPort := range p.AllowedPorts {
		if port == allowedPort {
			return true
		}
	}
	return false
}


func isDomainListed(domain string, list []string) bool {
	if len(list) == 0 {
		return false
	}

	domain = strings.ToLower(domain)
	for _, item := range list {
		item = strings.ToLower(strings.TrimSpace(item))
		if item == "" {
			continue
		}
		
		if domain == item {
			return true
		}
		
		if strings.HasPrefix(item, "*.") {
			suffix := strings.TrimPrefix(item, "*.")
			if strings.HasSuffix(domain, "."+suffix) || domain == suffix {
				return true
			}
		}
	}
	return false
}

func (p *SSRFProtection) isDomainAllowed(domain string) bool {
	listed := isDomainListed(domain, p.DomainList)
	if p.DomainFilterMode { 
		return listed
	}
	
	return !listed
}



func isIPListed(ip net.IP, list []string) bool {
	if len(list) == 0 {
		return false
	}

	for _, whitelistCIDR := range list {
		_, network, err := net.ParseCIDR(whitelistCIDR)
		if err != nil {
			
			if whitelistIP := net.ParseIP(whitelistCIDR); whitelistIP != nil {
				if ip.Equal(whitelistIP) {
					return true
				}
			}
			continue
		}

		if network.Contains(ip) {
			return true
		}
	}
	return false
}


func (p *SSRFProtection) IsIPAccessAllowed(ip net.IP) bool {
	
	if isPrivateIP(ip) && !p.AllowPrivateIp {
		return false
	}

	listed := isIPListed(ip, p.IpList)
	if p.IpFilterMode { 
		return listed
	}
	
	return !listed
}


func (p *SSRFProtection) ValidateURL(urlStr string) error {
	
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("invalid URL format: %v", err)
	}

	
	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("unsupported protocol: %s (only http/https allowed)", u.Scheme)
	}

	
	host, portStr, err := net.SplitHostPort(u.Host)
	if err != nil {
		
		host = u.Hostname()
		if u.Scheme == "https" {
			portStr = "443"
		} else {
			portStr = "80"
		}
	}

	
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("invalid port: %s", portStr)
	}

	if !p.isAllowedPort(port) {
		return fmt.Errorf("port %d is not allowed", port)
	}

	
	if ip := net.ParseIP(host); ip != nil {
		if !p.IsIPAccessAllowed(ip) {
			if isPrivateIP(ip) {
				return fmt.Errorf("private IP address not allowed: %s", ip.String())
			}
			if p.IpFilterMode {
				return fmt.Errorf("ip not in whitelist: %s", ip.String())
			}
			return fmt.Errorf("ip in blacklist: %s", ip.String())
		}
		return nil
	}

	
	if !p.isDomainAllowed(host) {
		if p.DomainFilterMode {
			return fmt.Errorf("domain not in whitelist: %s", host)
		}
		return fmt.Errorf("domain in blacklist: %s", host)
	}

	
	if !p.ApplyIPFilterForDomain {
		return nil
	}

	
	ips, err := net.LookupIP(host)
	if err != nil {
		return fmt.Errorf("DNS resolution failed for %s: %v", host, err)
	}
	for _, ip := range ips {
		if !p.IsIPAccessAllowed(ip) {
			if isPrivateIP(ip) && !p.AllowPrivateIp {
				return fmt.Errorf("private IP address not allowed: %s resolves to %s", host, ip.String())
			}
			if p.IpFilterMode {
				return fmt.Errorf("ip not in whitelist: %s resolves to %s", host, ip.String())
			}
			return fmt.Errorf("ip in blacklist: %s resolves to %s", host, ip.String())
		}
	}
	return nil
}


func ValidateURLWithFetchSetting(urlStr string, enableSSRFProtection, allowPrivateIp bool, domainFilterMode bool, ipFilterMode bool, domainList, ipList, allowedPorts []string, applyIPFilterForDomain bool) error {
	
	if !enableSSRFProtection {
		return nil
	}

	
	allowedPortInts, err := parsePortRanges(allowedPorts)
	if err != nil {
		return fmt.Errorf("request reject - invalid port configuration: %v", err)
	}

	protection := &SSRFProtection{
		AllowPrivateIp:         allowPrivateIp,
		DomainFilterMode:       domainFilterMode,
		DomainList:             domainList,
		IpFilterMode:           ipFilterMode,
		IpList:                 ipList,
		AllowedPorts:           allowedPortInts,
		ApplyIPFilterForDomain: applyIPFilterForDomain,
	}
	return protection.ValidateURL(urlStr)
}
