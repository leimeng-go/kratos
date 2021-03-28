package add

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/mod/modfile"
)

// CmdAdd represents the add command.
var CmdAdd = &cobra.Command{
	Use:   "add",
	Short: "Add a proto API template",
	Long:  "Add a proto API template. Example: kratos add helloworld/v1/hello.proto",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	// kratos add api/user/v1/user.proto
	input := args[0]
	n := strings.LastIndex(input, "/")
	// api/user/v1
	path := input[:n]
	// user.proto
	fileName := input[n+1:]
	// api.user.v1
	pkgName := strings.ReplaceAll(path, "/", ".")

	p := &Proto{
		// user.proto
		Name:        fileName,
		// api/user/v1
		Path:        path,
		// api.user
		Package:     pkgName,
		// helloworld/api/user/v1;v1
		GoPackage:   goPackage(path),
		// api/user/v1
		JavaPackage: javaPackage(pkgName),
		// user
		Service:     serviceName(fileName),
	}
	if err := p.Generate(); err != nil {
		fmt.Println(err)
		return
	}
}

func modName() string {
	modBytes, err := ioutil.ReadFile("go.mod")
	if err != nil {
		if modBytes, err = ioutil.ReadFile("../go.mod"); err != nil {
			return ""
		}
	}
	return modfile.ModulePath(modBytes)
}

func goPackage(path string) string {
	s := strings.Split(path, "/")
	return modName() + "/" + path + ";" + s[len(s)-1]
}

func javaPackage(name string) string {
	return name
}

func serviceName(name string) string {
	return unexport(strings.Split(name, ".")[0])
}

func unexport(s string) string { return strings.ToUpper(s[:1]) + s[1:] }
