package service

import "net"

func (sv *MetricService) GetTrustedSubnet() (network *net.IPNet) {
	return sv.trustedSubnet
}
