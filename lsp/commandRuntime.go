package lsp

import (
	"fmt"
	"github.com/Sora233/DDBOT/lsp/concern"
	"github.com/Sora233/DDBOT/lsp/concern_type"
	"github.com/Sora233/DDBOT/lsp/parser"
	localutils "github.com/Sora233/DDBOT/utils"
	"github.com/alecthomas/kong"
	"io"
	"strings"
)

type Runtime struct {
	bot *localutils.HackedBot
	l   *Lsp
	*parser.Parser

	debug   bool
	exit    bool
	silence bool
}

func (r *Runtime) Exit(int) {
	r.exit = true
}
func (r *Runtime) Debug() {
	r.debug = true
}

func (r *Runtime) parseCommandSyntax(ast interface{}, name string, options ...kong.Option) (*kong.Context, string) {
	args := r.GetArgs()
	cmdOut := &strings.Builder{}
	// kong 错误信息不太友好
	options = append(options, kong.Name(name), kong.UsageOnError(), kong.Exit(r.Exit))
	k, err := kong.New(ast, options...)
	if err != nil {
		logger.Errorf("kong new failed %v", err)
		r.Exit(0)
		return nil, ""
	}
	if r.silence {
		k.Stdout = io.Discard
	} else {
		k.Stdout = cmdOut
	}
	ctx, err := k.Parse(args)
	if r.exit {
		logger.WithField("content", args).Debug("exit")
		return ctx, cmdOut.String()
	}
	if err != nil {
		logger.WithField("content", args).Errorf("kong parse failed %v", err)
		r.Exit(0)
		var out string
		if !r.silence {
			out = fmt.Sprintf("参数解析失败 - %v", err)
		}
		return nil, out
	}
	return ctx, ""
}

func (r *Runtime) ParseRawSiteAndType(rawSite string, rawType string) (string, concern_type.Type, error) {
	return concern.ParseRawSiteAndType(rawSite, rawType)
}

func (r *Runtime) ParseRawSite(rawSite string) (string, error) {
	return concern.ParseRawSite(rawSite)
}

func NewRuntime(l *Lsp, silence ...bool) *Runtime {
	r := &Runtime{
		bot:    localutils.GetBot(),
		l:      l,
		Parser: parser.NewParser(),
	}
	if len(silence) > 0 {
		r.silence = silence[0]
	}
	return r
}
