// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build ingore

// Package unicode provides data and functions to test some properties of
// Unicode code points.

// unicode 包提供了一些测试Unicode码点属性的数据和函数.
package unicode


const (
	MaxRune         = '\U0010FFFF' // Maximum valid Unicode code point. // Unicode码点的最大值
	ReplacementChar = '\uFFFD'     // Represents invalid code points.   // 无效码点的表示
	MaxASCII        = '\u007F'     // maximum ASCII value.              // ASCII的最大值
	MaxLatin1       = '\u00FF'     // maximum Latin-1 value.            // Latin-1的最大值

)


// Indices into the Delta arrays inside CaseRanges for case mapping.

// CaseRange 中 Delta 数组的下标，以用于写法映射。
const (
	UpperCase = iota
	LowerCase
	TitleCase
	MaxCase
)


// If the Delta field of a CaseRange is UpperLower, it means
// this CaseRange represents a sequence of the form (say)
// Upper Lower Upper Lower.

// 若 CaseRange 的 Delta 字段为 UpperLower 或 LowerUpper，则该 CaseRange 即表示
// （所谓的）“Upper Lower Upper Lower”序列。
const (
	UpperLower = MaxRune + 1 // (Cannot be a valid delta.) // 不能算有效的 delta。

)


// Version is the Unicode edition from which the tables are derived.

// Version 为得到此表所用的 Unicode 版本。
const Version = "8.0.0"


// These variables have type *RangeTable.

// 这些变量的类型为 *RangeTable。
var (
	ASCII_Hex_Digit                    = _ASCII_Hex_Digit                    // ASCII_Hex_Digit 为带属性 ASCII_Hex_Digit 的 Unicode 字符集合。
	Bidi_Control                       = _Bidi_Control                       // Bidi_Control 为带属性 Bidi_Control 的 Unicode 字符集合。
	Dash                               = _Dash                               // Dash 为带属性 Dash 的 Unicode 字符集合。
	Deprecated                         = _Deprecated                         // Deprecated 为带属性 Deprecated 的 Unicode 字符集合。
	Diacritic                          = _Diacritic                          // Diacritic 为带属性 Diacritic 的 Unicode 字符集合。
	Extender                           = _Extender                           // Extender 为带属性 Extender 的 Unicode 字符集合。
	Hex_Digit                          = _Hex_Digit                          // Hex_Digit 为带属性 Hex_Digit 的 Unicode 字符集合。
	Hyphen                             = _Hyphen                             // Hyphen 为带属性 Hyphen 的 Unicode 字符集合。
	IDS_Binary_Operator                = _IDS_Binary_Operator                // IDS_Binary_Operator 为带属性 IDS_Binary_Operator 的 Unicode 字符集合。
	IDS_Trinary_Operator               = _IDS_Trinary_Operator               // IDS_Trinary_Operator 为带属性 IDS_Trinary_Operator 的 Unicode 字符集合。
	Ideographic                        = _Ideographic                        // Ideographic 为带属性 Ideographic 的 Unicode 字符集合。
	Join_Control                       = _Join_Control                       // Join_Control 为带属性 Join_Control 的 Unicode 字符集合。
	Logical_Order_Exception            = _Logical_Order_Exception            // Logical_Order_Exception 为带属性 Logical_Order_Exception 的 Unicode 字符集合。
	Noncharacter_Code_Point            = _Noncharacter_Code_Point            // Noncharacter_Code_Point 为带属性 Noncharacter_Code_Point 的 Unicode 字符集合。
	Other_Alphabetic                   = _Other_Alphabetic                   // Other_Alphabetic 为带属性 Other_Alphabetic 的 Unicode 字符集合。
	Other_Default_Ignorable_Code_Point = _Other_Default_Ignorable_Code_Point // Other_Default_Ignorable_Code_Point 为带属性 Other_Default_Ignorable_Code_Point 的 Unicode 字符集合。
	Other_Grapheme_Extend              = _Other_Grapheme_Extend              // Other_Grapheme_Extend 为带属性 Other_Grapheme_Extend 的 Unicode 字符集合。
	Other_ID_Continue                  = _Other_ID_Continue                  // Other_ID_Continue 为带属性 Other_ID_Continue 的 Unicode 字符集合。
	Other_ID_Start                     = _Other_ID_Start                     // Other_ID_Start 为带属性 Other_ID_Start 的 Unicode 字符集合。
	Other_Lowercase                    = _Other_Lowercase                    // Other_Lowercase 为带属性 Other_Lowercase 的 Unicode 字符集合。
	Other_Math                         = _Other_Math                         // Other_Math 为带属性 Other_Math 的 Unicode 字符集合。
	Other_Uppercase                    = _Other_Uppercase                    // Other_Uppercase 为带属性 Other_Uppercase 的 Unicode 字符集合。
	Pattern_Syntax                     = _Pattern_Syntax                     // Pattern_Syntax 为带属性 Pattern_Syntax 的 Unicode 字符集合。
	Pattern_White_Space                = _Pattern_White_Space                // Pattern_White_Space 为带属性 Pattern_White_Space 的 Unicode 字符集合。
	Quotation_Mark                     = _Quotation_Mark                     // Quotation_Mark 为带属性 Quotation_Mark 的 Unicode 字符集合。
	Radical                            = _Radical                            // Radical 为带属性 Radical 的 Unicode 字符集合。
	STerm                              = _STerm                              // STerm 为带属性 STerm 的 Unicode 字符集合。
	Soft_Dotted                        = _Soft_Dotted                        // Soft_Dotted 为带属性 Soft_Dotted 的 Unicode 字符集合。
	Terminal_Punctuation               = _Terminal_Punctuation               // Terminal_Punctuation 为带属性 Terminal_Punctuation 的 Unicode 字符集合。
	Unified_Ideograph                  = _Unified_Ideograph                  // Unified_Ideograph 为带属性 Unified_Ideograph 的 Unicode 字符集合。
	Variation_Selector                 = _Variation_Selector                 // Variation_Selector 为带属性 Variation_Selector 的 Unicode 字符集合。
	White_Space                        = _White_Space                        // White_Space 为带属性 White_Space 的 Unicode 字符集合。

)


