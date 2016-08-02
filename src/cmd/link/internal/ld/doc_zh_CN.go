// Copyright 2015 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

package ld // import "cmd/link/internal/ld"

import (
    "bufio"
    "bytes"
    "cmd/internal/gcprog"
    "cmd/internal/obj"
    "crypto/sha1"
    "debug/elf"
    "debug/macho"
    "encoding/binary"
    "flag"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "os"
    "os/exec"
    "path"
    "path/filepath"
    "reflect"
    "runtime"
    "runtime/pprof"
    "sort"
    "strconv"
    "strings"
    "sync"
    "time"
    "unsafe"
)

const (
    ARMAG = "!<arch>\n"
)

const (
    BUCKETSIZE    = 256 * MINFUNC
    SUBBUCKETS    = 16
    SUBBUCKETSIZE = BUCKETSIZE / SUBBUCKETS
    NOIDX         = 0x7fffffff
)

const (
    BuildmodeUnset BuildMode = iota
    BuildmodeExe
    BuildmodePIE
    BuildmodeCArchive
    BuildmodeCShared
    BuildmodeShared
)

//  *  Emit .debug_frame
const (
    CIERESERVE          = 16
    DATAALIGNMENTFACTOR = -4
)

//  *  Walk DWarfDebugInfoEntries, and emit .debug_info
const (
    COMPUNITHEADERSIZE = 4 + 2 + 4 + 1
)

// for dynexport field of LSym
const (
    CgoExportDynamic = 1 << 0
    CgoExportStatic  = 1 << 1
)

// Index into the abbrevs table below.
// Keep in sync with ispubname() and ispubtype() below.
// ispubtype considers >= NULLTYPE public
const (
    DW_ABRV_NULL = iota
    DW_ABRV_COMPUNIT
    DW_ABRV_FUNCTION
    DW_ABRV_VARIABLE
    DW_ABRV_AUTO
    DW_ABRV_PARAM
    DW_ABRV_STRUCTFIELD
    DW_ABRV_FUNCTYPEPARAM
    DW_ABRV_DOTDOTDOT
    DW_ABRV_ARRAYRANGE
    DW_ABRV_NULLTYPE
    DW_ABRV_BASETYPE
    DW_ABRV_ARRAYTYPE
    DW_ABRV_CHANTYPE
    DW_ABRV_FUNCTYPE
    DW_ABRV_IFACETYPE
    DW_ABRV_MAPTYPE
    DW_ABRV_PTRTYPE
    DW_ABRV_BARE_PTRTYPE // only for void*, no DW_AT_type attr to please gdb 6.
    DW_ABRV_SLICETYPE
    DW_ABRV_STRINGTYPE
    DW_ABRV_STRUCTTYPE
    DW_ABRV_TYPEDECL
    DW_NABRV
)

// Table 28
const (
    DW_ACCESS_public    = 0x01
    DW_ACCESS_protected = 0x02
    DW_ACCESS_private   = 0x03
)

// Table 25
const (
    DW_ATE_address         = 0x01
    DW_ATE_boolean         = 0x02
    DW_ATE_complex_float   = 0x03
    DW_ATE_float           = 0x04
    DW_ATE_signed          = 0x05
    DW_ATE_signed_char     = 0x06
    DW_ATE_unsigned        = 0x07
    DW_ATE_unsigned_char   = 0x08
    DW_ATE_imaginary_float = 0x09
    DW_ATE_packed_decimal  = 0x0a
    DW_ATE_numeric_string  = 0x0b
    DW_ATE_edited          = 0x0c
    DW_ATE_signed_fixed    = 0x0d
    DW_ATE_unsigned_fixed  = 0x0e
    DW_ATE_decimal_float   = 0x0f
    DW_ATE_lo_user         = 0x80
    DW_ATE_hi_user         = 0xff
)

// Go-specific type attributes.
const (
    DW_AT_go_kind = 0x2900
    DW_AT_go_key  = 0x2901
    DW_AT_go_elem = 0x2902

    DW_AT_internal_location = 253 // params and locals; not emitted
)

// Table 20
const (
    DW_AT_sibling              = 0x01 // reference
    DW_AT_location             = 0x02 // block, loclistptr
    DW_AT_name                 = 0x03 // string
    DW_AT_ordering             = 0x09 // constant
    DW_AT_byte_size            = 0x0b // block, constant, reference
    DW_AT_bit_offset           = 0x0c // block, constant, reference
    DW_AT_bit_size             = 0x0d // block, constant, reference
    DW_AT_stmt_list            = 0x10 // lineptr
    DW_AT_low_pc               = 0x11 // address
    DW_AT_high_pc              = 0x12 // address
    DW_AT_language             = 0x13 // constant
    DW_AT_discr                = 0x15 // reference
    DW_AT_discr_value          = 0x16 // constant
    DW_AT_visibility           = 0x17 // constant
    DW_AT_import               = 0x18 // reference
    DW_AT_string_length        = 0x19 // block, loclistptr
    DW_AT_common_reference     = 0x1a // reference
    DW_AT_comp_dir             = 0x1b // string
    DW_AT_const_value          = 0x1c // block, constant, string
    DW_AT_containing_type      = 0x1d // reference
    DW_AT_default_value        = 0x1e // reference
    DW_AT_inline               = 0x20 // constant
    DW_AT_is_optional          = 0x21 // flag
    DW_AT_lower_bound          = 0x22 // block, constant, reference
    DW_AT_producer             = 0x25 // string
    DW_AT_prototyped           = 0x27 // flag
    DW_AT_return_addr          = 0x2a // block, loclistptr
    DW_AT_start_scope          = 0x2c // constant
    DW_AT_bit_stride           = 0x2e // constant
    DW_AT_upper_bound          = 0x2f // block, constant, reference
    DW_AT_abstract_origin      = 0x31 // reference
    DW_AT_accessibility        = 0x32 // constant
    DW_AT_address_class        = 0x33 // constant
    DW_AT_artificial           = 0x34 // flag
    DW_AT_base_types           = 0x35 // reference
    DW_AT_calling_convention   = 0x36 // constant
    DW_AT_count                = 0x37 // block, constant, reference
    DW_AT_data_member_location = 0x38 // block, constant, loclistptr
    DW_AT_decl_column          = 0x39 // constant
    DW_AT_decl_file            = 0x3a // constant
    DW_AT_decl_line            = 0x3b // constant
    DW_AT_declaration          = 0x3c // flag
    DW_AT_discr_list           = 0x3d // block
    DW_AT_encoding             = 0x3e // constant
    DW_AT_external             = 0x3f // flag
    DW_AT_frame_base           = 0x40 // block, loclistptr
    DW_AT_friend               = 0x41 // reference
    DW_AT_identifier_case      = 0x42 // constant
    DW_AT_macro_info           = 0x43 // macptr
    DW_AT_namelist_item        = 0x44 // block
    DW_AT_priority             = 0x45 // reference
    DW_AT_segment              = 0x46 // block, loclistptr
    DW_AT_specification        = 0x47 // reference
    DW_AT_static_link          = 0x48 // block, loclistptr
    DW_AT_type                 = 0x49 // reference
    DW_AT_use_location         = 0x4a // block, loclistptr
    DW_AT_variable_parameter   = 0x4b // flag
    DW_AT_virtuality           = 0x4c // constant
    DW_AT_vtable_elem_location = 0x4d // block, loclistptr
    // Dwarf3
    DW_AT_allocated      = 0x4e // block, constant, reference
    DW_AT_associated     = 0x4f // block, constant, reference
    DW_AT_data_location  = 0x50 // block
    DW_AT_byte_stride    = 0x51 // block, constant, reference
    DW_AT_entry_pc       = 0x52 // address
    DW_AT_use_UTF8       = 0x53 // flag
    DW_AT_extension      = 0x54 // reference
    DW_AT_ranges         = 0x55 // rangelistptr
    DW_AT_trampoline     = 0x56 // address, flag, reference, string
    DW_AT_call_column    = 0x57 // constant
    DW_AT_call_file      = 0x58 // constant
    DW_AT_call_line      = 0x59 // constant
    DW_AT_description    = 0x5a // string
    DW_AT_binary_scale   = 0x5b // constant
    DW_AT_decimal_scale  = 0x5c // constant
    DW_AT_small          = 0x5d // reference
    DW_AT_decimal_sign   = 0x5e // constant
    DW_AT_digit_count    = 0x5f // constant
    DW_AT_picture_string = 0x60 // string
    DW_AT_mutable        = 0x61 // flag
    DW_AT_threads_scaled = 0x62 // flag
    DW_AT_explicit       = 0x63 // flag
    DW_AT_object_pointer = 0x64 // reference
    DW_AT_endianity      = 0x65 // constant
    DW_AT_elemental      = 0x66 // flag
    DW_AT_pure           = 0x67 // flag
    DW_AT_recursive      = 0x68 // flag

    DW_AT_lo_user = 0x2000 // ---
    DW_AT_hi_user = 0x3fff // ---
)

// Table 33
const (
    DW_CC_normal  = 0x01
    DW_CC_program = 0x02
    DW_CC_nocall  = 0x03
    DW_CC_lo_user = 0x40
    DW_CC_hi_user = 0xff
)

// Table 40.
const (
    // operand,...
    DW_CFA_nop              = 0x00
    DW_CFA_set_loc          = 0x01 // address
    DW_CFA_advance_loc1     = 0x02 // 1-byte delta
    DW_CFA_advance_loc2     = 0x03 // 2-byte delta
    DW_CFA_advance_loc4     = 0x04 // 4-byte delta
    DW_CFA_offset_extended  = 0x05 // ULEB128 register, ULEB128 offset
    DW_CFA_restore_extended = 0x06 // ULEB128 register
    DW_CFA_undefined        = 0x07 // ULEB128 register
    DW_CFA_same_value       = 0x08 // ULEB128 register
    DW_CFA_register         = 0x09 // ULEB128 register, ULEB128 register
    DW_CFA_remember_state   = 0x0a
    DW_CFA_restore_state    = 0x0b

    DW_CFA_def_cfa            = 0x0c // ULEB128 register, ULEB128 offset
    DW_CFA_def_cfa_register   = 0x0d // ULEB128 register
    DW_CFA_def_cfa_offset     = 0x0e // ULEB128 offset
    DW_CFA_def_cfa_expression = 0x0f // BLOCK
    DW_CFA_expression         = 0x10 // ULEB128 register, BLOCK
    DW_CFA_offset_extended_sf = 0x11 // ULEB128 register, SLEB128 offset
    DW_CFA_def_cfa_sf         = 0x12 // ULEB128 register, SLEB128 offset
    DW_CFA_def_cfa_offset_sf  = 0x13 // SLEB128 offset
    DW_CFA_val_offset         = 0x14 // ULEB128, ULEB128
    DW_CFA_val_offset_sf      = 0x15 // ULEB128, SLEB128
    DW_CFA_val_expression     = 0x16 // ULEB128, BLOCK

    DW_CFA_lo_user = 0x1c
    DW_CFA_hi_user = 0x3f

    // Opcodes that take an addend operand.
    DW_CFA_advance_loc = 0x1 << 6 // +delta
    DW_CFA_offset      = 0x2 << 6 // +register (ULEB128 offset)
    DW_CFA_restore     = 0x3 << 6 // +register
)

// Table 19
const (
    DW_CHILDREN_no  = 0x00
    DW_CHILDREN_yes = 0x01
)

// Not from the spec, but logicaly belongs here
const (
    DW_CLS_ADDRESS = 0x01 + iota
    DW_CLS_BLOCK
    DW_CLS_CONSTANT
    DW_CLS_FLAG
    DW_CLS_PTR // lineptr, loclistptr, macptr, rangelistptr
    DW_CLS_REFERENCE
    DW_CLS_ADDRLOC
    DW_CLS_STRING
)

// Table 36
const (
    DW_DSC_label = 0x00
    DW_DSC_range = 0x01
)

