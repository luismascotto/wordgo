package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Dictionary representa o dicionário de palavras para busca
type Dictionary struct {
	words map[string]bool
	trie  *TrieNode
}

// TrieNode representa um nó na árvore trie para busca de prefixos
type TrieNode struct {
	children map[rune]*TrieNode
	isWord   bool
}

// NewDictionary cria um novo dicionário a partir de um arquivo
func NewDictionary(filename string) (*Dictionary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("%s do dicionário: %w", ErrFileOpen, err)
	}
	defer file.Close()

	dict := &Dictionary{
		words: make(map[string]bool),
		trie:  &TrieNode{children: make(map[rune]*TrieNode)},
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		word := strings.TrimSpace(strings.ToUpper(scanner.Text()))
		if word != "" && len(word) >= MIN_WORD_LENGTH {
			dict.words[word] = true
			dict.insertIntoTrie(word)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("%s do dicionário: %w", ErrFileRead, err)
	}

	return dict, nil
}

// insertIntoTrie insere uma palavra na árvore trie
func (d *Dictionary) insertIntoTrie(word string) {
	node := d.trie
	for _, char := range word {
		if node.children[char] == nil {
			node.children[char] = &TrieNode{children: make(map[rune]*TrieNode)}
		}
		node = node.children[char]
	}
	node.isWord = true
}

// Contains verifica se uma palavra existe no dicionário
func (d *Dictionary) Contains(word string) bool {
	return d.words[strings.ToUpper(word)]
}

// Contains verifica se uma palavra existe no dicionário
func (d *Dictionary) ContainsUpped(uppedWord string) bool {
	return d.words[uppedWord]
}

// IsPrefix verifica se uma sequência é prefixo de alguma palavra válida
func (d *Dictionary) IsPrefix(sequence string) bool {
	node := d.trie
	for _, char := range sequence {
		if node.children[char] == nil {
			return false // Não é prefixo
		}
		node = node.children[char]
	}
	return true // É prefixo (pode ter filhos)
}

// IsWord verifica se uma sequência é uma palavra completa
func (d *Dictionary) IsWord(sequence string) bool {
	node := d.trie
	for _, char := range sequence {
		if node.children[char] == nil {
			return false // Não existe
		}
		node = node.children[char]
	}
	return node.isWord // É uma palavra completa
}

// PrintDictionaryStats imprime estatísticas do dicionário
func (d *Dictionary) PrintDictionaryStats() {
	fmt.Printf("Dicionário carregado com %d palavras\n", len(d.words))
}
