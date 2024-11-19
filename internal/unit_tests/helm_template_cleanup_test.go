package unit_tests

import (
	"fmt"
	"testing"

	"github.com/neo4j/helm-charts/internal/model"
	"github.com/stretchr/testify/assert"
	batchv1 "k8s.io/api/batch/v1"
)

func TestCleanupJobPodAnnotations(t *testing.T) {
	t.Parallel()

	forEachPrimaryChart(t, andEachSupportedEdition(func(t *testing.T, chart model.Neo4jHelmChartBuilder, edition string) {
		baseArgs := []string{
			"--set", "services.neo4j.enabled=true",
			"--set", "services.neo4j.cleanup.enabled=true",
			"--set", "neo4j.name=neo4j",
			"--set", "volumes.data.mode=defaultStorageClass",
		}

		// Add edition-specific args
		if edition == "enterprise" {
			baseArgs = append(baseArgs, "--set", "neo4j.edition=enterprise", "--set", "neo4j.acceptLicenseAgreement=yes")
		} else {
			baseArgs = append(baseArgs, "--set", "neo4j.edition=community")
		}

		// Test default annotations
		manifest, err := model.HelmTemplate(t, chart, baseArgs)
		if !assert.NoError(t, err) {
			return
		}

		cleanupJob := manifest.OfTypeWithName(
			&batchv1.Job{},
			fmt.Sprintf("%s-cleanup", model.DefaultHelmTemplateReleaseName.String()),
		).(*batchv1.Job)

		if !assert.NotNil(t, cleanupJob, "cleanup job not found") {
			return
		}

		// Check default annotation
		podAnnotations := cleanupJob.Spec.Template.ObjectMeta.Annotations
		assert.Equal(t, "false", podAnnotations["sidecar.istio.io/inject"], "Default Istio sidecar injection should be disabled")

		// Test custom annotations
		customAnnotationsArgs := append(
			baseArgs,
			"--set", "services.neo4j.cleanup.podAnnotations.sidecar\\.istio\\.io/inject=true",
			"--set", "services.neo4j.cleanup.podAnnotations.custom\\.annotation/test=value",
		)

		manifest, err = model.HelmTemplate(t, chart, customAnnotationsArgs)
		if !assert.NoError(t, err) {
			return
		}

		cleanupJob = manifest.OfTypeWithName(
			&batchv1.Job{},
			fmt.Sprintf("%s-cleanup", model.DefaultHelmTemplateReleaseName.String()),
		).(*batchv1.Job)

		// Check custom annotations
		podAnnotations = cleanupJob.Spec.Template.ObjectMeta.Annotations
		assert.Equal(t, "true", podAnnotations["sidecar.istio.io/inject"], "Custom Istio sidecar injection setting should be respected")
		assert.Equal(t, "value", podAnnotations["custom.annotation/test"], "Custom annotation should be present")
	}))
}
