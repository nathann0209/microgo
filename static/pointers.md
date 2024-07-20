### <u> What is a pointer: </u> 
- A pointer **holds the location in memory where a value (that a variable is assigned to) is stored**.

### <u> Using pointers in Go: </u>
The & is the address operator. It precedes a value type and returns
the address of the memory location where the value is stored: 

```
x := "hello"
pointerToX := &x
```

'&' must be used before a variable and not a constant. This is because constants do not have memory adresses and exist only at compile time (more on this later).

---
The '\*' operator can "dereference" a pointer and return the pointed to variable. 

```
var x int = 5
var xPointer *int = &x

*xPointer = 10

fmt.Println(x)

// This code prints '10'
```

Here we access 'x' by dereferencing \*xPointer. 

(so  \*xPointer = 10 is the same as x = 10 in result)

The main difference though is that \*xPointer = 10 is a more fundamental assignment in that it changes the value at the address that the value of the variable x is stored. This allows for byRef changes in functions. 

(see Functions: Part 3):


**NOTE: Dereferencing does not work when the pointer is nil**. 

---
'new' creates a pointer to a zero value instance of the type that is passed: 

```
var x = new(int)
fmt.Println(x == nil) // prints false
fmt.Println(*x)
// prints 0
```

Here 'x' is the pointer, so 'x' itself is not nil. 'x' simply points to the memory adress of an integer 0.

Therfeore, dereferencing 'x' returns 0. 

### <u> Pointing to a constant </u>
Earlier I stated that pointing to a constant does not work because constant values do not have memory adresses. 

```
PointerToEight := &8

fmt.Println(*PointerToEight)

// This returns: ./prog.go:8:21: invalid operation: cannot take address of 8 (untyped int constant)

```

There are some ways to get around this though: 

1. Assign a variable to your constant of choice (a variable always has a memory address)

```
var eight int = 8

PointerToEight := &eight

fmt.Println(*PointerToEight)

// This prints: 8
```

2. Write a function that takes in a variable and returns the pointer to the variable.
```
func stringp(s string) *string {
	return &s
}

var PointerToString *string

PointerToString = stringp("hello")

fmt.Println(*PointerToString)

// This prints: "hello"


```

This works because constant arguments are copied to a parameter, which is a variable, so can have a pointer to the address of the value it holds. 

NOTE: "Constant" values usually come in the form of a string, numeric or boolean. 

### <u> Why use pointers </u>
1. Functions which modify parameters ByRef.
2. When we have large amounts of data, making copies to pass between functions is very inefficient.  

### <u> Slices and Maps are already pointers </u>

When you pass a slice to a function, you're passing a descriptor that inludes:

- **Pointer to the underlying array:** This points to the actual data.
- **Length:** The number of elements in the slice.
- **Capacity:** The size of the underlying array starting from the slice's start index.
```
func incrementPeterAge(m map[string]int) {
    m["Peter"] += 1
}

func main() {
    ages := map[string]int{
        "Peter": 21,
    }
    incrementPeterAge(ages)
    fmt.Println(ages)  // Output: map[Peter:22]
}
```

(similar to a map).

This is why byRef changes to a map/slic can be made when passing them as arguments to a function. 