// Table 26
const (
    DW_DS_unsigned           = 0x01
    DW_DS_leading_overpunch  = 0x02
    DW_DS_trailing_overpunch = 0x03
    DW_DS_leading_separate   = 0x04
    DW_DS_trailing_separate  = 0x05
)

// Table 27
const (
    DW_END_default = 0x00
    DW_END_big     = 0x01
    DW_END_little  = 0x02
    DW_END_lo_user = 0x40
    DW_END_hi_user = 0xff
)

// Table 21
const (
    DW_FORM_addr      = 0x01 // address
    DW_FORM_block2    = 0x03 // block
    DW_FORM_block4    = 0x04 // block
    DW_FORM_data2     = 0x05 // constant
    DW_FORM_data4     = 0x06 // constant, lineptr, loclistptr, macptr, rangelistptr
    DW_FORM_data8     = 0x07 // constant, lineptr, loclistptr, macptr, rangelistptr
    DW_FORM_string    = 0x08 // string
    DW_FORM_block     = 0x09 // block
    DW_FORM_block1    = 0x0a // block
    DW_FORM_data1     = 0x0b // constant
    DW_FORM_flag      = 0x0c // flag
    DW_FORM_sdata     = 0x0d // constant
    DW_FORM_strp      = 0x0e // string
    DW_FORM_udata     = 0x0f // constant
    DW_FORM_ref_addr  = 0x10 // reference
    DW_FORM_ref1      = 0x11 // reference
    DW_FORM_ref2      = 0x12 // reference
    DW_FORM_ref4      = 0x13 // reference
    DW_FORM_ref8      = 0x14 // reference
    DW_FORM_ref_udata = 0x15 // reference
    DW_FORM_indirect  = 0x16 // (see Section 7.5.3)
)

// Table 32
const (
    DW_ID_case_sensitive   = 0x00
    DW_ID_up_case          = 0x01
    DW_ID_down_case        = 0x02
    DW_ID_case_insensitive = 0x03
)

// Table 34
const (
    DW_INL_not_inlined          = 0x00
    DW_INL_inlined              = 0x01
    DW_INL_declared_not_inlined = 0x02
    DW_INL_declared_inlined     = 0x03
)

// Table 31
const (
    DW_LANG_C89         = 0x0001
    DW_LANG_C           = 0x0002
    DW_LANG_Ada83       = 0x0003
    DW_LANG_C_plus_plus = 0x0004
    DW_LANG_Cobol74     = 0x0005
    DW_LANG_Cobol85     = 0x0006
    DW_LANG_Fortran77   = 0x0007
    DW_LANG_Fortran90   = 0x0008
    DW_LANG_Pascal83    = 0x0009
    DW_LANG_Modula2     = 0x000a
    // Dwarf3
    DW_LANG_Java           = 0x000b
    DW_LANG_C99            = 0x000c
    DW_LANG_Ada95          = 0x000d
    DW_LANG_Fortran95      = 0x000e
    DW_LANG_PLI            = 0x000f
    DW_LANG_ObjC           = 0x0010
    DW_LANG_ObjC_plus_plus = 0x0011
    DW_LANG_UPC            = 0x0012
    DW_LANG_D              = 0x0013
    // Dwarf4
    DW_LANG_Python = 0x0014
    // Dwarf5
    DW_LANG_Go = 0x0016

    DW_LANG_lo_user = 0x8000
    DW_LANG_hi_user = 0xffff
)

// Table 38
const (
    DW_LNE_end_sequence = 0x01
    DW_LNE_set_address  = 0x02
    DW_LNE_define_file  = 0x03
    DW_LNE_lo_user      = 0x80
    DW_LNE_hi_user      = 0xff
)

// Table 37
const (
    DW_LNS_copy             = 0x01
    DW_LNS_advance_pc       = 0x02
    DW_LNS_advance_line     = 0x03
    DW_LNS_set_file         = 0x04
    DW_LNS_set_column       = 0x05
    DW_LNS_negate_stmt      = 0x06
    DW_LNS_set_basic_block  = 0x07
    DW_LNS_const_add_pc     = 0x08
    DW_LNS_fixed_advance_pc = 0x09
    // Dwarf3
    DW_LNS_set_prologue_end   = 0x0a
    DW_LNS_set_epilogue_begin = 0x0b
    DW_LNS_set_isa            = 0x0c
)

// Table 39
const (
    DW_MACINFO_define     = 0x01
    DW_MACINFO_undef      = 0x02
    DW_MACINFO_start_file = 0x03
    DW_MACINFO_end_file   = 0x04
    DW_MACINFO_vendor_ext = 0xff
)

// Table 24 (#operands, notes)
const (
    DW_OP_addr                = 0x03 // 1 constant address (size target specific)
    DW_OP_deref               = 0x06 // 0
    DW_OP_const1u             = 0x08 // 1 1-byte constant
    DW_OP_const1s             = 0x09 // 1 1-byte constant
    DW_OP_const2u             = 0x0a // 1 2-byte constant
    DW_OP_const2s             = 0x0b // 1 2-byte constant
    DW_OP_const4u             = 0x0c // 1 4-byte constant
    DW_OP_const4s             = 0x0d // 1 4-byte constant
    DW_OP_const8u             = 0x0e // 1 8-byte constant
    DW_OP_const8s             = 0x0f // 1 8-byte constant
    DW_OP_constu              = 0x10 // 1 ULEB128 constant
    DW_OP_consts              = 0x11 // 1 SLEB128 constant
    DW_OP_dup                 = 0x12 // 0
    DW_OP_drop                = 0x13 // 0
    DW_OP_over                = 0x14 // 0
    DW_OP_pick                = 0x15 // 1 1-byte stack index
    DW_OP_swap                = 0x16 // 0
    DW_OP_rot                 = 0x17 // 0
    DW_OP_xderef              = 0x18 // 0
    DW_OP_abs                 = 0x19 // 0
    DW_OP_and                 = 0x1a // 0
    DW_OP_div                 = 0x1b // 0
    DW_OP_minus               = 0x1c // 0
    DW_OP_mod                 = 0x1d // 0
    DW_OP_mul                 = 0x1e // 0
    DW_OP_neg                 = 0x1f // 0
    DW_OP_not                 = 0x20 // 0
    DW_OP_or                  = 0x21 // 0
    DW_OP_plus                = 0x22 // 0
    DW_OP_plus_uconst         = 0x23 // 1 ULEB128 addend
    DW_OP_shl                 = 0x24 // 0
    DW_OP_shr                 = 0x25 // 0
    DW_OP_shra                = 0x26 // 0
    DW_OP_xor                 = 0x27 // 0
    DW_OP_skip                = 0x2f // 1 signed 2-byte constant
    DW_OP_bra                 = 0x28 // 1 signed 2-byte constant
    DW_OP_eq                  = 0x29 // 0
    DW_OP_ge                  = 0x2a // 0
    DW_OP_gt                  = 0x2b // 0
    DW_OP_le                  = 0x2c // 0
    DW_OP_lt                  = 0x2d // 0
    DW_OP_ne                  = 0x2e // 0
    DW_OP_lit0                = 0x30 // 0 ...
    DW_OP_lit31               = 0x4f // 0 literals 0..31 = (DW_OP_lit0 + literal)
    DW_OP_reg0                = 0x50 // 0 ..
    DW_OP_reg31               = 0x6f // 0 reg 0..31 = (DW_OP_reg0 + regnum)
    DW_OP_breg0               = 0x70 // 1 ...
    DW_OP_breg31              = 0x8f // 1 SLEB128 offset base register 0..31 = (DW_OP_breg0 + regnum)
    DW_OP_regx                = 0x90 // 1 ULEB128 register
    DW_OP_fbreg               = 0x91 // 1 SLEB128 offset
    DW_OP_bregx               = 0x92 // 2 ULEB128 register followed by SLEB128 offset
    DW_OP_piece               = 0x93 // 1 ULEB128 size of piece addressed
    DW_OP_deref_size          = 0x94 // 1 1-byte size of data retrieved
    DW_OP_xderef_size         = 0x95 // 1 1-byte size of data retrieved
    DW_OP_nop                 = 0x96 // 0
    DW_OP_push_object_address = 0x97 // 0
    DW_OP_call2               = 0x98 // 1 2-byte offset of DIE
    DW_OP_call4               = 0x99 // 1 4-byte offset of DIE
    DW_OP_call_ref            = 0x9a // 1 4- or 8-byte offset of DIE
    DW_OP_form_tls_address    = 0x9b // 0
    DW_OP_call_frame_cfa      = 0x9c // 0
    DW_OP_bit_piece           = 0x9d // 2
    DW_OP_lo_user             = 0xe0
    DW_OP_hi_user             = 0xff
)

// Table 35
const (
    DW_ORD_row_major = 0x00
    DW_ORD_col_major = 0x01
)

// Table 18
const (
    DW_TAG_array_type               = 0x01
    DW_TAG_class_type               = 0x02
    DW_TAG_entry_point              = 0x03
    DW_TAG_enumeration_type         = 0x04
    DW_TAG_formal_parameter         = 0x05
    DW_TAG_imported_declaration     = 0x08
    DW_TAG_label                    = 0x0a
    DW_TAG_lexical_block            = 0x0b
    DW_TAG_member                   = 0x0d
    DW_TAG_pointer_type             = 0x0f
    DW_TAG_reference_type           = 0x10
    DW_TAG_compile_unit             = 0x11
    DW_TAG_string_type              = 0x12
    DW_TAG_structure_type           = 0x13
    DW_TAG_subroutine_type          = 0x15
    DW_TAG_typedef                  = 0x16
    DW_TAG_union_type               = 0x17
    DW_TAG_unspecified_parameters   = 0x18
    DW_TAG_variant                  = 0x19
    DW_TAG_common_block             = 0x1a
    DW_TAG_common_inclusion         = 0x1b
    DW_TAG_inheritance              = 0x1c
    DW_TAG_inlined_subroutine       = 0x1d
    DW_TAG_module                   = 0x1e
    DW_TAG_ptr_to_member_type       = 0x1f
    DW_TAG_set_type                 = 0x20
    DW_TAG_subrange_type            = 0x21
    DW_TAG_with_stmt                = 0x22
    DW_TAG_access_declaration       = 0x23
    DW_TAG_base_type                = 0x24
    DW_TAG_catch_block              = 0x25
    DW_TAG_const_type               = 0x26
    DW_TAG_constant                 = 0x27
    DW_TAG_enumerator               = 0x28
    DW_TAG_file_type                = 0x29
    DW_TAG_friend                   = 0x2a
    DW_TAG_namelist                 = 0x2b
    DW_TAG_namelist_item            = 0x2c
    DW_TAG_packed_type              = 0x2d
    DW_TAG_subprogram               = 0x2e
    DW_TAG_template_type_parameter  = 0x2f
    DW_TAG_template_value_parameter = 0x30
    DW_TAG_thrown_type              = 0x31
    DW_TAG_try_block                = 0x32
    DW_TAG_variant_part             = 0x33
    DW_TAG_variable                 = 0x34
    DW_TAG_volatile_type            = 0x35
    // Dwarf3
    DW_TAG_dwarf_procedure  = 0x36
    DW_TAG_restrict_type    = 0x37
    DW_TAG_interface_type   = 0x38
    DW_TAG_namespace        = 0x39
    DW_TAG_imported_module  = 0x3a
    DW_TAG_unspecified_type = 0x3b
    DW_TAG_partial_unit     = 0x3c
    DW_TAG_imported_unit    = 0x3d
    DW_TAG_condition        = 0x3f
    DW_TAG_shared_type      = 0x40
    // Dwarf4
    DW_TAG_type_unit             = 0x41
    DW_TAG_rvalue_reference_type = 0x42
    DW_TAG_template_alias        = 0x43

    // User defined
    DW_TAG_lo_user = 0x4080
    DW_TAG_hi_user = 0xffff
)

// Table 30
const (
    DW_VIRTUALITY_none         = 0x00
    DW_VIRTUALITY_virtual      = 0x01
    DW_VIRTUALITY_pure_virtual = 0x02
)