// These variables have type *RangeTable.

// 这些变量的类型为 *RangeTable。
var (
	Ahom                   = _Ahom                   // Ahom 为书写系统 Ahom 中的 Unicode 字符集合。
	Anatolian_Hieroglyphs  = _Anatolian_Hieroglyphs  // Anatolian_Hieroglyphs 为书写系统 Anatolian_Hieroglyphs 中的 Unicode 字符集合。
	Arabic                 = _Arabic                 // Arabic 为书写系统 Arabic 中的 Unicode 字符集合。
	Armenian               = _Armenian               // Armenian 为书写系统 Armenian 中的 Unicode 字符集合。
	Avestan                = _Avestan                // Avestan 为书写系统 Avestan 中的 Unicode 字符集合。
	Balinese               = _Balinese               // Balinese 为书写系统 Balinese 中的 Unicode 字符集合。
	Bamum                  = _Bamum                  // Bamum 为书写系统 Bamum 中的 Unicode 字符集合。
	Bassa_Vah              = _Bassa_Vah              // Bassa_Vah 为书写系统 Bassa_Vah 中的 Unicode 字符集合。
	Batak                  = _Batak                  // Batak 为书写系统 Batak 中的 Unicode 字符集合。
	Bengali                = _Bengali                // Bengali 为书写系统 Bengali 中的 Unicode 字符集合。
	Bopomofo               = _Bopomofo               // Bopomofo 为书写系统 Bopomofo 中的 Unicode 字符集合。
	Brahmi                 = _Brahmi                 // Brahmi 为书写系统 Brahmi 中的 Unicode 字符集合。
	Braille                = _Braille                // Braille 为书写系统 Braille 中的 Unicode 字符集合。
	Buginese               = _Buginese               // Buginese 为书写系统 Buginese 中的 Unicode 字符集合。
	Buhid                  = _Buhid                  // Buhid 为书写系统 Buhid 中的 Unicode 字符集合。
	Canadian_Aboriginal    = _Canadian_Aboriginal    // Canadian_Aboriginal 为书写系统 Canadian_Aboriginal 中的 Unicode 字符集合。
	Carian                 = _Carian                 // Carian 为书写系统 Carian 中的 Unicode 字符集合。
	Caucasian_Albanian     = _Caucasian_Albanian     // Caucasian_Albanian 为书写系统 Caucasian_Albanian 中的 Unicode 字符集合。
	Chakma                 = _Chakma                 // Chakma 为书写系统 Chakma 中的 Unicode 字符集合。
	Cham                   = _Cham                   // Cham 为书写系统 Cham 中的 Unicode 字符集合。
	Cherokee               = _Cherokee               // Cherokee 为书写系统 Cherokee 中的 Unicode 字符集合。
	Common                 = _Common                 // Common 为书写系统 Common 中的 Unicode 字符集合。
	Coptic                 = _Coptic                 // Coptic 为书写系统 Coptic 中的 Unicode 字符集合。
	Cuneiform              = _Cuneiform              // Cuneiform 为书写系统 Cuneiform 中的 Unicode 字符集合。
	Cypriot                = _Cypriot                // Cypriot 为书写系统 Cypriot 中的 Unicode 字符集合。
	Cyrillic               = _Cyrillic               // Cyrillic 为书写系统 Cyrillic 中的 Unicode 字符集合。
	Deseret                = _Deseret                // Deseret 为书写系统 Deseret 中的 Unicode 字符集合。
	Devanagari             = _Devanagari             // Devanagari 为书写系统 Devanagari 中的 Unicode 字符集合。
	Duployan               = _Duployan               // Duployan 为书写系统 Duployan 中的 Unicode 字符集合。
	Egyptian_Hieroglyphs   = _Egyptian_Hieroglyphs   // Egyptian_Hieroglyphs 为书写系统 Egyptian_Hieroglyphs 中的 Unicode 字符集合。
	Elbasan                = _Elbasan                // Elbasan 为书写系统 Elbasan 中的 Unicode 字符集合。
	Ethiopic               = _Ethiopic               // Ethiopic 为书写系统 Ethiopic 中的 Unicode 字符集合。
	Georgian               = _Georgian               // Georgian 为书写系统 Georgian 中的 Unicode 字符集合。
	Glagolitic             = _Glagolitic             // Glagolitic 为书写系统 Glagolitic 中的 Unicode 字符集合。
	Gothic                 = _Gothic                 // Gothic 为书写系统 Gothic 中的 Unicode 字符集合。
	Grantha                = _Grantha                // Grantha 为书写系统 Grantha 中的 Unicode 字符集合。
	Greek                  = _Greek                  // Greek 为书写系统 Greek 中的 Unicode 字符集合。
	Gujarati               = _Gujarati               // Gujarati 为书写系统 Gujarati 中的 Unicode 字符集合。
	Gurmukhi               = _Gurmukhi               // Gurmukhi 为书写系统 Gurmukhi 中的 Unicode 字符集合。
	Han                    = _Han                    // Han 为书写系统 Han 中的 Unicode 字符集合。
	Hangul                 = _Hangul                 // Hangul 为书写系统 Hangul 中的 Unicode 字符集合。
	Hanunoo                = _Hanunoo                // Hanunoo 为书写系统 Hanunoo 中的 Unicode 字符集合。
	Hatran                 = _Hatran                 // Hatran 为书写系统 Hatran 中的 Unicode 字符集合。
	Hebrew                 = _Hebrew                 // Hebrew 为书写系统 Hebrew 中的 Unicode 字符集合。
	Hiragana               = _Hiragana               // Hiragana 为书写系统 Hiragana 中的 Unicode 字符集合。
	Imperial_Aramaic       = _Imperial_Aramaic       // Imperial_Aramaic 为书写系统 Imperial_Aramaic 中的 Unicode 字符集合。
	Inherited              = _Inherited              // Inherited 为书写系统 Inherited 中的 Unicode 字符集合。
	Inscriptional_Pahlavi  = _Inscriptional_Pahlavi  // Inscriptional_Pahlavi 为书写系统 Inscriptional_Pahlavi 中的 Unicode 字符集合。
	Inscriptional_Parthian = _Inscriptional_Parthian // Inscriptional_Parthian 为书写系统 Inscriptional_Parthian 中的 Unicode 字符集合。
	Javanese               = _Javanese               // Javanese 为书写系统 Javanese 中的 Unicode 字符集合。
	Kaithi                 = _Kaithi                 // Kaithi 为书写系统 Kaithi 中的 Unicode 字符集合。
	Kannada                = _Kannada                // Kannada 为书写系统 Kannada 中的 Unicode 字符集合。
	Katakana               = _Katakana               // Katakana 为书写系统 Katakana 中的 Unicode 字符集合。
	Kayah_Li               = _Kayah_Li               // Kayah_Li 为书写系统 Kayah_Li 中的 Unicode 字符集合。
	Kharoshthi             = _Kharoshthi             // Kharoshthi 为书写系统 Kharoshthi 中的 Unicode 字符集合。
	Khmer                  = _Khmer                  // Khmer 为书写系统 Khmer 中的 Unicode 字符集合。
	Khojki                 = _Khojki                 // Khojki 为书写系统 Khojki 中的 Unicode 字符集合。
	Khudawadi              = _Khudawadi              // Khudawadi 为书写系统 Khudawadi 中的 Unicode 字符集合。
	Lao                    = _Lao                    // Lao 为书写系统 Lao 中的 Unicode 字符集合。
	Latin                  = _Latin                  // Latin 为书写系统 Latin 中的 Unicode 字符集合。
	Lepcha                 = _Lepcha                 // Lepcha 为书写系统 Lepcha 中的 Unicode 字符集合。
	Limbu                  = _Limbu                  // Limbu 为书写系统 Limbu 中的 Unicode 字符集合。
	Linear_A               = _Linear_A               // Linear_A 为书写系统 Linear_A 中的 Unicode 字符集合。
	Linear_B               = _Linear_B               // Linear_B 为书写系统 Linear_B 中的 Unicode 字符集合。
	Lisu                   = _Lisu                   // Lisu 为书写系统 Lisu 中的 Unicode 字符集合。
	Lycian                 = _Lycian                 // Lycian 为书写系统 Lycian 中的 Unicode 字符集合。
	Lydian                 = _Lydian                 // Lydian 为书写系统 Lydian 中的 Unicode 字符集合。
	Mahajani               = _Mahajani               // Mahajani 为书写系统 Mahajani 中的 Unicode 字符集合。
	Malayalam              = _Malayalam              // Malayalam 为书写系统 Malayalam 中的 Unicode 字符集合。
	Mandaic                = _Mandaic                // Mandaic 为书写系统 Mandaic 中的 Unicode 字符集合。
	Manichaean             = _Manichaean             // Manichaean 为书写系统 Manichaean 中的 Unicode 字符集合。
	Meetei_Mayek           = _Meetei_Mayek           // Meetei_Mayek 为书写系统 Meetei_Mayek 中的 Unicode 字符集合。
	Mende_Kikakui          = _Mende_Kikakui          // Mende_Kikakui 为书写系统 Mende_Kikakui 中的 Unicode 字符集合。
	Meroitic_Cursive       = _Meroitic_Cursive       // Meroitic_Cursive 为书写系统 Meroitic_Cursive 中的 Unicode 字符集合。
	Meroitic_Hieroglyphs   = _Meroitic_Hieroglyphs   // Meroitic_Hieroglyphs 为书写系统 Meroitic_Hieroglyphs 中的 Unicode 字符集合。
	Miao                   = _Miao                   // Miao 为书写系统 Miao 中的 Unicode 字符集合。
	Modi                   = _Modi                   // Modi 为书写系统 Modi 中的 Unicode 字符集合。
	Mongolian              = _Mongolian              // Mongolian 为书写系统 Mongolian 中的 Unicode 字符集合。
	Mro                    = _Mro                    // Mro 为书写系统 Mro 中的 Unicode 字符集合。
	Multani                = _Multani                // Multani 为书写系统 Multani 中的 Unicode 字符集合。
	Myanmar                = _Myanmar                // Myanmar 为书写系统 Myanmar 中的 Unicode 字符集合。
	Nabataean              = _Nabataean              // Nabataean 为书写系统 Nabataean 中的 Unicode 字符集合。
	New_Tai_Lue            = _New_Tai_Lue            // New_Tai_Lue 为书写系统 New_Tai_Lue 中的 Unicode 字符集合。
	Nko                    = _Nko                    // Nko 为书写系统 Nko 中的 Unicode 字符集合。
	Ogham                  = _Ogham                  // Ogham 为书写系统 Ogham 中的 Unicode 字符集合。
	Ol_Chiki               = _Ol_Chiki               // Ol_Chiki 为书写系统 Ol_Chiki 中的 Unicode 字符集合。
	Old_Hungarian          = _Old_Hungarian          // Old_Hungarian 为书写系统 Old_Hungarian 中的 Unicode 字符集合。
	Old_Italic             = _Old_Italic             // Old_Italic 为书写系统 Old_Italic 中的 Unicode 字符集合。
	Old_North_Arabian      = _Old_North_Arabian      // Old_North_Arabian 为书写系统 Old_North_Arabian 中的 Unicode 字符集合。
	Old_Permic             = _Old_Permic             // Old_Permic 为书写系统 Old_Permic 中的 Unicode 字符集合。
	Old_Persian            = _Old_Persian            // Old_Persian 为书写系统 Old_Persian 中的 Unicode 字符集合。
	Old_South_Arabian      = _Old_South_Arabian      // Old_South_Arabian 为书写系统 Old_South_Arabian 中的 Unicode 字符集合。
	Old_Turkic             = _Old_Turkic             // Old_Turkic 为书写系统 Old_Turkic 中的 Unicode 字符集合。
	Oriya                  = _Oriya                  // Oriya 为书写系统 Oriya 中的 Unicode 字符集合。
	Osmanya                = _Osmanya                // Osmanya 为书写系统 Osmanya 中的 Unicode 字符集合。
	Pahawh_Hmong           = _Pahawh_Hmong           // Pahawh_Hmong 为书写系统 Pahawh_Hmong 中的 Unicode 字符集合。
	Palmyrene              = _Palmyrene              // Palmyrene 为书写系统 Palmyrene 中的 Unicode 字符集合。
	Pau_Cin_Hau            = _Pau_Cin_Hau            // Pau_Cin_Hau 为书写系统 Pau_Cin_Hau 中的 Unicode 字符集合。
	Phags_Pa               = _Phags_Pa               // Phags_Pa 为书写系统 Phags_Pa 中的 Unicode 字符集合。
	Phoenician             = _Phoenician             // Phoenician 为书写系统 Phoenician 中的 Unicode 字符集合。
	Psalter_Pahlavi        = _Psalter_Pahlavi        // Psalter_Pahlavi 为书写系统 Psalter_Pahlavi 中的 Unicode 字符集合。
	Rejang                 = _Rejang                 // Rejang 为书写系统 Rejang 中的 Unicode 字符集合。
	Runic                  = _Runic                  // Runic 为书写系统 Runic 中的 Unicode 字符集合。
	Samaritan              = _Samaritan              // Samaritan 为书写系统 Samaritan 中的 Unicode 字符集合。
	Saurashtra             = _Saurashtra             // Saurashtra 为书写系统 Saurashtra 中的 Unicode 字符集合。
	Sharada                = _Sharada                // Sharada 为书写系统 Sharada 中的 Unicode 字符集合。
	Shavian                = _Shavian                // Shavian 为书写系统 Shavian 中的 Unicode 字符集合。
	Siddham                = _Siddham                // Siddham 为书写系统 Siddham 中的 Unicode 字符集合。
	SignWriting            = _SignWriting            // SignWriting 为书写系统 SignWriting 中的 Unicode 字符集合。
	Sinhala                = _Sinhala                // Sinhala 为书写系统 Sinhala 中的 Unicode 字符集合。
	Sora_Sompeng           = _Sora_Sompeng           // Sora_Sompeng 为书写系统 Sora_Sompeng 中的 Unicode 字符集合。
	Sundanese              = _Sundanese              // Sundanese 为书写系统 Sundanese 中的 Unicode 字符集合。
	Syloti_Nagri           = _Syloti_Nagri           // Syloti_Nagri 为书写系统 Syloti_Nagri 中的 Unicode 字符集合。
	Syriac                 = _Syriac                 // Syriac 为书写系统 Syriac 中的 Unicode 字符集合。
	Tagalog                = _Tagalog                // Tagalog 为书写系统 Tagalog 中的 Unicode 字符集合。
	Tagbanwa               = _Tagbanwa               // Tagbanwa 为书写系统 Tagbanwa 中的 Unicode 字符集合。
	Tai_Le                 = _Tai_Le                 // Tai_Le 为书写系统 Tai_Le 中的 Unicode 字符集合。
	Tai_Tham               = _Tai_Tham               // Tai_Tham 为书写系统 Tai_Tham 中的 Unicode 字符集合。
	Tai_Viet               = _Tai_Viet               // Tai_Viet 为书写系统 Tai_Viet 中的 Unicode 字符集合。
	Takri                  = _Takri                  // Takri 为书写系统 Takri 中的 Unicode 字符集合。
	Tamil                  = _Tamil                  // Tamil 为书写系统 Tamil 中的 Unicode 字符集合。
	Telugu                 = _Telugu                 // Telugu 为书写系统 Telugu 中的 Unicode 字符集合。
	Thaana                 = _Thaana                 // Thaana 为书写系统 Thaana 中的 Unicode 字符集合。
	Thai                   = _Thai                   // Thai 为书写系统 Thai 中的 Unicode 字符集合。
	Tibetan                = _Tibetan                // Tibetan 为书写系统 Tibetan 中的 Unicode 字符集合。
	Tifinagh               = _Tifinagh               // Tifinagh 为书写系统 Tifinagh 中的 Unicode 字符集合。
	Tirhuta                = _Tirhuta                // Tirhuta 为书写系统 Tirhuta 中的 Unicode 字符集合。
	Ugaritic               = _Ugaritic               // Ugaritic 为书写系统 Ugaritic 中的 Unicode 字符集合。
	Vai                    = _Vai                    // Vai 为书写系统 Vai 中的 Unicode 字符集合。
	Warang_Citi            = _Warang_Citi            // Warang_Citi 为书写系统 Warang_Citi 中的 Unicode 字符集合。
	Yi                     = _Yi                     // Yi 为书写系统 Yi 中的 Unicode 字符集合。

)



