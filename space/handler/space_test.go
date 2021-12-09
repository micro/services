package handler

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	sthree "github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3iface"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/memory"
	pb "github.com/micro/services/space/proto"

	. "github.com/onsi/gomega"
)

type mockS3Client struct {
	s3iface.S3API
	head   func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error)
	put    func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error)
	delete func(input *sthree.DeleteObjectInput) (*sthree.DeleteObjectOutput, error)
	list   func(input *sthree.ListObjectsInput) (*sthree.ListObjectsOutput, error)
	get    func(input *sthree.GetObjectInput) (*sthree.GetObjectOutput, error)
}

func (m mockS3Client) HeadObject(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
	if m.head != nil {
		return m.head(input)
	}
	return &sthree.HeadObjectOutput{}, nil
}

func (m mockS3Client) PutObject(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
	if m.put != nil {
		return m.put(input)
	}
	return &sthree.PutObjectOutput{}, nil
}

func (m mockS3Client) DeleteObject(input *sthree.DeleteObjectInput) (*sthree.DeleteObjectOutput, error) {
	if m.delete != nil {
		return m.delete(input)
	}
	return &sthree.DeleteObjectOutput{}, nil
}

func (m mockS3Client) ListObjects(input *sthree.ListObjectsInput) (*sthree.ListObjectsOutput, error) {
	if m.list != nil {
		return m.list(input)
	}
	return &sthree.ListObjectsOutput{}, nil
}

func (m mockS3Client) GetObject(input *sthree.GetObjectInput) (*sthree.GetObjectOutput, error) {
	if m.get != nil {
		return m.get(input)
	}
	return &sthree.GetObjectOutput{}, nil
}

type mockError struct {
	code    string
	message string
	err     string
}

func (m mockError) Error() string {
	return m.err
}

func (m mockError) Code() string {
	return m.code
}

func (m mockError) Message() string {
	return m.message
}

func (m mockError) OrigErr() error {
	return fmt.Errorf(m.err)
}

func TestCreate(t *testing.T) {
	g := NewWithT(t)
	tcs := []struct {
		name       string
		objName    string
		visibility string
		err        error
		url        string
		head       func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error)
		put        func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error)
	}{
		{
			name:    "Simple case",
			objName: "foo.jpg",
			url:     "",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				return nil, mockError{code: "NotFound"}
			},
			put: func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(input.ACL).To(BeNil())
				g.Expect(*input.Metadata[mdVisibility]).To(Equal(visibilityPrivate))
				g.Expect(input.Metadata[mdCreated]).To(Not(BeNil()))
				return &sthree.PutObjectOutput{}, nil
			},
		},
		{
			name:       "Public object",
			objName:    "bar/baz/foo.jpg",
			visibility: "public",
			url:        "https://my-space.ams3.example.com/micro/123/bar/baz/foo.jpg",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				return nil, mockError{code: "NotFound"}
			},
			put: func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.ACL).To(Equal(mdACLPublic))
				g.Expect(*input.Metadata[mdVisibility]).To(Equal(visibilityPublic))
				g.Expect(input.Metadata[mdCreated]).To(Not(BeNil()))
				return &sthree.PutObjectOutput{}, nil
			},
		},
		{
			name:    "Missing name",
			objName: "",
			err:     errors.BadRequest("space.Create", "Missing name param"),
		},
		{
			name:    "Already exists",
			objName: "foo.jpg",
			err:     errors.BadRequest("space.Create", "Object already exists"),
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.Key).To(Equal("micro/123/foo.jpg"))
				return &sthree.HeadObjectOutput{}, nil
			},
		},
	}
	store.DefaultStore = memory.NewStore()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			handler := Space{
				conf: conf{
					AccessKey: "access",
					SecretKey: "secret",
					Endpoint:  "example.com",
					SpaceName: "my-space",
					SSL:       true,
					Region:    "ams3",
					BaseURL:   "https://my-space.ams3.example.com",
				},
				client: &mockS3Client{head: tc.head, put: tc.put},
			}
			ctx := context.Background()
			ctx = auth.ContextWithAccount(ctx, &auth.Account{
				ID:       "123",
				Type:     "user",
				Issuer:   "micro",
				Metadata: map[string]string{},
				Scopes:   []string{"space"},
				Name:     "john@example.com",
			})
			rsp := pb.CreateResponse{}
			err := handler.Create(ctx, &pb.CreateRequest{
				Object:     []byte("foobar"),
				Name:       tc.objName,
				Visibility: tc.visibility,
			}, &rsp)
			if tc.err != nil {
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(rsp.Url).To(Equal(tc.url))
			}

		})
	}

}

