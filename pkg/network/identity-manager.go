package network

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/executor"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/internal/file"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/compose"
	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg/config"
)

const (
	caTlsCertPath  = "/etc/hyperledger/fabric-ca-server/ca-cert.pem"
	peerOrgPath    = "/etc/hyperledger/organizations/peerOrganizations"
	ordererOrgPath = "/etc/hyperledger/organizations/ordererOrganizations"
)

type IdentityManager struct {
	config   config.Config
	executor executor.Executor
}

func NewIdentityManager(config config.Config, executor executor.Executor) *IdentityManager {
	return &IdentityManager{
		config:   config,
		executor: executor,
	}
}

func (im *IdentityManager) GenerateAll() error {
	for _, organization := range im.config.Organizations {
		steps := []struct {
			name string
			fn   func(organization config.Organization) error
		}{
			{"Enroll CA Admin", im.enrollCAadmin},
			{"Generate Config YAML", im.generateConfigYaml},
			{"Copy Peer CA Certs", im.copyPeersCACertificates},
			{"Copy Orderer CA Certs", im.copyOrderersCACertificates},
			{"Register Peers", im.registerPeers},
			{"Register Orderers", im.registerOrderers},
			{"Register User", im.registerUser},
			{"Register Org Admin", im.registerOrgAdmin},

			{"Generate Peers MSP", im.generatePeersMSP},
			{"Generate Peer User MSP", im.generatePeerUserMSP},
			{"Generate Peer Org Admin MSP", im.generatePeerOrgAdminMSP},

			{"Generate Orderers MSP", im.generateOrderersMSP},
			{"Generate Orderer Org Admin MSP", im.generateOrdererOrgAdminMSP},
			{"Generate Orderer User MSP", im.generateOrdererUserMSP},

			{"Generate Peer TLS Certs", im.generatePeerTlsCertificates},
			{"Generate Orderer TLS Certs", im.generateOrdererTlsCertificates},
		}

		for _, step := range steps {
			fmt.Printf(">>> Step: %s\n", step.name)
			if err := step.fn(organization); err != nil {
				return fmt.Errorf("failed at step %s: %w", step.name, err)
			}
		}
	}

	if err := im.shareTlsCertificates(); err != nil {
		return err
	}

	return nil
}

func (im *IdentityManager) execInCA(domain string, script string) error {
	caName := compose.ResolveCertificateAuthorityContainerName(domain)
	args := []string{"exec", caName, "sh", "-c", script}
	return im.executor.ExecCommand("docker", args...)
}

func getOrgBaseDir(domain string, orgType string) string {
	if orgType == "peer" {
		return fmt.Sprintf("%s/%s", peerOrgPath, domain)
	}
	return fmt.Sprintf("%s/%s", ordererOrgPath, domain)
}

func (im *IdentityManager) enrollCAadmin(organization config.Organization) error {
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"
	u := "https://admin:adminpw@localhost:7054"
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"-u", u,
		"--caname", caName,
		"--tls.certfiles", tls,
	}

	if err := im.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when enrolling the ca admin for the organization %s: %v", organization.Name, err)
	}

	return nil
}

func (im *IdentityManager) generateConfigYaml(organization config.Organization) error {
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

		if err := im.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when creating the config.yaml for organization %s: %v", organization.Name, err)
		}
	}

	return nil
}

func (im *IdentityManager) copyCACertificates(organization config.Organization, orgType string) error {
	basePath := getOrgBaseDir(organization.Domain, orgType)
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)

	tlscacerts := fmt.Sprintf("%[1]s/msp/tlscacerts", basePath)
	tlsca := fmt.Sprintf("%[1]s/tlsca", basePath)
	ca := fmt.Sprintf("%[1]s/ca", basePath)
	cacerts := fmt.Sprintf("%[1]s/msp/cacerts", basePath)

	for _, dir := range []string{tlscacerts, tlsca, ca, cacerts} {
		var value bool
		var err error

		if value, err = im.skipIfExists(caName, dir); err != nil {
			return err
		}

		if value == true {
			return nil
		}
	}

	scripts := []string{
		fmt.Sprintf("mkdir -p '%[1]s' && cp '%[2]s' '%[1]s/%[3]s-ca.crt'", tlscacerts, caTlsCertPath, organization.Domain),
		fmt.Sprintf("mkdir -p '%[1]s' && cp '%[2]s' '%[1]s/tlsca.%[3]s-cert.pem'", tlsca, caTlsCertPath, organization.Domain),
		fmt.Sprintf("mkdir -p '%[1]s' && cp '%[2]s' '%[1]s/ca.%[3]s-cert.pem'", ca, caTlsCertPath, organization.Domain),
		fmt.Sprintf("mkdir -p '%[1]s' && cp '%[2]s' '%[1]s/ca.%[3]s-cert.pem'", cacerts, caTlsCertPath, organization.Domain),
	}

	for _, script := range scripts {
		if err := im.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when executing the script %s for organization %s: %v", script, organization.Name, err)
		}
	}
	return nil
}

