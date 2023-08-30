package model

import (
	"context"
	"github.com/bishopfox/sliver/client/command/generate"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"github.com/bishopfox/sliver/protobuf/commonpb"
	"github.com/bishopfox/sliver/util"
	"net"
	"pear-admin-go/app/util/generateConf"
	"strconv"
	"strings"
	"time"
)

type BuildConf struct {
	Name   string `json:"name"`
	Os     string `json:"os"`
	Arch   string `json:"arch"`
	Dns    string `json:"dns"`
	Http   string `json:"http"`
	Mtls   string `json:"mtls"`
	Format string `json:"format"`
}

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
	limitdomainjoined string `json:"limit-domainjoined"`
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

// 初始化implant配置
func ParseCompileFlags(con *console.SliverConsoleClient, userconf BuildConf) *clientpb.ImplantConfig {
	var name string
	var configs generateConf.Conf
	{
		configs.Os = userconf.Os
		configs.Arch = userconf.Arch
		configs.Name = userconf.Name
		configs.Mtls = userconf.Mtls
		configs.Dns = userconf.Dns
		configs.Format = userconf.Format
	}

	if nameF := configs.Name; nameF != "" {
		name = strings.ToLower(nameF)

		if err := util.AllowedName(name); err != nil {
			con.PrintErrorf("%s\n", err)
			return nil
		}
	}

	c2s := []*clientpb.ImplantC2{}

	mtlsC2F := configs.Mtls
	mtlsC2, err := generate.ParseMTLSc2(mtlsC2F)
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}
	c2s = append(c2s, mtlsC2...)

	wgC2F := configs.Wg
	wgC2, err := generate.ParseWGc2(wgC2F)
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}
	wgKeyExchangePort := uint32(configs.Keyexchange)
	wgTcpCommsPort := uint32(configs.Tcpcomms)

	c2s = append(c2s, wgC2...)

	httpC2F := configs.Http
	httpC2, err := generate.ParseHTTPc2(httpC2F)
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}
	c2s = append(c2s, httpC2...)

	dnsC2F := configs.Dns
	dnsC2, err := generate.ParseDNSc2(dnsC2F)
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}
	c2s = append(c2s, dnsC2...)

	namedPipeC2F := configs.Namepipe
	namedPipeC2, err := generate.ParseNamedPipec2(namedPipeC2F)
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}
	c2s = append(c2s, namedPipeC2...)

	tcpPivotC2F := configs.Tcppivot
	tcpPivotC2, err := generate.ParseTCPPivotc2(tcpPivotC2F)
	if err != nil {
		con.PrintErrorf("%s\n", err.Error())
		return nil
	}
	c2s = append(c2s, tcpPivotC2...)

	var symbolObfuscation bool
	if debug := configs.Debug; debug {
		symbolObfuscation = false
	} else {
		symbolObfuscation = configs.Skipsymbols
		symbolObfuscation = !symbolObfuscation
	}

	if len(mtlsC2) == 0 && len(wgC2) == 0 && len(httpC2) == 0 && len(dnsC2) == 0 && len(namedPipeC2) == 0 && len(tcpPivotC2) == 0 {
		con.PrintErrorf("Must specify at least one of --mtls, --wg, --http, --dns, --named-pipe, or --tcp-pivot\n")
		return nil
	}

	rawCanaries := configs.Canary
	canaryDomains := []string{}
	if 0 < len(rawCanaries) {
		for _, canaryDomain := range strings.Split(rawCanaries, ",") {
			if !strings.HasSuffix(canaryDomain, ".") {
				canaryDomain += "." // Ensure we have the FQDN
			}
			canaryDomains = append(canaryDomains, canaryDomain)
		}
	}

	debug := configs.Debug
	evasion := configs.Evasion
	templateName := configs.Template

	reconnectInterval := int64(configs.Reconnect)
	pollTimeout := int64(configs.Polltimeout)
	maxConnectionErrors := uint32(configs.Maxerrors)

	limitDomainJoined, _ := strconv.ParseBool(configs.Limitdomainjoined)
	limitHostname := configs.Limithostname
	limitUsername := configs.Limitusername
	limitDatetime := configs.Limitdatetime
	limitFileExists := configs.Limitfileexists
	limitLocale := configs.Limitlocale
	debugFile := configs.Debugfile

	isSharedLib := false
	isService := false
	isShellcode := false
	sgnEnabled := false

	format := configs.Format
	runAtLoad := false
	var configFormat clientpb.OutputFormat
	switch format {
	case "exe":
		configFormat = clientpb.OutputFormat_EXECUTABLE
	case "shared":
		configFormat = clientpb.OutputFormat_SHARED_LIB
		isSharedLib = true
		runAtLoad, _ = strconv.ParseBool(configs.Runatload)
	case "shellcode":
		configFormat = clientpb.OutputFormat_SHELLCODE
		isShellcode = true
		sgnEnabled, _ = strconv.ParseBool(configs.Disablesgn)
		sgnEnabled = !sgnEnabled
	case "service":
		configFormat = clientpb.OutputFormat_SERVICE
		isService = true
	default:
		// Default to exe
		configFormat = clientpb.OutputFormat_EXECUTABLE
	}

	targetOSF := configs.Os
	targetOS := strings.ToLower(targetOSF)
	targetArchF := configs.Arch
	targetArch := strings.ToLower(targetArchF)
	targetOS, targetArch = generateConf.GetTargets(targetOS, targetArch, con)
	if targetOS == "" || targetArch == "" {
		return nil
	}
	if configFormat == clientpb.OutputFormat_SHELLCODE && targetOS != "windows" {
		con.PrintErrorf("Shellcode format is currently only supported on Windows\n")
		return nil
	}
	if len(namedPipeC2) > 0 && targetOS != "windows" {
		con.PrintErrorf("Named pipe pivoting can only be used in Windows.")
		return nil
	}

	// Check to see if we can *probably* build the target binary
	if !generateConf.CheckBuildTargetCompatibility(configFormat, targetOS, targetArch, con) {
		return nil
	}

	var tunIP net.IP
	if wg := configs.Wg; wg != "" {
		uniqueWGIP, err := con.Rpc.GenerateUniqueIP(context.Background(), &commonpb.Empty{})
		tunIP = net.ParseIP(uniqueWGIP.IP)
		if err != nil {
			con.PrintErrorf("Failed to generate unique ip for wg peer tun interface")
			return nil
		}
		con.PrintInfof("Generated unique ip for wg peer tun interface: %s\n", tunIP.String())
	}

	netGo := configs.Netgo

	// TODO: Use generics or something to check in a slice
	connectionStrategy := configs.Strategy
	if connectionStrategy != "" && connectionStrategy != "s" && connectionStrategy != "r" && connectionStrategy != "rd" {
		con.PrintErrorf("Invalid connection strategy: %s\n", connectionStrategy)
		return nil
	}

	// Parse Traffic Encoder Args
	httpC2Enabled := 0 < len(httpC2)
	trafficEncodersEnabled, trafficEncoderAssets := generateConf.ParseTrafficEncoderArgs(configs, httpC2Enabled, con)

	config := &clientpb.ImplantConfig{
		GOOS:             targetOS,
		GOARCH:           targetArch,
		Name:             name,
		Debug:            debug,
		Evasion:          evasion,
		SGNEnabled:       sgnEnabled,
		ObfuscateSymbols: symbolObfuscation,
		C2:               c2s,
		CanaryDomains:    canaryDomains,
		TemplateName:     templateName,

		WGPeerTunIP:       tunIP.String(),
		WGKeyExchangePort: wgKeyExchangePort,
		WGTcpCommsPort:    wgTcpCommsPort,

		ConnectionStrategy:  connectionStrategy,
		ReconnectInterval:   reconnectInterval * int64(time.Second),
		PollTimeout:         pollTimeout * int64(time.Second),
		MaxConnectionErrors: maxConnectionErrors,

		LimitDomainJoined: limitDomainJoined,
		LimitHostname:     limitHostname,
		LimitUsername:     limitUsername,
		LimitDatetime:     limitDatetime,
		LimitFileExists:   limitFileExists,
		LimitLocale:       limitLocale,

		Format:      configFormat,
		IsSharedLib: isSharedLib,
		IsService:   isService,
		IsShellcode: isShellcode,

		RunAtLoad:              runAtLoad,
		NetGoEnabled:           netGo,
		TrafficEncodersEnabled: trafficEncodersEnabled,
		Assets:                 trafficEncoderAssets,

		DebugFile: debugFile,
	}

	return config
}
