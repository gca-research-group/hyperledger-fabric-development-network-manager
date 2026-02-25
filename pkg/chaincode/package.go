package chaincode

import (
	"fmt"
	"math"
	"strconv"
)

func (c *Chaincode) Package() error {

	organization := c.config.Organizations[0]

	for _, chaincode := range c.config.Chaincodes {
		version := LoadVersion(chaincode)

		filename := ResolveFilename(chaincode)
		name := ResolveFilenameWithoutExtension(chaincode)
		label := ResolveLabel(name, version)
		basePath := ResolveChaincodePath(name)
		tarfile := ResolveChaincodeTar(chaincode, version)

		isChaincodeUpToDate := c.IsChaincodeUpToDate(organization, basePath, name)
		chaincodeFileExists := c.ChaincodeFileExists(organization, tarfile)

		if chaincodeFileExists && isChaincodeUpToDate {
			continue
		}

		if chaincodeFileExists && !isChaincodeUpToDate {
			currentVersion, err := strconv.ParseFloat(version, 64)
			if err != nil {
				return err
			}

			version = fmt.Sprintf("%d.0", int(math.Round(currentVersion))+1)

			label = ResolveLabel(name, version)
			tarfile = ResolveChaincodeTar(chaincode, version)
		}

		steps := []struct {
			name    string
			message string
			args    []string
		}{
			{"Initialize", "Error when initializing the chaincode module %s: %v", []string{
				"sh", "-c", fmt.Sprintf("cd %s && [ -f go.mod ] || go mod init %s; go mod tidy", basePath, name),
			}},
			{"Package", "Error when packaging the chaincode %s: %v", []string{
				"peer", "lifecycle", "chaincode", "package", tarfile,
				"--path", basePath,
				"--lang", "golang",
				"--label", label,
			}},
			{"Compute Checksum", "Error when computing the chaincode checksum %s: %v", []string{
				"sh", "-c", fmt.Sprintf("sha256sum %[1]s/%[2]s > %[3]s", basePath, filename, ResolveChecksum(name)),
			}},
			{"Compute Version", "Error when computing the chaincode version %s: %v", []string{"sh", "-c", fmt.Sprintf("echo %s > %s/version", version, basePath)}},
		}

		for _, step := range steps {
			fmt.Printf(">>> Step: %s\n", step.name)
			_, err := c.ExecInTools(organization, step.args)
			if err != nil {
				return fmt.Errorf(step.message, name, err)
			}
		}
	}

	return nil
}