// Table 29
const (
    DW_VIS_local     = 0x01
    DW_VIS_exported  = 0x02
    DW_VIS_qualified = 0x03
)

const (
    EI_MAG0              = 0
    EI_MAG1              = 1
    EI_MAG2              = 2
    EI_MAG3              = 3
    EI_CLASS             = 4
    EI_DATA              = 5
    EI_VERSION           = 6
    EI_OSABI             = 7
    EI_ABIVERSION        = 8
    OLD_EI_BRAND         = 8
    EI_PAD               = 9
    EI_NIDENT            = 16
    ELFMAG0              = 0x7f
    ELFMAG1              = 'E'
    ELFMAG2              = 'L'
    ELFMAG3              = 'F'
    SELFMAG              = 4
    EV_NONE              = 0
    EV_CURRENT           = 1
    ELFCLASSNONE         = 0
    ELFCLASS32           = 1
    ELFCLASS64           = 2
    ELFDATANONE          = 0
    ELFDATA2LSB          = 1
    ELFDATA2MSB          = 2
    ELFOSABI_NONE        = 0
    ELFOSABI_HPUX        = 1
    ELFOSABI_NETBSD      = 2
    ELFOSABI_LINUX       = 3
    ELFOSABI_HURD        = 4
    ELFOSABI_86OPEN      = 5
    ELFOSABI_SOLARIS     = 6
    ELFOSABI_AIX         = 7
    ELFOSABI_IRIX        = 8
    ELFOSABI_FREEBSD     = 9
    ELFOSABI_TRU64       = 10
    ELFOSABI_MODESTO     = 11
    ELFOSABI_OPENBSD     = 12
    ELFOSABI_OPENVMS     = 13
    ELFOSABI_NSK         = 14
    ELFOSABI_ARM         = 97
    ELFOSABI_STANDALONE  = 255
    ELFOSABI_SYSV        = ELFOSABI_NONE
    ELFOSABI_MONTEREY    = ELFOSABI_AIX
    ET_NONE              = 0
    ET_REL               = 1
    ET_EXEC              = 2
    ET_DYN               = 3
    ET_CORE              = 4
    ET_LOOS              = 0xfe00
    ET_HIOS              = 0xfeff
    ET_LOPROC            = 0xff00
    ET_HIPROC            = 0xffff
    EM_NONE              = 0
    EM_M32               = 1
    EM_SPARC             = 2
    EM_386               = 3
    EM_68K               = 4
    EM_88K               = 5
    EM_860               = 7
    EM_MIPS              = 8
    EM_S370              = 9
    EM_MIPS_RS3_LE       = 10
    EM_PARISC            = 15
    EM_VPP500            = 17
    EM_SPARC32PLUS       = 18
    EM_960               = 19
    EM_PPC               = 20
    EM_PPC64             = 21
    EM_S390              = 22
    EM_V800              = 36
    EM_FR20              = 37
    EM_RH32              = 38
    EM_RCE               = 39
    EM_ARM               = 40
    EM_SH                = 42
    EM_SPARCV9           = 43
    EM_TRICORE           = 44
    EM_ARC               = 45
    EM_H8_300            = 46
    EM_H8_300H           = 47
    EM_H8S               = 48
    EM_H8_500            = 49
    EM_IA_64             = 50
    EM_MIPS_X            = 51
    EM_COLDFIRE          = 52
    EM_68HC12            = 53
    EM_MMA               = 54
    EM_PCP               = 55
    EM_NCPU              = 56
    EM_NDR1              = 57
    EM_STARCORE          = 58
    EM_ME16              = 59
    EM_ST100             = 60
    EM_TINYJ             = 61
    EM_X86_64            = 62
    EM_AARCH64           = 183
    EM_486               = 6
    EM_MIPS_RS4_BE       = 10
    EM_ALPHA_STD         = 41
    EM_ALPHA             = 0x9026
    SHN_UNDEF            = 0
    SHN_LORESERVE        = 0xff00
    SHN_LOPROC           = 0xff00
    SHN_HIPROC           = 0xff1f
    SHN_LOOS             = 0xff20
    SHN_HIOS             = 0xff3f
    SHN_ABS              = 0xfff1
    SHN_COMMON           = 0xfff2
    SHN_XINDEX           = 0xffff
    SHN_HIRESERVE        = 0xffff
    SHT_NULL             = 0
    SHT_PROGBITS         = 1
    SHT_SYMTAB           = 2
    SHT_STRTAB           = 3
    SHT_RELA             = 4
    SHT_HASH             = 5
    SHT_DYNAMIC          = 6
    SHT_NOTE             = 7
    SHT_NOBITS           = 8
    SHT_REL              = 9
    SHT_SHLIB            = 10
    SHT_DYNSYM           = 11
    SHT_INIT_ARRAY       = 14
    SHT_FINI_ARRAY       = 15
    SHT_PREINIT_ARRAY    = 16
    SHT_GROUP            = 17
    SHT_SYMTAB_SHNDX     = 18
    SHT_LOOS             = 0x60000000
    SHT_HIOS             = 0x6fffffff
    SHT_GNU_VERDEF       = 0x6ffffffd
    SHT_GNU_VERNEED      = 0x6ffffffe
    SHT_GNU_VERSYM       = 0x6fffffff
    SHT_LOPROC           = 0x70000000
    SHT_ARM_ATTRIBUTES   = 0x70000003
    SHT_HIPROC           = 0x7fffffff
    SHT_LOUSER           = 0x80000000
    SHT_HIUSER           = 0xffffffff
    SHF_WRITE            = 0x1
    SHF_ALLOC            = 0x2
    SHF_EXECINSTR        = 0x4
    SHF_MERGE            = 0x10
    SHF_STRINGS          = 0x20
    SHF_INFO_LINK        = 0x40
    SHF_LINK_ORDER       = 0x80
    SHF_OS_NONCONFORMING = 0x100
    SHF_GROUP            = 0x200
    SHF_TLS              = 0x400
    SHF_MASKOS           = 0x0ff00000
    SHF_MASKPROC         = 0xf0000000
    PT_NULL              = 0
    PT_LOAD              = 1
    PT_DYNAMIC           = 2
    PT_INTERP            = 3
    PT_NOTE              = 4
    PT_SHLIB             = 5
    PT_PHDR              = 6
    PT_TLS               = 7
    PT_LOOS              = 0x60000000
    PT_HIOS              = 0x6fffffff
    PT_LOPROC            = 0x70000000
    PT_HIPROC            = 0x7fffffff
    PT_GNU_STACK         = 0x6474e551
    PT_PAX_FLAGS         = 0x65041580
    PF_X                 = 0x1
    PF_W                 = 0x2
    PF_R                 = 0x4
    PF_MASKOS            = 0x0ff00000
    PF_MASKPROC          = 0xf0000000
    DT_NULL              = 0
    DT_NEEDED            = 1
    DT_PLTRELSZ          = 2
    DT_PLTGOT            = 3
    DT_HASH              = 4
    DT_STRTAB            = 5
    DT_SYMTAB            = 6
    DT_RELA              = 7
    DT_RELASZ            = 8
    DT_RELAENT           = 9
    DT_STRSZ             = 10
    DT_SYMENT            = 11
    DT_INIT              = 12
    DT_FINI              = 13
    DT_SONAME            = 14
    DT_RPATH             = 15
    DT_SYMBOLIC          = 16
    DT_REL               = 17
    DT_RELSZ             = 18
    DT_RELENT            = 19
    DT_PLTREL            = 20
    DT_DEBUG             = 21
    DT_TEXTREL           = 22
    DT_JMPREL            = 23
    DT_BIND_NOW          = 24
    DT_INIT_ARRAY        = 25
    DT_FINI_ARRAY        = 26
    DT_INIT_ARRAYSZ      = 27
    DT_FINI_ARRAYSZ      = 28
    DT_RUNPATH           = 29
    DT_FLAGS             = 30
    DT_ENCODING          = 32
    DT_PREINIT_ARRAY     = 32
    DT_PREINIT_ARRAYSZ   = 33
    DT_LOOS              = 0x6000000d
    DT_HIOS              = 0x6ffff000
    DT_LOPROC            = 0x70000000
    DT_HIPROC            = 0x7fffffff
    DT_VERNEED           = 0x6ffffffe
    DT_VERNEEDNUM        = 0x6fffffff
    DT_VERSYM            = 0x6ffffff0
    DT_PPC64_GLINK       = DT_LOPROC + 0
    DT_PPC64_OPT         = DT_LOPROC + 3
    DF_ORIGIN            = 0x0001
    DF_SYMBOLIC          = 0x0002
    DF_TEXTREL           = 0x0004
    DF_BIND_NOW          = 0x0008
    DF_STATIC_TLS        = 0x0010
    NT_PRSTATUS          = 1
    NT_FPREGSET          = 2
    NT_PRPSINFO          = 3
    STB_LOCAL            = 0
    STB_GLOBAL           = 1
    STB_WEAK             = 2
    STB_LOOS             = 10
    STB_HIOS             = 12
    STB_LOPROC           = 13
    STB_HIPROC           = 15
    STT_NOTYPE           = 0
    STT_OBJECT           = 1
    STT_FUNC             = 2
    STT_SECTION          = 3
    STT_FILE             = 4
    STT_COMMON           = 5
    STT_TLS              = 6
    STT_LOOS             = 10
    STT_HIOS             = 12
    STT_LOPROC           = 13
    STT_HIPROC           = 15
    STV_DEFAULT          = 0x0
    STV_INTERNAL         = 0x1
    STV_HIDDEN           = 0x2
    STV_PROTECTED        = 0x3
    STN_UNDEF            = 0
)

//  * Go linker interface
const (
    ELF64HDRSIZE  = 64
    ELF64PHDRSIZE = 56
    ELF64SHDRSIZE = 64
    ELF64RELSIZE  = 16
    ELF64RELASIZE = 24
    ELF64SYMSIZE  = 24
    ELF32HDRSIZE  = 52
    ELF32PHDRSIZE = 32
    ELF32SHDRSIZE = 40
    ELF32SYMSIZE  = 16
    ELF32RELSIZE  = 8
)

//  * Total amount of space to reserve at the start of the file
//  * for Header, PHeaders, SHeaders, and interp.
//  * May waste some.
//  * On FreeBSD, cannot be larger than a page.
const (
    ELFRESERVE = 4096
)

// Build info note
const (
    ELF_NOTE_BUILDINFO_NAMESZ = 4
    ELF_NOTE_BUILDINFO_TAG    = 3
)

// Go specific notes
const (
    ELF_NOTE_GOPKGLIST_TAG = 1
    ELF_NOTE_GOABIHASH_TAG = 2
    ELF_NOTE_GODEPS_TAG    = 3
    ELF_NOTE_GOBUILDID_TAG = 4
)

// NetBSD Signature (as per sys/exec_elf.h)
const (
    ELF_NOTE_NETBSD_NAMESZ  = 7
    ELF_NOTE_NETBSD_DESCSZ  = 4
    ELF_NOTE_NETBSD_TAG     = 1
    ELF_NOTE_NETBSD_VERSION = 599000000 /* NetBSD 5.99 */
)

// OpenBSD Signature
const (
    ELF_NOTE_OPENBSD_NAMESZ  = 8
    ELF_NOTE_OPENBSD_DESCSZ  = 4
    ELF_NOTE_OPENBSD_TAG     = 1
    ELF_NOTE_OPENBSD_VERSION = 0
)

const (
    ElfAbiNone     = 0
    ElfAbiSystemV  = 0
    ElfAbiHPUX     = 1
    ElfAbiNetBSD   = 2
    ElfAbiLinux    = 3
    ElfAbiSolaris  = 6
    ElfAbiAix      = 7
    ElfAbiIrix     = 8
    ElfAbiFreeBSD  = 9
    ElfAbiTru64    = 10
    ElfAbiModesto  = 11
    ElfAbiOpenBSD  = 12
    ElfAbiARM      = 97
    ElfAbiEmbedded = 255
)

