# SDQORM
**SDQORM** is a **S**omething **D**evided **Q**uery **O**bject **R**elation **M**apper made in golang.
This library link string `query` into golang struct.

# Background
- In some cases of solving programming problems.

# Usage
Please refer to [example/example.go](./example/example.go)

# Techniques
- `reflect` 

# TODO: 
- structに関して、CustomParserを作成せずともTagを指定するだけでパースできるようにする
- indexのみではなくたとえば1-7番目などの指定ができるようにする
- 事前定義しない処理(例:intについてはstrconv.Atoiによる変換のみサポートしているが、それ以外の処理も必要になる場合があり得そう)も実行できるようにしたい
