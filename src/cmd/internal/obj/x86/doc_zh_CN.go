// +build ingore

package x86

import (
    "cmd/internal/obj"
    "cmd/internal/sys"
    "encoding/binary"
    "fmt"
    "log"
    "math"
    "strings"
)

//  *	amd64

// *	amd64
const (
	AAAA = obj.ABaseAMD64 + obj.A_ARCHSPECIFIC + iota
	AAAD
	AAAM
	AAAS
	AADCB
	AADCL
	AADCW
	AADDB
	AADDL
	AADDW
	AADJSP
	AANDB
	AANDL
	AANDW
	AARPL
	ABOUNDL
	ABOUNDW
	ABSFL
	ABSFW
	ABSRL
	ABSRW
	ABTL
	ABTW
	ABTCL
	ABTCW
	ABTRL
	ABTRW
	ABTSL
	ABTSW
	ABYTE
	ACLC
	ACLD
	ACLI
	ACLTS
	ACMC
	ACMPB
	ACMPL
	ACMPW
	ACMPSB
	ACMPSL
	ACMPSW
	ADAA
	ADAS
	ADECB
	ADECL
	ADECQ
	ADECW
	ADIVB
	ADIVL
	ADIVW
	AENTER
	AHADDPD
	AHADDPS
	AHLT
	AHSUBPD
	AHSUBPS
	AIDIVB
	AIDIVL
	AIDIVW
	AIMULB
	AIMULL
	AIMULW
	AINB
	AINL
	AINW
	AINCB
	AINCL
	AINCQ
	AINCW
	AINSB
	AINSL
	AINSW
	AINT
	AINTO
	AIRETL
	AIRETW
	AJCC // >= unsigned
	AJCS // < unsigned
	AJCXZL
	AJEQ // == (zero)
	AJGE // >= signed
	AJGT // > signed
	AJHI // > unsigned
	AJLE // <= signed
	AJLS // <= unsigned
	AJLT // < signed
	AJMI // sign bit set (negative)
	AJNE // != (nonzero)
	AJOC // overflow clear
	AJOS // overflow set
	AJPC // parity clear
	AJPL // sign bit clear (positive)
	AJPS // parity set
	ALAHF
	ALARL
	ALARW
	ALEAL
	ALEAW
	ALEAVEL
	ALEAVEW
	ALOCK
	ALODSB
	ALODSL
	ALODSW
	ALONG
	ALOOP
	ALOOPEQ
	ALOOPNE
	ALSLL
	ALSLW
	AMOVB
	AMOVL
	AMOVW
	AMOVBLSX
	AMOVBLZX
	AMOVBQSX
	AMOVBQZX
	AMOVBWSX
	AMOVBWZX
	AMOVWLSX
	AMOVWLZX
	AMOVWQSX
	AMOVWQZX
	AMOVSB
	AMOVSL
	AMOVSW
	AMULB
	AMULL
	AMULW
	ANEGB
	ANEGL
	ANEGW
	ANOTB
	ANOTL
	ANOTW
	AORB
	AORL
	AORW
	AOUTB
	AOUTL
	AOUTW
	AOUTSB
	AOUTSL
	AOUTSW
	APAUSE
	APOPAL
	APOPAW
	APOPCNTW
	APOPCNTL
	APOPCNTQ
	APOPFL
	APOPFW
	APOPL
	APOPW
	APUSHAL
	APUSHAW
	APUSHFL
	APUSHFW
	APUSHL
	APUSHW
	ARCLB
	ARCLL
	ARCLW
	ARCRB
	ARCRL
	ARCRW
	AREP
	AREPN
	AROLB
	AROLL
	AROLW
	ARORB
	ARORL
	ARORW
	ASAHF
	ASALB
	ASALL
	ASALW
	ASARB
	ASARL
	ASARW
	ASBBB
	ASBBL
	ASBBW
	ASCASB
	ASCASL
	ASCASW
	ASETCC
	ASETCS
	ASETEQ
	ASETGE
	ASETGT
	ASETHI
	ASETLE
	ASETLS
	ASETLT
	ASETMI
	ASETNE
	ASETOC
	ASETOS
	ASETPC
	ASETPL
	ASETPS
	ACDQ
	ACWD
	ASHLB
	ASHLL
	ASHLW
	ASHRB
	ASHRL
	ASHRW
	ASTC
	ASTD
	ASTI
	ASTOSB
	ASTOSL
	ASTOSW
	ASUBB
	ASUBL
	ASUBW
	ASYSCALL
	ATESTB
	ATESTL
	ATESTW
	AVERR
	AVERW
	AWAIT
	AWORD
	AXCHGB
	AXCHGL
	AXCHGW
	AXLAT
	AXORB
	AXORL
	AXORW
	AFMOVB
	AFMOVBP
	AFMOVD
	AFMOVDP
	AFMOVF
	AFMOVFP
	AFMOVL
	AFMOVLP
	AFMOVV
	AFMOVVP
	AFMOVW
	AFMOVWP
	AFMOVX
	AFMOVXP
	AFCOMD
	AFCOMDP
	AFCOMDPP
	AFCOMF
	AFCOMFP
	AFCOML
	AFCOMLP
	AFCOMW
	AFCOMWP
	AFUCOM
	AFUCOMP
	AFUCOMPP
	AFADDDP
	AFADDW
	AFADDL
	AFADDF
	AFADDD
	AFMULDP
	AFMULW
	AFMULL
	AFMULF
	AFMULD
	AFSUBDP
	AFSUBW
	AFSUBL
	AFSUBF
	AFSUBD
	AFSUBRDP
	AFSUBRW
	AFSUBRL
	AFSUBRF
	AFSUBRD
	AFDIVDP
	AFDIVW
	AFDIVL
	AFDIVF
	AFDIVD
	AFDIVRDP
	AFDIVRW
	AFDIVRL
	AFDIVRF
	AFDIVRD
	AFXCHD
	AFFREE
	AFLDCW
	AFLDENV
	AFRSTOR
	AFSAVE
	AFSTCW
	AFSTENV
	AFSTSW
	AF2XM1
	AFABS
	AFCHS
	AFCLEX
	AFCOS
	AFDECSTP
	AFINCSTP
	AFINIT
	AFLD1
	AFLDL2E
	AFLDL2T
	AFLDLG2
	AFLDLN2
	AFLDPI
	AFLDZ
	AFNOP
	AFPATAN
	AFPREM
	AFPREM1
	AFPTAN
	AFRNDINT
	AFSCALE
	AFSIN
	AFSINCOS
	AFSQRT
	AFTST
	AFXAM
	AFXTRACT
	AFYL2X
	AFYL2XP1
	// extra 32-bit operations
	ACMPXCHGB
	ACMPXCHGL
	ACMPXCHGW
	ACMPXCHG8B
	ACPUID
	AINVD
	AINVLPG
	ALFENCE
	AMFENCE
	AMOVNTIL
	ARDMSR
	ARDPMC
	ARDTSC
	ARSM
	ASFENCE
	ASYSRET
	AWBINVD
	AWRMSR
	AXADDB
	AXADDL
	AXADDW
	// conditional move
	ACMOVLCC
	ACMOVLCS
	ACMOVLEQ
	ACMOVLGE
	ACMOVLGT
	ACMOVLHI
	ACMOVLLE
	ACMOVLLS
	ACMOVLLT
	ACMOVLMI
	ACMOVLNE
	ACMOVLOC
	ACMOVLOS
	ACMOVLPC
	ACMOVLPL
	ACMOVLPS
	ACMOVQCC
	ACMOVQCS
	ACMOVQEQ
	ACMOVQGE
	ACMOVQGT
	ACMOVQHI
	ACMOVQLE
	ACMOVQLS
	ACMOVQLT
	ACMOVQMI
	ACMOVQNE
	ACMOVQOC
	ACMOVQOS
	ACMOVQPC
	ACMOVQPL
	ACMOVQPS
	ACMOVWCC
	ACMOVWCS
	ACMOVWEQ
	ACMOVWGE
	ACMOVWGT
	ACMOVWHI
	ACMOVWLE
	ACMOVWLS
	ACMOVWLT
	ACMOVWMI
	ACMOVWNE
	ACMOVWOC
	ACMOVWOS
	ACMOVWPC
	ACMOVWPL
	ACMOVWPS
	// 64-bit
	AADCQ
	AADDQ
	AANDQ
	ABSFQ
	ABSRQ
	ABTCQ
	ABTQ
	ABTRQ
	ABTSQ
	ACMPQ
	ACMPSQ
	ACMPXCHGQ
	ACQO
	ADIVQ
	AIDIVQ
	AIMULQ
	AIRETQ
	AJCXZQ
	ALEAQ
	ALEAVEQ
	ALODSQ
	AMOVQ
	AMOVLQSX
	AMOVLQZX
	AMOVNTIQ
	AMOVSQ
	AMULQ
	ANEGQ
	ANOTQ
	AORQ
	APOPFQ
	APOPQ
	APUSHFQ
	APUSHQ
	ARCLQ
	ARCRQ
	AROLQ
	ARORQ
	AQUAD
	ASALQ
	ASARQ
	ASBBQ
	ASCASQ
	ASHLQ
	ASHRQ
	ASTOSQ
	ASUBQ
	ATESTQ
	AXADDQ
	AXCHGQ
	AXORQ
	AXGETBV
	// media
	AADDPD
	AADDPS
	AADDSD
	AADDSS
	AANDNL
	AANDNQ
	AANDNPD
	AANDNPS
	AANDPD
	AANDPS
	ABEXTRL
	ABEXTRQ
	ABLSIL
	ABLSIQ
	ABLSMSKL
	ABLSMSKQ
	ABLSRL
	ABLSRQ
	ABZHIL
	ABZHIQ
	ACMPPD
	ACMPPS
	ACMPSD
	ACMPSS
	ACOMISD
	ACOMISS
	ACVTPD2PL
	ACVTPD2PS
	ACVTPL2PD
	ACVTPL2PS
	ACVTPS2PD
	ACVTPS2PL
	ACVTSD2SL
	ACVTSD2SQ
	ACVTSD2SS
	ACVTSL2SD
	ACVTSL2SS
	ACVTSQ2SD
	ACVTSQ2SS
	ACVTSS2SD
	ACVTSS2SL
	ACVTSS2SQ
	ACVTTPD2PL
	ACVTTPS2PL
	ACVTTSD2SL
	ACVTTSD2SQ
	ACVTTSS2SL
	ACVTTSS2SQ
	ADIVPD
	ADIVPS
	ADIVSD
	ADIVSS
	AEMMS
	AFXRSTOR
	AFXRSTOR64
	AFXSAVE
	AFXSAVE64
	ALDDQU
	ALDMXCSR
	AMASKMOVOU
	AMASKMOVQ
	AMAXPD
	AMAXPS
	AMAXSD
	AMAXSS
	AMINPD
	AMINPS
	AMINSD
	AMINSS
	AMOVAPD
	AMOVAPS
	AMOVOU
	AMOVHLPS
	AMOVHPD
	AMOVHPS
	AMOVLHPS
	AMOVLPD
	AMOVLPS
	AMOVMSKPD
	AMOVMSKPS
	AMOVNTO
	AMOVNTPD
	AMOVNTPS
	AMOVNTQ
	AMOVO
	AMOVQOZX
	AMOVSD
	AMOVSS
	AMOVUPD
	AMOVUPS
	AMULPD
	AMULPS
	AMULSD
	AMULSS
	AMULXL
	AMULXQ
	AORPD
	AORPS
	APACKSSLW
	APACKSSWB
	APACKUSWB
	APADDB
	APADDL
	APADDQ
	APADDSB
	APADDSW
	APADDUSB
	APADDUSW
	APADDW
	APAND
	APANDN
	APAVGB
	APAVGW
	APCMPEQB
	APCMPEQL
	APCMPEQW
	APCMPGTB
	APCMPGTL
	APCMPGTW
	APDEPL
	APDEPQ
	APEXTL
	APEXTQ
	APEXTRB
	APEXTRD
	APEXTRQ
	APEXTRW
	APHADDD
	APHADDSW
	APHADDW
	APHMINPOSUW
	APHSUBD
	APHSUBSW
	APHSUBW
	APINSRB
	APINSRD
	APINSRQ
	APINSRW
	APMADDWL
	APMAXSW
	APMAXUB
	APMINSW
	APMINUB
	APMOVMSKB
	APMOVSXBD
	APMOVSXBQ
	APMOVSXBW
	APMOVSXDQ
	APMOVSXWD
	APMOVSXWQ
	APMOVZXBD
	APMOVZXBQ
	APMOVZXBW
	APMOVZXDQ
	APMOVZXWD
	APMOVZXWQ
	APMULDQ
	APMULHUW
	APMULHW
	APMULLD
	APMULLW
	APMULULQ
	APOR
	APSADBW
	APSHUFB
	APSHUFHW
	APSHUFL
	APSHUFLW
	APSHUFW
	APSLLL
	APSLLO
	APSLLQ
	APSLLW
	APSRAL
	APSRAW
	APSRLL
	APSRLO
	APSRLQ
	APSRLW
	APSUBB
	APSUBL
	APSUBQ
	APSUBSB
	APSUBSW
	APSUBUSB
	APSUBUSW
	APSUBW
	APUNPCKHBW
	APUNPCKHLQ
	APUNPCKHQDQ
	APUNPCKHWL
	APUNPCKLBW
	APUNPCKLLQ
	APUNPCKLQDQ
	APUNPCKLWL
	APXOR
	ARCPPS
	ARCPSS
	ARSQRTPS
	ARSQRTSS
	ASARXL
	ASARXQ
	ASHLXL
	ASHLXQ
	ASHRXL
	ASHRXQ
	ASHUFPD
	ASHUFPS
	ASQRTPD
	ASQRTPS
	ASQRTSD
	ASQRTSS
	ASTMXCSR
	ASUBPD
	ASUBPS
	ASUBSD
	ASUBSS
	AUCOMISD
	AUCOMISS
	AUNPCKHPD
	AUNPCKHPS
	AUNPCKLPD
	AUNPCKLPS
	AXORPD
	AXORPS
	APCMPESTRI
	ARETFW
	ARETFL
	ARETFQ
	ASWAPGS
	AMODE
	ACRC32B
	ACRC32Q
	AIMUL3Q
	APREFETCHT0
	APREFETCHT1
	APREFETCHT2
	APREFETCHNTA
	AMOVQL
	ABSWAPL
	ABSWAPQ
	AAESENC
	AAESENCLAST
	AAESDEC
	AAESDECLAST
	AAESIMC
	AAESKEYGENASSIST
	AROUNDPS
	AROUNDSS
	AROUNDPD
	AROUNDSD
	APSHUFD
	APCLMULQDQ
	AVZEROUPPER
	AVMOVDQU
	AVMOVNTDQ
	AVMOVDQA
	AVPCMPEQB
	AVPXOR
	AVPMOVMSKB
	AVPAND
	AVPTEST
	AVPBROADCASTB
	AVPSHUFB
	AVPSHUFD
	AVPERM2F128
	AVPALIGNR
	AVPADDQ
	AVPADDD
	AVPSRLDQ
	AVPSLLDQ
	AVPSRLQ
	AVPSLLQ
	AVPSRLD
	AVPSLLD
	AVPOR
	AVPBLENDD
	AVINSERTI128
	AVPERM2I128
	ARORXL
	ARORXQ
	// from 386
	AJCXZW
	AFCMOVCC
	AFCMOVCS
	AFCMOVEQ
	AFCMOVHI
	AFCMOVLS
	AFCMOVNE
	AFCMOVNU
	AFCMOVUN
	AFCOMI
	AFCOMIP
	AFUCOMI
	AFUCOMIP
	// TSX
	AXACQUIRE
	AXRELEASE
	AXBEGIN
	AXEND
	AXABORT
	AXTEST
	ALAST
)



