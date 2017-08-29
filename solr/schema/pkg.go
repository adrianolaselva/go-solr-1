/*
Package schema provides core schema management and schema generation using reflection.
Schema related structs (Schema, Field etc.) and filedtype are defined in common package

Examples of struct filed tags and their meanings

	// First value is treated as field name, if it is "-", this field is ignored
	Foo      string `solr:"foo"`
	IgnoreMe string `solr:"-"`

	// Multiple values are separated by ",", you can leave first value empty, we will use the field name or json tag (if presented)
	Foo      string `json:"foo" solr:",type=string,docValues=true,indexed=false,stored=true,multiValued=false,required=true"`
*/
package schema

import (
	"context"
	"fmt"

	"github.com/at15/go-solr/pkg/common"
	"github.com/at15/go-solr/pkg/internal"
	"github.com/at15/go-solr/pkg/util"
	"github.com/pkg/errors"
)

var log = util.Logger.RegisterPkg()

const (
	baseURLTmpl = "/solr/%s/schema" // TODO: it seems this does not need trailing slash
)

type Service struct {
	client *internal.Client
	meta   *common.Schema

	core    common.Core
	baseURL string
}

func New(client *internal.Client, core common.Core) *Service {
	s := &Service{
		client: client,
	}
	s.SetCore(core)
	return s
}

func (svc *Service) SetCore(core common.Core) {
	svc.core = core
	svc.baseURL = fmt.Sprintf(baseURLTmpl, core.Name)
}

// GET localhost:8983/solr/demo/schema?wt=json
func (svc *Service) Get(ctx context.Context) (*common.Schema, error) {
	res := &common.SchemaResponse{}
	if _, err := svc.client.Get(ctx, svc.baseURL, res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("solr: can't get core %s schema", svc.core.Name))
	}
	// cache the schema
	svc.meta = res.Schema
	return res.Schema, nil
}