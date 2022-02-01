The `ExternalSecret` describes what data should be fetched, how the data should
be transformed and saved as a `Kind=Secret`:

* tells the operator what secrets should be synced by using `spec.data` to
  explicitly sync individual keys or use `spec.dataFrom` to get **all values**
  from the external API.
* you can also use `spec.dataFrom` to sync many secrets at once, based on a
  regular expression of their name or on tags/attributes. You need to check if your provider implements this behavior or not.
* you can specify how the secret should look like by specifying a
  `spec.target.template`

## Template

When the controller reconciles the `ExternalSecret` it will use the `spec.template` as a blueprint to construct a new `Kind=Secret`. You can use golang templates to define the blueprint and use template functions to transform secret values. You can also pull in `ConfigMaps` that contain golang-template data using `templateFrom`. See [advanced templating](guides-templating.md) for details.

## Update Behavior

The `Kind=Secret` is updated when:

* the `spec.refreshInterval` has passed and is not `0`
* the `ExternalSecret`'s `labels` or `annotations` are changed
* the `ExternalSecret`'s `spec` has been changed

You can trigger a secret refresh by using kubectl or any other kubernetes api client:

```
kubectl annotate es my-es force-sync=$(date +%s) --overwrite
```

## Example

Take a look at an annotated example to understand the design behind the
`ExternalSecret`.

``` yaml
{% include 'full-external-secret.yaml' %}
```