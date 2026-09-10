package cpmissing

import (
	"fmt"
	"sdmm/internal/app/command"
	"sdmm/internal/dmapi/dmmap" //for the undefinedvars struct
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/imguiext/icon"
	w "sdmm/internal/imguiext/widget"
	"sdmm/internal/util"

	"github.com/SpaiR/imgui-go"
)

func (m *Missing) Process(int32) {

	m.ShowControls()
}

func (m *Missing) removeUndefinedVariableFromPrefab(undef dmmap.UndefinedVar, i int) {
	tile := m.app.CurrentEditor().Dmm().GetTile(util.Point{X: undef.X, Y: undef.Y, Z: undef.Z})
	for _, instance := range tile.Instances() {
		if instance.Id() == undef.PrefabInfo {
			prefab := instance.Prefab()
			vars := prefab.Vars()
			for _, varName := range prefab.Vars().Iterate() {
				if undef.VarName == varName {

					//the doing part
					wsName := m.app.CurrentEditor().Dmm().Name
					vars = dmvars.Delete(vars, varName)
					instance.SetPrefab(dmmprefab.New(dmmprefab.IdNone, prefab.Path(), vars))
					//m.app.SyncPrefabs()
					m.app.SyncVarEditor()

					m.UndefinedVars = append(m.UndefinedVars[:i], m.UndefinedVars[i+1:]...)

					if wsVars, exists := m.WorkSpaceVars[wsName]; exists {
						// match variable inside the WorkSpaceVars-Undefined vars to what we just removed
						for idx, v := range wsVars {
							// match?
							if v.Path == undef.Path && v.VarName == undef.VarName && v.PrefabInfo == undef.PrefabInfo {
								wsVars = append(wsVars[:idx], wsVars[idx+1:]...)
								break
							}
						}

						if len(wsVars) == 0 {
							delete(m.WorkSpaceVars, wsName)
						} else {
							m.WorkSpaceVars[wsName] = wsVars
						}
					}

					//registers the redo and undo
					m.app.CommandStorage().Push(command.Make("Removed Undefined", func() {
						vars = dmvars.Set(vars, varName, undef.VarValue)
						instance.SetPrefab(dmmprefab.New(dmmprefab.IdNone, prefab.Path(), vars))
						//m.app.SyncPrefabs()
						m.app.SyncVarEditor()
						m.UndefinedVars = append(m.UndefinedVars[:i], append([]dmmap.UndefinedVar{undef}, m.UndefinedVars[i:]...)...)

						if m.WorkSpaceVars == nil {
							m.WorkSpaceVars = make(map[string][]dmmap.UndefinedVar)
						}

						m.WorkSpaceVars[wsName] = append(m.WorkSpaceVars[wsName], undef)

					}, func() {
						vars = dmvars.Delete(vars, varName)
						instance.SetPrefab(dmmprefab.New(dmmprefab.IdNone, prefab.Path(), vars))
						//m.app.SyncPrefabs()
						m.app.SyncVarEditor()
						m.UndefinedVars = append(m.UndefinedVars[:i], m.UndefinedVars[i+1:]...)

						if wsVars, exists := m.WorkSpaceVars[wsName]; exists {
							for idx, v := range wsVars {
								if v.Path == undef.Path && v.VarName == undef.VarName && v.PrefabInfo == undef.PrefabInfo {
									wsVars = append(wsVars[:idx], wsVars[idx+1:]...)
									break
								}
							}
							if len(wsVars) == 0 {
								delete(m.WorkSpaceVars, wsName)
							} else {
								m.WorkSpaceVars[wsName] = wsVars
							}
						}
					}))

				}
			}
		}
	}
}

func (m *Missing) ShowControls() {
	if imgui.BeginChild("results") {
		tableFlags := imgui.TableFlagsResizable |
			imgui.TableFlagsSizingStretchProp |
			imgui.TableFlagsBordersInnerH
		imgui.BeginTableV("Undefined Variables", 4, tableFlags, imgui.Vec2{X: 0, Y: 0}, 0.0)
		imgui.TableSetupColumn("")
		imgui.TableSetupColumn("")
		imgui.TableSetupColumn("")
		imgui.TableSetupColumnV("", imgui.TableColumnFlagsWidthFixed, 0.0, 0)
		for i, undef := range m.UndefinedVars {
			imgui.TableNextRow()
			imgui.TableSetColumnIndex(0)
			imgui.Text(undef.Path)

			imgui.TableSetColumnIndex(1)
			coords := fmt.Sprintf("X:%d Y:%d Z:%d", undef.X, undef.Y, undef.Z)
			imgui.Text(coords)

			imgui.TableSetColumnIndex(2)
			imgui.Text(undef.VarName)

			//imgui.TableSetColumnIndex(3)
			//imgui.Text(undef.VarValue)

			imgui.TableSetColumnIndex(3)
			imgui.BeginGroup()
			buttonLabelJump := fmt.Sprintf(icon.Search+"##%d", i)
			buttonLabelNull := fmt.Sprintf(icon.Clear+"##%d", i)
			w.Layout{
				w.Button(buttonLabelJump, func() {
					m.app.CurrentEditor().FocusCameraOnPosition(util.Point{X: undef.X, Y: undef.Y, Z: undef.Z})
				}).Small(true),
				w.SameLine(),
				w.Button(buttonLabelNull, func() {
					m.removeUndefinedVariableFromPrefab(undef, i)
				}).Small(true),
			}.Build()
			imgui.EndGroup()

		}
		imgui.EndTable()
		if imgui.Button("Dismiss") {
			m.UndefinedVars = nil

			wsName := m.app.CurrentEditor().Dmm().Name
			delete(m.WorkSpaceVars, wsName)
		}
		imgui.EndChild()
	}
}
