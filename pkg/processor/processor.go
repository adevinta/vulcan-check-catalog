package processor

import (
	"fmt"
	"io/ioutil"

	"github.com/adevinta/vulcan-check-catalog/pkg/helpers"
	"github.com/adevinta/vulcan-check-catalog/pkg/manifest"
	"github.com/adevinta/vulcan-check-catalog/pkg/model"
	log "github.com/sirupsen/logrus"
)

const (
	manifestFile    = "manifest.toml"
	DefaultImageTag = "stable"
)

type Runner struct {
	logger          *log.Logger
	path            string
	registryBaseURL string
	imageTag        string
}

type Processor interface {
	Run() (model.Checktypes, error)
}

func New(path, registryBaseURL, imageTag string, logger *log.Logger) Runner {
	tag := imageTag
	if tag == "" {
		tag = DefaultImageTag
	}
	return Runner{
		logger:          logger,
		path:            path,
		registryBaseURL: registryBaseURL,
		imageTag:        tag,
	}
}

func (r *Runner) Run() (model.Checktypes, error) {
	files, err := ioutil.ReadDir(r.path)
	if err != nil {
		return model.Checktypes{}, err
	}

	checktypes := model.Checktypes{}
	for _, f := range files {
		fpath := fmt.Sprintf("%s/%s", r.path, f.Name())
		isDir, err := helpers.IsExistingDirOrFile(fpath)
		if err != nil {
			return model.Checktypes{}, err
		}
		// We expect a path to a directory with one or more directories which
		// each one of these dectories is a check.
		if !isDir {
			continue
		}
		manifestFilePath := fmt.Sprintf("%s/%s", fpath, manifestFile)
		isDir, err = helpers.IsExistingDirOrFile(manifestFilePath)
		if err != nil {
			r.logger.Warnf("path [%s] does not look like a check folder", fpath)
			continue
		}
		if isDir {
			continue
		}
		// Potential check folder found.
		checkName := f.Name()
		md, err := manifest.Read(manifestFilePath)
		checktype := model.Checktype{
			Name:         checkName,
			Image:        fmt.Sprintf("%s/%s:%s", r.registryBaseURL, checkName, r.imageTag),
			Description:  md.Description,
			Timeout:      md.Timeout,
			RequiredVars: md.RequiredVars,
			Assets:       md.AssetTypes,
		}
		if md.Options != "" {
			checktype.Options = md.Options
		}
		if checktype.Assets == nil {
			checktype.Assets = []string{}
		}
		checktypes.Checktype = append(checktypes.Checktype, checktype)
	}

	return checktypes, nil
}