// Derived from Plan 9 from User Space's src/libmach/elf.h, elf.c
// http://code.swtch.com/plan9port/src/tip/src/libmach/
//
//     Copyright © 2004 Russ Cox.
//     Portions Copyright © 2008-2010 Google Inc.
//     Portions Copyright © 2010 The Go Authors.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
const (
    ElfClassNone = 0
    ElfClass32   = 1
    ElfClass64   = 2
)

const (
    ElfDataNone = 0
    ElfDataLsb  = 1
    ElfDataMsb  = 2
)

const (
    ElfMachNone        = 0
    ElfMach32100       = 1
    ElfMachSparc       = 2
    ElfMach386         = 3
    ElfMach68000       = 4
    ElfMach88000       = 5
    ElfMach486         = 6
    ElfMach860         = 7
    ElfMachMips        = 8
    ElfMachS370        = 9
    ElfMachMipsLe      = 10
    ElfMachParisc      = 15
    ElfMachVpp500      = 17
    ElfMachSparc32Plus = 18
    ElfMach960         = 19
    ElfMachPower       = 20
    ElfMachPower64     = 21
    ElfMachS390        = 22
    ElfMachV800        = 36
    ElfMachFr20        = 37
    ElfMachRh32        = 38
    ElfMachRce         = 39
    ElfMachArm         = 40
    ElfMachAlpha       = 41
    ElfMachSH          = 42
    ElfMachSparc9      = 43
    ElfMachAmd64       = 62
    ElfMachArm64       = 183
)

const (
    ElfNotePrStatus     = 1
    ElfNotePrFpreg      = 2
    ElfNotePrPsinfo     = 3
    ElfNotePrTaskstruct = 4
    ElfNotePrAuxv       = 6
    ElfNotePrXfpreg     = 0x46e62b7f
)

const (
    ElfProgNone      = 0
    ElfProgLoad      = 1
    ElfProgDynamic   = 2
    ElfProgInterp    = 3
    ElfProgNote      = 4
    ElfProgShlib     = 5
    ElfProgPhdr      = 6
    ElfProgFlagExec  = 0x1
    ElfProgFlagWrite = 0x2
    ElfProgFlagRead  = 0x4
)

const (
    ElfSectNone      = 0
    ElfSectProgbits  = 1
    ElfSectSymtab    = 2
    ElfSectStrtab    = 3
    ElfSectRela      = 4
    ElfSectHash      = 5
    ElfSectDynamic   = 6
    ElfSectNote      = 7
    ElfSectNobits    = 8
    ElfSectRel       = 9
    ElfSectShlib     = 10
    ElfSectDynsym    = 11
    ElfSectFlagWrite = 0x1
    ElfSectFlagAlloc = 0x2
    ElfSectFlagExec  = 0x4
)

//  *  Elf.
const (
    ElfStrDebugAbbrev = iota
    ElfStrDebugAranges
    ElfStrDebugFrame
    ElfStrDebugInfo
    ElfStrDebugLine
    ElfStrDebugLoc
    ElfStrDebugMacinfo
    ElfStrDebugPubNames
    ElfStrDebugPubTypes
    ElfStrDebugRanges
    ElfStrDebugStr
    ElfStrGDBScripts
    ElfStrRelDebugInfo
    ElfStrRelDebugAranges
    ElfStrRelDebugLine
    ElfStrRelDebugFrame
    NElfStrDbg
)

const (
    ElfSymBindLocal  = 0
    ElfSymBindGlobal = 1
    ElfSymBindWeak   = 2
)

const (
    ElfSymShnNone   = 0
    ElfSymShnAbs    = 0xFFF1
    ElfSymShnCommon = 0xFFF2
)

const (
    ElfSymTypeNone    = 0
    ElfSymTypeObject  = 1
    ElfSymTypeFunc    = 2
    ElfSymTypeSection = 3
    ElfSymTypeFile    = 4
)

const (
    ElfTypeNone         = 0
    ElfTypeRelocatable  = 1
    ElfTypeExecutable   = 2
    ElfTypeSharedObject = 3
    ElfTypeCore         = 4
)

//  whence for ldpkg
const (
    FileObj = 0 + iota
    ArchiveObj
    Pkgdef
)

//  * Debugging Information Entries and their attributes.
const (
    HASHSIZE = 107
)

const (
    IMAGE_FILE_MACHINE_I386              = 0x14c
    IMAGE_FILE_MACHINE_AMD64             = 0x8664
    IMAGE_FILE_RELOCS_STRIPPED           = 0x0001
    IMAGE_FILE_EXECUTABLE_IMAGE          = 0x0002
    IMAGE_FILE_LINE_NUMS_STRIPPED        = 0x0004
    IMAGE_FILE_LARGE_ADDRESS_AWARE       = 0x0020
    IMAGE_FILE_32BIT_MACHINE             = 0x0100
    IMAGE_FILE_DEBUG_STRIPPED            = 0x0200
    IMAGE_SCN_CNT_CODE                   = 0x00000020
    IMAGE_SCN_CNT_INITIALIZED_DATA       = 0x00000040
    IMAGE_SCN_CNT_UNINITIALIZED_DATA     = 0x00000080
    IMAGE_SCN_MEM_EXECUTE                = 0x20000000
    IMAGE_SCN_MEM_READ                   = 0x40000000
    IMAGE_SCN_MEM_WRITE                  = 0x80000000
    IMAGE_SCN_MEM_DISCARDABLE            = 0x2000000
    IMAGE_SCN_LNK_NRELOC_OVFL            = 0x1000000
    IMAGE_SCN_ALIGN_32BYTES              = 0x600000
    IMAGE_DIRECTORY_ENTRY_EXPORT         = 0
    IMAGE_DIRECTORY_ENTRY_IMPORT         = 1
    IMAGE_DIRECTORY_ENTRY_RESOURCE       = 2
    IMAGE_DIRECTORY_ENTRY_EXCEPTION      = 3
    IMAGE_DIRECTORY_ENTRY_SECURITY       = 4
    IMAGE_DIRECTORY_ENTRY_BASERELOC      = 5
    IMAGE_DIRECTORY_ENTRY_DEBUG          = 6
    IMAGE_DIRECTORY_ENTRY_COPYRIGHT      = 7
    IMAGE_DIRECTORY_ENTRY_ARCHITECTURE   = 7
    IMAGE_DIRECTORY_ENTRY_GLOBALPTR      = 8
    IMAGE_DIRECTORY_ENTRY_TLS            = 9
    IMAGE_DIRECTORY_ENTRY_LOAD_CONFIG    = 10
    IMAGE_DIRECTORY_ENTRY_BOUND_IMPORT   = 11
    IMAGE_DIRECTORY_ENTRY_IAT            = 12
    IMAGE_DIRECTORY_ENTRY_DELAY_IMPORT   = 13
    IMAGE_DIRECTORY_ENTRY_COM_DESCRIPTOR = 14
    IMAGE_SUBSYSTEM_WINDOWS_GUI          = 2
    IMAGE_SUBSYSTEM_WINDOWS_CUI          = 3
)

const (
    IMAGE_SYM_UNDEFINED              = 0
    IMAGE_SYM_ABSOLUTE               = -1
    IMAGE_SYM_DEBUG                  = -2
    IMAGE_SYM_TYPE_NULL              = 0
    IMAGE_SYM_TYPE_VOID              = 1
    IMAGE_SYM_TYPE_CHAR              = 2
    IMAGE_SYM_TYPE_SHORT             = 3
    IMAGE_SYM_TYPE_INT               = 4
    IMAGE_SYM_TYPE_LONG              = 5
    IMAGE_SYM_TYPE_FLOAT             = 6
    IMAGE_SYM_TYPE_DOUBLE            = 7
    IMAGE_SYM_TYPE_STRUCT            = 8
    IMAGE_SYM_TYPE_UNION             = 9
    IMAGE_SYM_TYPE_ENUM              = 10
    IMAGE_SYM_TYPE_MOE               = 11
    IMAGE_SYM_TYPE_BYTE              = 12
    IMAGE_SYM_TYPE_WORD              = 13
    IMAGE_SYM_TYPE_UINT              = 14
    IMAGE_SYM_TYPE_DWORD             = 15
    IMAGE_SYM_TYPE_PCODE             = 32768
    IMAGE_SYM_DTYPE_NULL             = 0
    IMAGE_SYM_DTYPE_POINTER          = 0x10
    IMAGE_SYM_DTYPE_FUNCTION         = 0x20
    IMAGE_SYM_DTYPE_ARRAY            = 0x30
    IMAGE_SYM_CLASS_END_OF_FUNCTION  = -1
    IMAGE_SYM_CLASS_NULL             = 0
    IMAGE_SYM_CLASS_AUTOMATIC        = 1
    IMAGE_SYM_CLASS_EXTERNAL         = 2
    IMAGE_SYM_CLASS_STATIC           = 3
    IMAGE_SYM_CLASS_REGISTER         = 4
    IMAGE_SYM_CLASS_EXTERNAL_DEF     = 5
    IMAGE_SYM_CLASS_LABEL            = 6
    IMAGE_SYM_CLASS_UNDEFINED_LABEL  = 7
    IMAGE_SYM_CLASS_MEMBER_OF_STRUCT = 8
    IMAGE_SYM_CLASS_ARGUMENT         = 9
    IMAGE_SYM_CLASS_STRUCT_TAG       = 10
    IMAGE_SYM_CLASS_MEMBER_OF_UNION  = 11
    IMAGE_SYM_CLASS_UNION_TAG        = 12
    IMAGE_SYM_CLASS_TYPE_DEFINITION  = 13
    IMAGE_SYM_CLASS_UNDEFINED_STATIC = 14
    IMAGE_SYM_CLASS_ENUM_TAG         = 15
    IMAGE_SYM_CLASS_MEMBER_OF_ENUM   = 16
    IMAGE_SYM_CLASS_REGISTER_PARAM   = 17
    IMAGE_SYM_CLASS_BIT_FIELD        = 18
    IMAGE_SYM_CLASS_FAR_EXTERNAL     = 68 /* Not in PECOFF v8 spec */
    IMAGE_SYM_CLASS_BLOCK            = 100
    IMAGE_SYM_CLASS_FUNCTION         = 101
    IMAGE_SYM_CLASS_END_OF_STRUCT    = 102
    IMAGE_SYM_CLASS_FILE             = 103
    IMAGE_SYM_CLASS_SECTION          = 104
    IMAGE_SYM_CLASS_WEAK_EXTERNAL    = 105
    IMAGE_SYM_CLASS_CLR_TOKEN        = 107
    IMAGE_REL_I386_ABSOLUTE          = 0x0000
    IMAGE_REL_I386_DIR16             = 0x0001
    IMAGE_REL_I386_REL16             = 0x0002
    IMAGE_REL_I386_DIR32             = 0x0006
    IMAGE_REL_I386_DIR32NB           = 0x0007
    IMAGE_REL_I386_SEG12             = 0x0009
    IMAGE_REL_I386_SECTION           = 0x000A
    IMAGE_REL_I386_SECREL            = 0x000B
    IMAGE_REL_I386_TOKEN             = 0x000C
    IMAGE_REL_I386_SECREL7           = 0x000D
    IMAGE_REL_I386_REL32             = 0x0014
    IMAGE_REL_AMD64_ABSOLUTE         = 0x0000
    IMAGE_REL_AMD64_ADDR64           = 0x0001
    IMAGE_REL_AMD64_ADDR32           = 0x0002
    IMAGE_REL_AMD64_ADDR32NB         = 0x0003
    IMAGE_REL_AMD64_REL32            = 0x0004
    IMAGE_REL_AMD64_REL32_1          = 0x0005
    IMAGE_REL_AMD64_REL32_2          = 0x0006
    IMAGE_REL_AMD64_REL32_3          = 0x0007
    IMAGE_REL_AMD64_REL32_4          = 0x0008
    IMAGE_REL_AMD64_REL32_5          = 0x0009
    IMAGE_REL_AMD64_SECTION          = 0x000A
    IMAGE_REL_AMD64_SECREL           = 0x000B
    IMAGE_REL_AMD64_SECREL7          = 0x000C
    IMAGE_REL_AMD64_TOKEN            = 0x000D
    IMAGE_REL_AMD64_SREL32           = 0x000E
    IMAGE_REL_AMD64_PAIR             = 0x000F
    IMAGE_REL_AMD64_SSPAN32          = 0x0010
)

