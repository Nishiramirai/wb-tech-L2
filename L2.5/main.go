package main

// Что выведет программа?

// Объяснить вывод программы.

type customError struct {
	msg string
}

func (e *customError) Error() string {
	return e.msg
}

func test() *customError {
	// ... do something
	return nil
}

func main() {
	var err error
	err = test()

	// так как функция test возвращает nil указатель типа *customError,
	// то внутри интерфейса err лежит nil значение типа *customError
	// и сам интерфейс err равен nil
	if err != nil {
		println("error")
		return
	}
	println("ok")
}
