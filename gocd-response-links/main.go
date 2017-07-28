// Adapted from `stringer`:
// - https://blog.golang.org/generate
// - http://godoc.org/golang.org/x/tools/cmd/stringer

package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/build"
	"go/format"
	"go/parser"
	"go/token"
	"go/types"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	typeNames = flag.String("type", "", "comma-separated list of type names; must be set")
	output    = flag.String("output", "", "output file name; default srcdir/<type>_string.go")
)

// Usage is a replacement usage function for the flags package.
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\tgocd-response-links [flags] -type T [directory]\n")
	fmt.Fprintf(os.Stderr, "\tgocd-response-links [flags] -type T files... # Must be a single package\n")
	fmt.Fprintf(os.Stderr, "Flags:\n")
	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("gocd-response-links: ")
	flag.Usage = Usage
	flag.Parse()
	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}
	types := strings.Split(*typeNames, ",")

	// We accept either one directory or a list of files. Which do we have?
	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	// Parse the package once.
	var (
		dir string
		g   Generator
	)
	if len(args) == 1 && isDirectory(args[0]) {
		dir = args[0]
		g.parsePackageDir(args[0])
	} else {
		dir = filepath.Dir(args[0])
		g.parsePackageFiles(args)
	}

	// Print the header and package clause.
	g.Printf("// Code generated by \"gocd-response-links %s\"; DO NOT EDIT.\n", strings.Join(os.Args[1:], " "))
	g.Printf("\n")
	g.Printf("package %s", g.pkg.name)
	g.Printf("\n")
	g.Printf("import (\"encoding/json\"\n\"net/url\"\n)\n") // Used by all methods.

	// Run generate for each type.
	for _, typeName := range types {
		g.generate(typeName)
	}

	// Format the output.
	src, err := g.format()
	if err != nil {
		// Should never happen, but can arise when developing this code.
		// The user can compile the output to see the error.
		log.Printf("warning: internal error: invalid Go generated: %s", err)
		log.Printf("warning: compile the package to analyze the error")
		src = g.buf.Bytes()
	}

	// Write to file.
	outputName := *output
	if outputName == "" {
		baseName := fmt.Sprintf("gen_%s.go", strings.ToLower(types[0]))
		outputName = filepath.Join(dir, strings.ToLower(baseName))
	}
	err = ioutil.WriteFile(outputName, src, 0644)
	if err != nil {
		log.Fatalf("writing output: %s", err)
	}
}

// isDirectory reports whether the named file is a directory.
func isDirectory(name string) bool {
	info, err := os.Stat(name)
	if err != nil {
		log.Fatal(err)
	}
	return info.IsDir()
}

// Generator holds the state of the analysis. Primarily used to buffer
// the output for format.Source.
type Generator struct {
	buf bytes.Buffer // Accumulated output.
	pkg *Package     // Package we are scanning.
}

func (g *Generator) Printf(format string, args ...interface{}) {
	fmt.Fprintf(&g.buf, format, args...)
}

// File holds a single parsed file and associated data.
type File struct {
	pkg  *Package  // Package to which this file belongs.
	file *ast.File // Parsed AST.
	// These fields are reset for each type being generated.
	typeName string   // Name of the constant type.
	values   []string // Accumulator for constant values of that type.
}

type Package struct {
	dir      string
	name     string
	defs     map[*ast.Ident]types.Object
	files    []*File
	typesPkg *types.Package
}

// parsePackageDir parses the package residing in the directory.
func (g *Generator) parsePackageDir(directory string) {
	pkg, err := build.Default.ImportDir(directory, 0)
	if err != nil {
		log.Fatalf("cannot process directory %s: %s", directory, err)
	}
	var names []string
	names = append(names, pkg.GoFiles...)
	names = append(names, pkg.CgoFiles...)
	// TODO: Need to think about constants in test files. Maybe write type_string_test.go
	// in a separate pass? For later.
	// names = append(names, pkg.TestGoFiles...) // These are also in the "foo" package.
	names = append(names, pkg.SFiles...)
	names = prefixDirectory(directory, names)
	g.parsePackage(directory, names, nil)
}

// parsePackageFiles parses the package occupying the named files.
func (g *Generator) parsePackageFiles(names []string) {
	g.parsePackage(".", names, nil)
}

// prefixDirectory places the directory name on the beginning of each name in the list.
func prefixDirectory(directory string, names []string) []string {
	if directory == "." {
		return names
	}
	ret := make([]string, len(names))
	for i, name := range names {
		ret[i] = filepath.Join(directory, name)
	}
	return ret
}

// parsePackage analyzes the single package constructed from the named files.
// If text is non-nil, it is a string to be used instead of the content of the file,
// to be used for testing. parsePackage exits if there is an error.
func (g *Generator) parsePackage(directory string, names []string, text interface{}) {
	var files []*File
	var astFiles []*ast.File
	g.pkg = new(Package)
	fs := token.NewFileSet()
	for _, name := range names {
		if !strings.HasSuffix(name, ".go") {
			continue
		}
		parsedFile, err := parser.ParseFile(fs, name, text, 0)
		if err != nil {
			log.Fatalf("parsing package: %s: %s", name, err)
		}
		astFiles = append(astFiles, parsedFile)
		files = append(files, &File{
			file: parsedFile,
			pkg:  g.pkg,
		})
	}
	if len(astFiles) == 0 {
		log.Fatalf("%s: no buildable Go files", directory)
	}
	g.pkg.name = astFiles[0].Name.Name
	g.pkg.files = files
	g.pkg.dir = directory
	// Type check the package.
	g.pkg.check(fs, astFiles)
}

