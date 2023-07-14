[![CI + Unit Test](https://github.com/dyammarcano/utils/actions/workflows/ci.yml/badge.svg)](https://github.com/dyammarcano/utils/actions/workflows/ci.yml)

# Binary Serialize

### TODO:

- [ ] unit test 

### Usage

```go
package main

import (
    "fmt"
    "github.com/dyammarcano/utils/"
)

type (
	MyStruct struct {
		Field1      int32
		Field2      string
		OtherStruct OtherStruct
	}

	OtherStruct struct {
		FieldA float64
		FieldB bool
	}
)

func main() {
    // struct to serialize
    inputData := MyStruct{
		Field1: 42,
		Field2: "Hello, World!",
		OtherStruct: OtherStruct{
			FieldA: 3.14159,
			FieldB: true,
		},
	}

    // serialize a struct
	serializedData, err := Marshal(inputData)
    if err != nil {
        panic(err)
    }

    // deserialize a string
    var deserializedData MyStruct
	if err := Unmarshal(serializedData, &deserializedData); err != nil {
        panic(err)
    }

    // compare structs
    if reflect.DeepEqual(inputData, deserializedData) {
        fmt.Println("Structs are equal")
    }
}
```

### Inport package in your project

```go
import "github.com/dyammarcano/utils"
```