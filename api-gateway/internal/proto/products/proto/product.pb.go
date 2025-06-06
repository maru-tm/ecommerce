// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.24.4
// source: internal/proto/products/product.proto

package proto

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Category struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Category) Reset() {
	*x = Category{}
	mi := &file_internal_proto_products_product_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Category) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Category) ProtoMessage() {}

func (x *Category) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_products_product_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Category.ProtoReflect.Descriptor instead.
func (*Category) Descriptor() ([]byte, []int) {
	return file_internal_proto_products_product_proto_rawDescGZIP(), []int{0}
}

func (x *Category) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Category) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

type Product struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Name          string                 `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Category      *Category              `protobuf:"bytes,3,opt,name=category,proto3" json:"category,omitempty"`
	Price         float64                `protobuf:"fixed64,4,opt,name=price,proto3" json:"price,omitempty"`
	Stock         int32                  `protobuf:"varint,5,opt,name=stock,proto3" json:"stock,omitempty"`
	Description   string                 `protobuf:"bytes,6,opt,name=description,proto3" json:"description,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Product) Reset() {
	*x = Product{}
	mi := &file_internal_proto_products_product_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Product) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Product) ProtoMessage() {}

func (x *Product) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_products_product_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Product.ProtoReflect.Descriptor instead.
func (*Product) Descriptor() ([]byte, []int) {
	return file_internal_proto_products_product_proto_rawDescGZIP(), []int{1}
}

func (x *Product) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *Product) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Product) GetCategory() *Category {
	if x != nil {
		return x.Category
	}
	return nil
}

func (x *Product) GetPrice() float64 {
	if x != nil {
		return x.Price
	}
	return 0
}

func (x *Product) GetStock() int32 {
	if x != nil {
		return x.Stock
	}
	return 0
}

func (x *Product) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

type ProductId struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Id            string                 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductId) Reset() {
	*x = ProductId{}
	mi := &file_internal_proto_products_product_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductId) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductId) ProtoMessage() {}

func (x *ProductId) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_products_product_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductId.ProtoReflect.Descriptor instead.
func (*ProductId) Descriptor() ([]byte, []int) {
	return file_internal_proto_products_product_proto_rawDescGZIP(), []int{2}
}

func (x *ProductId) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

type ProductList struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Products      []*Product             `protobuf:"bytes,1,rep,name=products,proto3" json:"products,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ProductList) Reset() {
	*x = ProductList{}
	mi := &file_internal_proto_products_product_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ProductList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProductList) ProtoMessage() {}

func (x *ProductList) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_products_product_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProductList.ProtoReflect.Descriptor instead.
func (*ProductList) Descriptor() ([]byte, []int) {
	return file_internal_proto_products_product_proto_rawDescGZIP(), []int{3}
}

func (x *ProductList) GetProducts() []*Product {
	if x != nil {
		return x.Products
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Empty) Reset() {
	*x = Empty{}
	mi := &file_internal_proto_products_product_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_internal_proto_products_product_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_internal_proto_products_product_proto_rawDescGZIP(), []int{4}
}

var File_internal_proto_products_product_proto protoreflect.FileDescriptor

const file_internal_proto_products_product_proto_rawDesc = "" +
	"\n" +
	"%internal/proto/products/product.proto\x12\tinventory\".\n" +
	"\bCategory\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\"\xac\x01\n" +
	"\aProduct\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\x12\x12\n" +
	"\x04name\x18\x02 \x01(\tR\x04name\x12/\n" +
	"\bcategory\x18\x03 \x01(\v2\x13.inventory.CategoryR\bcategory\x12\x14\n" +
	"\x05price\x18\x04 \x01(\x01R\x05price\x12\x14\n" +
	"\x05stock\x18\x05 \x01(\x05R\x05stock\x12 \n" +
	"\vdescription\x18\x06 \x01(\tR\vdescription\"\x1b\n" +
	"\tProductId\x12\x0e\n" +
	"\x02id\x18\x01 \x01(\tR\x02id\"=\n" +
	"\vProductList\x12.\n" +
	"\bproducts\x18\x01 \x03(\v2\x12.inventory.ProductR\bproducts\"\a\n" +
	"\x05Empty2\xb1\x02\n" +
	"\x0eProductService\x127\n" +
	"\rCreateProduct\x12\x12.inventory.Product\x1a\x12.inventory.Product\x12:\n" +
	"\x0eGetProductByID\x12\x14.inventory.ProductId\x1a\x12.inventory.Product\x128\n" +
	"\fListProducts\x12\x10.inventory.Empty\x1a\x16.inventory.ProductList\x127\n" +
	"\rUpdateProduct\x12\x12.inventory.Product\x1a\x12.inventory.Product\x127\n" +
	"\rDeleteProduct\x12\x14.inventory.ProductId\x1a\x10.inventory.EmptyB\x16Z\x14internal/proto;protob\x06proto3"

var (
	file_internal_proto_products_product_proto_rawDescOnce sync.Once
	file_internal_proto_products_product_proto_rawDescData []byte
)

func file_internal_proto_products_product_proto_rawDescGZIP() []byte {
	file_internal_proto_products_product_proto_rawDescOnce.Do(func() {
		file_internal_proto_products_product_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_internal_proto_products_product_proto_rawDesc), len(file_internal_proto_products_product_proto_rawDesc)))
	})
	return file_internal_proto_products_product_proto_rawDescData
}

var file_internal_proto_products_product_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_proto_products_product_proto_goTypes = []any{
	(*Category)(nil),    // 0: inventory.Category
	(*Product)(nil),     // 1: inventory.Product
	(*ProductId)(nil),   // 2: inventory.ProductId
	(*ProductList)(nil), // 3: inventory.ProductList
	(*Empty)(nil),       // 4: inventory.Empty
}
var file_internal_proto_products_product_proto_depIdxs = []int32{
	0, // 0: inventory.Product.category:type_name -> inventory.Category
	1, // 1: inventory.ProductList.products:type_name -> inventory.Product
	1, // 2: inventory.ProductService.CreateProduct:input_type -> inventory.Product
	2, // 3: inventory.ProductService.GetProductByID:input_type -> inventory.ProductId
	4, // 4: inventory.ProductService.ListProducts:input_type -> inventory.Empty
	1, // 5: inventory.ProductService.UpdateProduct:input_type -> inventory.Product
	2, // 6: inventory.ProductService.DeleteProduct:input_type -> inventory.ProductId
	1, // 7: inventory.ProductService.CreateProduct:output_type -> inventory.Product
	1, // 8: inventory.ProductService.GetProductByID:output_type -> inventory.Product
	3, // 9: inventory.ProductService.ListProducts:output_type -> inventory.ProductList
	1, // 10: inventory.ProductService.UpdateProduct:output_type -> inventory.Product
	4, // 11: inventory.ProductService.DeleteProduct:output_type -> inventory.Empty
	7, // [7:12] is the sub-list for method output_type
	2, // [2:7] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_internal_proto_products_product_proto_init() }
func file_internal_proto_products_product_proto_init() {
	if File_internal_proto_products_product_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_internal_proto_products_product_proto_rawDesc), len(file_internal_proto_products_product_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_proto_products_product_proto_goTypes,
		DependencyIndexes: file_internal_proto_products_product_proto_depIdxs,
		MessageInfos:      file_internal_proto_products_product_proto_msgTypes,
	}.Build()
	File_internal_proto_products_product_proto = out.File
	file_internal_proto_products_product_proto_goTypes = nil
	file_internal_proto_products_product_proto_depIdxs = nil
}
