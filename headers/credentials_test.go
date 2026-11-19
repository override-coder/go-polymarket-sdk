package headers

import (
	"testing"

	clobtypes "github.com/override-coder/go-polymarket-sdk/clob/types"
	sdktypes "github.com/override-coder/go-polymarket-sdk/types"
)

func TestMissingCredentials(t *testing.T) {
	args := clobtypes.L2HeaderArgs{Method: "GET", RequestPath: "/orders"}
	if h, err := CreateL2Headers("address", nil, args, nil); err == nil || h != nil {
		t.Fatalf("missing L2 credentials: headers=%v err=%v", h, err)
	}
	if h, err := CreateL2BuilderHeaders(nil, args, nil); err == nil || h != nil {
		t.Fatalf("missing builder credentials: headers=%v err=%v", h, err)
	}
}

func TestCredentialHeaders(t *testing.T) {
	ts := int64(1700000000)
	args := clobtypes.L2HeaderArgs{Method: "POST", RequestPath: "/order", Body: "{}"}
	l2, err := CreateL2Headers("address", &sdktypes.ApiKeyCreds{
		ApiKey: "key", Secret: "c2VjcmV0", Passphrase: "pass",
	}, args, &ts)
	if err != nil {
		t.Fatal(err)
	}
	builder, err := CreateL2BuilderHeaders(&sdktypes.BuilderApiKeyCreds{
		Key: "key", Secret: "c2VjcmV0", Passphrase: "pass",
	}, args, &ts)
	if err != nil {
		t.Fatal(err)
	}
	if l2["POLY_SIGNATURE"] == "" || l2["POLY_SIGNATURE"] != builder["POLY_BUILDER_SIGNATURE"] ||
		l2["POLY_TIMESTAMP"] != "1700000000" || builder["POLY_BUILDER_TIMESTAMP"] != "1700000000" ||
		l2["POLY_API_KEY"] != "key" || builder["POLY_BUILDER_API_KEY"] != "key" {
		t.Fatal("credential header contract changed")
	}
}
