// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v5.26.1
// source: pay.proto

package payments

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Helper type for 128-bit integers, as protobuf does not natively support them.
type Uint128 struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	High uint64 `protobuf:"varint,1,opt,name=high,proto3" json:"high,omitempty"` // Higher 64 bits
	Low  uint64 `protobuf:"varint,2,opt,name=low,proto3" json:"low,omitempty"`   // Lower 64 bits
}

func (x *Uint128) Reset() {
	*x = Uint128{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Uint128) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Uint128) ProtoMessage() {}

func (x *Uint128) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Uint128.ProtoReflect.Descriptor instead.
func (*Uint128) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{0}
}

func (x *Uint128) GetHigh() uint64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *Uint128) GetLow() uint64 {
	if x != nil {
		return x.Low
	}
	return 0
}

// Updated Account message to include all fields required by TigerBeetle
type AccountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id             *Uint128 `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	DebitsPending  *Uint128 `protobuf:"bytes,2,opt,name=debitsPending,proto3" json:"debitsPending,omitempty"`
	DebitsPosted   *Uint128 `protobuf:"bytes,3,opt,name=debitsPosted,proto3" json:"debitsPosted,omitempty"`
	CreditsPending *Uint128 `protobuf:"bytes,4,opt,name=creditsPending,proto3" json:"creditsPending,omitempty"`
	CreditsPosted  *Uint128 `protobuf:"bytes,5,opt,name=creditsPosted,proto3" json:"creditsPosted,omitempty"`
	UserData128    *Uint128 `protobuf:"bytes,6,opt,name=userData128,proto3" json:"userData128,omitempty"`
	UserData64     uint64   `protobuf:"varint,7,opt,name=userData64,proto3" json:"userData64,omitempty"`
	UserData32     uint32   `protobuf:"varint,8,opt,name=userData32,proto3" json:"userData32,omitempty"`
	Reserved       uint32   `protobuf:"varint,9,opt,name=reserved,proto3" json:"reserved,omitempty"`
	Ledger         uint32   `protobuf:"varint,10,opt,name=ledger,proto3" json:"ledger,omitempty"`
	Code           uint32   `protobuf:"varint,11,opt,name=code,proto3" json:"code,omitempty"`
	Flags          uint32   `protobuf:"varint,12,opt,name=flags,proto3" json:"flags,omitempty"`
	Timestamp      uint64   `protobuf:"varint,13,opt,name=timestamp,proto3" json:"timestamp,omitempty"`
}

func (x *AccountResponse) Reset() {
	*x = AccountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AccountResponse) ProtoMessage() {}

func (x *AccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AccountResponse.ProtoReflect.Descriptor instead.
func (*AccountResponse) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{1}
}

func (x *AccountResponse) GetId() *Uint128 {
	if x != nil {
		return x.Id
	}
	return nil
}

func (x *AccountResponse) GetDebitsPending() *Uint128 {
	if x != nil {
		return x.DebitsPending
	}
	return nil
}

func (x *AccountResponse) GetDebitsPosted() *Uint128 {
	if x != nil {
		return x.DebitsPosted
	}
	return nil
}

func (x *AccountResponse) GetCreditsPending() *Uint128 {
	if x != nil {
		return x.CreditsPending
	}
	return nil
}

func (x *AccountResponse) GetCreditsPosted() *Uint128 {
	if x != nil {
		return x.CreditsPosted
	}
	return nil
}

func (x *AccountResponse) GetUserData128() *Uint128 {
	if x != nil {
		return x.UserData128
	}
	return nil
}

func (x *AccountResponse) GetUserData64() uint64 {
	if x != nil {
		return x.UserData64
	}
	return 0
}

func (x *AccountResponse) GetUserData32() uint32 {
	if x != nil {
		return x.UserData32
	}
	return 0
}

func (x *AccountResponse) GetReserved() uint32 {
	if x != nil {
		return x.Reserved
	}
	return 0
}

func (x *AccountResponse) GetLedger() uint32 {
	if x != nil {
		return x.Ledger
	}
	return 0
}

func (x *AccountResponse) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *AccountResponse) GetFlags() uint32 {
	if x != nil {
		return x.Flags
	}
	return 0
}

func (x *AccountResponse) GetTimestamp() uint64 {
	if x != nil {
		return x.Timestamp
	}
	return 0
}

type Account struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id     uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Ledger uint32 `protobuf:"varint,2,opt,name=ledger,proto3" json:"ledger,omitempty"`
	Code   uint32 `protobuf:"varint,3,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *Account) Reset() {
	*x = Account{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Account) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Account) ProtoMessage() {}

