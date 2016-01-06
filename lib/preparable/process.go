package preparable

import "os/exec"
import "strings"

import "github.com/topfreegames/apm/lib/process"

// ProcPreparable is a preparable with all the necessary informations to run
// a process. To actually run a process, call the Start() method.
type ProcPreparable struct {
	Name       string
	SourcePath string
	Cmd        string
	SysFolder  string
	Language   string
	KeepAlive  bool
	Args       []string
}

// PrepareBin will compile the Golang project from SourcePath and populate Cmd with the proper
// command for the process to be executed.
// Returns the compile command output.
func (proc_preparable *ProcPreparable) PrepareBin() ([]byte, error) {
	// Remove the last character '/' if present
	if proc_preparable.SourcePath[len(proc_preparable.SourcePath)-1] == '/' {
		proc_preparable.SourcePath = strings.TrimSuffix(proc_preparable.SourcePath, "/")
	}
	cmd := ""
	cmdArgs := []string{}
	binPath := proc_preparable.getBinPath()
	if proc_preparable.Language == "go" {
		cmd = "go"
		cmdArgs = []string{"build", "-o", binPath, proc_preparable.SourcePath + "/."}
	}

	proc_preparable.Cmd = proc_preparable.getBinPath()
	return exec.Command(cmd, cmdArgs...).Output()
}

// Start will execute the process based on the information presented on the preparable.
// This function should be called from inside the master to make sure
// all the watchers and process handling are done correctly.
// Returns a tuple with the process and an error in case there's any.
func (proc_preparable *ProcPreparable) Start() (*process.Proc, error) {
	proc := &process.Proc{
		Name:      proc_preparable.Name,
		Cmd:       proc_preparable.Cmd,
		Args:      proc_preparable.Args,
		Path:      proc_preparable.getPath(),
		Pidfile:   proc_preparable.getPidPath(),
		Outfile:   proc_preparable.getOutPath(),
		Errfile:   proc_preparable.getErrPath(),
		KeepAlive: proc_preparable.KeepAlive,
		Status:    &process.ProcStatus{},
	}

	err := proc.Start()
	return proc, err
}

func (proc_preparable *ProcPreparable) getPath() string {
	if proc_preparable.SysFolder[len(proc_preparable.SysFolder)-1] == '/' {
		proc_preparable.SysFolder = strings.TrimSuffix(proc_preparable.SysFolder, "/")
	}
	return proc_preparable.SysFolder + "/" + proc_preparable.Name
}

func (proc_preparable *ProcPreparable) getBinPath() string {
	return proc_preparable.getPath() + "/" + proc_preparable.Name
}

func (proc_preparable *ProcPreparable) getPidPath() string {
	return proc_preparable.getBinPath() + ".pid"
}

func (proc_preparable *ProcPreparable) getOutPath() string {
	return proc_preparable.getBinPath() + ".out"
}

func (proc_preparable *ProcPreparable) getErrPath() string {
	return proc_preparable.getBinPath() + ".err"
}
