package play

import "github.com/spf13/cobra"

// addSubCommands 添加二级命令
func (p *Playground) addSubCommands() []*cobra.Command {
	subCmds := []*cobra.Command{
		p.playDemo(),
		p.playGorm(),
		p.playFeishu(),
		p.playPassword(),
		p.playCasbin(),
	}

	return subCmds
}
