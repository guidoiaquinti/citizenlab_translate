package main

import "errors"

var ErrInvalidInputType = errors.New("Invalid input type provided. Available types are: json,yml")
