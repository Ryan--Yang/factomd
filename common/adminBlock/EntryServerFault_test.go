package adminBlock_test

import (
	"testing"

	. "github.com/FactomProject/factomd/common/adminBlock"
	//"github.com/FactomProject/factomd/common/interfaces"
	"github.com/FactomProject/factomd/common/primitives"
	//"github.com/FactomProject/factomd/state"
	"github.com/FactomProject/factomd/testHelper"
)

func TestServerFaultMarshalUnmarshal(t *testing.T) {
	sf := new(ServerFault)

	sf.Timestamp = primitives.NewTimestampNow()
	sf.ServerID = testHelper.NewRepeatingHash(1)
	sf.AuditServerID = testHelper.NewRepeatingHash(2)

	sf.VMIndex = 0x33
	sf.DBHeight = 0x44556677
	sf.Height = 0x88990011

	core, err := sf.MarshalCore()
	if err != nil {
		t.Errorf("%v", err)
	}
	for i := 0; i < 10; i++ {
		priv := testHelper.NewPrimitivesPrivateKey(uint64(i))
		sig := priv.Sign(core)
		sf.SignatureList.List = append(sf.SignatureList.List, sig)
	}
	sf.SignatureList.Length = uint32(len(sf.SignatureList.List))

	bin, err := sf.MarshalBinary()
	if err != nil {
		t.Errorf("%v", err)
	}

	sf2 := new(ServerFault)
	rest, err := sf2.UnmarshalBinaryData(bin)
	if err != nil {
		t.Errorf("%v", err)
	}
	if len(rest) > 0 {
		t.Errorf("Unexpected extra piece of data - %x", rest)
	}
	t.Logf("%v", sf.String())
	t.Logf("%v", sf2.String())

	if sf.Timestamp.GetTimeMilliUInt64() !=
		sf2.Timestamp.GetTimeMilliUInt64() {
		t.Errorf("Invalid Timestamp")
	}
	if sf.ServerID.IsSameAs(sf2.ServerID) == false {
		t.Errorf("Invalid ServerID")
	}
	if sf.AuditServerID.IsSameAs(sf2.AuditServerID) == false {
		t.Errorf("Invalid AuditServerID")
	}
	if sf.VMIndex != sf2.VMIndex {
		t.Errorf("Invalid VMIndex")
	}
	if sf.DBHeight != sf2.DBHeight {
		t.Errorf("Invalid DBHeight")
	}
	if sf.Height != sf2.Height {
		t.Errorf("Invalid Height")
	}

	if sf.SignatureList.Length != sf2.SignatureList.Length {
		t.Errorf("Invalid SignatureList.Length")
	}
	if len(sf.SignatureList.List) != len(sf2.SignatureList.List) {
		t.Errorf("Invalid len of SignatureList.List")
	} else {
		for i := range sf.SignatureList.List {
			if sf.SignatureList.List[i].IsSameAs(sf2.SignatureList.List[i]) == false {
				t.Errorf("Invalid SignatureList.List at %v", i)
			}
		}
	}
}

func TestServerFaultUpdateState(t *testing.T) {
	sigs := 10
	sf := new(ServerFault)

	sf.Timestamp = primitives.NewTimestampNow()
	sf.ServerID = testHelper.NewRepeatingHash(1)
	sf.AuditServerID = testHelper.NewRepeatingHash(2)

	sf.VMIndex = 0x33
	sf.DBHeight = 0x44556677
	sf.Height = 0x88990011

	core, err := sf.MarshalCore()
	if err != nil {
		t.Errorf("%v", err)
	}
	for i := 0; i < sigs; i++ {
		priv := testHelper.NewPrimitivesPrivateKey(uint64(i))
		sig := priv.Sign(core)
		sf.SignatureList.List = append(sf.SignatureList.List, sig)
	}
	sf.SignatureList.Length = uint32(len(sf.SignatureList.List))

	s := testHelper.CreateAndPopulateTestState()
	idindex := s.CreateBlankFactomIdentity(primitives.NewZeroHash())
	s.Identities[idindex].ManagementChainID = primitives.NewZeroHash()
	for i := 0; i < sigs; i++ {
		//Federated Server
		index := s.AddAuthorityFromChainID(testHelper.NewRepeatingHash(byte(i)))
		s.Authorities[index].SigningKey = *testHelper.NewPrimitivesPrivateKey(uint64(i)).Pub
		s.Authorities[index].Status = 1

		s.AddFedServer(s.GetLeaderHeight(), testHelper.NewRepeatingHash(byte(i)))

		//Audit Server
		index = s.AddAuthorityFromChainID(testHelper.NewRepeatingHash(byte(i + sigs)))
		s.Authorities[index].SigningKey = *testHelper.NewPrimitivesPrivateKey(uint64(i + sigs)).Pub
		s.Authorities[index].Status = 0

		s.AddAuditServer(s.GetLeaderHeight(), testHelper.NewRepeatingHash(byte(i+sigs)))
	}

	err = sf.UpdateState(s)
	if err != nil {
		t.Errorf("%v", err)
	}

}

/*
func TestAuthoritySignature(t *testing.T) {
	s := testHelper.CreateAndPopulateTestState()
	idindex := CreateBlankFactomIdentity(s, primitives.NewZeroHash())
	s.Identities[idindex].ManagementChainID = primitives.NewZeroHash()

	index := s.AddAuthorityFromChainID(primitives.NewZeroHash())
	s.Authorities[index].SigningKey = *(s.GetServerPublicKey())
	s.Authorities[index].Status = 1

	ack := new(messages.Ack)
	ack.DBHeight = s.LLeaderHeight
	ack.VMIndex = 1
	ack.Minute = byte(5)
	ack.Timestamp = s.GetTimestamp()
	ack.MessageHash = primitives.NewZeroHash()
	ack.LeaderChainID = s.IdentityChainID
	ack.SerialHash = primitives.NewZeroHash()

	err := ack.Sign(s)
	if err != nil {
		t.Error("Authority Test Failed when signing message")
	}

	msg, err := ack.MarshalForSignature()
	if err != nil {
		t.Error("Authority Test Failed when marshalling for sig")
	}

	sig := ack.GetSignature()
	server, err := s.Authorities[0].VerifySignature(msg, sig.GetSignature())
	if !server || err != nil {
		t.Error("Authority Test Failed when checking sigs")
	}
}

*/