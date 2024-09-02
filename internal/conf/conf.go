package conf

import (
	"github.com/Zkeai/go_template/common/database"
	chttp "github.com/Zkeai/go_template/common/net/cttp"
	"github.com/Zkeai/go_template/common/redis"
)

type Conf struct {
	Server *chttp.Config    `yaml:"server"`
	DB     *database.Config `yaml:"db"`
	Rdb    *redis.Config    `yaml:"redis"`
}
