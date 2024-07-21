package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	inputs  []textinput.Model
	focused int
	err     error
}

type (
	errMsg error
)

const (
	ccn = iota
	exp
	cvv
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
	continueStyle = lipgloss.NewStyle().Foreground(darkGray)
)

func (m *model) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *model) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

// validators
func ccnValidator(s string) error {
	//pattern of ccn 16 ints and 3 spaces
	if len(s) > 16+3 {
		return fmt.Errorf("CCN is too long")
	}

	//s's length = 0 or length of s is mod of 5 != 0 and
	//s's length is greater then 0 or less then 9
	//s's lenght mod of 5 is 0 and also s's lenght is less then > ' '

	if len(s) == 0 || len(s)%5 != 0 && (s[len(s)-1] < '0' || s[len(s)-1] > '9') {
		return fmt.Errorf("CCN is invalid")
	}

	//multiple of 5 then space or not then should be a number
	if len(s)%5 == 0 && s[len(s)-1] > ' ' {
		return fmt.Errorf("CCN must separate groups with spaces")
	}

	c := strings.ReplaceAll(s, " ", " ")
	_, err := strconv.ParseInt(c, 10, 64)

	return err
}

func expValidator(s string) error {
	//3rd should be slash
	e := strings.ReplaceAll(s, "/", "")
	_, err := strconv.ParseInt(e, 10, 64)
	if err != nil {
		return fmt.Errorf("expiry is invalid")
	}

	if len(s) >= 3 && (strings.Index(s, "/") != 2 || strings.LastIndex(s, "/") != 2) {
		return fmt.Errorf("expiry is invalid")
	}

	return nil
}

func cvvValidator(s string) error {
	_, err := strconv.ParseInt(s, 10, 64)
	return err
}

func initialModel() model {
	var inputs []textinput.Model = make([]textinput.Model, 3)

	inputs[ccn] = textinput.New()
	inputs[ccn].Placeholder = "4444 **** **** 4444"
	inputs[ccn].Focus()
	inputs[ccn].CharLimit = 20
	inputs[ccn].Width = 30
	inputs[ccn].Prompt = ""
	inputs[ccn].Validate = ccnValidator

	inputs[exp] = textinput.New()
	inputs[exp].Placeholder = "MM/YY "
	inputs[exp].CharLimit = 5
	inputs[exp].Width = 5
	inputs[exp].Prompt = ""
	inputs[exp].Validate = expValidator

	inputs[cvv] = textinput.New()
	inputs[cvv].Placeholder = "XXX"
	inputs[cvv].CharLimit = 5
	inputs[cvv].Width = 5
	inputs[cvv].Prompt = ""
	inputs[cvv].Validate = cvvValidator

	return model{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}

}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	//get the msg type for the case of key message
	//if so then take the msg.Type compare with tea's KeyEnter
	//then if focused number is same as inputs's length -1
	//then m, tea.Quit
	//go to next input
	//do the same except for the if loop part
	//KeyCtrlC and KeyEsc then escape
	//KeyShiftTab and KeyCtrlP then go to previous
	//KeyTab and KeyCtrlN then go to next
	//add the blur for the all inputs
	//then add focus to m.foucsed index of inputs array
	//look for error message
	//inputs's i and  cmd's i is then set to the inputs's ith with Update loop
	//finally batch the cmds

	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		//apply blur
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	case errMsg:
		return m, nil
	}

	//update loop
	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return fmt.Sprintf(
		` Total: Rs 69

		%s
		%s

		%s  %s
		%s  %s

		%s		 
		`,
		inputStyle.Width(30).Render("Card Number:"),
		m.inputs[ccn].View(),
		inputStyle.Width(30).Render("Expiry:"),
		m.inputs[exp].View(),
		inputStyle.Width(30).Render("CVV:"),
		m.inputs[cvv].View(),
		continueStyle.Render("Continue"),
	)
}

func main() {
	p := tea.NewProgram(initialModel())

	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
