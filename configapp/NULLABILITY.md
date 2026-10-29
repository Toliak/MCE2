# Nullability conventions

Go has no built-in null safety. These conventions make pointer and optional-value usage explicit and consistent.

## Pointers

- **Pointers are non-nil by default.** A function parameter or struct field declared as `*T` must not be `nil` unless its documentation explicitly says otherwise.
- Return `nil` only together with an `error`, or document a specific exception in the function/field comment.
- Public constructors must validate required pointers and return an error instead of storing `nil`.
- Avoid pointer-to-pointer (`**T`) solely to model optionality. Use `optional.Optional[T]` instead.

## Optional values

Use `optional.Optional[T]` (package `configapp/optional`) for values that may legitimately be absent:

```go
var value optional.Optional[string] = optional.None[string]()
value = optional.Some("hello")

if v, ok := value.Get(); ok {
    fmt.Println(v)
}
```

Do **not** use pointers as an Optional substitute. If a value can be absent, model it with `optional.Optional[T]`.

## Nil maps and slices

Nil maps and slices are allowed where Go or JSON zero-value semantics require them. Document such fields explicitly, for example:

```go
// May be nil; nil means "absent".
Params map[string]string `json:"p,omitempty"`
```

## Documentation

Add an explicit comment to every pointer or optional field/parameter:

- `// Must not be nil.`
- `// May be nil; nil means ...`
- `// Optional; absence means ...`