//  * Total amount of space to reserve at the start of the file
//  * for Header, PHeaders, and SHeaders.
//  * May waste some.
const (
    INITIAL_MACHO_HEADR = 4 * 1024
)

const (
    LC_ID_DYLIB             = 0xd
    LC_LOAD_DYLINKER        = 0xe
    LC_PREBOUND_DYLIB       = 0x10
    LC_LOAD_WEAK_DYLIB      = 0x18
    LC_UUID                 = 0x1b
    LC_RPATH                = 0x8000001c
    LC_CODE_SIGNATURE       = 0x1d
    LC_SEGMENT_SPLIT_INFO   = 0x1e
    LC_REEXPORT_DYLIB       = 0x8000001f
    LC_ENCRYPTION_INFO      = 0x21
    LC_DYLD_INFO            = 0x22
    LC_DYLD_INFO_ONLY       = 0x80000022
    LC_VERSION_MIN_MACOSX   = 0x24
    LC_VERSION_MIN_IPHONEOS = 0x25
    LC_FUNCTION_STARTS      = 0x26
    LC_MAIN                 = 0x80000028
    LC_DATA_IN_CODE         = 0x29
    LC_SOURCE_VERSION       = 0x2A
    LC_DYLIB_CODE_SIGN_DRS  = 0x2B
    LC_ENCRYPTION_INFO_64   = 0x2C
)

//  * Generate short opcodes when possible, long ones when necessary.
//  * See section 6.2.5
const (
    LINE_BASE   = -1
    LINE_RANGE  = 4
    OPCODE_BASE = 10
)

const (
    LdMachoCpuVax         = 1
    LdMachoCpu68000       = 6
    LdMachoCpu386         = 7
    LdMachoCpuAmd64       = 0x1000007
    LdMachoCpuMips        = 8
    LdMachoCpu98000       = 10
    LdMachoCpuHppa        = 11
    LdMachoCpuArm         = 12
    LdMachoCpu88000       = 13
    LdMachoCpuSparc       = 14
    LdMachoCpu860         = 15
    LdMachoCpuAlpha       = 16
    LdMachoCpuPower       = 18
    LdMachoCmdSegment     = 1
    LdMachoCmdSymtab      = 2
    LdMachoCmdSymseg      = 3
    LdMachoCmdThread      = 4
    LdMachoCmdDysymtab    = 11
    LdMachoCmdSegment64   = 25
    LdMachoFileObject     = 1
    LdMachoFileExecutable = 2
    LdMachoFileFvmlib     = 3
    LdMachoFileCore       = 4
    LdMachoFilePreload    = 5
)

const (
    LinkAuto = 0 + iota
    LinkInternal
    LinkExternal
)

const (
    MACHO_CPU_AMD64               = 1<<24 | 7
    MACHO_CPU_386                 = 7
    MACHO_SUBCPU_X86              = 3
    MACHO_CPU_ARM                 = 12
    MACHO_SUBCPU_ARM              = 0
    MACHO_SUBCPU_ARMV7            = 9
    MACHO_CPU_ARM64               = 1<<24 | 12
    MACHO_SUBCPU_ARM64_ALL        = 0
    MACHO32SYMSIZE                = 12
    MACHO64SYMSIZE                = 16
    MACHO_X86_64_RELOC_UNSIGNED   = 0
    MACHO_X86_64_RELOC_SIGNED     = 1
    MACHO_X86_64_RELOC_BRANCH     = 2
    MACHO_X86_64_RELOC_GOT_LOAD   = 3
    MACHO_X86_64_RELOC_GOT        = 4
    MACHO_X86_64_RELOC_SUBTRACTOR = 5
    MACHO_X86_64_RELOC_SIGNED_1   = 6
    MACHO_X86_64_RELOC_SIGNED_2   = 7
    MACHO_X86_64_RELOC_SIGNED_4   = 8
    MACHO_ARM_RELOC_VANILLA       = 0
    MACHO_ARM_RELOC_BR24          = 5
    MACHO_ARM64_RELOC_UNSIGNED    = 0
    MACHO_ARM64_RELOC_BRANCH26    = 2
    MACHO_ARM64_RELOC_PAGE21      = 3
    MACHO_ARM64_RELOC_PAGEOFF12   = 4
    MACHO_ARM64_RELOC_ADDEND      = 10
    MACHO_GENERIC_RELOC_VANILLA   = 0
    MACHO_FAKE_GOTPCREL           = 100
)

const (
    MINFUNC = 16 // minimum size for a function
)

// synthesizemaptypes is way too closely married to runtime/hashmap.c
const (
    MaxKeySize = 128
    MaxValSize = 128
    BucketSize = 8
)

//  * We use the 64-bit data structures on both 32- and 64-bit machines
//  * in order to write the code just once.  The 64-bit data structure is
//  * written in the 32-bit format on the 32-bit machines.
const (
    NSECT = 48
)

// Derived from Plan 9 from User Space's src/libmach/elf.h, elf.c
// http://code.swtch.com/plan9port/src/tip/src/libmach/
//
//     Copyright © 2004 Russ Cox.
//     Portions Copyright © 2008-2010 Google Inc.
//     Portions Copyright © 2010 The Go Authors.
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
// of the Software, and to permit persons to whom the Software is furnished to do
// so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED,
// INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A
// PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT
// HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
// SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
const (
    N_EXT  = 0x01
    N_TYPE = 0x1e
    N_STAB = 0xe0
)

const (
    PEBASE = 0x00400000

    // SectionAlignment must be greater than or equal to FileAlignment.
    // The default is the page size for the architecture.
    PESECTALIGN = 0x1000

    // FileAlignment should be a power of 2 between 512 and 64 K, inclusive.
    // The default is 512. If the SectionAlignment is less than
    // the architecture's page size, then FileAlignment must match SectionAlignment.
    PEFILEALIGN = 2 << 8
)

// Reloc.variant
const (
    RV_NONE = iota
    RV_POWER_LO
    RV_POWER_HI
    RV_POWER_HA
    RV_POWER_DS
    RV_CHECK_OVERFLOW = 1 << 8
    RV_TYPE_MASK      = RV_CHECK_OVERFLOW - 1
)

