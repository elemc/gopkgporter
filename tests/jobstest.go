package tests

import (
	"gopkgporter/app/jobs"

	"github.com/revel/revel/testing"
)

type JobsTest struct {
	testing.TestSuite
}

func (t *JobsTest) Before() {
	println("Set up")
}

func (t *JobsTest) TestGetFromKoji() {
	getFromKojiJob := jobs.GetFromKoji{}
	getFromKojiJob.Run()
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *JobsTest) After() {
	println("Tear down")
}
