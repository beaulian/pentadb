package pentadb

import (
	"crypto/md5"
	"time"
	"net"
)

func Md5Hash(key []byte) []byte {
	md := md5.New()
	md.Write(key)
	return md.Sum(nil)
}

func KemataHash(digest []byte, i int) uint32 {
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
	conn, err := net.DialTimeout("tcp", ipaddr, timeout)
	if err != nil {
		LOG.Errorf("node %s is unreachable: %s", ipaddr, err)
		return false
	}
	defer conn.Close()
	return true
}