var AzeriCase SpecialCase = _TurkishCase


// CaseRanges is the table describing case mappings for all letters with
// non-self mappings.

// CaseRanges 是描述所有“非自映射字母”的写法映射表。
var CaseRanges = _CaseRanges


// Categories is the set of Unicode category tables.

// Categories 为 Unicode 类别表的集合。
var Categories = map[string]*RangeTable{
	"C":  C,
	"Cc": Cc,
	"Cf": Cf,
	"Co": Co,
	"Cs": Cs,
	"L":  L,
	"Ll": Ll,
	"Lm": Lm,
	"Lo": Lo,
	"Lt": Lt,
	"Lu": Lu,
	"M":  M,
	"Mc": Mc,
	"Me": Me,
	"Mn": Mn,
	"N":  N,
	"Nd": Nd,
	"Nl": Nl,
	"No": No,
	"P":  P,
	"Pc": Pc,
	"Pd": Pd,
	"Pe": Pe,
	"Pf": Pf,
	"Pi": Pi,
	"Po": Po,
	"Ps": Ps,
	"S":  S,
	"Sc": Sc,
	"Sk": Sk,
	"Sm": Sm,
	"So": So,
	"Z":  Z,
	"Zl": Zl,
	"Zp": Zp,
	"Zs": Zs,
}