func (im *IdentityManager) copyPeersCACertificates(organization config.Organization) error {
	return im.copyCACertificates(organization, "peer")
}

func (im *IdentityManager) copyOrderersCACertificates(organization config.Organization) error {
	return im.copyCACertificates(organization, "orderer")
}

func (im *IdentityManager) isRegistered(caName string, id string) (bool, error) {
	args := []string{
		"exec", caName,
		"fabric-ca-client", "identity", "list",
		"--caname", caName,
		"--tls.certfiles", caTlsCertPath,
	}

	output, err := im.executor.OutputCommand("docker", args...)

	if err != nil {
		return false, err
	}

	return strings.Contains(strings.TrimSpace(string(output)), id), nil
}

func (im *IdentityManager) registerPeers(organization config.Organization) error {
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)

	for _, peer := range organization.Peers {
		id := peer.Subdomain

		if alreadyExists, _ := im.isRegistered(caName, id); alreadyExists {
			continue
		}

		args := []string{
			"exec", caName,
			"fabric-ca-client", "register",
			"--caname", caName,
			"--id.name", id,
			"--id.secret", id + "pw",
			"--id.type", "peer",
			"--tls.certfiles", caTlsCertPath,
		}

		if err := im.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering the peer %s for organization %s: %v", peer.Subdomain, organization.Name, err)
		}
	}

	return nil
}

func (im *IdentityManager) registerOrderers(organization config.Organization) error {
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)

	for _, orderer := range organization.Orderers {
		id := strings.ToLower(orderer.Subdomain)

		if alreadyExists, _ := im.isRegistered(caName, id); alreadyExists {
			continue
		}

		args := []string{
			"exec", caName,
			"fabric-ca-client", "register",
			"--caname", caName,
			"--id.name", id,
			"--id.secret", id + "pw",
			"--id.type", "orderer",
			"--tls.certfiles", caTlsCertPath,
		}

		if err := im.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering the orderer %s for organization %s: %v", orderer.Subdomain, organization.Name, err)
		}
	}

	return nil
}

func (im *IdentityManager) registerUser(organization config.Organization) error {
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)
	id := "user1"

	if alreadyExists, _ := im.isRegistered(caName, id); alreadyExists {
		return nil
	}

	args := []string{
		"exec", caName,
		"fabric-ca-client", "register",
		"--caname", caName,
		"--id.name", id,
		"--id.secret", id + "pw",
		"--id.type", "client",
		"--tls.certfiles", caTlsCertPath,
	}

	if err := im.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when registering the user for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (im *IdentityManager) registerOrgAdmin(organization config.Organization) error {
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)
	id := "orgadmin"

	if alreadyExists, _ := im.isRegistered(caName, id); alreadyExists {
		return nil
	}

	args := []string{
		"exec", caName,
		"fabric-ca-client", "register",
		"--caname", caName,
		"--id.name", id,
		"--id.secret", id + "pw",
		"--id.type", "admin",
		"--tls.certfiles", caTlsCertPath,
	}

	if err := im.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when registering the admin for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (im *IdentityManager) generateMSP(organization config.Organization, origin, destination string, id string) error {
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)
	u := fmt.Sprintf("https://%[1]s:%[1]spw@localhost:7054", id)

	for _, dir := range []string{origin, destination} {
		var value bool
		var err error

		value, err = im.skipIfExists(caName, dir)

		if err != nil {
			return err
		}

		if value == true {
			return nil
		}
	}

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"-u", u,
		"--caname", caName,
		"--tls.certfiles", caTlsCertPath,
		"-M", origin,
	}

	if err := im.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when enrolling the %s for organization %s: %v", id, organization.Name, err)
	}

	scripts := []string{
		fmt.Sprintf("mkdir -p '%s/cacerts'", destination),
		fmt.Sprintf("mkdir -p '%s/keystore'", destination),
		fmt.Sprintf("mkdir -p '%s/signcerts'", destination),
		fmt.Sprintf("mkdir -p '%s/tlscacerts'", destination),

		fmt.Sprintf("cp '%s'* '%s'", fmt.Sprintf("%s/cacerts/", origin), fmt.Sprintf("%s/cacerts/ca.%s-cert.pem", destination, organization.Domain)),
		fmt.Sprintf("cp '%s'* '%s'", fmt.Sprintf("%s/keystore/", origin), fmt.Sprintf("%s/keystore/priv_sk", destination)),
		fmt.Sprintf("cp '%s'* '%s'", fmt.Sprintf("%s/signcerts/", origin), fmt.Sprintf("%s/signcerts/cert.pem", destination)),
		fmt.Sprintf("cp '%s' '%s/tlscacerts/%s-ca.crt'", caTlsCertPath, destination, organization.Domain),

		fmt.Sprintf("cp '%s' '%s'", fmt.Sprintf("%s/%s/msp/config.yaml", peerOrgPath, organization.Domain), fmt.Sprintf("%s/config.yaml", destination)),
	}

	for _, script := range scripts {
		if err := im.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when copying the config.yaml to the %s for organization %s: %v", id, organization.Name, err)
		}
	}

	return nil
}

