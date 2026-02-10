package fabric

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

const (
	caTlsCertPath  = "/etc/hyperledger/fabric-ca-server/ca-cert.pem"
	peerOrgPath    = "/etc/hyperledger/organizations/peerOrganizations"
	ordererOrgPath = "/etc/hyperledger/organizations/ordererOrganizations"
)

func (f *Fabric) GenerateIdentityCertificates() error {
	for _, organization := range f.config.Organizations {
		steps := []struct {
			name string
			fn   func(organization pkg.Organization) error
		}{
			{"Enroll CA Admin", f.enrollCAadmin},
			{"Generate Config YAML", f.generateConfigYaml},
			{"Copy Peer CA Certs", f.copyPeersCACertificates},
			{"Copy Orderer CA Certs", f.copyOrderersCACertificates},
			{"Register Peers", f.registerPeers},
			{"Register Orderers", f.registerOrderes},
			{"Register User", f.registerUser},
			{"Register Org Admin", f.registerOrgAdmin},

			{"Generate Peers MSP", f.generatePeersMSP},
			{"Generate Peer User MSP", f.generatePeerUserMSP},
			{"Generate Peer Org Admin MSP", f.generatePeerOrgAdminMSP},

			{"Generate Orderers MSP", f.generateOrderersMSP},
			{"Generate Orderer Org Admin MSP", f.generateOrdererOrgAdminMSP},
			{"Generate Orderer User MSP", f.generateOrdererUserMSP},

			{"Generate Peer TLS Certs", f.generatePeerTlsCertificates},
			{"Generate Orderer TLS Certs", f.generateOrdererTlsCertificates},
		}

		for _, step := range steps {
			fmt.Printf(">>> Step: %s\n", step.name)
			if err := step.fn(organization); err != nil {
				return fmt.Errorf("failed at step %s: %w", step.name, err)
			}
		}
	}

	return nil
}

func (f *Fabric) execInCA(domain string, script string) error {
	caContainer := fmt.Sprintf("ca.%s", domain)
	args := []string{"exec", caContainer, "sh", "-c", script}
	return f.executor.ExecCommand("docker", args...)
}

func getOrgBaseDir(domain string, orgType string) string {
	if orgType == "peer" {
		return fmt.Sprintf("%s/%s", peerOrgPath, domain)
	}
	return fmt.Sprintf("%s/%s", ordererOrgPath, domain)
}

func (f *Fabric) enrollCAadmin(organization pkg.Organization) error {
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"
	u := "https://admin:adminpw@localhost:7054"
	caName := fmt.Sprintf("ca.%s", organization.Domain)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"-u", u,
		"--caname", caName,
		"--tls.certfiles", tls,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when enrolling the ca admin for the organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) generateConfigYaml(organization pkg.Organization) error {
	template := `
NodeOUs:
  Enable: true
  ClientOUIdentifier:
    Certificate: cacerts/ca.%[1]s-cert.pem
    OrganizationalUnitIdentifier: client
  PeerOUIdentifier:
    Certificate: cacerts/ca.%[1]s-cert.pem
    OrganizationalUnitIdentifier: peer
  AdminOUIdentifier:
    Certificate: cacerts/ca.%[1]s-cert.pem
    OrganizationalUnitIdentifier: admin
  OrdererOUIdentifier:
    Certificate: cacerts/ca.%[1]s-cert.pem
    OrganizationalUnitIdentifier: orderer`

	content := fmt.Sprintf(template, organization.Domain)

	for _, path := range []string{fmt.Sprintf("%s/%s", peerOrgPath, organization.Domain), fmt.Sprintf("%s/%s", ordererOrgPath, organization.Domain)} {
		script := fmt.Sprintf("mkdir -p '%[1]s/msp' && cat <<EOF > %[1]s/msp/config.yaml\n%[2]s\nEOF", path, content)

		if err := f.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when creating the config.yaml for organization %s: %v", organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) copyCACertificates(organization pkg.Organization, orgType string) error {
	basePath := getOrgBaseDir(organization.Domain, orgType)
	scripts := []string{
		fmt.Sprintf("mkdir -p '%[1]s/msp/tlscacerts' && cp '%[2]s' '%[1]s/msp/tlscacerts/ca.crt'", basePath, caTlsCertPath),
		fmt.Sprintf("mkdir -p '%[1]s/tlsca' && cp '%[2]s' '%[1]s/tlsca/tlsca.%[3]s-cert.pem'", basePath, caTlsCertPath, organization.Domain),
		fmt.Sprintf("mkdir -p '%[1]s/ca' && cp '%[2]s' '%[1]s/ca/ca.%[3]s-cert.pem'", basePath, caTlsCertPath, organization.Domain),
		fmt.Sprintf("mkdir -p '%[1]s/msp/cacerts' && cp '%[2]s' '%[1]s/msp/cacerts/ca.%[3]s-cert.pem'", basePath, caTlsCertPath, organization.Domain),
	}

	for _, script := range scripts {
		if err := f.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when executing the script %s for organization %s: %v", script, organization.Name, err)
		}
	}
	return nil
}

