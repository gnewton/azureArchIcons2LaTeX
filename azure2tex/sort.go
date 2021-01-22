package main

import(
        "strings"
        )

type ByName []*Entry

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return strings.Compare(a[i].englishName, a[j].englishName) == -1 }
