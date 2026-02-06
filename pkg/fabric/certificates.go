package fabric

import (
	"fmt"
	"strings"
)

func (f *Fabric) GenerateCertificates() error {
	for _, organization := range f.config.Organizations {
		fmt.Printf("\n=========== Generating certificates for %s ===========\n", organization.Name)

		name := fmt.Sprintf("%sadmin", strings.ToLower(organization.Name))
		password := fmt.Sprintf("%sadminpw", strings.ToLower(organization.Name))
		caName := fmt.Sprintf("ca.%s", organization.Domain)
		tls := fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server")
		u := "https://admin:adminpw@localhost:7054"

		var args []string

		args = []string{
			"exec", caName, "fabric-ca-client", "enroll",
			"-u", u,
			"--caname", caName,
			"--tls.certfiles", tls,
			"-M", "/etc/hyperledger/fabric-ca-client",
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when generating the certificates for the organization %s: %v", organization.Name, err)
		}

		fmt.Printf("\n=========== Generating user admin certificates for %s ===========\n", organization.Name)

		args = []string{
			"exec", caName, "fabric-ca-client", "register",
			"-u", "https://localhost:7054",
			"--caname", caName,
			"--id.name", name,
			"--id.secret", password,
			"--id.type", "admin",
			"--tls.certfiles", tls,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering admin for the organization %s: %v", organization.Name, err)
		}

		args = []string{
			"exec", caName, "fabric-ca-client", "enroll",
			"-u", u,
			"--caname", caName,
			"--tls.certfiles", tls,
			"-M", fmt.Sprintf("/etc/hyperledger/fabric-ca-client/users/Admin@%s/msp", organization.Domain),
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when enrolling admin for the organization %s: %v", organization.Name, err)
		}

		for i := range organization.Peers {
			fmt.Printf("\n=========== Generating certificates for peer %d of %s ===========\n", i, organization.Name)

			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain), "fabric-ca-client", "register",
				"-u", fmt.Sprintf("https://%s:%s@localhost:7054", name, password),
				"--caname", caName,
				"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server"),
				"--id.name", fmt.Sprintf("peer%d", i),
				"--id.secret", fmt.Sprintf("peer%dpw", i),
				"--id.type", "peer",
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when generating the certificates for the peer %d of organization %s: %v", i, organization.Name, err)
			}

			fmt.Printf("\n=========== Enrolling peer %d of %s ===========\n", i, organization.Name)

			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain), "fabric-ca-client", "enroll",
				"-u", fmt.Sprintf("https://peer%d:peer%dpw@localhost:7054", i, i),
				"--caname", fmt.Sprintf("ca.%s", organization.Domain),
				"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server"),
				"-M", fmt.Sprintf("%s/peers/peer%d.%s/msp", "/etc/hyperledger/fabric-ca-client", i, organization.Domain),
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when enrolling the peer %d of organization %s: %v", i, organization.Name, err)
			}

			fmt.Printf("\n=========== Enrolling TLS peer %d of %s ===========\n", i, organization.Name)

			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain), "fabric-ca-client", "enroll",
				"-u", fmt.Sprintf("https://peer%d:peer%dpw@localhost:7054", i, i),
				"--caname", fmt.Sprintf("ca.%s", organization.Domain),
				"--enrollment.profile", "tls",
				"--csr.hosts", fmt.Sprintf("%s,localhost", fmt.Sprintf("peer%d.%s", i, organization.Domain)),
				"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server"),
				"-M", fmt.Sprintf("%s/peers/peer%d.%s/tls", "/etc/hyperledger/fabric-ca-client", i, organization.Domain),
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when enrolling the peer %d of organization %s: %v", i, organization.Name, err)
			}
		}

		for _, orderer := range organization.Orderers {
			fmt.Printf("\n=========== Generating certificates for the orderer %s of %s ===========\n", orderer.Name, organization.Name)

			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain), "fabric-ca-client", "register",
				"-u", u,
				"--caname", caName,
				"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server"),
				"--id.name", fmt.Sprintf("orderer%s", orderer.Name),
				"--id.secret", fmt.Sprintf("orderer%spw", orderer.Name),
				"--id.type", "orderer",
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when generating the certificates for the orderer %s of organization %s: %v", orderer.Name, organization.Name, err)
			}

			fmt.Printf("\n=========== Enrolling orderer %s of %s ===========\n", orderer.Name, organization.Name)

			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain), "fabric-ca-client", "enroll",
				"-u", fmt.Sprintf("https://orderer%s:orderer%spw@localhost:7054", orderer.Name, orderer.Name),
				"--caname", fmt.Sprintf("ca.%s", organization.Domain),
				"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server"),
				"-M", fmt.Sprintf("%s/orderers/%s.%s/msp", "/etc/hyperledger/fabric-ca-client", orderer.Hostname, organization.Domain),
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when enrolling the orderer %s of organization %s: %v", orderer.Name, organization.Name, err)
			}

			fmt.Printf("\n=========== Enrolling TLS orderer %s of %s ===========\n", orderer.Name, organization.Name)

			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain), "fabric-ca-client", "enroll",
				"-u", fmt.Sprintf("https://orderer%s:orderer%spw@localhost:7054", orderer.Name, orderer.Name),
				"--caname", fmt.Sprintf("ca.%s", organization.Domain),
				"--enrollment.profile", "tls",
				"--csr.hosts", fmt.Sprintf("%s,localhost", fmt.Sprintf("%s.%s", orderer.Hostname, organization.Domain)),
				"--tls.certfiles", fmt.Sprintf("%s/ca-cert.pem", "/etc/hyperledger/fabric-ca-server"),
				"-M", fmt.Sprintf("%s/orderers/%s.%s/tls", "/etc/hyperledger/fabric-ca-client", orderer.Hostname, organization.Domain),
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when enrolling the orderer %s of organization %s: %v", orderer.Name, organization.Name, err)
			}
		}
	}

	return nil
}

