package main

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"syscall/js"

	"veritas/pkg"
)

func createElement(doc js.Value, tag, innerHTML string) js.Value {
	element := doc.Call("createElement", tag)
	element.Set("innerHTML", innerHTML)
	return element
}

func appendChildren(parent js.Value, children ...js.Value) {
	for _, child := range children {
		parent.Call("appendChild", child)
	}
}

func calculate(doc js.Value, exp veritas.Expression) ([]int, error) {
	countItems := make(map[string]int)
	exp.Len(countItems)

	keys := make([]string, 0, len(countItems))
	for k := range countItems {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var headerContent bytes.Buffer
	for _, key := range keys {
		headerContent.WriteString(fmt.Sprintf("%v ", key))
	}

	variableMap := make(map[byte]int)
	result := make([]int, 1<<len(countItems))

	content := doc.Call("getElementById", "content")
	if content.IsNull() {
		content = initializeContent(doc)
	}

	table := createElement(doc, "table", "")
	headerRow := createElement(doc, "tr", "")
	headerCell := createElement(doc, "td", headerContent.String())
	appendChildren(headerRow, headerCell)
	appendChildren(table, headerRow)

	for index, combination := range generateCombinations(len(countItems)) {
		for idx, key := range keys {
			variableMap[key[0]] = combination[idx]
		}
		result[index] = veritas.Eval(exp, variableMap)

		rowText := fmt.Sprintf("%v = %v", combination, result[index])
		row := createElement(doc, "tr", "")
		cell := createElement(doc, "td", rowText)
		appendChildren(row, cell)
		appendChildren(table, row)
	}

	appendChildren(content, table)

	return result, nil
}

func generateCombinations(length int) [][]int {
	count := 1 << length
	combinations := make([][]int, count)

	for i := 0; i < count; i++ {
		combination := make([]int, length)
		for j := 0; j < length; j++ {
			combination[j] = (i >> j) & 1
		}
		combinations[i] = combination
	}

	return combinations
}

func parseInputExpression(input string) (veritas.Expression, error) {
	var buffer bytes.Buffer
	first := true
	args := strings.Split(input, ",")
	if len(args) < 2 {
		return nil, fmt.Errorf("Unexpected input: at least two expressions are required")
	}
	buffer.Grow(len(args) * 2)

	buffer.WriteString("(")
	for idx, arg := range args {
		if idx < len(args)-1 {
			buffer.WriteString(arg)
			if !first {
				buffer.WriteString(")")
				if len(args)-2 > idx {
					buffer.WriteString("&(")
				}
			} else {
				first = false
				buffer.WriteString("&(")
			}
		} else {
			buffer.WriteString(")@(")
			buffer.WriteString(arg)
			buffer.WriteString(")")
		}
	}
	lexer := veritas.NewLex(buffer.String())
	parser := veritas.NewParser(lexer)
	return parser.ParseExpression(veritas.LOWEST), nil
}

func initializeContent(doc js.Value) js.Value {
	content := createElement(doc, "div", "")
	content.Set("id", "content")
	content.Set("style", "padding:5px;display:table")
	doc.Get("body").Call("appendChild", content)
	return content
}

func displayResults(doc js.Value, expString string, results []int) {
	content := doc.Call("getElementById", "content")

	header := createElement(doc, "div", expString)
	appendChildren(content, header)

	contradictionFound := false
	for _, result := range results {
		if result == 0 {
			contradiction := createElement(doc, "div", "Противоречие")
			appendChildren(content, contradiction)
			contradictionFound = true
			break
		}
	}

	if !contradictionFound {
		success := createElement(doc, "div", "Тавтология")
		appendChildren(content, success)
	}
}

func main() {
	document := js.Global().Get("document")
	input := document.Call("getElementById", "conc").Get("name").String()
	exp, err := parseInputExpression(input)
	if err != nil {
		fmt.Println("Error parsing input expression:", err)
		return
	}

	results, err := calculate(document, exp)
	if err != nil {
		fmt.Println("Error calculating results:", err)
		return
	}
	displayResults(document, exp.String(), results)
}
