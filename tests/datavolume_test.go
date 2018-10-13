package tests_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"kubevirt.io/kubevirt-ansible/tests"
	kubev1 "kubevirt.io/kubevirt/pkg/api/v1"
	ktests "kubevirt.io/kubevirt/tests"
)

// template parameters
const rawDataVolumeVMFilePath = "tests/manifests/template/datavolume-vm.yml"
const rawDataVolumeVMIFilePath = "tests/manifests/template/datavolume-vmi.yml"
const rawDataVolumeFilePath = "tests/manifests/template/datavolume.yml"

var _ = Describe("DataVolume Integration Test", func() {
	var dataVolumeName, vmName, dstDataVolumeFilePath, url, dstVMIFilePath string

	Context("Datavolume with VM", func() {
		BeforeEach(func() {
			dataVolumeName = "datavolume1"
			vmName = "test-vm-i"
			dstDataVolumeFilePath = "/tmp/test-datavolume-vm.json"
			url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
		})

		It("Creating VM and start VMI will be success", func() {
			tests.ProcessTemplateWithParameters(rawDataVolumeVMFilePath, dstDataVolumeFilePath, "VM_APIVERSION="+kubev1.GroupVersion.String(), "VM_NAME="+vmName, "IMG_URL="+tests.ReplaceImageURL(url), "DATAVOLUME_NAME="+dataVolumeName)
			tests.CreateResourceWithFilePathTestNamespace(dstDataVolumeFilePath)
			tests.WaitUntilResourceReadyByNameTestNamespace("pvc", dataVolumeName, "-o=jsonpath='{.metadata.annotations}'", "pv.kubernetes.io/bind-completed:yes")
			By("Start VM with virtctl")
			args := []string{"start", vmName, "-n", tests.NamespaceTestDefault}
			_, _, err := ktests.RunCommand("virtctl", args...)
			Expect(err).ToNot(HaveOccurred())
			tests.WaitUntilResourceReadyByNameTestNamespace("vmi", vmName, "-o=jsonpath='{.status.phase}'", "Running")
		})
	})

	Context("Datavolume with VMI", func() {
		BeforeEach(func() {
			dataVolumeName = "datavolume2"
			vmName = "test-vmi"
			dstDataVolumeFilePath = "/tmp/test-datavolume.json"
			dstVMIFilePath = "/tmp/test-datavolume-vmi.json"
			url = "https://download.cirros-cloud.net/0.4.0/cirros-0.4.0-x86_64-disk.img"
		})

		It("Pre creating datavolume then create VMI will be success", func() {
			tests.ProcessTemplateWithParameters(rawDataVolumeFilePath, dstDataVolumeFilePath, "DATAVOLUME_NAME="+dataVolumeName, "IMG_URL="+tests.ReplaceImageURL(url))
			tests.CreateResourceWithFilePathTestNamespace(dstDataVolumeFilePath)
			tests.WaitUntilResourceReadyByNameTestNamespace("pvc", dataVolumeName, "-o=jsonpath='{.metadata.annotations}'", "pv.kubernetes.io/bind-completed:yes")

			tests.ProcessTemplateWithParameters(rawDataVolumeVMIFilePath, dstVMIFilePath, "VM_APIVERSION="+kubev1.GroupVersion.String(), "VM_NAME="+vmName, "DATAVOLUME_NAME="+dataVolumeName)
			tests.CreateResourceWithFilePathTestNamespace(dstVMIFilePath)
			tests.WaitUntilResourceReadyByNameTestNamespace("vmi", vmName, "-o=jsonpath='{.status.phase}'", "Running")
		})
	})

})
