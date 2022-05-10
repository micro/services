# Go-Wordle

This is a CLI version of popular game [Wordle](https://www.powerlanguage.co.uk/wordle/) powered by golang.

<img src="https://raw.githubusercontent.com/XiaoMengXinX/go-wordle/master/img/play.gif" width="60%" height="60%">

The word dictionary is extracted from the [javascript file of the official wordle](https://www.powerlanguage.co.uk/wordle/main.e65ce0a5.js).

### How to play

- Download the binary file from the release page

- Or build the binary file through `go install`

```shell  
go install github.com/XiaoMengXinX/go-wordle@latest
```

- Or compile it from the source code:

```shell
git clone https://github.com/XiaoMengXinX/go-wordle
cd go-wordle
go build wordle.go
```

Run:

```shell
./go-wordle
```

### Custom dictionaries

You can use the `-w` option to specify the custom word dictionary file like this:

```shell
./go-wordle -w your_wordlist.txt
```

The word list file must be in this format:

```
word1
word2
word3
...
wordN
```
