package miniocli

type cli struct {
}

func NewMinIOCli() MinIOClient {
	return &cli{}
}

func (c *cli) LoadFile() {

}