// check type-checks the package. The package must be OK to proceed.
func (pkg *Package) check(fs *token.FileSet, astFiles []*ast.File) {
	pkg.defs = make(map[*ast.Ident]types.Object)
	config := types.Config{Importer: defaultImporter(), FakeImportC: true}
	info := &types.Info{
		Defs: pkg.defs,
	}
	typesPkg, err := config.Check(pkg.dir, fs, astFiles, info)
	if err != nil {
		log.Fatalf("checking package: %s", err)
	}
	pkg.typesPkg = typesPkg
}

// generate produces the String method for the named type.
func (g *Generator) generate(typeName string) {
	values := make([]string, 0, 100)
	for _, file := range g.pkg.files {
		// Set the state for this run of the walker.
		file.typeName = typeName
		file.values = nil
		if file.file != nil {
			ast.Inspect(file.file, file.genDecl)
			values = append(values, file.values...)
		}
	}

	if len(values) == 0 {
		log.Fatalf("no values defined for type %s", typeName)
	}

	g.buildMarshalling(values, typeName)
	g.buildUnmarshalling(values, typeName)
}

const marshall_header = `func (l %s) MarshalJSON() ([]byte, error) {
type h struct {
	h string ` + "`json:\"href\"`" + `
}
ls := struct {
`

const marshall_footer = `j, e := json.Marshal(ls)
if e != nil {
	return nil, e
}
return j, nil
	`

func (g *Generator) buildMarshalling(values []string, typeName string) {
	g.Printf(marshall_header, typeName)
	for _, field := range values {
		g.Printf(fmt.Sprintf("%s *h `json:\"%s,omitempty\"`\n", field, strings.ToLower(field)))
	}
	g.Printf("}{}\n")

	for _, field := range values {
		g.Printf(fmt.Sprintf("if l.%s != nil {ls.%s = &h{h:l.%s.String()}}\n", field, field, field))
	}
	g.Printf(marshall_footer)
	g.Printf("}\n")
}

const unmarshall_header = `
func (l *%s) UnmarshalJSON(j []byte) error {
	var d map[string]map[string]string
	e := json.Unmarshal(j, &d)
	if e != nil {
		return e
	}
`

const unmarshall_field = `
if d["%s"]["href"] != "" {
	l.%s, e = url.Parse(d["%s"]["href"])
	if e != nil {
		return e
	}
}`

const unmarshall_footer = `
return nil
}
`

func (g *Generator) buildUnmarshalling(values []string, typeName string) {
	g.Printf(unmarshall_header, typeName)
	for _, field := range values {
		g.Printf(fmt.Sprintf(unmarshall_field, strings.ToLower(field), field, strings.ToLower(field)))
	}
	g.Printf(unmarshall_footer)
}


// format returns the gofmt-ed contents of the Generator's buffer.
func (g *Generator) format() ([]byte, error) {
	return format.Source(g.buf.Bytes())
}

// genDecl processes one declaration clause.
func (f *File) genDecl(node ast.Node) bool {
	decl, ok := node.(*ast.GenDecl)
	if !ok || decl.Tok != token.TYPE {
		// We only care about const declarations.
		return true
	}
	// The name of the type of the constants we are declaring.
	// Can change if this is a multi-element declaration.
	typ := ""

	// Loop over the elements of the declaration. Each element is a TypeSpec:
	// a list of names possibly followed by a type, possibly followed by values.
	// If the type and value are both missing, we carry down the type (and value,
	// but the "go/types" package takes care of that).
	for _, spec := range decl.Specs {
		tspec := spec.(*ast.TypeSpec)
		//fmt.Printf(tspec.Name.Name)
		// Guaranteed to succeed as this is TYPE.
		stype, ok := tspec.Type.(*ast.StructType)
		if !ok || tspec.Type == nil && len(stype.Fields.List) > 0 {
			continue
		}
		if tspec.Type != nil {
			// "X T". We have a type. Remember it.
			typ = tspec.Name.Name
		}

		if typ != f.typeName {
			// This is not the type we're looking for.
			continue
		}

		// We now have a list of names (from one line of source code) all being
		// declared with the desired type.
		// Grab their names and actual values and store them in f.values.
		for _, field := range stype.Fields.List {
			name := field.Names[0]
			// This dance lets the type checker find the values for us. It's a
			// bit tricky: look up the object declared by the name, find its
			// types.Const, and extract its value.
			obj, ok := f.pkg.defs[name]
			if !ok {
				log.Fatalf("no value for constant %s", name)
			}

			link_name := obj.Name()
			f.values = append(f.values, link_name)


		}
	}
	return false
}