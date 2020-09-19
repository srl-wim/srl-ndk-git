/*
Package ndkygit is a generated package which contains definitions
of structs which represent a YANG schema. The generated schema can be
compressed by a series of transformations (compression was false
in this case).

This package was generated by /Users/henderiw/CodeProjects/go/pkg/mod/github.com/openconfig/ygot@v0.8.3/genutil/names.go
using the following YANG input files:
	- ../yang/ndk-git.yang
	- ../yang/srl_nokia-common.yang
Imported modules were sourced from:
	- ../yang/...
*/
package ndkygit

var (
	// ySchema is a byte slice contain a gzip compressed representation of the
	// YANG schema from which the Go code was generated. When uncompressed the
	// contents of the byte slice is a JSON document containing an object, keyed
	// on the name of the generated struct, and containing the JSON marshalled
	// contents of a goyang yang.Entry struct, which defines the schema for the
	// fields within the struct.
	ySchema = []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xec, 0x9b, 0x6b, 0x6f, 0xd3, 0x3c,
		0x14, 0xc7, 0xdf, 0xf7, 0x53, 0xe4, 0xf1, 0x03, 0xe2, 0xd6, 0xb0, 0x74, 0xb4, 0xdd, 0x56, 0x84,
		0x60, 0x62, 0x0c, 0x24, 0x2e, 0x9a, 0x36, 0xc4, 0x0b, 0x46, 0x19, 0x5e, 0xe6, 0xa6, 0xd6, 0x5a,
		0xbb, 0x72, 0x6c, 0xb1, 0x0d, 0xca, 0x67, 0x47, 0xb9, 0xf4, 0x92, 0xe6, 0x62, 0x27, 0xdd, 0x45,
		0x0a, 0xe7, 0x45, 0x25, 0x6a, 0xff, 0x1d, 0x1f, 0x9f, 0x73, 0xf2, 0x3b, 0xc9, 0x29, 0xfb, 0xd5,
		0xb0, 0x2c, 0xcb, 0x42, 0x9f, 0xf0, 0x98, 0xa0, 0x9e, 0x85, 0x50, 0x33, 0xfa, 0xfe, 0x9e, 0xb2,
		0x33, 0xd4, 0xb3, 0x9c, 0xf8, 0xeb, 0x6b, 0xce, 0x06, 0xd4, 0x5b, 0x1a, 0xd8, 0xa3, 0x02, 0xf5,
		0xac, 0x68, 0x71, 0x38, 0xe0, 0x51, 0x99, 0x18, 0x48, 0x5c, 0x35, 0x98, 0x6c, 0x26, 0xa7, 0xe2,
		0x0d, 0x5a, 0x2b, 0xc3, 0xab, 0x1b, 0xcd, 0x27, 0x0e, 0x04, 0x19, 0xd0, 0x8b, 0xd4, 0x16, 0x89,
		0x6d, 0x7c, 0x31, 0xb2, 0xd3, 0x5b, 0x85, 0x92, 0x23, 0xae, 0x84, 0x4b, 0x32, 0x97, 0x47, 0xe6,
		0x90, 0xcb, 0x9f, 0x5c, 0x04, 0x16, 0xa1, 0x49, 0xb4, 0x53, 0x33, 0x5b, 0xf8, 0x0e, 0xfb, 0xbb,
		0xc2, 0x53, 0x63, 0xc2, 0x82, 0xf3, 0x4a, 0xa1, 0x48, 0x8e, 0x70, 0x49, 0x35, 0x37, 0x2c, 0xa5,
		0x9c, 0x26, 0x46, 0xa6, 0x2b, 0x67, 0x5e, 0x75, 0xf2, 0x7c, 0x02, 0xbb, 0x92, 0x72, 0x96, 0x7f,
		0x9a, 0x99, 0x43, 0x62, 0x5d, 0x8e, 0x85, 0xc9, 0x18, 0xa7, 0xa6, 0xf3, 0x42, 0x61, 0x12, 0x92,
		0x12, 0xa1, 0x31, 0x0d, 0x51, 0xe9, 0x50, 0x95, 0x0e, 0x59, 0xb9, 0xd0, 0x65, 0x87, 0x30, 0x27,
		0x94, 0xf3, 0x4b, 0x7f, 0xbe, 0x9c, 0x10, 0x43, 0x8f, 0x49, 0x41, 0x99, 0x57, 0xe4, 0xb0, 0xd9,
		0x2d, 0xb4, 0x5d, 0xa0, 0xf9, 0x40, 0x98, 0x27, 0x87, 0xa8, 0x67, 0x1d, 0x17, 0x1e, 0xb9, 0xd8,
		0xe5, 0xe1, 0x95, 0x3e, 0x52, 0xa6, 0x8d, 0x8d, 0x61, 0x6a, 0xa5, 0xe4, 0x5f, 0xf0, 0x48, 0x91,
		0x12, 0xfa, 0x7d, 0x11, 0x65, 0xf6, 0x1e, 0xf5, 0xa8, 0xf4, 0x83, 0x85, 0xda, 0x75, 0xd3, 0xa6,
		0xc1, 0x11, 0xf1, 0xc5, 0x8d, 0x1f, 0x71, 0xb3, 0xd3, 0xb9, 0xc1, 0x43, 0x36, 0xaa, 0xcd, 0xf6,
		0x0b, 0x32, 0xe8, 0x00, 0x4b, 0x49, 0x04, 0xd3, 0xa6, 0x10, 0x3a, 0xde, 0xb5, 0xbf, 0x62, 0xfb,
		0xca, 0xb1, 0x77, 0xac, 0xff, 0x5e, 0xfd, 0x7f, 0xef, 0xfe, 0xf7, 0x6f, 0xca, 0x71, 0x36, 0xbb,
		0x0f, 0x1f, 0xfd, 0x7e, 0xf2, 0xe2, 0xc7, 0x9f, 0xa7, 0xcd, 0x07, 0x1b, 0x27, 0xbd, 0xe7, 0x2f,
		0xed, 0xfe, 0xe3, 0xfc, 0x7b, 0xa9, 0xdf, 0x30, 0xb3, 0x3b, 0x23, 0x98, 0x08, 0x2b, 0x39, 0xe4,
		0xc2, 0x80, 0x89, 0x91, 0x0e, 0x98, 0x08, 0x4c, 0x04, 0x26, 0x02, 0x13, 0xff, 0x01, 0x26, 0xda,
		0x64, 0x8c, 0xe9, 0xc8, 0x94, 0x8c, 0xb1, 0x1a, 0xf8, 0x08, 0x7c, 0x04, 0x3e, 0x02, 0x1f, 0x6b,
		0xcd, 0xc7, 0x53, 0x81, 0x99, 0x3b, 0xd4, 0x93, 0x31, 0xd6, 0x01, 0x13, 0x81, 0x89, 0xc0, 0x44,
		0x60, 0x62, 0xad, 0x99, 0x38, 0xa0, 0x23, 0xa2, 0x27, 0x62, 0xa8, 0x02, 0x1e, 0x02, 0x0f, 0x81,
		0x87, 0xc0, 0xc3, 0x5a, 0xf3, 0x90, 0x4f, 0x88, 0xb0, 0x7d, 0x89, 0xa5, 0x01, 0x15, 0x97, 0xb4,
		0x6b, 0xb2, 0x71, 0x13, 0xd8, 0x78, 0x07, 0x6c, 0xd4, 0xc6, 0x2f, 0xc1, 0xc7, 0x76, 0x81, 0xe6,
		0x0d, 0x53, 0xe3, 0x60, 0xcf, 0xe9, 0x3a, 0x99, 0x27, 0x3c, 0xcc, 0xe8, 0x15, 0x36, 0xfb, 0xad,
		0x2f, 0xa1, 0x86, 0xca, 0x0c, 0x95, 0x19, 0x2a, 0x33, 0x54, 0xe6, 0x5a, 0x57, 0x66, 0x41, 0x26,
		0x5c, 0xcf, 0xc5, 0x50, 0x05, 0x3c, 0x04, 0x1e, 0x02, 0x0f, 0x81, 0x87, 0xb5, 0xe6, 0x61, 0xf0,
		0xe0, 0x4a, 0x7d, 0x49, 0x5d, 0x5f, 0x4f, 0xc5, 0x25, 0x6d, 0x31, 0x1b, 0x5b, 0xf0, 0xa6, 0x72,
		0x7b, 0x6c, 0xcc, 0xfb, 0x0f, 0x7f, 0x8b, 0xe6, 0x1c, 0xa6, 0x23, 0x25, 0x0c, 0xdc, 0x30, 0xef,
		0xd3, 0xc5, 0x0b, 0x34, 0x67, 0x32, 0xbb, 0x95, 0xb5, 0x05, 0xb1, 0x4c, 0xf0, 0x2b, 0x24, 0x41,
		0xd9, 0x64, 0xa8, 0x9c, 0x14, 0x95, 0x93, 0xa3, 0x5a, 0x92, 0x18, 0xc2, 0x4b, 0xe3, 0x73, 0x6d,
		0x61, 0x4d, 0x79, 0xfc, 0x8a, 0x08, 0x6e, 0x9f, 0x62, 0x9f, 0x9c, 0xd9, 0x2e, 0x57, 0x4c, 0x12,
		0xd1, 0x6d, 0x9b, 0xb8, 0x3f, 0xce, 0x96, 0x6d, 0x03, 0xe9, 0x1e, 0x19, 0x60, 0x35, 0x0a, 0xdd,
		0xe0, 0x98, 0x5c, 0xfa, 0x10, 0x33, 0x8f, 0x68, 0x49, 0x6b, 0x5e, 0xb4, 0x2b, 0x15, 0xef, 0x8a,
		0x15, 0xae, 0x6a, 0x31, 0x5f, 0xa7, 0xde, 0x95, 0x28, 0xee, 0x95, 0x8a, 0xfc, 0x75, 0xb9, 0xa2,
		0xb5, 0xdd, 0x6e, 0x77, 0xb7, 0xda, 0x6d, 0x67, 0xeb, 0xd9, 0x96, 0xb3, 0xd3, 0xe9, 0xb4, 0xba,
		0xad, 0xce, 0x2d, 0x7a, 0xa7, 0x71, 0x3d, 0xaa, 0x7e, 0xc5, 0xa7, 0x8b, 0x82, 0xe8, 0x20, 0x5f,
		0xb9, 0x2e, 0xf1, 0x7d, 0x73, 0xa4, 0xcf, 0x16, 0x00, 0xd2, 0x01, 0xe9, 0x80, 0x74, 0x40, 0x3a,
		0x20, 0xfd, 0x4e, 0x90, 0x5e, 0xea, 0xc1, 0x7e, 0x97, 0x31, 0x2e, 0x8b, 0x1b, 0xf9, 0x51, 0x31,
		0x70, 0x87, 0x64, 0x8c, 0x27, 0x38, 0x6c, 0x58, 0xa0, 0x0d, 0x8f, 0xca, 0xf0, 0xa3, 0x7d, 0x5b,
		0x8b, 0xdf, 0xff, 0x84, 0x72, 0x25, 0x8b, 0x49, 0xf0, 0x96, 0xca, 0x93, 0xe0, 0x73, 0xb4, 0x58,
		0xbb, 0xc6, 0xbb, 0xa5, 0xe4, 0xe7, 0xc4, 0xe0, 0x47, 0x88, 0x48, 0x06, 0xdd, 0x36, 0xe8, 0xb6,
		0x41, 0xb7, 0x0d, 0xba, 0x6d, 0xf5, 0xe9, 0xb6, 0x15, 0xfe, 0x99, 0xa6, 0x06, 0xee, 0x39, 0x50,
		0x47, 0x99, 0x1d, 0xbc, 0x0c, 0x82, 0x27, 0xcd, 0x5f, 0x98, 0x12, 0xfd, 0x2b, 0x36, 0x26, 0xcf,
		0x08, 0x44, 0xfd, 0x7d, 0x7c, 0x4e, 0x0e, 0x39, 0x9f, 0xb1, 0x28, 0x5a, 0xd5, 0x98, 0xfe, 0x05,
		0x00, 0x00, 0xff, 0xff, 0x01, 0x00, 0x00, 0xff, 0xff, 0xd8, 0x62, 0xf3, 0x1e, 0xf9, 0x3b, 0x00,
		0x00,
	}
)