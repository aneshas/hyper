package hyper

import (
	"context"
	"path"
)

func NewBoilStoreGen(tools GoTools, fs FS) *BoilStoreGen {
	return &BoilStoreGen{
		tools: tools,
		fs:    fs,
	}
}

type BoilStoreGen struct {
	tools GoTools
	fs    FS
}

func (b *BoilStoreGen) GenStore(ctx context.Context, app App, uc UC) error {
	// TODO - Check if the embedded store is boiler store and generate
	// eg. we could provide array of different generators and all of them could check what to generate or multiple can be generated
	// with eg gen --boiler --gorm --empty ...

	// TODO - Check if file exists already - if yes skip

	storeDir := path.Join(app.Dir(), "internal", "boiler")

	err := b.fs.MkAppDir(storeDir)
	if err != nil {
		return err
	}

	return nil
}
