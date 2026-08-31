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
	pageCustomize
	pageAppearance
	pageThemes
	pageAccents
	pageTransparency
	pageLayout
	pageReference
	pageCalculator
	pageDensity
	pageReset
)

type model struct {
	teethInput      textinput.Model
	compressorInput textinput.Model

	page   page
	width  int
	height int

	cfg config

	theme  string
	accent string

	themeIndex        int
	accentIndex       int
	customIndex       int
	appearanceIndex   int
	transparencyIndex int
	layoutIndex       int
	referenceIndex    int
	calculatorIndex   int
	densityIndex      int
	resetIndex        int
}

func newInput(value string) textinput.Model {
	input := textinput.New()
	input.SetValue(value)
	input.SetWidth(14)
	input.CharLimit = 8
	input.Prompt = ""
	input.Blur()
	return input
}

func initialModel() model {
	cfg := loadConfig()

	teeth := 6
	_, _, offset := Calculate(teeth)
	compressors := AutoCompressors(offset)

	m := model{
		teethInput:      newInput(strconv.Itoa(teeth)),
		compressorInput: newInput(strconv.Itoa(compressors)),
		page:            pageMain,
		width:           100,
		height:          24,
		cfg:             cfg,
		theme:           cfg.Appearance.Theme,
		accent:          cfg.Appearance.Accent,
	}

	m.updateThemeIndex()
	m.updateAccentIndex()

	return m
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

		case pageCustomize:
			return m.updateCustomize(msg)

		case pageAppearance:
			return m.updateAppearance(msg)

		case pageThemes:
			return m.updateThemes(msg)

		case pageAccents:
			return m.updateAccents(msg)

		case pageTransparency:
			return m.updateTransparency(msg)

		case pageLayout:
			return m.updateLayout(msg)

		case pageReference:
			return m.updateReference(msg)

		case pageCalculator:
			return m.updateCalculator(msg)

		case pageDensity:
			return m.updateDensity(msg)

		case pageReset:
			return m.updateReset(msg)
		}
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Main input handling
// ─────────────────────────────────────────────────────────────

func (m model) updateMain(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
		m.stopEditing()
		m.page = pageCustomize
		m.customIndex = 0
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
		m.compressorInput, cmd = m.compressorInput.Update(msg)
		return m, cmd
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
	teeth, err := strconv.Atoi(m.teethInput.Value())
	if err != nil || teeth <= 0 {
		return
	}

	_, _, offset := Calculate(teeth)

	m.compressorInput.SetValue(
		strconv.Itoa(AutoCompressors(offset)),
	)
}

func (m model) currentTeeth() int {
	value, err := strconv.Atoi(m.teethInput.Value())
	if err != nil || value <= 0 {
		return 6
	}

	return value
}

func (m model) currentCompressors() int {
	value, err := strconv.Atoi(m.compressorInput.Value())
	if err != nil || value <= 0 {
		return 1
	}

	return value
}

// ─────────────────────────────────────────────────────────────
// Customization
// ─────────────────────────────────────────────────────────────