// These variables have type *RangeTable.

// These variables have type *RangeTable.
// 这些变量的类型为 *RangeTable。
var (
	Cc     = _Cc // Cc 为类别 Cc 中的 Unicode 字符集合。
	Cf     = _Cf // Cf 为类别 Cf 中的 Unicode 字符集合。
	Co     = _Co // Co 为类别 Co 中的 Unicode 字符集合。
	Cs     = _Cs // Cs 为类别 Cs 中的 Unicode 字符集合。
	Digit  = _Nd // Digit 为带属性“十进制数字”的 Unicode 字符集合。
	Nd     = _Nd // Nd 为类别 Nd 中的 Unicode 字符集合。
	Letter = _L  // Letter/L 为类别 L 中的 Unicode 字母字符集合。
	L      = _L
	Lm     = _Lm // Lm 为类别 Lm 中的 Unicode 字符集合。
	Lo     = _Lo // Lo 为类别 Lo 中的 Unicode 字符集合。
	Lower  = _Ll // Lower 为 Unicode 小写字母集合。
	Ll     = _Ll // Ll 为类别 Ll 中的 Unicode 字符集合。
	Mark   = _M  // Mark/M 为类别 M 中的 Unicode 标记字符集合。
	M      = _M
	Mc     = _Mc // Mc 为类别 Mc 中的 Unicode 字符集合。
	Me     = _Me // Me 为类别 Me 中的 Unicode 字符集合。
	Mn     = _Mn // Mn 为类别 Mn 中的 Unicode 字符集合。
	Nl     = _Nl // Nl 为类别 Nl 中的 Unicode 字符集合。
	No     = _No // No 为类别 No 中的 Unicode 字符集合。
	Number = _N  // Number/N 为类别 N 中的 Unicode 数字字符集合。
	N      = _N
	Other  = _C // Other/C 为类别 C 中的 Unicode 控制和特殊字符集合。
	C      = _C
	Pc     = _Pc // Pc 为类别 Pc 中的 Unicode 字符集合。
	Pd     = _Pd // Pd 为类别 Pd 中的 Unicode 字符集合。
	Pe     = _Pe // Pe 为类别 Pe 中的 Unicode 字符集合。
	Pf     = _Pf // Pf 为类别 Pf 中的 Unicode 字符集合。
	Pi     = _Pi // Pi 为类别 Pi 中的 Unicode 字符集合。
	Po     = _Po // Po 为类别 Po 中的 Unicode 字符集合。
	Ps     = _Ps // Ps 为类别 Ps 中的 Unicode 字符集合。
	Punct  = _P  // Punct/P 为类别 P 中的 Unicode 标点字符集合。
	P      = _P
	Sc     = _Sc // Sc 为类别 Sc 中的 Unicode 字符集合。
	Sk     = _Sk // Sk 为类别 Sk 中的 Unicode 字符集合。
	Sm     = _Sm // Sm 为类别 Sm 中的 Unicode 字符集合。
	So     = _So // So 为类别 So 中的 Unicode 字符集合。
	Space  = _Z  // Space/Z 为类别 Z 中的 Unicode 空白字符集合。
	Z      = _Z
	Symbol = _S // Symbol/S 为类别 S 中的 Unicode 符号字符集合。
	S      = _S
	Title  = _Lt // Title 为 Unicode 标题字母集合。
	Lt     = _Lt // Lt 为类别 Lt 中的 Unicode 字符集合。
	Upper  = _Lu // Upper 为 Unicode 大写字母集合。
	Lu     = _Lu // Lu 为类别 Lu 中的 Unicode 字符集合。
	Zl     = _Zl // Zl 为类别 Zl 中的 Unicode 字符集合。
	Zp     = _Zp // Zp 为类别 Zp 中的 Unicode 字符集合。
	Zs     = _Zs // Zs 为类别 Zs 中的 Unicode 字符集合。

)


