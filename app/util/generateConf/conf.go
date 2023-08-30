package generateConf

import (
	"context"
	"fmt"
	"github.com/bishopfox/sliver/client/command/generate"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"

	"runtime"
	"strings"
)

type Conf struct {
	Os                string `json:"os"`
	Arch              string `json:"arch"`
	Wg                string `json:"wg"`
	Name              string `json:"name,omitempty"`
	Mtls              string `json:"mtls,omitempty"`
	Http              string `json:"http"`
	Dns               string `json:"dns"`
	Namepipe          string `json:"name-pipe"`
	Tcppivot          string `json:"tcp-pivot"`
	Canary            string `json:"canary"`
	Limithostname     string `json:"limit-hostname"`
	Limitusername     string `json:"limit-username"`
	Limitdatetime     string `json:"limit-datetime"`
	Limitfileexists   string `json:"limit-fileexists"`
	Limitlocale       string `json:"limit-locale"`
	Debugfile         string `json:"debug-file"`
	Format            string `json:"format"`
	Strategy          string `json:"strategy"`
	Runatload         string `json:"run-at-load"`
	Disablesgn        string `json:"disable-sgn"`
	Template          string `json:"template"`
	Limitdomainjoined string `json:"limit-domainjoined"`
	Trafficencoders   string `json:"traffic-encoders"`

	Netgo       bool `json:"netgo"`
	Evasion     bool `json:"evasion"`
	Debug       bool `json:"debug"`
	Skipsymbols bool `json:"skip-symbols"`

	Reconnect   int64 `json:"reconnect"`
	Polltimeout int64 `json:"poll-timeout"`

	Maxerrors   uint32 `json:"max-errors"`
	Keyexchange uint32 `json:"key-exchange,omitempty"`
	Tcpcomms    uint32 `json:"tcp-comms,omitempty"`
}

func GetTargets(targetOS string, targetArch string, con *console.SliverConsoleClient) (string, string) {
	/* For UX we convert some synonymous terms */
	if targetOS == "darwin" || targetOS == "mac" || targetOS == "macos" || targetOS == "osx" {
		targetOS = "darwin"
	}
	if targetOS == "windows" || targetOS == "win" || targetOS == "shit" {
		targetOS = "windows"
	}
	if targetOS == "linux" || targetOS == "lin" {
		targetOS = "linux"
	}

	if targetArch == "amd64" || targetArch == "x64" || strings.HasPrefix(targetArch, "64") {
		targetArch = "amd64"
	}
	if targetArch == "386" || targetArch == "x86" || strings.HasPrefix(targetArch, "32") {
		targetArch = "386"
	}

	target := fmt.Sprintf("%s/%s", targetOS, targetArch)
	if _, ok := generate.SupportedCompilerTargets[target]; !ok {
		con.Printf("⚠️  Unsupported compiler target %s%s%s, but we can try to compile a generic implant.\n",
			console.Bold, target, console.Normal,
		)
		con.Printf("⚠️  Generic implants do not support all commands/features.\n")
		//prompt := &survey.Confirm{Message: "Attempt to build generic implant?"}
		//var confirm bool
		//survey.AskOne(prompt, &confirm)
		//if !confirm {
		//	return "", ""
		//}
		return "", ""
	}

	return targetOS, targetArch
}

func CheckBuildTargetCompatibility(format clientpb.OutputFormat, targetOS string, targetArch string, con *console.SliverConsoleClient) bool {
	if format == clientpb.OutputFormat_EXECUTABLE {
		return true // We don't need cross-compilers when targeting EXECUTABLE formats
	}

	compilers, err := con.Rpc.GetCompiler(context.Background(), &commonpb.Empty{})
	if err != nil {
		fmt.Printf("Failed to check target compatibility: %s\n", err)
		return true
	}

	if runtime.GOOS != "windows" && targetOS == "windows" {
		if !hasCC(targetOS, targetArch, compilers.CrossCompilers) {
			return warnMissingCrossCompiler(format, targetOS, targetArch, con)
		}
	}

	if runtime.GOOS != "darwin" && targetOS == "darwin" {
		if !hasCC(targetOS, targetArch, compilers.CrossCompilers) {
			return warnMissingCrossCompiler(format, targetOS, targetArch, con)
		}
	}

	if runtime.GOOS != "linux" && targetOS == "linux" {
		if !hasCC(targetOS, targetArch, compilers.CrossCompilers) {
			return warnMissingCrossCompiler(format, targetOS, targetArch, con)
		}
	}

	return true
}

func hasCC(targetOS string, targetArch string, crossCompilers []*clientpb.CrossCompiler) bool {
	for _, cc := range crossCompilers {
		if cc.GetTargetGOOS() == targetOS && cc.GetTargetGOARCH() == targetArch {
			return true
		}
	}
	return false
}

const (
	crossCompilerInfoURL = "https://github.com/BishopFox/sliver/wiki/Cross-Compiling-Implants"
)

func warnMissingCrossCompiler(format clientpb.OutputFormat, targetOS string, targetArch string, con *console.SliverConsoleClient) bool {
	fmt.Printf("Missing cross-compiler for %s on %s/%s\n", nameOfOutputFormat(format), targetOS, targetArch)
	switch targetOS {
	case "windows":
		fmt.Printf("The server cannot find an installation of mingw")
	case "darwin":
		fmt.Printf("The server cannot find an installation of osxcross")
	case "linux":
		fmt.Printf("The server cannot find an installation of musl-cross")
	}
	con.PrintWarnf("For more information please read %s\n", crossCompilerInfoURL)

	confirm := false
	//prompt := &survey.Confirm{Message: "Try to compile anyways (will likely fail)?"}
	//survey.AskOne(prompt, &confirm, nil)
	return confirm
}
func nameOfOutputFormat(value clientpb.OutputFormat) string {
	switch value {
	case clientpb.OutputFormat_EXECUTABLE:
		return "Executable"
	case clientpb.OutputFormat_SERVICE:
		return "Service"
	case clientpb.OutputFormat_SHARED_LIB:
		return "Shared Library"
	case clientpb.OutputFormat_SHELLCODE:
		return "Shellcode"
	default:
		return "Unknown"
	}
}

// parseTrafficEncoderArgs - parses the traffic encoder args and returns a bool indicating if traffic encoders are enabled
func ParseTrafficEncoderArgs(cmd Conf, httpC2Enabled bool, con *console.SliverConsoleClient) (bool, []*commonpb.File) {
	trafficEncoders := cmd.Trafficencoders
	encoders := []*commonpb.File{}
	if trafficEncoders != "" {
		if !httpC2Enabled {
			con.PrintWarnf("Traffic encoders are only supported with HTTP C2, flag will be ignored\n")
			return false, encoders
		}
		enabledEncoders := strings.Split(trafficEncoders, ",")
		for _, encoder := range enabledEncoders {
			if !strings.HasSuffix(encoder, ".wasm") {
				encoder += ".wasm"
			}
			encoders = append(encoders, &commonpb.File{Name: encoder})
		}
		return true, encoders
	}
	return false, encoders
}