func (m model) updateCustomize(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Appearance",
		"Layout",
		"Reference Chart",
		"Calculator",
		"Density",
		"Reset to Defaults",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.customIndex--

		if m.customIndex < 0 {
			m.customIndex = len(items) - 1
		}

	case "down":
		m.customIndex++

		if m.customIndex >= len(items) {
			m.customIndex = 0
		}

	case "enter":
		switch m.customIndex {

		case 0:
			m.page = pageAppearance
			m.appearanceIndex = 0

		case 1:
			m.page = pageLayout
			m.layoutIndex = 0

		case 2:
			m.page = pageReference
			m.referenceIndex = 0

		case 3:
			m.page = pageCalculator
			m.calculatorIndex = 0

		case 4:
			m.page = pageDensity
			m.densityIndex =
				densityIndexFromValue(m.cfg.Density)

		case 5:
			m.page = pageReset
			m.resetIndex = 1

		case 6:
			m.page = pageMain
		}

	case "esc":
		m.page = pageMain
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Appearance
// ─────────────────────────────────────────────────────────────

func (m model) updateAppearance(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Theme",
		"Accent",
		"Background",
		"Panel Transparency",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.appearanceIndex--

		if m.appearanceIndex < 0 {
			m.appearanceIndex = len(items) - 1
		}

	case "down":
		m.appearanceIndex++

		if m.appearanceIndex >= len(items) {
			m.appearanceIndex = 0
		}

	case "enter":
		switch m.appearanceIndex {

		case 0:
			m.updateThemeIndex()
			m.page = pageThemes

		case 1:
			m.updateAccentIndex()
			m.page = pageAccents

		case 2:
			if m.cfg.Appearance.Background == "theme" {
				m.cfg.Appearance.Background = "transparent"
			} else {
				m.cfg.Appearance.Background = "theme"
			}
			m.saveSettings()

		case 3:
			m.transparencyIndex = 0
			m.page = pageTransparency

		case 4:
			m.page = pageCustomize
		}

	case "esc":
		m.page = pageCustomize
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Themes
// ─────────────────────────────────────────────────────────────

func (m model) updateThemes(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
		m.theme = themeKeys[m.themeIndex]
		m.cfg.Appearance.Theme = m.theme
		m.saveSettings()
		m.page = pageAppearance

	case "esc":
		m.page = pageAppearance
	}

	return m, nil
}

func (m model) updateAccents(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
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
		m.accent = accentKeys[m.accentIndex]
		m.cfg.Appearance.Accent = m.accent
		m.saveSettings()
		m.page = pageAppearance

	case "esc":
		m.page = pageAppearance
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Transparency
// ─────────────────────────────────────────────────────────────

func (m model) updateTransparency(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Top Bar",
		"Reference Panel",
		"Calculator Panel",
		"Bottom Bar",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.transparencyIndex--

		if m.transparencyIndex < 0 {
			m.transparencyIndex = len(items) - 1
		}

	case "down":
		m.transparencyIndex++

		if m.transparencyIndex >= len(items) {
			m.transparencyIndex = 0
		}

	case "space", "enter":
		switch m.transparencyIndex {

		case 0:
			m.cfg.Appearance.TopBarTransparent =
				!m.cfg.Appearance.TopBarTransparent

		case 1:
			m.cfg.Appearance.ReferenceTransparent =
				!m.cfg.Appearance.ReferenceTransparent

		case 2:
			m.cfg.Appearance.CalculatorTransparent =
				!m.cfg.Appearance.CalculatorTransparent

		case 3:
			m.cfg.Appearance.BottomBarTransparent =
				!m.cfg.Appearance.BottomBarTransparent

		case 4:
			m.page = pageAppearance
			return m, nil
		}

		m.saveSettings()

	case "esc":
		m.page = pageAppearance
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Layout
// ─────────────────────────────────────────────────────────────

func (m model) updateLayout(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Automatic",
		"Calculator left",
		"Calculator right",
		"Calculator only",
		"Reference only",
		"Stacked",
		"Reference width",
		"Stacked order",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.layoutIndex--

		if m.layoutIndex < 0 {
			m.layoutIndex = len(items) - 1
		}

	case "down":
		m.layoutIndex++

		if m.layoutIndex >= len(items) {
			m.layoutIndex = 0
		}

	case "left":
		switch m.layoutIndex {

		case 6:
			m.cycleReferenceWidth(-1)

		case 7:
			m.toggleOrder()
		}

	case "right":
		switch m.layoutIndex {

		case 6:
			m.cycleReferenceWidth(1)

		case 7:
			m.toggleOrder()
		}

	case "enter":
		switch m.layoutIndex {

		case 0:
			m.setLayout("automatic")

		case 1:
			m.setLayout("calculator-left")

		case 2:
			m.setLayout("calculator-right")

		case 3:
			m.setLayout("calculator-only")

		case 4:
			m.setLayout("reference-only")

		case 5:
			m.setLayout("stacked")

		case 6:
			m.cycleReferenceWidth(1)

		case 7:
			m.toggleOrder()

		case 8:
			m.page = pageCustomize
		}

	case "esc":
		m.page = pageCustomize
	}

	return m, nil
}

func (m *model) setLayout(mode string) {
	m.cfg.Layout.Mode = mode
	m.saveSettings()
}

func (m *model) toggleOrder() {
	if m.cfg.Layout.Order == "reference-first" {
		m.cfg.Layout.Order = "calculator-first"
	} else {
		m.cfg.Layout.Order = "reference-first"
	}

	m.saveSettings()
}

func (m *model) cycleReferenceWidth(delta int) {
	values := []string{
		"balanced",
		"50",
		"60",
		"70",
		"80",
	}

	current := 0

	for i, value := range values {
		if value == m.cfg.Layout.ReferenceWidth {
			current = i
			break
		}
	}

	current += delta

	for current < 0 {
		current += len(values)
	}

	current %= len(values)

	m.cfg.Layout.ReferenceWidth = values[current]
	m.saveSettings()
}

func (m model) panelWidths(width int) (int, int) {
	left := int(float64(width) * 0.60)

	switch m.cfg.Layout.ReferenceWidth {
	case "50":
		left = width / 2

	case "60":
		left = int(float64(width) * 0.60)

	case "70":
		left = int(float64(width) * 0.70)

	case "80":
		left = int(float64(width) * 0.80)

	case "balanced":
		left = int(float64(width) * 0.60)
	}

	if left < 1 {
		left = 1
	}

	right := width - left - 1

	if right < 1 {
		right = 1
	}

	return left, right
}

func (m model) effectiveLayout() string {
	mode := m.cfg.Layout.Mode

	if mode != "automatic" {
		return mode
	}

	// Automatic should be predictable:
	// use the default calculator-right layout on normal terminals,
	// and simplify to calculator-only when the terminal is too narrow.
	if m.width < 80 {
		return "calculator-only"
	}

	return "calculator-right"
}

// ─────────────────────────────────────────────────────────────
// Reference settings
// ─────────────────────────────────────────────────────────────

func (m model) updateReference(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Enabled",
		"Teeth",
		"Full Angle",
		"Half Angle",
		"Offset",
		"Compressor Value",
		"Compressors",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.referenceIndex--

		if m.referenceIndex < 0 {
			m.referenceIndex = len(items) - 1
		}

	case "down":
		m.referenceIndex++

		if m.referenceIndex >= len(items) {
			m.referenceIndex = 0
		}

	case "space", "enter":

		switch m.referenceIndex {

		case 0:
			m.cfg.Reference.Enabled =
				!m.cfg.Reference.Enabled

		case 1:
			m.cfg.Reference.Teeth =
				!m.cfg.Reference.Teeth

		case 2:
			m.cfg.Reference.FullAngle =
				!m.cfg.Reference.FullAngle

		case 3:
			m.cfg.Reference.HalfAngle =
				!m.cfg.Reference.HalfAngle

		case 4:
			m.cfg.Reference.Offset =
				!m.cfg.Reference.Offset

		case 5:
			m.cfg.Reference.CompressorValue =
				!m.cfg.Reference.CompressorValue

		case 6:
			m.cfg.Reference.Compressors =
				!m.cfg.Reference.Compressors

		case 7:
			m.page = pageCustomize
			return m, nil
		}

		m.saveSettings()

	case "esc":
		m.page = pageCustomize
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Calculator settings
// ─────────────────────────────────────────────────────────────

func (m model) updateCalculator(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Enabled",
		"Number of Teeth",
		"Compressors",
		"Full Angle",
		"Half Angle",
		"Offset",
		"Compressor Value",
		"Warnings",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.calculatorIndex--

		if m.calculatorIndex < 0 {
			m.calculatorIndex = len(items) - 1
		}

	case "down":
		m.calculatorIndex++

		if m.calculatorIndex >= len(items) {
			m.calculatorIndex = 0
		}

	case "space", "enter":

		switch m.calculatorIndex {

		case 0:
			m.cfg.Calculator.Enabled =
				!m.cfg.Calculator.Enabled

		case 1:
			m.cfg.Calculator.Teeth =
				!m.cfg.Calculator.Teeth

		case 2:
			m.cfg.Calculator.Compressors =
				!m.cfg.Calculator.Compressors

		case 3:
			m.cfg.Calculator.FullAngle =
				!m.cfg.Calculator.FullAngle

		case 4:
			m.cfg.Calculator.HalfAngle =
				!m.cfg.Calculator.HalfAngle

		case 5:
			m.cfg.Calculator.Offset =
				!m.cfg.Calculator.Offset

		case 6:
			m.cfg.Calculator.CompressorValue =
				!m.cfg.Calculator.CompressorValue

		case 7:
			m.cfg.Calculator.Warnings =
				!m.cfg.Calculator.Warnings

		case 8:
			m.page = pageCustomize
			return m, nil
		}

		m.saveSettings()

	case "esc":
		m.page = pageCustomize
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Density
// ─────────────────────────────────────────────────────────────

func (m model) updateDensity(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {
	items := []string{
		"Normal",
		"Compact",
		"Minimal",
		"Back",
	}

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up":
		m.densityIndex--

		if m.densityIndex < 0 {
			m.densityIndex = len(items) - 1
		}

	case "down":
		m.densityIndex++

		if m.densityIndex >= len(items) {
			m.densityIndex = 0
		}

	case "enter":

		switch m.densityIndex {

		case 0:
			m.cfg.Density = "normal"

		case 1:
			m.cfg.Density = "compact"

		case 2:
			m.cfg.Density = "minimal"

		case 3:
			m.page = pageCustomize
			return m, nil
		}

		m.saveSettings()

	case "esc":
		m.page = pageCustomize
	}

	return m, nil
}

func densityIndexFromValue(value string) int {
	switch value {

	case "compact":
		return 1

	case "minimal":
		return 2

	default:
		return 0
	}
}

// ─────────────────────────────────────────────────────────────
// Reset
// ─────────────────────────────────────────────────────────────

func (m model) updateReset(msg tea.KeyPressMsg) (tea.Model, tea.Cmd) {

	switch msg.String() {

	case "q", "ctrl+c":
		return m, tea.Quit

	case "up", "down", "left", "right":
		m.resetIndex ^= 1

	case "enter":

		if m.resetIndex == 0 {

			m.cfg = defaultConfig()

			m.theme =
				m.cfg.Appearance.Theme

			m.accent =
				m.cfg.Appearance.Accent

			m.updateThemeIndex()
			m.updateAccentIndex()

			m.saveSettings()
		}

		m.page = pageCustomize

	case "esc":
		m.page = pageCustomize
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// Settings helpers
// ─────────────────────────────────────────────────────────────

func (m *model) saveSettings() {
	saveConfig(m.cfg)

	m.theme = m.cfg.Appearance.Theme
	m.accent = m.cfg.Appearance.Accent
}

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
// Mouse
// ─────────────────────────────────────────────────────────────

func (m model) handleMouse(
	msg tea.MouseClickMsg,
) (tea.Model, tea.Cmd) {

	if msg.Button != tea.MouseLeft {
		return m, nil
	}

	if m.page != pageMain {
		return m, nil
	}

	width := m.width

	if width < 1 {
		width = 1
	}

	mode := m.effectiveLayout()

	// Compact / stacked modes.
	if mode == "calculator-only" ||
		mode == "reference-only" ||
		mode == "stacked" {

		teethY := 6
		compressorY := 9

		if msg.X >= 2 &&
			msg.X < 26 &&
			msg.Y == teethY {

			m.modeTeeth()
			return m, nil
		}

		if msg.X >= 2 &&
			msg.X < 26 &&
			msg.Y == compressorY {

			m.modeCompressors()
			return m, nil
		}

		return m, nil
	}

	leftWidth, _ :=
		m.panelWidths(width)

	calculatorLeft :=
		mode == "calculator-left"

	calculatorX := leftWidth + 1

	if calculatorLeft {
		calculatorX = 0
	}

	fieldX := calculatorX + 2

	if msg.X >= fieldX &&
		msg.X < fieldX+24 &&
		msg.Y == 6 {

		m.modeTeeth()
		return m, nil
	}

	if msg.X >= fieldX &&
		msg.X < fieldX+24 &&
		msg.Y == 9 {

		m.modeCompressors()
		return m, nil
	}

	return m, nil
}

// ─────────────────────────────────────────────────────────────
// View routing
// ─────────────────────────────────────────────────────────────

func (m model) View() tea.View {

	switch m.page {

	case pageCustomize:
		return m.viewCustomize()

	case pageAppearance:
		return m.viewAppearance()

	case pageThemes:
		return m.viewThemes()

	case pageAccents:
		return m.viewAccents()

	case pageTransparency:
		return m.viewTransparency()

	case pageLayout:
		return m.viewLayout()

	case pageReference:
		return m.viewReference()

	case pageCalculator:
		return m.viewCalculator()

	case pageDensity:
		return m.viewDensity()

	case pageReset:
		return m.viewReset()

	default:
		return m.viewMain()
	}
}

// ─────────────────────────────────────────────────────────────
// Menu renderer
// ─────────────────────────────────────────────────────────────

func (m model) renderMenu(
	title string,
	items []string,
	selected int,
	p palette,
) tea.View {

	width := m.width
	height := m.height

	if width < 1 {
		width = 1
	}

	if height < 1 {
		height = 1
	}

	var lines []string

	lines = append(
		lines,
		lipgloss.NewStyle().
			Bold(true).
			Foreground(p.accent).
			Render(title),
		"",
	)

	for i, item := range items {

		prefix := "  "

		if i == selected {
			prefix = "> "
		}

		style :=
			lipgloss.NewStyle().
				Foreground(p.text)

		if i == selected {

			style =
				style.
					Foreground(p.accent).
					Background(p.surface).
					Bold(true)
		}

		lines = append(
			lines,
			style.Render(prefix+item),
		)
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

	box :=
		lipgloss.NewStyle().
			Width(62).
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

	output :=
		lipgloss.Place(
			width,
			height,
			lipgloss.Center,
			lipgloss.Center,
			box,
		)

	view :=
		tea.NewView(output)

	view.AltScreen = true
	view.MouseMode =
		tea.MouseModeCellMotion

	if m.cfg.Appearance.Background != "transparent" {
		view.BackgroundColor = p.background
	}

	return view
}

// ─────────────────────────────────────────────────────────────
// Menu views
// ─────────────────────────────────────────────────────────────

func (m model) viewCustomize() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"CUSTOMIZATION",
		[]string{
			"Appearance",
			"Layout",
			"Reference Chart",
			"Calculator",
			"Density",
			"Reset to Defaults",
			"Back",
		},
		m.customIndex,
		p,
	)
}

func (m model) viewAppearance() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	background :=
		"Theme"

	if m.cfg.Appearance.Background ==
		"transparent" {

		background =
			"Transparent"
	}

	return m.renderMenu(
		"APPEARANCE",
		[]string{
			"Theme: " +
				themeDisplayName(
					m.cfg.Appearance.Theme,
				),

			"Accent: " +
				accentDisplayName(
					m.cfg.Appearance.Accent,
				),

			"Background: " +
				background,

			"Panel Transparency",

			"Back",
		},
		m.appearanceIndex,
		p,
	)
}

func (m model) viewThemes() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"THEME",
		themeNames,
		m.themeIndex,
		p,
	)
}

func (m model) viewAccents() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"ACCENT COLOR",
		accentNames,
		m.accentIndex,
		p,
	)
}

