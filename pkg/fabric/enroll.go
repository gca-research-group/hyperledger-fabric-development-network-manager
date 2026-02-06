package fabric

import (
	"fmt"
	"strings"

	"github.com/gca-research-group/hyperledger-fabric-development-network-manager/pkg"
)

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
		return fmt.Errorf("Error when generating the certificates for the organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) generateConfigYaml(organization pkg.Organization) error {
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

	folders := []string{
		fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/msp", organization.Domain),
		fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/msp", organization.Domain),
	}

	for _, folder := range folders {
		args := []string{
			"exec", fmt.Sprintf("ca.%s", organization.Domain),
			"sh", "-c", fmt.Sprintf("cat <<EOF > %s/config.yaml %s", folder, config),
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when creating config.yaml for organization %s: %v", organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) copyPeersCACertificates(organization pkg.Organization) error {
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	tlscacerts := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/msp/tlscacerts", organization.Domain)
	tlsca := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/tlsca", organization.Domain)
	ca := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/ca", organization.Domain)

	scripts := []string{
		fmt.Sprintf("mkdir -p '%s'", tlscacerts),
		fmt.Sprintf("cp '%s' '%s/ca.crt'", tls, tlscacerts),

		fmt.Sprintf("mkdir -p '%s'", tlsca),
		fmt.Sprintf("cp '%s' '%s/tlsca.%s-cert.pem'", tls, tlsca, organization.Domain),

		fmt.Sprintf("mkdir -p '%s'", ca),
		fmt.Sprintf("cp '%s' '%s/ca.%s-cert.pem'", tls, ca, organization.Domain),
	}

	for _, script := range scripts {
		args := []string{
			"exec", fmt.Sprintf("ca.%s", organization.Domain),
			"sh", "-c", script,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when executing the script %s for organization %s: %v", script, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) copyOrderersCACertificates(organization pkg.Organization) error {
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	tlscacerts := fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/msp/tlscacerts", organization.Domain)
	tlsca := fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/tlsca", organization.Domain)
	ca := fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/ca", organization.Domain)

	scripts := []string{
		fmt.Sprintf("mkdir -p '%s'", tlscacerts),
		fmt.Sprintf("cp '%s' '%s/ca.crt'", tls, tlscacerts),

		fmt.Sprintf("mkdir -p '%s'", tlsca),
		fmt.Sprintf("cp '%s' '%s/tlsca.%s-cert.pem'", tls, tlsca, organization.Domain),

		fmt.Sprintf("mkdir -p '%s'", ca),
		fmt.Sprintf("cp '%s' '%s/ca.%s-cert.pem'", tls, ca, organization.Domain),
	}

	for _, script := range scripts {
		args := []string{
			"exec", fmt.Sprintf("ca.%s", organization.Domain),
			"sh", "-c", script,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when executing the script %s for organization %s: %v", script, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) registerPeers(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	for i := range organization.Peers {
		args := []string{
			"exec", caName,
			"fabric-ca-client", "register",
			"--caname", caName,
			"--id.name", fmt.Sprintf("peer%d", i),
			"--id.secret", fmt.Sprintf("peer%dpw", i),
			"--id.type", "peer",
			"--tls.certfiles", tls,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering the peer %d for organization %s: %v", i, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) registerOrderes(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	for _, orderer := range organization.Orderers {
		args := []string{
			"exec", caName,
			"fabric-ca-client", "register",
			"--caname", caName,
			"--id.name", strings.ToLower(orderer.Hostname),
			"--id.secret", fmt.Sprintf("%spw", strings.ToLower(orderer.Hostname)),
			"--id.type", "orderer",
			"--tls.certfiles", tls,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when registering the orderer %s for organization %s: %v", orderer.Hostname, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) registerUser(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	args := []string{
		"exec", caName,
		"fabric-ca-client", "register",
		"--caname", caName,
		"--id.name", "user1",
		"--id.secret", "user1pw",
		"--id.type", "client",
		"--tls.certfiles", tls,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when registering the user for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) registerOrgAdmin(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	args := []string{
		"exec", caName,
		"fabric-ca-client", "register",
		"--caname", caName,
		"--id.name", "org1admin",
		"--id.secret", "org1adminpw",
		"--id.type", "admin",
		"--tls.certfiles", tls,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when registering the admin for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) generatePeersMSP(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	for i := range organization.Peers {
		u := fmt.Sprintf("https://peer%d:peer%dpw@localhost:7054", i, i)
		m := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/peers/peer%d.%s/msp", organization.Domain, i, organization.Domain)

		args := []string{
			"exec", caName,
			"fabric-ca-client", "enroll",
			"-u", u,
			"--caname", caName,
			"--tls.certfiles", tls,
			"-M", m,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when enrolling the peer %d for organization %s: %v", i, organization.Name, err)
		}

		from := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/msp/config.yaml", organization.Domain)
		to := fmt.Sprintf("%s/config.yaml", m)

		script := fmt.Sprintf("cp '%s' '%s'", from, to)

		args = []string{
			"exec", fmt.Sprintf("ca.%s", organization.Domain),
			"sh", "-c", script,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when copying the config.yaml to the peer %d for organization %s: %v", i, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) generateOrderersMSP(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	for _, orderer := range organization.Orderers {
		u := fmt.Sprintf("https://%s:%spw@localhost:7054", orderer.Hostname, orderer.Hostname)
		m := fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/orderers/%s.%s/msp", organization.Domain, orderer.Hostname, organization.Domain)

		args := []string{
			"exec", caName,
			"fabric-ca-client", "enroll",
			"-u", u,
			"--caname", caName,
			"--tls.certfiles", tls,
			"-M", m,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when enrolling the orderer %s for organization %s: %v", orderer.Hostname, organization.Name, err)
		}

		from := fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/msp/config.yaml", organization.Domain)
		to := fmt.Sprintf("%s/config.yaml", m)

		script := fmt.Sprintf("cp '%s' '%s'", from, to)

		args = []string{
			"exec", fmt.Sprintf("ca.%s", organization.Domain),
			"sh", "-c", script,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when copying the config.yaml to the orderer %s for organization %s: %v", orderer.Hostname, organization.Name, err)
		}
	}

	return nil
}

func (f *Fabric) generatePeerTlsCertificates(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	for i := range organization.Peers {
		u := fmt.Sprintf("https://peer%d:peer%dpw@localhost:7054", i, i)
		m := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/peers/peer%d.%s/tls", organization.Domain, i, organization.Domain)

		args := []string{
			"exec", caName,
			"fabric-ca-client", "enroll",
			"--caname", caName,
			"-u", u,
			"-M", m,
			"--enrollment.profile", "tls",
			"--csr.hosts", fmt.Sprintf("peer%d.%s", i, organization.Domain), "--csr.hosts", "localhost",
			"--tls.certfiles", tls,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when genereting the tls certificates of the peer %d for organization %s: %v", i, organization.Name, err)
		}

		scripts := []string{
			fmt.Sprintf("cp '%s/tlscacerts/'* '%s/ca.crt'", m, m),
			fmt.Sprintf("cp '%s/signcerts/'* '%s/server.crt'", m, m),
			fmt.Sprintf("cp '%s/keystore/'* '%s/server.key'", m, m),
		}

		for _, script := range scripts {
			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain),
				"sh", "-c", script,
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when copying the tls certificate of the peer %d for organization %s: %v", i, organization.Name, err)
			}
		}
	}

	return nil
}

func (f *Fabric) generateOrdererTlsCertificates(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	for _, orderer := range organization.Orderers {
		u := fmt.Sprintf("https://%s:%spw@localhost:7054", orderer.Hostname, orderer.Hostname)
		m := fmt.Sprintf("/etc/hyperledger/fabric/organizations/ordererOrganizations/%s/orderers/%s.%s/tls", organization.Domain, orderer.Hostname, organization.Domain)

		args := []string{
			"exec", caName,
			"fabric-ca-client", "enroll",
			"--caname", caName,
			"-u", u,
			"-M", m,
			"--enrollment.profile", "tls",
			"--csr.hosts", fmt.Sprintf("%s.%s", orderer.Hostname, organization.Domain), "--csr.hosts", "localhost",
			"--tls.certfiles", tls,
		}

		if err := f.executor.ExecCommand("docker", args...); err != nil {
			return fmt.Errorf("Error when genereting the tls certificates of the orderer %s for organization %s: %v", orderer.Hostname, organization.Name, err)
		}

		scripts := []string{
			fmt.Sprintf("cp '%s/tlscacerts/'* '%s/ca.crt'", m, m),
			fmt.Sprintf("cp '%s/signcerts/'* '%s/server.crt'", m, m),
			fmt.Sprintf("cp '%s/keystore/'* '%s/server.key'", m, m),
		}

		for _, script := range scripts {
			args = []string{
				"exec", fmt.Sprintf("ca.%s", organization.Domain),
				"sh", "-c", script,
			}

			if err := f.executor.ExecCommand("docker", args...); err != nil {
				return fmt.Errorf("Error when copying the tls certificate of the orderer %s for organization %s: %v", orderer.Hostname, organization.Name, err)
			}
		}
	}

	return nil
}

func (f *Fabric) generateUserMSP(organization pkg.Organization) error {
	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	u := "https://user1:user1pw@localhost:7054"
	m := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/users/User1@%s/msp", organization.Domain, organization.Domain)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"-u", u,
		"--caname", caName,
		"-M", m,
		"--tls.certfiles", tls,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when enrolling the user for organization %s: %v", organization.Name, err)
	}

	from := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/msp/config.yaml", organization.Domain)
	to := fmt.Sprintf("%s/config.yaml", m)

	script := fmt.Sprintf("cp '%s' '%s'", from, to)

	args = []string{
		"exec", fmt.Sprintf("ca.%s", organization.Domain),
		"sh", "-c", script,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when copying the config.yaml to the user for organization %s: %v", organization.Name, err)
	}

	return nil
}

func (f *Fabric) generateOrgAdminMSP(organization pkg.Organization) error {

	caName := fmt.Sprintf("ca.%s", organization.Domain)
	tls := "/etc/hyperledger/fabric-ca-server/ca-cert.pem"

	u := "https://org1admin:org1adminpw@localhost:7054"
	m := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/users/Admin@%s/msp", organization.Domain, organization.Domain)

	args := []string{
		"exec", caName,
		"fabric-ca-client", "enroll",
		"-u", u,
		"--caname", caName,
		"-M", m,
		"--tls.certfiles", tls,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when enrolling the user for organization %s: %v", organization.Name, err)
	}

	from := fmt.Sprintf("/etc/hyperledger/fabric/organizations/peerOrganizations/%s/msp/config.yaml", organization.Domain)
	to := fmt.Sprintf("%s/config.yaml", m)

	script := fmt.Sprintf("cp '%s' '%s'", from, to)

	args = []string{
		"exec", fmt.Sprintf("ca.%s", organization.Domain),
		"sh", "-c", script,
	}

	if err := f.executor.ExecCommand("docker", args...); err != nil {
		return fmt.Errorf("Error when copying the config.yaml to the user for organization %s: %v", organization.Name, err)
	}

	return nil
}
