package core

import (
	"fmt"
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	help   = "2048 GAME\nPress Esc to exit, press hjkl or arrow keys to move."
	youWin = "You win!!!!!"
)

type core struct {
	board [][]int
	cols  int
	rows  int

	app   *tview.Application
	txt   *tview.TextView
	table *tview.Table
	won   bool

	steps int
}

func NewCore(cols, rows int) *core {
	c := &core{
		board: func() [][]int {
			board := make([][]int, cols)
			for i := 0; i < cols; i++ {
				board[i] = make([]int, rows)
			}
			return board
		}(),
		cols: cols,
		rows: rows,
		app:  tview.NewApplication(),
	}

	txt := tview.NewTextView()
	c.txt = txt

	table := tview.NewTable()
	table.SetBorders(true)
	table.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		defer c.refreshTable()

		changed := false
		switch event.Key() {
		case tcell.KeyUp:
			changed = c.up()
		case tcell.KeyDown:
			changed = c.down()
		case tcell.KeyLeft:
			changed = c.left()
		case tcell.KeyRight:
			changed = c.right()
		case tcell.KeyEsc:
			c.app.Stop()
		case tcell.KeyRune:
			switch event.Rune() {
			case 'h':
				changed = c.left()
			case 'j':
				changed = c.down()
			case 'k':
				changed = c.up()
			case 'l':
				changed = c.right()
			}
		}

		if changed {
			c.randomInsert(1)
			c.steps++
			c.refreshTxt()
		}

		return event
	})
	c.table = table

	return c
}

func (c *core) Run() {
	// init board
	c.randomInsert(4)
	c.refreshTable()
	c.refreshTxt()

	// init layout
	rootFlex := tview.NewFlex()
	rootFlex.AddItem(tview.NewBox(), 0, 1, false)
	rootFlex.AddItem(tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewBox(), 0, 1, false).
		AddItem(c.table, 0, 1, true).
		AddItem(c.txt, 0, 1, false),
		0, 1, true)
	rootFlex.AddItem(tview.NewBox(), 0, 1, false)

	c.app.SetRoot(rootFlex, true).EnableMouse(false).Run()
}

func (c *core) randomInsert(n int) {
	// use map to achieve random
	emptyCell := make(map[int][]int)
	cnt := 0
	for i := 0; i < c.cols; i++ {
		for j := 0; j < c.rows; j++ {
			if c.board[i][j] == 0 {
				emptyCell[cnt] = []int{i, j}
				cnt++
			}
		}
	}

	// gen random 2 or 4 into n random empty cells
	for _, v := range emptyCell {
		i, j := v[0], v[1]
		c.board[i][j] = 2 << uint(rand.Intn(2))
		if n--; n == 0 {
			return
		}
	}
}

func (c *core) refreshTable() {
	table := c.table
	for i := 0; i < c.cols; i++ {
		for j := 0; j < c.rows; j++ {
			table.SetCell(i, j, c.numToCell(i, j))
		}
	}
}

func (c *core) refreshTxt() {
	text := fmt.Sprintf("%s\n\nSteps: %d", help, c.steps)

	if c.won {
		text += "\n\n" + youWin
		c.txt.SetText(text)
		return
	}

	for i := 0; i < c.cols; i++ {
		for j := 0; j < c.rows; j++ {
			if c.board[i][j] >= 2048 {
				c.won = true
				text += "\n\n" + youWin
			}
		}
	}

	c.txt.SetText(text)
}