//  * Relocation types.
const (
    R_X86_64_NONE           = 0
    R_X86_64_64             = 1
    R_X86_64_PC32           = 2
    R_X86_64_GOT32          = 3
    R_X86_64_PLT32          = 4
    R_X86_64_COPY           = 5
    R_X86_64_GLOB_DAT       = 6
    R_X86_64_JMP_SLOT       = 7
    R_X86_64_RELATIVE       = 8
    R_X86_64_GOTPCREL       = 9
    R_X86_64_32             = 10
    R_X86_64_32S            = 11
    R_X86_64_16             = 12
    R_X86_64_PC16           = 13
    R_X86_64_8              = 14
    R_X86_64_PC8            = 15
    R_X86_64_DTPMOD64       = 16
    R_X86_64_DTPOFF64       = 17
    R_X86_64_TPOFF64        = 18
    R_X86_64_TLSGD          = 19
    R_X86_64_TLSLD          = 20
    R_X86_64_DTPOFF32       = 21
    R_X86_64_GOTTPOFF       = 22
    R_X86_64_TPOFF32        = 23
    R_X86_64_PC64           = 24
    R_X86_64_GOTOFF64       = 25
    R_X86_64_GOTPC32        = 26
    R_X86_64_GOT64          = 27
    R_X86_64_GOTPCREL64     = 28
    R_X86_64_GOTPC64        = 29
    R_X86_64_GOTPLT64       = 30
    R_X86_64_PLTOFF64       = 31
    R_X86_64_SIZE32         = 32
    R_X86_64_SIZE64         = 33
    R_X86_64_GOTPC32_TLSDEC = 34
    R_X86_64_TLSDESC_CALL   = 35
    R_X86_64_TLSDESC        = 36
    R_X86_64_IRELATIVE      = 37
    R_X86_64_PC32_BND       = 40
    R_X86_64_GOTPCRELX      = 41
    R_X86_64_REX_GOTPCRELX  = 42

    R_AARCH64_ABS64                       = 257
    R_AARCH64_ABS32                       = 258
    R_AARCH64_CALL26                      = 283
    R_AARCH64_ADR_PREL_PG_HI21            = 275
    R_AARCH64_ADD_ABS_LO12_NC             = 277
    R_AARCH64_LDST8_ABS_LO12_NC           = 278
    R_AARCH64_LDST16_ABS_LO12_NC          = 284
    R_AARCH64_LDST32_ABS_LO12_NC          = 285
    R_AARCH64_LDST64_ABS_LO12_NC          = 286
    R_AARCH64_ADR_GOT_PAGE                = 311
    R_AARCH64_LD64_GOT_LO12_NC            = 312
    R_AARCH64_TLSIE_ADR_GOTTPREL_PAGE21   = 541
    R_AARCH64_TLSIE_LD64_GOTTPREL_LO12_NC = 542
    R_AARCH64_TLSLE_MOVW_TPREL_G0         = 547

    R_ALPHA_NONE           = 0
    R_ALPHA_REFLONG        = 1
    R_ALPHA_REFQUAD        = 2
    R_ALPHA_GPREL32        = 3
    R_ALPHA_LITERAL        = 4
    R_ALPHA_LITUSE         = 5
    R_ALPHA_GPDISP         = 6
    R_ALPHA_BRADDR         = 7
    R_ALPHA_HINT           = 8
    R_ALPHA_SREL16         = 9
    R_ALPHA_SREL32         = 10
    R_ALPHA_SREL64         = 11
    R_ALPHA_OP_PUSH        = 12
    R_ALPHA_OP_STORE       = 13
    R_ALPHA_OP_PSUB        = 14
    R_ALPHA_OP_PRSHIFT     = 15
    R_ALPHA_GPVALUE        = 16
    R_ALPHA_GPRELHIGH      = 17
    R_ALPHA_GPRELLOW       = 18
    R_ALPHA_IMMED_GP_16    = 19
    R_ALPHA_IMMED_GP_HI32  = 20
    R_ALPHA_IMMED_SCN_HI32 = 21
    R_ALPHA_IMMED_BR_HI32  = 22
    R_ALPHA_IMMED_LO32     = 23
    R_ALPHA_COPY           = 24
    R_ALPHA_GLOB_DAT       = 25
    R_ALPHA_JMP_SLOT       = 26
    R_ALPHA_RELATIVE       = 27

    R_ARM_NONE          = 0
    R_ARM_PC24          = 1
    R_ARM_ABS32         = 2
    R_ARM_REL32         = 3
    R_ARM_PC13          = 4
    R_ARM_ABS16         = 5
    R_ARM_ABS12         = 6
    R_ARM_THM_ABS5      = 7
    R_ARM_ABS8          = 8
    R_ARM_SBREL32       = 9
    R_ARM_THM_PC22      = 10
    R_ARM_THM_PC8       = 11
    R_ARM_AMP_VCALL9    = 12
    R_ARM_SWI24         = 13
    R_ARM_THM_SWI8      = 14
    R_ARM_XPC25         = 15
    R_ARM_THM_XPC22     = 16
    R_ARM_COPY          = 20
    R_ARM_GLOB_DAT      = 21
    R_ARM_JUMP_SLOT     = 22
    R_ARM_RELATIVE      = 23
    R_ARM_GOTOFF        = 24
    R_ARM_GOTPC         = 25
    R_ARM_GOT32         = 26
    R_ARM_PLT32         = 27
    R_ARM_CALL          = 28
    R_ARM_JUMP24        = 29
    R_ARM_V4BX          = 40
    R_ARM_GOT_PREL      = 96
    R_ARM_GNU_VTENTRY   = 100
    R_ARM_GNU_VTINHERIT = 101
    R_ARM_TLS_IE32      = 107
    R_ARM_TLS_LE32      = 108
    R_ARM_RSBREL32      = 250
    R_ARM_THM_RPC22     = 251
    R_ARM_RREL32        = 252
    R_ARM_RABS32        = 253
    R_ARM_RPC24         = 254
    R_ARM_RBASE         = 255

    R_386_NONE          = 0
    R_386_32            = 1
    R_386_PC32          = 2
    R_386_GOT32         = 3
    R_386_PLT32         = 4
    R_386_COPY          = 5
    R_386_GLOB_DAT      = 6
    R_386_JMP_SLOT      = 7
    R_386_RELATIVE      = 8
    R_386_GOTOFF        = 9
    R_386_GOTPC         = 10
    R_386_TLS_TPOFF     = 14
    R_386_TLS_IE        = 15
    R_386_TLS_GOTIE     = 16
    R_386_TLS_LE        = 17
    R_386_TLS_GD        = 18
    R_386_TLS_LDM       = 19
    R_386_TLS_GD_32     = 24
    R_386_TLS_GD_PUSH   = 25
    R_386_TLS_GD_CALL   = 26
    R_386_TLS_GD_POP    = 27
    R_386_TLS_LDM_32    = 28
    R_386_TLS_LDM_PUSH  = 29
    R_386_TLS_LDM_CALL  = 30
    R_386_TLS_LDM_POP   = 31
    R_386_TLS_LDO_32    = 32
    R_386_TLS_IE_32     = 33
    R_386_TLS_LE_32     = 34
    R_386_TLS_DTPMOD32  = 35
    R_386_TLS_DTPOFF32  = 36
    R_386_TLS_TPOFF32   = 37
    R_386_TLS_GOTDESC   = 39
    R_386_TLS_DESC_CALL = 40
    R_386_TLS_DESC      = 41
    R_386_IRELATIVE     = 42
    R_386_GOT32X        = 43

    R_PPC_NONE            = 0
    R_PPC_ADDR32          = 1
    R_PPC_ADDR24          = 2
    R_PPC_ADDR16          = 3
    R_PPC_ADDR16_LO       = 4
    R_PPC_ADDR16_HI       = 5
    R_PPC_ADDR16_HA       = 6
    R_PPC_ADDR14          = 7
    R_PPC_ADDR14_BRTAKEN  = 8
    R_PPC_ADDR14_BRNTAKEN = 9
    R_PPC_REL24           = 10
    R_PPC_REL14           = 11
    R_PPC_REL14_BRTAKEN   = 12
    R_PPC_REL14_BRNTAKEN  = 13
    R_PPC_GOT16           = 14
    R_PPC_GOT16_LO        = 15
    R_PPC_GOT16_HI        = 16
    R_PPC_GOT16_HA        = 17
    R_PPC_PLTREL24        = 18
    R_PPC_COPY            = 19
    R_PPC_GLOB_DAT        = 20
    R_PPC_JMP_SLOT        = 21
    R_PPC_RELATIVE        = 22
    R_PPC_LOCAL24PC       = 23
    R_PPC_UADDR32         = 24
    R_PPC_UADDR16         = 25
    R_PPC_REL32           = 26
    R_PPC_PLT32           = 27
    R_PPC_PLTREL32        = 28
    R_PPC_PLT16_LO        = 29
    R_PPC_PLT16_HI        = 30
    R_PPC_PLT16_HA        = 31
    R_PPC_SDAREL16        = 32
    R_PPC_SECTOFF         = 33
    R_PPC_SECTOFF_LO      = 34
    R_PPC_SECTOFF_HI      = 35
    R_PPC_SECTOFF_HA      = 36
    R_PPC_TLS             = 67
    R_PPC_DTPMOD32        = 68
    R_PPC_TPREL16         = 69
    R_PPC_TPREL16_LO      = 70
    R_PPC_TPREL16_HI      = 71
    R_PPC_TPREL16_HA      = 72
    R_PPC_TPREL32         = 73
    R_PPC_DTPREL16        = 74
    R_PPC_DTPREL16_LO     = 75
    R_PPC_DTPREL16_HI     = 76
    R_PPC_DTPREL16_HA     = 77
    R_PPC_DTPREL32        = 78
    R_PPC_GOT_TLSGD16     = 79
    R_PPC_GOT_TLSGD16_LO  = 80
    R_PPC_GOT_TLSGD16_HI  = 81
    R_PPC_GOT_TLSGD16_HA  = 82
    R_PPC_GOT_TLSLD16     = 83
    R_PPC_GOT_TLSLD16_LO  = 84
    R_PPC_GOT_TLSLD16_HI  = 85
    R_PPC_GOT_TLSLD16_HA  = 86
    R_PPC_GOT_TPREL16     = 87
    R_PPC_GOT_TPREL16_LO  = 88
    R_PPC_GOT_TPREL16_HI  = 89
    R_PPC_GOT_TPREL16_HA  = 90
    R_PPC_EMB_NADDR32     = 101
    R_PPC_EMB_NADDR16     = 102
    R_PPC_EMB_NADDR16_LO  = 103
    R_PPC_EMB_NADDR16_HI  = 104
    R_PPC_EMB_NADDR16_HA  = 105
    R_PPC_EMB_SDAI16      = 106
    R_PPC_EMB_SDA2I16     = 107
    R_PPC_EMB_SDA2REL     = 108
    R_PPC_EMB_SDA21       = 109
    R_PPC_EMB_MRKREF      = 110
    R_PPC_EMB_RELSEC16    = 111
    R_PPC_EMB_RELST_LO    = 112
    R_PPC_EMB_RELST_HI    = 113
    R_PPC_EMB_RELST_HA    = 114
    R_PPC_EMB_BIT_FLD     = 115
    R_PPC_EMB_RELSDA      = 116

    R_PPC64_ADDR32            = R_PPC_ADDR32
    R_PPC64_ADDR16_LO         = R_PPC_ADDR16_LO
    R_PPC64_ADDR16_HA         = R_PPC_ADDR16_HA
    R_PPC64_REL24             = R_PPC_REL24
    R_PPC64_GOT16_HA          = R_PPC_GOT16_HA
    R_PPC64_JMP_SLOT          = R_PPC_JMP_SLOT
    R_PPC64_TPREL16           = R_PPC_TPREL16
    R_PPC64_ADDR64            = 38
    R_PPC64_TOC16             = 47
    R_PPC64_TOC16_LO          = 48
    R_PPC64_TOC16_HI          = 49
    R_PPC64_TOC16_HA          = 50
    R_PPC64_ADDR16_LO_DS      = 57
    R_PPC64_GOT16_LO_DS       = 59
    R_PPC64_TOC16_DS          = 63
    R_PPC64_TOC16_LO_DS       = 64
    R_PPC64_TLS               = 67
    R_PPC64_GOT_TPREL16_LO_DS = 88
    R_PPC64_GOT_TPREL16_HA    = 90
    R_PPC64_REL16_LO          = 250
    R_PPC64_REL16_HI          = 251
    R_PPC64_REL16_HA          = 252

    R_SPARC_NONE     = 0
    R_SPARC_8        = 1
    R_SPARC_16       = 2
    R_SPARC_32       = 3
    R_SPARC_DISP8    = 4
    R_SPARC_DISP16   = 5
    R_SPARC_DISP32   = 6
    R_SPARC_WDISP30  = 7
    R_SPARC_WDISP22  = 8
    R_SPARC_HI22     = 9
    R_SPARC_22       = 10
    R_SPARC_13       = 11
    R_SPARC_LO10     = 12
    R_SPARC_GOT10    = 13
    R_SPARC_GOT13    = 14
    R_SPARC_GOT22    = 15
    R_SPARC_PC10     = 16
    R_SPARC_PC22     = 17
    R_SPARC_WPLT30   = 18
    R_SPARC_COPY     = 19
    R_SPARC_GLOB_DAT = 20
    R_SPARC_JMP_SLOT = 21
    R_SPARC_RELATIVE = 22
    R_SPARC_UA32     = 23
    R_SPARC_PLT32    = 24
    R_SPARC_HIPLT22  = 25
    R_SPARC_LOPLT10  = 26
    R_SPARC_PCPLT32  = 27
    R_SPARC_PCPLT22  = 28
    R_SPARC_PCPLT10  = 29
    R_SPARC_10       = 30
    R_SPARC_11       = 31
    R_SPARC_64       = 32
    R_SPARC_OLO10    = 33
    R_SPARC_HH22     = 34
    R_SPARC_HM10     = 35
    R_SPARC_LM22     = 36
    R_SPARC_PC_HH22  = 37
    R_SPARC_PC_HM10  = 38
    R_SPARC_PC_LM22  = 39
    R_SPARC_WDISP16  = 40
    R_SPARC_WDISP19  = 41
    R_SPARC_GLOB_JMP = 42
    R_SPARC_7        = 43
    R_SPARC_5        = 44
    R_SPARC_6        = 45
    R_SPARC_DISP64   = 46
    R_SPARC_PLT64    = 47
    R_SPARC_HIX22    = 48
    R_SPARC_LOX10    = 49
    R_SPARC_H44      = 50
    R_SPARC_M44      = 51
    R_SPARC_L44      = 52
    R_SPARC_REGISTER = 53
    R_SPARC_UA64     = 54
    R_SPARC_UA16     = 55

    ARM_MAGIC_TRAMP_NUMBER = 0x5c000003
)

const (
    SARMAG  = 8
    SAR_HDR = 16 + 44
)

const (
    SymKindLocal = 0 + iota
    SymKindExtdef
    SymKindUndef
    NumSymKind
)

const (
    Tag_file                 = 1
    Tag_CPU_name             = 4
    Tag_CPU_raw_name         = 5
    Tag_compatibility        = 32
    Tag_nodefaults           = 64
    Tag_also_compatible_with = 65
    Tag_ABI_VFP_args         = 28
)

var ELF_NOTE_BUILDINFO_NAME = []byte("GNU\x00")

var ELF_NOTE_GO_NAME = []byte("Go\x00\x00")

var ELF_NOTE_NETBSD_NAME = []byte("NetBSD\x00")

var ELF_NOTE_OPENBSD_NAME = []byte("OpenBSD\x00")

var ElfMagic = [4]uint8{0x7F, 'E', 'L', 'F'}

var Elfstrdat []byte

var Iself bool

var Link386 = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "386",
    Thechar:   '8',
    Minlc:     1,
    Ptrsize:   4,
    Regsize:   4,
}

var Linkamd64 = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "amd64",
    Thechar:   '6',
    Minlc:     1,
    Ptrsize:   8,
    Regsize:   8,
}

var Linkamd64p32 = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "amd64p32",
    Thechar:   '6',
    Minlc:     1,
    Ptrsize:   4,
    Regsize:   8,
}

var Linkarm = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "arm",
    Thechar:   '5',
    Minlc:     4,
    Ptrsize:   4,
    Regsize:   4,
}

var Linkarm64 = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "arm64",
    Thechar:   '7',
    Minlc:     4,
    Ptrsize:   8,
    Regsize:   8,
}

var Linkmips64 = LinkArch{
    ByteOrder: binary.BigEndian,
    Name:      "mips64",
    Thechar:   '0',
    Minlc:     4,
    Ptrsize:   8,
    Regsize:   8,
}

var Linkmips64le = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "mips64le",
    Thechar:   '0',
    Minlc:     4,
    Ptrsize:   8,
    Regsize:   8,
}

var Linkppc64 = LinkArch{
    ByteOrder: binary.BigEndian,
    Name:      "ppc64",
    Thechar:   '9',
    Minlc:     4,
    Ptrsize:   8,
    Regsize:   8,
}

var Linkppc64le = LinkArch{
    ByteOrder: binary.LittleEndian,
    Name:      "ppc64le",
    Thechar:   '9',
    Minlc:     4,
    Ptrsize:   8,
    Regsize:   8,
}

var Nelfsym int = 1

var PEFILEHEADR int32

var PESECTHEADR int32

