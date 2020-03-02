package checkcm7

import (
	"github.com/stackrox/rox/central/compliance/checks/common"
	"github.com/stackrox/rox/central/compliance/framework"
	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/features"
	"github.com/stackrox/rox/pkg/stringutils"
)

const (
	controlID = "NIST_SP_800_53:CM_7"

	interpretationText = `This control requires that unnecessary features be prohibited or restricted.

For this control, StackRox validates that at least one policy is enabled and enforced based on: 1) port exposure or service exposure level, and 2) runtime behavior.`
)

func isPortExposeOrExposureLevelPolicy(p *storage.Policy) bool {
	if p.GetFields().GetPortPolicy() != nil {
		return true
	}
	if len(p.GetFields().GetPortExposurePolicy().GetExposureLevels()) > 0 {
		return true
	}

	return false
}

func init() {
	framework.MustRegisterNewCheckIfFlagEnabled(
		framework.CheckMetadata{
			ID:                 controlID,
			Scope:              framework.ClusterKind,
			DataDependencies:   []string{"Policies"},
			InterpretationText: interpretationText,
		},
		func(ctx framework.ComplianceContext) {
			policies := ctx.Data().Policies()
			var portExposePolicy, runtimePolicy string
			for name, p := range policies {
				if !common.IsPolicyEnabled(p) || !common.IsPolicyEnforced(p) {
					continue
				}
				if portExposePolicy == "" && isPortExposeOrExposureLevelPolicy(p) {
					portExposePolicy = name
				}
				if runtimePolicy == "" && common.PolicyIsInLifecycleStage(p, storage.LifecycleStage_RUNTIME) {
					runtimePolicy = name
				}
			}
			if stringutils.AllNotEmpty(portExposePolicy, runtimePolicy) {
				framework.Passf(ctx, "At least one policy regarding port exposure/exposure level (%q) and at least one runtime policy (%q) are enabled and enforced", portExposePolicy, runtimePolicy)
				return
			}
			framework.Fail(ctx, "Required, but could not find, at least one policy regarding port exposure and at least one runtime policy that is enabled and enforced")
		}, features.NistSP800_53)
}
