/*
 *   Copyright 2021 SAP SE
 *
 *   Licensed under the Apache License, Version 2.0 (the "License");
 *   you may not use this file except in compliance with the License.
 *   You may obtain a copy of the License at
 *
 *       http://www.apache.org/licenses/LICENSE-2.0
 *
 *   Unless required by applicable law or agreed to in writing, software
 *   distributed under the License is distributed on an "AS IS" BASIS,
 *   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *   See the License for the specific language governing permissions and
 *   limitations under the License.
 */

package utils

import (
	"errors"
	"net"
)

func InetNtoa(val uint32) net.IP {
	return net.IPv4(byte(val>>24), byte(val>>16&0xFF), byte(val>>8)&0xFF, byte(val&0xFF))
}

func InetAton(ip net.IP) (uint32, error) {
	ip = ip.To4()
	if ip == nil {
		return 0, errors.New("not an IPv4")
	}
	return uint32(ip[3]) | uint32(ip[2])<<8 | uint32(ip[1])<<16 | uint32(ip[0])<<24, nil
}