// FoldCategory maps a category name to a table of
// code points outside the category that are equivalent under
// simple case folding to code points inside the category.
// If there is no entry for a category name, there are no such points.

// FoldCategory 将一个类别名映射到该类别外的码点表上，
// 这相当于在简单的情况下对该类别内的码点进行转换。
// 若一个类别名没有对应的条目，则该码点不存在。
var FoldCategory = map[string]*RangeTable{
	"Common":    foldCommon,
	"Greek":     foldGreek,
	"Inherited": foldInherited,
	"L":         foldL,
	"Ll":        foldLl,
	"Lt":        foldLt,
	"Lu":        foldLu,
	"M":         foldM,
	"Mn":        foldMn,
}


// FoldScript maps a script name to a table of
// code points outside the script that are equivalent under
// simple case folding to code points inside the script.
// If there is no entry for a script name, there are no such points.

// FoldCategory 将一个书写系统名映射到该书写系统外的码点表上，
// 这相当于在简单的情况下对该书写系统内的码点进行转换。
// 若一个书写系统名没有对应的条目，则该码点不存在。
var FoldScript = map[string]*RangeTable{}


// GraphicRanges defines the set of graphic characters according to Unicode.

// GraphicRanges 根据Unicode定义了可显示字符的集合。
var GraphicRanges = []*RangeTable{
	L, M, N, P, S, Zs,
}


// PrintRanges defines the set of printable characters according to Go.
// ASCII space, U+0020, is handled separately.

