package main

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

type page int

const (
	pageMain page = iota
	pageThemes
	pageAccents
)

func newInput(value string) textinput.Model {
	input := textinput.New()
	input.SetValue(value)
	input.SetWidth(14)
	input.CharLimit = 8
	input.Prompt = ""
	input.Blur()

	return input
}

type model struct {
	teethInput      textinput.Model
	compressorInput textinput.Model

	page   page
	width  int
	height int

	theme  string
	accent string

	themeIndex  int
	accentIndex int
}

func initialModel() model {
	cfg := loadConfig()

	teeth := 6

	_, _, offset := Calculate(teeth)
	compressors := AutoCompressors(offset)

	return model{
		teethInput: newInput(
			strconv.Itoa(teeth),
		),
		compressorInput: newInput(
			strconv.Itoa(compressors),
		),
		page:   pageMain,
		width:  100,
		height: 24,
		theme:  cfg.Theme,
		accent: cfg.Accent,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.MouseClickMsg:
		return m.handleMouse(msg)

	case tea.KeyPressMsg:
		switch m.page {

		case pageMain:
			return m.updateMain(msg)

		case pageThemes:
			return m.updateThemes(msg)

		case pageAccents:
			return m.updateAccents(msg)
		}
	}

	return m, nil
}

func (m model) updateMain(
	msg tea.KeyPressMsg,
) (tea.Model, tea.Cmd) {

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "e":
		m.modeTeeth()
		return m, nil

	case "r":
		m.modeCompressors()
		return m, nil

	case "c":
		m.page = pageThemes
		m.teethInput.Blur()
		m.compressorInput.Blur()
		m.updateThemeIndex()
		return m, nil

	case "esc":
		m.stopEditing()
		return m, nil
	}

	var cmd tea.Cmd

	if m.teethInput.Focused() {
		m.teethInput, cmd = m.teethInput.Update(msg)

		if m.teethInput.Value() != "" {
			m.updateAutomaticCompressor()
		}

		return m, cmd
	}

	if m.compressorInput.Focused() {
		m.compressorInput, cmd =
			m.compressorInput.Update(msg)

		return m, cmd
	}

	return m, nil
}

func (m model) updateThemes(
	msg tea.KeyPressMsg,
) (tea.Model, tea.Cmd) {

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.themeIndex--

		if m.themeIndex < 0 {
			m.themeIndex = len(themeKeys) - 1
		}

	case "down":
		m.themeIndex++

		if m.themeIndex >= len(themeKeys) {
			m.themeIndex = 0
		}

	case "enter":
		m.page = pageAccents
		m.updateAccentIndex()

	case "esc", "c":
		m.page = pageMain
	}

	return m, nil
}

func (m model) updateAccents(
	msg tea.KeyPressMsg,
) (tea.Model, tea.Cmd) {

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.accentIndex--

		if m.accentIndex < 0 {
			m.accentIndex = len(accentKeys) - 1
		}

	case "down":
		m.accentIndex++

		if m.accentIndex >= len(accentKeys) {
			m.accentIndex = 0
		}

	case "enter":
		m.applyCustomization()

	case "esc":
		m.page = pageThemes
	}

	return m, nil
}

func (m *model) modeTeeth() {
	m.compressorInput.Blur()
	m.teethInput.Focus()
	m.teethInput.CursorEnd()
}

func (m *model) modeCompressors() {
	m.teethInput.Blur()
	m.compressorInput.Focus()
	m.compressorInput.CursorEnd()
}

func (m *model) stopEditing() {
	m.teethInput.Blur()
	m.compressorInput.Blur()
}

func (m *model) updateAutomaticCompressor() {
	teeth, err := strconv.Atoi(
		m.teethInput.Value(),
	)

	if err != nil || teeth <= 0 {
		return
	}

	_, _, offset := Calculate(teeth)

	compressors := AutoCompressors(offset)

	m.compressorInput.SetValue(
		strconv.Itoa(compressors),
	)
}

