package runner

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	"github.com/thehunter365/gorunner/utils"
)

//Lang type
type Lang int

const (
	//GO golang code type
	GO Lang = iota
	//JAVA code
	JAVA
	//PYTHON code
	PYTHON
)

//RawCode type
type RawCode struct {
	CodeLines []string `json:"codeLines"`
}

//Runner type
type Runner struct {
	RawCode string

	Lang      Lang
	CodeLines []string
	Return    string

	TimeOut time.Duration
}

//NewRunner func
func NewRunner(lang Lang, jsonCode string) *Runner {
	return &Runner{
		RawCode: jsonCode,
		Lang:    lang,
		TimeOut: 5,
	}
}

//ParseCode from json to plain old text
func (r *Runner) ParseCode() (code []string) {
	var rc RawCode
	err := json.Unmarshal([]byte(r.RawCode), &rc)
	handleErr(err)
	code = rc.CodeLines
	r.CodeLines = code
	return
}

//StartRunner func
func (r *Runner) StartRunner() (out string) {
	out = r.execCode()
	r.Return = out
	return
}

//ExecCode func
func (r *Runner) execCode() (stdout string) {
	c1 := make(chan []byte, 1)

	if len(r.CodeLines) == 0 {
		r.ParseCode()
	}

	utils.FileWrite("tmp.go", r.CodeLines)
	cmd := exec.Command("go", "run", "../tmp.go")

	go func() {
		out, err := cmd.CombinedOutput()
		c1 <- out
		if err != nil {
			log.Fatalln(err)
		}
	}()
	select {
	case <-time.After(r.TimeOut * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			log.Fatal("failed to kill process: ", err)
		}
		log.Println("process timed out")
	case o := <-c1:
		stdout = string(o)
		log.Print("process finished successfully")
	}
	utils.FileDelete("tmp.go")
	return
}

func handleErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
