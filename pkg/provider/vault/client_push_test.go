/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package vault

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"testing"

	corev1 "k8s.io/api/core/v1"

	esv1beta1 "github.com/external-secrets/external-secrets/apis/externalsecrets/v1beta1"
	testingfake "github.com/external-secrets/external-secrets/pkg/provider/testing/fake"
	"github.com/external-secrets/external-secrets/pkg/provider/vault/fake"
	"github.com/external-secrets/external-secrets/pkg/provider/vault/util"
)

func TestDeleteSecret(t *testing.T) {
	type args struct {
		store    *esv1beta1.VaultProvider
		vLogical util.Logical
	}

	type want struct {
		err error
	}
	tests := map[string]struct {
		reason string
		args   args
		ref    *testingfake.PushSecretData
		want   want
		value  []byte
	}{
		"DeleteSecretNoOpKV1": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, nil),
					WriteWithContextFn:        fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn:       fake.ExpectDeleteWithContextNoCall(),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretNoOpKV2": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, nil),
					WriteWithContextFn:        fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn:       fake.ExpectDeleteWithContextNoCall(),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretFailIfErrorKV1": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, fmt.Errorf("failed to read")),
					WriteWithContextFn:        fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn:       fake.ExpectDeleteWithContextNoCall(),
				},
			},
			want: want{
				err: fmt.Errorf("failed to read"),
			},
		},
		"DeleteSecretFailIfErrorKV2": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, fmt.Errorf("failed to read")),
					WriteWithContextFn:        fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn:       fake.ExpectDeleteWithContextNoCall(),
				},
			},
			want: want{
				err: fmt.Errorf("failed to read"),
			},
		},
		"DeleteSecretNotManagedKV1": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "another-secret-tool",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretNotManagedKV2": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "another-secret-tool",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretSuccessKV1": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretSuccessKV2": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretErrorKV1": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, fmt.Errorf("failed to delete")),
				},
			},
			want: want{
				err: fmt.Errorf("failed to delete"),
			},
		},
		"DeleteSecretErrorKV2": {
			reason: "No secret is because it does not exist",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, fmt.Errorf("failed to delete")),
				},
			},
			want: want{
				err: fmt.Errorf("failed to delete"),
			},
		},
		"DeleteSecretUpdatePropertyKV1": {
			reason: "Secret should only be updated if Property is set",
			ref:    &testingfake.PushSecretData{RemoteKey: "secret", Property: "fake-key"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value",
						"foo":      "bar",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextValue(map[string]interface{}{
						"foo": "bar",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						}}),
					DeleteWithContextFn: fake.ExpectDeleteWithContextNoCall(),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretUpdatePropertyKV2": {
			reason: "Secret should only be updated if Property is set",
			ref:    &testingfake.PushSecretData{RemoteKey: "secret", Property: "fake-key"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value",
							"foo":      "bar",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextValue(map[string]interface{}{"data": map[string]interface{}{"foo": "bar"}}),
					DeleteWithContextFn: fake.ExpectDeleteWithContextNoCall(),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretIfNoOtherPropertiesKV1": {
			reason: "Secret should only be deleted if no other properties are set",
			ref:    &testingfake.PushSecretData{RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"foo": "bar",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"DeleteSecretIfNoOtherPropertiesKV2": {
			reason: "Secret should only be deleted if no other properties are set",
			ref:    &testingfake.PushSecretData{RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"foo": "bar",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn:  fake.ExpectWriteWithContextNoCall(),
					DeleteWithContextFn: fake.NewDeleteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			ref := testingfake.PushSecretData{RemoteKey: "secret", Property: ""}
			if tc.ref != nil {
				ref = *tc.ref
			}
			client := &client{
				logical: tc.args.vLogical,
				store:   tc.args.store,
			}
			err := client.DeleteSecret(context.Background(), ref)

			// Error nil XOR tc.want.err nil
			if ((err == nil) || (tc.want.err == nil)) && !((err == nil) && (tc.want.err == nil)) {
				t.Errorf("\nTesting DeleteSecret:\nName: %v\nReason: %v\nWant error: %v\nGot error: %v", name, tc.reason, tc.want.err, err)
			}

			// if errors are the same type but their contents do not match
			if err != nil && tc.want.err != nil {
				if !strings.Contains(err.Error(), tc.want.err.Error()) {
					t.Errorf("\nTesting DeleteSecret:\nName: %v\nReason: %v\nWant error: %v\nGot error got nil", name, tc.reason, tc.want.err)
				}
			}
		})
	}
}
func TestPushSecret(t *testing.T) {
	secretKey := "secret-key"
	noPermission := errors.New("no permission")
	type args struct {
		store    *esv1beta1.VaultProvider
		vLogical util.Logical
	}

	type want struct {
		err error
	}
	tests := map[string]struct {
		reason string
		args   args
		want   want
		data   *testingfake.PushSecretData
		value  []byte
	}{
		"SetSecretKV1": {
			reason: "secret is successfully set, with no existing vault secret",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, nil),
					WriteWithContextFn:        fake.NewWriteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"SetSecretKV2": {
			reason: "secret is successfully set, with no existing vault secret",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, nil),
					WriteWithContextFn:        fake.NewWriteWithContextFn(nil, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"SetSecretWithWriteErrorKV1": {
			reason: "secret cannot be pushed if write fails",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, nil),
					WriteWithContextFn:        fake.NewWriteWithContextFn(nil, noPermission),
				},
			},
			want: want{
				err: noPermission,
			},
		},
		"SetSecretWithWriteErrorKV2": {
			reason: "secret cannot be pushed if write fails",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, nil),
					WriteWithContextFn:        fake.NewWriteWithContextFn(nil, noPermission),
				},
			},
			want: want{
				err: noPermission,
			},
		},
		"SetSecretEqualsPushSecretV1": {
			reason: "vault secret kv equals secret to push kv",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"SetSecretEqualsPushSecretV2": {
			reason: "vault secret kv equals secret to push kv",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
				},
			},
			want: want{
				err: nil,
			},
		},
		"PushSecretPropertyKV1": {
			reason: "push secret with property adds the property",
			value:  []byte("fake-value"),
			data:   &testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextValue(map[string]interface{}{
						"fake-key": "fake-value",
						"custom_metadata": map[string]string{
							"managed-by": "external-secrets",
						},
						"foo": "fake-value",
					}),
				},
			},
			want: want{
				err: nil,
			},
		},
		"PushSecretPropertyKV2": {
			reason: "push secret with property adds the property",
			value:  []byte("fake-value"),
			data:   &testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextValue(map[string]interface{}{"data": map[string]interface{}{"fake-key": "fake-value", "foo": "fake-value"}}),
				},
			},
			want: want{
				err: nil,
			},
		},
		"PushSecretUpdatePropertyKV1": {
			reason: "push secret with property only updates the property",
			value:  []byte("new-value"),
			data:   &testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"foo": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextValue(map[string]interface{}{
						"foo": "new-value",
						"custom_metadata": map[string]string{
							"managed-by": "external-secrets",
						},
					}),
				},
			},
			want: want{
				err: nil,
			},
		},
		"PushSecretUpdatePropertyKV2": {
			reason: "push secret with property only updates the property",
			value:  []byte("new-value"),
			data:   &testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"foo": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextValue(map[string]interface{}{"data": map[string]interface{}{"foo": "new-value"}}),
				},
			},
			want: want{
				err: nil,
			},
		},
		"PushSecretPropertyNoUpdateKV1": {
			reason: "push secret with property only updates the property",
			value:  []byte("fake-value"),
			data:   &testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"foo": "fake-value",
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextNoCall(),
				},
			},
			want: want{
				err: nil,
			},
		},
		"PushSecretPropertyNoUpdateKV2": {
			reason: "push secret with property only updates the property",
			value:  []byte("fake-value"),
			data:   &testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: "foo"},
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"foo": "fake-value",
						},
						"custom_metadata": map[string]interface{}{
							"managed-by": "external-secrets",
						},
					}, nil),
					WriteWithContextFn: fake.ExpectWriteWithContextNoCall(),
				},
			},
			want: want{
				err: nil,
			},
		},
		"SetSecretErrorReadingSecretKV1": {
			reason: "error occurs if secret cannot be read",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, noPermission),
				},
			},
			want: want{
				err: fmt.Errorf(errReadSecret, noPermission),
			},
		},
		"SetSecretErrorReadingSecretKV2": {
			reason: "error occurs if secret cannot be read",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(nil, noPermission),
				},
			},
			want: want{
				err: fmt.Errorf(errReadSecret, noPermission),
			},
		},
		"SetSecretNotManagedByESOV1": {
			reason: "a secret not managed by ESO cannot be updated",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV1).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"fake-key": "fake-value2",
						"custom_metadata": map[string]interface{}{
							"managed-by": "not-external-secrets",
						},
					}, nil),
				},
			},
			want: want{
				err: errors.New("secret not managed by external-secrets"),
			},
		},
		"SetSecretNotManagedByESOV2": {
			reason: "a secret not managed by ESO cannot be updated",
			args: args{
				store: makeValidSecretStoreWithVersion(esv1beta1.VaultKVStoreV2).Spec.Provider.Vault,
				vLogical: &fake.Logical{
					ReadWithDataWithContextFn: fake.NewReadWithContextFn(map[string]interface{}{
						"data": map[string]interface{}{
							"fake-key": "fake-value2",
							"custom_metadata": map[string]interface{}{
								"managed-by": "not-external-secrets",
							},
						},
					}, nil),
				},
			},
			want: want{
				err: errors.New("secret not managed by external-secrets"),
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			data := testingfake.PushSecretData{SecretKey: secretKey, RemoteKey: "secret", Property: ""}
			if tc.data != nil {
				data = *tc.data
			}
			client := &client{
				logical: tc.args.vLogical,
				store:   tc.args.store,
			}
			val := tc.value
			if val == nil {
				val = []byte(`{"fake-key":"fake-value"}`)
			}
			s := &corev1.Secret{Data: map[string][]byte{secretKey: val}}
			err := client.PushSecret(context.Background(), s, data)

			// Error nil XOR tc.want.err nil
			if ((err == nil) || (tc.want.err == nil)) && !((err == nil) && (tc.want.err == nil)) {
				t.Errorf("\nTesting SetSecret:\nName: %v\nReason: %v\nWant error: %v\nGot error: %v", name, tc.reason, tc.want.err, err)
			}

			// if errors are the same type but their contents do not match
			if err != nil && tc.want.err != nil {
				if !strings.Contains(err.Error(), tc.want.err.Error()) {
					t.Errorf("\nTesting SetSecret:\nName: %v\nReason: %v\nWant error: %v\nGot error got nil", name, tc.reason, tc.want.err)
				}
			}
		})
	}
}
