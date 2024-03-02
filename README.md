# Taggo

A dead simple library for parsing discriminated unions in Go.

## Motivation

If you ever parsed json, chances are you stumbled into a schema like the
following:

```json
{
    "shape": "square",
    "side": 10
}
```

```json
{
    "shape": "circle",
    "radius": 5
}
```

This is commonly reffered as a [tagged or discriminated
union](https://en.wikipedia.org/wiki/Tagged_union). And is usually a pain to
deal with in Go. It's so bad that the search for "golang discriminated union"
returns both of these in the first page:

 - [Reddit: Idiomatic way in Go to represent a Tagged Union?](https://www.reddit.com/r/golang/comments/13hjevf/idiomatic_way_in_go_to_represent_a_tagged_union/)
 - [Dax on X: "golang people - how do you deal..."](https://twitter.com/thdxr/status/1726269236577010031)

And pretty much all the suggestions are either quite bad or a simple "don't do
it".

So this library is meant both as a way to deal with it in the most general
cases and as a reference to a pattern that can be applied in other
circumstances (like parsing json with another library that not
"encoding/json").

## How to use

This lib provides two types: `Discriminator` and `Discriminated`.
`Discriminator` is an interface that you should implement that tells apart the
target types based on the tag:

```go
type ShapeDiscriminator struct {
	Shape string `json:"shape"`
}

func (sd *ShapeDiscriminator) GetType() (any, error) {
	switch sd.Shape {
	case "square":
		return &Square{}, nil
	case "circle":
		return &Circle{}, nil
	}
	return nil, errors.New("unknown shape")
}
```

If you wish, you can also return another interface that is not `any`:


```go
type Shape interface {
	Area() float64
}

type ShapeDiscriminator struct {
	Shape string `json:"shape"`
}

func (sd *ShapeDiscriminator) GetType() (Shape, error) {
	switch sd.Shape {
	case "square":
		return &Square{}, nil
	case "circle":
		return &Circle{}, nil
	}
	return nil, errors.New("unknown shape")
}
```

`Discriminated` then "marks" a type to be parsed using a `Discriminator`:

```go
type ShapeUnion = taggo.Discriminated[Shape, *ShapeDiscriminator]
```

Then you can parse a json into the `Discriminated` type

```go
var shapeUnion ShapeUnion

err := json.Unmarshal(bytes, &shapeUnion)
if err != nil {
	return err
}

shape := shapeUnion.Value // Either Circle or Square
```


## Full example


```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"

	"github.com/eduardomosko/taggo"
)

type Circle struct {
	Radius float64 `json:"radius"`
}

func (c *Circle) Area() float64 {
	return math.Pi * c.Radius * c.Radius
}

type Square struct {
	Side float64 `json:"side"`
}

func (s *Square) Area() float64 {
	return s.Side * s.Side
}

// Shape represents a 2D shape
type Shape interface {
	Area() float64
}

// ShapeDiscriminator tells the shapes apart
type ShapeDiscriminator struct {
	Shape string `json:"shape"`
}

func (sd *ShapeDiscriminator) GetType() (Shape, error) {
	switch sd.Shape {
	case "square":
		return &Square{}, nil
	case "circle":
		return &Circle{}, nil
	}
	return nil, errors.New("unknown shape")
}

// ShapeUnion parses a Shape by it's tag
type ShapeUnion = taggo.Discriminated[Shape, *ShapeDiscriminator]

func main() {
	inputSquare := []byte(`{"shape":"square","side":10}`)
	inputCircle := []byte(`{"shape":"circle","radius":5}`)

	shape1 := ShapeUnion{}
	shape2 := ShapeUnion{}

	err := json.Unmarshal(inputSquare, &shape1)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(inputCircle, &shape2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%#v\n", shape1.Value) // => &main.Square{Side:10}
	fmt.Printf("%#v\n", shape2.Value) // => &main.Circle{Radius:5}
}
```

## Marshaling

If you also need to marshal the types, my suggestion is just to
repeat the tag in the target type:

```go
type Circle struct {
	Shape  string `json:"shape"`
	Radius int    `json:"radius"`
}

type Square struct {
	Shape string `json:"shape"`
	Side  int    `json:"side"`
}
```

If you need something more fail-proof than that, contact me because I have some
ideas that are not mature enough to be turned into a library, but could be
useful.
