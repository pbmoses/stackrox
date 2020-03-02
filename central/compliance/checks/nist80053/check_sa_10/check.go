package checksa10

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
)

const (
	controlID = "NIST_SP_800_53:SA_10"

	phase = storage.LifecycleStage_DEPLOY
)

var (
	interpretationText = `This control requires definition of acceptable configuration practices and ongoing monitoring of compliance with these practices.

For this control, ` + common.AnyPolicyInLifeCycleInterpretation(phase)
)

func init() {
	framework.MustRegisterNewCheckIfFlagEnabled(
		framework.CheckMetadata{
			ID:                 controlID,
			Scope:              framework.ClusterKind,
			DataDependencies:   []string{"Policies"},
			InterpretationText: interpretationText,
		},
		func(ctx framework.ComplianceContext) {
			common.CheckAnyPolicyInLifeCycle(ctx, phase)
		}, features.NistSP800_53)
}
