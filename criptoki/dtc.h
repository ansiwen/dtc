/* Code generated by cmd/cgo; DO NOT EDIT. */

/* package dtcmaster/criptoki */


#line 1 "cgo-builtin-export-prolog"

#include <stddef.h> /* for ptrdiff_t below */

#ifndef GO_CGO_EXPORT_PROLOGUE_H
#define GO_CGO_EXPORT_PROLOGUE_H

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef struct { const char *p; ptrdiff_t n; } _GoString_;
#endif

#endif

/* Start of preamble from import "C" comments.  */


#line 3 "pkcs11.go"

#include "pkcs11go.h"

#line 1 "cgo-generated-wrapper"


/* End of preamble from import "C" comments.  */


/* Start of boilerplate cgo prologue.  */
#line 1 "cgo-gcc-export-header-prolog"

#ifndef GO_CGO_PROLOGUE_H
#define GO_CGO_PROLOGUE_H

typedef signed char GoInt8;
typedef unsigned char GoUint8;
typedef short GoInt16;
typedef unsigned short GoUint16;
typedef int GoInt32;
typedef unsigned int GoUint32;
typedef long long GoInt64;
typedef unsigned long long GoUint64;
typedef GoInt64 GoInt;
typedef GoUint64 GoUint;
typedef __SIZE_TYPE__ GoUintptr;
typedef float GoFloat32;
typedef double GoFloat64;
typedef float _Complex GoComplex64;
typedef double _Complex GoComplex128;

/*
  static assertion to make sure the file is being used on architecture
  at least with matching size of GoInt.
*/
typedef char _check_for_64_bit_pointer_matching_GoInt[sizeof(void*)==64/8 ? 1:-1];

#ifndef GO_CGO_GOSTRING_TYPEDEF
typedef _GoString_ GoString;
#endif
typedef void *GoMap;
typedef void *GoChan;
typedef struct { void *t; void *v; } GoInterface;
typedef struct { void *data; GoInt len; GoInt cap; } GoSlice;

#endif

/* End of boilerplate cgo prologue.  */

#ifdef __cplusplus
extern "C" {
#endif


extern CK_RV C_Initialize(CK_VOID_PTR p0);

extern CK_RV C_Finalize(CK_VOID_PTR p0);

extern CK_RV C_InitToken(CK_SLOT_ID p0, CK_UTF8CHAR_PTR p1, CK_ULONG p2, CK_UTF8CHAR_PTR p3);

extern CK_RV C_InitPIN(CK_SESSION_HANDLE p0, CK_UTF8CHAR_PTR p1, CK_ULONG p2);

extern CK_RV C_SetPIN(CK_SESSION_HANDLE p0, CK_UTF8CHAR_PTR p1, CK_ULONG p2, CK_UTF8CHAR_PTR p3, CK_ULONG p4);

extern CK_RV C_GetInfo(CK_INFO_PTR p0);

extern CK_RV C_GetFunctionList(CK_FUNCTION_LIST_PTR_PTR p0);

extern CK_RV C_GetSlotList(CK_BBOOL p0, CK_SLOT_ID_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_GetSlotInfo(CK_SLOT_ID p0, CK_SLOT_INFO_PTR p1);

extern CK_RV C_GetTokenInfo(CK_SLOT_ID p0, CK_TOKEN_INFO_PTR p1);

extern CK_RV C_OpenSession(CK_SLOT_ID p0, CK_FLAGS p1, CK_VOID_PTR p2, CK_NOTIFY p3, CK_SESSION_HANDLE_PTR p4);

extern CK_RV C_CloseSession(CK_SESSION_HANDLE p0);

extern CK_RV C_CloseAllSessions(CK_SLOT_ID p0);

extern CK_RV C_GetSessionInfo(CK_SESSION_HANDLE p0, CK_SESSION_INFO_PTR p1);

extern CK_RV C_Login(CK_SESSION_HANDLE p0, CK_USER_TYPE p1, CK_UTF8CHAR_PTR p2, CK_ULONG p3);

extern CK_RV C_Logout(CK_SESSION_HANDLE p0);

extern CK_RV C_CreateObject(CK_SESSION_HANDLE p0, CK_ATTRIBUTE_PTR p1, CK_ULONG p2, CK_OBJECT_HANDLE_PTR p3);

extern CK_RV C_DestroyObject(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE p1);

extern CK_RV C_FindObjectsInit(CK_SESSION_HANDLE p0, CK_ATTRIBUTE_PTR p1, CK_ULONG p2);

extern CK_RV C_FindObjects(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE_PTR p1, CK_ULONG p2, CK_ULONG_PTR p3);

extern CK_RV C_FindObjectsFinal(CK_SESSION_HANDLE p0);

extern CK_RV C_GetAttributeValue(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE p1, CK_ATTRIBUTE_PTR p2, CK_ULONG p3);

extern CK_RV C_GenerateKeyPair(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_ATTRIBUTE_PTR p2, CK_ULONG p3, CK_ATTRIBUTE_PTR p4, CK_ULONG p5, CK_OBJECT_HANDLE_PTR p6, CK_OBJECT_HANDLE_PTR p7);

extern CK_RV C_SignInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2);

extern CK_RV C_SignUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2);

