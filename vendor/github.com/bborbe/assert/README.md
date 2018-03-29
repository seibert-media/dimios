# Assert

Assert ala Java Hamcrest in Go

## Usage

```
import (
	"testing"
	. "github.com/bborbe/assert"
)

func TestEquals(t *testing.T) {
  value := ...
  if err := AssertThat(value, Is("foo")); err != nil {
    t.Fatal(err)
  }
}
```