const (
	//  mark flags
	DONE          = 1 << iota
	PRESERVEFLAGS // not allowed to clobber flags

)



const (
	E = 0xff
)



const (
	// Loop alignment constants:
	// want to align loop entry to LoopAlign-byte boundary,
	// and willing to insert at most MaxLoopPad bytes of NOP to do so.
	// We define a loop entry as the target of a backward jump.
	//
	// gcc uses MaxLoopPad = 10 for its 'generic x86-64' config,
	// and it aligns all jump targets, not just backward jump targets.
	//
	// As of 6/1/2012, the effect of setting MaxLoopPad = 10 here
	// is very slight but negative, so the alignment is disabled by
	// setting MaxLoopPad = 0. The code is here for reference and
	// for future experiments.
	LoopAlign  = 16
	MaxLoopPad = 0
	FuncAlign  = 16
)



const (
	Px   = 0
	Px1  = 1      // symbolic; exact value doesn't matter
	P32  = 0x32   /* 32-bit only */
	Pe   = 0x66   /* operand escape */
	Pm   = 0x0f   /* 2byte opcode escape */
	Pq   = 0xff   /* both escapes: 66 0f */
	Pb   = 0xfe   /* byte operands */
	Pf2  = 0xf2   /* xmm escape 1: f2 0f */
	Pf3  = 0xf3   /* xmm escape 2: f3 0f */
	Pef3 = 0xf5   /* xmm escape 2 with 16-bit prefix: 66 f3 0f */
	Pq3  = 0x67   /* xmm escape 3: 66 48 0f */
	Pq4  = 0x68   /* xmm escape 4: 66 0F 38 */
	Pfw  = 0xf4   /* Pf3 with Rex.w: f3 48 0f */
	Pw   = 0x48   /* Rex.w */
	Pw8  = 0x90   // symbolic; exact value doesn't matter
	Py   = 0x80   /* defaults to 64-bit mode */
	Py1  = 0x81   // symbolic; exact value doesn't matter
	Py3  = 0x83   // symbolic; exact value doesn't matter
	Pvex = 0x84   // symbolic: exact value doesn't matter
	Rxw  = 1 << 3 /* =1, 64-bit operand size */
	Rxr  = 1 << 2 /* extend modrm reg */
	Rxx  = 1 << 1 /* extend sib index */
	Rxb  = 1 << 0 /* extend modrm r/m, sib base, or opcode reg */

)