var (
    Segtext   Segment
    Segrodata Segment
    Segdata   Segment
    Segdwarf  Segment
)

var (
    Thearch Arch

    Debug  [128]int
    Lcsize int32

    Spsize  int32
    Symsize int32
)

var (
    Thestring   string
    Thelinkarch *LinkArch

    Funcalign int

    Buildmode  BuildMode
    Linkshared bool

    Ctxt      *Link
    HEADR     int32
    HEADTYPE  int32
    INITRND   int32
    INITTEXT  int64
    INITDAT   int64
    INITENTRY string /* entry point */

    Linkmode int
)

var (

    // buffered output
    Bso obj.Biobuf
)

type ArHdr struct {
    name string
    date string
    uid  string
    gid  string
    mode string
    size string
    fmag string
}

type Arch struct {
    Thechar          int
    Ptrsize          int
    Intsize          int
    Regsize          int
    Funcalign        int
    Maxalign         int
    Minlc            int
    Dwarfregsp       int
    Dwarfreglr       int
    Linuxdynld       string
    Freebsddynld     string
    Netbsddynld      string
    Openbsddynld     string
    Dragonflydynld   string
    Solarisdynld     string
    Adddynrel        func(*LSym, *Reloc)
    Archinit         func()
    Archreloc        func(*Reloc, *LSym, *int64) int
    Archrelocvariant func(*Reloc, *LSym, int64) int64
    Asmb             func()
    Elfreloc1        func(*Reloc, int64) int
    Elfsetupplt      func()
    Gentext          func()
    Machoreloc1      func(*Reloc, int64) int
    PEreloc1         func(*Reloc, int64) bool
    Lput             func(uint32)
    Wput             func(uint16)
    Vput             func(uint64)
}

type Auto struct {
    Asym    *LSym
    Link    *Auto
    Aoffset int32
    Name    int16
    Gotype  *LSym
}

// A BuildMode indicates the sort of object we are building:
//
//     "exe": build a main package and everything it imports into an executable.
//     "c-shared": build a main package, plus all packages that it imports, into a
//       single C shared library. The only callable symbols will be those functions
//       marked as exported.
//     "shared": combine all packages passed on the command line, and their
//       dependencies, into a single shared library that will be used when
//       building with the -linkshared option.
type BuildMode uint8

type COFFSym struct {
    sym       *LSym
    strtbloff int
    sect      int
    value     int64
    typ       uint16
}

type Chain struct {
    sym   *LSym
    up    *Chain
    limit int // limit on entry to sym
}

type DWAbbrev struct {
    tag      uint8
    children uint8
    attr     []DWAttrForm
}

type DWAttr struct {
    link  *DWAttr
    atr   uint16 // DW_AT_
    cls   uint8  // DW_CLS_
    value int64
    data  interface{}
}

//  * Defining Abbrevs.  This is hardcoded, and there will be
//  * only a handful of them.  The DWARF spec places no restriction on
//  * the ordering of attributes in the Abbrevs and DIEs, and we will
//  * always write them out in the order of declaration in the abbrev.
type DWAttrForm struct {
    attr uint16
    form uint8
}

type DWDie struct {
    abbrev int
    link   *DWDie
    child  *DWDie
    attr   *DWAttr
    // offset into .debug_info section, i.e relative to
    // infoo. only valid after call to putdie()
    offs  int64
    hash  []*DWDie // optional index of children by name, enabled by mkindex()
    hlink *DWDie   // bucket chain in parent's index
}

type Dll struct {
    name     string
    nameoff  uint64
    thunkoff uint64
    ms       *Imp
    next     *Dll
}

//  * ELF header.
type ElfEhdr struct {
    ident     [EI_NIDENT]uint8
    type_     uint16
    machine   uint16
    version   uint32
    entry     uint64
    phoff     uint64
    shoff     uint64
    flags     uint32
    ehsize    uint16
    phentsize uint16
    phnum     uint16
    shentsize uint16
    shnum     uint16
    shstrndx  uint16
}

type ElfHdrBytes struct {
    Ident     [16]uint8
    Type      [2]uint8
    Machine   [2]uint8
    Version   [4]uint8
    Entry     [4]uint8
    Phoff     [4]uint8
    Shoff     [4]uint8
    Flags     [4]uint8
    Ehsize    [2]uint8
    Phentsize [2]uint8
    Phnum     [2]uint8
    Shentsize [2]uint8
    Shnum     [2]uint8
    Shstrndx  [2]uint8
}

type ElfHdrBytes64 struct {
    Ident     [16]uint8
    Type      [2]uint8
    Machine   [2]uint8
    Version   [4]uint8
    Entry     [8]uint8
    Phoff     [8]uint8
    Shoff     [8]uint8
    Flags     [4]uint8
    Ehsize    [2]uint8
    Phentsize [2]uint8
    Phnum     [2]uint8
    Shentsize [2]uint8
    Shnum     [2]uint8
    Shstrndx  [2]uint8
}

type ElfObj struct {
    f         *obj.Biobuf
    base      int64 // offset in f where ELF begins
    length    int64 // length of ELF
    is64      int
    name      string
    e         binary.ByteOrder
    sect      []ElfSect
    nsect     uint
    shstrtab  string
    nsymtab   int
    symtab    *ElfSect
    symstr    *ElfSect
    type_     uint32
    machine   uint32
    version   uint32
    entry     uint64
    phoff     uint64
    shoff     uint64
    flags     uint32
    ehsize    uint32
    phentsize uint32
    phnum     uint32
    shentsize uint32
    shnum     uint32
    shstrndx  uint32
}

//  * Program header.
type ElfPhdr struct {
    type_  uint32
    flags  uint32
    off    uint64
    vaddr  uint64
    paddr  uint64
    filesz uint64
    memsz  uint64
    align  uint64
}

type ElfProgBytes struct {
}

type ElfProgBytes64 struct {
}

type ElfSect struct {
    name    string
    nameoff uint32
    type_   uint32
    flags   uint64
    addr    uint64
    off     uint64
    size    uint64
    link    uint32
    info    uint32
    align   uint64
    entsize uint64
    base    []byte
    sym     *LSym
}

type ElfSectBytes struct {
    Name    [4]uint8
    Type    [4]uint8
    Flags   [4]uint8
    Addr    [4]uint8
    Off     [4]uint8
    Size    [4]uint8
    Link    [4]uint8
    Info    [4]uint8
    Align   [4]uint8
    Entsize [4]uint8
}

type ElfSectBytes64 struct {
    Name    [4]uint8
    Type    [4]uint8
    Flags   [8]uint8
    Addr    [8]uint8
    Off     [8]uint8
    Size    [8]uint8
    Link    [4]uint8
    Info    [4]uint8
    Align   [8]uint8
    Entsize [8]uint8
}

//  * Section header.
type ElfShdr struct {
    name      uint32
    type_     uint32
    flags     uint64
    addr      uint64
    off       uint64
    size      uint64
    link      uint32
    info      uint32
    addralign uint64
    entsize   uint64
    shnum     int
    secsym    *LSym
}

type ElfSym struct {
    name  string
    value uint64
    size  uint64
    bind  uint8
    type_ uint8
    other uint8
    shndx uint16
    sym   *LSym
}

type ElfSymBytes struct {
    Name  [4]uint8
    Value [4]uint8
    Size  [4]uint8
    Info  uint8
    Other uint8
    Shndx [2]uint8
}

type ElfSymBytes64 struct {
    Name  [4]uint8
    Info  uint8
    Other uint8
    Shndx [2]uint8
    Value [8]uint8
    Size  [8]uint8
}

//  * Note header.  The ".note" section contains an array of notes.  Each
//  * begins with this header, aligned to a word boundary.  Immediately
//  * following the note header is n_namesz bytes of name, padded to the
//  * next word boundary.  Then comes n_descsz bytes of descriptor, again
//  * padded to a word boundary.  The values of n_namesz and n_descsz do
//  * not include the padding.
type Elf_Note struct {
    n_namesz uint32
    n_descsz uint32
    n_type   uint32
}

type Elfaux struct {
    next *Elfaux
    num  int
    vers string
}

type Elflib struct {
    next *Elflib
    aux  *Elfaux
    file string
}

type Elfstring struct {
    s   string
    off int
}

type GCProg struct {
    sym *LSym
    w   gcprog.Writer
}

type Hostobj struct {
    ld     func(*obj.Biobuf, string, int64, string)
    pkg    string
    pn     string
    file   string
    off    int64
    length int64
}

type IMAGE_DATA_DIRECTORY struct {
    VirtualAddress uint32
    Size           uint32
}

type IMAGE_EXPORT_DIRECTORY struct {
    Characteristics       uint32
    TimeDateStamp         uint32
    MajorVersion          uint16
    MinorVersion          uint16
    Name                  uint32
    Base                  uint32
    NumberOfFunctions     uint32
    NumberOfNames         uint32
    AddressOfFunctions    uint32
    AddressOfNames        uint32
    AddressOfNameOrdinals uint32
}

type IMAGE_FILE_HEADER struct {
    Machine              uint16
    NumberOfSections     uint16
    TimeDateStamp        uint32
    PointerToSymbolTable uint32
    NumberOfSymbols      uint32
    SizeOfOptionalHeader uint16
    Characteristics      uint16
}

type IMAGE_IMPORT_DESCRIPTOR struct {
    OriginalFirstThunk uint32
    TimeDateStamp      uint32
    ForwarderChain     uint32
    Name               uint32
    FirstThunk         uint32
}

type IMAGE_OPTIONAL_HEADER struct {
    Magic                       uint16
    MajorLinkerVersion          uint8
    MinorLinkerVersion          uint8
    SizeOfCode                  uint32
    SizeOfInitializedData       uint32
    SizeOfUninitializedData     uint32
    AddressOfEntryPoint         uint32
    BaseOfCode                  uint32
    BaseOfData                  uint32
    ImageBase                   uint32
    SectionAlignment            uint32
    FileAlignment               uint32
    MajorOperatingSystemVersion uint16
    MinorOperatingSystemVersion uint16
    MajorImageVersion           uint16
    MinorImageVersion           uint16
    MajorSubsystemVersion       uint16
    MinorSubsystemVersion       uint16
    Win32VersionValue           uint32
    SizeOfImage                 uint32
    SizeOfHeaders               uint32
    CheckSum                    uint32
    Subsystem                   uint16
    DllCharacteristics          uint16
    SizeOfStackReserve          uint32
    SizeOfStackCommit           uint32
    SizeOfHeapReserve           uint32
    SizeOfHeapCommit            uint32
    LoaderFlags                 uint32
    NumberOfRvaAndSizes         uint32
    DataDirectory               [16]IMAGE_DATA_DIRECTORY
}

type IMAGE_SECTION_HEADER struct {
    Name                 [8]uint8
    VirtualSize          uint32
    VirtualAddress       uint32
    SizeOfRawData        uint32
    PointerToRawData     uint32
    PointerToRelocations uint32
    PointerToLineNumbers uint32
    NumberOfRelocations  uint16
    NumberOfLineNumbers  uint16
    Characteristics      uint32
}

type Imp struct {
    s       *LSym
    off     uint64
    next    *Imp
    argsize int
}

type LSym struct {
    Name       string
    Extname    string
    Type       int16
    Version    int16
    Dupok      uint8
    Cfunc      uint8
    External   uint8
    Nosplit    uint8
    Reachable  bool
    Cgoexport  uint8
    Special    uint8
    Stkcheck   uint8
    Hide       uint8
    Leaf       uint8
    Localentry uint8
    Onlist     uint8
    // ElfType is set for symbols read from shared libraries by ldshlibsyms. It
    // is not set for symbols defined by the packages being linked or by symbols
    // read by ldelf (and so is left as elf.STT_NOTYPE).
    ElfType     elf.SymType
    Dynid       int32
    Plt         int32
    Got         int32
    Align       int32
    Elfsym      int32
    LocalElfsym int32
    Args        int32
    Locals      int32
    Value       int64
    Size        int64
    Allsym      *LSym
    Next        *LSym
    Sub         *LSym
    Outer       *LSym
    Gotype      *LSym
    Reachparent *LSym
    Queue       *LSym
    File        string
    Dynimplib   string
    Dynimpvers  string
    Sect        *Section
    Autom       *Auto
    Pcln        *Pcln
    P           []byte
    R           []Reloc
    Local       bool
}

