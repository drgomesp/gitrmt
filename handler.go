package gitrmt

type Handler interface {
	Capabilities() []string
	List(forPush bool) ([]string, error)
	Push(localRef string, remoteRef string, force bool) (string, error)
	Finish() error
}
