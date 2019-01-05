package integration_tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"os"
	"os/exec"
	"testing"
	"time"

	"gopkg.in/h2non/baloo.v3"
)

var test = baloo.New("http://localhost:5005")

func TestMain(m *testing.M) {
	go start()
	time.Sleep(1 * time.Second)
	output := m.Run()
	out, _ := exec.Command("./end.sh").CombinedOutput()
	fmt.Printf("\nOutput:\n%s\n", string(out))
	os.Exit(output)
}

func TestHealth(t *testing.T) {
	assert.Nil(t, test.Get("/api/health").
		Expect(t).
		Status(200).
		Type("json").
		Done())
}


func start() {
	exec.Command("./start.sh", ).Run()
}