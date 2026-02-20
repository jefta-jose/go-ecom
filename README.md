### variable naming convention
variable type

### installing and removing missing dependencies 
go mod tidy

## Understanding the Mount Method

### Function Definition with Receiver
```go
func (app *application) mount() http.Handler
```

**Breaking it down:**

1. **`func`** - Keyword to define a function

2. **`(app *application)`** - The **receiver** (makes this a method)
   - This is like `this` or `self` in other languages
   - Means: "this function belongs to the `application` type"
   - You call it like: `app.mount()`

3. **`mount`** - The function name

4. **`()`** - Parameters (empty = no parameters needed)

5. **`http.Handler`** - Return type (what the function gives back)

### The `*` Asterisk (Pointer)

In Go, `*` means "pointer" - a reference to where something lives in memory.

- `*application` = pointer to the original application struct
- Without `*`, you'd get a copy of the struct
- With `*`, you work with the actual original (can modify it)

### Two Sets of Parentheses Explained

- `(app *application)` = **receiver** (who this method belongs to)
- `mount()` = **function name + parameters**

**Example with parameters:**
```go
func (app *application) mount(port string) http.Handler
                              // ^ parameter would go here
```

### How to Call It
```go
app := application{}  // create an application
router := app.mount()  // call the mount method
```