const (
	REG_AL = obj.RBaseAMD64 + iota
	REG_CL
	REG_DL
	REG_BL
	REG_SPB
	REG_BPB
	REG_SIB
	REG_DIB
	REG_R8B
	REG_R9B
	REG_R10B
	REG_R11B
	REG_R12B
	REG_R13B
	REG_R14B
	REG_R15B
	REG_AX
	REG_CX
	REG_DX
	REG_BX
	REG_SP
	REG_BP
	REG_SI
	REG_DI
	REG_R8
	REG_R9
	REG_R10
	REG_R11
	REG_R12
	REG_R13
	REG_R14
	REG_R15
	REG_AH
	REG_CH
	REG_DH
	REG_BH
	REG_F0
	REG_F1
	REG_F2
	REG_F3
	REG_F4
	REG_F5
	REG_F6
	REG_F7
	REG_M0
	REG_M1
	REG_M2
	REG_M3
	REG_M4
	REG_M5
	REG_M6
	REG_M7
	REG_X0
	REG_X1
	REG_X2
	REG_X3
	REG_X4
	REG_X5
	REG_X6
	REG_X7
	REG_X8
	REG_X9
	REG_X10
	REG_X11
	REG_X12
	REG_X13
	REG_X14
	REG_X15
	REG_Y0
	REG_Y1
	REG_Y2
	REG_Y3
	REG_Y4
	REG_Y5
	REG_Y6
	REG_Y7
	REG_Y8
	REG_Y9
	REG_Y10
	REG_Y11
	REG_Y12
	REG_Y13
	REG_Y14
	REG_Y15
	REG_CS
	REG_SS
	REG_DS
	REG_ES
	REG_FS
	REG_GS
	REG_GDTR /* global descriptor table register */
	REG_IDTR /* interrupt descriptor table register */
	REG_LDTR /* local descriptor table register */
	REG_MSW  /* machine status word */
	REG_TASK /* task register */
	REG_CR0
	REG_CR1
	REG_CR2
	REG_CR3
	REG_CR4
	REG_CR5
	REG_CR6
	REG_CR7
	REG_CR8
	REG_CR9
	REG_CR10
	REG_CR11
	REG_CR12
	REG_CR13
	REG_CR14
	REG_CR15
	REG_DR0
	REG_DR1
	REG_DR2
	REG_DR3
	REG_DR4
	REG_DR5
	REG_DR6
	REG_DR7
	REG_TR0
	REG_TR1
	REG_TR2
	REG_TR3
	REG_TR4
	REG_TR5
	REG_TR6
	REG_TR7
	REG_TLS
	MAXREG
	REG_CR   = REG_CR0
	REG_DR   = REG_DR0
	REG_TR   = REG_TR0
	REGARG   = -1
	REGRET   = REG_AX
	FREGRET  = REG_X0
	REGSP    = REG_SP
	REGCTXT  = REG_DX
	REGEXT   = REG_R15     /* compiler allocates external registers R15 down */
	FREGMIN  = REG_X0 + 5  /* first register variable */
	FREGEXT  = REG_X0 + 15 /* first external register */
	T_TYPE   = 1 << 0
	T_INDEX  = 1 << 1
	T_OFFSET = 1 << 2
	T_FCONST = 1 << 3
	T_SYM    = 1 << 4
	T_SCONST = 1 << 5
	T_64     = 1 << 6
	T_GOTYPE = 1 << 7
)