func (x *Account) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Account.ProtoReflect.Descriptor instead.
func (*Account) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{2}
}

func (x *Account) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Account) GetLedger() uint32 {
	if x != nil {
		return x.Ledger
	}
	return 0
}

func (x *Account) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

// Define messages for transfer operations
type Transfer struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id              uint64 `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DebitAccountId  uint64 `protobuf:"varint,2,opt,name=debitAccountId,proto3" json:"debitAccountId,omitempty"`
	CreditAccountId uint64 `protobuf:"varint,3,opt,name=creditAccountId,proto3" json:"creditAccountId,omitempty"`
	Amount          uint64 `protobuf:"varint,4,opt,name=amount,proto3" json:"amount,omitempty"`
	Ledger          uint32 `protobuf:"varint,5,opt,name=ledger,proto3" json:"ledger,omitempty"`
	Code            uint32 `protobuf:"varint,6,opt,name=code,proto3" json:"code,omitempty"`
}

func (x *Transfer) Reset() {
	*x = Transfer{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Transfer) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Transfer) ProtoMessage() {}

func (x *Transfer) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Transfer.ProtoReflect.Descriptor instead.
func (*Transfer) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{3}
}

func (x *Transfer) GetId() uint64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Transfer) GetDebitAccountId() uint64 {
	if x != nil {
		return x.DebitAccountId
	}
	return 0
}

func (x *Transfer) GetCreditAccountId() uint64 {
	if x != nil {
		return x.CreditAccountId
	}
	return 0
}

func (x *Transfer) GetAmount() uint64 {
	if x != nil {
		return x.Amount
	}
	return 0
}

func (x *Transfer) GetLedger() uint32 {
	if x != nil {
		return x.Ledger
	}
	return 0
}

func (x *Transfer) GetCode() uint32 {
	if x != nil {
		return x.Code
	}
	return 0
}

type CreateAccountRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Account *Account `protobuf:"bytes,1,opt,name=account,proto3" json:"account,omitempty"`
}

func (x *CreateAccountRequest) Reset() {
	*x = CreateAccountRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAccountRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountRequest) ProtoMessage() {}

func (x *CreateAccountRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountRequest.ProtoReflect.Descriptor instead.
func (*CreateAccountRequest) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{4}
}

func (x *CreateAccountRequest) GetAccount() *Account {
	if x != nil {
		return x.Account
	}
	return nil
}

type CreateAccountResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results string `protobuf:"bytes,1,opt,name=results,proto3" json:"results,omitempty"`
}

func (x *CreateAccountResponse) Reset() {
	*x = CreateAccountResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAccountResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountResponse) ProtoMessage() {}

func (x *CreateAccountResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountResponse.ProtoReflect.Descriptor instead.
func (*CreateAccountResponse) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{5}
}

func (x *CreateAccountResponse) GetResults() string {
	if x != nil {
		return x.Results
	}
	return ""
}

type CreateAccountBatchRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accounts []*Account `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts,omitempty"`
}

func (x *CreateAccountBatchRequest) Reset() {
	*x = CreateAccountBatchRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAccountBatchRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountBatchRequest) ProtoMessage() {}

func (x *CreateAccountBatchRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountBatchRequest.ProtoReflect.Descriptor instead.
func (*CreateAccountBatchRequest) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{6}
}

func (x *CreateAccountBatchRequest) GetAccounts() []*Account {
	if x != nil {
		return x.Accounts
	}
	return nil
}

type CreateAccountBatchResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []string `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *CreateAccountBatchResponse) Reset() {
	*x = CreateAccountBatchResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateAccountBatchResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateAccountBatchResponse) ProtoMessage() {}

func (x *CreateAccountBatchResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateAccountBatchResponse.ProtoReflect.Descriptor instead.
func (*CreateAccountBatchResponse) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{7}
}

func (x *CreateAccountBatchResponse) GetResults() []string {
	if x != nil {
		return x.Results
	}
	return nil
}

type LookupAccountsRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AccountIds []uint64 `protobuf:"varint,1,rep,packed,name=accountIds,proto3" json:"accountIds,omitempty"`
}

func (x *LookupAccountsRequest) Reset() {
	*x = LookupAccountsRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupAccountsRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupAccountsRequest) ProtoMessage() {}

func (x *LookupAccountsRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupAccountsRequest.ProtoReflect.Descriptor instead.
func (*LookupAccountsRequest) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{8}
}

func (x *LookupAccountsRequest) GetAccountIds() []uint64 {
	if x != nil {
		return x.AccountIds
	}
	return nil
}

type LookupAccountsResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Accounts []*AccountResponse `protobuf:"bytes,1,rep,name=accounts,proto3" json:"accounts,omitempty"`
}

func (x *LookupAccountsResponse) Reset() {
	*x = LookupAccountsResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LookupAccountsResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LookupAccountsResponse) ProtoMessage() {}

func (x *LookupAccountsResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LookupAccountsResponse.ProtoReflect.Descriptor instead.
func (*LookupAccountsResponse) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{9}
}

func (x *LookupAccountsResponse) GetAccounts() []*AccountResponse {
	if x != nil {
		return x.Accounts
	}
	return nil
}

type CreateTransfersRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Transfers []*Transfer `protobuf:"bytes,1,rep,name=transfers,proto3" json:"transfers,omitempty"`
}

func (x *CreateTransfersRequest) Reset() {
	*x = CreateTransfersRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransfersRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransfersRequest) ProtoMessage() {}

func (x *CreateTransfersRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransfersRequest.ProtoReflect.Descriptor instead.
func (*CreateTransfersRequest) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{10}
}

func (x *CreateTransfersRequest) GetTransfers() []*Transfer {
	if x != nil {
		return x.Transfers
	}
	return nil
}

type CreateTransfersResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Results []string `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
}

func (x *CreateTransfersResponse) Reset() {
	*x = CreateTransfersResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pay_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreateTransfersResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreateTransfersResponse) ProtoMessage() {}

