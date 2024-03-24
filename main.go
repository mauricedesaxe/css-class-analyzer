package main

import "css-class-analyzer/analyzer"

func main() {
	analyzer.Analyze("./analyzer/example-pages", "./analyzer/classes.log")
}
