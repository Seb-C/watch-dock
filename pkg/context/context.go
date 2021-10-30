package context

import (
	"context"
	"flag"
)

// Compile-time check that the type implements the interface
var _ context.Context = Context{}

type Context struct{
	context.Context
	dockerComposeArgs []string
}

func ParseContext() Context {
	flag.Parse()

	dockerComposeArgs := []string{}
	for argIndex, arg := range flag.Args() {
		if arg == "-" || arg == "--" {
			dockerComposeArgs = flag.Args()[argIndex+1:]
		}
	}

	return Context{
		context.Background(),
		dockerComposeArgs,
	}
}

func (ctx Context) DockerComposeArgs(args... string) []string {
	out := make([]string, 0, len(args)+len(ctx.dockerComposeArgs))
	out = append(out, ctx.dockerComposeArgs...)
	out = append(out, args...)
	return out
}