func (m model) viewTransparency() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"PANEL TRANSPARENCY",
		[]string{
			"Top Bar: " +
				boolDisplay(
					m.cfg.Appearance.TopBarTransparent,
				),

			"Reference Panel: " +
				boolDisplay(
					m.cfg.Appearance.ReferenceTransparent,
				),

			"Calculator Panel: " +
				boolDisplay(
					m.cfg.Appearance.CalculatorTransparent,
				),

			"Bottom Bar: " +
				boolDisplay(
					m.cfg.Appearance.BottomBarTransparent,
				),

			"Back",
		},
		m.transparencyIndex,
		p,
	)
}

func (m model) viewLayout() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	mode :=
		m.cfg.Layout.Mode

	if mode == "" {
		mode =
			"calculator-right"
	}

	return m.renderMenu(
		"LAYOUT",
		[]string{
			"Automatic" +
				selectedSuffix(
					mode == "automatic",
				),

			"Calculator left" +
				selectedSuffix(
					mode == "calculator-left",
				),

			"Calculator right" +
				selectedSuffix(
					mode == "calculator-right",
				),

			"Calculator only" +
				selectedSuffix(
					mode == "calculator-only",
				),

			"Reference only" +
				selectedSuffix(
					mode == "reference-only",
				),

			"Stacked" +
				selectedSuffix(
					mode == "stacked",
				),

			"Reference width: " +
				layoutWidthDisplay(
					m.cfg.Layout.ReferenceWidth,
				),

			"Stacked order: " +
				orderDisplay(
					m.cfg.Layout.Order,
				),

			"Back",
		},
		m.layoutIndex,
		p,
	)
}

