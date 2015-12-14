package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/service/ec2"
)

type InstanceInfo struct {
	Name, Profile string
	*ec2.Instance
}

func (i *InstanceInfo) LowerName() string {
	return strings.ToLower(i.Name)
}

func (i *InstanceInfo) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
		i.Name,
		p2s(i.InstanceId),
		p2s(i.PublicIpAddress),
		p2s(i.PrivateIpAddress),
		i.Profile,
		p2s(i.Placement.AvailabilityZone),
		p2s(i.InstanceType),
		p2s(i.State.Name),
	)
}
