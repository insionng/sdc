package main

import (
	"fmt"
	"strings"
)

type trieTree struct {
	key        rune
	word       []rune
	isword     bool
	childNodes map[rune]*trieTree
}

func (t *trieTree) Insert(word []rune, index int) {
	//fmt.Println("插入的关键词是", string(word))
	//获取第index个字
	key := word[index]
	//fmt.Println("保存的key是:", string(key))
	index++
	//如果以前保存这个字,则找下一个节点
	if _, found := t.childNodes[key]; found {
		t.childNodes[key].Insert(word, index)
	} else {
		//如果没有保存过,则保存这个key
		trie := new(trieTree)
		t.childNodes[key] = trie
		trie.key = key
		//判断是否到词的最后
		if index == len(word) {
			trie.word = word
			trie.isword = true
		} else {
			trie.childNodes = make(map[rune]*trieTree)
			trie.Insert(word, index)
		}
	}

}

func (t *trieTree) Replace(msg []rune) string {
	//保存根节点指针
	var root *trieTree = t
	for index, s := range msg {
		if _, found := t.childNodes[s]; found {
			if t.childNodes[s].isword {
				wordlen := len(t.childNodes[s].word)
				copy(msg[index-wordlen+1:index+1], []rune(strings.Repeat("*", wordlen)))
				fmt.Println("找到过滤词", string(t.childNodes[s].word), "===>", string(msg[index-wordlen+1:index+1]))
				continue
			}
			t = t.childNodes[s]
		} else {
			t = root

		}
	}
	return string(msg)
}

func main() {
	//初始化根节点
	trie := new(trieTree)
	trie.childNodes = make(map[rune]*trieTree)

	var wordList = []string{"我擦", "我草", "我操", "fuck", "傻逼", "SB"}

	var msg string = "昨天见到你妈,逼我买房子,你说我怎么办呢,fuck了,我擦"

	for _, word := range wordList {
		//子节点开始保存数据
		word = strings.ToLower(word)
		trie.Insert([]rune(word), 0)
	}

	msg = trie.Replace([]rune(msg))
	fmt.Println(msg)
	//fmt.Printf("%s\n%v\n", trie.childNodes, trie.childNodes)
	//printMap(trie)
}

func printMap(trie *trieTree) {
	for _, v := range trie.childNodes {
		fmt.Println("key:", string(v.key), "word:", string(v.word), "isword:", v.isword)
		printMap(v)
	}

}
