package client

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"gopkg.in/resty.v1"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	iamsign "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type IamResty struct {
	configs []*aws.Config
	client  *resty.Client
}

func NewRestClient(configs ...*aws.Config) *IamResty {
	c := &IamResty{}
	c.configs = configs
	c.client = resty.New()
	c.client.SetPreRequestHook(c.iamAuthHook)
	return c
}

func (c *IamResty) Client() *resty.Client {
	return c.client
}

func (c *IamResty) iamAuthHook(_ *resty.Client, r *resty.Request) error {
	return c.iamAuth(r, time.Now())
}

func (c *IamResty) iamAuth(r *resty.Request, signTime time.Time) error {
	var body io.ReadSeeker
	if r.Body != nil {
		bodyBytes, err := ioutil.ReadAll(r.RawRequest.Body)
		if err != nil {
			return err
		}
		body = bytes.NewReader(bodyBytes)
	} else {
		body = bytes.NewReader([]byte{})
	}

	sess := session.New(c.configs...)
	region := *sess.Config.Region
	service := "execute-api"
	signer := iamsign.NewSigner(sess.Config.Credentials)
	_, err := signer.Sign(r.RawRequest, body, service, region, signTime)
	return err
}
