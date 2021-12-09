package translator

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/enterprise/options/ratelimit"
	rltypes "github.com/solo-io/solo-apis/pkg/api/ratelimit.solo.io/v1alpha1"
	"github.com/solo-io/solo-kit/pkg/utils/prototime"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/api/v1/options/retries"
)

var _ = Describe("MergeRoutePlugins", func() {
	It("merges top-level route plugins fields", func() {
		dst := &v1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{Value: "preserve-me"},
			Retries: &retries.RetryPolicy{
				RetryOn:    "5XX",
				NumRetries: 0, // should not overwrite this field
			},
		}
		d := prototime.DurationToProto(time.Minute)
		src := &v1.RouteOptions{
			Timeout: d,
			Retries: &retries.RetryPolicy{
				RetryOn:    "5XX",
				NumRetries: 3, // do not overwrite 0 value above
			},
			RatelimitBasic: &ratelimit.IngressRateLimit{
				AuthorizedLimits: &rltypes.RateLimit{
					Unit:            1,
					RequestsPerUnit: 2,
				},
			},
		}
		expected := &v1.RouteOptions{
			PrefixRewrite: &wrappers.StringValue{Value: "preserve-me"},
			Timeout:       d,
			Retries: &retries.RetryPolicy{
				RetryOn:    "5XX",
				NumRetries: 0,
			},
			RatelimitBasic: &ratelimit.IngressRateLimit{
				AuthorizedLimits: &rltypes.RateLimit{
					Unit:            1,
					RequestsPerUnit: 2,
				},
			},
		}

		actual := mergeRouteOptions(dst, src)
		Expect(actual).To(Equal(expected))
	})
})