func TestUpdate(t *testing.T) {
	g := NewWithT(t)
	tcs := []struct {
		name       string
		objName    string
		visibility string
		err        error
		url        string
		head       func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error)
		put        func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error)
	}{
		{
			name:    "Does not exist",
			objName: "foo.jpg",
			url:     "",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				return nil, mockError{code: "NotFound"}
			},
			put: func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(input.ACL).To(BeNil())
				g.Expect(*input.Metadata[mdVisibility]).To(Equal(visibilityPrivate))
				g.Expect(input.Metadata[mdCreated]).To(Not(BeNil()))
				return &sthree.PutObjectOutput{}, nil
			},
		},
		{
			name:       "Does not exist. Public object",
			objName:    "bar/baz/foo.jpg",
			visibility: "public",
			url:        "https://my-space.ams3.example.com/micro/123/bar/baz/foo.jpg",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				return nil, mockError{code: "NotFound"}
			},
			put: func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.ACL).To(Equal(mdACLPublic))
				g.Expect(*input.Metadata[mdVisibility]).To(Equal(visibilityPublic))
				g.Expect(input.Metadata[mdCreated]).To(Not(BeNil()))
				return &sthree.PutObjectOutput{}, nil
			},
		},
		{
			name:    "Missing name",
			objName: "",
			err:     errors.BadRequest("space.Update", "Missing name param"),
		},
		{
			name:    "Already exists",
			objName: "foo.jpg",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.Key).To(Equal("micro/123/foo.jpg"))
				return &sthree.HeadObjectOutput{
					Metadata: map[string]*string{
						mdCreated:    aws.String("1638541918"),
						mdVisibility: aws.String(visibilityPrivate),
					},
				}, nil
			},
			put: func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(input.ACL).To(BeNil())
				g.Expect(*input.Metadata[mdVisibility]).To(Equal(visibilityPrivate))
				// created shouuld be copied from the previous
				g.Expect(*input.Metadata[mdCreated]).To(Equal("1638541918"))
				return &sthree.PutObjectOutput{}, nil
			},
			url: "",
		},
		{
			name:    "Already exists public",
			objName: "foo.jpg",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.Key).To(Equal("micro/123/foo.jpg"))
				return &sthree.HeadObjectOutput{
					Metadata: map[string]*string{
						mdCreated:    aws.String("1638541918"),
						mdVisibility: aws.String(visibilityPrivate),
					},
				}, nil
			},
			put: func(input *sthree.PutObjectInput) (*sthree.PutObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.ACL).To(Equal(mdACLPublic))
				g.Expect(*input.Metadata[mdVisibility]).To(Equal(visibilityPublic))
				// created shouuld be copied from the previous
				g.Expect(*input.Metadata[mdCreated]).To(Equal("1638541918"))
				return &sthree.PutObjectOutput{}, nil
			},
			url:        "https://my-space.ams3.example.com/micro/123/foo.jpg",
			visibility: "public",
		},
	}
	store.DefaultStore = memory.NewStore()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			handler := Space{
				conf: conf{
					AccessKey: "access",
					SecretKey: "secret",
					Endpoint:  "example.com",
					SpaceName: "my-space",
					SSL:       true,
					Region:    "ams3",
					BaseURL:   "https://my-space.ams3.example.com",
				},
				client: &mockS3Client{head: tc.head, put: tc.put},
			}
			ctx := context.Background()
			ctx = auth.ContextWithAccount(ctx, &auth.Account{
				ID:       "123",
				Type:     "user",
				Issuer:   "micro",
				Metadata: map[string]string{},
				Scopes:   []string{"space"},
				Name:     "john@example.com",
			})

			rsp := pb.UpdateResponse{}
			err := handler.Update(ctx, &pb.UpdateRequest{
				Object:     []byte("foobar"),
				Name:       tc.objName,
				Visibility: tc.visibility,
			}, &rsp)
			if tc.err != nil {
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(rsp.Url).To(Equal(tc.url))
			}

		})
	}

}

