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

func (p *pluginContainer) DoDo(t PluginType, v ...interface{}) []error {

	list := p.plugins[t]
	rsp := make([]error, 0)
	for _, plugin := range list {
		err := plugin.OnDo(v...)
		if err != nil {
			rsp = append(rsp, err)
		}
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
	OnDo(v ...interface{}) error
	Type() PluginType
	Init(v ...interface{}) Plugin
}