func (m model) viewReference() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"REFERENCE CHART",
		[]string{
			"Enabled: " +
				boolDisplay(
					m.cfg.Reference.Enabled,
				),

			"Teeth: " +
				boolDisplay(
					m.cfg.Reference.Teeth,
				),

			"Full Angle: " +
				boolDisplay(
					m.cfg.Reference.FullAngle,
				),

			"Half Angle: " +
				boolDisplay(
					m.cfg.Reference.HalfAngle,
				),

			"Offset: " +
				boolDisplay(
					m.cfg.Reference.Offset,
				),

			"Compressor Value: " +
				boolDisplay(
					m.cfg.Reference.CompressorValue,
				),

			"Compressors: " +
				boolDisplay(
					m.cfg.Reference.Compressors,
				),

			"Back",
		},
		m.referenceIndex,
		p,
	)
}

func (m model) viewCalculator() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"CALCULATOR",
		[]string{
			"Enabled: " +
				boolDisplay(
					m.cfg.Calculator.Enabled,
				),

			"Number of Teeth: " +
				boolDisplay(
					m.cfg.Calculator.Teeth,
				),

			"Compressors: " +
				boolDisplay(
					m.cfg.Calculator.Compressors,
				),

			"Full Angle: " +
				boolDisplay(
					m.cfg.Calculator.FullAngle,
				),

			"Half Angle: " +
				boolDisplay(
					m.cfg.Calculator.HalfAngle,
				),

			"Offset: " +
				boolDisplay(
					m.cfg.Calculator.Offset,
				),

			"Compressor Value: " +
				boolDisplay(
					m.cfg.Calculator.CompressorValue,
				),

			"Warnings: " +
				boolDisplay(
					m.cfg.Calculator.Warnings,
				),

			"Back",
		},
		m.calculatorIndex,
		p,
	)
}

