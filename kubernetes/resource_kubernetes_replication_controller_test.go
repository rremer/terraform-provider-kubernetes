package kubernetes

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	api "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestAccKubernetesReplicationController_basic(t *testing.T) {
	var conf api.ReplicationController
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_replication_controller.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfig_basic(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.TestAnnotationTwo", "two"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one", "TestAnnotationTwo": "two"}),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.%", "3"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.TestLabelTwo", "two"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.TestLabelThree", "three"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelTwo": "two", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.image", "nginx:1.7.8"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.name", "tf-acc-test"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.metadata.0.annotations.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.metadata.0.annotations.TestAnnotationFive", "five"),
				),
			},
			{
				Config: testAccKubernetesReplicationControllerConfig_modified(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.TestAnnotationOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.Different", "1234"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{"TestAnnotationOne": "one", "Different": "1234"}),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.%", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.TestLabelOne", "one"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.TestLabelThree", "three"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.name", name),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.uid"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.image", "nginx:1.7.9"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.name", "tf-acc-test"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.metadata.0.annotations.%", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.metadata.0.annotations.TestAnnotationSix", "six"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_initContainer(t *testing.T) {
	var conf api.ReplicationController
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_replication_controller.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfig_initContainer(name),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.image", "busybox"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.name", "install"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.command.0", "wget"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.command.1", "-O"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.command.2", "/work-dir/index.html"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.command.3", "http://kubernetes.io"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.volume_mount.0.name", "workdir"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.init_container.0.volume_mount.0.mount_path", "/work-dir"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.nameservers.#", "3"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.nameservers.0", "1.1.1.1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.nameservers.1", "8.8.8.8"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.nameservers.2", "9.9.9.9"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.searches.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.searches.0", "kubernetes.io"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.option.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.option.0.name", "ndots"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.option.0.value", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.option.1.name", "use-vc"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_config.0.option.1.value", ""),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.dns_policy", "Default"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_importBasic(t *testing.T) {
	resourceName := "kubernetes_replication_controller.test"
	name := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfig_basic(name),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
		},
	})
}

func TestAccKubernetesReplicationController_generatedName(t *testing.T) {
	var conf api.ReplicationController
	prefix := "tf-acc-test-gen-"

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "kubernetes_replication_controller.test",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfig_generatedName(prefix),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.annotations.%", "0"),
					testAccCheckMetaAnnotations(&conf.ObjectMeta, map[string]string{}),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.labels.%", "3"),
					testAccCheckMetaLabels(&conf.ObjectMeta, map[string]string{"TestLabelOne": "one", "TestLabelTwo": "two", "TestLabelThree": "three"}),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "metadata.0.generate_name", prefix),
					resource.TestMatchResourceAttr("kubernetes_replication_controller.test", "metadata.0.name", regexp.MustCompile("^"+prefix)),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.generation"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.resource_version"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.self_link"),
					resource.TestCheckResourceAttrSet("kubernetes_replication_controller.test", "metadata.0.uid"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_importGeneratedName(t *testing.T) {
	resourceName := "kubernetes_replication_controller.test"
	prefix := "tf-acc-test-gen-import-"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfig_generatedName(prefix),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"metadata.0.resource_version"},
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_security_context(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithSecurityContext(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.fs_group", "100"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.run_as_non_root", "true"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.run_as_user", "101"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.supplemental_groups.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.supplemental_groups.988695518", "101"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_security_context_run_as_group(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t); skipIfUnsupportedSecurityContextRunAsGroup(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithSecurityContextRunAsGroup(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.fs_group", "100"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.run_as_group", "100"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.run_as_non_root", "true"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.run_as_user", "101"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.supplemental_groups.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.security_context.0.supplemental_groups.988695518", "101"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_container_liveness_probe_using_exec(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/busybox"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithLivenessProbeUsingExec(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.args.#", "3"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.exec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.exec.0.command.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.exec.0.command.0", "cat"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.exec.0.command.1", "/tmp/healthy"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.failure_threshold", "3"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.initial_delay_seconds", "5"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_container_liveness_probe_using_http_get(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/liveness"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithLivenessProbeUsingHTTPGet(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.args.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.http_get.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.http_get.0.path", "/healthz"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.http_get.0.port", "8080"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.http_get.0.http_header.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.http_get.0.http_header.0.name", "X-Custom-Header"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.http_get.0.http_header.0.value", "Awesome"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.initial_delay_seconds", "3"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_container_liveness_probe_using_tcp(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/liveness"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithLivenessProbeUsingTCP(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.args.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.tcp_socket.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.liveness_probe.0.tcp_socket.0.port", "8080"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_container_lifecycle(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "gcr.io/google_containers/liveness"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithLifeCycle(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.post_start.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.post_start.0.exec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.post_start.0.exec.0.command.#", "2"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.post_start.0.exec.0.command.0", "ls"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.post_start.0.exec.0.command.1", "-al"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.pre_stop.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.pre_stop.0.exec.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.pre_stop.0.exec.0.command.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.lifecycle.0.pre_stop.0.exec.0.command.0", "date"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_container_security_context(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithContainerSecurityContext(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.security_context.#", "1"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_volume_mount(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	secretName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithVolumeMounts(secretName, rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.image", imageName),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.0.mount_path", "/tmp/my_path"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.0.name", "db"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.0.read_only", "false"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.0.sub_path", ""),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_resource_requirements(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))

	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithResourceRequirements(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.image", imageName),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.resources.0.requests.0.memory", "50Mi"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.resources.0.requests.0.cpu", "250m"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.resources.0.limits.0.memory", "512Mi"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.resources.0.limits.0.cpu", "500m"),
				),
			},
		},
	})
}

func TestAccKubernetesReplicationController_with_empty_dir_volume(t *testing.T) {
	var conf api.ReplicationController

	rcName := fmt.Sprintf("tf-acc-test-%s", acctest.RandStringFromCharSet(10, acctest.CharSetAlphaNum))
	imageName := "nginx:1.7.9"

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckKubernetesReplicationControllerDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccKubernetesReplicationControllerConfigWithEmptyDirVolumes(rcName, imageName),
				Check: resource.ComposeAggregateTestCheckFunc(
					testAccCheckKubernetesReplicationControllerExists("kubernetes_replication_controller.test", &conf),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.image", imageName),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.#", "1"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.0.mount_path", "/cache"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.container.0.volume_mount.0.name", "cache-volume"),
					resource.TestCheckResourceAttr("kubernetes_replication_controller.test", "spec.0.template.0.spec.0.volume.0.empty_dir.0.medium", "Memory"),
				),
			},
		},
	})
}

func testAccCheckKubernetesReplicationControllerDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*KubeClientsets).MainClientset

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "kubernetes_replication_controller" {
			continue
		}

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		resp, err := conn.CoreV1().ReplicationControllers(namespace).Get(name, meta_v1.GetOptions{})
		if err == nil {
			if resp.Name == rs.Primary.ID {
				return fmt.Errorf("Replication Controller still exists: %s", rs.Primary.ID)
			}
		}
	}

	return nil
}

func testAccCheckKubernetesReplicationControllerExists(n string, obj *api.ReplicationController) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		conn := testAccProvider.Meta().(*KubeClientsets).MainClientset

		namespace, name, err := idParts(rs.Primary.ID)
		if err != nil {
			return err
		}

		out, err := conn.CoreV1().ReplicationControllers(namespace).Get(name, meta_v1.GetOptions{})
		if err != nil {
			return err
		}

		*obj = *out
		return nil
	}
}

func testAccKubernetesReplicationControllerConfig_basic(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    replicas = 1000 # This is intentionally high to exercise the waiter

    selector = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    template {
      metadata {
        labels = {
          TestLabelOne   = "one"
          TestLabelTwo   = "two"
          TestLabelThree = "three"
        }

        annotations = {
          TestAnnotationFive = "five"
        }
      }

      spec {
        container {
          image = "nginx:1.7.8"
          name  = "tf-acc-test"
        }
      }
    }
  }
}
`, name)
}

func testAccKubernetesReplicationControllerConfig_initContainer(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      TestAnnotationTwo = "two"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    replicas = 1000 # This is intentionally high to exercise the waiter

    selector = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    template {
      metadata {
        labels = {
          TestLabelOne   = "one"
          TestLabelTwo   = "two"
          TestLabelThree = "three"
        }
      }

      spec {
        container {
          name  = "nginx"
          image = "nginx"

          port {
            container_port = 80
          }

          volume_mount {
            name       = "workdir"
            mount_path = "/usr/share/nginx/html"
          }
        }

        init_container {
          name    = "install"
          image   = "busybox"
          command = ["wget", "-O", "/work-dir/index.html", "http://kubernetes.io"]

          volume_mount {
            name       = "workdir"
            mount_path = "/work-dir"
          }
        }

        dns_config {
          nameservers = ["1.1.1.1", "8.8.8.8", "9.9.9.9"]
          searches    = ["kubernetes.io"]

          option {
            name  = "ndots"
            value = 1
          }

          option {
            name = "use-vc"
          }
        }

        dns_policy = "Default"

        volume {
          name      = "workdir"
          empty_dir {}
        }
      }
    }
  }
}
`, name)
}

