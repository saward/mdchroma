# mdchroma

Integrating [Chroma](https://github.com/alecthomas/chroma) syntax highlighter as
a [Go Markdown](https://github.com/gomarkdown/markdown) renderer.

## Install and prerequisites

This project requires and uses the
[Go Markdown](https://github.com/gomarkdown/markdown) package.

```
$ go get -u github.com/saward/mdchroma
```

_This project uses the module approach of go 1.11_ 

## Features

This renderer integrates chroma to highlight code with triple backtick notation.
It will try to use the given language when available otherwise it will try to
detect the language. If none of these two method works it will fallback to sane
defaults.

This is a fork of the [bfchroma](https://github.com/Depado/bfchroma) package,
modified to instead work with
[Go Markdown](https://github.com/gomarkdown/markdown).

## Usage

```go
output := markdown.ToHTML(md, nil, mdchroma.NewRenderer())
```

## Examples

```go
package main

import (
	"fmt"

	chroma "github.com/saward/mdchroma"
	"github.com/gomarkdown/markdown"
)

var md = "This is some sample code.\n\n```go\n" +
	`func main() {
	fmt.Println("Hi")
}
` + "```"

func main() {
	html := markdown.ToHTML([]byte(md), nil, chroma.NewRenderer())
	fmt.Println(string(html))
}
```


Will output :

```html
<p>This is some sample code.</p>
<pre style="color:#f8f8f2;background-color:#272822"><span style="color:#66d9ef">func</span> <span style="color:#a6e22e">main</span>() {
<span style="color:#a6e22e">fmt</span>.<span style="color:#a6e22e">Println</span>(<span style="color:#e6db74">&#34;Hi&#34;</span>)
}
</pre>
```
