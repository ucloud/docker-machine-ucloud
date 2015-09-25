package ucloud

import (
	"encoding/base64"
	"fmt"

	"github.com/ucloud/ucloud-sdk-go/service/uhost"
	"github.com/ucloud/ucloud-sdk-go/service/unet"
	"github.com/ucloud/ucloud-sdk-go/ucloud"
	"github.com/ucloud/ucloud-sdk-go/ucloud/auth"
	"github.com/docker/machine/libmachine/log"
	"github.com/xiaohui/goucloud/ucloud/utils"
)

var (
	hostsvc *uhost.UHost
	unetsvc *unet.UNet
)

func init() {
	config := &ucloud.Config{
		Credentials: &auth.KeyPair{
			PublicKey:  "ucloudsomeone@example.com1296235120854146120",
			PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		},
		Region:    "cn-north-03",
		ProjectID: "",
	}

	hostsvc = uhost.New(config)
	unetsvc = unet.New(config)
}

func createUHost(region, imageId, password string) (string, error) {
	password = base64.StdEncoding.EncodeToString([]byte(password))

	createUhostParams := uhost.CreateUHostInstanceParams{

		Region:    region,
		ImageId:   imageId,
		LoginMode: "Password",
		Password:  password,
		CPU:       1,
		Memory:    2048,
		Quantity:  1,
		Count:     1,
	}

	resp, err := hostsvc.CreateUHostInstance(&createUhostParams)
	utils.DumpVal(resp)
	if err != nil {
		return "", err
	}
	if resp == nil {
		return "", fmt.Errorf("response is empty")
	}

	if resp.RetCode != 0 {
		return "", fmt.Errorf("Create UHost error, RetCode:%d, err:%s", resp.RetCode, err)
	}

	return resp.UHostIds[0], nil
}

func startUHost(region, hostId string) error {

	startUhostParams := uhost.StartUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.StartUHostInstance(&startUhostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func killUHost(region, hostId string) error {
	killUHostParams := uhost.PoweroffUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.PoweroffUHostInstance(&killUHostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func rebootUHost(region, hostId string) error {

	killUHostParams := uhost.PoweroffUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.PoweroffUHostInstance(&killUHostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func terminateUHost(region, hostId string) error {

	terminateUHostParams := uhost.TerminateUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.TerminateUHostInstance(&terminateUHostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Start UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

func stopUHost(region, hostId string) error {
	stopUhostParams := uhost.StopUHostInstanceParams{
		Region:  region,
		UHostId: hostId,
	}

	resp, err := hostsvc.StopUHostInstance(&stopUhostParams)
	if err != nil {
		return err
	}

	if resp.RetCode != 0 {
		return fmt.Errorf("Stop UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	return nil
}

type UHostDetail struct {
	region string
	hostID string

	state     string
	ipAddress string
	cpu       int
	memory    int
}

func getHostDescription(region, hostId string) (*UHostDetail, error) {

	describeParams := uhost.DescribeUHostInstanceParams{
		Region: region,
		UHostIds:[]string{hostId},
		Offset: 0,
		Limit:  10,
	}

	log.Debug(hostsvc)
	resp, err := hostsvc.DescribeUHostInstance(&describeParams)
	if err != nil {
		return nil, err
	}

	if resp.RetCode != 0 {
		return nil, fmt.Errorf("Describe UHost error, Retcode:%d, err:%s", resp.RetCode, err)
	}

	if &resp.UHostSet[0] == nil {
		return nil, fmt.Errorf("UHostSet is empty")
	}
	hostState := resp.UHostSet[0].State

	// TODO: Now we get eip,later we will return the type of ip by options
	var hostIpAddress string
	for _, ip := range resp.UHostSet[0].IPSet {
		if ip.Type == "Private" {
			continue
		} else {
			hostIpAddress = ip.IP
		}
	}

	details := &UHostDetail{
		region:    region,
		hostID:    hostId,
		state:     hostState,
		ipAddress: hostIpAddress,
		cpu:       resp.UHostSet[0].CPU,
		memory:    resp.UHostSet[0].Memory,
	}

	return details, nil
}


// createUNet create network for uhost
func (d *Driver) createUNet() error {
	createEIPParams := unet.AllocateEIPParams{
		Region: d.Region,
		OperatorName: "Bgp",
		Bandwidth: 2,
		ChargeType: "Dynamic",
		Quantity: 1,
	}

	resp, err := unetsvc.AllocateEIP(&createEIPParams)
	if err != nil {
		return fmt.Errorf("Allocate EIP failed:%s", err)
	}
	utils.DumpVal(resp)
	if resp.RetCode != 0 {
		return fmt.Errorf("Allocate EIP failed")
	}

	// FIXME: it is ugly here to get eip
	eipId := (*resp.EIPSet)[0].EIPId
	d.IPAddress = (*(*resp.EIPSet)[0].EIPAddr)[0].IP

	bindHostParams := unet.BindEIPParams{
		Region: d.Region,
		EIPId: eipId,
		ResourceType: "uhost",
		ResourceId: d.UhostID,
	}

	bindEIPResp, err := unetsvc.BindEIP(&bindHostParams)
	if err != nil {
		return fmt.Errorf("Bind EIP failed:%s", err)
	}
	utils.DumpVal(bindEIPResp)

	// create security group
	securityGroupParams := unet.CreateSecurityGroupParams{
		Region: d.Region,
		GroupName: "docker-machine",
		Description: "docker machine to open 2379 and 22 port of tcp",
		Rule: []string{"TCP|22|0.0.0.0/0|ACCEPT|50", "TCP|3389|0.0.0.0/0|ACCEPT|50","TCP|2379|0.0.0.0/0|ACCEPT|50"},
	}
	createSecurityGroupResp, err := unetsvc.CreateSecurityGroup(&securityGroupParams)
	if err != nil {
		return fmt.Errorf("create security group failed:%s", err)
	}
	utils.DumpVal(createSecurityGroupResp)

	// TODO： because CreateSecurityGroup don't return GroupId, so we have to
	// iterate all security groups to find the right one, Here we use the given
	// security group for testing.
	grantSecurityGroupParams := unet.GrantSecurityGroupParams{
		Region: d.Region,
		GroupId: 33149, //"docker machine"
		ResourceType: "uhost",
		ResourceId: d.UhostID,
	}
	grantSecurityResp, err := unetsvc.GrantSecurityGroup(&grantSecurityGroupParams)
	if err != nil {
		return fmt.Errorf("grant security group failed:%s", err)
	}
	utils.DumpVal(grantSecurityResp)



	return nil
}