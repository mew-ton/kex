package ui

import (
	"fmt"
	"strings"

	"atomicgo.dev/cursor"
	"atomicgo.dev/keyboard"
	"atomicgo.dev/keyboard/keys"
	"github.com/pterm/pterm"
)

// Multiselect displays a custom interactive multiselect menu.
// It returns the list of selected options.
// Pre-selected options can be passed in the selected map (true = selected).
func Multiselect(prompt string, options []string, preSelected map[string]bool) ([]string, error) {
	if preSelected == nil {
		preSelected = make(map[string]bool)
	}

	// Local copy of selection state
	selected := make(map[string]bool)
	for k, v := range preSelected {
		selected[k] = v
	}

	currIdx := 0
	if len(options) == 0 {
		return nil, fmt.Errorf("no options provided")
	}

	area, err := pterm.DefaultArea.Start()
	if err != nil {
		return nil, err
	}
	defer area.Stop()

	// Initial Render
	render := func() string {
		var b strings.Builder
		b.WriteString(pterm.Cyan("? ") + pterm.DefaultInteractiveMultiselect.TextStyle.Sprint(prompt) + "\n")

		for i, opt := range options {
			// Selector
			selector := "  "
			if i == currIdx {
				selector = pterm.DefaultInteractiveMultiselect.SelectorStyle.Sprint("> ")
			}

			// Checkbox
			checked := selected[opt]
			checkbox := "[ ]"
			if checked {
				checkbox = "[x]"
			}
			// Style checkbox: Green if checked? pterm default uses Checkmark struct.
			// Let's stick to simple text or pterm style.
			if checked {
				checkbox = pterm.Green(checkbox)
			} else {
				checkbox = pterm.Gray(checkbox)
			}

			// Option Text
			text := opt
			if i == currIdx {
				text = pterm.DefaultInteractiveMultiselect.OptionStyle.Sprint(text)
			}

			b.WriteString(fmt.Sprintf("%s%s %s\n", selector, checkbox, text))
		}

		// Custom Help Text
		help := pterm.ThemeDefault.SecondaryStyle.Sprint("enter: confirm | space: toggle | up/down: move")
		b.WriteString(help)
		return b.String()
	}

	area.Update(render())

	cursor.Hide()
	defer cursor.Show()

	err = keyboard.Listen(func(key keys.Key) (stop bool, err error) {
		switch key.Code {
		case keys.CtrlC:
			return true, fmt.Errorf("keyboard interrupt")
		case keys.Enter:
			return true, nil
		case keys.Up, keys.CtrlP:
			if currIdx > 0 {
				currIdx--
				area.Update(render())
			}
		case keys.Down, keys.CtrlN:
			if currIdx < len(options)-1 {
				currIdx++
				area.Update(render())
			}
		case keys.Space:
			opt := options[currIdx]
			selected[opt] = !selected[opt]
			area.Update(render())
		}
		return false, nil // Continue listening
	})

	if err != nil {
		return nil, err
	}

	var result []string
	for _, opt := range options {
		if selected[opt] {
			result = append(result, opt)
		}
	}

	return result, nil
}
