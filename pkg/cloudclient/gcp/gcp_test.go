package gcp

import (
	"testing"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/cloud-ingress-operator/config"
	"github.com/openshift/cloud-ingress-operator/pkg/testutils"
	computev1 "google.golang.org/api/compute/v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

func TestNewClient(t *testing.T) {
	dummySA := `{
		"type": "service_account",
		"private_key_id": "abc",
		"private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvgIBADANBgkqhkiG9w0BAQEFAKECBKgwggSkAgEAAoIBAQDY3E8o1NEFcjFAKEHW/5ZfFJw29/8NEqpViNjQIx95Xx5KDtJ+nWFAKEW0uqsSqKlKGhAdAo+Q6bjx2cFAKEVsXTu7XrZUY5Kltvj94DvUa1wjNXs606r/RxWTJ58bfdC+gLLxBfGnB6CwK0YQ\nxnfpjNbkUfVVzO0MQD7UP0Hl5ZcY0Puvxd/yHuONQn/rIAieTHH1pqgW+zrH/y3c\n59IGThC9PPtugI9ea8RSnVj3PWz1bX2UkCDpy9IRh9LzJLaYYX9RUd7++dULUlat\nAaXBh1U6emUDzhrIsgApjDVtimOPbmQWmX1S60mqQikRpVYZ8u+NDD+LNw+/Eovn\nxCj2Y3z1AgMBAAECggEAWDBzoqO1IvVXjBA2lqId10T6hXmN3j1ifyH+aAqK+FVl\nGjyWjDj0xWQcJ9ync7bQ6fSeTeNGzP0M6kzDU1+w6FgyZqwdmXWI2VmEizRjwk+/\n/uLQUcL7I55Dxn7KUoZs/rZPmQDxmGLoue60Gg6z3yLzVcKiDc7cnhzhdBgDc8vd\nQorNAlqGPRnm3EqKQ6VQp6fyQmCAxrr45kspRXNLddat3AMsuqImDkqGKBmF3Q1y\nxWGe81LphUiRqvqbyUlh6cdSZ8pLBpc9m0c3qWPKs9paqBIvgUPlvOZMqec6x4S6\nChbdkkTRLnbsRr0Yg/nDeEPlkhRBhasXpxpMUBgPywKBgQDs2axNkFjbU94uXvd5\nznUhDVxPFBuxyUHtsJNqW4p/ujLNimGet5E/YthCnQeC2P3Ym7c3fiz68amM6hiA\nOnW7HYPZ+jKFnefpAtjyOOs46AkftEg07T9XjwWNPt8+8l0DYawPoJgbM5iE0L2O\nx8TU1Vs4mXc+ql9F90GzI0x3VwKBgQDqZOOqWw3hTnNT07Ixqnmd3dugV9S7eW6o\nU9OoUgJB4rYTpG+yFqNqbRT8bkx37iKBMEReppqonOqGm4wtuRR6LSLlgcIU9Iwx\nyfH12UWqVmFSHsgZFqM/cK3wGev38h1WBIOx3/djKn7BdlKVh8kWyx6uC8bmV+E6\nOoK0vJD6kwKBgHAySOnROBZlqzkiKW8c+uU2VATtzJSydrWm0J4wUPJifNBa/hVW\ndcqmAzXC9xznt5AVa3wxHBOfyKaE+ig8CSsjNyNZ3vbmr0X04FoV1m91k2TeXNod\njMTobkPThaNm4eLJMN2SQJuaHGTGERWC0l3T18t+/zrDMDCPiSLX1NAvAoGBAN1T\nVLJYdjvIMxf1bm59VYcepbK7HLHFkRq6xMJMZbtG0ryraZjUzYvB4q4VjHk2UDiC\nlhx13tXWDZH7MJtABzjyg+AI7XWSEQs2cBXACos0M4Myc6lU+eL+iA+OuoUOhmrh\nqmT8YYGu76/IBWUSqWuvcpHPpwl7871i4Ga/I3qnAoGBANNkKAcMoeAbJQK7a/Rn\nwPEJB+dPgNDIaboAsh1nZhVhN5cvdvCWuEYgOGCPQLYQF0zmTLcM+sVxOYgfy8mV\nfbNgPgsP5xmu6dw2COBKdtozw0HrWSRjACd1N4yGu75+wPCcX/gQarcjRcXXZeEa\nNtBLSfcqPULqD+h7br9lEJio\n-----END PRIVATE KEY-----\n",
		"client_email": "123-abc@developer.gserviceaccount.com",
		"client_id": "123-abc.apps.googleusercontent.com",
		"auth_uri": "https://accounts.google.com/o/oauth2/auth",
		"token_uri": "http://localhost:8080/token"
	  }`

	infra := &configv1.Infrastructure{
		ObjectMeta: v1.ObjectMeta{
			Name:      "cluster",
			Namespace: "",
		},
		Status: configv1.InfrastructureStatus{
			APIServerURL: "https://api.dummy.fsrl.s1.testing.org:6443",
			PlatformStatus: &configv1.PlatformStatus{
				GCP: &configv1.GCPPlatformStatus{
					Region: "eu-west-1",
				},
			},
		},
	}

	fakeSecret := &corev1.Secret{
		ObjectMeta: v1.ObjectMeta{
			Name:      config.GCPSecretName,
			Namespace: config.OperatorNamespace,
		},
		Data: make(map[string][]byte),
	}
	fakeSecret.Data["service_account.json"] = []byte(dummySA)

	objs := []runtime.Object{infra, fakeSecret}
	mocks := testutils.NewTestMock(t, objs)
	cli, err := NewClient(mocks.FakeKubeClient)

	if err != nil {
		t.Error("err occured while creating cli:", err)
	}

	if cli == nil {
		t.Error("cli should have been initialized")
	}
}

// testing only performHealthCheck
// mocking c.computeService.RegionBackendServices.List(c.projectID, c.region).Do() is hard coz they're not bound to an interface
func TestPerformHealthCheck(t *testing.T) {
	l := &computev1.BackendServiceList{
		Items: []*computev1.BackendService{
			{
				Name: "test-api-internal",
			},
		}}

	err := performHealthCheck(l, "test")
	if err != nil {
		t.Error("tests shouldn't have failed:", err)
	}
}
