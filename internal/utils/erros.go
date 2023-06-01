package utils

import "errors"

var ErrSomethingWentWrong = errors.New("something went wrong")
var ErrUserAlreadyExists = errors.New("user already exists")
var ErrWrongCredentials = errors.New("wrong credentials")
var ErrUsersOrderAlreadyExists = errors.New("order already exists 1")
var ErrOtherOrderAlreadyExists = errors.New("order already exists 2")
var ErrNotEnoughBalance = errors.New("not enough balance")
var ErrWrongOrderNumber = errors.New("wrong order number")
var ErrWrongAPIKey = errors.New("wrong API key")
