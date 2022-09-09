# tealfmt

tealfmt is a formatter for TEAL, the language used to write Algorand smart contracts.

## Rules

### Indent All Opcodes
Every opcode is indented. Note the emphasis on opcode. Lines with `#pragma` and labels are not indented:

```c
#pragma version 1

int 1

two:
int 2
```

```c
#pragma version 1
    
    int 1
    
two:
    int 2
```

### Indent comments
All lines with comments are indented, unless the are above a label. If there are multiple comment in a row above a label, they will not be indented

```c
// one
int 1

// foo
// bar
two:
//two
int 2
```

```c
    // one
    int 1

// foo
// bar
two:
    //two
    int 2
```

### Add whitespace before comment line(s)
Every line that is solely a comment will get whitespace above it, unless the line above it is a comment

```c
int 1
// two
int 2
// three
// four
int 3
int 4
```

```c
int 1

    // two
    int 2

    // three
    // four
    int 3
    int 4
```

### Add whitespace after void ops
Whitepsace is added after *most* opcodes that don't put anything on the stack. Whitespace won't be added back-to-back.


```c
store 9
store 8
store 7
load 7
gtxns Receiver
global CurrentApplicationAddress
==
assert
load 7
gtxns Amount
pushint 100000 // 100000
==
assert
bytec_2 // "auction_end"
global LatestTimestamp
load 9
+
app_global_put
```

```c
    store 9
    store 8
    store 7

    load 7
    gtxns Receiver
    global CurrentApplicationAddress
    ==
    assert

    load 7
    gtxns Amount
    pushint 100000 // 100000
    ==
    assert

    bytec_2 // "auction_end"
    global LatestTimestamp
    load 9
    +
    app_global_put
```
