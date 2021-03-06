// Code generated by protoc-gen-micro. DO NOT EDIT.
// source: storage.proto

/*
Package proto is a generated protocol buffer package.

It is generated from these files:
	storage.proto

It has these top-level messages:
	File
	Response
*/
package proto

import proto1 "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	client "github.com/micro/go-micro/client"
	server "github.com/micro/go-micro/server"
	context "context"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto1.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto1.ProtoPackageIsVersion2 // please upgrade the proto package

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ client.Option
var _ server.Option

// Client API for Storage service

type StorageService interface {
	Save(ctx context.Context, in *File, opts ...client.CallOption) (*Response, error)
}

type storageService struct {
	c    client.Client
	name string
}

func NewStorageService(name string, c client.Client) StorageService {
	if c == nil {
		c = client.NewClient()
	}
	if len(name) == 0 {
		name = "proto"
	}
	return &storageService{
		c:    c,
		name: name,
	}
}

func (c *storageService) Save(ctx context.Context, in *File, opts ...client.CallOption) (*Response, error) {
	req := c.c.NewRequest(c.name, "Storage.Save", in)
	out := new(Response)
	err := c.c.Call(ctx, req, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Storage service

type StorageHandler interface {
	Save(context.Context, *File, *Response) error
}

func RegisterStorageHandler(s server.Server, hdlr StorageHandler, opts ...server.HandlerOption) error {
	type storage interface {
		Save(ctx context.Context, in *File, out *Response) error
	}
	type Storage struct {
		storage
	}
	h := &storageHandler{hdlr}
	return s.Handle(s.NewHandler(&Storage{h}, opts...))
}

type storageHandler struct {
	StorageHandler
}

func (h *storageHandler) Save(ctx context.Context, in *File, out *Response) error {
	return h.StorageHandler.Save(ctx, in, out)
}
