package branding

const (
	productNameRHACS    = "Red Hat Advanced Cluster Security for Kubernetes"
	productNameStackrox = "StackRox"

	productNameRHACSShort = "RHACS"
)

// GetProductName returns the product name based on ProductBranding env variable
func GetProductName() string {
	if getProductBrandingEnv() == ProductBrandingRHACS {
		return productNameRHACS
	}
	return productNameStackrox
}

// GetProductNameShort returns the short form of the product name based on ProductBranding env variable
func GetProductNameShort() string {
	if getProductBrandingEnv() == ProductBrandingRHACS {
		return productNameRHACSShort
	}
	return productNameStackrox
}