// PrintRanges 根据Go定义了可打印字符的集合。ASCII空格（即U+0020）另作处理。
var PrintRanges = []*RangeTable{
	L, M, N, P, S,
}


// Properties is the set of Unicode property tables.

// Properties 为 Unicode 属性表的集合。
var Properties = map[string]*RangeTable{
	"ASCII_Hex_Digit":                    ASCII_Hex_Digit,
	"Bidi_Control":                       Bidi_Control,
	"Dash":                               Dash,
	"Deprecated":                         Deprecated,
	"Diacritic":                          Diacritic,
	"Extender":                           Extender,
	"Hex_Digit":                          Hex_Digit,
	"Hyphen":                             Hyphen,
	"IDS_Binary_Operator":                IDS_Binary_Operator,
	"IDS_Trinary_Operator":               IDS_Trinary_Operator,
	"Ideographic":                        Ideographic,
	"Join_Control":                       Join_Control,
	"Logical_Order_Exception":            Logical_Order_Exception,
	"Noncharacter_Code_Point":            Noncharacter_Code_Point,
	"Other_Alphabetic":                   Other_Alphabetic,
	"Other_Default_Ignorable_Code_Point": Other_Default_Ignorable_Code_Point,
	"Other_Grapheme_Extend":              Other_Grapheme_Extend,
	"Other_ID_Continue":                  Other_ID_Continue,
	"Other_ID_Start":                     Other_ID_Start,
	"Other_Lowercase":                    Other_Lowercase,
	"Other_Math":                         Other_Math,
	"Other_Uppercase":                    Other_Uppercase,
	"Pattern_Syntax":                     Pattern_Syntax,
	"Pattern_White_Space":                Pattern_White_Space,
	"Quotation_Mark":                     Quotation_Mark,
	"Radical":                            Radical,
	"STerm":                              STerm,
	"Soft_Dotted":                        Soft_Dotted,
	"Terminal_Punctuation":               Terminal_Punctuation,
	"Unified_Ideograph":                  Unified_Ideograph,
	"Variation_Selector":                 Variation_Selector,
	"White_Space":                        White_Space,
}


// Scripts is the set of Unicode script tables.

// Scripts 为 Unicode 书写表的集合。
var Scripts = map[string]*RangeTable{
	"Ahom":                   Ahom,
	"Anatolian_Hieroglyphs":  Anatolian_Hieroglyphs,
	"Arabic":                 Arabic,
	"Armenian":               Armenian,
	"Avestan":                Avestan,
	"Balinese":               Balinese,
	"Bamum":                  Bamum,
	"Bassa_Vah":              Bassa_Vah,
	"Batak":                  Batak,
	"Bengali":                Bengali,
	"Bopomofo":               Bopomofo,
	"Brahmi":                 Brahmi,
	"Braille":                Braille,
	"Buginese":               Buginese,
	"Buhid":                  Buhid,
	"Canadian_Aboriginal":    Canadian_Aboriginal,
	"Carian":                 Carian,
	"Caucasian_Albanian":     Caucasian_Albanian,
	"Chakma":                 Chakma,
	"Cham":                   Cham,
	"Cherokee":               Cherokee,
	"Common":                 Common,
	"Coptic":                 Coptic,
	"Cuneiform":              Cuneiform,
	"Cypriot":                Cypriot,
	"Cyrillic":               Cyrillic,
	"Deseret":                Deseret,
	"Devanagari":             Devanagari,
	"Duployan":               Duployan,
	"Egyptian_Hieroglyphs":   Egyptian_Hieroglyphs,
	"Elbasan":                Elbasan,
	"Ethiopic":               Ethiopic,
	"Georgian":               Georgian,
	"Glagolitic":             Glagolitic,
	"Gothic":                 Gothic,
	"Grantha":                Grantha,
	"Greek":                  Greek,
	"Gujarati":               Gujarati,
	"Gurmukhi":               Gurmukhi,
	"Han":                    Han,
	"Hangul":                 Hangul,
	"Hanunoo":                Hanunoo,
	"Hatran":                 Hatran,
	"Hebrew":                 Hebrew,
	"Hiragana":               Hiragana,
	"Imperial_Aramaic":       Imperial_Aramaic,
	"Inherited":              Inherited,
	"Inscriptional_Pahlavi":  Inscriptional_Pahlavi,
	"Inscriptional_Parthian": Inscriptional_Parthian,
	"Javanese":               Javanese,
	"Kaithi":                 Kaithi,
	"Kannada":                Kannada,
	"Katakana":               Katakana,
	"Kayah_Li":               Kayah_Li,
	"Kharoshthi":             Kharoshthi,
	"Khmer":                  Khmer,
	"Khojki":                 Khojki,
	"Khudawadi":              Khudawadi,
	"Lao":                    Lao,
	"Latin":                  Latin,
	"Lepcha":                 Lepcha,
	"Limbu":                  Limbu,
	"Linear_A":               Linear_A,
	"Linear_B":               Linear_B,
	"Lisu":                   Lisu,
	"Lycian":                 Lycian,
	"Lydian":                 Lydian,
	"Mahajani":               Mahajani,
	"Malayalam":              Malayalam,
	"Mandaic":                Mandaic,
	"Manichaean":             Manichaean,
	"Meetei_Mayek":           Meetei_Mayek,
	"Mende_Kikakui":          Mende_Kikakui,
	"Meroitic_Cursive":       Meroitic_Cursive,
	"Meroitic_Hieroglyphs":   Meroitic_Hieroglyphs,
	"Miao":                   Miao,
	"Modi":                   Modi,
	"Mongolian":              Mongolian,
	"Mro":                    Mro,
	"Multani":                Multani,
	"Myanmar":                Myanmar,
	"Nabataean":              Nabataean,
	"New_Tai_Lue":            New_Tai_Lue,
	"Nko":                    Nko,
	"Ogham":                  Ogham,
	"Ol_Chiki":               Ol_Chiki,
	"Old_Hungarian":          Old_Hungarian,
	"Old_Italic":             Old_Italic,
	"Old_North_Arabian":      Old_North_Arabian,
	"Old_Permic":             Old_Permic,
	"Old_Persian":            Old_Persian,
	"Old_South_Arabian":      Old_South_Arabian,
	"Old_Turkic":             Old_Turkic,
	"Oriya":                  Oriya,
	"Osmanya":                Osmanya,
	"Pahawh_Hmong":           Pahawh_Hmong,
	"Palmyrene":              Palmyrene,
	"Pau_Cin_Hau":            Pau_Cin_Hau,
	"Phags_Pa":               Phags_Pa,
	"Phoenician":             Phoenician,
	"Psalter_Pahlavi":        Psalter_Pahlavi,
	"Rejang":                 Rejang,
	"Runic":                  Runic,
	"Samaritan":              Samaritan,
	"Saurashtra":             Saurashtra,
	"Sharada":                Sharada,
	"Shavian":                Shavian,
	"Siddham":                Siddham,
	"SignWriting":            SignWriting,
	"Sinhala":                Sinhala,
	"Sora_Sompeng":           Sora_Sompeng,
	"Sundanese":              Sundanese,
	"Syloti_Nagri":           Syloti_Nagri,
	"Syriac":                 Syriac,
	"Tagalog":                Tagalog,
	"Tagbanwa":               Tagbanwa,
	"Tai_Le":                 Tai_Le,
	"Tai_Tham":               Tai_Tham,
	"Tai_Viet":               Tai_Viet,
	"Takri":                  Takri,
	"Tamil":                  Tamil,
	"Telugu":                 Telugu,
	"Thaana":                 Thaana,
	"Thai":                   Thai,
	"Tibetan":                Tibetan,
	"Tifinagh":               Tifinagh,
	"Tirhuta":                Tirhuta,
	"Ugaritic":               Ugaritic,
	"Vai":                    Vai,
	"Warang_Citi":            Warang_Citi,
	"Yi":                     Yi,
}



