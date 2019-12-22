package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/golang-migrate/migrate/v4/source"
	"github.com/markbates/pkger"
)

type PkgerDriver struct {
	migrations *source.Migrations
}

func init() {
	source.Register("pkger", &PkgerDriver{})
}

func WithInstance() (source.Driver, error) {
	d := &PkgerDriver{
		migrations: source.NewMigrations(),
	}
	err := pkger.Walk("/migrations", func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !strings.HasSuffix(path, ".sql") {
			return nil
		}
		mig, err := source.DefaultParse(info.Name())
		if err != nil {
			fmt.Printf("Failed to parse %s: %s", path, err.Error())
			return err
		}
		d.migrations.Append(mig)
		return nil
	})
	return d, err
}

func (d *PkgerDriver) Open(path string) (source.Driver, error) {
	return WithInstance()
}

func (d *PkgerDriver) Close() error {
	return nil
}

func (d *PkgerDriver) First() (version uint, err error) {
	v, ok := d.migrations.First()
	if ok {
		return v, nil
	}
	return 0, &os.PathError{Op: "first", Err: os.ErrNotExist}
}

func (d *PkgerDriver) Next(curr uint) (version uint, err error) {
	v, ok := d.migrations.Next(curr)
	if ok {
		return v, nil
	}
	return 0, &os.PathError{Op: "next", Err: os.ErrNotExist}
}

func (d *PkgerDriver) Prev(curr uint) (version uint, err error) {
	v, ok := d.migrations.Prev(curr)
	if ok {
		return v, nil
	}
	return 0, &os.PathError{Op: "prev", Err: os.ErrNotExist}
}

func (d *PkgerDriver) ReadUp(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := d.migrations.Up(version); ok {
		fp, err := pkger.Open("/migrations/" + m.Raw)
		if err != nil {
			return nil, "", err
		}
		return ioutil.NopCloser(fp), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Err: os.ErrNotExist}
}

func (d *PkgerDriver) ReadDown(version uint) (r io.ReadCloser, identifier string, err error) {
	if m, ok := d.migrations.Down(version); ok {
		fp, err := pkger.Open("/migrations/" + m.Raw)
		if err != nil {
			return nil, "", err
		}
		return ioutil.NopCloser(fp), m.Identifier, nil
	}
	return nil, "", &os.PathError{Op: fmt.Sprintf("read version %v", version), Err: os.ErrNotExist}
}
