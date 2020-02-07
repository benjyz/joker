// +build !fast_init

package core

var privateMeta Map = EmptyArrayMap().Assoc(KEYWORDS.private, Boolean{B: true}).(Map)

func intern(name string, proc ProcFn, procName string) {
	vr := GLOBAL_ENV.CoreNamespace.Intern(MakeSymbol(name))
	vr.Value = Proc{Fn: proc, Name: procName}
	vr.isPrivate = true
	vr.meta = privateMeta
}

func init() {
	GLOBAL_ENV.CoreNamespace.InternVar("*assert*", Boolean{B: true},
		MakeMeta(nil, "When set to logical false, assert is a noop. Defaults to true.", "1.0"))

	intern("list__", procList, "procList")
	intern("cons__", procCons, "procCons")
	intern("first__", procFirst, "procFirst")
	intern("next__", procNext, "procNext")
	intern("rest__", procRest, "procRest")
	intern("conj__", procConj, "procConj")
	intern("seq__", procSeq, "procSeq")
	intern("instance?__", procIsInstance, "procIsInstance")
	intern("assoc__", procAssoc, "procAssoc")
	intern("meta__", procMeta, "procMeta")
	intern("with-meta__", procWithMeta, "procWithMeta")
	intern("=__", procEquals, "procEquals")
	intern("count__", procCount, "procCount")
	intern("subvec__", procSubvec, "procSubvec")
	intern("cast__", procCast, "procCast")
	intern("vec__", procVec, "procVec")
	intern("hash-map__", procHashMap, "procHashMap")
	intern("hash-set__", procHashSet, "procHashSet")
	intern("str__", procStr, "procStr")
	intern("symbol__", procSymbol, "procSymbol")
	intern("gensym__", procGensym, "procGensym")
	intern("keyword__", procKeyword, "procKeyword")
	intern("apply__", procApply, "procApply")
	intern("lazy-seq__", procLazySeq, "procLazySeq")
	intern("delay__", procDelay, "procDelay")
	intern("force__", procForce, "procForce")
	intern("identical__", procIdentical, "procIdentical")
	intern("compare__", procCompare, "procCompare")
	intern("zero?__", procIsZero, "procIsZero")
	intern("int__", procInt, "procInt")
	intern("nth__", procNth, "procNth")
	intern("<__", procLt, "procLt")
	intern("<=__", procLte, "procLte")
	intern(">__", procGt, "procGt")
	intern(">=__", procGte, "procGte")
	intern("==__", procEq, "procEq")
	intern("inc'__", procIncEx, "procIncEx")
	intern("inc__", procInc, "procInc")
	intern("dec'__", procDecEx, "procDecEx")
	intern("dec__", procDec, "procDec")
	intern("add'__", procAddEx, "procAddEx")
	intern("add__", procAdd, "procAdd")
	intern("multiply'__", procMultiplyEx, "procMultiplyEx")
	intern("multiply__", procMultiply, "procMultiply")
	intern("divide__", procDivide, "procDivide")
	intern("subtract'__", procSubtractEx, "procSubtractEx")
	intern("subtract__", procSubtract, "procSubtract")
	intern("max__", procMax, "procMax")
	intern("min__", procMin, "procMin")
	intern("pos__", procIsPos, "procIsPos")
	intern("neg__", procIsNeg, "procIsNeg")
	intern("quot__", procQuot, "procQuot")
	intern("rem__", procRem, "procRem")
	intern("bit-not__", procBitNot, "procBitNot")
	intern("bit-and__", procBitAnd, "procBitAnd")
	intern("bit-or__", procBitOr, "procBitOr")
	intern("bit-xor_", procBitXor, "procBitXor")
	intern("bit-and-not__", procBitAndNot, "procBitAndNot")
	intern("bit-clear__", procBitClear, "procBitClear")
	intern("bit-set__", procBitSet, "procBitSet")
	intern("bit-flip__", procBitFlip, "procBitFlip")
	intern("bit-test__", procBitTest, "procBitTest")
	intern("bit-shift-left__", procBitShiftLeft, "procBitShiftLeft")
	intern("bit-shift-right__", procBitShiftRight, "procBitShiftRight")
	intern("unsigned-bit-shift-right__", procUnsignedBitShiftRight, "procUnsignedBitShiftRight")
	intern("peek__", procPeek, "procPeek")
	intern("pop__", procPop, "procPop")
	intern("contains?__", procContains, "procContains")
	intern("get__", procGet, "procGet")
	intern("dissoc__", procDissoc, "procDissoc")
	intern("disj__", procDisj, "procDisj")
	intern("find__", procFind, "procFind")
	intern("keys__", procKeys, "procKeys")
	intern("vals__", procVals, "procVals")
	intern("rseq__", procRseq, "procRseq")
	intern("name__", procName, "procName")
	intern("namespace__", procNamespace, "procNamespace")
	intern("find-var__", procFindVar, "procFindVar")
	intern("sort__", procSort, "procSort")
	intern("eval__", procEval, "procEval")
	intern("type__", procType, "procType")
	intern("num__", procNumber, "procNumber")
	intern("double__", procDouble, "procDouble")
	intern("char__", procChar, "procChar")
	intern("boolean__", procBoolean, "procBoolean")
	intern("numerator__", procNumerator, "procNumerator")
	intern("denominator__", procDenominator, "procDenominator")
	intern("bigint__", procBigInt, "procBigInt")
	intern("bigfloat__", procBigFloat, "procBigFloat")
	intern("pr__", procPr, "procPr")
	intern("pprint__", procPprint, "procPprint")
	intern("newline__", procNewline, "procNewline")
	intern("flush__", procFlush, "procFlush")
	intern("read__", procRead, "procRead")
	intern("read-line__", procReadLine, "procReadLine")
	intern("reader-read-line__", procReaderReadLine, "procReaderReadLine")
	intern("read-string__", procReadString, "procReadString")
	intern("nano-time__", procNanoTime, "procNanoTime")
	intern("macroexpand-1__", procMacroexpand1, "procMacroexpand1")
	intern("load-string__", procLoadString, "procLoadString")
	intern("find-ns__", procFindNamespace, "procFindNamespace")
	intern("create-ns__", procCreateNamespace, "procCreateNamespace")
	intern("inject-ns__", procInjectNamespace, "procInjectNamespace")
	intern("remove-ns__", procRemoveNamespace, "procRemoveNamespace")
	intern("all-ns__", procAllNamespaces, "procAllNamespaces")
	intern("ns-name__", procNamespaceName, "procNamespaceName")
	intern("ns-map__", procNamespaceMap, "procNamespaceMap")
	intern("ns-unmap__", procNamespaceUnmap, "procNamespaceUnmap")
	intern("var-ns__", procVarNamespace, "procVarNamespace")
	intern("ns-initialized?__", procIsNamespaceInitialized, "procIsNamespaceInitialized")
	intern("refer__", procRefer, "procRefer")
	intern("alias__", procAlias, "procAlias")
	intern("ns-aliases__", procNamespaceAliases, "procNamespaceAliases")
	intern("ns-unalias__", procNamespaceUnalias, "procNamespaceUnalias")
	intern("var-get__", procVarGet, "procVarGet")
	intern("var-set__", procVarSet, "procVarSet")
	intern("ns-resolve__", procNsResolve, "procNsResolve")
	intern("array-map__", procArrayMap, "procArrayMap")
	intern("buffer__", procBuffer, "procBuffer")
	intern("buffered-reader__", procBufferedReader, "procBufferedReader")
	intern("ex-info__", procExInfo, "procExInfo")
	intern("ex-data__", procExData, "procExData")
	intern("ex-cause__", procExCause, "procExCause")
	intern("ex-message__", procExMessage, "procExMessage")
	intern("regex__", procRegex, "procRegex")
	intern("re-seq__", procReSeq, "procReSeq")
	intern("re-find__", procReFind, "procReFind")
	intern("rand__", procRand, "procRand")
	intern("special-symbol?__", procIsSpecialSymbol, "procIsSpecialSymbol")
	intern("subs__", procSubs, "procSubs")
	intern("intern__", procIntern, "procIntern")
	intern("set-meta__", procSetMeta, "procSetMeta")
	intern("atom__", procAtom, "procAtom")
	intern("deref__", procDeref, "procDeref")
	intern("swap__", procSwap, "procSwap")
	intern("swap-vals__", procSwapVals, "procSwapVals")
	intern("reset__", procReset, "procReset")
	intern("reset-vals__", procResetVals, "procResetVals")
	intern("alter-meta__", procAlterMeta, "procAlterMeta")
	intern("reset-meta__", procResetMeta, "procResetMeta")
	intern("empty__", procEmpty, "procEmpty")
	intern("bound?__", procIsBound, "procIsBound")
	intern("format__", procFormat, "procFormat")
	intern("load-file__", procLoadFile, "procLoadFile")
	intern("load-lib-from-path__", procLoadLibFromPath, "procLoadLibFromPath")
	intern("reduce-kv__", procReduceKv, "procReduceKv")
	intern("slurp__", procSlurp, "procSlurp")
	intern("spit__", procSpit, "procSpit")
	intern("shuffle__", procShuffle, "procShuffle")
	intern("realized?__", procIsRealized, "procIsRealized")
	intern("derive-info__", procDeriveInfo, "procDeriveInfo")
	intern("joker-version__", procJokerVersion, "procJokerVersion")

	intern("hash__", procHash, "procHash")

	intern("index-of__", procIndexOf, "procIndexOf")
	intern("lib-path__", procLibPath, "procLibPath")
	intern("intern-fake-var__", procInternFakeVar, "procInternFakeVar")
	intern("parse__", procParse, "procParse")
	intern("inc-problem-count__", procIncProblemCount, "procIncProblemCount")
	intern("types__", procTypes, "procTypes")
	intern("go__", procGo, "procGo")
	intern("<!__", procReceive, "procReceive")
	intern(">!__", procSend, "procSend")
	intern("chan__", procCreateChan, "procCreateChan")
	intern("close!__", procCloseChan, "procCloseChan")

	intern("go-spew__", procGoSpew, "procGoSpew")
	intern("verbosity-level__", procVerbosityLevel, "procVerbosityLevel")

	intern("goobject?__", procGoObject, "procGoObject")
	intern("Go__", proc_Go, "proc_Go")
	intern("new__", procNew, "procNew")
	intern("GoTypeOf__", procGoTypeOf, "procGoTypeOf")
	intern("GoTypeOfAsString__", procGoTypeOfAsString, "procGoTypeOfAsString")
	intern("ref__", procRef, "procRef")
}

func lateInitializations() {
	// none needed for !fast_init
}
