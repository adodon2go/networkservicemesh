// +build basic

package integration

import (
	"testing"

	. "github.com/onsi/gomega"
	"github.com/sirupsen/logrus"

	"cisco-app-networking.github.io/networkservicemesh/test/kubetest"
	"cisco-app-networking.github.io/networkservicemesh/test/kubetest/pods"
)

func TestExec(t *testing.T) {
	g := NewWithT(t)

	if testing.Short() {
		t.Skip("Skip, please run without -short")
		return
	}

	k8s, err := kubetest.NewK8sWithoutRoles(g, true)
	defer k8s.Cleanup()
	g.Expect(err).To(BeNil())
	defer k8s.SaveTestArtifacts(t)

	k8s.DeletePodsByName("alpine-pod")

	alpinePod := k8s.CreatePod(pods.AlpinePod("alpine-pod", nil))

	ipResponse, errResponse, error := k8s.Exec(alpinePod, alpinePod.Spec.Containers[0].Name, "ip", "addr")
	g.Expect(error).To(BeNil())
	g.Expect(errResponse).To(Equal(""))
	logrus.Printf("NSC IP status:%s", ipResponse)
	logrus.Printf("End of test")
	k8s.DeletePods(alpinePod)
}
