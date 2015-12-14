package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/go-ini/ini"
)

var (
	regions = []string{
		"us-east-1",
		"us-west-1",
		"us-west-2",
		"eu-west-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-northeast-1",
		"sa-east-1",
	}
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

type EC2List struct {
	sync.RWMutex
	list []InstanceInfo
}

func (l *EC2List) Add(infos ...InstanceInfo) {
	l.Lock()
	defer l.Unlock()
	l.list = append(l.list, infos...)
}

func (l *EC2List) Sort() {
	sort.Sort(l)
}

func (l *EC2List) Output(w io.Writer) {
	for _, info := range l.list {
		fmt.Fprintf(w, info.String())
	}
}

func (l *EC2List) Len() int {
	return len(l.list)
}

func (l *EC2List) Swap(i, j int) {
	l.list[i], l.list[j] = l.list[j], l.list[i]
}

func (l *EC2List) Less(i, j int) bool {
	return l.list[i].LowerName() < l.list[j].LowerName()
}

//p2s return string value from pointer of string.
//if the pointer is nil, p2s returns empty string.
func p2s(ps *string) string {
	if ps == nil {
		return ""
	}
	return *ps
}

func nameTag(tags []*ec2.Tag) string {
	var name string
	for _, tag := range tags {
		if *tag.Key == "Name" {
			name = *tag.Value
		}
	}
	return name
}

func getCredentialFile() (string, error) {
	var provider credentials.SharedCredentialsProvider
	_, err := provider.Retrieve()
	if err != nil {
		return "", err
	}
	return provider.Filename, nil
}

func getProfileNames(filePath string) ([]string, error) {
	cfg, err := ini.Load(filePath)
	if err != nil {
		return nil, err
	}
	sections := cfg.SectionStrings()
	return sections, nil
}

func validateCredential(creds *credentials.Credentials) bool {
	_, err := creds.Get()
	if err != nil {
		return false
	}
	return true
}

func retrieve(cfg *aws.Config, profile, region string) ([]InstanceInfo, error) {
	svc := ec2.New(session.New(cfg), &aws.Config{Region: aws.String(region)})
	c := make(chan bool, 1)
	var resp *ec2.DescribeInstancesOutput
	var err error
	go func() {
		resp, err = svc.DescribeInstances(nil)
		if err != nil {
			log.Fatalf("profile=%s region=%s err=%v\n", profile, region, err)
		}
		c <- true
	}()

	select {
	case _ = <-c:
	case <-time.After(time.Second * 10):
		log.Fatalf("timeout error. profile=%s region=%s\n", profile, region)
	}

	var infos []InstanceInfo

	for idx := range resp.Reservations {
		for _, inst := range resp.Reservations[idx].Instances {
			name := nameTag(inst.Tags)
			info := InstanceInfo{
				Name:     name,
				Profile:  profile,
				Instance: inst,
			}
			infos = append(infos, info)
		}
	}

	return infos, nil
}

func main() {
	credfile, err := getCredentialFile()
	if err != nil {
		return
	}

	profiles, err := getProfileNames(credfile)
	if err != nil {
		return
	}

	var list EC2List
	var wg sync.WaitGroup

	for _, profile := range profiles {
		profile := profile
		creds := credentials.NewSharedCredentials(credfile, profile)
		if ok := validateCredential(creds); !ok {
			continue
		}

		cfg := aws.NewConfig().WithCredentials(creds)

		for _, region := range regions {
			region := region
			wg.Add(1)
			go func() {
				defer wg.Done()
				infos, err := retrieve(cfg, profile, region)
				if err != nil {
					log.Fatal(err)
				}
				list.Add(infos...)
			}()
		}
	}

	wg.Wait()
	list.Sort()
	list.Output(os.Stdout)
}
