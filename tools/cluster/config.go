package cluster

import (
	"encoding/json"
	"fmt"
	"os"
	"path"

	"github.com/nnikolash/wasp-types-exported/packages/cryptolib"
	"github.com/nnikolash/wasp-types-exported/packages/l1connection"
	"github.com/nnikolash/wasp-types-exported/tools/cluster/templates"
)

type WaspConfig struct {
	NumNodes int

	// node ports are calculated as these values + node index
	FirstAPIPort       int
	FirstPeeringPort   int
	FirstProfilingPort int
	FirstMetricsPort   int
}

func (w *WaspConfig) WaspConfigTemplateParams(i int) templates.WaspConfigParams {
	return templates.WaspConfigParams{
		APIPort:                w.FirstAPIPort + i,
		PeeringPort:            w.FirstPeeringPort + i,
		ProfilingPort:          w.FirstProfilingPort + i,
		MetricsPort:            w.FirstMetricsPort + i,
		PruningMinStatesToKeep: 10000,
		AuthScheme:             "none",
	}
}

type ClusterConfig struct {
	Wasp []templates.WaspConfigParams
	L1   l1connection.Config
}

func DefaultWaspConfig() WaspConfig {
	return WaspConfig{
		NumNodes:           4,
		FirstAPIPort:       19090,
		FirstPeeringPort:   14000,
		FirstProfilingPort: 11060,
		FirstMetricsPort:   12112,
	}
}

func ConfigExists(dataPath string) (bool, error) {
	return fileExists(configPath(dataPath))
}

func NewConfig(waspConfig WaspConfig, l1Config l1connection.Config, modifyConfig ...templates.ModifyNodesConfigFn) *ClusterConfig {
	nodesConfigs := make([]templates.WaspConfigParams, waspConfig.NumNodes)
	for i := 0; i < waspConfig.NumNodes; i++ {
		// generate template from waspconfigs
		nodesConfigs[i] = waspConfig.WaspConfigTemplateParams(i)
		// set L1 part of the template
		nodesConfigs[i].L1INXAddress = l1Config.INXAddress
		// modify the template if needed
		if len(modifyConfig) > 0 && modifyConfig[0] != nil {
			nodesConfigs[i] = modifyConfig[0](i, nodesConfigs[i])
		}
	}

	return &ClusterConfig{
		Wasp: nodesConfigs,
		L1:   l1Config,
	}
}

func LoadConfig(dataPath string) (*ClusterConfig, error) {
	b, err := os.ReadFile(configPath(dataPath))
	if err != nil {
		return nil, err
	}
	var c ClusterConfig
	err = json.Unmarshal(b, &c)
	return &c, err
}

func (c *ClusterConfig) Save(dataPath string) error {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath(dataPath), b, 0o600)
}

func configPath(dataPath string) string {
	return path.Join(dataPath, "cluster.json")
}

func (c *ClusterConfig) setValidatorAddressIfNotSet() {
	for i := range c.Wasp {
		if c.Wasp[i].ValidatorKeyPair == nil {
			kp := cryptolib.NewKeyPair()
			c.Wasp[i].ValidatorKeyPair = kp
			c.Wasp[i].ValidatorAddress = kp.Address().Bech32("atoi") // privtangle bech32
		}
	}
}

func (c *ClusterConfig) waspHosts(nodeIndexes []int, getHost func(i int) string) []string {
	hosts := make([]string, 0)
	for _, i := range nodeIndexes {
		if i < 0 || i > len(c.Wasp)-1 {
			panic(fmt.Sprintf("Node index out of bounds in smart contract configuration: %d/%d", i, len(c.Wasp)))
		}
		hosts = append(hosts, getHost(i))
	}
	return hosts
}

func (c *ClusterConfig) AllNodes() []int {
	nodes := make([]int, len(c.Wasp))
	for i := 0; i < len(c.Wasp); i++ {
		nodes[i] = i
	}
	return nodes
}

func (c *ClusterConfig) APIHosts(nodeIndexes ...[]int) []string {
	nodes := c.AllNodes()
	if len(nodeIndexes) == 1 {
		nodes = nodeIndexes[0]
	}
	return c.waspHosts(nodes, func(i int) string { return c.APIHost(i) })
}

func (c *ClusterConfig) APIHost(nodeIndex int) string {
	return fmt.Sprintf("http://127.0.0.1:%d", c.APIPort(nodeIndex))
}

func (c *ClusterConfig) APIPort(nodeIndex int) int {
	return c.Wasp[nodeIndex].APIPort
}

func (c *ClusterConfig) ValidatorAddress(nodeIndex int) string {
	return c.Wasp[nodeIndex].ValidatorAddress
}

func (c *ClusterConfig) ValidatorKeyPair(nodeIndex int) *cryptolib.KeyPair {
	return c.Wasp[nodeIndex].ValidatorKeyPair
}

func (c *ClusterConfig) PeeringHosts(nodeIndexes ...[]int) []string {
	nodes := c.AllNodes()
	if len(nodeIndexes) == 1 {
		nodes = nodeIndexes[0]
	}
	return c.waspHosts(nodes, func(i int) string { return c.PeeringHost(i) })
}

func (c *ClusterConfig) PeeringHost(nodeIndex int) string {
	return fmt.Sprintf("127.0.0.1:%d", c.PeeringPort(nodeIndex))
}

func (c *ClusterConfig) PeeringPort(nodeIndex int) int {
	return c.Wasp[nodeIndex].PeeringPort
}

func (c *ClusterConfig) L1APIAddress(nodeIndex int) string {
	return c.L1.APIAddress
}

func (c *ClusterConfig) ProfilingPort(nodeIndex int) int {
	return c.Wasp[nodeIndex].ProfilingPort
}

func (c *ClusterConfig) PrometheusPort(nodeIndex int) int {
	return c.Wasp[nodeIndex].MetricsPort
}