func (m model) viewDensity() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"DENSITY",
		[]string{
			"Normal" +
				selectedSuffix(
					m.cfg.Density == "normal",
				),

			"Compact" +
				selectedSuffix(
					m.cfg.Density == "compact",
				),

			"Minimal" +
				selectedSuffix(
					m.cfg.Density == "minimal",
				),

			"Back",
		},
		m.densityIndex,
		p,
	)
}

func (m model) viewReset() tea.View {

	p := getPalette(
		m.theme,
		m.accent,
	)

	return m.renderMenu(
		"RESET TO DEFAULTS",
		[]string{
			"YES - Reset everything",
			"NO - Cancel",
		},
		m.resetIndex,
		p,
	)
}

// ─────────────────────────────────────────────────────────────
// Main rendering
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

	teeth :=
		m.currentTeeth()

	compressors :=
		m.currentCompressors()

	full, half, offset :=
		Calculate(teeth)

	compressorValue :=
		CompressorValue(
			offset,
			compressors,
		)

	roundedValue :=
		Round3(compressorValue)

	sectionStyle :=
		lipgloss.NewStyle().
			Bold(true).
			Foreground(p.accent)

	labelStyle :=
		lipgloss.NewStyle().
			Foreground(p.muted)

	valueStyle :=
		lipgloss.NewStyle().
			Bold(true).
			Foreground(p.text)

	fieldStyle :=
		lipgloss.NewStyle().
			Width(20).
			Foreground(p.text).
			Padding(0, 1)

	if !m.cfg.Appearance.CalculatorTransparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		fieldStyle =
			fieldStyle.Background(p.surface)
	}

	focusedFieldStyle :=
		fieldStyle.
			Border(
				lipgloss.NormalBorder(),
			).
			BorderForeground(p.accent)

	// ─────────────────────────────────────────────────────────
	// Top bar
	// ─────────────────────────────────────────────────────────

	titleText :=
		"PC GEAR CALCULATOR"

	authorText :=
		"Made by Xad0"

	titleStyle :=
		lipgloss.NewStyle().
			Bold(true).
			Foreground(p.accent)

	authorStyle :=
		lipgloss.NewStyle().
			Foreground(p.muted)

	if !m.cfg.Appearance.TopBarTransparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		titleStyle =
			titleStyle.Background(p.surface)

		authorStyle =
			authorStyle.Background(p.surface)
	}

	title :=
		titleStyle.Render(titleText)

	author :=
		authorStyle.Render(authorText)

	gap :=
		width -
			lipgloss.Width(titleText) -
			lipgloss.Width(authorText)

	if gap < 1 {
		gap = 1
	}

	gapStyle :=
		lipgloss.NewStyle()

	if !m.cfg.Appearance.TopBarTransparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		gapStyle =
			gapStyle.Background(p.surface)
	}

	topContent :=
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			title,
			gapStyle.Render(
				strings.Repeat(
					" ",
					gap,
				),
			),
			author,
		)

	topRowStyle :=
		lipgloss.NewStyle().
			Width(width).
			Height(1)

	if !m.cfg.Appearance.TopBarTransparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		topRowStyle =
			topRowStyle.Background(
				p.surface,
			)
	}

	topRow :=
		topRowStyle.Render(
			topContent,
		)

	topSecondStyle :=
		lipgloss.NewStyle().
			Width(width).
			Height(1)

	if !m.cfg.Appearance.TopBarTransparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		topSecondStyle =
			topSecondStyle.Background(
				p.surface,
			)
	}

	topSecond :=
		topSecondStyle.Render(
			strings.Repeat(
				" ",
				width,
			),
		)

	topBar :=
		lipgloss.JoinVertical(
			lipgloss.Left,
			topRow,
			topSecond,
		)

	// ─────────────────────────────────────────────────────────
	// Reference
	// ─────────────────────────────────────────────────────────

	var reference strings.Builder

	if m.cfg.Reference.Enabled {

		reference.WriteString(
			sectionStyle.Render(
				"REFERENCE / 4–20 TEETH",
			),
		)

		reference.WriteString("\n\n")

		columns :=
			referenceColumns(
				m.cfg.Reference,
			)

		if len(columns) > 0 {

			reference.WriteString(
				buildReferenceHeader(
					columns,
				),
			)

			reference.WriteString("\n")

			reference.WriteString(
				buildReferenceDivider(
					columns,
				),
			)

			reference.WriteString("\n")

			for n := 4; n <= 20; n++ {

				f, h, o :=
					Calculate(n)

				c :=
					AutoCompressors(o)

				cv :=
					Round3(
						CompressorValue(
							o,
							c,
						),
					)

				reference.WriteString(
					buildReferenceRow(
						n,
						f,
						h,
						o,
						cv,
						c,
						columns,
					),
				)

				if n != 20 {
					reference.WriteString("\n")
				}
			}
		}
	}

	// ─────────────────────────────────────────────────────────
	// Calculator
	// ─────────────────────────────────────────────────────────

	var calculatorLines []string

	if m.cfg.Calculator.Enabled {

		calculatorLines =
			append(
				calculatorLines,
				sectionStyle.Render(
					"CALCULATOR",
				),
			)

		if m.cfg.Density != "minimal" {
			calculatorLines =
				append(
					calculatorLines,
					"",
				)
		}

		addGap :=
			func() {

				if m.cfg.Density ==
					"normal" {

					calculatorLines =
						append(
							calculatorLines,
							"",
						)
				}
			}

		if m.cfg.Calculator.Teeth {

			field :=
				fieldStyle.Render(
					m.teethInput.Value(),
				)

			if m.teethInput.Focused() {

				field =
					focusedFieldStyle.Render(
						m.teethInput.View(),
					)
			}

			label :=
				"NUMBER OF TEETH"

			if m.cfg.Density ==
				"minimal" {

				label = "TEETH"
			}

			calculatorLines =
				append(
					calculatorLines,
					labelStyle.Render(label),
					field,
				)

			addGap()
		}

		if m.cfg.Calculator.Compressors {

			field :=
				fieldStyle.Render(
					m.compressorInput.Value(),
				)

			if m.compressorInput.Focused() {

				field =
					focusedFieldStyle.Render(
						m.compressorInput.View(),
					)
			}

			calculatorLines =
				append(
					calculatorLines,
					labelStyle.Render(
						"COMPRESSORS",
					),
					field,
				)

			addGap()
		}

		appendValue :=
			func(
				label string,
				value string,
			) {

				calculatorLines =
					append(
						calculatorLines,
						labelStyle.Render(
							label,
						),
						valueStyle.Render(
							value,
						),
					)

				addGap()
			}

		if m.cfg.Calculator.FullAngle {

			label :=
				"FULL ANGLE"

			if m.cfg.Density ==
				"minimal" {

				label = "FULL"
			}

			appendValue(
				label,
				fmt.Sprintf(
					"%.6f°",
					full,
				),
			)
		}

		if m.cfg.Calculator.HalfAngle {

			label :=
				"HALF ANGLE"

			if m.cfg.Density ==
				"minimal" {

				label = "HALF"
			}

			appendValue(
				label,
				fmt.Sprintf(
					"%.6f°",
					half,
				),
			)
		}

		if m.cfg.Calculator.Offset {

			appendValue(
				"OFFSET",
				fmt.Sprintf(
					"%.8f",
					offset,
				),
			)
		}

		if m.cfg.Calculator.CompressorValue {

			label :=
				"COMPRESSOR VALUE"

			if m.cfg.Density ==
				"minimal" {

				label = "COMP"
			}

			appendValue(
				label,
				fmt.Sprintf(
					"%.3f",
					roundedValue,
				),
			)
		}

		if m.cfg.Calculator.Warnings {

			switch {

			case compressorValue < 0:

				calculatorLines =
					append(
						calculatorLines,
						lipgloss.NewStyle().
							Foreground(
								lipgloss.Color(
									"#f38ba8",
								),
							).
							Bold(true).
							Render(
								"OUT OF RANGE",
							),
					)

			case compressorValue > 1:

				calculatorLines =
					append(
						calculatorLines,
						lipgloss.NewStyle().
							Foreground(
								lipgloss.Color(
									"#f9e2af",
								),
							).
							Bold(true).
							Render(
								"MORE COMPRESSORS NEEDED",
							),
					)
			}
		}
	}

	calculator :=
		strings.Join(
			calculatorLines,
			"\n",
		)

	contentHeight :=
		height - 4

	if contentHeight < 1 {
		contentHeight = 1
	}

	mode :=
		m.effectiveLayout()

	// ─────────────────────────────────────────────────────────
	// Both disabled
	// ─────────────────────────────────────────────────────────

	if !m.cfg.Reference.Enabled &&
		!m.cfg.Calculator.Enabled {

		message :=
			lipgloss.NewStyle().
				Foreground(p.muted).
				Bold(true).
				Render(
					"Both the reference chart and calculator are disabled.\n\nPress C to open customization.",
				)

		return m.renderSinglePanel(
			topBar,
			message,
			width,
			height,
			contentHeight,
			p,
			false,
		)
	}

	// ─────────────────────────────────────────────────────────
	// Calculator only
	// ─────────────────────────────────────────────────────────

	if mode == "calculator-only" ||
		!m.cfg.Reference.Enabled {

		return m.renderSinglePanel(
			topBar,
			calculator,
			width,
			height,
			contentHeight,
			p,
			m.cfg.Appearance.CalculatorTransparent,
		)
	}

	// ─────────────────────────────────────────────────────────
	// Reference only
	// ─────────────────────────────────────────────────────────

	if mode == "reference-only" ||
		!m.cfg.Calculator.Enabled {

		return m.renderSinglePanel(
			topBar,
			reference.String(),
			width,
			height,
			contentHeight,
			p,
			m.cfg.Appearance.ReferenceTransparent,
		)
	}

	// ─────────────────────────────────────────────────────────
	// Stacked
	// ─────────────────────────────────────────────────────────

	if mode == "stacked" {

		// Keep the reference large enough to see most/all
		// of the table, while reserving four lines for the
		// horizontal calculator.
		calculatorHeight := 5

		if contentHeight < 8 {
			calculatorHeight =
				contentHeight / 2
		}

		if calculatorHeight < 1 {
			calculatorHeight = 1
		}

		referenceHeight :=
			contentHeight - calculatorHeight

		if referenceHeight < 1 {
			referenceHeight = 1
		}

		referencePanel :=
			m.renderPanel(
				reference.String(),
				width,
				referenceHeight,
				p,
				m.cfg.Appearance.ReferenceTransparent,
			)

		stackedCalculator :=
			m.renderStackedCalculator(
				width,
				calculatorHeight,
				p,
				full,
				half,
				offset,
				compressorValue,
				roundedValue,
				fieldStyle,
				focusedFieldStyle,
				labelStyle,
				valueStyle,
			)

		var main string

		if m.cfg.Layout.Order ==
			"calculator-first" {

			main =
				lipgloss.JoinVertical(
					lipgloss.Left,
					stackedCalculator,
					referencePanel,
				)

		} else {

			main =
				lipgloss.JoinVertical(
					lipgloss.Left,
					referencePanel,
					stackedCalculator,
				)
		}

		return m.renderMainWithBottom(
			topBar,
			main,
			width,
			p,
		)
	}

	// ─────────────────────────────────────────────────────────
	// Side by side
	// ─────────────────────────────────────────────────────────

	leftWidth, rightWidth :=
		m.panelWidths(width)

	leftText :=
		reference.String()

	rightText :=
		calculator

	leftTransparent :=
		m.cfg.Appearance.ReferenceTransparent

	rightTransparent :=
		m.cfg.Appearance.CalculatorTransparent

	if mode == "calculator-left" {

		leftText =
			calculator

		rightText =
			reference.String()

		leftTransparent =
			m.cfg.Appearance.CalculatorTransparent

		rightTransparent =
			m.cfg.Appearance.ReferenceTransparent

		leftWidth, rightWidth =
			rightWidth, leftWidth
	}

	leftPanel :=
		m.renderPanel(
			leftText,
			leftWidth,
			contentHeight,
			p,
			leftTransparent,
		)

	rightPanel :=
		m.renderPanel(
			rightText,
			rightWidth,
			contentHeight,
			p,
			rightTransparent,
		)

	dividerLines :=
		make(
			[]string,
			contentHeight,
		)

	for i := range dividerLines {

		dividerLines[i] =
			"│"
	}

	divider :=
		lipgloss.NewStyle().
			Width(1).
			Height(contentHeight).
			Foreground(p.border).
			Render(
				strings.Join(
					dividerLines,
					"\n",
				),
			)

	main :=
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			leftPanel,
			divider,
			rightPanel,
		)

	return m.renderMainWithBottom(
		topBar,
		main,
		width,
		p,
	)
}

