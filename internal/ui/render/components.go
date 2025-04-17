package render

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/cpaluszek/pipeye/internal/ui"
)

func NewStyledTable(headers []string, width, selectedIndex int) *table.Table {
	return table.New().
		Border(lipgloss.NormalBorder()).
		BorderHeader(true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderBottom(false).
		BorderColumn(false).
		Headers(headers...).
		Width(width).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch row {
			case table.HeaderRow:
				return ui.TableHeaderStyle
			case selectedIndex:
				return ui.SelectedRowStyle
			default:
				return ui.RowStyle
			}
		})
}