func (c *core) numToCell(i, j int) *tview.TableCell {
	content := ""
	color := tcell.ColorDefault
	switch num := c.board[i][j]; num {
	case 0:
		content = "         "
	case 2:
		content = fmt.Sprintf("    %d    ", num)
	case 4:
		color = tcell.ColorGreen
		content = fmt.Sprintf("    %d    ", num)
	case 8:
		color = tcell.ColorSeaGreen
		content = fmt.Sprintf("    %d    ", num)
	case 16:
		color = tcell.ColorBlue
		content = fmt.Sprintf("   %d   ", num)
	case 32:
		color = tcell.Color100
		content = fmt.Sprintf("   %d   ", num)
	case 64:
		color = tcell.ColorYellow
		content = fmt.Sprintf("   %d   ", num)
	case 128:
		color = tcell.ColorYellowGreen
		content = fmt.Sprintf("   %d  ", num)
	case 256:
		color = tcell.ColorOrange
		content = fmt.Sprintf("   %d  ", num)
	case 512:
		color = tcell.ColorOrangeRed
		content = fmt.Sprintf("   %d  ", num)
	case 1024:
		color = tcell.ColorRed
		content = fmt.Sprintf("  %d  ", num)
	case 2048:
		color = tcell.ColorIndianRed
		content = fmt.Sprintf("  %d  ", num)
	case 4096:
		color = tcell.ColorGold
		content = fmt.Sprintf("  %d  ", num)
	case 8192:
		color = tcell.ColorGoldenrod
		content = fmt.Sprintf("  %d  ", num)
	default:
		color = tcell.ColorRed
		content = fmt.Sprintf(" %d", num)
	}
	cell := tview.NewTableCell(content)
	cell.SetTextColor(color)
	cell.SetAlign(tview.AlignCenter)
	return cell
}

func (c *core) up() (changed bool) {
	board := c.board

	mv := func(row int) {
		nums := make([]int, c.rows)
		for i := 0; i < c.cols; i++ {
			nums[i] = board[i][row]
		}
		merged, diff := c.merge(nums)
		if !changed {
			changed = diff
		}
		for i := 0; i < c.cols; i++ {
			board[i][row] = merged[i]
		}
	}

	for i := 0; i < 4; i++ {
		mv(i)
	}

	return
}

func (c *core) down() (changed bool) {
	board := c.board

	mv := func(row int) {
		nums := make([]int, 0, c.rows)
		for i := c.cols - 1; i >= 0; i-- {
			nums = append(nums, board[i][row])
		}
		merged, diff := c.merge(nums)
		if !changed {
			changed = diff
		}
		for i := c.cols - 1; i >= 0; i-- {
			board[i][row] = merged[c.cols-1-i]
		}
	}

	for i := 0; i < 4; i++ {
		mv(i)
	}

	return
}

func (c *core) left() (changed bool) {
	board := c.board

	mv := func(col int) {
		nums := make([]int, 0, c.cols)
		for i := 0; i < c.rows; i++ {
			nums = append(nums, board[col][i])
		}
		merged, diff := c.merge(nums)
		if !changed {
			changed = diff
		}
		for i := 0; i < c.rows; i++ {
			board[col][i] = merged[i]
		}
	}

	for i := 0; i < 4; i++ {
		mv(i)
	}

	return
}

func (c *core) right() (changed bool) {
	board := c.board

	mv := func(col int) {
		nums := make([]int, 0, c.cols)
		for i := c.rows - 1; i >= 0; i-- {
			nums = append(nums, board[col][i])
		}
		merged, diff := c.merge(nums)
		if !changed {
			changed = diff
		}
		for i := c.rows - 1; i >= 0; i-- {
			board[col][i] = merged[c.rows-1-i]
		}
	}

	for i := 0; i < 4; i++ {
		mv(i)
	}

	return
}

func (c *core) merge(nums []int) (merged []int, changed bool) {
	q := newQueue(nums)
	merged = make([]int, 0, len(nums))

	for {
		// pop num 1
		n, ok := q.popNonZero()
		if !ok {
			break
		}

		// pop num 2
		m, ok := q.popNonZero()
		if !ok {
			merged = append(merged, n)
			break
		}

		// addable
		if n == m {
			merged = append(merged, n+m)
		} else {
			merged = append(merged, n)
			q.addFirst(m)
		}
	}

	// fill zeros
	if len(merged) < len(nums) {
		merged = append(merged, make([]int, len(nums)-len(merged))...)
	}

	// check if changed
	for i := 0; i < len(nums); i++ {
		if nums[i] != merged[i] {
			changed = true
			break
		}
	}

	return
}