func (f *Fabric) BuildFabricFolderStructure() error {
	for _, organization := range f.config.Organizations {
		var args []string

		fmt.Printf("\n=========== Building folder structure for %s ===========\n", organization.Name)

		var folders []string

		basePeerPath := "/etc/hyperledger/crypto-material/peerOrganizations"
		baseOrdererPath := "/etc/hyperledger/crypto-material/ordererOrganizations"

		folders = append(folders, fmt.Sprintf("%s/%s/msp/cacerts", basePeerPath, organization.Domain))
		folders = append(folders, fmt.Sprintf("%s/%s/msp/tlscacerts", basePeerPath, organization.Domain))

		folders = append(folders, fmt.Sprintf("%s/%s/users/Admin@%s/msp/signcerts", basePeerPath, organization.Domain, organization.Domain))
		folders = append(folders, fmt.Sprintf("%s/%s/users/Admin@%s/msp/keystore", basePeerPath, organization.Domain, organization.Domain))
		folders = append(folders, fmt.Sprintf("%s/%s/users/Admin@%s/msp/cacerts", basePeerPath, organization.Domain, organization.Domain))

		for i := range organization.Peers {
			folders = append(folders, fmt.Sprintf("%s/%s/peers/peer%d.%s/msp/signcerts", basePeerPath, organization.Domain, i, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/peers/peer%d.%s/msp/keystore", basePeerPath, organization.Domain, i, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/peers/peer%d.%s/msp/cacerts", basePeerPath, organization.Domain, i, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/peers/peer%d.%s/msp/tlscacerts", basePeerPath, organization.Domain, i, organization.Domain))

			folders = append(folders, fmt.Sprintf("%s/%s/peers/peer%d.%s/tls", basePeerPath, organization.Domain, i, organization.Domain))
		}

		if len(organization.Orderers) > 0 {

			folders = append(folders, fmt.Sprintf("%s/%s/msp/signcerts", baseOrdererPath, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/msp/keystore", baseOrdererPath, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/msp/cacerts", baseOrdererPath, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/msp/tlscacerts", baseOrdererPath, organization.Domain))

			folders = append(folders, fmt.Sprintf("%s/%s/users/Admin@%s/msp/signcerts", baseOrdererPath, organization.Domain, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/users/Admin@%s/msp/keystore", baseOrdererPath, organization.Domain, organization.Domain))
			folders = append(folders, fmt.Sprintf("%s/%s/users/Admin@%s/msp/cacerts", baseOrdererPath, organization.Domain, organization.Domain))

			for _, orderer := range organization.Orderers {
				folders = append(folders, fmt.Sprintf("%s/%s/orderers/%s.%s/msp/signcerts", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
				folders = append(folders, fmt.Sprintf("%s/%s/orderers/%s.%s/msp/keystore", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
				folders = append(folders, fmt.Sprintf("%s/%s/orderers/%s.%s/msp/cacerts", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
				folders = append(folders, fmt.Sprintf("%s/%s/orderers/%s.%s/msp/tlscacerts", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))

				folders = append(folders, fmt.Sprintf("%s/%s/orderers/%s.%s/tls", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
			}

		}

		for _, folder := range folders {
			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain),
				"sh", "-c", fmt.Sprintf("mkdir -p '%s'", folder),
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when creating folder %s for organization %s: %v", folder, organization.Name, err)
			}
		}

		config := `
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.%s-cert.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca.%s-cert.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/ca.%s-cert.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/ca.%s-cert.pem
    OrganizationalUnitIdentifier: orderer`

		config = fmt.Sprintf(config, organization.Domain, organization.Domain, organization.Domain, organization.Domain)

		var scripts []string

		scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/msp/cacerts/ca.%s-cert.pem", basePeerPath, organization.Domain, organization.Domain))
		scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/msp/tlscacerts/tlsca.%s-cert.pem", basePeerPath, organization.Domain, organization.Domain))

		scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/users/Admin@%s/msp/signcerts/* %s/%s/users/Admin@%s/msp/signcerts/peer.%s.pem", organization.Domain, basePeerPath, organization.Domain, organization.Domain, organization.Domain))
		scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/users/Admin@%s/msp/keystore/* %s/%s/users/Admin@%s/msp/keystore/priv_sk", organization.Domain, basePeerPath, organization.Domain, organization.Domain))
		scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/users/Admin@%s/msp/cacerts/ca.%s-cert.pem", basePeerPath, organization.Domain, organization.Domain, organization.Domain))

		for i := range organization.Peers {
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/peers/peer%d.%s/msp/signcerts/* %s/%s/peers/peer%d.%s/msp/signcerts/peer.%s.pem", i, organization.Domain, basePeerPath, organization.Domain, i, organization.Domain, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/peers/peer%d.%s/msp/keystore/* %s/%s/peers/peer%d.%s/msp/keystore/priv_sk", i, organization.Domain, basePeerPath, organization.Domain, i, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-server/ca-cert.pem %s/%s/peers/peer%d.%s/msp/cacerts/ca.%s-cert.pem", basePeerPath, organization.Domain, i, organization.Domain, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-server/ca-cert.pem %s/%s/peers/peer%d.%s/msp/tlscacerts/tlsca.%s-cert.pem", basePeerPath, organization.Domain, i, organization.Domain, organization.Domain))

			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/peers/peer%d.%s/tls/signcerts/* %s/%s/peers/peer%d.%s/tls/server.crt", i, organization.Domain, basePeerPath, organization.Domain, i, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/peers/peer%d.%s/tls/keystore/* %s/%s/peers/peer%d.%s/tls/server.key", i, organization.Domain, basePeerPath, organization.Domain, i, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/peers/peer%d.%s/tls/tlscacerts/* %s/%s/peers/peer%d.%s/tls/ca.crt", i, organization.Domain, basePeerPath, organization.Domain, i, organization.Domain))

			folders := []string{
				fmt.Sprintf("%s/%s/msp/", basePeerPath, organization.Domain),
				fmt.Sprintf("%s/%s/peers/peer%d.%s/msp", basePeerPath, organization.Domain, i, organization.Domain),
				fmt.Sprintf("%s/%s/users/Admin@%s/msp", basePeerPath, organization.Domain, organization.Domain),
			}

			for _, folder := range folders {
				args = []string{
					"exec", fmt.Sprintf("ca.%s", organization.Domain),
					"sh", "-c", fmt.Sprintf("cat <<EOF > %s/config.yaml %s", folder, config),
				}

				if err := f.executor.ExecCommand("docker", args...); err != nil {
					return fmt.Errorf("Error when creating config.yaml for organization %s: %v", organization.Name, err)
				}
			}
		}

		if len(organization.Orderers) > 0 {
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/users/Admin@%s/msp/signcerts/* %s/%s/msp/signcerts/ca.%s-cert.pem", organization.Domain, baseOrdererPath, organization.Domain, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/msp/signcerts/priv_sk", baseOrdererPath, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/msp/cacerts/ca.%s-cert.pem", baseOrdererPath, organization.Domain, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/msp/tlscacerts/tlsca.%s-cert.pem", baseOrdererPath, organization.Domain, organization.Domain))

			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/users/Admin@%s/msp/signcerts/* %s/%s/users/Admin@%s/msp/signcerts/ca.%s.pem", organization.Domain, baseOrdererPath, organization.Domain, organization.Domain, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/users/Admin@%s/msp/keystore/* %s/%s/users/Admin@%s/msp/keystore/priv_sk", organization.Domain, baseOrdererPath, organization.Domain, organization.Domain))
			scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/cacerts/* %s/%s/users/Admin@%s/msp/cacerts/ca.%s-cert.pem", baseOrdererPath, organization.Domain, organization.Domain, organization.Domain))

			for _, orderer := range organization.Orderers {
				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/orderers/%s.%s/msp/signcerts/* %s/%s/orderers/%s.%s/msp/signcerts/%s.%s.pem", orderer.Hostname, organization.Domain, baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain, orderer.Hostname, organization.Domain))
				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/orderers/%s.%s/msp/keystore/* %s/%s/orderers/%s.%s/msp/keystore/priv_sk", orderer.Hostname, organization.Domain, baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-server/ca-cert.pem %s/%s/orderers/%s.%s/msp/cacerts/ca.%s-cert.pem", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain, organization.Domain))
				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-server/ca-cert.pem %s/%s/orderers/%s.%s/msp/tlscacerts/tlsca.%s-cert.pem", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain, organization.Domain))

				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/orderers/%s.%s/tls/signcerts/* %s/%s/orderers/%s.%s/tls/server.crt", orderer.Hostname, organization.Domain, baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/orderers/%s.%s/tls/keystore/* %s/%s/orderers/%s.%s/tls/server.key", orderer.Hostname, organization.Domain, baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))
				scripts = append(scripts, fmt.Sprintf("cp /etc/hyperledger/fabric-ca-client/orderers/%s.%s/tls/tlscacerts/* %s/%s/orderers/%s.%s/tls/ca.crt", orderer.Hostname, organization.Domain, baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain))

				folders := []string{
					fmt.Sprintf("%s/%s/msp/", baseOrdererPath, organization.Domain),
					fmt.Sprintf("%s/%s/orderers/%s.%s/msp", baseOrdererPath, organization.Domain, orderer.Hostname, organization.Domain),
					fmt.Sprintf("%s/%s/users/Admin@%s/msp", baseOrdererPath, organization.Domain, organization.Domain),
				}

				for _, folder := range folders {
					args = []string{
						"exec", fmt.Sprintf("ca.%s", organization.Domain),
						"sh", "-c", fmt.Sprintf("cat <<EOF > %s/config.yaml %s", folder, config),
					}

					if err := f.executor.ExecCommand("docker", args...); err != nil {
						return fmt.Errorf("Error when creating config.yaml for organization %s: %v", organization.Name, err)
					}
				}

			}
		}

		for _, script := range scripts {
			args := []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain),
				"sh", "-c", script,
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when copying certificates for organization %s: %v", organization.Name, err)
			}
		}
	}

	return nil
}
