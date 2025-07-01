/*
Copyright 2021 Upbound Inc.
*/

package clients

import (
	"context"
	"encoding/json"

	"github.com/crossplane/crossplane-runtime/pkg/resource"
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/crossplane/upjet/pkg/terraform"

	"github.com/supahlab/provider-tailscale/apis/v1beta1"
)

const (
	// error messages
	errNoProviderConfig     = "no providerConfigRef provided"
	errGetProviderConfig    = "cannot get referenced ProviderConfig"
	errTrackUsage           = "cannot track ProviderConfig usage"
	errExtractCredentials   = "cannot extract credentials"
	errUnmarshalCredentials = "cannot unmarshal tailscale credentials as JSON"
)

const (
	keyBaseURL = "base_url" // (String) The base URL of the Tailscale API. Defaults to https://api.tailscale.com. Can be set via the TAILSCALE_BASE_URL environment variable.
	keyAPIKey = "api_key" // (String, Sensitive) The API key to use for authenticating requests to the API. Can be set via the TAILSCALE_API_KEY environment variable. Conflicts with 'oauth_client_id' and 'oauth_client_secret'.
	keyOAuthClientID = "oauth_client_id" // (String) The OAuth application's ID when using OAuth client credentials. Can be set via the TAILSCALE_OAUTH_CLIENT_ID environment variable. Both 'oauth_client_id' and 'oauth_client_secret' must be set. Conflicts with 'api_key'.
	keyOAuthClientSecret = "oauth_client_secret" // (String, Sensitive) The OAuth application's secret when using OAuth client credentials. Can be set via the TAILSCALE_OAUTH_CLIENT_SECRET environment variable. Both 'oauth_client_id' and 'oauth_client_secret' must be set. Conflicts with 'api_key'.
	keyOAuthScopes = "scopes" // (List of String) The OAuth 2.0 scopes to request for the access token generated using the supplied OAuth client credentials. See https://tailscale.com/kb/1215/oauth-clients/#scopes for available scopes. Only valid when both 'oauth_client_id' and 'oauth_client_secret' are set.
	keyTailnet = "tailnet" // (String) The organization name of the Tailnet in which to perform actions. Can be set via the TAILSCALE_TAILNET environment variable. Default is the tailnet that owns API credentials passed to the provider.
	keyUserAgent = "user_agent" // user_agent (String) User-Agent header for API requests.
)

// TerraformSetupBuilder builds Terraform a terraform.SetupFn function which
// returns Terraform provider setup configuration
func TerraformSetupBuilder(version, providerSource, providerVersion string) terraform.SetupFn {
	return func(ctx context.Context, client client.Client, mg resource.Managed) (terraform.Setup, error) {
		ps := terraform.Setup{
			Version: version,
			Requirement: terraform.ProviderRequirement{
				Source:  providerSource,
				Version: providerVersion,
			},
		}

		configRef := mg.GetProviderConfigReference()
		if configRef == nil {
			return ps, errors.New(errNoProviderConfig)
		}
		pc := &v1beta1.ProviderConfig{}
		if err := client.Get(ctx, types.NamespacedName{Name: configRef.Name}, pc); err != nil {
			return ps, errors.Wrap(err, errGetProviderConfig)
		}

		t := resource.NewProviderConfigUsageTracker(client, &v1beta1.ProviderConfigUsage{})
		if err := t.Track(ctx, mg); err != nil {
			return ps, errors.Wrap(err, errTrackUsage)
		}

		data, err := resource.CommonCredentialExtractor(ctx, pc.Spec.Credentials.Source, client, pc.Spec.Credentials.CommonCredentialSelectors)
		if err != nil {
			return ps, errors.Wrap(err, errExtractCredentials)
		}
		creds := map[string]string{}
		if err := json.Unmarshal(data, &creds); err != nil {
			return ps, errors.Wrap(err, errUnmarshalCredentials)
		}
		
		ps.Configuration = map[string]any{}
		if v, ok := creds[keyAPIKey]; ok {
		  ps.Configuration[keyAPIKey] = v
		}
		if v, ok := creds[keyBaseURL]; ok {
		  ps.Configuration[keyBaseURL] = v
		}
		if v, ok := creds[keyOAuthClientID]; ok {
		  ps.Configuration[keyOAuthClientID] = v
		}
		if v, ok := creds[keyOAuthClientSecret]; ok {
		  ps.Configuration[keyOAuthClientSecret] = v
		}
		if v, ok := creds[keyOAuthScopes]; ok {
		  ps.Configuration[keyOAuthScopes] = v
		}
		if v, ok := creds[keyTailnet]; ok {
		  ps.Configuration[keyTailnet] = v
		}
		if v, ok := creds[keyUserAgent]; ok {
		  ps.Configuration[keyUserAgent] = v
		}
		return ps, nil
	}
}
