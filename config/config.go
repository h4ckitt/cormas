package config

const secretkey = "Mr.RavandraIsTheB3sTUpw0rKCl13nt"

func GetConfig() Config {
	return Config{
		SecretKey: secretkey,
	}
}