func TestDelete(t *testing.T) {
	g := NewWithT(t)
	tcs := []struct {
		name    string
		objName string
		err     error
		delete  func(input *sthree.DeleteObjectInput) (*sthree.DeleteObjectOutput, error)
	}{
		{
			name:    "Simple case",
			objName: "foo.jpg",
		},
		{
			name:    "Missing name",
			objName: "",
			err:     errors.BadRequest("space.Delete", "Missing name param"),
		},
	}
	store.DefaultStore = memory.NewStore()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			handler := Space{
				conf: conf{
					AccessKey: "access",
					SecretKey: "secret",
					Endpoint:  "example.com",
					SpaceName: "my-space",
					SSL:       true,
					Region:    "ams3",
					BaseURL:   "https://my-space.ams3.example.com",
				},
				client: &mockS3Client{
					delete: func(input *sthree.DeleteObjectInput) (*sthree.DeleteObjectOutput, error) {
						g.Expect(input.Bucket).To(Equal(aws.String("my-space")))
						g.Expect(input.Key).To(Equal(aws.String("micro/123/" + tc.objName)))
						return &sthree.DeleteObjectOutput{}, nil
					}},
			}
			ctx := context.Background()
			ctx = auth.ContextWithAccount(ctx, &auth.Account{
				ID:       "123",
				Type:     "user",
				Issuer:   "micro",
				Metadata: map[string]string{},
				Scopes:   []string{"space"},
				Name:     "john@example.com",
			})
			rsp := pb.DeleteResponse{}
			err := handler.Delete(ctx, &pb.DeleteRequest{
				Name: tc.objName,
			}, &rsp)
			if tc.err != nil {
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(err).To(BeNil())
			}

		})
	}

}

func TestList(t *testing.T) {
	g := NewWithT(t)
	tcs := []struct {
		name       string
		prefix     string
		err        error
		list       func(input *sthree.ListObjectsInput) (*sthree.ListObjectsOutput, error)
		visibility string
	}{
		{
			name:   "Simple case",
			prefix: "file",
		},
		{
			name: "Empty prefix",
		},
	}
	store.DefaultStore = memory.NewStore()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			store.Write(
				store.NewRecord(fmt.Sprintf("%s/micro/123/file.jpg", prefixByUser),
					meta{
						Visibility:   "public",
						CreateTime:   "2009-11-10T23:00:00Z",
						ModifiedTime: "2009-11-10T23:00:00Z",
					}))
			store.Write(
				store.NewRecord(fmt.Sprintf("%s/micro/123/file2.jpg", prefixByUser),
					meta{
						Visibility:   "private",
						CreateTime:   "2009-11-10T23:00:01Z",
						ModifiedTime: "2009-11-10T23:00:01Z",
					}))
			handler := Space{
				conf: conf{
					AccessKey: "access",
					SecretKey: "secret",
					Endpoint:  "example.com",
					SpaceName: "my-space",
					SSL:       true,
					Region:    "ams3",
					BaseURL:   "https://my-space.ams3.example.com",
				},
				client: &mockS3Client{
					list: func(input *sthree.ListObjectsInput) (*sthree.ListObjectsOutput, error) {
						g.Expect(input.Bucket).To(Equal(aws.String("my-space")))
						g.Expect(input.Prefix).To(Equal(aws.String("micro/123/" + tc.prefix)))
						return &sthree.ListObjectsOutput{
							Contents: []*sthree.Object{
								{
									Key:          aws.String("micro/123/file.jpg"),
									LastModified: aws.Time(time.Unix(1257894000, 0)),
								},
								{
									Key:          aws.String("micro/123/file2.jpg"),
									LastModified: aws.Time(time.Unix(1257894000, 0)),
								},
							},
						}, nil
					}},
			}
			ctx := context.Background()
			ctx = auth.ContextWithAccount(ctx, &auth.Account{
				ID:       "123",
				Type:     "user",
				Issuer:   "micro",
				Metadata: map[string]string{},
				Scopes:   []string{"space"},
				Name:     "john@example.com",
			})
			rsp := pb.ListResponse{}
			err := handler.List(ctx, &pb.ListRequest{
				Prefix: tc.prefix,
			}, &rsp)
			if tc.err != nil {
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(rsp.Objects).To(HaveLen(2))
				g.Expect(rsp.Objects[0].Name).To(Equal("file.jpg"))
				g.Expect(rsp.Objects[0].Url).To(Equal("https://my-space.ams3.example.com/micro/123/file.jpg"))
				g.Expect(rsp.Objects[1].Name).To(Equal("file2.jpg"))
				g.Expect(rsp.Objects[1].Url).To(Equal(""))
			}

		})
	}

}

