package enricher

import (
	"context"
	"time"

	"github.com/stackrox/rox/generated/storage"
	"github.com/stackrox/rox/pkg/images/integration"
	"github.com/stackrox/rox/pkg/integrationhealth"
	"github.com/stackrox/rox/pkg/logging"
	pkgMetrics "github.com/stackrox/rox/pkg/metrics"
	registryTypes "github.com/stackrox/rox/pkg/registries/types"
	scannerTypes "github.com/stackrox/rox/pkg/scanners/types"
	"golang.org/x/time/rate"
)

var (
	log = logging.LoggerForModule()
)

// FetchOption determines what attempts should be made to retrieve the metadata
type FetchOption int

// These are all the possible fetch options for the enricher
const (
	Default FetchOption = iota
	NoInlineScan
	RefetchScans
)

// EnrichmentContext is used to pass options through the enricher without exploding the number of function arguments
type EnrichmentContext struct {
	// FetchOpt define constraints about using external data
	FetchOpt FetchOption

	// EnforcementOnly indicates that we don't care about any violations unless they have enforcement enabled.
	EnforcementOnly bool

	// UseNonBlockingCallsWherePossible tells the enricher to make non-blocking calls to image scanners where that is
	// possible. Note that, if NoExternalMetadata is true, this param is irrelevant since no external calls are made at all.
	UseNonBlockingCallsWherePossible bool

	// Internal is used to indicate when the caller is internal.
	// This is used to indicate that we do not want to fail upon failing to find integrations.
	Internal bool
}

// EnrichmentResult denotes possible return values of the EnrichImage function.
type EnrichmentResult struct {
	// ImageUpdated returns whether or not the image was updated, either with metadata or with a scan.
	ImageUpdated bool

	ScanResult ScanResult
}

// A ScanResult denotes the result of an attempt to scan an image.
//go:generate stringer -type=ScanResult
type ScanResult int

const (
	// ScanNotDone denotes that the image was not scanned.
	ScanNotDone ScanResult = iota
	// ScanTriggered denotes that the image was not scanned, but that non-blocking API requests were made
	// to request scanning.
	ScanTriggered
	// ScanSucceeded denotes that the image was successfully scanned.
	ScanSucceeded
)

// ImageEnricher provides functions for enriching images with integrations.
//go:generate mockgen-wrapper
type ImageEnricher interface {
	EnrichImage(ctx EnrichmentContext, image *storage.Image) (EnrichmentResult, error)
}

type cveSuppressor interface {
	EnrichImageWithSuppressedCVEs(image *storage.Image)
}

type imageGetter interface {
	GetImage(ctx context.Context, id string) (*storage.Image, bool, error)
}

// New returns a new ImageEnricher instance for the given subsystem.
// (The subsystem is just used for Prometheus metrics.)
func New(cvesSuppressor cveSuppressor, cvesSuppressorV2 cveSuppressor, is integration.Set, subsystem pkgMetrics.Subsystem, imageGetter imageGetter, healthReporter integrationhealth.Reporter) ImageEnricher {
	enricher := &enricherImpl{
		cvesSuppressor:   cvesSuppressor,
		cvesSuppressorV2: cvesSuppressorV2,
		integrations:     is,

		// number of consecutive errors per registry or scanner to ascertain health of the integration
		errorsPerRegistry:         make(map[registryTypes.ImageRegistry]int32, len(is.RegistrySet().GetAll())),
		errorsPerScanner:          make(map[scannerTypes.ImageScanner]int32, len(is.ScannerSet().GetAll())),
		integrationHealthReporter: healthReporter,

		metadataLimiter: rate.NewLimiter(rate.Every(50*time.Millisecond), 1),
		asyncRateLimiter: rate.NewLimiter(rate.Every(1*time.Second), 5),

		imageGetter: imageGetter,

		metrics: newMetrics(subsystem),
	}
	return enricher
}