var TurkishCase SpecialCase = _TurkishCase


// CaseRange represents a range of Unicode code points for simple (one
// code point to one code point) case conversion.
// The range runs from Lo to Hi inclusive, with a fixed stride of 1.  Deltas
// are the number to add to the code point to reach the code point for a
// different case for that character.  They may be negative.  If zero, it
// means the character is in the corresponding case. There is a special
// case representing sequences of alternating corresponding Upper and Lower
// pairs.  It appears with a fixed Delta of
//     {UpperLower, UpperLower, UpperLower}
// The constant UpperLower has an otherwise impossible delta value.

// CaseRange 表示Unicode码点中，简单的（即一对一的）大小写转换的范围。该范围从
// Lo 连续到 Hi，包括一个固定的间距。Delta 为添加的码点数量，
// 以便于该字符不同写法间的转换。它们可为负数。若为零，即表示该字符的写法一致。
// 还有种特殊的写法，表示一对大小写交替对应的序列。它会与像
//     {UpperLower, UpperLower, UpperLower}
// 这样固定的 Delta 一同出现。常量 UpperLower 可能拥有其它的 delta 值。
type CaseRange struct {
	Lo    uint32
	Hi    uint32
	Delta d
}


// Range16 represents of a range of 16-bit Unicode code points. The range runs
// from Lo to Hi inclusive and has the specified stride.

// Range16 表示16位Unicode码点的范围。该范围从 Lo 连续到 Hi 且包括两端，
// 还有一个指定的间距。
type Range16 struct {
	Lo     uint16
	Hi     uint16
	Stride uint16
}


// Range32 represents of a range of Unicode code points and is used when one or
// more of the values will not fit in 16 bits.  The range runs from Lo to Hi
// inclusive and has the specified stride. Lo and Hi must always be >= 1<<16.

// Range32 表示Unicode码点的范围，它在一个或多个值不能用16位容纳时使用。该范围从
// Lo 连续到 Hi 且包括两端，还有一个指定的间距。Lo 和 Hi 都必须满足 >= 1<<16。
type Range32 struct {
	Lo     uint32
	Hi     uint32
	Stride uint32
}


// RangeTable defines a set of Unicode code points by listing the ranges of
// code points within the set. The ranges are listed in two slices
// to save space: a slice of 16-bit ranges and a slice of 32-bit ranges.
// The two slices must be in sorted order and non-overlapping.
// Also, R32 should contain only values >= 0x10000 (1<<16).

// RangeTable 通过列出码点范围，定义了Unicode码点的集合。为了节省空间， 其范围分
// 别在16位、32位这两个切片中列出。这两个切片必须已经排序且无重叠的部分。 此外，
// R32只包含 >= 0x10000 (1<<16) 的值。
type RangeTable struct {
	R16         []Range16
	R32         []Range32
	LatinOffset int // number of entries in R16 with Hi <= MaxLatin1 // R16 中满足 Hi <= MaxLatin1 的条目数
}


// SpecialCase represents language-specific case mappings such as Turkish.
// Methods of SpecialCase customize (by overriding) the standard mappings.

// SpecialCase 表示语言相关的写法映射，例如土耳其语。SpecialCase 的方法（通过覆
// 盖） 来定制标准的映射。
type SpecialCase []CaseRange


// In reports whether the rune is a member of one of the ranges.

// In 报告该符文是否为该范围中的一员。
func In(r rune, ranges ...*RangeTable) bool

// Is reports whether the rune is in the specified table of ranges.

// Is 报告该符文是否在指定范围的表中。
func Is(rangeTab *RangeTable, r rune) bool

// IsControl reports whether the rune is a control character.
// The C (Other) Unicode category includes more code points
// such as surrogates; use Is(C, r) to test for them.

// IsControl 报告该字符是否为控制字符。Unicode的C（其它）
// 类别包括了更多像替代值这样的码点；请使用 Is(C, r) 来测试它们。
func IsControl(r rune) bool

