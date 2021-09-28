package core

import (
	"fmt"
	. "getitle/src/utils"
	"math"
	"net"
	"os"
	"strings"
)

func ip2int(ip string) uint {
	s2ip := net.ParseIP(ip).To4()
	return uint(s2ip[3]) | uint(s2ip[2])<<8 | uint(s2ip[1])<<16 | uint(s2ip[0])<<24
}

func int2ip(ipint uint) string {
	ip := make(net.IP, net.IPv4len)
	ip[0] = byte(ipint >> 24)
	ip[1] = byte(ipint >> 16)
	ip[2] = byte(ipint >> 8)
	ip[3] = byte(ipint)
	return ip.String()
}

func mask2ipuint(mask int) uint64 {
	return ((uint64(4294967296) >> uint(32-mask)) - 1) << uint(32-mask)
}

func ip2superip(ip string, mask int) string {
	ipint := ip2int(ip)
	return int2ip(ipint & uint(mask2ipuint(mask)))
}

func getMaskRange(mask int) (before uint, after uint) {

	before = uint(math.Pow(2, 32) - math.Pow(2, float64(32-mask)))
	after = uint(math.Pow(2, float64(32-mask)) - 1)
	return before, after
}

func getIpRange(target string) (start uint, fin uint) {
	_, cidr, _ := net.ParseCIDR(target)
	mask, _ := cidr.Mask.Size()
	before, after := getMaskRange(mask)
	ipint := ip2int(cidr.IP.String())

	start = ipint & before
	fin = ipint | after
	return start, fin
}

func getIp(target string) string {
	if isIPv4(target) {
		return target
	}
	iprecords, err := net.LookupIP(target)
	if err != nil {
		fmt.Println("[-] error IPv4 or bad domain:" + target + ". JUMPED!")
		return ""
	}
	for _, ip := range iprecords {
		if ip.To4() != nil {
			fmt.Println("[*] parse domain SUCCESS, map " + target + " to " + ip.String())
			return ip.String()
		}
	}
	return ""
}

func ipForamt(target string) string {
	target = strings.Replace(target, "http://", "", -1)
	target = strings.Replace(target, "https://", "", -1)
	target = strings.Trim(target, "/")
	if strings.Contains(target, "/") {
		ip := strings.Split(target, "/")[0]
		mask := strings.Split(target, "/")[1]
		target = getIp(ip) + "/" + mask
	} else {
		target = getIp(target) + "/32"
	}
	return target
}

func isIPv4(ip string) bool {
	address := net.ParseIP(ip).To4()
	if address != nil {
		return true
	}
	return false
}

func ipInit(config Config) Config {
	// 如果输入的是文件,则格式化所有输入值.如果无有效ip
	if config.ListFile != "" {
		var iplist []string
		for _, ip := range config.IPlist {
			tmpip := ipForamt(ip)
			if !strings.HasPrefix(tmpip, "/") {
				iplist = append(iplist, tmpip)
			} else {
				fmt.Println("[-] " + tmpip + " ip format error")
			}
		}
		config.IPlist = SliceUnique(iplist) // 去重
		if len(config.IPlist) == 0 {
			fmt.Println("[-] all IP error")
			os.Exit(0)
		}
	} else if config.IP != "" {
		config.IP = ipForamt(config.IP)
		if strings.HasPrefix(config.IP, "/") {
			fmt.Println("[-] IP format error")
			os.Exit(0)
		}
	}
	return config
}
