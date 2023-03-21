package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"os"

	"github.com/fatih/color"
	"github.com/go-faster/errors"

	"example/internal/oas"
)

func run(ctx context.Context) error {
	var arg struct {
		BaseURL string
		ID      int64
	}
	flag.StringVar(&arg.BaseURL, "url", "http://localhost:8080", "target server url")
	flag.Int64Var(&arg.ID, "id", 1337, "pet id to request")
	flag.Parse()

	client, err := oas.NewClient(arg.BaseURL)
	if err != nil {
		return errors.Wrap(err, "client")
	}

	res, err := client.GetPetById(ctx, oas.GetPetByIdParams{
		PetId: arg.ID,
	})
	if err != nil {
		return errors.Wrap(err, "get pet")
	}

	switch p := res.(type) {
	case *oas.Pet:
		color.New(color.FgGreen, color.Bold).Print("pet: ")
		data, err := p.MarshalJSON()
		if err != nil {
			return errors.Wrap(err, "marshal")
		}
		var out bytes.Buffer
		if err := json.Indent(&out, data, "", "  "); err != nil {
			return errors.Wrap(err, "indent")
		}
		color.New(color.FgGreen).Println(out.String())
	case *oas.GetPetByIdNotFound:
		return errors.New("not found")
	}

	return nil
}

func main() {
	if err := run(context.Background()); err != nil {
		color.New(color.FgRed, color.Bold).Print("error: ")
		color.New(color.FgRed).Printf("%+v\n", err)
		os.Exit(2)
	} else {
		color.New(color.FgGreen, color.Bold).Print("Success")
	}
}
