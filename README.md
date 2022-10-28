# atomicfile

Write files to a temp file first and atomically rename them on close.

```go
import "github.com/marwatk/atomicfile"

func main(args []string) {
  f, err := atomicfile.New("path/to/final.txt", 0644)
  if err != nil {
    fmt.Printf("Error opening atomic file: %v", err)
  }
  _, err = fmt.Fprintf(f, "New content")
  if err != nil {
    fmt.Printf("Error writing to file: %v", err)
  }
  err = f.Close()
  if err != nil {
    fmt.Printf("Error renaming to file to final name: %v", err)
  }
}
```

Based on [github.com/natefinch/atomic](github.com/natefinch/atomic)

By default, writing to a file in go (and generally any language) can fail
partway through... you then have a partially written file, which probably was
truncated when the write began, and bam, now you've lost data.

This go package avoids this problem, by writing first to a temp file, and then
overwriting the target file in an atomic way.  This is easy on linux, os.Rename
just is atomic.  However, on Windows, os.Rename is not atomic, and so bad things
can happen.  By wrapping the windows API moveFileEx, we can ensure that the move
is atomic, and we can be safe in knowing that either the move succeeds entirely,
or neither file will be modified.