extern CK_RV C_SignFinal(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_Sign(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_VerifyInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2);

extern CK_RV C_Verify(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG p4);

extern CK_RV C_VerifyUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2);

extern CK_RV C_VerifyFinal(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2);

extern CK_RV C_Digest(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_SeedRandom(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2);

extern CK_RV C_GenerateRandom(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2);

extern CK_RV C_GetMechanismList(CK_SLOT_ID p0, CK_MECHANISM_TYPE_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_GetMechanismInfo(CK_SLOT_ID p0, CK_MECHANISM_TYPE p1, CK_MECHANISM_INFO_PTR p2);

extern CK_RV C_GetOperationState(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_SetOperationState(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_OBJECT_HANDLE p3, CK_OBJECT_HANDLE p4);

extern CK_RV C_CopyObject(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE p1, CK_ATTRIBUTE_PTR p2, CK_ULONG p3, CK_OBJECT_HANDLE_PTR p4);

extern CK_RV C_GetObjectSize(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE p1, CK_ULONG_PTR p2);

extern CK_RV C_SetAttributeValue(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE p1, CK_ATTRIBUTE_PTR p2, CK_ULONG p3);

extern CK_RV C_EncryptInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2);

extern CK_RV C_Encrypt(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_EncryptUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_EncryptFinal(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_DecryptInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2);

extern CK_RV C_Decrypt(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_DecryptUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_DecryptFinal(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_DigestInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1);

extern CK_RV C_DigestUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2);

extern CK_RV C_DigestKey(CK_SESSION_HANDLE p0, CK_OBJECT_HANDLE p1);

extern CK_RV C_DigestFinal(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG_PTR p2);

extern CK_RV C_SignRecoverInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2);

extern CK_RV C_SignRecover(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_VerifyRecoverInit(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2);

extern CK_RV C_VerifyRecover(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_DigestEncryptUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_DecryptDigestUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_SignEncryptUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_DecryptVerifyUpdate(CK_SESSION_HANDLE p0, CK_BYTE_PTR p1, CK_ULONG p2, CK_BYTE_PTR p3, CK_ULONG_PTR p4);

extern CK_RV C_GenerateKey(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_ATTRIBUTE_PTR p2, CK_ULONG p3, CK_OBJECT_HANDLE_PTR p4);

extern CK_RV C_WrapKey(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2, CK_OBJECT_HANDLE p3, CK_BYTE_PTR p4, CK_ULONG_PTR p5);

extern CK_RV C_UnwrapKey(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2, CK_BYTE_PTR p3, CK_ULONG p4, CK_ATTRIBUTE_PTR p5, CK_ULONG p6, CK_OBJECT_HANDLE_PTR p7);

extern CK_RV C_DeriveKey(CK_SESSION_HANDLE p0, CK_MECHANISM_PTR p1, CK_OBJECT_HANDLE p2, CK_ATTRIBUTE_PTR p3, CK_ULONG p4, CK_OBJECT_HANDLE_PTR p5);

extern CK_RV C_GetFunctionStatus(CK_SESSION_HANDLE p0);

extern CK_RV C_CancelFunction(CK_SESSION_HANDLE p0);

extern CK_RV C_WaitForSlotEvent(CK_FLAGS p0, CK_SLOT_ID_PTR p1, CK_VOID_PTR p2);

#ifdef __cplusplus
}
#endif
