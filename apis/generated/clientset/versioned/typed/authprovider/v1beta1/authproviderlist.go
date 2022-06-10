// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/stackrox/stackrox/apis/authprovider/v1beta1"
	scheme "github.com/stackrox/stackrox/apis/generated/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AuthProviderListsGetter has a method to return a AuthProviderListInterface.
// A group's client should implement this interface.
type AuthProviderListsGetter interface {
	AuthProviderLists(namespace string) AuthProviderListInterface
}

// AuthProviderListInterface has methods to work with AuthProviderList resources.
type AuthProviderListInterface interface {
	Create(ctx context.Context, authProviderList *v1beta1.AuthProviderList, opts v1.CreateOptions) (*v1beta1.AuthProviderList, error)
	Update(ctx context.Context, authProviderList *v1beta1.AuthProviderList, opts v1.UpdateOptions) (*v1beta1.AuthProviderList, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.AuthProviderList, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.AuthProviderListList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.AuthProviderList, err error)
	AuthProviderListExpansion
}

// authProviderLists implements AuthProviderListInterface
type authProviderLists struct {
	client rest.Interface
	ns     string
}

// newAuthProviderLists returns a AuthProviderLists
func newAuthProviderLists(c *AuthproviderV1beta1Client, namespace string) *authProviderLists {
	return &authProviderLists{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the authProviderList, and returns the corresponding authProviderList object, and an error if there is any.
func (c *authProviderLists) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.AuthProviderList, err error) {
	result = &v1beta1.AuthProviderList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("authproviderlists").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AuthProviderLists that match those selectors.
func (c *authProviderLists) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.AuthProviderListList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.AuthProviderListList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("authproviderlists").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested authProviderLists.
func (c *authProviderLists) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("authproviderlists").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a authProviderList and creates it.  Returns the server's representation of the authProviderList, and an error, if there is any.
func (c *authProviderLists) Create(ctx context.Context, authProviderList *v1beta1.AuthProviderList, opts v1.CreateOptions) (result *v1beta1.AuthProviderList, err error) {
	result = &v1beta1.AuthProviderList{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("authproviderlists").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(authProviderList).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a authProviderList and updates it. Returns the server's representation of the authProviderList, and an error, if there is any.
func (c *authProviderLists) Update(ctx context.Context, authProviderList *v1beta1.AuthProviderList, opts v1.UpdateOptions) (result *v1beta1.AuthProviderList, err error) {
	result = &v1beta1.AuthProviderList{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("authproviderlists").
		Name(authProviderList.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(authProviderList).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the authProviderList and deletes it. Returns an error if one occurs.
func (c *authProviderLists) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("authproviderlists").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *authProviderLists) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("authproviderlists").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched authProviderList.
func (c *authProviderLists) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.AuthProviderList, err error) {
	result = &v1beta1.AuthProviderList{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("authproviderlists").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
