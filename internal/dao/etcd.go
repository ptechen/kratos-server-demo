package dao

var EtcdConf struct{
	Endpoints []string
}

//func init() {
//	if err := paladin.Get("etcd.toml").UnmarshalTOML(&EtcdConf); err != nil {
//		panic(err)
//	}
//}