func (m model) currentTeeth() int {
	value, err := strconv.Atoi(
		m.teethInput.Value(),
	)

	if err != nil || value <= 0 {
		return 6
	}

	return value
}

func (m model) currentCompressors() int {
	value, err := strconv.Atoi(
		m.compressorInput.Value(),
	)

	if err != nil || value <= 0 {
		return 1
	}

	return value
}

// ─────────────────────────────────────────────────────────────
// Mouse
// ─────────────────────────────────────────────────────────────

func (m model) handleMouse(
	msg tea.MouseClickMsg,
) (tea.Model, tea.Cmd) {

	if msg.Button != tea.MouseLeft {
		return m, nil
	}

	switch m.page {

	case pageMain:
		return m.handleMainMouse(msg)

	case pageThemes:
		return m.handleThemeMouse(msg)

	case pageAccents:
		return m.handleAccentMouse(msg)
	}

	return m, nil
}

func (m model) handleMainMouse(
	msg tea.MouseClickMsg,
) (tea.Model, tea.Cmd) {

	width := m.width

	if width < 1 {
		width = 1
	}

	// Compact layout.
	if width < 110 {

		fieldX := 2

		mainTop := 2

		teethY := mainTop + 4
		compressorY := mainTop + 7

		if msg.X >= fieldX &&
			msg.X < fieldX+24 &&
			msg.Y == teethY {

			m.modeTeeth()
			return m, nil
		}

		if msg.X >= fieldX &&
			msg.X < fieldX+24 &&
			msg.Y == compressorY {

			m.modeCompressors()
			return m, nil
		}

		m.stopEditing()

		return m, nil
	}

	// Normal layout.
	leftWidth := int(
		float64(width) * 0.60,
	)

	const referencePanelMinWidth = 78

	if leftWidth < referencePanelMinWidth {
		leftWidth = referencePanelMinWidth
	}

	calculatorX := leftWidth + 1
	fieldX := calculatorX + 2

	fieldWidth := 24

	mainTop := 2

	teethY := mainTop + 4
	compressorY := mainTop + 7

	if msg.X >= fieldX &&
		msg.X < fieldX+fieldWidth &&
		msg.Y == teethY {

		m.modeTeeth()
		return m, nil
	}

	if msg.X >= fieldX &&
		msg.X < fieldX+fieldWidth &&
		msg.Y == compressorY {

		m.modeCompressors()
		return m, nil
	}

	m.stopEditing()

	return m, nil
}

func (m model) handleThemeMouse(
	msg tea.MouseClickMsg,
) (tea.Model, tea.Cmd) {

	boxWidth := 58
	boxHeight := 17

	boxX := (m.width - boxWidth) / 2
	boxY := (m.height - boxHeight) / 2

	firstOptionY := boxY + 7

	for i := range themeKeys {

		optionY := firstOptionY + i

		if msg.X >= boxX+1 &&
			msg.X < boxX+boxWidth-1 &&
			msg.Y == optionY {

			m.themeIndex = i
			m.page = pageAccents
			m.updateAccentIndex()

			return m, nil
		}
	}

	if msg.X < boxX ||
		msg.X >= boxX+boxWidth ||
		msg.Y < boxY ||
		msg.Y >= boxY+boxHeight {

		m.page = pageMain
	}

	return m, nil
}

func (m model) handleAccentMouse(
	msg tea.MouseClickMsg,
) (tea.Model, tea.Cmd) {

	boxWidth := 58
	boxHeight := 20

	boxX := (m.width - boxWidth) / 2
	boxY := (m.height - boxHeight) / 2

	firstOptionY := boxY + 7

	for i := range accentKeys {

		optionY := firstOptionY + i

		if msg.X >= boxX+1 &&
			msg.X < boxX+boxWidth-1 &&
			msg.Y == optionY {

			m.accentIndex = i
			m.applyCustomization()

			return m, nil
		}
	}

	if msg.X < boxX ||
		msg.X >= boxX+boxWidth ||
		msg.Y < boxY ||
		msg.Y >= boxY+boxHeight {

		m.page = pageThemes
	}

	return m, nil
}