// ─────────────────────────────────────────────────────────────
// Horizontal stacked calculator
// ─────────────────────────────────────────────────────────────

func (m model) renderStackedCalculator(
	width int,
	height int,
	p palette,
	full float64,
	half float64,
	offset float64,
	compressorValue float64,
	roundedValue float64,
	fieldStyle lipgloss.Style,
	focusedFieldStyle lipgloss.Style,
	labelStyle lipgloss.Style,
	valueStyle lipgloss.Style,
) string {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(p.accent).
		Render("CALCULATOR")

	type calcItem struct {
		label string
		value string
		input bool
		kind  string
	}

	var items []calcItem

	if m.cfg.Calculator.Teeth {
		items = append(items, calcItem{
			label: "TEETH",
			input: true,
			kind:  "teeth",
		})
	}

	if m.cfg.Calculator.Compressors {
		items = append(items, calcItem{
			label: "COMPRESSORS",
			input: true,
			kind:  "compressors",
		})
	}

	if m.cfg.Calculator.FullAngle {
		items = append(items, calcItem{
			label: "FULL",
			value: fmt.Sprintf("%.6f°", full),
		})
	}

	if m.cfg.Calculator.HalfAngle {
		items = append(items, calcItem{
			label: "HALF",
			value: fmt.Sprintf("%.6f°", half),
		})
	}

	if m.cfg.Calculator.Offset {
		items = append(items, calcItem{
			label: "OFFSET",
			value: fmt.Sprintf("%.8f", offset),
		})
	}

	if m.cfg.Calculator.CompressorValue {
		items = append(items, calcItem{
			label: "COMP",
			value: fmt.Sprintf("%.3f", roundedValue),
		})
	}

	if len(items) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			labelStyle.Render("No calculator elements enabled."),
		)

		return m.renderPanel(
			content,
			width,
			height,
			p,
			m.cfg.Appearance.CalculatorTransparent,
		)
	}

	var blocks []string

	for _, item := range items {
		var value string

		switch item.kind {
		case "teeth":
			if m.teethInput.Focused() {
				value = focusedFieldStyle.Render(
					m.teethInput.View(),
				)
			} else {
				value = fieldStyle.Render(
					m.teethInput.Value(),
				)
			}

		case "compressors":
			if m.compressorInput.Focused() {
				value = focusedFieldStyle.Render(
					m.compressorInput.View(),
				)
			} else {
				value = fieldStyle.Render(
					m.compressorInput.Value(),
				)
			}

		default:
			value = valueStyle.Render(item.value)
		}

		block := lipgloss.JoinVertical(
			lipgloss.Left,
			lipgloss.NewStyle().
				Bold(true).
				Foreground(p.muted).
				Render(item.label),
			value,
		)

		blocks = append(blocks, block)
	}

	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		blocks...,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		row,
	)

	if m.cfg.Calculator.Warnings {
		switch {
		case compressorValue < 0:
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				content,
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("#f38ba8")).
					Bold(true).
					Render("OUT OF RANGE"),
			)

		case compressorValue > 1:
			content = lipgloss.JoinVertical(
				lipgloss.Left,
				content,
				lipgloss.NewStyle().
					Foreground(lipgloss.Color("#f9e2af")).
					Bold(true).
					Render("MORE COMPRESSORS NEEDED"),
			)
		}
	}

	return m.renderPanel(
		content,
		width,
		height,
		p,
		m.cfg.Appearance.CalculatorTransparent,
	)
}

