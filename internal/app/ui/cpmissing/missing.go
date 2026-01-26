package cpmissing

import (
	"sdmm/internal/app/ui/component"
	"sdmm/internal/app/ui/cpwsarea/wsmap/pmap/editor"

	"sdmm/internal/dmapi/dmmap/dmminstance"

	"sdmm/internal/dmapi/dmmap" //for the undefinedvars struct
)

func (m *Missing) Init(app App) {
	m.app = app
}

type App interface {
	CurrentEditor() *editor.Editor
	DoEditInstance(*dmminstance.Instance)
	ShowLayout(name string, focus bool)
}

type Missing struct {
	component.Component

	app           App
	index         int
	UndefinedVars []dmmap.UndefinedVar
}
