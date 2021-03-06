// Copyright 2019 Michael Li <alimy@gility.net>. All rights reserved.
// Use of this source code is governed by Apache License 2.0 that
// can be found in the LICENSE file.

package core

import (
	"context"
	"log"
)

const (
	// run mode list
	InSerialMode runMode = iota
	InConcurrentMode
	InSerialDebugMode
	InConcurrentDebugMode

	// generator Names
	GeneratorGin        = "gin"
	GeneratorChi        = "chi"
	GeneratorMux        = "mux"
	GeneratorEcho       = "echo"
	GeneratorIris       = "iris"
	GeneratorFiber      = "fiber"
	GeneratorMacaron    = "macaron"
	GeneratorHttpRouter = "httprouter"

	// parser Names
	ParserStructTag = "structTag"
)

var (
	// generators generator list
	generators = make(map[string]Generator, 4)

	// parsers parser list
	parsers = make(map[string]Parser, 1)

	// inDebug whether in debug mode
	inDebug bool
)

// runMode indicate process mode (InSerialMode | InSerialDebugMode | InConcurrentMode | InConcurrentDebugMode)
type runMode uint8

// InitOpts use for generator or parser init
type InitOpts struct {
	RunMode       runMode
	GeneratorName string
	ParserName    string
	SinkPath      string
	DefaultTag    string
	NoneQuery     bool
	Cleanup       bool
}

// ParserOpts used for initial parser
type ParserOpts struct {
	DefaultTag string
	NoneQuery  bool
}

// GeneratorOpts used for initial generator
type GeneratorOpts struct {
	SinkPath string
	Cleanup  bool
}

// Option pass option to custom run behavior
type Option interface {
	apply(opts *InitOpts)
}

// Options generator options
type Options []Option

// InitOpts return an initOpts instance
func (opts Options) InitOpts() *InitOpts {
	res := defaultInitOpts()
	for _, opt := range opts {
		opt.apply(res)
	}
	return res
}

// ParserOpts return a ParserOpts instance
func (opts *InitOpts) ParserOpts() *ParserOpts {
	return &ParserOpts{
		DefaultTag: opts.DefaultTag,
		NoneQuery:  opts.NoneQuery,
	}
}

// GeneratorOpts return a GeneratorOpts
func (opts *InitOpts) GeneratorOpts() *GeneratorOpts {
	return &GeneratorOpts{
		SinkPath: opts.SinkPath,
		Cleanup:  opts.Cleanup,
	}
}

// optFunc used for convert function to Option interface
type optFunc func(opts *InitOpts)

func (f optFunc) apply(opts *InitOpts) {
	f(opts)
}

// Parser parse entries
type Parser interface {
	Name() string
	Init(opts *ParserOpts) error
	Parse(entries []interface{}) (Descriptors, error)
	ParseContext(ctx MirCtx, entries []interface{})
	Clone() Parser
}

// Generator generate interface code for engine
type Generator interface {
	Name() string
	Init(opts *GeneratorOpts) error
	Generate(Descriptors) error
	GenerateContext(ctx MirCtx)
	Clone() Generator
}

// MirCtx mir's concurrent parser/generator context
type MirCtx interface {
	context.Context
	Cancel(err error)
	ParserDone()
	GeneratorDone()
	Wait() error
	Capcity() int
	Pipe() (<-chan *IfaceDescriptor, chan<- *IfaceDescriptor)
}

// String runMode describe
func (m runMode) String() string {
	res := "not support mode"
	switch m {
	case InSerialMode:
		res = "serial mode"
	case InSerialDebugMode:
		res = "serial debug mode"
	case InConcurrentMode:
		res = "concurrent mode"
	case InConcurrentDebugMode:
		res = "concurrent debug mode"
	}
	return res
}

// RunMode set run mode option
func RunMode(mode runMode) Option {
	return optFunc(func(opts *InitOpts) {
		opts.RunMode = mode
	})
}

// GeneratorName set generator name option
func GeneratorName(name string) Option {
	return optFunc(func(opts *InitOpts) {
		opts.GeneratorName = name
	})
}

// ParserName set parser name option
func ParserName(name string) Option {
	return optFunc(func(opts *InitOpts) {
		opts.ParserName = name
	})
}

// SinkPath set generated code out directory
func SinkPath(path string) Option {
	return optFunc(func(opts *InitOpts) {
		opts.SinkPath = path
	})
}

// Cleanup set generator cleanup out first when re-generate code
func Cleanup(enable bool) Option {
	return optFunc(func(opts *InitOpts) {
		opts.Cleanup = enable
	})
}

// NoneQuery set parser whether parse query
func NoneQuery(enable bool) Option {
	return optFunc(func(opts *InitOpts) {
		opts.NoneQuery = enable
	})
}

// DefaultTag set parser's default struct field tag string key
func DefaultTag(tag string) Option {
	return optFunc(func(opts *InitOpts) {
		opts.DefaultTag = tag
	})
}

// RegisterGenerators register generators
func RegisterGenerators(gs ...Generator) {
	for _, g := range gs {
		if g != nil && g.Name() != "" {
			generators[g.Name()] = g
		}
	}
}

// RegisterParsers register parsers
func RegisterParsers(ps ...Parser) {
	for _, p := range ps {
		if p != nil && p.Name() != "" {
			parsers[p.Name()] = p
		}
	}
}

// GeneratorByName get a generator by name
func GeneratorByName(name string) Generator {
	return generators[name]
}

// DefaultGenerator get a default generator
func DefaultGenerator() Generator {
	return generators[GeneratorGin]
}

// ParserByName get a parser by name
func ParserByName(name string) Parser {
	return parsers[name]
}

// DefaultParser get a default parser
func DefaultParser() Parser {
	return parsers[ParserStructTag]
}

// Logus print log info
func Logus(format string, v ...interface{}) {
	if inDebug {
		log.Printf("[mir] "+format, v...)
	}
}

// InitFrom initial from Options and return an InitOpts instance
func InitFrom(opts Options) *InitOpts {
	var initOpts *InitOpts
	if opts == nil {
		initOpts = defaultInitOpts()
	} else {
		initOpts = opts.InitOpts()
	}

	switch initOpts.RunMode {
	case InSerialDebugMode, InConcurrentDebugMode:
		inDebug = true
	default:
		inDebug = false
	}

	return initOpts
}

func defaultInitOpts() *InitOpts {
	return &InitOpts{
		RunMode:       InSerialMode,
		GeneratorName: GeneratorGin,
		ParserName:    ParserStructTag,
		SinkPath:      ".gen",
		DefaultTag:    "mir",
		Cleanup:       true,
	}
}