// ─────────────────────────────────────────────────────────────
// Shared panel helpers
// ─────────────────────────────────────────────────────────────

func (m model) renderMainWithBottom(
	topBar string,
	main string,
	width int,
	p palette,
) tea.View {

	bottomBar :=
		m.renderBottomBar(
			width,
			p,
		)

	output :=
		strings.Join(
			[]string{
				topBar,
				main,
				bottomBar,
			},
			"\n",
		)

	view :=
		tea.NewView(output)

	view.AltScreen = true
	view.MouseMode =
		tea.MouseModeCellMotion

	if m.cfg.Appearance.Background != "transparent" {
		view.BackgroundColor = p.background
	}

	return view
}

func (m model) renderPanel(
	content string,
	width int,
	height int,
	p palette,
	transparent bool,
) string {

	if width < 1 {
		width = 1
	}

	if height < 1 {
		height = 1
	}

	// Width passed into this helper represents the TOTAL rendered
	// panel width. Account for the panel's horizontal padding so
	// the rendered result does not grow beyond the space reserved
	// for it by the layout code.
	style :=
		lipgloss.NewStyle().
			Width(width).
			Height(height).
			Padding(0, 2).
			Foreground(p.text)

	if !transparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		style =
			style.Background(p.panel)
	}

	return style.Render(content)
}

