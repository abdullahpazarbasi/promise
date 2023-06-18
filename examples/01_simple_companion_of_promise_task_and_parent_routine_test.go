package examples

import (
	"fmt"
	"github.com/abdullahpazarbasi/promise"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

func Test_Simple_Companion_of_Promise_Task_and_Parent_Routine(t *testing.T) {
	f := promise.New(getHostname)

	fmt.Println(">  ", "sub-task has not started yet")

	p := f.TimeOutLimit(500 * time.Millisecond).Commit()
	fmt.Println(">  ", "sub-task has just started")

	fmt.Println(">  ", "parent routine tasks ...")

	hostname, err := p.Await()
	require.NoError(t, err)
	fmt.Println(">  ", "hostname from promise:", hostname)
}

func getHostname() (string, error) {
	time.Sleep(200 * time.Millisecond)

	return os.Hostname()
}