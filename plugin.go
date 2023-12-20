package gnetrpc

type PluginType uint16

const (
	PluginTypeOnBoot PluginType = iota
	PluginTypeOnShutdown
	PluginTypeOnOpen
	PluginTypeOnClose
	PluginTypeOnTraffic
	PluginTypeOnTick
)

type pluginContainer struct {
	plugins map[PluginType][]Plugin
}

func (p *pluginContainer) DoDo(t PluginType, v ...interface{}) []interface{} {

	list := p.plugins[t]
	rsp := make([]interface{}, 0, len(list))
	for _, plugin := range list {
		rsp = append(rsp, plugin.OnDo(v...))
	}
	return rsp
}
func (p *pluginContainer) Plugins(t PluginType) []Plugin {
	return p.plugins[t]
}
func (p *pluginContainer) Add(t PluginType, plugin Plugin) {
	if list, ok := p.plugins[t]; ok {
		list = append(list, plugin)
		p.plugins[t] = list
	} else {
		p.plugins[t] = []Plugin{plugin}
	}
}

type Plugin interface {
	OnDo(v ...interface{}) interface{}
	Type() PluginType
}
