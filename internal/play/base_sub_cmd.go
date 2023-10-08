package play

// addSubCommands 添加二级命令
func (p *Playground) addSubCommands() {
	p.subCmds = append(p.subCmds,
		p.playDemo(),
		p.playGorm(),
		p.playFeishu(),
	)
}
