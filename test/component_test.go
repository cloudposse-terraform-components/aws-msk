package test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/aws-component-helper"
	"github.com/gruntwork-io/terratest/modules/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestComponent(t *testing.T) {
	// Define the AWS region to use for the tests
	awsRegion := "us-east-2"

	// Initialize the test fixture
	fixture := helper.NewFixture(t, "../", awsRegion, "test/fixtures")

	// Ensure teardown is executed after the test
	defer fixture.TearDown()
	fixture.SetUp(&atmos.Options{})

	// Define the test suite
	fixture.Suite("default", func(t *testing.T, suite *helper.Suite) {
		suite.AddDependency("vpc", "default-test")

		// Setup phase: Create DNS zones for testing
		suite.Setup(t, func(t *testing.T, atm *helper.Atmos) {
			basicDomain := "components.cptest.test-automation.app"

			// Deploy the delegated DNS zone
			inputs := map[string]interface{}{
				"zone_config": []map[string]interface{}{
					{
						"subdomain": suite.GetRandomIdentifier(),
						"zone_name": basicDomain,
					},
				},
			}
			atm.GetAndDeploy("dns-delegated", "default-test", inputs)
		})

		// Teardown phase: Destroy the DNS zones created during setup
		suite.TearDown(t, func(t *testing.T, atm *helper.Atmos) {
			// Deploy the delegated DNS zone
			inputs := map[string]interface{}{
				"zone_config": []map[string]interface{}{
					{
						"subdomain": suite.GetRandomIdentifier(),
						"zone_name": "components.cptest.test-automation.app",
					},
				},
			}
			atm.GetAndDestroy("dns-delegated", "default-test", inputs)
		})

		// Test phase: Validate the functionality of the ALB component
		suite.Test(t, "basic", func(t *testing.T, atm *helper.Atmos) {
			suffix := strings.ToLower(random.UniqueId())
			inputs := map[string]interface{}{
				"name": "msk-" + suffix,
			}

			defer atm.GetAndDestroy("msk/basic", "default-test", inputs)
			component := atm.GetAndDeploy("msk/basic", "default-test", inputs)
			assert.NotNil(t, component)

			clusterName := atm.Output(component, "cluster_name")
			assert.NotEmpty(t, clusterName)

			clusterArn := atm.Output(component, "cluster_arn")
			assert.Contains(t, clusterArn, clusterName)

			storageMode := atm.Output(component, "storage_mode")
			assert.Equal(t, "LOCAL", storageMode)

			bootstrapBrokers := atm.Output(component, "bootstrap_brokers")
			bootstrapBrokersList := strings.Split(bootstrapBrokers, ",")
			assert.Equal(t, 2, len(bootstrapBrokersList))

			bootstrapBrokersTLS := atm.Output(component, "bootstrap_brokers_tls")
			bootstrapBrokersTLSList := strings.Split(bootstrapBrokersTLS, ",")
			assert.Equal(t, 2, len(bootstrapBrokersTLSList))

			bootstrapBrokersPublicTLS := atm.Output(component, "bootstrap_brokers_public_tls")
			assert.Empty(t, bootstrapBrokersPublicTLS)

			bootstrapBrokersSaslScram := atm.Output(component, "bootstrap_brokers_sasl_scram")
			assert.Equal(t, "", bootstrapBrokersSaslScram)

			bootstrapBrokersPublicSaslScram := atm.Output(component, "bootstrap_brokers_public_sasl_scram")
			assert.Equal(t, "", bootstrapBrokersPublicSaslScram)

			bootstrapBrokersSaslIam := atm.Output(component, "bootstrap_brokers_sasl_iam")
			assert.Equal(t, "", bootstrapBrokersSaslIam)

			bootstrapBrokersPublicSaslIam := atm.Output(component, "bootstrap_brokers_public_sasl_iam")
			assert.Equal(t, "", bootstrapBrokersPublicSaslIam)

			zookeeperConnectString := atm.Output(component, "zookeeper_connect_string")
			assert.Equal(t, 3, len(strings.Split(zookeeperConnectString, ",")))

			zookeeperConnectStringTLS := atm.Output(component, "zookeeper_connect_string_tls")
			assert.Equal(t, 3, len(strings.Split(zookeeperConnectStringTLS, ",")))

			brokerEndpoints := atm.OutputList(component, "broker_endpoints")
			assert.Equal(t, 2, len(brokerEndpoints))

			currentVersion := atm.Output(component, "current_version")
			assert.NotEmpty(t, currentVersion)

			configArn := atm.Output(component, "config_arn")
			assert.NotEmpty(t, configArn)

			var latestRevision int
			atm.OutputStruct(component, "latest_revision", &latestRevision)
			assert.NotEmpty(t, latestRevision)

			hostnames := atm.OutputList(component, "hostnames")
			assert.Equal(t, 2, len(hostnames))

			securityGroupId := atm.Output(component, "security_group_id")
			assert.True(t, strings.HasPrefix(securityGroupId, "sg-"))

			securityGroupArn := atm.Output(component, "security_group_arn")
			assert.Contains(t, securityGroupArn, securityGroupId)

			securityGroupName := atm.Output(component, "security_group_name")
			assert.Contains(t, securityGroupName, clusterName)

			client := NewMSKClient(t, awsRegion)
			describeClusterOutput, err := client.DescribeCluster(context.Background(), &kafka.DescribeClusterInput{
				ClusterArn: &clusterArn,
			})
			require.NoError(t, err)
			cluster := describeClusterOutput.ClusterInfo
			assert.EqualValues(t, "ACTIVE", cluster.State)

			assert.EqualValues(t, clusterName, *cluster.ClusterName)
			assert.EqualValues(t, clusterArn, *cluster.ClusterArn)
			assert.EqualValues(t, storageMode, cluster.StorageMode)

			awsBootstrapBrokers, err := client.GetBootstrapBrokers(context.Background(), &kafka.GetBootstrapBrokersInput{
				ClusterArn: &clusterArn,
			})
			require.NoError(t, err)

			assert.ElementsMatch(t, bootstrapBrokersList, strings.Split(*awsBootstrapBrokers.BootstrapBrokerString, ","))
			assert.ElementsMatch(t, bootstrapBrokersTLSList, strings.Split(*awsBootstrapBrokers.BootstrapBrokerStringTls, ","))
			assert.Nil(t, awsBootstrapBrokers.BootstrapBrokerStringPublicTls)
			assert.Nil(t, awsBootstrapBrokers.BootstrapBrokerStringSaslScram)
			assert.Nil(t, awsBootstrapBrokers.BootstrapBrokerStringPublicSaslScram)

			assert.EqualValues(t, zookeeperConnectString, *cluster.ZookeeperConnectString)
			assert.EqualValues(t, zookeeperConnectStringTLS, *cluster.ZookeeperConnectStringTls)

			// assert.EqualValues(t, brokerEndpoints, cluster.BrokerNodeGroupInfo.BrokerEndpoints)
			assert.EqualValues(t, currentVersion, *cluster.CurrentVersion)
			// assert.EqualValues(t, latestRevision, cluster)
			// assert.EqualValues(t, hostnames, cluster)
			assert.EqualValues(t, securityGroupId, cluster.BrokerNodeGroupInfo.SecurityGroups[0])
			// assert.EqualValues(t, securityGroupArn, cluster.SecurityGroupArn)
			// assert.EqualValues(t, securityGroupName, cluster.SecurityGroupName)
		})
	})
}

func NewMSKClient(t *testing.T, region string) *kafka.Client {
	client, err := NewMSKClientE(t, region)
	require.NoError(t, err)

	return client
}

func NewMSKClientE(t *testing.T, region string) (*kafka.Client, error) {
	sess, err := aws.NewAuthenticatedSession(region)
	if err != nil {
		return nil, err
	}
	return kafka.NewFromConfig(*sess), nil
}