const (
	REG_NONE = 0
)



const (
	// Combinations used in the manual.
	VEX_128_0F_WIG      = vex128 | vex0F | vexWIG
	VEX_128_66_0F_W0    = vex128 | vex66 | vex0F | vexW0
	VEX_128_66_0F_W1    = vex128 | vex66 | vex0F | vexW1
	VEX_128_66_0F_WIG   = vex128 | vex66 | vex0F | vexWIG
	VEX_128_66_0F38_W0  = vex128 | vex66 | vex0F38 | vexW0
	VEX_128_66_0F38_W1  = vex128 | vex66 | vex0F38 | vexW1
	VEX_128_66_0F38_WIG = vex128 | vex66 | vex0F38 | vexWIG
	VEX_128_66_0F3A_W0  = vex128 | vex66 | vex0F3A | vexW0
	VEX_128_66_0F3A_W1  = vex128 | vex66 | vex0F3A | vexW1
	VEX_128_66_0F3A_WIG = vex128 | vex66 | vex0F3A | vexWIG
	VEX_128_F2_0F_WIG   = vex128 | vexF2 | vex0F | vexWIG
	VEX_128_F3_0F_WIG   = vex128 | vexF3 | vex0F | vexWIG
	VEX_256_66_0F_WIG   = vex256 | vex66 | vex0F | vexWIG
	VEX_256_66_0F38_W0  = vex256 | vex66 | vex0F38 | vexW0
	VEX_256_66_0F38_W1  = vex256 | vex66 | vex0F38 | vexW1
	VEX_256_66_0F38_WIG = vex256 | vex66 | vex0F38 | vexWIG
	VEX_256_66_0F3A_W0  = vex256 | vex66 | vex0F3A | vexW0
	VEX_256_66_0F3A_W1  = vex256 | vex66 | vex0F3A | vexW1
	VEX_256_66_0F3A_WIG = vex256 | vex66 | vex0F3A | vexWIG
	VEX_256_F2_0F_WIG   = vex256 | vexF2 | vex0F | vexWIG
	VEX_256_F3_0F_WIG   = vex256 | vexF3 | vex0F | vexWIG
	VEX_LIG_0F_WIG      = vexLIG | vex0F | vexWIG
	VEX_LIG_66_0F_WIG   = vexLIG | vex66 | vex0F | vexWIG
	VEX_LIG_66_0F38_W0  = vexLIG | vex66 | vex0F38 | vexW0
	VEX_LIG_66_0F38_W1  = vexLIG | vex66 | vex0F38 | vexW1
	VEX_LIG_66_0F3A_WIG = vexLIG | vex66 | vex0F3A | vexWIG
	VEX_LIG_F2_0F_W0    = vexLIG | vexF2 | vex0F | vexW0
	VEX_LIG_F2_0F_W1    = vexLIG | vexF2 | vex0F | vexW1
	VEX_LIG_F2_0F_WIG   = vexLIG | vexF2 | vex0F | vexWIG
	VEX_LIG_F3_0F_W0    = vexLIG | vexF3 | vex0F | vexW0
	VEX_LIG_F3_0F_W1    = vexLIG | vexF3 | vex0F | vexW1
	VEX_LIG_F3_0F_WIG   = vexLIG | vexF3 | vex0F | vexWIG
	VEX_LZ_0F_WIG       = vexLZ | vex0F | vexWIG
	VEX_LZ_0F38_W0      = vexLZ | vex0F38 | vexW0
	VEX_LZ_0F38_W1      = vexLZ | vex0F38 | vexW1
	VEX_LZ_66_0F38_W0   = vexLZ | vex66 | vex0F38 | vexW0
	VEX_LZ_66_0F38_W1   = vexLZ | vex66 | vex0F38 | vexW1
	VEX_LZ_F2_0F38_W0   = vexLZ | vexF2 | vex0F38 | vexW0
	VEX_LZ_F2_0F38_W1   = vexLZ | vexF2 | vex0F38 | vexW1
	VEX_LZ_F2_0F3A_W0   = vexLZ | vexF2 | vex0F3A | vexW0
	VEX_LZ_F2_0F3A_W1   = vexLZ | vexF2 | vex0F3A | vexW1
	VEX_LZ_F3_0F38_W0   = vexLZ | vexF3 | vex0F38 | vexW0
	VEX_LZ_F3_0F38_W1   = vexLZ | vexF3 | vex0F38 | vexW1
)



