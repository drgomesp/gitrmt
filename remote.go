package gitrmt

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type RemoteHandler interface {
	Capabilities() string
	List(forPush bool) ([]string, error)
	Push(localRef string, remoteRef string, force bool) (string, error)
	Finish() error
}

type Remote struct {
	in       io.Reader
	out      io.Writer
	handler  RemoteHandler
	lazyWork []func() (string, error)
}

func NewRemote(in io.Reader, out io.Writer, handler RemoteHandler) *Remote {
	log.Printf("$GIT_DIR=%v\n", os.Getenv("GIT_DIR"))

	return &Remote{
		in:       in,
		out:      out,
		handler:  handler,
		lazyWork: make([]func() (string, error), 0),
	}
}

func (r *Remote) Run() error {
	reader := bufio.NewReader(r.in)

loop:
	for {
		command, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		command = strings.Trim(command, "\n")

		switch {
		case command == "capabilities":
			fmt.Fprintf(r.out, "%s\n", r.handler.Capabilities())
		case strings.HasPrefix(command, "list"):
			list, err := r.handler.List(strings.HasPrefix(command, "list for-push"))
			if err != nil {
				return err
			}

			for _, e := range list {
				fmt.Fprintf(r.out, "%s\n", e)
			}

			_, _ = fmt.Fprint(r.out, "\n")
		case strings.HasPrefix(command, "push "):
			refs := strings.Split(command[5:], ":")
			isForce := strings.HasPrefix(refs[0], "+")
			r.push(refs, isForce)
		case strings.HasPrefix(command, "fetch "):
			parts := strings.Split(command, " ")
			if parts[1] != "0000000000000000000000000000000000000000" {
				r.fetch(parts[1], parts[2])
			}
		case command == "":
			fallthrough
		case command == "\n":
			log.Println("processing tasks...")
			for _, task := range r.lazyWork {
				resp, err := task()
				if err != nil {
					return fmt.Errorf("error processing task: %w", err)
				}
				fmt.Fprintf(r.out, "%s", resp)
			}
			_, _ = fmt.Fprintf(r.out, "\n")
			r.lazyWork = nil
			break loop
		default:
			return fmt.Errorf("received unknown command %q", command)
		}
	}

	return r.handler.Finish()
}
func (r *Remote) fetch(sha string, ref string) {
	r.lazyWork = append(r.lazyWork, func() (string, error) {
		return "", nil
	})
}

func (r *Remote) push(refs []string, force bool) {
	src, dst := refs[0], refs[1]

	r.lazyWork = append(r.lazyWork, func() (string, error) {
		log.Println("push")
		done, err := r.handler.Push(src, dst, force)
		if err != nil {
			return "", err
		}

		return fmt.Sprintf("ok %s\n", done), nil
	})
}
