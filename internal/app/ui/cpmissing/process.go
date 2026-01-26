package cpmissing

import (
	"fmt"
	"sdmm/internal/dmapi/dmmap/dmmdata/dmmprefab"
	"sdmm/internal/dmapi/dmvars"
	"sdmm/internal/imguiext/icon"
	"sdmm/internal/util"

	"github.com/SpaiR/imgui-go"
)

func (m *Missing) Process(int32) {

	m.ShowControls()
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
			if imgui.Button(buttonLabelJump) {
				m.app.CurrentEditor().FocusCameraOnPosition(util.Point{X: undef.X, Y: undef.Y, Z: undef.Z})
			}

			buttonLabelNull := fmt.Sprintf(icon.Clear+"##%d", i)
			if imgui.Button(buttonLabelNull) {
				tile := m.app.CurrentEditor().Dmm().GetTile(util.Point{X: undef.X, Y: undef.Y, Z: undef.Z})
				for _, instance := range tile.Instances() {
					if instance.Id() == undef.PrefabInfo {
						prefab := instance.Prefab()
						vars := prefab.Vars()
						for _, varName := range prefab.Vars().Iterate() {
							if undef.VarName == varName {
								vars = dmvars.Delete(vars, varName)
								instance.SetPrefab(dmmprefab.New(dmmprefab.IdNone, prefab.Path(), vars))
								m.app.CurrentEditor().CommitChanges("Undef Variable bombed")
								// i is your index number
								m.UndefinedVars = append(m.UndefinedVars[:i], m.UndefinedVars[i+1:]...)
							}
						}
					}
				}
			}
			imgui.EndGroup()

		}
		imgui.EndTable()
		imgui.EndChild()
	}
}