func TestHead(t *testing.T) {
	g := NewWithT(t)
	tcs := []struct {
		name       string
		objectName string
		url        string
		visibility string
		modified   string
		created    string
		err        error
		head       func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error)
	}{
		{
			name:       "Simple case",
			objectName: "foo.jpg",
			visibility: "public",
			url:        "https://my-space.ams3.example.com/micro/123/foo.jpg",
			created:    "2009-11-10T23:00:00Z",
			modified:   "2009-11-10T23:00:00Z",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.Key).To(Equal("micro/123/foo.jpg"))

				return &sthree.HeadObjectOutput{
					LastModified: aws.Time(time.Unix(1257894000, 0)),
					Metadata: map[string]*string{
						mdCreated:    aws.String("1257894000"),
						mdVisibility: aws.String(visibilityPublic),
					},
				}, nil
			},
		},
		{
			name:       "Simple case private",
			objectName: "foo.jpg",
			visibility: "private",
			url:        "",
			created:    "2009-11-10T23:00:00Z",
			modified:   "2009-11-10T23:00:00Z",
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				g.Expect(*input.Bucket).To(Equal("my-space"))
				g.Expect(*input.Key).To(Equal("micro/123/foo.jpg"))

				return &sthree.HeadObjectOutput{
					LastModified: aws.Time(time.Unix(1257894000, 0)),
					Metadata: map[string]*string{
						mdCreated:    aws.String("2009-11-10T23:00:00Z"),
						mdVisibility: aws.String("private"),
					},
				}, nil
			},
		},
		{
			name: "Empty prefix",
			err:  errors.BadRequest("space.Head", "Missing name param"),
		},
		{
			name:       "Not found",
			objectName: "foo.jpg",
			err:        errors.BadRequest("space.Head", "Object not found"),
			head: func(input *sthree.HeadObjectInput) (*sthree.HeadObjectOutput, error) {
				return nil, mockError{code: "NotFound"}
			},
		},
	}
	store.DefaultStore = memory.NewStore()

	for _, tc := range tcs {
		t.Run(tc.name, func(t *testing.T) {
			handler := Space{
				conf: conf{
					AccessKey: "access",
					SecretKey: "secret",
					Endpoint:  "example.com",
					SpaceName: "my-space",
					SSL:       true,
					Region:    "ams3",
					BaseURL:   "https://my-space.ams3.example.com",
				},
				client: &mockS3Client{
					head: tc.head,
				},
			}
			ctx := context.Background()
			ctx = auth.ContextWithAccount(ctx, &auth.Account{
				ID:       "123",
				Type:     "user",
				Issuer:   "micro",
				Metadata: map[string]string{},
				Scopes:   []string{"space"},
				Name:     "john@example.com",
			})
			rsp := pb.HeadResponse{}
			err := handler.Head(ctx, &pb.HeadRequest{
				Name: tc.objectName,
			}, &rsp)
			if tc.err != nil {
				g.Expect(err).To(Equal(tc.err))
			} else {
				g.Expect(err).To(BeNil())
				g.Expect(rsp.Object.Name).To(Equal(tc.objectName))
				g.Expect(rsp.Object.Url).To(Equal(tc.url))
				g.Expect(rsp.Object.Visibility).To(Equal(tc.visibility))
				g.Expect(rsp.Object.Created).To(Equal(tc.created))
				g.Expect(rsp.Object.Modified).To(Equal(tc.modified))
			}

		})
	}

}