func (m model) renderSinglePanel(
	topBar string,
	content string,
	width int,
	height int,
	contentHeight int,
	p palette,
	transparent bool,
) tea.View {

	panel :=
		m.renderPanel(
			content,
			width,
			contentHeight,
			p,
			transparent,
		)

	bottomBar :=
		m.renderBottomBar(
			width,
			p,
		)

	output :=
		strings.Join(
			[]string{
				topBar,
				panel,
				bottomBar,
			},
			"\n",
		)

	view :=
		tea.NewView(output)

	view.AltScreen = true
	view.MouseMode =
		tea.MouseModeCellMotion

	if m.cfg.Appearance.Background != "transparent" {
		view.BackgroundColor = p.background
	}

	return view
}

func (m model) renderBottomBar(
	width int,
	p palette,
) string {

	text :=
		"[E] Edit teeth    [R] Edit compressors    [C] Customize    [Q] Quit"

	style :=
		lipgloss.NewStyle().
			Width(width).
			Height(2).
			Foreground(p.muted)

	if !m.cfg.Appearance.BottomBarTransparent &&
		m.cfg.Appearance.Background !=
			"transparent" {

		style =
			style.Background(p.surface)
	}

	return style.Render(text)
}

// ─────────────────────────────────────────────────────────────
// Reference table
// ─────────────────────────────────────────────────────────────

type refColumn int

const (
	refTeeth refColumn = iota
	refFull
	refHalf
	refOffset
	refCompValue
	refCompressors
)

func referenceColumns(
	cfg referenceConfig,
) []refColumn {

	var columns []refColumn

	if cfg.Teeth {
		columns =
			append(columns, refTeeth)
	}

	if cfg.FullAngle {
		columns =
			append(columns, refFull)
	}

	if cfg.HalfAngle {
		columns =
			append(columns, refHalf)
	}

	if cfg.Offset {
		columns =
			append(columns, refOffset)
	}

	if cfg.CompressorValue {
		columns =
			append(columns, refCompValue)
	}

	if cfg.Compressors {
		columns =
			append(columns, refCompressors)
	}

	return columns
}

func buildReferenceHeader(
	columns []refColumn,
) string {

	var parts []string

	for _, column := range columns {

		switch column {

		case refTeeth:
			parts =
				append(parts, "TEETH")

		case refFull:
			parts =
				append(parts, "FULL ANGLE")

		case refHalf:
			parts =
				append(parts, "HALF ANGLE")

		case refOffset:
			parts =
				append(parts, "OFFSET")

		case refCompValue:
			parts =
				append(parts, "COMP VALUE")

		case refCompressors:
			parts =
				append(parts, "COMPRESSORS")
		}
	}

	return strings.Join(
		parts,
		"    ",
	)
}

func buildReferenceDivider(
	columns []refColumn,
) string {

	var parts []string

	for _, column := range columns {

		switch column {

		case refTeeth:
			parts =
				append(parts, "─────")

		case refFull:
			parts =
				append(parts, "──────────")

		case refHalf:
			parts =
				append(parts, "──────────")

		case refOffset:
			parts =
				append(parts, "────────")

		case refCompValue:
			parts =
				append(parts, "──────────")

		case refCompressors:
			parts =
				append(parts, "───────────")
		}
	}

	return strings.Join(
		parts,
		"    ",
	)
}

func buildReferenceRow(
	teeth int,
	full float64,
	half float64,
	offset float64,
	compValue float64,
	compressors int,
	columns []refColumn,
) string {

	var parts []string

	for _, column := range columns {

		switch column {

		case refTeeth:
			parts =
				append(
					parts,
					fmt.Sprintf(
						"%5d",
						teeth,
					),
				)

		case refFull:
			parts =
				append(
					parts,
					fmt.Sprintf(
						"%9.3f°",
						full,
					),
				)

		case refHalf:
			parts =
				append(
					parts,
					fmt.Sprintf(
						"%9.3f°",
						half,
					),
				)

		case refOffset:
			parts =
				append(
					parts,
					fmt.Sprintf(
						"%8.3f",
						offset,
					),
				)

		case refCompValue:
			parts =
				append(
					parts,
					fmt.Sprintf(
						"%10.3f",
						compValue,
					),
				)

		case refCompressors:
			parts =
				append(
					parts,
					fmt.Sprintf(
						"%10d",
						compressors,
					),
				)
		}
	}

	return strings.Join(
		parts,
		"    ",
	)
}

// ─────────────────────────────────────────────────────────────
// Display helpers
// ─────────────────────────────────────────────────────────────

func boolDisplay(value bool) string {
	if value {
		return "[ON]"
	}

	return "[OFF]"
}

func selectedSuffix(value bool) string {
	if value {
		return " [SELECTED]"
	}

	return ""
}

func layoutWidthDisplay(value string) string {

	switch value {

	case "50":
		return "50%"

	case "60":
		return "60%"

	case "70":
		return "70%"

	case "80":
		return "80%"

	default:
		return "Balanced"
	}
}

func orderDisplay(value string) string {
	if value == "calculator-first" {
		return "Calculator first"
	}

	return "Reference first"
}

func themeDisplayName(key string) string {

	for i, value := range themeKeys {

		if value == key {
			return themeNames[i]
		}
	}

	return "Catppuccin Mocha"
}

func accentDisplayName(key string) string {

	for i, value := range accentKeys {

		if value == key {
			return accentNames[i]
		}
	}

	return "Green"
}