const (
	Yxxx = iota
	Ynone
	Yi0 // $0
	Yi1 // $1
	Yi8 // $x, x fits in int8
	Yu8 // $x, x fits in uint8
	Yu7 // $x, x in 0..127 (fits in both int8 and uint8)
	Ys32
	Yi32
	Yi64
	Yiauto
	Yal
	Ycl
	Yax
	Ycx
	Yrb
	Yrl
	Yrl32 // Yrl on 32-bit system
	Yrf
	Yf0
	Yrx
	Ymb
	Yml
	Ym
	Ybr
	Ycs
	Yss
	Yds
	Yes
	Yfs
	Ygs
	Ygdtr
	Yidtr
	Yldtr
	Ymsw
	Ytask
	Ycr0
	Ycr1
	Ycr2
	Ycr3
	Ycr4
	Ycr5
	Ycr6
	Ycr7
	Ycr8
	Ydr0
	Ydr1
	Ydr2
	Ydr3
	Ydr4
	Ydr5
	Ydr6
	Ydr7
	Ytr0
	Ytr1
	Ytr2
	Ytr3
	Ytr4
	Ytr5
	Ytr6
	Ytr7
	Ymr
	Ymm
	Yxr
	Yxm
	Yyr
	Yym
	Ytls
	Ytextsize
	Yindir
	Ymax
)



const (
	Zxxx = iota
	Zlit
	Zlitm_r
	Z_rp
	Zbr
	Zcall
	Zcallcon
	Zcallduff
	Zcallind
	Zcallindreg
	Zib_
	Zib_rp
	Zibo_m
	Zibo_m_xm
	Zil_
	Zil_rp
	Ziq_rp
	Zilo_m
	Zjmp
	Zjmpcon
	Zloop
	Zo_iw
	Zm_o
	Zm_r
	Zm2_r
	Zm_r_xm
	Zm_r_i_xm
	Zm_r_xm_nr
	Zr_m_xm_nr
	Zibm_r /* mmx1,mmx2/mem64,imm8 */
	Zibr_m
	Zmb_r
	Zaut_r
	Zo_m
	Zo_m64
	Zpseudo
	Zr_m
	Zr_m_xm
	Zrp_
	Z_ib
	Z_il
	Zm_ibo
	Zm_ilo
	Zib_rr
	Zil_rr
	Zclr
	Zbyte
	Zvex_rm_v_r
	Zvex_r_v_rm
	Zvex_v_rm_r
	Zvex_i_rm_r
	Zvex_i_r_v
	Zvex_i_rm_v_r
	Zmax
)



