package client

import (
	"log"
	"net"
	"time"
	"crypto/md5"
)

func KemataHash(key string, i int) uint32 {
	md := md5.New()
	md.Write([]byte(key))
	digest := md.Sum(nil)
	// calculate the hash value
	// each four bytes constitute a 32-bit integer
	// then add the four 32-bit integers to the final hash value
	var hash uint32 = 0
	hash += (uint32(digest[(i << 2) + 3] & 0xff) << 24) |
		(uint32(digest[(i << 2) + 2] & 0xff) << 16) |
		(uint32(digest[(i << 2) + 1] & 0xff) << 8) |
		uint32(digest[i << 2] & 0xff)

	return hash
}

func Reachable(ipaddr string, timeout time.Duration) bool {
	_, err := net.DialTimeout("tcp", ipaddr, timeout)
	if err != nil {
		log.Fatalf("node: %s is unreachable due to error: %s", ipaddr, err)
		return false
	}
	return true
}