func (f *Fabric) copyPeersCACertificates(organization pkg.Organization) error {
	return f.copyCACertificates(organization, "peer")
}

func (f *Fabric) copyOrderersCACertificates(organization pkg.Organization) error {
	return f.copyCACertificates(organization, "orderer")
}

func (f *Fabric) registerPeers(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)

	for i := range organization.Peers {
		id := fmt.Sprintf("peer%d", i)

		args := []string{
			"exec", caName,
			"fabric-ca-client", "register",
			"--caname", caName,
			"--id.name", id,
			"--id.secret", id + "pw",
			"--id.type", "peer",
			"--tls.certfiles", caTlsCertPath,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering the peer %d for organization %s: %v", i, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) registerOrderes(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)

	for _, orderer := range organization.Orderers {
		id := strings.ToLower(orderer.Hostname)
		args := []string{
			"exec", caName,
			"fabric-ca-client", "register",
			"--caname", caName,
			"--id.name", id,
			"--id.secret", id + "pw",
			"--id.type", "orderer",
			"--tls.certfiles", caTlsCertPath,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering the orderer %s for organization %s: %v", orderer.Hostname, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) registerUser(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "register",
		"--caname", caName,
		"--id.name", "user1",
		"--id.secret", "user1pw",
		"--id.type", "client",
		"--tls.certfiles", caTlsCertPath,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when registering the user for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) registerOrgAdmin(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "register",
		"--caname", caName,
		"--id.name", "orgadmin",
		"--id.secret", "orgadminpw",
		"--id.type", "admin",
		"--tls.certfiles", caTlsCertPath,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when registering the admin for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) generateMSP(organization pkg.Organization, mspPath string, id string) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	u := fmt.Sprintf("https://%[1]s:%[1]spw@localhost:7054", id)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"-u", u,
		"--caname", caName,
		"--tls.certfiles", caTlsCertPath,
		"-M", mspPath,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when enrolling the %s for organization %s: %v", id, organization.Name, err)
	}

	scripts := []string{
		fmt.Sprintf("mv '%s'* '%s'", fmt.Sprintf("%s/cacerts/", mspPath), fmt.Sprintf("%s/cacerts/ca.%s-cert.pem", mspPath, organization.Domain)),
		fmt.Sprintf("mv '%s'* '%s'", fmt.Sprintf("%s/keystore/", mspPath), fmt.Sprintf("%s/keystore/priv_sk", mspPath)),
		fmt.Sprintf("cp '%s' '%s'", fmt.Sprintf("%s/%s/msp/config.yaml", peerOrgPath, organization.Domain), fmt.Sprintf("%s/config.yaml", mspPath)),
	}

	for _, script := range scripts {
		if err := f.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when copying the config.yaml to the %s for organization %s: %v", id, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) generatePeersMSP(organization pkg.Organization) error {
	for i := range organization.Peers {
		id := fmt.Sprintf("peer%d", i)
		mspPath := fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s.%[2]s/msp", peerOrgPath, organization.Domain, id)

		if err := f.generateMSP(organization, mspPath, id); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fabric) generatePeerUserMSP(organization pkg.Organization) error {
	id := "user1"
	mspPath := fmt.Sprintf("%[1]s/%[2]s/users/User1@%[2]s/msp", peerOrgPath, organization.Domain)

	if err := f.generateMSP(organization, mspPath, id); err != nil {
		return err
	}

	return nil
}

func (f *Fabric) generatePeerOrgAdminMSP(organization pkg.Organization) error {
	id := "orgadmin"
	mspPath := fmt.Sprintf("%[1]s/%[2]s/users/Admin@%[2]s/msp", peerOrgPath, organization.Domain)

	if err := f.generateMSP(organization, mspPath, id); err != nil {
		return err
	}

	return nil
}

func (f *Fabric) generateOrderersMSP(organization pkg.Organization) error {
	for _, orderer := range organization.Orderers {
		id := strings.ToLower(orderer.Hostname)
		mspPath := fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s.%[2]s/msp", ordererOrgPath, organization.Domain, id)

		if err := f.generateMSP(organization, mspPath, id); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fabric) generateOrdererUserMSP(organization pkg.Organization) error {
	id := "user1"
	mspPath := fmt.Sprintf("%[1]s/%[2]s/users/User1@%[2]s/msp", ordererOrgPath, organization.Domain)

	if err := f.generateMSP(organization, mspPath, id); err != nil {
		return err
	}

	return nil
}

func (f *Fabric) generateOrdererOrgAdminMSP(organization pkg.Organization) error {
	id := "orgadmin"
	mspPath := fmt.Sprintf("%[1]s/%[2]s/users/Admin@%[2]s/msp", ordererOrgPath, organization.Domain)

	if err := f.generateMSP(organization, mspPath, id); err != nil {
		return err
	}

	return nil
}

func (f *Fabric) generateTLS(organization pkg.Organization, tlsPath string, id string) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	u := fmt.Sprintf("https://%[1]s:%[1]spw@localhost:7054", id)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"--caname", caName,
		"-u", u,
		"-M", tlsPath,
		"--enrollment.profile", "tls",
		"--csr.hosts", fmt.Sprintf("%s.%s", id, organization.Domain), "--csr.hosts", "localhost",
		"--tls.certfiles", caTlsCertPath,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when genereting the tls certificates of the %s for organization %s: %v", id, organization.Name, err)
	}

	scripts := []string{
		fmt.Sprintf("cp '%[1]s/tlscacerts/'* '%[1]s/ca.crt'", tlsPath),
		fmt.Sprintf("cp '%[1]s/signcerts/'* '%[1]s/server.crt'", tlsPath),
		fmt.Sprintf("cp '%[1]s/keystore/'* '%[1]s/server.key'", tlsPath),
	}

	for _, script := range scripts {
		if err := f.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when copying the tls certificate of the %s for organization %s: %v", id, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) generatePeerTlsCertificates(organization pkg.Organization) error {
	for i := range organization.Peers {
		id := fmt.Sprintf("peer%d", i)
		tlsPath := fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s.%[2]s/tls", peerOrgPath, organization.Domain, id)

		if err := f.generateTLS(organization, tlsPath, id); err != nil {
			return err
		}
	}

	return nil
}

func (f *Fabric) generateOrdererTlsCertificates(organization pkg.Organization) error {
	for _, orderer := range organization.Orderers {
		id := strings.ToLower(orderer.Hostname)
		tlsPath := fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s.%[2]s/tls", ordererOrgPath, organization.Domain, id)

		if err := f.generateTLS(organization, tlsPath, id); err != nil {
			return err
		}
	}

	return nil
}
