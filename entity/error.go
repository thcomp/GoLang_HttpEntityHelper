package entity

import "fmt"

var ErrNotExistFormData error = fmt.Errorf("not exist form data")
var ErrUnsupportEntity error = fmt.Errorf("unsupport entity")