var Anames = []string{
	obj.A_ARCHSPECIFIC: "AAA",
	"AAD",
	"AAM",
	"AAS",
	"ADCB",
	"ADCL",
	"ADCW",
	"ADDB",
	"ADDL",
	"ADDW",
	"ADJSP",
	"ANDB",
	"ANDL",
	"ANDW",
	"ARPL",
	"BOUNDL",
	"BOUNDW",
	"BSFL",
	"BSFW",
	"BSRL",
	"BSRW",
	"BTL",
	"BTW",
	"BTCL",
	"BTCW",
	"BTRL",
	"BTRW",
	"BTSL",
	"BTSW",
	"BYTE",
	"CLC",
	"CLD",
	"CLI",
	"CLTS",
	"CMC",
	"CMPB",
	"CMPL",
	"CMPW",
	"CMPSB",
	"CMPSL",
	"CMPSW",
	"DAA",
	"DAS",
	"DECB",
	"DECL",
	"DECQ",
	"DECW",
	"DIVB",
	"DIVL",
	"DIVW",
	"ENTER",
	"HADDPD",
	"HADDPS",
	"HLT",
	"HSUBPD",
	"HSUBPS",
	"IDIVB",
	"IDIVL",
	"IDIVW",
	"IMULB",
	"IMULL",
	"IMULW",
	"INB",
	"INL",
	"INW",
	"INCB",
	"INCL",
	"INCQ",
	"INCW",
	"INSB",
	"INSL",
	"INSW",
	"INT",
	"INTO",
	"IRETL",
	"IRETW",
	"JCC",
	"JCS",
	"JCXZL",
	"JEQ",
	"JGE",
	"JGT",
	"JHI",
	"JLE",
	"JLS",
	"JLT",
	"JMI",
	"JNE",
	"JOC",
	"JOS",
	"JPC",
	"JPL",
	"JPS",
	"LAHF",
	"LARL",
	"LARW",
	"LEAL",
	"LEAW",
	"LEAVEL",
	"LEAVEW",
	"LOCK",
	"LODSB",
	"LODSL",
	"LODSW",
	"LONG",
	"LOOP",
	"LOOPEQ",
	"LOOPNE",
	"LSLL",
	"LSLW",
	"MOVB",
	"MOVL",
	"MOVW",
	"MOVBLSX",
	"MOVBLZX",
	"MOVBQSX",
	"MOVBQZX",
	"MOVBWSX",
	"MOVBWZX",
	"MOVWLSX",
	"MOVWLZX",
	"MOVWQSX",
	"MOVWQZX",
	"MOVSB",
	"MOVSL",
	"MOVSW",
	"MULB",
	"MULL",
	"MULW",
	"NEGB",
	"NEGL",
	"NEGW",
	"NOTB",
	"NOTL",
	"NOTW",
	"ORB",
	"ORL",
	"ORW",
	"OUTB",
	"OUTL",
	"OUTW",
	"OUTSB",
	"OUTSL",
	"OUTSW",
	"PAUSE",
	"POPAL",
	"POPAW",
	"POPCNTW",
	"POPCNTL",
	"POPCNTQ",
	"POPFL",
	"POPFW",
	"POPL",
	"POPW",
	"PUSHAL",
	"PUSHAW",
	"PUSHFL",
	"PUSHFW",
	"PUSHL",
	"PUSHW",
	"RCLB",
	"RCLL",
	"RCLW",
	"RCRB",
	"RCRL",
	"RCRW",
	"REP",
	"REPN",
	"ROLB",
	"ROLL",
	"ROLW",
	"RORB",
	"RORL",
	"RORW",
	"SAHF",
	"SALB",
	"SALL",
	"SALW",
	"SARB",
	"SARL",
	"SARW",
	"SBBB",
	"SBBL",
	"SBBW",
	"SCASB",
	"SCASL",
	"SCASW",
	"SETCC",
	"SETCS",
	"SETEQ",
	"SETGE",
	"SETGT",
	"SETHI",
	"SETLE",
	"SETLS",
	"SETLT",
	"SETMI",
	"SETNE",
	"SETOC",
	"SETOS",
	"SETPC",
	"SETPL",
	"SETPS",
	"CDQ",
	"CWD",
	"SHLB",
	"SHLL",
	"SHLW",
	"SHRB",
	"SHRL",
	"SHRW",
	"STC",
	"STD",
	"STI",
	"STOSB",
	"STOSL",
	"STOSW",
	"SUBB",
	"SUBL",
	"SUBW",
	"SYSCALL",
	"TESTB",
	"TESTL",
	"TESTW",
	"VERR",
	"VERW",
	"WAIT",
	"WORD",
	"XCHGB",
	"XCHGL",
	"XCHGW",
	"XLAT",
	"XORB",
	"XORL",
	"XORW",
	"FMOVB",
	"FMOVBP",
	"FMOVD",
	"FMOVDP",
	"FMOVF",
	"FMOVFP",
	"FMOVL",
	"FMOVLP",
	"FMOVV",
	"FMOVVP",
	"FMOVW",
	"FMOVWP",
	"FMOVX",
	"FMOVXP",
	"FCOMD",
	"FCOMDP",
	"FCOMDPP",
	"FCOMF",
	"FCOMFP",
	"FCOML",
	"FCOMLP",
	"FCOMW",
	"FCOMWP",
	"FUCOM",
	"FUCOMP",
	"FUCOMPP",
	"FADDDP",
	"FADDW",
	"FADDL",
	"FADDF",
	"FADDD",
	"FMULDP",
	"FMULW",
	"FMULL",
	"FMULF",
	"FMULD",
	"FSUBDP",
	"FSUBW",
	"FSUBL",
	"FSUBF",
	"FSUBD",
	"FSUBRDP",
	"FSUBRW",
	"FSUBRL",
	"FSUBRF",
	"FSUBRD",
	"FDIVDP",
	"FDIVW",
	"FDIVL",
	"FDIVF",
	"FDIVD",
	"FDIVRDP",
	"FDIVRW",
	"FDIVRL",
	"FDIVRF",
	"FDIVRD",
	"FXCHD",
	"FFREE",
	"FLDCW",
	"FLDENV",
	"FRSTOR",
	"FSAVE",
	"FSTCW",
	"FSTENV",
	"FSTSW",
	"F2XM1",
	"FABS",
	"FCHS",
	"FCLEX",
	"FCOS",
	"FDECSTP",
	"FINCSTP",
	"FINIT",
	"FLD1",
	"FLDL2E",
	"FLDL2T",
	"FLDLG2",
	"FLDLN2",
	"FLDPI",
	"FLDZ",
	"FNOP",
	"FPATAN",
	"FPREM",
	"FPREM1",
	"FPTAN",
	"FRNDINT",
	"FSCALE",
	"FSIN",
	"FSINCOS",
	"FSQRT",
	"FTST",
	"FXAM",
	"FXTRACT",
	"FYL2X",
	"FYL2XP1",
	"CMPXCHGB",
	"CMPXCHGL",
	"CMPXCHGW",
	"CMPXCHG8B",
	"CPUID",
	"INVD",
	"INVLPG",
	"LFENCE",
	"MFENCE",
	"MOVNTIL",
	"RDMSR",
	"RDPMC",
	"RDTSC",
	"RSM",
	"SFENCE",
	"SYSRET",
	"WBINVD",
	"WRMSR",
	"XADDB",
	"XADDL",
	"XADDW",
	"CMOVLCC",
	"CMOVLCS",
	"CMOVLEQ",
	"CMOVLGE",
	"CMOVLGT",
	"CMOVLHI",
	"CMOVLLE",
	"CMOVLLS",
	"CMOVLLT",
	"CMOVLMI",
	"CMOVLNE",
	"CMOVLOC",
	"CMOVLOS",
	"CMOVLPC",
	"CMOVLPL",
	"CMOVLPS",
	"CMOVQCC",
	"CMOVQCS",
	"CMOVQEQ",
	"CMOVQGE",
	"CMOVQGT",
	"CMOVQHI",
	"CMOVQLE",
	"CMOVQLS",
	"CMOVQLT",
	"CMOVQMI",
	"CMOVQNE",
	"CMOVQOC",
	"CMOVQOS",
	"CMOVQPC",
	"CMOVQPL",
	"CMOVQPS",
	"CMOVWCC",
	"CMOVWCS",
	"CMOVWEQ",
	"CMOVWGE",
	"CMOVWGT",
	"CMOVWHI",
	"CMOVWLE",
	"CMOVWLS",
	"CMOVWLT",
	"CMOVWMI",
	"CMOVWNE",
	"CMOVWOC",
	"CMOVWOS",
	"CMOVWPC",
	"CMOVWPL",
	"CMOVWPS",
	"ADCQ",
	"ADDQ",
	"ANDQ",
	"BSFQ",
	"BSRQ",
	"BTCQ",
	"BTQ",
	"BTRQ",
	"BTSQ",
	"CMPQ",
	"CMPSQ",
	"CMPXCHGQ",
	"CQO",
	"DIVQ",
	"IDIVQ",
	"IMULQ",
	"IRETQ",
	"JCXZQ",
	"LEAQ",
	"LEAVEQ",
	"LODSQ",
	"MOVQ",
	"MOVLQSX",
	"MOVLQZX",
	"MOVNTIQ",
	"MOVSQ",
	"MULQ",
	"NEGQ",
	"NOTQ",
	"ORQ",
	"POPFQ",
	"POPQ",
	"PUSHFQ",
	"PUSHQ",
	"RCLQ",
	"RCRQ",
	"ROLQ",
	"RORQ",
	"QUAD",
	"SALQ",
	"SARQ",
	"SBBQ",
	"SCASQ",
	"SHLQ",
	"SHRQ",
	"STOSQ",
	"SUBQ",
	"TESTQ",
	"XADDQ",
	"XCHGQ",
	"XORQ",
	"XGETBV",
	"ADDPD",
	"ADDPS",
	"ADDSD",
	"ADDSS",
	"ANDNL",
	"ANDNQ",
	"ANDNPD",
	"ANDNPS",
	"ANDPD",
	"ANDPS",
	"BEXTRL",
	"BEXTRQ",
	"BLSIL",
	"BLSIQ",
	"BLSMSKL",
	"BLSMSKQ",
	"BLSRL",
	"BLSRQ",
	"BZHIL",
	"BZHIQ",
	"CMPPD",
	"CMPPS",
	"CMPSD",
	"CMPSS",
	"COMISD",
	"COMISS",
	"CVTPD2PL",
	"CVTPD2PS",
	"CVTPL2PD",
	"CVTPL2PS",
	"CVTPS2PD",
	"CVTPS2PL",
	"CVTSD2SL",
	"CVTSD2SQ",
	"CVTSD2SS",
	"CVTSL2SD",
	"CVTSL2SS",
	"CVTSQ2SD",
	"CVTSQ2SS",
	"CVTSS2SD",
	"CVTSS2SL",
	"CVTSS2SQ",
	"CVTTPD2PL",
	"CVTTPS2PL",
	"CVTTSD2SL",
	"CVTTSD2SQ",
	"CVTTSS2SL",
	"CVTTSS2SQ",
	"DIVPD",
	"DIVPS",
	"DIVSD",
	"DIVSS",
	"EMMS",
	"FXRSTOR",
	"FXRSTOR64",
	"FXSAVE",
	"FXSAVE64",
	"LDDQU",
	"LDMXCSR",
	"MASKMOVOU",
	"MASKMOVQ",
	"MAXPD",
	"MAXPS",
	"MAXSD",
	"MAXSS",
	"MINPD",
	"MINPS",
	"MINSD",
	"MINSS",
	"MOVAPD",
	"MOVAPS",
	"MOVOU",
	"MOVHLPS",
	"MOVHPD",
	"MOVHPS",
	"MOVLHPS",
	"MOVLPD",
	"MOVLPS",
	"MOVMSKPD",
	"MOVMSKPS",
	"MOVNTO",
	"MOVNTPD",
	"MOVNTPS",
	"MOVNTQ",
	"MOVO",
	"MOVQOZX",
	"MOVSD",
	"MOVSS",
	"MOVUPD",
	"MOVUPS",
	"MULPD",
	"MULPS",
	"MULSD",
	"MULSS",
	"MULXL",
	"MULXQ",
	"ORPD",
	"ORPS",
	"PACKSSLW",
	"PACKSSWB",
	"PACKUSWB",
	"PADDB",
	"PADDL",
	"PADDQ",
	"PADDSB",
	"PADDSW",
	"PADDUSB",
	"PADDUSW",
	"PADDW",
	"PAND",
	"PANDN",
	"PAVGB",
	"PAVGW",
	"PCMPEQB",
	"PCMPEQL",
	"PCMPEQW",
	"PCMPGTB",
	"PCMPGTL",
	"PCMPGTW",
	"PDEPL",
	"PDEPQ",
	"PEXTL",
	"PEXTQ",
	"PEXTRB",
	"PEXTRD",
	"PEXTRQ",
	"PEXTRW",
	"PHADDD",
	"PHADDSW",
	"PHADDW",
	"PHMINPOSUW",
	"PHSUBD",
	"PHSUBSW",
	"PHSUBW",
	"PINSRB",
	"PINSRD",
	"PINSRQ",
	"PINSRW",
	"PMADDWL",
	"PMAXSW",
	"PMAXUB",
	"PMINSW",
	"PMINUB",
	"PMOVMSKB",
	"PMOVSXBD",
	"PMOVSXBQ",
	"PMOVSXBW",
	"PMOVSXDQ",
	"PMOVSXWD",
	"PMOVSXWQ",
	"PMOVZXBD",
	"PMOVZXBQ",
	"PMOVZXBW",
	"PMOVZXDQ",
	"PMOVZXWD",
	"PMOVZXWQ",
	"PMULDQ",
	"PMULHUW",
	"PMULHW",
	"PMULLD",
	"PMULLW",
	"PMULULQ",
	"POR",
	"PSADBW",
	"PSHUFB",
	"PSHUFHW",
	"PSHUFL",
	"PSHUFLW",
	"PSHUFW",
	"PSLLL",
	"PSLLO",
	"PSLLQ",
	"PSLLW",
	"PSRAL",
	"PSRAW",
	"PSRLL",
	"PSRLO",
	"PSRLQ",
	"PSRLW",
	"PSUBB",
	"PSUBL",
	"PSUBQ",
	"PSUBSB",
	"PSUBSW",
	"PSUBUSB",
	"PSUBUSW",
	"PSUBW",
	"PUNPCKHBW",
	"PUNPCKHLQ",
	"PUNPCKHQDQ",
	"PUNPCKHWL",
	"PUNPCKLBW",
	"PUNPCKLLQ",
	"PUNPCKLQDQ",
	"PUNPCKLWL",
	"PXOR",
	"RCPPS",
	"RCPSS",
	"RSQRTPS",
	"RSQRTSS",
	"SARXL",
	"SARXQ",
	"SHLXL",
	"SHLXQ",
	"SHRXL",
	"SHRXQ",
	"SHUFPD",
	"SHUFPS",
	"SQRTPD",
	"SQRTPS",
	"SQRTSD",
	"SQRTSS",
	"STMXCSR",
	"SUBPD",
	"SUBPS",
	"SUBSD",
	"SUBSS",
	"UCOMISD",
	"UCOMISS",
	"UNPCKHPD",
	"UNPCKHPS",
	"UNPCKLPD",
	"UNPCKLPS",
	"XORPD",
	"XORPS",
	"PCMPESTRI",
	"RETFW",
	"RETFL",
	"RETFQ",
	"SWAPGS",
	"MODE",
	"CRC32B",
	"CRC32Q",
	"IMUL3Q",
	"PREFETCHT0",
	"PREFETCHT1",
	"PREFETCHT2",
	"PREFETCHNTA",
	"MOVQL",
	"BSWAPL",
	"BSWAPQ",
	"AESENC",
	"AESENCLAST",
	"AESDEC",
	"AESDECLAST",
	"AESIMC",
	"AESKEYGENASSIST",
	"ROUNDPS",
	"ROUNDSS",
	"ROUNDPD",
	"ROUNDSD",
	"PSHUFD",
	"PCLMULQDQ",
	"VZEROUPPER",
	"VMOVDQU",
	"VMOVNTDQ",
	"VMOVDQA",
	"VPCMPEQB",
	"VPXOR",
	"VPMOVMSKB",
	"VPAND",
	"VPTEST",
	"VPBROADCASTB",
	"VPSHUFB",
	"VPSHUFD",
	"VPERM2F128",
	"VPALIGNR",
	"VPADDQ",
	"VPADDD",
	"VPSRLDQ",
	"VPSLLDQ",
	"VPSRLQ",
	"VPSLLQ",
	"VPSRLD",
	"VPSLLD",
	"VPOR",
	"VPBLENDD",
	"VINSERTI128",
	"VPERM2I128",
	"RORXL",
	"RORXQ",
	"JCXZW",
	"FCMOVCC",
	"FCMOVCS",
	"FCMOVEQ",
	"FCMOVHI",
	"FCMOVLS",
	"FCMOVNE",
	"FCMOVNU",
	"FCMOVUN",
	"FCOMI",
	"FCOMIP",
	"FUCOMI",
	"FUCOMIP",
	"XACQUIRE",
	"XRELEASE",
	"XBEGIN",
	"XEND",
	"XABORT",
	"XTEST",
	"LAST",
}