func (x *CreateTransfersResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pay_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreateTransfersResponse.ProtoReflect.Descriptor instead.
func (*CreateTransfersResponse) Descriptor() ([]byte, []int) {
	return file_pay_proto_rawDescGZIP(), []int{11}
}

func (x *CreateTransfersResponse) GetResults() []string {
	if x != nil {
		return x.Results
	}
	return nil
}

var File_pay_proto protoreflect.FileDescriptor

var file_pay_proto_rawDesc = []byte{
	0x0a, 0x09, 0x70, 0x61, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x22, 0x2f, 0x0a, 0x07, 0x55, 0x69, 0x6e, 0x74, 0x31, 0x32, 0x38, 0x12,
	0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x04, 0x68,
	0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x03, 0x6c, 0x6f, 0x77, 0x22, 0x83, 0x04, 0x0a, 0x0f, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x20, 0x0a, 0x02, 0x69, 0x64, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x55, 0x69, 0x6e, 0x74, 0x31, 0x32, 0x38, 0x52, 0x02, 0x69, 0x64, 0x12, 0x36, 0x0a, 0x0d, 0x64,
	0x65, 0x62, 0x69, 0x74, 0x73, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x69, 0x6e,
	0x74, 0x31, 0x32, 0x38, 0x52, 0x0d, 0x64, 0x65, 0x62, 0x69, 0x74, 0x73, 0x50, 0x65, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x12, 0x34, 0x0a, 0x0c, 0x64, 0x65, 0x62, 0x69, 0x74, 0x73, 0x50, 0x6f, 0x73,
	0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x55, 0x69, 0x6e, 0x74, 0x31, 0x32, 0x38, 0x52, 0x0c, 0x64, 0x65, 0x62,
	0x69, 0x74, 0x73, 0x50, 0x6f, 0x73, 0x74, 0x65, 0x64, 0x12, 0x38, 0x0a, 0x0e, 0x63, 0x72, 0x65,
	0x64, 0x69, 0x74, 0x73, 0x50, 0x65, 0x6e, 0x64, 0x69, 0x6e, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x69, 0x6e, 0x74,
	0x31, 0x32, 0x38, 0x52, 0x0e, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x73, 0x50, 0x65, 0x6e, 0x64,
	0x69, 0x6e, 0x67, 0x12, 0x36, 0x0a, 0x0d, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x73, 0x50, 0x6f,
	0x73, 0x74, 0x65, 0x64, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61, 0x79,
	0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x69, 0x6e, 0x74, 0x31, 0x32, 0x38, 0x52, 0x0d, 0x63, 0x72,
	0x65, 0x64, 0x69, 0x74, 0x73, 0x50, 0x6f, 0x73, 0x74, 0x65, 0x64, 0x12, 0x32, 0x0a, 0x0b, 0x75,
	0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x31, 0x32, 0x38, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x10, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x55, 0x69, 0x6e, 0x74, 0x31,
	0x32, 0x38, 0x52, 0x0b, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x31, 0x32, 0x38, 0x12,
	0x1e, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x36, 0x34, 0x18, 0x07, 0x20,
	0x01, 0x28, 0x04, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x36, 0x34, 0x12,
	0x1e, 0x0a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x33, 0x32, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x0d, 0x52, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x44, 0x61, 0x74, 0x61, 0x33, 0x32, 0x12,
	0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x18, 0x09, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x08, 0x72, 0x65, 0x73, 0x65, 0x72, 0x76, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c,
	0x65, 0x64, 0x67, 0x65, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x0b, 0x20, 0x01, 0x28,
	0x0d, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x05, 0x66, 0x6c, 0x61, 0x67, 0x73, 0x12, 0x1c, 0x0a,
	0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x09, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x22, 0x45, 0x0a, 0x07, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x65, 0x64, 0x67, 0x65, 0x72, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x22, 0xb0, 0x01, 0x0a, 0x08, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x12,
	0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x04, 0x52, 0x02, 0x69, 0x64, 0x12,
	0x26, 0x0a, 0x0e, 0x64, 0x65, 0x62, 0x69, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x04, 0x52, 0x0e, 0x64, 0x65, 0x62, 0x69, 0x74, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x28, 0x0a, 0x0f, 0x63, 0x72, 0x65, 0x64, 0x69,
	0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x04,
	0x52, 0x0f, 0x63, 0x72, 0x65, 0x64, 0x69, 0x74, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49,
	0x64, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x04, 0x52, 0x06, 0x61, 0x6d, 0x6f, 0x75, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x65, 0x64,
	0x67, 0x65, 0x72, 0x18, 0x05, 0x20, 0x01, 0x28, 0x0d, 0x52, 0x06, 0x6c, 0x65, 0x64, 0x67, 0x65,
	0x72, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x0d, 0x52,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x22, 0x42, 0x0a, 0x14, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2a, 0x0a,
	0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10,
	0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0x31, 0x0a, 0x15, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22, 0x49, 0x0a, 0x19,
	0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x74,
	0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x2c, 0x0a, 0x08, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x61,
	0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x08, 0x61,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x22, 0x36, 0x0a, 0x1a, 0x43, 0x72, 0x65, 0x61, 0x74,
	0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x22,
	0x37, 0x0a, 0x15, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x73, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x1e, 0x0a, 0x0a, 0x61, 0x63, 0x63, 0x6f,
	0x75, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0a, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x64, 0x73, 0x22, 0x4e, 0x0a, 0x16, 0x4c, 0x6f, 0x6f, 0x6b,
	0x75, 0x70, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x12, 0x34, 0x0a, 0x08, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x18, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x41,
	0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x52, 0x08,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x22, 0x49, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x2f, 0x0a, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x52, 0x09, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x66,
	0x65, 0x72, 0x73, 0x22, 0x33, 0x0a, 0x17, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61,
	0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18,
	0x0a, 0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x07, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x73, 0x32, 0xc9, 0x01, 0x0a, 0x14, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x50, 0x0a, 0x0d, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x12, 0x1d, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x1a, 0x1e, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x22, 0x00, 0x12, 0x5f, 0x0a, 0x12, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x12, 0x22, 0x2e, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e,
	0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x23, 0x2e,
	0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x42, 0x61, 0x74, 0x63, 0x68, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e,
	0x73, 0x65, 0x22, 0x00, 0x32, 0x70, 0x0a, 0x19, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x73, 0x4c, 0x6f, 0x6f, 0x6b, 0x55, 0x70, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63,
	0x65, 0x12, 0x53, 0x0a, 0x0e, 0x4c, 0x6f, 0x6f, 0x6b, 0x75, 0x70, 0x41, 0x63, 0x63, 0x6f, 0x75,
	0x6e, 0x74, 0x73, 0x12, 0x1e, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x1f, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x4c, 0x6f,
	0x6f, 0x6b, 0x75, 0x70, 0x41, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x73, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x6f, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12,
	0x56, 0x0a, 0x0f, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65,
	0x72, 0x73, 0x12, 0x1f, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72, 0x65,
	0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x20, 0x2e, 0x70, 0x61, 0x79, 0x6d, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x72,
	0x65, 0x61, 0x74, 0x65, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x66, 0x65, 0x72, 0x73, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x0b, 0x5a, 0x09, 0x2f, 0x70, 0x61, 0x79, 0x6d,
	0x65, 0x6e, 0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pay_proto_rawDescOnce sync.Once
	file_pay_proto_rawDescData = file_pay_proto_rawDesc
)

func file_pay_proto_rawDescGZIP() []byte {
	file_pay_proto_rawDescOnce.Do(func() {
		file_pay_proto_rawDescData = protoimpl.X.CompressGZIP(file_pay_proto_rawDescData)
	})
	return file_pay_proto_rawDescData
}

var file_pay_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_pay_proto_goTypes = []interface{}{
	(*Uint128)(nil),                    // 0: payment.Uint128
	(*AccountResponse)(nil),            // 1: payment.AccountResponse
	(*Account)(nil),                    // 2: payment.Account
	(*Transfer)(nil),                   // 3: payment.Transfer
	(*CreateAccountRequest)(nil),       // 4: payment.CreateAccountRequest
	(*CreateAccountResponse)(nil),      // 5: payment.CreateAccountResponse
	(*CreateAccountBatchRequest)(nil),  // 6: payment.CreateAccountBatchRequest
	(*CreateAccountBatchResponse)(nil), // 7: payment.CreateAccountBatchResponse
	(*LookupAccountsRequest)(nil),      // 8: payment.LookupAccountsRequest
	(*LookupAccountsResponse)(nil),     // 9: payment.LookupAccountsResponse
	(*CreateTransfersRequest)(nil),     // 10: payment.CreateTransfersRequest
	(*CreateTransfersResponse)(nil),    // 11: payment.CreateTransfersResponse
}
var file_pay_proto_depIdxs = []int32{
	0,  // 0: payment.AccountResponse.id:type_name -> payment.Uint128
	0,  // 1: payment.AccountResponse.debitsPending:type_name -> payment.Uint128
	0,  // 2: payment.AccountResponse.debitsPosted:type_name -> payment.Uint128
	0,  // 3: payment.AccountResponse.creditsPending:type_name -> payment.Uint128
	0,  // 4: payment.AccountResponse.creditsPosted:type_name -> payment.Uint128
	0,  // 5: payment.AccountResponse.userData128:type_name -> payment.Uint128
	2,  // 6: payment.CreateAccountRequest.account:type_name -> payment.Account
	2,  // 7: payment.CreateAccountBatchRequest.accounts:type_name -> payment.Account
	1,  // 8: payment.LookupAccountsResponse.accounts:type_name -> payment.AccountResponse
	3,  // 9: payment.CreateTransfersRequest.transfers:type_name -> payment.Transfer
	4,  // 10: payment.CreateAccountService.CreateAccount:input_type -> payment.CreateAccountRequest
	6,  // 11: payment.CreateAccountService.CreateAccountBatch:input_type -> payment.CreateAccountBatchRequest
	8,  // 12: payment.TransactionsLookUpService.LookupAccounts:input_type -> payment.LookupAccountsRequest
	10, // 13: payment.CreateTransferService.CreateTransfers:input_type -> payment.CreateTransfersRequest
	5,  // 14: payment.CreateAccountService.CreateAccount:output_type -> payment.CreateAccountResponse
	7,  // 15: payment.CreateAccountService.CreateAccountBatch:output_type -> payment.CreateAccountBatchResponse
	9,  // 16: payment.TransactionsLookUpService.LookupAccounts:output_type -> payment.LookupAccountsResponse
	11, // 17: payment.CreateTransferService.CreateTransfers:output_type -> payment.CreateTransfersResponse
	14, // [14:18] is the sub-list for method output_type
	10, // [10:14] is the sub-list for method input_type
	10, // [10:10] is the sub-list for extension type_name
	10, // [10:10] is the sub-list for extension extendee
	0,  // [0:10] is the sub-list for field type_name
}

func init() { file_pay_proto_init() }
func file_pay_proto_init() {
	if File_pay_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pay_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Uint128); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AccountResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Account); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Transfer); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAccountRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAccountResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAccountBatchRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateAccountBatchResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupAccountsRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LookupAccountsResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTransfersRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pay_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreateTransfersResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pay_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_pay_proto_goTypes,
		DependencyIndexes: file_pay_proto_depIdxs,
		MessageInfos:      file_pay_proto_msgTypes,
	}.Build()
	File_pay_proto = out.File
	file_pay_proto_rawDesc = nil
	file_pay_proto_goTypes = nil
	file_pay_proto_depIdxs = nil
}