func (m *model) applyCustomization() {
	m.theme = themeKeys[m.themeIndex]
	m.accent = accentKeys[m.accentIndex]

	saveConfig(config{
		Theme:  m.theme,
		Accent: m.accent,
	})

	m.page = pageMain
}

// ─────────────────────────────────────────────────────────────
// Theme helpers
// ─────────────────────────────────────────────────────────────

func (m *model) updateThemeIndex() {
	for i, key := range themeKeys {

		if key == m.theme {
			m.themeIndex = i
			return
		}
	}

	m.themeIndex = 0
}

func (m *model) updateAccentIndex() {
	for i, key := range accentKeys {

		if key == m.accent {
			m.accentIndex = i
			return
		}
	}

	m.accentIndex = 0
}

// ─────────────────────────────────────────────────────────────
// View selection
// ─────────────────────────────────────────────────────────────

func (m model) View() tea.View {

	switch m.page {

	case pageThemes:
		return m.viewThemes()

	case pageAccents:
		return m.viewAccents()

	default:
		return m.viewMain()
	}
}

// ─────────────────────────────────────────────────────────────
// Main view
// ─────────────────────────────────────────────────────────────

func (m model) viewMain() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	width := m.width
	height := m.height

	if width < 1 {
		width = 1
	}

	if height < 1 {
		height = 1
	}

	teeth := m.currentTeeth()
	compressors := m.currentCompressors()

	full, half, offset := Calculate(teeth)

	compressorValue := CompressorValue(
		offset,
		compressors,
	)

	roundedValue := Round3(
		compressorValue,
	)

	sectionStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(p.accent)

	labelStyle := lipgloss.NewStyle().
		Foreground(p.muted)

	valueStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(p.text)

	fieldStyle := lipgloss.NewStyle().
		Width(20).
		Background(p.surface).
		Foreground(p.text).
		Padding(0, 1)

	focusedFieldStyle := fieldStyle.
		Border(lipgloss.NormalBorder()).
		BorderForeground(p.accent)

	// ───────── TOP BAR ─────────

	titleText := "PC GEAR CALCULATOR"
	authorText := "Made by Xad0"

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(p.accent).
		Background(p.surface).
		Render(titleText)

	author := lipgloss.NewStyle().
		Foreground(p.muted).
		Background(p.surface).
		Render(authorText)

	gap := width -
		lipgloss.Width(titleText) -
		lipgloss.Width(authorText)

	if gap < 1 {
		gap = 1
	}

	gapText := lipgloss.NewStyle().
		Background(p.surface).
		Render(
			strings.Repeat(" ", gap),
		)

	topContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		title,
		gapText,
		author,
	)

	topRow := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Background(p.surface).
		Render(topContent)

	topSecondRow := lipgloss.NewStyle().
		Width(width).
		Height(1).
		Background(p.surface).
		Render(
			strings.Repeat(" ", width),
		)

	topBar := lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		topSecondRow,
	)

	// ───────── REFERENCE CHART ─────────

	var reference strings.Builder

	reference.WriteString(
		sectionStyle.Render(
			"REFERENCE / 4–20 TEETH",
		),
	)

	reference.WriteString("\n\n")

	reference.WriteString(
		"TEETH    FULL ANGLE    HALF ANGLE      OFFSET    COMP VALUE    COMPRESSORS",
	)

	reference.WriteString("\n")

	reference.WriteString(
		"────────────────────────────────────────────────────────────────────────",
	)

	reference.WriteString("\n")

	for n := 4; n <= 20; n++ {

		f, h, o := Calculate(n)

		c := AutoCompressors(o)

		cv := Round3(
			CompressorValue(o, c),
		)

		reference.WriteString(
			fmt.Sprintf(
				"%5d    %9.3f°    %9.3f°    %8.3f    %10.3f    %10d",
				n,
				f,
				h,
				o,
				cv,
				c,
			),
		)

		if n != 20 {
			reference.WriteString("\n")
		}
	}

	// ───────── CALCULATOR ─────────

	teethField := fieldStyle.Render(
		m.teethInput.Value(),
	)

	if m.teethInput.Focused() {
		teethField = focusedFieldStyle.Render(
			m.teethInput.View(),
		)
	}

	compressorField := fieldStyle.Render(
		m.compressorInput.Value(),
	)

	if m.compressorInput.Focused() {
		compressorField = focusedFieldStyle.Render(
			m.compressorInput.View(),
		)
	}

	calculatorLines := []string{
		sectionStyle.Render("CALCULATOR"),
		"",
		labelStyle.Render("NUMBER OF TEETH"),
		teethField,
		"",
		labelStyle.Render("COMPRESSORS"),
		compressorField,
		"",
		labelStyle.Render("FULL ANGLE"),
		valueStyle.Render(
			fmt.Sprintf("%.6f°", full),
		),
		"",
		labelStyle.Render("HALF ANGLE"),
		valueStyle.Render(
			fmt.Sprintf("%.6f°", half),
		),
		"",
		labelStyle.Render("OFFSET"),
		valueStyle.Render(
			fmt.Sprintf("%.8f", offset),
		),
		"",
		labelStyle.Render("COMPRESSOR VALUE"),
		valueStyle.Render(
			fmt.Sprintf("%.3f", roundedValue),
		),
	}

	if compressorValue < 0 {

		calculatorLines = append(
			calculatorLines,
			"",
			lipgloss.NewStyle().
				Foreground(
					lipgloss.Color("#f38ba8"),
				).
				Bold(true).
				Render(
					"OUT OF RANGE",
				),
		)

	} else if compressorValue > 1 {

		calculatorLines = append(
			calculatorLines,
			"",
			lipgloss.NewStyle().
				Foreground(
					lipgloss.Color("#f9e2af"),
				).
				Bold(true).
				Render(
					"MORE COMPRESSORS NEEDED",
				),
		)
	}

	calculator := strings.Join(
		calculatorLines,
		"\n",
	)

	// ───────── MAIN PANELS ─────────

	contentHeight := height - 4

	if contentHeight < 1 {
		contentHeight = 1
	}

	// Normal two-panel layout.
	if width >= 110 {

		leftWidth := int(
			float64(width) * 0.60,
		)

		const referencePanelMinWidth = 78

		if leftWidth < referencePanelMinWidth {
			leftWidth = referencePanelMinWidth
		}

		rightWidth := width -
			leftWidth -
			1

		if rightWidth < 1 {
			rightWidth = 1
		}

		leftPanel := lipgloss.NewStyle().
			Width(leftWidth).
			Height(contentHeight).
			Padding(1, 2).
			Background(p.background).
			Foreground(p.text).
			Render(
				reference.String(),
			)

		dividerLines := make(
			[]string,
			contentHeight,
		)

		for i := range dividerLines {
			dividerLines[i] = "│"
		}

		divider := lipgloss.NewStyle().
			Width(1).
			Height(contentHeight).
			Foreground(p.border).
			Render(
				strings.Join(
					dividerLines,
					"\n",
				),
			)

		rightPanel := lipgloss.NewStyle().
			Width(rightWidth).
			Height(contentHeight).
			Padding(1, 2).
			Background(p.panel).
			Render(
				calculator,
			)

		main := lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftPanel,
			divider,
			rightPanel,
		)

		// ───────── BOTTOM BAR ─────────

		bottomText :=
			"[E] Edit teeth    [R] Edit compressors    [C] Customize    [Q] Quit"

		bottomBar := lipgloss.NewStyle().
			Width(width).
			Height(2).
			Background(p.surface).
			Foreground(p.muted).
			Render(bottomText)

		output := strings.Join(
			[]string{
				topBar,
				main,
				bottomBar,
			},
			"\n",
		)

		view := tea.NewView(output)

		view.AltScreen = true
		view.MouseMode = tea.MouseModeCellMotion

		return view
	}

	// ───────── COMPACT LAYOUT ─────────

	compactCalculator := lipgloss.NewStyle().
		Width(width).
		Height(contentHeight).
		Padding(1, 2).
		Background(p.panel).
		Foreground(p.text).
		Render(
			calculator,
		)

	bottomText :=
		"[E] Edit teeth    [R] Edit compressors    [C] Customize    [Q] Quit"

	bottomBar := lipgloss.NewStyle().
		Width(width).
		Height(2).
		Background(p.surface).
		Foreground(p.muted).
		Render(bottomText)

	output := strings.Join(
		[]string{
			topBar,
			compactCalculator,
			bottomBar,
		},
		"\n",
	)

	view := tea.NewView(output)

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