var Link386 = obj.LinkArch{
	Arch:       sys.Arch386,
	Preprocess: preprocess,
	Assemble:   span6,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
}



var Linkamd64 = obj.LinkArch{
	Arch:       sys.ArchAMD64,
	Preprocess: preprocess,
	Assemble:   span6,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
}



var Linkamd64p32 = obj.LinkArch{
	Arch:       sys.ArchAMD64P32,
	Preprocess: preprocess,
	Assemble:   span6,
	Follow:     follow,
	Progedit:   progedit,
	UnaryDst:   unaryDst,
}



var Register = []string{
	"AL",
	"CL",
	"DL",
	"BL",
	"SPB",
	"BPB",
	"SIB",
	"DIB",
	"R8B",
	"R9B",
	"R10B",
	"R11B",
	"R12B",
	"R13B",
	"R14B",
	"R15B",
	"AX",
	"CX",
	"DX",
	"BX",
	"SP",
	"BP",
	"SI",
	"DI",
	"R8",
	"R9",
	"R10",
	"R11",
	"R12",
	"R13",
	"R14",
	"R15",
	"AH",
	"CH",
	"DH",
	"BH",
	"F0",
	"F1",
	"F2",
	"F3",
	"F4",
	"F5",
	"F6",
	"F7",
	"M0",
	"M1",
	"M2",
	"M3",
	"M4",
	"M5",
	"M6",
	"M7",
	"X0",
	"X1",
	"X2",
	"X3",
	"X4",
	"X5",
	"X6",
	"X7",
	"X8",
	"X9",
	"X10",
	"X11",
	"X12",
	"X13",
	"X14",
	"X15",
	"Y0",
	"Y1",
	"Y2",
	"Y3",
	"Y4",
	"Y5",
	"Y6",
	"Y7",
	"Y8",
	"Y9",
	"Y10",
	"Y11",
	"Y12",
	"Y13",
	"Y14",
	"Y15",
	"CS",
	"SS",
	"DS",
	"ES",
	"FS",
	"GS",
	"GDTR",
	"IDTR",
	"LDTR",
	"MSW",
	"TASK",
	"CR0",
	"CR1",
	"CR2",
	"CR3",
	"CR4",
	"CR5",
	"CR6",
	"CR7",
	"CR8",
	"CR9",
	"CR10",
	"CR11",
	"CR12",
	"CR13",
	"CR14",
	"CR15",
	"DR0",
	"DR1",
	"DR2",
	"DR3",
	"DR4",
	"DR5",
	"DR6",
	"DR7",
	"TR0",
	"TR1",
	"TR2",
	"TR3",
	"TR4",
	"TR5",
	"TR6",
	"TR7",
	"TLS",
	"MAXREG",
}



type Movtab struct {
	as   obj.As
	ft   uint8
	f3t  uint8
	tt   uint8
	code uint8
	op   [4]uint8
}



type Optab struct {
	as     obj.As
	ytab   []ytab
	prefix uint8
	op     [23]uint8
}


func CanUse1InsnTLS(ctxt *obj.Link) bool

func Rconv(r int) string

