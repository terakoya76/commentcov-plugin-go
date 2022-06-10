package main

import (
	"github.com/commentcov/commentcov/pkg/pluggable"
	"github.com/commentcov/commentcov/proto"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-plugin"

	"github.com/commentcov/commentcov-plugin-go/ast"
)

// pluginImpl implements pluggable.Pluggable.
type pluginImpl struct {
	logger hclog.Logger
}

// MeasureCoverage is the implementation of pluggable.Pluggable.
func (i *pluginImpl) MeasureCoverage(files []string) ([]*proto.CoverageItem, error) {
	items := make([]*proto.CoverageItem, 0)

	for _, file := range files {
		cis, err := ast.FileToCoverageItems(i.logger, file)
		if err != nil {
			i.logger.Trace(err.Error())
			return []*proto.CoverageItem{}, err
		}

		items = append(items, cis...)
	}

	return items, nil
}

// main is entrypoint as plugin.
// Serving MeasureCoverage as gRPC Server.
func main() {
	logger := hclog.New(&hclog.LoggerOptions{
		Level:      hclog.Trace,
		Name:       "commmentcov-plugin-go",
		JSONFormat: true,
	})

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: pluggable.PluginHandshakeConfig,
		VersionedPlugins: map[int]plugin.PluginSet{
			1: {
				"commentcov": &pluggable.CommentcovPlugin{
					Impl: &pluginImpl{
						logger: logger,
					},
				},
			},
		},
		GRPCServer: plugin.DefaultGRPCServer,
	})
}