// ─────────────────────────────────────────────────────────────
// Theme screen
// ─────────────────────────────────────────────────────────────

func (m model) viewThemes() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	var lines []string

	lines = append(
		lines,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(p.accent).
			Render(
				"CUSTOMIZATION",
			),
	)

	lines = append(
		lines,
		"",
		"BASE THEME",
		"",
	)

	for i, name := range themeNames {

		if i == m.themeIndex {

			lines = append(
				lines,
				lipgloss.NewStyle().
					Foreground(p.accent).
					Background(p.surface).
					Bold(true).
					Render(
						"> "+name,
					),
			)

		} else {

			lines = append(
				lines,
				lipgloss.NewStyle().
					Foreground(p.text).
					Render(
						"  "+name,
					),
			)
		}
	}

	lines = append(
		lines,
		"",
		lipgloss.NewStyle().
			Foreground(p.muted).
			Render(
				"[↑↓] Select    [ENTER] Choose    [ESC] Back",
			),
	)

	box := lipgloss.NewStyle().
		Width(58).
		Padding(2, 3).
		Background(p.panel).
		Border(lipgloss.NormalBorder()).
		BorderForeground(p.border).
		Render(
			strings.Join(
				lines,
				"\n",
			),
		)

	output := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	view := tea.NewView(output)

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}