func (im *IdentityManager) generatePeersMSP(organization config.Organization) error {
	for _, peer := range organization.Peers {
		id := peer.Subdomain
		origin := fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s.%[2]s/msp", "/var/hyperledger", organization.Domain, id)
		destination := fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s.%[2]s/msp", peerOrgPath, organization.Domain, id)

		if err := im.generateMSP(organization, origin, destination, id); err != nil {
			return err
		}
	}

	return nil
}

func (im *IdentityManager) generatePeerUserMSP(organization config.Organization) error {
	id := "user1"
	origin := fmt.Sprintf("%[1]s/%[2]s/peers/users/User1@%[2]s/msp", "/var/hyperledger", organization.Domain)
	destination := fmt.Sprintf("%[1]s/%[2]s/users/User1@%[2]s/msp", peerOrgPath, organization.Domain)

	if err := im.generateMSP(organization, origin, destination, id); err != nil {
		return err
	}

	return nil
}

func (im *IdentityManager) generatePeerOrgAdminMSP(organization config.Organization) error {
	id := "orgadmin"
	origin := fmt.Sprintf("%[1]s/%[2]s/peers/users/Admin@%[2]s/msp", "/var/hyperledger", organization.Domain)
	destination := fmt.Sprintf("%[1]s/%[2]s/users/Admin@%[2]s/msp", peerOrgPath, organization.Domain)

	if err := im.generateMSP(organization, origin, destination, id); err != nil {
		return err
	}

	return nil
}

func (im *IdentityManager) generateOrderersMSP(organization config.Organization) error {
	for _, orderer := range organization.Orderers {
		id := strings.ToLower(orderer.Subdomain)
		origin := fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s.%[2]s/msp", "/var/hyperledger", organization.Domain, id)
		destination := fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s.%[2]s/msp", ordererOrgPath, organization.Domain, id)

		if err := im.generateMSP(organization, origin, destination, id); err != nil {
			return err
		}
	}

	return nil
}

func (im *IdentityManager) generateOrdererUserMSP(organization config.Organization) error {
	id := "user1"
	origin := fmt.Sprintf("%[1]s/%[2]s/orderers/users/User1@%[2]s/msp", "/var/hyperledger", organization.Domain)
	destination := fmt.Sprintf("%[1]s/%[2]s/users/User1@%[2]s/msp", ordererOrgPath, organization.Domain)

	if err := im.generateMSP(organization, origin, destination, id); err != nil {
		return err
	}

	return nil
}

func (im *IdentityManager) generateOrdererOrgAdminMSP(organization config.Organization) error {
	id := "orgadmin"
	origin := fmt.Sprintf("%[1]s/%[2]s/orderers/users/Admin@%[2]s/msp", "/var/hyperledger", organization.Domain)
	destination := fmt.Sprintf("%[1]s/%[2]s/users/Admin@%[2]s/msp", ordererOrgPath, organization.Domain)

	if err := im.generateMSP(organization, origin, destination, id); err != nil {
		return err
	}

	return nil
}

