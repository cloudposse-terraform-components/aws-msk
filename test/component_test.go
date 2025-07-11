package test

import (
	"context"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/kafka"
	"github.com/cloudposse/test-helpers/pkg/atmos"
	helper "github.com/cloudposse/test-helpers/pkg/atmos/component-helper"
	awshelper "github.com/cloudposse/test-helpers/pkg/aws"
	"github.com/gruntwork-io/terratest/modules/random"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type ComponentSuite struct {
	helper.TestSuite
}

func (s *ComponentSuite) TestBasic() {
	const component = "msk/basic"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	suffix := strings.ToLower(random.UniqueId())
	inputs := map[string]any{
		"name": "msk-" + suffix,
	}

	defer s.DestroyAtmosComponent(s.T(), component, stack, &inputs)
	options, _ := s.DeployAtmosComponent(s.T(), component, stack, &inputs)
	assert.NotNil(s.T(), options)

	clusterName := atmos.Output(s.T(), options, "cluster_name")
	assert.NotEmpty(s.T(), clusterName)

	clusterArn := atmos.Output(s.T(), options, "cluster_arn")
	assert.Contains(s.T(), clusterArn, clusterName)

	storageMode := atmos.Output(s.T(), options, "storage_mode")
	assert.Equal(s.T(), "LOCAL", storageMode)

	bootstrapBrokers := atmos.Output(s.T(), options, "bootstrap_brokers")
	bootstrapBrokersList := strings.Split(bootstrapBrokers, ",")
	assert.Equal(s.T(), 2, len(bootstrapBrokersList))

	bootstrapBrokersTLS := atmos.Output(s.T(), options, "bootstrap_brokers_tls")
	bootstrapBrokersTLSList := strings.Split(bootstrapBrokersTLS, ",")
	assert.Equal(s.T(), 2, len(bootstrapBrokersTLSList))

	bootstrapBrokersPublicTLS := atmos.Output(s.T(), options, "bootstrap_brokers_public_tls")
	assert.Empty(s.T(), bootstrapBrokersPublicTLS)

	// Additional assertions for SASL and Zookeeper
	// ...

	ssmKeyPaths := atmos.OutputList(s.T(), options, "ssm_key_paths")
	assert.Equal(s.T(), 5, len(ssmKeyPaths))

	client := awshelper.NewMSKClient(s.T(), awsRegion)
	describeClusterOutput, err := client.DescribeCluster(context.Background(), &kafka.DescribeClusterInput{
		ClusterArn: &clusterArn,
	})
	require.NoError(s.T(), err)
	cluster := describeClusterOutput.ClusterInfo
	assert.EqualValues(s.T(), "ACTIVE", cluster.State)

	assert.EqualValues(s.T(), clusterName, *cluster.ClusterName)
	assert.EqualValues(s.T(), clusterArn, *cluster.ClusterArn)
	assert.EqualValues(s.T(), storageMode, cluster.StorageMode)

	s.DriftTest(component, stack, &inputs)
}

func (s *ComponentSuite) TestEnabledFlag() {
	const component = "msk/disabled"
	const stack = "default-test"
	const awsRegion = "us-east-2"

	suffix := strings.ToLower(random.UniqueId())
	inputs := map[string]interface{}{
		"name": "msk-" + suffix,
	}

	s.VerifyEnabledFlag(component, stack, &inputs)
}

func TestRunSuite(t *testing.T) {
	suite := new(ComponentSuite)

	suite.AddDependency(t, "vpc", "default-test", nil)

	subdomain := strings.ToLower(random.UniqueId())
	inputs := map[string]interface{}{
		"zone_config": []map[string]interface{}{
			{
				"subdomain": subdomain,
				"zone_name": "components.cptest.test-automation.app",
			},
		},
	}
	suite.AddDependency(t, "dns-delegated", "default-test", &inputs)
	helper.Run(t, suite)
}
