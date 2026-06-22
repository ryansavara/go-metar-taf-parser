package metartafparser

import "strings"

type remarkParser struct {
	supplier *remarkCommandSupplier
}

func newRemarkParser(locale Locale) *remarkParser {
	return &remarkParser{
		supplier: newRemarkCommandSupplier(locale),
	}
}

func (p *remarkParser) Parse(code string) []Remark {
	code = strings.TrimSpace(code)
	var remarks []Remark

	for code != "" {
		cmd, err := p.supplier.Get(code)
		if err != nil {
			defaultCmd := p.supplier.defaultCommand
			code, remarks, err = defaultCmd.Execute(code, remarks)
			if err != nil {
				break
			}
		} else {
			var cmdErr error
			code, remarks, cmdErr = cmd.Execute(code, remarks)
			if cmdErr != nil {
				defaultCmd := p.supplier.defaultCommand
				code, remarks, cmdErr = defaultCmd.Execute(code, remarks)
				if cmdErr != nil {
					break
				}
			}
		}
		code = strings.TrimSpace(code)
	}

	return remarks
}