func (im *IdentityManager) generateTLS(organization config.Organization, origin string, destination string, id string) error {
	caName := compose.ResolveCertificateAuthorityContainerName(organization.Domain)
	u := fmt.Sprintf("https://%[1]s:%[1]spw@localhost:7054", id)

	for _, dir := range []string{origin, destination} {
		var value bool
		var err error

		if value, err = im.skipIfExists(caName, dir); err != nil {
			return err
		}

		if value == true {
			return nil
		}
	}

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"--caname", caName,
		"-u", u,
		"-M", origin,
		"--enrollment.profile", "tls",
		"--csr.hosts", fmt.Sprintf("%s.%s", id, organization.Domain),
		"--csr.hosts", "localhost",
		"--tls.certfiles", caTlsCertPath,
	}

	if err := im.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when genereting the tls certificates of the %s for organization %s: %v", id, organization.Name, err)
	}

	scripts := []string{
		fmt.Sprintf("mkdir -p '%s'", destination),
		fmt.Sprintf("cp '%s/tlscacerts/'* '%s/ca.crt'", origin, destination),
		fmt.Sprintf("cp '%s/signcerts/'* '%s/server.crt'", origin, destination),
		fmt.Sprintf("cp '%s/keystore/'* '%s/server.key'", origin, destination),
	}

	for _, script := range scripts {
		if err := im.execInCA(organization.Domain, script); err != nil {
			return fmt.Errorf("Error when copying the tls certificate of the %s for organization %s: %v", id, organization.Name, err)
		}
	}

	return nil
}

func (im *IdentityManager) generatePeerTlsCertificates(organization config.Organization) error {
	for _, peer := range organization.Peers {
		id := peer.Subdomain
		origin := fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s.%[2]s/tls", "/var/hyperledger", organization.Domain, id)
		destination := fmt.Sprintf("%[1]s/%[2]s/peers/%[3]s.%[2]s/tls", peerOrgPath, organization.Domain, id)

		if err := im.generateTLS(organization, origin, destination, id); err != nil {
			return err
		}
	}

	return nil
}

func (im *IdentityManager) generateOrdererTlsCertificates(organization config.Organization) error {
	for _, orderer := range organization.Orderers {
		id := strings.ToLower(orderer.Subdomain)

		origin := fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s.%[2]s/tls", "/var/hyperledger", organization.Domain, id)
		destination := fmt.Sprintf("%[1]s/%[2]s/orderers/%[3]s.%[2]s/tls", ordererOrgPath, organization.Domain, id)

		if err := im.generateTLS(organization, origin, destination, id); err != nil {
			return err
		}
	}

	return nil
}

func (im *IdentityManager) shareTlsCertificates() error {

	for _, sourceOrganization := range im.config.Organizations {
		folder := "%[1]s/%[2]s/certificate-authority/organizations/peerOrganizations/%[2]s"

		for _, targetOrganization := range im.config.Organizations {
			if targetOrganization.Domain == sourceOrganization.Domain {
				continue
			}

			origin := fmt.Sprintf("%[1]s/msp/tlscacerts/%[2]s-ca.crt", fmt.Sprintf(folder, im.config.Output, sourceOrganization.Domain), sourceOrganization.Domain)
			destination := fmt.Sprintf("%[1]s/msp/tlscacerts/%[2]s-ca.crt", fmt.Sprintf(folder, im.config.Output, targetOrganization.Domain), sourceOrganization.Domain)

			if err := file.Copy(origin, destination); err != nil {
				return err
			}

			for _, peer := range targetOrganization.Peers {
				destination := fmt.Sprintf("%[1]s/peers/%[2]s.%[3]s/msp/tlscacerts/%[4]s-ca.crt", fmt.Sprintf(folder, im.config.Output, targetOrganization.Domain), peer.Subdomain, targetOrganization.Domain, sourceOrganization.Domain)

				if err := file.Copy(origin, destination); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (im *IdentityManager) skipIfExists(caName string, item string) (bool, error) {
	args := []string{
		"exec", caName, "sh", "-c", fmt.Sprintf("test -d %s && echo 1 || echo 0", item),
	}
	output, err := im.executor.OutputCommand("docker", args...)

	if err != nil {
		return false, err
	}

	if strings.TrimSpace(string(output)) == "1" {
		fmt.Printf("Skipping: Directory already exists\n")
		return true, nil
	}

	return false, nil
}
