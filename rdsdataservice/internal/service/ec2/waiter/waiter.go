package waiter

import (
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

const (
	// Maximum amount of time to wait for a LocalGatewayRouteTableVpcAssociation to return Associated
	LocalGatewayRouteTableVpcAssociationAssociatedTimeout = 5 * time.Minute

	// Maximum amount of time to wait for a LocalGatewayRouteTableVpcAssociation to return Disassociated
	LocalGatewayRouteTableVpcAssociationDisassociatedTimeout = 5 * time.Minute
)

// LocalGatewayRouteTableVpcAssociationAssociated waits for a LocalGatewayRouteTableVpcAssociation to return Associated
func LocalGatewayRouteTableVpcAssociationAssociated(conn *ec2.EC2, localGatewayRouteTableVpcAssociationID string) (*ec2.LocalGatewayRouteTableVpcAssociation, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.RouteTableAssociationStateCodeAssociating},
		Target:  []string{ec2.RouteTableAssociationStateCodeAssociated},
		Refresh: LocalGatewayRouteTableVpcAssociationState(conn, localGatewayRouteTableVpcAssociationID),
		Timeout: LocalGatewayRouteTableVpcAssociationAssociatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.LocalGatewayRouteTableVpcAssociation); ok {
		return output, err
	}

	return nil, err
}

// LocalGatewayRouteTableVpcAssociationDisassociated waits for a LocalGatewayRouteTableVpcAssociation to return Disassociated
func LocalGatewayRouteTableVpcAssociationDisassociated(conn *ec2.EC2, localGatewayRouteTableVpcAssociationID string) (*ec2.LocalGatewayRouteTableVpcAssociation, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.RouteTableAssociationStateCodeDisassociating},
		Target:  []string{ec2.RouteTableAssociationStateCodeDisassociated},
		Refresh: LocalGatewayRouteTableVpcAssociationState(conn, localGatewayRouteTableVpcAssociationID),
		Timeout: LocalGatewayRouteTableVpcAssociationAssociatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.LocalGatewayRouteTableVpcAssociation); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnEndpointDeletedTimout = 5 * time.Minute
)

func ClientVpnEndpointDeleted(conn *ec2.EC2, id string) (*ec2.ClientVpnEndpoint, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnEndpointStatusCodeDeleting},
		Target:  []string{},
		Refresh: ClientVpnEndpointStatus(conn, id),
		Timeout: ClientVpnEndpointDeletedTimout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.ClientVpnEndpoint); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnAuthorizationRuleActiveTimeout = 10 * time.Minute

	ClientVpnAuthorizationRuleRevokedTimeout = 10 * time.Minute
)

func ClientVpnAuthorizationRuleAuthorized(conn *ec2.EC2, authorizationRuleID string) (*ec2.AuthorizationRule, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnAuthorizationRuleStatusCodeAuthorizing},
		Target:  []string{ec2.ClientVpnAuthorizationRuleStatusCodeActive},
		Refresh: ClientVpnAuthorizationRuleStatus(conn, authorizationRuleID),
		Timeout: ClientVpnAuthorizationRuleActiveTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.AuthorizationRule); ok {
		return output, err
	}

	return nil, err
}

func ClientVpnAuthorizationRuleRevoked(conn *ec2.EC2, authorizationRuleID string) (*ec2.AuthorizationRule, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnAuthorizationRuleStatusCodeRevoking},
		Target:  []string{},
		Refresh: ClientVpnAuthorizationRuleStatus(conn, authorizationRuleID),
		Timeout: ClientVpnAuthorizationRuleRevokedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.AuthorizationRule); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnNetworkAssociationAssociatedTimeout = 10 * time.Minute

	ClientVpnNetworkAssociationAssociatedDelay = 4 * time.Minute

	ClientVpnNetworkAssociationDisassociatedTimeout = 10 * time.Minute

	ClientVpnNetworkAssociationDisassociatedDelay = 4 * time.Minute

	ClientVpnNetworkAssociationStatusPollInterval = 10 * time.Second
)

func ClientVpnNetworkAssociationAssociated(conn *ec2.EC2, networkAssociationID, clientVpnEndpointID string) (*ec2.TargetNetwork, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{ec2.AssociationStatusCodeAssociating},
		Target:       []string{ec2.AssociationStatusCodeAssociated},
		Refresh:      ClientVpnNetworkAssociationStatus(conn, networkAssociationID, clientVpnEndpointID),
		Timeout:      ClientVpnNetworkAssociationAssociatedTimeout,
		Delay:        ClientVpnNetworkAssociationAssociatedDelay,
		PollInterval: ClientVpnNetworkAssociationStatusPollInterval,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.TargetNetwork); ok {
		return output, err
	}

	return nil, err
}

func ClientVpnNetworkAssociationDisassociated(conn *ec2.EC2, networkAssociationID, clientVpnEndpointID string) (*ec2.TargetNetwork, error) {
	stateConf := &resource.StateChangeConf{
		Pending:      []string{ec2.AssociationStatusCodeDisassociating},
		Target:       []string{},
		Refresh:      ClientVpnNetworkAssociationStatus(conn, networkAssociationID, clientVpnEndpointID),
		Timeout:      ClientVpnNetworkAssociationDisassociatedTimeout,
		Delay:        ClientVpnNetworkAssociationDisassociatedDelay,
		PollInterval: ClientVpnNetworkAssociationStatusPollInterval,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.TargetNetwork); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnRouteDeletedTimeout = 1 * time.Minute
)

func ClientVpnRouteDeleted(conn *ec2.EC2, routeID string) (*ec2.ClientVpnRoute, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnRouteStatusCodeActive, ec2.ClientVpnRouteStatusCodeDeleting},
		Target:  []string{},
		Refresh: ClientVpnRouteStatus(conn, routeID),
		Timeout: ClientVpnRouteDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.ClientVpnRoute); ok {
		return output, err
	}

	return nil, err
}

func SecurityGroupCreated(conn *ec2.EC2, id string, timeout time.Duration) (*ec2.SecurityGroup, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{SecurityGroupStatusNotFound},
		Target:  []string{SecurityGroupStatusCreated},
		Refresh: SecurityGroupStatus(conn, id),
		Timeout: timeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.SecurityGroup); ok {
		return output, err
	}

	return nil, err
}

const (
	VpnGatewayVpcAttachmentAttachedTimeout = 15 * time.Minute

	VpnGatewayVpcAttachmentDetachedTimeout = 30 * time.Minute
)

func VpnGatewayVpcAttachmentAttached(conn *ec2.EC2, vpnGatewayID, vpcID string) (*ec2.VpcAttachment, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.AttachmentStatusDetached, ec2.AttachmentStatusAttaching},
		Target:  []string{ec2.AttachmentStatusAttached},
		Refresh: VpnGatewayVpcAttachmentState(conn, vpnGatewayID, vpcID),
		Timeout: VpnGatewayVpcAttachmentAttachedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.VpcAttachment); ok {
		return output, err
	}

	return nil, err
}

func VpnGatewayVpcAttachmentDetached(conn *ec2.EC2, vpnGatewayID, vpcID string) (*ec2.VpcAttachment, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.AttachmentStatusAttached, ec2.AttachmentStatusDetaching},
		Target:  []string{ec2.AttachmentStatusDetached},
		Refresh: VpnGatewayVpcAttachmentState(conn, vpnGatewayID, vpcID),
		Timeout: VpnGatewayVpcAttachmentDetachedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.VpcAttachment); ok {
		return output, err
	}

	return nil, err
}
