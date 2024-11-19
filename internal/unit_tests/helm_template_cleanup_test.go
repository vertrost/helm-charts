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
			"--set", "neo4j.name=neo4j",
			"--set", "volumes.data.mode=defaultStorageClass",
			"--set", "podSpec.annotations.sidecar\\.istio\\.io/inject=true",
			"--set", "podSpec.annotations.custom\\.annotation/test=value",
		}

		if edition == "enterprise" {
			baseArgs = append(baseArgs,
				"--set", "neo4j.edition=enterprise",
				"--set", "neo4j.acceptLicenseAgreement=yes")
		} else {
			baseArgs = append(baseArgs, "--set", "neo4j.edition=community")
		}

		// Add the hook-delete flag to generate cleanup job
		baseArgs = append(baseArgs, "--hooks")

		manifest, err := model.HelmTemplate(t, chart, baseArgs)
		if !assert.NoError(t, err) {
			t.Logf("Template error: %v", err)
			return
		}

		// Debug: Print all jobs in the manifest
		jobs := manifest.OfType(&batchv1.Job{})
		t.Logf("Found %d jobs in manifest", len(jobs))
		for _, obj := range jobs {
			if job, ok := obj.(*batchv1.Job); ok {
				t.Logf("Job name: %s", job.ObjectMeta.Name)
			}
		}

		cleanupJob := manifest.OfTypeWithName(
			&batchv1.Job{},
			fmt.Sprintf("%s-cleanup", model.DefaultHelmTemplateReleaseName.String()),
		)

		if !assert.NotNil(t, cleanupJob, "cleanup job not found") {
			return
		}

		// Check annotations
		podAnnotations := cleanupJob.(*batchv1.Job).Spec.Template.ObjectMeta.Annotations
		assert.Equal(t, "true", podAnnotations["sidecar.istio.io/inject"], "Custom Istio sidecar injection setting should be respected")
		assert.Equal(t, "value", podAnnotations["custom.annotation/test"], "Custom annotation should be present")
	}))
}