// IsDigit reports whether the rune is a decimal digit.

// IsDigit 报告该符文是否为十进制数字。
func IsDigit(r rune) bool

// IsGraphic reports whether the rune is defined as a Graphic by Unicode.
// Such characters include letters, marks, numbers, punctuation, symbols, and
// spaces, from categories L, M, N, P, S, Zs.

// IsGraphic 报告该符文是否为Unicode定义的可显示字符。包括字母、标记、数字、
// 标点、符号和空白这样的，类别为L、M、N、P,、S和Zs的字符。
func IsGraphic(r rune) bool

// IsLetter reports whether the rune is a letter (category L).

// IsLetter 报告该符文是否为字母（类别L）。
func IsLetter(r rune) bool

// IsLower reports whether the rune is a lower case letter.

// IsLower 报告该符文是否为小写字母。
func IsLower(r rune) bool

// IsMark reports whether the rune is a mark character (category M).

// IsMark 报告该符文是否为标记字符（类别M）。
func IsMark(r rune) bool

// IsNumber reports whether the rune is a number (category N).

// IsNumber 报告该符文是否为数字（类别N）。
func IsNumber(r rune) bool

// IsOneOf reports whether the rune is a member of one of the ranges. The
// function "In" provides a nicer signature and should be used in preference to
// IsOneOf.

// IsOneOf 报告该符文是否为该范围中的一员。
// 函数“In”提供了一个更好的签名，比起 IsOneOf 来我们更倾向于使用它。
func IsOneOf(ranges []*RangeTable, r rune) bool

// IsPrint reports whether the rune is defined as printable by Go. Such
// characters include letters, marks, numbers, punctuation, symbols, and the
// ASCII space character, from categories L, M, N, P, S and the ASCII space
// character.  This categorization is the same as IsGraphic except that the
// only spacing character is ASCII space, U+0020.

// IsPrint 报告该符文是否为Go定义的可打印字符。包括字母、标记、数字、标点、
// 符号和ASCII空格这样的，类别为L、M、N、P、S和ASCII空格的字符。
// 除空白字符只有ASCII空格（即U+0020）外，其它的类别与 IsGraphic 相同。
func IsPrint(r rune) bool

// IsPunct reports whether the rune is a Unicode punctuation character
// (category P).

// IsPunct 报告该符文是否为Unicode标点字符（类别P）。
func IsPunct(r rune) bool

// IsSpace reports whether the rune is a space character as defined
// by Unicode's White Space property; in the Latin-1 space
// this is
//     '\t', '\n', '\v', '\f', '\r', ' ', U+0085 (NEL), U+00A0 (NBSP).
// Other definitions of spacing characters are set by category
// Z and property Pattern_White_Space.

// IsSpace 报告该符文是否为Unicode空白字符属性定义的空白符；在Latin-1中的空白为
//     '\t'、'\n'、'\v'、'\f'、'\r'、' '、U+0085 (NEL) 和 U+00A0 (NBSP)。
// 其它空白字符的定义由类别Z和属性 Pattern_White_Space 设置。
func IsSpace(r rune) bool

// IsSymbol reports whether the rune is a symbolic character.

// IsSymbol 报告该符文是否为符号字符。
func IsSymbol(r rune) bool

// IsTitle reports whether the rune is a title case letter.

// IsTitle 报告该符文是否为标题字母。
func IsTitle(r rune) bool

// IsUpper reports whether the rune is an upper case letter.

// IsUpper 报告该符文是否为大写字母。
func IsUpper(r rune) bool

// SimpleFold iterates over Unicode code points equivalent under
// the Unicode-defined simple case folding.  Among the code points
// equivalent to rune (including rune itself), SimpleFold returns the
// smallest rune > r if one exists, or else the smallest rune >= 0.
//
// For example:
//     SimpleFold('A') = 'a'
//     SimpleFold('a') = 'A'
//
//     SimpleFold('K') = 'k'
//     SimpleFold('k') = '\u212A' (Kelvin symbol, K)
//     SimpleFold('\u212A') = 'K'
//
//     SimpleFold('1') = '1'

// SimpleFold 遍历Unicode码点，等价于Unicode定义下的简单写法转换。
// 其中的码点等价于符文（包括符文自身），若存在最小的 >= r 的符文，SimpleFold
// 返回就会返回它，否则就会返回最小的 >= 0 的符文。
//
// 例如：
//     SimpleFold('A') = 'a'
//     SimpleFold('a') = 'A'
//
//     SimpleFold('K') = 'k'
//     SimpleFold('k') = '\u212A' （开尔文符号，K)
//     SimpleFold('\u212A') = 'K'
//
//     SimpleFold('1') = '1'
func SimpleFold(r rune) rune

// To maps the rune to the specified case: UpperCase, LowerCase, or TitleCase.

// To 将该符文映射为指定的写法：UpperCase、LowerCase、或 TitleCase。
func To(_case int, r rune) rune

// ToLower maps the rune to lower case.

// ToUpper 将该符文映射为小写形式。
func ToLower(r rune) rune

// ToTitle maps the rune to title case.

// ToTitle 将该符文映射为标题形式。
func ToTitle(r rune) rune

// ToUpper maps the rune to upper case.

// ToUpper 将该符文映射为大写形式。
func ToUpper(r rune) rune

// ToLower maps the rune to lower case giving priority to the special mapping.

// ToLower 将该符文映射为大写形式，优先考虑特殊的映射。
func (SpecialCase) ToLower(r rune) rune

// ToTitle maps the rune to title case giving priority to the special mapping.

// ToTitle 将该符文映射为标题形式，优先考虑特殊的映射。
func (SpecialCase) ToTitle(r rune) rune

// ToUpper maps the rune to upper case giving priority to the special mapping.

// ToUpper 将该符文映射为大写形式，优先考虑特殊的映射。
func (SpecialCase) ToUpper(r rune) rune

