package service

import (
	"context"
	"github.com/bishopfox/sliver/client/console"
	"github.com/bishopfox/sliver/protobuf/clientpb"
	"os"
	"path/filepath"
	rpcs "pear-admin-go/app/core/rpc"
	"pear-admin-go/app/model"
	"strings"
)

func Build(conf model.BuildConf) (string, error) {
	var e error

	/// conf
	con := rpcs.GetInstance().GetCon()
	config := model.ParseCompileFlags(con, conf)

	generated, err := con.Rpc.Generate(context.Background(), &clientpb.GenerateReq{
		Config: config,
	})
	if err != nil {
		return "", e
	}

	if len(generated.File.Data) == 0 {
		return "", e
	}
	fileData := generated.File.Data

	if len(generated.File.Data) == 0 {
		//con.PrintErrorf("Build failed, no file data\n")
		return conf.Name, e
	}

	if config.IsShellcode {
		if !config.SGNEnabled {
			//con.PrintErrorf("Shikata ga nai encoder is %sdisabled%s\n", console.Bold, console.Normal)
			return "", nil
		} else {
			//con.PrintInfof("Encoding shellcode with shikata ga nai ... ")
			resp, err := con.Rpc.ShellcodeEncoder(context.Background(), &clientpb.ShellcodeEncodeReq{
				Encoder:      clientpb.ShellcodeEncoder_SHIKATA_GA_NAI,
				Architecture: config.GOARCH,
				Iterations:   1,
				BadChars:     []byte{},
				Data:         fileData,
			})
			if err != nil {
				//con.PrintErrorf("%s\n", err)
				return "", e
			} else {
				//con.Printf("success!\n")
				fileData = resp.GetData()
			}
		}
	}

	saveTo, err := saveLocation(conf.Name, generated.File.Name, con)
	if err != nil {
		return config.Name, err
	}

	err = os.WriteFile(saveTo, fileData, 0o700)
	if err != nil {
		con.PrintErrorf("Failed to write to: %s\n", saveTo)
		return config.Name, err
	}
	con.PrintInfof("Implant saved to %s\n", saveTo)

	return saveTo, nil
}

func expandPath(path string) string {
	// unless path starts with ~
	if len(path) == 0 || path[0] != 126 {
		return path
	}

	return filepath.Join(os.Getenv("HOME"), path[1:])
}

func saveLocation(save, DefaultName string, con *console.SliverConsoleClient) (string, error) {
	var saveTo string

	//BasePath := "static/implants"

	if save == "" {
		save, _ = os.Getwd()
	}
	save = expandPath(save)
	fi, err := os.Stat(save)
	if os.IsNotExist(err) {
		con.Printf("%s does not exist\n", save)
		if strings.HasSuffix(save, "/") {
			con.Printf("%s is dir\n", save)
			os.MkdirAll(save, 0o700)
			saveTo, _ = filepath.Abs(filepath.Join(saveTo, "implants", DefaultName))
		} else {
			con.Printf("%s is not dir\n", save)
			saveDir := filepath.Dir(save)
			_, err := os.Stat(saveTo)
			if os.IsNotExist(err) {
				os.MkdirAll(saveDir, 0o700)
			}
			saveTo, _ = filepath.Abs(filepath.Join("implants", save))
		}
	} else {
		if fi.IsDir() {
			saveTo, _ = filepath.Abs(filepath.Join(save, "implants", DefaultName))
		} else {
			//con.PrintInfof("%s is not dir\n", save)
			//prompt := &survey.Confirm{Message: "Overwrite existing file?"}
			//var confirm bool
			//survey.AskOne(prompt, &confirm)
			//if !confirm {
			//	return "", errors.New("file already exists")
			//}
			saveTo, _ = filepath.Abs(save)
		}
	}
	return saveTo, nil
}