type LdMachoCmd struct {
    type_ int
    off   uint32
    size  uint32
    seg   LdMachoSeg
    sym   LdMachoSymtab
    dsym  LdMachoDysymtab
}

type LdMachoDysymtab struct {
    ilocalsym      uint32
    nlocalsym      uint32
    iextdefsym     uint32
    nextdefsym     uint32
    iundefsym      uint32
    nundefsym      uint32
    tocoff         uint32
    ntoc           uint32
    modtaboff      uint32
    nmodtab        uint32
    extrefsymoff   uint32
    nextrefsyms    uint32
    indirectsymoff uint32
    nindirectsyms  uint32
    extreloff      uint32
    nextrel        uint32
    locreloff      uint32
    nlocrel        uint32
    indir          []uint32
}

type LdMachoObj struct {
    f          *obj.Biobuf
    base       int64 // off in f where Mach-O begins
    length     int64 // length of Mach-O
    is64       bool
    name       string
    e          binary.ByteOrder
    cputype    uint
    subcputype uint
    filetype   uint32
    flags      uint32
    cmd        []LdMachoCmd
    ncmd       uint
}

type LdMachoRel struct {
    addr      uint32
    symnum    uint32
    pcrel     uint8
    length    uint8
    extrn     uint8
    type_     uint8
    scattered uint8
    value     uint32
}

type LdMachoSect struct {
    name    string
    segname string
    addr    uint64
    size    uint64
    off     uint32
    align   uint32
    reloff  uint32
    nreloc  uint32
    flags   uint32
    res1    uint32
    res2    uint32
    sym     *LSym
    rel     []LdMachoRel
}

type LdMachoSeg struct {
    name     string
    vmaddr   uint64
    vmsize   uint64
    fileoff  uint32
    filesz   uint32
    maxprot  uint32
    initprot uint32
    nsect    uint32
    flags    uint32
    sect     []LdMachoSect
}

type LdMachoSym struct {
    name    string
    type_   uint8
    sectnum uint8
    desc    uint16
    kind    int8
    value   uint64
    sym     *LSym
}

type LdMachoSymtab struct {
    symoff  uint32
    nsym    uint32
    stroff  uint32
    strsize uint32
    str     []byte
    sym     []LdMachoSym
}

type Library struct {
    Objref string
    Srcref string
    File   string
    Pkg    string
    Shlib  string
    hash   []byte
}

type Link struct {
    Thechar    int32
    Thestring  string
    Goarm      int32
    Headtype   int
    Arch       *LinkArch
    Debugasm   int32
    Debugvlog  int32
    Bso        *obj.Biobuf
    Windows    int32
    Goroot     string
    Hash       map[symVer]*LSym
    Allsym     *LSym
    Nsymbol    int32
    Tlsg       *LSym
    Libdir     []string
    Library    []*Library
    Shlibs     []Shlib
    Tlsoffset  int
    Diag       func(string, ...interface{})
    Cursym     *LSym
    Version    int
    Textp      *LSym
    Etextp     *LSym
    Nhistfile  int32
    Filesyms   *LSym
    Moduledata *LSym
}

type LinkArch struct {
    ByteOrder binary.ByteOrder
    Name      string
    Thechar   int
    Minlc     int
    Ptrsize   int
    Regsize   int
}

type MachoHdr struct {
    cpu    uint32
    subcpu uint32
}

type MachoLoad struct {
    type_ uint32
    data  []uint32
}

type MachoSect struct {
    name    string
    segname string
    addr    uint64
    size    uint64
    off     uint32
    align   uint32
    reloc   uint32
    nreloc  uint32
    flag    uint32
    res1    uint32
    res2    uint32
}

type MachoSeg struct {
    name       string
    vsize      uint64
    vaddr      uint64
    fileoffset uint64
    filesize   uint64
    prot1      uint32
    prot2      uint32
    nsect      uint32
    msect      uint32
    sect       []MachoSect
    flag       uint32
}

// X64
type PE64_IMAGE_OPTIONAL_HEADER struct {
    Magic                       uint16
    MajorLinkerVersion          uint8
    MinorLinkerVersion          uint8
    SizeOfCode                  uint32
    SizeOfInitializedData       uint32
    SizeOfUninitializedData     uint32
    AddressOfEntryPoint         uint32
    BaseOfCode                  uint32
    ImageBase                   uint64
    SectionAlignment            uint32
    FileAlignment               uint32
    MajorOperatingSystemVersion uint16
    MinorOperatingSystemVersion uint16
    MajorImageVersion           uint16
    MinorImageVersion           uint16
    MajorSubsystemVersion       uint16
    MinorSubsystemVersion       uint16
    Win32VersionValue           uint32
    SizeOfImage                 uint32
    SizeOfHeaders               uint32
    CheckSum                    uint32
    Subsystem                   uint16
    DllCharacteristics          uint16
    SizeOfStackReserve          uint64
    SizeOfStackCommit           uint64
    SizeOfHeapReserve           uint64
    SizeOfHeapCommit            uint64
    LoaderFlags                 uint32
    NumberOfRvaAndSizes         uint32
    DataDirectory               [16]IMAGE_DATA_DIRECTORY
}

type Pcdata struct {
    P []byte
}

type Pciter struct {
    d       Pcdata
    p       []byte
    pc      uint32
    nextpc  uint32
    pcscale uint32
    value   int32
    start   int
    done    int
}

type Pcln struct {
    Pcsp        Pcdata
    Pcfile      Pcdata
    Pcline      Pcdata
    Pcdata      []Pcdata
    Npcdata     int
    Funcdata    []*LSym
    Funcdataoff []int64
    Nfuncdata   int
    File        []*LSym
    Nfile       int
    Mfile       int
    Lastfile    *LSym
    Lastindex   int
}

type PeObj struct {
    f      *obj.Biobuf
    name   string
    base   uint32
    sect   []PeSect
    nsect  uint
    pesym  []PeSym
    npesym uint
    fh     IMAGE_FILE_HEADER
    snames []byte
}

type PeSect struct {
    name string
    base []byte
    size uint64
    sym  *LSym
    sh   IMAGE_SECTION_HEADER
}

type PeSym struct {
    name    string
    value   uint32
    sectnum uint16
    type_   uint16
    sclass  uint8
    aux     uint8
    sym     *LSym
}

type Pkg struct {
    mark    bool
    checked bool
    path    string
    impby   []*Pkg
}

type Reloc struct {
    Off     int32
    Siz     uint8
    Done    uint8
    Type    int32
    Variant int32
    Add     int64
    Xadd    int64
    Sym     *LSym
    Xsym    *LSym
}

type Rpath struct {
    set bool
    val string
}

type Section struct {
    Rwx     uint8
    Extnum  int16
    Align   int32
    Name    string
    Vaddr   uint64
    Length  uint64
    Next    *Section
    Seg     *Segment
    Elfsect *ElfShdr
    Reloff  uint64
    Rellen  uint64
}

type Segment struct {
    Rwx     uint8  // permission as usual unix bits (5 = r-x etc)
    Vaddr   uint64 // virtual address
    Length  uint64 // length in memory
    Fileoff uint64 // file offset
    Filelen uint64 // length on disk
    Sect    *Section
}

type Shlib struct {
    Path             string
    Hash             []byte
    Deps             []string
    File             *elf.File
    gcdata_addresses map[*LSym]uint64
}

func Addaddr(ctxt *Link, s *LSym, t *LSym) int64

func Addaddrplus(ctxt *Link, s *LSym, t *LSym, add int64) int64

func Adddynsym(ctxt *Link, s *LSym)

func Addpcrelplus(ctxt *Link, s *LSym, t *LSym, add int64) int64

func Addrel(s *LSym) *Reloc

func Addstring(s *LSym, str string) int64

func Adduint16(ctxt *Link, s *LSym, v uint16) int64

func Adduint32(ctxt *Link, s *LSym, v uint32) int64

func Adduint64(ctxt *Link, s *LSym, v uint64) int64

func Adduint8(ctxt *Link, s *LSym, v uint8) int64

func Asmbelf(symo int64)

func Asmbelfsetup()

func Asmbmacho()

func Asmbpe()

func Asmelfsym()

func Asmplan9sym()

func AtExit(f func())

func Be16(b []byte) uint16

func Be32(b []byte) uint32

func Cflush()

func Codeblk(addr int64, size int64)

func Cpos() int64

func Cput(c uint8)

func Cseek(p int64)

func Cwrite(p []byte)

func Datblk(addr int64, size int64)

func Diag(format string, args ...interface{})

func Domacholink() int64

//  * This is the main entry point for generating dwarf.  After emitting
//  * the mandatory debug_abbrev section, it calls writelines() to set up
//  * the per-compilation unit part of the DIE tree, while simultaneously
//  * emitting the debug_line section.  When the final tree contains
//  * forward references, it will write the debug_info section in 2
//  * passes.
//  *
func Dwarfemitdebugsections()

// DynlinkingGo returns whether we are producing Go code that can live
// in separate shared libraries linked together at runtime.
func DynlinkingGo() bool

func ELF32_R_INFO(sym uint32, type_ uint32) uint32

func ELF32_R_SYM(info uint32) uint32

func ELF32_R_TYPE(info uint32) uint32

func ELF32_ST_BIND(info uint8) uint8

func ELF32_ST_INFO(bind uint8, type_ uint8) uint8

func ELF32_ST_TYPE(info uint8) uint8

func ELF32_ST_VISIBILITY(oth uint8) uint8

func ELF64_R_INFO(sym uint32, type_ uint32) uint64

func ELF64_R_SYM(info uint64) uint32

func ELF64_R_TYPE(info uint64) uint32

func ELF64_ST_BIND(info uint8) uint8

func ELF64_ST_INFO(bind uint8, type_ uint8) uint8

func ELF64_ST_TYPE(info uint8) uint8

func ELF64_ST_VISIBILITY(oth uint8) uint8

func Elfadddynsym(ctxt *Link, s *LSym)

func Elfemitreloc()

// Initialize the global variable that describes the ELF header. It will be updated
// as we write section and prog headers.
func Elfinit()

func Elfwritedynent(s *LSym, tag int, val uint64)

func Elfwritedynentsymplus(s *LSym, tag int, t *LSym, add int64)

func Entryvalue() int64

func Exit(code int)

func Exitf(format string, a ...interface{})

func Headstr(v int) string

func Ldmain()

func Le16(b []byte) uint16

func Le32(b []byte) uint32

func Le64(b []byte) uint64

func Lflag(arg string)

func Linklookup(ctxt *Link, name string, v int) *LSym

// read-only lookup
func Linkrlookup(ctxt *Link, name string, v int) *LSym

func Lputb(l uint32)

func Lputl(l uint32)

func Machoadddynlib(lib string)

func Machoemitreloc()

func Machoinit()

func Peinit()

func Rnd(v int64, r int64) int64

func Symaddr(s *LSym) int64

func Symgrow(ctxt *Link, s *LSym, siz int64)

// UseRelro returns whether to make use of "read only relocations" aka
// relro.
func UseRelro() bool

func Vputb(v uint64)

func Vputl(v uint64)

func Wputb(w uint16)

func Wputl(w uint16)

func (*BuildMode) Set(s string) error

func (*BuildMode) String() string

func (*GCProg) AddSym(s *LSym)

func (*GCProg) End(size int64)

func (*GCProg) Init(name string)

func (*LSym) ElfsymForReloc() int32

func (*LSym) String() string

// The smallest possible offset from the hardware stack pointer to a local
// variable on the stack. Architectures that use a link register save its value
// on the stack in the function prologue and so always have a pointer between
// the hardware stack pointer and the local variable area.
func (*Link) FixedFrameSize() int64

func (*Rpath) Set(val string) error

func (*Rpath) String() string

