package cardanocli

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"strings"
)

// Hook hook for command running
type Hook interface {
	BeforeRun(x *exec.Cmd)
	AfterRun(x *exec.Cmd, err *error)
}

type hookCollectStdout struct {
	buf *bytes.Buffer
	out *string
}

func (h *hookCollectStdout) BeforeRun(x *exec.Cmd) {
	x.Stdout = h.buf
}

func (h *hookCollectStdout) AfterRun(x *exec.Cmd, err *error) {
	if *err != nil {
		return
	}
	*h.out = strings.TrimSpace(h.buf.String())
}

// CollectStdout create a hook that collects stdout as string
func CollectStdout(out *string) Hook {
	return &hookCollectStdout{
		buf: &bytes.Buffer{},
		out: out,
	}
}

type hookCollectStderr struct {
	buf *bytes.Buffer
	out *string
}

func (h *hookCollectStderr) BeforeRun(x *exec.Cmd) {
	x.Stderr = h.buf
}

func (h *hookCollectStderr) AfterRun(x *exec.Cmd, err *error) {
	if *err != nil {
		return
	}
	*h.out = strings.TrimSpace(h.buf.String())
}

// CollectStderr create a hook that collects stderr as string
func CollectStderr(out *string) Hook {
	return &hookCollectStderr{
		buf: &bytes.Buffer{},
		out: out,
	}
}

// CollectStdoutJSON create a hook that collects stdout as JSON
type hookCollectStdoutJSON struct {
	buf *bytes.Buffer
	out interface{}
}

func (h *hookCollectStdoutJSON) BeforeRun(x *exec.Cmd) {
	x.Stdout = h.buf
}

func (h *hookCollectStdoutJSON) AfterRun(x *exec.Cmd, err *error) {
	if *err != nil {
		return
	}
	*err = json.Unmarshal(h.buf.Bytes(), h.out)
}

func CollectStdoutJSON(out interface{}) Hook {
	return &hookCollectStdoutJSON{
		buf: &bytes.Buffer{},
		out: out,
	}
}
