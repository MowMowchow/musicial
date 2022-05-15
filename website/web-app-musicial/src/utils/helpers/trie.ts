export class TrieNode {
  children: { [key: string]: TrieNode };
  isWord: boolean;
  constructor() {
    this.children = {};
    this.isWord = false;
  }
}

export class Trie {
  root: TrieNode;
  suggestionResults: Set<string>;
  constructor() {
    this.root = new TrieNode();
    this.suggestionResults = new Set();
  }

  build(keys: string[]): void {
    for (let word of keys) {
      this.insert(word);
    }
    this.getSuggestions("");
  }

  insert(key: string): void {
    let prefix = this.root;
    for (let layer = 0, child; layer < key.length; layer++) {
      child = key[layer];
      if (!(child in prefix.children)) {
        prefix.children[child] = new TrieNode();
      }
      prefix = prefix.children[child];
    }
    prefix.isWord = true;
  }

  find(key: string): boolean {
    let prefix = this.root;
    for (let layer = 0, child; layer < key.length; layer++) {
      child = key[layer];
      if (!(child in prefix.children)) {
        return false;
      }
      prefix = prefix.children[child];
    }
    return prefix.isWord;
  }

  getAllPossibilities(node: TrieNode, word: string) {
    if (node.isWord) {
      this.suggestionResults.add(word);
    }
    for (let letter in node.children) {
      this.getAllPossibilities(node.children[letter], word + letter);
    }
  }

  getSuggestions(word: string) {
    this.suggestionResults.clear();
    let prefix = this.root;
    for (let letter of word) {
      if (!(letter in prefix.children)) {
        return;
      }
      prefix = prefix.children[letter];
    }
    if (Object.keys(prefix.children).length == 0) {
      return;
    }
    this.getAllPossibilities(prefix, word);
  }
}