// ─────────────────────────────────────────────────────────────
// Accent screen
// ─────────────────────────────────────────────────────────────

func (m model) viewAccents() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	var lines []string

	lines = append(
		lines,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(p.accent).
			Render(
				"CUSTOMIZATION",
			),
	)

	lines = append(
		lines,
		"",
		lipgloss.NewStyle().
			Foreground(p.muted).
			Render(
				themeNames[m.themeIndex]+
					" / ACCENT COLOR",
			),
		"",
	)

	for i, name := range accentNames {

		accentColor :=
			accentColors[accentKeys[i]]

		if i == m.accentIndex {

			lines = append(
				lines,
				lipgloss.NewStyle().
					Foreground(accentColor).
					Background(p.surface).
					Bold(true).
					Render(
						"> "+name,
					),
			)

		} else {

			lines = append(
				lines,
				lipgloss.NewStyle().
					Foreground(accentColor).
					Render(
						"  "+name,
					),
			)
		}
	}

	lines = append(
		lines,
		"",
		lipgloss.NewStyle().
			Foreground(p.muted).
			Render(
				"[↑↓] Select    [ENTER] Apply    [ESC] Back",
			),
	)

	box := lipgloss.NewStyle().
		Width(58).
		Padding(2, 3).
		Background(p.panel).
		Border(lipgloss.NormalBorder()).
		BorderForeground(p.border).
		Render(
			strings.Join(
				lines,
				"\n",
			),
		)

	output := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)

	view := tea.NewView(output)

	view.AltScreen = true
	view.MouseMode = tea.MouseModeCellMotion

	return view
}
