package tests

import (
	"fmt"

	"github.com/revel"
)

type AppTest struct {
	revel.TestSuite
}

func (t *AppTest) Before() {
	fmt.Println("Set up")
}

func (t *AppTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) TestUserRegister() {
	t.Get("/register")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *AppTest) After() {
	fmt.Println("Tear down")
}