func testAccKubernetesReplicationControllerConfig_modified(name string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    annotations = {
      TestAnnotationOne = "one"
      Different         = "1234"
    }

    labels = {
      TestLabelOne   = "one"
      TestLabelThree = "three"
    }

    name = "%s"
  }

  spec {
    selector = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    template {
      metadata {
        labels = {
          TestLabelOne   = "one"
          TestLabelTwo   = "two"
          TestLabelThree = "three"
        }

        annotations = {
          TestAnnotationSix = "six"
        }
      }

      spec {
        container {
          image = "nginx:1.7.9"
          name  = "tf-acc-test"
        }
      }
    }
  }
}
`, name)
}

func testAccKubernetesReplicationControllerConfig_generatedName(prefix string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    labels = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    generate_name = "%s"
  }

  spec {
    selector = {
      TestLabelOne   = "one"
      TestLabelTwo   = "two"
      TestLabelThree = "three"
    }

    template {
      metadata {
        labels = {
          TestLabelOne   = "one"
          TestLabelTwo   = "two"
          TestLabelThree = "three"
        }
      }

      spec {
        container {
          image = "nginx:1.7.9"
          name  = "tf-acc-test"
        }
      }
    }
  }
}
`, prefix)
}

func testAccKubernetesReplicationControllerConfigWithSecurityContext(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        security_context {
          fs_group            = 100
          run_as_non_root     = true
          run_as_user         = 101
          supplemental_groups = [101]
        }

        container {
          image = "%s"
          name  = "containername"
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithSecurityContextRunAsGroup(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        security_context {
          fs_group            = 100
          run_as_group        = 100
          run_as_non_root     = true
          run_as_user         = 101
          supplemental_groups = [101]
        }

        container {
          image = "%s"
          name  = "containername"
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithLivenessProbeUsingExec(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"
          args  = ["/bin/sh", "-c", "touch /tmp/healthy; sleep 300; rm -rf /tmp/healthy; sleep 600"]

          liveness_probe {
            exec {
              command = ["cat", "/tmp/healthy"]
            }

            initial_delay_seconds = 5
            period_seconds        = 5
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithLivenessProbeUsingHTTPGet(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"
          args  = ["/server"]

          liveness_probe {
            http_get {
              path = "/healthz"
              port = 8080

              http_header {
                name  = "X-Custom-Header"
                value = "Awesome"
              }
            }

            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithLivenessProbeUsingTCP(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"
          args  = ["/server"]

          liveness_probe {
            tcp_socket {
              port = 8080
            }

            initial_delay_seconds = 3
            period_seconds        = 3
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithLifeCycle(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"
          args  = ["/server"]

          lifecycle {
            post_start {
              exec {
                command = ["ls", "-al"]
              }
            }

            pre_stop {
              exec {
                command = ["date"]
              }
            }
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithContainerSecurityContext(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"

          security_context {
            privileged  = true
            run_as_user = 1

            se_linux_options {
              level = "s0:c123,c456"
            }
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithVolumeMounts(secretName, rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_secret" "test" {
  metadata {
    name = "%s"
  }

  data = {
    one = "first"
  }
}

resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"

          volume_mount {
            mount_path = "/tmp/my_path"
            name       = "db"
          }
        }

        volume {
          name = "db"

          secret {
            secret_name = "${kubernetes_secret.test.metadata.0.name}"
          }
        }
      }
    }
  }
}
`, secretName, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithResourceRequirements(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"

          resources {
            limits {
              cpu    = "0.5"
              memory = "512Mi"
            }

            requests {
              cpu    = "250m"
              memory = "50Mi"
            }
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}

func testAccKubernetesReplicationControllerConfigWithEmptyDirVolumes(rcName, imageName string) string {
	return fmt.Sprintf(`
resource "kubernetes_replication_controller" "test" {
  metadata {
    name = "%s"

    labels = {
      Test = "TfAcceptanceTest"
    }
  }

  spec {
    selector = {
      Test = "TfAcceptanceTest"
    }

    template {
      metadata {
        labels = {
          Test = "TfAcceptanceTest"
        }
      }

      spec {
        container {
          image = "%s"
          name  = "containername"

          volume_mount {
            mount_path = "/cache"
            name       = "cache-volume"
          }
        }

        volume {
          name = "cache-volume"

          empty_dir {
            medium = "Memory"
          }
        }
      }
    }
  }
}
`, rcName, imageName)
}
