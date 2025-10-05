# goodice
A good, minimal, efficient go rpg-like dice library.
Straight to the point.

`go get github.com/5000k/goodice`

## Basic Usage

```go
package main

import (
	"github.com/5000K/goodice"
)

func main(){
	result, err := goodice.Generate("D20")
	println(result.Result)
	
	result2, err := goodice.GenerateSeeded("3D20-2D8+12", 1337)
	println(result2.Result) // total
	println(result2.Parts) // contains the result of each dice roll plus all constants
}
```



## Supported dice notation
The parser accepts the standard TTRPG dice notation. It has two main parts: constants and dice rolls

A constant is an integer number (e.g. `7` or `42` - I think you know what an integer is)

Dice rolls are in the form of XdY, where:
- X is an optional integer, **How many dice to roll**. If not provided, 1 die is assumed.
- d simply is "d" or "D"
- Y is the amount of sides each die has (e.g. 20 for a D20)

So `D20`, `d20` and `1d20` will all do the same.

The parts can be chained together with the operators + and - (add and subtract). The whole string will be evaluated from left to right, just like you would when calculating it by hand.

## Reusing parsed scripts
If you are going to use a script many times, it will be more efficient to reuse a Goodice instance.
For this, you'll create an instance with `goodice.New(script: string)`. This will parse your script. You can then repeatedly call `Generate()` or `GenerateSeeded(seed:int)` on the instance.
