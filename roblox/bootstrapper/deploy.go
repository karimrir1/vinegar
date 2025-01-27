package bootstrapper

import (
	"errors"
	"fmt"
	"log"

	"github.com/vinegarhq/vinegar/util"
)

func (p *Package) Verify(src string) error {
	log.Printf("Verifying Package %s (%s)", p.Name, p.Checksum)

	if err := util.VerifyFileMD5(src, p.Checksum); err != nil {
		return fmt.Errorf("verify package %s: %w", p.Name, err)
	}

	return nil
}

func (p *Package) Download(dest, deployURL string) error {
	if err := p.Verify(dest); err == nil {
		log.Printf("Package %s is already downloaded", p.Name)
		return nil
	}

	log.Printf("Downloading Package %s (%s)", p.Name, dest)

	if err := util.Download(deployURL+"-"+p.Name, dest); err != nil {
		return fmt.Errorf("download package %s (%s): %w", p.Name, dest, err)
	}

	return p.Verify(dest)
}

func (p *Package) Fetch(dest, deployURL string) error {
	err := p.Download(dest, deployURL)
	if err == nil {
		return nil
	}

	log.Printf("Failed to fetch package %s: %s, retrying...", p.Name, errors.Unwrap(err))

	return p.Download(dest, deployURL)
}

func (p *Package) Extract(src, dest string) error {
	if err := extract(src, dest); err != nil {
		return fmt.Errorf("extract package %s (%s): %w", p.Name, src, err)
	}

	log.Printf("Extracted Package %s (%s)", p.Name, p.Checksum)
	return nil
}
