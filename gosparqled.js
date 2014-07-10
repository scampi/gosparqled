"use strict";
(function() {

Error.stackTraceLimit = -1;

var $global, $module;
if (typeof window !== "undefined") { /* web page */
	$global = window;
} else if (typeof self !== "undefined") { /* web worker */
	$global = self;
} else if (typeof global !== "undefined") { /* Node.js */
	$global = global;
	$global.require = require;
} else {
	console.log("warning: no global object found");
}
if (typeof module !== "undefined") {
	$module = module;
}

var $idCounter = 0;
var $keys = function(m) { return m ? Object.keys(m) : []; };
var $min = Math.min;
var $parseInt = parseInt;
var $parseFloat = function(f) {
	if (f.constructor === Number) {
		return f;
	}
	return parseFloat(f);
};
var $mod = function(x, y) { return x % y; };
var $toString = String;
var $reflect, $newStringPtr;
var $Array = Array;

var $floatKey = function(f) {
	if (f !== f) {
		$idCounter++;
		return "NaN$" + $idCounter;
	}
	return String(f);
};

var $mapArray = function(array, f) {
	var newArray = new array.constructor(array.length), i;
	for (i = 0; i < array.length; i++) {
		newArray[i] = f(array[i]);
	}
	return newArray;
};

var $newType = function(size, kind, string, name, pkgPath, constructor) {
	var typ;
	switch(kind) {
	case "Bool":
	case "Int":
	case "Int8":
	case "Int16":
	case "Int32":
	case "Uint":
	case "Uint8" :
	case "Uint16":
	case "Uint32":
	case "Uintptr":
	case "String":
	case "UnsafePointer":
		typ = function(v) { this.$val = v; };
		typ.prototype.$key = function() { return string + "$" + this.$val; };
		break;

	case "Float32":
	case "Float64":
		typ = function(v) { this.$val = v; };
		typ.prototype.$key = function() { return string + "$" + $floatKey(this.$val); };
		break;

	case "Int64":
		typ = function(high, low) {
			this.$high = (high + Math.floor(Math.ceil(low) / 4294967296)) >> 0;
			this.$low = low >>> 0;
			this.$val = this;
		};
		typ.prototype.$key = function() { return string + "$" + this.$high + "$" + this.$low; };
		break;

	case "Uint64":
		typ = function(high, low) {
			this.$high = (high + Math.floor(Math.ceil(low) / 4294967296)) >>> 0;
			this.$low = low >>> 0;
			this.$val = this;
		};
		typ.prototype.$key = function() { return string + "$" + this.$high + "$" + this.$low; };
		break;

	case "Complex64":
	case "Complex128":
		typ = function(real, imag) {
			this.$real = real;
			this.$imag = imag;
			this.$val = this;
		};
		typ.prototype.$key = function() { return string + "$" + this.$real + "$" + this.$imag; };
		break;

	case "Array":
		typ = function(v) { this.$val = v; };
		typ.Ptr = $newType(4, "Ptr", "*" + string, "", "", function(array) {
			this.$get = function() { return array; };
			this.$val = array;
		});
		typ.init = function(elem, len) {
			typ.elem = elem;
			typ.len = len;
			typ.prototype.$key = function() {
				return string + "$" + Array.prototype.join.call($mapArray(this.$val, function(e) {
					var key = e.$key ? e.$key() : String(e);
					return key.replace(/\\/g, "\\\\").replace(/\$/g, "\\$");
				}), "$");
			};
			typ.extendReflectType = function(rt) {
				rt.arrayType = new $reflect.arrayType.Ptr(rt, elem.reflectType(), undefined, len);
			};
			typ.Ptr.init(typ);
			for (var i = 0; i < len; i++) {
				Object.defineProperty(typ.Ptr.nil, i, { get: $throwNilPointerError, set: $throwNilPointerError });
			}
		};
		break;

	case "Chan":
		typ = function(capacity) {
			this.$val = this;
			this.$capacity = capacity;
			this.$buffer = [];
			this.$sendQueue = [];
			this.$recvQueue = [];
			this.$closed = false;
		};
		typ.prototype.$key = function() {
			if (this.$id === undefined) {
				$idCounter++;
				this.$id = $idCounter;
			}
			return String(this.$id);
		};
		typ.init = function(elem, sendOnly, recvOnly) {
			typ.elem = elem;
			typ.sendOnly = sendOnly;
			typ.recvOnly = recvOnly;
			typ.nil = new typ(0);
			typ.nil.$sendQueue = typ.nil.$recvQueue = { length: 0, push: function() {}, shift: function() { return undefined; } };
			typ.extendReflectType = function(rt) {
				rt.chanType = new $reflect.chanType.Ptr(rt, elem.reflectType(), sendOnly ? $reflect.SendDir : (recvOnly ? $reflect.RecvDir : $reflect.BothDir));
			};
		};
		break;

	case "Func":
		typ = function(v) { this.$val = v; };
		typ.init = function(params, results, variadic) {
			typ.params = params;
			typ.results = results;
			typ.variadic = variadic;
			typ.extendReflectType = function(rt) {
				var typeSlice = ($sliceType($ptrType($reflect.rtype.Ptr)));
				rt.funcType = new $reflect.funcType.Ptr(rt, variadic, new typeSlice($mapArray(params, function(p) { return p.reflectType(); })), new typeSlice($mapArray(results, function(p) { return p.reflectType(); })));
			};
		};
		break;

	case "Interface":
		typ = { implementedBy: [] };
		typ.init = function(methods) {
			typ.methods = methods;
			typ.extendReflectType = function(rt) {
				var imethods = $mapArray(methods, function(m) {
					return new $reflect.imethod.Ptr($newStringPtr(m[1]), $newStringPtr(m[2]), $funcType(m[3], m[4], m[5]).reflectType());
				});
				var methodSlice = ($sliceType($ptrType($reflect.imethod.Ptr)));
				rt.interfaceType = new $reflect.interfaceType.Ptr(rt, new methodSlice(imethods));
			};
		};
		break;

	case "Map":
		typ = function(v) { this.$val = v; };
		typ.init = function(key, elem) {
			typ.key = key;
			typ.elem = elem;
			typ.extendReflectType = function(rt) {
				rt.mapType = new $reflect.mapType.Ptr(rt, key.reflectType(), elem.reflectType(), undefined, undefined);
			};
		};
		break;

	case "Ptr":
		typ = constructor || function(getter, setter, target) {
			this.$get = getter;
			this.$set = setter;
			this.$target = target;
			this.$val = this;
		};
		typ.prototype.$key = function() {
			if (this.$id === undefined) {
				$idCounter++;
				this.$id = $idCounter;
			}
			return String(this.$id);
		};
		typ.init = function(elem) {
			typ.nil = new typ($throwNilPointerError, $throwNilPointerError);
			typ.extendReflectType = function(rt) {
				rt.ptrType = new $reflect.ptrType.Ptr(rt, elem.reflectType());
			};
		};
		break;

	case "Slice":
		var nativeArray;
		typ = function(array) {
			if (array.constructor !== nativeArray) {
				array = new nativeArray(array);
			}
			this.$array = array;
			this.$offset = 0;
			this.$length = array.length;
			this.$capacity = array.length;
			this.$val = this;
		};
		typ.make = function(length, capacity) {
			capacity = capacity || length;
			var array = new nativeArray(capacity), i;
			if (nativeArray === Array) {
				for (i = 0; i < capacity; i++) {
					array[i] = typ.elem.zero();
				}
			}
			var slice = new typ(array);
			slice.$length = length;
			return slice;
		};
		typ.init = function(elem) {
			typ.elem = elem;
			nativeArray = $nativeArray(elem.kind);
			typ.nil = new typ([]);
			typ.extendReflectType = function(rt) {
				rt.sliceType = new $reflect.sliceType.Ptr(rt, elem.reflectType());
			};
		};
		break;

	case "Struct":
		typ = function(v) { this.$val = v; };
		typ.prototype.$key = function() { $throwRuntimeError("hash of unhashable type " + string); };
		typ.Ptr = $newType(4, "Ptr", "*" + string, "", "", constructor);
		typ.Ptr.Struct = typ;
		typ.Ptr.prototype.$get = function() { return this; };
		typ.init = function(fields) {
			var i;
			typ.fields = fields;
			typ.Ptr.extendReflectType = function(rt) {
				rt.ptrType = new $reflect.ptrType.Ptr(rt, typ.reflectType());
			};
			/* nil value */
			typ.Ptr.nil = Object.create(constructor.prototype);
			typ.Ptr.nil.$val = typ.Ptr.nil;
			for (i = 0; i < fields.length; i++) {
				var field = fields[i];
				Object.defineProperty(typ.Ptr.nil, field[0], { get: $throwNilPointerError, set: $throwNilPointerError });
			}
			/* methods for embedded fields */
			for (i = 0; i < typ.methods.length; i++) {
				var method = typ.methods[i];
				if (method[6] != -1) {
					(function(field, methodName) {
						typ.prototype[methodName] = function() {
							var v = this.$val[field[0]];
							return v[methodName].apply(v, arguments);
						};
					})(fields[method[6]], method[0]);
				}
			}
			for (i = 0; i < typ.Ptr.methods.length; i++) {
				var method = typ.Ptr.methods[i];
				if (method[6] != -1) {
					(function(field, methodName) {
						typ.Ptr.prototype[methodName] = function() {
							var v = this[field[0]];
							if (v.$val === undefined) {
								v = new field[3](v);
							}
							return v[methodName].apply(v, arguments);
						};
					})(fields[method[6]], method[0]);
				}
			}
			/* reflect type */
			typ.extendReflectType = function(rt) {
				var reflectFields = new Array(fields.length), i;
				for (i = 0; i < fields.length; i++) {
					var field = fields[i];
					reflectFields[i] = new $reflect.structField.Ptr($newStringPtr(field[1]), $newStringPtr(field[2]), field[3].reflectType(), $newStringPtr(field[4]), i);
				}
				rt.structType = new $reflect.structType.Ptr(rt, new ($sliceType($reflect.structField.Ptr))(reflectFields));
			};
		};
		break;

	default:
		$panic(new $String("invalid kind: " + kind));
	}

	switch(kind) {
	case "Bool":
	case "Map":
		typ.zero = function() { return false; };
		break;

	case "Int":
	case "Int8":
	case "Int16":
	case "Int32":
	case "Uint":
	case "Uint8" :
	case "Uint16":
	case "Uint32":
	case "Uintptr":
	case "UnsafePointer":
	case "Float32":
	case "Float64":
		typ.zero = function() { return 0; };
		break;

	case "String":
		typ.zero = function() { return ""; };
		break;

	case "Int64":
	case "Uint64":
	case "Complex64":
	case "Complex128":
		var zero = new typ(0, 0);
		typ.zero = function() { return zero; };
		break;

	case "Chan":
	case "Ptr":
	case "Slice":
		typ.zero = function() { return typ.nil; };
		break;

	case "Func":
		typ.zero = function() { return $throwNilPointerError; };
		break;

	case "Interface":
		typ.zero = function() { return null; };
		break;

	case "Array":
		typ.zero = function() {
			var arrayClass = $nativeArray(typ.elem.kind);
			if (arrayClass !== Array) {
				return new arrayClass(typ.len);
			}
			var array = new Array(typ.len), i;
			for (i = 0; i < typ.len; i++) {
				array[i] = typ.elem.zero();
			}
			return array;
		};
		break;

	case "Struct":
		typ.zero = function() { return new typ.Ptr(); };
		break;

	default:
		$panic(new $String("invalid kind: " + kind));
	}

	typ.kind = kind;
	typ.string = string;
	typ.typeName = name;
	typ.pkgPath = pkgPath;
	typ.methods = [];
	var rt = null;
	typ.reflectType = function() {
		if (rt === null) {
			rt = new $reflect.rtype.Ptr(size, 0, 0, 0, 0, $reflect.kinds[kind], undefined, undefined, $newStringPtr(string), undefined, undefined);
			rt.jsType = typ;

			var methods = [];
			if (typ.methods !== undefined) {
				var i;
				for (i = 0; i < typ.methods.length; i++) {
					var m = typ.methods[i];
					methods.push(new $reflect.method.Ptr($newStringPtr(m[1]), $newStringPtr(m[2]), $funcType(m[3], m[4], m[5]).reflectType(), $funcType([typ].concat(m[3]), m[4], m[5]).reflectType(), undefined, undefined));
				}
			}
			if (name !== "" || methods.length !== 0) {
				var methodSlice = ($sliceType($ptrType($reflect.method.Ptr)));
				rt.uncommonType = new $reflect.uncommonType.Ptr($newStringPtr(name), $newStringPtr(pkgPath), new methodSlice(methods));
				rt.uncommonType.jsType = typ;
			}

			if (typ.extendReflectType !== undefined) {
				typ.extendReflectType(rt);
			}
		}
		return rt;
	};
	return typ;
};

var $Bool          = $newType( 1, "Bool",          "bool",           "bool",       "", null);
var $Int           = $newType( 4, "Int",           "int",            "int",        "", null);
var $Int8          = $newType( 1, "Int8",          "int8",           "int8",       "", null);
var $Int16         = $newType( 2, "Int16",         "int16",          "int16",      "", null);
var $Int32         = $newType( 4, "Int32",         "int32",          "int32",      "", null);
var $Int64         = $newType( 8, "Int64",         "int64",          "int64",      "", null);
var $Uint          = $newType( 4, "Uint",          "uint",           "uint",       "", null);
var $Uint8         = $newType( 1, "Uint8",         "uint8",          "uint8",      "", null);
var $Uint16        = $newType( 2, "Uint16",        "uint16",         "uint16",     "", null);
var $Uint32        = $newType( 4, "Uint32",        "uint32",         "uint32",     "", null);
var $Uint64        = $newType( 8, "Uint64",        "uint64",         "uint64",     "", null);
var $Uintptr       = $newType( 4, "Uintptr",       "uintptr",        "uintptr",    "", null);
var $Float32       = $newType( 4, "Float32",       "float32",        "float32",    "", null);
var $Float64       = $newType( 8, "Float64",       "float64",        "float64",    "", null);
var $Complex64     = $newType( 8, "Complex64",     "complex64",      "complex64",  "", null);
var $Complex128    = $newType(16, "Complex128",    "complex128",     "complex128", "", null);
var $String        = $newType( 8, "String",        "string",         "string",     "", null);
var $UnsafePointer = $newType( 4, "UnsafePointer", "unsafe.Pointer", "Pointer",    "", null);

var $nativeArray = function(elemKind) {
	return ({ Int: Int32Array, Int8: Int8Array, Int16: Int16Array, Int32: Int32Array, Uint: Uint32Array, Uint8: Uint8Array, Uint16: Uint16Array, Uint32: Uint32Array, Uintptr: Uint32Array, Float32: Float32Array, Float64: Float64Array })[elemKind] || Array;
};
var $toNativeArray = function(elemKind, array) {
	var nativeArray = $nativeArray(elemKind);
	if (nativeArray === Array) {
		return array;
	}
	return new nativeArray(array);
};
var $arrayTypes = {};
var $arrayType = function(elem, len) {
	var string = "[" + len + "]" + elem.string;
	var typ = $arrayTypes[string];
	if (typ === undefined) {
		typ = $newType(12, "Array", string, "", "", null);
		typ.init(elem, len);
		$arrayTypes[string] = typ;
	}
	return typ;
};

var $chanType = function(elem, sendOnly, recvOnly) {
	var string = (recvOnly ? "<-" : "") + "chan" + (sendOnly ? "<- " : " ") + elem.string;
	var field = sendOnly ? "SendChan" : (recvOnly ? "RecvChan" : "Chan");
	var typ = elem[field];
	if (typ === undefined) {
		typ = $newType(4, "Chan", string, "", "", null);
		typ.init(elem, sendOnly, recvOnly);
		elem[field] = typ;
	}
	return typ;
};

var $funcSig = function(params, results, variadic) {
	var paramTypes = $mapArray(params, function(p) { return p.string; });
	if (variadic) {
		paramTypes[paramTypes.length - 1] = "..." + paramTypes[paramTypes.length - 1].substr(2);
	}
	var string = "(" + paramTypes.join(", ") + ")";
	if (results.length === 1) {
		string += " " + results[0].string;
	} else if (results.length > 1) {
		string += " (" + $mapArray(results, function(r) { return r.string; }).join(", ") + ")";
	}
	return string;
};

var $funcTypes = {};
var $funcType = function(params, results, variadic) {
	var string = "func" + $funcSig(params, results, variadic);
	var typ = $funcTypes[string];
	if (typ === undefined) {
		typ = $newType(4, "Func", string, "", "", null);
		typ.init(params, results, variadic);
		$funcTypes[string] = typ;
	}
	return typ;
};

var $interfaceTypes = {};
var $interfaceType = function(methods) {
	var string = "interface {}";
	if (methods.length !== 0) {
		string = "interface { " + $mapArray(methods, function(m) {
			return (m[2] !== "" ? m[2] + "." : "") + m[1] + $funcSig(m[3], m[4], m[5]);
		}).join("; ") + " }";
	}
	var typ = $interfaceTypes[string];
	if (typ === undefined) {
		typ = $newType(8, "Interface", string, "", "", null);
		typ.init(methods);
		$interfaceTypes[string] = typ;
	}
	return typ;
};
var $emptyInterface = $interfaceType([]);
var $interfaceNil = { $key: function() { return "nil"; } };
var $error = $newType(8, "Interface", "error", "error", "", null);
$error.init([["Error", "Error", "", [], [$String], false]]);

var $Map = function() {};
(function() {
	var names = Object.getOwnPropertyNames(Object.prototype), i;
	for (i = 0; i < names.length; i++) {
		$Map.prototype[names[i]] = undefined;
	}
})();
var $mapTypes = {};
var $mapType = function(key, elem) {
	var string = "map[" + key.string + "]" + elem.string;
	var typ = $mapTypes[string];
	if (typ === undefined) {
		typ = $newType(4, "Map", string, "", "", null);
		typ.init(key, elem);
		$mapTypes[string] = typ;
	}
	return typ;
};

var $throwNilPointerError = function() { $throwRuntimeError("invalid memory address or nil pointer dereference"); };
var $ptrType = function(elem) {
	var typ = elem.Ptr;
	if (typ === undefined) {
		typ = $newType(4, "Ptr", "*" + elem.string, "", "", null);
		typ.init(elem);
		elem.Ptr = typ;
	}
	return typ;
};

var $sliceType = function(elem) {
	var typ = elem.Slice;
	if (typ === undefined) {
		typ = $newType(12, "Slice", "[]" + elem.string, "", "", null);
		typ.init(elem);
		elem.Slice = typ;
	}
	return typ;
};

var $structTypes = {};
var $structType = function(fields) {
	var string = "struct { " + $mapArray(fields, function(f) {
		return f[1] + " " + f[3].string + (f[4] !== "" ? (" \"" + f[4].replace(/\\/g, "\\\\").replace(/"/g, "\\\"") + "\"") : "");
	}).join("; ") + " }";
  if (fields.length === 0) {
  	string = "struct {}";
  }
	var typ = $structTypes[string];
	if (typ === undefined) {
		typ = $newType(0, "Struct", string, "", "", function() {
			this.$val = this;
			var i;
			for (i = 0; i < fields.length; i++) {
				var field = fields[i];
				var arg = arguments[i];
				this[field[0]] = arg !== undefined ? arg : field[3].zero();
			}
		});
		/* collect methods for anonymous fields */
		var i, j;
		for (i = 0; i < fields.length; i++) {
			var field = fields[i];
			if (field[1] === "") {
				var methods = field[3].methods;
				for (j = 0; j < methods.length; j++) {
					var m = methods[j].slice(0, 6).concat([i]);
					typ.methods.push(m);
					typ.Ptr.methods.push(m);
				}
				if (field[3].kind === "Struct") {
					var methods = field[3].Ptr.methods;
					for (j = 0; j < methods.length; j++) {
						typ.Ptr.methods.push(methods[j].slice(0, 6).concat([i]));
					}
				}
			}
		}
		typ.init(fields);
		$structTypes[string] = typ;
	}
	return typ;
};

var $stringPtrMap = new $Map();
$newStringPtr = function(str) {
	if (str === undefined || str === "") {
		return $ptrType($String).nil;
	}
	var ptr = $stringPtrMap[str];
	if (ptr === undefined) {
		ptr = new ($ptrType($String))(function() { return str; }, function(v) { str = v; });
		$stringPtrMap[str] = ptr;
	}
	return ptr;
};
var $newDataPointer = function(data, constructor) {
	if (constructor.Struct) {
		return data;
	}
	return new constructor(function() { return data; }, function(v) { data = v; });
};

var $coerceFloat32 = function(f) {
	var math = $packages["math"];
	if (math === undefined) {
		return f;
	}
	return math.Float32frombits(math.Float32bits(f));
};
var $flatten64 = function(x) {
	return x.$high * 4294967296 + x.$low;
};
var $shiftLeft64 = function(x, y) {
	if (y === 0) {
		return x;
	}
	if (y < 32) {
		return new x.constructor(x.$high << y | x.$low >>> (32 - y), (x.$low << y) >>> 0);
	}
	if (y < 64) {
		return new x.constructor(x.$low << (y - 32), 0);
	}
	return new x.constructor(0, 0);
};
var $shiftRightInt64 = function(x, y) {
	if (y === 0) {
		return x;
	}
	if (y < 32) {
		return new x.constructor(x.$high >> y, (x.$low >>> y | x.$high << (32 - y)) >>> 0);
	}
	if (y < 64) {
		return new x.constructor(x.$high >> 31, (x.$high >> (y - 32)) >>> 0);
	}
	if (x.$high < 0) {
		return new x.constructor(-1, 4294967295);
	}
	return new x.constructor(0, 0);
};
var $shiftRightUint64 = function(x, y) {
	if (y === 0) {
		return x;
	}
	if (y < 32) {
		return new x.constructor(x.$high >>> y, (x.$low >>> y | x.$high << (32 - y)) >>> 0);
	}
	if (y < 64) {
		return new x.constructor(0, x.$high >>> (y - 32));
	}
	return new x.constructor(0, 0);
};
var $mul64 = function(x, y) {
	var high = 0, low = 0, i;
	if ((y.$low & 1) !== 0) {
		high = x.$high;
		low = x.$low;
	}
	for (i = 1; i < 32; i++) {
		if ((y.$low & 1<<i) !== 0) {
			high += x.$high << i | x.$low >>> (32 - i);
			low += (x.$low << i) >>> 0;
		}
	}
	for (i = 0; i < 32; i++) {
		if ((y.$high & 1<<i) !== 0) {
			high += x.$low << i;
		}
	}
	return new x.constructor(high, low);
};
var $div64 = function(x, y, returnRemainder) {
	if (y.$high === 0 && y.$low === 0) {
		$throwRuntimeError("integer divide by zero");
	}

	var s = 1;
	var rs = 1;

	var xHigh = x.$high;
	var xLow = x.$low;
	if (xHigh < 0) {
		s = -1;
		rs = -1;
		xHigh = -xHigh;
		if (xLow !== 0) {
			xHigh--;
			xLow = 4294967296 - xLow;
		}
	}

	var yHigh = y.$high;
	var yLow = y.$low;
	if (y.$high < 0) {
		s *= -1;
		yHigh = -yHigh;
		if (yLow !== 0) {
			yHigh--;
			yLow = 4294967296 - yLow;
		}
	}

	var high = 0, low = 0, n = 0, i;
	while (yHigh < 2147483648 && ((xHigh > yHigh) || (xHigh === yHigh && xLow > yLow))) {
		yHigh = (yHigh << 1 | yLow >>> 31) >>> 0;
		yLow = (yLow << 1) >>> 0;
		n++;
	}
	for (i = 0; i <= n; i++) {
		high = high << 1 | low >>> 31;
		low = (low << 1) >>> 0;
		if ((xHigh > yHigh) || (xHigh === yHigh && xLow >= yLow)) {
			xHigh = xHigh - yHigh;
			xLow = xLow - yLow;
			if (xLow < 0) {
				xHigh--;
				xLow += 4294967296;
			}
			low++;
			if (low === 4294967296) {
				high++;
				low = 0;
			}
		}
		yLow = (yLow >>> 1 | yHigh << (32 - 1)) >>> 0;
		yHigh = yHigh >>> 1;
	}

	if (returnRemainder) {
		return new x.constructor(xHigh * rs, xLow * rs);
	}
	return new x.constructor(high * s, low * s);
};

var $divComplex = function(n, d) {
	var ninf = n.$real === 1/0 || n.$real === -1/0 || n.$imag === 1/0 || n.$imag === -1/0;
	var dinf = d.$real === 1/0 || d.$real === -1/0 || d.$imag === 1/0 || d.$imag === -1/0;
	var nnan = !ninf && (n.$real !== n.$real || n.$imag !== n.$imag);
	var dnan = !dinf && (d.$real !== d.$real || d.$imag !== d.$imag);
	if(nnan || dnan) {
		return new n.constructor(0/0, 0/0);
	}
	if (ninf && !dinf) {
		return new n.constructor(1/0, 1/0);
	}
	if (!ninf && dinf) {
		return new n.constructor(0, 0);
	}
	if (d.$real === 0 && d.$imag === 0) {
		if (n.$real === 0 && n.$imag === 0) {
			return new n.constructor(0/0, 0/0);
		}
		return new n.constructor(1/0, 1/0);
	}
	var a = Math.abs(d.$real);
	var b = Math.abs(d.$imag);
	if (a <= b) {
		var ratio = d.$real / d.$imag;
		var denom = d.$real * ratio + d.$imag;
		return new n.constructor((n.$real * ratio + n.$imag) / denom, (n.$imag * ratio - n.$real) / denom);
	}
	var ratio = d.$imag / d.$real;
	var denom = d.$imag * ratio + d.$real;
	return new n.constructor((n.$imag * ratio + n.$real) / denom, (n.$imag - n.$real * ratio) / denom);
};

var $subslice = function(slice, low, high, max) {
	if (low < 0 || high < low || max < high || high > slice.$capacity || max > slice.$capacity) {
		$throwRuntimeError("slice bounds out of range");
	}
	var s = new slice.constructor(slice.$array);
	s.$offset = slice.$offset + low;
	s.$length = slice.$length - low;
	s.$capacity = slice.$capacity - low;
	if (high !== undefined) {
		s.$length = high - low;
	}
	if (max !== undefined) {
		s.$capacity = max - low;
	}
	return s;
};

var $sliceToArray = function(slice) {
	if (slice.$length === 0) {
		return [];
	}
	if (slice.$array.constructor !== Array) {
		return slice.$array.subarray(slice.$offset, slice.$offset + slice.$length);
	}
	return slice.$array.slice(slice.$offset, slice.$offset + slice.$length);
};

var $decodeRune = function(str, pos) {
	var c0 = str.charCodeAt(pos);

	if (c0 < 0x80) {
		return [c0, 1];
	}

	if (c0 !== c0 || c0 < 0xC0) {
		return [0xFFFD, 1];
	}

	var c1 = str.charCodeAt(pos + 1);
	if (c1 !== c1 || c1 < 0x80 || 0xC0 <= c1) {
		return [0xFFFD, 1];
	}

	if (c0 < 0xE0) {
		var r = (c0 & 0x1F) << 6 | (c1 & 0x3F);
		if (r <= 0x7F) {
			return [0xFFFD, 1];
		}
		return [r, 2];
	}

	var c2 = str.charCodeAt(pos + 2);
	if (c2 !== c2 || c2 < 0x80 || 0xC0 <= c2) {
		return [0xFFFD, 1];
	}

	if (c0 < 0xF0) {
		var r = (c0 & 0x0F) << 12 | (c1 & 0x3F) << 6 | (c2 & 0x3F);
		if (r <= 0x7FF) {
			return [0xFFFD, 1];
		}
		if (0xD800 <= r && r <= 0xDFFF) {
			return [0xFFFD, 1];
		}
		return [r, 3];
	}

	var c3 = str.charCodeAt(pos + 3);
	if (c3 !== c3 || c3 < 0x80 || 0xC0 <= c3) {
		return [0xFFFD, 1];
	}

	if (c0 < 0xF8) {
		var r = (c0 & 0x07) << 18 | (c1 & 0x3F) << 12 | (c2 & 0x3F) << 6 | (c3 & 0x3F);
		if (r <= 0xFFFF || 0x10FFFF < r) {
			return [0xFFFD, 1];
		}
		return [r, 4];
	}

	return [0xFFFD, 1];
};

var $encodeRune = function(r) {
	if (r < 0 || r > 0x10FFFF || (0xD800 <= r && r <= 0xDFFF)) {
		r = 0xFFFD;
	}
	if (r <= 0x7F) {
		return String.fromCharCode(r);
	}
	if (r <= 0x7FF) {
		return String.fromCharCode(0xC0 | r >> 6, 0x80 | (r & 0x3F));
	}
	if (r <= 0xFFFF) {
		return String.fromCharCode(0xE0 | r >> 12, 0x80 | (r >> 6 & 0x3F), 0x80 | (r & 0x3F));
	}
	return String.fromCharCode(0xF0 | r >> 18, 0x80 | (r >> 12 & 0x3F), 0x80 | (r >> 6 & 0x3F), 0x80 | (r & 0x3F));
};

var $stringToBytes = function(str) {
	var array = new Uint8Array(str.length), i;
	for (i = 0; i < str.length; i++) {
		array[i] = str.charCodeAt(i);
	}
	return array;
};

var $bytesToString = function(slice) {
	if (slice.$length === 0) {
		return "";
	}
	var str = "", i;
	for (i = 0; i < slice.$length; i += 10000) {
		str += String.fromCharCode.apply(null, slice.$array.subarray(slice.$offset + i, slice.$offset + Math.min(slice.$length, i + 10000)));
	}
	return str;
};

var $stringToRunes = function(str) {
	var array = new Int32Array(str.length);
	var rune, i, j = 0;
	for (i = 0; i < str.length; i += rune[1], j++) {
		rune = $decodeRune(str, i);
		array[j] = rune[0];
	}
	return array.subarray(0, j);
};

var $runesToString = function(slice) {
	if (slice.$length === 0) {
		return "";
	}
	var str = "", i;
	for (i = 0; i < slice.$length; i++) {
		str += $encodeRune(slice.$array[slice.$offset + i]);
	}
	return str;
};

var $needsExternalization = function(t) {
	switch (t.kind) {
		case "Bool":
		case "Int":
		case "Int8":
		case "Int16":
		case "Int32":
		case "Uint":
		case "Uint8":
		case "Uint16":
		case "Uint32":
		case "Uintptr":
		case "Float32":
		case "Float64":
			return false;
		case "Interface":
			return t !== $packages["github.com/gopherjs/gopherjs/js"].Object;
		default:
			return true;
	}
};

var $externalize = function(v, t) {
	switch (t.kind) {
	case "Bool":
	case "Int":
	case "Int8":
	case "Int16":
	case "Int32":
	case "Uint":
	case "Uint8":
	case "Uint16":
	case "Uint32":
	case "Uintptr":
	case "Float32":
	case "Float64":
		return v;
	case "Int64":
	case "Uint64":
		return $flatten64(v);
	case "Array":
		if ($needsExternalization(t.elem)) {
			return $mapArray(v, function(e) { return $externalize(e, t.elem); });
		}
		return v;
	case "Func":
		if (v === $throwNilPointerError) {
			return null;
		}
		$checkForDeadlock = false;
		var convert = false;
		var i;
		for (i = 0; i < t.params.length; i++) {
			convert = convert || (t.params[i] !== $packages["github.com/gopherjs/gopherjs/js"].Object);
		}
		for (i = 0; i < t.results.length; i++) {
			convert = convert || $needsExternalization(t.results[i]);
		}
		if (!convert) {
			return v;
		}
		return function() {
			var args = [], i;
			for (i = 0; i < t.params.length; i++) {
				if (t.variadic && i === t.params.length - 1) {
					var vt = t.params[i].elem, varargs = [], j;
					for (j = i; j < arguments.length; j++) {
						varargs.push($internalize(arguments[j], vt));
					}
					args.push(new (t.params[i])(varargs));
					break;
				}
				args.push($internalize(arguments[i], t.params[i]));
			}
			var result = v.apply(this, args);
			switch (t.results.length) {
			case 0:
				return;
			case 1:
				return $externalize(result, t.results[0]);
			default:
				for (i = 0; i < t.results.length; i++) {
					result[i] = $externalize(result[i], t.results[i]);
				}
				return result;
			}
		};
	case "Interface":
		if (v === null) {
			return null;
		}
		if (t === $packages["github.com/gopherjs/gopherjs/js"].Object || v.constructor.kind === undefined) {
			return v;
		}
		return $externalize(v.$val, v.constructor);
	case "Map":
		var m = {};
		var keys = $keys(v), i;
		for (i = 0; i < keys.length; i++) {
			var entry = v[keys[i]];
			m[$externalize(entry.k, t.key)] = $externalize(entry.v, t.elem);
		}
		return m;
	case "Ptr":
		var o = {}, i;
		for (i = 0; i < t.methods.length; i++) {
			var m = t.methods[i];
			if (m[2] !== "") { /* not exported */
				continue;
			}
			(function(m) {
				o[m[1]] = $externalize(function() {
					return v[m[0]].apply(v, arguments);
				}, $funcType(m[3], m[4], m[5]));
			})(m);
		}
		return o;
	case "Slice":
		if ($needsExternalization(t.elem)) {
			return $mapArray($sliceToArray(v), function(e) { return $externalize(e, t.elem); });
		}
		return $sliceToArray(v);
	case "String":
		var s = "", r, i, j = 0;
		for (i = 0; i < v.length; i += r[1], j++) {
			r = $decodeRune(v, i);
			s += String.fromCharCode(r[0]);
		}
		return s;
	case "Struct":
		var timePkg = $packages["time"];
		if (timePkg && v.constructor === timePkg.Time.Ptr) {
			var milli = $div64(v.UnixNano(), new $Int64(0, 1000000));
			return new Date($flatten64(milli));
		}
		var o = {}, i;
		for (i = 0; i < t.fields.length; i++) {
			var f = t.fields[i];
			if (f[2] !== "") { /* not exported */
				continue;
			}
			o[f[1]] = $externalize(v[f[0]], f[3]);
		}
		return o;
	}
	$panic(new $String("cannot externalize " + t.string));
};

var $internalize = function(v, t, recv) {
	switch (t.kind) {
	case "Bool":
		return !!v;
	case "Int":
		return parseInt(v);
	case "Int8":
		return parseInt(v) << 24 >> 24;
	case "Int16":
		return parseInt(v) << 16 >> 16;
	case "Int32":
		return parseInt(v) >> 0;
	case "Uint":
		return parseInt(v);
	case "Uint8":
		return parseInt(v) << 24 >>> 24;
	case "Uint16":
		return parseInt(v) << 16 >>> 16;
	case "Uint32":
	case "Uintptr":
		return parseInt(v) >>> 0;
	case "Int64":
	case "Uint64":
		return new t(0, v);
	case "Float32":
	case "Float64":
		return parseFloat(v);
	case "Array":
		if (v.length !== t.len) {
			$throwRuntimeError("got array with wrong size from JavaScript native");
		}
		return $mapArray(v, function(e) { return $internalize(e, t.elem); });
	case "Func":
		return function() {
			var args = [], i;
			for (i = 0; i < t.params.length; i++) {
				if (t.variadic && i === t.params.length - 1) {
					var vt = t.params[i].elem, varargs = arguments[i], j;
					for (j = 0; j < varargs.$length; j++) {
						args.push($externalize(varargs.$array[varargs.$offset + j], vt));
					}
					break;
				}
				args.push($externalize(arguments[i], t.params[i]));
			}
			var result = v.apply(recv, args);
			switch (t.results.length) {
			case 0:
				return;
			case 1:
				return $internalize(result, t.results[0]);
			default:
				for (i = 0; i < t.results.length; i++) {
					result[i] = $internalize(result[i], t.results[i]);
				}
				return result;
			}
		};
	case "Interface":
		if (v === null || t === $packages["github.com/gopherjs/gopherjs/js"].Object) {
			return v;
		}
		switch (v.constructor) {
		case Int8Array:
			return new ($sliceType($Int8))(v);
		case Int16Array:
			return new ($sliceType($Int16))(v);
		case Int32Array:
			return new ($sliceType($Int))(v);
		case Uint8Array:
			return new ($sliceType($Uint8))(v);
		case Uint16Array:
			return new ($sliceType($Uint16))(v);
		case Uint32Array:
			return new ($sliceType($Uint))(v);
		case Float32Array:
			return new ($sliceType($Float32))(v);
		case Float64Array:
			return new ($sliceType($Float64))(v);
		case Array:
			return $internalize(v, $sliceType($emptyInterface));
		case Boolean:
			return new $Bool(!!v);
		case Date:
			var timePkg = $packages["time"];
			if (timePkg) {
				return new timePkg.Time(timePkg.Unix(new $Int64(0, 0), new $Int64(0, v.getTime() * 1000000)));
			}
		case Function:
			var funcType = $funcType([$sliceType($emptyInterface)], [$packages["github.com/gopherjs/gopherjs/js"].Object], true);
			return new funcType($internalize(v, funcType));
		case Number:
			return new $Float64(parseFloat(v));
		case String:
			return new $String($internalize(v, $String));
		default:
			var mapType = $mapType($String, $emptyInterface);
			return new mapType($internalize(v, mapType));
		}
	case "Map":
		var m = new $Map();
		var keys = $keys(v), i;
		for (i = 0; i < keys.length; i++) {
			var key = $internalize(keys[i], t.key);
			m[key.$key ? key.$key() : key] = { k: key, v: $internalize(v[keys[i]], t.elem) };
		}
		return m;
	case "Slice":
		return new t($mapArray(v, function(e) { return $internalize(e, t.elem); }));
	case "String":
		v = String(v);
		var s = "", i;
		for (i = 0; i < v.length; i++) {
			s += $encodeRune(v.charCodeAt(i));
		}
		return s;
	default:
		$panic(new $String("cannot internalize " + t.string));
	}
};

var $copyString = function(dst, src) {
	var n = Math.min(src.length, dst.$length), i;
	for (i = 0; i < n; i++) {
		dst.$array[dst.$offset + i] = src.charCodeAt(i);
	}
	return n;
};

var $copySlice = function(dst, src) {
	var n = Math.min(src.$length, dst.$length), i;
	$internalCopy(dst.$array, src.$array, dst.$offset, src.$offset, n, dst.constructor.elem);
	return n;
};

var $copy = function(dst, src, type) {
	var i;
	switch (type.kind) {
	case "Array":
		$internalCopy(dst, src, 0, 0, src.length, type.elem);
		return true;
	case "Struct":
		for (i = 0; i < type.fields.length; i++) {
			var field = type.fields[i];
			var name = field[0];
			if (!$copy(dst[name], src[name], field[3])) {
				dst[name] = src[name];
			}
		}
		return true;
	default:
		return false;
	}
};

var $internalCopy = function(dst, src, dstOffset, srcOffset, n, elem) {
	var i;
	if (n === 0) {
		return;
	}

	if (src.subarray) {
		dst.set(src.subarray(srcOffset, srcOffset + n), dstOffset);
		return;
	}

	switch (elem.kind) {
	case "Array":
	case "Struct":
		for (i = 0; i < n; i++) {
			$copy(dst[dstOffset + i], src[srcOffset + i], elem);
		}
		return;
	}

	for (i = 0; i < n; i++) {
		dst[dstOffset + i] = src[srcOffset + i];
	}
};

var $clone = function(src, type) {
	var clone = type.zero();
	$copy(clone, src, type);
	return clone;
};

var $append = function(slice) {
	return $internalAppend(slice, arguments, 1, arguments.length - 1);
};

var $appendSlice = function(slice, toAppend) {
	return $internalAppend(slice, toAppend.$array, toAppend.$offset, toAppend.$length);
};

var $internalAppend = function(slice, array, offset, length) {
	if (length === 0) {
		return slice;
	}

	var newArray = slice.$array;
	var newOffset = slice.$offset;
	var newLength = slice.$length + length;
	var newCapacity = slice.$capacity;

	if (newLength > newCapacity) {
		newOffset = 0;
		newCapacity = Math.max(newLength, slice.$capacity < 1024 ? slice.$capacity * 2 : Math.floor(slice.$capacity * 5 / 4));

		if (slice.$array.constructor === Array) {
			newArray = slice.$array.slice(slice.$offset, slice.$offset + slice.$length);
			newArray.length = newCapacity;
			var zero = slice.constructor.elem.zero, i;
			for (i = slice.$length; i < newCapacity; i++) {
				newArray[i] = zero();
			}
		} else {
			newArray = new slice.$array.constructor(newCapacity);
			newArray.set(slice.$array.subarray(slice.$offset, slice.$offset + slice.$length));
		}
	}

	$internalCopy(newArray, array, newOffset + slice.$length, offset, length, slice.constructor.elem);

	var newSlice = new slice.constructor(newArray);
	newSlice.$offset = newOffset;
	newSlice.$length = newLength;
	newSlice.$capacity = newCapacity;
	return newSlice;
};

var $getStack = function() {
	return (new Error()).stack.split("\n");
};
var $stackDepthOffset = 0;
var $getStackDepth = function() {
	return $stackDepthOffset + $getStack().length;
};

var $deferFrames = [], $skippedDeferFrames = 0, $jumpToDefer = false, $panicStackDepth = null, $panicValue;
var $callDeferred = function(deferred, jsErr) {
	if ($skippedDeferFrames !== 0) {
		$skippedDeferFrames--;
		throw jsErr;
	}
	if ($jumpToDefer) {
		$jumpToDefer = false;
		throw jsErr;
	}

	$stackDepthOffset--;
	var outerPanicStackDepth = $panicStackDepth;
	var outerPanicValue = $panicValue;

	var localPanicValue = $curGoroutine.panicStack.pop();
	if (jsErr) {
		localPanicValue = new $packages["github.com/gopherjs/gopherjs/js"].Error.Ptr(jsErr);
	}
	if (localPanicValue !== undefined) {
		$panicStackDepth = $getStackDepth();
		$panicValue = localPanicValue;
	}

	var call;
	try {
		while (true) {
			if (deferred === null) {
				deferred = $deferFrames[$deferFrames.length - 1 - $skippedDeferFrames];
				if (deferred === undefined) {
					if (localPanicValue.constructor === $String) {
						throw new Error(localPanicValue.$val);
					} else if (localPanicValue.Error !== undefined) {
						throw new Error(localPanicValue.Error());
					} else if (localPanicValue.String !== undefined) {
						throw new Error(localPanicValue.String());
					} else {
						throw new Error(localPanicValue);
					}
				}
			}
			var call = deferred.pop();
			if (call === undefined) {
				if (localPanicValue !== undefined) {
					$skippedDeferFrames++;
					deferred = null;
					continue;
				}
				return;
			}
			var r = call[0].apply(undefined, call[1]);
		  if (r && r.constructor === Function) {
				deferred.push([r, []]);
			}

			if (localPanicValue !== undefined && $panicStackDepth === null) {
				throw null; /* error was recovered */
			}
		}
	} finally {
		if ($curGoroutine.asleep) {
			deferred.push(call);
			$jumpToDefer = true;
		}
		if (localPanicValue !== undefined) {
			if ($panicStackDepth !== null) {
				$curGoroutine.panicStack.push(localPanicValue);
			}
			$panicStackDepth = outerPanicStackDepth;
			$panicValue = outerPanicValue;
		}
		$stackDepthOffset++;
	}
};

var $panic = function(value) {
	$curGoroutine.panicStack.push(value);
	$callDeferred(null, null);
};
var $recover = function() {
	if ($panicStackDepth === null || $panicStackDepth !== $getStackDepth() - 2) {
		return null;
	}
	$panicStackDepth = null;
	return $panicValue;
};
var $nonblockingCall = function() {
	$panic(new $packages["runtime"].NotSupportedError.Ptr("non-blocking call to blocking function (mark call with \"//gopherjs:blocking\" to fix)"));
};
var $throw = function(err) { throw err; };
var $throwRuntimeError; /* set by package "runtime" */

var $dummyGoroutine = { asleep: false, exit: false, panicStack: [] };
var $curGoroutine = $dummyGoroutine, $totalGoroutines = 0, $awakeGoroutines = 0, $checkForDeadlock = true;
var $go = function(fun, args, direct) {
	$totalGoroutines++;
	$awakeGoroutines++;
	args.push(true);
  var goroutine = function() {
	  try {
			$curGoroutine = goroutine;
			$skippedDeferFrames = 0;
			$jumpToDefer = false;
			var r = fun.apply(undefined, args);
			if (r !== undefined) {
				fun = r;
				args = [];
				$schedule(goroutine, direct);
				return;
			}
			goroutine.exit = true;
		} catch (err) {
			if (!$curGoroutine.asleep) {
				goroutine.exit = true;
				throw err;
			}
		} finally {
			$curGoroutine = $dummyGoroutine;
			if (goroutine.exit) { /* also set by runtime.Goexit() */
				$totalGoroutines--;
				goroutine.asleep = true;
			}
			if (goroutine.asleep) {
				$awakeGoroutines--;
				if ($awakeGoroutines === 0 && $totalGoroutines !== 0 && $checkForDeadlock) {
					$panic(new $String("fatal error: all goroutines are asleep - deadlock!"));
				}
			}
		}
	};
	goroutine.asleep = false;
	goroutine.exit = false;
	goroutine.panicStack = [];
	$schedule(goroutine, direct);
};

var $scheduled = [], $schedulerLoopActive = false;
var $schedule = function(goroutine, direct) {
	if (goroutine.asleep) {
		goroutine.asleep = false;
		$awakeGoroutines++;
	}

	if (direct) {
		goroutine();
		return;
	}

	$scheduled.push(goroutine);
	if (!$schedulerLoopActive) {
		$schedulerLoopActive = true;
		setTimeout(function() {
			while (true) {
				var r = $scheduled.shift();
				if (r === undefined) {
					$schedulerLoopActive = false;
					break;
				}
				r();
			};
		}, 0);
	}
};

var $send = function(chan, value) {
	if (chan.$closed) {
		$throwRuntimeError("send on closed channel");
	}
	var queuedRecv = chan.$recvQueue.shift();
	if (queuedRecv !== undefined) {
		queuedRecv.chanValue = [value, true];
		$schedule(queuedRecv);
		return;
	}
	if (chan.$buffer.length < chan.$capacity) {
		chan.$buffer.push(value);
		return;
	}

	chan.$sendQueue.push([$curGoroutine, value]);
	var blocked = false;
	return function() {
		if (blocked) {
			if (chan.$closed) {
				$throwRuntimeError("send on closed channel");
			}
			return;
		};
		blocked = true;
		$curGoroutine.asleep = true;
		throw null;
	};
};
var $recv = function(chan) {
	var queuedSend = chan.$sendQueue.shift();
	if (queuedSend !== undefined) {
		$schedule(queuedSend[0]);
		chan.$buffer.push(queuedSend[1]);
	}
	var bufferedValue = chan.$buffer.shift();
	if (bufferedValue !== undefined) {
		return [bufferedValue, true];
	}
	if (chan.$closed) {
		return [chan.constructor.elem.zero(), false];
	}

	chan.$recvQueue.push($curGoroutine);
	var blocked = false;
	return function() {
		if (blocked) {
			var value = $curGoroutine.chanValue;
			$curGoroutine.chanValue = undefined;
			return value;
		};
		blocked = true;
		$curGoroutine.asleep = true;
		throw null;
	};
};
var $close = function(chan) {
	if (chan.$closed) {
		$throwRuntimeError("close of closed channel");
	}
	chan.$closed = true;
	while (true) {
		var queuedSend = chan.$sendQueue.shift();
		if (queuedSend === undefined) {
			break;
		}
		$schedule(queuedSend[0]); /* will panic because of closed channel */
	}
	while (true) {
		var queuedRecv = chan.$recvQueue.shift();
		if (queuedRecv === undefined) {
			break;
		}
		queuedRecv.chanValue = [chan.constructor.elem.zero(), false];
		$schedule(queuedRecv);
	}
};
var $select = function(comms) {
	var ready = [], i;
	var selection = -1;
	for (i = 0; i < comms.length; i++) {
		var comm = comms[i];
		var chan = comm[0];
		switch (comm.length) {
		case 0: /* default */
			selection = i;
			break;
		case 1: /* recv */
			if (chan.$sendQueue.length !== 0 || chan.$buffer.length !== 0 || chan.$closed) {
				ready.push(i);
			}
			break;
		case 2: /* send */
			if (chan.$closed) {
				$throwRuntimeError("send on closed channel");
			}
			if (chan.$recvQueue.length !== 0 || chan.$buffer.length < chan.$capacity) {
				ready.push(i);
			}
			break;
		}
	}

	if (ready.length !== 0) {
		selection = ready[Math.floor(Math.random() * ready.length)];
	}
	if (selection !== -1) {
		var comm = comms[selection];
		switch (comm.length) {
		case 0: /* default */
			return [selection];
		case 1: /* recv */
			return [selection, $recv(comm[0])];
		case 2: /* send */
			$send(comm[0], comm[1]);
			return [selection];
		}
	}

	for (i = 0; i < comms.length; i++) {
		var comm = comms[i];
		switch (comm.length) {
		case 1: /* recv */
			comm[0].$recvQueue.push($curGoroutine);
			break;
		case 2: /* send */
			var queueEntry = [$curGoroutine, comm[1]];
			comm.push(queueEntry);
			comm[0].$sendQueue.push(queueEntry);
			break;
		}
	}
	var blocked = false;
	return function() {
		if (blocked) {
			var selection;
			for (i = 0; i < comms.length; i++) {
				var comm = comms[i];
				switch (comm.length) {
				case 1: /* recv */
					var queue = comm[0].$recvQueue;
					var index = queue.indexOf($curGoroutine);
					if (index !== -1) {
						queue.splice(index, 1);
						break;
					}
					var value = $curGoroutine.chanValue;
					$curGoroutine.chanValue = undefined;
					selection = [i, value];
					break;
				case 3: /* send */
					var queue = comm[0].$sendQueue;
					var index = queue.indexOf(comm[2]);
					if (index !== -1) {
						queue.splice(index, 1);
						break;
					}
					if (comm[0].$closed) {
						$throwRuntimeError("send on closed channel");
					}
					selection = [i];
					break;
				}
			}
			return selection;
		};
		blocked = true;
		$curGoroutine.asleep = true;
		throw null;
	};
};

var $equal = function(a, b, type) {
	if (a === b) {
		return true;
	}
	var i;
	switch (type.kind) {
	case "Float32":
		return $float32IsEqual(a, b);
	case "Complex64":
		return $float32IsEqual(a.$real, b.$real) && $float32IsEqual(a.$imag, b.$imag);
	case "Complex128":
		return a.$real === b.$real && a.$imag === b.$imag;
	case "Int64":
	case "Uint64":
		return a.$high === b.$high && a.$low === b.$low;
	case "Ptr":
		if (a.constructor.Struct) {
			return false;
		}
		return $pointerIsEqual(a, b);
	case "Array":
		if (a.length != b.length) {
			return false;
		}
		var i;
		for (i = 0; i < a.length; i++) {
			if (!$equal(a[i], b[i], type.elem)) {
				return false;
			}
		}
		return true;
	case "Struct":
		for (i = 0; i < type.fields.length; i++) {
			var field = type.fields[i];
			var name = field[0];
			if (!$equal(a[name], b[name], field[3])) {
				return false;
			}
		}
		return true;
	default:
		return false;
	}
};
var $interfaceIsEqual = function(a, b) {
	if (a === null || b === null || a === undefined || b === undefined || a.constructor !== b.constructor) {
		return a === b;
	}
	switch (a.constructor.kind) {
	case "Func":
	case "Map":
	case "Slice":
	case "Struct":
		$throwRuntimeError("comparing uncomparable type " + a.constructor.string);
	case undefined: /* js.Object */
		return a === b;
	default:
		return $equal(a.$val, b.$val, a.constructor);
	}
};
var $float32IsEqual = function(a, b) {
	if (a === b) {
		return true;
	}
	if (a === 0 || b === 0 || a === 1/0 || b === 1/0 || a === -1/0 || b === -1/0 || a !== a || b !== b) {
		return false;
	}
	var math = $packages["math"];
	return math !== undefined && math.Float32bits(a) === math.Float32bits(b);
};
var $sliceIsEqual = function(a, ai, b, bi) {
	return a.$array === b.$array && a.$offset + ai === b.$offset + bi;
};
var $pointerIsEqual = function(a, b) {
	if (a === b) {
		return true;
	}
	if (a.$get === $throwNilPointerError || b.$get === $throwNilPointerError) {
		return a.$get === $throwNilPointerError && b.$get === $throwNilPointerError;
	}
	var old = a.$get();
	var dummy = new Object();
	a.$set(dummy);
	var equal = b.$get() === dummy;
	a.$set(old);
	return equal;
};

var $typeAssertionFailed = function(obj, expected) {
	var got = "";
	if (obj !== null) {
		got = obj.constructor.string;
	}
	$panic(new $packages["runtime"].TypeAssertionError.Ptr("", got, expected.string, ""));
};

var $packages = {};
$packages["github.com/gopherjs/gopherjs/js"] = (function() {
	var $pkg = {}, Object, Error, init;
	Object = $pkg.Object = $newType(8, "Interface", "js.Object", "Object", "github.com/gopherjs/gopherjs/js", null);
	Error = $pkg.Error = $newType(0, "Struct", "js.Error", "Error", "github.com/gopherjs/gopherjs/js", function(Object_) {
		this.$val = this;
		this.Object = Object_ !== undefined ? Object_ : null;
	});
	Error.Ptr.prototype.Error = function() {
		var err;
		err = this;
		return "JavaScript error: " + $internalize(err.Object.message, $String);
	};
	Error.prototype.Error = function() { return this.$val.Error(); };
	init = function() {
		var e;
		e = new Error.Ptr(null);
	};
	$pkg.$init = function() {
		Object.init([["Bool", "Bool", "", [], [$Bool], false], ["Call", "Call", "", [$String, ($sliceType($emptyInterface))], [Object], true], ["Delete", "Delete", "", [$String], [], false], ["Float", "Float", "", [], [$Float64], false], ["Get", "Get", "", [$String], [Object], false], ["Index", "Index", "", [$Int], [Object], false], ["Int", "Int", "", [], [$Int], false], ["Int64", "Int64", "", [], [$Int64], false], ["Interface", "Interface", "", [], [$emptyInterface], false], ["Invoke", "Invoke", "", [($sliceType($emptyInterface))], [Object], true], ["IsNull", "IsNull", "", [], [$Bool], false], ["IsUndefined", "IsUndefined", "", [], [$Bool], false], ["Length", "Length", "", [], [$Int], false], ["New", "New", "", [($sliceType($emptyInterface))], [Object], true], ["Set", "Set", "", [$String, $emptyInterface], [], false], ["SetIndex", "SetIndex", "", [$Int, $emptyInterface], [], false], ["Str", "Str", "", [], [$String], false], ["Uint64", "Uint64", "", [], [$Uint64], false], ["Unsafe", "Unsafe", "", [], [$Uintptr], false]]);
		Error.methods = [["Bool", "Bool", "", [], [$Bool], false, 0], ["Call", "Call", "", [$String, ($sliceType($emptyInterface))], [Object], true, 0], ["Delete", "Delete", "", [$String], [], false, 0], ["Float", "Float", "", [], [$Float64], false, 0], ["Get", "Get", "", [$String], [Object], false, 0], ["Index", "Index", "", [$Int], [Object], false, 0], ["Int", "Int", "", [], [$Int], false, 0], ["Int64", "Int64", "", [], [$Int64], false, 0], ["Interface", "Interface", "", [], [$emptyInterface], false, 0], ["Invoke", "Invoke", "", [($sliceType($emptyInterface))], [Object], true, 0], ["IsNull", "IsNull", "", [], [$Bool], false, 0], ["IsUndefined", "IsUndefined", "", [], [$Bool], false, 0], ["Length", "Length", "", [], [$Int], false, 0], ["New", "New", "", [($sliceType($emptyInterface))], [Object], true, 0], ["Set", "Set", "", [$String, $emptyInterface], [], false, 0], ["SetIndex", "SetIndex", "", [$Int, $emptyInterface], [], false, 0], ["Str", "Str", "", [], [$String], false, 0], ["Uint64", "Uint64", "", [], [$Uint64], false, 0], ["Unsafe", "Unsafe", "", [], [$Uintptr], false, 0]];
		($ptrType(Error)).methods = [["Bool", "Bool", "", [], [$Bool], false, 0], ["Call", "Call", "", [$String, ($sliceType($emptyInterface))], [Object], true, 0], ["Delete", "Delete", "", [$String], [], false, 0], ["Error", "Error", "", [], [$String], false, -1], ["Float", "Float", "", [], [$Float64], false, 0], ["Get", "Get", "", [$String], [Object], false, 0], ["Index", "Index", "", [$Int], [Object], false, 0], ["Int", "Int", "", [], [$Int], false, 0], ["Int64", "Int64", "", [], [$Int64], false, 0], ["Interface", "Interface", "", [], [$emptyInterface], false, 0], ["Invoke", "Invoke", "", [($sliceType($emptyInterface))], [Object], true, 0], ["IsNull", "IsNull", "", [], [$Bool], false, 0], ["IsUndefined", "IsUndefined", "", [], [$Bool], false, 0], ["Length", "Length", "", [], [$Int], false, 0], ["New", "New", "", [($sliceType($emptyInterface))], [Object], true, 0], ["Set", "Set", "", [$String, $emptyInterface], [], false, 0], ["SetIndex", "SetIndex", "", [$Int, $emptyInterface], [], false, 0], ["Str", "Str", "", [], [$String], false, 0], ["Uint64", "Uint64", "", [], [$Uint64], false, 0], ["Unsafe", "Unsafe", "", [], [$Uintptr], false, 0]];
		Error.init([["Object", "", "", Object, ""]]);
		init();
	};
	return $pkg;
})();
$packages["runtime"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], NotSupportedError, TypeAssertionError, errorString, MemStats, sizeof_C_MStats, init, getgoroot, SetFinalizer, GOROOT, init$1;
	NotSupportedError = $pkg.NotSupportedError = $newType(0, "Struct", "runtime.NotSupportedError", "NotSupportedError", "runtime", function(Feature_) {
		this.$val = this;
		this.Feature = Feature_ !== undefined ? Feature_ : "";
	});
	TypeAssertionError = $pkg.TypeAssertionError = $newType(0, "Struct", "runtime.TypeAssertionError", "TypeAssertionError", "runtime", function(interfaceString_, concreteString_, assertedString_, missingMethod_) {
		this.$val = this;
		this.interfaceString = interfaceString_ !== undefined ? interfaceString_ : "";
		this.concreteString = concreteString_ !== undefined ? concreteString_ : "";
		this.assertedString = assertedString_ !== undefined ? assertedString_ : "";
		this.missingMethod = missingMethod_ !== undefined ? missingMethod_ : "";
	});
	errorString = $pkg.errorString = $newType(8, "String", "runtime.errorString", "errorString", "runtime", null);
	MemStats = $pkg.MemStats = $newType(0, "Struct", "runtime.MemStats", "MemStats", "runtime", function(Alloc_, TotalAlloc_, Sys_, Lookups_, Mallocs_, Frees_, HeapAlloc_, HeapSys_, HeapIdle_, HeapInuse_, HeapReleased_, HeapObjects_, StackInuse_, StackSys_, MSpanInuse_, MSpanSys_, MCacheInuse_, MCacheSys_, BuckHashSys_, GCSys_, OtherSys_, NextGC_, LastGC_, PauseTotalNs_, PauseNs_, NumGC_, EnableGC_, DebugGC_, BySize_) {
		this.$val = this;
		this.Alloc = Alloc_ !== undefined ? Alloc_ : new $Uint64(0, 0);
		this.TotalAlloc = TotalAlloc_ !== undefined ? TotalAlloc_ : new $Uint64(0, 0);
		this.Sys = Sys_ !== undefined ? Sys_ : new $Uint64(0, 0);
		this.Lookups = Lookups_ !== undefined ? Lookups_ : new $Uint64(0, 0);
		this.Mallocs = Mallocs_ !== undefined ? Mallocs_ : new $Uint64(0, 0);
		this.Frees = Frees_ !== undefined ? Frees_ : new $Uint64(0, 0);
		this.HeapAlloc = HeapAlloc_ !== undefined ? HeapAlloc_ : new $Uint64(0, 0);
		this.HeapSys = HeapSys_ !== undefined ? HeapSys_ : new $Uint64(0, 0);
		this.HeapIdle = HeapIdle_ !== undefined ? HeapIdle_ : new $Uint64(0, 0);
		this.HeapInuse = HeapInuse_ !== undefined ? HeapInuse_ : new $Uint64(0, 0);
		this.HeapReleased = HeapReleased_ !== undefined ? HeapReleased_ : new $Uint64(0, 0);
		this.HeapObjects = HeapObjects_ !== undefined ? HeapObjects_ : new $Uint64(0, 0);
		this.StackInuse = StackInuse_ !== undefined ? StackInuse_ : new $Uint64(0, 0);
		this.StackSys = StackSys_ !== undefined ? StackSys_ : new $Uint64(0, 0);
		this.MSpanInuse = MSpanInuse_ !== undefined ? MSpanInuse_ : new $Uint64(0, 0);
		this.MSpanSys = MSpanSys_ !== undefined ? MSpanSys_ : new $Uint64(0, 0);
		this.MCacheInuse = MCacheInuse_ !== undefined ? MCacheInuse_ : new $Uint64(0, 0);
		this.MCacheSys = MCacheSys_ !== undefined ? MCacheSys_ : new $Uint64(0, 0);
		this.BuckHashSys = BuckHashSys_ !== undefined ? BuckHashSys_ : new $Uint64(0, 0);
		this.GCSys = GCSys_ !== undefined ? GCSys_ : new $Uint64(0, 0);
		this.OtherSys = OtherSys_ !== undefined ? OtherSys_ : new $Uint64(0, 0);
		this.NextGC = NextGC_ !== undefined ? NextGC_ : new $Uint64(0, 0);
		this.LastGC = LastGC_ !== undefined ? LastGC_ : new $Uint64(0, 0);
		this.PauseTotalNs = PauseTotalNs_ !== undefined ? PauseTotalNs_ : new $Uint64(0, 0);
		this.PauseNs = PauseNs_ !== undefined ? PauseNs_ : ($arrayType($Uint64, 256)).zero();
		this.NumGC = NumGC_ !== undefined ? NumGC_ : 0;
		this.EnableGC = EnableGC_ !== undefined ? EnableGC_ : false;
		this.DebugGC = DebugGC_ !== undefined ? DebugGC_ : false;
		this.BySize = BySize_ !== undefined ? BySize_ : ($arrayType(($structType([["Size", "Size", "", $Uint32, ""], ["Mallocs", "Mallocs", "", $Uint64, ""], ["Frees", "Frees", "", $Uint64, ""]])), 61)).zero();
	});
	NotSupportedError.Ptr.prototype.Error = function() {
		var err;
		err = this;
		return "not supported by GopherJS: " + err.Feature;
	};
	NotSupportedError.prototype.Error = function() { return this.$val.Error(); };
	init = function() {
		var e;
		$throwRuntimeError = $externalize((function(msg) {
			$panic(new errorString(msg));
		}), ($funcType([$String], [], false)));
		e = null;
		e = new TypeAssertionError.Ptr("", "", "", "");
		e = new NotSupportedError.Ptr("");
	};
	getgoroot = function() {
		var process, goroot;
		process = $global.process;
		if (process === undefined) {
			return "/";
		}
		goroot = process.env.GOROOT;
		if (goroot === undefined) {
			return "";
		}
		return $internalize(goroot, $String);
	};
	SetFinalizer = $pkg.SetFinalizer = function(x, f) {
	};
	TypeAssertionError.Ptr.prototype.RuntimeError = function() {
	};
	TypeAssertionError.prototype.RuntimeError = function() { return this.$val.RuntimeError(); };
	TypeAssertionError.Ptr.prototype.Error = function() {
		var e, inter;
		e = this;
		inter = e.interfaceString;
		if (inter === "") {
			inter = "interface";
		}
		if (e.concreteString === "") {
			return "interface conversion: " + inter + " is nil, not " + e.assertedString;
		}
		if (e.missingMethod === "") {
			return "interface conversion: " + inter + " is " + e.concreteString + ", not " + e.assertedString;
		}
		return "interface conversion: " + e.concreteString + " is not " + e.assertedString + ": missing method " + e.missingMethod;
	};
	TypeAssertionError.prototype.Error = function() { return this.$val.Error(); };
	errorString.prototype.RuntimeError = function() {
		var e;
		e = this.$val;
	};
	$ptrType(errorString).prototype.RuntimeError = function() { return new errorString(this.$get()).RuntimeError(); };
	errorString.prototype.Error = function() {
		var e;
		e = this.$val;
		return "runtime error: " + e;
	};
	$ptrType(errorString).prototype.Error = function() { return new errorString(this.$get()).Error(); };
	GOROOT = $pkg.GOROOT = function() {
		var s;
		s = getgoroot();
		if (!(s === "")) {
			return s;
		}
		return "/usr/lib/go";
	};
	init$1 = function() {
		var memStats;
		memStats = new MemStats.Ptr(); $copy(memStats, new MemStats.Ptr(), MemStats);
		if (!((sizeof_C_MStats === 3712))) {
			console.log(sizeof_C_MStats, 3712);
			$panic(new $String("MStats vs MemStatsType size mismatch"));
		}
	};
	$pkg.$init = function() {
		($ptrType(NotSupportedError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		NotSupportedError.init([["Feature", "Feature", "", $String, ""]]);
		($ptrType(TypeAssertionError)).methods = [["Error", "Error", "", [], [$String], false, -1], ["RuntimeError", "RuntimeError", "", [], [], false, -1]];
		TypeAssertionError.init([["interfaceString", "interfaceString", "runtime", $String, ""], ["concreteString", "concreteString", "runtime", $String, ""], ["assertedString", "assertedString", "runtime", $String, ""], ["missingMethod", "missingMethod", "runtime", $String, ""]]);
		errorString.methods = [["Error", "Error", "", [], [$String], false, -1], ["RuntimeError", "RuntimeError", "", [], [], false, -1]];
		($ptrType(errorString)).methods = [["Error", "Error", "", [], [$String], false, -1], ["RuntimeError", "RuntimeError", "", [], [], false, -1]];
		MemStats.init([["Alloc", "Alloc", "", $Uint64, ""], ["TotalAlloc", "TotalAlloc", "", $Uint64, ""], ["Sys", "Sys", "", $Uint64, ""], ["Lookups", "Lookups", "", $Uint64, ""], ["Mallocs", "Mallocs", "", $Uint64, ""], ["Frees", "Frees", "", $Uint64, ""], ["HeapAlloc", "HeapAlloc", "", $Uint64, ""], ["HeapSys", "HeapSys", "", $Uint64, ""], ["HeapIdle", "HeapIdle", "", $Uint64, ""], ["HeapInuse", "HeapInuse", "", $Uint64, ""], ["HeapReleased", "HeapReleased", "", $Uint64, ""], ["HeapObjects", "HeapObjects", "", $Uint64, ""], ["StackInuse", "StackInuse", "", $Uint64, ""], ["StackSys", "StackSys", "", $Uint64, ""], ["MSpanInuse", "MSpanInuse", "", $Uint64, ""], ["MSpanSys", "MSpanSys", "", $Uint64, ""], ["MCacheInuse", "MCacheInuse", "", $Uint64, ""], ["MCacheSys", "MCacheSys", "", $Uint64, ""], ["BuckHashSys", "BuckHashSys", "", $Uint64, ""], ["GCSys", "GCSys", "", $Uint64, ""], ["OtherSys", "OtherSys", "", $Uint64, ""], ["NextGC", "NextGC", "", $Uint64, ""], ["LastGC", "LastGC", "", $Uint64, ""], ["PauseTotalNs", "PauseTotalNs", "", $Uint64, ""], ["PauseNs", "PauseNs", "", ($arrayType($Uint64, 256)), ""], ["NumGC", "NumGC", "", $Uint32, ""], ["EnableGC", "EnableGC", "", $Bool, ""], ["DebugGC", "DebugGC", "", $Bool, ""], ["BySize", "BySize", "", ($arrayType(($structType([["Size", "Size", "", $Uint32, ""], ["Mallocs", "Mallocs", "", $Uint64, ""], ["Frees", "Frees", "", $Uint64, ""]])), 61)), ""]]);
		sizeof_C_MStats = 3712;
		init();
		init$1();
	};
	return $pkg;
})();
$packages["errors"] = (function() {
	var $pkg = {}, errorString, New;
	errorString = $pkg.errorString = $newType(0, "Struct", "errors.errorString", "errorString", "errors", function(s_) {
		this.$val = this;
		this.s = s_ !== undefined ? s_ : "";
	});
	New = $pkg.New = function(text) {
		return new errorString.Ptr(text);
	};
	errorString.Ptr.prototype.Error = function() {
		var e;
		e = this;
		return e.s;
	};
	errorString.prototype.Error = function() { return this.$val.Error(); };
	$pkg.$init = function() {
		($ptrType(errorString)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		errorString.init([["s", "s", "errors", $String, ""]]);
	};
	return $pkg;
})();
$packages["sync/atomic"] = (function() {
	var $pkg = {}, CompareAndSwapInt32, AddInt32, LoadUint32, StoreInt32, StoreUint32;
	CompareAndSwapInt32 = $pkg.CompareAndSwapInt32 = function(addr, old, new$1) {
		if (addr.$get() === old) {
			addr.$set(new$1);
			return true;
		}
		return false;
	};
	AddInt32 = $pkg.AddInt32 = function(addr, delta) {
		var new$1;
		new$1 = addr.$get() + delta >> 0;
		addr.$set(new$1);
		return new$1;
	};
	LoadUint32 = $pkg.LoadUint32 = function(addr) {
		return addr.$get();
	};
	StoreInt32 = $pkg.StoreInt32 = function(addr, val) {
		addr.$set(val);
	};
	StoreUint32 = $pkg.StoreUint32 = function(addr, val) {
		addr.$set(val);
	};
	$pkg.$init = function() {
	};
	return $pkg;
})();
$packages["sync"] = (function() {
	var $pkg = {}, atomic = $packages["sync/atomic"], runtime = $packages["runtime"], Pool, Mutex, Locker, Once, poolLocal, syncSema, RWMutex, rlocker, allPools, runtime_registerPoolCleanup, runtime_Syncsemcheck, poolCleanup, init, indexLocal, runtime_Semacquire, runtime_Semrelease, init$1;
	Pool = $pkg.Pool = $newType(0, "Struct", "sync.Pool", "Pool", "sync", function(local_, localSize_, store_, New_) {
		this.$val = this;
		this.local = local_ !== undefined ? local_ : 0;
		this.localSize = localSize_ !== undefined ? localSize_ : 0;
		this.store = store_ !== undefined ? store_ : ($sliceType($emptyInterface)).nil;
		this.New = New_ !== undefined ? New_ : $throwNilPointerError;
	});
	Mutex = $pkg.Mutex = $newType(0, "Struct", "sync.Mutex", "Mutex", "sync", function(state_, sema_) {
		this.$val = this;
		this.state = state_ !== undefined ? state_ : 0;
		this.sema = sema_ !== undefined ? sema_ : 0;
	});
	Locker = $pkg.Locker = $newType(8, "Interface", "sync.Locker", "Locker", "sync", null);
	Once = $pkg.Once = $newType(0, "Struct", "sync.Once", "Once", "sync", function(m_, done_) {
		this.$val = this;
		this.m = m_ !== undefined ? m_ : new Mutex.Ptr();
		this.done = done_ !== undefined ? done_ : 0;
	});
	poolLocal = $pkg.poolLocal = $newType(0, "Struct", "sync.poolLocal", "poolLocal", "sync", function(private$0_, shared_, Mutex_, pad_) {
		this.$val = this;
		this.private$0 = private$0_ !== undefined ? private$0_ : null;
		this.shared = shared_ !== undefined ? shared_ : ($sliceType($emptyInterface)).nil;
		this.Mutex = Mutex_ !== undefined ? Mutex_ : new Mutex.Ptr();
		this.pad = pad_ !== undefined ? pad_ : ($arrayType($Uint8, 128)).zero();
	});
	syncSema = $pkg.syncSema = $newType(12, "Array", "sync.syncSema", "syncSema", "sync", null);
	RWMutex = $pkg.RWMutex = $newType(0, "Struct", "sync.RWMutex", "RWMutex", "sync", function(w_, writerSem_, readerSem_, readerCount_, readerWait_) {
		this.$val = this;
		this.w = w_ !== undefined ? w_ : new Mutex.Ptr();
		this.writerSem = writerSem_ !== undefined ? writerSem_ : 0;
		this.readerSem = readerSem_ !== undefined ? readerSem_ : 0;
		this.readerCount = readerCount_ !== undefined ? readerCount_ : 0;
		this.readerWait = readerWait_ !== undefined ? readerWait_ : 0;
	});
	rlocker = $pkg.rlocker = $newType(0, "Struct", "sync.rlocker", "rlocker", "sync", function(w_, writerSem_, readerSem_, readerCount_, readerWait_) {
		this.$val = this;
		this.w = w_ !== undefined ? w_ : new Mutex.Ptr();
		this.writerSem = writerSem_ !== undefined ? writerSem_ : 0;
		this.readerSem = readerSem_ !== undefined ? readerSem_ : 0;
		this.readerCount = readerCount_ !== undefined ? readerCount_ : 0;
		this.readerWait = readerWait_ !== undefined ? readerWait_ : 0;
	});
	Pool.Ptr.prototype.Get = function() {
		var p, x, x$1, x$2;
		p = this;
		if (p.store.$length === 0) {
			if (!(p.New === $throwNilPointerError)) {
				return p.New();
			}
			return null;
		}
		x$2 = (x = p.store, x$1 = p.store.$length - 1 >> 0, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1]));
		p.store = $subslice(p.store, 0, (p.store.$length - 1 >> 0));
		return x$2;
	};
	Pool.prototype.Get = function() { return this.$val.Get(); };
	Pool.Ptr.prototype.Put = function(x) {
		var p;
		p = this;
		if ($interfaceIsEqual(x, null)) {
			return;
		}
		p.store = $append(p.store, x);
	};
	Pool.prototype.Put = function(x) { return this.$val.Put(x); };
	runtime_registerPoolCleanup = function(cleanup) {
	};
	runtime_Syncsemcheck = function(size) {
	};
	Mutex.Ptr.prototype.Lock = function() {
		var m, awoke, old, new$1;
		m = this;
		if (atomic.CompareAndSwapInt32(new ($ptrType($Int32))(function() { return this.$target.state; }, function($v) { this.$target.state = $v; }, m), 0, 1)) {
			return;
		}
		awoke = false;
		while (true) {
			old = m.state;
			new$1 = old | 1;
			if (!(((old & 1) === 0))) {
				new$1 = old + 4 >> 0;
			}
			if (awoke) {
				new$1 = new$1 & ~(2);
			}
			if (atomic.CompareAndSwapInt32(new ($ptrType($Int32))(function() { return this.$target.state; }, function($v) { this.$target.state = $v; }, m), old, new$1)) {
				if ((old & 1) === 0) {
					break;
				}
				runtime_Semacquire(new ($ptrType($Uint32))(function() { return this.$target.sema; }, function($v) { this.$target.sema = $v; }, m));
				awoke = true;
			}
		}
	};
	Mutex.prototype.Lock = function() { return this.$val.Lock(); };
	Mutex.Ptr.prototype.Unlock = function() {
		var m, new$1, old;
		m = this;
		new$1 = atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.state; }, function($v) { this.$target.state = $v; }, m), -1);
		if ((((new$1 + 1 >> 0)) & 1) === 0) {
			$panic(new $String("sync: unlock of unlocked mutex"));
		}
		old = new$1;
		while (true) {
			if (((old >> 2 >> 0) === 0) || !(((old & 3) === 0))) {
				return;
			}
			new$1 = ((old - 4 >> 0)) | 2;
			if (atomic.CompareAndSwapInt32(new ($ptrType($Int32))(function() { return this.$target.state; }, function($v) { this.$target.state = $v; }, m), old, new$1)) {
				runtime_Semrelease(new ($ptrType($Uint32))(function() { return this.$target.sema; }, function($v) { this.$target.sema = $v; }, m));
				return;
			}
			old = m.state;
		}
	};
	Mutex.prototype.Unlock = function() { return this.$val.Unlock(); };
	Once.Ptr.prototype.Do = function(f) {
		var $deferred = [], $err = null, o, _recv;
		/* */ try { $deferFrames.push($deferred);
		o = this;
		if (atomic.LoadUint32(new ($ptrType($Uint32))(function() { return this.$target.done; }, function($v) { this.$target.done = $v; }, o)) === 1) {
			return;
		}
		o.m.Lock();
		$deferred.push([(_recv = o.m, function() { $stackDepthOffset--; try { return _recv.Unlock(); } finally { $stackDepthOffset++; } }), []]);
		if (o.done === 0) {
			f();
			atomic.StoreUint32(new ($ptrType($Uint32))(function() { return this.$target.done; }, function($v) { this.$target.done = $v; }, o), 1);
		}
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	Once.prototype.Do = function(f) { return this.$val.Do(f); };
	poolCleanup = function() {
		var _ref, _i, i, p, i$1, l, _ref$1, _i$1, j, x;
		_ref = allPools;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			p = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			(i < 0 || i >= allPools.$length) ? $throwRuntimeError("index out of range") : allPools.$array[allPools.$offset + i] = ($ptrType(Pool)).nil;
			i$1 = 0;
			while (i$1 < (p.localSize >> 0)) {
				l = indexLocal(p.local, i$1);
				l.private$0 = null;
				_ref$1 = l.shared;
				_i$1 = 0;
				while (_i$1 < _ref$1.$length) {
					j = _i$1;
					(x = l.shared, (j < 0 || j >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + j] = null);
					_i$1++;
				}
				l.shared = ($sliceType($emptyInterface)).nil;
				i$1 = i$1 + (1) >> 0;
			}
			_i++;
		}
		allPools = new ($sliceType(($ptrType(Pool))))([]);
	};
	init = function() {
		runtime_registerPoolCleanup(poolCleanup);
	};
	indexLocal = function(l, i) {
		var x;
		return (x = l, ((i < 0 || i >= x.length) ? $throwRuntimeError("index out of range") : x[i]));
	};
	runtime_Semacquire = function() {
		$panic("Native function not implemented: sync.runtime_Semacquire");
	};
	runtime_Semrelease = function() {
		$panic("Native function not implemented: sync.runtime_Semrelease");
	};
	init$1 = function() {
		var s;
		s = syncSema.zero(); $copy(s, syncSema.zero(), syncSema);
		runtime_Syncsemcheck(12);
	};
	RWMutex.Ptr.prototype.RLock = function() {
		var rw;
		rw = this;
		if (atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.readerCount; }, function($v) { this.$target.readerCount = $v; }, rw), 1) < 0) {
			runtime_Semacquire(new ($ptrType($Uint32))(function() { return this.$target.readerSem; }, function($v) { this.$target.readerSem = $v; }, rw));
		}
	};
	RWMutex.prototype.RLock = function() { return this.$val.RLock(); };
	RWMutex.Ptr.prototype.RUnlock = function() {
		var rw;
		rw = this;
		if (atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.readerCount; }, function($v) { this.$target.readerCount = $v; }, rw), -1) < 0) {
			if (atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.readerWait; }, function($v) { this.$target.readerWait = $v; }, rw), -1) === 0) {
				runtime_Semrelease(new ($ptrType($Uint32))(function() { return this.$target.writerSem; }, function($v) { this.$target.writerSem = $v; }, rw));
			}
		}
	};
	RWMutex.prototype.RUnlock = function() { return this.$val.RUnlock(); };
	RWMutex.Ptr.prototype.Lock = function() {
		var rw, r;
		rw = this;
		rw.w.Lock();
		r = atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.readerCount; }, function($v) { this.$target.readerCount = $v; }, rw), -1073741824) + 1073741824 >> 0;
		if (!((r === 0)) && !((atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.readerWait; }, function($v) { this.$target.readerWait = $v; }, rw), r) === 0))) {
			runtime_Semacquire(new ($ptrType($Uint32))(function() { return this.$target.writerSem; }, function($v) { this.$target.writerSem = $v; }, rw));
		}
	};
	RWMutex.prototype.Lock = function() { return this.$val.Lock(); };
	RWMutex.Ptr.prototype.Unlock = function() {
		var rw, r, i;
		rw = this;
		r = atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.readerCount; }, function($v) { this.$target.readerCount = $v; }, rw), 1073741824);
		i = 0;
		while (i < (r >> 0)) {
			runtime_Semrelease(new ($ptrType($Uint32))(function() { return this.$target.readerSem; }, function($v) { this.$target.readerSem = $v; }, rw));
			i = i + (1) >> 0;
		}
		rw.w.Unlock();
	};
	RWMutex.prototype.Unlock = function() { return this.$val.Unlock(); };
	RWMutex.Ptr.prototype.RLocker = function() {
		var rw;
		rw = this;
		return $clone(rw, rlocker);
	};
	RWMutex.prototype.RLocker = function() { return this.$val.RLocker(); };
	rlocker.Ptr.prototype.Lock = function() {
		var r;
		r = this;
		$clone(r, RWMutex).RLock();
	};
	rlocker.prototype.Lock = function() { return this.$val.Lock(); };
	rlocker.Ptr.prototype.Unlock = function() {
		var r;
		r = this;
		$clone(r, RWMutex).RUnlock();
	};
	rlocker.prototype.Unlock = function() { return this.$val.Unlock(); };
	$pkg.$init = function() {
		($ptrType(Pool)).methods = [["Get", "Get", "", [], [$emptyInterface], false, -1], ["Put", "Put", "", [$emptyInterface], [], false, -1], ["getSlow", "getSlow", "sync", [], [$emptyInterface], false, -1], ["pin", "pin", "sync", [], [($ptrType(poolLocal))], false, -1], ["pinSlow", "pinSlow", "sync", [], [($ptrType(poolLocal))], false, -1]];
		Pool.init([["local", "local", "sync", $UnsafePointer, ""], ["localSize", "localSize", "sync", $Uintptr, ""], ["store", "store", "sync", ($sliceType($emptyInterface)), ""], ["New", "New", "", ($funcType([], [$emptyInterface], false)), ""]]);
		($ptrType(Mutex)).methods = [["Lock", "Lock", "", [], [], false, -1], ["Unlock", "Unlock", "", [], [], false, -1]];
		Mutex.init([["state", "state", "sync", $Int32, ""], ["sema", "sema", "sync", $Uint32, ""]]);
		Locker.init([["Lock", "Lock", "", [], [], false], ["Unlock", "Unlock", "", [], [], false]]);
		($ptrType(Once)).methods = [["Do", "Do", "", [($funcType([], [], false))], [], false, -1]];
		Once.init([["m", "m", "sync", Mutex, ""], ["done", "done", "sync", $Uint32, ""]]);
		($ptrType(poolLocal)).methods = [["Lock", "Lock", "", [], [], false, 2], ["Unlock", "Unlock", "", [], [], false, 2]];
		poolLocal.init([["private$0", "private", "sync", $emptyInterface, ""], ["shared", "shared", "sync", ($sliceType($emptyInterface)), ""], ["Mutex", "", "", Mutex, ""], ["pad", "pad", "sync", ($arrayType($Uint8, 128)), ""]]);
		syncSema.init($Uintptr, 3);
		($ptrType(RWMutex)).methods = [["Lock", "Lock", "", [], [], false, -1], ["RLock", "RLock", "", [], [], false, -1], ["RLocker", "RLocker", "", [], [Locker], false, -1], ["RUnlock", "RUnlock", "", [], [], false, -1], ["Unlock", "Unlock", "", [], [], false, -1]];
		RWMutex.init([["w", "w", "sync", Mutex, ""], ["writerSem", "writerSem", "sync", $Uint32, ""], ["readerSem", "readerSem", "sync", $Uint32, ""], ["readerCount", "readerCount", "sync", $Int32, ""], ["readerWait", "readerWait", "sync", $Int32, ""]]);
		($ptrType(rlocker)).methods = [["Lock", "Lock", "", [], [], false, -1], ["Unlock", "Unlock", "", [], [], false, -1]];
		rlocker.init([["w", "w", "sync", Mutex, ""], ["writerSem", "writerSem", "sync", $Uint32, ""], ["readerSem", "readerSem", "sync", $Uint32, ""], ["readerCount", "readerCount", "sync", $Int32, ""], ["readerWait", "readerWait", "sync", $Int32, ""]]);
		allPools = ($sliceType(($ptrType(Pool)))).nil;
		init();
		init$1();
	};
	return $pkg;
})();
$packages["io"] = (function() {
	var $pkg = {}, runtime = $packages["runtime"], errors = $packages["errors"], sync = $packages["sync"], Reader, Writer, RuneReader, errWhence, errOffset;
	Reader = $pkg.Reader = $newType(8, "Interface", "io.Reader", "Reader", "io", null);
	Writer = $pkg.Writer = $newType(8, "Interface", "io.Writer", "Writer", "io", null);
	RuneReader = $pkg.RuneReader = $newType(8, "Interface", "io.RuneReader", "RuneReader", "io", null);
	$pkg.$init = function() {
		Reader.init([["Read", "Read", "", [($sliceType($Uint8))], [$Int, $error], false]]);
		Writer.init([["Write", "Write", "", [($sliceType($Uint8))], [$Int, $error], false]]);
		RuneReader.init([["ReadRune", "ReadRune", "", [], [$Int32, $Int, $error], false]]);
		$pkg.ErrShortWrite = errors.New("short write");
		$pkg.ErrShortBuffer = errors.New("short buffer");
		$pkg.EOF = errors.New("EOF");
		$pkg.ErrUnexpectedEOF = errors.New("unexpected EOF");
		$pkg.ErrNoProgress = errors.New("multiple Read calls return no data or error");
		errWhence = errors.New("Seek: invalid whence");
		errOffset = errors.New("Seek: invalid offset");
		$pkg.ErrClosedPipe = errors.New("io: read/write on closed pipe");
	};
	return $pkg;
})();
$packages["unicode"] = (function() {
	var $pkg = {}, RangeTable, Range16, Range32, _White_Space, IsSpace, is16, is32, isExcludingLatin;
	RangeTable = $pkg.RangeTable = $newType(0, "Struct", "unicode.RangeTable", "RangeTable", "unicode", function(R16_, R32_, LatinOffset_) {
		this.$val = this;
		this.R16 = R16_ !== undefined ? R16_ : ($sliceType(Range16)).nil;
		this.R32 = R32_ !== undefined ? R32_ : ($sliceType(Range32)).nil;
		this.LatinOffset = LatinOffset_ !== undefined ? LatinOffset_ : 0;
	});
	Range16 = $pkg.Range16 = $newType(0, "Struct", "unicode.Range16", "Range16", "unicode", function(Lo_, Hi_, Stride_) {
		this.$val = this;
		this.Lo = Lo_ !== undefined ? Lo_ : 0;
		this.Hi = Hi_ !== undefined ? Hi_ : 0;
		this.Stride = Stride_ !== undefined ? Stride_ : 0;
	});
	Range32 = $pkg.Range32 = $newType(0, "Struct", "unicode.Range32", "Range32", "unicode", function(Lo_, Hi_, Stride_) {
		this.$val = this;
		this.Lo = Lo_ !== undefined ? Lo_ : 0;
		this.Hi = Hi_ !== undefined ? Hi_ : 0;
		this.Stride = Stride_ !== undefined ? Stride_ : 0;
	});
	IsSpace = $pkg.IsSpace = function(r) {
		var _ref;
		if ((r >>> 0) <= 255) {
			_ref = r;
			if (_ref === 9 || _ref === 10 || _ref === 11 || _ref === 12 || _ref === 13 || _ref === 32 || _ref === 133 || _ref === 160) {
				return true;
			}
			return false;
		}
		return isExcludingLatin($pkg.White_Space, r);
	};
	is16 = function(ranges, r) {
		var _ref, _i, i, range_, _r, lo, hi, _q, m, range_$1, _r$1;
		if (ranges.$length <= 18 || r <= 255) {
			_ref = ranges;
			_i = 0;
			while (_i < _ref.$length) {
				i = _i;
				range_ = ((i < 0 || i >= ranges.$length) ? $throwRuntimeError("index out of range") : ranges.$array[ranges.$offset + i]);
				if (r < range_.Lo) {
					return false;
				}
				if (r <= range_.Hi) {
					return (_r = ((r - range_.Lo << 16 >>> 16)) % range_.Stride, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) === 0;
				}
				_i++;
			}
			return false;
		}
		lo = 0;
		hi = ranges.$length;
		while (lo < hi) {
			m = lo + (_q = ((hi - lo >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			range_$1 = ((m < 0 || m >= ranges.$length) ? $throwRuntimeError("index out of range") : ranges.$array[ranges.$offset + m]);
			if (range_$1.Lo <= r && r <= range_$1.Hi) {
				return (_r$1 = ((r - range_$1.Lo << 16 >>> 16)) % range_$1.Stride, _r$1 === _r$1 ? _r$1 : $throwRuntimeError("integer divide by zero")) === 0;
			}
			if (r < range_$1.Lo) {
				hi = m;
			} else {
				lo = m + 1 >> 0;
			}
		}
		return false;
	};
	is32 = function(ranges, r) {
		var _ref, _i, i, range_, _r, lo, hi, _q, m, range_$1, _r$1;
		if (ranges.$length <= 18) {
			_ref = ranges;
			_i = 0;
			while (_i < _ref.$length) {
				i = _i;
				range_ = ((i < 0 || i >= ranges.$length) ? $throwRuntimeError("index out of range") : ranges.$array[ranges.$offset + i]);
				if (r < range_.Lo) {
					return false;
				}
				if (r <= range_.Hi) {
					return (_r = ((r - range_.Lo >>> 0)) % range_.Stride, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) === 0;
				}
				_i++;
			}
			return false;
		}
		lo = 0;
		hi = ranges.$length;
		while (lo < hi) {
			m = lo + (_q = ((hi - lo >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			range_$1 = new Range32.Ptr(); $copy(range_$1, ((m < 0 || m >= ranges.$length) ? $throwRuntimeError("index out of range") : ranges.$array[ranges.$offset + m]), Range32);
			if (range_$1.Lo <= r && r <= range_$1.Hi) {
				return (_r$1 = ((r - range_$1.Lo >>> 0)) % range_$1.Stride, _r$1 === _r$1 ? _r$1 : $throwRuntimeError("integer divide by zero")) === 0;
			}
			if (r < range_$1.Lo) {
				hi = m;
			} else {
				lo = m + 1 >> 0;
			}
		}
		return false;
	};
	isExcludingLatin = function(rangeTab, r) {
		var r16, off, x, r32;
		r16 = rangeTab.R16;
		off = rangeTab.LatinOffset;
		if (r16.$length > off && r <= ((x = r16.$length - 1 >> 0, ((x < 0 || x >= r16.$length) ? $throwRuntimeError("index out of range") : r16.$array[r16.$offset + x])).Hi >> 0)) {
			return is16($subslice(r16, off), (r << 16 >>> 16));
		}
		r32 = rangeTab.R32;
		if (r32.$length > 0 && r >= (((0 < 0 || 0 >= r32.$length) ? $throwRuntimeError("index out of range") : r32.$array[r32.$offset + 0]).Lo >> 0)) {
			return is32(r32, (r >>> 0));
		}
		return false;
	};
	$pkg.$init = function() {
		RangeTable.init([["R16", "R16", "", ($sliceType(Range16)), ""], ["R32", "R32", "", ($sliceType(Range32)), ""], ["LatinOffset", "LatinOffset", "", $Int, ""]]);
		Range16.init([["Lo", "Lo", "", $Uint16, ""], ["Hi", "Hi", "", $Uint16, ""], ["Stride", "Stride", "", $Uint16, ""]]);
		Range32.init([["Lo", "Lo", "", $Uint32, ""], ["Hi", "Hi", "", $Uint32, ""], ["Stride", "Stride", "", $Uint32, ""]]);
		_White_Space = new RangeTable.Ptr(new ($sliceType(Range16))([new Range16.Ptr(9, 13, 1), new Range16.Ptr(32, 32, 1), new Range16.Ptr(133, 133, 1), new Range16.Ptr(160, 160, 1), new Range16.Ptr(5760, 5760, 1), new Range16.Ptr(8192, 8202, 1), new Range16.Ptr(8232, 8233, 1), new Range16.Ptr(8239, 8239, 1), new Range16.Ptr(8287, 8287, 1), new Range16.Ptr(12288, 12288, 1)]), ($sliceType(Range32)).nil, 4);
		$pkg.White_Space = _White_Space;
	};
	return $pkg;
})();
$packages["unicode/utf8"] = (function() {
	var $pkg = {}, decodeRuneInternal, decodeRuneInStringInternal, DecodeRune, DecodeRuneInString, DecodeLastRune, DecodeLastRuneInString, RuneLen, EncodeRune, RuneCountInString, RuneStart;
	decodeRuneInternal = function(p) {
		var r = 0, size = 0, short$1 = false, n, _tmp, _tmp$1, _tmp$2, c0, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10, _tmp$11, c1, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, _tmp$23, c2, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, _tmp$38, c3, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46, _tmp$47, _tmp$48, _tmp$49, _tmp$50;
		n = p.$length;
		if (n < 1) {
			_tmp = 65533; _tmp$1 = 0; _tmp$2 = true; r = _tmp; size = _tmp$1; short$1 = _tmp$2;
			return [r, size, short$1];
		}
		c0 = ((0 < 0 || 0 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 0]);
		if (c0 < 128) {
			_tmp$3 = (c0 >> 0); _tmp$4 = 1; _tmp$5 = false; r = _tmp$3; size = _tmp$4; short$1 = _tmp$5;
			return [r, size, short$1];
		}
		if (c0 < 192) {
			_tmp$6 = 65533; _tmp$7 = 1; _tmp$8 = false; r = _tmp$6; size = _tmp$7; short$1 = _tmp$8;
			return [r, size, short$1];
		}
		if (n < 2) {
			_tmp$9 = 65533; _tmp$10 = 1; _tmp$11 = true; r = _tmp$9; size = _tmp$10; short$1 = _tmp$11;
			return [r, size, short$1];
		}
		c1 = ((1 < 0 || 1 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 1]);
		if (c1 < 128 || 192 <= c1) {
			_tmp$12 = 65533; _tmp$13 = 1; _tmp$14 = false; r = _tmp$12; size = _tmp$13; short$1 = _tmp$14;
			return [r, size, short$1];
		}
		if (c0 < 224) {
			r = ((((c0 & 31) >>> 0) >> 0) << 6 >> 0) | (((c1 & 63) >>> 0) >> 0);
			if (r <= 127) {
				_tmp$15 = 65533; _tmp$16 = 1; _tmp$17 = false; r = _tmp$15; size = _tmp$16; short$1 = _tmp$17;
				return [r, size, short$1];
			}
			_tmp$18 = r; _tmp$19 = 2; _tmp$20 = false; r = _tmp$18; size = _tmp$19; short$1 = _tmp$20;
			return [r, size, short$1];
		}
		if (n < 3) {
			_tmp$21 = 65533; _tmp$22 = 1; _tmp$23 = true; r = _tmp$21; size = _tmp$22; short$1 = _tmp$23;
			return [r, size, short$1];
		}
		c2 = ((2 < 0 || 2 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 2]);
		if (c2 < 128 || 192 <= c2) {
			_tmp$24 = 65533; _tmp$25 = 1; _tmp$26 = false; r = _tmp$24; size = _tmp$25; short$1 = _tmp$26;
			return [r, size, short$1];
		}
		if (c0 < 240) {
			r = (((((c0 & 15) >>> 0) >> 0) << 12 >> 0) | ((((c1 & 63) >>> 0) >> 0) << 6 >> 0)) | (((c2 & 63) >>> 0) >> 0);
			if (r <= 2047) {
				_tmp$27 = 65533; _tmp$28 = 1; _tmp$29 = false; r = _tmp$27; size = _tmp$28; short$1 = _tmp$29;
				return [r, size, short$1];
			}
			if (55296 <= r && r <= 57343) {
				_tmp$30 = 65533; _tmp$31 = 1; _tmp$32 = false; r = _tmp$30; size = _tmp$31; short$1 = _tmp$32;
				return [r, size, short$1];
			}
			_tmp$33 = r; _tmp$34 = 3; _tmp$35 = false; r = _tmp$33; size = _tmp$34; short$1 = _tmp$35;
			return [r, size, short$1];
		}
		if (n < 4) {
			_tmp$36 = 65533; _tmp$37 = 1; _tmp$38 = true; r = _tmp$36; size = _tmp$37; short$1 = _tmp$38;
			return [r, size, short$1];
		}
		c3 = ((3 < 0 || 3 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 3]);
		if (c3 < 128 || 192 <= c3) {
			_tmp$39 = 65533; _tmp$40 = 1; _tmp$41 = false; r = _tmp$39; size = _tmp$40; short$1 = _tmp$41;
			return [r, size, short$1];
		}
		if (c0 < 248) {
			r = ((((((c0 & 7) >>> 0) >> 0) << 18 >> 0) | ((((c1 & 63) >>> 0) >> 0) << 12 >> 0)) | ((((c2 & 63) >>> 0) >> 0) << 6 >> 0)) | (((c3 & 63) >>> 0) >> 0);
			if (r <= 65535 || 1114111 < r) {
				_tmp$42 = 65533; _tmp$43 = 1; _tmp$44 = false; r = _tmp$42; size = _tmp$43; short$1 = _tmp$44;
				return [r, size, short$1];
			}
			_tmp$45 = r; _tmp$46 = 4; _tmp$47 = false; r = _tmp$45; size = _tmp$46; short$1 = _tmp$47;
			return [r, size, short$1];
		}
		_tmp$48 = 65533; _tmp$49 = 1; _tmp$50 = false; r = _tmp$48; size = _tmp$49; short$1 = _tmp$50;
		return [r, size, short$1];
	};
	decodeRuneInStringInternal = function(s) {
		var r = 0, size = 0, short$1 = false, n, _tmp, _tmp$1, _tmp$2, c0, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10, _tmp$11, c1, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, _tmp$23, c2, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, _tmp$38, c3, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46, _tmp$47, _tmp$48, _tmp$49, _tmp$50;
		n = s.length;
		if (n < 1) {
			_tmp = 65533; _tmp$1 = 0; _tmp$2 = true; r = _tmp; size = _tmp$1; short$1 = _tmp$2;
			return [r, size, short$1];
		}
		c0 = s.charCodeAt(0);
		if (c0 < 128) {
			_tmp$3 = (c0 >> 0); _tmp$4 = 1; _tmp$5 = false; r = _tmp$3; size = _tmp$4; short$1 = _tmp$5;
			return [r, size, short$1];
		}
		if (c0 < 192) {
			_tmp$6 = 65533; _tmp$7 = 1; _tmp$8 = false; r = _tmp$6; size = _tmp$7; short$1 = _tmp$8;
			return [r, size, short$1];
		}
		if (n < 2) {
			_tmp$9 = 65533; _tmp$10 = 1; _tmp$11 = true; r = _tmp$9; size = _tmp$10; short$1 = _tmp$11;
			return [r, size, short$1];
		}
		c1 = s.charCodeAt(1);
		if (c1 < 128 || 192 <= c1) {
			_tmp$12 = 65533; _tmp$13 = 1; _tmp$14 = false; r = _tmp$12; size = _tmp$13; short$1 = _tmp$14;
			return [r, size, short$1];
		}
		if (c0 < 224) {
			r = ((((c0 & 31) >>> 0) >> 0) << 6 >> 0) | (((c1 & 63) >>> 0) >> 0);
			if (r <= 127) {
				_tmp$15 = 65533; _tmp$16 = 1; _tmp$17 = false; r = _tmp$15; size = _tmp$16; short$1 = _tmp$17;
				return [r, size, short$1];
			}
			_tmp$18 = r; _tmp$19 = 2; _tmp$20 = false; r = _tmp$18; size = _tmp$19; short$1 = _tmp$20;
			return [r, size, short$1];
		}
		if (n < 3) {
			_tmp$21 = 65533; _tmp$22 = 1; _tmp$23 = true; r = _tmp$21; size = _tmp$22; short$1 = _tmp$23;
			return [r, size, short$1];
		}
		c2 = s.charCodeAt(2);
		if (c2 < 128 || 192 <= c2) {
			_tmp$24 = 65533; _tmp$25 = 1; _tmp$26 = false; r = _tmp$24; size = _tmp$25; short$1 = _tmp$26;
			return [r, size, short$1];
		}
		if (c0 < 240) {
			r = (((((c0 & 15) >>> 0) >> 0) << 12 >> 0) | ((((c1 & 63) >>> 0) >> 0) << 6 >> 0)) | (((c2 & 63) >>> 0) >> 0);
			if (r <= 2047) {
				_tmp$27 = 65533; _tmp$28 = 1; _tmp$29 = false; r = _tmp$27; size = _tmp$28; short$1 = _tmp$29;
				return [r, size, short$1];
			}
			if (55296 <= r && r <= 57343) {
				_tmp$30 = 65533; _tmp$31 = 1; _tmp$32 = false; r = _tmp$30; size = _tmp$31; short$1 = _tmp$32;
				return [r, size, short$1];
			}
			_tmp$33 = r; _tmp$34 = 3; _tmp$35 = false; r = _tmp$33; size = _tmp$34; short$1 = _tmp$35;
			return [r, size, short$1];
		}
		if (n < 4) {
			_tmp$36 = 65533; _tmp$37 = 1; _tmp$38 = true; r = _tmp$36; size = _tmp$37; short$1 = _tmp$38;
			return [r, size, short$1];
		}
		c3 = s.charCodeAt(3);
		if (c3 < 128 || 192 <= c3) {
			_tmp$39 = 65533; _tmp$40 = 1; _tmp$41 = false; r = _tmp$39; size = _tmp$40; short$1 = _tmp$41;
			return [r, size, short$1];
		}
		if (c0 < 248) {
			r = ((((((c0 & 7) >>> 0) >> 0) << 18 >> 0) | ((((c1 & 63) >>> 0) >> 0) << 12 >> 0)) | ((((c2 & 63) >>> 0) >> 0) << 6 >> 0)) | (((c3 & 63) >>> 0) >> 0);
			if (r <= 65535 || 1114111 < r) {
				_tmp$42 = 65533; _tmp$43 = 1; _tmp$44 = false; r = _tmp$42; size = _tmp$43; short$1 = _tmp$44;
				return [r, size, short$1];
			}
			_tmp$45 = r; _tmp$46 = 4; _tmp$47 = false; r = _tmp$45; size = _tmp$46; short$1 = _tmp$47;
			return [r, size, short$1];
		}
		_tmp$48 = 65533; _tmp$49 = 1; _tmp$50 = false; r = _tmp$48; size = _tmp$49; short$1 = _tmp$50;
		return [r, size, short$1];
	};
	DecodeRune = $pkg.DecodeRune = function(p) {
		var r = 0, size = 0, _tuple;
		_tuple = decodeRuneInternal(p); r = _tuple[0]; size = _tuple[1];
		return [r, size];
	};
	DecodeRuneInString = $pkg.DecodeRuneInString = function(s) {
		var r = 0, size = 0, _tuple;
		_tuple = decodeRuneInStringInternal(s); r = _tuple[0]; size = _tuple[1];
		return [r, size];
	};
	DecodeLastRune = $pkg.DecodeLastRune = function(p) {
		var r = 0, size = 0, end, _tmp, _tmp$1, start, _tmp$2, _tmp$3, lim, _tuple, _tmp$4, _tmp$5, _tmp$6, _tmp$7;
		end = p.$length;
		if (end === 0) {
			_tmp = 65533; _tmp$1 = 0; r = _tmp; size = _tmp$1;
			return [r, size];
		}
		start = end - 1 >> 0;
		r = (((start < 0 || start >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + start]) >> 0);
		if (r < 128) {
			_tmp$2 = r; _tmp$3 = 1; r = _tmp$2; size = _tmp$3;
			return [r, size];
		}
		lim = end - 4 >> 0;
		if (lim < 0) {
			lim = 0;
		}
		start = start - (1) >> 0;
		while (start >= lim) {
			if (RuneStart(((start < 0 || start >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + start]))) {
				break;
			}
			start = start - (1) >> 0;
		}
		if (start < 0) {
			start = 0;
		}
		_tuple = DecodeRune($subslice(p, start, end)); r = _tuple[0]; size = _tuple[1];
		if (!(((start + size >> 0) === end))) {
			_tmp$4 = 65533; _tmp$5 = 1; r = _tmp$4; size = _tmp$5;
			return [r, size];
		}
		_tmp$6 = r; _tmp$7 = size; r = _tmp$6; size = _tmp$7;
		return [r, size];
	};
	DecodeLastRuneInString = $pkg.DecodeLastRuneInString = function(s) {
		var r = 0, size = 0, end, _tmp, _tmp$1, start, _tmp$2, _tmp$3, lim, _tuple, _tmp$4, _tmp$5, _tmp$6, _tmp$7;
		end = s.length;
		if (end === 0) {
			_tmp = 65533; _tmp$1 = 0; r = _tmp; size = _tmp$1;
			return [r, size];
		}
		start = end - 1 >> 0;
		r = (s.charCodeAt(start) >> 0);
		if (r < 128) {
			_tmp$2 = r; _tmp$3 = 1; r = _tmp$2; size = _tmp$3;
			return [r, size];
		}
		lim = end - 4 >> 0;
		if (lim < 0) {
			lim = 0;
		}
		start = start - (1) >> 0;
		while (start >= lim) {
			if (RuneStart(s.charCodeAt(start))) {
				break;
			}
			start = start - (1) >> 0;
		}
		if (start < 0) {
			start = 0;
		}
		_tuple = DecodeRuneInString(s.substring(start, end)); r = _tuple[0]; size = _tuple[1];
		if (!(((start + size >> 0) === end))) {
			_tmp$4 = 65533; _tmp$5 = 1; r = _tmp$4; size = _tmp$5;
			return [r, size];
		}
		_tmp$6 = r; _tmp$7 = size; r = _tmp$6; size = _tmp$7;
		return [r, size];
	};
	RuneLen = $pkg.RuneLen = function(r) {
		if (r < 0) {
			return -1;
		} else if (r <= 127) {
			return 1;
		} else if (r <= 2047) {
			return 2;
		} else if (55296 <= r && r <= 57343) {
			return -1;
		} else if (r <= 65535) {
			return 3;
		} else if (r <= 1114111) {
			return 4;
		}
		return -1;
	};
	EncodeRune = $pkg.EncodeRune = function(p, r) {
		var i;
		i = (r >>> 0);
		if (i <= 127) {
			(0 < 0 || 0 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 0] = (r << 24 >>> 24);
			return 1;
		} else if (i <= 2047) {
			(0 < 0 || 0 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 0] = (192 | ((r >> 6 >> 0) << 24 >>> 24)) >>> 0;
			(1 < 0 || 1 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 1] = (128 | (((r << 24 >>> 24) & 63) >>> 0)) >>> 0;
			return 2;
		} else if (i > 1114111 || 55296 <= i && i <= 57343) {
			r = 65533;
			(0 < 0 || 0 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 0] = (224 | ((r >> 12 >> 0) << 24 >>> 24)) >>> 0;
			(1 < 0 || 1 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 1] = (128 | ((((r >> 6 >> 0) << 24 >>> 24) & 63) >>> 0)) >>> 0;
			(2 < 0 || 2 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 2] = (128 | (((r << 24 >>> 24) & 63) >>> 0)) >>> 0;
			return 3;
		} else if (i <= 65535) {
			(0 < 0 || 0 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 0] = (224 | ((r >> 12 >> 0) << 24 >>> 24)) >>> 0;
			(1 < 0 || 1 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 1] = (128 | ((((r >> 6 >> 0) << 24 >>> 24) & 63) >>> 0)) >>> 0;
			(2 < 0 || 2 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 2] = (128 | (((r << 24 >>> 24) & 63) >>> 0)) >>> 0;
			return 3;
		} else {
			(0 < 0 || 0 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 0] = (240 | ((r >> 18 >> 0) << 24 >>> 24)) >>> 0;
			(1 < 0 || 1 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 1] = (128 | ((((r >> 12 >> 0) << 24 >>> 24) & 63) >>> 0)) >>> 0;
			(2 < 0 || 2 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 2] = (128 | ((((r >> 6 >> 0) << 24 >>> 24) & 63) >>> 0)) >>> 0;
			(3 < 0 || 3 >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + 3] = (128 | (((r << 24 >>> 24) & 63) >>> 0)) >>> 0;
			return 4;
		}
	};
	RuneCountInString = $pkg.RuneCountInString = function(s) {
		var n = 0, _ref, _i, _rune;
		_ref = s;
		_i = 0;
		while (_i < _ref.length) {
			_rune = $decodeRune(_ref, _i);
			n = n + (1) >> 0;
			_i += _rune[1];
		}
		return n;
	};
	RuneStart = $pkg.RuneStart = function(b) {
		return !((((b & 192) >>> 0) === 128));
	};
	$pkg.$init = function() {
	};
	return $pkg;
})();
$packages["bytes"] = (function() {
	var $pkg = {}, errors = $packages["errors"], io = $packages["io"], utf8 = $packages["unicode/utf8"], unicode = $packages["unicode"], Buffer, readOp, IndexByte, makeSlice, NewBuffer;
	Buffer = $pkg.Buffer = $newType(0, "Struct", "bytes.Buffer", "Buffer", "bytes", function(buf_, off_, runeBytes_, bootstrap_, lastRead_) {
		this.$val = this;
		this.buf = buf_ !== undefined ? buf_ : ($sliceType($Uint8)).nil;
		this.off = off_ !== undefined ? off_ : 0;
		this.runeBytes = runeBytes_ !== undefined ? runeBytes_ : ($arrayType($Uint8, 4)).zero();
		this.bootstrap = bootstrap_ !== undefined ? bootstrap_ : ($arrayType($Uint8, 64)).zero();
		this.lastRead = lastRead_ !== undefined ? lastRead_ : 0;
	});
	readOp = $pkg.readOp = $newType(4, "Int", "bytes.readOp", "readOp", "bytes", null);
	IndexByte = $pkg.IndexByte = function(s, c) {
		var _ref, _i, i, b;
		_ref = s;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			b = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (b === c) {
				return i;
			}
			_i++;
		}
		return -1;
	};
	Buffer.Ptr.prototype.Bytes = function() {
		var b;
		b = this;
		return $subslice(b.buf, b.off);
	};
	Buffer.prototype.Bytes = function() { return this.$val.Bytes(); };
	Buffer.Ptr.prototype.String = function() {
		var b;
		b = this;
		if (b === ($ptrType(Buffer)).nil) {
			return "<nil>";
		}
		return $bytesToString($subslice(b.buf, b.off));
	};
	Buffer.prototype.String = function() { return this.$val.String(); };
	Buffer.Ptr.prototype.Len = function() {
		var b;
		b = this;
		return b.buf.$length - b.off >> 0;
	};
	Buffer.prototype.Len = function() { return this.$val.Len(); };
	Buffer.Ptr.prototype.Truncate = function(n) {
		var b;
		b = this;
		b.lastRead = 0;
		if (n < 0 || n > b.Len()) {
			$panic(new $String("bytes.Buffer: truncation out of range"));
		} else if (n === 0) {
			b.off = 0;
		}
		b.buf = $subslice(b.buf, 0, (b.off + n >> 0));
	};
	Buffer.prototype.Truncate = function(n) { return this.$val.Truncate(n); };
	Buffer.Ptr.prototype.Reset = function() {
		var b;
		b = this;
		b.Truncate(0);
	};
	Buffer.prototype.Reset = function() { return this.$val.Reset(); };
	Buffer.Ptr.prototype.grow = function(n) {
		var b, m, buf, _q, x;
		b = this;
		m = b.Len();
		if ((m === 0) && !((b.off === 0))) {
			b.Truncate(0);
		}
		if ((b.buf.$length + n >> 0) > b.buf.$capacity) {
			buf = ($sliceType($Uint8)).nil;
			if (b.buf === ($sliceType($Uint8)).nil && n <= 64) {
				buf = $subslice(new ($sliceType($Uint8))(b.bootstrap), 0);
			} else if ((m + n >> 0) <= (_q = b.buf.$capacity / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"))) {
				$copySlice(b.buf, $subslice(b.buf, b.off));
				buf = $subslice(b.buf, 0, m);
			} else {
				buf = makeSlice((x = b.buf.$capacity, (((2 >>> 16 << 16) * x >> 0) + (2 << 16 >>> 16) * x) >> 0) + n >> 0);
				$copySlice(buf, $subslice(b.buf, b.off));
			}
			b.buf = buf;
			b.off = 0;
		}
		b.buf = $subslice(b.buf, 0, ((b.off + m >> 0) + n >> 0));
		return b.off + m >> 0;
	};
	Buffer.prototype.grow = function(n) { return this.$val.grow(n); };
	Buffer.Ptr.prototype.Grow = function(n) {
		var b, m;
		b = this;
		if (n < 0) {
			$panic(new $String("bytes.Buffer.Grow: negative count"));
		}
		m = b.grow(n);
		b.buf = $subslice(b.buf, 0, m);
	};
	Buffer.prototype.Grow = function(n) { return this.$val.Grow(n); };
	Buffer.Ptr.prototype.Write = function(p) {
		var n = 0, err = null, b, m, _tmp, _tmp$1;
		b = this;
		b.lastRead = 0;
		m = b.grow(p.$length);
		_tmp = $copySlice($subslice(b.buf, m), p); _tmp$1 = null; n = _tmp; err = _tmp$1;
		return [n, err];
	};
	Buffer.prototype.Write = function(p) { return this.$val.Write(p); };
	Buffer.Ptr.prototype.WriteString = function(s) {
		var n = 0, err = null, b, m, _tmp, _tmp$1;
		b = this;
		b.lastRead = 0;
		m = b.grow(s.length);
		_tmp = $copyString($subslice(b.buf, m), s); _tmp$1 = null; n = _tmp; err = _tmp$1;
		return [n, err];
	};
	Buffer.prototype.WriteString = function(s) { return this.$val.WriteString(s); };
	Buffer.Ptr.prototype.ReadFrom = function(r) {
		var n = new $Int64(0, 0), err = null, b, free, newBuf, x, _tuple, m, e, x$1, _tmp, _tmp$1, _tmp$2, _tmp$3;
		b = this;
		b.lastRead = 0;
		if (b.off >= b.buf.$length) {
			b.Truncate(0);
		}
		while (true) {
			free = b.buf.$capacity - b.buf.$length >> 0;
			if (free < 512) {
				newBuf = b.buf;
				if ((b.off + free >> 0) < 512) {
					newBuf = makeSlice((x = b.buf.$capacity, (((2 >>> 16 << 16) * x >> 0) + (2 << 16 >>> 16) * x) >> 0) + 512 >> 0);
				}
				$copySlice(newBuf, $subslice(b.buf, b.off));
				b.buf = $subslice(newBuf, 0, (b.buf.$length - b.off >> 0));
				b.off = 0;
			}
			_tuple = r.Read($subslice(b.buf, b.buf.$length, b.buf.$capacity)); m = _tuple[0]; e = _tuple[1];
			b.buf = $subslice(b.buf, 0, (b.buf.$length + m >> 0));
			n = (x$1 = new $Int64(0, m), new $Int64(n.$high + x$1.$high, n.$low + x$1.$low));
			if ($interfaceIsEqual(e, io.EOF)) {
				break;
			}
			if (!($interfaceIsEqual(e, null))) {
				_tmp = n; _tmp$1 = e; n = _tmp; err = _tmp$1;
				return [n, err];
			}
		}
		_tmp$2 = n; _tmp$3 = null; n = _tmp$2; err = _tmp$3;
		return [n, err];
	};
	Buffer.prototype.ReadFrom = function(r) { return this.$val.ReadFrom(r); };
	makeSlice = function(n) {
		var $deferred = [], $err = null;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			if (!($interfaceIsEqual($recover(), null))) {
				$panic($pkg.ErrTooLarge);
			}
		}), []]);
		return ($sliceType($Uint8)).make(n);
		/* */ } catch(err) { $err = err; return ($sliceType($Uint8)).nil; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	Buffer.Ptr.prototype.WriteTo = function(w) {
		var n = new $Int64(0, 0), err = null, b, nBytes, _tuple, m, e, _tmp, _tmp$1, _tmp$2, _tmp$3;
		b = this;
		b.lastRead = 0;
		if (b.off < b.buf.$length) {
			nBytes = b.Len();
			_tuple = w.Write($subslice(b.buf, b.off)); m = _tuple[0]; e = _tuple[1];
			if (m > nBytes) {
				$panic(new $String("bytes.Buffer.WriteTo: invalid Write count"));
			}
			b.off = b.off + (m) >> 0;
			n = new $Int64(0, m);
			if (!($interfaceIsEqual(e, null))) {
				_tmp = n; _tmp$1 = e; n = _tmp; err = _tmp$1;
				return [n, err];
			}
			if (!((m === nBytes))) {
				_tmp$2 = n; _tmp$3 = io.ErrShortWrite; n = _tmp$2; err = _tmp$3;
				return [n, err];
			}
		}
		b.Truncate(0);
		return [n, err];
	};
	Buffer.prototype.WriteTo = function(w) { return this.$val.WriteTo(w); };
	Buffer.Ptr.prototype.WriteByte = function(c) {
		var b, m, x;
		b = this;
		b.lastRead = 0;
		m = b.grow(1);
		(x = b.buf, (m < 0 || m >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + m] = c);
		return null;
	};
	Buffer.prototype.WriteByte = function(c) { return this.$val.WriteByte(c); };
	Buffer.Ptr.prototype.WriteRune = function(r) {
		var n = 0, err = null, b, _tmp, _tmp$1, _tmp$2, _tmp$3;
		b = this;
		if (r < 128) {
			b.WriteByte((r << 24 >>> 24));
			_tmp = 1; _tmp$1 = null; n = _tmp; err = _tmp$1;
			return [n, err];
		}
		n = utf8.EncodeRune($subslice(new ($sliceType($Uint8))(b.runeBytes), 0), r);
		b.Write($subslice(new ($sliceType($Uint8))(b.runeBytes), 0, n));
		_tmp$2 = n; _tmp$3 = null; n = _tmp$2; err = _tmp$3;
		return [n, err];
	};
	Buffer.prototype.WriteRune = function(r) { return this.$val.WriteRune(r); };
	Buffer.Ptr.prototype.Read = function(p) {
		var n = 0, err = null, b, _tmp, _tmp$1;
		b = this;
		b.lastRead = 0;
		if (b.off >= b.buf.$length) {
			b.Truncate(0);
			if (p.$length === 0) {
				return [n, err];
			}
			_tmp = 0; _tmp$1 = io.EOF; n = _tmp; err = _tmp$1;
			return [n, err];
		}
		n = $copySlice(p, $subslice(b.buf, b.off));
		b.off = b.off + (n) >> 0;
		if (n > 0) {
			b.lastRead = 2;
		}
		return [n, err];
	};
	Buffer.prototype.Read = function(p) { return this.$val.Read(p); };
	Buffer.Ptr.prototype.Next = function(n) {
		var b, m, data;
		b = this;
		b.lastRead = 0;
		m = b.Len();
		if (n > m) {
			n = m;
		}
		data = $subslice(b.buf, b.off, (b.off + n >> 0));
		b.off = b.off + (n) >> 0;
		if (n > 0) {
			b.lastRead = 2;
		}
		return data;
	};
	Buffer.prototype.Next = function(n) { return this.$val.Next(n); };
	Buffer.Ptr.prototype.ReadByte = function() {
		var c = 0, err = null, b, _tmp, _tmp$1, x, x$1, _tmp$2, _tmp$3;
		b = this;
		b.lastRead = 0;
		if (b.off >= b.buf.$length) {
			b.Truncate(0);
			_tmp = 0; _tmp$1 = io.EOF; c = _tmp; err = _tmp$1;
			return [c, err];
		}
		c = (x = b.buf, x$1 = b.off, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1]));
		b.off = b.off + (1) >> 0;
		b.lastRead = 2;
		_tmp$2 = c; _tmp$3 = null; c = _tmp$2; err = _tmp$3;
		return [c, err];
	};
	Buffer.prototype.ReadByte = function() { return this.$val.ReadByte(); };
	Buffer.Ptr.prototype.ReadRune = function() {
		var r = 0, size = 0, err = null, b, _tmp, _tmp$1, _tmp$2, x, x$1, c, _tmp$3, _tmp$4, _tmp$5, _tuple, n, _tmp$6, _tmp$7, _tmp$8;
		b = this;
		b.lastRead = 0;
		if (b.off >= b.buf.$length) {
			b.Truncate(0);
			_tmp = 0; _tmp$1 = 0; _tmp$2 = io.EOF; r = _tmp; size = _tmp$1; err = _tmp$2;
			return [r, size, err];
		}
		b.lastRead = 1;
		c = (x = b.buf, x$1 = b.off, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1]));
		if (c < 128) {
			b.off = b.off + (1) >> 0;
			_tmp$3 = (c >> 0); _tmp$4 = 1; _tmp$5 = null; r = _tmp$3; size = _tmp$4; err = _tmp$5;
			return [r, size, err];
		}
		_tuple = utf8.DecodeRune($subslice(b.buf, b.off)); r = _tuple[0]; n = _tuple[1];
		b.off = b.off + (n) >> 0;
		_tmp$6 = r; _tmp$7 = n; _tmp$8 = null; r = _tmp$6; size = _tmp$7; err = _tmp$8;
		return [r, size, err];
	};
	Buffer.prototype.ReadRune = function() { return this.$val.ReadRune(); };
	Buffer.Ptr.prototype.UnreadRune = function() {
		var b, _tuple, n;
		b = this;
		if (!((b.lastRead === 1))) {
			return errors.New("bytes.Buffer: UnreadRune: previous operation was not ReadRune");
		}
		b.lastRead = 0;
		if (b.off > 0) {
			_tuple = utf8.DecodeLastRune($subslice(b.buf, 0, b.off)); n = _tuple[1];
			b.off = b.off - (n) >> 0;
		}
		return null;
	};
	Buffer.prototype.UnreadRune = function() { return this.$val.UnreadRune(); };
	Buffer.Ptr.prototype.UnreadByte = function() {
		var b;
		b = this;
		if (!((b.lastRead === 1)) && !((b.lastRead === 2))) {
			return errors.New("bytes.Buffer: UnreadByte: previous operation was not a read");
		}
		b.lastRead = 0;
		if (b.off > 0) {
			b.off = b.off - (1) >> 0;
		}
		return null;
	};
	Buffer.prototype.UnreadByte = function() { return this.$val.UnreadByte(); };
	Buffer.Ptr.prototype.ReadBytes = function(delim) {
		var line = ($sliceType($Uint8)).nil, err = null, b, _tuple, slice;
		b = this;
		_tuple = b.readSlice(delim); slice = _tuple[0]; err = _tuple[1];
		line = $appendSlice(line, slice);
		return [line, err];
	};
	Buffer.prototype.ReadBytes = function(delim) { return this.$val.ReadBytes(delim); };
	Buffer.Ptr.prototype.readSlice = function(delim) {
		var line = ($sliceType($Uint8)).nil, err = null, b, i, end, _tmp, _tmp$1;
		b = this;
		i = IndexByte($subslice(b.buf, b.off), delim);
		end = (b.off + i >> 0) + 1 >> 0;
		if (i < 0) {
			end = b.buf.$length;
			err = io.EOF;
		}
		line = $subslice(b.buf, b.off, end);
		b.off = end;
		b.lastRead = 2;
		_tmp = line; _tmp$1 = err; line = _tmp; err = _tmp$1;
		return [line, err];
	};
	Buffer.prototype.readSlice = function(delim) { return this.$val.readSlice(delim); };
	Buffer.Ptr.prototype.ReadString = function(delim) {
		var line = "", err = null, b, _tuple, slice, _tmp, _tmp$1;
		b = this;
		_tuple = b.readSlice(delim); slice = _tuple[0]; err = _tuple[1];
		_tmp = $bytesToString(slice); _tmp$1 = err; line = _tmp; err = _tmp$1;
		return [line, err];
	};
	Buffer.prototype.ReadString = function(delim) { return this.$val.ReadString(delim); };
	NewBuffer = $pkg.NewBuffer = function(buf) {
		return new Buffer.Ptr(buf, 0, ($arrayType($Uint8, 4)).zero(), ($arrayType($Uint8, 64)).zero(), 0);
	};
	$pkg.$init = function() {
		($ptrType(Buffer)).methods = [["Bytes", "Bytes", "", [], [($sliceType($Uint8))], false, -1], ["Grow", "Grow", "", [$Int], [], false, -1], ["Len", "Len", "", [], [$Int], false, -1], ["Next", "Next", "", [$Int], [($sliceType($Uint8))], false, -1], ["Read", "Read", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["ReadByte", "ReadByte", "", [], [$Uint8, $error], false, -1], ["ReadBytes", "ReadBytes", "", [$Uint8], [($sliceType($Uint8)), $error], false, -1], ["ReadFrom", "ReadFrom", "", [io.Reader], [$Int64, $error], false, -1], ["ReadRune", "ReadRune", "", [], [$Int32, $Int, $error], false, -1], ["ReadString", "ReadString", "", [$Uint8], [$String, $error], false, -1], ["Reset", "Reset", "", [], [], false, -1], ["String", "String", "", [], [$String], false, -1], ["Truncate", "Truncate", "", [$Int], [], false, -1], ["UnreadByte", "UnreadByte", "", [], [$error], false, -1], ["UnreadRune", "UnreadRune", "", [], [$error], false, -1], ["Write", "Write", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["WriteByte", "WriteByte", "", [$Uint8], [$error], false, -1], ["WriteRune", "WriteRune", "", [$Int32], [$Int, $error], false, -1], ["WriteString", "WriteString", "", [$String], [$Int, $error], false, -1], ["WriteTo", "WriteTo", "", [io.Writer], [$Int64, $error], false, -1], ["grow", "grow", "bytes", [$Int], [$Int], false, -1], ["readSlice", "readSlice", "bytes", [$Uint8], [($sliceType($Uint8)), $error], false, -1]];
		Buffer.init([["buf", "buf", "bytes", ($sliceType($Uint8)), ""], ["off", "off", "bytes", $Int, ""], ["runeBytes", "runeBytes", "bytes", ($arrayType($Uint8, 4)), ""], ["bootstrap", "bootstrap", "bytes", ($arrayType($Uint8, 64)), ""], ["lastRead", "lastRead", "bytes", readOp, ""]]);
		$pkg.ErrTooLarge = errors.New("bytes.Buffer: too large");
	};
	return $pkg;
})();
$packages["math"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], math, zero, posInf, negInf, nan, pow10tab, init, IsInf, Ldexp, Float32bits, Float32frombits, Float64bits, init$1;
	init = function() {
		Float32bits(0);
		Float32frombits(0);
	};
	IsInf = $pkg.IsInf = function(f, sign) {
		if (f === posInf) {
			return sign >= 0;
		}
		if (f === negInf) {
			return sign <= 0;
		}
		return false;
	};
	Ldexp = $pkg.Ldexp = function(frac, exp$1) {
		if (frac === 0) {
			return frac;
		}
		if (exp$1 >= 1024) {
			return frac * $parseFloat(math.pow(2, 1023)) * $parseFloat(math.pow(2, exp$1 - 1023 >> 0));
		}
		if (exp$1 <= -1024) {
			return frac * $parseFloat(math.pow(2, -1023)) * $parseFloat(math.pow(2, exp$1 + 1023 >> 0));
		}
		return frac * $parseFloat(math.pow(2, exp$1));
	};
	Float32bits = $pkg.Float32bits = function(f) {
		var s, e, r;
		if ($float32IsEqual(f, 0)) {
			if ($float32IsEqual(1 / f, negInf)) {
				return 2147483648;
			}
			return 0;
		}
		if (!(($float32IsEqual(f, f)))) {
			return 2143289344;
		}
		s = 0;
		if (f < 0) {
			s = 2147483648;
			f = -f;
		}
		e = 150;
		while (f >= 1.6777216e+07) {
			f = f / (2);
			if (e === 255) {
				break;
			}
			e = e + (1) >>> 0;
		}
		while (f < 8.388608e+06) {
			e = e - (1) >>> 0;
			if (e === 0) {
				break;
			}
			f = f * (2);
		}
		r = $parseFloat($mod(f, 2));
		if ((r > 0.5 && r < 1) || r >= 1.5) {
			f = f + (1);
		}
		return (((s | (e << 23 >>> 0)) >>> 0) | (((f >> 0) & ~8388608))) >>> 0;
	};
	Float32frombits = $pkg.Float32frombits = function(b) {
		var s, e, m;
		s = 1;
		if (!((((b & 2147483648) >>> 0) === 0))) {
			s = -1;
		}
		e = (((b >>> 23 >>> 0)) & 255) >>> 0;
		m = (b & 8388607) >>> 0;
		if (e === 255) {
			if (m === 0) {
				return s / 0;
			}
			return nan;
		}
		if (!((e === 0))) {
			m = m + (8388608) >>> 0;
		}
		if (e === 0) {
			e = 1;
		}
		return Ldexp(m, ((e >> 0) - 127 >> 0) - 23 >> 0) * s;
	};
	Float64bits = $pkg.Float64bits = function(f) {
		var s, e, x, x$1, x$2, x$3;
		if (f === 0) {
			if (1 / f === negInf) {
				return new $Uint64(2147483648, 0);
			}
			return new $Uint64(0, 0);
		}
		if (!((f === f))) {
			return new $Uint64(2146959360, 1);
		}
		s = new $Uint64(0, 0);
		if (f < 0) {
			s = new $Uint64(2147483648, 0);
			f = -f;
		}
		e = 1075;
		while (f >= 9.007199254740992e+15) {
			f = f / (2);
			if (e === 2047) {
				break;
			}
			e = e + (1) >>> 0;
		}
		while (f < 4.503599627370496e+15) {
			e = e - (1) >>> 0;
			if (e === 0) {
				break;
			}
			f = f * (2);
		}
		return (x = (x$1 = $shiftLeft64(new $Uint64(0, e), 52), new $Uint64(s.$high | x$1.$high, (s.$low | x$1.$low) >>> 0)), x$2 = (x$3 = new $Uint64(0, f), new $Uint64(x$3.$high &~ 1048576, (x$3.$low &~ 0) >>> 0)), new $Uint64(x.$high | x$2.$high, (x.$low | x$2.$low) >>> 0));
	};
	init$1 = function() {
		var i, _q, m, x;
		pow10tab[0] = 1;
		pow10tab[1] = 10;
		i = 2;
		while (i < 70) {
			m = (_q = i / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
			(i < 0 || i >= pow10tab.length) ? $throwRuntimeError("index out of range") : pow10tab[i] = ((m < 0 || m >= pow10tab.length) ? $throwRuntimeError("index out of range") : pow10tab[m]) * (x = i - m >> 0, ((x < 0 || x >= pow10tab.length) ? $throwRuntimeError("index out of range") : pow10tab[x]));
			i = i + (1) >> 0;
		}
	};
	$pkg.$init = function() {
		pow10tab = ($arrayType($Float64, 70)).zero();
		math = $global.Math;
		zero = 0;
		posInf = 1 / zero;
		negInf = -1 / zero;
		nan = 0 / zero;
		init();
		init$1();
	};
	return $pkg;
})();
$packages["syscall"] = (function() {
	var $pkg = {}, bytes = $packages["bytes"], js = $packages["github.com/gopherjs/gopherjs/js"], sync = $packages["sync"], runtime = $packages["runtime"], mmapper, Errno, Timespec, Stat_t, Dirent, warningPrinted, lineBuffer, syscallModule, alreadyTriedToLoad, minusOne, envOnce, envLock, env, envs, mapper, errors, printWarning, printToConsole, init, syscall, Syscall, Syscall6, BytePtrFromString, copyenv, Getenv, itoa, Open, clen, ReadDirent, ParseDirent, Read, Write, open, Close, Fchdir, Fchmod, Fsync, Getdents, read, write, munmap, Fchown, Fstat, Ftruncate, Lstat, Pread, Pwrite, Seek, mmap;
	mmapper = $pkg.mmapper = $newType(0, "Struct", "syscall.mmapper", "mmapper", "syscall", function(Mutex_, active_, mmap_, munmap_) {
		this.$val = this;
		this.Mutex = Mutex_ !== undefined ? Mutex_ : new sync.Mutex.Ptr();
		this.active = active_ !== undefined ? active_ : false;
		this.mmap = mmap_ !== undefined ? mmap_ : $throwNilPointerError;
		this.munmap = munmap_ !== undefined ? munmap_ : $throwNilPointerError;
	});
	Errno = $pkg.Errno = $newType(4, "Uintptr", "syscall.Errno", "Errno", "syscall", null);
	Timespec = $pkg.Timespec = $newType(0, "Struct", "syscall.Timespec", "Timespec", "syscall", function(Sec_, Nsec_) {
		this.$val = this;
		this.Sec = Sec_ !== undefined ? Sec_ : new $Int64(0, 0);
		this.Nsec = Nsec_ !== undefined ? Nsec_ : new $Int64(0, 0);
	});
	Stat_t = $pkg.Stat_t = $newType(0, "Struct", "syscall.Stat_t", "Stat_t", "syscall", function(Dev_, Ino_, Nlink_, Mode_, Uid_, Gid_, X__pad0_, Rdev_, Size_, Blksize_, Blocks_, Atim_, Mtim_, Ctim_, X__unused_) {
		this.$val = this;
		this.Dev = Dev_ !== undefined ? Dev_ : new $Uint64(0, 0);
		this.Ino = Ino_ !== undefined ? Ino_ : new $Uint64(0, 0);
		this.Nlink = Nlink_ !== undefined ? Nlink_ : new $Uint64(0, 0);
		this.Mode = Mode_ !== undefined ? Mode_ : 0;
		this.Uid = Uid_ !== undefined ? Uid_ : 0;
		this.Gid = Gid_ !== undefined ? Gid_ : 0;
		this.X__pad0 = X__pad0_ !== undefined ? X__pad0_ : 0;
		this.Rdev = Rdev_ !== undefined ? Rdev_ : new $Uint64(0, 0);
		this.Size = Size_ !== undefined ? Size_ : new $Int64(0, 0);
		this.Blksize = Blksize_ !== undefined ? Blksize_ : new $Int64(0, 0);
		this.Blocks = Blocks_ !== undefined ? Blocks_ : new $Int64(0, 0);
		this.Atim = Atim_ !== undefined ? Atim_ : new Timespec.Ptr();
		this.Mtim = Mtim_ !== undefined ? Mtim_ : new Timespec.Ptr();
		this.Ctim = Ctim_ !== undefined ? Ctim_ : new Timespec.Ptr();
		this.X__unused = X__unused_ !== undefined ? X__unused_ : ($arrayType($Int64, 3)).zero();
	});
	Dirent = $pkg.Dirent = $newType(0, "Struct", "syscall.Dirent", "Dirent", "syscall", function(Ino_, Off_, Reclen_, Type_, Name_, Pad_cgo_0_) {
		this.$val = this;
		this.Ino = Ino_ !== undefined ? Ino_ : new $Uint64(0, 0);
		this.Off = Off_ !== undefined ? Off_ : new $Int64(0, 0);
		this.Reclen = Reclen_ !== undefined ? Reclen_ : 0;
		this.Type = Type_ !== undefined ? Type_ : 0;
		this.Name = Name_ !== undefined ? Name_ : ($arrayType($Int8, 256)).zero();
		this.Pad_cgo_0 = Pad_cgo_0_ !== undefined ? Pad_cgo_0_ : ($arrayType($Uint8, 5)).zero();
	});
	printWarning = function() {
		if (!warningPrinted) {
			console.log("warning: system calls not available, see https://github.com/gopherjs/gopherjs/blob/master/doc/syscalls.md");
		}
		warningPrinted = true;
	};
	printToConsole = function(b) {
		var goPrintToConsole, i;
		goPrintToConsole = $global.goPrintToConsole;
		if (!(goPrintToConsole === undefined)) {
			goPrintToConsole(b);
			return;
		}
		lineBuffer = $appendSlice(lineBuffer, b);
		while (true) {
			i = bytes.IndexByte(lineBuffer, 10);
			if (i === -1) {
				break;
			}
			$global.console.log($externalize($bytesToString($subslice(lineBuffer, 0, i)), $String));
			lineBuffer = $subslice(lineBuffer, (i + 1 >> 0));
		}
	};
	init = function() {
		var process, jsEnv, envkeys, i, key;
		process = $global.process;
		if (!(process === undefined)) {
			jsEnv = process.env;
			envkeys = $global.Object.keys(jsEnv);
			envs = ($sliceType($String)).make($parseInt(envkeys.length));
			i = 0;
			while (i < $parseInt(envkeys.length)) {
				key = $internalize(envkeys[i], $String);
				(i < 0 || i >= envs.$length) ? $throwRuntimeError("index out of range") : envs.$array[envs.$offset + i] = key + "=" + $internalize(jsEnv[$externalize(key, $String)], $String);
				i = i + (1) >> 0;
			}
		}
	};
	syscall = function(name) {
		var $deferred = [], $err = null, require;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			$recover();
		}), []]);
		if (syscallModule === null) {
			if (alreadyTriedToLoad) {
				return null;
			}
			alreadyTriedToLoad = true;
			require = $global.require;
			if (require === undefined) {
				$panic(new $String(""));
			}
			syscallModule = require($externalize("syscall", $String));
		}
		return syscallModule[$externalize(name, $String)];
		/* */ } catch(err) { $err = err; return null; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	Syscall = $pkg.Syscall = function(trap, a1, a2, a3) {
		var r1 = 0, r2 = 0, err = 0, f, r, _tmp, _tmp$1, _tmp$2, array, slice, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8;
		f = syscall("Syscall");
		if (!(f === null)) {
			r = f(trap, a1, a2, a3);
			_tmp = (($parseInt(r[0]) >> 0) >>> 0); _tmp$1 = (($parseInt(r[1]) >> 0) >>> 0); _tmp$2 = (($parseInt(r[2]) >> 0) >>> 0); r1 = _tmp; r2 = _tmp$1; err = _tmp$2;
			return [r1, r2, err];
		}
		if ((trap === 1) && ((a1 === 1) || (a1 === 2))) {
			array = a2;
			slice = ($sliceType($Uint8)).make($parseInt(array.length));
			slice.$array = array;
			printToConsole(slice);
			_tmp$3 = ($parseInt(array.length) >>> 0); _tmp$4 = 0; _tmp$5 = 0; r1 = _tmp$3; r2 = _tmp$4; err = _tmp$5;
			return [r1, r2, err];
		}
		printWarning();
		_tmp$6 = (minusOne >>> 0); _tmp$7 = 0; _tmp$8 = 13; r1 = _tmp$6; r2 = _tmp$7; err = _tmp$8;
		return [r1, r2, err];
	};
	Syscall6 = $pkg.Syscall6 = function(trap, a1, a2, a3, a4, a5, a6) {
		var r1 = 0, r2 = 0, err = 0, f, r, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5;
		f = syscall("Syscall6");
		if (!(f === null)) {
			r = f(trap, a1, a2, a3, a4, a5, a6);
			_tmp = (($parseInt(r[0]) >> 0) >>> 0); _tmp$1 = (($parseInt(r[1]) >> 0) >>> 0); _tmp$2 = (($parseInt(r[2]) >> 0) >>> 0); r1 = _tmp; r2 = _tmp$1; err = _tmp$2;
			return [r1, r2, err];
		}
		if (!((trap === 202))) {
			printWarning();
		}
		_tmp$3 = (minusOne >>> 0); _tmp$4 = 0; _tmp$5 = 13; r1 = _tmp$3; r2 = _tmp$4; err = _tmp$5;
		return [r1, r2, err];
	};
	BytePtrFromString = $pkg.BytePtrFromString = function(s) {
		var array, _ref, _i, i, b;
		array = new ($global.Uint8Array)(s.length + 1 >> 0);
		_ref = new ($sliceType($Uint8))($stringToBytes(s));
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			b = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (b === 0) {
				return [($ptrType($Uint8)).nil, new Errno(22)];
			}
			array[i] = b;
			_i++;
		}
		array[s.length] = 0;
		return [array, null];
	};
	copyenv = function() {
		var _ref, _i, i, s, j, key, _tuple, _entry, ok, _key;
		env = new $Map();
		_ref = envs;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			s = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			j = 0;
			while (j < s.length) {
				if (s.charCodeAt(j) === 61) {
					key = s.substring(0, j);
					_tuple = (_entry = env[key], _entry !== undefined ? [_entry.v, true] : [0, false]); ok = _tuple[1];
					if (!ok) {
						_key = key; (env || $throwRuntimeError("assignment to entry in nil map"))[_key] = { k: _key, v: i };
					}
					break;
				}
				j = j + (1) >> 0;
			}
			_i++;
		}
	};
	Getenv = $pkg.Getenv = function(key) {
		var value = "", found = false, $deferred = [], $err = null, _tmp, _tmp$1, _recv, _tuple, _entry, i, ok, _tmp$2, _tmp$3, s, i$1, _tmp$4, _tmp$5, _tmp$6, _tmp$7;
		/* */ try { $deferFrames.push($deferred);
		envOnce.Do(copyenv);
		if (key.length === 0) {
			_tmp = ""; _tmp$1 = false; value = _tmp; found = _tmp$1;
			return [value, found];
		}
		envLock.RLock();
		$deferred.push([(_recv = envLock, function() { $stackDepthOffset--; try { return _recv.RUnlock(); } finally { $stackDepthOffset++; } }), []]);
		_tuple = (_entry = env[key], _entry !== undefined ? [_entry.v, true] : [0, false]); i = _tuple[0]; ok = _tuple[1];
		if (!ok) {
			_tmp$2 = ""; _tmp$3 = false; value = _tmp$2; found = _tmp$3;
			return [value, found];
		}
		s = ((i < 0 || i >= envs.$length) ? $throwRuntimeError("index out of range") : envs.$array[envs.$offset + i]);
		i$1 = 0;
		while (i$1 < s.length) {
			if (s.charCodeAt(i$1) === 61) {
				_tmp$4 = s.substring((i$1 + 1 >> 0)); _tmp$5 = true; value = _tmp$4; found = _tmp$5;
				return [value, found];
			}
			i$1 = i$1 + (1) >> 0;
		}
		_tmp$6 = ""; _tmp$7 = false; value = _tmp$6; found = _tmp$7;
		return [value, found];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [value, found]; }
	};
	itoa = function(val) {
		var buf, i, _r, _q;
		if (val < 0) {
			return "-" + itoa(-val);
		}
		buf = ($arrayType($Uint8, 32)).zero(); $copy(buf, ($arrayType($Uint8, 32)).zero(), ($arrayType($Uint8, 32)));
		i = 31;
		while (val >= 10) {
			(i < 0 || i >= buf.length) ? $throwRuntimeError("index out of range") : buf[i] = (((_r = val % 10, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) + 48 >> 0) << 24 >>> 24);
			i = i - (1) >> 0;
			val = (_q = val / (10), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		}
		(i < 0 || i >= buf.length) ? $throwRuntimeError("index out of range") : buf[i] = ((val + 48 >> 0) << 24 >>> 24);
		return $bytesToString($subslice(new ($sliceType($Uint8))(buf), i));
	};
	Timespec.Ptr.prototype.Unix = function() {
		var sec = new $Int64(0, 0), nsec = new $Int64(0, 0), ts, _tmp, _tmp$1;
		ts = this;
		_tmp = ts.Sec; _tmp$1 = ts.Nsec; sec = _tmp; nsec = _tmp$1;
		return [sec, nsec];
	};
	Timespec.prototype.Unix = function() { return this.$val.Unix(); };
	Timespec.Ptr.prototype.Nano = function() {
		var ts, x, x$1;
		ts = this;
		return (x = $mul64(ts.Sec, new $Int64(0, 1000000000)), x$1 = ts.Nsec, new $Int64(x.$high + x$1.$high, x.$low + x$1.$low));
	};
	Timespec.prototype.Nano = function() { return this.$val.Nano(); };
	Open = $pkg.Open = function(path, mode, perm) {
		var fd = 0, err = null, _tuple;
		_tuple = open(path, mode | 0, perm); fd = _tuple[0]; err = _tuple[1];
		return [fd, err];
	};
	clen = function(n) {
		var i;
		i = 0;
		while (i < n.$length) {
			if (((i < 0 || i >= n.$length) ? $throwRuntimeError("index out of range") : n.$array[n.$offset + i]) === 0) {
				return i;
			}
			i = i + (1) >> 0;
		}
		return n.$length;
	};
	ReadDirent = $pkg.ReadDirent = function(fd, buf) {
		var n = 0, err = null, _tuple;
		_tuple = Getdents(fd, buf); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	ParseDirent = $pkg.ParseDirent = function(buf, max, names) {
		var consumed = 0, count = 0, newnames = ($sliceType($String)).nil, origlen, dirent, _array, _struct, _view, x, bytes$1, name, _tmp, _tmp$1, _tmp$2;
		origlen = buf.$length;
		count = 0;
		while (!((max === 0)) && buf.$length > 0) {
			dirent = [undefined];
			dirent[0] = (_array = $sliceToArray(buf), _struct = new Dirent.Ptr(), _view = new DataView(_array.buffer, _array.byteOffset), _struct.Ino = new $Uint64(_view.getUint32(4, true), _view.getUint32(0, true)), _struct.Off = new $Int64(_view.getUint32(12, true), _view.getUint32(8, true)), _struct.Reclen = _view.getUint16(16, true), _struct.Type = _view.getUint8(18, true), _struct.Name = new ($nativeArray("Int8"))(_array.buffer, $min(_array.byteOffset + 19, _array.buffer.byteLength)), _struct.Pad_cgo_0 = new ($nativeArray("Uint8"))(_array.buffer, $min(_array.byteOffset + 275, _array.buffer.byteLength)), _struct);
			buf = $subslice(buf, dirent[0].Reclen);
			if ((x = dirent[0].Ino, (x.$high === 0 && x.$low === 0))) {
				continue;
			}
			bytes$1 = $sliceToArray(new ($sliceType($Uint8))(dirent[0].Name));
			name = $bytesToString($subslice(new ($sliceType($Uint8))(bytes$1), 0, clen(new ($sliceType($Uint8))(bytes$1))));
			if (name === "." || name === "..") {
				continue;
			}
			max = max - (1) >> 0;
			count = count + (1) >> 0;
			names = $append(names, name);
		}
		_tmp = origlen - buf.$length >> 0; _tmp$1 = count; _tmp$2 = names; consumed = _tmp; count = _tmp$1; newnames = _tmp$2;
		return [consumed, count, newnames];
	};
	mmapper.Ptr.prototype.Mmap = function(fd, offset, length, prot, flags) {
		var data = ($sliceType($Uint8)).nil, err = null, $deferred = [], $err = null, m, _tmp, _tmp$1, _tuple, addr, errno, _tmp$2, _tmp$3, sl, b, x, x$1, p, _recv, _key, _tmp$4, _tmp$5;
		/* */ try { $deferFrames.push($deferred);
		m = this;
		if (length <= 0) {
			_tmp = ($sliceType($Uint8)).nil; _tmp$1 = new Errno(22); data = _tmp; err = _tmp$1;
			return [data, err];
		}
		_tuple = m.mmap(0, (length >>> 0), prot, flags, fd, offset); addr = _tuple[0]; errno = _tuple[1];
		if (!($interfaceIsEqual(errno, null))) {
			_tmp$2 = ($sliceType($Uint8)).nil; _tmp$3 = errno; data = _tmp$2; err = _tmp$3;
			return [data, err];
		}
		sl = new ($structType([["addr", "addr", "syscall", $Uintptr, ""], ["len", "len", "syscall", $Int, ""], ["cap", "cap", "syscall", $Int, ""]])).Ptr(addr, length, length);
		b = sl;
		p = new ($ptrType($Uint8))(function() { return (x$1 = b.$capacity - 1 >> 0, ((x$1 < 0 || x$1 >= this.$target.$length) ? $throwRuntimeError("index out of range") : this.$target.$array[this.$target.$offset + x$1])); }, function($v) { (x = b.$capacity - 1 >> 0, (x < 0 || x >= this.$target.$length) ? $throwRuntimeError("index out of range") : this.$target.$array[this.$target.$offset + x] = $v); }, b);
		m.Mutex.Lock();
		$deferred.push([(_recv = m, function() { $stackDepthOffset--; try { return _recv.Unlock(); } finally { $stackDepthOffset++; } }), []]);
		_key = p; (m.active || $throwRuntimeError("assignment to entry in nil map"))[_key.$key()] = { k: _key, v: b };
		_tmp$4 = b; _tmp$5 = null; data = _tmp$4; err = _tmp$5;
		return [data, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [data, err]; }
	};
	mmapper.prototype.Mmap = function(fd, offset, length, prot, flags) { return this.$val.Mmap(fd, offset, length, prot, flags); };
	mmapper.Ptr.prototype.Munmap = function(data) {
		var err = null, $deferred = [], $err = null, m, x, x$1, p, _recv, _entry, b, errno;
		/* */ try { $deferFrames.push($deferred);
		m = this;
		if ((data.$length === 0) || !((data.$length === data.$capacity))) {
			err = new Errno(22);
			return err;
		}
		p = new ($ptrType($Uint8))(function() { return (x$1 = data.$capacity - 1 >> 0, ((x$1 < 0 || x$1 >= this.$target.$length) ? $throwRuntimeError("index out of range") : this.$target.$array[this.$target.$offset + x$1])); }, function($v) { (x = data.$capacity - 1 >> 0, (x < 0 || x >= this.$target.$length) ? $throwRuntimeError("index out of range") : this.$target.$array[this.$target.$offset + x] = $v); }, data);
		m.Mutex.Lock();
		$deferred.push([(_recv = m, function() { $stackDepthOffset--; try { return _recv.Unlock(); } finally { $stackDepthOffset++; } }), []]);
		b = (_entry = m.active[p.$key()], _entry !== undefined ? _entry.v : ($sliceType($Uint8)).nil);
		if (b === ($sliceType($Uint8)).nil || !($sliceIsEqual(b, 0, data, 0))) {
			err = new Errno(22);
			return err;
		}
		errno = m.munmap($sliceToArray(b), (b.$length >>> 0));
		if (!($interfaceIsEqual(errno, null))) {
			err = errno;
			return err;
		}
		delete m.active[p.$key()];
		err = null;
		return err;
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return err; }
	};
	mmapper.prototype.Munmap = function(data) { return this.$val.Munmap(data); };
	Errno.prototype.Error = function() {
		var e, s;
		e = this.$val;
		if (0 <= (e >> 0) && (e >> 0) < 133) {
			s = ((e < 0 || e >= errors.length) ? $throwRuntimeError("index out of range") : errors[e]);
			if (!(s === "")) {
				return s;
			}
		}
		return "errno " + itoa((e >> 0));
	};
	$ptrType(Errno).prototype.Error = function() { return new Errno(this.$get()).Error(); };
	Errno.prototype.Temporary = function() {
		var e;
		e = this.$val;
		return (e === 4) || (e === 24) || (new Errno(e)).Timeout();
	};
	$ptrType(Errno).prototype.Temporary = function() { return new Errno(this.$get()).Temporary(); };
	Errno.prototype.Timeout = function() {
		var e;
		e = this.$val;
		return (e === 11) || (e === 11) || (e === 110);
	};
	$ptrType(Errno).prototype.Timeout = function() { return new Errno(this.$get()).Timeout(); };
	Read = $pkg.Read = function(fd, p) {
		var n = 0, err = null, _tuple;
		_tuple = read(fd, p); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	Write = $pkg.Write = function(fd, p) {
		var n = 0, err = null, _tuple;
		_tuple = write(fd, p); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	open = function(path, mode, perm) {
		var fd = 0, err = null, _p0, _tuple, _tuple$1, r0, e1;
		_p0 = ($ptrType($Uint8)).nil;
		_tuple = BytePtrFromString(path); _p0 = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, null))) {
			return [fd, err];
		}
		_tuple$1 = Syscall(2, _p0, (mode >>> 0), (perm >>> 0)); r0 = _tuple$1[0]; e1 = _tuple$1[2];
		fd = (r0 >> 0);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [fd, err];
	};
	Close = $pkg.Close = function(fd) {
		var err = null, _tuple, e1;
		_tuple = Syscall(3, (fd >>> 0), 0, 0); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Fchdir = $pkg.Fchdir = function(fd) {
		var err = null, _tuple, e1;
		_tuple = Syscall(81, (fd >>> 0), 0, 0); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Fchmod = $pkg.Fchmod = function(fd, mode) {
		var err = null, _tuple, e1;
		_tuple = Syscall(91, (fd >>> 0), (mode >>> 0), 0); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Fsync = $pkg.Fsync = function(fd) {
		var err = null, _tuple, e1;
		_tuple = Syscall(74, (fd >>> 0), 0, 0); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Getdents = $pkg.Getdents = function(fd, buf) {
		var n = 0, err = null, _p0, _tuple, r0, e1;
		_p0 = 0;
		if (buf.$length > 0) {
			_p0 = $sliceToArray(buf);
		} else {
			_p0 = new Uint8Array(0);
		}
		_tuple = Syscall(217, (fd >>> 0), _p0, (buf.$length >>> 0)); r0 = _tuple[0]; e1 = _tuple[2];
		n = (r0 >> 0);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [n, err];
	};
	read = function(fd, p) {
		var n = 0, err = null, _p0, _tuple, r0, e1;
		_p0 = 0;
		if (p.$length > 0) {
			_p0 = $sliceToArray(p);
		} else {
			_p0 = new Uint8Array(0);
		}
		_tuple = Syscall(0, (fd >>> 0), _p0, (p.$length >>> 0)); r0 = _tuple[0]; e1 = _tuple[2];
		n = (r0 >> 0);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [n, err];
	};
	write = function(fd, p) {
		var n = 0, err = null, _p0, _tuple, r0, e1;
		_p0 = 0;
		if (p.$length > 0) {
			_p0 = $sliceToArray(p);
		} else {
			_p0 = new Uint8Array(0);
		}
		_tuple = Syscall(1, (fd >>> 0), _p0, (p.$length >>> 0)); r0 = _tuple[0]; e1 = _tuple[2];
		n = (r0 >> 0);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [n, err];
	};
	munmap = function(addr, length) {
		var err = null, _tuple, e1;
		_tuple = Syscall(11, addr, length, 0); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Fchown = $pkg.Fchown = function(fd, uid, gid) {
		var err = null, _tuple, e1;
		_tuple = Syscall(93, (fd >>> 0), (uid >>> 0), (gid >>> 0)); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Fstat = $pkg.Fstat = function(fd, stat) {
		var err = null, _tuple, _array, _struct, _view, e1;
		_array = new Uint8Array(144);
		_tuple = Syscall(5, (fd >>> 0), _array, 0); e1 = _tuple[2];
		_struct = stat, _view = new DataView(_array.buffer, _array.byteOffset), _struct.Dev = new $Uint64(_view.getUint32(4, true), _view.getUint32(0, true)), _struct.Ino = new $Uint64(_view.getUint32(12, true), _view.getUint32(8, true)), _struct.Nlink = new $Uint64(_view.getUint32(20, true), _view.getUint32(16, true)), _struct.Mode = _view.getUint32(24, true), _struct.Uid = _view.getUint32(28, true), _struct.Gid = _view.getUint32(32, true), _struct.X__pad0 = _view.getInt32(36, true), _struct.Rdev = new $Uint64(_view.getUint32(44, true), _view.getUint32(40, true)), _struct.Size = new $Int64(_view.getUint32(52, true), _view.getUint32(48, true)), _struct.Blksize = new $Int64(_view.getUint32(60, true), _view.getUint32(56, true)), _struct.Blocks = new $Int64(_view.getUint32(68, true), _view.getUint32(64, true)), _struct.Atim.Sec = new $Int64(_view.getUint32(76, true), _view.getUint32(72, true)), _struct.Atim.Nsec = new $Int64(_view.getUint32(84, true), _view.getUint32(80, true)), _struct.Mtim.Sec = new $Int64(_view.getUint32(92, true), _view.getUint32(88, true)), _struct.Mtim.Nsec = new $Int64(_view.getUint32(100, true), _view.getUint32(96, true)), _struct.Ctim.Sec = new $Int64(_view.getUint32(108, true), _view.getUint32(104, true)), _struct.Ctim.Nsec = new $Int64(_view.getUint32(116, true), _view.getUint32(112, true)), _struct.X__unused = new ($nativeArray("Int64"))(_array.buffer, $min(_array.byteOffset + 120, _array.buffer.byteLength));
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Ftruncate = $pkg.Ftruncate = function(fd, length) {
		var err = null, _tuple, e1;
		_tuple = Syscall(77, (fd >>> 0), (length.$low >>> 0), 0); e1 = _tuple[2];
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Lstat = $pkg.Lstat = function(path, stat) {
		var err = null, _p0, _tuple, _tuple$1, _array, _struct, _view, e1;
		_p0 = ($ptrType($Uint8)).nil;
		_tuple = BytePtrFromString(path); _p0 = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, null))) {
			return err;
		}
		_array = new Uint8Array(144);
		_tuple$1 = Syscall(6, _p0, _array, 0); e1 = _tuple$1[2];
		_struct = stat, _view = new DataView(_array.buffer, _array.byteOffset), _struct.Dev = new $Uint64(_view.getUint32(4, true), _view.getUint32(0, true)), _struct.Ino = new $Uint64(_view.getUint32(12, true), _view.getUint32(8, true)), _struct.Nlink = new $Uint64(_view.getUint32(20, true), _view.getUint32(16, true)), _struct.Mode = _view.getUint32(24, true), _struct.Uid = _view.getUint32(28, true), _struct.Gid = _view.getUint32(32, true), _struct.X__pad0 = _view.getInt32(36, true), _struct.Rdev = new $Uint64(_view.getUint32(44, true), _view.getUint32(40, true)), _struct.Size = new $Int64(_view.getUint32(52, true), _view.getUint32(48, true)), _struct.Blksize = new $Int64(_view.getUint32(60, true), _view.getUint32(56, true)), _struct.Blocks = new $Int64(_view.getUint32(68, true), _view.getUint32(64, true)), _struct.Atim.Sec = new $Int64(_view.getUint32(76, true), _view.getUint32(72, true)), _struct.Atim.Nsec = new $Int64(_view.getUint32(84, true), _view.getUint32(80, true)), _struct.Mtim.Sec = new $Int64(_view.getUint32(92, true), _view.getUint32(88, true)), _struct.Mtim.Nsec = new $Int64(_view.getUint32(100, true), _view.getUint32(96, true)), _struct.Ctim.Sec = new $Int64(_view.getUint32(108, true), _view.getUint32(104, true)), _struct.Ctim.Nsec = new $Int64(_view.getUint32(116, true), _view.getUint32(112, true)), _struct.X__unused = new ($nativeArray("Int64"))(_array.buffer, $min(_array.byteOffset + 120, _array.buffer.byteLength));
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return err;
	};
	Pread = $pkg.Pread = function(fd, p, offset) {
		var n = 0, err = null, _p0, _tuple, r0, e1;
		_p0 = 0;
		if (p.$length > 0) {
			_p0 = $sliceToArray(p);
		} else {
			_p0 = new Uint8Array(0);
		}
		_tuple = Syscall6(17, (fd >>> 0), _p0, (p.$length >>> 0), (offset.$low >>> 0), 0, 0); r0 = _tuple[0]; e1 = _tuple[2];
		n = (r0 >> 0);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [n, err];
	};
	Pwrite = $pkg.Pwrite = function(fd, p, offset) {
		var n = 0, err = null, _p0, _tuple, r0, e1;
		_p0 = 0;
		if (p.$length > 0) {
			_p0 = $sliceToArray(p);
		} else {
			_p0 = new Uint8Array(0);
		}
		_tuple = Syscall6(18, (fd >>> 0), _p0, (p.$length >>> 0), (offset.$low >>> 0), 0, 0); r0 = _tuple[0]; e1 = _tuple[2];
		n = (r0 >> 0);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [n, err];
	};
	Seek = $pkg.Seek = function(fd, offset, whence) {
		var off = new $Int64(0, 0), err = null, _tuple, r0, e1;
		_tuple = Syscall(8, (fd >>> 0), (offset.$low >>> 0), (whence >>> 0)); r0 = _tuple[0]; e1 = _tuple[2];
		off = new $Int64(0, r0.constructor === Number ? r0 : 1);
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [off, err];
	};
	mmap = function(addr, length, prot, flags, fd, offset) {
		var xaddr = 0, err = null, _tuple, r0, e1;
		_tuple = Syscall6(9, addr, length, (prot >>> 0), (flags >>> 0), (fd >>> 0), (offset.$low >>> 0)); r0 = _tuple[0]; e1 = _tuple[2];
		xaddr = r0;
		if (!((e1 === 0))) {
			err = new Errno(e1);
		}
		return [xaddr, err];
	};
	$pkg.$init = function() {
		($ptrType(mmapper)).methods = [["Lock", "Lock", "", [], [], false, 0], ["Mmap", "Mmap", "", [$Int, $Int64, $Int, $Int, $Int], [($sliceType($Uint8)), $error], false, -1], ["Munmap", "Munmap", "", [($sliceType($Uint8))], [$error], false, -1], ["Unlock", "Unlock", "", [], [], false, 0]];
		mmapper.init([["Mutex", "", "", sync.Mutex, ""], ["active", "active", "syscall", ($mapType(($ptrType($Uint8)), ($sliceType($Uint8)))), ""], ["mmap", "mmap", "syscall", ($funcType([$Uintptr, $Uintptr, $Int, $Int, $Int, $Int64], [$Uintptr, $error], false)), ""], ["munmap", "munmap", "syscall", ($funcType([$Uintptr, $Uintptr], [$error], false)), ""]]);
		Errno.methods = [["Error", "Error", "", [], [$String], false, -1], ["Temporary", "Temporary", "", [], [$Bool], false, -1], ["Timeout", "Timeout", "", [], [$Bool], false, -1]];
		($ptrType(Errno)).methods = [["Error", "Error", "", [], [$String], false, -1], ["Temporary", "Temporary", "", [], [$Bool], false, -1], ["Timeout", "Timeout", "", [], [$Bool], false, -1]];
		($ptrType(Timespec)).methods = [["Nano", "Nano", "", [], [$Int64], false, -1], ["Unix", "Unix", "", [], [$Int64, $Int64], false, -1]];
		Timespec.init([["Sec", "Sec", "", $Int64, ""], ["Nsec", "Nsec", "", $Int64, ""]]);
		Stat_t.init([["Dev", "Dev", "", $Uint64, ""], ["Ino", "Ino", "", $Uint64, ""], ["Nlink", "Nlink", "", $Uint64, ""], ["Mode", "Mode", "", $Uint32, ""], ["Uid", "Uid", "", $Uint32, ""], ["Gid", "Gid", "", $Uint32, ""], ["X__pad0", "X__pad0", "", $Int32, ""], ["Rdev", "Rdev", "", $Uint64, ""], ["Size", "Size", "", $Int64, ""], ["Blksize", "Blksize", "", $Int64, ""], ["Blocks", "Blocks", "", $Int64, ""], ["Atim", "Atim", "", Timespec, ""], ["Mtim", "Mtim", "", Timespec, ""], ["Ctim", "Ctim", "", Timespec, ""], ["X__unused", "X__unused", "", ($arrayType($Int64, 3)), ""]]);
		Dirent.init([["Ino", "Ino", "", $Uint64, ""], ["Off", "Off", "", $Int64, ""], ["Reclen", "Reclen", "", $Uint16, ""], ["Type", "Type", "", $Uint8, ""], ["Name", "Name", "", ($arrayType($Int8, 256)), ""], ["Pad_cgo_0", "Pad_cgo_0", "", ($arrayType($Uint8, 5)), ""]]);
		lineBuffer = ($sliceType($Uint8)).nil;
		syscallModule = null;
		envOnce = new sync.Once.Ptr();
		envLock = new sync.RWMutex.Ptr();
		env = false;
		envs = ($sliceType($String)).nil;
		warningPrinted = false;
		alreadyTriedToLoad = false;
		minusOne = -1;
		$pkg.Stdin = 0;
		$pkg.Stdout = 1;
		$pkg.Stderr = 2;
		errors = $toNativeArray("String", ["", "operation not permitted", "no such file or directory", "no such process", "interrupted system call", "input/output error", "no such device or address", "argument list too long", "exec format error", "bad file descriptor", "no child processes", "resource temporarily unavailable", "cannot allocate memory", "permission denied", "bad address", "block device required", "device or resource busy", "file exists", "invalid cross-device link", "no such device", "not a directory", "is a directory", "invalid argument", "too many open files in system", "too many open files", "inappropriate ioctl for device", "text file busy", "file too large", "no space left on device", "illegal seek", "read-only file system", "too many links", "broken pipe", "numerical argument out of domain", "numerical result out of range", "resource deadlock avoided", "file name too long", "no locks available", "function not implemented", "directory not empty", "too many levels of symbolic links", "", "no message of desired type", "identifier removed", "channel number out of range", "level 2 not synchronized", "level 3 halted", "level 3 reset", "link number out of range", "protocol driver not attached", "no CSI structure available", "level 2 halted", "invalid exchange", "invalid request descriptor", "exchange full", "no anode", "invalid request code", "invalid slot", "", "bad font file format", "device not a stream", "no data available", "timer expired", "out of streams resources", "machine is not on the network", "package not installed", "object is remote", "link has been severed", "advertise error", "srmount error", "communication error on send", "protocol error", "multihop attempted", "RFS specific error", "bad message", "value too large for defined data type", "name not unique on network", "file descriptor in bad state", "remote address changed", "can not access a needed shared library", "accessing a corrupted shared library", ".lib section in a.out corrupted", "attempting to link in too many shared libraries", "cannot exec a shared library directly", "invalid or incomplete multibyte or wide character", "interrupted system call should be restarted", "streams pipe error", "too many users", "socket operation on non-socket", "destination address required", "message too long", "protocol wrong type for socket", "protocol not available", "protocol not supported", "socket type not supported", "operation not supported", "protocol family not supported", "address family not supported by protocol", "address already in use", "cannot assign requested address", "network is down", "network is unreachable", "network dropped connection on reset", "software caused connection abort", "connection reset by peer", "no buffer space available", "transport endpoint is already connected", "transport endpoint is not connected", "cannot send after transport endpoint shutdown", "too many references: cannot splice", "connection timed out", "connection refused", "host is down", "no route to host", "operation already in progress", "operation now in progress", "stale NFS file handle", "structure needs cleaning", "not a XENIX named type file", "no XENIX semaphores available", "is a named type file", "remote I/O error", "disk quota exceeded", "no medium found", "wrong medium type", "operation canceled", "required key not available", "key has expired", "key has been revoked", "key was rejected by service", "owner died", "state not recoverable", "operation not possible due to RF-kill"]);
		mapper = new mmapper.Ptr(new sync.Mutex.Ptr(), new $Map(), mmap, munmap);
		init();
	};
	return $pkg;
})();
$packages["strings"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], errors = $packages["errors"], io = $packages["io"], utf8 = $packages["unicode/utf8"], unicode = $packages["unicode"], IndexByte, explode, hashstr, Count, LastIndex, genSplit, SplitN, Join, TrimLeftFunc, TrimRightFunc, TrimFunc, indexFunc, lastIndexFunc, TrimSpace;
	IndexByte = $pkg.IndexByte = function(s, c) {
		return $parseInt(s.indexOf($global.String.fromCharCode(c))) >> 0;
	};
	explode = function(s, n) {
		var l, a, size, ch, _tmp, _tmp$1, i, cur, _tuple;
		if (n === 0) {
			return ($sliceType($String)).nil;
		}
		l = utf8.RuneCountInString(s);
		if (n <= 0 || n > l) {
			n = l;
		}
		a = ($sliceType($String)).make(n);
		size = 0;
		ch = 0;
		_tmp = 0; _tmp$1 = 0; i = _tmp; cur = _tmp$1;
		while ((i + 1 >> 0) < n) {
			_tuple = utf8.DecodeRuneInString(s.substring(cur)); ch = _tuple[0]; size = _tuple[1];
			if (ch === 65533) {
				(i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i] = "\xEF\xBF\xBD";
			} else {
				(i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i] = s.substring(cur, (cur + size >> 0));
			}
			cur = cur + (size) >> 0;
			i = i + (1) >> 0;
		}
		if (cur < s.length) {
			(i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i] = s.substring(cur);
		}
		return a;
	};
	hashstr = function(sep) {
		var hash, i, _tmp, _tmp$1, pow, sq, i$1, x, x$1;
		hash = 0;
		i = 0;
		while (i < sep.length) {
			hash = ((((hash >>> 16 << 16) * 16777619 >>> 0) + (hash << 16 >>> 16) * 16777619) >>> 0) + (sep.charCodeAt(i) >>> 0) >>> 0;
			i = i + (1) >> 0;
		}
		_tmp = 1; _tmp$1 = 16777619; pow = _tmp; sq = _tmp$1;
		i$1 = sep.length;
		while (i$1 > 0) {
			if (!(((i$1 & 1) === 0))) {
				pow = (x = sq, (((pow >>> 16 << 16) * x >>> 0) + (pow << 16 >>> 16) * x) >>> 0);
			}
			sq = (x$1 = sq, (((sq >>> 16 << 16) * x$1 >>> 0) + (sq << 16 >>> 16) * x$1) >>> 0);
			i$1 = (i$1 >> $min((1), 31)) >> 0;
		}
		return [hash, pow];
	};
	Count = $pkg.Count = function(s, sep) {
		var n, c, i, _tuple, hashsep, pow, h, i$1, lastmatch, i$2, x, x$1;
		n = 0;
		if (sep.length === 0) {
			return utf8.RuneCountInString(s) + 1 >> 0;
		} else if (sep.length === 1) {
			c = sep.charCodeAt(0);
			i = 0;
			while (i < s.length) {
				if (s.charCodeAt(i) === c) {
					n = n + (1) >> 0;
				}
				i = i + (1) >> 0;
			}
			return n;
		} else if (sep.length > s.length) {
			return 0;
		} else if (sep.length === s.length) {
			if (sep === s) {
				return 1;
			}
			return 0;
		}
		_tuple = hashstr(sep); hashsep = _tuple[0]; pow = _tuple[1];
		h = 0;
		i$1 = 0;
		while (i$1 < sep.length) {
			h = ((((h >>> 16 << 16) * 16777619 >>> 0) + (h << 16 >>> 16) * 16777619) >>> 0) + (s.charCodeAt(i$1) >>> 0) >>> 0;
			i$1 = i$1 + (1) >> 0;
		}
		lastmatch = 0;
		if ((h === hashsep) && s.substring(0, sep.length) === sep) {
			n = n + (1) >> 0;
			lastmatch = sep.length;
		}
		i$2 = sep.length;
		while (i$2 < s.length) {
			h = (x = 16777619, (((h >>> 16 << 16) * x >>> 0) + (h << 16 >>> 16) * x) >>> 0);
			h = h + ((s.charCodeAt(i$2) >>> 0)) >>> 0;
			h = h - ((x$1 = (s.charCodeAt((i$2 - sep.length >> 0)) >>> 0), (((pow >>> 16 << 16) * x$1 >>> 0) + (pow << 16 >>> 16) * x$1) >>> 0)) >>> 0;
			i$2 = i$2 + (1) >> 0;
			if ((h === hashsep) && lastmatch <= (i$2 - sep.length >> 0) && s.substring((i$2 - sep.length >> 0), i$2) === sep) {
				n = n + (1) >> 0;
				lastmatch = i$2;
			}
		}
		return n;
	};
	LastIndex = $pkg.LastIndex = function(s, sep) {
		var n, c, i, i$1;
		n = sep.length;
		if (n === 0) {
			return s.length;
		}
		c = sep.charCodeAt(0);
		if (n === 1) {
			i = s.length - 1 >> 0;
			while (i >= 0) {
				if (s.charCodeAt(i) === c) {
					return i;
				}
				i = i - (1) >> 0;
			}
			return -1;
		}
		i$1 = s.length - n >> 0;
		while (i$1 >= 0) {
			if ((s.charCodeAt(i$1) === c) && s.substring(i$1, (i$1 + n >> 0)) === sep) {
				return i$1;
			}
			i$1 = i$1 - (1) >> 0;
		}
		return -1;
	};
	genSplit = function(s, sep, sepSave, n) {
		var c, start, a, na, i;
		if (n === 0) {
			return ($sliceType($String)).nil;
		}
		if (sep === "") {
			return explode(s, n);
		}
		if (n < 0) {
			n = Count(s, sep) + 1 >> 0;
		}
		c = sep.charCodeAt(0);
		start = 0;
		a = ($sliceType($String)).make(n);
		na = 0;
		i = 0;
		while ((i + sep.length >> 0) <= s.length && (na + 1 >> 0) < n) {
			if ((s.charCodeAt(i) === c) && ((sep.length === 1) || s.substring(i, (i + sep.length >> 0)) === sep)) {
				(na < 0 || na >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + na] = s.substring(start, (i + sepSave >> 0));
				na = na + (1) >> 0;
				start = i + sep.length >> 0;
				i = i + ((sep.length - 1 >> 0)) >> 0;
			}
			i = i + (1) >> 0;
		}
		(na < 0 || na >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + na] = s.substring(start);
		return $subslice(a, 0, (na + 1 >> 0));
	};
	SplitN = $pkg.SplitN = function(s, sep, n) {
		return genSplit(s, sep, 0, n);
	};
	Join = $pkg.Join = function(a, sep) {
		var x, x$1, n, i, b, bp, _ref, _i, s;
		if (a.$length === 0) {
			return "";
		}
		if (a.$length === 1) {
			return ((0 < 0 || 0 >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + 0]);
		}
		n = (x = sep.length, x$1 = (a.$length - 1 >> 0), (((x >>> 16 << 16) * x$1 >> 0) + (x << 16 >>> 16) * x$1) >> 0);
		i = 0;
		while (i < a.$length) {
			n = n + (((i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i]).length) >> 0;
			i = i + (1) >> 0;
		}
		b = ($sliceType($Uint8)).make(n);
		bp = $copyString(b, ((0 < 0 || 0 >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + 0]));
		_ref = $subslice(a, 1);
		_i = 0;
		while (_i < _ref.$length) {
			s = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			bp = bp + ($copyString($subslice(b, bp), sep)) >> 0;
			bp = bp + ($copyString($subslice(b, bp), s)) >> 0;
			_i++;
		}
		return $bytesToString(b);
	};
	TrimLeftFunc = $pkg.TrimLeftFunc = function(s, f) {
		var i;
		i = indexFunc(s, f, false);
		if (i === -1) {
			return "";
		}
		return s.substring(i);
	};
	TrimRightFunc = $pkg.TrimRightFunc = function(s, f) {
		var i, _tuple, wid;
		i = lastIndexFunc(s, f, false);
		if (i >= 0 && s.charCodeAt(i) >= 128) {
			_tuple = utf8.DecodeRuneInString(s.substring(i)); wid = _tuple[1];
			i = i + (wid) >> 0;
		} else {
			i = i + (1) >> 0;
		}
		return s.substring(0, i);
	};
	TrimFunc = $pkg.TrimFunc = function(s, f) {
		return TrimRightFunc(TrimLeftFunc(s, f), f);
	};
	indexFunc = function(s, f, truth) {
		var start, wid, r, _tuple;
		start = 0;
		while (start < s.length) {
			wid = 1;
			r = (s.charCodeAt(start) >> 0);
			if (r >= 128) {
				_tuple = utf8.DecodeRuneInString(s.substring(start)); r = _tuple[0]; wid = _tuple[1];
			}
			if (f(r) === truth) {
				return start;
			}
			start = start + (wid) >> 0;
		}
		return -1;
	};
	lastIndexFunc = function(s, f, truth) {
		var i, _tuple, r, size;
		i = s.length;
		while (i > 0) {
			_tuple = utf8.DecodeLastRuneInString(s.substring(0, i)); r = _tuple[0]; size = _tuple[1];
			i = i - (size) >> 0;
			if (f(r) === truth) {
				return i;
			}
		}
		return -1;
	};
	TrimSpace = $pkg.TrimSpace = function(s) {
		return TrimFunc(s, unicode.IsSpace);
	};
	$pkg.$init = function() {
	};
	return $pkg;
})();
$packages["time"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], strings = $packages["strings"], errors = $packages["errors"], syscall = $packages["syscall"], sync = $packages["sync"], runtime = $packages["runtime"], ParseError, Time, Month, Weekday, Duration, Location, zone, zoneTrans, std0x, longDayNames, shortDayNames, shortMonthNames, longMonthNames, atoiError, errBad, errLeadingInt, months, days, daysBefore, utcLoc, localLoc, localOnce, zoneinfo, badData, zoneDirs, _tuple, initLocal, startsWithLowerCase, nextStdChunk, match, lookup, appendUint, atoi, formatNano, quote, isDigit, getnum, cutspace, skip, Parse, parse, parseTimeZone, parseGMT, parseNanoseconds, leadingInt, absWeekday, absClock, fmtFrac, fmtInt, absDate, Unix, isLeap, norm, Date, div, FixedZone;
	ParseError = $pkg.ParseError = $newType(0, "Struct", "time.ParseError", "ParseError", "time", function(Layout_, Value_, LayoutElem_, ValueElem_, Message_) {
		this.$val = this;
		this.Layout = Layout_ !== undefined ? Layout_ : "";
		this.Value = Value_ !== undefined ? Value_ : "";
		this.LayoutElem = LayoutElem_ !== undefined ? LayoutElem_ : "";
		this.ValueElem = ValueElem_ !== undefined ? ValueElem_ : "";
		this.Message = Message_ !== undefined ? Message_ : "";
	});
	Time = $pkg.Time = $newType(0, "Struct", "time.Time", "Time", "time", function(sec_, nsec_, loc_) {
		this.$val = this;
		this.sec = sec_ !== undefined ? sec_ : new $Int64(0, 0);
		this.nsec = nsec_ !== undefined ? nsec_ : 0;
		this.loc = loc_ !== undefined ? loc_ : ($ptrType(Location)).nil;
	});
	Month = $pkg.Month = $newType(4, "Int", "time.Month", "Month", "time", null);
	Weekday = $pkg.Weekday = $newType(4, "Int", "time.Weekday", "Weekday", "time", null);
	Duration = $pkg.Duration = $newType(8, "Int64", "time.Duration", "Duration", "time", null);
	Location = $pkg.Location = $newType(0, "Struct", "time.Location", "Location", "time", function(name_, zone_, tx_, cacheStart_, cacheEnd_, cacheZone_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : "";
		this.zone = zone_ !== undefined ? zone_ : ($sliceType(zone)).nil;
		this.tx = tx_ !== undefined ? tx_ : ($sliceType(zoneTrans)).nil;
		this.cacheStart = cacheStart_ !== undefined ? cacheStart_ : new $Int64(0, 0);
		this.cacheEnd = cacheEnd_ !== undefined ? cacheEnd_ : new $Int64(0, 0);
		this.cacheZone = cacheZone_ !== undefined ? cacheZone_ : ($ptrType(zone)).nil;
	});
	zone = $pkg.zone = $newType(0, "Struct", "time.zone", "zone", "time", function(name_, offset_, isDST_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : "";
		this.offset = offset_ !== undefined ? offset_ : 0;
		this.isDST = isDST_ !== undefined ? isDST_ : false;
	});
	zoneTrans = $pkg.zoneTrans = $newType(0, "Struct", "time.zoneTrans", "zoneTrans", "time", function(when_, index_, isstd_, isutc_) {
		this.$val = this;
		this.when = when_ !== undefined ? when_ : new $Int64(0, 0);
		this.index = index_ !== undefined ? index_ : 0;
		this.isstd = isstd_ !== undefined ? isstd_ : false;
		this.isutc = isutc_ !== undefined ? isutc_ : false;
	});
	initLocal = function() {
		var d, s, i, j, x;
		d = new ($global.Date)();
		s = $internalize(d, $String);
		i = strings.IndexByte(s, 40);
		j = strings.IndexByte(s, 41);
		if ((i === -1) || (j === -1)) {
			localLoc.name = "UTC";
			return;
		}
		localLoc.name = s.substring((i + 1 >> 0), j);
		localLoc.zone = new ($sliceType(zone))([new zone.Ptr(localLoc.name, (x = $parseInt(d.getTimezoneOffset()) >> 0, (((x >>> 16 << 16) * -60 >> 0) + (x << 16 >>> 16) * -60) >> 0), false)]);
	};
	startsWithLowerCase = function(str) {
		var c;
		if (str.length === 0) {
			return false;
		}
		c = str.charCodeAt(0);
		return 97 <= c && c <= 122;
	};
	nextStdChunk = function(layout) {
		var prefix = "", std = 0, suffix = "", i, c, _ref, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, x, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, _tmp$23, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46, _tmp$47, _tmp$48, _tmp$49, _tmp$50, _tmp$51, _tmp$52, _tmp$53, _tmp$54, _tmp$55, _tmp$56, _tmp$57, _tmp$58, _tmp$59, _tmp$60, _tmp$61, _tmp$62, _tmp$63, _tmp$64, _tmp$65, _tmp$66, _tmp$67, _tmp$68, _tmp$69, _tmp$70, _tmp$71, _tmp$72, _tmp$73, _tmp$74, ch, j, std$1, _tmp$75, _tmp$76, _tmp$77, _tmp$78, _tmp$79, _tmp$80;
		i = 0;
		while (i < layout.length) {
			c = (layout.charCodeAt(i) >> 0);
			_ref = c;
			if (_ref === 74) {
				if (layout.length >= (i + 3 >> 0) && layout.substring(i, (i + 3 >> 0)) === "Jan") {
					if (layout.length >= (i + 7 >> 0) && layout.substring(i, (i + 7 >> 0)) === "January") {
						_tmp = layout.substring(0, i); _tmp$1 = 257; _tmp$2 = layout.substring((i + 7 >> 0)); prefix = _tmp; std = _tmp$1; suffix = _tmp$2;
						return [prefix, std, suffix];
					}
					if (!startsWithLowerCase(layout.substring((i + 3 >> 0)))) {
						_tmp$3 = layout.substring(0, i); _tmp$4 = 258; _tmp$5 = layout.substring((i + 3 >> 0)); prefix = _tmp$3; std = _tmp$4; suffix = _tmp$5;
						return [prefix, std, suffix];
					}
				}
			} else if (_ref === 77) {
				if (layout.length >= (i + 3 >> 0)) {
					if (layout.substring(i, (i + 3 >> 0)) === "Mon") {
						if (layout.length >= (i + 6 >> 0) && layout.substring(i, (i + 6 >> 0)) === "Monday") {
							_tmp$6 = layout.substring(0, i); _tmp$7 = 261; _tmp$8 = layout.substring((i + 6 >> 0)); prefix = _tmp$6; std = _tmp$7; suffix = _tmp$8;
							return [prefix, std, suffix];
						}
						if (!startsWithLowerCase(layout.substring((i + 3 >> 0)))) {
							_tmp$9 = layout.substring(0, i); _tmp$10 = 262; _tmp$11 = layout.substring((i + 3 >> 0)); prefix = _tmp$9; std = _tmp$10; suffix = _tmp$11;
							return [prefix, std, suffix];
						}
					}
					if (layout.substring(i, (i + 3 >> 0)) === "MST") {
						_tmp$12 = layout.substring(0, i); _tmp$13 = 21; _tmp$14 = layout.substring((i + 3 >> 0)); prefix = _tmp$12; std = _tmp$13; suffix = _tmp$14;
						return [prefix, std, suffix];
					}
				}
			} else if (_ref === 48) {
				if (layout.length >= (i + 2 >> 0) && 49 <= layout.charCodeAt((i + 1 >> 0)) && layout.charCodeAt((i + 1 >> 0)) <= 54) {
					_tmp$15 = layout.substring(0, i); _tmp$16 = (x = layout.charCodeAt((i + 1 >> 0)) - 49 << 24 >>> 24, ((x < 0 || x >= std0x.length) ? $throwRuntimeError("index out of range") : std0x[x])); _tmp$17 = layout.substring((i + 2 >> 0)); prefix = _tmp$15; std = _tmp$16; suffix = _tmp$17;
					return [prefix, std, suffix];
				}
			} else if (_ref === 49) {
				if (layout.length >= (i + 2 >> 0) && (layout.charCodeAt((i + 1 >> 0)) === 53)) {
					_tmp$18 = layout.substring(0, i); _tmp$19 = 522; _tmp$20 = layout.substring((i + 2 >> 0)); prefix = _tmp$18; std = _tmp$19; suffix = _tmp$20;
					return [prefix, std, suffix];
				}
				_tmp$21 = layout.substring(0, i); _tmp$22 = 259; _tmp$23 = layout.substring((i + 1 >> 0)); prefix = _tmp$21; std = _tmp$22; suffix = _tmp$23;
				return [prefix, std, suffix];
			} else if (_ref === 50) {
				if (layout.length >= (i + 4 >> 0) && layout.substring(i, (i + 4 >> 0)) === "2006") {
					_tmp$24 = layout.substring(0, i); _tmp$25 = 273; _tmp$26 = layout.substring((i + 4 >> 0)); prefix = _tmp$24; std = _tmp$25; suffix = _tmp$26;
					return [prefix, std, suffix];
				}
				_tmp$27 = layout.substring(0, i); _tmp$28 = 263; _tmp$29 = layout.substring((i + 1 >> 0)); prefix = _tmp$27; std = _tmp$28; suffix = _tmp$29;
				return [prefix, std, suffix];
			} else if (_ref === 95) {
				if (layout.length >= (i + 2 >> 0) && (layout.charCodeAt((i + 1 >> 0)) === 50)) {
					_tmp$30 = layout.substring(0, i); _tmp$31 = 264; _tmp$32 = layout.substring((i + 2 >> 0)); prefix = _tmp$30; std = _tmp$31; suffix = _tmp$32;
					return [prefix, std, suffix];
				}
			} else if (_ref === 51) {
				_tmp$33 = layout.substring(0, i); _tmp$34 = 523; _tmp$35 = layout.substring((i + 1 >> 0)); prefix = _tmp$33; std = _tmp$34; suffix = _tmp$35;
				return [prefix, std, suffix];
			} else if (_ref === 52) {
				_tmp$36 = layout.substring(0, i); _tmp$37 = 525; _tmp$38 = layout.substring((i + 1 >> 0)); prefix = _tmp$36; std = _tmp$37; suffix = _tmp$38;
				return [prefix, std, suffix];
			} else if (_ref === 53) {
				_tmp$39 = layout.substring(0, i); _tmp$40 = 527; _tmp$41 = layout.substring((i + 1 >> 0)); prefix = _tmp$39; std = _tmp$40; suffix = _tmp$41;
				return [prefix, std, suffix];
			} else if (_ref === 80) {
				if (layout.length >= (i + 2 >> 0) && (layout.charCodeAt((i + 1 >> 0)) === 77)) {
					_tmp$42 = layout.substring(0, i); _tmp$43 = 531; _tmp$44 = layout.substring((i + 2 >> 0)); prefix = _tmp$42; std = _tmp$43; suffix = _tmp$44;
					return [prefix, std, suffix];
				}
			} else if (_ref === 112) {
				if (layout.length >= (i + 2 >> 0) && (layout.charCodeAt((i + 1 >> 0)) === 109)) {
					_tmp$45 = layout.substring(0, i); _tmp$46 = 532; _tmp$47 = layout.substring((i + 2 >> 0)); prefix = _tmp$45; std = _tmp$46; suffix = _tmp$47;
					return [prefix, std, suffix];
				}
			} else if (_ref === 45) {
				if (layout.length >= (i + 7 >> 0) && layout.substring(i, (i + 7 >> 0)) === "-070000") {
					_tmp$48 = layout.substring(0, i); _tmp$49 = 27; _tmp$50 = layout.substring((i + 7 >> 0)); prefix = _tmp$48; std = _tmp$49; suffix = _tmp$50;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 9 >> 0) && layout.substring(i, (i + 9 >> 0)) === "-07:00:00") {
					_tmp$51 = layout.substring(0, i); _tmp$52 = 30; _tmp$53 = layout.substring((i + 9 >> 0)); prefix = _tmp$51; std = _tmp$52; suffix = _tmp$53;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 5 >> 0) && layout.substring(i, (i + 5 >> 0)) === "-0700") {
					_tmp$54 = layout.substring(0, i); _tmp$55 = 26; _tmp$56 = layout.substring((i + 5 >> 0)); prefix = _tmp$54; std = _tmp$55; suffix = _tmp$56;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 6 >> 0) && layout.substring(i, (i + 6 >> 0)) === "-07:00") {
					_tmp$57 = layout.substring(0, i); _tmp$58 = 29; _tmp$59 = layout.substring((i + 6 >> 0)); prefix = _tmp$57; std = _tmp$58; suffix = _tmp$59;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 3 >> 0) && layout.substring(i, (i + 3 >> 0)) === "-07") {
					_tmp$60 = layout.substring(0, i); _tmp$61 = 28; _tmp$62 = layout.substring((i + 3 >> 0)); prefix = _tmp$60; std = _tmp$61; suffix = _tmp$62;
					return [prefix, std, suffix];
				}
			} else if (_ref === 90) {
				if (layout.length >= (i + 7 >> 0) && layout.substring(i, (i + 7 >> 0)) === "Z070000") {
					_tmp$63 = layout.substring(0, i); _tmp$64 = 23; _tmp$65 = layout.substring((i + 7 >> 0)); prefix = _tmp$63; std = _tmp$64; suffix = _tmp$65;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 9 >> 0) && layout.substring(i, (i + 9 >> 0)) === "Z07:00:00") {
					_tmp$66 = layout.substring(0, i); _tmp$67 = 25; _tmp$68 = layout.substring((i + 9 >> 0)); prefix = _tmp$66; std = _tmp$67; suffix = _tmp$68;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 5 >> 0) && layout.substring(i, (i + 5 >> 0)) === "Z0700") {
					_tmp$69 = layout.substring(0, i); _tmp$70 = 22; _tmp$71 = layout.substring((i + 5 >> 0)); prefix = _tmp$69; std = _tmp$70; suffix = _tmp$71;
					return [prefix, std, suffix];
				}
				if (layout.length >= (i + 6 >> 0) && layout.substring(i, (i + 6 >> 0)) === "Z07:00") {
					_tmp$72 = layout.substring(0, i); _tmp$73 = 24; _tmp$74 = layout.substring((i + 6 >> 0)); prefix = _tmp$72; std = _tmp$73; suffix = _tmp$74;
					return [prefix, std, suffix];
				}
			} else if (_ref === 46) {
				if ((i + 1 >> 0) < layout.length && ((layout.charCodeAt((i + 1 >> 0)) === 48) || (layout.charCodeAt((i + 1 >> 0)) === 57))) {
					ch = layout.charCodeAt((i + 1 >> 0));
					j = i + 1 >> 0;
					while (j < layout.length && (layout.charCodeAt(j) === ch)) {
						j = j + (1) >> 0;
					}
					if (!isDigit(layout, j)) {
						std$1 = 31;
						if (layout.charCodeAt((i + 1 >> 0)) === 57) {
							std$1 = 32;
						}
						std$1 = std$1 | ((((j - ((i + 1 >> 0)) >> 0)) << 16 >> 0));
						_tmp$75 = layout.substring(0, i); _tmp$76 = std$1; _tmp$77 = layout.substring(j); prefix = _tmp$75; std = _tmp$76; suffix = _tmp$77;
						return [prefix, std, suffix];
					}
				}
			}
			i = i + (1) >> 0;
		}
		_tmp$78 = layout; _tmp$79 = 0; _tmp$80 = ""; prefix = _tmp$78; std = _tmp$79; suffix = _tmp$80;
		return [prefix, std, suffix];
	};
	match = function(s1, s2) {
		var i, c1, c2;
		i = 0;
		while (i < s1.length) {
			c1 = s1.charCodeAt(i);
			c2 = s2.charCodeAt(i);
			if (!((c1 === c2))) {
				c1 = (c1 | (32)) >>> 0;
				c2 = (c2 | (32)) >>> 0;
				if (!((c1 === c2)) || c1 < 97 || c1 > 122) {
					return false;
				}
			}
			i = i + (1) >> 0;
		}
		return true;
	};
	lookup = function(tab, val) {
		var _ref, _i, i, v;
		_ref = tab;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			v = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (val.length >= v.length && match(val.substring(0, v.length), v)) {
				return [i, val.substring(v.length), null];
			}
			_i++;
		}
		return [-1, val, errBad];
	};
	appendUint = function(b, x, pad) {
		var _q, _r, buf, n, _r$1, _q$1;
		if (x < 10) {
			if (!((pad === 0))) {
				b = $append(b, pad);
			}
			return $append(b, ((48 + x >>> 0) << 24 >>> 24));
		}
		if (x < 100) {
			b = $append(b, ((48 + (_q = x / 10, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >>> 0 : $throwRuntimeError("integer divide by zero")) >>> 0) << 24 >>> 24));
			b = $append(b, ((48 + (_r = x % 10, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) >>> 0) << 24 >>> 24));
			return b;
		}
		buf = ($arrayType($Uint8, 32)).zero(); $copy(buf, ($arrayType($Uint8, 32)).zero(), ($arrayType($Uint8, 32)));
		n = 32;
		if (x === 0) {
			return $append(b, 48);
		}
		while (x >= 10) {
			n = n - (1) >> 0;
			(n < 0 || n >= buf.length) ? $throwRuntimeError("index out of range") : buf[n] = (((_r$1 = x % 10, _r$1 === _r$1 ? _r$1 : $throwRuntimeError("integer divide by zero")) + 48 >>> 0) << 24 >>> 24);
			x = (_q$1 = x / (10), (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >>> 0 : $throwRuntimeError("integer divide by zero"));
		}
		n = n - (1) >> 0;
		(n < 0 || n >= buf.length) ? $throwRuntimeError("index out of range") : buf[n] = ((x + 48 >>> 0) << 24 >>> 24);
		return $appendSlice(b, $subslice(new ($sliceType($Uint8))(buf), n));
	};
	atoi = function(s) {
		var x = 0, err = null, neg, _tuple$1, q, rem, _tmp, _tmp$1, _tmp$2, _tmp$3;
		neg = false;
		if (!(s === "") && ((s.charCodeAt(0) === 45) || (s.charCodeAt(0) === 43))) {
			neg = s.charCodeAt(0) === 45;
			s = s.substring(1);
		}
		_tuple$1 = leadingInt(s); q = _tuple$1[0]; rem = _tuple$1[1]; err = _tuple$1[2];
		x = ((q.$low + ((q.$high >> 31) * 4294967296)) >> 0);
		if (!($interfaceIsEqual(err, null)) || !(rem === "")) {
			_tmp = 0; _tmp$1 = atoiError; x = _tmp; err = _tmp$1;
			return [x, err];
		}
		if (neg) {
			x = -x;
		}
		_tmp$2 = x; _tmp$3 = null; x = _tmp$2; err = _tmp$3;
		return [x, err];
	};
	formatNano = function(b, nanosec, n, trim) {
		var u, buf, start, _r, _q, x;
		u = nanosec;
		buf = ($arrayType($Uint8, 9)).zero(); $copy(buf, ($arrayType($Uint8, 9)).zero(), ($arrayType($Uint8, 9)));
		start = 9;
		while (start > 0) {
			start = start - (1) >> 0;
			(start < 0 || start >= buf.length) ? $throwRuntimeError("index out of range") : buf[start] = (((_r = u % 10, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) + 48 >>> 0) << 24 >>> 24);
			u = (_q = u / (10), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >>> 0 : $throwRuntimeError("integer divide by zero"));
		}
		if (n > 9) {
			n = 9;
		}
		if (trim) {
			while (n > 0 && ((x = n - 1 >> 0, ((x < 0 || x >= buf.length) ? $throwRuntimeError("index out of range") : buf[x])) === 48)) {
				n = n - (1) >> 0;
			}
			if (n === 0) {
				return b;
			}
		}
		b = $append(b, 46);
		return $appendSlice(b, $subslice(new ($sliceType($Uint8))(buf), 0, n));
	};
	Time.Ptr.prototype.String = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		return t.Format("2006-01-02 15:04:05.999999999 -0700 MST");
	};
	Time.prototype.String = function() { return this.$val.String(); };
	Time.Ptr.prototype.Format = function(layout) {
		var t, _tuple$1, name, offset, abs, year, month, day, hour, min, sec, b, buf, max, _tuple$2, prefix, std, suffix, _tuple$3, _tuple$4, _ref, y, _r, y$1, m, s, _r$1, hr, _r$2, hr$1, _q, zone$1, absoffset, _q$1, _r$3, _r$4, _q$2, zone$2, _q$3, _r$5;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.locabs(); name = _tuple$1[0]; offset = _tuple$1[1]; abs = _tuple$1[2];
		year = -1;
		month = 0;
		day = 0;
		hour = -1;
		min = 0;
		sec = 0;
		b = ($sliceType($Uint8)).nil;
		buf = ($arrayType($Uint8, 64)).zero(); $copy(buf, ($arrayType($Uint8, 64)).zero(), ($arrayType($Uint8, 64)));
		max = layout.length + 10 >> 0;
		if (max <= 64) {
			b = $subslice(new ($sliceType($Uint8))(buf), 0, 0);
		} else {
			b = ($sliceType($Uint8)).make(0, max);
		}
		while (!(layout === "")) {
			_tuple$2 = nextStdChunk(layout); prefix = _tuple$2[0]; std = _tuple$2[1]; suffix = _tuple$2[2];
			if (!(prefix === "")) {
				b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes(prefix)));
			}
			if (std === 0) {
				break;
			}
			layout = suffix;
			if (year < 0 && !(((std & 256) === 0))) {
				_tuple$3 = absDate(abs, true); year = _tuple$3[0]; month = _tuple$3[1]; day = _tuple$3[2];
			}
			if (hour < 0 && !(((std & 512) === 0))) {
				_tuple$4 = absClock(abs); hour = _tuple$4[0]; min = _tuple$4[1]; sec = _tuple$4[2];
			}
			_ref = std & 65535;
			switch (0) { default: if (_ref === 274) {
				y = year;
				if (y < 0) {
					y = -y;
				}
				b = appendUint(b, ((_r = y % 100, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) >>> 0), 48);
			} else if (_ref === 273) {
				y$1 = year;
				if (year <= -1000) {
					b = $append(b, 45);
					y$1 = -y$1;
				} else if (year <= -100) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("-0")));
					y$1 = -y$1;
				} else if (year <= -10) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("-00")));
					y$1 = -y$1;
				} else if (year < 0) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("-000")));
					y$1 = -y$1;
				} else if (year < 10) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("000")));
				} else if (year < 100) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("00")));
				} else if (year < 1000) {
					b = $append(b, 48);
				}
				b = appendUint(b, (y$1 >>> 0), 0);
			} else if (_ref === 258) {
				b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes((new Month(month)).String().substring(0, 3))));
			} else if (_ref === 257) {
				m = (new Month(month)).String();
				b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes(m)));
			} else if (_ref === 259) {
				b = appendUint(b, (month >>> 0), 0);
			} else if (_ref === 260) {
				b = appendUint(b, (month >>> 0), 48);
			} else if (_ref === 262) {
				b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes((new Weekday(absWeekday(abs))).String().substring(0, 3))));
			} else if (_ref === 261) {
				s = (new Weekday(absWeekday(abs))).String();
				b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes(s)));
			} else if (_ref === 263) {
				b = appendUint(b, (day >>> 0), 0);
			} else if (_ref === 264) {
				b = appendUint(b, (day >>> 0), 32);
			} else if (_ref === 265) {
				b = appendUint(b, (day >>> 0), 48);
			} else if (_ref === 522) {
				b = appendUint(b, (hour >>> 0), 48);
			} else if (_ref === 523) {
				hr = (_r$1 = hour % 12, _r$1 === _r$1 ? _r$1 : $throwRuntimeError("integer divide by zero"));
				if (hr === 0) {
					hr = 12;
				}
				b = appendUint(b, (hr >>> 0), 0);
			} else if (_ref === 524) {
				hr$1 = (_r$2 = hour % 12, _r$2 === _r$2 ? _r$2 : $throwRuntimeError("integer divide by zero"));
				if (hr$1 === 0) {
					hr$1 = 12;
				}
				b = appendUint(b, (hr$1 >>> 0), 48);
			} else if (_ref === 525) {
				b = appendUint(b, (min >>> 0), 0);
			} else if (_ref === 526) {
				b = appendUint(b, (min >>> 0), 48);
			} else if (_ref === 527) {
				b = appendUint(b, (sec >>> 0), 0);
			} else if (_ref === 528) {
				b = appendUint(b, (sec >>> 0), 48);
			} else if (_ref === 531) {
				if (hour >= 12) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("PM")));
				} else {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("AM")));
				}
			} else if (_ref === 532) {
				if (hour >= 12) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("pm")));
				} else {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes("am")));
				}
			} else if (_ref === 22 || _ref === 24 || _ref === 23 || _ref === 25 || _ref === 26 || _ref === 29 || _ref === 27 || _ref === 30) {
				if ((offset === 0) && ((std === 22) || (std === 24) || (std === 23) || (std === 25))) {
					b = $append(b, 90);
					break;
				}
				zone$1 = (_q = offset / 60, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
				absoffset = offset;
				if (zone$1 < 0) {
					b = $append(b, 45);
					zone$1 = -zone$1;
					absoffset = -absoffset;
				} else {
					b = $append(b, 43);
				}
				b = appendUint(b, ((_q$1 = zone$1 / 60, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >> 0 : $throwRuntimeError("integer divide by zero")) >>> 0), 48);
				if ((std === 24) || (std === 29)) {
					b = $append(b, 58);
				}
				b = appendUint(b, ((_r$3 = zone$1 % 60, _r$3 === _r$3 ? _r$3 : $throwRuntimeError("integer divide by zero")) >>> 0), 48);
				if ((std === 23) || (std === 27) || (std === 30) || (std === 25)) {
					if ((std === 30) || (std === 25)) {
						b = $append(b, 58);
					}
					b = appendUint(b, ((_r$4 = absoffset % 60, _r$4 === _r$4 ? _r$4 : $throwRuntimeError("integer divide by zero")) >>> 0), 48);
				}
			} else if (_ref === 21) {
				if (!(name === "")) {
					b = $appendSlice(b, new ($sliceType($Uint8))($stringToBytes(name)));
					break;
				}
				zone$2 = (_q$2 = offset / 60, (_q$2 === _q$2 && _q$2 !== 1/0 && _q$2 !== -1/0) ? _q$2 >> 0 : $throwRuntimeError("integer divide by zero"));
				if (zone$2 < 0) {
					b = $append(b, 45);
					zone$2 = -zone$2;
				} else {
					b = $append(b, 43);
				}
				b = appendUint(b, ((_q$3 = zone$2 / 60, (_q$3 === _q$3 && _q$3 !== 1/0 && _q$3 !== -1/0) ? _q$3 >> 0 : $throwRuntimeError("integer divide by zero")) >>> 0), 48);
				b = appendUint(b, ((_r$5 = zone$2 % 60, _r$5 === _r$5 ? _r$5 : $throwRuntimeError("integer divide by zero")) >>> 0), 48);
			} else if (_ref === 31 || _ref === 32) {
				b = formatNano(b, (t.Nanosecond() >>> 0), std >> 16 >> 0, (std & 65535) === 32);
			} }
		}
		return $bytesToString(b);
	};
	Time.prototype.Format = function(layout) { return this.$val.Format(layout); };
	quote = function(s) {
		return "\"" + s + "\"";
	};
	ParseError.Ptr.prototype.Error = function() {
		var e;
		e = this;
		if (e.Message === "") {
			return "parsing time " + quote(e.Value) + " as " + quote(e.Layout) + ": cannot parse " + quote(e.ValueElem) + " as " + quote(e.LayoutElem);
		}
		return "parsing time " + quote(e.Value) + e.Message;
	};
	ParseError.prototype.Error = function() { return this.$val.Error(); };
	isDigit = function(s, i) {
		var c;
		if (s.length <= i) {
			return false;
		}
		c = s.charCodeAt(i);
		return 48 <= c && c <= 57;
	};
	getnum = function(s, fixed) {
		var x;
		if (!isDigit(s, 0)) {
			return [0, s, errBad];
		}
		if (!isDigit(s, 1)) {
			if (fixed) {
				return [0, s, errBad];
			}
			return [((s.charCodeAt(0) - 48 << 24 >>> 24) >> 0), s.substring(1), null];
		}
		return [(x = ((s.charCodeAt(0) - 48 << 24 >>> 24) >> 0), (((x >>> 16 << 16) * 10 >> 0) + (x << 16 >>> 16) * 10) >> 0) + ((s.charCodeAt(1) - 48 << 24 >>> 24) >> 0) >> 0, s.substring(2), null];
	};
	cutspace = function(s) {
		while (s.length > 0 && (s.charCodeAt(0) === 32)) {
			s = s.substring(1);
		}
		return s;
	};
	skip = function(value, prefix) {
		while (prefix.length > 0) {
			if (prefix.charCodeAt(0) === 32) {
				if (value.length > 0 && !((value.charCodeAt(0) === 32))) {
					return [value, errBad];
				}
				prefix = cutspace(prefix);
				value = cutspace(value);
				continue;
			}
			if ((value.length === 0) || !((value.charCodeAt(0) === prefix.charCodeAt(0)))) {
				return [value, errBad];
			}
			prefix = prefix.substring(1);
			value = value.substring(1);
		}
		return [value, null];
	};
	Parse = $pkg.Parse = function(layout, value) {
		return parse(layout, value, $pkg.UTC, $pkg.Local);
	};
	parse = function(layout, value, defaultLocation, local) {
		var _tmp, _tmp$1, alayout, avalue, rangeErrString, amSet, pmSet, year, month, day, hour, min, sec, nsec, z, zoneOffset, zoneName, err, _tuple$1, prefix, std, suffix, stdstr, _tuple$2, p, _ref, _tmp$2, _tmp$3, _tuple$3, _tmp$4, _tmp$5, _tuple$4, _tuple$5, _tuple$6, _tuple$7, _tuple$8, _tuple$9, _tuple$10, _tuple$11, _tuple$12, _tuple$13, _tuple$14, _tuple$15, n, _tuple$16, _tmp$6, _tmp$7, _ref$1, _tmp$8, _tmp$9, _ref$2, _tmp$10, _tmp$11, _tmp$12, _tmp$13, sign, hour$1, min$1, seconds, _tmp$14, _tmp$15, _tmp$16, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, _tmp$23, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40, _tmp$41, hr, mm, ss, _tuple$17, _tuple$18, _tuple$19, x, _ref$3, _tuple$20, n$1, ok, _tmp$42, _tmp$43, ndigit, _tuple$21, i, _tuple$22, t, x$1, x$2, _tuple$23, x$3, name, offset, t$1, _tuple$24, x$4, offset$1, ok$1, x$5, x$6, _tuple$25, x$7;
		_tmp = layout; _tmp$1 = value; alayout = _tmp; avalue = _tmp$1;
		rangeErrString = "";
		amSet = false;
		pmSet = false;
		year = 0;
		month = 1;
		day = 1;
		hour = 0;
		min = 0;
		sec = 0;
		nsec = 0;
		z = ($ptrType(Location)).nil;
		zoneOffset = -1;
		zoneName = "";
		while (true) {
			err = null;
			_tuple$1 = nextStdChunk(layout); prefix = _tuple$1[0]; std = _tuple$1[1]; suffix = _tuple$1[2];
			stdstr = layout.substring(prefix.length, (layout.length - suffix.length >> 0));
			_tuple$2 = skip(value, prefix); value = _tuple$2[0]; err = _tuple$2[1];
			if (!($interfaceIsEqual(err, null))) {
				return [new Time.Ptr(new $Int64(0, 0), 0, ($ptrType(Location)).nil), new ParseError.Ptr(alayout, avalue, prefix, value, "")];
			}
			if (std === 0) {
				if (!((value.length === 0))) {
					return [new Time.Ptr(new $Int64(0, 0), 0, ($ptrType(Location)).nil), new ParseError.Ptr(alayout, avalue, "", value, ": extra text: " + value)];
				}
				break;
			}
			layout = suffix;
			p = "";
			_ref = std & 65535;
			switch (0) { default: if (_ref === 274) {
				if (value.length < 2) {
					err = errBad;
					break;
				}
				_tmp$2 = value.substring(0, 2); _tmp$3 = value.substring(2); p = _tmp$2; value = _tmp$3;
				_tuple$3 = atoi(p); year = _tuple$3[0]; err = _tuple$3[1];
				if (year >= 69) {
					year = year + (1900) >> 0;
				} else {
					year = year + (2000) >> 0;
				}
			} else if (_ref === 273) {
				if (value.length < 4 || !isDigit(value, 0)) {
					err = errBad;
					break;
				}
				_tmp$4 = value.substring(0, 4); _tmp$5 = value.substring(4); p = _tmp$4; value = _tmp$5;
				_tuple$4 = atoi(p); year = _tuple$4[0]; err = _tuple$4[1];
			} else if (_ref === 258) {
				_tuple$5 = lookup(shortMonthNames, value); month = _tuple$5[0]; value = _tuple$5[1]; err = _tuple$5[2];
			} else if (_ref === 257) {
				_tuple$6 = lookup(longMonthNames, value); month = _tuple$6[0]; value = _tuple$6[1]; err = _tuple$6[2];
			} else if (_ref === 259 || _ref === 260) {
				_tuple$7 = getnum(value, std === 260); month = _tuple$7[0]; value = _tuple$7[1]; err = _tuple$7[2];
				if (month <= 0 || 12 < month) {
					rangeErrString = "month";
				}
			} else if (_ref === 262) {
				_tuple$8 = lookup(shortDayNames, value); value = _tuple$8[1]; err = _tuple$8[2];
			} else if (_ref === 261) {
				_tuple$9 = lookup(longDayNames, value); value = _tuple$9[1]; err = _tuple$9[2];
			} else if (_ref === 263 || _ref === 264 || _ref === 265) {
				if ((std === 264) && value.length > 0 && (value.charCodeAt(0) === 32)) {
					value = value.substring(1);
				}
				_tuple$10 = getnum(value, std === 265); day = _tuple$10[0]; value = _tuple$10[1]; err = _tuple$10[2];
				if (day < 0 || 31 < day) {
					rangeErrString = "day";
				}
			} else if (_ref === 522) {
				_tuple$11 = getnum(value, false); hour = _tuple$11[0]; value = _tuple$11[1]; err = _tuple$11[2];
				if (hour < 0 || 24 <= hour) {
					rangeErrString = "hour";
				}
			} else if (_ref === 523 || _ref === 524) {
				_tuple$12 = getnum(value, std === 524); hour = _tuple$12[0]; value = _tuple$12[1]; err = _tuple$12[2];
				if (hour < 0 || 12 < hour) {
					rangeErrString = "hour";
				}
			} else if (_ref === 525 || _ref === 526) {
				_tuple$13 = getnum(value, std === 526); min = _tuple$13[0]; value = _tuple$13[1]; err = _tuple$13[2];
				if (min < 0 || 60 <= min) {
					rangeErrString = "minute";
				}
			} else if (_ref === 527 || _ref === 528) {
				_tuple$14 = getnum(value, std === 528); sec = _tuple$14[0]; value = _tuple$14[1]; err = _tuple$14[2];
				if (sec < 0 || 60 <= sec) {
					rangeErrString = "second";
				}
				if (value.length >= 2 && (value.charCodeAt(0) === 46) && isDigit(value, 1)) {
					_tuple$15 = nextStdChunk(layout); std = _tuple$15[1];
					std = std & (65535);
					if ((std === 31) || (std === 32)) {
						break;
					}
					n = 2;
					while (n < value.length && isDigit(value, n)) {
						n = n + (1) >> 0;
					}
					_tuple$16 = parseNanoseconds(value, n); nsec = _tuple$16[0]; rangeErrString = _tuple$16[1]; err = _tuple$16[2];
					value = value.substring(n);
				}
			} else if (_ref === 531) {
				if (value.length < 2) {
					err = errBad;
					break;
				}
				_tmp$6 = value.substring(0, 2); _tmp$7 = value.substring(2); p = _tmp$6; value = _tmp$7;
				_ref$1 = p;
				if (_ref$1 === "PM") {
					pmSet = true;
				} else if (_ref$1 === "AM") {
					amSet = true;
				} else {
					err = errBad;
				}
			} else if (_ref === 532) {
				if (value.length < 2) {
					err = errBad;
					break;
				}
				_tmp$8 = value.substring(0, 2); _tmp$9 = value.substring(2); p = _tmp$8; value = _tmp$9;
				_ref$2 = p;
				if (_ref$2 === "pm") {
					pmSet = true;
				} else if (_ref$2 === "am") {
					amSet = true;
				} else {
					err = errBad;
				}
			} else if (_ref === 22 || _ref === 24 || _ref === 23 || _ref === 25 || _ref === 26 || _ref === 28 || _ref === 29 || _ref === 27 || _ref === 30) {
				if (((std === 22) || (std === 24)) && value.length >= 1 && (value.charCodeAt(0) === 90)) {
					value = value.substring(1);
					z = $pkg.UTC;
					break;
				}
				_tmp$10 = ""; _tmp$11 = ""; _tmp$12 = ""; _tmp$13 = ""; sign = _tmp$10; hour$1 = _tmp$11; min$1 = _tmp$12; seconds = _tmp$13;
				if ((std === 24) || (std === 29)) {
					if (value.length < 6) {
						err = errBad;
						break;
					}
					if (!((value.charCodeAt(3) === 58))) {
						err = errBad;
						break;
					}
					_tmp$14 = value.substring(0, 1); _tmp$15 = value.substring(1, 3); _tmp$16 = value.substring(4, 6); _tmp$17 = "00"; _tmp$18 = value.substring(6); sign = _tmp$14; hour$1 = _tmp$15; min$1 = _tmp$16; seconds = _tmp$17; value = _tmp$18;
				} else if (std === 28) {
					if (value.length < 3) {
						err = errBad;
						break;
					}
					_tmp$19 = value.substring(0, 1); _tmp$20 = value.substring(1, 3); _tmp$21 = "00"; _tmp$22 = "00"; _tmp$23 = value.substring(3); sign = _tmp$19; hour$1 = _tmp$20; min$1 = _tmp$21; seconds = _tmp$22; value = _tmp$23;
				} else if ((std === 25) || (std === 30)) {
					if (value.length < 9) {
						err = errBad;
						break;
					}
					if (!((value.charCodeAt(3) === 58)) || !((value.charCodeAt(6) === 58))) {
						err = errBad;
						break;
					}
					_tmp$24 = value.substring(0, 1); _tmp$25 = value.substring(1, 3); _tmp$26 = value.substring(4, 6); _tmp$27 = value.substring(7, 9); _tmp$28 = value.substring(9); sign = _tmp$24; hour$1 = _tmp$25; min$1 = _tmp$26; seconds = _tmp$27; value = _tmp$28;
				} else if ((std === 23) || (std === 27)) {
					if (value.length < 7) {
						err = errBad;
						break;
					}
					_tmp$29 = value.substring(0, 1); _tmp$30 = value.substring(1, 3); _tmp$31 = value.substring(3, 5); _tmp$32 = value.substring(5, 7); _tmp$33 = value.substring(7); sign = _tmp$29; hour$1 = _tmp$30; min$1 = _tmp$31; seconds = _tmp$32; value = _tmp$33;
				} else {
					if (value.length < 5) {
						err = errBad;
						break;
					}
					_tmp$34 = value.substring(0, 1); _tmp$35 = value.substring(1, 3); _tmp$36 = value.substring(3, 5); _tmp$37 = "00"; _tmp$38 = value.substring(5); sign = _tmp$34; hour$1 = _tmp$35; min$1 = _tmp$36; seconds = _tmp$37; value = _tmp$38;
				}
				_tmp$39 = 0; _tmp$40 = 0; _tmp$41 = 0; hr = _tmp$39; mm = _tmp$40; ss = _tmp$41;
				_tuple$17 = atoi(hour$1); hr = _tuple$17[0]; err = _tuple$17[1];
				if ($interfaceIsEqual(err, null)) {
					_tuple$18 = atoi(min$1); mm = _tuple$18[0]; err = _tuple$18[1];
				}
				if ($interfaceIsEqual(err, null)) {
					_tuple$19 = atoi(seconds); ss = _tuple$19[0]; err = _tuple$19[1];
				}
				zoneOffset = (x = (((((hr >>> 16 << 16) * 60 >> 0) + (hr << 16 >>> 16) * 60) >> 0) + mm >> 0), (((x >>> 16 << 16) * 60 >> 0) + (x << 16 >>> 16) * 60) >> 0) + ss >> 0;
				_ref$3 = sign.charCodeAt(0);
				if (_ref$3 === 43) {
				} else if (_ref$3 === 45) {
					zoneOffset = -zoneOffset;
				} else {
					err = errBad;
				}
			} else if (_ref === 21) {
				if (value.length >= 3 && value.substring(0, 3) === "UTC") {
					z = $pkg.UTC;
					value = value.substring(3);
					break;
				}
				_tuple$20 = parseTimeZone(value); n$1 = _tuple$20[0]; ok = _tuple$20[1];
				if (!ok) {
					err = errBad;
					break;
				}
				_tmp$42 = value.substring(0, n$1); _tmp$43 = value.substring(n$1); zoneName = _tmp$42; value = _tmp$43;
			} else if (_ref === 31) {
				ndigit = 1 + ((std >> 16 >> 0)) >> 0;
				if (value.length < ndigit) {
					err = errBad;
					break;
				}
				_tuple$21 = parseNanoseconds(value, ndigit); nsec = _tuple$21[0]; rangeErrString = _tuple$21[1]; err = _tuple$21[2];
				value = value.substring(ndigit);
			} else if (_ref === 32) {
				if (value.length < 2 || !((value.charCodeAt(0) === 46)) || value.charCodeAt(1) < 48 || 57 < value.charCodeAt(1)) {
					break;
				}
				i = 0;
				while (i < 9 && (i + 1 >> 0) < value.length && 48 <= value.charCodeAt((i + 1 >> 0)) && value.charCodeAt((i + 1 >> 0)) <= 57) {
					i = i + (1) >> 0;
				}
				_tuple$22 = parseNanoseconds(value, 1 + i >> 0); nsec = _tuple$22[0]; rangeErrString = _tuple$22[1]; err = _tuple$22[2];
				value = value.substring((1 + i >> 0));
			} }
			if (!(rangeErrString === "")) {
				return [new Time.Ptr(new $Int64(0, 0), 0, ($ptrType(Location)).nil), new ParseError.Ptr(alayout, avalue, stdstr, value, ": " + rangeErrString + " out of range")];
			}
			if (!($interfaceIsEqual(err, null))) {
				return [new Time.Ptr(new $Int64(0, 0), 0, ($ptrType(Location)).nil), new ParseError.Ptr(alayout, avalue, stdstr, value, "")];
			}
		}
		if (pmSet && hour < 12) {
			hour = hour + (12) >> 0;
		} else if (amSet && (hour === 12)) {
			hour = 0;
		}
		if (!(z === ($ptrType(Location)).nil)) {
			return [Date(year, (month >> 0), day, hour, min, sec, nsec, z), null];
		}
		if (!((zoneOffset === -1))) {
			t = new Time.Ptr(); $copy(t, Date(year, (month >> 0), day, hour, min, sec, nsec, $pkg.UTC), Time);
			t.sec = (x$1 = t.sec, x$2 = new $Int64(0, zoneOffset), new $Int64(x$1.$high - x$2.$high, x$1.$low - x$2.$low));
			_tuple$23 = local.lookup((x$3 = t.sec, new $Int64(x$3.$high + -15, x$3.$low + 2288912640))); name = _tuple$23[0]; offset = _tuple$23[1];
			if ((offset === zoneOffset) && (zoneName === "" || name === zoneName)) {
				t.loc = local;
				return [t, null];
			}
			t.loc = FixedZone(zoneName, zoneOffset);
			return [t, null];
		}
		if (!(zoneName === "")) {
			t$1 = new Time.Ptr(); $copy(t$1, Date(year, (month >> 0), day, hour, min, sec, nsec, $pkg.UTC), Time);
			_tuple$24 = local.lookupName(zoneName, (x$4 = t$1.sec, new $Int64(x$4.$high + -15, x$4.$low + 2288912640))); offset$1 = _tuple$24[0]; ok$1 = _tuple$24[2];
			if (ok$1) {
				t$1.sec = (x$5 = t$1.sec, x$6 = new $Int64(0, offset$1), new $Int64(x$5.$high - x$6.$high, x$5.$low - x$6.$low));
				t$1.loc = local;
				return [t$1, null];
			}
			if (zoneName.length > 3 && zoneName.substring(0, 3) === "GMT") {
				_tuple$25 = atoi(zoneName.substring(3)); offset$1 = _tuple$25[0];
				offset$1 = (x$7 = 3600, (((offset$1 >>> 16 << 16) * x$7 >> 0) + (offset$1 << 16 >>> 16) * x$7) >> 0);
			}
			t$1.loc = FixedZone(zoneName, offset$1);
			return [t$1, null];
		}
		return [Date(year, (month >> 0), day, hour, min, sec, nsec, defaultLocation), null];
	};
	parseTimeZone = function(value) {
		var length = 0, ok = false, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5, nUpper, c, _ref, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15;
		if (value.length < 3) {
			_tmp = 0; _tmp$1 = false; length = _tmp; ok = _tmp$1;
			return [length, ok];
		}
		if (value.length >= 4 && (value.substring(0, 4) === "ChST" || value.substring(0, 4) === "MeST")) {
			_tmp$2 = 4; _tmp$3 = true; length = _tmp$2; ok = _tmp$3;
			return [length, ok];
		}
		if (value.substring(0, 3) === "GMT") {
			length = parseGMT(value);
			_tmp$4 = length; _tmp$5 = true; length = _tmp$4; ok = _tmp$5;
			return [length, ok];
		}
		nUpper = 0;
		nUpper = 0;
		while (nUpper < 6) {
			if (nUpper >= value.length) {
				break;
			}
			c = value.charCodeAt(nUpper);
			if (c < 65 || 90 < c) {
				break;
			}
			nUpper = nUpper + (1) >> 0;
		}
		_ref = nUpper;
		if (_ref === 0 || _ref === 1 || _ref === 2 || _ref === 6) {
			_tmp$6 = 0; _tmp$7 = false; length = _tmp$6; ok = _tmp$7;
			return [length, ok];
		} else if (_ref === 5) {
			if (value.charCodeAt(4) === 84) {
				_tmp$8 = 5; _tmp$9 = true; length = _tmp$8; ok = _tmp$9;
				return [length, ok];
			}
		} else if (_ref === 4) {
			if (value.charCodeAt(3) === 84) {
				_tmp$10 = 4; _tmp$11 = true; length = _tmp$10; ok = _tmp$11;
				return [length, ok];
			}
		} else if (_ref === 3) {
			_tmp$12 = 3; _tmp$13 = true; length = _tmp$12; ok = _tmp$13;
			return [length, ok];
		}
		_tmp$14 = 0; _tmp$15 = false; length = _tmp$14; ok = _tmp$15;
		return [length, ok];
	};
	parseGMT = function(value) {
		var sign, _tuple$1, x, rem, err;
		value = value.substring(3);
		if (value.length === 0) {
			return 3;
		}
		sign = value.charCodeAt(0);
		if (!((sign === 45)) && !((sign === 43))) {
			return 3;
		}
		_tuple$1 = leadingInt(value.substring(1)); x = _tuple$1[0]; rem = _tuple$1[1]; err = _tuple$1[2];
		if (!($interfaceIsEqual(err, null))) {
			return 3;
		}
		if (sign === 45) {
			x = new $Int64(-x.$high, -x.$low);
		}
		if ((x.$high === 0 && x.$low === 0) || (x.$high < -1 || (x.$high === -1 && x.$low < 4294967282)) || (0 < x.$high || (0 === x.$high && 12 < x.$low))) {
			return 3;
		}
		return (3 + value.length >> 0) - rem.length >> 0;
	};
	parseNanoseconds = function(value, nbytes) {
		var ns = 0, rangeErrString = "", err = null, _tuple$1, scaleDigits, i, x;
		if (!((value.charCodeAt(0) === 46))) {
			err = errBad;
			return [ns, rangeErrString, err];
		}
		_tuple$1 = atoi(value.substring(1, nbytes)); ns = _tuple$1[0]; err = _tuple$1[1];
		if (!($interfaceIsEqual(err, null))) {
			return [ns, rangeErrString, err];
		}
		if (ns < 0 || 1000000000 <= ns) {
			rangeErrString = "fractional second";
			return [ns, rangeErrString, err];
		}
		scaleDigits = 10 - nbytes >> 0;
		i = 0;
		while (i < scaleDigits) {
			ns = (x = 10, (((ns >>> 16 << 16) * x >> 0) + (ns << 16 >>> 16) * x) >> 0);
			i = i + (1) >> 0;
		}
		return [ns, rangeErrString, err];
	};
	leadingInt = function(s) {
		var x = new $Int64(0, 0), rem = "", err = null, i, c, _tmp, _tmp$1, _tmp$2, x$1, x$2, x$3, _tmp$3, _tmp$4, _tmp$5;
		i = 0;
		while (i < s.length) {
			c = s.charCodeAt(i);
			if (c < 48 || c > 57) {
				break;
			}
			if ((x.$high > 214748364 || (x.$high === 214748364 && x.$low >= 3435973835))) {
				_tmp = new $Int64(0, 0); _tmp$1 = ""; _tmp$2 = errLeadingInt; x = _tmp; rem = _tmp$1; err = _tmp$2;
				return [x, rem, err];
			}
			x = (x$1 = (x$2 = $mul64(x, new $Int64(0, 10)), x$3 = new $Int64(0, c), new $Int64(x$2.$high + x$3.$high, x$2.$low + x$3.$low)), new $Int64(x$1.$high - 0, x$1.$low - 48));
			i = i + (1) >> 0;
		}
		_tmp$3 = x; _tmp$4 = s.substring(i); _tmp$5 = null; x = _tmp$3; rem = _tmp$4; err = _tmp$5;
		return [x, rem, err];
	};
	Time.Ptr.prototype.After = function(u) {
		var t, x, x$1, x$2, x$3;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (x = t.sec, x$1 = u.sec, (x.$high > x$1.$high || (x.$high === x$1.$high && x.$low > x$1.$low))) || (x$2 = t.sec, x$3 = u.sec, (x$2.$high === x$3.$high && x$2.$low === x$3.$low)) && t.nsec > u.nsec;
	};
	Time.prototype.After = function(u) { return this.$val.After(u); };
	Time.Ptr.prototype.Before = function(u) {
		var t, x, x$1, x$2, x$3;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (x = t.sec, x$1 = u.sec, (x.$high < x$1.$high || (x.$high === x$1.$high && x.$low < x$1.$low))) || (x$2 = t.sec, x$3 = u.sec, (x$2.$high === x$3.$high && x$2.$low === x$3.$low)) && t.nsec < u.nsec;
	};
	Time.prototype.Before = function(u) { return this.$val.Before(u); };
	Time.Ptr.prototype.Equal = function(u) {
		var t, x, x$1;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (x = t.sec, x$1 = u.sec, (x.$high === x$1.$high && x.$low === x$1.$low)) && (t.nsec === u.nsec);
	};
	Time.prototype.Equal = function(u) { return this.$val.Equal(u); };
	Month.prototype.String = function() {
		var m, x;
		m = this.$val;
		return (x = m - 1 >> 0, ((x < 0 || x >= months.length) ? $throwRuntimeError("index out of range") : months[x]));
	};
	$ptrType(Month).prototype.String = function() { return new Month(this.$get()).String(); };
	Weekday.prototype.String = function() {
		var d;
		d = this.$val;
		return ((d < 0 || d >= days.length) ? $throwRuntimeError("index out of range") : days[d]);
	};
	$ptrType(Weekday).prototype.String = function() { return new Weekday(this.$get()).String(); };
	Time.Ptr.prototype.IsZero = function() {
		var t, x;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (x = t.sec, (x.$high === 0 && x.$low === 0)) && (t.nsec === 0);
	};
	Time.prototype.IsZero = function() { return this.$val.IsZero(); };
	Time.Ptr.prototype.abs = function() {
		var t, l, x, sec, x$1, x$2, x$3, _tuple$1, offset, x$4, x$5;
		t = new Time.Ptr(); $copy(t, this, Time);
		l = t.loc;
		if (l === ($ptrType(Location)).nil || l === localLoc) {
			l = l.get();
		}
		sec = (x = t.sec, new $Int64(x.$high + -15, x.$low + 2288912640));
		if (!(l === utcLoc)) {
			if (!(l.cacheZone === ($ptrType(zone)).nil) && (x$1 = l.cacheStart, (x$1.$high < sec.$high || (x$1.$high === sec.$high && x$1.$low <= sec.$low))) && (x$2 = l.cacheEnd, (sec.$high < x$2.$high || (sec.$high === x$2.$high && sec.$low < x$2.$low)))) {
				sec = (x$3 = new $Int64(0, l.cacheZone.offset), new $Int64(sec.$high + x$3.$high, sec.$low + x$3.$low));
			} else {
				_tuple$1 = l.lookup(sec); offset = _tuple$1[1];
				sec = (x$4 = new $Int64(0, offset), new $Int64(sec.$high + x$4.$high, sec.$low + x$4.$low));
			}
		}
		return (x$5 = new $Int64(sec.$high + 2147483646, sec.$low + 450480384), new $Uint64(x$5.$high, x$5.$low));
	};
	Time.prototype.abs = function() { return this.$val.abs(); };
	Time.Ptr.prototype.locabs = function() {
		var name = "", offset = 0, abs = new $Uint64(0, 0), t, l, x, sec, x$1, x$2, _tuple$1, x$3, x$4;
		t = new Time.Ptr(); $copy(t, this, Time);
		l = t.loc;
		if (l === ($ptrType(Location)).nil || l === localLoc) {
			l = l.get();
		}
		sec = (x = t.sec, new $Int64(x.$high + -15, x.$low + 2288912640));
		if (!(l === utcLoc)) {
			if (!(l.cacheZone === ($ptrType(zone)).nil) && (x$1 = l.cacheStart, (x$1.$high < sec.$high || (x$1.$high === sec.$high && x$1.$low <= sec.$low))) && (x$2 = l.cacheEnd, (sec.$high < x$2.$high || (sec.$high === x$2.$high && sec.$low < x$2.$low)))) {
				name = l.cacheZone.name;
				offset = l.cacheZone.offset;
			} else {
				_tuple$1 = l.lookup(sec); name = _tuple$1[0]; offset = _tuple$1[1];
			}
			sec = (x$3 = new $Int64(0, offset), new $Int64(sec.$high + x$3.$high, sec.$low + x$3.$low));
		} else {
			name = "UTC";
		}
		abs = (x$4 = new $Int64(sec.$high + 2147483646, sec.$low + 450480384), new $Uint64(x$4.$high, x$4.$low));
		return [name, offset, abs];
	};
	Time.prototype.locabs = function() { return this.$val.locabs(); };
	Time.Ptr.prototype.Date = function() {
		var year = 0, month = 0, day = 0, t, _tuple$1;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.date(true); year = _tuple$1[0]; month = _tuple$1[1]; day = _tuple$1[2];
		return [year, month, day];
	};
	Time.prototype.Date = function() { return this.$val.Date(); };
	Time.Ptr.prototype.Year = function() {
		var t, _tuple$1, year;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.date(false); year = _tuple$1[0];
		return year;
	};
	Time.prototype.Year = function() { return this.$val.Year(); };
	Time.Ptr.prototype.Month = function() {
		var t, _tuple$1, month;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.date(true); month = _tuple$1[1];
		return month;
	};
	Time.prototype.Month = function() { return this.$val.Month(); };
	Time.Ptr.prototype.Day = function() {
		var t, _tuple$1, day;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.date(true); day = _tuple$1[2];
		return day;
	};
	Time.prototype.Day = function() { return this.$val.Day(); };
	Time.Ptr.prototype.Weekday = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		return absWeekday(t.abs());
	};
	Time.prototype.Weekday = function() { return this.$val.Weekday(); };
	absWeekday = function(abs) {
		var sec, _q;
		sec = $div64((new $Uint64(abs.$high + 0, abs.$low + 86400)), new $Uint64(0, 604800), true);
		return ((_q = (sec.$low >> 0) / 86400, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0);
	};
	Time.Ptr.prototype.ISOWeek = function() {
		var year = 0, week = 0, t, _tuple$1, month, day, yday, _r, wday, _q, _r$1, jan1wday, _r$2, dec31wday;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.date(true); year = _tuple$1[0]; month = _tuple$1[1]; day = _tuple$1[2]; yday = _tuple$1[3];
		wday = (_r = ((t.Weekday() + 6 >> 0) >> 0) % 7, _r === _r ? _r : $throwRuntimeError("integer divide by zero"));
		week = (_q = (((yday - wday >> 0) + 7 >> 0)) / 7, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		jan1wday = (_r$1 = (((wday - yday >> 0) + 371 >> 0)) % 7, _r$1 === _r$1 ? _r$1 : $throwRuntimeError("integer divide by zero"));
		if (1 <= jan1wday && jan1wday <= 3) {
			week = week + (1) >> 0;
		}
		if (week === 0) {
			year = year - (1) >> 0;
			week = 52;
			if ((jan1wday === 4) || ((jan1wday === 5) && isLeap(year))) {
				week = week + (1) >> 0;
			}
		}
		if ((month === 12) && day >= 29 && wday < 3) {
			dec31wday = (_r$2 = (((wday + 31 >> 0) - day >> 0)) % 7, _r$2 === _r$2 ? _r$2 : $throwRuntimeError("integer divide by zero"));
			if (0 <= dec31wday && dec31wday <= 2) {
				year = year + (1) >> 0;
				week = 1;
			}
		}
		return [year, week];
	};
	Time.prototype.ISOWeek = function() { return this.$val.ISOWeek(); };
	Time.Ptr.prototype.Clock = function() {
		var hour = 0, min = 0, sec = 0, t, _tuple$1;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = absClock(t.abs()); hour = _tuple$1[0]; min = _tuple$1[1]; sec = _tuple$1[2];
		return [hour, min, sec];
	};
	Time.prototype.Clock = function() { return this.$val.Clock(); };
	absClock = function(abs) {
		var hour = 0, min = 0, sec = 0, _q, _q$1;
		sec = ($div64(abs, new $Uint64(0, 86400), true).$low >> 0);
		hour = (_q = sec / 3600, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		sec = sec - (((((hour >>> 16 << 16) * 3600 >> 0) + (hour << 16 >>> 16) * 3600) >> 0)) >> 0;
		min = (_q$1 = sec / 60, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >> 0 : $throwRuntimeError("integer divide by zero"));
		sec = sec - (((((min >>> 16 << 16) * 60 >> 0) + (min << 16 >>> 16) * 60) >> 0)) >> 0;
		return [hour, min, sec];
	};
	Time.Ptr.prototype.Hour = function() {
		var t, _q;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (_q = ($div64(t.abs(), new $Uint64(0, 86400), true).$low >> 0) / 3600, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
	};
	Time.prototype.Hour = function() { return this.$val.Hour(); };
	Time.Ptr.prototype.Minute = function() {
		var t, _q;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (_q = ($div64(t.abs(), new $Uint64(0, 3600), true).$low >> 0) / 60, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
	};
	Time.prototype.Minute = function() { return this.$val.Minute(); };
	Time.Ptr.prototype.Second = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		return ($div64(t.abs(), new $Uint64(0, 60), true).$low >> 0);
	};
	Time.prototype.Second = function() { return this.$val.Second(); };
	Time.Ptr.prototype.Nanosecond = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (t.nsec >> 0);
	};
	Time.prototype.Nanosecond = function() { return this.$val.Nanosecond(); };
	Time.Ptr.prototype.YearDay = function() {
		var t, _tuple$1, yday;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.date(false); yday = _tuple$1[3];
		return yday + 1 >> 0;
	};
	Time.prototype.YearDay = function() { return this.$val.YearDay(); };
	Duration.prototype.String = function() {
		var d, buf, w, u, neg, prec, unit, x, _tuple$1, _tuple$2;
		d = this;
		buf = ($arrayType($Uint8, 32)).zero(); $copy(buf, ($arrayType($Uint8, 32)).zero(), ($arrayType($Uint8, 32)));
		w = 32;
		u = new $Uint64(d.$high, d.$low);
		neg = (d.$high < 0 || (d.$high === 0 && d.$low < 0));
		if (neg) {
			u = new $Uint64(-u.$high, -u.$low);
		}
		if ((u.$high < 0 || (u.$high === 0 && u.$low < 1000000000))) {
			prec = 0;
			unit = 0;
			if ((u.$high === 0 && u.$low === 0)) {
				return "0";
			} else if ((u.$high < 0 || (u.$high === 0 && u.$low < 1000))) {
				prec = 0;
				unit = 110;
			} else if ((u.$high < 0 || (u.$high === 0 && u.$low < 1000000))) {
				prec = 3;
				unit = 117;
			} else {
				prec = 6;
				unit = 109;
			}
			w = w - (2) >> 0;
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = unit;
			(x = w + 1 >> 0, (x < 0 || x >= buf.length) ? $throwRuntimeError("index out of range") : buf[x] = 115);
			_tuple$1 = fmtFrac($subslice(new ($sliceType($Uint8))(buf), 0, w), u, prec); w = _tuple$1[0]; u = _tuple$1[1];
			w = fmtInt($subslice(new ($sliceType($Uint8))(buf), 0, w), u);
		} else {
			w = w - (1) >> 0;
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 115;
			_tuple$2 = fmtFrac($subslice(new ($sliceType($Uint8))(buf), 0, w), u, 9); w = _tuple$2[0]; u = _tuple$2[1];
			w = fmtInt($subslice(new ($sliceType($Uint8))(buf), 0, w), $div64(u, new $Uint64(0, 60), true));
			u = $div64(u, (new $Uint64(0, 60)), false);
			if ((u.$high > 0 || (u.$high === 0 && u.$low > 0))) {
				w = w - (1) >> 0;
				(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 109;
				w = fmtInt($subslice(new ($sliceType($Uint8))(buf), 0, w), $div64(u, new $Uint64(0, 60), true));
				u = $div64(u, (new $Uint64(0, 60)), false);
				if ((u.$high > 0 || (u.$high === 0 && u.$low > 0))) {
					w = w - (1) >> 0;
					(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 104;
					w = fmtInt($subslice(new ($sliceType($Uint8))(buf), 0, w), u);
				}
			}
		}
		if (neg) {
			w = w - (1) >> 0;
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 45;
		}
		return $bytesToString($subslice(new ($sliceType($Uint8))(buf), w));
	};
	$ptrType(Duration).prototype.String = function() { return this.$get().String(); };
	fmtFrac = function(buf, v, prec) {
		var nw = 0, nv = new $Uint64(0, 0), w, print, i, digit, _tmp, _tmp$1;
		w = buf.$length;
		print = false;
		i = 0;
		while (i < prec) {
			digit = $div64(v, new $Uint64(0, 10), true);
			print = print || !((digit.$high === 0 && digit.$low === 0));
			if (print) {
				w = w - (1) >> 0;
				(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = (digit.$low << 24 >>> 24) + 48 << 24 >>> 24;
			}
			v = $div64(v, (new $Uint64(0, 10)), false);
			i = i + (1) >> 0;
		}
		if (print) {
			w = w - (1) >> 0;
			(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = 46;
		}
		_tmp = w; _tmp$1 = v; nw = _tmp; nv = _tmp$1;
		return [nw, nv];
	};
	fmtInt = function(buf, v) {
		var w;
		w = buf.$length;
		if ((v.$high === 0 && v.$low === 0)) {
			w = w - (1) >> 0;
			(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = 48;
		} else {
			while ((v.$high > 0 || (v.$high === 0 && v.$low > 0))) {
				w = w - (1) >> 0;
				(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = ($div64(v, new $Uint64(0, 10), true).$low << 24 >>> 24) + 48 << 24 >>> 24;
				v = $div64(v, (new $Uint64(0, 10)), false);
			}
		}
		return w;
	};
	Duration.prototype.Nanoseconds = function() {
		var d;
		d = this;
		return new $Int64(d.$high, d.$low);
	};
	$ptrType(Duration).prototype.Nanoseconds = function() { return this.$get().Nanoseconds(); };
	Duration.prototype.Seconds = function() {
		var d, sec, nsec;
		d = this;
		sec = $div64(d, new Duration(0, 1000000000), false);
		nsec = $div64(d, new Duration(0, 1000000000), true);
		return $flatten64(sec) + $flatten64(nsec) * 1e-09;
	};
	$ptrType(Duration).prototype.Seconds = function() { return this.$get().Seconds(); };
	Duration.prototype.Minutes = function() {
		var d, min, nsec;
		d = this;
		min = $div64(d, new Duration(13, 4165425152), false);
		nsec = $div64(d, new Duration(13, 4165425152), true);
		return $flatten64(min) + $flatten64(nsec) * 1.6666666666666667e-11;
	};
	$ptrType(Duration).prototype.Minutes = function() { return this.$get().Minutes(); };
	Duration.prototype.Hours = function() {
		var d, hour, nsec;
		d = this;
		hour = $div64(d, new Duration(838, 817405952), false);
		nsec = $div64(d, new Duration(838, 817405952), true);
		return $flatten64(hour) + $flatten64(nsec) * 2.777777777777778e-13;
	};
	$ptrType(Duration).prototype.Hours = function() { return this.$get().Hours(); };
	Time.Ptr.prototype.Add = function(d) {
		var t, x, x$1, x$2, x$3, nsec, x$4, x$5, x$6, x$7;
		t = new Time.Ptr(); $copy(t, this, Time);
		t.sec = (x = t.sec, x$1 = (x$2 = $div64(d, new Duration(0, 1000000000), false), new $Int64(x$2.$high, x$2.$low)), new $Int64(x.$high + x$1.$high, x.$low + x$1.$low));
		nsec = (t.nsec >> 0) + ((x$3 = $div64(d, new Duration(0, 1000000000), true), x$3.$low + ((x$3.$high >> 31) * 4294967296)) >> 0) >> 0;
		if (nsec >= 1000000000) {
			t.sec = (x$4 = t.sec, x$5 = new $Int64(0, 1), new $Int64(x$4.$high + x$5.$high, x$4.$low + x$5.$low));
			nsec = nsec - (1000000000) >> 0;
		} else if (nsec < 0) {
			t.sec = (x$6 = t.sec, x$7 = new $Int64(0, 1), new $Int64(x$6.$high - x$7.$high, x$6.$low - x$7.$low));
			nsec = nsec + (1000000000) >> 0;
		}
		t.nsec = (nsec >>> 0);
		return t;
	};
	Time.prototype.Add = function(d) { return this.$val.Add(d); };
	Time.Ptr.prototype.Sub = function(u) {
		var t, x, x$1, x$2, x$3, x$4, d;
		t = new Time.Ptr(); $copy(t, this, Time);
		d = (x = $mul64((x$1 = (x$2 = t.sec, x$3 = u.sec, new $Int64(x$2.$high - x$3.$high, x$2.$low - x$3.$low)), new Duration(x$1.$high, x$1.$low)), new Duration(0, 1000000000)), x$4 = new Duration(0, ((t.nsec >> 0) - (u.nsec >> 0) >> 0)), new Duration(x.$high + x$4.$high, x.$low + x$4.$low));
		if (u.Add(d).Equal($clone(t, Time))) {
			return d;
		} else if (t.Before($clone(u, Time))) {
			return new Duration(-2147483648, 0);
		} else {
			return new Duration(2147483647, 4294967295);
		}
	};
	Time.prototype.Sub = function(u) { return this.$val.Sub(u); };
	Time.Ptr.prototype.AddDate = function(years, months$1, days$1) {
		var t, _tuple$1, year, month, day, _tuple$2, hour, min, sec;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.Date(); year = _tuple$1[0]; month = _tuple$1[1]; day = _tuple$1[2];
		_tuple$2 = t.Clock(); hour = _tuple$2[0]; min = _tuple$2[1]; sec = _tuple$2[2];
		return Date(year + years >> 0, month + (months$1 >> 0) >> 0, day + days$1 >> 0, hour, min, sec, (t.nsec >> 0), t.loc);
	};
	Time.prototype.AddDate = function(years, months$1, days$1) { return this.$val.AddDate(years, months$1, days$1); };
	Time.Ptr.prototype.date = function(full) {
		var year = 0, month = 0, day = 0, yday = 0, t, _tuple$1;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = absDate(t.abs(), full); year = _tuple$1[0]; month = _tuple$1[1]; day = _tuple$1[2]; yday = _tuple$1[3];
		return [year, month, day, yday];
	};
	Time.prototype.date = function(full) { return this.$val.date(full); };
	absDate = function(abs, full) {
		var year = 0, month = 0, day = 0, yday = 0, d, n, y, x, x$1, x$2, x$3, x$4, x$5, x$6, x$7, x$8, x$9, x$10, _q, x$11, end, begin;
		d = $div64(abs, new $Uint64(0, 86400), false);
		n = $div64(d, new $Uint64(0, 146097), false);
		y = $mul64(new $Uint64(0, 400), n);
		d = (x = $mul64(new $Uint64(0, 146097), n), new $Uint64(d.$high - x.$high, d.$low - x.$low));
		n = $div64(d, new $Uint64(0, 36524), false);
		n = (x$1 = $shiftRightUint64(n, 2), new $Uint64(n.$high - x$1.$high, n.$low - x$1.$low));
		y = (x$2 = $mul64(new $Uint64(0, 100), n), new $Uint64(y.$high + x$2.$high, y.$low + x$2.$low));
		d = (x$3 = $mul64(new $Uint64(0, 36524), n), new $Uint64(d.$high - x$3.$high, d.$low - x$3.$low));
		n = $div64(d, new $Uint64(0, 1461), false);
		y = (x$4 = $mul64(new $Uint64(0, 4), n), new $Uint64(y.$high + x$4.$high, y.$low + x$4.$low));
		d = (x$5 = $mul64(new $Uint64(0, 1461), n), new $Uint64(d.$high - x$5.$high, d.$low - x$5.$low));
		n = $div64(d, new $Uint64(0, 365), false);
		n = (x$6 = $shiftRightUint64(n, 2), new $Uint64(n.$high - x$6.$high, n.$low - x$6.$low));
		y = (x$7 = n, new $Uint64(y.$high + x$7.$high, y.$low + x$7.$low));
		d = (x$8 = $mul64(new $Uint64(0, 365), n), new $Uint64(d.$high - x$8.$high, d.$low - x$8.$low));
		year = ((x$9 = (x$10 = new $Int64(y.$high, y.$low), new $Int64(x$10.$high + -69, x$10.$low + 4075721025)), x$9.$low + ((x$9.$high >> 31) * 4294967296)) >> 0);
		yday = (d.$low >> 0);
		if (!full) {
			return [year, month, day, yday];
		}
		day = yday;
		if (isLeap(year)) {
			if (day > 59) {
				day = day - (1) >> 0;
			} else if (day === 59) {
				month = 2;
				day = 29;
				return [year, month, day, yday];
			}
		}
		month = ((_q = day / 31, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0);
		end = ((x$11 = month + 1 >> 0, ((x$11 < 0 || x$11 >= daysBefore.length) ? $throwRuntimeError("index out of range") : daysBefore[x$11])) >> 0);
		begin = 0;
		if (day >= end) {
			month = month + (1) >> 0;
			begin = end;
		} else {
			begin = (((month < 0 || month >= daysBefore.length) ? $throwRuntimeError("index out of range") : daysBefore[month]) >> 0);
		}
		month = month + (1) >> 0;
		day = (day - begin >> 0) + 1 >> 0;
		return [year, month, day, yday];
	};
	Time.Ptr.prototype.UTC = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		t.loc = $pkg.UTC;
		return t;
	};
	Time.prototype.UTC = function() { return this.$val.UTC(); };
	Time.Ptr.prototype.Local = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		t.loc = $pkg.Local;
		return t;
	};
	Time.prototype.Local = function() { return this.$val.Local(); };
	Time.Ptr.prototype.In = function(loc) {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		if (loc === ($ptrType(Location)).nil) {
			$panic(new $String("time: missing Location in call to Time.In"));
		}
		t.loc = loc;
		return t;
	};
	Time.prototype.In = function(loc) { return this.$val.In(loc); };
	Time.Ptr.prototype.Location = function() {
		var t, l;
		t = new Time.Ptr(); $copy(t, this, Time);
		l = t.loc;
		if (l === ($ptrType(Location)).nil) {
			l = $pkg.UTC;
		}
		return l;
	};
	Time.prototype.Location = function() { return this.$val.Location(); };
	Time.Ptr.prototype.Zone = function() {
		var name = "", offset = 0, t, _tuple$1, x;
		t = new Time.Ptr(); $copy(t, this, Time);
		_tuple$1 = t.loc.lookup((x = t.sec, new $Int64(x.$high + -15, x.$low + 2288912640))); name = _tuple$1[0]; offset = _tuple$1[1];
		return [name, offset];
	};
	Time.prototype.Zone = function() { return this.$val.Zone(); };
	Time.Ptr.prototype.Unix = function() {
		var t, x;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (x = t.sec, new $Int64(x.$high + -15, x.$low + 2288912640));
	};
	Time.prototype.Unix = function() { return this.$val.Unix(); };
	Time.Ptr.prototype.UnixNano = function() {
		var t, x, x$1, x$2, x$3;
		t = new Time.Ptr(); $copy(t, this, Time);
		return (x = $mul64(((x$1 = t.sec, new $Int64(x$1.$high + -15, x$1.$low + 2288912640))), new $Int64(0, 1000000000)), x$2 = (x$3 = t.nsec, new $Int64(0, x$3.constructor === Number ? x$3 : 1)), new $Int64(x.$high + x$2.$high, x.$low + x$2.$low));
	};
	Time.prototype.UnixNano = function() { return this.$val.UnixNano(); };
	Time.Ptr.prototype.MarshalBinary = function() {
		var t, offsetMin, _tuple$1, offset, _r, _q, enc;
		t = new Time.Ptr(); $copy(t, this, Time);
		offsetMin = 0;
		if (t.Location() === utcLoc) {
			offsetMin = -1;
		} else {
			_tuple$1 = t.Zone(); offset = _tuple$1[1];
			if (!(((_r = offset % 60, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) === 0))) {
				return [($sliceType($Uint8)).nil, errors.New("Time.MarshalBinary: zone offset has fractional minute")];
			}
			offset = (_q = offset / (60), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
			if (offset < -32768 || (offset === -1) || offset > 32767) {
				return [($sliceType($Uint8)).nil, errors.New("Time.MarshalBinary: unexpected zone offset")];
			}
			offsetMin = (offset << 16 >> 16);
		}
		enc = new ($sliceType($Uint8))([1, ($shiftRightInt64(t.sec, 56).$low << 24 >>> 24), ($shiftRightInt64(t.sec, 48).$low << 24 >>> 24), ($shiftRightInt64(t.sec, 40).$low << 24 >>> 24), ($shiftRightInt64(t.sec, 32).$low << 24 >>> 24), ($shiftRightInt64(t.sec, 24).$low << 24 >>> 24), ($shiftRightInt64(t.sec, 16).$low << 24 >>> 24), ($shiftRightInt64(t.sec, 8).$low << 24 >>> 24), (t.sec.$low << 24 >>> 24), ((t.nsec >>> 24 >>> 0) << 24 >>> 24), ((t.nsec >>> 16 >>> 0) << 24 >>> 24), ((t.nsec >>> 8 >>> 0) << 24 >>> 24), (t.nsec << 24 >>> 24), ((offsetMin >> 8 << 16 >> 16) << 24 >>> 24), (offsetMin << 24 >>> 24)]);
		return [enc, null];
	};
	Time.prototype.MarshalBinary = function() { return this.$val.MarshalBinary(); };
	Time.Ptr.prototype.UnmarshalBinary = function(data$1) {
		var t, buf, x, x$1, x$2, x$3, x$4, x$5, x$6, x$7, x$8, x$9, x$10, x$11, x$12, x$13, x$14, offset, _tuple$1, x$15, localoff;
		t = this;
		buf = data$1;
		if (buf.$length === 0) {
			return errors.New("Time.UnmarshalBinary: no data");
		}
		if (!((((0 < 0 || 0 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 0]) === 1))) {
			return errors.New("Time.UnmarshalBinary: unsupported version");
		}
		if (!((buf.$length === 15))) {
			return errors.New("Time.UnmarshalBinary: invalid length");
		}
		buf = $subslice(buf, 1);
		t.sec = (x = (x$1 = (x$2 = (x$3 = (x$4 = (x$5 = (x$6 = new $Int64(0, ((7 < 0 || 7 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 7])), x$7 = $shiftLeft64(new $Int64(0, ((6 < 0 || 6 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 6])), 8), new $Int64(x$6.$high | x$7.$high, (x$6.$low | x$7.$low) >>> 0)), x$8 = $shiftLeft64(new $Int64(0, ((5 < 0 || 5 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 5])), 16), new $Int64(x$5.$high | x$8.$high, (x$5.$low | x$8.$low) >>> 0)), x$9 = $shiftLeft64(new $Int64(0, ((4 < 0 || 4 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 4])), 24), new $Int64(x$4.$high | x$9.$high, (x$4.$low | x$9.$low) >>> 0)), x$10 = $shiftLeft64(new $Int64(0, ((3 < 0 || 3 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 3])), 32), new $Int64(x$3.$high | x$10.$high, (x$3.$low | x$10.$low) >>> 0)), x$11 = $shiftLeft64(new $Int64(0, ((2 < 0 || 2 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 2])), 40), new $Int64(x$2.$high | x$11.$high, (x$2.$low | x$11.$low) >>> 0)), x$12 = $shiftLeft64(new $Int64(0, ((1 < 0 || 1 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 1])), 48), new $Int64(x$1.$high | x$12.$high, (x$1.$low | x$12.$low) >>> 0)), x$13 = $shiftLeft64(new $Int64(0, ((0 < 0 || 0 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 0])), 56), new $Int64(x.$high | x$13.$high, (x.$low | x$13.$low) >>> 0));
		buf = $subslice(buf, 8);
		t.nsec = (((((((3 < 0 || 3 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 3]) >> 0) | ((((2 < 0 || 2 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 2]) >> 0) << 8 >> 0)) | ((((1 < 0 || 1 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 1]) >> 0) << 16 >> 0)) | ((((0 < 0 || 0 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 0]) >> 0) << 24 >> 0)) >>> 0);
		buf = $subslice(buf, 4);
		offset = (x$14 = (((((1 < 0 || 1 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 1]) << 16 >> 16) | ((((0 < 0 || 0 >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + 0]) << 16 >> 16) << 8 << 16 >> 16)) >> 0), (((x$14 >>> 16 << 16) * 60 >> 0) + (x$14 << 16 >>> 16) * 60) >> 0);
		if (offset === -60) {
			t.loc = utcLoc;
		} else {
			_tuple$1 = $pkg.Local.lookup((x$15 = t.sec, new $Int64(x$15.$high + -15, x$15.$low + 2288912640))); localoff = _tuple$1[1];
			if (offset === localoff) {
				t.loc = $pkg.Local;
			} else {
				t.loc = FixedZone("", offset);
			}
		}
		return null;
	};
	Time.prototype.UnmarshalBinary = function(data$1) { return this.$val.UnmarshalBinary(data$1); };
	Time.Ptr.prototype.GobEncode = function() {
		var t;
		t = new Time.Ptr(); $copy(t, this, Time);
		return t.MarshalBinary();
	};
	Time.prototype.GobEncode = function() { return this.$val.GobEncode(); };
	Time.Ptr.prototype.GobDecode = function(data$1) {
		var t;
		t = this;
		return t.UnmarshalBinary(data$1);
	};
	Time.prototype.GobDecode = function(data$1) { return this.$val.GobDecode(data$1); };
	Time.Ptr.prototype.MarshalJSON = function() {
		var t, y;
		t = new Time.Ptr(); $copy(t, this, Time);
		y = t.Year();
		if (y < 0 || y >= 10000) {
			return [($sliceType($Uint8)).nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")];
		}
		return [new ($sliceType($Uint8))($stringToBytes(t.Format("\"2006-01-02T15:04:05.999999999Z07:00\""))), null];
	};
	Time.prototype.MarshalJSON = function() { return this.$val.MarshalJSON(); };
	Time.Ptr.prototype.UnmarshalJSON = function(data$1) {
		var err = null, t, _tuple$1;
		t = this;
		_tuple$1 = Parse("\"2006-01-02T15:04:05Z07:00\"", $bytesToString(data$1)); $copy(t, _tuple$1[0], Time); err = _tuple$1[1];
		return err;
	};
	Time.prototype.UnmarshalJSON = function(data$1) { return this.$val.UnmarshalJSON(data$1); };
	Time.Ptr.prototype.MarshalText = function() {
		var t, y;
		t = new Time.Ptr(); $copy(t, this, Time);
		y = t.Year();
		if (y < 0 || y >= 10000) {
			return [($sliceType($Uint8)).nil, errors.New("Time.MarshalText: year outside of range [0,9999]")];
		}
		return [new ($sliceType($Uint8))($stringToBytes(t.Format("2006-01-02T15:04:05.999999999Z07:00"))), null];
	};
	Time.prototype.MarshalText = function() { return this.$val.MarshalText(); };
	Time.Ptr.prototype.UnmarshalText = function(data$1) {
		var err = null, t, _tuple$1;
		t = this;
		_tuple$1 = Parse("2006-01-02T15:04:05Z07:00", $bytesToString(data$1)); $copy(t, _tuple$1[0], Time); err = _tuple$1[1];
		return err;
	};
	Time.prototype.UnmarshalText = function(data$1) { return this.$val.UnmarshalText(data$1); };
	Unix = $pkg.Unix = function(sec, nsec) {
		var n, x, x$1, x$2, x$3;
		if ((nsec.$high < 0 || (nsec.$high === 0 && nsec.$low < 0)) || (nsec.$high > 0 || (nsec.$high === 0 && nsec.$low >= 1000000000))) {
			n = $div64(nsec, new $Int64(0, 1000000000), false);
			sec = (x = n, new $Int64(sec.$high + x.$high, sec.$low + x.$low));
			nsec = (x$1 = $mul64(n, new $Int64(0, 1000000000)), new $Int64(nsec.$high - x$1.$high, nsec.$low - x$1.$low));
			if ((nsec.$high < 0 || (nsec.$high === 0 && nsec.$low < 0))) {
				nsec = (x$2 = new $Int64(0, 1000000000), new $Int64(nsec.$high + x$2.$high, nsec.$low + x$2.$low));
				sec = (x$3 = new $Int64(0, 1), new $Int64(sec.$high - x$3.$high, sec.$low - x$3.$low));
			}
		}
		return new Time.Ptr(new $Int64(sec.$high + 14, sec.$low + 2006054656), (nsec.$low >>> 0), $pkg.Local);
	};
	isLeap = function(year) {
		var _r, _r$1, _r$2;
		return ((_r = year % 4, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) === 0) && (!(((_r$1 = year % 100, _r$1 === _r$1 ? _r$1 : $throwRuntimeError("integer divide by zero")) === 0)) || ((_r$2 = year % 400, _r$2 === _r$2 ? _r$2 : $throwRuntimeError("integer divide by zero")) === 0));
	};
	norm = function(hi, lo, base) {
		var nhi = 0, nlo = 0, _q, n, _q$1, n$1, _tmp, _tmp$1;
		if (lo < 0) {
			n = (_q = ((-lo - 1 >> 0)) / base, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) + 1 >> 0;
			hi = hi - (n) >> 0;
			lo = lo + (((((n >>> 16 << 16) * base >> 0) + (n << 16 >>> 16) * base) >> 0)) >> 0;
		}
		if (lo >= base) {
			n$1 = (_q$1 = lo / base, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >> 0 : $throwRuntimeError("integer divide by zero"));
			hi = hi + (n$1) >> 0;
			lo = lo - (((((n$1 >>> 16 << 16) * base >> 0) + (n$1 << 16 >>> 16) * base) >> 0)) >> 0;
		}
		_tmp = hi; _tmp$1 = lo; nhi = _tmp; nlo = _tmp$1;
		return [nhi, nlo];
	};
	Date = $pkg.Date = function(year, month, day, hour, min, sec, nsec, loc) {
		var m, _tuple$1, _tuple$2, _tuple$3, _tuple$4, _tuple$5, x, x$1, y, n, x$2, d, x$3, x$4, x$5, x$6, x$7, x$8, x$9, x$10, x$11, abs, x$12, x$13, unix, _tuple$6, offset, start, end, x$14, utc, _tuple$7, _tuple$8, x$15;
		if (loc === ($ptrType(Location)).nil) {
			$panic(new $String("time: missing Location in call to Date"));
		}
		m = (month >> 0) - 1 >> 0;
		_tuple$1 = norm(year, m, 12); year = _tuple$1[0]; m = _tuple$1[1];
		month = (m >> 0) + 1 >> 0;
		_tuple$2 = norm(sec, nsec, 1000000000); sec = _tuple$2[0]; nsec = _tuple$2[1];
		_tuple$3 = norm(min, sec, 60); min = _tuple$3[0]; sec = _tuple$3[1];
		_tuple$4 = norm(hour, min, 60); hour = _tuple$4[0]; min = _tuple$4[1];
		_tuple$5 = norm(day, hour, 24); day = _tuple$5[0]; hour = _tuple$5[1];
		y = (x = (x$1 = new $Int64(0, year), new $Int64(x$1.$high - -69, x$1.$low - 4075721025)), new $Uint64(x.$high, x.$low));
		n = $div64(y, new $Uint64(0, 400), false);
		y = (x$2 = $mul64(new $Uint64(0, 400), n), new $Uint64(y.$high - x$2.$high, y.$low - x$2.$low));
		d = $mul64(new $Uint64(0, 146097), n);
		n = $div64(y, new $Uint64(0, 100), false);
		y = (x$3 = $mul64(new $Uint64(0, 100), n), new $Uint64(y.$high - x$3.$high, y.$low - x$3.$low));
		d = (x$4 = $mul64(new $Uint64(0, 36524), n), new $Uint64(d.$high + x$4.$high, d.$low + x$4.$low));
		n = $div64(y, new $Uint64(0, 4), false);
		y = (x$5 = $mul64(new $Uint64(0, 4), n), new $Uint64(y.$high - x$5.$high, y.$low - x$5.$low));
		d = (x$6 = $mul64(new $Uint64(0, 1461), n), new $Uint64(d.$high + x$6.$high, d.$low + x$6.$low));
		n = y;
		d = (x$7 = $mul64(new $Uint64(0, 365), n), new $Uint64(d.$high + x$7.$high, d.$low + x$7.$low));
		d = (x$8 = new $Uint64(0, (x$9 = month - 1 >> 0, ((x$9 < 0 || x$9 >= daysBefore.length) ? $throwRuntimeError("index out of range") : daysBefore[x$9]))), new $Uint64(d.$high + x$8.$high, d.$low + x$8.$low));
		if (isLeap(year) && month >= 3) {
			d = (x$10 = new $Uint64(0, 1), new $Uint64(d.$high + x$10.$high, d.$low + x$10.$low));
		}
		d = (x$11 = new $Uint64(0, (day - 1 >> 0)), new $Uint64(d.$high + x$11.$high, d.$low + x$11.$low));
		abs = $mul64(d, new $Uint64(0, 86400));
		abs = (x$12 = new $Uint64(0, ((((((hour >>> 16 << 16) * 3600 >> 0) + (hour << 16 >>> 16) * 3600) >> 0) + ((((min >>> 16 << 16) * 60 >> 0) + (min << 16 >>> 16) * 60) >> 0) >> 0) + sec >> 0)), new $Uint64(abs.$high + x$12.$high, abs.$low + x$12.$low));
		unix = (x$13 = new $Int64(abs.$high, abs.$low), new $Int64(x$13.$high + -2147483647, x$13.$low + 3844486912));
		_tuple$6 = loc.lookup(unix); offset = _tuple$6[1]; start = _tuple$6[3]; end = _tuple$6[4];
		if (!((offset === 0))) {
			utc = (x$14 = new $Int64(0, offset), new $Int64(unix.$high - x$14.$high, unix.$low - x$14.$low));
			if ((utc.$high < start.$high || (utc.$high === start.$high && utc.$low < start.$low))) {
				_tuple$7 = loc.lookup(new $Int64(start.$high - 0, start.$low - 1)); offset = _tuple$7[1];
			} else if ((utc.$high > end.$high || (utc.$high === end.$high && utc.$low >= end.$low))) {
				_tuple$8 = loc.lookup(end); offset = _tuple$8[1];
			}
			unix = (x$15 = new $Int64(0, offset), new $Int64(unix.$high - x$15.$high, unix.$low - x$15.$low));
		}
		return new Time.Ptr(new $Int64(unix.$high + 14, unix.$low + 2006054656), (nsec >>> 0), loc);
	};
	Time.Ptr.prototype.Truncate = function(d) {
		var t, _tuple$1, r;
		t = new Time.Ptr(); $copy(t, this, Time);
		if ((d.$high < 0 || (d.$high === 0 && d.$low <= 0))) {
			return t;
		}
		_tuple$1 = div($clone(t, Time), d); r = _tuple$1[1];
		return t.Add(new Duration(-r.$high, -r.$low));
	};
	Time.prototype.Truncate = function(d) { return this.$val.Truncate(d); };
	Time.Ptr.prototype.Round = function(d) {
		var t, _tuple$1, r, x;
		t = new Time.Ptr(); $copy(t, this, Time);
		if ((d.$high < 0 || (d.$high === 0 && d.$low <= 0))) {
			return t;
		}
		_tuple$1 = div($clone(t, Time), d); r = _tuple$1[1];
		if ((x = new Duration(r.$high + r.$high, r.$low + r.$low), (x.$high < d.$high || (x.$high === d.$high && x.$low < d.$low)))) {
			return t.Add(new Duration(-r.$high, -r.$low));
		}
		return t.Add(new Duration(d.$high - r.$high, d.$low - r.$low));
	};
	Time.prototype.Round = function(d) { return this.$val.Round(d); };
	div = function(t, d) {
		var qmod2 = 0, r = new Duration(0, 0), neg, nsec, x, x$1, x$2, x$3, x$4, x$5, _q, _r, x$6, d1, x$7, x$8, x$9, x$10, x$11, sec, tmp, u1, u0, _tmp, _tmp$1, u0x, x$12, _tmp$2, _tmp$3, x$13, x$14, d1$1, x$15, d0, _tmp$4, _tmp$5, x$16, x$17, x$18, x$19;
		neg = false;
		nsec = (t.nsec >> 0);
		if ((x = t.sec, (x.$high < 0 || (x.$high === 0 && x.$low < 0)))) {
			neg = true;
			t.sec = (x$1 = t.sec, new $Int64(-x$1.$high, -x$1.$low));
			nsec = -nsec;
			if (nsec < 0) {
				nsec = nsec + (1000000000) >> 0;
				t.sec = (x$2 = t.sec, x$3 = new $Int64(0, 1), new $Int64(x$2.$high - x$3.$high, x$2.$low - x$3.$low));
			}
		}
		if ((d.$high < 0 || (d.$high === 0 && d.$low < 1000000000)) && (x$4 = $div64(new Duration(0, 1000000000), (new Duration(d.$high + d.$high, d.$low + d.$low)), true), (x$4.$high === 0 && x$4.$low === 0))) {
			qmod2 = ((_q = nsec / ((d.$low + ((d.$high >> 31) * 4294967296)) >> 0), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0) & 1;
			r = new Duration(0, (_r = nsec % ((d.$low + ((d.$high >> 31) * 4294967296)) >> 0), _r === _r ? _r : $throwRuntimeError("integer divide by zero")));
		} else if ((x$5 = $div64(d, new Duration(0, 1000000000), true), (x$5.$high === 0 && x$5.$low === 0))) {
			d1 = (x$6 = $div64(d, new Duration(0, 1000000000), false), new $Int64(x$6.$high, x$6.$low));
			qmod2 = ((x$7 = $div64(t.sec, d1, false), x$7.$low + ((x$7.$high >> 31) * 4294967296)) >> 0) & 1;
			r = (x$8 = $mul64((x$9 = $div64(t.sec, d1, true), new Duration(x$9.$high, x$9.$low)), new Duration(0, 1000000000)), x$10 = new Duration(0, nsec), new Duration(x$8.$high + x$10.$high, x$8.$low + x$10.$low));
		} else {
			sec = (x$11 = t.sec, new $Uint64(x$11.$high, x$11.$low));
			tmp = $mul64(($shiftRightUint64(sec, 32)), new $Uint64(0, 1000000000));
			u1 = $shiftRightUint64(tmp, 32);
			u0 = $shiftLeft64(tmp, 32);
			tmp = $mul64(new $Uint64(sec.$high & 0, (sec.$low & 4294967295) >>> 0), new $Uint64(0, 1000000000));
			_tmp = u0; _tmp$1 = new $Uint64(u0.$high + tmp.$high, u0.$low + tmp.$low); u0x = _tmp; u0 = _tmp$1;
			if ((u0.$high < u0x.$high || (u0.$high === u0x.$high && u0.$low < u0x.$low))) {
				u1 = (x$12 = new $Uint64(0, 1), new $Uint64(u1.$high + x$12.$high, u1.$low + x$12.$low));
			}
			_tmp$2 = u0; _tmp$3 = (x$13 = new $Uint64(0, nsec), new $Uint64(u0.$high + x$13.$high, u0.$low + x$13.$low)); u0x = _tmp$2; u0 = _tmp$3;
			if ((u0.$high < u0x.$high || (u0.$high === u0x.$high && u0.$low < u0x.$low))) {
				u1 = (x$14 = new $Uint64(0, 1), new $Uint64(u1.$high + x$14.$high, u1.$low + x$14.$low));
			}
			d1$1 = new $Uint64(d.$high, d.$low);
			while (!((x$15 = $shiftRightUint64(d1$1, 63), (x$15.$high === 0 && x$15.$low === 1)))) {
				d1$1 = $shiftLeft64(d1$1, (1));
			}
			d0 = new $Uint64(0, 0);
			while (true) {
				qmod2 = 0;
				if ((u1.$high > d1$1.$high || (u1.$high === d1$1.$high && u1.$low > d1$1.$low)) || (u1.$high === d1$1.$high && u1.$low === d1$1.$low) && (u0.$high > d0.$high || (u0.$high === d0.$high && u0.$low >= d0.$low))) {
					qmod2 = 1;
					_tmp$4 = u0; _tmp$5 = new $Uint64(u0.$high - d0.$high, u0.$low - d0.$low); u0x = _tmp$4; u0 = _tmp$5;
					if ((u0.$high > u0x.$high || (u0.$high === u0x.$high && u0.$low > u0x.$low))) {
						u1 = (x$16 = new $Uint64(0, 1), new $Uint64(u1.$high - x$16.$high, u1.$low - x$16.$low));
					}
					u1 = (x$17 = d1$1, new $Uint64(u1.$high - x$17.$high, u1.$low - x$17.$low));
				}
				if ((d1$1.$high === 0 && d1$1.$low === 0) && (x$18 = new $Uint64(d.$high, d.$low), (d0.$high === x$18.$high && d0.$low === x$18.$low))) {
					break;
				}
				d0 = $shiftRightUint64(d0, (1));
				d0 = (x$19 = $shiftLeft64((new $Uint64(d1$1.$high & 0, (d1$1.$low & 1) >>> 0)), 63), new $Uint64(d0.$high | x$19.$high, (d0.$low | x$19.$low) >>> 0));
				d1$1 = $shiftRightUint64(d1$1, (1));
			}
			r = new Duration(u0.$high, u0.$low);
		}
		if (neg && !((r.$high === 0 && r.$low === 0))) {
			qmod2 = (qmod2 ^ (1)) >> 0;
			r = new Duration(d.$high - r.$high, d.$low - r.$low);
		}
		return [qmod2, r];
	};
	Location.Ptr.prototype.get = function() {
		var l;
		l = this;
		if (l === ($ptrType(Location)).nil) {
			return utcLoc;
		}
		if (l === localLoc) {
			localOnce.Do(initLocal);
		}
		return l;
	};
	Location.prototype.get = function() { return this.$val.get(); };
	Location.Ptr.prototype.String = function() {
		var l;
		l = this;
		return l.get().name;
	};
	Location.prototype.String = function() { return this.$val.String(); };
	FixedZone = $pkg.FixedZone = function(name, offset) {
		var l, x;
		l = new Location.Ptr(name, new ($sliceType(zone))([new zone.Ptr(name, offset, false)]), new ($sliceType(zoneTrans))([new zoneTrans.Ptr(new $Int64(-2147483648, 0), 0, false, false)]), new $Int64(-2147483648, 0), new $Int64(2147483647, 4294967295), ($ptrType(zone)).nil);
		l.cacheZone = (x = l.zone, ((0 < 0 || 0 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + 0]));
		return l;
	};
	Location.Ptr.prototype.lookup = function(sec) {
		var name = "", offset = 0, isDST = false, start = new $Int64(0, 0), end = new $Int64(0, 0), l, zone$1, x, x$1, x$2, x$3, x$4, x$5, zone$2, x$6, tx, lo, hi, _q, m, lim, x$7, x$8, zone$3;
		l = this;
		l = l.get();
		if (l.zone.$length === 0) {
			name = "UTC";
			offset = 0;
			isDST = false;
			start = new $Int64(-2147483648, 0);
			end = new $Int64(2147483647, 4294967295);
			return [name, offset, isDST, start, end];
		}
		zone$1 = l.cacheZone;
		if (!(zone$1 === ($ptrType(zone)).nil) && (x = l.cacheStart, (x.$high < sec.$high || (x.$high === sec.$high && x.$low <= sec.$low))) && (x$1 = l.cacheEnd, (sec.$high < x$1.$high || (sec.$high === x$1.$high && sec.$low < x$1.$low)))) {
			name = zone$1.name;
			offset = zone$1.offset;
			isDST = zone$1.isDST;
			start = l.cacheStart;
			end = l.cacheEnd;
			return [name, offset, isDST, start, end];
		}
		if ((l.tx.$length === 0) || (x$2 = (x$3 = l.tx, ((0 < 0 || 0 >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + 0])).when, (sec.$high < x$2.$high || (sec.$high === x$2.$high && sec.$low < x$2.$low)))) {
			zone$2 = (x$4 = l.zone, x$5 = l.lookupFirstZone(), ((x$5 < 0 || x$5 >= x$4.$length) ? $throwRuntimeError("index out of range") : x$4.$array[x$4.$offset + x$5]));
			name = zone$2.name;
			offset = zone$2.offset;
			isDST = zone$2.isDST;
			start = new $Int64(-2147483648, 0);
			if (l.tx.$length > 0) {
				end = (x$6 = l.tx, ((0 < 0 || 0 >= x$6.$length) ? $throwRuntimeError("index out of range") : x$6.$array[x$6.$offset + 0])).when;
			} else {
				end = new $Int64(2147483647, 4294967295);
			}
			return [name, offset, isDST, start, end];
		}
		tx = l.tx;
		end = new $Int64(2147483647, 4294967295);
		lo = 0;
		hi = tx.$length;
		while ((hi - lo >> 0) > 1) {
			m = lo + (_q = ((hi - lo >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			lim = ((m < 0 || m >= tx.$length) ? $throwRuntimeError("index out of range") : tx.$array[tx.$offset + m]).when;
			if ((sec.$high < lim.$high || (sec.$high === lim.$high && sec.$low < lim.$low))) {
				end = lim;
				hi = m;
			} else {
				lo = m;
			}
		}
		zone$3 = (x$7 = l.zone, x$8 = ((lo < 0 || lo >= tx.$length) ? $throwRuntimeError("index out of range") : tx.$array[tx.$offset + lo]).index, ((x$8 < 0 || x$8 >= x$7.$length) ? $throwRuntimeError("index out of range") : x$7.$array[x$7.$offset + x$8]));
		name = zone$3.name;
		offset = zone$3.offset;
		isDST = zone$3.isDST;
		start = ((lo < 0 || lo >= tx.$length) ? $throwRuntimeError("index out of range") : tx.$array[tx.$offset + lo]).when;
		return [name, offset, isDST, start, end];
	};
	Location.prototype.lookup = function(sec) { return this.$val.lookup(sec); };
	Location.Ptr.prototype.lookupFirstZone = function() {
		var l, x, x$1, x$2, x$3, zi, x$4, _ref, _i, zi$1, x$5;
		l = this;
		if (!l.firstZoneUsed()) {
			return 0;
		}
		if (l.tx.$length > 0 && (x = l.zone, x$1 = (x$2 = l.tx, ((0 < 0 || 0 >= x$2.$length) ? $throwRuntimeError("index out of range") : x$2.$array[x$2.$offset + 0])).index, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1])).isDST) {
			zi = ((x$3 = l.tx, ((0 < 0 || 0 >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + 0])).index >> 0) - 1 >> 0;
			while (zi >= 0) {
				if (!(x$4 = l.zone, ((zi < 0 || zi >= x$4.$length) ? $throwRuntimeError("index out of range") : x$4.$array[x$4.$offset + zi])).isDST) {
					return zi;
				}
				zi = zi - (1) >> 0;
			}
		}
		_ref = l.zone;
		_i = 0;
		while (_i < _ref.$length) {
			zi$1 = _i;
			if (!(x$5 = l.zone, ((zi$1 < 0 || zi$1 >= x$5.$length) ? $throwRuntimeError("index out of range") : x$5.$array[x$5.$offset + zi$1])).isDST) {
				return zi$1;
			}
			_i++;
		}
		return 0;
	};
	Location.prototype.lookupFirstZone = function() { return this.$val.lookupFirstZone(); };
	Location.Ptr.prototype.firstZoneUsed = function() {
		var l, _ref, _i, tx;
		l = this;
		_ref = l.tx;
		_i = 0;
		while (_i < _ref.$length) {
			tx = new zoneTrans.Ptr(); $copy(tx, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), zoneTrans);
			if (tx.index === 0) {
				return true;
			}
			_i++;
		}
		return false;
	};
	Location.prototype.firstZoneUsed = function() { return this.$val.firstZoneUsed(); };
	Location.Ptr.prototype.lookupName = function(name, unix) {
		var offset = 0, isDST = false, ok = false, l, _ref, _i, i, x, zone$1, _tuple$1, x$1, nam, offset$1, isDST$1, _tmp, _tmp$1, _tmp$2, _ref$1, _i$1, i$1, x$2, zone$2, _tmp$3, _tmp$4, _tmp$5;
		l = this;
		l = l.get();
		_ref = l.zone;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			zone$1 = (x = l.zone, ((i < 0 || i >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + i]));
			if (zone$1.name === name) {
				_tuple$1 = l.lookup((x$1 = new $Int64(0, zone$1.offset), new $Int64(unix.$high - x$1.$high, unix.$low - x$1.$low))); nam = _tuple$1[0]; offset$1 = _tuple$1[1]; isDST$1 = _tuple$1[2];
				if (nam === zone$1.name) {
					_tmp = offset$1; _tmp$1 = isDST$1; _tmp$2 = true; offset = _tmp; isDST = _tmp$1; ok = _tmp$2;
					return [offset, isDST, ok];
				}
			}
			_i++;
		}
		_ref$1 = l.zone;
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			i$1 = _i$1;
			zone$2 = (x$2 = l.zone, ((i$1 < 0 || i$1 >= x$2.$length) ? $throwRuntimeError("index out of range") : x$2.$array[x$2.$offset + i$1]));
			if (zone$2.name === name) {
				_tmp$3 = zone$2.offset; _tmp$4 = zone$2.isDST; _tmp$5 = true; offset = _tmp$3; isDST = _tmp$4; ok = _tmp$5;
				return [offset, isDST, ok];
			}
			_i$1++;
		}
		return [offset, isDST, ok];
	};
	Location.prototype.lookupName = function(name, unix) { return this.$val.lookupName(name, unix); };
	$pkg.$init = function() {
		($ptrType(ParseError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		ParseError.init([["Layout", "Layout", "", $String, ""], ["Value", "Value", "", $String, ""], ["LayoutElem", "LayoutElem", "", $String, ""], ["ValueElem", "ValueElem", "", $String, ""], ["Message", "Message", "", $String, ""]]);
		Time.methods = [["Add", "Add", "", [Duration], [Time], false, -1], ["AddDate", "AddDate", "", [$Int, $Int, $Int], [Time], false, -1], ["After", "After", "", [Time], [$Bool], false, -1], ["Before", "Before", "", [Time], [$Bool], false, -1], ["Clock", "Clock", "", [], [$Int, $Int, $Int], false, -1], ["Date", "Date", "", [], [$Int, Month, $Int], false, -1], ["Day", "Day", "", [], [$Int], false, -1], ["Equal", "Equal", "", [Time], [$Bool], false, -1], ["Format", "Format", "", [$String], [$String], false, -1], ["GobEncode", "GobEncode", "", [], [($sliceType($Uint8)), $error], false, -1], ["Hour", "Hour", "", [], [$Int], false, -1], ["ISOWeek", "ISOWeek", "", [], [$Int, $Int], false, -1], ["In", "In", "", [($ptrType(Location))], [Time], false, -1], ["IsZero", "IsZero", "", [], [$Bool], false, -1], ["Local", "Local", "", [], [Time], false, -1], ["Location", "Location", "", [], [($ptrType(Location))], false, -1], ["MarshalBinary", "MarshalBinary", "", [], [($sliceType($Uint8)), $error], false, -1], ["MarshalJSON", "MarshalJSON", "", [], [($sliceType($Uint8)), $error], false, -1], ["MarshalText", "MarshalText", "", [], [($sliceType($Uint8)), $error], false, -1], ["Minute", "Minute", "", [], [$Int], false, -1], ["Month", "Month", "", [], [Month], false, -1], ["Nanosecond", "Nanosecond", "", [], [$Int], false, -1], ["Round", "Round", "", [Duration], [Time], false, -1], ["Second", "Second", "", [], [$Int], false, -1], ["String", "String", "", [], [$String], false, -1], ["Sub", "Sub", "", [Time], [Duration], false, -1], ["Truncate", "Truncate", "", [Duration], [Time], false, -1], ["UTC", "UTC", "", [], [Time], false, -1], ["Unix", "Unix", "", [], [$Int64], false, -1], ["UnixNano", "UnixNano", "", [], [$Int64], false, -1], ["Weekday", "Weekday", "", [], [Weekday], false, -1], ["Year", "Year", "", [], [$Int], false, -1], ["YearDay", "YearDay", "", [], [$Int], false, -1], ["Zone", "Zone", "", [], [$String, $Int], false, -1], ["abs", "abs", "time", [], [$Uint64], false, -1], ["date", "date", "time", [$Bool], [$Int, Month, $Int, $Int], false, -1], ["locabs", "locabs", "time", [], [$String, $Int, $Uint64], false, -1]];
		($ptrType(Time)).methods = [["Add", "Add", "", [Duration], [Time], false, -1], ["AddDate", "AddDate", "", [$Int, $Int, $Int], [Time], false, -1], ["After", "After", "", [Time], [$Bool], false, -1], ["Before", "Before", "", [Time], [$Bool], false, -1], ["Clock", "Clock", "", [], [$Int, $Int, $Int], false, -1], ["Date", "Date", "", [], [$Int, Month, $Int], false, -1], ["Day", "Day", "", [], [$Int], false, -1], ["Equal", "Equal", "", [Time], [$Bool], false, -1], ["Format", "Format", "", [$String], [$String], false, -1], ["GobDecode", "GobDecode", "", [($sliceType($Uint8))], [$error], false, -1], ["GobEncode", "GobEncode", "", [], [($sliceType($Uint8)), $error], false, -1], ["Hour", "Hour", "", [], [$Int], false, -1], ["ISOWeek", "ISOWeek", "", [], [$Int, $Int], false, -1], ["In", "In", "", [($ptrType(Location))], [Time], false, -1], ["IsZero", "IsZero", "", [], [$Bool], false, -1], ["Local", "Local", "", [], [Time], false, -1], ["Location", "Location", "", [], [($ptrType(Location))], false, -1], ["MarshalBinary", "MarshalBinary", "", [], [($sliceType($Uint8)), $error], false, -1], ["MarshalJSON", "MarshalJSON", "", [], [($sliceType($Uint8)), $error], false, -1], ["MarshalText", "MarshalText", "", [], [($sliceType($Uint8)), $error], false, -1], ["Minute", "Minute", "", [], [$Int], false, -1], ["Month", "Month", "", [], [Month], false, -1], ["Nanosecond", "Nanosecond", "", [], [$Int], false, -1], ["Round", "Round", "", [Duration], [Time], false, -1], ["Second", "Second", "", [], [$Int], false, -1], ["String", "String", "", [], [$String], false, -1], ["Sub", "Sub", "", [Time], [Duration], false, -1], ["Truncate", "Truncate", "", [Duration], [Time], false, -1], ["UTC", "UTC", "", [], [Time], false, -1], ["Unix", "Unix", "", [], [$Int64], false, -1], ["UnixNano", "UnixNano", "", [], [$Int64], false, -1], ["UnmarshalBinary", "UnmarshalBinary", "", [($sliceType($Uint8))], [$error], false, -1], ["UnmarshalJSON", "UnmarshalJSON", "", [($sliceType($Uint8))], [$error], false, -1], ["UnmarshalText", "UnmarshalText", "", [($sliceType($Uint8))], [$error], false, -1], ["Weekday", "Weekday", "", [], [Weekday], false, -1], ["Year", "Year", "", [], [$Int], false, -1], ["YearDay", "YearDay", "", [], [$Int], false, -1], ["Zone", "Zone", "", [], [$String, $Int], false, -1], ["abs", "abs", "time", [], [$Uint64], false, -1], ["date", "date", "time", [$Bool], [$Int, Month, $Int, $Int], false, -1], ["locabs", "locabs", "time", [], [$String, $Int, $Uint64], false, -1]];
		Time.init([["sec", "sec", "time", $Int64, ""], ["nsec", "nsec", "time", $Uintptr, ""], ["loc", "loc", "time", ($ptrType(Location)), ""]]);
		Month.methods = [["String", "String", "", [], [$String], false, -1]];
		($ptrType(Month)).methods = [["String", "String", "", [], [$String], false, -1]];
		Weekday.methods = [["String", "String", "", [], [$String], false, -1]];
		($ptrType(Weekday)).methods = [["String", "String", "", [], [$String], false, -1]];
		Duration.methods = [["Hours", "Hours", "", [], [$Float64], false, -1], ["Minutes", "Minutes", "", [], [$Float64], false, -1], ["Nanoseconds", "Nanoseconds", "", [], [$Int64], false, -1], ["Seconds", "Seconds", "", [], [$Float64], false, -1], ["String", "String", "", [], [$String], false, -1]];
		($ptrType(Duration)).methods = [["Hours", "Hours", "", [], [$Float64], false, -1], ["Minutes", "Minutes", "", [], [$Float64], false, -1], ["Nanoseconds", "Nanoseconds", "", [], [$Int64], false, -1], ["Seconds", "Seconds", "", [], [$Float64], false, -1], ["String", "String", "", [], [$String], false, -1]];
		($ptrType(Location)).methods = [["String", "String", "", [], [$String], false, -1], ["firstZoneUsed", "firstZoneUsed", "time", [], [$Bool], false, -1], ["get", "get", "time", [], [($ptrType(Location))], false, -1], ["lookup", "lookup", "time", [$Int64], [$String, $Int, $Bool, $Int64, $Int64], false, -1], ["lookupFirstZone", "lookupFirstZone", "time", [], [$Int], false, -1], ["lookupName", "lookupName", "time", [$String, $Int64], [$Int, $Bool, $Bool], false, -1]];
		Location.init([["name", "name", "time", $String, ""], ["zone", "zone", "time", ($sliceType(zone)), ""], ["tx", "tx", "time", ($sliceType(zoneTrans)), ""], ["cacheStart", "cacheStart", "time", $Int64, ""], ["cacheEnd", "cacheEnd", "time", $Int64, ""], ["cacheZone", "cacheZone", "time", ($ptrType(zone)), ""]]);
		zone.init([["name", "name", "time", $String, ""], ["offset", "offset", "time", $Int, ""], ["isDST", "isDST", "time", $Bool, ""]]);
		zoneTrans.init([["when", "when", "time", $Int64, ""], ["index", "index", "time", $Uint8, ""], ["isstd", "isstd", "time", $Bool, ""], ["isutc", "isutc", "time", $Bool, ""]]);
		localLoc = new Location.Ptr();
		localOnce = new sync.Once.Ptr();
		std0x = $toNativeArray("Int", [260, 265, 524, 526, 528, 274]);
		longDayNames = new ($sliceType($String))(["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]);
		shortDayNames = new ($sliceType($String))(["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"]);
		shortMonthNames = new ($sliceType($String))(["---", "Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"]);
		longMonthNames = new ($sliceType($String))(["---", "January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]);
		atoiError = errors.New("time: invalid number");
		errBad = errors.New("bad value for field");
		errLeadingInt = errors.New("time: bad [0-9]*");
		months = $toNativeArray("String", ["January", "February", "March", "April", "May", "June", "July", "August", "September", "October", "November", "December"]);
		days = $toNativeArray("String", ["Sunday", "Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday"]);
		daysBefore = $toNativeArray("Int32", [0, 31, 59, 90, 120, 151, 181, 212, 243, 273, 304, 334, 365]);
		utcLoc = new Location.Ptr("UTC", ($sliceType(zone)).nil, ($sliceType(zoneTrans)).nil, new $Int64(0, 0), new $Int64(0, 0), ($ptrType(zone)).nil);
		$pkg.UTC = utcLoc;
		$pkg.Local = localLoc;
		_tuple = syscall.Getenv("ZONEINFO"); zoneinfo = _tuple[0];
		badData = errors.New("malformed time zone information");
		zoneDirs = new ($sliceType($String))(["/usr/share/zoneinfo/", "/usr/share/lib/zoneinfo/", "/usr/lib/locale/TZ/", runtime.GOROOT() + "/lib/time/zoneinfo.zip"]);
	};
	return $pkg;
})();
$packages["os"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], io = $packages["io"], syscall = $packages["syscall"], time = $packages["time"], errors = $packages["errors"], runtime = $packages["runtime"], atomic = $packages["sync/atomic"], sync = $packages["sync"], PathError, SyscallError, LinkError, File, file, dirInfo, FileInfo, FileMode, fileStat, lstat, init, Getenv, NewSyscallError, IsNotExist, isNotExist, Open, sigpipe, syscallMode, NewFile, epipecheck, OpenFile, Lstat, basename, fileInfoFromStat, timespecToTime;
	PathError = $pkg.PathError = $newType(0, "Struct", "os.PathError", "PathError", "os", function(Op_, Path_, Err_) {
		this.$val = this;
		this.Op = Op_ !== undefined ? Op_ : "";
		this.Path = Path_ !== undefined ? Path_ : "";
		this.Err = Err_ !== undefined ? Err_ : null;
	});
	SyscallError = $pkg.SyscallError = $newType(0, "Struct", "os.SyscallError", "SyscallError", "os", function(Syscall_, Err_) {
		this.$val = this;
		this.Syscall = Syscall_ !== undefined ? Syscall_ : "";
		this.Err = Err_ !== undefined ? Err_ : null;
	});
	LinkError = $pkg.LinkError = $newType(0, "Struct", "os.LinkError", "LinkError", "os", function(Op_, Old_, New_, Err_) {
		this.$val = this;
		this.Op = Op_ !== undefined ? Op_ : "";
		this.Old = Old_ !== undefined ? Old_ : "";
		this.New = New_ !== undefined ? New_ : "";
		this.Err = Err_ !== undefined ? Err_ : null;
	});
	File = $pkg.File = $newType(0, "Struct", "os.File", "File", "os", function(file_) {
		this.$val = this;
		this.file = file_ !== undefined ? file_ : ($ptrType(file)).nil;
	});
	file = $pkg.file = $newType(0, "Struct", "os.file", "file", "os", function(fd_, name_, dirinfo_, nepipe_) {
		this.$val = this;
		this.fd = fd_ !== undefined ? fd_ : 0;
		this.name = name_ !== undefined ? name_ : "";
		this.dirinfo = dirinfo_ !== undefined ? dirinfo_ : ($ptrType(dirInfo)).nil;
		this.nepipe = nepipe_ !== undefined ? nepipe_ : 0;
	});
	dirInfo = $pkg.dirInfo = $newType(0, "Struct", "os.dirInfo", "dirInfo", "os", function(buf_, nbuf_, bufp_) {
		this.$val = this;
		this.buf = buf_ !== undefined ? buf_ : ($sliceType($Uint8)).nil;
		this.nbuf = nbuf_ !== undefined ? nbuf_ : 0;
		this.bufp = bufp_ !== undefined ? bufp_ : 0;
	});
	FileInfo = $pkg.FileInfo = $newType(8, "Interface", "os.FileInfo", "FileInfo", "os", null);
	FileMode = $pkg.FileMode = $newType(4, "Uint32", "os.FileMode", "FileMode", "os", null);
	fileStat = $pkg.fileStat = $newType(0, "Struct", "os.fileStat", "fileStat", "os", function(name_, size_, mode_, modTime_, sys_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : "";
		this.size = size_ !== undefined ? size_ : new $Int64(0, 0);
		this.mode = mode_ !== undefined ? mode_ : 0;
		this.modTime = modTime_ !== undefined ? modTime_ : new time.Time.Ptr();
		this.sys = sys_ !== undefined ? sys_ : null;
	});
	init = function() {
		var process, args, i;
		process = $global.process;
		if (process === undefined) {
			$pkg.Args = new ($sliceType($String))(["browser"]);
			return;
		}
		args = process.argv;
		$pkg.Args = ($sliceType($String)).make(($parseInt(args.length) - 1 >> 0));
		i = 0;
		while (i < ($parseInt(args.length) - 1 >> 0)) {
			(i < 0 || i >= $pkg.Args.$length) ? $throwRuntimeError("index out of range") : $pkg.Args.$array[$pkg.Args.$offset + i] = $internalize(args[(i + 1 >> 0)], $String);
			i = i + (1) >> 0;
		}
	};
	File.Ptr.prototype.readdirnames = function(n) {
		var names = ($sliceType($String)).nil, err = null, f, d, size, errno, _tuple, _tmp, _tmp$1, _tmp$2, _tmp$3, nb, nc, _tuple$1, _tmp$4, _tmp$5, _tmp$6, _tmp$7;
		f = this;
		if (f.file.dirinfo === ($ptrType(dirInfo)).nil) {
			f.file.dirinfo = new dirInfo.Ptr();
			f.file.dirinfo.buf = ($sliceType($Uint8)).make(4096);
		}
		d = f.file.dirinfo;
		size = n;
		if (size <= 0) {
			size = 100;
			n = -1;
		}
		names = ($sliceType($String)).make(0, size);
		while (!((n === 0))) {
			if (d.bufp >= d.nbuf) {
				d.bufp = 0;
				errno = null;
				_tuple = syscall.ReadDirent(f.file.fd, d.buf); d.nbuf = _tuple[0]; errno = _tuple[1];
				if (!($interfaceIsEqual(errno, null))) {
					_tmp = names; _tmp$1 = NewSyscallError("readdirent", errno); names = _tmp; err = _tmp$1;
					return [names, err];
				}
				if (d.nbuf <= 0) {
					break;
				}
			}
			_tmp$2 = 0; _tmp$3 = 0; nb = _tmp$2; nc = _tmp$3;
			_tuple$1 = syscall.ParseDirent($subslice(d.buf, d.bufp, d.nbuf), n, names); nb = _tuple$1[0]; nc = _tuple$1[1]; names = _tuple$1[2];
			d.bufp = d.bufp + (nb) >> 0;
			n = n - (nc) >> 0;
		}
		if (n >= 0 && (names.$length === 0)) {
			_tmp$4 = names; _tmp$5 = io.EOF; names = _tmp$4; err = _tmp$5;
			return [names, err];
		}
		_tmp$6 = names; _tmp$7 = null; names = _tmp$6; err = _tmp$7;
		return [names, err];
	};
	File.prototype.readdirnames = function(n) { return this.$val.readdirnames(n); };
	File.Ptr.prototype.Readdir = function(n) {
		var fi = ($sliceType(FileInfo)).nil, err = null, f, _tmp, _tmp$1, _tuple;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = ($sliceType(FileInfo)).nil; _tmp$1 = $pkg.ErrInvalid; fi = _tmp; err = _tmp$1;
			return [fi, err];
		}
		_tuple = f.readdir(n); fi = _tuple[0]; err = _tuple[1];
		return [fi, err];
	};
	File.prototype.Readdir = function(n) { return this.$val.Readdir(n); };
	File.Ptr.prototype.Readdirnames = function(n) {
		var names = ($sliceType($String)).nil, err = null, f, _tmp, _tmp$1, _tuple;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = ($sliceType($String)).nil; _tmp$1 = $pkg.ErrInvalid; names = _tmp; err = _tmp$1;
			return [names, err];
		}
		_tuple = f.readdirnames(n); names = _tuple[0]; err = _tuple[1];
		return [names, err];
	};
	File.prototype.Readdirnames = function(n) { return this.$val.Readdirnames(n); };
	Getenv = $pkg.Getenv = function(key) {
		var _tuple, v;
		_tuple = syscall.Getenv(key); v = _tuple[0];
		return v;
	};
	PathError.Ptr.prototype.Error = function() {
		var e;
		e = this;
		return e.Op + " " + e.Path + ": " + e.Err.Error();
	};
	PathError.prototype.Error = function() { return this.$val.Error(); };
	SyscallError.Ptr.prototype.Error = function() {
		var e;
		e = this;
		return e.Syscall + ": " + e.Err.Error();
	};
	SyscallError.prototype.Error = function() { return this.$val.Error(); };
	NewSyscallError = $pkg.NewSyscallError = function(syscall$1, err) {
		if ($interfaceIsEqual(err, null)) {
			return null;
		}
		return new SyscallError.Ptr(syscall$1, err);
	};
	IsNotExist = $pkg.IsNotExist = function(err) {
		return isNotExist(err);
	};
	isNotExist = function(err) {
		var pe, _ref, _type;
		_ref = err;
		_type = _ref !== null ? _ref.constructor : null;
		if (_type === null) {
			pe = _ref;
			return false;
		} else if (_type === ($ptrType(PathError))) {
			pe = _ref.$val;
			err = pe.Err;
		} else if (_type === ($ptrType(LinkError))) {
			pe = _ref.$val;
			err = pe.Err;
		}
		return $interfaceIsEqual(err, new syscall.Errno(2)) || $interfaceIsEqual(err, $pkg.ErrNotExist);
	};
	File.Ptr.prototype.Name = function() {
		var f;
		f = this;
		return f.file.name;
	};
	File.prototype.Name = function() { return this.$val.Name(); };
	LinkError.Ptr.prototype.Error = function() {
		var e;
		e = this;
		return e.Op + " " + e.Old + " " + e.New + ": " + e.Err.Error();
	};
	LinkError.prototype.Error = function() { return this.$val.Error(); };
	File.Ptr.prototype.Read = function(b) {
		var n = 0, err = null, f, _tmp, _tmp$1, _tuple, e, _tmp$2, _tmp$3, _tmp$4, _tmp$5;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = 0; _tmp$1 = $pkg.ErrInvalid; n = _tmp; err = _tmp$1;
			return [n, err];
		}
		_tuple = f.read(b); n = _tuple[0]; e = _tuple[1];
		if (n < 0) {
			n = 0;
		}
		if ((n === 0) && b.$length > 0 && $interfaceIsEqual(e, null)) {
			_tmp$2 = 0; _tmp$3 = io.EOF; n = _tmp$2; err = _tmp$3;
			return [n, err];
		}
		if (!($interfaceIsEqual(e, null))) {
			err = new PathError.Ptr("read", f.file.name, e);
		}
		_tmp$4 = n; _tmp$5 = err; n = _tmp$4; err = _tmp$5;
		return [n, err];
	};
	File.prototype.Read = function(b) { return this.$val.Read(b); };
	File.Ptr.prototype.ReadAt = function(b, off) {
		var n = 0, err = null, f, _tmp, _tmp$1, _tuple, m, e, _tmp$2, _tmp$3, x;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = 0; _tmp$1 = $pkg.ErrInvalid; n = _tmp; err = _tmp$1;
			return [n, err];
		}
		while (b.$length > 0) {
			_tuple = f.pread(b, off); m = _tuple[0]; e = _tuple[1];
			if ((m === 0) && $interfaceIsEqual(e, null)) {
				_tmp$2 = n; _tmp$3 = io.EOF; n = _tmp$2; err = _tmp$3;
				return [n, err];
			}
			if (!($interfaceIsEqual(e, null))) {
				err = new PathError.Ptr("read", f.file.name, e);
				break;
			}
			n = n + (m) >> 0;
			b = $subslice(b, m);
			off = (x = new $Int64(0, m), new $Int64(off.$high + x.$high, off.$low + x.$low));
		}
		return [n, err];
	};
	File.prototype.ReadAt = function(b, off) { return this.$val.ReadAt(b, off); };
	File.Ptr.prototype.Write = function(b) {
		var n = 0, err = null, f, _tmp, _tmp$1, _tuple, e, _tmp$2, _tmp$3;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = 0; _tmp$1 = $pkg.ErrInvalid; n = _tmp; err = _tmp$1;
			return [n, err];
		}
		_tuple = f.write(b); n = _tuple[0]; e = _tuple[1];
		if (n < 0) {
			n = 0;
		}
		if (!((n === b.$length))) {
			err = io.ErrShortWrite;
		}
		epipecheck(f, e);
		if (!($interfaceIsEqual(e, null))) {
			err = new PathError.Ptr("write", f.file.name, e);
		}
		_tmp$2 = n; _tmp$3 = err; n = _tmp$2; err = _tmp$3;
		return [n, err];
	};
	File.prototype.Write = function(b) { return this.$val.Write(b); };
	File.Ptr.prototype.WriteAt = function(b, off) {
		var n = 0, err = null, f, _tmp, _tmp$1, _tuple, m, e, x;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = 0; _tmp$1 = $pkg.ErrInvalid; n = _tmp; err = _tmp$1;
			return [n, err];
		}
		while (b.$length > 0) {
			_tuple = f.pwrite(b, off); m = _tuple[0]; e = _tuple[1];
			if (!($interfaceIsEqual(e, null))) {
				err = new PathError.Ptr("write", f.file.name, e);
				break;
			}
			n = n + (m) >> 0;
			b = $subslice(b, m);
			off = (x = new $Int64(0, m), new $Int64(off.$high + x.$high, off.$low + x.$low));
		}
		return [n, err];
	};
	File.prototype.WriteAt = function(b, off) { return this.$val.WriteAt(b, off); };
	File.Ptr.prototype.Seek = function(offset, whence) {
		var ret = new $Int64(0, 0), err = null, f, _tmp, _tmp$1, _tuple, r, e, _tmp$2, _tmp$3, _tmp$4, _tmp$5;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = new $Int64(0, 0); _tmp$1 = $pkg.ErrInvalid; ret = _tmp; err = _tmp$1;
			return [ret, err];
		}
		_tuple = f.seek(offset, whence); r = _tuple[0]; e = _tuple[1];
		if ($interfaceIsEqual(e, null) && !(f.file.dirinfo === ($ptrType(dirInfo)).nil) && !((r.$high === 0 && r.$low === 0))) {
			e = new syscall.Errno(21);
		}
		if (!($interfaceIsEqual(e, null))) {
			_tmp$2 = new $Int64(0, 0); _tmp$3 = new PathError.Ptr("seek", f.file.name, e); ret = _tmp$2; err = _tmp$3;
			return [ret, err];
		}
		_tmp$4 = r; _tmp$5 = null; ret = _tmp$4; err = _tmp$5;
		return [ret, err];
	};
	File.prototype.Seek = function(offset, whence) { return this.$val.Seek(offset, whence); };
	File.Ptr.prototype.WriteString = function(s) {
		var ret = 0, err = null, f, _tmp, _tmp$1, _tuple;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = 0; _tmp$1 = $pkg.ErrInvalid; ret = _tmp; err = _tmp$1;
			return [ret, err];
		}
		_tuple = f.Write(new ($sliceType($Uint8))($stringToBytes(s))); ret = _tuple[0]; err = _tuple[1];
		return [ret, err];
	};
	File.prototype.WriteString = function(s) { return this.$val.WriteString(s); };
	File.Ptr.prototype.Chdir = function() {
		var f, e;
		f = this;
		if (f === ($ptrType(File)).nil) {
			return $pkg.ErrInvalid;
		}
		e = syscall.Fchdir(f.file.fd);
		if (!($interfaceIsEqual(e, null))) {
			return new PathError.Ptr("chdir", f.file.name, e);
		}
		return null;
	};
	File.prototype.Chdir = function() { return this.$val.Chdir(); };
	Open = $pkg.Open = function(name) {
		var file$1 = ($ptrType(File)).nil, err = null, _tuple;
		_tuple = OpenFile(name, 0, 0); file$1 = _tuple[0]; err = _tuple[1];
		return [file$1, err];
	};
	sigpipe = function() {
		$panic("Native function not implemented: os.sigpipe");
	};
	syscallMode = function(i) {
		var o = 0;
		o = (o | (((new FileMode(i)).Perm() >>> 0))) >>> 0;
		if (!((((i & 8388608) >>> 0) === 0))) {
			o = (o | (2048)) >>> 0;
		}
		if (!((((i & 4194304) >>> 0) === 0))) {
			o = (o | (1024)) >>> 0;
		}
		if (!((((i & 1048576) >>> 0) === 0))) {
			o = (o | (512)) >>> 0;
		}
		return o;
	};
	File.Ptr.prototype.Chmod = function(mode) {
		var f, e;
		f = this;
		if (f === ($ptrType(File)).nil) {
			return $pkg.ErrInvalid;
		}
		e = syscall.Fchmod(f.file.fd, syscallMode(mode));
		if (!($interfaceIsEqual(e, null))) {
			return new PathError.Ptr("chmod", f.file.name, e);
		}
		return null;
	};
	File.prototype.Chmod = function(mode) { return this.$val.Chmod(mode); };
	File.Ptr.prototype.Chown = function(uid, gid) {
		var f, e;
		f = this;
		if (f === ($ptrType(File)).nil) {
			return $pkg.ErrInvalid;
		}
		e = syscall.Fchown(f.file.fd, uid, gid);
		if (!($interfaceIsEqual(e, null))) {
			return new PathError.Ptr("chown", f.file.name, e);
		}
		return null;
	};
	File.prototype.Chown = function(uid, gid) { return this.$val.Chown(uid, gid); };
	File.Ptr.prototype.Truncate = function(size) {
		var f, e;
		f = this;
		if (f === ($ptrType(File)).nil) {
			return $pkg.ErrInvalid;
		}
		e = syscall.Ftruncate(f.file.fd, size);
		if (!($interfaceIsEqual(e, null))) {
			return new PathError.Ptr("truncate", f.file.name, e);
		}
		return null;
	};
	File.prototype.Truncate = function(size) { return this.$val.Truncate(size); };
	File.Ptr.prototype.Sync = function() {
		var err = null, f, e;
		f = this;
		if (f === ($ptrType(File)).nil) {
			err = $pkg.ErrInvalid;
			return err;
		}
		e = syscall.Fsync(f.file.fd);
		if (!($interfaceIsEqual(e, null))) {
			err = NewSyscallError("fsync", e);
			return err;
		}
		err = null;
		return err;
	};
	File.prototype.Sync = function() { return this.$val.Sync(); };
	File.Ptr.prototype.Fd = function() {
		var f;
		f = this;
		if (f === ($ptrType(File)).nil) {
			return 4294967295;
		}
		return (f.file.fd >>> 0);
	};
	File.prototype.Fd = function() { return this.$val.Fd(); };
	NewFile = $pkg.NewFile = function(fd, name) {
		var fdi, f;
		fdi = (fd >> 0);
		if (fdi < 0) {
			return ($ptrType(File)).nil;
		}
		f = new File.Ptr(new file.Ptr(fdi, name, ($ptrType(dirInfo)).nil, 0));
		runtime.SetFinalizer(f.file, new ($funcType([($ptrType(file))], [$error], false))((function(recv) { $stackDepthOffset--; try { return recv.close(); } finally { $stackDepthOffset++; } })));
		return f;
	};
	epipecheck = function(file$1, e) {
		if ($interfaceIsEqual(e, new syscall.Errno(32))) {
			if (atomic.AddInt32(new ($ptrType($Int32))(function() { return this.$target.file.nepipe; }, function($v) { this.$target.file.nepipe = $v; }, file$1), 1) >= 10) {
				sigpipe();
			}
		} else {
			atomic.StoreInt32(new ($ptrType($Int32))(function() { return this.$target.file.nepipe; }, function($v) { this.$target.file.nepipe = $v; }, file$1), 0);
		}
	};
	OpenFile = $pkg.OpenFile = function(name, flag, perm) {
		var file$1 = ($ptrType(File)).nil, err = null, _tuple, r, e, _tmp, _tmp$1, _tmp$2, _tmp$3;
		_tuple = syscall.Open(name, flag | 524288, syscallMode(perm)); r = _tuple[0]; e = _tuple[1];
		if (!($interfaceIsEqual(e, null))) {
			_tmp = ($ptrType(File)).nil; _tmp$1 = new PathError.Ptr("open", name, e); file$1 = _tmp; err = _tmp$1;
			return [file$1, err];
		}
		_tmp$2 = NewFile((r >>> 0), name); _tmp$3 = null; file$1 = _tmp$2; err = _tmp$3;
		return [file$1, err];
	};
	File.Ptr.prototype.Close = function() {
		var f;
		f = this;
		if (f === ($ptrType(File)).nil) {
			return $pkg.ErrInvalid;
		}
		return f.file.close();
	};
	File.prototype.Close = function() { return this.$val.Close(); };
	file.Ptr.prototype.close = function() {
		var file$1, err, e;
		file$1 = this;
		if (file$1 === ($ptrType(file)).nil || file$1.fd < 0) {
			return new syscall.Errno(22);
		}
		err = null;
		e = syscall.Close(file$1.fd);
		if (!($interfaceIsEqual(e, null))) {
			err = new PathError.Ptr("close", file$1.name, e);
		}
		file$1.fd = -1;
		runtime.SetFinalizer(file$1, null);
		return err;
	};
	file.prototype.close = function() { return this.$val.close(); };
	File.Ptr.prototype.Stat = function() {
		var fi = null, err = null, f, _tmp, _tmp$1, stat, _tmp$2, _tmp$3, _tmp$4, _tmp$5;
		f = this;
		if (f === ($ptrType(File)).nil) {
			_tmp = null; _tmp$1 = $pkg.ErrInvalid; fi = _tmp; err = _tmp$1;
			return [fi, err];
		}
		stat = new syscall.Stat_t.Ptr(); $copy(stat, new syscall.Stat_t.Ptr(), syscall.Stat_t);
		err = syscall.Fstat(f.file.fd, stat);
		if (!($interfaceIsEqual(err, null))) {
			_tmp$2 = null; _tmp$3 = new PathError.Ptr("stat", f.file.name, err); fi = _tmp$2; err = _tmp$3;
			return [fi, err];
		}
		_tmp$4 = fileInfoFromStat(stat, f.file.name); _tmp$5 = null; fi = _tmp$4; err = _tmp$5;
		return [fi, err];
	};
	File.prototype.Stat = function() { return this.$val.Stat(); };
	Lstat = $pkg.Lstat = function(name) {
		var fi = null, err = null, stat, _tmp, _tmp$1, _tmp$2, _tmp$3;
		stat = new syscall.Stat_t.Ptr(); $copy(stat, new syscall.Stat_t.Ptr(), syscall.Stat_t);
		err = syscall.Lstat(name, stat);
		if (!($interfaceIsEqual(err, null))) {
			_tmp = null; _tmp$1 = new PathError.Ptr("lstat", name, err); fi = _tmp; err = _tmp$1;
			return [fi, err];
		}
		_tmp$2 = fileInfoFromStat(stat, name); _tmp$3 = null; fi = _tmp$2; err = _tmp$3;
		return [fi, err];
	};
	File.Ptr.prototype.readdir = function(n) {
		var fi = ($sliceType(FileInfo)).nil, err = null, f, dirname, _tuple, names, _ref, _i, filename, _tuple$1, fip, lerr, _tmp, _tmp$1, _tmp$2, _tmp$3;
		f = this;
		dirname = f.file.name;
		if (dirname === "") {
			dirname = ".";
		}
		_tuple = f.Readdirnames(n); names = _tuple[0]; err = _tuple[1];
		fi = ($sliceType(FileInfo)).make(0, names.$length);
		_ref = names;
		_i = 0;
		while (_i < _ref.$length) {
			filename = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			_tuple$1 = lstat(dirname + "/" + filename); fip = _tuple$1[0]; lerr = _tuple$1[1];
			if (IsNotExist(lerr)) {
				_i++;
				continue;
			}
			if (!($interfaceIsEqual(lerr, null))) {
				_tmp = fi; _tmp$1 = lerr; fi = _tmp; err = _tmp$1;
				return [fi, err];
			}
			fi = $append(fi, fip);
			_i++;
		}
		_tmp$2 = fi; _tmp$3 = err; fi = _tmp$2; err = _tmp$3;
		return [fi, err];
	};
	File.prototype.readdir = function(n) { return this.$val.readdir(n); };
	File.Ptr.prototype.read = function(b) {
		var n = 0, err = null, f, _tuple;
		f = this;
		_tuple = syscall.Read(f.file.fd, b); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	File.prototype.read = function(b) { return this.$val.read(b); };
	File.Ptr.prototype.pread = function(b, off) {
		var n = 0, err = null, f, _tuple;
		f = this;
		_tuple = syscall.Pread(f.file.fd, b, off); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	File.prototype.pread = function(b, off) { return this.$val.pread(b, off); };
	File.Ptr.prototype.write = function(b) {
		var n = 0, err = null, f, bcap, _tuple, m, err$1, _tmp, _tmp$1;
		f = this;
		while (true) {
			bcap = b;
			_tuple = syscall.Write(f.file.fd, bcap); m = _tuple[0]; err$1 = _tuple[1];
			n = n + (m) >> 0;
			if (0 < m && m < bcap.$length || $interfaceIsEqual(err$1, new syscall.Errno(4))) {
				b = $subslice(b, m);
				continue;
			}
			_tmp = n; _tmp$1 = err$1; n = _tmp; err = _tmp$1;
			return [n, err];
		}
	};
	File.prototype.write = function(b) { return this.$val.write(b); };
	File.Ptr.prototype.pwrite = function(b, off) {
		var n = 0, err = null, f, _tuple;
		f = this;
		_tuple = syscall.Pwrite(f.file.fd, b, off); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	File.prototype.pwrite = function(b, off) { return this.$val.pwrite(b, off); };
	File.Ptr.prototype.seek = function(offset, whence) {
		var ret = new $Int64(0, 0), err = null, f, _tuple;
		f = this;
		_tuple = syscall.Seek(f.file.fd, offset, whence); ret = _tuple[0]; err = _tuple[1];
		return [ret, err];
	};
	File.prototype.seek = function(offset, whence) { return this.$val.seek(offset, whence); };
	basename = function(name) {
		var i;
		i = name.length - 1 >> 0;
		while (i > 0 && (name.charCodeAt(i) === 47)) {
			name = name.substring(0, i);
			i = i - (1) >> 0;
		}
		i = i - (1) >> 0;
		while (i >= 0) {
			if (name.charCodeAt(i) === 47) {
				name = name.substring((i + 1 >> 0));
				break;
			}
			i = i - (1) >> 0;
		}
		return name;
	};
	fileInfoFromStat = function(st, name) {
		var fs, _ref;
		fs = new fileStat.Ptr(basename(name), st.Size, 0, timespecToTime($clone(st.Mtim, syscall.Timespec)), st);
		fs.mode = (((st.Mode & 511) >>> 0) >>> 0);
		_ref = (st.Mode & 61440) >>> 0;
		if (_ref === 24576) {
			fs.mode = (fs.mode | (67108864)) >>> 0;
		} else if (_ref === 8192) {
			fs.mode = (fs.mode | (69206016)) >>> 0;
		} else if (_ref === 16384) {
			fs.mode = (fs.mode | (2147483648)) >>> 0;
		} else if (_ref === 4096) {
			fs.mode = (fs.mode | (33554432)) >>> 0;
		} else if (_ref === 40960) {
			fs.mode = (fs.mode | (134217728)) >>> 0;
		} else if (_ref === 32768) {
		} else if (_ref === 49152) {
			fs.mode = (fs.mode | (16777216)) >>> 0;
		}
		if (!((((st.Mode & 1024) >>> 0) === 0))) {
			fs.mode = (fs.mode | (4194304)) >>> 0;
		}
		if (!((((st.Mode & 2048) >>> 0) === 0))) {
			fs.mode = (fs.mode | (8388608)) >>> 0;
		}
		if (!((((st.Mode & 512) >>> 0) === 0))) {
			fs.mode = (fs.mode | (1048576)) >>> 0;
		}
		return fs;
	};
	timespecToTime = function(ts) {
		return time.Unix(ts.Sec, ts.Nsec);
	};
	FileMode.prototype.String = function() {
		var m, buf, w, _ref, _i, _rune, i, c, y, _ref$1, _i$1, _rune$1, i$1, c$1, y$1;
		m = this.$val;
		buf = ($arrayType($Uint8, 32)).zero(); $copy(buf, ($arrayType($Uint8, 32)).zero(), ($arrayType($Uint8, 32)));
		w = 0;
		_ref = "dalTLDpSugct";
		_i = 0;
		while (_i < _ref.length) {
			_rune = $decodeRune(_ref, _i);
			i = _i;
			c = _rune[0];
			if (!((((m & (((y = ((31 - i >> 0) >>> 0), y < 32 ? (1 << y) : 0) >>> 0))) >>> 0) === 0))) {
				(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = (c << 24 >>> 24);
				w = w + (1) >> 0;
			}
			_i += _rune[1];
		}
		if (w === 0) {
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 45;
			w = w + (1) >> 0;
		}
		_ref$1 = "rwxrwxrwx";
		_i$1 = 0;
		while (_i$1 < _ref$1.length) {
			_rune$1 = $decodeRune(_ref$1, _i$1);
			i$1 = _i$1;
			c$1 = _rune$1[0];
			if (!((((m & (((y$1 = ((8 - i$1 >> 0) >>> 0), y$1 < 32 ? (1 << y$1) : 0) >>> 0))) >>> 0) === 0))) {
				(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = (c$1 << 24 >>> 24);
			} else {
				(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 45;
			}
			w = w + (1) >> 0;
			_i$1 += _rune$1[1];
		}
		return $bytesToString($subslice(new ($sliceType($Uint8))(buf), 0, w));
	};
	$ptrType(FileMode).prototype.String = function() { return new FileMode(this.$get()).String(); };
	FileMode.prototype.IsDir = function() {
		var m;
		m = this.$val;
		return !((((m & 2147483648) >>> 0) === 0));
	};
	$ptrType(FileMode).prototype.IsDir = function() { return new FileMode(this.$get()).IsDir(); };
	FileMode.prototype.IsRegular = function() {
		var m;
		m = this.$val;
		return ((m & 2399141888) >>> 0) === 0;
	};
	$ptrType(FileMode).prototype.IsRegular = function() { return new FileMode(this.$get()).IsRegular(); };
	FileMode.prototype.Perm = function() {
		var m;
		m = this.$val;
		return (m & 511) >>> 0;
	};
	$ptrType(FileMode).prototype.Perm = function() { return new FileMode(this.$get()).Perm(); };
	fileStat.Ptr.prototype.Name = function() {
		var fs;
		fs = this;
		return fs.name;
	};
	fileStat.prototype.Name = function() { return this.$val.Name(); };
	fileStat.Ptr.prototype.IsDir = function() {
		var fs;
		fs = this;
		return (new FileMode(fs.Mode())).IsDir();
	};
	fileStat.prototype.IsDir = function() { return this.$val.IsDir(); };
	fileStat.Ptr.prototype.Size = function() {
		var fs;
		fs = this;
		return fs.size;
	};
	fileStat.prototype.Size = function() { return this.$val.Size(); };
	fileStat.Ptr.prototype.Mode = function() {
		var fs;
		fs = this;
		return fs.mode;
	};
	fileStat.prototype.Mode = function() { return this.$val.Mode(); };
	fileStat.Ptr.prototype.ModTime = function() {
		var fs;
		fs = this;
		return fs.modTime;
	};
	fileStat.prototype.ModTime = function() { return this.$val.ModTime(); };
	fileStat.Ptr.prototype.Sys = function() {
		var fs;
		fs = this;
		return fs.sys;
	};
	fileStat.prototype.Sys = function() { return this.$val.Sys(); };
	$pkg.$init = function() {
		($ptrType(PathError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		PathError.init([["Op", "Op", "", $String, ""], ["Path", "Path", "", $String, ""], ["Err", "Err", "", $error, ""]]);
		($ptrType(SyscallError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		SyscallError.init([["Syscall", "Syscall", "", $String, ""], ["Err", "Err", "", $error, ""]]);
		($ptrType(LinkError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		LinkError.init([["Op", "Op", "", $String, ""], ["Old", "Old", "", $String, ""], ["New", "New", "", $String, ""], ["Err", "Err", "", $error, ""]]);
		File.methods = [["close", "close", "os", [], [$error], false, 0]];
		($ptrType(File)).methods = [["Chdir", "Chdir", "", [], [$error], false, -1], ["Chmod", "Chmod", "", [FileMode], [$error], false, -1], ["Chown", "Chown", "", [$Int, $Int], [$error], false, -1], ["Close", "Close", "", [], [$error], false, -1], ["Fd", "Fd", "", [], [$Uintptr], false, -1], ["Name", "Name", "", [], [$String], false, -1], ["Read", "Read", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["ReadAt", "ReadAt", "", [($sliceType($Uint8)), $Int64], [$Int, $error], false, -1], ["Readdir", "Readdir", "", [$Int], [($sliceType(FileInfo)), $error], false, -1], ["Readdirnames", "Readdirnames", "", [$Int], [($sliceType($String)), $error], false, -1], ["Seek", "Seek", "", [$Int64, $Int], [$Int64, $error], false, -1], ["Stat", "Stat", "", [], [FileInfo, $error], false, -1], ["Sync", "Sync", "", [], [$error], false, -1], ["Truncate", "Truncate", "", [$Int64], [$error], false, -1], ["Write", "Write", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["WriteAt", "WriteAt", "", [($sliceType($Uint8)), $Int64], [$Int, $error], false, -1], ["WriteString", "WriteString", "", [$String], [$Int, $error], false, -1], ["close", "close", "os", [], [$error], false, 0], ["pread", "pread", "os", [($sliceType($Uint8)), $Int64], [$Int, $error], false, -1], ["pwrite", "pwrite", "os", [($sliceType($Uint8)), $Int64], [$Int, $error], false, -1], ["read", "read", "os", [($sliceType($Uint8))], [$Int, $error], false, -1], ["readdir", "readdir", "os", [$Int], [($sliceType(FileInfo)), $error], false, -1], ["readdirnames", "readdirnames", "os", [$Int], [($sliceType($String)), $error], false, -1], ["seek", "seek", "os", [$Int64, $Int], [$Int64, $error], false, -1], ["write", "write", "os", [($sliceType($Uint8))], [$Int, $error], false, -1]];
		File.init([["file", "", "os", ($ptrType(file)), ""]]);
		($ptrType(file)).methods = [["close", "close", "os", [], [$error], false, -1]];
		file.init([["fd", "fd", "os", $Int, ""], ["name", "name", "os", $String, ""], ["dirinfo", "dirinfo", "os", ($ptrType(dirInfo)), ""], ["nepipe", "nepipe", "os", $Int32, ""]]);
		dirInfo.init([["buf", "buf", "os", ($sliceType($Uint8)), ""], ["nbuf", "nbuf", "os", $Int, ""], ["bufp", "bufp", "os", $Int, ""]]);
		FileInfo.init([["IsDir", "IsDir", "", [], [$Bool], false], ["ModTime", "ModTime", "", [], [time.Time], false], ["Mode", "Mode", "", [], [FileMode], false], ["Name", "Name", "", [], [$String], false], ["Size", "Size", "", [], [$Int64], false], ["Sys", "Sys", "", [], [$emptyInterface], false]]);
		FileMode.methods = [["IsDir", "IsDir", "", [], [$Bool], false, -1], ["IsRegular", "IsRegular", "", [], [$Bool], false, -1], ["Perm", "Perm", "", [], [FileMode], false, -1], ["String", "String", "", [], [$String], false, -1]];
		($ptrType(FileMode)).methods = [["IsDir", "IsDir", "", [], [$Bool], false, -1], ["IsRegular", "IsRegular", "", [], [$Bool], false, -1], ["Perm", "Perm", "", [], [FileMode], false, -1], ["String", "String", "", [], [$String], false, -1]];
		($ptrType(fileStat)).methods = [["IsDir", "IsDir", "", [], [$Bool], false, -1], ["ModTime", "ModTime", "", [], [time.Time], false, -1], ["Mode", "Mode", "", [], [FileMode], false, -1], ["Name", "Name", "", [], [$String], false, -1], ["Size", "Size", "", [], [$Int64], false, -1], ["Sys", "Sys", "", [], [$emptyInterface], false, -1]];
		fileStat.init([["name", "name", "os", $String, ""], ["size", "size", "os", $Int64, ""], ["mode", "mode", "os", FileMode, ""], ["modTime", "modTime", "os", time.Time, ""], ["sys", "sys", "os", $emptyInterface, ""]]);
		$pkg.Args = ($sliceType($String)).nil;
		$pkg.ErrInvalid = errors.New("invalid argument");
		$pkg.ErrPermission = errors.New("permission denied");
		$pkg.ErrExist = errors.New("file already exists");
		$pkg.ErrNotExist = errors.New("file does not exist");
		$pkg.Stdin = NewFile((syscall.Stdin >>> 0), "/dev/stdin");
		$pkg.Stdout = NewFile((syscall.Stdout >>> 0), "/dev/stdout");
		$pkg.Stderr = NewFile((syscall.Stderr >>> 0), "/dev/stderr");
		lstat = Lstat;
		init();
	};
	return $pkg;
})();
$packages["strconv"] = (function() {
	var $pkg = {}, math = $packages["math"], errors = $packages["errors"], utf8 = $packages["unicode/utf8"], decimal, leftCheat, extFloat, floatInfo, decimalSlice, optimize, leftcheats, smallPowersOfTen, powersOfTen, uint64pow10, float32info, float64info, isPrint16, isNotPrint16, isPrint32, isNotPrint32, shifts, digitZero, trim, rightShift, prefixIsLessThan, leftShift, shouldRoundUp, frexp10Many, adjustLastDigitFixed, adjustLastDigit, AppendFloat, genericFtoa, bigFtoa, formatDigits, roundShortest, fmtE, fmtF, fmtB, max, FormatInt, Itoa, formatBits, quoteWith, Quote, QuoteToASCII, QuoteRune, AppendQuoteRune, QuoteRuneToASCII, AppendQuoteRuneToASCII, CanBackquote, unhex, UnquoteChar, Unquote, contains, bsearch16, bsearch32, IsPrint;
	decimal = $pkg.decimal = $newType(0, "Struct", "strconv.decimal", "decimal", "strconv", function(d_, nd_, dp_, neg_, trunc_) {
		this.$val = this;
		this.d = d_ !== undefined ? d_ : ($arrayType($Uint8, 800)).zero();
		this.nd = nd_ !== undefined ? nd_ : 0;
		this.dp = dp_ !== undefined ? dp_ : 0;
		this.neg = neg_ !== undefined ? neg_ : false;
		this.trunc = trunc_ !== undefined ? trunc_ : false;
	});
	leftCheat = $pkg.leftCheat = $newType(0, "Struct", "strconv.leftCheat", "leftCheat", "strconv", function(delta_, cutoff_) {
		this.$val = this;
		this.delta = delta_ !== undefined ? delta_ : 0;
		this.cutoff = cutoff_ !== undefined ? cutoff_ : "";
	});
	extFloat = $pkg.extFloat = $newType(0, "Struct", "strconv.extFloat", "extFloat", "strconv", function(mant_, exp_, neg_) {
		this.$val = this;
		this.mant = mant_ !== undefined ? mant_ : new $Uint64(0, 0);
		this.exp = exp_ !== undefined ? exp_ : 0;
		this.neg = neg_ !== undefined ? neg_ : false;
	});
	floatInfo = $pkg.floatInfo = $newType(0, "Struct", "strconv.floatInfo", "floatInfo", "strconv", function(mantbits_, expbits_, bias_) {
		this.$val = this;
		this.mantbits = mantbits_ !== undefined ? mantbits_ : 0;
		this.expbits = expbits_ !== undefined ? expbits_ : 0;
		this.bias = bias_ !== undefined ? bias_ : 0;
	});
	decimalSlice = $pkg.decimalSlice = $newType(0, "Struct", "strconv.decimalSlice", "decimalSlice", "strconv", function(d_, nd_, dp_, neg_) {
		this.$val = this;
		this.d = d_ !== undefined ? d_ : ($sliceType($Uint8)).nil;
		this.nd = nd_ !== undefined ? nd_ : 0;
		this.dp = dp_ !== undefined ? dp_ : 0;
		this.neg = neg_ !== undefined ? neg_ : false;
	});
	decimal.Ptr.prototype.String = function() {
		var a, n, buf, w;
		a = this;
		n = 10 + a.nd >> 0;
		if (a.dp > 0) {
			n = n + (a.dp) >> 0;
		}
		if (a.dp < 0) {
			n = n + (-a.dp) >> 0;
		}
		buf = ($sliceType($Uint8)).make(n);
		w = 0;
		if (a.nd === 0) {
			return "0";
		} else if (a.dp <= 0) {
			(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = 48;
			w = w + (1) >> 0;
			(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = 46;
			w = w + (1) >> 0;
			w = w + (digitZero($subslice(buf, w, (w + -a.dp >> 0)))) >> 0;
			w = w + ($copySlice($subslice(buf, w), $subslice(new ($sliceType($Uint8))(a.d), 0, a.nd))) >> 0;
		} else if (a.dp < a.nd) {
			w = w + ($copySlice($subslice(buf, w), $subslice(new ($sliceType($Uint8))(a.d), 0, a.dp))) >> 0;
			(w < 0 || w >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + w] = 46;
			w = w + (1) >> 0;
			w = w + ($copySlice($subslice(buf, w), $subslice(new ($sliceType($Uint8))(a.d), a.dp, a.nd))) >> 0;
		} else {
			w = w + ($copySlice($subslice(buf, w), $subslice(new ($sliceType($Uint8))(a.d), 0, a.nd))) >> 0;
			w = w + (digitZero($subslice(buf, w, ((w + a.dp >> 0) - a.nd >> 0)))) >> 0;
		}
		return $bytesToString($subslice(buf, 0, w));
	};
	decimal.prototype.String = function() { return this.$val.String(); };
	digitZero = function(dst) {
		var _ref, _i, i;
		_ref = dst;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			(i < 0 || i >= dst.$length) ? $throwRuntimeError("index out of range") : dst.$array[dst.$offset + i] = 48;
			_i++;
		}
		return dst.$length;
	};
	trim = function(a) {
		var x, x$1;
		while (a.nd > 0 && ((x = a.d, x$1 = a.nd - 1 >> 0, ((x$1 < 0 || x$1 >= x.length) ? $throwRuntimeError("index out of range") : x[x$1])) === 48)) {
			a.nd = a.nd - (1) >> 0;
		}
		if (a.nd === 0) {
			a.dp = 0;
		}
	};
	decimal.Ptr.prototype.Assign = function(v) {
		var a, buf, n, v1, x, x$1, x$2;
		a = this;
		buf = ($arrayType($Uint8, 24)).zero(); $copy(buf, ($arrayType($Uint8, 24)).zero(), ($arrayType($Uint8, 24)));
		n = 0;
		while ((v.$high > 0 || (v.$high === 0 && v.$low > 0))) {
			v1 = $div64(v, new $Uint64(0, 10), false);
			v = (x = $mul64(new $Uint64(0, 10), v1), new $Uint64(v.$high - x.$high, v.$low - x.$low));
			(n < 0 || n >= buf.length) ? $throwRuntimeError("index out of range") : buf[n] = (new $Uint64(v.$high + 0, v.$low + 48).$low << 24 >>> 24);
			n = n + (1) >> 0;
			v = v1;
		}
		a.nd = 0;
		n = n - (1) >> 0;
		while (n >= 0) {
			(x$1 = a.d, x$2 = a.nd, (x$2 < 0 || x$2 >= x$1.length) ? $throwRuntimeError("index out of range") : x$1[x$2] = ((n < 0 || n >= buf.length) ? $throwRuntimeError("index out of range") : buf[n]));
			a.nd = a.nd + (1) >> 0;
			n = n - (1) >> 0;
		}
		a.dp = a.nd;
		trim(a);
	};
	decimal.prototype.Assign = function(v) { return this.$val.Assign(v); };
	rightShift = function(a, k) {
		var r, w, n, x, c, x$1, c$1, dig, y, x$2, dig$1, y$1, x$3;
		r = 0;
		w = 0;
		n = 0;
		while (((n >> $min(k, 31)) >> 0) === 0) {
			if (r >= a.nd) {
				if (n === 0) {
					a.nd = 0;
					return;
				}
				while (((n >> $min(k, 31)) >> 0) === 0) {
					n = (((n >>> 16 << 16) * 10 >> 0) + (n << 16 >>> 16) * 10) >> 0;
					r = r + (1) >> 0;
				}
				break;
			}
			c = ((x = a.d, ((r < 0 || r >= x.length) ? $throwRuntimeError("index out of range") : x[r])) >> 0);
			n = (((((n >>> 16 << 16) * 10 >> 0) + (n << 16 >>> 16) * 10) >> 0) + c >> 0) - 48 >> 0;
			r = r + (1) >> 0;
		}
		a.dp = a.dp - ((r - 1 >> 0)) >> 0;
		while (r < a.nd) {
			c$1 = ((x$1 = a.d, ((r < 0 || r >= x$1.length) ? $throwRuntimeError("index out of range") : x$1[r])) >> 0);
			dig = (n >> $min(k, 31)) >> 0;
			n = n - (((y = k, y < 32 ? (dig << y) : 0) >> 0)) >> 0;
			(x$2 = a.d, (w < 0 || w >= x$2.length) ? $throwRuntimeError("index out of range") : x$2[w] = ((dig + 48 >> 0) << 24 >>> 24));
			w = w + (1) >> 0;
			n = (((((n >>> 16 << 16) * 10 >> 0) + (n << 16 >>> 16) * 10) >> 0) + c$1 >> 0) - 48 >> 0;
			r = r + (1) >> 0;
		}
		while (n > 0) {
			dig$1 = (n >> $min(k, 31)) >> 0;
			n = n - (((y$1 = k, y$1 < 32 ? (dig$1 << y$1) : 0) >> 0)) >> 0;
			if (w < 800) {
				(x$3 = a.d, (w < 0 || w >= x$3.length) ? $throwRuntimeError("index out of range") : x$3[w] = ((dig$1 + 48 >> 0) << 24 >>> 24));
				w = w + (1) >> 0;
			} else if (dig$1 > 0) {
				a.trunc = true;
			}
			n = (((n >>> 16 << 16) * 10 >> 0) + (n << 16 >>> 16) * 10) >> 0;
		}
		a.nd = w;
		trim(a);
	};
	prefixIsLessThan = function(b, s) {
		var i;
		i = 0;
		while (i < s.length) {
			if (i >= b.$length) {
				return true;
			}
			if (!((((i < 0 || i >= b.$length) ? $throwRuntimeError("index out of range") : b.$array[b.$offset + i]) === s.charCodeAt(i)))) {
				return ((i < 0 || i >= b.$length) ? $throwRuntimeError("index out of range") : b.$array[b.$offset + i]) < s.charCodeAt(i);
			}
			i = i + (1) >> 0;
		}
		return false;
	};
	leftShift = function(a, k) {
		var delta, r, w, n, y, x, _q, quo, rem, x$1, _q$1, quo$1, rem$1, x$2;
		delta = ((k < 0 || k >= leftcheats.$length) ? $throwRuntimeError("index out of range") : leftcheats.$array[leftcheats.$offset + k]).delta;
		if (prefixIsLessThan($subslice(new ($sliceType($Uint8))(a.d), 0, a.nd), ((k < 0 || k >= leftcheats.$length) ? $throwRuntimeError("index out of range") : leftcheats.$array[leftcheats.$offset + k]).cutoff)) {
			delta = delta - (1) >> 0;
		}
		r = a.nd;
		w = a.nd + delta >> 0;
		n = 0;
		r = r - (1) >> 0;
		while (r >= 0) {
			n = n + (((y = k, y < 32 ? (((((x = a.d, ((r < 0 || r >= x.length) ? $throwRuntimeError("index out of range") : x[r])) >> 0) - 48 >> 0)) << y) : 0) >> 0)) >> 0;
			quo = (_q = n / 10, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
			rem = n - ((((10 >>> 16 << 16) * quo >> 0) + (10 << 16 >>> 16) * quo) >> 0) >> 0;
			w = w - (1) >> 0;
			if (w < 800) {
				(x$1 = a.d, (w < 0 || w >= x$1.length) ? $throwRuntimeError("index out of range") : x$1[w] = ((rem + 48 >> 0) << 24 >>> 24));
			} else if (!((rem === 0))) {
				a.trunc = true;
			}
			n = quo;
			r = r - (1) >> 0;
		}
		while (n > 0) {
			quo$1 = (_q$1 = n / 10, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >> 0 : $throwRuntimeError("integer divide by zero"));
			rem$1 = n - ((((10 >>> 16 << 16) * quo$1 >> 0) + (10 << 16 >>> 16) * quo$1) >> 0) >> 0;
			w = w - (1) >> 0;
			if (w < 800) {
				(x$2 = a.d, (w < 0 || w >= x$2.length) ? $throwRuntimeError("index out of range") : x$2[w] = ((rem$1 + 48 >> 0) << 24 >>> 24));
			} else if (!((rem$1 === 0))) {
				a.trunc = true;
			}
			n = quo$1;
		}
		a.nd = a.nd + (delta) >> 0;
		if (a.nd >= 800) {
			a.nd = 800;
		}
		a.dp = a.dp + (delta) >> 0;
		trim(a);
	};
	decimal.Ptr.prototype.Shift = function(k) {
		var a;
		a = this;
		if (a.nd === 0) {
		} else if (k > 0) {
			while (k > 27) {
				leftShift(a, 27);
				k = k - (27) >> 0;
			}
			leftShift(a, (k >>> 0));
		} else if (k < 0) {
			while (k < -27) {
				rightShift(a, 27);
				k = k + (27) >> 0;
			}
			rightShift(a, (-k >>> 0));
		}
	};
	decimal.prototype.Shift = function(k) { return this.$val.Shift(k); };
	shouldRoundUp = function(a, nd) {
		var x, _r, x$1, x$2, x$3;
		if (nd < 0 || nd >= a.nd) {
			return false;
		}
		if (((x = a.d, ((nd < 0 || nd >= x.length) ? $throwRuntimeError("index out of range") : x[nd])) === 53) && ((nd + 1 >> 0) === a.nd)) {
			if (a.trunc) {
				return true;
			}
			return nd > 0 && !(((_r = (((x$1 = a.d, x$2 = nd - 1 >> 0, ((x$2 < 0 || x$2 >= x$1.length) ? $throwRuntimeError("index out of range") : x$1[x$2])) - 48 << 24 >>> 24)) % 2, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) === 0));
		}
		return (x$3 = a.d, ((nd < 0 || nd >= x$3.length) ? $throwRuntimeError("index out of range") : x$3[nd])) >= 53;
	};
	decimal.Ptr.prototype.Round = function(nd) {
		var a;
		a = this;
		if (nd < 0 || nd >= a.nd) {
			return;
		}
		if (shouldRoundUp(a, nd)) {
			a.RoundUp(nd);
		} else {
			a.RoundDown(nd);
		}
	};
	decimal.prototype.Round = function(nd) { return this.$val.Round(nd); };
	decimal.Ptr.prototype.RoundDown = function(nd) {
		var a;
		a = this;
		if (nd < 0 || nd >= a.nd) {
			return;
		}
		a.nd = nd;
		trim(a);
	};
	decimal.prototype.RoundDown = function(nd) { return this.$val.RoundDown(nd); };
	decimal.Ptr.prototype.RoundUp = function(nd) {
		var a, i, x, c, _lhs, _index;
		a = this;
		if (nd < 0 || nd >= a.nd) {
			return;
		}
		i = nd - 1 >> 0;
		while (i >= 0) {
			c = (x = a.d, ((i < 0 || i >= x.length) ? $throwRuntimeError("index out of range") : x[i]));
			if (c < 57) {
				_lhs = a.d; _index = i; (_index < 0 || _index >= _lhs.length) ? $throwRuntimeError("index out of range") : _lhs[_index] = ((_index < 0 || _index >= _lhs.length) ? $throwRuntimeError("index out of range") : _lhs[_index]) + (1) << 24 >>> 24;
				a.nd = i + 1 >> 0;
				return;
			}
			i = i - (1) >> 0;
		}
		a.d[0] = 49;
		a.nd = 1;
		a.dp = a.dp + (1) >> 0;
	};
	decimal.prototype.RoundUp = function(nd) { return this.$val.RoundUp(nd); };
	decimal.Ptr.prototype.RoundedInteger = function() {
		var a, i, n, x, x$1, x$2, x$3;
		a = this;
		if (a.dp > 20) {
			return new $Uint64(4294967295, 4294967295);
		}
		i = 0;
		n = new $Uint64(0, 0);
		i = 0;
		while (i < a.dp && i < a.nd) {
			n = (x = $mul64(n, new $Uint64(0, 10)), x$1 = new $Uint64(0, ((x$2 = a.d, ((i < 0 || i >= x$2.length) ? $throwRuntimeError("index out of range") : x$2[i])) - 48 << 24 >>> 24)), new $Uint64(x.$high + x$1.$high, x.$low + x$1.$low));
			i = i + (1) >> 0;
		}
		while (i < a.dp) {
			n = $mul64(n, (new $Uint64(0, 10)));
			i = i + (1) >> 0;
		}
		if (shouldRoundUp(a, a.dp)) {
			n = (x$3 = new $Uint64(0, 1), new $Uint64(n.$high + x$3.$high, n.$low + x$3.$low));
		}
		return n;
	};
	decimal.prototype.RoundedInteger = function() { return this.$val.RoundedInteger(); };
	extFloat.Ptr.prototype.AssignComputeBounds = function(mant, exp, neg, flt) {
		var lower = new extFloat.Ptr(), upper = new extFloat.Ptr(), f, x, _tmp, _tmp$1, expBiased, x$1, x$2, x$3, x$4;
		f = this;
		f.mant = mant;
		f.exp = exp - (flt.mantbits >> 0) >> 0;
		f.neg = neg;
		if (f.exp <= 0 && (x = $shiftLeft64(($shiftRightUint64(mant, (-f.exp >>> 0))), (-f.exp >>> 0)), (mant.$high === x.$high && mant.$low === x.$low))) {
			f.mant = $shiftRightUint64(f.mant, ((-f.exp >>> 0)));
			f.exp = 0;
			_tmp = new extFloat.Ptr(); $copy(_tmp, f, extFloat); _tmp$1 = new extFloat.Ptr(); $copy(_tmp$1, f, extFloat); $copy(lower, _tmp, extFloat); $copy(upper, _tmp$1, extFloat);
			return [lower, upper];
		}
		expBiased = exp - flt.bias >> 0;
		$copy(upper, new extFloat.Ptr((x$1 = $mul64(new $Uint64(0, 2), f.mant), new $Uint64(x$1.$high + 0, x$1.$low + 1)), f.exp - 1 >> 0, f.neg), extFloat);
		if (!((x$2 = $shiftLeft64(new $Uint64(0, 1), flt.mantbits), (mant.$high === x$2.$high && mant.$low === x$2.$low))) || (expBiased === 1)) {
			$copy(lower, new extFloat.Ptr((x$3 = $mul64(new $Uint64(0, 2), f.mant), new $Uint64(x$3.$high - 0, x$3.$low - 1)), f.exp - 1 >> 0, f.neg), extFloat);
		} else {
			$copy(lower, new extFloat.Ptr((x$4 = $mul64(new $Uint64(0, 4), f.mant), new $Uint64(x$4.$high - 0, x$4.$low - 1)), f.exp - 2 >> 0, f.neg), extFloat);
		}
		return [lower, upper];
	};
	extFloat.prototype.AssignComputeBounds = function(mant, exp, neg, flt) { return this.$val.AssignComputeBounds(mant, exp, neg, flt); };
	extFloat.Ptr.prototype.Normalize = function() {
		var shift = 0, f, _tmp, _tmp$1, mant, exp, x, x$1, x$2, x$3, x$4, x$5, _tmp$2, _tmp$3;
		f = this;
		_tmp = f.mant; _tmp$1 = f.exp; mant = _tmp; exp = _tmp$1;
		if ((mant.$high === 0 && mant.$low === 0)) {
			shift = 0;
			return shift;
		}
		if ((x = $shiftRightUint64(mant, 32), (x.$high === 0 && x.$low === 0))) {
			mant = $shiftLeft64(mant, (32));
			exp = exp - (32) >> 0;
		}
		if ((x$1 = $shiftRightUint64(mant, 48), (x$1.$high === 0 && x$1.$low === 0))) {
			mant = $shiftLeft64(mant, (16));
			exp = exp - (16) >> 0;
		}
		if ((x$2 = $shiftRightUint64(mant, 56), (x$2.$high === 0 && x$2.$low === 0))) {
			mant = $shiftLeft64(mant, (8));
			exp = exp - (8) >> 0;
		}
		if ((x$3 = $shiftRightUint64(mant, 60), (x$3.$high === 0 && x$3.$low === 0))) {
			mant = $shiftLeft64(mant, (4));
			exp = exp - (4) >> 0;
		}
		if ((x$4 = $shiftRightUint64(mant, 62), (x$4.$high === 0 && x$4.$low === 0))) {
			mant = $shiftLeft64(mant, (2));
			exp = exp - (2) >> 0;
		}
		if ((x$5 = $shiftRightUint64(mant, 63), (x$5.$high === 0 && x$5.$low === 0))) {
			mant = $shiftLeft64(mant, (1));
			exp = exp - (1) >> 0;
		}
		shift = ((f.exp - exp >> 0) >>> 0);
		_tmp$2 = mant; _tmp$3 = exp; f.mant = _tmp$2; f.exp = _tmp$3;
		return shift;
	};
	extFloat.prototype.Normalize = function() { return this.$val.Normalize(); };
	extFloat.Ptr.prototype.Multiply = function(g) {
		var f, _tmp, _tmp$1, fhi, flo, _tmp$2, _tmp$3, ghi, glo, cross1, cross2, x, x$1, x$2, x$3, x$4, x$5, x$6, x$7, rem, x$8, x$9, x$10;
		f = this;
		_tmp = $shiftRightUint64(f.mant, 32); _tmp$1 = new $Uint64(0, (f.mant.$low >>> 0)); fhi = _tmp; flo = _tmp$1;
		_tmp$2 = $shiftRightUint64(g.mant, 32); _tmp$3 = new $Uint64(0, (g.mant.$low >>> 0)); ghi = _tmp$2; glo = _tmp$3;
		cross1 = $mul64(fhi, glo);
		cross2 = $mul64(flo, ghi);
		f.mant = (x = (x$1 = $mul64(fhi, ghi), x$2 = $shiftRightUint64(cross1, 32), new $Uint64(x$1.$high + x$2.$high, x$1.$low + x$2.$low)), x$3 = $shiftRightUint64(cross2, 32), new $Uint64(x.$high + x$3.$high, x.$low + x$3.$low));
		rem = (x$4 = (x$5 = new $Uint64(0, (cross1.$low >>> 0)), x$6 = new $Uint64(0, (cross2.$low >>> 0)), new $Uint64(x$5.$high + x$6.$high, x$5.$low + x$6.$low)), x$7 = $shiftRightUint64(($mul64(flo, glo)), 32), new $Uint64(x$4.$high + x$7.$high, x$4.$low + x$7.$low));
		rem = (x$8 = new $Uint64(0, 2147483648), new $Uint64(rem.$high + x$8.$high, rem.$low + x$8.$low));
		f.mant = (x$9 = f.mant, x$10 = ($shiftRightUint64(rem, 32)), new $Uint64(x$9.$high + x$10.$high, x$9.$low + x$10.$low));
		f.exp = (f.exp + g.exp >> 0) + 64 >> 0;
	};
	extFloat.prototype.Multiply = function(g) { return this.$val.Multiply(g); };
	extFloat.Ptr.prototype.AssignDecimal = function(mantissa, exp10, neg, trunc, flt) {
		var ok = false, f, errors$1, _q, i, _r, adjExp, x, x$1, shift, y, denormalExp, extrabits, halfway, x$2, x$3, x$4, mant_extra, x$5, x$6, x$7, x$8, x$9, x$10, x$11, x$12;
		f = this;
		errors$1 = 0;
		if (trunc) {
			errors$1 = errors$1 + (4) >> 0;
		}
		f.mant = mantissa;
		f.exp = 0;
		f.neg = neg;
		i = (_q = ((exp10 - -348 >> 0)) / 8, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		if (exp10 < -348 || i >= 87) {
			ok = false;
			return ok;
		}
		adjExp = (_r = ((exp10 - -348 >> 0)) % 8, _r === _r ? _r : $throwRuntimeError("integer divide by zero"));
		if (adjExp < 19 && (x = (x$1 = 19 - adjExp >> 0, ((x$1 < 0 || x$1 >= uint64pow10.length) ? $throwRuntimeError("index out of range") : uint64pow10[x$1])), (mantissa.$high < x.$high || (mantissa.$high === x.$high && mantissa.$low < x.$low)))) {
			f.mant = $mul64(f.mant, (((adjExp < 0 || adjExp >= uint64pow10.length) ? $throwRuntimeError("index out of range") : uint64pow10[adjExp])));
			f.Normalize();
		} else {
			f.Normalize();
			f.Multiply($clone(((adjExp < 0 || adjExp >= smallPowersOfTen.length) ? $throwRuntimeError("index out of range") : smallPowersOfTen[adjExp]), extFloat));
			errors$1 = errors$1 + (4) >> 0;
		}
		f.Multiply($clone(((i < 0 || i >= powersOfTen.length) ? $throwRuntimeError("index out of range") : powersOfTen[i]), extFloat));
		if (errors$1 > 0) {
			errors$1 = errors$1 + (1) >> 0;
		}
		errors$1 = errors$1 + (4) >> 0;
		shift = f.Normalize();
		errors$1 = (y = (shift), y < 32 ? (errors$1 << y) : 0) >> 0;
		denormalExp = flt.bias - 63 >> 0;
		extrabits = 0;
		if (f.exp <= denormalExp) {
			extrabits = (((63 - flt.mantbits >>> 0) + 1 >>> 0) + ((denormalExp - f.exp >> 0) >>> 0) >>> 0);
		} else {
			extrabits = (63 - flt.mantbits >>> 0);
		}
		halfway = $shiftLeft64(new $Uint64(0, 1), ((extrabits - 1 >>> 0)));
		mant_extra = (x$2 = f.mant, x$3 = (x$4 = $shiftLeft64(new $Uint64(0, 1), extrabits), new $Uint64(x$4.$high - 0, x$4.$low - 1)), new $Uint64(x$2.$high & x$3.$high, (x$2.$low & x$3.$low) >>> 0));
		if ((x$5 = (x$6 = new $Int64(halfway.$high, halfway.$low), x$7 = new $Int64(0, errors$1), new $Int64(x$6.$high - x$7.$high, x$6.$low - x$7.$low)), x$8 = new $Int64(mant_extra.$high, mant_extra.$low), (x$5.$high < x$8.$high || (x$5.$high === x$8.$high && x$5.$low < x$8.$low))) && (x$9 = new $Int64(mant_extra.$high, mant_extra.$low), x$10 = (x$11 = new $Int64(halfway.$high, halfway.$low), x$12 = new $Int64(0, errors$1), new $Int64(x$11.$high + x$12.$high, x$11.$low + x$12.$low)), (x$9.$high < x$10.$high || (x$9.$high === x$10.$high && x$9.$low < x$10.$low)))) {
			ok = false;
			return ok;
		}
		ok = true;
		return ok;
	};
	extFloat.prototype.AssignDecimal = function(mantissa, exp10, neg, trunc, flt) { return this.$val.AssignDecimal(mantissa, exp10, neg, trunc, flt); };
	extFloat.Ptr.prototype.frexp10 = function() {
		var exp10 = 0, index = 0, f, _q, x, approxExp10, _q$1, i, exp, _tmp, _tmp$1;
		f = this;
		approxExp10 = (_q = (x = (-46 - f.exp >> 0), (((x >>> 16 << 16) * 28 >> 0) + (x << 16 >>> 16) * 28) >> 0) / 93, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		i = (_q$1 = ((approxExp10 - -348 >> 0)) / 8, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >> 0 : $throwRuntimeError("integer divide by zero"));
		Loop:
		while (true) {
			exp = (f.exp + ((i < 0 || i >= powersOfTen.length) ? $throwRuntimeError("index out of range") : powersOfTen[i]).exp >> 0) + 64 >> 0;
			if (exp < -60) {
				i = i + (1) >> 0;
			} else if (exp > -32) {
				i = i - (1) >> 0;
			} else {
				break Loop;
			}
		}
		f.Multiply($clone(((i < 0 || i >= powersOfTen.length) ? $throwRuntimeError("index out of range") : powersOfTen[i]), extFloat));
		_tmp = -((-348 + ((((i >>> 16 << 16) * 8 >> 0) + (i << 16 >>> 16) * 8) >> 0) >> 0)); _tmp$1 = i; exp10 = _tmp; index = _tmp$1;
		return [exp10, index];
	};
	extFloat.prototype.frexp10 = function() { return this.$val.frexp10(); };
	frexp10Many = function(a, b, c) {
		var exp10 = 0, _tuple, i;
		_tuple = c.frexp10(); exp10 = _tuple[0]; i = _tuple[1];
		a.Multiply($clone(((i < 0 || i >= powersOfTen.length) ? $throwRuntimeError("index out of range") : powersOfTen[i]), extFloat));
		b.Multiply($clone(((i < 0 || i >= powersOfTen.length) ? $throwRuntimeError("index out of range") : powersOfTen[i]), extFloat));
		return exp10;
	};
	extFloat.Ptr.prototype.FixedDecimal = function(d, n) {
		var f, x, _tuple, exp10, shift, integer, x$1, x$2, fraction, nonAsciiName, needed, integerDigits, pow10, _tmp, _tmp$1, i, pow, x$3, rest, x$4, _q, x$5, buf, pos, v, _q$1, v1, i$1, x$6, x$7, nd, x$8, x$9, digit, x$10, x$11, x$12, ok, i$2, x$13;
		f = this;
		if ((x = f.mant, (x.$high === 0 && x.$low === 0))) {
			d.nd = 0;
			d.dp = 0;
			d.neg = f.neg;
			return true;
		}
		if (n === 0) {
			$panic(new $String("strconv: internal error: extFloat.FixedDecimal called with n == 0"));
		}
		f.Normalize();
		_tuple = f.frexp10(); exp10 = _tuple[0];
		shift = (-f.exp >>> 0);
		integer = ($shiftRightUint64(f.mant, shift).$low >>> 0);
		fraction = (x$1 = f.mant, x$2 = $shiftLeft64(new $Uint64(0, integer), shift), new $Uint64(x$1.$high - x$2.$high, x$1.$low - x$2.$low));
		nonAsciiName = new $Uint64(0, 1);
		needed = n;
		integerDigits = 0;
		pow10 = new $Uint64(0, 1);
		_tmp = 0; _tmp$1 = new $Uint64(0, 1); i = _tmp; pow = _tmp$1;
		while (i < 20) {
			if ((x$3 = new $Uint64(0, integer), (pow.$high > x$3.$high || (pow.$high === x$3.$high && pow.$low > x$3.$low)))) {
				integerDigits = i;
				break;
			}
			pow = $mul64(pow, (new $Uint64(0, 10)));
			i = i + (1) >> 0;
		}
		rest = integer;
		if (integerDigits > needed) {
			pow10 = (x$4 = integerDigits - needed >> 0, ((x$4 < 0 || x$4 >= uint64pow10.length) ? $throwRuntimeError("index out of range") : uint64pow10[x$4]));
			integer = (_q = integer / ((pow10.$low >>> 0)), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >>> 0 : $throwRuntimeError("integer divide by zero"));
			rest = rest - ((x$5 = (pow10.$low >>> 0), (((integer >>> 16 << 16) * x$5 >>> 0) + (integer << 16 >>> 16) * x$5) >>> 0)) >>> 0;
		} else {
			rest = 0;
		}
		buf = ($arrayType($Uint8, 32)).zero(); $copy(buf, ($arrayType($Uint8, 32)).zero(), ($arrayType($Uint8, 32)));
		pos = 32;
		v = integer;
		while (v > 0) {
			v1 = (_q$1 = v / 10, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >>> 0 : $throwRuntimeError("integer divide by zero"));
			v = v - (((((10 >>> 16 << 16) * v1 >>> 0) + (10 << 16 >>> 16) * v1) >>> 0)) >>> 0;
			pos = pos - (1) >> 0;
			(pos < 0 || pos >= buf.length) ? $throwRuntimeError("index out of range") : buf[pos] = ((v + 48 >>> 0) << 24 >>> 24);
			v = v1;
		}
		i$1 = pos;
		while (i$1 < 32) {
			(x$6 = d.d, x$7 = i$1 - pos >> 0, (x$7 < 0 || x$7 >= x$6.$length) ? $throwRuntimeError("index out of range") : x$6.$array[x$6.$offset + x$7] = ((i$1 < 0 || i$1 >= buf.length) ? $throwRuntimeError("index out of range") : buf[i$1]));
			i$1 = i$1 + (1) >> 0;
		}
		nd = 32 - pos >> 0;
		d.nd = nd;
		d.dp = integerDigits + exp10 >> 0;
		needed = needed - (nd) >> 0;
		if (needed > 0) {
			if (!((rest === 0)) || !((pow10.$high === 0 && pow10.$low === 1))) {
				$panic(new $String("strconv: internal error, rest != 0 but needed > 0"));
			}
			while (needed > 0) {
				fraction = $mul64(fraction, (new $Uint64(0, 10)));
				nonAsciiName = $mul64(nonAsciiName, (new $Uint64(0, 10)));
				if ((x$8 = $mul64(new $Uint64(0, 2), nonAsciiName), x$9 = $shiftLeft64(new $Uint64(0, 1), shift), (x$8.$high > x$9.$high || (x$8.$high === x$9.$high && x$8.$low > x$9.$low)))) {
					return false;
				}
				digit = $shiftRightUint64(fraction, shift);
				(x$10 = d.d, (nd < 0 || nd >= x$10.$length) ? $throwRuntimeError("index out of range") : x$10.$array[x$10.$offset + nd] = (new $Uint64(digit.$high + 0, digit.$low + 48).$low << 24 >>> 24));
				fraction = (x$11 = $shiftLeft64(digit, shift), new $Uint64(fraction.$high - x$11.$high, fraction.$low - x$11.$low));
				nd = nd + (1) >> 0;
				needed = needed - (1) >> 0;
			}
			d.nd = nd;
		}
		ok = adjustLastDigitFixed(d, (x$12 = $shiftLeft64(new $Uint64(0, rest), shift), new $Uint64(x$12.$high | fraction.$high, (x$12.$low | fraction.$low) >>> 0)), pow10, shift, nonAsciiName);
		if (!ok) {
			return false;
		}
		i$2 = d.nd - 1 >> 0;
		while (i$2 >= 0) {
			if (!(((x$13 = d.d, ((i$2 < 0 || i$2 >= x$13.$length) ? $throwRuntimeError("index out of range") : x$13.$array[x$13.$offset + i$2])) === 48))) {
				d.nd = i$2 + 1 >> 0;
				break;
			}
			i$2 = i$2 - (1) >> 0;
		}
		return true;
	};
	extFloat.prototype.FixedDecimal = function(d, n) { return this.$val.FixedDecimal(d, n); };
	adjustLastDigitFixed = function(d, num, den, shift, nonAsciiName) {
		var x, x$1, x$2, x$3, x$4, x$5, x$6, i, x$7, x$8, _lhs, _index;
		if ((x = $shiftLeft64(den, shift), (num.$high > x.$high || (num.$high === x.$high && num.$low > x.$low)))) {
			$panic(new $String("strconv: num > den<<shift in adjustLastDigitFixed"));
		}
		if ((x$1 = $mul64(new $Uint64(0, 2), nonAsciiName), x$2 = $shiftLeft64(den, shift), (x$1.$high > x$2.$high || (x$1.$high === x$2.$high && x$1.$low > x$2.$low)))) {
			$panic(new $String("strconv: \xCE\xB5 > (den<<shift)/2"));
		}
		if ((x$3 = $mul64(new $Uint64(0, 2), (new $Uint64(num.$high + nonAsciiName.$high, num.$low + nonAsciiName.$low))), x$4 = $shiftLeft64(den, shift), (x$3.$high < x$4.$high || (x$3.$high === x$4.$high && x$3.$low < x$4.$low)))) {
			return true;
		}
		if ((x$5 = $mul64(new $Uint64(0, 2), (new $Uint64(num.$high - nonAsciiName.$high, num.$low - nonAsciiName.$low))), x$6 = $shiftLeft64(den, shift), (x$5.$high > x$6.$high || (x$5.$high === x$6.$high && x$5.$low > x$6.$low)))) {
			i = d.nd - 1 >> 0;
			while (i >= 0) {
				if ((x$7 = d.d, ((i < 0 || i >= x$7.$length) ? $throwRuntimeError("index out of range") : x$7.$array[x$7.$offset + i])) === 57) {
					d.nd = d.nd - (1) >> 0;
				} else {
					break;
				}
				i = i - (1) >> 0;
			}
			if (i < 0) {
				(x$8 = d.d, (0 < 0 || 0 >= x$8.$length) ? $throwRuntimeError("index out of range") : x$8.$array[x$8.$offset + 0] = 49);
				d.nd = 1;
				d.dp = d.dp + (1) >> 0;
			} else {
				_lhs = d.d; _index = i; (_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index] = ((_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index]) + (1) << 24 >>> 24;
			}
			return true;
		}
		return false;
	};
	extFloat.Ptr.prototype.ShortestDecimal = function(d, lower, upper) {
		var f, x, buf, n, v, v1, x$1, nd, i, x$2, x$3, _tmp, _tmp$1, x$4, x$5, exp10, x$6, x$7, x$8, x$9, shift, integer, x$10, x$11, fraction, x$12, x$13, allowance, x$14, x$15, targetDiff, integerDigits, _tmp$2, _tmp$3, i$1, pow, x$16, i$2, x$17, pow$1, _q, digit, x$18, x$19, x$20, currentDiff, digit$1, multiplier, x$21, x$22, x$23, x$24;
		f = this;
		if ((x = f.mant, (x.$high === 0 && x.$low === 0))) {
			d.nd = 0;
			d.dp = 0;
			d.neg = f.neg;
			return true;
		}
		if ((f.exp === 0) && $equal(lower, f, extFloat) && $equal(lower, upper, extFloat)) {
			buf = ($arrayType($Uint8, 24)).zero(); $copy(buf, ($arrayType($Uint8, 24)).zero(), ($arrayType($Uint8, 24)));
			n = 23;
			v = f.mant;
			while ((v.$high > 0 || (v.$high === 0 && v.$low > 0))) {
				v1 = $div64(v, new $Uint64(0, 10), false);
				v = (x$1 = $mul64(new $Uint64(0, 10), v1), new $Uint64(v.$high - x$1.$high, v.$low - x$1.$low));
				(n < 0 || n >= buf.length) ? $throwRuntimeError("index out of range") : buf[n] = (new $Uint64(v.$high + 0, v.$low + 48).$low << 24 >>> 24);
				n = n - (1) >> 0;
				v = v1;
			}
			nd = (24 - n >> 0) - 1 >> 0;
			i = 0;
			while (i < nd) {
				(x$3 = d.d, (i < 0 || i >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + i] = (x$2 = (n + 1 >> 0) + i >> 0, ((x$2 < 0 || x$2 >= buf.length) ? $throwRuntimeError("index out of range") : buf[x$2])));
				i = i + (1) >> 0;
			}
			_tmp = nd; _tmp$1 = nd; d.nd = _tmp; d.dp = _tmp$1;
			while (d.nd > 0 && ((x$4 = d.d, x$5 = d.nd - 1 >> 0, ((x$5 < 0 || x$5 >= x$4.$length) ? $throwRuntimeError("index out of range") : x$4.$array[x$4.$offset + x$5])) === 48)) {
				d.nd = d.nd - (1) >> 0;
			}
			if (d.nd === 0) {
				d.dp = 0;
			}
			d.neg = f.neg;
			return true;
		}
		upper.Normalize();
		if (f.exp > upper.exp) {
			f.mant = $shiftLeft64(f.mant, (((f.exp - upper.exp >> 0) >>> 0)));
			f.exp = upper.exp;
		}
		if (lower.exp > upper.exp) {
			lower.mant = $shiftLeft64(lower.mant, (((lower.exp - upper.exp >> 0) >>> 0)));
			lower.exp = upper.exp;
		}
		exp10 = frexp10Many(lower, f, upper);
		upper.mant = (x$6 = upper.mant, x$7 = new $Uint64(0, 1), new $Uint64(x$6.$high + x$7.$high, x$6.$low + x$7.$low));
		lower.mant = (x$8 = lower.mant, x$9 = new $Uint64(0, 1), new $Uint64(x$8.$high - x$9.$high, x$8.$low - x$9.$low));
		shift = (-upper.exp >>> 0);
		integer = ($shiftRightUint64(upper.mant, shift).$low >>> 0);
		fraction = (x$10 = upper.mant, x$11 = $shiftLeft64(new $Uint64(0, integer), shift), new $Uint64(x$10.$high - x$11.$high, x$10.$low - x$11.$low));
		allowance = (x$12 = upper.mant, x$13 = lower.mant, new $Uint64(x$12.$high - x$13.$high, x$12.$low - x$13.$low));
		targetDiff = (x$14 = upper.mant, x$15 = f.mant, new $Uint64(x$14.$high - x$15.$high, x$14.$low - x$15.$low));
		integerDigits = 0;
		_tmp$2 = 0; _tmp$3 = new $Uint64(0, 1); i$1 = _tmp$2; pow = _tmp$3;
		while (i$1 < 20) {
			if ((x$16 = new $Uint64(0, integer), (pow.$high > x$16.$high || (pow.$high === x$16.$high && pow.$low > x$16.$low)))) {
				integerDigits = i$1;
				break;
			}
			pow = $mul64(pow, (new $Uint64(0, 10)));
			i$1 = i$1 + (1) >> 0;
		}
		i$2 = 0;
		while (i$2 < integerDigits) {
			pow$1 = (x$17 = (integerDigits - i$2 >> 0) - 1 >> 0, ((x$17 < 0 || x$17 >= uint64pow10.length) ? $throwRuntimeError("index out of range") : uint64pow10[x$17]));
			digit = (_q = integer / (pow$1.$low >>> 0), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >>> 0 : $throwRuntimeError("integer divide by zero"));
			(x$18 = d.d, (i$2 < 0 || i$2 >= x$18.$length) ? $throwRuntimeError("index out of range") : x$18.$array[x$18.$offset + i$2] = ((digit + 48 >>> 0) << 24 >>> 24));
			integer = integer - ((x$19 = (pow$1.$low >>> 0), (((digit >>> 16 << 16) * x$19 >>> 0) + (digit << 16 >>> 16) * x$19) >>> 0)) >>> 0;
			currentDiff = (x$20 = $shiftLeft64(new $Uint64(0, integer), shift), new $Uint64(x$20.$high + fraction.$high, x$20.$low + fraction.$low));
			if ((currentDiff.$high < allowance.$high || (currentDiff.$high === allowance.$high && currentDiff.$low < allowance.$low))) {
				d.nd = i$2 + 1 >> 0;
				d.dp = integerDigits + exp10 >> 0;
				d.neg = f.neg;
				return adjustLastDigit(d, currentDiff, targetDiff, allowance, $shiftLeft64(pow$1, shift), new $Uint64(0, 2));
			}
			i$2 = i$2 + (1) >> 0;
		}
		d.nd = integerDigits;
		d.dp = d.nd + exp10 >> 0;
		d.neg = f.neg;
		digit$1 = 0;
		multiplier = new $Uint64(0, 1);
		while (true) {
			fraction = $mul64(fraction, (new $Uint64(0, 10)));
			multiplier = $mul64(multiplier, (new $Uint64(0, 10)));
			digit$1 = ($shiftRightUint64(fraction, shift).$low >> 0);
			(x$21 = d.d, x$22 = d.nd, (x$22 < 0 || x$22 >= x$21.$length) ? $throwRuntimeError("index out of range") : x$21.$array[x$21.$offset + x$22] = ((digit$1 + 48 >> 0) << 24 >>> 24));
			d.nd = d.nd + (1) >> 0;
			fraction = (x$23 = $shiftLeft64(new $Uint64(0, digit$1), shift), new $Uint64(fraction.$high - x$23.$high, fraction.$low - x$23.$low));
			if ((x$24 = $mul64(allowance, multiplier), (fraction.$high < x$24.$high || (fraction.$high === x$24.$high && fraction.$low < x$24.$low)))) {
				return adjustLastDigit(d, fraction, $mul64(targetDiff, multiplier), $mul64(allowance, multiplier), $shiftLeft64(new $Uint64(0, 1), shift), $mul64(multiplier, new $Uint64(0, 2)));
			}
		}
	};
	extFloat.prototype.ShortestDecimal = function(d, lower, upper) { return this.$val.ShortestDecimal(d, lower, upper); };
	adjustLastDigit = function(d, currentDiff, targetDiff, maxDiff, ulpDecimal, ulpBinary) {
		var x, x$1, x$2, x$3, _lhs, _index, x$4, x$5, x$6, x$7, x$8, x$9, x$10;
		if ((x = $mul64(new $Uint64(0, 2), ulpBinary), (ulpDecimal.$high < x.$high || (ulpDecimal.$high === x.$high && ulpDecimal.$low < x.$low)))) {
			return false;
		}
		while ((x$1 = (x$2 = (x$3 = $div64(ulpDecimal, new $Uint64(0, 2), false), new $Uint64(currentDiff.$high + x$3.$high, currentDiff.$low + x$3.$low)), new $Uint64(x$2.$high + ulpBinary.$high, x$2.$low + ulpBinary.$low)), (x$1.$high < targetDiff.$high || (x$1.$high === targetDiff.$high && x$1.$low < targetDiff.$low)))) {
			_lhs = d.d; _index = d.nd - 1 >> 0; (_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index] = ((_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index]) - (1) << 24 >>> 24;
			currentDiff = (x$4 = ulpDecimal, new $Uint64(currentDiff.$high + x$4.$high, currentDiff.$low + x$4.$low));
		}
		if ((x$5 = new $Uint64(currentDiff.$high + ulpDecimal.$high, currentDiff.$low + ulpDecimal.$low), x$6 = (x$7 = (x$8 = $div64(ulpDecimal, new $Uint64(0, 2), false), new $Uint64(targetDiff.$high + x$8.$high, targetDiff.$low + x$8.$low)), new $Uint64(x$7.$high + ulpBinary.$high, x$7.$low + ulpBinary.$low)), (x$5.$high < x$6.$high || (x$5.$high === x$6.$high && x$5.$low <= x$6.$low)))) {
			return false;
		}
		if ((currentDiff.$high < ulpBinary.$high || (currentDiff.$high === ulpBinary.$high && currentDiff.$low < ulpBinary.$low)) || (x$9 = new $Uint64(maxDiff.$high - ulpBinary.$high, maxDiff.$low - ulpBinary.$low), (currentDiff.$high > x$9.$high || (currentDiff.$high === x$9.$high && currentDiff.$low > x$9.$low)))) {
			return false;
		}
		if ((d.nd === 1) && ((x$10 = d.d, ((0 < 0 || 0 >= x$10.$length) ? $throwRuntimeError("index out of range") : x$10.$array[x$10.$offset + 0])) === 48)) {
			d.nd = 0;
			d.dp = 0;
		}
		return true;
	};
	AppendFloat = $pkg.AppendFloat = function(dst, f, fmt, prec, bitSize) {
		return genericFtoa(dst, f, fmt, prec, bitSize);
	};
	genericFtoa = function(dst, val, fmt, prec, bitSize) {
		var bits, flt, _ref, x, neg, y, exp, x$1, x$2, mant, _ref$1, y$1, s, x$3, digs, ok, shortest, f, _tuple, lower, upper, buf, _ref$2, digits, _ref$3, buf$1, f$1;
		bits = new $Uint64(0, 0);
		flt = ($ptrType(floatInfo)).nil;
		_ref = bitSize;
		if (_ref === 32) {
			bits = new $Uint64(0, math.Float32bits(val));
			flt = float32info;
		} else if (_ref === 64) {
			bits = math.Float64bits(val);
			flt = float64info;
		} else {
			$panic(new $String("strconv: illegal AppendFloat/FormatFloat bitSize"));
		}
		neg = !((x = $shiftRightUint64(bits, ((flt.expbits + flt.mantbits >>> 0))), (x.$high === 0 && x.$low === 0)));
		exp = ($shiftRightUint64(bits, flt.mantbits).$low >> 0) & ((((y = flt.expbits, y < 32 ? (1 << y) : 0) >> 0) - 1 >> 0));
		mant = (x$1 = (x$2 = $shiftLeft64(new $Uint64(0, 1), flt.mantbits), new $Uint64(x$2.$high - 0, x$2.$low - 1)), new $Uint64(bits.$high & x$1.$high, (bits.$low & x$1.$low) >>> 0));
		_ref$1 = exp;
		if (_ref$1 === (((y$1 = flt.expbits, y$1 < 32 ? (1 << y$1) : 0) >> 0) - 1 >> 0)) {
			s = "";
			if (!((mant.$high === 0 && mant.$low === 0))) {
				s = "NaN";
			} else if (neg) {
				s = "-Inf";
			} else {
				s = "+Inf";
			}
			return $appendSlice(dst, new ($sliceType($Uint8))($stringToBytes(s)));
		} else if (_ref$1 === 0) {
			exp = exp + (1) >> 0;
		} else {
			mant = (x$3 = $shiftLeft64(new $Uint64(0, 1), flt.mantbits), new $Uint64(mant.$high | x$3.$high, (mant.$low | x$3.$low) >>> 0));
		}
		exp = exp + (flt.bias) >> 0;
		if (fmt === 98) {
			return fmtB(dst, neg, mant, exp, flt);
		}
		if (!optimize) {
			return bigFtoa(dst, prec, fmt, neg, mant, exp, flt);
		}
		digs = new decimalSlice.Ptr(); $copy(digs, new decimalSlice.Ptr(), decimalSlice);
		ok = false;
		shortest = prec < 0;
		if (shortest) {
			f = new extFloat.Ptr();
			_tuple = f.AssignComputeBounds(mant, exp, neg, flt); lower = new extFloat.Ptr(); $copy(lower, _tuple[0], extFloat); upper = new extFloat.Ptr(); $copy(upper, _tuple[1], extFloat);
			buf = ($arrayType($Uint8, 32)).zero(); $copy(buf, ($arrayType($Uint8, 32)).zero(), ($arrayType($Uint8, 32)));
			digs.d = new ($sliceType($Uint8))(buf);
			ok = f.ShortestDecimal(digs, lower, upper);
			if (!ok) {
				return bigFtoa(dst, prec, fmt, neg, mant, exp, flt);
			}
			_ref$2 = fmt;
			if (_ref$2 === 101 || _ref$2 === 69) {
				prec = digs.nd - 1 >> 0;
			} else if (_ref$2 === 102) {
				prec = max(digs.nd - digs.dp >> 0, 0);
			} else if (_ref$2 === 103 || _ref$2 === 71) {
				prec = digs.nd;
			}
		} else if (!((fmt === 102))) {
			digits = prec;
			_ref$3 = fmt;
			if (_ref$3 === 101 || _ref$3 === 69) {
				digits = digits + (1) >> 0;
			} else if (_ref$3 === 103 || _ref$3 === 71) {
				if (prec === 0) {
					prec = 1;
				}
				digits = prec;
			}
			if (digits <= 15) {
				buf$1 = ($arrayType($Uint8, 24)).zero(); $copy(buf$1, ($arrayType($Uint8, 24)).zero(), ($arrayType($Uint8, 24)));
				digs.d = new ($sliceType($Uint8))(buf$1);
				f$1 = new extFloat.Ptr(mant, exp - (flt.mantbits >> 0) >> 0, neg);
				ok = f$1.FixedDecimal(digs, digits);
			}
		}
		if (!ok) {
			return bigFtoa(dst, prec, fmt, neg, mant, exp, flt);
		}
		return formatDigits(dst, shortest, neg, $clone(digs, decimalSlice), prec, fmt);
	};
	bigFtoa = function(dst, prec, fmt, neg, mant, exp, flt) {
		var d, digs, shortest, _ref, _ref$1;
		d = new decimal.Ptr();
		d.Assign(mant);
		d.Shift(exp - (flt.mantbits >> 0) >> 0);
		digs = new decimalSlice.Ptr(); $copy(digs, new decimalSlice.Ptr(), decimalSlice);
		shortest = prec < 0;
		if (shortest) {
			roundShortest(d, mant, exp, flt);
			$copy(digs, new decimalSlice.Ptr(new ($sliceType($Uint8))(d.d), d.nd, d.dp, false), decimalSlice);
			_ref = fmt;
			if (_ref === 101 || _ref === 69) {
				prec = digs.nd - 1 >> 0;
			} else if (_ref === 102) {
				prec = max(digs.nd - digs.dp >> 0, 0);
			} else if (_ref === 103 || _ref === 71) {
				prec = digs.nd;
			}
		} else {
			_ref$1 = fmt;
			if (_ref$1 === 101 || _ref$1 === 69) {
				d.Round(prec + 1 >> 0);
			} else if (_ref$1 === 102) {
				d.Round(d.dp + prec >> 0);
			} else if (_ref$1 === 103 || _ref$1 === 71) {
				if (prec === 0) {
					prec = 1;
				}
				d.Round(prec);
			}
			$copy(digs, new decimalSlice.Ptr(new ($sliceType($Uint8))(d.d), d.nd, d.dp, false), decimalSlice);
		}
		return formatDigits(dst, shortest, neg, $clone(digs, decimalSlice), prec, fmt);
	};
	formatDigits = function(dst, shortest, neg, digs, prec, fmt) {
		var _ref, eprec, exp;
		_ref = fmt;
		if (_ref === 101 || _ref === 69) {
			return fmtE(dst, neg, $clone(digs, decimalSlice), prec, fmt);
		} else if (_ref === 102) {
			return fmtF(dst, neg, $clone(digs, decimalSlice), prec);
		} else if (_ref === 103 || _ref === 71) {
			eprec = prec;
			if (eprec > digs.nd && digs.nd >= digs.dp) {
				eprec = digs.nd;
			}
			if (shortest) {
				eprec = 6;
			}
			exp = digs.dp - 1 >> 0;
			if (exp < -4 || exp >= eprec) {
				if (prec > digs.nd) {
					prec = digs.nd;
				}
				return fmtE(dst, neg, $clone(digs, decimalSlice), prec - 1 >> 0, (fmt + 101 << 24 >>> 24) - 103 << 24 >>> 24);
			}
			if (prec > digs.dp) {
				prec = digs.nd;
			}
			return fmtF(dst, neg, $clone(digs, decimalSlice), max(prec - digs.dp >> 0, 0));
		}
		return $append(dst, 37, fmt);
	};
	roundShortest = function(d, mant, exp, flt) {
		var minexp, x, x$1, upper, x$2, mantlo, explo, x$3, x$4, lower, x$5, x$6, inclusive, i, _tmp, _tmp$1, _tmp$2, l, m, u, x$7, x$8, x$9, okdown, okup;
		if ((mant.$high === 0 && mant.$low === 0)) {
			d.nd = 0;
			return;
		}
		minexp = flt.bias + 1 >> 0;
		if (exp > minexp && (x = (d.dp - d.nd >> 0), (((332 >>> 16 << 16) * x >> 0) + (332 << 16 >>> 16) * x) >> 0) >= (x$1 = (exp - (flt.mantbits >> 0) >> 0), (((100 >>> 16 << 16) * x$1 >> 0) + (100 << 16 >>> 16) * x$1) >> 0)) {
			return;
		}
		upper = new decimal.Ptr();
		upper.Assign((x$2 = $mul64(mant, new $Uint64(0, 2)), new $Uint64(x$2.$high + 0, x$2.$low + 1)));
		upper.Shift((exp - (flt.mantbits >> 0) >> 0) - 1 >> 0);
		mantlo = new $Uint64(0, 0);
		explo = 0;
		if ((x$3 = $shiftLeft64(new $Uint64(0, 1), flt.mantbits), (mant.$high > x$3.$high || (mant.$high === x$3.$high && mant.$low > x$3.$low))) || (exp === minexp)) {
			mantlo = new $Uint64(mant.$high - 0, mant.$low - 1);
			explo = exp;
		} else {
			mantlo = (x$4 = $mul64(mant, new $Uint64(0, 2)), new $Uint64(x$4.$high - 0, x$4.$low - 1));
			explo = exp - 1 >> 0;
		}
		lower = new decimal.Ptr();
		lower.Assign((x$5 = $mul64(mantlo, new $Uint64(0, 2)), new $Uint64(x$5.$high + 0, x$5.$low + 1)));
		lower.Shift((explo - (flt.mantbits >> 0) >> 0) - 1 >> 0);
		inclusive = (x$6 = $div64(mant, new $Uint64(0, 2), true), (x$6.$high === 0 && x$6.$low === 0));
		i = 0;
		while (i < d.nd) {
			_tmp = 0; _tmp$1 = 0; _tmp$2 = 0; l = _tmp; m = _tmp$1; u = _tmp$2;
			if (i < lower.nd) {
				l = (x$7 = lower.d, ((i < 0 || i >= x$7.length) ? $throwRuntimeError("index out of range") : x$7[i]));
			} else {
				l = 48;
			}
			m = (x$8 = d.d, ((i < 0 || i >= x$8.length) ? $throwRuntimeError("index out of range") : x$8[i]));
			if (i < upper.nd) {
				u = (x$9 = upper.d, ((i < 0 || i >= x$9.length) ? $throwRuntimeError("index out of range") : x$9[i]));
			} else {
				u = 48;
			}
			okdown = !((l === m)) || (inclusive && (l === m) && ((i + 1 >> 0) === lower.nd));
			okup = !((m === u)) && (inclusive || (m + 1 << 24 >>> 24) < u || (i + 1 >> 0) < upper.nd);
			if (okdown && okup) {
				d.Round(i + 1 >> 0);
				return;
			} else if (okdown) {
				d.RoundDown(i + 1 >> 0);
				return;
			} else if (okup) {
				d.RoundUp(i + 1 >> 0);
				return;
			}
			i = i + (1) >> 0;
		}
	};
	fmtE = function(dst, neg, d, prec, fmt) {
		var ch, x, i, m, x$1, exp, buf, i$1, _r, _q, _ref;
		if (neg) {
			dst = $append(dst, 45);
		}
		ch = 48;
		if (!((d.nd === 0))) {
			ch = (x = d.d, ((0 < 0 || 0 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + 0]));
		}
		dst = $append(dst, ch);
		if (prec > 0) {
			dst = $append(dst, 46);
			i = 1;
			m = ((d.nd + prec >> 0) + 1 >> 0) - max(d.nd, prec + 1 >> 0) >> 0;
			while (i < m) {
				dst = $append(dst, (x$1 = d.d, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i])));
				i = i + (1) >> 0;
			}
			while (i <= prec) {
				dst = $append(dst, 48);
				i = i + (1) >> 0;
			}
		}
		dst = $append(dst, fmt);
		exp = d.dp - 1 >> 0;
		if (d.nd === 0) {
			exp = 0;
		}
		if (exp < 0) {
			ch = 45;
			exp = -exp;
		} else {
			ch = 43;
		}
		dst = $append(dst, ch);
		buf = ($arrayType($Uint8, 3)).zero(); $copy(buf, ($arrayType($Uint8, 3)).zero(), ($arrayType($Uint8, 3)));
		i$1 = 3;
		while (exp >= 10) {
			i$1 = i$1 - (1) >> 0;
			(i$1 < 0 || i$1 >= buf.length) ? $throwRuntimeError("index out of range") : buf[i$1] = (((_r = exp % 10, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) + 48 >> 0) << 24 >>> 24);
			exp = (_q = exp / (10), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		}
		i$1 = i$1 - (1) >> 0;
		(i$1 < 0 || i$1 >= buf.length) ? $throwRuntimeError("index out of range") : buf[i$1] = ((exp + 48 >> 0) << 24 >>> 24);
		_ref = i$1;
		if (_ref === 0) {
			dst = $append(dst, buf[0], buf[1], buf[2]);
		} else if (_ref === 1) {
			dst = $append(dst, buf[1], buf[2]);
		} else if (_ref === 2) {
			dst = $append(dst, 48, buf[2]);
		}
		return dst;
	};
	fmtF = function(dst, neg, d, prec) {
		var i, x, i$1, ch, j, x$1;
		if (neg) {
			dst = $append(dst, 45);
		}
		if (d.dp > 0) {
			i = 0;
			i = 0;
			while (i < d.dp && i < d.nd) {
				dst = $append(dst, (x = d.d, ((i < 0 || i >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + i])));
				i = i + (1) >> 0;
			}
			while (i < d.dp) {
				dst = $append(dst, 48);
				i = i + (1) >> 0;
			}
		} else {
			dst = $append(dst, 48);
		}
		if (prec > 0) {
			dst = $append(dst, 46);
			i$1 = 0;
			while (i$1 < prec) {
				ch = 48;
				j = d.dp + i$1 >> 0;
				if (0 <= j && j < d.nd) {
					ch = (x$1 = d.d, ((j < 0 || j >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + j]));
				}
				dst = $append(dst, ch);
				i$1 = i$1 + (1) >> 0;
			}
		}
		return dst;
	};
	fmtB = function(dst, neg, mant, exp, flt) {
		var buf, w, esign, n, _r, _q, x;
		buf = ($arrayType($Uint8, 50)).zero(); $copy(buf, ($arrayType($Uint8, 50)).zero(), ($arrayType($Uint8, 50)));
		w = 50;
		exp = exp - ((flt.mantbits >> 0)) >> 0;
		esign = 43;
		if (exp < 0) {
			esign = 45;
			exp = -exp;
		}
		n = 0;
		while (exp > 0 || n < 1) {
			n = n + (1) >> 0;
			w = w - (1) >> 0;
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = (((_r = exp % 10, _r === _r ? _r : $throwRuntimeError("integer divide by zero")) + 48 >> 0) << 24 >>> 24);
			exp = (_q = exp / (10), (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		}
		w = w - (1) >> 0;
		(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = esign;
		w = w - (1) >> 0;
		(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 112;
		n = 0;
		while ((mant.$high > 0 || (mant.$high === 0 && mant.$low > 0)) || n < 1) {
			n = n + (1) >> 0;
			w = w - (1) >> 0;
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = ((x = $div64(mant, new $Uint64(0, 10), true), new $Uint64(x.$high + 0, x.$low + 48)).$low << 24 >>> 24);
			mant = $div64(mant, (new $Uint64(0, 10)), false);
		}
		if (neg) {
			w = w - (1) >> 0;
			(w < 0 || w >= buf.length) ? $throwRuntimeError("index out of range") : buf[w] = 45;
		}
		return $appendSlice(dst, $subslice(new ($sliceType($Uint8))(buf), w));
	};
	max = function(a, b) {
		if (a > b) {
			return a;
		}
		return b;
	};
	FormatInt = $pkg.FormatInt = function(i, base) {
		var _tuple, s;
		_tuple = formatBits(($sliceType($Uint8)).nil, new $Uint64(i.$high, i.$low), base, (i.$high < 0 || (i.$high === 0 && i.$low < 0)), false); s = _tuple[1];
		return s;
	};
	Itoa = $pkg.Itoa = function(i) {
		return FormatInt(new $Int64(0, i), 10);
	};
	formatBits = function(dst, u, base, neg, append_) {
		var d = ($sliceType($Uint8)).nil, s = "", a, i, q, x, j, x$1, x$2, q$1, x$3, s$1, b, m, b$1;
		if (base < 2 || base > 36) {
			$panic(new $String("strconv: illegal AppendInt/FormatInt base"));
		}
		a = ($arrayType($Uint8, 65)).zero(); $copy(a, ($arrayType($Uint8, 65)).zero(), ($arrayType($Uint8, 65)));
		i = 65;
		if (neg) {
			u = new $Uint64(-u.$high, -u.$low);
		}
		if (base === 10) {
			while ((u.$high > 0 || (u.$high === 0 && u.$low >= 100))) {
				i = i - (2) >> 0;
				q = $div64(u, new $Uint64(0, 100), false);
				j = ((x = $mul64(q, new $Uint64(0, 100)), new $Uint64(u.$high - x.$high, u.$low - x.$low)).$low >>> 0);
				(x$1 = i + 1 >> 0, (x$1 < 0 || x$1 >= a.length) ? $throwRuntimeError("index out of range") : a[x$1] = "0123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789".charCodeAt(j));
				(x$2 = i + 0 >> 0, (x$2 < 0 || x$2 >= a.length) ? $throwRuntimeError("index out of range") : a[x$2] = "0000000000111111111122222222223333333333444444444455555555556666666666777777777788888888889999999999".charCodeAt(j));
				u = q;
			}
			if ((u.$high > 0 || (u.$high === 0 && u.$low >= 10))) {
				i = i - (1) >> 0;
				q$1 = $div64(u, new $Uint64(0, 10), false);
				(i < 0 || i >= a.length) ? $throwRuntimeError("index out of range") : a[i] = "0123456789abcdefghijklmnopqrstuvwxyz".charCodeAt(((x$3 = $mul64(q$1, new $Uint64(0, 10)), new $Uint64(u.$high - x$3.$high, u.$low - x$3.$low)).$low >>> 0));
				u = q$1;
			}
		} else {
			s$1 = ((base < 0 || base >= shifts.length) ? $throwRuntimeError("index out of range") : shifts[base]);
			if (s$1 > 0) {
				b = new $Uint64(0, base);
				m = (b.$low >>> 0) - 1 >>> 0;
				while ((u.$high > b.$high || (u.$high === b.$high && u.$low >= b.$low))) {
					i = i - (1) >> 0;
					(i < 0 || i >= a.length) ? $throwRuntimeError("index out of range") : a[i] = "0123456789abcdefghijklmnopqrstuvwxyz".charCodeAt((((u.$low >>> 0) & m) >>> 0));
					u = $shiftRightUint64(u, (s$1));
				}
			} else {
				b$1 = new $Uint64(0, base);
				while ((u.$high > b$1.$high || (u.$high === b$1.$high && u.$low >= b$1.$low))) {
					i = i - (1) >> 0;
					(i < 0 || i >= a.length) ? $throwRuntimeError("index out of range") : a[i] = "0123456789abcdefghijklmnopqrstuvwxyz".charCodeAt(($div64(u, b$1, true).$low >>> 0));
					u = $div64(u, (b$1), false);
				}
			}
		}
		i = i - (1) >> 0;
		(i < 0 || i >= a.length) ? $throwRuntimeError("index out of range") : a[i] = "0123456789abcdefghijklmnopqrstuvwxyz".charCodeAt((u.$low >>> 0));
		if (neg) {
			i = i - (1) >> 0;
			(i < 0 || i >= a.length) ? $throwRuntimeError("index out of range") : a[i] = 45;
		}
		if (append_) {
			d = $appendSlice(dst, $subslice(new ($sliceType($Uint8))(a), i));
			return [d, s];
		}
		s = $bytesToString($subslice(new ($sliceType($Uint8))(a), i));
		return [d, s];
	};
	quoteWith = function(s, quote, ASCIIonly) {
		var runeTmp, _q, x, buf, width, r, _tuple, n, _ref, s$1, s$2;
		runeTmp = ($arrayType($Uint8, 4)).zero(); $copy(runeTmp, ($arrayType($Uint8, 4)).zero(), ($arrayType($Uint8, 4)));
		buf = ($sliceType($Uint8)).make(0, (_q = (x = s.length, (((3 >>> 16 << 16) * x >> 0) + (3 << 16 >>> 16) * x) >> 0) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")));
		buf = $append(buf, quote);
		width = 0;
		while (s.length > 0) {
			r = (s.charCodeAt(0) >> 0);
			width = 1;
			if (r >= 128) {
				_tuple = utf8.DecodeRuneInString(s); r = _tuple[0]; width = _tuple[1];
			}
			if ((width === 1) && (r === 65533)) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\x")));
				buf = $append(buf, "0123456789abcdef".charCodeAt((s.charCodeAt(0) >>> 4 << 24 >>> 24)));
				buf = $append(buf, "0123456789abcdef".charCodeAt(((s.charCodeAt(0) & 15) >>> 0)));
				s = s.substring(width);
				continue;
			}
			if ((r === (quote >> 0)) || (r === 92)) {
				buf = $append(buf, 92);
				buf = $append(buf, (r << 24 >>> 24));
				s = s.substring(width);
				continue;
			}
			if (ASCIIonly) {
				if (r < 128 && IsPrint(r)) {
					buf = $append(buf, (r << 24 >>> 24));
					s = s.substring(width);
					continue;
				}
			} else if (IsPrint(r)) {
				n = utf8.EncodeRune(new ($sliceType($Uint8))(runeTmp), r);
				buf = $appendSlice(buf, $subslice(new ($sliceType($Uint8))(runeTmp), 0, n));
				s = s.substring(width);
				continue;
			}
			_ref = r;
			if (_ref === 7) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\a")));
			} else if (_ref === 8) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\b")));
			} else if (_ref === 12) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\f")));
			} else if (_ref === 10) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\n")));
			} else if (_ref === 13) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\r")));
			} else if (_ref === 9) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\t")));
			} else if (_ref === 11) {
				buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\v")));
			} else {
				if (r < 32) {
					buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\x")));
					buf = $append(buf, "0123456789abcdef".charCodeAt((s.charCodeAt(0) >>> 4 << 24 >>> 24)));
					buf = $append(buf, "0123456789abcdef".charCodeAt(((s.charCodeAt(0) & 15) >>> 0)));
				} else if (r > 1114111) {
					r = 65533;
					buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\u")));
					s$1 = 12;
					while (s$1 >= 0) {
						buf = $append(buf, "0123456789abcdef".charCodeAt((((r >> $min((s$1 >>> 0), 31)) >> 0) & 15)));
						s$1 = s$1 - (4) >> 0;
					}
				} else if (r < 65536) {
					buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\u")));
					s$1 = 12;
					while (s$1 >= 0) {
						buf = $append(buf, "0123456789abcdef".charCodeAt((((r >> $min((s$1 >>> 0), 31)) >> 0) & 15)));
						s$1 = s$1 - (4) >> 0;
					}
				} else {
					buf = $appendSlice(buf, new ($sliceType($Uint8))($stringToBytes("\\U")));
					s$2 = 28;
					while (s$2 >= 0) {
						buf = $append(buf, "0123456789abcdef".charCodeAt((((r >> $min((s$2 >>> 0), 31)) >> 0) & 15)));
						s$2 = s$2 - (4) >> 0;
					}
				}
			}
			s = s.substring(width);
		}
		buf = $append(buf, quote);
		return $bytesToString(buf);
	};
	Quote = $pkg.Quote = function(s) {
		return quoteWith(s, 34, false);
	};
	QuoteToASCII = $pkg.QuoteToASCII = function(s) {
		return quoteWith(s, 34, true);
	};
	QuoteRune = $pkg.QuoteRune = function(r) {
		return quoteWith($encodeRune(r), 39, false);
	};
	AppendQuoteRune = $pkg.AppendQuoteRune = function(dst, r) {
		return $appendSlice(dst, new ($sliceType($Uint8))($stringToBytes(QuoteRune(r))));
	};
	QuoteRuneToASCII = $pkg.QuoteRuneToASCII = function(r) {
		return quoteWith($encodeRune(r), 39, true);
	};
	AppendQuoteRuneToASCII = $pkg.AppendQuoteRuneToASCII = function(dst, r) {
		return $appendSlice(dst, new ($sliceType($Uint8))($stringToBytes(QuoteRuneToASCII(r))));
	};
	CanBackquote = $pkg.CanBackquote = function(s) {
		var i, c;
		i = 0;
		while (i < s.length) {
			c = s.charCodeAt(i);
			if ((c < 32 && !((c === 9))) || (c === 96) || (c === 127)) {
				return false;
			}
			i = i + (1) >> 0;
		}
		return true;
	};
	unhex = function(b) {
		var v = 0, ok = false, c, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5;
		c = (b >> 0);
		if (48 <= c && c <= 57) {
			_tmp = c - 48 >> 0; _tmp$1 = true; v = _tmp; ok = _tmp$1;
			return [v, ok];
		} else if (97 <= c && c <= 102) {
			_tmp$2 = (c - 97 >> 0) + 10 >> 0; _tmp$3 = true; v = _tmp$2; ok = _tmp$3;
			return [v, ok];
		} else if (65 <= c && c <= 70) {
			_tmp$4 = (c - 65 >> 0) + 10 >> 0; _tmp$5 = true; v = _tmp$4; ok = _tmp$5;
			return [v, ok];
		}
		return [v, ok];
	};
	UnquoteChar = $pkg.UnquoteChar = function(s, quote) {
		var value = 0, multibyte = false, tail = "", err = null, c, _tuple, r, size, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, c$1, _ref, n, _ref$1, v, j, _tuple$1, x, ok, v$1, j$1, x$1;
		c = s.charCodeAt(0);
		if ((c === quote) && ((quote === 39) || (quote === 34))) {
			err = $pkg.ErrSyntax;
			return [value, multibyte, tail, err];
		} else if (c >= 128) {
			_tuple = utf8.DecodeRuneInString(s); r = _tuple[0]; size = _tuple[1];
			_tmp = r; _tmp$1 = true; _tmp$2 = s.substring(size); _tmp$3 = null; value = _tmp; multibyte = _tmp$1; tail = _tmp$2; err = _tmp$3;
			return [value, multibyte, tail, err];
		} else if (!((c === 92))) {
			_tmp$4 = (s.charCodeAt(0) >> 0); _tmp$5 = false; _tmp$6 = s.substring(1); _tmp$7 = null; value = _tmp$4; multibyte = _tmp$5; tail = _tmp$6; err = _tmp$7;
			return [value, multibyte, tail, err];
		}
		if (s.length <= 1) {
			err = $pkg.ErrSyntax;
			return [value, multibyte, tail, err];
		}
		c$1 = s.charCodeAt(1);
		s = s.substring(2);
		_ref = c$1;
		switch (0) { default: if (_ref === 97) {
			value = 7;
		} else if (_ref === 98) {
			value = 8;
		} else if (_ref === 102) {
			value = 12;
		} else if (_ref === 110) {
			value = 10;
		} else if (_ref === 114) {
			value = 13;
		} else if (_ref === 116) {
			value = 9;
		} else if (_ref === 118) {
			value = 11;
		} else if (_ref === 120 || _ref === 117 || _ref === 85) {
			n = 0;
			_ref$1 = c$1;
			if (_ref$1 === 120) {
				n = 2;
			} else if (_ref$1 === 117) {
				n = 4;
			} else if (_ref$1 === 85) {
				n = 8;
			}
			v = 0;
			if (s.length < n) {
				err = $pkg.ErrSyntax;
				return [value, multibyte, tail, err];
			}
			j = 0;
			while (j < n) {
				_tuple$1 = unhex(s.charCodeAt(j)); x = _tuple$1[0]; ok = _tuple$1[1];
				if (!ok) {
					err = $pkg.ErrSyntax;
					return [value, multibyte, tail, err];
				}
				v = (v << 4 >> 0) | x;
				j = j + (1) >> 0;
			}
			s = s.substring(n);
			if (c$1 === 120) {
				value = v;
				break;
			}
			if (v > 1114111) {
				err = $pkg.ErrSyntax;
				return [value, multibyte, tail, err];
			}
			value = v;
			multibyte = true;
		} else if (_ref === 48 || _ref === 49 || _ref === 50 || _ref === 51 || _ref === 52 || _ref === 53 || _ref === 54 || _ref === 55) {
			v$1 = (c$1 >> 0) - 48 >> 0;
			if (s.length < 2) {
				err = $pkg.ErrSyntax;
				return [value, multibyte, tail, err];
			}
			j$1 = 0;
			while (j$1 < 2) {
				x$1 = (s.charCodeAt(j$1) >> 0) - 48 >> 0;
				if (x$1 < 0 || x$1 > 7) {
					err = $pkg.ErrSyntax;
					return [value, multibyte, tail, err];
				}
				v$1 = ((v$1 << 3 >> 0)) | x$1;
				j$1 = j$1 + (1) >> 0;
			}
			s = s.substring(2);
			if (v$1 > 255) {
				err = $pkg.ErrSyntax;
				return [value, multibyte, tail, err];
			}
			value = v$1;
		} else if (_ref === 92) {
			value = 92;
		} else if (_ref === 39 || _ref === 34) {
			if (!((c$1 === quote))) {
				err = $pkg.ErrSyntax;
				return [value, multibyte, tail, err];
			}
			value = (c$1 >> 0);
		} else {
			err = $pkg.ErrSyntax;
			return [value, multibyte, tail, err];
		} }
		tail = s;
		return [value, multibyte, tail, err];
	};
	Unquote = $pkg.Unquote = function(s) {
		var t = "", err = null, n, _tmp, _tmp$1, quote, _tmp$2, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10, _tmp$11, _ref, _tmp$12, _tmp$13, _tuple, r, size, _tmp$14, _tmp$15, runeTmp, _q, x, buf, _tuple$1, c, multibyte, ss, err$1, _tmp$16, _tmp$17, n$1, _tmp$18, _tmp$19, _tmp$20, _tmp$21;
		n = s.length;
		if (n < 2) {
			_tmp = ""; _tmp$1 = $pkg.ErrSyntax; t = _tmp; err = _tmp$1;
			return [t, err];
		}
		quote = s.charCodeAt(0);
		if (!((quote === s.charCodeAt((n - 1 >> 0))))) {
			_tmp$2 = ""; _tmp$3 = $pkg.ErrSyntax; t = _tmp$2; err = _tmp$3;
			return [t, err];
		}
		s = s.substring(1, (n - 1 >> 0));
		if (quote === 96) {
			if (contains(s, 96)) {
				_tmp$4 = ""; _tmp$5 = $pkg.ErrSyntax; t = _tmp$4; err = _tmp$5;
				return [t, err];
			}
			_tmp$6 = s; _tmp$7 = null; t = _tmp$6; err = _tmp$7;
			return [t, err];
		}
		if (!((quote === 34)) && !((quote === 39))) {
			_tmp$8 = ""; _tmp$9 = $pkg.ErrSyntax; t = _tmp$8; err = _tmp$9;
			return [t, err];
		}
		if (contains(s, 10)) {
			_tmp$10 = ""; _tmp$11 = $pkg.ErrSyntax; t = _tmp$10; err = _tmp$11;
			return [t, err];
		}
		if (!contains(s, 92) && !contains(s, quote)) {
			_ref = quote;
			if (_ref === 34) {
				_tmp$12 = s; _tmp$13 = null; t = _tmp$12; err = _tmp$13;
				return [t, err];
			} else if (_ref === 39) {
				_tuple = utf8.DecodeRuneInString(s); r = _tuple[0]; size = _tuple[1];
				if ((size === s.length) && (!((r === 65533)) || !((size === 1)))) {
					_tmp$14 = s; _tmp$15 = null; t = _tmp$14; err = _tmp$15;
					return [t, err];
				}
			}
		}
		runeTmp = ($arrayType($Uint8, 4)).zero(); $copy(runeTmp, ($arrayType($Uint8, 4)).zero(), ($arrayType($Uint8, 4)));
		buf = ($sliceType($Uint8)).make(0, (_q = (x = s.length, (((3 >>> 16 << 16) * x >> 0) + (3 << 16 >>> 16) * x) >> 0) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")));
		while (s.length > 0) {
			_tuple$1 = UnquoteChar(s, quote); c = _tuple$1[0]; multibyte = _tuple$1[1]; ss = _tuple$1[2]; err$1 = _tuple$1[3];
			if (!($interfaceIsEqual(err$1, null))) {
				_tmp$16 = ""; _tmp$17 = err$1; t = _tmp$16; err = _tmp$17;
				return [t, err];
			}
			s = ss;
			if (c < 128 || !multibyte) {
				buf = $append(buf, (c << 24 >>> 24));
			} else {
				n$1 = utf8.EncodeRune(new ($sliceType($Uint8))(runeTmp), c);
				buf = $appendSlice(buf, $subslice(new ($sliceType($Uint8))(runeTmp), 0, n$1));
			}
			if ((quote === 39) && !((s.length === 0))) {
				_tmp$18 = ""; _tmp$19 = $pkg.ErrSyntax; t = _tmp$18; err = _tmp$19;
				return [t, err];
			}
		}
		_tmp$20 = $bytesToString(buf); _tmp$21 = null; t = _tmp$20; err = _tmp$21;
		return [t, err];
	};
	contains = function(s, c) {
		var i;
		i = 0;
		while (i < s.length) {
			if (s.charCodeAt(i) === c) {
				return true;
			}
			i = i + (1) >> 0;
		}
		return false;
	};
	bsearch16 = function(a, x) {
		var _tmp, _tmp$1, i, j, _q, h;
		_tmp = 0; _tmp$1 = a.$length; i = _tmp; j = _tmp$1;
		while (i < j) {
			h = i + (_q = ((j - i >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			if (((h < 0 || h >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + h]) < x) {
				i = h + 1 >> 0;
			} else {
				j = h;
			}
		}
		return i;
	};
	bsearch32 = function(a, x) {
		var _tmp, _tmp$1, i, j, _q, h;
		_tmp = 0; _tmp$1 = a.$length; i = _tmp; j = _tmp$1;
		while (i < j) {
			h = i + (_q = ((j - i >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			if (((h < 0 || h >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + h]) < x) {
				i = h + 1 >> 0;
			} else {
				j = h;
			}
		}
		return i;
	};
	IsPrint = $pkg.IsPrint = function(r) {
		var _tmp, _tmp$1, _tmp$2, rr, isPrint, isNotPrint, i, x, x$1, j, _tmp$3, _tmp$4, _tmp$5, rr$1, isPrint$1, isNotPrint$1, i$1, x$2, x$3, j$1;
		if (r <= 255) {
			if (32 <= r && r <= 126) {
				return true;
			}
			if (161 <= r && r <= 255) {
				return !((r === 173));
			}
			return false;
		}
		if (0 <= r && r < 65536) {
			_tmp = (r << 16 >>> 16); _tmp$1 = isPrint16; _tmp$2 = isNotPrint16; rr = _tmp; isPrint = _tmp$1; isNotPrint = _tmp$2;
			i = bsearch16(isPrint, rr);
			if (i >= isPrint.$length || rr < (x = i & ~1, ((x < 0 || x >= isPrint.$length) ? $throwRuntimeError("index out of range") : isPrint.$array[isPrint.$offset + x])) || (x$1 = i | 1, ((x$1 < 0 || x$1 >= isPrint.$length) ? $throwRuntimeError("index out of range") : isPrint.$array[isPrint.$offset + x$1])) < rr) {
				return false;
			}
			j = bsearch16(isNotPrint, rr);
			return j >= isNotPrint.$length || !((((j < 0 || j >= isNotPrint.$length) ? $throwRuntimeError("index out of range") : isNotPrint.$array[isNotPrint.$offset + j]) === rr));
		}
		_tmp$3 = (r >>> 0); _tmp$4 = isPrint32; _tmp$5 = isNotPrint32; rr$1 = _tmp$3; isPrint$1 = _tmp$4; isNotPrint$1 = _tmp$5;
		i$1 = bsearch32(isPrint$1, rr$1);
		if (i$1 >= isPrint$1.$length || rr$1 < (x$2 = i$1 & ~1, ((x$2 < 0 || x$2 >= isPrint$1.$length) ? $throwRuntimeError("index out of range") : isPrint$1.$array[isPrint$1.$offset + x$2])) || (x$3 = i$1 | 1, ((x$3 < 0 || x$3 >= isPrint$1.$length) ? $throwRuntimeError("index out of range") : isPrint$1.$array[isPrint$1.$offset + x$3])) < rr$1) {
			return false;
		}
		if (r >= 131072) {
			return true;
		}
		r = r - (65536) >> 0;
		j$1 = bsearch16(isNotPrint$1, (r << 16 >>> 16));
		return j$1 >= isNotPrint$1.$length || !((((j$1 < 0 || j$1 >= isNotPrint$1.$length) ? $throwRuntimeError("index out of range") : isNotPrint$1.$array[isNotPrint$1.$offset + j$1]) === (r << 16 >>> 16)));
	};
	$pkg.$init = function() {
		($ptrType(decimal)).methods = [["Assign", "Assign", "", [$Uint64], [], false, -1], ["Round", "Round", "", [$Int], [], false, -1], ["RoundDown", "RoundDown", "", [$Int], [], false, -1], ["RoundUp", "RoundUp", "", [$Int], [], false, -1], ["RoundedInteger", "RoundedInteger", "", [], [$Uint64], false, -1], ["Shift", "Shift", "", [$Int], [], false, -1], ["String", "String", "", [], [$String], false, -1], ["floatBits", "floatBits", "strconv", [($ptrType(floatInfo))], [$Uint64, $Bool], false, -1], ["set", "set", "strconv", [$String], [$Bool], false, -1]];
		decimal.init([["d", "d", "strconv", ($arrayType($Uint8, 800)), ""], ["nd", "nd", "strconv", $Int, ""], ["dp", "dp", "strconv", $Int, ""], ["neg", "neg", "strconv", $Bool, ""], ["trunc", "trunc", "strconv", $Bool, ""]]);
		leftCheat.init([["delta", "delta", "strconv", $Int, ""], ["cutoff", "cutoff", "strconv", $String, ""]]);
		($ptrType(extFloat)).methods = [["AssignComputeBounds", "AssignComputeBounds", "", [$Uint64, $Int, $Bool, ($ptrType(floatInfo))], [extFloat, extFloat], false, -1], ["AssignDecimal", "AssignDecimal", "", [$Uint64, $Int, $Bool, $Bool, ($ptrType(floatInfo))], [$Bool], false, -1], ["FixedDecimal", "FixedDecimal", "", [($ptrType(decimalSlice)), $Int], [$Bool], false, -1], ["Multiply", "Multiply", "", [extFloat], [], false, -1], ["Normalize", "Normalize", "", [], [$Uint], false, -1], ["ShortestDecimal", "ShortestDecimal", "", [($ptrType(decimalSlice)), ($ptrType(extFloat)), ($ptrType(extFloat))], [$Bool], false, -1], ["floatBits", "floatBits", "strconv", [($ptrType(floatInfo))], [$Uint64, $Bool], false, -1], ["frexp10", "frexp10", "strconv", [], [$Int, $Int], false, -1]];
		extFloat.init([["mant", "mant", "strconv", $Uint64, ""], ["exp", "exp", "strconv", $Int, ""], ["neg", "neg", "strconv", $Bool, ""]]);
		floatInfo.init([["mantbits", "mantbits", "strconv", $Uint, ""], ["expbits", "expbits", "strconv", $Uint, ""], ["bias", "bias", "strconv", $Int, ""]]);
		decimalSlice.init([["d", "d", "strconv", ($sliceType($Uint8)), ""], ["nd", "nd", "strconv", $Int, ""], ["dp", "dp", "strconv", $Int, ""], ["neg", "neg", "strconv", $Bool, ""]]);
		optimize = true;
		$pkg.ErrRange = errors.New("value out of range");
		$pkg.ErrSyntax = errors.New("invalid syntax");
		leftcheats = new ($sliceType(leftCheat))([new leftCheat.Ptr(0, ""), new leftCheat.Ptr(1, "5"), new leftCheat.Ptr(1, "25"), new leftCheat.Ptr(1, "125"), new leftCheat.Ptr(2, "625"), new leftCheat.Ptr(2, "3125"), new leftCheat.Ptr(2, "15625"), new leftCheat.Ptr(3, "78125"), new leftCheat.Ptr(3, "390625"), new leftCheat.Ptr(3, "1953125"), new leftCheat.Ptr(4, "9765625"), new leftCheat.Ptr(4, "48828125"), new leftCheat.Ptr(4, "244140625"), new leftCheat.Ptr(4, "1220703125"), new leftCheat.Ptr(5, "6103515625"), new leftCheat.Ptr(5, "30517578125"), new leftCheat.Ptr(5, "152587890625"), new leftCheat.Ptr(6, "762939453125"), new leftCheat.Ptr(6, "3814697265625"), new leftCheat.Ptr(6, "19073486328125"), new leftCheat.Ptr(7, "95367431640625"), new leftCheat.Ptr(7, "476837158203125"), new leftCheat.Ptr(7, "2384185791015625"), new leftCheat.Ptr(7, "11920928955078125"), new leftCheat.Ptr(8, "59604644775390625"), new leftCheat.Ptr(8, "298023223876953125"), new leftCheat.Ptr(8, "1490116119384765625"), new leftCheat.Ptr(9, "7450580596923828125")]);
		smallPowersOfTen = $toNativeArray("Struct", [new extFloat.Ptr(new $Uint64(2147483648, 0), -63, false), new extFloat.Ptr(new $Uint64(2684354560, 0), -60, false), new extFloat.Ptr(new $Uint64(3355443200, 0), -57, false), new extFloat.Ptr(new $Uint64(4194304000, 0), -54, false), new extFloat.Ptr(new $Uint64(2621440000, 0), -50, false), new extFloat.Ptr(new $Uint64(3276800000, 0), -47, false), new extFloat.Ptr(new $Uint64(4096000000, 0), -44, false), new extFloat.Ptr(new $Uint64(2560000000, 0), -40, false)]);
		powersOfTen = $toNativeArray("Struct", [new extFloat.Ptr(new $Uint64(4203730336, 136053384), -1220, false), new extFloat.Ptr(new $Uint64(3132023167, 2722021238), -1193, false), new extFloat.Ptr(new $Uint64(2333539104, 810921078), -1166, false), new extFloat.Ptr(new $Uint64(3477244234, 1573795306), -1140, false), new extFloat.Ptr(new $Uint64(2590748842, 1432697645), -1113, false), new extFloat.Ptr(new $Uint64(3860516611, 1025131999), -1087, false), new extFloat.Ptr(new $Uint64(2876309015, 3348809418), -1060, false), new extFloat.Ptr(new $Uint64(4286034428, 3200048207), -1034, false), new extFloat.Ptr(new $Uint64(3193344495, 1097586188), -1007, false), new extFloat.Ptr(new $Uint64(2379227053, 2424306748), -980, false), new extFloat.Ptr(new $Uint64(3545324584, 827693699), -954, false), new extFloat.Ptr(new $Uint64(2641472655, 2913388981), -927, false), new extFloat.Ptr(new $Uint64(3936100983, 602835915), -901, false), new extFloat.Ptr(new $Uint64(2932623761, 1081627501), -874, false), new extFloat.Ptr(new $Uint64(2184974969, 1572261463), -847, false), new extFloat.Ptr(new $Uint64(3255866422, 1308317239), -821, false), new extFloat.Ptr(new $Uint64(2425809519, 944281679), -794, false), new extFloat.Ptr(new $Uint64(3614737867, 629291719), -768, false), new extFloat.Ptr(new $Uint64(2693189581, 2545915892), -741, false), new extFloat.Ptr(new $Uint64(4013165208, 388672741), -715, false), new extFloat.Ptr(new $Uint64(2990041083, 708162190), -688, false), new extFloat.Ptr(new $Uint64(2227754207, 3536207675), -661, false), new extFloat.Ptr(new $Uint64(3319612455, 450088378), -635, false), new extFloat.Ptr(new $Uint64(2473304014, 3139815830), -608, false), new extFloat.Ptr(new $Uint64(3685510180, 2103616900), -582, false), new extFloat.Ptr(new $Uint64(2745919064, 224385782), -555, false), new extFloat.Ptr(new $Uint64(4091738259, 3737383206), -529, false), new extFloat.Ptr(new $Uint64(3048582568, 2868871352), -502, false), new extFloat.Ptr(new $Uint64(2271371013, 1820084875), -475, false), new extFloat.Ptr(new $Uint64(3384606560, 885076051), -449, false), new extFloat.Ptr(new $Uint64(2521728396, 2444895829), -422, false), new extFloat.Ptr(new $Uint64(3757668132, 1881767613), -396, false), new extFloat.Ptr(new $Uint64(2799680927, 3102062735), -369, false), new extFloat.Ptr(new $Uint64(4171849679, 2289335700), -343, false), new extFloat.Ptr(new $Uint64(3108270227, 2410191823), -316, false), new extFloat.Ptr(new $Uint64(2315841784, 3205436779), -289, false), new extFloat.Ptr(new $Uint64(3450873173, 1697722806), -263, false), new extFloat.Ptr(new $Uint64(2571100870, 3497754540), -236, false), new extFloat.Ptr(new $Uint64(3831238852, 707476230), -210, false), new extFloat.Ptr(new $Uint64(2854495385, 1769181907), -183, false), new extFloat.Ptr(new $Uint64(4253529586, 2197867022), -157, false), new extFloat.Ptr(new $Uint64(3169126500, 2450594539), -130, false), new extFloat.Ptr(new $Uint64(2361183241, 1867548876), -103, false), new extFloat.Ptr(new $Uint64(3518437208, 3793315116), -77, false), new extFloat.Ptr(new $Uint64(2621440000, 0), -50, false), new extFloat.Ptr(new $Uint64(3906250000, 0), -24, false), new extFloat.Ptr(new $Uint64(2910383045, 2892103680), 3, false), new extFloat.Ptr(new $Uint64(2168404344, 4170451332), 30, false), new extFloat.Ptr(new $Uint64(3231174267, 3372684723), 56, false), new extFloat.Ptr(new $Uint64(2407412430, 2078956656), 83, false), new extFloat.Ptr(new $Uint64(3587324068, 2884206696), 109, false), new extFloat.Ptr(new $Uint64(2672764710, 395977285), 136, false), new extFloat.Ptr(new $Uint64(3982729777, 3569679143), 162, false), new extFloat.Ptr(new $Uint64(2967364920, 2361961896), 189, false), new extFloat.Ptr(new $Uint64(2210859150, 447440347), 216, false), new extFloat.Ptr(new $Uint64(3294436857, 1114709402), 242, false), new extFloat.Ptr(new $Uint64(2454546732, 2786846552), 269, false), new extFloat.Ptr(new $Uint64(3657559652, 443583978), 295, false), new extFloat.Ptr(new $Uint64(2725094297, 2599384906), 322, false), new extFloat.Ptr(new $Uint64(4060706939, 3028118405), 348, false), new extFloat.Ptr(new $Uint64(3025462433, 2044532855), 375, false), new extFloat.Ptr(new $Uint64(2254145170, 1536935362), 402, false), new extFloat.Ptr(new $Uint64(3358938053, 3365297469), 428, false), new extFloat.Ptr(new $Uint64(2502603868, 4204241075), 455, false), new extFloat.Ptr(new $Uint64(3729170365, 2577424355), 481, false), new extFloat.Ptr(new $Uint64(2778448436, 3677981733), 508, false), new extFloat.Ptr(new $Uint64(4140210802, 2744688476), 534, false), new extFloat.Ptr(new $Uint64(3084697427, 1424604878), 561, false), new extFloat.Ptr(new $Uint64(2298278679, 4062331362), 588, false), new extFloat.Ptr(new $Uint64(3424702107, 3546052773), 614, false), new extFloat.Ptr(new $Uint64(2551601907, 2065781727), 641, false), new extFloat.Ptr(new $Uint64(3802183132, 2535403578), 667, false), new extFloat.Ptr(new $Uint64(2832847187, 1558426518), 694, false), new extFloat.Ptr(new $Uint64(4221271257, 2762425404), 720, false), new extFloat.Ptr(new $Uint64(3145092172, 2812560400), 747, false), new extFloat.Ptr(new $Uint64(2343276271, 3057687578), 774, false), new extFloat.Ptr(new $Uint64(3491753744, 2790753324), 800, false), new extFloat.Ptr(new $Uint64(2601559269, 3918606633), 827, false), new extFloat.Ptr(new $Uint64(3876625403, 2711358621), 853, false), new extFloat.Ptr(new $Uint64(2888311001, 1648096297), 880, false), new extFloat.Ptr(new $Uint64(2151959390, 2057817989), 907, false), new extFloat.Ptr(new $Uint64(3206669376, 61660461), 933, false), new extFloat.Ptr(new $Uint64(2389154863, 1581580175), 960, false), new extFloat.Ptr(new $Uint64(3560118173, 2626467905), 986, false), new extFloat.Ptr(new $Uint64(2652494738, 3034782633), 1013, false), new extFloat.Ptr(new $Uint64(3952525166, 3135207385), 1039, false), new extFloat.Ptr(new $Uint64(2944860731, 2616258155), 1066, false)]);
		uint64pow10 = $toNativeArray("Uint64", [new $Uint64(0, 1), new $Uint64(0, 10), new $Uint64(0, 100), new $Uint64(0, 1000), new $Uint64(0, 10000), new $Uint64(0, 100000), new $Uint64(0, 1000000), new $Uint64(0, 10000000), new $Uint64(0, 100000000), new $Uint64(0, 1000000000), new $Uint64(2, 1410065408), new $Uint64(23, 1215752192), new $Uint64(232, 3567587328), new $Uint64(2328, 1316134912), new $Uint64(23283, 276447232), new $Uint64(232830, 2764472320), new $Uint64(2328306, 1874919424), new $Uint64(23283064, 1569325056), new $Uint64(232830643, 2808348672), new $Uint64(2328306436, 2313682944)]);
		float32info = new floatInfo.Ptr(23, 8, -127);
		float64info = new floatInfo.Ptr(52, 11, -1023);
		isPrint16 = new ($sliceType($Uint16))([32, 126, 161, 887, 890, 894, 900, 1319, 1329, 1366, 1369, 1418, 1423, 1479, 1488, 1514, 1520, 1524, 1542, 1563, 1566, 1805, 1808, 1866, 1869, 1969, 1984, 2042, 2048, 2093, 2096, 2139, 2142, 2142, 2208, 2220, 2276, 2444, 2447, 2448, 2451, 2482, 2486, 2489, 2492, 2500, 2503, 2504, 2507, 2510, 2519, 2519, 2524, 2531, 2534, 2555, 2561, 2570, 2575, 2576, 2579, 2617, 2620, 2626, 2631, 2632, 2635, 2637, 2641, 2641, 2649, 2654, 2662, 2677, 2689, 2745, 2748, 2765, 2768, 2768, 2784, 2787, 2790, 2801, 2817, 2828, 2831, 2832, 2835, 2873, 2876, 2884, 2887, 2888, 2891, 2893, 2902, 2903, 2908, 2915, 2918, 2935, 2946, 2954, 2958, 2965, 2969, 2975, 2979, 2980, 2984, 2986, 2990, 3001, 3006, 3010, 3014, 3021, 3024, 3024, 3031, 3031, 3046, 3066, 3073, 3129, 3133, 3149, 3157, 3161, 3168, 3171, 3174, 3183, 3192, 3199, 3202, 3257, 3260, 3277, 3285, 3286, 3294, 3299, 3302, 3314, 3330, 3386, 3389, 3406, 3415, 3415, 3424, 3427, 3430, 3445, 3449, 3455, 3458, 3478, 3482, 3517, 3520, 3526, 3530, 3530, 3535, 3551, 3570, 3572, 3585, 3642, 3647, 3675, 3713, 3716, 3719, 3722, 3725, 3725, 3732, 3751, 3754, 3773, 3776, 3789, 3792, 3801, 3804, 3807, 3840, 3948, 3953, 4058, 4096, 4295, 4301, 4301, 4304, 4685, 4688, 4701, 4704, 4749, 4752, 4789, 4792, 4805, 4808, 4885, 4888, 4954, 4957, 4988, 4992, 5017, 5024, 5108, 5120, 5788, 5792, 5872, 5888, 5908, 5920, 5942, 5952, 5971, 5984, 6003, 6016, 6109, 6112, 6121, 6128, 6137, 6144, 6157, 6160, 6169, 6176, 6263, 6272, 6314, 6320, 6389, 6400, 6428, 6432, 6443, 6448, 6459, 6464, 6464, 6468, 6509, 6512, 6516, 6528, 6571, 6576, 6601, 6608, 6618, 6622, 6683, 6686, 6780, 6783, 6793, 6800, 6809, 6816, 6829, 6912, 6987, 6992, 7036, 7040, 7155, 7164, 7223, 7227, 7241, 7245, 7295, 7360, 7367, 7376, 7414, 7424, 7654, 7676, 7957, 7960, 7965, 7968, 8005, 8008, 8013, 8016, 8061, 8064, 8147, 8150, 8175, 8178, 8190, 8208, 8231, 8240, 8286, 8304, 8305, 8308, 8348, 8352, 8378, 8400, 8432, 8448, 8585, 8592, 9203, 9216, 9254, 9280, 9290, 9312, 11084, 11088, 11097, 11264, 11507, 11513, 11559, 11565, 11565, 11568, 11623, 11631, 11632, 11647, 11670, 11680, 11835, 11904, 12019, 12032, 12245, 12272, 12283, 12289, 12438, 12441, 12543, 12549, 12589, 12593, 12730, 12736, 12771, 12784, 19893, 19904, 40908, 40960, 42124, 42128, 42182, 42192, 42539, 42560, 42647, 42655, 42743, 42752, 42899, 42912, 42922, 43000, 43051, 43056, 43065, 43072, 43127, 43136, 43204, 43214, 43225, 43232, 43259, 43264, 43347, 43359, 43388, 43392, 43481, 43486, 43487, 43520, 43574, 43584, 43597, 43600, 43609, 43612, 43643, 43648, 43714, 43739, 43766, 43777, 43782, 43785, 43790, 43793, 43798, 43808, 43822, 43968, 44013, 44016, 44025, 44032, 55203, 55216, 55238, 55243, 55291, 63744, 64109, 64112, 64217, 64256, 64262, 64275, 64279, 64285, 64449, 64467, 64831, 64848, 64911, 64914, 64967, 65008, 65021, 65024, 65049, 65056, 65062, 65072, 65131, 65136, 65276, 65281, 65470, 65474, 65479, 65482, 65487, 65490, 65495, 65498, 65500, 65504, 65518, 65532, 65533]);
		isNotPrint16 = new ($sliceType($Uint16))([173, 907, 909, 930, 1376, 1416, 1424, 1757, 2111, 2209, 2303, 2424, 2432, 2436, 2473, 2481, 2526, 2564, 2601, 2609, 2612, 2615, 2621, 2653, 2692, 2702, 2706, 2729, 2737, 2740, 2758, 2762, 2820, 2857, 2865, 2868, 2910, 2948, 2961, 2971, 2973, 3017, 3076, 3085, 3089, 3113, 3124, 3141, 3145, 3159, 3204, 3213, 3217, 3241, 3252, 3269, 3273, 3295, 3312, 3332, 3341, 3345, 3397, 3401, 3460, 3506, 3516, 3541, 3543, 3715, 3721, 3736, 3744, 3748, 3750, 3756, 3770, 3781, 3783, 3912, 3992, 4029, 4045, 4294, 4681, 4695, 4697, 4745, 4785, 4799, 4801, 4823, 4881, 5760, 5901, 5997, 6001, 6751, 8024, 8026, 8028, 8030, 8117, 8133, 8156, 8181, 8335, 9984, 11311, 11359, 11558, 11687, 11695, 11703, 11711, 11719, 11727, 11735, 11743, 11930, 12352, 12687, 12831, 13055, 42895, 43470, 43815, 64311, 64317, 64319, 64322, 64325, 65107, 65127, 65141, 65511]);
		isPrint32 = new ($sliceType($Uint32))([65536, 65613, 65616, 65629, 65664, 65786, 65792, 65794, 65799, 65843, 65847, 65930, 65936, 65947, 66000, 66045, 66176, 66204, 66208, 66256, 66304, 66339, 66352, 66378, 66432, 66499, 66504, 66517, 66560, 66717, 66720, 66729, 67584, 67589, 67592, 67640, 67644, 67644, 67647, 67679, 67840, 67867, 67871, 67897, 67903, 67903, 67968, 68023, 68030, 68031, 68096, 68102, 68108, 68147, 68152, 68154, 68159, 68167, 68176, 68184, 68192, 68223, 68352, 68405, 68409, 68437, 68440, 68466, 68472, 68479, 68608, 68680, 69216, 69246, 69632, 69709, 69714, 69743, 69760, 69825, 69840, 69864, 69872, 69881, 69888, 69955, 70016, 70088, 70096, 70105, 71296, 71351, 71360, 71369, 73728, 74606, 74752, 74850, 74864, 74867, 77824, 78894, 92160, 92728, 93952, 94020, 94032, 94078, 94095, 94111, 110592, 110593, 118784, 119029, 119040, 119078, 119081, 119154, 119163, 119261, 119296, 119365, 119552, 119638, 119648, 119665, 119808, 119967, 119970, 119970, 119973, 119974, 119977, 120074, 120077, 120134, 120138, 120485, 120488, 120779, 120782, 120831, 126464, 126500, 126503, 126523, 126530, 126530, 126535, 126548, 126551, 126564, 126567, 126619, 126625, 126651, 126704, 126705, 126976, 127019, 127024, 127123, 127136, 127150, 127153, 127166, 127169, 127199, 127232, 127242, 127248, 127339, 127344, 127386, 127462, 127490, 127504, 127546, 127552, 127560, 127568, 127569, 127744, 127776, 127792, 127868, 127872, 127891, 127904, 127946, 127968, 127984, 128000, 128252, 128256, 128317, 128320, 128323, 128336, 128359, 128507, 128576, 128581, 128591, 128640, 128709, 128768, 128883, 131072, 173782, 173824, 177972, 177984, 178205, 194560, 195101, 917760, 917999]);
		isNotPrint32 = new ($sliceType($Uint16))([12, 39, 59, 62, 799, 926, 2057, 2102, 2134, 2564, 2580, 2584, 4285, 4405, 54357, 54429, 54445, 54458, 54460, 54468, 54534, 54549, 54557, 54586, 54591, 54597, 54609, 60932, 60960, 60963, 60968, 60979, 60984, 60986, 61000, 61002, 61004, 61008, 61011, 61016, 61018, 61020, 61022, 61024, 61027, 61035, 61043, 61048, 61053, 61055, 61066, 61092, 61098, 61648, 61743, 62262, 62405, 62527, 62529, 62712]);
		shifts = $toNativeArray("Uint", [0, 0, 1, 0, 2, 0, 0, 0, 3, 0, 0, 0, 0, 0, 0, 0, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 0, 0, 0, 0]);
	};
	return $pkg;
})();
$packages["reflect"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], runtime = $packages["runtime"], strconv = $packages["strconv"], sync = $packages["sync"], math = $packages["math"], mapIter, Type, Kind, rtype, method, uncommonType, ChanDir, arrayType, chanType, funcType, imethod, interfaceType, mapType, ptrType, sliceType, structField, structType, Method, StructField, StructTag, fieldScan, Value, flag, ValueError, iword, nonEmptyInterface, initialized, kindNames, uint8Type, x, init, jsType, reflectType, isWrapped, copyStruct, makeValue, MakeSlice, jsObject, TypeOf, ValueOf, SliceOf, Zero, unsafe_New, makeInt, memmove, loadScalar, chanclose, chanrecv, chansend, mapaccess, mapassign, mapdelete, mapiterinit, mapiterkey, mapiternext, maplen, cvtDirect, methodReceiver, valueInterface, ifaceE2I, methodName, makeMethodValue, PtrTo, implements$1, directlyAssignable, haveIdenticalUnderlyingType, toType, overflowFloat32, New, convertOp, makeFloat, makeComplex, makeString, makeBytes, makeRunes, cvtInt, cvtUint, cvtFloatInt, cvtFloatUint, cvtIntFloat, cvtUintFloat, cvtFloat, cvtComplex, cvtIntString, cvtUintString, cvtBytesString, cvtStringBytes, cvtRunesString, cvtStringRunes, cvtT2I, cvtI2I, call;
	mapIter = $pkg.mapIter = $newType(0, "Struct", "reflect.mapIter", "mapIter", "reflect", function(t_, m_, keys_, i_) {
		this.$val = this;
		this.t = t_ !== undefined ? t_ : null;
		this.m = m_ !== undefined ? m_ : null;
		this.keys = keys_ !== undefined ? keys_ : null;
		this.i = i_ !== undefined ? i_ : 0;
	});
	Type = $pkg.Type = $newType(8, "Interface", "reflect.Type", "Type", "reflect", null);
	Kind = $pkg.Kind = $newType(4, "Uint", "reflect.Kind", "Kind", "reflect", null);
	rtype = $pkg.rtype = $newType(0, "Struct", "reflect.rtype", "rtype", "reflect", function(size_, hash_, _$2_, align_, fieldAlign_, kind_, alg_, gc_, string_, uncommonType_, ptrToThis_, zero_) {
		this.$val = this;
		this.size = size_ !== undefined ? size_ : 0;
		this.hash = hash_ !== undefined ? hash_ : 0;
		this._$2 = _$2_ !== undefined ? _$2_ : 0;
		this.align = align_ !== undefined ? align_ : 0;
		this.fieldAlign = fieldAlign_ !== undefined ? fieldAlign_ : 0;
		this.kind = kind_ !== undefined ? kind_ : 0;
		this.alg = alg_ !== undefined ? alg_ : ($ptrType($Uintptr)).nil;
		this.gc = gc_ !== undefined ? gc_ : 0;
		this.string = string_ !== undefined ? string_ : ($ptrType($String)).nil;
		this.uncommonType = uncommonType_ !== undefined ? uncommonType_ : ($ptrType(uncommonType)).nil;
		this.ptrToThis = ptrToThis_ !== undefined ? ptrToThis_ : ($ptrType(rtype)).nil;
		this.zero = zero_ !== undefined ? zero_ : 0;
	});
	method = $pkg.method = $newType(0, "Struct", "reflect.method", "method", "reflect", function(name_, pkgPath_, mtyp_, typ_, ifn_, tfn_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : ($ptrType($String)).nil;
		this.pkgPath = pkgPath_ !== undefined ? pkgPath_ : ($ptrType($String)).nil;
		this.mtyp = mtyp_ !== undefined ? mtyp_ : ($ptrType(rtype)).nil;
		this.typ = typ_ !== undefined ? typ_ : ($ptrType(rtype)).nil;
		this.ifn = ifn_ !== undefined ? ifn_ : 0;
		this.tfn = tfn_ !== undefined ? tfn_ : 0;
	});
	uncommonType = $pkg.uncommonType = $newType(0, "Struct", "reflect.uncommonType", "uncommonType", "reflect", function(name_, pkgPath_, methods_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : ($ptrType($String)).nil;
		this.pkgPath = pkgPath_ !== undefined ? pkgPath_ : ($ptrType($String)).nil;
		this.methods = methods_ !== undefined ? methods_ : ($sliceType(method)).nil;
	});
	ChanDir = $pkg.ChanDir = $newType(4, "Int", "reflect.ChanDir", "ChanDir", "reflect", null);
	arrayType = $pkg.arrayType = $newType(0, "Struct", "reflect.arrayType", "arrayType", "reflect", function(rtype_, elem_, slice_, len_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.elem = elem_ !== undefined ? elem_ : ($ptrType(rtype)).nil;
		this.slice = slice_ !== undefined ? slice_ : ($ptrType(rtype)).nil;
		this.len = len_ !== undefined ? len_ : 0;
	});
	chanType = $pkg.chanType = $newType(0, "Struct", "reflect.chanType", "chanType", "reflect", function(rtype_, elem_, dir_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.elem = elem_ !== undefined ? elem_ : ($ptrType(rtype)).nil;
		this.dir = dir_ !== undefined ? dir_ : 0;
	});
	funcType = $pkg.funcType = $newType(0, "Struct", "reflect.funcType", "funcType", "reflect", function(rtype_, dotdotdot_, in$2_, out_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.dotdotdot = dotdotdot_ !== undefined ? dotdotdot_ : false;
		this.in$2 = in$2_ !== undefined ? in$2_ : ($sliceType(($ptrType(rtype)))).nil;
		this.out = out_ !== undefined ? out_ : ($sliceType(($ptrType(rtype)))).nil;
	});
	imethod = $pkg.imethod = $newType(0, "Struct", "reflect.imethod", "imethod", "reflect", function(name_, pkgPath_, typ_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : ($ptrType($String)).nil;
		this.pkgPath = pkgPath_ !== undefined ? pkgPath_ : ($ptrType($String)).nil;
		this.typ = typ_ !== undefined ? typ_ : ($ptrType(rtype)).nil;
	});
	interfaceType = $pkg.interfaceType = $newType(0, "Struct", "reflect.interfaceType", "interfaceType", "reflect", function(rtype_, methods_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.methods = methods_ !== undefined ? methods_ : ($sliceType(imethod)).nil;
	});
	mapType = $pkg.mapType = $newType(0, "Struct", "reflect.mapType", "mapType", "reflect", function(rtype_, key_, elem_, bucket_, hmap_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.key = key_ !== undefined ? key_ : ($ptrType(rtype)).nil;
		this.elem = elem_ !== undefined ? elem_ : ($ptrType(rtype)).nil;
		this.bucket = bucket_ !== undefined ? bucket_ : ($ptrType(rtype)).nil;
		this.hmap = hmap_ !== undefined ? hmap_ : ($ptrType(rtype)).nil;
	});
	ptrType = $pkg.ptrType = $newType(0, "Struct", "reflect.ptrType", "ptrType", "reflect", function(rtype_, elem_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.elem = elem_ !== undefined ? elem_ : ($ptrType(rtype)).nil;
	});
	sliceType = $pkg.sliceType = $newType(0, "Struct", "reflect.sliceType", "sliceType", "reflect", function(rtype_, elem_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.elem = elem_ !== undefined ? elem_ : ($ptrType(rtype)).nil;
	});
	structField = $pkg.structField = $newType(0, "Struct", "reflect.structField", "structField", "reflect", function(name_, pkgPath_, typ_, tag_, offset_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : ($ptrType($String)).nil;
		this.pkgPath = pkgPath_ !== undefined ? pkgPath_ : ($ptrType($String)).nil;
		this.typ = typ_ !== undefined ? typ_ : ($ptrType(rtype)).nil;
		this.tag = tag_ !== undefined ? tag_ : ($ptrType($String)).nil;
		this.offset = offset_ !== undefined ? offset_ : 0;
	});
	structType = $pkg.structType = $newType(0, "Struct", "reflect.structType", "structType", "reflect", function(rtype_, fields_) {
		this.$val = this;
		this.rtype = rtype_ !== undefined ? rtype_ : new rtype.Ptr();
		this.fields = fields_ !== undefined ? fields_ : ($sliceType(structField)).nil;
	});
	Method = $pkg.Method = $newType(0, "Struct", "reflect.Method", "Method", "reflect", function(Name_, PkgPath_, Type_, Func_, Index_) {
		this.$val = this;
		this.Name = Name_ !== undefined ? Name_ : "";
		this.PkgPath = PkgPath_ !== undefined ? PkgPath_ : "";
		this.Type = Type_ !== undefined ? Type_ : null;
		this.Func = Func_ !== undefined ? Func_ : new Value.Ptr();
		this.Index = Index_ !== undefined ? Index_ : 0;
	});
	StructField = $pkg.StructField = $newType(0, "Struct", "reflect.StructField", "StructField", "reflect", function(Name_, PkgPath_, Type_, Tag_, Offset_, Index_, Anonymous_) {
		this.$val = this;
		this.Name = Name_ !== undefined ? Name_ : "";
		this.PkgPath = PkgPath_ !== undefined ? PkgPath_ : "";
		this.Type = Type_ !== undefined ? Type_ : null;
		this.Tag = Tag_ !== undefined ? Tag_ : "";
		this.Offset = Offset_ !== undefined ? Offset_ : 0;
		this.Index = Index_ !== undefined ? Index_ : ($sliceType($Int)).nil;
		this.Anonymous = Anonymous_ !== undefined ? Anonymous_ : false;
	});
	StructTag = $pkg.StructTag = $newType(8, "String", "reflect.StructTag", "StructTag", "reflect", null);
	fieldScan = $pkg.fieldScan = $newType(0, "Struct", "reflect.fieldScan", "fieldScan", "reflect", function(typ_, index_) {
		this.$val = this;
		this.typ = typ_ !== undefined ? typ_ : ($ptrType(structType)).nil;
		this.index = index_ !== undefined ? index_ : ($sliceType($Int)).nil;
	});
	Value = $pkg.Value = $newType(0, "Struct", "reflect.Value", "Value", "reflect", function(typ_, ptr_, scalar_, flag_) {
		this.$val = this;
		this.typ = typ_ !== undefined ? typ_ : ($ptrType(rtype)).nil;
		this.ptr = ptr_ !== undefined ? ptr_ : 0;
		this.scalar = scalar_ !== undefined ? scalar_ : 0;
		this.flag = flag_ !== undefined ? flag_ : 0;
	});
	flag = $pkg.flag = $newType(4, "Uintptr", "reflect.flag", "flag", "reflect", null);
	ValueError = $pkg.ValueError = $newType(0, "Struct", "reflect.ValueError", "ValueError", "reflect", function(Method_, Kind_) {
		this.$val = this;
		this.Method = Method_ !== undefined ? Method_ : "";
		this.Kind = Kind_ !== undefined ? Kind_ : 0;
	});
	iword = $pkg.iword = $newType(4, "UnsafePointer", "reflect.iword", "iword", "reflect", null);
	nonEmptyInterface = $pkg.nonEmptyInterface = $newType(0, "Struct", "reflect.nonEmptyInterface", "nonEmptyInterface", "reflect", function(itab_, word_) {
		this.$val = this;
		this.itab = itab_ !== undefined ? itab_ : ($ptrType(($structType([["ityp", "ityp", "reflect", ($ptrType(rtype)), ""], ["typ", "typ", "reflect", ($ptrType(rtype)), ""], ["link", "link", "reflect", $UnsafePointer, ""], ["bad", "bad", "reflect", $Int32, ""], ["unused", "unused", "reflect", $Int32, ""], ["fun", "fun", "reflect", ($arrayType($UnsafePointer, 100000)), ""]])))).nil;
		this.word = word_ !== undefined ? word_ : 0;
	});
	init = function() {
		var used, x$1, x$2, x$3, x$4, x$5, x$6, x$7, x$8, x$9, x$10, x$11, x$12, x$13, pkg, _map, _key, x$14;
		used = (function(i) {
		});
		used((x$1 = new rtype.Ptr(0, 0, 0, 0, 0, 0, ($ptrType($Uintptr)).nil, 0, ($ptrType($String)).nil, ($ptrType(uncommonType)).nil, ($ptrType(rtype)).nil, 0), new x$1.constructor.Struct(x$1)));
		used((x$2 = new uncommonType.Ptr(($ptrType($String)).nil, ($ptrType($String)).nil, ($sliceType(method)).nil), new x$2.constructor.Struct(x$2)));
		used((x$3 = new method.Ptr(($ptrType($String)).nil, ($ptrType($String)).nil, ($ptrType(rtype)).nil, ($ptrType(rtype)).nil, 0, 0), new x$3.constructor.Struct(x$3)));
		used((x$4 = new arrayType.Ptr(new rtype.Ptr(), ($ptrType(rtype)).nil, ($ptrType(rtype)).nil, 0), new x$4.constructor.Struct(x$4)));
		used((x$5 = new chanType.Ptr(new rtype.Ptr(), ($ptrType(rtype)).nil, 0), new x$5.constructor.Struct(x$5)));
		used((x$6 = new funcType.Ptr(new rtype.Ptr(), false, ($sliceType(($ptrType(rtype)))).nil, ($sliceType(($ptrType(rtype)))).nil), new x$6.constructor.Struct(x$6)));
		used((x$7 = new interfaceType.Ptr(new rtype.Ptr(), ($sliceType(imethod)).nil), new x$7.constructor.Struct(x$7)));
		used((x$8 = new mapType.Ptr(new rtype.Ptr(), ($ptrType(rtype)).nil, ($ptrType(rtype)).nil, ($ptrType(rtype)).nil, ($ptrType(rtype)).nil), new x$8.constructor.Struct(x$8)));
		used((x$9 = new ptrType.Ptr(new rtype.Ptr(), ($ptrType(rtype)).nil), new x$9.constructor.Struct(x$9)));
		used((x$10 = new sliceType.Ptr(new rtype.Ptr(), ($ptrType(rtype)).nil), new x$10.constructor.Struct(x$10)));
		used((x$11 = new structType.Ptr(new rtype.Ptr(), ($sliceType(structField)).nil), new x$11.constructor.Struct(x$11)));
		used((x$12 = new imethod.Ptr(($ptrType($String)).nil, ($ptrType($String)).nil, ($ptrType(rtype)).nil), new x$12.constructor.Struct(x$12)));
		used((x$13 = new structField.Ptr(($ptrType($String)).nil, ($ptrType($String)).nil, ($ptrType(rtype)).nil, ($ptrType($String)).nil, 0), new x$13.constructor.Struct(x$13)));
		pkg = $pkg;
		pkg.kinds = $externalize((_map = new $Map(), _key = "Bool", _map[_key] = { k: _key, v: 1 }, _key = "Int", _map[_key] = { k: _key, v: 2 }, _key = "Int8", _map[_key] = { k: _key, v: 3 }, _key = "Int16", _map[_key] = { k: _key, v: 4 }, _key = "Int32", _map[_key] = { k: _key, v: 5 }, _key = "Int64", _map[_key] = { k: _key, v: 6 }, _key = "Uint", _map[_key] = { k: _key, v: 7 }, _key = "Uint8", _map[_key] = { k: _key, v: 8 }, _key = "Uint16", _map[_key] = { k: _key, v: 9 }, _key = "Uint32", _map[_key] = { k: _key, v: 10 }, _key = "Uint64", _map[_key] = { k: _key, v: 11 }, _key = "Uintptr", _map[_key] = { k: _key, v: 12 }, _key = "Float32", _map[_key] = { k: _key, v: 13 }, _key = "Float64", _map[_key] = { k: _key, v: 14 }, _key = "Complex64", _map[_key] = { k: _key, v: 15 }, _key = "Complex128", _map[_key] = { k: _key, v: 16 }, _key = "Array", _map[_key] = { k: _key, v: 17 }, _key = "Chan", _map[_key] = { k: _key, v: 18 }, _key = "Func", _map[_key] = { k: _key, v: 19 }, _key = "Interface", _map[_key] = { k: _key, v: 20 }, _key = "Map", _map[_key] = { k: _key, v: 21 }, _key = "Ptr", _map[_key] = { k: _key, v: 22 }, _key = "Slice", _map[_key] = { k: _key, v: 23 }, _key = "String", _map[_key] = { k: _key, v: 24 }, _key = "Struct", _map[_key] = { k: _key, v: 25 }, _key = "UnsafePointer", _map[_key] = { k: _key, v: 26 }, _map), ($mapType($String, Kind)));
		pkg.RecvDir = 1;
		pkg.SendDir = 2;
		pkg.BothDir = 3;
		$reflect = pkg;
		initialized = true;
		uint8Type = (x$14 = TypeOf(new $Uint8(0)), (x$14 !== null && x$14.constructor === ($ptrType(rtype)) ? x$14.$val : $typeAssertionFailed(x$14, ($ptrType(rtype)))));
	};
	jsType = function(typ) {
		return typ.jsType;
	};
	reflectType = function(typ) {
		return typ.reflectType();
	};
	isWrapped = function(typ) {
		var _ref;
		_ref = typ.Kind();
		if (_ref === 1 || _ref === 2 || _ref === 3 || _ref === 4 || _ref === 5 || _ref === 7 || _ref === 8 || _ref === 9 || _ref === 10 || _ref === 12 || _ref === 13 || _ref === 14 || _ref === 17 || _ref === 21 || _ref === 19 || _ref === 24 || _ref === 25) {
			return true;
		} else if (_ref === 22) {
			return typ.Elem().Kind() === 17;
		}
		return false;
	};
	copyStruct = function(dst, src, typ) {
		var fields, i, name;
		fields = jsType(typ).fields;
		i = 0;
		while (i < $parseInt(fields.length)) {
			name = $internalize(fields[i][0], $String);
			dst[$externalize(name, $String)] = src[$externalize(name, $String)];
			i = i + (1) >> 0;
		}
	};
	makeValue = function(t, v, fl) {
		var rt;
		rt = t.common();
		if ((t.Kind() === 17) || (t.Kind() === 25) || rt.pointers()) {
			return new Value.Ptr(rt, v, 0, (fl | ((t.Kind() >>> 0) << 4 >>> 0)) >>> 0);
		}
		if (t.Size() > 4 || (t.Kind() === 24)) {
			return new Value.Ptr(rt, $newDataPointer(v, jsType(rt.ptrTo())), 0, (((fl | ((t.Kind() >>> 0) << 4 >>> 0)) >>> 0) | 2) >>> 0);
		}
		return new Value.Ptr(rt, 0, v, (fl | ((t.Kind() >>> 0) << 4 >>> 0)) >>> 0);
	};
	MakeSlice = $pkg.MakeSlice = function(typ, len, cap) {
		if (!((typ.Kind() === 23))) {
			$panic(new $String("reflect.MakeSlice of non-slice type"));
		}
		if (len < 0) {
			$panic(new $String("reflect.MakeSlice: negative len"));
		}
		if (cap < 0) {
			$panic(new $String("reflect.MakeSlice: negative cap"));
		}
		if (len > cap) {
			$panic(new $String("reflect.MakeSlice: len > cap"));
		}
		return makeValue(typ, jsType(typ).make(len, cap, $externalize((function() {
			return jsType(typ.Elem()).zero();
		}), ($funcType([], [js.Object], false)))), 0);
	};
	jsObject = function() {
		return reflectType($packages[$externalize("github.com/gopherjs/gopherjs/js", $String)].Object);
	};
	TypeOf = $pkg.TypeOf = function(i) {
		var c;
		if (!initialized) {
			return new rtype.Ptr(0, 0, 0, 0, 0, 0, ($ptrType($Uintptr)).nil, 0, ($ptrType($String)).nil, ($ptrType(uncommonType)).nil, ($ptrType(rtype)).nil, 0);
		}
		if ($interfaceIsEqual(i, null)) {
			return null;
		}
		c = i.constructor;
		if (c.kind === undefined) {
			return jsObject();
		}
		return reflectType(c);
	};
	ValueOf = $pkg.ValueOf = function(i) {
		var c;
		if ($interfaceIsEqual(i, null)) {
			return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
		}
		c = i.constructor;
		if (c.kind === undefined) {
			return new Value.Ptr(jsObject(), 0, i, 320);
		}
		return makeValue(reflectType(c), i.$val, 0);
	};
	rtype.Ptr.prototype.ptrTo = function() {
		var t;
		t = this;
		return reflectType($ptrType(jsType(t)));
	};
	rtype.prototype.ptrTo = function() { return this.$val.ptrTo(); };
	SliceOf = $pkg.SliceOf = function(t) {
		return reflectType($sliceType(jsType(t)));
	};
	Zero = $pkg.Zero = function(typ) {
		return makeValue(typ, jsType(typ).zero(), 0);
	};
	unsafe_New = function(typ) {
		var _ref;
		_ref = typ.Kind();
		if (_ref === 25) {
			return new (jsType(typ).Ptr)();
		} else if (_ref === 17) {
			return jsType(typ).zero();
		} else {
			return $newDataPointer(jsType(typ).zero(), jsType(typ.ptrTo()));
		}
	};
	makeInt = function(f, bits, t) {
		var typ, ptr, s, _ref;
		typ = t.common();
		if (typ.size > 4) {
			ptr = unsafe_New(typ);
			ptr.$set(bits);
			return new Value.Ptr(typ, ptr, 0, (((f | 2) >>> 0) | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
		}
		s = 0;
		_ref = typ.Kind();
		if (_ref === 3) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set((bits.$low << 24 >> 24));
		} else if (_ref === 4) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set((bits.$low << 16 >> 16));
		} else if (_ref === 2 || _ref === 5) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set((bits.$low >> 0));
		} else if (_ref === 8) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set((bits.$low << 24 >>> 24));
		} else if (_ref === 9) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set((bits.$low << 16 >>> 16));
		} else if (_ref === 7 || _ref === 10 || _ref === 12) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set((bits.$low >>> 0));
		}
		return new Value.Ptr(typ, 0, s, (f | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
	};
	memmove = function(adst, asrc, n) {
		adst.$set(asrc.$get());
	};
	loadScalar = function(p, n) {
		return p.$get();
	};
	chanclose = function(ch) {
		$panic(new runtime.NotSupportedError.Ptr("channels"));
	};
	chanrecv = function(t, ch, nb, val) {
		var selected = false, received = false;
		$panic(new runtime.NotSupportedError.Ptr("channels"));
	};
	chansend = function(t, ch, val, nb) {
		$panic(new runtime.NotSupportedError.Ptr("channels"));
	};
	mapaccess = function(t, m, key) {
		var k, entry;
		k = key.$get();
		if (!(k.$key === undefined)) {
			k = k.$key();
		}
		entry = m[$externalize($internalize(k, $String), $String)];
		if (entry === undefined) {
			return 0;
		}
		return $newDataPointer(entry.v, jsType(PtrTo(t.Elem())));
	};
	mapassign = function(t, m, key, val) {
		var kv, k, jsVal, et, newVal, entry;
		kv = key.$get();
		k = kv;
		if (!(k.$key === undefined)) {
			k = k.$key();
		}
		jsVal = val.$get();
		et = t.Elem();
		if (et.Kind() === 25) {
			newVal = jsType(et).zero();
			copyStruct(newVal, jsVal, et);
			jsVal = newVal;
		}
		entry = new ($global.Object)();
		entry.k = kv;
		entry.v = jsVal;
		m[$externalize($internalize(k, $String), $String)] = entry;
	};
	mapdelete = function(t, m, key) {
		var k;
		k = key.$get();
		if (!(k.$key === undefined)) {
			k = k.$key();
		}
		delete m[$externalize($internalize(k, $String), $String)];
	};
	mapiterinit = function(t, m) {
		return new mapIter.Ptr(t, m, $keys(m), 0);
	};
	mapiterkey = function(it) {
		var iter, k;
		iter = it;
		k = iter.keys[iter.i];
		return $newDataPointer(iter.m[$externalize($internalize(k, $String), $String)].k, jsType(PtrTo(iter.t.Key())));
	};
	mapiternext = function(it) {
		var iter;
		iter = it;
		iter.i = iter.i + (1) >> 0;
	};
	maplen = function(m) {
		return $parseInt($keys(m).length);
	};
	cvtDirect = function(v, typ) {
		var srcVal, val, k, _ref, slice;
		srcVal = v.iword();
		if (srcVal === jsType(v.typ).nil) {
			return makeValue(typ, jsType(typ).nil, v.flag);
		}
		val = null;
		k = typ.Kind();
		_ref = k;
		switch (0) { default: if (_ref === 18) {
			val = new (jsType(typ))();
		} else if (_ref === 23) {
			slice = new (jsType(typ))(srcVal.$array);
			slice.$offset = srcVal.$offset;
			slice.$length = srcVal.$length;
			slice.$capacity = srcVal.$capacity;
			val = $newDataPointer(slice, jsType(PtrTo(typ)));
		} else if (_ref === 22) {
			if (typ.Elem().Kind() === 25) {
				if ($interfaceIsEqual(typ.Elem(), v.typ.Elem())) {
					val = srcVal;
					break;
				}
				val = new (jsType(typ))();
				copyStruct(val, srcVal, typ.Elem());
				break;
			}
			val = new (jsType(typ))(srcVal.$get, srcVal.$set);
		} else if (_ref === 25) {
			val = new (jsType(typ).Ptr)();
			copyStruct(val, srcVal, typ);
		} else if (_ref === 17 || _ref === 19 || _ref === 20 || _ref === 21 || _ref === 24) {
			val = v.ptr;
		} else {
			$panic(new ValueError.Ptr("reflect.Convert", k));
		} }
		return new Value.Ptr(typ.common(), val, 0, (((v.flag & 3) >>> 0) | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
	};
	methodReceiver = function(op, v, i) {
		var rcvrtype = ($ptrType(rtype)).nil, t = ($ptrType(rtype)).nil, fn = 0, name, tt, x$1, m, iface, ut, x$2, m$1, rcvr;
		name = "";
		if (v.typ.Kind() === 20) {
			tt = v.typ.interfaceType;
			if (i < 0 || i >= tt.methods.$length) {
				$panic(new $String("reflect: internal error: invalid method index"));
			}
			m = (x$1 = tt.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
			if (!($pointerIsEqual(m.pkgPath, ($ptrType($String)).nil))) {
				$panic(new $String("reflect: " + op + " of unexported method"));
			}
			iface = $clone(v.ptr, nonEmptyInterface);
			if (iface.itab === ($ptrType(($structType([["ityp", "ityp", "reflect", ($ptrType(rtype)), ""], ["typ", "typ", "reflect", ($ptrType(rtype)), ""], ["link", "link", "reflect", $UnsafePointer, ""], ["bad", "bad", "reflect", $Int32, ""], ["unused", "unused", "reflect", $Int32, ""], ["fun", "fun", "reflect", ($arrayType($UnsafePointer, 100000)), ""]])))).nil) {
				$panic(new $String("reflect: " + op + " of method on nil interface value"));
			}
			t = m.typ;
			name = m.name.$get();
		} else {
			ut = v.typ.uncommonType.uncommon();
			if (ut === ($ptrType(uncommonType)).nil || i < 0 || i >= ut.methods.$length) {
				$panic(new $String("reflect: internal error: invalid method index"));
			}
			m$1 = (x$2 = ut.methods, ((i < 0 || i >= x$2.$length) ? $throwRuntimeError("index out of range") : x$2.$array[x$2.$offset + i]));
			if (!($pointerIsEqual(m$1.pkgPath, ($ptrType($String)).nil))) {
				$panic(new $String("reflect: " + op + " of unexported method"));
			}
			t = m$1.mtyp;
			name = $internalize(jsType(v.typ).methods[i][0], $String);
		}
		rcvr = v.iword();
		if (isWrapped(v.typ)) {
			rcvr = new (jsType(v.typ))(rcvr);
		}
		fn = rcvr[$externalize(name, $String)];
		return [rcvrtype, t, fn];
	};
	valueInterface = function(v, safe) {
		if (v.flag === 0) {
			$panic(new ValueError.Ptr("reflect.Value.Interface", 0));
		}
		if (safe && !((((v.flag & 1) >>> 0) === 0))) {
			$panic(new $String("reflect.Value.Interface: cannot return value obtained from unexported field or method"));
		}
		if (!((((v.flag & 8) >>> 0) === 0))) {
			$copy(v, makeMethodValue("Interface", $clone(v, Value)), Value);
		}
		if (isWrapped(v.typ)) {
			return new (jsType(v.typ))(v.iword());
		}
		return v.iword();
	};
	ifaceE2I = function(t, src, dst) {
		dst.$set(src);
	};
	methodName = function() {
		return "?FIXME?";
	};
	makeMethodValue = function(op, v) {
		var _tuple, fn, rcvr, fv;
		if (((v.flag & 8) >>> 0) === 0) {
			$panic(new $String("reflect: internal error: invalid use of makePartialFunc"));
		}
		_tuple = methodReceiver(op, $clone(v, Value), (v.flag >> 0) >> 9 >> 0); fn = _tuple[2];
		rcvr = v.iword();
		if (isWrapped(v.typ)) {
			rcvr = new (jsType(v.typ))(rcvr);
		}
		fv = (function() {
			return fn.apply(rcvr, $externalize(new ($sliceType(js.Object))($global.Array.prototype.slice.call(arguments, [])), ($sliceType(js.Object))));
		});
		return new Value.Ptr(v.Type().common(), fv, 0, (((v.flag & 1) >>> 0) | 304) >>> 0);
	};
	rtype.Ptr.prototype.pointers = function() {
		var t, _ref;
		t = this;
		_ref = t.Kind();
		if (_ref === 22 || _ref === 21 || _ref === 18 || _ref === 19 || _ref === 25 || _ref === 17) {
			return true;
		} else {
			return false;
		}
	};
	rtype.prototype.pointers = function() { return this.$val.pointers(); };
	uncommonType.Ptr.prototype.Method = function(i) {
		var m = new Method.Ptr(), t, x$1, p, fl, mt, name, fn;
		t = this;
		if (t === ($ptrType(uncommonType)).nil || i < 0 || i >= t.methods.$length) {
			$panic(new $String("reflect: Method index out of range"));
		}
		p = (x$1 = t.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
		if (!($pointerIsEqual(p.name, ($ptrType($String)).nil))) {
			m.Name = p.name.$get();
		}
		fl = 304;
		if (!($pointerIsEqual(p.pkgPath, ($ptrType($String)).nil))) {
			m.PkgPath = p.pkgPath.$get();
			fl = (fl | (1)) >>> 0;
		}
		mt = p.typ;
		m.Type = mt;
		name = $internalize(t.jsType.methods[i][0], $String);
		fn = (function(rcvr) {
			return rcvr[$externalize(name, $String)].apply(rcvr, $externalize($subslice(new ($sliceType(js.Object))($global.Array.prototype.slice.call(arguments, [])), 1), ($sliceType(js.Object))));
		});
		$copy(m.Func, new Value.Ptr(mt, fn, 0, fl), Value);
		m.Index = i;
		return m;
	};
	uncommonType.prototype.Method = function(i) { return this.$val.Method(i); };
	Value.Ptr.prototype.iword = function() {
		var v, val, _ref, newVal;
		v = new Value.Ptr(); $copy(v, this, Value);
		if ((v.typ.Kind() === 17) || (v.typ.Kind() === 25)) {
			return v.ptr;
		}
		if (!((((v.flag & 2) >>> 0) === 0))) {
			val = v.ptr.$get();
			if (!(val === null) && !(val.constructor === jsType(v.typ))) {
				_ref = v.typ.Kind();
				switch (0) { default: if (_ref === 11 || _ref === 6) {
					val = new (jsType(v.typ))(val.$high, val.$low);
				} else if (_ref === 15 || _ref === 16) {
					val = new (jsType(v.typ))(val.$real, val.$imag);
				} else if (_ref === 23) {
					if (val === val.constructor.nil) {
						val = jsType(v.typ).nil;
						break;
					}
					newVal = new (jsType(v.typ))(val.$array);
					newVal.$offset = val.$offset;
					newVal.$length = val.$length;
					newVal.$capacity = val.$capacity;
					val = newVal;
				} }
			}
			return val;
		}
		if (v.typ.pointers()) {
			return v.ptr;
		}
		return v.scalar;
	};
	Value.prototype.iword = function() { return this.$val.iword(); };
	Value.Ptr.prototype.call = function(op, in$1) {
		var v, t, fn, rcvr, _tuple, isSlice, n, _ref, _i, x$1, i, _tmp, _tmp$1, xt, targ, m, slice, elem, i$1, x$2, x$3, xt$1, origIn, nin, nout, argsArray, _ref$1, _i$1, i$2, arg, results, _ref$2, ret, _ref$3, _i$2, i$3;
		v = new Value.Ptr(); $copy(v, this, Value);
		t = v.typ;
		fn = 0;
		rcvr = null;
		if (!((((v.flag & 8) >>> 0) === 0))) {
			_tuple = methodReceiver(op, $clone(v, Value), (v.flag >> 0) >> 9 >> 0); t = _tuple[1]; fn = _tuple[2];
			rcvr = v.iword();
			if (isWrapped(v.typ)) {
				rcvr = new (jsType(v.typ))(rcvr);
			}
		} else {
			fn = v.iword();
		}
		if (fn === 0) {
			$panic(new $String("reflect.Value.Call: call of nil function"));
		}
		isSlice = op === "CallSlice";
		n = t.NumIn();
		if (isSlice) {
			if (!t.IsVariadic()) {
				$panic(new $String("reflect: CallSlice of non-variadic function"));
			}
			if (in$1.$length < n) {
				$panic(new $String("reflect: CallSlice with too few input arguments"));
			}
			if (in$1.$length > n) {
				$panic(new $String("reflect: CallSlice with too many input arguments"));
			}
		} else {
			if (t.IsVariadic()) {
				n = n - (1) >> 0;
			}
			if (in$1.$length < n) {
				$panic(new $String("reflect: Call with too few input arguments"));
			}
			if (!t.IsVariadic() && in$1.$length > n) {
				$panic(new $String("reflect: Call with too many input arguments"));
			}
		}
		_ref = in$1;
		_i = 0;
		while (_i < _ref.$length) {
			x$1 = new Value.Ptr(); $copy(x$1, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), Value);
			if (x$1.Kind() === 0) {
				$panic(new $String("reflect: " + op + " using zero Value argument"));
			}
			_i++;
		}
		i = 0;
		while (i < n) {
			_tmp = ((i < 0 || i >= in$1.$length) ? $throwRuntimeError("index out of range") : in$1.$array[in$1.$offset + i]).Type(); _tmp$1 = t.In(i); xt = _tmp; targ = _tmp$1;
			if (!xt.AssignableTo(targ)) {
				$panic(new $String("reflect: " + op + " using " + xt.String() + " as type " + targ.String()));
			}
			i = i + (1) >> 0;
		}
		if (!isSlice && t.IsVariadic()) {
			m = in$1.$length - n >> 0;
			slice = new Value.Ptr(); $copy(slice, MakeSlice(t.In(n), m, m), Value);
			elem = t.In(n).Elem();
			i$1 = 0;
			while (i$1 < m) {
				x$3 = new Value.Ptr(); $copy(x$3, (x$2 = n + i$1 >> 0, ((x$2 < 0 || x$2 >= in$1.$length) ? $throwRuntimeError("index out of range") : in$1.$array[in$1.$offset + x$2])), Value);
				xt$1 = x$3.Type();
				if (!xt$1.AssignableTo(elem)) {
					$panic(new $String("reflect: cannot use " + xt$1.String() + " as type " + elem.String() + " in " + op));
				}
				slice.Index(i$1).Set($clone(x$3, Value));
				i$1 = i$1 + (1) >> 0;
			}
			origIn = in$1;
			in$1 = ($sliceType(Value)).make((n + 1 >> 0));
			$copySlice($subslice(in$1, 0, n), origIn);
			$copy(((n < 0 || n >= in$1.$length) ? $throwRuntimeError("index out of range") : in$1.$array[in$1.$offset + n]), slice, Value);
		}
		nin = in$1.$length;
		if (!((nin === t.NumIn()))) {
			$panic(new $String("reflect.Value.Call: wrong argument count"));
		}
		nout = t.NumOut();
		argsArray = new ($global.Array)(t.NumIn());
		_ref$1 = in$1;
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			i$2 = _i$1;
			arg = new Value.Ptr(); $copy(arg, ((_i$1 < 0 || _i$1 >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i$1]), Value);
			argsArray[i$2] = arg.assignTo("reflect.Value.Call", t.In(i$2).common(), ($ptrType($emptyInterface)).nil).iword();
			_i$1++;
		}
		results = fn.apply(rcvr, argsArray);
		_ref$2 = nout;
		if (_ref$2 === 0) {
			return ($sliceType(Value)).nil;
		} else if (_ref$2 === 1) {
			return new ($sliceType(Value))([$clone(makeValue(t.Out(0), results, 0), Value)]);
		} else {
			ret = ($sliceType(Value)).make(nout);
			_ref$3 = ret;
			_i$2 = 0;
			while (_i$2 < _ref$3.$length) {
				i$3 = _i$2;
				$copy(((i$3 < 0 || i$3 >= ret.$length) ? $throwRuntimeError("index out of range") : ret.$array[ret.$offset + i$3]), makeValue(t.Out(i$3), results[i$3], 0), Value);
				_i$2++;
			}
			return ret;
		}
	};
	Value.prototype.call = function(op, in$1) { return this.$val.call(op, in$1); };
	Value.Ptr.prototype.Cap = function() {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 17) {
			return v.typ.Len();
		} else if (_ref === 23) {
			return $parseInt(v.iword().$capacity) >> 0;
		}
		$panic(new ValueError.Ptr("reflect.Value.Cap", k));
	};
	Value.prototype.Cap = function() { return this.$val.Cap(); };
	Value.Ptr.prototype.Elem = function() {
		var v, k, _ref, val, typ, val$1, tt, fl;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 20) {
			val = v.iword();
			if (val === null) {
				return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
			}
			typ = reflectType(val.constructor);
			return makeValue(typ, val.$val, (v.flag & 1) >>> 0);
		} else if (_ref === 22) {
			if (v.IsNil()) {
				return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
			}
			val$1 = v.iword();
			tt = v.typ.ptrType;
			fl = (((((v.flag & 1) >>> 0) | 2) >>> 0) | 4) >>> 0;
			fl = (fl | (((tt.elem.Kind() >>> 0) << 4 >>> 0))) >>> 0;
			return new Value.Ptr(tt.elem, val$1, 0, fl);
		} else {
			$panic(new ValueError.Ptr("reflect.Value.Elem", k));
		}
	};
	Value.prototype.Elem = function() { return this.$val.Elem(); };
	Value.Ptr.prototype.Field = function(i) {
		var v, tt, x$1, field, name, typ, fl, s;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(25);
		tt = v.typ.structType;
		if (i < 0 || i >= tt.fields.$length) {
			$panic(new $String("reflect: Field index out of range"));
		}
		field = (x$1 = tt.fields, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
		name = $internalize(jsType(v.typ).fields[i][0], $String);
		typ = field.typ;
		fl = (v.flag & 7) >>> 0;
		if (!($pointerIsEqual(field.pkgPath, ($ptrType($String)).nil))) {
			fl = (fl | (1)) >>> 0;
		}
		fl = (fl | (((typ.Kind() >>> 0) << 4 >>> 0))) >>> 0;
		s = v.ptr;
		if (!((((fl & 2) >>> 0) === 0)) && !((typ.Kind() === 17)) && !((typ.Kind() === 25))) {
			return new Value.Ptr(typ, new (jsType(PtrTo(typ)))($externalize((function() {
				return s[$externalize(name, $String)];
			}), ($funcType([], [js.Object], false))), $externalize((function(v$1) {
				s[$externalize(name, $String)] = v$1;
			}), ($funcType([js.Object], [], false)))), 0, fl);
		}
		return makeValue(typ, s[$externalize(name, $String)], fl);
	};
	Value.prototype.Field = function(i) { return this.$val.Field(i); };
	Value.Ptr.prototype.Index = function(i) {
		var v, k, _ref, tt, typ, fl, a, s, tt$1, typ$1, fl$1, a$1, str, fl$2;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 17) {
			tt = v.typ.arrayType;
			if (i < 0 || i > (tt.len >> 0)) {
				$panic(new $String("reflect: array index out of range"));
			}
			typ = tt.elem;
			fl = (v.flag & 7) >>> 0;
			fl = (fl | (((typ.Kind() >>> 0) << 4 >>> 0))) >>> 0;
			a = v.ptr;
			if (!((((fl & 2) >>> 0) === 0)) && !((typ.Kind() === 17)) && !((typ.Kind() === 25))) {
				return new Value.Ptr(typ, new (jsType(PtrTo(typ)))($externalize((function() {
					return a[i];
				}), ($funcType([], [js.Object], false))), $externalize((function(v$1) {
					a[i] = v$1;
				}), ($funcType([js.Object], [], false)))), 0, fl);
			}
			return makeValue(typ, a[i], fl);
		} else if (_ref === 23) {
			s = v.iword();
			if (i < 0 || i >= ($parseInt(s.$length) >> 0)) {
				$panic(new $String("reflect: slice index out of range"));
			}
			tt$1 = v.typ.sliceType;
			typ$1 = tt$1.elem;
			fl$1 = (6 | ((v.flag & 1) >>> 0)) >>> 0;
			fl$1 = (fl$1 | (((typ$1.Kind() >>> 0) << 4 >>> 0))) >>> 0;
			i = i + (($parseInt(s.$offset) >> 0)) >> 0;
			a$1 = s.$array;
			if (!((((fl$1 & 2) >>> 0) === 0)) && !((typ$1.Kind() === 17)) && !((typ$1.Kind() === 25))) {
				return new Value.Ptr(typ$1, new (jsType(PtrTo(typ$1)))($externalize((function() {
					return a$1[i];
				}), ($funcType([], [js.Object], false))), $externalize((function(v$1) {
					a$1[i] = v$1;
				}), ($funcType([js.Object], [], false)))), 0, fl$1);
			}
			return makeValue(typ$1, a$1[i], fl$1);
		} else if (_ref === 24) {
			str = v.ptr.$get();
			if (i < 0 || i >= str.length) {
				$panic(new $String("reflect: string index out of range"));
			}
			fl$2 = (((v.flag & 1) >>> 0) | 128) >>> 0;
			return new Value.Ptr(uint8Type, 0, (str.charCodeAt(i) >>> 0), fl$2);
		} else {
			$panic(new ValueError.Ptr("reflect.Value.Index", k));
		}
	};
	Value.prototype.Index = function(i) { return this.$val.Index(i); };
	Value.Ptr.prototype.IsNil = function() {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 18 || _ref === 22 || _ref === 23) {
			return v.iword() === jsType(v.typ).nil;
		} else if (_ref === 19) {
			return v.iword() === $throwNilPointerError;
		} else if (_ref === 21) {
			return v.iword() === false;
		} else if (_ref === 20) {
			return v.iword() === null;
		} else {
			$panic(new ValueError.Ptr("reflect.Value.IsNil", k));
		}
	};
	Value.prototype.IsNil = function() { return this.$val.IsNil(); };
	Value.Ptr.prototype.Len = function() {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 17 || _ref === 24) {
			return $parseInt(v.iword().length);
		} else if (_ref === 23) {
			return $parseInt(v.iword().$length) >> 0;
		} else if (_ref === 21) {
			return $parseInt($keys(v.iword()).length);
		} else {
			$panic(new ValueError.Ptr("reflect.Value.Len", k));
		}
	};
	Value.prototype.Len = function() { return this.$val.Len(); };
	Value.Ptr.prototype.Pointer = function() {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 18 || _ref === 21 || _ref === 22 || _ref === 23 || _ref === 26) {
			if (v.IsNil()) {
				return 0;
			}
			return v.iword();
		} else if (_ref === 19) {
			if (v.IsNil()) {
				return 0;
			}
			return 1;
		} else {
			$panic(new ValueError.Ptr("reflect.Value.Pointer", k));
		}
	};
	Value.prototype.Pointer = function() { return this.$val.Pointer(); };
	Value.Ptr.prototype.Set = function(x$1) {
		var v, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(x$1.flag)).mustBeExported();
		$copy(x$1, x$1.assignTo("reflect.Set", v.typ, ($ptrType($emptyInterface)).nil), Value);
		if (!((((v.flag & 2) >>> 0) === 0))) {
			_ref = v.typ.Kind();
			if (_ref === 17) {
				$copy(v.ptr, x$1.ptr, jsType(v.typ));
			} else if (_ref === 20) {
				v.ptr.$set(valueInterface($clone(x$1, Value), false));
			} else if (_ref === 25) {
				copyStruct(v.ptr, x$1.ptr, v.typ);
			} else {
				v.ptr.$set(x$1.iword());
			}
			return;
		}
		v.ptr = x$1.ptr;
	};
	Value.prototype.Set = function(x$1) { return this.$val.Set(x$1); };
	Value.Ptr.prototype.SetCap = function(n) {
		var v, s, newSlice;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(23);
		s = v.ptr.$get();
		if (n < ($parseInt(s.$length) >> 0) || n > ($parseInt(s.$capacity) >> 0)) {
			$panic(new $String("reflect: slice capacity out of range in SetCap"));
		}
		newSlice = new (jsType(v.typ))(s.$array);
		newSlice.$offset = s.$offset;
		newSlice.$length = s.$length;
		newSlice.$capacity = n;
		v.ptr.$set(newSlice);
	};
	Value.prototype.SetCap = function(n) { return this.$val.SetCap(n); };
	Value.Ptr.prototype.SetLen = function(n) {
		var v, s, newSlice;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(23);
		s = v.ptr.$get();
		if (n < 0 || n > ($parseInt(s.$capacity) >> 0)) {
			$panic(new $String("reflect: slice length out of range in SetLen"));
		}
		newSlice = new (jsType(v.typ))(s.$array);
		newSlice.$offset = s.$offset;
		newSlice.$length = n;
		newSlice.$capacity = s.$capacity;
		v.ptr.$set(newSlice);
	};
	Value.prototype.SetLen = function(n) { return this.$val.SetLen(n); };
	Value.Ptr.prototype.Slice = function(i, j) {
		var v, cap, typ, s, kind, _ref, tt, str;
		v = new Value.Ptr(); $copy(v, this, Value);
		cap = 0;
		typ = null;
		s = null;
		kind = (new flag(v.flag)).kind();
		_ref = kind;
		if (_ref === 17) {
			if (((v.flag & 4) >>> 0) === 0) {
				$panic(new $String("reflect.Value.Slice: slice of unaddressable array"));
			}
			tt = v.typ.arrayType;
			cap = (tt.len >> 0);
			typ = SliceOf(tt.elem);
			s = new (jsType(typ))(v.iword());
		} else if (_ref === 23) {
			typ = v.typ;
			s = v.iword();
			cap = $parseInt(s.$capacity) >> 0;
		} else if (_ref === 24) {
			str = v.ptr.$get();
			if (i < 0 || j < i || j > str.length) {
				$panic(new $String("reflect.Value.Slice: string slice index out of bounds"));
			}
			return ValueOf(new $String(str.substring(i, j)));
		} else {
			$panic(new ValueError.Ptr("reflect.Value.Slice", kind));
		}
		if (i < 0 || j < i || j > cap) {
			$panic(new $String("reflect.Value.Slice: slice index out of bounds"));
		}
		return makeValue(typ, $subslice(s, i, j), (v.flag & 1) >>> 0);
	};
	Value.prototype.Slice = function(i, j) { return this.$val.Slice(i, j); };
	Value.Ptr.prototype.Slice3 = function(i, j, k) {
		var v, cap, typ, s, kind, _ref, tt;
		v = new Value.Ptr(); $copy(v, this, Value);
		cap = 0;
		typ = null;
		s = null;
		kind = (new flag(v.flag)).kind();
		_ref = kind;
		if (_ref === 17) {
			if (((v.flag & 4) >>> 0) === 0) {
				$panic(new $String("reflect.Value.Slice: slice of unaddressable array"));
			}
			tt = v.typ.arrayType;
			cap = (tt.len >> 0);
			typ = SliceOf(tt.elem);
			s = new (jsType(typ))(v.iword());
		} else if (_ref === 23) {
			typ = v.typ;
			s = v.iword();
			cap = $parseInt(s.$capacity) >> 0;
		} else {
			$panic(new ValueError.Ptr("reflect.Value.Slice3", kind));
		}
		if (i < 0 || j < i || k < j || k > cap) {
			$panic(new $String("reflect.Value.Slice3: slice index out of bounds"));
		}
		return makeValue(typ, $subslice(s, i, j, k), (v.flag & 1) >>> 0);
	};
	Value.prototype.Slice3 = function(i, j, k) { return this.$val.Slice3(i, j, k); };
	Kind.prototype.String = function() {
		var k;
		k = this.$val;
		if ((k >> 0) < kindNames.$length) {
			return ((k < 0 || k >= kindNames.$length) ? $throwRuntimeError("index out of range") : kindNames.$array[kindNames.$offset + k]);
		}
		return "kind" + strconv.Itoa((k >> 0));
	};
	$ptrType(Kind).prototype.String = function() { return new Kind(this.$get()).String(); };
	uncommonType.Ptr.prototype.uncommon = function() {
		var t;
		t = this;
		return t;
	};
	uncommonType.prototype.uncommon = function() { return this.$val.uncommon(); };
	uncommonType.Ptr.prototype.PkgPath = function() {
		var t;
		t = this;
		if (t === ($ptrType(uncommonType)).nil || $pointerIsEqual(t.pkgPath, ($ptrType($String)).nil)) {
			return "";
		}
		return t.pkgPath.$get();
	};
	uncommonType.prototype.PkgPath = function() { return this.$val.PkgPath(); };
	uncommonType.Ptr.prototype.Name = function() {
		var t;
		t = this;
		if (t === ($ptrType(uncommonType)).nil || $pointerIsEqual(t.name, ($ptrType($String)).nil)) {
			return "";
		}
		return t.name.$get();
	};
	uncommonType.prototype.Name = function() { return this.$val.Name(); };
	rtype.Ptr.prototype.String = function() {
		var t;
		t = this;
		return t.string.$get();
	};
	rtype.prototype.String = function() { return this.$val.String(); };
	rtype.Ptr.prototype.Size = function() {
		var t;
		t = this;
		return t.size;
	};
	rtype.prototype.Size = function() { return this.$val.Size(); };
	rtype.Ptr.prototype.Bits = function() {
		var t, k, x$1;
		t = this;
		if (t === ($ptrType(rtype)).nil) {
			$panic(new $String("reflect: Bits of nil Type"));
		}
		k = t.Kind();
		if (k < 2 || k > 16) {
			$panic(new $String("reflect: Bits of non-arithmetic Type " + t.String()));
		}
		return (x$1 = (t.size >> 0), (((x$1 >>> 16 << 16) * 8 >> 0) + (x$1 << 16 >>> 16) * 8) >> 0);
	};
	rtype.prototype.Bits = function() { return this.$val.Bits(); };
	rtype.Ptr.prototype.Align = function() {
		var t;
		t = this;
		return (t.align >> 0);
	};
	rtype.prototype.Align = function() { return this.$val.Align(); };
	rtype.Ptr.prototype.FieldAlign = function() {
		var t;
		t = this;
		return (t.fieldAlign >> 0);
	};
	rtype.prototype.FieldAlign = function() { return this.$val.FieldAlign(); };
	rtype.Ptr.prototype.Kind = function() {
		var t;
		t = this;
		return (((t.kind & 127) >>> 0) >>> 0);
	};
	rtype.prototype.Kind = function() { return this.$val.Kind(); };
	rtype.Ptr.prototype.common = function() {
		var t;
		t = this;
		return t;
	};
	rtype.prototype.common = function() { return this.$val.common(); };
	uncommonType.Ptr.prototype.NumMethod = function() {
		var t;
		t = this;
		if (t === ($ptrType(uncommonType)).nil) {
			return 0;
		}
		return t.methods.$length;
	};
	uncommonType.prototype.NumMethod = function() { return this.$val.NumMethod(); };
	uncommonType.Ptr.prototype.MethodByName = function(name) {
		var m = new Method.Ptr(), ok = false, t, p, _ref, _i, i, x$1, _tmp, _tmp$1;
		t = this;
		if (t === ($ptrType(uncommonType)).nil) {
			return [m, ok];
		}
		p = ($ptrType(method)).nil;
		_ref = t.methods;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			p = (x$1 = t.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
			if (!($pointerIsEqual(p.name, ($ptrType($String)).nil)) && p.name.$get() === name) {
				_tmp = new Method.Ptr(); $copy(_tmp, t.Method(i), Method); _tmp$1 = true; $copy(m, _tmp, Method); ok = _tmp$1;
				return [m, ok];
			}
			_i++;
		}
		return [m, ok];
	};
	uncommonType.prototype.MethodByName = function(name) { return this.$val.MethodByName(name); };
	rtype.Ptr.prototype.NumMethod = function() {
		var t, tt;
		t = this;
		if (t.Kind() === 20) {
			tt = t.interfaceType;
			return tt.NumMethod();
		}
		return t.uncommonType.NumMethod();
	};
	rtype.prototype.NumMethod = function() { return this.$val.NumMethod(); };
	rtype.Ptr.prototype.Method = function(i) {
		var m = new Method.Ptr(), t, tt;
		t = this;
		if (t.Kind() === 20) {
			tt = t.interfaceType;
			$copy(m, tt.Method(i), Method);
			return m;
		}
		$copy(m, t.uncommonType.Method(i), Method);
		return m;
	};
	rtype.prototype.Method = function(i) { return this.$val.Method(i); };
	rtype.Ptr.prototype.MethodByName = function(name) {
		var m = new Method.Ptr(), ok = false, t, tt, _tuple, _tuple$1;
		t = this;
		if (t.Kind() === 20) {
			tt = t.interfaceType;
			_tuple = tt.MethodByName(name); $copy(m, _tuple[0], Method); ok = _tuple[1];
			return [m, ok];
		}
		_tuple$1 = t.uncommonType.MethodByName(name); $copy(m, _tuple$1[0], Method); ok = _tuple$1[1];
		return [m, ok];
	};
	rtype.prototype.MethodByName = function(name) { return this.$val.MethodByName(name); };
	rtype.Ptr.prototype.PkgPath = function() {
		var t;
		t = this;
		return t.uncommonType.PkgPath();
	};
	rtype.prototype.PkgPath = function() { return this.$val.PkgPath(); };
	rtype.Ptr.prototype.Name = function() {
		var t;
		t = this;
		return t.uncommonType.Name();
	};
	rtype.prototype.Name = function() { return this.$val.Name(); };
	rtype.Ptr.prototype.ChanDir = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 18))) {
			$panic(new $String("reflect: ChanDir of non-chan type"));
		}
		tt = t.chanType;
		return (tt.dir >> 0);
	};
	rtype.prototype.ChanDir = function() { return this.$val.ChanDir(); };
	rtype.Ptr.prototype.IsVariadic = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 19))) {
			$panic(new $String("reflect: IsVariadic of non-func type"));
		}
		tt = t.funcType;
		return tt.dotdotdot;
	};
	rtype.prototype.IsVariadic = function() { return this.$val.IsVariadic(); };
	rtype.Ptr.prototype.Elem = function() {
		var t, _ref, tt, tt$1, tt$2, tt$3, tt$4;
		t = this;
		_ref = t.Kind();
		if (_ref === 17) {
			tt = t.arrayType;
			return toType(tt.elem);
		} else if (_ref === 18) {
			tt$1 = t.chanType;
			return toType(tt$1.elem);
		} else if (_ref === 21) {
			tt$2 = t.mapType;
			return toType(tt$2.elem);
		} else if (_ref === 22) {
			tt$3 = t.ptrType;
			return toType(tt$3.elem);
		} else if (_ref === 23) {
			tt$4 = t.sliceType;
			return toType(tt$4.elem);
		}
		$panic(new $String("reflect: Elem of invalid type"));
	};
	rtype.prototype.Elem = function() { return this.$val.Elem(); };
	rtype.Ptr.prototype.Field = function(i) {
		var t, tt;
		t = this;
		if (!((t.Kind() === 25))) {
			$panic(new $String("reflect: Field of non-struct type"));
		}
		tt = t.structType;
		return tt.Field(i);
	};
	rtype.prototype.Field = function(i) { return this.$val.Field(i); };
	rtype.Ptr.prototype.FieldByIndex = function(index) {
		var t, tt;
		t = this;
		if (!((t.Kind() === 25))) {
			$panic(new $String("reflect: FieldByIndex of non-struct type"));
		}
		tt = t.structType;
		return tt.FieldByIndex(index);
	};
	rtype.prototype.FieldByIndex = function(index) { return this.$val.FieldByIndex(index); };
	rtype.Ptr.prototype.FieldByName = function(name) {
		var t, tt;
		t = this;
		if (!((t.Kind() === 25))) {
			$panic(new $String("reflect: FieldByName of non-struct type"));
		}
		tt = t.structType;
		return tt.FieldByName(name);
	};
	rtype.prototype.FieldByName = function(name) { return this.$val.FieldByName(name); };
	rtype.Ptr.prototype.FieldByNameFunc = function(match) {
		var t, tt;
		t = this;
		if (!((t.Kind() === 25))) {
			$panic(new $String("reflect: FieldByNameFunc of non-struct type"));
		}
		tt = t.structType;
		return tt.FieldByNameFunc(match);
	};
	rtype.prototype.FieldByNameFunc = function(match) { return this.$val.FieldByNameFunc(match); };
	rtype.Ptr.prototype.In = function(i) {
		var t, tt, x$1;
		t = this;
		if (!((t.Kind() === 19))) {
			$panic(new $String("reflect: In of non-func type"));
		}
		tt = t.funcType;
		return toType((x$1 = tt.in$2, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i])));
	};
	rtype.prototype.In = function(i) { return this.$val.In(i); };
	rtype.Ptr.prototype.Key = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 21))) {
			$panic(new $String("reflect: Key of non-map type"));
		}
		tt = t.mapType;
		return toType(tt.key);
	};
	rtype.prototype.Key = function() { return this.$val.Key(); };
	rtype.Ptr.prototype.Len = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 17))) {
			$panic(new $String("reflect: Len of non-array type"));
		}
		tt = t.arrayType;
		return (tt.len >> 0);
	};
	rtype.prototype.Len = function() { return this.$val.Len(); };
	rtype.Ptr.prototype.NumField = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 25))) {
			$panic(new $String("reflect: NumField of non-struct type"));
		}
		tt = t.structType;
		return tt.fields.$length;
	};
	rtype.prototype.NumField = function() { return this.$val.NumField(); };
	rtype.Ptr.prototype.NumIn = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 19))) {
			$panic(new $String("reflect: NumIn of non-func type"));
		}
		tt = t.funcType;
		return tt.in$2.$length;
	};
	rtype.prototype.NumIn = function() { return this.$val.NumIn(); };
	rtype.Ptr.prototype.NumOut = function() {
		var t, tt;
		t = this;
		if (!((t.Kind() === 19))) {
			$panic(new $String("reflect: NumOut of non-func type"));
		}
		tt = t.funcType;
		return tt.out.$length;
	};
	rtype.prototype.NumOut = function() { return this.$val.NumOut(); };
	rtype.Ptr.prototype.Out = function(i) {
		var t, tt, x$1;
		t = this;
		if (!((t.Kind() === 19))) {
			$panic(new $String("reflect: Out of non-func type"));
		}
		tt = t.funcType;
		return toType((x$1 = tt.out, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i])));
	};
	rtype.prototype.Out = function(i) { return this.$val.Out(i); };
	ChanDir.prototype.String = function() {
		var d, _ref;
		d = this.$val;
		_ref = d;
		if (_ref === 2) {
			return "chan<-";
		} else if (_ref === 1) {
			return "<-chan";
		} else if (_ref === 3) {
			return "chan";
		}
		return "ChanDir" + strconv.Itoa((d >> 0));
	};
	$ptrType(ChanDir).prototype.String = function() { return new ChanDir(this.$get()).String(); };
	interfaceType.Ptr.prototype.Method = function(i) {
		var m = new Method.Ptr(), t, x$1, p;
		t = this;
		if (i < 0 || i >= t.methods.$length) {
			return m;
		}
		p = (x$1 = t.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
		m.Name = p.name.$get();
		if (!($pointerIsEqual(p.pkgPath, ($ptrType($String)).nil))) {
			m.PkgPath = p.pkgPath.$get();
		}
		m.Type = toType(p.typ);
		m.Index = i;
		return m;
	};
	interfaceType.prototype.Method = function(i) { return this.$val.Method(i); };
	interfaceType.Ptr.prototype.NumMethod = function() {
		var t;
		t = this;
		return t.methods.$length;
	};
	interfaceType.prototype.NumMethod = function() { return this.$val.NumMethod(); };
	interfaceType.Ptr.prototype.MethodByName = function(name) {
		var m = new Method.Ptr(), ok = false, t, p, _ref, _i, i, x$1, _tmp, _tmp$1;
		t = this;
		if (t === ($ptrType(interfaceType)).nil) {
			return [m, ok];
		}
		p = ($ptrType(imethod)).nil;
		_ref = t.methods;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			p = (x$1 = t.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
			if (p.name.$get() === name) {
				_tmp = new Method.Ptr(); $copy(_tmp, t.Method(i), Method); _tmp$1 = true; $copy(m, _tmp, Method); ok = _tmp$1;
				return [m, ok];
			}
			_i++;
		}
		return [m, ok];
	};
	interfaceType.prototype.MethodByName = function(name) { return this.$val.MethodByName(name); };
	StructTag.prototype.Get = function(key) {
		var tag, i, name, qvalue, _tuple, value;
		tag = this.$val;
		while (!(tag === "")) {
			i = 0;
			while (i < tag.length && (tag.charCodeAt(i) === 32)) {
				i = i + (1) >> 0;
			}
			tag = tag.substring(i);
			if (tag === "") {
				break;
			}
			i = 0;
			while (i < tag.length && !((tag.charCodeAt(i) === 32)) && !((tag.charCodeAt(i) === 58)) && !((tag.charCodeAt(i) === 34))) {
				i = i + (1) >> 0;
			}
			if ((i + 1 >> 0) >= tag.length || !((tag.charCodeAt(i) === 58)) || !((tag.charCodeAt((i + 1 >> 0)) === 34))) {
				break;
			}
			name = tag.substring(0, i);
			tag = tag.substring((i + 1 >> 0));
			i = 1;
			while (i < tag.length && !((tag.charCodeAt(i) === 34))) {
				if (tag.charCodeAt(i) === 92) {
					i = i + (1) >> 0;
				}
				i = i + (1) >> 0;
			}
			if (i >= tag.length) {
				break;
			}
			qvalue = tag.substring(0, (i + 1 >> 0));
			tag = tag.substring((i + 1 >> 0));
			if (key === name) {
				_tuple = strconv.Unquote(qvalue); value = _tuple[0];
				return value;
			}
		}
		return "";
	};
	$ptrType(StructTag).prototype.Get = function(key) { return new StructTag(this.$get()).Get(key); };
	structType.Ptr.prototype.Field = function(i) {
		var f = new StructField.Ptr(), t, x$1, p, t$1;
		t = this;
		if (i < 0 || i >= t.fields.$length) {
			return f;
		}
		p = (x$1 = t.fields, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
		f.Type = toType(p.typ);
		if (!($pointerIsEqual(p.name, ($ptrType($String)).nil))) {
			f.Name = p.name.$get();
		} else {
			t$1 = f.Type;
			if (t$1.Kind() === 22) {
				t$1 = t$1.Elem();
			}
			f.Name = t$1.Name();
			f.Anonymous = true;
		}
		if (!($pointerIsEqual(p.pkgPath, ($ptrType($String)).nil))) {
			f.PkgPath = p.pkgPath.$get();
		}
		if (!($pointerIsEqual(p.tag, ($ptrType($String)).nil))) {
			f.Tag = p.tag.$get();
		}
		f.Offset = p.offset;
		f.Index = new ($sliceType($Int))([i]);
		return f;
	};
	structType.prototype.Field = function(i) { return this.$val.Field(i); };
	structType.Ptr.prototype.FieldByIndex = function(index) {
		var f = new StructField.Ptr(), t, _ref, _i, i, x$1, ft;
		t = this;
		f.Type = toType(t.rtype);
		_ref = index;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			x$1 = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (i > 0) {
				ft = f.Type;
				if ((ft.Kind() === 22) && (ft.Elem().Kind() === 25)) {
					ft = ft.Elem();
				}
				f.Type = ft;
			}
			$copy(f, f.Type.Field(x$1), StructField);
			_i++;
		}
		return f;
	};
	structType.prototype.FieldByIndex = function(index) { return this.$val.FieldByIndex(index); };
	structType.Ptr.prototype.FieldByNameFunc = function(match) {
		var result = new StructField.Ptr(), ok = false, t, current, next, nextCount, visited, _map, _key, _tmp, _tmp$1, count, _ref, _i, scan, t$1, _entry, _key$1, _ref$1, _i$1, i, x$1, f, fname, ntyp, _entry$1, _tmp$2, _tmp$3, styp, _entry$2, _key$2, _map$1, _key$3, _key$4, _entry$3, _key$5, index;
		t = this;
		current = new ($sliceType(fieldScan))([]);
		next = new ($sliceType(fieldScan))([new fieldScan.Ptr(t, ($sliceType($Int)).nil)]);
		nextCount = false;
		visited = (_map = new $Map(), _map);
		while (next.$length > 0) {
			_tmp = next; _tmp$1 = $subslice(current, 0, 0); current = _tmp; next = _tmp$1;
			count = nextCount;
			nextCount = false;
			_ref = current;
			_i = 0;
			while (_i < _ref.$length) {
				scan = new fieldScan.Ptr(); $copy(scan, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), fieldScan);
				t$1 = scan.typ;
				if ((_entry = visited[t$1.$key()], _entry !== undefined ? _entry.v : false)) {
					_i++;
					continue;
				}
				_key$1 = t$1; (visited || $throwRuntimeError("assignment to entry in nil map"))[_key$1.$key()] = { k: _key$1, v: true };
				_ref$1 = t$1.fields;
				_i$1 = 0;
				while (_i$1 < _ref$1.$length) {
					i = _i$1;
					f = (x$1 = t$1.fields, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
					fname = "";
					ntyp = ($ptrType(rtype)).nil;
					if (!($pointerIsEqual(f.name, ($ptrType($String)).nil))) {
						fname = f.name.$get();
					} else {
						ntyp = f.typ;
						if (ntyp.Kind() === 22) {
							ntyp = ntyp.Elem().common();
						}
						fname = ntyp.Name();
					}
					if (match(fname)) {
						if ((_entry$1 = count[t$1.$key()], _entry$1 !== undefined ? _entry$1.v : 0) > 1 || ok) {
							_tmp$2 = new StructField.Ptr("", "", null, "", 0, ($sliceType($Int)).nil, false); _tmp$3 = false; $copy(result, _tmp$2, StructField); ok = _tmp$3;
							return [result, ok];
						}
						$copy(result, t$1.Field(i), StructField);
						result.Index = ($sliceType($Int)).nil;
						result.Index = $appendSlice(result.Index, scan.index);
						result.Index = $append(result.Index, i);
						ok = true;
						_i$1++;
						continue;
					}
					if (ok || ntyp === ($ptrType(rtype)).nil || !((ntyp.Kind() === 25))) {
						_i$1++;
						continue;
					}
					styp = ntyp.structType;
					if ((_entry$2 = nextCount[styp.$key()], _entry$2 !== undefined ? _entry$2.v : 0) > 0) {
						_key$2 = styp; (nextCount || $throwRuntimeError("assignment to entry in nil map"))[_key$2.$key()] = { k: _key$2, v: 2 };
						_i$1++;
						continue;
					}
					if (nextCount === false) {
						nextCount = (_map$1 = new $Map(), _map$1);
					}
					_key$4 = styp; (nextCount || $throwRuntimeError("assignment to entry in nil map"))[_key$4.$key()] = { k: _key$4, v: 1 };
					if ((_entry$3 = count[t$1.$key()], _entry$3 !== undefined ? _entry$3.v : 0) > 1) {
						_key$5 = styp; (nextCount || $throwRuntimeError("assignment to entry in nil map"))[_key$5.$key()] = { k: _key$5, v: 2 };
					}
					index = ($sliceType($Int)).nil;
					index = $appendSlice(index, scan.index);
					index = $append(index, i);
					next = $append(next, new fieldScan.Ptr(styp, index));
					_i$1++;
				}
				_i++;
			}
			if (ok) {
				break;
			}
		}
		return [result, ok];
	};
	structType.prototype.FieldByNameFunc = function(match) { return this.$val.FieldByNameFunc(match); };
	structType.Ptr.prototype.FieldByName = function(name) {
		var f = new StructField.Ptr(), present = false, t, hasAnon, _ref, _i, i, x$1, tf, _tmp, _tmp$1, _tuple;
		t = this;
		hasAnon = false;
		if (!(name === "")) {
			_ref = t.fields;
			_i = 0;
			while (_i < _ref.$length) {
				i = _i;
				tf = (x$1 = t.fields, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
				if ($pointerIsEqual(tf.name, ($ptrType($String)).nil)) {
					hasAnon = true;
					_i++;
					continue;
				}
				if (tf.name.$get() === name) {
					_tmp = new StructField.Ptr(); $copy(_tmp, t.Field(i), StructField); _tmp$1 = true; $copy(f, _tmp, StructField); present = _tmp$1;
					return [f, present];
				}
				_i++;
			}
		}
		if (!hasAnon) {
			return [f, present];
		}
		_tuple = t.FieldByNameFunc((function(s) {
			return s === name;
		})); $copy(f, _tuple[0], StructField); present = _tuple[1];
		return [f, present];
	};
	structType.prototype.FieldByName = function(name) { return this.$val.FieldByName(name); };
	PtrTo = $pkg.PtrTo = function(t) {
		return (t !== null && t.constructor === ($ptrType(rtype)) ? t.$val : $typeAssertionFailed(t, ($ptrType(rtype)))).ptrTo();
	};
	rtype.Ptr.prototype.Implements = function(u) {
		var t;
		t = this;
		if ($interfaceIsEqual(u, null)) {
			$panic(new $String("reflect: nil type passed to Type.Implements"));
		}
		if (!((u.Kind() === 20))) {
			$panic(new $String("reflect: non-interface type passed to Type.Implements"));
		}
		return implements$1((u !== null && u.constructor === ($ptrType(rtype)) ? u.$val : $typeAssertionFailed(u, ($ptrType(rtype)))), t);
	};
	rtype.prototype.Implements = function(u) { return this.$val.Implements(u); };
	rtype.Ptr.prototype.AssignableTo = function(u) {
		var t, uu;
		t = this;
		if ($interfaceIsEqual(u, null)) {
			$panic(new $String("reflect: nil type passed to Type.AssignableTo"));
		}
		uu = (u !== null && u.constructor === ($ptrType(rtype)) ? u.$val : $typeAssertionFailed(u, ($ptrType(rtype))));
		return directlyAssignable(uu, t) || implements$1(uu, t);
	};
	rtype.prototype.AssignableTo = function(u) { return this.$val.AssignableTo(u); };
	rtype.Ptr.prototype.ConvertibleTo = function(u) {
		var t, uu;
		t = this;
		if ($interfaceIsEqual(u, null)) {
			$panic(new $String("reflect: nil type passed to Type.ConvertibleTo"));
		}
		uu = (u !== null && u.constructor === ($ptrType(rtype)) ? u.$val : $typeAssertionFailed(u, ($ptrType(rtype))));
		return !(convertOp(uu, t) === $throwNilPointerError);
	};
	rtype.prototype.ConvertibleTo = function(u) { return this.$val.ConvertibleTo(u); };
	implements$1 = function(T, V) {
		var t, v, i, j, x$1, tm, x$2, vm, v$1, i$1, j$1, x$3, tm$1, x$4, vm$1;
		if (!((T.Kind() === 20))) {
			return false;
		}
		t = T.interfaceType;
		if (t.methods.$length === 0) {
			return true;
		}
		if (V.Kind() === 20) {
			v = V.interfaceType;
			i = 0;
			j = 0;
			while (j < v.methods.$length) {
				tm = (x$1 = t.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
				vm = (x$2 = v.methods, ((j < 0 || j >= x$2.$length) ? $throwRuntimeError("index out of range") : x$2.$array[x$2.$offset + j]));
				if ($pointerIsEqual(vm.name, tm.name) && $pointerIsEqual(vm.pkgPath, tm.pkgPath) && vm.typ === tm.typ) {
					i = i + (1) >> 0;
					if (i >= t.methods.$length) {
						return true;
					}
				}
				j = j + (1) >> 0;
			}
			return false;
		}
		v$1 = V.uncommonType.uncommon();
		if (v$1 === ($ptrType(uncommonType)).nil) {
			return false;
		}
		i$1 = 0;
		j$1 = 0;
		while (j$1 < v$1.methods.$length) {
			tm$1 = (x$3 = t.methods, ((i$1 < 0 || i$1 >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + i$1]));
			vm$1 = (x$4 = v$1.methods, ((j$1 < 0 || j$1 >= x$4.$length) ? $throwRuntimeError("index out of range") : x$4.$array[x$4.$offset + j$1]));
			if ($pointerIsEqual(vm$1.name, tm$1.name) && $pointerIsEqual(vm$1.pkgPath, tm$1.pkgPath) && vm$1.mtyp === tm$1.typ) {
				i$1 = i$1 + (1) >> 0;
				if (i$1 >= t.methods.$length) {
					return true;
				}
			}
			j$1 = j$1 + (1) >> 0;
		}
		return false;
	};
	directlyAssignable = function(T, V) {
		if (T === V) {
			return true;
		}
		if (!(T.Name() === "") && !(V.Name() === "") || !((T.Kind() === V.Kind()))) {
			return false;
		}
		return haveIdenticalUnderlyingType(T, V);
	};
	haveIdenticalUnderlyingType = function(T, V) {
		var kind, _ref, t, v, _ref$1, _i, i, typ, x$1, _ref$2, _i$1, i$1, typ$1, x$2, t$1, v$1, t$2, v$2, _ref$3, _i$2, i$2, x$3, tf, x$4, vf;
		if (T === V) {
			return true;
		}
		kind = T.Kind();
		if (!((kind === V.Kind()))) {
			return false;
		}
		if (1 <= kind && kind <= 16 || (kind === 24) || (kind === 26)) {
			return true;
		}
		_ref = kind;
		if (_ref === 17) {
			return $interfaceIsEqual(T.Elem(), V.Elem()) && (T.Len() === V.Len());
		} else if (_ref === 18) {
			if ((V.ChanDir() === 3) && $interfaceIsEqual(T.Elem(), V.Elem())) {
				return true;
			}
			return (V.ChanDir() === T.ChanDir()) && $interfaceIsEqual(T.Elem(), V.Elem());
		} else if (_ref === 19) {
			t = T.funcType;
			v = V.funcType;
			if (!(t.dotdotdot === v.dotdotdot) || !((t.in$2.$length === v.in$2.$length)) || !((t.out.$length === v.out.$length))) {
				return false;
			}
			_ref$1 = t.in$2;
			_i = 0;
			while (_i < _ref$1.$length) {
				i = _i;
				typ = ((_i < 0 || _i >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i]);
				if (!(typ === (x$1 = v.in$2, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i])))) {
					return false;
				}
				_i++;
			}
			_ref$2 = t.out;
			_i$1 = 0;
			while (_i$1 < _ref$2.$length) {
				i$1 = _i$1;
				typ$1 = ((_i$1 < 0 || _i$1 >= _ref$2.$length) ? $throwRuntimeError("index out of range") : _ref$2.$array[_ref$2.$offset + _i$1]);
				if (!(typ$1 === (x$2 = v.out, ((i$1 < 0 || i$1 >= x$2.$length) ? $throwRuntimeError("index out of range") : x$2.$array[x$2.$offset + i$1])))) {
					return false;
				}
				_i$1++;
			}
			return true;
		} else if (_ref === 20) {
			t$1 = T.interfaceType;
			v$1 = V.interfaceType;
			if ((t$1.methods.$length === 0) && (v$1.methods.$length === 0)) {
				return true;
			}
			return false;
		} else if (_ref === 21) {
			return $interfaceIsEqual(T.Key(), V.Key()) && $interfaceIsEqual(T.Elem(), V.Elem());
		} else if (_ref === 22 || _ref === 23) {
			return $interfaceIsEqual(T.Elem(), V.Elem());
		} else if (_ref === 25) {
			t$2 = T.structType;
			v$2 = V.structType;
			if (!((t$2.fields.$length === v$2.fields.$length))) {
				return false;
			}
			_ref$3 = t$2.fields;
			_i$2 = 0;
			while (_i$2 < _ref$3.$length) {
				i$2 = _i$2;
				tf = (x$3 = t$2.fields, ((i$2 < 0 || i$2 >= x$3.$length) ? $throwRuntimeError("index out of range") : x$3.$array[x$3.$offset + i$2]));
				vf = (x$4 = v$2.fields, ((i$2 < 0 || i$2 >= x$4.$length) ? $throwRuntimeError("index out of range") : x$4.$array[x$4.$offset + i$2]));
				if (!($pointerIsEqual(tf.name, vf.name)) && ($pointerIsEqual(tf.name, ($ptrType($String)).nil) || $pointerIsEqual(vf.name, ($ptrType($String)).nil) || !(tf.name.$get() === vf.name.$get()))) {
					return false;
				}
				if (!($pointerIsEqual(tf.pkgPath, vf.pkgPath)) && ($pointerIsEqual(tf.pkgPath, ($ptrType($String)).nil) || $pointerIsEqual(vf.pkgPath, ($ptrType($String)).nil) || !(tf.pkgPath.$get() === vf.pkgPath.$get()))) {
					return false;
				}
				if (!(tf.typ === vf.typ)) {
					return false;
				}
				if (!($pointerIsEqual(tf.tag, vf.tag)) && ($pointerIsEqual(tf.tag, ($ptrType($String)).nil) || $pointerIsEqual(vf.tag, ($ptrType($String)).nil) || !(tf.tag.$get() === vf.tag.$get()))) {
					return false;
				}
				if (!((tf.offset === vf.offset))) {
					return false;
				}
				_i$2++;
			}
			return true;
		}
		return false;
	};
	toType = function(t) {
		if (t === ($ptrType(rtype)).nil) {
			return null;
		}
		return t;
	};
	flag.prototype.kind = function() {
		var f;
		f = this.$val;
		return (((((f >>> 4 >>> 0)) & 31) >>> 0) >>> 0);
	};
	$ptrType(flag).prototype.kind = function() { return new flag(this.$get()).kind(); };
	Value.Ptr.prototype.pointer = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (!((v.typ.size === 4)) || !v.typ.pointers()) {
			$panic(new $String("can't call pointer on a non-pointer Value"));
		}
		if (!((((v.flag & 2) >>> 0) === 0))) {
			return v.ptr.$get();
		}
		return v.ptr;
	};
	Value.prototype.pointer = function() { return this.$val.pointer(); };
	ValueError.Ptr.prototype.Error = function() {
		var e;
		e = this;
		if (e.Kind === 0) {
			return "reflect: call of " + e.Method + " on zero Value";
		}
		return "reflect: call of " + e.Method + " on " + (new Kind(e.Kind)).String() + " Value";
	};
	ValueError.prototype.Error = function() { return this.$val.Error(); };
	flag.prototype.mustBe = function(expected) {
		var f, k;
		f = this.$val;
		k = (new flag(f)).kind();
		if (!((k === expected))) {
			$panic(new ValueError.Ptr(methodName(), k));
		}
	};
	$ptrType(flag).prototype.mustBe = function(expected) { return new flag(this.$get()).mustBe(expected); };
	flag.prototype.mustBeExported = function() {
		var f;
		f = this.$val;
		if (f === 0) {
			$panic(new ValueError.Ptr(methodName(), 0));
		}
		if (!((((f & 1) >>> 0) === 0))) {
			$panic(new $String("reflect: " + methodName() + " using value obtained using unexported field"));
		}
	};
	$ptrType(flag).prototype.mustBeExported = function() { return new flag(this.$get()).mustBeExported(); };
	flag.prototype.mustBeAssignable = function() {
		var f;
		f = this.$val;
		if (f === 0) {
			$panic(new ValueError.Ptr(methodName(), 0));
		}
		if (!((((f & 1) >>> 0) === 0))) {
			$panic(new $String("reflect: " + methodName() + " using value obtained using unexported field"));
		}
		if (((f & 4) >>> 0) === 0) {
			$panic(new $String("reflect: " + methodName() + " using unaddressable value"));
		}
	};
	$ptrType(flag).prototype.mustBeAssignable = function() { return new flag(this.$get()).mustBeAssignable(); };
	Value.Ptr.prototype.Addr = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (((v.flag & 4) >>> 0) === 0) {
			$panic(new $String("reflect.Value.Addr of unaddressable value"));
		}
		return new Value.Ptr(v.typ.ptrTo(), v.ptr, 0, ((((v.flag & 1) >>> 0)) | 352) >>> 0);
	};
	Value.prototype.Addr = function() { return this.$val.Addr(); };
	Value.Ptr.prototype.Bool = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(1);
		if (!((((v.flag & 2) >>> 0) === 0))) {
			return v.ptr.$get();
		}
		return v.scalar;
	};
	Value.prototype.Bool = function() { return this.$val.Bool(); };
	Value.Ptr.prototype.Bytes = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(23);
		if (!((v.typ.Elem().Kind() === 8))) {
			$panic(new $String("reflect.Value.Bytes of non-byte slice"));
		}
		return v.ptr.$get();
	};
	Value.prototype.Bytes = function() { return this.$val.Bytes(); };
	Value.Ptr.prototype.runes = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(23);
		if (!((v.typ.Elem().Kind() === 5))) {
			$panic(new $String("reflect.Value.Bytes of non-rune slice"));
		}
		return v.ptr.$get();
	};
	Value.prototype.runes = function() { return this.$val.runes(); };
	Value.Ptr.prototype.CanAddr = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		return !((((v.flag & 4) >>> 0) === 0));
	};
	Value.prototype.CanAddr = function() { return this.$val.CanAddr(); };
	Value.Ptr.prototype.CanSet = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		return ((v.flag & 5) >>> 0) === 4;
	};
	Value.prototype.CanSet = function() { return this.$val.CanSet(); };
	Value.Ptr.prototype.Call = function(in$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(19);
		(new flag(v.flag)).mustBeExported();
		return v.call("Call", in$1);
	};
	Value.prototype.Call = function(in$1) { return this.$val.Call(in$1); };
	Value.Ptr.prototype.CallSlice = function(in$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(19);
		(new flag(v.flag)).mustBeExported();
		return v.call("CallSlice", in$1);
	};
	Value.prototype.CallSlice = function(in$1) { return this.$val.CallSlice(in$1); };
	Value.Ptr.prototype.Close = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(18);
		(new flag(v.flag)).mustBeExported();
		chanclose(v.pointer());
	};
	Value.prototype.Close = function() { return this.$val.Close(); };
	Value.Ptr.prototype.Complex = function() {
		var v, k, _ref, x$1, x$2;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 15) {
			if (!((((v.flag & 2) >>> 0) === 0))) {
				return (x$1 = v.ptr.$get(), new $Complex128(x$1.$real, x$1.$imag));
			}
			return (x$2 = v.scalar, new $Complex128(x$2.$real, x$2.$imag));
		} else if (_ref === 16) {
			return v.ptr.$get();
		}
		$panic(new ValueError.Ptr("reflect.Value.Complex", k));
	};
	Value.prototype.Complex = function() { return this.$val.Complex(); };
	Value.Ptr.prototype.FieldByIndex = function(index) {
		var v, _ref, _i, i, x$1;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(25);
		_ref = index;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			x$1 = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (i > 0) {
				if ((v.Kind() === 22) && (v.typ.Elem().Kind() === 25)) {
					if (v.IsNil()) {
						$panic(new $String("reflect: indirection through nil pointer to embedded struct"));
					}
					$copy(v, v.Elem(), Value);
				}
			}
			$copy(v, v.Field(x$1), Value);
			_i++;
		}
		return v;
	};
	Value.prototype.FieldByIndex = function(index) { return this.$val.FieldByIndex(index); };
	Value.Ptr.prototype.FieldByName = function(name) {
		var v, _tuple, f, ok;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(25);
		_tuple = v.typ.FieldByName(name); f = new StructField.Ptr(); $copy(f, _tuple[0], StructField); ok = _tuple[1];
		if (ok) {
			return v.FieldByIndex(f.Index);
		}
		return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
	};
	Value.prototype.FieldByName = function(name) { return this.$val.FieldByName(name); };
	Value.Ptr.prototype.FieldByNameFunc = function(match) {
		var v, _tuple, f, ok;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(25);
		_tuple = v.typ.FieldByNameFunc(match); f = new StructField.Ptr(); $copy(f, _tuple[0], StructField); ok = _tuple[1];
		if (ok) {
			return v.FieldByIndex(f.Index);
		}
		return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
	};
	Value.prototype.FieldByNameFunc = function(match) { return this.$val.FieldByNameFunc(match); };
	Value.Ptr.prototype.Float = function() {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 13) {
			if (!((((v.flag & 2) >>> 0) === 0))) {
				return $coerceFloat32(v.ptr.$get());
			}
			return $coerceFloat32(v.scalar);
		} else if (_ref === 14) {
			if (!((((v.flag & 2) >>> 0) === 0))) {
				return v.ptr.$get();
			}
			return v.scalar;
		}
		$panic(new ValueError.Ptr("reflect.Value.Float", k));
	};
	Value.prototype.Float = function() { return this.$val.Float(); };
	Value.Ptr.prototype.Int = function() {
		var v, k, p, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		p = 0;
		if (!((((v.flag & 2) >>> 0) === 0))) {
			p = v.ptr;
		} else {
			p = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, v);
		}
		_ref = k;
		if (_ref === 2) {
			return new $Int64(0, p.$get());
		} else if (_ref === 3) {
			return new $Int64(0, p.$get());
		} else if (_ref === 4) {
			return new $Int64(0, p.$get());
		} else if (_ref === 5) {
			return new $Int64(0, p.$get());
		} else if (_ref === 6) {
			return p.$get();
		}
		$panic(new ValueError.Ptr("reflect.Value.Int", k));
	};
	Value.prototype.Int = function() { return this.$val.Int(); };
	Value.Ptr.prototype.CanInterface = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (v.flag === 0) {
			$panic(new ValueError.Ptr("reflect.Value.CanInterface", 0));
		}
		return ((v.flag & 1) >>> 0) === 0;
	};
	Value.prototype.CanInterface = function() { return this.$val.CanInterface(); };
	Value.Ptr.prototype.Interface = function() {
		var i = null, v;
		v = new Value.Ptr(); $copy(v, this, Value);
		i = valueInterface($clone(v, Value), true);
		return i;
	};
	Value.prototype.Interface = function() { return this.$val.Interface(); };
	Value.Ptr.prototype.InterfaceData = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(20);
		return v.ptr;
	};
	Value.prototype.InterfaceData = function() { return this.$val.InterfaceData(); };
	Value.Ptr.prototype.IsValid = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		return !((v.flag === 0));
	};
	Value.prototype.IsValid = function() { return this.$val.IsValid(); };
	Value.Ptr.prototype.Kind = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		return (new flag(v.flag)).kind();
	};
	Value.prototype.Kind = function() { return this.$val.Kind(); };
	Value.Ptr.prototype.MapIndex = function(key) {
		var v, tt, k, e, typ, fl, c;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(21);
		tt = v.typ.mapType;
		$copy(key, key.assignTo("reflect.Value.MapIndex", tt.key, ($ptrType($emptyInterface)).nil), Value);
		k = 0;
		if (!((((key.flag & 2) >>> 0) === 0))) {
			k = key.ptr;
		} else if (key.typ.pointers()) {
			k = new ($ptrType($UnsafePointer))(function() { return this.$target.ptr; }, function($v) { this.$target.ptr = $v; }, key);
		} else {
			k = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, key);
		}
		e = mapaccess(v.typ, v.pointer(), k);
		if (e === 0) {
			return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
		}
		typ = tt.elem;
		fl = ((((v.flag | key.flag) >>> 0)) & 1) >>> 0;
		fl = (fl | (((typ.Kind() >>> 0) << 4 >>> 0))) >>> 0;
		if (typ.size > 4) {
			c = unsafe_New(typ);
			memmove(c, e, typ.size);
			return new Value.Ptr(typ, c, 0, (fl | 2) >>> 0);
		} else if (typ.pointers()) {
			return new Value.Ptr(typ, e.$get(), 0, fl);
		} else {
			return new Value.Ptr(typ, 0, loadScalar(e, typ.size), fl);
		}
	};
	Value.prototype.MapIndex = function(key) { return this.$val.MapIndex(key); };
	Value.Ptr.prototype.MapKeys = function() {
		var v, tt, keyType, fl, m, mlen, it, a, i, key, c;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(21);
		tt = v.typ.mapType;
		keyType = tt.key;
		fl = (((v.flag & 1) >>> 0) | ((keyType.Kind() >>> 0) << 4 >>> 0)) >>> 0;
		m = v.pointer();
		mlen = 0;
		if (!(m === 0)) {
			mlen = maplen(m);
		}
		it = mapiterinit(v.typ, m);
		a = ($sliceType(Value)).make(mlen);
		i = 0;
		i = 0;
		while (i < a.$length) {
			key = mapiterkey(it);
			if (key === 0) {
				break;
			}
			if (keyType.size > 4) {
				c = unsafe_New(keyType);
				memmove(c, key, keyType.size);
				$copy(((i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i]), new Value.Ptr(keyType, c, 0, (fl | 2) >>> 0), Value);
			} else if (keyType.pointers()) {
				$copy(((i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i]), new Value.Ptr(keyType, key.$get(), 0, fl), Value);
			} else {
				$copy(((i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i]), new Value.Ptr(keyType, 0, loadScalar(key, keyType.size), fl), Value);
			}
			mapiternext(it);
			i = i + (1) >> 0;
		}
		return $subslice(a, 0, i);
	};
	Value.prototype.MapKeys = function() { return this.$val.MapKeys(); };
	Value.Ptr.prototype.Method = function(i) {
		var v, fl;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (v.typ === ($ptrType(rtype)).nil) {
			$panic(new ValueError.Ptr("reflect.Value.Method", 0));
		}
		if (!((((v.flag & 8) >>> 0) === 0)) || i < 0 || i >= v.typ.NumMethod()) {
			$panic(new $String("reflect: Method index out of range"));
		}
		if ((v.typ.Kind() === 20) && v.IsNil()) {
			$panic(new $String("reflect: Method on nil interface value"));
		}
		fl = (v.flag & 3) >>> 0;
		fl = (fl | (304)) >>> 0;
		fl = (fl | (((((i >>> 0) << 9 >>> 0) | 8) >>> 0))) >>> 0;
		return new Value.Ptr(v.typ, v.ptr, v.scalar, fl);
	};
	Value.prototype.Method = function(i) { return this.$val.Method(i); };
	Value.Ptr.prototype.NumMethod = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (v.typ === ($ptrType(rtype)).nil) {
			$panic(new ValueError.Ptr("reflect.Value.NumMethod", 0));
		}
		if (!((((v.flag & 8) >>> 0) === 0))) {
			return 0;
		}
		return v.typ.NumMethod();
	};
	Value.prototype.NumMethod = function() { return this.$val.NumMethod(); };
	Value.Ptr.prototype.MethodByName = function(name) {
		var v, _tuple, m, ok;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (v.typ === ($ptrType(rtype)).nil) {
			$panic(new ValueError.Ptr("reflect.Value.MethodByName", 0));
		}
		if (!((((v.flag & 8) >>> 0) === 0))) {
			return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
		}
		_tuple = v.typ.MethodByName(name); m = new Method.Ptr(); $copy(m, _tuple[0], Method); ok = _tuple[1];
		if (!ok) {
			return new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0);
		}
		return v.Method(m.Index);
	};
	Value.prototype.MethodByName = function(name) { return this.$val.MethodByName(name); };
	Value.Ptr.prototype.NumField = function() {
		var v, tt;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(25);
		tt = v.typ.structType;
		return tt.fields.$length;
	};
	Value.prototype.NumField = function() { return this.$val.NumField(); };
	Value.Ptr.prototype.OverflowComplex = function(x$1) {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 15) {
			return overflowFloat32(x$1.$real) || overflowFloat32(x$1.$imag);
		} else if (_ref === 16) {
			return false;
		}
		$panic(new ValueError.Ptr("reflect.Value.OverflowComplex", k));
	};
	Value.prototype.OverflowComplex = function(x$1) { return this.$val.OverflowComplex(x$1); };
	Value.Ptr.prototype.OverflowFloat = function(x$1) {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 13) {
			return overflowFloat32(x$1);
		} else if (_ref === 14) {
			return false;
		}
		$panic(new ValueError.Ptr("reflect.Value.OverflowFloat", k));
	};
	Value.prototype.OverflowFloat = function(x$1) { return this.$val.OverflowFloat(x$1); };
	overflowFloat32 = function(x$1) {
		if (x$1 < 0) {
			x$1 = -x$1;
		}
		return 3.4028234663852886e+38 < x$1 && x$1 <= 1.7976931348623157e+308;
	};
	Value.Ptr.prototype.OverflowInt = function(x$1) {
		var v, k, _ref, x$2, bitSize, trunc;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 2 || _ref === 3 || _ref === 4 || _ref === 5 || _ref === 6) {
			bitSize = (x$2 = v.typ.size, (((x$2 >>> 16 << 16) * 8 >>> 0) + (x$2 << 16 >>> 16) * 8) >>> 0);
			trunc = $shiftRightInt64(($shiftLeft64(x$1, ((64 - bitSize >>> 0)))), ((64 - bitSize >>> 0)));
			return !((x$1.$high === trunc.$high && x$1.$low === trunc.$low));
		}
		$panic(new ValueError.Ptr("reflect.Value.OverflowInt", k));
	};
	Value.prototype.OverflowInt = function(x$1) { return this.$val.OverflowInt(x$1); };
	Value.Ptr.prototype.OverflowUint = function(x$1) {
		var v, k, _ref, x$2, bitSize, trunc;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 7 || _ref === 12 || _ref === 8 || _ref === 9 || _ref === 10 || _ref === 11) {
			bitSize = (x$2 = v.typ.size, (((x$2 >>> 16 << 16) * 8 >>> 0) + (x$2 << 16 >>> 16) * 8) >>> 0);
			trunc = $shiftRightUint64(($shiftLeft64(x$1, ((64 - bitSize >>> 0)))), ((64 - bitSize >>> 0)));
			return !((x$1.$high === trunc.$high && x$1.$low === trunc.$low));
		}
		$panic(new ValueError.Ptr("reflect.Value.OverflowUint", k));
	};
	Value.prototype.OverflowUint = function(x$1) { return this.$val.OverflowUint(x$1); };
	Value.Ptr.prototype.Recv = function() {
		var x$1 = new Value.Ptr(), ok = false, v, _tuple;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(18);
		(new flag(v.flag)).mustBeExported();
		_tuple = v.recv(false); $copy(x$1, _tuple[0], Value); ok = _tuple[1];
		return [x$1, ok];
	};
	Value.prototype.Recv = function() { return this.$val.Recv(); };
	Value.Ptr.prototype.recv = function(nb) {
		var val = new Value.Ptr(), ok = false, v, tt, t, p, _tuple, selected;
		v = new Value.Ptr(); $copy(v, this, Value);
		tt = v.typ.chanType;
		if (((tt.dir >> 0) & 1) === 0) {
			$panic(new $String("reflect: recv on send-only channel"));
		}
		t = tt.elem;
		$copy(val, new Value.Ptr(t, 0, 0, (t.Kind() >>> 0) << 4 >>> 0), Value);
		p = 0;
		if (t.size > 4) {
			p = unsafe_New(t);
			val.ptr = p;
			val.flag = (val.flag | (2)) >>> 0;
		} else if (t.pointers()) {
			p = new ($ptrType($UnsafePointer))(function() { return this.$target.ptr; }, function($v) { this.$target.ptr = $v; }, val);
		} else {
			p = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, val);
		}
		_tuple = chanrecv(v.typ, v.pointer(), nb, p); selected = _tuple[0]; ok = _tuple[1];
		if (!selected) {
			$copy(val, new Value.Ptr(($ptrType(rtype)).nil, 0, 0, 0), Value);
		}
		return [val, ok];
	};
	Value.prototype.recv = function(nb) { return this.$val.recv(nb); };
	Value.Ptr.prototype.Send = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(18);
		(new flag(v.flag)).mustBeExported();
		v.send($clone(x$1, Value), false);
	};
	Value.prototype.Send = function(x$1) { return this.$val.Send(x$1); };
	Value.Ptr.prototype.send = function(x$1, nb) {
		var selected = false, v, tt, p;
		v = new Value.Ptr(); $copy(v, this, Value);
		tt = v.typ.chanType;
		if (((tt.dir >> 0) & 2) === 0) {
			$panic(new $String("reflect: send on recv-only channel"));
		}
		(new flag(x$1.flag)).mustBeExported();
		$copy(x$1, x$1.assignTo("reflect.Value.Send", tt.elem, ($ptrType($emptyInterface)).nil), Value);
		p = 0;
		if (!((((x$1.flag & 2) >>> 0) === 0))) {
			p = x$1.ptr;
		} else if (x$1.typ.pointers()) {
			p = new ($ptrType($UnsafePointer))(function() { return this.$target.ptr; }, function($v) { this.$target.ptr = $v; }, x$1);
		} else {
			p = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, x$1);
		}
		selected = chansend(v.typ, v.pointer(), p, nb);
		return selected;
	};
	Value.prototype.send = function(x$1, nb) { return this.$val.send(x$1, nb); };
	Value.Ptr.prototype.SetBool = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(1);
		v.ptr.$set(x$1);
	};
	Value.prototype.SetBool = function(x$1) { return this.$val.SetBool(x$1); };
	Value.Ptr.prototype.SetBytes = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(23);
		if (!((v.typ.Elem().Kind() === 8))) {
			$panic(new $String("reflect.Value.SetBytes of non-byte slice"));
		}
		v.ptr.$set(x$1);
	};
	Value.prototype.SetBytes = function(x$1) { return this.$val.SetBytes(x$1); };
	Value.Ptr.prototype.setRunes = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(23);
		if (!((v.typ.Elem().Kind() === 5))) {
			$panic(new $String("reflect.Value.setRunes of non-rune slice"));
		}
		v.ptr.$set(x$1);
	};
	Value.prototype.setRunes = function(x$1) { return this.$val.setRunes(x$1); };
	Value.Ptr.prototype.SetComplex = function(x$1) {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 15) {
			v.ptr.$set(new $Complex64(x$1.$real, x$1.$imag));
		} else if (_ref === 16) {
			v.ptr.$set(x$1);
		} else {
			$panic(new ValueError.Ptr("reflect.Value.SetComplex", k));
		}
	};
	Value.prototype.SetComplex = function(x$1) { return this.$val.SetComplex(x$1); };
	Value.Ptr.prototype.SetFloat = function(x$1) {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 13) {
			v.ptr.$set(x$1);
		} else if (_ref === 14) {
			v.ptr.$set(x$1);
		} else {
			$panic(new ValueError.Ptr("reflect.Value.SetFloat", k));
		}
	};
	Value.prototype.SetFloat = function(x$1) { return this.$val.SetFloat(x$1); };
	Value.Ptr.prototype.SetInt = function(x$1) {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 2) {
			v.ptr.$set(((x$1.$low + ((x$1.$high >> 31) * 4294967296)) >> 0));
		} else if (_ref === 3) {
			v.ptr.$set(((x$1.$low + ((x$1.$high >> 31) * 4294967296)) << 24 >> 24));
		} else if (_ref === 4) {
			v.ptr.$set(((x$1.$low + ((x$1.$high >> 31) * 4294967296)) << 16 >> 16));
		} else if (_ref === 5) {
			v.ptr.$set(((x$1.$low + ((x$1.$high >> 31) * 4294967296)) >> 0));
		} else if (_ref === 6) {
			v.ptr.$set(x$1);
		} else {
			$panic(new ValueError.Ptr("reflect.Value.SetInt", k));
		}
	};
	Value.prototype.SetInt = function(x$1) { return this.$val.SetInt(x$1); };
	Value.Ptr.prototype.SetMapIndex = function(key, val) {
		var v, tt, k, e;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(21);
		(new flag(v.flag)).mustBeExported();
		(new flag(key.flag)).mustBeExported();
		tt = v.typ.mapType;
		$copy(key, key.assignTo("reflect.Value.SetMapIndex", tt.key, ($ptrType($emptyInterface)).nil), Value);
		k = 0;
		if (!((((key.flag & 2) >>> 0) === 0))) {
			k = key.ptr;
		} else if (key.typ.pointers()) {
			k = new ($ptrType($UnsafePointer))(function() { return this.$target.ptr; }, function($v) { this.$target.ptr = $v; }, key);
		} else {
			k = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, key);
		}
		if (val.typ === ($ptrType(rtype)).nil) {
			mapdelete(v.typ, v.pointer(), k);
			return;
		}
		(new flag(val.flag)).mustBeExported();
		$copy(val, val.assignTo("reflect.Value.SetMapIndex", tt.elem, ($ptrType($emptyInterface)).nil), Value);
		e = 0;
		if (!((((val.flag & 2) >>> 0) === 0))) {
			e = val.ptr;
		} else if (val.typ.pointers()) {
			e = new ($ptrType($UnsafePointer))(function() { return this.$target.ptr; }, function($v) { this.$target.ptr = $v; }, val);
		} else {
			e = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, val);
		}
		mapassign(v.typ, v.pointer(), k, e);
	};
	Value.prototype.SetMapIndex = function(key, val) { return this.$val.SetMapIndex(key, val); };
	Value.Ptr.prototype.SetUint = function(x$1) {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 7) {
			v.ptr.$set((x$1.$low >>> 0));
		} else if (_ref === 8) {
			v.ptr.$set((x$1.$low << 24 >>> 24));
		} else if (_ref === 9) {
			v.ptr.$set((x$1.$low << 16 >>> 16));
		} else if (_ref === 10) {
			v.ptr.$set((x$1.$low >>> 0));
		} else if (_ref === 11) {
			v.ptr.$set(x$1);
		} else if (_ref === 12) {
			v.ptr.$set((x$1.$low >>> 0));
		} else {
			$panic(new ValueError.Ptr("reflect.Value.SetUint", k));
		}
	};
	Value.prototype.SetUint = function(x$1) { return this.$val.SetUint(x$1); };
	Value.Ptr.prototype.SetPointer = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(26);
		v.ptr.$set(x$1);
	};
	Value.prototype.SetPointer = function(x$1) { return this.$val.SetPointer(x$1); };
	Value.Ptr.prototype.SetString = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBeAssignable();
		(new flag(v.flag)).mustBe(24);
		v.ptr.$set(x$1);
	};
	Value.prototype.SetString = function(x$1) { return this.$val.SetString(x$1); };
	Value.Ptr.prototype.String = function() {
		var v, k, _ref;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		_ref = k;
		if (_ref === 0) {
			return "<invalid Value>";
		} else if (_ref === 24) {
			return v.ptr.$get();
		}
		return "<" + v.typ.String() + " Value>";
	};
	Value.prototype.String = function() { return this.$val.String(); };
	Value.Ptr.prototype.TryRecv = function() {
		var x$1 = new Value.Ptr(), ok = false, v, _tuple;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(18);
		(new flag(v.flag)).mustBeExported();
		_tuple = v.recv(true); $copy(x$1, _tuple[0], Value); ok = _tuple[1];
		return [x$1, ok];
	};
	Value.prototype.TryRecv = function() { return this.$val.TryRecv(); };
	Value.Ptr.prototype.TrySend = function(x$1) {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		(new flag(v.flag)).mustBe(18);
		(new flag(v.flag)).mustBeExported();
		return v.send($clone(x$1, Value), true);
	};
	Value.prototype.TrySend = function(x$1) { return this.$val.TrySend(x$1); };
	Value.Ptr.prototype.Type = function() {
		var v, f, i, tt, x$1, m, ut, x$2, m$1;
		v = new Value.Ptr(); $copy(v, this, Value);
		f = v.flag;
		if (f === 0) {
			$panic(new ValueError.Ptr("reflect.Value.Type", 0));
		}
		if (((f & 8) >>> 0) === 0) {
			return v.typ;
		}
		i = (v.flag >> 0) >> 9 >> 0;
		if (v.typ.Kind() === 20) {
			tt = v.typ.interfaceType;
			if (i < 0 || i >= tt.methods.$length) {
				$panic(new $String("reflect: internal error: invalid method index"));
			}
			m = (x$1 = tt.methods, ((i < 0 || i >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + i]));
			return m.typ;
		}
		ut = v.typ.uncommonType.uncommon();
		if (ut === ($ptrType(uncommonType)).nil || i < 0 || i >= ut.methods.$length) {
			$panic(new $String("reflect: internal error: invalid method index"));
		}
		m$1 = (x$2 = ut.methods, ((i < 0 || i >= x$2.$length) ? $throwRuntimeError("index out of range") : x$2.$array[x$2.$offset + i]));
		return m$1.mtyp;
	};
	Value.prototype.Type = function() { return this.$val.Type(); };
	Value.Ptr.prototype.Uint = function() {
		var v, k, p, _ref, x$1;
		v = new Value.Ptr(); $copy(v, this, Value);
		k = (new flag(v.flag)).kind();
		p = 0;
		if (!((((v.flag & 2) >>> 0) === 0))) {
			p = v.ptr;
		} else {
			p = new ($ptrType($Uintptr))(function() { return this.$target.scalar; }, function($v) { this.$target.scalar = $v; }, v);
		}
		_ref = k;
		if (_ref === 7) {
			return new $Uint64(0, p.$get());
		} else if (_ref === 8) {
			return new $Uint64(0, p.$get());
		} else if (_ref === 9) {
			return new $Uint64(0, p.$get());
		} else if (_ref === 10) {
			return new $Uint64(0, p.$get());
		} else if (_ref === 11) {
			return p.$get();
		} else if (_ref === 12) {
			return (x$1 = p.$get(), new $Uint64(0, x$1.constructor === Number ? x$1 : 1));
		}
		$panic(new ValueError.Ptr("reflect.Value.Uint", k));
	};
	Value.prototype.Uint = function() { return this.$val.Uint(); };
	Value.Ptr.prototype.UnsafeAddr = function() {
		var v;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (v.typ === ($ptrType(rtype)).nil) {
			$panic(new ValueError.Ptr("reflect.Value.UnsafeAddr", 0));
		}
		if (((v.flag & 4) >>> 0) === 0) {
			$panic(new $String("reflect.Value.UnsafeAddr of unaddressable value"));
		}
		return v.ptr;
	};
	Value.prototype.UnsafeAddr = function() { return this.$val.UnsafeAddr(); };
	New = $pkg.New = function(typ) {
		var ptr, fl;
		if ($interfaceIsEqual(typ, null)) {
			$panic(new $String("reflect: New(nil)"));
		}
		ptr = unsafe_New((typ !== null && typ.constructor === ($ptrType(rtype)) ? typ.$val : $typeAssertionFailed(typ, ($ptrType(rtype)))));
		fl = 352;
		return new Value.Ptr(typ.common().ptrTo(), ptr, 0, fl);
	};
	Value.Ptr.prototype.assignTo = function(context, dst, target) {
		var v, fl, x$1;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (!((((v.flag & 8) >>> 0) === 0))) {
			$copy(v, makeMethodValue(context, $clone(v, Value)), Value);
		}
		if (directlyAssignable(dst, v.typ)) {
			v.typ = dst;
			fl = (v.flag & 7) >>> 0;
			fl = (fl | (((dst.Kind() >>> 0) << 4 >>> 0))) >>> 0;
			return new Value.Ptr(dst, v.ptr, v.scalar, fl);
		} else if (implements$1(dst, v.typ)) {
			if (target === ($ptrType($emptyInterface)).nil) {
				target = $newDataPointer(null, ($ptrType($emptyInterface)));
			}
			x$1 = valueInterface($clone(v, Value), false);
			if (dst.NumMethod() === 0) {
				target.$set(x$1);
			} else {
				ifaceE2I(dst, x$1, target);
			}
			return new Value.Ptr(dst, target, 0, 322);
		}
		$panic(new $String(context + ": value of type " + v.typ.String() + " is not assignable to type " + dst.String()));
	};
	Value.prototype.assignTo = function(context, dst, target) { return this.$val.assignTo(context, dst, target); };
	Value.Ptr.prototype.Convert = function(t) {
		var v, op;
		v = new Value.Ptr(); $copy(v, this, Value);
		if (!((((v.flag & 8) >>> 0) === 0))) {
			$copy(v, makeMethodValue("Convert", $clone(v, Value)), Value);
		}
		op = convertOp(t.common(), v.typ);
		if (op === $throwNilPointerError) {
			$panic(new $String("reflect.Value.Convert: value of type " + v.typ.String() + " cannot be converted to type " + t.String()));
		}
		return op($clone(v, Value), t);
	};
	Value.prototype.Convert = function(t) { return this.$val.Convert(t); };
	convertOp = function(dst, src) {
		var _ref, _ref$1, _ref$2, _ref$3, _ref$4, _ref$5, _ref$6;
		_ref = src.Kind();
		if (_ref === 2 || _ref === 3 || _ref === 4 || _ref === 5 || _ref === 6) {
			_ref$1 = dst.Kind();
			if (_ref$1 === 2 || _ref$1 === 3 || _ref$1 === 4 || _ref$1 === 5 || _ref$1 === 6 || _ref$1 === 7 || _ref$1 === 8 || _ref$1 === 9 || _ref$1 === 10 || _ref$1 === 11 || _ref$1 === 12) {
				return cvtInt;
			} else if (_ref$1 === 13 || _ref$1 === 14) {
				return cvtIntFloat;
			} else if (_ref$1 === 24) {
				return cvtIntString;
			}
		} else if (_ref === 7 || _ref === 8 || _ref === 9 || _ref === 10 || _ref === 11 || _ref === 12) {
			_ref$2 = dst.Kind();
			if (_ref$2 === 2 || _ref$2 === 3 || _ref$2 === 4 || _ref$2 === 5 || _ref$2 === 6 || _ref$2 === 7 || _ref$2 === 8 || _ref$2 === 9 || _ref$2 === 10 || _ref$2 === 11 || _ref$2 === 12) {
				return cvtUint;
			} else if (_ref$2 === 13 || _ref$2 === 14) {
				return cvtUintFloat;
			} else if (_ref$2 === 24) {
				return cvtUintString;
			}
		} else if (_ref === 13 || _ref === 14) {
			_ref$3 = dst.Kind();
			if (_ref$3 === 2 || _ref$3 === 3 || _ref$3 === 4 || _ref$3 === 5 || _ref$3 === 6) {
				return cvtFloatInt;
			} else if (_ref$3 === 7 || _ref$3 === 8 || _ref$3 === 9 || _ref$3 === 10 || _ref$3 === 11 || _ref$3 === 12) {
				return cvtFloatUint;
			} else if (_ref$3 === 13 || _ref$3 === 14) {
				return cvtFloat;
			}
		} else if (_ref === 15 || _ref === 16) {
			_ref$4 = dst.Kind();
			if (_ref$4 === 15 || _ref$4 === 16) {
				return cvtComplex;
			}
		} else if (_ref === 24) {
			if ((dst.Kind() === 23) && dst.Elem().PkgPath() === "") {
				_ref$5 = dst.Elem().Kind();
				if (_ref$5 === 8) {
					return cvtStringBytes;
				} else if (_ref$5 === 5) {
					return cvtStringRunes;
				}
			}
		} else if (_ref === 23) {
			if ((dst.Kind() === 24) && src.Elem().PkgPath() === "") {
				_ref$6 = src.Elem().Kind();
				if (_ref$6 === 8) {
					return cvtBytesString;
				} else if (_ref$6 === 5) {
					return cvtRunesString;
				}
			}
		}
		if (haveIdenticalUnderlyingType(dst, src)) {
			return cvtDirect;
		}
		if ((dst.Kind() === 22) && dst.Name() === "" && (src.Kind() === 22) && src.Name() === "" && haveIdenticalUnderlyingType(dst.Elem().common(), src.Elem().common())) {
			return cvtDirect;
		}
		if (implements$1(dst, src)) {
			if (src.Kind() === 20) {
				return cvtI2I;
			}
			return cvtT2I;
		}
		return $throwNilPointerError;
	};
	makeFloat = function(f, v, t) {
		var typ, ptr, s, _ref;
		typ = t.common();
		if (typ.size > 4) {
			ptr = unsafe_New(typ);
			ptr.$set(v);
			return new Value.Ptr(typ, ptr, 0, (((f | 2) >>> 0) | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
		}
		s = 0;
		_ref = typ.size;
		if (_ref === 4) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set(v);
		} else if (_ref === 8) {
			new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set(v);
		}
		return new Value.Ptr(typ, 0, s, (f | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
	};
	makeComplex = function(f, v, t) {
		var typ, ptr, _ref, s;
		typ = t.common();
		if (typ.size > 4) {
			ptr = unsafe_New(typ);
			_ref = typ.size;
			if (_ref === 8) {
				ptr.$set(new $Complex64(v.$real, v.$imag));
			} else if (_ref === 16) {
				ptr.$set(v);
			}
			return new Value.Ptr(typ, ptr, 0, (((f | 2) >>> 0) | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
		}
		s = 0;
		new ($ptrType($Uintptr))(function() { return s; }, function($v) { s = $v; }).$set(new $Complex64(v.$real, v.$imag));
		return new Value.Ptr(typ, 0, s, (f | ((typ.Kind() >>> 0) << 4 >>> 0)) >>> 0);
	};
	makeString = function(f, v, t) {
		var ret;
		ret = new Value.Ptr(); $copy(ret, New(t).Elem(), Value);
		ret.SetString(v);
		ret.flag = ((ret.flag & ~4) | f) >>> 0;
		return ret;
	};
	makeBytes = function(f, v, t) {
		var ret;
		ret = new Value.Ptr(); $copy(ret, New(t).Elem(), Value);
		ret.SetBytes(v);
		ret.flag = ((ret.flag & ~4) | f) >>> 0;
		return ret;
	};
	makeRunes = function(f, v, t) {
		var ret;
		ret = new Value.Ptr(); $copy(ret, New(t).Elem(), Value);
		ret.setRunes(v);
		ret.flag = ((ret.flag & ~4) | f) >>> 0;
		return ret;
	};
	cvtInt = function(v, t) {
		var x$1;
		return makeInt((v.flag & 1) >>> 0, (x$1 = v.Int(), new $Uint64(x$1.$high, x$1.$low)), t);
	};
	cvtUint = function(v, t) {
		return makeInt((v.flag & 1) >>> 0, v.Uint(), t);
	};
	cvtFloatInt = function(v, t) {
		var x$1;
		return makeInt((v.flag & 1) >>> 0, (x$1 = new $Int64(0, v.Float()), new $Uint64(x$1.$high, x$1.$low)), t);
	};
	cvtFloatUint = function(v, t) {
		return makeInt((v.flag & 1) >>> 0, new $Uint64(0, v.Float()), t);
	};
	cvtIntFloat = function(v, t) {
		return makeFloat((v.flag & 1) >>> 0, $flatten64(v.Int()), t);
	};
	cvtUintFloat = function(v, t) {
		return makeFloat((v.flag & 1) >>> 0, $flatten64(v.Uint()), t);
	};
	cvtFloat = function(v, t) {
		return makeFloat((v.flag & 1) >>> 0, v.Float(), t);
	};
	cvtComplex = function(v, t) {
		return makeComplex((v.flag & 1) >>> 0, v.Complex(), t);
	};
	cvtIntString = function(v, t) {
		return makeString((v.flag & 1) >>> 0, $encodeRune(v.Int().$low), t);
	};
	cvtUintString = function(v, t) {
		return makeString((v.flag & 1) >>> 0, $encodeRune(v.Uint().$low), t);
	};
	cvtBytesString = function(v, t) {
		return makeString((v.flag & 1) >>> 0, $bytesToString(v.Bytes()), t);
	};
	cvtStringBytes = function(v, t) {
		return makeBytes((v.flag & 1) >>> 0, new ($sliceType($Uint8))($stringToBytes(v.String())), t);
	};
	cvtRunesString = function(v, t) {
		return makeString((v.flag & 1) >>> 0, $runesToString(v.runes()), t);
	};
	cvtStringRunes = function(v, t) {
		return makeRunes((v.flag & 1) >>> 0, new ($sliceType($Int32))($stringToRunes(v.String())), t);
	};
	cvtT2I = function(v, typ) {
		var target, x$1;
		target = $newDataPointer(null, ($ptrType($emptyInterface)));
		x$1 = valueInterface($clone(v, Value), false);
		if (typ.NumMethod() === 0) {
			target.$set(x$1);
		} else {
			ifaceE2I((typ !== null && typ.constructor === ($ptrType(rtype)) ? typ.$val : $typeAssertionFailed(typ, ($ptrType(rtype)))), x$1, target);
		}
		return new Value.Ptr(typ.common(), target, 0, (((((v.flag & 1) >>> 0) | 2) >>> 0) | 320) >>> 0);
	};
	cvtI2I = function(v, typ) {
		var ret;
		if (v.IsNil()) {
			ret = new Value.Ptr(); $copy(ret, Zero(typ), Value);
			ret.flag = (ret.flag | (((v.flag & 1) >>> 0))) >>> 0;
			return ret;
		}
		return cvtT2I($clone(v.Elem(), Value), typ);
	};
	call = function() {
		$panic("Native function not implemented: reflect.call");
	};
	$pkg.$init = function() {
		mapIter.init([["t", "t", "reflect", Type, ""], ["m", "m", "reflect", js.Object, ""], ["keys", "keys", "reflect", js.Object, ""], ["i", "i", "reflect", $Int, ""]]);
		Type.init([["Align", "Align", "", [], [$Int], false], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false], ["Bits", "Bits", "", [], [$Int], false], ["ChanDir", "ChanDir", "", [], [ChanDir], false], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false], ["Elem", "Elem", "", [], [Type], false], ["Field", "Field", "", [$Int], [StructField], false], ["FieldAlign", "FieldAlign", "", [], [$Int], false], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false], ["Implements", "Implements", "", [Type], [$Bool], false], ["In", "In", "", [$Int], [Type], false], ["IsVariadic", "IsVariadic", "", [], [$Bool], false], ["Key", "Key", "", [], [Type], false], ["Kind", "Kind", "", [], [Kind], false], ["Len", "Len", "", [], [$Int], false], ["Method", "Method", "", [$Int], [Method], false], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false], ["Name", "Name", "", [], [$String], false], ["NumField", "NumField", "", [], [$Int], false], ["NumIn", "NumIn", "", [], [$Int], false], ["NumMethod", "NumMethod", "", [], [$Int], false], ["NumOut", "NumOut", "", [], [$Int], false], ["Out", "Out", "", [$Int], [Type], false], ["PkgPath", "PkgPath", "", [], [$String], false], ["Size", "Size", "", [], [$Uintptr], false], ["String", "String", "", [], [$String], false], ["common", "common", "reflect", [], [($ptrType(rtype))], false], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false]]);
		Kind.methods = [["String", "String", "", [], [$String], false, -1]];
		($ptrType(Kind)).methods = [["String", "String", "", [], [$String], false, -1]];
		rtype.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 9]];
		($ptrType(rtype)).methods = [["Align", "Align", "", [], [$Int], false, -1], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, -1], ["Bits", "Bits", "", [], [$Int], false, -1], ["ChanDir", "ChanDir", "", [], [ChanDir], false, -1], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, -1], ["Elem", "Elem", "", [], [Type], false, -1], ["Field", "Field", "", [$Int], [StructField], false, -1], ["FieldAlign", "FieldAlign", "", [], [$Int], false, -1], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, -1], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, -1], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, -1], ["Implements", "Implements", "", [Type], [$Bool], false, -1], ["In", "In", "", [$Int], [Type], false, -1], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, -1], ["Key", "Key", "", [], [Type], false, -1], ["Kind", "Kind", "", [], [Kind], false, -1], ["Len", "Len", "", [], [$Int], false, -1], ["Method", "Method", "", [$Int], [Method], false, -1], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, -1], ["Name", "Name", "", [], [$String], false, -1], ["NumField", "NumField", "", [], [$Int], false, -1], ["NumIn", "NumIn", "", [], [$Int], false, -1], ["NumMethod", "NumMethod", "", [], [$Int], false, -1], ["NumOut", "NumOut", "", [], [$Int], false, -1], ["Out", "Out", "", [$Int], [Type], false, -1], ["PkgPath", "PkgPath", "", [], [$String], false, -1], ["Size", "Size", "", [], [$Uintptr], false, -1], ["String", "String", "", [], [$String], false, -1], ["common", "common", "reflect", [], [($ptrType(rtype))], false, -1], ["pointers", "pointers", "reflect", [], [$Bool], false, -1], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, -1], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 9]];
		rtype.init([["size", "size", "reflect", $Uintptr, ""], ["hash", "hash", "reflect", $Uint32, ""], ["_$2", "_", "reflect", $Uint8, ""], ["align", "align", "reflect", $Uint8, ""], ["fieldAlign", "fieldAlign", "reflect", $Uint8, ""], ["kind", "kind", "reflect", $Uint8, ""], ["alg", "alg", "reflect", ($ptrType($Uintptr)), ""], ["gc", "gc", "reflect", $UnsafePointer, ""], ["string", "string", "reflect", ($ptrType($String)), ""], ["uncommonType", "", "reflect", ($ptrType(uncommonType)), ""], ["ptrToThis", "ptrToThis", "reflect", ($ptrType(rtype)), ""], ["zero", "zero", "reflect", $UnsafePointer, ""]]);
		method.init([["name", "name", "reflect", ($ptrType($String)), ""], ["pkgPath", "pkgPath", "reflect", ($ptrType($String)), ""], ["mtyp", "mtyp", "reflect", ($ptrType(rtype)), ""], ["typ", "typ", "reflect", ($ptrType(rtype)), ""], ["ifn", "ifn", "reflect", $UnsafePointer, ""], ["tfn", "tfn", "reflect", $UnsafePointer, ""]]);
		($ptrType(uncommonType)).methods = [["Method", "Method", "", [$Int], [Method], false, -1], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, -1], ["Name", "Name", "", [], [$String], false, -1], ["NumMethod", "NumMethod", "", [], [$Int], false, -1], ["PkgPath", "PkgPath", "", [], [$String], false, -1], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, -1]];
		uncommonType.init([["name", "name", "reflect", ($ptrType($String)), ""], ["pkgPath", "pkgPath", "reflect", ($ptrType($String)), ""], ["methods", "methods", "reflect", ($sliceType(method)), ""]]);
		ChanDir.methods = [["String", "String", "", [], [$String], false, -1]];
		($ptrType(ChanDir)).methods = [["String", "String", "", [], [$String], false, -1]];
		arrayType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(arrayType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		arrayType.init([["rtype", "", "reflect", rtype, "reflect:\"array\""], ["elem", "elem", "reflect", ($ptrType(rtype)), ""], ["slice", "slice", "reflect", ($ptrType(rtype)), ""], ["len", "len", "reflect", $Uintptr, ""]]);
		chanType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(chanType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		chanType.init([["rtype", "", "reflect", rtype, "reflect:\"chan\""], ["elem", "elem", "reflect", ($ptrType(rtype)), ""], ["dir", "dir", "reflect", $Uintptr, ""]]);
		funcType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(funcType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		funcType.init([["rtype", "", "reflect", rtype, "reflect:\"func\""], ["dotdotdot", "dotdotdot", "reflect", $Bool, ""], ["in$2", "in", "reflect", ($sliceType(($ptrType(rtype)))), ""], ["out", "out", "reflect", ($sliceType(($ptrType(rtype)))), ""]]);
		imethod.init([["name", "name", "reflect", ($ptrType($String)), ""], ["pkgPath", "pkgPath", "reflect", ($ptrType($String)), ""], ["typ", "typ", "reflect", ($ptrType(rtype)), ""]]);
		interfaceType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(interfaceType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, -1], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, -1], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, -1], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		interfaceType.init([["rtype", "", "reflect", rtype, "reflect:\"interface\""], ["methods", "methods", "reflect", ($sliceType(imethod)), ""]]);
		mapType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(mapType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		mapType.init([["rtype", "", "reflect", rtype, "reflect:\"map\""], ["key", "key", "reflect", ($ptrType(rtype)), ""], ["elem", "elem", "reflect", ($ptrType(rtype)), ""], ["bucket", "bucket", "reflect", ($ptrType(rtype)), ""], ["hmap", "hmap", "reflect", ($ptrType(rtype)), ""]]);
		ptrType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(ptrType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		ptrType.init([["rtype", "", "reflect", rtype, "reflect:\"ptr\""], ["elem", "elem", "reflect", ($ptrType(rtype)), ""]]);
		sliceType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(sliceType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, 0], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, 0], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, 0], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, 0], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		sliceType.init([["rtype", "", "reflect", rtype, "reflect:\"slice\""], ["elem", "elem", "reflect", ($ptrType(rtype)), ""]]);
		structField.init([["name", "name", "reflect", ($ptrType($String)), ""], ["pkgPath", "pkgPath", "reflect", ($ptrType($String)), ""], ["typ", "typ", "reflect", ($ptrType(rtype)), ""], ["tag", "tag", "reflect", ($ptrType($String)), ""], ["offset", "offset", "reflect", $Uintptr, ""]]);
		structType.methods = [["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		($ptrType(structType)).methods = [["Align", "Align", "", [], [$Int], false, 0], ["AssignableTo", "AssignableTo", "", [Type], [$Bool], false, 0], ["Bits", "Bits", "", [], [$Int], false, 0], ["ChanDir", "ChanDir", "", [], [ChanDir], false, 0], ["ConvertibleTo", "ConvertibleTo", "", [Type], [$Bool], false, 0], ["Elem", "Elem", "", [], [Type], false, 0], ["Field", "Field", "", [$Int], [StructField], false, -1], ["FieldAlign", "FieldAlign", "", [], [$Int], false, 0], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [StructField], false, -1], ["FieldByName", "FieldByName", "", [$String], [StructField, $Bool], false, -1], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [StructField, $Bool], false, -1], ["Implements", "Implements", "", [Type], [$Bool], false, 0], ["In", "In", "", [$Int], [Type], false, 0], ["IsVariadic", "IsVariadic", "", [], [$Bool], false, 0], ["Key", "Key", "", [], [Type], false, 0], ["Kind", "Kind", "", [], [Kind], false, 0], ["Len", "Len", "", [], [$Int], false, 0], ["Method", "Method", "", [$Int], [Method], false, 0], ["MethodByName", "MethodByName", "", [$String], [Method, $Bool], false, 0], ["Name", "Name", "", [], [$String], false, 0], ["NumField", "NumField", "", [], [$Int], false, 0], ["NumIn", "NumIn", "", [], [$Int], false, 0], ["NumMethod", "NumMethod", "", [], [$Int], false, 0], ["NumOut", "NumOut", "", [], [$Int], false, 0], ["Out", "Out", "", [$Int], [Type], false, 0], ["PkgPath", "PkgPath", "", [], [$String], false, 0], ["Size", "Size", "", [], [$Uintptr], false, 0], ["String", "String", "", [], [$String], false, 0], ["common", "common", "reflect", [], [($ptrType(rtype))], false, 0], ["pointers", "pointers", "reflect", [], [$Bool], false, 0], ["ptrTo", "ptrTo", "reflect", [], [($ptrType(rtype))], false, 0], ["uncommon", "uncommon", "reflect", [], [($ptrType(uncommonType))], false, 0]];
		structType.init([["rtype", "", "reflect", rtype, "reflect:\"struct\""], ["fields", "fields", "reflect", ($sliceType(structField)), ""]]);
		Method.init([["Name", "Name", "", $String, ""], ["PkgPath", "PkgPath", "", $String, ""], ["Type", "Type", "", Type, ""], ["Func", "Func", "", Value, ""], ["Index", "Index", "", $Int, ""]]);
		StructField.init([["Name", "Name", "", $String, ""], ["PkgPath", "PkgPath", "", $String, ""], ["Type", "Type", "", Type, ""], ["Tag", "Tag", "", StructTag, ""], ["Offset", "Offset", "", $Uintptr, ""], ["Index", "Index", "", ($sliceType($Int)), ""], ["Anonymous", "Anonymous", "", $Bool, ""]]);
		StructTag.methods = [["Get", "Get", "", [$String], [$String], false, -1]];
		($ptrType(StructTag)).methods = [["Get", "Get", "", [$String], [$String], false, -1]];
		fieldScan.init([["typ", "typ", "reflect", ($ptrType(structType)), ""], ["index", "index", "reflect", ($sliceType($Int)), ""]]);
		Value.methods = [["Addr", "Addr", "", [], [Value], false, -1], ["Bool", "Bool", "", [], [$Bool], false, -1], ["Bytes", "Bytes", "", [], [($sliceType($Uint8))], false, -1], ["Call", "Call", "", [($sliceType(Value))], [($sliceType(Value))], false, -1], ["CallSlice", "CallSlice", "", [($sliceType(Value))], [($sliceType(Value))], false, -1], ["CanAddr", "CanAddr", "", [], [$Bool], false, -1], ["CanInterface", "CanInterface", "", [], [$Bool], false, -1], ["CanSet", "CanSet", "", [], [$Bool], false, -1], ["Cap", "Cap", "", [], [$Int], false, -1], ["Close", "Close", "", [], [], false, -1], ["Complex", "Complex", "", [], [$Complex128], false, -1], ["Convert", "Convert", "", [Type], [Value], false, -1], ["Elem", "Elem", "", [], [Value], false, -1], ["Field", "Field", "", [$Int], [Value], false, -1], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [Value], false, -1], ["FieldByName", "FieldByName", "", [$String], [Value], false, -1], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [Value], false, -1], ["Float", "Float", "", [], [$Float64], false, -1], ["Index", "Index", "", [$Int], [Value], false, -1], ["Int", "Int", "", [], [$Int64], false, -1], ["Interface", "Interface", "", [], [$emptyInterface], false, -1], ["InterfaceData", "InterfaceData", "", [], [($arrayType($Uintptr, 2))], false, -1], ["IsNil", "IsNil", "", [], [$Bool], false, -1], ["IsValid", "IsValid", "", [], [$Bool], false, -1], ["Kind", "Kind", "", [], [Kind], false, -1], ["Len", "Len", "", [], [$Int], false, -1], ["MapIndex", "MapIndex", "", [Value], [Value], false, -1], ["MapKeys", "MapKeys", "", [], [($sliceType(Value))], false, -1], ["Method", "Method", "", [$Int], [Value], false, -1], ["MethodByName", "MethodByName", "", [$String], [Value], false, -1], ["NumField", "NumField", "", [], [$Int], false, -1], ["NumMethod", "NumMethod", "", [], [$Int], false, -1], ["OverflowComplex", "OverflowComplex", "", [$Complex128], [$Bool], false, -1], ["OverflowFloat", "OverflowFloat", "", [$Float64], [$Bool], false, -1], ["OverflowInt", "OverflowInt", "", [$Int64], [$Bool], false, -1], ["OverflowUint", "OverflowUint", "", [$Uint64], [$Bool], false, -1], ["Pointer", "Pointer", "", [], [$Uintptr], false, -1], ["Recv", "Recv", "", [], [Value, $Bool], false, -1], ["Send", "Send", "", [Value], [], false, -1], ["Set", "Set", "", [Value], [], false, -1], ["SetBool", "SetBool", "", [$Bool], [], false, -1], ["SetBytes", "SetBytes", "", [($sliceType($Uint8))], [], false, -1], ["SetCap", "SetCap", "", [$Int], [], false, -1], ["SetComplex", "SetComplex", "", [$Complex128], [], false, -1], ["SetFloat", "SetFloat", "", [$Float64], [], false, -1], ["SetInt", "SetInt", "", [$Int64], [], false, -1], ["SetLen", "SetLen", "", [$Int], [], false, -1], ["SetMapIndex", "SetMapIndex", "", [Value, Value], [], false, -1], ["SetPointer", "SetPointer", "", [$UnsafePointer], [], false, -1], ["SetString", "SetString", "", [$String], [], false, -1], ["SetUint", "SetUint", "", [$Uint64], [], false, -1], ["Slice", "Slice", "", [$Int, $Int], [Value], false, -1], ["Slice3", "Slice3", "", [$Int, $Int, $Int], [Value], false, -1], ["String", "String", "", [], [$String], false, -1], ["TryRecv", "TryRecv", "", [], [Value, $Bool], false, -1], ["TrySend", "TrySend", "", [Value], [$Bool], false, -1], ["Type", "Type", "", [], [Type], false, -1], ["Uint", "Uint", "", [], [$Uint64], false, -1], ["UnsafeAddr", "UnsafeAddr", "", [], [$Uintptr], false, -1], ["assignTo", "assignTo", "reflect", [$String, ($ptrType(rtype)), ($ptrType($emptyInterface))], [Value], false, -1], ["call", "call", "reflect", [$String, ($sliceType(Value))], [($sliceType(Value))], false, -1], ["iword", "iword", "reflect", [], [iword], false, -1], ["kind", "kind", "reflect", [], [Kind], false, 3], ["mustBe", "mustBe", "reflect", [Kind], [], false, 3], ["mustBeAssignable", "mustBeAssignable", "reflect", [], [], false, 3], ["mustBeExported", "mustBeExported", "reflect", [], [], false, 3], ["pointer", "pointer", "reflect", [], [$UnsafePointer], false, -1], ["recv", "recv", "reflect", [$Bool], [Value, $Bool], false, -1], ["runes", "runes", "reflect", [], [($sliceType($Int32))], false, -1], ["send", "send", "reflect", [Value, $Bool], [$Bool], false, -1], ["setRunes", "setRunes", "reflect", [($sliceType($Int32))], [], false, -1]];
		($ptrType(Value)).methods = [["Addr", "Addr", "", [], [Value], false, -1], ["Bool", "Bool", "", [], [$Bool], false, -1], ["Bytes", "Bytes", "", [], [($sliceType($Uint8))], false, -1], ["Call", "Call", "", [($sliceType(Value))], [($sliceType(Value))], false, -1], ["CallSlice", "CallSlice", "", [($sliceType(Value))], [($sliceType(Value))], false, -1], ["CanAddr", "CanAddr", "", [], [$Bool], false, -1], ["CanInterface", "CanInterface", "", [], [$Bool], false, -1], ["CanSet", "CanSet", "", [], [$Bool], false, -1], ["Cap", "Cap", "", [], [$Int], false, -1], ["Close", "Close", "", [], [], false, -1], ["Complex", "Complex", "", [], [$Complex128], false, -1], ["Convert", "Convert", "", [Type], [Value], false, -1], ["Elem", "Elem", "", [], [Value], false, -1], ["Field", "Field", "", [$Int], [Value], false, -1], ["FieldByIndex", "FieldByIndex", "", [($sliceType($Int))], [Value], false, -1], ["FieldByName", "FieldByName", "", [$String], [Value], false, -1], ["FieldByNameFunc", "FieldByNameFunc", "", [($funcType([$String], [$Bool], false))], [Value], false, -1], ["Float", "Float", "", [], [$Float64], false, -1], ["Index", "Index", "", [$Int], [Value], false, -1], ["Int", "Int", "", [], [$Int64], false, -1], ["Interface", "Interface", "", [], [$emptyInterface], false, -1], ["InterfaceData", "InterfaceData", "", [], [($arrayType($Uintptr, 2))], false, -1], ["IsNil", "IsNil", "", [], [$Bool], false, -1], ["IsValid", "IsValid", "", [], [$Bool], false, -1], ["Kind", "Kind", "", [], [Kind], false, -1], ["Len", "Len", "", [], [$Int], false, -1], ["MapIndex", "MapIndex", "", [Value], [Value], false, -1], ["MapKeys", "MapKeys", "", [], [($sliceType(Value))], false, -1], ["Method", "Method", "", [$Int], [Value], false, -1], ["MethodByName", "MethodByName", "", [$String], [Value], false, -1], ["NumField", "NumField", "", [], [$Int], false, -1], ["NumMethod", "NumMethod", "", [], [$Int], false, -1], ["OverflowComplex", "OverflowComplex", "", [$Complex128], [$Bool], false, -1], ["OverflowFloat", "OverflowFloat", "", [$Float64], [$Bool], false, -1], ["OverflowInt", "OverflowInt", "", [$Int64], [$Bool], false, -1], ["OverflowUint", "OverflowUint", "", [$Uint64], [$Bool], false, -1], ["Pointer", "Pointer", "", [], [$Uintptr], false, -1], ["Recv", "Recv", "", [], [Value, $Bool], false, -1], ["Send", "Send", "", [Value], [], false, -1], ["Set", "Set", "", [Value], [], false, -1], ["SetBool", "SetBool", "", [$Bool], [], false, -1], ["SetBytes", "SetBytes", "", [($sliceType($Uint8))], [], false, -1], ["SetCap", "SetCap", "", [$Int], [], false, -1], ["SetComplex", "SetComplex", "", [$Complex128], [], false, -1], ["SetFloat", "SetFloat", "", [$Float64], [], false, -1], ["SetInt", "SetInt", "", [$Int64], [], false, -1], ["SetLen", "SetLen", "", [$Int], [], false, -1], ["SetMapIndex", "SetMapIndex", "", [Value, Value], [], false, -1], ["SetPointer", "SetPointer", "", [$UnsafePointer], [], false, -1], ["SetString", "SetString", "", [$String], [], false, -1], ["SetUint", "SetUint", "", [$Uint64], [], false, -1], ["Slice", "Slice", "", [$Int, $Int], [Value], false, -1], ["Slice3", "Slice3", "", [$Int, $Int, $Int], [Value], false, -1], ["String", "String", "", [], [$String], false, -1], ["TryRecv", "TryRecv", "", [], [Value, $Bool], false, -1], ["TrySend", "TrySend", "", [Value], [$Bool], false, -1], ["Type", "Type", "", [], [Type], false, -1], ["Uint", "Uint", "", [], [$Uint64], false, -1], ["UnsafeAddr", "UnsafeAddr", "", [], [$Uintptr], false, -1], ["assignTo", "assignTo", "reflect", [$String, ($ptrType(rtype)), ($ptrType($emptyInterface))], [Value], false, -1], ["call", "call", "reflect", [$String, ($sliceType(Value))], [($sliceType(Value))], false, -1], ["iword", "iword", "reflect", [], [iword], false, -1], ["kind", "kind", "reflect", [], [Kind], false, 3], ["mustBe", "mustBe", "reflect", [Kind], [], false, 3], ["mustBeAssignable", "mustBeAssignable", "reflect", [], [], false, 3], ["mustBeExported", "mustBeExported", "reflect", [], [], false, 3], ["pointer", "pointer", "reflect", [], [$UnsafePointer], false, -1], ["recv", "recv", "reflect", [$Bool], [Value, $Bool], false, -1], ["runes", "runes", "reflect", [], [($sliceType($Int32))], false, -1], ["send", "send", "reflect", [Value, $Bool], [$Bool], false, -1], ["setRunes", "setRunes", "reflect", [($sliceType($Int32))], [], false, -1]];
		Value.init([["typ", "typ", "reflect", ($ptrType(rtype)), ""], ["ptr", "ptr", "reflect", $UnsafePointer, ""], ["scalar", "scalar", "reflect", $Uintptr, ""], ["flag", "", "reflect", flag, ""]]);
		flag.methods = [["kind", "kind", "reflect", [], [Kind], false, -1], ["mustBe", "mustBe", "reflect", [Kind], [], false, -1], ["mustBeAssignable", "mustBeAssignable", "reflect", [], [], false, -1], ["mustBeExported", "mustBeExported", "reflect", [], [], false, -1]];
		($ptrType(flag)).methods = [["kind", "kind", "reflect", [], [Kind], false, -1], ["mustBe", "mustBe", "reflect", [Kind], [], false, -1], ["mustBeAssignable", "mustBeAssignable", "reflect", [], [], false, -1], ["mustBeExported", "mustBeExported", "reflect", [], [], false, -1]];
		($ptrType(ValueError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		ValueError.init([["Method", "Method", "", $String, ""], ["Kind", "Kind", "", Kind, ""]]);
		nonEmptyInterface.init([["itab", "itab", "reflect", ($ptrType(($structType([["ityp", "ityp", "reflect", ($ptrType(rtype)), ""], ["typ", "typ", "reflect", ($ptrType(rtype)), ""], ["link", "link", "reflect", $UnsafePointer, ""], ["bad", "bad", "reflect", $Int32, ""], ["unused", "unused", "reflect", $Int32, ""], ["fun", "fun", "reflect", ($arrayType($UnsafePointer, 100000)), ""]])))), ""], ["word", "word", "reflect", iword, ""]]);
		initialized = false;
		kindNames = new ($sliceType($String))(["invalid", "bool", "int", "int8", "int16", "int32", "int64", "uint", "uint8", "uint16", "uint32", "uint64", "uintptr", "float32", "float64", "complex64", "complex128", "array", "chan", "func", "interface", "map", "ptr", "slice", "string", "struct", "unsafe.Pointer"]);
		uint8Type = (x = TypeOf(new $Uint8(0)), (x !== null && x.constructor === ($ptrType(rtype)) ? x.$val : $typeAssertionFailed(x, ($ptrType(rtype)))));
		init();
	};
	return $pkg;
})();
$packages["fmt"] = (function() {
	var $pkg = {}, math = $packages["math"], strconv = $packages["strconv"], utf8 = $packages["unicode/utf8"], errors = $packages["errors"], io = $packages["io"], os = $packages["os"], reflect = $packages["reflect"], sync = $packages["sync"], fmt, State, Formatter, Stringer, GoStringer, buffer, pp, runeUnreader, scanError, ss, ssave, padZeroBytes, padSpaceBytes, trueBytes, falseBytes, commaSpaceBytes, nilAngleBytes, nilParenBytes, nilBytes, mapBytes, percentBangBytes, missingBytes, badIndexBytes, panicBytes, extraBytes, irparenBytes, bytesBytes, badWidthBytes, badPrecBytes, noVerbBytes, ppFree, intBits, uintptrBits, space, ssFree, complexError, boolError, init, doPrec, newPrinter, Fprintf, Printf, Sprintf, Fprint, Sprint, Fprintln, Println, getField, parsenum, intFromArg, parseArgNumber, isSpace, notSpace, indexRune;
	fmt = $pkg.fmt = $newType(0, "Struct", "fmt.fmt", "fmt", "fmt", function(intbuf_, buf_, wid_, prec_, widPresent_, precPresent_, minus_, plus_, sharp_, space_, unicode_, uniQuote_, zero_) {
		this.$val = this;
		this.intbuf = intbuf_ !== undefined ? intbuf_ : ($arrayType($Uint8, 65)).zero();
		this.buf = buf_ !== undefined ? buf_ : ($ptrType(buffer)).nil;
		this.wid = wid_ !== undefined ? wid_ : 0;
		this.prec = prec_ !== undefined ? prec_ : 0;
		this.widPresent = widPresent_ !== undefined ? widPresent_ : false;
		this.precPresent = precPresent_ !== undefined ? precPresent_ : false;
		this.minus = minus_ !== undefined ? minus_ : false;
		this.plus = plus_ !== undefined ? plus_ : false;
		this.sharp = sharp_ !== undefined ? sharp_ : false;
		this.space = space_ !== undefined ? space_ : false;
		this.unicode = unicode_ !== undefined ? unicode_ : false;
		this.uniQuote = uniQuote_ !== undefined ? uniQuote_ : false;
		this.zero = zero_ !== undefined ? zero_ : false;
	});
	State = $pkg.State = $newType(8, "Interface", "fmt.State", "State", "fmt", null);
	Formatter = $pkg.Formatter = $newType(8, "Interface", "fmt.Formatter", "Formatter", "fmt", null);
	Stringer = $pkg.Stringer = $newType(8, "Interface", "fmt.Stringer", "Stringer", "fmt", null);
	GoStringer = $pkg.GoStringer = $newType(8, "Interface", "fmt.GoStringer", "GoStringer", "fmt", null);
	buffer = $pkg.buffer = $newType(12, "Slice", "fmt.buffer", "buffer", "fmt", null);
	pp = $pkg.pp = $newType(0, "Struct", "fmt.pp", "pp", "fmt", function(n_, panicking_, erroring_, buf_, arg_, value_, reordered_, goodArgNum_, runeBuf_, fmt_) {
		this.$val = this;
		this.n = n_ !== undefined ? n_ : 0;
		this.panicking = panicking_ !== undefined ? panicking_ : false;
		this.erroring = erroring_ !== undefined ? erroring_ : false;
		this.buf = buf_ !== undefined ? buf_ : buffer.nil;
		this.arg = arg_ !== undefined ? arg_ : null;
		this.value = value_ !== undefined ? value_ : new reflect.Value.Ptr();
		this.reordered = reordered_ !== undefined ? reordered_ : false;
		this.goodArgNum = goodArgNum_ !== undefined ? goodArgNum_ : false;
		this.runeBuf = runeBuf_ !== undefined ? runeBuf_ : ($arrayType($Uint8, 4)).zero();
		this.fmt = fmt_ !== undefined ? fmt_ : new fmt.Ptr();
	});
	runeUnreader = $pkg.runeUnreader = $newType(8, "Interface", "fmt.runeUnreader", "runeUnreader", "fmt", null);
	scanError = $pkg.scanError = $newType(0, "Struct", "fmt.scanError", "scanError", "fmt", function(err_) {
		this.$val = this;
		this.err = err_ !== undefined ? err_ : null;
	});
	ss = $pkg.ss = $newType(0, "Struct", "fmt.ss", "ss", "fmt", function(rr_, buf_, peekRune_, prevRune_, count_, atEOF_, ssave_) {
		this.$val = this;
		this.rr = rr_ !== undefined ? rr_ : null;
		this.buf = buf_ !== undefined ? buf_ : buffer.nil;
		this.peekRune = peekRune_ !== undefined ? peekRune_ : 0;
		this.prevRune = prevRune_ !== undefined ? prevRune_ : 0;
		this.count = count_ !== undefined ? count_ : 0;
		this.atEOF = atEOF_ !== undefined ? atEOF_ : false;
		this.ssave = ssave_ !== undefined ? ssave_ : new ssave.Ptr();
	});
	ssave = $pkg.ssave = $newType(0, "Struct", "fmt.ssave", "ssave", "fmt", function(validSave_, nlIsEnd_, nlIsSpace_, argLimit_, limit_, maxWid_) {
		this.$val = this;
		this.validSave = validSave_ !== undefined ? validSave_ : false;
		this.nlIsEnd = nlIsEnd_ !== undefined ? nlIsEnd_ : false;
		this.nlIsSpace = nlIsSpace_ !== undefined ? nlIsSpace_ : false;
		this.argLimit = argLimit_ !== undefined ? argLimit_ : 0;
		this.limit = limit_ !== undefined ? limit_ : 0;
		this.maxWid = maxWid_ !== undefined ? maxWid_ : 0;
	});
	init = function() {
		var i;
		i = 0;
		while (i < 65) {
			(i < 0 || i >= padZeroBytes.$length) ? $throwRuntimeError("index out of range") : padZeroBytes.$array[padZeroBytes.$offset + i] = 48;
			(i < 0 || i >= padSpaceBytes.$length) ? $throwRuntimeError("index out of range") : padSpaceBytes.$array[padSpaceBytes.$offset + i] = 32;
			i = i + (1) >> 0;
		}
	};
	fmt.Ptr.prototype.clearflags = function() {
		var f;
		f = this;
		f.wid = 0;
		f.widPresent = false;
		f.prec = 0;
		f.precPresent = false;
		f.minus = false;
		f.plus = false;
		f.sharp = false;
		f.space = false;
		f.unicode = false;
		f.uniQuote = false;
		f.zero = false;
	};
	fmt.prototype.clearflags = function() { return this.$val.clearflags(); };
	fmt.Ptr.prototype.init = function(buf) {
		var f;
		f = this;
		f.buf = buf;
		f.clearflags();
	};
	fmt.prototype.init = function(buf) { return this.$val.init(buf); };
	fmt.Ptr.prototype.computePadding = function(width) {
		var padding = ($sliceType($Uint8)).nil, leftWidth = 0, rightWidth = 0, f, left, w, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8;
		f = this;
		left = !f.minus;
		w = f.wid;
		if (w < 0) {
			left = false;
			w = -w;
		}
		w = w - (width) >> 0;
		if (w > 0) {
			if (left && f.zero) {
				_tmp = padZeroBytes; _tmp$1 = w; _tmp$2 = 0; padding = _tmp; leftWidth = _tmp$1; rightWidth = _tmp$2;
				return [padding, leftWidth, rightWidth];
			}
			if (left) {
				_tmp$3 = padSpaceBytes; _tmp$4 = w; _tmp$5 = 0; padding = _tmp$3; leftWidth = _tmp$4; rightWidth = _tmp$5;
				return [padding, leftWidth, rightWidth];
			} else {
				_tmp$6 = padSpaceBytes; _tmp$7 = 0; _tmp$8 = w; padding = _tmp$6; leftWidth = _tmp$7; rightWidth = _tmp$8;
				return [padding, leftWidth, rightWidth];
			}
		}
		return [padding, leftWidth, rightWidth];
	};
	fmt.prototype.computePadding = function(width) { return this.$val.computePadding(width); };
	fmt.Ptr.prototype.writePadding = function(n, padding) {
		var f, m;
		f = this;
		while (n > 0) {
			m = n;
			if (m > 65) {
				m = 65;
			}
			f.buf.Write($subslice(padding, 0, m));
			n = n - (m) >> 0;
		}
	};
	fmt.prototype.writePadding = function(n, padding) { return this.$val.writePadding(n, padding); };
	fmt.Ptr.prototype.pad = function(b) {
		var f, _tuple, padding, left, right;
		f = this;
		if (!f.widPresent || (f.wid === 0)) {
			f.buf.Write(b);
			return;
		}
		_tuple = f.computePadding(b.$length); padding = _tuple[0]; left = _tuple[1]; right = _tuple[2];
		if (left > 0) {
			f.writePadding(left, padding);
		}
		f.buf.Write(b);
		if (right > 0) {
			f.writePadding(right, padding);
		}
	};
	fmt.prototype.pad = function(b) { return this.$val.pad(b); };
	fmt.Ptr.prototype.padString = function(s) {
		var f, _tuple, padding, left, right;
		f = this;
		if (!f.widPresent || (f.wid === 0)) {
			f.buf.WriteString(s);
			return;
		}
		_tuple = f.computePadding(utf8.RuneCountInString(s)); padding = _tuple[0]; left = _tuple[1]; right = _tuple[2];
		if (left > 0) {
			f.writePadding(left, padding);
		}
		f.buf.WriteString(s);
		if (right > 0) {
			f.writePadding(right, padding);
		}
	};
	fmt.prototype.padString = function(s) { return this.$val.padString(s); };
	fmt.Ptr.prototype.fmt_boolean = function(v) {
		var f;
		f = this;
		if (v) {
			f.pad(trueBytes);
		} else {
			f.pad(falseBytes);
		}
	};
	fmt.prototype.fmt_boolean = function(v) { return this.$val.fmt_boolean(v); };
	fmt.Ptr.prototype.integer = function(a, base, signedness, digits) {
		var f, buf, width, negative, prec, i, ua, _ref, runeWidth, width$1, j;
		f = this;
		if (f.precPresent && (f.prec === 0) && (a.$high === 0 && a.$low === 0)) {
			return;
		}
		buf = $subslice(new ($sliceType($Uint8))(f.intbuf), 0);
		if (f.widPresent) {
			width = f.wid;
			if ((base.$high === 0 && base.$low === 16) && f.sharp) {
				width = width + (2) >> 0;
			}
			if (width > 65) {
				buf = ($sliceType($Uint8)).make(width);
			}
		}
		negative = signedness === true && (a.$high < 0 || (a.$high === 0 && a.$low < 0));
		if (negative) {
			a = new $Int64(-a.$high, -a.$low);
		}
		prec = 0;
		if (f.precPresent) {
			prec = f.prec;
			f.zero = false;
		} else if (f.zero && f.widPresent && !f.minus && f.wid > 0) {
			prec = f.wid;
			if (negative || f.plus || f.space) {
				prec = prec - (1) >> 0;
			}
		}
		i = buf.$length;
		ua = new $Uint64(a.$high, a.$low);
		while ((ua.$high > base.$high || (ua.$high === base.$high && ua.$low >= base.$low))) {
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = digits.charCodeAt($flatten64($div64(ua, base, true)));
			ua = $div64(ua, (base), false);
		}
		i = i - (1) >> 0;
		(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = digits.charCodeAt($flatten64(ua));
		while (i > 0 && prec > (buf.$length - i >> 0)) {
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 48;
		}
		if (f.sharp) {
			_ref = base;
			if ((_ref.$high === 0 && _ref.$low === 8)) {
				if (!((((i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i]) === 48))) {
					i = i - (1) >> 0;
					(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 48;
				}
			} else if ((_ref.$high === 0 && _ref.$low === 16)) {
				i = i - (1) >> 0;
				(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = (120 + digits.charCodeAt(10) << 24 >>> 24) - 97 << 24 >>> 24;
				i = i - (1) >> 0;
				(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 48;
			}
		}
		if (f.unicode) {
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 43;
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 85;
		}
		if (negative) {
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 45;
		} else if (f.plus) {
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 43;
		} else if (f.space) {
			i = i - (1) >> 0;
			(i < 0 || i >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + i] = 32;
		}
		if (f.unicode && f.uniQuote && (a.$high > 0 || (a.$high === 0 && a.$low >= 0)) && (a.$high < 0 || (a.$high === 0 && a.$low <= 1114111)) && strconv.IsPrint(((a.$low + ((a.$high >> 31) * 4294967296)) >> 0))) {
			runeWidth = utf8.RuneLen(((a.$low + ((a.$high >> 31) * 4294967296)) >> 0));
			width$1 = (2 + runeWidth >> 0) + 1 >> 0;
			$copySlice($subslice(buf, (i - width$1 >> 0)), $subslice(buf, i));
			i = i - (width$1) >> 0;
			j = buf.$length - width$1 >> 0;
			(j < 0 || j >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + j] = 32;
			j = j + (1) >> 0;
			(j < 0 || j >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + j] = 39;
			j = j + (1) >> 0;
			utf8.EncodeRune($subslice(buf, j), ((a.$low + ((a.$high >> 31) * 4294967296)) >> 0));
			j = j + (runeWidth) >> 0;
			(j < 0 || j >= buf.$length) ? $throwRuntimeError("index out of range") : buf.$array[buf.$offset + j] = 39;
		}
		f.pad($subslice(buf, i));
	};
	fmt.prototype.integer = function(a, base, signedness, digits) { return this.$val.integer(a, base, signedness, digits); };
	fmt.Ptr.prototype.truncate = function(s) {
		var f, n, _ref, _i, _rune, i;
		f = this;
		if (f.precPresent && f.prec < utf8.RuneCountInString(s)) {
			n = f.prec;
			_ref = s;
			_i = 0;
			while (_i < _ref.length) {
				_rune = $decodeRune(_ref, _i);
				i = _i;
				if (n === 0) {
					s = s.substring(0, i);
					break;
				}
				n = n - (1) >> 0;
				_i += _rune[1];
			}
		}
		return s;
	};
	fmt.prototype.truncate = function(s) { return this.$val.truncate(s); };
	fmt.Ptr.prototype.fmt_s = function(s) {
		var f;
		f = this;
		s = f.truncate(s);
		f.padString(s);
	};
	fmt.prototype.fmt_s = function(s) { return this.$val.fmt_s(s); };
	fmt.Ptr.prototype.fmt_sbx = function(s, b, digits) {
		var f, n, x, buf, i, c;
		f = this;
		n = b.$length;
		if (b === ($sliceType($Uint8)).nil) {
			n = s.length;
		}
		x = (digits.charCodeAt(10) - 97 << 24 >>> 24) + 120 << 24 >>> 24;
		buf = ($sliceType($Uint8)).nil;
		i = 0;
		while (i < n) {
			if (i > 0 && f.space) {
				buf = $append(buf, 32);
			}
			if (f.sharp) {
				buf = $append(buf, 48, x);
			}
			c = 0;
			if (b === ($sliceType($Uint8)).nil) {
				c = s.charCodeAt(i);
			} else {
				c = ((i < 0 || i >= b.$length) ? $throwRuntimeError("index out of range") : b.$array[b.$offset + i]);
			}
			buf = $append(buf, digits.charCodeAt((c >>> 4 << 24 >>> 24)), digits.charCodeAt(((c & 15) >>> 0)));
			i = i + (1) >> 0;
		}
		f.pad(buf);
	};
	fmt.prototype.fmt_sbx = function(s, b, digits) { return this.$val.fmt_sbx(s, b, digits); };
	fmt.Ptr.prototype.fmt_sx = function(s, digits) {
		var f;
		f = this;
		f.fmt_sbx(s, ($sliceType($Uint8)).nil, digits);
	};
	fmt.prototype.fmt_sx = function(s, digits) { return this.$val.fmt_sx(s, digits); };
	fmt.Ptr.prototype.fmt_bx = function(b, digits) {
		var f;
		f = this;
		f.fmt_sbx("", b, digits);
	};
	fmt.prototype.fmt_bx = function(b, digits) { return this.$val.fmt_bx(b, digits); };
	fmt.Ptr.prototype.fmt_q = function(s) {
		var f, quoted;
		f = this;
		s = f.truncate(s);
		quoted = "";
		if (f.sharp && strconv.CanBackquote(s)) {
			quoted = "`" + s + "`";
		} else {
			if (f.plus) {
				quoted = strconv.QuoteToASCII(s);
			} else {
				quoted = strconv.Quote(s);
			}
		}
		f.padString(quoted);
	};
	fmt.prototype.fmt_q = function(s) { return this.$val.fmt_q(s); };
	fmt.Ptr.prototype.fmt_qc = function(c) {
		var f, quoted;
		f = this;
		quoted = ($sliceType($Uint8)).nil;
		if (f.plus) {
			quoted = strconv.AppendQuoteRuneToASCII($subslice(new ($sliceType($Uint8))(f.intbuf), 0, 0), ((c.$low + ((c.$high >> 31) * 4294967296)) >> 0));
		} else {
			quoted = strconv.AppendQuoteRune($subslice(new ($sliceType($Uint8))(f.intbuf), 0, 0), ((c.$low + ((c.$high >> 31) * 4294967296)) >> 0));
		}
		f.pad(quoted);
	};
	fmt.prototype.fmt_qc = function(c) { return this.$val.fmt_qc(c); };
	doPrec = function(f, def) {
		if (f.precPresent) {
			return f.prec;
		}
		return def;
	};
	fmt.Ptr.prototype.formatFloat = function(v, verb, prec, n) {
		var $deferred = [], $err = null, f, num;
		/* */ try { $deferFrames.push($deferred);
		f = this;
		num = strconv.AppendFloat($subslice(new ($sliceType($Uint8))(f.intbuf), 0, 1), v, verb, prec, n);
		if ((((1 < 0 || 1 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 1]) === 45) || (((1 < 0 || 1 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 1]) === 43)) {
			num = $subslice(num, 1);
		} else {
			(0 < 0 || 0 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 0] = 43;
		}
		if (math.IsInf(v, 0)) {
			if (f.zero) {
				$deferred.push([(function() {
					f.zero = true;
				}), []]);
				f.zero = false;
			}
		}
		if (f.zero && f.widPresent && f.wid > num.$length) {
			if (f.space && v >= 0) {
				f.buf.WriteByte(32);
				f.wid = f.wid - (1) >> 0;
			} else if (f.plus || v < 0) {
				f.buf.WriteByte(((0 < 0 || 0 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 0]));
				f.wid = f.wid - (1) >> 0;
			}
			f.pad($subslice(num, 1));
			return;
		}
		if (f.space && (((0 < 0 || 0 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 0]) === 43)) {
			(0 < 0 || 0 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 0] = 32;
			f.pad(num);
			return;
		}
		if (f.plus || (((0 < 0 || 0 >= num.$length) ? $throwRuntimeError("index out of range") : num.$array[num.$offset + 0]) === 45) || math.IsInf(v, 0)) {
			f.pad(num);
			return;
		}
		f.pad($subslice(num, 1));
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	fmt.prototype.formatFloat = function(v, verb, prec, n) { return this.$val.formatFloat(v, verb, prec, n); };
	fmt.Ptr.prototype.fmt_e64 = function(v) {
		var f;
		f = this;
		f.formatFloat(v, 101, doPrec(f, 6), 64);
	};
	fmt.prototype.fmt_e64 = function(v) { return this.$val.fmt_e64(v); };
	fmt.Ptr.prototype.fmt_E64 = function(v) {
		var f;
		f = this;
		f.formatFloat(v, 69, doPrec(f, 6), 64);
	};
	fmt.prototype.fmt_E64 = function(v) { return this.$val.fmt_E64(v); };
	fmt.Ptr.prototype.fmt_f64 = function(v) {
		var f;
		f = this;
		f.formatFloat(v, 102, doPrec(f, 6), 64);
	};
	fmt.prototype.fmt_f64 = function(v) { return this.$val.fmt_f64(v); };
	fmt.Ptr.prototype.fmt_g64 = function(v) {
		var f;
		f = this;
		f.formatFloat(v, 103, doPrec(f, -1), 64);
	};
	fmt.prototype.fmt_g64 = function(v) { return this.$val.fmt_g64(v); };
	fmt.Ptr.prototype.fmt_G64 = function(v) {
		var f;
		f = this;
		f.formatFloat(v, 71, doPrec(f, -1), 64);
	};
	fmt.prototype.fmt_G64 = function(v) { return this.$val.fmt_G64(v); };
	fmt.Ptr.prototype.fmt_fb64 = function(v) {
		var f;
		f = this;
		f.formatFloat(v, 98, 0, 64);
	};
	fmt.prototype.fmt_fb64 = function(v) { return this.$val.fmt_fb64(v); };
	fmt.Ptr.prototype.fmt_e32 = function(v) {
		var f;
		f = this;
		f.formatFloat($coerceFloat32(v), 101, doPrec(f, 6), 32);
	};
	fmt.prototype.fmt_e32 = function(v) { return this.$val.fmt_e32(v); };
	fmt.Ptr.prototype.fmt_E32 = function(v) {
		var f;
		f = this;
		f.formatFloat($coerceFloat32(v), 69, doPrec(f, 6), 32);
	};
	fmt.prototype.fmt_E32 = function(v) { return this.$val.fmt_E32(v); };
	fmt.Ptr.prototype.fmt_f32 = function(v) {
		var f;
		f = this;
		f.formatFloat($coerceFloat32(v), 102, doPrec(f, 6), 32);
	};
	fmt.prototype.fmt_f32 = function(v) { return this.$val.fmt_f32(v); };
	fmt.Ptr.prototype.fmt_g32 = function(v) {
		var f;
		f = this;
		f.formatFloat($coerceFloat32(v), 103, doPrec(f, -1), 32);
	};
	fmt.prototype.fmt_g32 = function(v) { return this.$val.fmt_g32(v); };
	fmt.Ptr.prototype.fmt_G32 = function(v) {
		var f;
		f = this;
		f.formatFloat($coerceFloat32(v), 71, doPrec(f, -1), 32);
	};
	fmt.prototype.fmt_G32 = function(v) { return this.$val.fmt_G32(v); };
	fmt.Ptr.prototype.fmt_fb32 = function(v) {
		var f;
		f = this;
		f.formatFloat($coerceFloat32(v), 98, 0, 32);
	};
	fmt.prototype.fmt_fb32 = function(v) { return this.$val.fmt_fb32(v); };
	fmt.Ptr.prototype.fmt_c64 = function(v, verb) {
		var f;
		f = this;
		f.fmt_complex($coerceFloat32(v.$real), $coerceFloat32(v.$imag), 32, verb);
	};
	fmt.prototype.fmt_c64 = function(v, verb) { return this.$val.fmt_c64(v, verb); };
	fmt.Ptr.prototype.fmt_c128 = function(v, verb) {
		var f;
		f = this;
		f.fmt_complex(v.$real, v.$imag, 64, verb);
	};
	fmt.prototype.fmt_c128 = function(v, verb) { return this.$val.fmt_c128(v, verb); };
	fmt.Ptr.prototype.fmt_complex = function(r, j, size, verb) {
		var f, oldPlus, oldSpace, oldWid, i, _ref;
		f = this;
		f.buf.WriteByte(40);
		oldPlus = f.plus;
		oldSpace = f.space;
		oldWid = f.wid;
		i = 0;
		while (true) {
			_ref = verb;
			if (_ref === 98) {
				f.formatFloat(r, 98, 0, size);
			} else if (_ref === 101) {
				f.formatFloat(r, 101, doPrec(f, 6), size);
			} else if (_ref === 69) {
				f.formatFloat(r, 69, doPrec(f, 6), size);
			} else if (_ref === 102 || _ref === 70) {
				f.formatFloat(r, 102, doPrec(f, 6), size);
			} else if (_ref === 103) {
				f.formatFloat(r, 103, doPrec(f, -1), size);
			} else if (_ref === 71) {
				f.formatFloat(r, 71, doPrec(f, -1), size);
			}
			if (!((i === 0))) {
				break;
			}
			f.plus = true;
			f.space = false;
			f.wid = oldWid;
			r = j;
			i = i + (1) >> 0;
		}
		f.space = oldSpace;
		f.plus = oldPlus;
		f.wid = oldWid;
		f.buf.Write(irparenBytes);
	};
	fmt.prototype.fmt_complex = function(r, j, size, verb) { return this.$val.fmt_complex(r, j, size, verb); };
	$ptrType(buffer).prototype.Write = function(p) {
		var n = 0, err = null, b, _tmp, _tmp$1;
		b = this;
		b.$set($appendSlice(b.$get(), p));
		_tmp = p.$length; _tmp$1 = null; n = _tmp; err = _tmp$1;
		return [n, err];
	};
	$ptrType(buffer).prototype.WriteString = function(s) {
		var n = 0, err = null, b, _tmp, _tmp$1;
		b = this;
		b.$set($appendSlice(b.$get(), new buffer($stringToBytes(s))));
		_tmp = s.length; _tmp$1 = null; n = _tmp; err = _tmp$1;
		return [n, err];
	};
	$ptrType(buffer).prototype.WriteByte = function(c) {
		var b;
		b = this;
		b.$set($append(b.$get(), c));
		return null;
	};
	$ptrType(buffer).prototype.WriteRune = function(r) {
		var bp, b, n, x, w;
		bp = this;
		if (r < 128) {
			bp.$set($append(bp.$get(), (r << 24 >>> 24)));
			return null;
		}
		b = bp.$get();
		n = b.$length;
		while ((n + 4 >> 0) > b.$capacity) {
			b = $append(b, 0);
		}
		w = utf8.EncodeRune((x = $subslice(b, n, (n + 4 >> 0)), $subslice(new ($sliceType($Uint8))(x.$array), x.$offset, x.$offset + x.$length)), r);
		bp.$set($subslice(b, 0, (n + w >> 0)));
		return null;
	};
	newPrinter = function() {
		var x, p;
		p = (x = ppFree.Get(), (x !== null && x.constructor === ($ptrType(pp)) ? x.$val : $typeAssertionFailed(x, ($ptrType(pp)))));
		p.panicking = false;
		p.erroring = false;
		p.fmt.init(new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p));
		return p;
	};
	pp.Ptr.prototype.free = function() {
		var p;
		p = this;
		if (p.buf.$capacity > 1024) {
			return;
		}
		p.buf = $subslice(p.buf, 0, 0);
		p.arg = null;
		$copy(p.value, new reflect.Value.Ptr(($ptrType(reflect.rtype)).nil, 0, 0, 0), reflect.Value);
		ppFree.Put(p);
	};
	pp.prototype.free = function() { return this.$val.free(); };
	pp.Ptr.prototype.Width = function() {
		var wid = 0, ok = false, p, _tmp, _tmp$1;
		p = this;
		_tmp = p.fmt.wid; _tmp$1 = p.fmt.widPresent; wid = _tmp; ok = _tmp$1;
		return [wid, ok];
	};
	pp.prototype.Width = function() { return this.$val.Width(); };
	pp.Ptr.prototype.Precision = function() {
		var prec = 0, ok = false, p, _tmp, _tmp$1;
		p = this;
		_tmp = p.fmt.prec; _tmp$1 = p.fmt.precPresent; prec = _tmp; ok = _tmp$1;
		return [prec, ok];
	};
	pp.prototype.Precision = function() { return this.$val.Precision(); };
	pp.Ptr.prototype.Flag = function(b) {
		var p, _ref;
		p = this;
		_ref = b;
		if (_ref === 45) {
			return p.fmt.minus;
		} else if (_ref === 43) {
			return p.fmt.plus;
		} else if (_ref === 35) {
			return p.fmt.sharp;
		} else if (_ref === 32) {
			return p.fmt.space;
		} else if (_ref === 48) {
			return p.fmt.zero;
		}
		return false;
	};
	pp.prototype.Flag = function(b) { return this.$val.Flag(b); };
	pp.Ptr.prototype.add = function(c) {
		var p;
		p = this;
		new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteRune(c);
	};
	pp.prototype.add = function(c) { return this.$val.add(c); };
	pp.Ptr.prototype.Write = function(b) {
		var ret = 0, err = null, p, _tuple;
		p = this;
		_tuple = new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(b); ret = _tuple[0]; err = _tuple[1];
		return [ret, err];
	};
	pp.prototype.Write = function(b) { return this.$val.Write(b); };
	Fprintf = $pkg.Fprintf = function(w, format, a) {
		var n = 0, err = null, p, _tuple, x;
		p = newPrinter();
		p.doPrintf(format, a);
		_tuple = w.Write((x = p.buf, $subslice(new ($sliceType($Uint8))(x.$array), x.$offset, x.$offset + x.$length))); n = _tuple[0]; err = _tuple[1];
		p.free();
		return [n, err];
	};
	Printf = $pkg.Printf = function(format, a) {
		var n = 0, err = null, _tuple;
		_tuple = Fprintf(os.Stdout, format, a); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	Sprintf = $pkg.Sprintf = function(format, a) {
		var p, s;
		p = newPrinter();
		p.doPrintf(format, a);
		s = $bytesToString(p.buf);
		p.free();
		return s;
	};
	Fprint = $pkg.Fprint = function(w, a) {
		var n = 0, err = null, p, _tuple, x;
		p = newPrinter();
		p.doPrint(a, false, false);
		_tuple = w.Write((x = p.buf, $subslice(new ($sliceType($Uint8))(x.$array), x.$offset, x.$offset + x.$length))); n = _tuple[0]; err = _tuple[1];
		p.free();
		return [n, err];
	};
	Sprint = $pkg.Sprint = function(a) {
		var p, s;
		p = newPrinter();
		p.doPrint(a, false, false);
		s = $bytesToString(p.buf);
		p.free();
		return s;
	};
	Fprintln = $pkg.Fprintln = function(w, a) {
		var n = 0, err = null, p, _tuple, x;
		p = newPrinter();
		p.doPrint(a, true, true);
		_tuple = w.Write((x = p.buf, $subslice(new ($sliceType($Uint8))(x.$array), x.$offset, x.$offset + x.$length))); n = _tuple[0]; err = _tuple[1];
		p.free();
		return [n, err];
	};
	Println = $pkg.Println = function(a) {
		var n = 0, err = null, _tuple;
		_tuple = Fprintln(os.Stdout, a); n = _tuple[0]; err = _tuple[1];
		return [n, err];
	};
	getField = function(v, i) {
		var val;
		val = new reflect.Value.Ptr(); $copy(val, v.Field(i), reflect.Value);
		if ((val.Kind() === 20) && !val.IsNil()) {
			$copy(val, val.Elem(), reflect.Value);
		}
		return val;
	};
	parsenum = function(s, start, end) {
		var num = 0, isnum = false, newi = 0, _tmp, _tmp$1, _tmp$2;
		if (start >= end) {
			_tmp = 0; _tmp$1 = false; _tmp$2 = end; num = _tmp; isnum = _tmp$1; newi = _tmp$2;
			return [num, isnum, newi];
		}
		newi = start;
		while (newi < end && 48 <= s.charCodeAt(newi) && s.charCodeAt(newi) <= 57) {
			num = ((((num >>> 16 << 16) * 10 >> 0) + (num << 16 >>> 16) * 10) >> 0) + ((s.charCodeAt(newi) - 48 << 24 >>> 24) >> 0) >> 0;
			isnum = true;
			newi = newi + (1) >> 0;
		}
		return [num, isnum, newi];
	};
	pp.Ptr.prototype.unknownType = function(v) {
		var p;
		p = this;
		if ($interfaceIsEqual(v, null)) {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilAngleBytes);
			return;
		}
		new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(63);
		new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(reflect.TypeOf(v).String());
		new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(63);
	};
	pp.prototype.unknownType = function(v) { return this.$val.unknownType(v); };
	pp.Ptr.prototype.badVerb = function(verb) {
		var p;
		p = this;
		p.erroring = true;
		p.add(37);
		p.add(33);
		p.add(verb);
		p.add(40);
		if (!($interfaceIsEqual(p.arg, null))) {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(reflect.TypeOf(p.arg).String());
			p.add(61);
			p.printArg(p.arg, 118, false, false, 0);
		} else if (p.value.IsValid()) {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(p.value.Type().String());
			p.add(61);
			p.printValue($clone(p.value, reflect.Value), 118, false, false, 0);
		} else {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilAngleBytes);
		}
		p.add(41);
		p.erroring = false;
	};
	pp.prototype.badVerb = function(verb) { return this.$val.badVerb(verb); };
	pp.Ptr.prototype.fmtBool = function(v, verb) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 116 || _ref === 118) {
			p.fmt.fmt_boolean(v);
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtBool = function(v, verb) { return this.$val.fmtBool(v, verb); };
	pp.Ptr.prototype.fmtC = function(c) {
		var p, r, x, w;
		p = this;
		r = ((c.$low + ((c.$high >> 31) * 4294967296)) >> 0);
		if (!((x = new $Int64(0, r), (x.$high === c.$high && x.$low === c.$low)))) {
			r = 65533;
		}
		w = utf8.EncodeRune($subslice(new ($sliceType($Uint8))(p.runeBuf), 0, 4), r);
		p.fmt.pad($subslice(new ($sliceType($Uint8))(p.runeBuf), 0, w));
	};
	pp.prototype.fmtC = function(c) { return this.$val.fmtC(c); };
	pp.Ptr.prototype.fmtInt64 = function(v, verb) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 98) {
			p.fmt.integer(v, new $Uint64(0, 2), true, "0123456789abcdef");
		} else if (_ref === 99) {
			p.fmtC(v);
		} else if (_ref === 100 || _ref === 118) {
			p.fmt.integer(v, new $Uint64(0, 10), true, "0123456789abcdef");
		} else if (_ref === 111) {
			p.fmt.integer(v, new $Uint64(0, 8), true, "0123456789abcdef");
		} else if (_ref === 113) {
			if ((0 < v.$high || (0 === v.$high && 0 <= v.$low)) && (v.$high < 0 || (v.$high === 0 && v.$low <= 1114111))) {
				p.fmt.fmt_qc(v);
			} else {
				p.badVerb(verb);
			}
		} else if (_ref === 120) {
			p.fmt.integer(v, new $Uint64(0, 16), true, "0123456789abcdef");
		} else if (_ref === 85) {
			p.fmtUnicode(v);
		} else if (_ref === 88) {
			p.fmt.integer(v, new $Uint64(0, 16), true, "0123456789ABCDEF");
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtInt64 = function(v, verb) { return this.$val.fmtInt64(v, verb); };
	pp.Ptr.prototype.fmt0x64 = function(v, leading0x) {
		var p, sharp;
		p = this;
		sharp = p.fmt.sharp;
		p.fmt.sharp = leading0x;
		p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 16), false, "0123456789abcdef");
		p.fmt.sharp = sharp;
	};
	pp.prototype.fmt0x64 = function(v, leading0x) { return this.$val.fmt0x64(v, leading0x); };
	pp.Ptr.prototype.fmtUnicode = function(v) {
		var p, precPresent, sharp, prec;
		p = this;
		precPresent = p.fmt.precPresent;
		sharp = p.fmt.sharp;
		p.fmt.sharp = false;
		prec = p.fmt.prec;
		if (!precPresent) {
			p.fmt.prec = 4;
			p.fmt.precPresent = true;
		}
		p.fmt.unicode = true;
		p.fmt.uniQuote = sharp;
		p.fmt.integer(v, new $Uint64(0, 16), false, "0123456789ABCDEF");
		p.fmt.unicode = false;
		p.fmt.uniQuote = false;
		p.fmt.prec = prec;
		p.fmt.precPresent = precPresent;
		p.fmt.sharp = sharp;
	};
	pp.prototype.fmtUnicode = function(v) { return this.$val.fmtUnicode(v); };
	pp.Ptr.prototype.fmtUint64 = function(v, verb, goSyntax) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 98) {
			p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 2), false, "0123456789abcdef");
		} else if (_ref === 99) {
			p.fmtC(new $Int64(v.$high, v.$low));
		} else if (_ref === 100) {
			p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 10), false, "0123456789abcdef");
		} else if (_ref === 118) {
			if (goSyntax) {
				p.fmt0x64(v, true);
			} else {
				p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 10), false, "0123456789abcdef");
			}
		} else if (_ref === 111) {
			p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 8), false, "0123456789abcdef");
		} else if (_ref === 113) {
			if ((0 < v.$high || (0 === v.$high && 0 <= v.$low)) && (v.$high < 0 || (v.$high === 0 && v.$low <= 1114111))) {
				p.fmt.fmt_qc(new $Int64(v.$high, v.$low));
			} else {
				p.badVerb(verb);
			}
		} else if (_ref === 120) {
			p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 16), false, "0123456789abcdef");
		} else if (_ref === 88) {
			p.fmt.integer(new $Int64(v.$high, v.$low), new $Uint64(0, 16), false, "0123456789ABCDEF");
		} else if (_ref === 85) {
			p.fmtUnicode(new $Int64(v.$high, v.$low));
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtUint64 = function(v, verb, goSyntax) { return this.$val.fmtUint64(v, verb, goSyntax); };
	pp.Ptr.prototype.fmtFloat32 = function(v, verb) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 98) {
			p.fmt.fmt_fb32(v);
		} else if (_ref === 101) {
			p.fmt.fmt_e32(v);
		} else if (_ref === 69) {
			p.fmt.fmt_E32(v);
		} else if (_ref === 102 || _ref === 70) {
			p.fmt.fmt_f32(v);
		} else if (_ref === 103 || _ref === 118) {
			p.fmt.fmt_g32(v);
		} else if (_ref === 71) {
			p.fmt.fmt_G32(v);
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtFloat32 = function(v, verb) { return this.$val.fmtFloat32(v, verb); };
	pp.Ptr.prototype.fmtFloat64 = function(v, verb) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 98) {
			p.fmt.fmt_fb64(v);
		} else if (_ref === 101) {
			p.fmt.fmt_e64(v);
		} else if (_ref === 69) {
			p.fmt.fmt_E64(v);
		} else if (_ref === 102 || _ref === 70) {
			p.fmt.fmt_f64(v);
		} else if (_ref === 103 || _ref === 118) {
			p.fmt.fmt_g64(v);
		} else if (_ref === 71) {
			p.fmt.fmt_G64(v);
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtFloat64 = function(v, verb) { return this.$val.fmtFloat64(v, verb); };
	pp.Ptr.prototype.fmtComplex64 = function(v, verb) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 98 || _ref === 101 || _ref === 69 || _ref === 102 || _ref === 70 || _ref === 103 || _ref === 71) {
			p.fmt.fmt_c64(v, verb);
		} else if (_ref === 118) {
			p.fmt.fmt_c64(v, 103);
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtComplex64 = function(v, verb) { return this.$val.fmtComplex64(v, verb); };
	pp.Ptr.prototype.fmtComplex128 = function(v, verb) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 98 || _ref === 101 || _ref === 69 || _ref === 102 || _ref === 70 || _ref === 103 || _ref === 71) {
			p.fmt.fmt_c128(v, verb);
		} else if (_ref === 118) {
			p.fmt.fmt_c128(v, 103);
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtComplex128 = function(v, verb) { return this.$val.fmtComplex128(v, verb); };
	pp.Ptr.prototype.fmtString = function(v, verb, goSyntax) {
		var p, _ref;
		p = this;
		_ref = verb;
		if (_ref === 118) {
			if (goSyntax) {
				p.fmt.fmt_q(v);
			} else {
				p.fmt.fmt_s(v);
			}
		} else if (_ref === 115) {
			p.fmt.fmt_s(v);
		} else if (_ref === 120) {
			p.fmt.fmt_sx(v, "0123456789abcdef");
		} else if (_ref === 88) {
			p.fmt.fmt_sx(v, "0123456789ABCDEF");
		} else if (_ref === 113) {
			p.fmt.fmt_q(v);
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtString = function(v, verb, goSyntax) { return this.$val.fmtString(v, verb, goSyntax); };
	pp.Ptr.prototype.fmtBytes = function(v, verb, goSyntax, typ, depth) {
		var p, _ref, _i, i, c, _ref$1;
		p = this;
		if ((verb === 118) || (verb === 100)) {
			if (goSyntax) {
				if (v === ($sliceType($Uint8)).nil) {
					if ($interfaceIsEqual(typ, null)) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString("[]byte(nil)");
					} else {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(typ.String());
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilParenBytes);
					}
					return;
				}
				if ($interfaceIsEqual(typ, null)) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(bytesBytes);
				} else {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(typ.String());
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(123);
				}
			} else {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(91);
			}
			_ref = v;
			_i = 0;
			while (_i < _ref.$length) {
				i = _i;
				c = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
				if (i > 0) {
					if (goSyntax) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(commaSpaceBytes);
					} else {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(32);
					}
				}
				p.printArg(new $Uint8(c), 118, p.fmt.plus, goSyntax, depth + 1 >> 0);
				_i++;
			}
			if (goSyntax) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(125);
			} else {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(93);
			}
			return;
		}
		_ref$1 = verb;
		if (_ref$1 === 115) {
			p.fmt.fmt_s($bytesToString(v));
		} else if (_ref$1 === 120) {
			p.fmt.fmt_bx(v, "0123456789abcdef");
		} else if (_ref$1 === 88) {
			p.fmt.fmt_bx(v, "0123456789ABCDEF");
		} else if (_ref$1 === 113) {
			p.fmt.fmt_q($bytesToString(v));
		} else {
			p.badVerb(verb);
		}
	};
	pp.prototype.fmtBytes = function(v, verb, goSyntax, typ, depth) { return this.$val.fmtBytes(v, verb, goSyntax, typ, depth); };
	pp.Ptr.prototype.fmtPointer = function(value, verb, goSyntax) {
		var p, use0x64, _ref, u, _ref$1;
		p = this;
		use0x64 = true;
		_ref = verb;
		if (_ref === 112 || _ref === 118) {
		} else if (_ref === 98 || _ref === 100 || _ref === 111 || _ref === 120 || _ref === 88) {
			use0x64 = false;
		} else {
			p.badVerb(verb);
			return;
		}
		u = 0;
		_ref$1 = value.Kind();
		if (_ref$1 === 18 || _ref$1 === 19 || _ref$1 === 21 || _ref$1 === 22 || _ref$1 === 23 || _ref$1 === 26) {
			u = value.Pointer();
		} else {
			p.badVerb(verb);
			return;
		}
		if (goSyntax) {
			p.add(40);
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(value.Type().String());
			p.add(41);
			p.add(40);
			if (u === 0) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilBytes);
			} else {
				p.fmt0x64(new $Uint64(0, u.constructor === Number ? u : 1), true);
			}
			p.add(41);
		} else if ((verb === 118) && (u === 0)) {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilAngleBytes);
		} else {
			if (use0x64) {
				p.fmt0x64(new $Uint64(0, u.constructor === Number ? u : 1), !p.fmt.sharp);
			} else {
				p.fmtUint64(new $Uint64(0, u.constructor === Number ? u : 1), verb, false);
			}
		}
	};
	pp.prototype.fmtPointer = function(value, verb, goSyntax) { return this.$val.fmtPointer(value, verb, goSyntax); };
	pp.Ptr.prototype.catchPanic = function(arg, verb) {
		var p, err, v;
		p = this;
		err = $recover();
		if (!($interfaceIsEqual(err, null))) {
			v = new reflect.Value.Ptr(); $copy(v, reflect.ValueOf(arg), reflect.Value);
			if ((v.Kind() === 22) && v.IsNil()) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilAngleBytes);
				return;
			}
			if (p.panicking) {
				$panic(err);
			}
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(percentBangBytes);
			p.add(verb);
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(panicBytes);
			p.panicking = true;
			p.printArg(err, 118, false, false, 0);
			p.panicking = false;
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(41);
		}
	};
	pp.prototype.catchPanic = function(arg, verb) { return this.$val.catchPanic(arg, verb); };
	pp.Ptr.prototype.handleMethods = function(verb, plus, goSyntax, depth) {
		var wasString = false, handled = false, $deferred = [], $err = null, p, _tuple, x, formatter, ok, arg, verb$1, _recv, _tuple$1, x$1, stringer, ok$1, arg$1, verb$2, _recv$1, _ref, v, _ref$1, _type, arg$2, verb$3, _recv$2, arg$3, verb$4, _recv$3;
		/* */ try { $deferFrames.push($deferred);
		p = this;
		if (p.erroring) {
			return [wasString, handled];
		}
		_tuple = (x = p.arg, (x !== null && Formatter.implementedBy.indexOf(x.constructor) !== -1 ? [x, true] : [null, false])); formatter = _tuple[0]; ok = _tuple[1];
		if (ok) {
			handled = true;
			wasString = false;
			$deferred.push([(_recv = p, function(arg, verb$1) { $stackDepthOffset--; try { return _recv.catchPanic(arg, verb$1); } finally { $stackDepthOffset++; } }), [p.arg, verb]]);
			formatter.Format(p, verb);
			return [wasString, handled];
		}
		if (plus) {
			p.fmt.plus = false;
		}
		if (goSyntax) {
			p.fmt.sharp = false;
			_tuple$1 = (x$1 = p.arg, (x$1 !== null && GoStringer.implementedBy.indexOf(x$1.constructor) !== -1 ? [x$1, true] : [null, false])); stringer = _tuple$1[0]; ok$1 = _tuple$1[1];
			if (ok$1) {
				wasString = false;
				handled = true;
				$deferred.push([(_recv$1 = p, function(arg$1, verb$2) { $stackDepthOffset--; try { return _recv$1.catchPanic(arg$1, verb$2); } finally { $stackDepthOffset++; } }), [p.arg, verb]]);
				p.fmtString(stringer.GoString(), 115, false);
				return [wasString, handled];
			}
		} else {
			_ref = verb;
			if (_ref === 118 || _ref === 115 || _ref === 120 || _ref === 88 || _ref === 113) {
				_ref$1 = p.arg;
				_type = _ref$1 !== null ? _ref$1.constructor : null;
				if ($error.implementedBy.indexOf(_type) !== -1) {
					v = _ref$1;
					wasString = false;
					handled = true;
					$deferred.push([(_recv$2 = p, function(arg$2, verb$3) { $stackDepthOffset--; try { return _recv$2.catchPanic(arg$2, verb$3); } finally { $stackDepthOffset++; } }), [p.arg, verb]]);
					p.printArg(new $String(v.Error()), verb, plus, false, depth);
					return [wasString, handled];
				} else if (Stringer.implementedBy.indexOf(_type) !== -1) {
					v = _ref$1;
					wasString = false;
					handled = true;
					$deferred.push([(_recv$3 = p, function(arg$3, verb$4) { $stackDepthOffset--; try { return _recv$3.catchPanic(arg$3, verb$4); } finally { $stackDepthOffset++; } }), [p.arg, verb]]);
					p.printArg(new $String(v.String()), verb, plus, false, depth);
					return [wasString, handled];
				}
			}
		}
		handled = false;
		return [wasString, handled];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [wasString, handled]; }
	};
	pp.prototype.handleMethods = function(verb, plus, goSyntax, depth) { return this.$val.handleMethods(verb, plus, goSyntax, depth); };
	pp.Ptr.prototype.printArg = function(arg, verb, plus, goSyntax, depth) {
		var wasString = false, p, _ref, oldPlus, oldSharp, f, _ref$1, _type, _tuple, isString, handled;
		p = this;
		p.arg = arg;
		$copy(p.value, new reflect.Value.Ptr(($ptrType(reflect.rtype)).nil, 0, 0, 0), reflect.Value);
		if ($interfaceIsEqual(arg, null)) {
			if ((verb === 84) || (verb === 118)) {
				p.fmt.pad(nilAngleBytes);
			} else {
				p.badVerb(verb);
			}
			wasString = false;
			return wasString;
		}
		_ref = verb;
		if (_ref === 84) {
			p.printArg(new $String(reflect.TypeOf(arg).String()), 115, false, false, 0);
			wasString = false;
			return wasString;
		} else if (_ref === 112) {
			p.fmtPointer($clone(reflect.ValueOf(arg), reflect.Value), verb, goSyntax);
			wasString = false;
			return wasString;
		}
		oldPlus = p.fmt.plus;
		oldSharp = p.fmt.sharp;
		if (plus) {
			p.fmt.plus = false;
		}
		if (goSyntax) {
			p.fmt.sharp = false;
		}
		_ref$1 = arg;
		_type = _ref$1 !== null ? _ref$1.constructor : null;
		if (_type === $Bool) {
			f = _ref$1.$val;
			p.fmtBool(f, verb);
		} else if (_type === $Float32) {
			f = _ref$1.$val;
			p.fmtFloat32(f, verb);
		} else if (_type === $Float64) {
			f = _ref$1.$val;
			p.fmtFloat64(f, verb);
		} else if (_type === $Complex64) {
			f = _ref$1.$val;
			p.fmtComplex64(f, verb);
		} else if (_type === $Complex128) {
			f = _ref$1.$val;
			p.fmtComplex128(f, verb);
		} else if (_type === $Int) {
			f = _ref$1.$val;
			p.fmtInt64(new $Int64(0, f), verb);
		} else if (_type === $Int8) {
			f = _ref$1.$val;
			p.fmtInt64(new $Int64(0, f), verb);
		} else if (_type === $Int16) {
			f = _ref$1.$val;
			p.fmtInt64(new $Int64(0, f), verb);
		} else if (_type === $Int32) {
			f = _ref$1.$val;
			p.fmtInt64(new $Int64(0, f), verb);
		} else if (_type === $Int64) {
			f = _ref$1.$val;
			p.fmtInt64(f, verb);
		} else if (_type === $Uint) {
			f = _ref$1.$val;
			p.fmtUint64(new $Uint64(0, f), verb, goSyntax);
		} else if (_type === $Uint8) {
			f = _ref$1.$val;
			p.fmtUint64(new $Uint64(0, f), verb, goSyntax);
		} else if (_type === $Uint16) {
			f = _ref$1.$val;
			p.fmtUint64(new $Uint64(0, f), verb, goSyntax);
		} else if (_type === $Uint32) {
			f = _ref$1.$val;
			p.fmtUint64(new $Uint64(0, f), verb, goSyntax);
		} else if (_type === $Uint64) {
			f = _ref$1.$val;
			p.fmtUint64(f, verb, goSyntax);
		} else if (_type === $Uintptr) {
			f = _ref$1.$val;
			p.fmtUint64(new $Uint64(0, f.constructor === Number ? f : 1), verb, goSyntax);
		} else if (_type === $String) {
			f = _ref$1.$val;
			p.fmtString(f, verb, goSyntax);
			wasString = (verb === 115) || (verb === 118);
		} else if (_type === ($sliceType($Uint8))) {
			f = _ref$1.$val;
			p.fmtBytes(f, verb, goSyntax, null, depth);
			wasString = verb === 115;
		} else {
			f = _ref$1;
			p.fmt.plus = oldPlus;
			p.fmt.sharp = oldSharp;
			_tuple = p.handleMethods(verb, plus, goSyntax, depth); isString = _tuple[0]; handled = _tuple[1];
			if (handled) {
				wasString = isString;
				return wasString;
			}
			wasString = p.printReflectValue($clone(reflect.ValueOf(arg), reflect.Value), verb, plus, goSyntax, depth);
			return wasString;
		}
		p.arg = null;
		return wasString;
	};
	pp.prototype.printArg = function(arg, verb, plus, goSyntax, depth) { return this.$val.printArg(arg, verb, plus, goSyntax, depth); };
	pp.Ptr.prototype.printValue = function(value, verb, plus, goSyntax, depth) {
		var wasString = false, p, _ref, _tuple, isString, handled;
		p = this;
		if (!value.IsValid()) {
			if ((verb === 84) || (verb === 118)) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilAngleBytes);
			} else {
				p.badVerb(verb);
			}
			wasString = false;
			return wasString;
		}
		_ref = verb;
		if (_ref === 84) {
			p.printArg(new $String(value.Type().String()), 115, false, false, 0);
			wasString = false;
			return wasString;
		} else if (_ref === 112) {
			p.fmtPointer($clone(value, reflect.Value), verb, goSyntax);
			wasString = false;
			return wasString;
		}
		p.arg = null;
		if (value.CanInterface()) {
			p.arg = value.Interface();
		}
		_tuple = p.handleMethods(verb, plus, goSyntax, depth); isString = _tuple[0]; handled = _tuple[1];
		if (handled) {
			wasString = isString;
			return wasString;
		}
		wasString = p.printReflectValue($clone(value, reflect.Value), verb, plus, goSyntax, depth);
		return wasString;
	};
	pp.prototype.printValue = function(value, verb, plus, goSyntax, depth) { return this.$val.printValue(value, verb, plus, goSyntax, depth); };
	pp.Ptr.prototype.printReflectValue = function(value, verb, plus, goSyntax, depth) {
		var wasString = false, p, oldValue, f, _ref, x, keys, _ref$1, _i, i, key, v, t, i$1, f$1, value$1, typ, bytes, _ref$2, _i$1, i$2, i$3, v$1, a, _ref$3;
		p = this;
		oldValue = new reflect.Value.Ptr(); $copy(oldValue, p.value, reflect.Value);
		$copy(p.value, value, reflect.Value);
		f = new reflect.Value.Ptr(); $copy(f, value, reflect.Value);
		_ref = f.Kind();
		BigSwitch:
		switch (0) { default: if (_ref === 1) {
			p.fmtBool(f.Bool(), verb);
		} else if (_ref === 2 || _ref === 3 || _ref === 4 || _ref === 5 || _ref === 6) {
			p.fmtInt64(f.Int(), verb);
		} else if (_ref === 7 || _ref === 8 || _ref === 9 || _ref === 10 || _ref === 11 || _ref === 12) {
			p.fmtUint64(f.Uint(), verb, goSyntax);
		} else if (_ref === 13 || _ref === 14) {
			if (f.Type().Size() === 4) {
				p.fmtFloat32(f.Float(), verb);
			} else {
				p.fmtFloat64(f.Float(), verb);
			}
		} else if (_ref === 15 || _ref === 16) {
			if (f.Type().Size() === 8) {
				p.fmtComplex64((x = f.Complex(), new $Complex64(x.$real, x.$imag)), verb);
			} else {
				p.fmtComplex128(f.Complex(), verb);
			}
		} else if (_ref === 24) {
			p.fmtString(f.String(), verb, goSyntax);
		} else if (_ref === 21) {
			if (goSyntax) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(f.Type().String());
				if (f.IsNil()) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString("(nil)");
					break;
				}
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(123);
			} else {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(mapBytes);
			}
			keys = f.MapKeys();
			_ref$1 = keys;
			_i = 0;
			while (_i < _ref$1.$length) {
				i = _i;
				key = new reflect.Value.Ptr(); $copy(key, ((_i < 0 || _i >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i]), reflect.Value);
				if (i > 0) {
					if (goSyntax) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(commaSpaceBytes);
					} else {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(32);
					}
				}
				p.printValue($clone(key, reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(58);
				p.printValue($clone(f.MapIndex($clone(key, reflect.Value)), reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
				_i++;
			}
			if (goSyntax) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(125);
			} else {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(93);
			}
		} else if (_ref === 25) {
			if (goSyntax) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(value.Type().String());
			}
			p.add(123);
			v = new reflect.Value.Ptr(); $copy(v, f, reflect.Value);
			t = v.Type();
			i$1 = 0;
			while (i$1 < v.NumField()) {
				if (i$1 > 0) {
					if (goSyntax) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(commaSpaceBytes);
					} else {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(32);
					}
				}
				if (plus || goSyntax) {
					f$1 = new reflect.StructField.Ptr(); $copy(f$1, t.Field(i$1), reflect.StructField);
					if (!(f$1.Name === "")) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(f$1.Name);
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(58);
					}
				}
				p.printValue($clone(getField($clone(v, reflect.Value), i$1), reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
				i$1 = i$1 + (1) >> 0;
			}
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(125);
		} else if (_ref === 20) {
			value$1 = new reflect.Value.Ptr(); $copy(value$1, f.Elem(), reflect.Value);
			if (!value$1.IsValid()) {
				if (goSyntax) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(f.Type().String());
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilParenBytes);
				} else {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(nilAngleBytes);
				}
			} else {
				wasString = p.printValue($clone(value$1, reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
			}
		} else if (_ref === 17 || _ref === 23) {
			typ = f.Type();
			if (typ.Elem().Kind() === 8) {
				bytes = ($sliceType($Uint8)).nil;
				if (f.Kind() === 23) {
					bytes = f.Bytes();
				} else if (f.CanAddr()) {
					bytes = f.Slice(0, f.Len()).Bytes();
				} else {
					bytes = ($sliceType($Uint8)).make(f.Len());
					_ref$2 = bytes;
					_i$1 = 0;
					while (_i$1 < _ref$2.$length) {
						i$2 = _i$1;
						(i$2 < 0 || i$2 >= bytes.$length) ? $throwRuntimeError("index out of range") : bytes.$array[bytes.$offset + i$2] = (f.Index(i$2).Uint().$low << 24 >>> 24);
						_i$1++;
					}
				}
				p.fmtBytes(bytes, verb, goSyntax, typ, depth);
				wasString = verb === 115;
				break;
			}
			if (goSyntax) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(value.Type().String());
				if ((f.Kind() === 23) && f.IsNil()) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString("(nil)");
					break;
				}
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(123);
			} else {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(91);
			}
			i$3 = 0;
			while (i$3 < f.Len()) {
				if (i$3 > 0) {
					if (goSyntax) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(commaSpaceBytes);
					} else {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(32);
					}
				}
				p.printValue($clone(f.Index(i$3), reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
				i$3 = i$3 + (1) >> 0;
			}
			if (goSyntax) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(125);
			} else {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(93);
			}
		} else if (_ref === 22) {
			v$1 = f.Pointer();
			if (!((v$1 === 0)) && (depth === 0)) {
				a = new reflect.Value.Ptr(); $copy(a, f.Elem(), reflect.Value);
				_ref$3 = a.Kind();
				if (_ref$3 === 17 || _ref$3 === 23) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(38);
					p.printValue($clone(a, reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
					break BigSwitch;
				} else if (_ref$3 === 25) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(38);
					p.printValue($clone(a, reflect.Value), verb, plus, goSyntax, depth + 1 >> 0);
					break BigSwitch;
				}
			}
			p.fmtPointer($clone(value, reflect.Value), verb, goSyntax);
		} else if (_ref === 18 || _ref === 19 || _ref === 26) {
			p.fmtPointer($clone(value, reflect.Value), verb, goSyntax);
		} else {
			p.unknownType(new f.constructor.Struct(f));
		} }
		$copy(p.value, oldValue, reflect.Value);
		wasString = wasString;
		return wasString;
	};
	pp.prototype.printReflectValue = function(value, verb, plus, goSyntax, depth) { return this.$val.printReflectValue(value, verb, plus, goSyntax, depth); };
	intFromArg = function(a, argNum) {
		var num = 0, isInt = false, newArgNum = 0, _tuple, x;
		newArgNum = argNum;
		if (argNum < a.$length) {
			_tuple = (x = ((argNum < 0 || argNum >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + argNum]), (x !== null && x.constructor === $Int ? [x.$val, true] : [0, false])); num = _tuple[0]; isInt = _tuple[1];
			newArgNum = argNum + 1 >> 0;
		}
		return [num, isInt, newArgNum];
	};
	parseArgNumber = function(format) {
		var index = 0, wid = 0, ok = false, i, _tuple, width, ok$1, newi, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8;
		i = 1;
		while (i < format.length) {
			if (format.charCodeAt(i) === 93) {
				_tuple = parsenum(format, 1, i); width = _tuple[0]; ok$1 = _tuple[1]; newi = _tuple[2];
				if (!ok$1 || !((newi === i))) {
					_tmp = 0; _tmp$1 = i + 1 >> 0; _tmp$2 = false; index = _tmp; wid = _tmp$1; ok = _tmp$2;
					return [index, wid, ok];
				}
				_tmp$3 = width - 1 >> 0; _tmp$4 = i + 1 >> 0; _tmp$5 = true; index = _tmp$3; wid = _tmp$4; ok = _tmp$5;
				return [index, wid, ok];
			}
			i = i + (1) >> 0;
		}
		_tmp$6 = 0; _tmp$7 = 1; _tmp$8 = false; index = _tmp$6; wid = _tmp$7; ok = _tmp$8;
		return [index, wid, ok];
	};
	pp.Ptr.prototype.argNumber = function(argNum, format, i, numArgs) {
		var newArgNum = 0, newi = 0, found = false, p, _tmp, _tmp$1, _tmp$2, _tuple, index, wid, ok, _tmp$3, _tmp$4, _tmp$5, _tmp$6, _tmp$7, _tmp$8;
		p = this;
		if (format.length <= i || !((format.charCodeAt(i) === 91))) {
			_tmp = argNum; _tmp$1 = i; _tmp$2 = false; newArgNum = _tmp; newi = _tmp$1; found = _tmp$2;
			return [newArgNum, newi, found];
		}
		p.reordered = true;
		_tuple = parseArgNumber(format.substring(i)); index = _tuple[0]; wid = _tuple[1]; ok = _tuple[2];
		if (ok && 0 <= index && index < numArgs) {
			_tmp$3 = index; _tmp$4 = i + wid >> 0; _tmp$5 = true; newArgNum = _tmp$3; newi = _tmp$4; found = _tmp$5;
			return [newArgNum, newi, found];
		}
		p.goodArgNum = false;
		_tmp$6 = argNum; _tmp$7 = i + wid >> 0; _tmp$8 = true; newArgNum = _tmp$6; newi = _tmp$7; found = _tmp$8;
		return [newArgNum, newi, found];
	};
	pp.prototype.argNumber = function(argNum, format, i, numArgs) { return this.$val.argNumber(argNum, format, i, numArgs); };
	pp.Ptr.prototype.doPrintf = function(format, a) {
		var p, end, argNum, afterIndex, i, lasti, _ref, _tuple, _tuple$1, _tuple$2, _tuple$3, _tuple$4, _tuple$5, _tuple$6, _tuple$7, c, w, arg, goSyntax, plus, arg$1;
		p = this;
		end = format.length;
		argNum = 0;
		afterIndex = false;
		p.reordered = false;
		i = 0;
		while (i < end) {
			p.goodArgNum = true;
			lasti = i;
			while (i < end && !((format.charCodeAt(i) === 37))) {
				i = i + (1) >> 0;
			}
			if (i > lasti) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(format.substring(lasti, i));
			}
			if (i >= end) {
				break;
			}
			i = i + (1) >> 0;
			p.fmt.clearflags();
			F:
			while (i < end) {
				_ref = format.charCodeAt(i);
				if (_ref === 35) {
					p.fmt.sharp = true;
				} else if (_ref === 48) {
					p.fmt.zero = true;
				} else if (_ref === 43) {
					p.fmt.plus = true;
				} else if (_ref === 45) {
					p.fmt.minus = true;
				} else if (_ref === 32) {
					p.fmt.space = true;
				} else {
					break F;
				}
				i = i + (1) >> 0;
			}
			_tuple = p.argNumber(argNum, format, i, a.$length); argNum = _tuple[0]; i = _tuple[1]; afterIndex = _tuple[2];
			if (i < end && (format.charCodeAt(i) === 42)) {
				i = i + (1) >> 0;
				_tuple$1 = intFromArg(a, argNum); p.fmt.wid = _tuple$1[0]; p.fmt.widPresent = _tuple$1[1]; argNum = _tuple$1[2];
				if (!p.fmt.widPresent) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(badWidthBytes);
				}
				afterIndex = false;
			} else {
				_tuple$2 = parsenum(format, i, end); p.fmt.wid = _tuple$2[0]; p.fmt.widPresent = _tuple$2[1]; i = _tuple$2[2];
				if (afterIndex && p.fmt.widPresent) {
					p.goodArgNum = false;
				}
			}
			if ((i + 1 >> 0) < end && (format.charCodeAt(i) === 46)) {
				i = i + (1) >> 0;
				if (afterIndex) {
					p.goodArgNum = false;
				}
				_tuple$3 = p.argNumber(argNum, format, i, a.$length); argNum = _tuple$3[0]; i = _tuple$3[1]; afterIndex = _tuple$3[2];
				if (format.charCodeAt(i) === 42) {
					i = i + (1) >> 0;
					_tuple$4 = intFromArg(a, argNum); p.fmt.prec = _tuple$4[0]; p.fmt.precPresent = _tuple$4[1]; argNum = _tuple$4[2];
					if (!p.fmt.precPresent) {
						new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(badPrecBytes);
					}
					afterIndex = false;
				} else {
					_tuple$5 = parsenum(format, i, end); p.fmt.prec = _tuple$5[0]; p.fmt.precPresent = _tuple$5[1]; i = _tuple$5[2];
					if (!p.fmt.precPresent) {
						p.fmt.prec = 0;
						p.fmt.precPresent = true;
					}
				}
			}
			if (!afterIndex) {
				_tuple$6 = p.argNumber(argNum, format, i, a.$length); argNum = _tuple$6[0]; i = _tuple$6[1]; afterIndex = _tuple$6[2];
			}
			if (i >= end) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(noVerbBytes);
				continue;
			}
			_tuple$7 = utf8.DecodeRuneInString(format.substring(i)); c = _tuple$7[0]; w = _tuple$7[1];
			i = i + (w) >> 0;
			if (c === 37) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(37);
				continue;
			}
			if (!p.goodArgNum) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(percentBangBytes);
				p.add(c);
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(badIndexBytes);
				continue;
			} else if (argNum >= a.$length) {
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(percentBangBytes);
				p.add(c);
				new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(missingBytes);
				continue;
			}
			arg = ((argNum < 0 || argNum >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + argNum]);
			argNum = argNum + (1) >> 0;
			goSyntax = (c === 118) && p.fmt.sharp;
			plus = (c === 118) && p.fmt.plus;
			p.printArg(arg, c, plus, goSyntax, 0);
		}
		if (!p.reordered && argNum < a.$length) {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(extraBytes);
			while (argNum < a.$length) {
				arg$1 = ((argNum < 0 || argNum >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + argNum]);
				if (!($interfaceIsEqual(arg$1, null))) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteString(reflect.TypeOf(arg$1).String());
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(61);
				}
				p.printArg(arg$1, 118, false, false, 0);
				if ((argNum + 1 >> 0) < a.$length) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).Write(commaSpaceBytes);
				}
				argNum = argNum + (1) >> 0;
			}
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(41);
		}
	};
	pp.prototype.doPrintf = function(format, a) { return this.$val.doPrintf(format, a); };
	pp.Ptr.prototype.doPrint = function(a, addspace, addnewline) {
		var p, prevString, argNum, arg, isString;
		p = this;
		prevString = false;
		argNum = 0;
		while (argNum < a.$length) {
			p.fmt.clearflags();
			arg = ((argNum < 0 || argNum >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + argNum]);
			if (argNum > 0) {
				isString = !($interfaceIsEqual(arg, null)) && (reflect.TypeOf(arg).Kind() === 24);
				if (addspace || !isString && !prevString) {
					new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(32);
				}
			}
			prevString = p.printArg(arg, 118, false, false, 0);
			argNum = argNum + (1) >> 0;
		}
		if (addnewline) {
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, p).WriteByte(10);
		}
	};
	pp.prototype.doPrint = function(a, addspace, addnewline) { return this.$val.doPrint(a, addspace, addnewline); };
	ss.Ptr.prototype.Read = function(buf) {
		var n = 0, err = null, s, _tmp, _tmp$1;
		s = this;
		_tmp = 0; _tmp$1 = errors.New("ScanState's Read should not be called. Use ReadRune"); n = _tmp; err = _tmp$1;
		return [n, err];
	};
	ss.prototype.Read = function(buf) { return this.$val.Read(buf); };
	ss.Ptr.prototype.ReadRune = function() {
		var r = 0, size = 0, err = null, s, _tuple;
		s = this;
		if (s.peekRune >= 0) {
			s.count = s.count + (1) >> 0;
			r = s.peekRune;
			size = utf8.RuneLen(r);
			s.prevRune = r;
			s.peekRune = -1;
			return [r, size, err];
		}
		if (s.atEOF || s.ssave.nlIsEnd && (s.prevRune === 10) || s.count >= s.ssave.argLimit) {
			err = io.EOF;
			return [r, size, err];
		}
		_tuple = s.rr.ReadRune(); r = _tuple[0]; size = _tuple[1]; err = _tuple[2];
		if ($interfaceIsEqual(err, null)) {
			s.count = s.count + (1) >> 0;
			s.prevRune = r;
		} else if ($interfaceIsEqual(err, io.EOF)) {
			s.atEOF = true;
		}
		return [r, size, err];
	};
	ss.prototype.ReadRune = function() { return this.$val.ReadRune(); };
	ss.Ptr.prototype.Width = function() {
		var wid = 0, ok = false, s, _tmp, _tmp$1, _tmp$2, _tmp$3;
		s = this;
		if (s.ssave.maxWid === 1073741824) {
			_tmp = 0; _tmp$1 = false; wid = _tmp; ok = _tmp$1;
			return [wid, ok];
		}
		_tmp$2 = s.ssave.maxWid; _tmp$3 = true; wid = _tmp$2; ok = _tmp$3;
		return [wid, ok];
	};
	ss.prototype.Width = function() { return this.$val.Width(); };
	ss.Ptr.prototype.getRune = function() {
		var r = 0, s, _tuple, err;
		s = this;
		_tuple = s.ReadRune(); r = _tuple[0]; err = _tuple[2];
		if (!($interfaceIsEqual(err, null))) {
			if ($interfaceIsEqual(err, io.EOF)) {
				r = -1;
				return r;
			}
			s.error(err);
		}
		return r;
	};
	ss.prototype.getRune = function() { return this.$val.getRune(); };
	ss.Ptr.prototype.UnreadRune = function() {
		var s, _tuple, x, u, ok;
		s = this;
		_tuple = (x = s.rr, (x !== null && runeUnreader.implementedBy.indexOf(x.constructor) !== -1 ? [x, true] : [null, false])); u = _tuple[0]; ok = _tuple[1];
		if (ok) {
			u.UnreadRune();
		} else {
			s.peekRune = s.prevRune;
		}
		s.prevRune = -1;
		s.count = s.count - (1) >> 0;
		return null;
	};
	ss.prototype.UnreadRune = function() { return this.$val.UnreadRune(); };
	ss.Ptr.prototype.error = function(err) {
		var s, x;
		s = this;
		$panic((x = new scanError.Ptr(err), new x.constructor.Struct(x)));
	};
	ss.prototype.error = function(err) { return this.$val.error(err); };
	ss.Ptr.prototype.errorString = function(err) {
		var s, x;
		s = this;
		$panic((x = new scanError.Ptr(errors.New(err)), new x.constructor.Struct(x)));
	};
	ss.prototype.errorString = function(err) { return this.$val.errorString(err); };
	ss.Ptr.prototype.Token = function(skipSpace, f) {
		var tok = ($sliceType($Uint8)).nil, err = null, $deferred = [], $err = null, s;
		/* */ try { $deferFrames.push($deferred);
		s = this;
		$deferred.push([(function() {
			var e, _tuple, se, ok;
			e = $recover();
			if (!($interfaceIsEqual(e, null))) {
				_tuple = (e !== null && e.constructor === scanError ? [e.$val, true] : [new scanError.Ptr(), false]); se = new scanError.Ptr(); $copy(se, _tuple[0], scanError); ok = _tuple[1];
				if (ok) {
					err = se.err;
				} else {
					$panic(e);
				}
			}
		}), []]);
		if (f === $throwNilPointerError) {
			f = notSpace;
		}
		s.buf = $subslice(s.buf, 0, 0);
		tok = s.token(skipSpace, f);
		return [tok, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [tok, err]; }
	};
	ss.prototype.Token = function(skipSpace, f) { return this.$val.Token(skipSpace, f); };
	isSpace = function(r) {
		var rx, _ref, _i, rng;
		if (r >= 65536) {
			return false;
		}
		rx = (r << 16 >>> 16);
		_ref = space;
		_i = 0;
		while (_i < _ref.$length) {
			rng = ($arrayType($Uint16, 2)).zero(); $copy(rng, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), ($arrayType($Uint16, 2)));
			if (rx < rng[0]) {
				return false;
			}
			if (rx <= rng[1]) {
				return true;
			}
			_i++;
		}
		return false;
	};
	notSpace = function(r) {
		return !isSpace(r);
	};
	ss.Ptr.prototype.SkipSpace = function() {
		var s;
		s = this;
		s.skipSpace(false);
	};
	ss.prototype.SkipSpace = function() { return this.$val.SkipSpace(); };
	ss.Ptr.prototype.free = function(old) {
		var s;
		s = this;
		if (old.validSave) {
			$copy(s.ssave, old, ssave);
			return;
		}
		if (s.buf.$capacity > 1024) {
			return;
		}
		s.buf = $subslice(s.buf, 0, 0);
		s.rr = null;
		ssFree.Put(s);
	};
	ss.prototype.free = function(old) { return this.$val.free(old); };
	ss.Ptr.prototype.skipSpace = function(stopAtNewline) {
		var s, r;
		s = this;
		while (true) {
			r = s.getRune();
			if (r === -1) {
				return;
			}
			if ((r === 13) && s.peek("\n")) {
				continue;
			}
			if (r === 10) {
				if (stopAtNewline) {
					break;
				}
				if (s.ssave.nlIsSpace) {
					continue;
				}
				s.errorString("unexpected newline");
				return;
			}
			if (!isSpace(r)) {
				s.UnreadRune();
				break;
			}
		}
	};
	ss.prototype.skipSpace = function(stopAtNewline) { return this.$val.skipSpace(stopAtNewline); };
	ss.Ptr.prototype.token = function(skipSpace, f) {
		var s, r, x;
		s = this;
		if (skipSpace) {
			s.skipSpace(false);
		}
		while (true) {
			r = s.getRune();
			if (r === -1) {
				break;
			}
			if (!f(r)) {
				s.UnreadRune();
				break;
			}
			new ($ptrType(buffer))(function() { return this.$target.buf; }, function($v) { this.$target.buf = $v; }, s).WriteRune(r);
		}
		return (x = s.buf, $subslice(new ($sliceType($Uint8))(x.$array), x.$offset, x.$offset + x.$length));
	};
	ss.prototype.token = function(skipSpace, f) { return this.$val.token(skipSpace, f); };
	indexRune = function(s, r) {
		var _ref, _i, _rune, i, c;
		_ref = s;
		_i = 0;
		while (_i < _ref.length) {
			_rune = $decodeRune(_ref, _i);
			i = _i;
			c = _rune[0];
			if (c === r) {
				return i;
			}
			_i += _rune[1];
		}
		return -1;
	};
	ss.Ptr.prototype.peek = function(ok) {
		var s, r;
		s = this;
		r = s.getRune();
		if (!((r === -1))) {
			s.UnreadRune();
		}
		return indexRune(ok, r) >= 0;
	};
	ss.prototype.peek = function(ok) { return this.$val.peek(ok); };
	$pkg.$init = function() {
		($ptrType(fmt)).methods = [["clearflags", "clearflags", "fmt", [], [], false, -1], ["computePadding", "computePadding", "fmt", [$Int], [($sliceType($Uint8)), $Int, $Int], false, -1], ["fmt_E32", "fmt_E32", "fmt", [$Float32], [], false, -1], ["fmt_E64", "fmt_E64", "fmt", [$Float64], [], false, -1], ["fmt_G32", "fmt_G32", "fmt", [$Float32], [], false, -1], ["fmt_G64", "fmt_G64", "fmt", [$Float64], [], false, -1], ["fmt_boolean", "fmt_boolean", "fmt", [$Bool], [], false, -1], ["fmt_bx", "fmt_bx", "fmt", [($sliceType($Uint8)), $String], [], false, -1], ["fmt_c128", "fmt_c128", "fmt", [$Complex128, $Int32], [], false, -1], ["fmt_c64", "fmt_c64", "fmt", [$Complex64, $Int32], [], false, -1], ["fmt_complex", "fmt_complex", "fmt", [$Float64, $Float64, $Int, $Int32], [], false, -1], ["fmt_e32", "fmt_e32", "fmt", [$Float32], [], false, -1], ["fmt_e64", "fmt_e64", "fmt", [$Float64], [], false, -1], ["fmt_f32", "fmt_f32", "fmt", [$Float32], [], false, -1], ["fmt_f64", "fmt_f64", "fmt", [$Float64], [], false, -1], ["fmt_fb32", "fmt_fb32", "fmt", [$Float32], [], false, -1], ["fmt_fb64", "fmt_fb64", "fmt", [$Float64], [], false, -1], ["fmt_g32", "fmt_g32", "fmt", [$Float32], [], false, -1], ["fmt_g64", "fmt_g64", "fmt", [$Float64], [], false, -1], ["fmt_q", "fmt_q", "fmt", [$String], [], false, -1], ["fmt_qc", "fmt_qc", "fmt", [$Int64], [], false, -1], ["fmt_s", "fmt_s", "fmt", [$String], [], false, -1], ["fmt_sbx", "fmt_sbx", "fmt", [$String, ($sliceType($Uint8)), $String], [], false, -1], ["fmt_sx", "fmt_sx", "fmt", [$String, $String], [], false, -1], ["formatFloat", "formatFloat", "fmt", [$Float64, $Uint8, $Int, $Int], [], false, -1], ["init", "init", "fmt", [($ptrType(buffer))], [], false, -1], ["integer", "integer", "fmt", [$Int64, $Uint64, $Bool, $String], [], false, -1], ["pad", "pad", "fmt", [($sliceType($Uint8))], [], false, -1], ["padString", "padString", "fmt", [$String], [], false, -1], ["truncate", "truncate", "fmt", [$String], [$String], false, -1], ["writePadding", "writePadding", "fmt", [$Int, ($sliceType($Uint8))], [], false, -1]];
		fmt.init([["intbuf", "intbuf", "fmt", ($arrayType($Uint8, 65)), ""], ["buf", "buf", "fmt", ($ptrType(buffer)), ""], ["wid", "wid", "fmt", $Int, ""], ["prec", "prec", "fmt", $Int, ""], ["widPresent", "widPresent", "fmt", $Bool, ""], ["precPresent", "precPresent", "fmt", $Bool, ""], ["minus", "minus", "fmt", $Bool, ""], ["plus", "plus", "fmt", $Bool, ""], ["sharp", "sharp", "fmt", $Bool, ""], ["space", "space", "fmt", $Bool, ""], ["unicode", "unicode", "fmt", $Bool, ""], ["uniQuote", "uniQuote", "fmt", $Bool, ""], ["zero", "zero", "fmt", $Bool, ""]]);
		State.init([["Flag", "Flag", "", [$Int], [$Bool], false], ["Precision", "Precision", "", [], [$Int, $Bool], false], ["Width", "Width", "", [], [$Int, $Bool], false], ["Write", "Write", "", [($sliceType($Uint8))], [$Int, $error], false]]);
		Formatter.init([["Format", "Format", "", [State, $Int32], [], false]]);
		Stringer.init([["String", "String", "", [], [$String], false]]);
		GoStringer.init([["GoString", "GoString", "", [], [$String], false]]);
		($ptrType(buffer)).methods = [["Write", "Write", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["WriteByte", "WriteByte", "", [$Uint8], [$error], false, -1], ["WriteRune", "WriteRune", "", [$Int32], [$error], false, -1], ["WriteString", "WriteString", "", [$String], [$Int, $error], false, -1]];
		buffer.init($Uint8);
		($ptrType(pp)).methods = [["Flag", "Flag", "", [$Int], [$Bool], false, -1], ["Precision", "Precision", "", [], [$Int, $Bool], false, -1], ["Width", "Width", "", [], [$Int, $Bool], false, -1], ["Write", "Write", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["add", "add", "fmt", [$Int32], [], false, -1], ["argNumber", "argNumber", "fmt", [$Int, $String, $Int, $Int], [$Int, $Int, $Bool], false, -1], ["badVerb", "badVerb", "fmt", [$Int32], [], false, -1], ["catchPanic", "catchPanic", "fmt", [$emptyInterface, $Int32], [], false, -1], ["doPrint", "doPrint", "fmt", [($sliceType($emptyInterface)), $Bool, $Bool], [], false, -1], ["doPrintf", "doPrintf", "fmt", [$String, ($sliceType($emptyInterface))], [], false, -1], ["fmt0x64", "fmt0x64", "fmt", [$Uint64, $Bool], [], false, -1], ["fmtBool", "fmtBool", "fmt", [$Bool, $Int32], [], false, -1], ["fmtBytes", "fmtBytes", "fmt", [($sliceType($Uint8)), $Int32, $Bool, reflect.Type, $Int], [], false, -1], ["fmtC", "fmtC", "fmt", [$Int64], [], false, -1], ["fmtComplex128", "fmtComplex128", "fmt", [$Complex128, $Int32], [], false, -1], ["fmtComplex64", "fmtComplex64", "fmt", [$Complex64, $Int32], [], false, -1], ["fmtFloat32", "fmtFloat32", "fmt", [$Float32, $Int32], [], false, -1], ["fmtFloat64", "fmtFloat64", "fmt", [$Float64, $Int32], [], false, -1], ["fmtInt64", "fmtInt64", "fmt", [$Int64, $Int32], [], false, -1], ["fmtPointer", "fmtPointer", "fmt", [reflect.Value, $Int32, $Bool], [], false, -1], ["fmtString", "fmtString", "fmt", [$String, $Int32, $Bool], [], false, -1], ["fmtUint64", "fmtUint64", "fmt", [$Uint64, $Int32, $Bool], [], false, -1], ["fmtUnicode", "fmtUnicode", "fmt", [$Int64], [], false, -1], ["free", "free", "fmt", [], [], false, -1], ["handleMethods", "handleMethods", "fmt", [$Int32, $Bool, $Bool, $Int], [$Bool, $Bool], false, -1], ["printArg", "printArg", "fmt", [$emptyInterface, $Int32, $Bool, $Bool, $Int], [$Bool], false, -1], ["printReflectValue", "printReflectValue", "fmt", [reflect.Value, $Int32, $Bool, $Bool, $Int], [$Bool], false, -1], ["printValue", "printValue", "fmt", [reflect.Value, $Int32, $Bool, $Bool, $Int], [$Bool], false, -1], ["unknownType", "unknownType", "fmt", [$emptyInterface], [], false, -1]];
		pp.init([["n", "n", "fmt", $Int, ""], ["panicking", "panicking", "fmt", $Bool, ""], ["erroring", "erroring", "fmt", $Bool, ""], ["buf", "buf", "fmt", buffer, ""], ["arg", "arg", "fmt", $emptyInterface, ""], ["value", "value", "fmt", reflect.Value, ""], ["reordered", "reordered", "fmt", $Bool, ""], ["goodArgNum", "goodArgNum", "fmt", $Bool, ""], ["runeBuf", "runeBuf", "fmt", ($arrayType($Uint8, 4)), ""], ["fmt", "fmt", "fmt", fmt, ""]]);
		runeUnreader.init([["UnreadRune", "UnreadRune", "", [], [$error], false]]);
		scanError.init([["err", "err", "fmt", $error, ""]]);
		($ptrType(ss)).methods = [["Read", "Read", "", [($sliceType($Uint8))], [$Int, $error], false, -1], ["ReadRune", "ReadRune", "", [], [$Int32, $Int, $error], false, -1], ["SkipSpace", "SkipSpace", "", [], [], false, -1], ["Token", "Token", "", [$Bool, ($funcType([$Int32], [$Bool], false))], [($sliceType($Uint8)), $error], false, -1], ["UnreadRune", "UnreadRune", "", [], [$error], false, -1], ["Width", "Width", "", [], [$Int, $Bool], false, -1], ["accept", "accept", "fmt", [$String], [$Bool], false, -1], ["advance", "advance", "fmt", [$String], [$Int], false, -1], ["complexTokens", "complexTokens", "fmt", [], [$String, $String], false, -1], ["consume", "consume", "fmt", [$String, $Bool], [$Bool], false, -1], ["convertFloat", "convertFloat", "fmt", [$String, $Int], [$Float64], false, -1], ["convertString", "convertString", "fmt", [$Int32], [$String], false, -1], ["doScan", "doScan", "fmt", [($sliceType($emptyInterface))], [$Int, $error], false, -1], ["doScanf", "doScanf", "fmt", [$String, ($sliceType($emptyInterface))], [$Int, $error], false, -1], ["error", "error", "fmt", [$error], [], false, -1], ["errorString", "errorString", "fmt", [$String], [], false, -1], ["floatToken", "floatToken", "fmt", [], [$String], false, -1], ["free", "free", "fmt", [ssave], [], false, -1], ["getBase", "getBase", "fmt", [$Int32], [$Int, $String], false, -1], ["getRune", "getRune", "fmt", [], [$Int32], false, -1], ["hexByte", "hexByte", "fmt", [], [$Uint8, $Bool], false, -1], ["hexDigit", "hexDigit", "fmt", [$Int32], [$Int], false, -1], ["hexString", "hexString", "fmt", [], [$String], false, -1], ["mustReadRune", "mustReadRune", "fmt", [], [$Int32], false, -1], ["notEOF", "notEOF", "fmt", [], [], false, -1], ["okVerb", "okVerb", "fmt", [$Int32, $String, $String], [$Bool], false, -1], ["peek", "peek", "fmt", [$String], [$Bool], false, -1], ["quotedString", "quotedString", "fmt", [], [$String], false, -1], ["scanBasePrefix", "scanBasePrefix", "fmt", [], [$Int, $String, $Bool], false, -1], ["scanBool", "scanBool", "fmt", [$Int32], [$Bool], false, -1], ["scanComplex", "scanComplex", "fmt", [$Int32, $Int], [$Complex128], false, -1], ["scanInt", "scanInt", "fmt", [$Int32, $Int], [$Int64], false, -1], ["scanNumber", "scanNumber", "fmt", [$String, $Bool], [$String], false, -1], ["scanOne", "scanOne", "fmt", [$Int32, $emptyInterface], [], false, -1], ["scanRune", "scanRune", "fmt", [$Int], [$Int64], false, -1], ["scanUint", "scanUint", "fmt", [$Int32, $Int], [$Uint64], false, -1], ["skipSpace", "skipSpace", "fmt", [$Bool], [], false, -1], ["token", "token", "fmt", [$Bool, ($funcType([$Int32], [$Bool], false))], [($sliceType($Uint8))], false, -1]];
		ss.init([["rr", "rr", "fmt", io.RuneReader, ""], ["buf", "buf", "fmt", buffer, ""], ["peekRune", "peekRune", "fmt", $Int32, ""], ["prevRune", "prevRune", "fmt", $Int32, ""], ["count", "count", "fmt", $Int, ""], ["atEOF", "atEOF", "fmt", $Bool, ""], ["ssave", "", "fmt", ssave, ""]]);
		ssave.init([["validSave", "validSave", "fmt", $Bool, ""], ["nlIsEnd", "nlIsEnd", "fmt", $Bool, ""], ["nlIsSpace", "nlIsSpace", "fmt", $Bool, ""], ["argLimit", "argLimit", "fmt", $Int, ""], ["limit", "limit", "fmt", $Int, ""], ["maxWid", "maxWid", "fmt", $Int, ""]]);
		padZeroBytes = ($sliceType($Uint8)).make(65);
		padSpaceBytes = ($sliceType($Uint8)).make(65);
		trueBytes = new ($sliceType($Uint8))($stringToBytes("true"));
		falseBytes = new ($sliceType($Uint8))($stringToBytes("false"));
		commaSpaceBytes = new ($sliceType($Uint8))($stringToBytes(", "));
		nilAngleBytes = new ($sliceType($Uint8))($stringToBytes("<nil>"));
		nilParenBytes = new ($sliceType($Uint8))($stringToBytes("(nil)"));
		nilBytes = new ($sliceType($Uint8))($stringToBytes("nil"));
		mapBytes = new ($sliceType($Uint8))($stringToBytes("map["));
		percentBangBytes = new ($sliceType($Uint8))($stringToBytes("%!"));
		missingBytes = new ($sliceType($Uint8))($stringToBytes("(MISSING)"));
		badIndexBytes = new ($sliceType($Uint8))($stringToBytes("(BADINDEX)"));
		panicBytes = new ($sliceType($Uint8))($stringToBytes("(PANIC="));
		extraBytes = new ($sliceType($Uint8))($stringToBytes("%!(EXTRA "));
		irparenBytes = new ($sliceType($Uint8))($stringToBytes("i)"));
		bytesBytes = new ($sliceType($Uint8))($stringToBytes("[]byte{"));
		badWidthBytes = new ($sliceType($Uint8))($stringToBytes("%!(BADWIDTH)"));
		badPrecBytes = new ($sliceType($Uint8))($stringToBytes("%!(BADPREC)"));
		noVerbBytes = new ($sliceType($Uint8))($stringToBytes("%!(NOVERB)"));
		ppFree = new sync.Pool.Ptr(0, 0, ($sliceType($emptyInterface)).nil, (function() {
			return new pp.Ptr();
		}));
		intBits = reflect.TypeOf(new $Int(0)).Bits();
		uintptrBits = reflect.TypeOf(new $Uintptr(0)).Bits();
		space = new ($sliceType(($arrayType($Uint16, 2))))([$toNativeArray("Uint16", [9, 13]), $toNativeArray("Uint16", [32, 32]), $toNativeArray("Uint16", [133, 133]), $toNativeArray("Uint16", [160, 160]), $toNativeArray("Uint16", [5760, 5760]), $toNativeArray("Uint16", [8192, 8202]), $toNativeArray("Uint16", [8232, 8233]), $toNativeArray("Uint16", [8239, 8239]), $toNativeArray("Uint16", [8287, 8287]), $toNativeArray("Uint16", [12288, 12288])]);
		ssFree = new sync.Pool.Ptr(0, 0, ($sliceType($emptyInterface)).nil, (function() {
			return new ss.Ptr();
		}));
		complexError = errors.New("syntax error scanning complex number");
		boolError = errors.New("syntax error scanning boolean");
		init();
	};
	return $pkg;
})();
$packages["sort"] = (function() {
	var $pkg = {}, IntSlice, Search, SearchInts, min, insertionSort, siftDown, heapSort, medianOfThree, swapRange, doPivot, quickSort, Sort, Ints;
	IntSlice = $pkg.IntSlice = $newType(12, "Slice", "sort.IntSlice", "IntSlice", "sort", null);
	Search = $pkg.Search = function(n, f) {
		var _tmp, _tmp$1, i, j, _q, h;
		_tmp = 0; _tmp$1 = n; i = _tmp; j = _tmp$1;
		while (i < j) {
			h = i + (_q = ((j - i >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
			if (!f(h)) {
				i = h + 1 >> 0;
			} else {
				j = h;
			}
		}
		return i;
	};
	SearchInts = $pkg.SearchInts = function(a, x) {
		return Search(a.$length, (function(i) {
			return ((i < 0 || i >= a.$length) ? $throwRuntimeError("index out of range") : a.$array[a.$offset + i]) >= x;
		}));
	};
	IntSlice.prototype.Search = function(x) {
		var p;
		p = this;
		return SearchInts($subslice(new ($sliceType($Int))(p.$array), p.$offset, p.$offset + p.$length), x);
	};
	$ptrType(IntSlice).prototype.Search = function(x) { return this.$get().Search(x); };
	min = function(a, b) {
		if (a < b) {
			return a;
		}
		return b;
	};
	insertionSort = function(data, a, b) {
		var i, j;
		i = a + 1 >> 0;
		while (i < b) {
			j = i;
			while (j > a && data.Less(j, j - 1 >> 0)) {
				data.Swap(j, j - 1 >> 0);
				j = j - (1) >> 0;
			}
			i = i + (1) >> 0;
		}
	};
	siftDown = function(data, lo, hi, first) {
		var root, child;
		root = lo;
		while (true) {
			child = ((((2 >>> 16 << 16) * root >> 0) + (2 << 16 >>> 16) * root) >> 0) + 1 >> 0;
			if (child >= hi) {
				break;
			}
			if ((child + 1 >> 0) < hi && data.Less(first + child >> 0, (first + child >> 0) + 1 >> 0)) {
				child = child + (1) >> 0;
			}
			if (!data.Less(first + root >> 0, first + child >> 0)) {
				return;
			}
			data.Swap(first + root >> 0, first + child >> 0);
			root = child;
		}
	};
	heapSort = function(data, a, b) {
		var first, lo, hi, _q, i, i$1;
		first = a;
		lo = 0;
		hi = b - a >> 0;
		i = (_q = ((hi - 1 >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero"));
		while (i >= 0) {
			siftDown(data, i, hi, first);
			i = i - (1) >> 0;
		}
		i$1 = hi - 1 >> 0;
		while (i$1 >= 0) {
			data.Swap(first, first + i$1 >> 0);
			siftDown(data, lo, i$1, first);
			i$1 = i$1 - (1) >> 0;
		}
	};
	medianOfThree = function(data, a, b, c) {
		var m0, m1, m2;
		m0 = b;
		m1 = a;
		m2 = c;
		if (data.Less(m1, m0)) {
			data.Swap(m1, m0);
		}
		if (data.Less(m2, m1)) {
			data.Swap(m2, m1);
		}
		if (data.Less(m1, m0)) {
			data.Swap(m1, m0);
		}
	};
	swapRange = function(data, a, b, n) {
		var i;
		i = 0;
		while (i < n) {
			data.Swap(a + i >> 0, b + i >> 0);
			i = i + (1) >> 0;
		}
	};
	doPivot = function(data, lo, hi) {
		var midlo = 0, midhi = 0, _q, m, _q$1, s, pivot, _tmp, _tmp$1, _tmp$2, _tmp$3, a, b, c, d, n, _tmp$4, _tmp$5;
		m = lo + (_q = ((hi - lo >> 0)) / 2, (_q === _q && _q !== 1/0 && _q !== -1/0) ? _q >> 0 : $throwRuntimeError("integer divide by zero")) >> 0;
		if ((hi - lo >> 0) > 40) {
			s = (_q$1 = ((hi - lo >> 0)) / 8, (_q$1 === _q$1 && _q$1 !== 1/0 && _q$1 !== -1/0) ? _q$1 >> 0 : $throwRuntimeError("integer divide by zero"));
			medianOfThree(data, lo, lo + s >> 0, lo + ((((2 >>> 16 << 16) * s >> 0) + (2 << 16 >>> 16) * s) >> 0) >> 0);
			medianOfThree(data, m, m - s >> 0, m + s >> 0);
			medianOfThree(data, hi - 1 >> 0, (hi - 1 >> 0) - s >> 0, (hi - 1 >> 0) - ((((2 >>> 16 << 16) * s >> 0) + (2 << 16 >>> 16) * s) >> 0) >> 0);
		}
		medianOfThree(data, lo, m, hi - 1 >> 0);
		pivot = lo;
		_tmp = lo + 1 >> 0; _tmp$1 = lo + 1 >> 0; _tmp$2 = hi; _tmp$3 = hi; a = _tmp; b = _tmp$1; c = _tmp$2; d = _tmp$3;
		while (true) {
			while (b < c) {
				if (data.Less(b, pivot)) {
					b = b + (1) >> 0;
				} else if (!data.Less(pivot, b)) {
					data.Swap(a, b);
					a = a + (1) >> 0;
					b = b + (1) >> 0;
				} else {
					break;
				}
			}
			while (b < c) {
				if (data.Less(pivot, c - 1 >> 0)) {
					c = c - (1) >> 0;
				} else if (!data.Less(c - 1 >> 0, pivot)) {
					data.Swap(c - 1 >> 0, d - 1 >> 0);
					c = c - (1) >> 0;
					d = d - (1) >> 0;
				} else {
					break;
				}
			}
			if (b >= c) {
				break;
			}
			data.Swap(b, c - 1 >> 0);
			b = b + (1) >> 0;
			c = c - (1) >> 0;
		}
		n = min(b - a >> 0, a - lo >> 0);
		swapRange(data, lo, b - n >> 0, n);
		n = min(hi - d >> 0, d - c >> 0);
		swapRange(data, c, hi - n >> 0, n);
		_tmp$4 = (lo + b >> 0) - a >> 0; _tmp$5 = hi - ((d - c >> 0)) >> 0; midlo = _tmp$4; midhi = _tmp$5;
		return [midlo, midhi];
	};
	quickSort = function(data, a, b, maxDepth) {
		var _tuple, mlo, mhi;
		while ((b - a >> 0) > 7) {
			if (maxDepth === 0) {
				heapSort(data, a, b);
				return;
			}
			maxDepth = maxDepth - (1) >> 0;
			_tuple = doPivot(data, a, b); mlo = _tuple[0]; mhi = _tuple[1];
			if ((mlo - a >> 0) < (b - mhi >> 0)) {
				quickSort(data, a, mlo, maxDepth);
				a = mhi;
			} else {
				quickSort(data, mhi, b, maxDepth);
				b = mlo;
			}
		}
		if ((b - a >> 0) > 1) {
			insertionSort(data, a, b);
		}
	};
	Sort = $pkg.Sort = function(data) {
		var n, maxDepth, i, x;
		n = data.Len();
		maxDepth = 0;
		i = n;
		while (i > 0) {
			maxDepth = maxDepth + (1) >> 0;
			i = (i >> $min((1), 31)) >> 0;
		}
		maxDepth = (x = 2, (((maxDepth >>> 16 << 16) * x >> 0) + (maxDepth << 16 >>> 16) * x) >> 0);
		quickSort(data, 0, n, maxDepth);
	};
	IntSlice.prototype.Len = function() {
		var p;
		p = this;
		return p.$length;
	};
	$ptrType(IntSlice).prototype.Len = function() { return this.$get().Len(); };
	IntSlice.prototype.Less = function(i, j) {
		var p;
		p = this;
		return ((i < 0 || i >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + i]) < ((j < 0 || j >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + j]);
	};
	$ptrType(IntSlice).prototype.Less = function(i, j) { return this.$get().Less(i, j); };
	IntSlice.prototype.Swap = function(i, j) {
		var p, _tmp, _tmp$1;
		p = this;
		_tmp = ((j < 0 || j >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + j]); _tmp$1 = ((i < 0 || i >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + i]); (i < 0 || i >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + i] = _tmp; (j < 0 || j >= p.$length) ? $throwRuntimeError("index out of range") : p.$array[p.$offset + j] = _tmp$1;
	};
	$ptrType(IntSlice).prototype.Swap = function(i, j) { return this.$get().Swap(i, j); };
	IntSlice.prototype.Sort = function() {
		var p;
		p = this;
		Sort(p);
	};
	$ptrType(IntSlice).prototype.Sort = function() { return this.$get().Sort(); };
	Ints = $pkg.Ints = function(a) {
		Sort($subslice(new IntSlice(a.$array), a.$offset, a.$offset + a.$length));
	};
	$pkg.$init = function() {
		IntSlice.methods = [["Len", "Len", "", [], [$Int], false, -1], ["Less", "Less", "", [$Int, $Int], [$Bool], false, -1], ["Search", "Search", "", [$Int], [$Int], false, -1], ["Sort", "Sort", "", [], [], false, -1], ["Swap", "Swap", "", [$Int, $Int], [], false, -1]];
		($ptrType(IntSlice)).methods = [["Len", "Len", "", [], [$Int], false, -1], ["Less", "Less", "", [$Int, $Int], [$Bool], false, -1], ["Search", "Search", "", [$Int], [$Int], false, -1], ["Sort", "Sort", "", [], [], false, -1], ["Swap", "Swap", "", [$Int, $Int], [], false, -1]];
		IntSlice.init($Int);
	};
	return $pkg;
})();
$packages["path/filepath"] = (function() {
	var $pkg = {}, errors = $packages["errors"], os = $packages["os"], runtime = $packages["runtime"], sort = $packages["sort"], strings = $packages["strings"], utf8 = $packages["unicode/utf8"], bytes = $packages["bytes"];
	$pkg.$init = function() {
		$pkg.ErrBadPattern = errors.New("syntax error in pattern");
		$pkg.SkipDir = errors.New("skip this directory");
	};
	return $pkg;
})();
$packages["io/ioutil"] = (function() {
	var $pkg = {}, bytes = $packages["bytes"], io = $packages["io"], os = $packages["os"], sort = $packages["sort"], sync = $packages["sync"], filepath = $packages["path/filepath"], strconv = $packages["strconv"], time = $packages["time"], blackHolePool, readAll, ReadFile;
	readAll = function(r, capacity) {
		var b = ($sliceType($Uint8)).nil, err = null, $deferred = [], $err = null, buf, _tuple, _tmp, _tmp$1;
		/* */ try { $deferFrames.push($deferred);
		buf = bytes.NewBuffer(($sliceType($Uint8)).make(0, $flatten64(capacity)));
		$deferred.push([(function() {
			var e, _tuple, panicErr, ok;
			e = $recover();
			if ($interfaceIsEqual(e, null)) {
				return;
			}
			_tuple = (e !== null && $error.implementedBy.indexOf(e.constructor) !== -1 ? [e, true] : [null, false]); panicErr = _tuple[0]; ok = _tuple[1];
			if (ok && $interfaceIsEqual(panicErr, bytes.ErrTooLarge)) {
				err = panicErr;
			} else {
				$panic(e);
			}
		}), []]);
		_tuple = buf.ReadFrom(r); err = _tuple[1];
		_tmp = buf.Bytes(); _tmp$1 = err; b = _tmp; err = _tmp$1;
		return [b, err];
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); return [b, err]; }
	};
	ReadFile = $pkg.ReadFile = function(filename) {
		var $deferred = [], $err = null, _tuple, f, err, _recv, n, _tuple$1, fi, err$1, size;
		/* */ try { $deferFrames.push($deferred);
		_tuple = os.Open(filename); f = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, null))) {
			return [($sliceType($Uint8)).nil, err];
		}
		$deferred.push([(_recv = f, function() { $stackDepthOffset--; try { return _recv.Close(); } finally { $stackDepthOffset++; } }), []]);
		n = new $Int64(0, 0);
		_tuple$1 = f.Stat(); fi = _tuple$1[0]; err$1 = _tuple$1[1];
		if ($interfaceIsEqual(err$1, null)) {
			size = fi.Size();
			if ((size.$high < 0 || (size.$high === 0 && size.$low < 1000000000))) {
				n = size;
			}
		}
		return readAll(f, new $Int64(n.$high + 0, n.$low + 512));
		/* */ } catch(err) { $err = err; return [($sliceType($Uint8)).nil, null]; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	$pkg.$init = function() {
		blackHolePool = new sync.Pool.Ptr(0, 0, ($sliceType($emptyInterface)).nil, (function() {
			var b;
			b = ($sliceType($Uint8)).make(8192);
			return new ($ptrType(($sliceType($Uint8))))(function() { return b; }, function($v) { b = $v; });
		}));
	};
	return $pkg;
})();
$packages["path"] = (function() {
	var $pkg = {}, errors = $packages["errors"], strings = $packages["strings"], utf8 = $packages["unicode/utf8"], lazybuf, Clean, Split, Join;
	lazybuf = $pkg.lazybuf = $newType(0, "Struct", "path.lazybuf", "lazybuf", "path", function(s_, buf_, w_) {
		this.$val = this;
		this.s = s_ !== undefined ? s_ : "";
		this.buf = buf_ !== undefined ? buf_ : ($sliceType($Uint8)).nil;
		this.w = w_ !== undefined ? w_ : 0;
	});
	lazybuf.Ptr.prototype.index = function(i) {
		var b, x;
		b = this;
		if (!(b.buf === ($sliceType($Uint8)).nil)) {
			return (x = b.buf, ((i < 0 || i >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + i]));
		}
		return b.s.charCodeAt(i);
	};
	lazybuf.prototype.index = function(i) { return this.$val.index(i); };
	lazybuf.Ptr.prototype.append = function(c) {
		var b, x, x$1;
		b = this;
		if (b.buf === ($sliceType($Uint8)).nil) {
			if (b.w < b.s.length && (b.s.charCodeAt(b.w) === c)) {
				b.w = b.w + (1) >> 0;
				return;
			}
			b.buf = ($sliceType($Uint8)).make(b.s.length);
			$copyString(b.buf, b.s.substring(0, b.w));
		}
		(x = b.buf, x$1 = b.w, (x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1] = c);
		b.w = b.w + (1) >> 0;
	};
	lazybuf.prototype.append = function(c) { return this.$val.append(c); };
	lazybuf.Ptr.prototype.string = function() {
		var b;
		b = this;
		if (b.buf === ($sliceType($Uint8)).nil) {
			return b.s.substring(0, b.w);
		}
		return $bytesToString($subslice(b.buf, 0, b.w));
	};
	lazybuf.prototype.string = function() { return this.$val.string(); };
	Clean = $pkg.Clean = function(path) {
		var rooted, n, out, _tmp, _tmp$1, r, dotdot, _tmp$2, _tmp$3;
		if (path === "") {
			return ".";
		}
		rooted = path.charCodeAt(0) === 47;
		n = path.length;
		out = new lazybuf.Ptr(path, ($sliceType($Uint8)).nil, 0);
		_tmp = 0; _tmp$1 = 0; r = _tmp; dotdot = _tmp$1;
		if (rooted) {
			out.append(47);
			_tmp$2 = 1; _tmp$3 = 1; r = _tmp$2; dotdot = _tmp$3;
		}
		while (r < n) {
			if (path.charCodeAt(r) === 47) {
				r = r + (1) >> 0;
			} else if ((path.charCodeAt(r) === 46) && (((r + 1 >> 0) === n) || (path.charCodeAt((r + 1 >> 0)) === 47))) {
				r = r + (1) >> 0;
			} else if ((path.charCodeAt(r) === 46) && (path.charCodeAt((r + 1 >> 0)) === 46) && (((r + 2 >> 0) === n) || (path.charCodeAt((r + 2 >> 0)) === 47))) {
				r = r + (2) >> 0;
				if (out.w > dotdot) {
					out.w = out.w - (1) >> 0;
					while (out.w > dotdot && !((out.index(out.w) === 47))) {
						out.w = out.w - (1) >> 0;
					}
				} else if (!rooted) {
					if (out.w > 0) {
						out.append(47);
					}
					out.append(46);
					out.append(46);
					dotdot = out.w;
				}
			} else {
				if (rooted && !((out.w === 1)) || !rooted && !((out.w === 0))) {
					out.append(47);
				}
				while (r < n && !((path.charCodeAt(r) === 47))) {
					out.append(path.charCodeAt(r));
					r = r + (1) >> 0;
				}
			}
		}
		if (out.w === 0) {
			return ".";
		}
		return out.string();
	};
	Split = $pkg.Split = function(path) {
		var dir = "", file = "", i, _tmp, _tmp$1;
		i = strings.LastIndex(path, "/");
		_tmp = path.substring(0, (i + 1 >> 0)); _tmp$1 = path.substring((i + 1 >> 0)); dir = _tmp; file = _tmp$1;
		return [dir, file];
	};
	Join = $pkg.Join = function(elem) {
		var _ref, _i, i, e;
		_ref = elem;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			e = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			if (!(e === "")) {
				return Clean(strings.Join($subslice(elem, i), "/"));
			}
			_i++;
		}
		return "";
	};
	$pkg.$init = function() {
		($ptrType(lazybuf)).methods = [["append", "append", "path", [$Uint8], [], false, -1], ["index", "index", "path", [$Int], [$Uint8], false, -1], ["string", "string", "path", [], [$String], false, -1]];
		lazybuf.init([["s", "s", "path", $String, ""], ["buf", "buf", "path", ($sliceType($Uint8)), ""], ["w", "w", "path", $Int, ""]]);
		$pkg.ErrBadPattern = errors.New("syntax error in pattern");
	};
	return $pkg;
})();
$packages["github.com/hoisie/mustache"] = (function() {
	var $pkg = {}, bytes = $packages["bytes"], errors = $packages["errors"], fmt = $packages["fmt"], io = $packages["io"], ioutil = $packages["io/ioutil"], os = $packages["os"], path = $packages["path"], reflect = $packages["reflect"], strings = $packages["strings"], textElement, varElement, sectionElement, Template, parseError, esc_quot, esc_apos, esc_amp, esc_lt, esc_gt, htmlEscape, lookup, isEmpty, indirect, renderSection, renderElement, ParseString, ParseFile;
	textElement = $pkg.textElement = $newType(0, "Struct", "mustache.textElement", "textElement", "github.com/hoisie/mustache", function(text_) {
		this.$val = this;
		this.text = text_ !== undefined ? text_ : ($sliceType($Uint8)).nil;
	});
	varElement = $pkg.varElement = $newType(0, "Struct", "mustache.varElement", "varElement", "github.com/hoisie/mustache", function(name_, raw_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : "";
		this.raw = raw_ !== undefined ? raw_ : false;
	});
	sectionElement = $pkg.sectionElement = $newType(0, "Struct", "mustache.sectionElement", "sectionElement", "github.com/hoisie/mustache", function(name_, inverted_, startline_, elems_) {
		this.$val = this;
		this.name = name_ !== undefined ? name_ : "";
		this.inverted = inverted_ !== undefined ? inverted_ : false;
		this.startline = startline_ !== undefined ? startline_ : 0;
		this.elems = elems_ !== undefined ? elems_ : ($sliceType($emptyInterface)).nil;
	});
	Template = $pkg.Template = $newType(0, "Struct", "mustache.Template", "Template", "github.com/hoisie/mustache", function(data_, otag_, ctag_, p_, curline_, dir_, elems_) {
		this.$val = this;
		this.data = data_ !== undefined ? data_ : "";
		this.otag = otag_ !== undefined ? otag_ : "";
		this.ctag = ctag_ !== undefined ? ctag_ : "";
		this.p = p_ !== undefined ? p_ : 0;
		this.curline = curline_ !== undefined ? curline_ : 0;
		this.dir = dir_ !== undefined ? dir_ : "";
		this.elems = elems_ !== undefined ? elems_ : ($sliceType($emptyInterface)).nil;
	});
	parseError = $pkg.parseError = $newType(0, "Struct", "mustache.parseError", "parseError", "github.com/hoisie/mustache", function(line_, message_) {
		this.$val = this;
		this.line = line_ !== undefined ? line_ : 0;
		this.message = message_ !== undefined ? message_ : "";
	});
	parseError.Ptr.prototype.Error = function() {
		var p;
		p = new parseError.Ptr(); $copy(p, this, parseError);
		return fmt.Sprintf("line %d: %s", new ($sliceType($emptyInterface))([new $Int(p.line), new $String(p.message)]));
	};
	parseError.prototype.Error = function() { return this.$val.Error(); };
	htmlEscape = function(w, s) {
		var esc, last, _ref, _i, i, c, _ref$1;
		esc = ($sliceType($Uint8)).nil;
		last = 0;
		_ref = s;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			c = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			_ref$1 = c;
			if (_ref$1 === 34) {
				esc = esc_quot;
			} else if (_ref$1 === 39) {
				esc = esc_apos;
			} else if (_ref$1 === 38) {
				esc = esc_amp;
			} else if (_ref$1 === 60) {
				esc = esc_lt;
			} else if (_ref$1 === 62) {
				esc = esc_gt;
			} else {
				_i++;
				continue;
			}
			w.Write($subslice(s, last, i));
			w.Write(esc);
			last = i + 1 >> 0;
			_i++;
		}
		w.Write($subslice(s, last));
	};
	Template.Ptr.prototype.readString = function(s) {
		var tmpl, i, newlines, match, j, e, text;
		tmpl = this;
		i = tmpl.p;
		newlines = 0;
		while (true) {
			if ((i + s.length >> 0) > tmpl.data.length) {
				return [tmpl.data.substring(tmpl.p), io.EOF];
			}
			if (tmpl.data.charCodeAt(i) === 10) {
				newlines = newlines + (1) >> 0;
			}
			if (!((tmpl.data.charCodeAt(i) === s.charCodeAt(0)))) {
				i = i + (1) >> 0;
				continue;
			}
			match = true;
			j = 1;
			while (j < s.length) {
				if (!((s.charCodeAt(j) === tmpl.data.charCodeAt((i + j >> 0))))) {
					match = false;
					break;
				}
				j = j + (1) >> 0;
			}
			if (match) {
				e = i + s.length >> 0;
				text = tmpl.data.substring(tmpl.p, e);
				tmpl.p = e;
				tmpl.curline = tmpl.curline + (newlines) >> 0;
				return [text, null];
			} else {
				i = i + (1) >> 0;
			}
		}
		return ["", null];
	};
	Template.prototype.readString = function(s) { return this.$val.readString(s); };
	Template.Ptr.prototype.parsePartial = function(name) {
		var tmpl, filenames, filename, _ref, _i, name$1, _tuple, f, err, _tuple$1, partial, err$1;
		tmpl = this;
		filenames = new ($sliceType($String))([path.Join(new ($sliceType($String))([tmpl.dir, name])), path.Join(new ($sliceType($String))([tmpl.dir, name + ".mustache"])), path.Join(new ($sliceType($String))([tmpl.dir, name + ".stache"])), name, name + ".mustache", name + ".stache"]);
		filename = "";
		_ref = filenames;
		_i = 0;
		while (_i < _ref.$length) {
			name$1 = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			_tuple = os.Open(name$1); f = _tuple[0]; err = _tuple[1];
			if ($interfaceIsEqual(err, null)) {
				filename = name$1;
				f.Close();
				break;
			}
			_i++;
		}
		if (filename === "") {
			return [($ptrType(Template)).nil, errors.New(fmt.Sprintf("Could not find partial %q", new ($sliceType($emptyInterface))([new $String(name)])))];
		}
		_tuple$1 = ParseFile(filename); partial = _tuple$1[0]; err$1 = _tuple$1[1];
		if (!($interfaceIsEqual(err$1, null))) {
			return [($ptrType(Template)).nil, err$1];
		}
		return [partial, null];
	};
	Template.prototype.parsePartial = function(name) { return this.$val.parsePartial(name); };
	Template.Ptr.prototype.parseSection = function(section) {
		var tmpl, _tuple, text, err, x, _tuple$1, _tuple$2, x$1, tag, x$2, _ref, name, se, err$1, name$1, x$3, name$2, _tuple$3, partial, err$2, x$4, newtags;
		tmpl = this;
		while (true) {
			_tuple = tmpl.readString(tmpl.otag); text = _tuple[0]; err = _tuple[1];
			if ($interfaceIsEqual(err, io.EOF)) {
				return (x = new parseError.Ptr(section.startline, "Section " + section.name + " has no closing tag"), new x.constructor.Struct(x));
			}
			text = text.substring(0, (text.length - tmpl.otag.length >> 0));
			section.elems = $append(section.elems, new textElement.Ptr(new ($sliceType($Uint8))($stringToBytes(text))));
			if (tmpl.p < tmpl.data.length && (tmpl.data.charCodeAt(tmpl.p) === 123)) {
				_tuple$1 = tmpl.readString("}" + tmpl.ctag); text = _tuple$1[0]; err = _tuple$1[1];
			} else {
				_tuple$2 = tmpl.readString(tmpl.ctag); text = _tuple$2[0]; err = _tuple$2[1];
			}
			if ($interfaceIsEqual(err, io.EOF)) {
				return (x$1 = new parseError.Ptr(tmpl.curline, "unmatched open tag"), new x$1.constructor.Struct(x$1));
			}
			tag = strings.TrimSpace(text.substring(0, (text.length - tmpl.ctag.length >> 0)));
			if (tag.length === 0) {
				return (x$2 = new parseError.Ptr(tmpl.curline, "empty tag"), new x$2.constructor.Struct(x$2));
			}
			_ref = tag.charCodeAt(0);
			switch (0) { default: if (_ref === 33) {
				break;
			} else if (_ref === 35 || _ref === 94) {
				name = strings.TrimSpace(tag.substring(1));
				if (tmpl.data.length > tmpl.p && (tmpl.data.charCodeAt(tmpl.p) === 10)) {
					tmpl.p = tmpl.p + (1) >> 0;
				} else if (tmpl.data.length > (tmpl.p + 1 >> 0) && (tmpl.data.charCodeAt(tmpl.p) === 13) && (tmpl.data.charCodeAt((tmpl.p + 1 >> 0)) === 10)) {
					tmpl.p = tmpl.p + (2) >> 0;
				}
				se = new sectionElement.Ptr(name, tag.charCodeAt(0) === 94, tmpl.curline, new ($sliceType($emptyInterface))([]));
				err$1 = tmpl.parseSection(se);
				if (!($interfaceIsEqual(err$1, null))) {
					return err$1;
				}
				section.elems = $append(section.elems, se);
			} else if (_ref === 47) {
				name$1 = strings.TrimSpace(tag.substring(1));
				if (!(name$1 === section.name)) {
					return (x$3 = new parseError.Ptr(tmpl.curline, "interleaved closing tag: " + name$1), new x$3.constructor.Struct(x$3));
				} else {
					return null;
				}
			} else if (_ref === 62) {
				name$2 = strings.TrimSpace(tag.substring(1));
				_tuple$3 = tmpl.parsePartial(name$2); partial = _tuple$3[0]; err$2 = _tuple$3[1];
				if (!($interfaceIsEqual(err$2, null))) {
					return err$2;
				}
				section.elems = $append(section.elems, partial);
			} else if (_ref === 61) {
				if (!((tag.charCodeAt((tag.length - 1 >> 0)) === 61))) {
					return (x$4 = new parseError.Ptr(tmpl.curline, "Invalid meta tag"), new x$4.constructor.Struct(x$4));
				}
				tag = strings.TrimSpace(tag.substring(1, (tag.length - 1 >> 0)));
				newtags = strings.SplitN(tag, " ", 2);
				if (newtags.$length === 2) {
					tmpl.otag = ((0 < 0 || 0 >= newtags.$length) ? $throwRuntimeError("index out of range") : newtags.$array[newtags.$offset + 0]);
					tmpl.ctag = ((1 < 0 || 1 >= newtags.$length) ? $throwRuntimeError("index out of range") : newtags.$array[newtags.$offset + 1]);
				}
			} else if (_ref === 123) {
				if (tag.charCodeAt((tag.length - 1 >> 0)) === 125) {
					section.elems = $append(section.elems, new varElement.Ptr(tag.substring(1, (tag.length - 1 >> 0)), true));
				}
			} else {
				section.elems = $append(section.elems, new varElement.Ptr(tag, false));
			} }
		}
		return null;
	};
	Template.prototype.parseSection = function(section) { return this.$val.parseSection(section); };
	Template.Ptr.prototype.parse = function() {
		var tmpl, _tuple, text, err, _tuple$1, _tuple$2, x, tag, x$1, _ref, name, se, err$1, x$2, name$1, _tuple$3, partial, err$2, x$3, newtags;
		tmpl = this;
		while (true) {
			_tuple = tmpl.readString(tmpl.otag); text = _tuple[0]; err = _tuple[1];
			if ($interfaceIsEqual(err, io.EOF)) {
				tmpl.elems = $append(tmpl.elems, new textElement.Ptr(new ($sliceType($Uint8))($stringToBytes(text))));
				return null;
			}
			text = text.substring(0, (text.length - tmpl.otag.length >> 0));
			tmpl.elems = $append(tmpl.elems, new textElement.Ptr(new ($sliceType($Uint8))($stringToBytes(text))));
			if (tmpl.p < tmpl.data.length && (tmpl.data.charCodeAt(tmpl.p) === 123)) {
				_tuple$1 = tmpl.readString("}" + tmpl.ctag); text = _tuple$1[0]; err = _tuple$1[1];
			} else {
				_tuple$2 = tmpl.readString(tmpl.ctag); text = _tuple$2[0]; err = _tuple$2[1];
			}
			if ($interfaceIsEqual(err, io.EOF)) {
				return (x = new parseError.Ptr(tmpl.curline, "unmatched open tag"), new x.constructor.Struct(x));
			}
			tag = strings.TrimSpace(text.substring(0, (text.length - tmpl.ctag.length >> 0)));
			if (tag.length === 0) {
				return (x$1 = new parseError.Ptr(tmpl.curline, "empty tag"), new x$1.constructor.Struct(x$1));
			}
			_ref = tag.charCodeAt(0);
			switch (0) { default: if (_ref === 33) {
				break;
			} else if (_ref === 35 || _ref === 94) {
				name = strings.TrimSpace(tag.substring(1));
				if (tmpl.data.length > tmpl.p && (tmpl.data.charCodeAt(tmpl.p) === 10)) {
					tmpl.p = tmpl.p + (1) >> 0;
				} else if (tmpl.data.length > (tmpl.p + 1 >> 0) && (tmpl.data.charCodeAt(tmpl.p) === 13) && (tmpl.data.charCodeAt((tmpl.p + 1 >> 0)) === 10)) {
					tmpl.p = tmpl.p + (2) >> 0;
				}
				se = new sectionElement.Ptr(name, tag.charCodeAt(0) === 94, tmpl.curline, new ($sliceType($emptyInterface))([]));
				err$1 = tmpl.parseSection(se);
				if (!($interfaceIsEqual(err$1, null))) {
					return err$1;
				}
				tmpl.elems = $append(tmpl.elems, se);
			} else if (_ref === 47) {
				return (x$2 = new parseError.Ptr(tmpl.curline, "unmatched close tag"), new x$2.constructor.Struct(x$2));
			} else if (_ref === 62) {
				name$1 = strings.TrimSpace(tag.substring(1));
				_tuple$3 = tmpl.parsePartial(name$1); partial = _tuple$3[0]; err$2 = _tuple$3[1];
				if (!($interfaceIsEqual(err$2, null))) {
					return err$2;
				}
				tmpl.elems = $append(tmpl.elems, partial);
			} else if (_ref === 61) {
				if (!((tag.charCodeAt((tag.length - 1 >> 0)) === 61))) {
					return (x$3 = new parseError.Ptr(tmpl.curline, "Invalid meta tag"), new x$3.constructor.Struct(x$3));
				}
				tag = strings.TrimSpace(tag.substring(1, (tag.length - 1 >> 0)));
				newtags = strings.SplitN(tag, " ", 2);
				if (newtags.$length === 2) {
					tmpl.otag = ((0 < 0 || 0 >= newtags.$length) ? $throwRuntimeError("index out of range") : newtags.$array[newtags.$offset + 0]);
					tmpl.ctag = ((1 < 0 || 1 >= newtags.$length) ? $throwRuntimeError("index out of range") : newtags.$array[newtags.$offset + 1]);
				}
			} else if (_ref === 123) {
				if (tag.charCodeAt((tag.length - 1 >> 0)) === 125) {
					tmpl.elems = $append(tmpl.elems, new varElement.Ptr(tag.substring(1, (tag.length - 1 >> 0)), true));
				}
			} else {
				tmpl.elems = $append(tmpl.elems, new varElement.Ptr(tag, false));
			} }
		}
		return null;
	};
	Template.prototype.parse = function() { return this.$val.parse(); };
	lookup = function(contextChain, name) {
		var $deferred = [], $err = null, _ref, _i, ctx, v, typ, n, i, m, mtyp, x, av, _ref$1, ret, ret$1;
		/* */ try { $deferFrames.push($deferred);
		$deferred.push([(function() {
			var r;
			r = $recover();
			if (!($interfaceIsEqual(r, null))) {
				fmt.Printf("Panic while looking up %q: %s\n", new ($sliceType($emptyInterface))([new $String(name), r]));
			}
		}), []]);
		_ref = contextChain;
		_i = 0;
		Outer:
		while (_i < _ref.$length) {
			ctx = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			v = new reflect.Value.Ptr(); $copy(v, (ctx !== null && ctx.constructor === reflect.Value ? ctx.$val : $typeAssertionFailed(ctx, reflect.Value)), reflect.Value);
			while (v.IsValid()) {
				typ = v.Type();
				n = v.Type().NumMethod();
				if (n > 0) {
					i = 0;
					while (i < n) {
						m = new reflect.Method.Ptr(); $copy(m, typ.Method(i), reflect.Method);
						mtyp = m.Type;
						if (m.Name === name && (mtyp.NumIn() === 1)) {
							return (x = v.Method(i).Call(($sliceType(reflect.Value)).nil), ((0 < 0 || 0 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + 0]));
						}
						i = i + (1) >> 0;
					}
				}
				if (name === ".") {
					return v;
				}
				av = new reflect.Value.Ptr(); $copy(av, v, reflect.Value);
				_ref$1 = av.Kind();
				if (_ref$1 === 22) {
					$copy(v, av.Elem(), reflect.Value);
				} else if (_ref$1 === 20) {
					$copy(v, av.Elem(), reflect.Value);
				} else if (_ref$1 === 25) {
					ret = new reflect.Value.Ptr(); $copy(ret, av.FieldByName(name), reflect.Value);
					if (ret.IsValid()) {
						return ret;
					} else {
						_i++;
						continue Outer;
					}
				} else if (_ref$1 === 21) {
					ret$1 = new reflect.Value.Ptr(); $copy(ret$1, av.MapIndex($clone(reflect.ValueOf(new $String(name)), reflect.Value)), reflect.Value);
					if (ret$1.IsValid()) {
						return ret$1;
					} else {
						_i++;
						continue Outer;
					}
				} else {
					_i++;
					continue Outer;
				}
			}
			_i++;
		}
		return new reflect.Value.Ptr(($ptrType(reflect.rtype)).nil, 0, 0, 0);
		/* */ } catch(err) { $err = err; return new reflect.Value.Ptr(); } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	isEmpty = function(v) {
		var valueInd, val, _ref;
		if (!v.IsValid() || $interfaceIsEqual(v.Interface(), null)) {
			return true;
		}
		valueInd = new reflect.Value.Ptr(); $copy(valueInd, indirect($clone(v, reflect.Value)), reflect.Value);
		if (!valueInd.IsValid()) {
			return true;
		}
		val = new reflect.Value.Ptr(); $copy(val, valueInd, reflect.Value);
		_ref = val.Kind();
		if (_ref === 1) {
			return !val.Bool();
		} else if (_ref === 23) {
			return val.Len() === 0;
		}
		return false;
	};
	indirect = function(v) {
		var av, _ref;
		loop:
		while (v.IsValid()) {
			av = new reflect.Value.Ptr(); $copy(av, v, reflect.Value);
			_ref = av.Kind();
			if (_ref === 22) {
				$copy(v, av.Elem(), reflect.Value);
			} else if (_ref === 20) {
				$copy(v, av.Elem(), reflect.Value);
			} else {
				break loop;
			}
		}
		return v;
	};
	renderSection = function(section, contextChain, buf) {
		var value, x, x$1, context, contexts, isEmpty$1, valueInd, val, _ref, i, x$2, i$1, x$3, chain2, _ref$1, _i, ctx, _ref$2, _i$1, elem;
		value = new reflect.Value.Ptr(); $copy(value, lookup(contextChain, section.name), reflect.Value);
		context = new reflect.Value.Ptr(); $copy(context, (x = (x$1 = contextChain.$length - 1 >> 0, ((x$1 < 0 || x$1 >= contextChain.$length) ? $throwRuntimeError("index out of range") : contextChain.$array[contextChain.$offset + x$1])), (x !== null && x.constructor === reflect.Value ? x.$val : $typeAssertionFailed(x, reflect.Value))), reflect.Value);
		contexts = new ($sliceType($emptyInterface))([]);
		isEmpty$1 = isEmpty($clone(value, reflect.Value));
		if (isEmpty$1 && !section.inverted || !isEmpty$1 && section.inverted) {
			return;
		} else if (!section.inverted) {
			valueInd = new reflect.Value.Ptr(); $copy(valueInd, indirect($clone(value, reflect.Value)), reflect.Value);
			val = new reflect.Value.Ptr(); $copy(val, valueInd, reflect.Value);
			_ref = val.Kind();
			if (_ref === 23) {
				i = 0;
				while (i < val.Len()) {
					contexts = $append(contexts, (x$2 = val.Index(i), new x$2.constructor.Struct(x$2)));
					i = i + (1) >> 0;
				}
			} else if (_ref === 17) {
				i$1 = 0;
				while (i$1 < val.Len()) {
					contexts = $append(contexts, (x$3 = val.Index(i$1), new x$3.constructor.Struct(x$3)));
					i$1 = i$1 + (1) >> 0;
				}
			} else if (_ref === 21 || _ref === 25) {
				contexts = $append(contexts, new value.constructor.Struct(value));
			} else {
				contexts = $append(contexts, new context.constructor.Struct(context));
			}
		} else if (section.inverted) {
			contexts = $append(contexts, new context.constructor.Struct(context));
		}
		chain2 = ($sliceType($emptyInterface)).make((contextChain.$length + 1 >> 0));
		$copySlice($subslice(chain2, 1), contextChain);
		_ref$1 = contexts;
		_i = 0;
		while (_i < _ref$1.$length) {
			ctx = ((_i < 0 || _i >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i]);
			(0 < 0 || 0 >= chain2.$length) ? $throwRuntimeError("index out of range") : chain2.$array[chain2.$offset + 0] = ctx;
			_ref$2 = section.elems;
			_i$1 = 0;
			while (_i$1 < _ref$2.$length) {
				elem = ((_i$1 < 0 || _i$1 >= _ref$2.$length) ? $throwRuntimeError("index out of range") : _ref$2.$array[_ref$2.$offset + _i$1]);
				renderElement(elem, chain2, buf);
				_i$1++;
			}
			_i++;
		}
	};
	renderElement = function(element, contextChain, buf) {
		var $deferred = [], $err = null, elem, _ref, _type, val, s;
		/* */ try { $deferFrames.push($deferred);
		_ref = element;
		_type = _ref !== null ? _ref.constructor : null;
		if (_type === ($ptrType(textElement))) {
			elem = _ref.$val;
			buf.Write(elem.text);
		} else if (_type === ($ptrType(varElement))) {
			elem = _ref.$val;
			$deferred.push([(function() {
				var r;
				r = $recover();
				if (!($interfaceIsEqual(r, null))) {
					fmt.Printf("Panic while looking up %q: %s\n", new ($sliceType($emptyInterface))([new $String(elem.name), r]));
				}
			}), []]);
			val = new reflect.Value.Ptr(); $copy(val, lookup(contextChain, elem.name), reflect.Value);
			if (val.IsValid()) {
				if (elem.raw) {
					fmt.Fprint(buf, new ($sliceType($emptyInterface))([val.Interface()]));
				} else {
					s = fmt.Sprint(new ($sliceType($emptyInterface))([val.Interface()]));
					htmlEscape(buf, new ($sliceType($Uint8))($stringToBytes(s)));
				}
			}
		} else if (_type === ($ptrType(sectionElement))) {
			elem = _ref.$val;
			renderSection(elem, contextChain, buf);
		} else if (_type === ($ptrType(Template))) {
			elem = _ref.$val;
			elem.renderTemplate(contextChain, buf);
		}
		/* */ } catch(err) { $err = err; } finally { $deferFrames.pop(); $callDeferred($deferred, $err); }
	};
	Template.Ptr.prototype.renderTemplate = function(contextChain, buf) {
		var tmpl, _ref, _i, elem;
		tmpl = this;
		_ref = tmpl.elems;
		_i = 0;
		while (_i < _ref.$length) {
			elem = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			renderElement(elem, contextChain, buf);
			_i++;
		}
	};
	Template.prototype.renderTemplate = function(contextChain, buf) { return this.$val.renderTemplate(contextChain, buf); };
	Template.Ptr.prototype.Render = function(context) {
		var tmpl, buf, contextChain, _ref, _i, c, val;
		tmpl = this;
		buf = new bytes.Buffer.Ptr(); $copy(buf, new bytes.Buffer.Ptr(), bytes.Buffer);
		contextChain = ($sliceType($emptyInterface)).nil;
		_ref = context;
		_i = 0;
		while (_i < _ref.$length) {
			c = ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]);
			val = new reflect.Value.Ptr(); $copy(val, reflect.ValueOf(c), reflect.Value);
			contextChain = $append(contextChain, new val.constructor.Struct(val));
			_i++;
		}
		tmpl.renderTemplate(contextChain, buf);
		return buf.String();
	};
	Template.prototype.Render = function(context) { return this.$val.Render(context); };
	Template.Ptr.prototype.RenderInLayout = function(layout, context) {
		var tmpl, content, allContext, _map, _key;
		tmpl = this;
		content = tmpl.Render(context);
		allContext = ($sliceType($emptyInterface)).make((context.$length + 1 >> 0));
		$copySlice($subslice(allContext, 1), context);
		(0 < 0 || 0 >= allContext.$length) ? $throwRuntimeError("index out of range") : allContext.$array[allContext.$offset + 0] = new ($mapType($String, $String))((_map = new $Map(), _key = "content", _map[_key] = { k: _key, v: content }, _map));
		return layout.Render(allContext);
	};
	Template.prototype.RenderInLayout = function(layout, context) { return this.$val.RenderInLayout(layout, context); };
	ParseString = $pkg.ParseString = function(data) {
		var cwd, tmpl, err;
		cwd = os.Getenv("CWD");
		tmpl = new Template.Ptr(data, "{{", "}}", 0, 1, cwd, new ($sliceType($emptyInterface))([]));
		err = tmpl.parse();
		if (!($interfaceIsEqual(err, null))) {
			return [($ptrType(Template)).nil, err];
		}
		return [tmpl, err];
	};
	ParseFile = $pkg.ParseFile = function(filename) {
		var _tuple, data, err, _tuple$1, dirname, tmpl;
		_tuple = ioutil.ReadFile(filename); data = _tuple[0]; err = _tuple[1];
		if (!($interfaceIsEqual(err, null))) {
			return [($ptrType(Template)).nil, err];
		}
		_tuple$1 = path.Split(filename); dirname = _tuple$1[0];
		tmpl = new Template.Ptr($bytesToString(data), "{{", "}}", 0, 1, dirname, new ($sliceType($emptyInterface))([]));
		err = tmpl.parse();
		if (!($interfaceIsEqual(err, null))) {
			return [($ptrType(Template)).nil, err];
		}
		return [tmpl, null];
	};
	$pkg.$init = function() {
		textElement.init([["text", "text", "github.com/hoisie/mustache", ($sliceType($Uint8)), ""]]);
		varElement.init([["name", "name", "github.com/hoisie/mustache", $String, ""], ["raw", "raw", "github.com/hoisie/mustache", $Bool, ""]]);
		sectionElement.init([["name", "name", "github.com/hoisie/mustache", $String, ""], ["inverted", "inverted", "github.com/hoisie/mustache", $Bool, ""], ["startline", "startline", "github.com/hoisie/mustache", $Int, ""], ["elems", "elems", "github.com/hoisie/mustache", ($sliceType($emptyInterface)), ""]]);
		($ptrType(Template)).methods = [["Render", "Render", "", [($sliceType($emptyInterface))], [$String], true, -1], ["RenderInLayout", "RenderInLayout", "", [($ptrType(Template)), ($sliceType($emptyInterface))], [$String], true, -1], ["parse", "parse", "github.com/hoisie/mustache", [], [$error], false, -1], ["parsePartial", "parsePartial", "github.com/hoisie/mustache", [$String], [($ptrType(Template)), $error], false, -1], ["parseSection", "parseSection", "github.com/hoisie/mustache", [($ptrType(sectionElement))], [$error], false, -1], ["readString", "readString", "github.com/hoisie/mustache", [$String], [$String, $error], false, -1], ["renderTemplate", "renderTemplate", "github.com/hoisie/mustache", [($sliceType($emptyInterface)), io.Writer], [], false, -1]];
		Template.init([["data", "data", "github.com/hoisie/mustache", $String, ""], ["otag", "otag", "github.com/hoisie/mustache", $String, ""], ["ctag", "ctag", "github.com/hoisie/mustache", $String, ""], ["p", "p", "github.com/hoisie/mustache", $Int, ""], ["curline", "curline", "github.com/hoisie/mustache", $Int, ""], ["dir", "dir", "github.com/hoisie/mustache", $String, ""], ["elems", "elems", "github.com/hoisie/mustache", ($sliceType($emptyInterface)), ""]]);
		parseError.methods = [["Error", "Error", "", [], [$String], false, -1]];
		($ptrType(parseError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		parseError.init([["line", "line", "github.com/hoisie/mustache", $Int, ""], ["message", "message", "github.com/hoisie/mustache", $String, ""]]);
		esc_quot = new ($sliceType($Uint8))($stringToBytes("&quot;"));
		esc_apos = new ($sliceType($Uint8))($stringToBytes("&apos;"));
		esc_amp = new ($sliceType($Uint8))($stringToBytes("&amp;"));
		esc_lt = new ($sliceType($Uint8))($stringToBytes("&lt;"));
		esc_gt = new ($sliceType($Uint8))($stringToBytes("&gt;"));
	};
	return $pkg;
})();
$packages["github.com/scampi/gosparqled/autocompletion"] = (function() {
	var $pkg = {}, strings = $packages["strings"], mustache = $packages["github.com/hoisie/mustache"], fmt = $packages["fmt"], math = $packages["math"], sort = $packages["sort"], strconv = $packages["strconv"], triplePattern, Bgp, pegRule, tokenTree, node32, element, token16, tokens16, state16, token32, tokens32, state32, Sparql, textPosition, parseError, rul3s, translatePositions;
	triplePattern = $pkg.triplePattern = $newType(0, "Struct", "autocompletion.triplePattern", "triplePattern", "github.com/scampi/gosparqled/autocompletion", function(S_, P_, O_) {
		this.$val = this;
		this.S = S_ !== undefined ? S_ : "";
		this.P = P_ !== undefined ? P_ : "";
		this.O = O_ !== undefined ? O_ : "";
	});
	Bgp = $pkg.Bgp = $newType(0, "Struct", "autocompletion.Bgp", "Bgp", "github.com/scampi/gosparqled/autocompletion", function(triplePattern_, Tps_, scope_, Template_) {
		this.$val = this;
		this.triplePattern = triplePattern_ !== undefined ? triplePattern_ : new triplePattern.Ptr();
		this.Tps = Tps_ !== undefined ? Tps_ : ($sliceType(triplePattern)).nil;
		this.scope = scope_ !== undefined ? scope_ : false;
		this.Template = Template_ !== undefined ? Template_ : ($ptrType(mustache.Template)).nil;
	});
	pegRule = $pkg.pegRule = $newType(1, "Uint8", "autocompletion.pegRule", "pegRule", "github.com/scampi/gosparqled/autocompletion", null);
	tokenTree = $pkg.tokenTree = $newType(8, "Interface", "autocompletion.tokenTree", "tokenTree", "github.com/scampi/gosparqled/autocompletion", null);
	node32 = $pkg.node32 = $newType(0, "Struct", "autocompletion.node32", "node32", "github.com/scampi/gosparqled/autocompletion", function(token32_, up_, next_) {
		this.$val = this;
		this.token32 = token32_ !== undefined ? token32_ : new token32.Ptr();
		this.up = up_ !== undefined ? up_ : ($ptrType(node32)).nil;
		this.next = next_ !== undefined ? next_ : ($ptrType(node32)).nil;
	});
	element = $pkg.element = $newType(0, "Struct", "autocompletion.element", "element", "github.com/scampi/gosparqled/autocompletion", function(node_, down_) {
		this.$val = this;
		this.node = node_ !== undefined ? node_ : ($ptrType(node32)).nil;
		this.down = down_ !== undefined ? down_ : ($ptrType(element)).nil;
	});
	token16 = $pkg.token16 = $newType(0, "Struct", "autocompletion.token16", "token16", "github.com/scampi/gosparqled/autocompletion", function(pegRule_, begin_, end_, next_) {
		this.$val = this;
		this.pegRule = pegRule_ !== undefined ? pegRule_ : 0;
		this.begin = begin_ !== undefined ? begin_ : 0;
		this.end = end_ !== undefined ? end_ : 0;
		this.next = next_ !== undefined ? next_ : 0;
	});
	tokens16 = $pkg.tokens16 = $newType(0, "Struct", "autocompletion.tokens16", "tokens16", "github.com/scampi/gosparqled/autocompletion", function(tree_, ordered_) {
		this.$val = this;
		this.tree = tree_ !== undefined ? tree_ : ($sliceType(token16)).nil;
		this.ordered = ordered_ !== undefined ? ordered_ : ($sliceType(($sliceType(token16)))).nil;
	});
	state16 = $pkg.state16 = $newType(0, "Struct", "autocompletion.state16", "state16", "github.com/scampi/gosparqled/autocompletion", function(token16_, depths_, leaf_) {
		this.$val = this;
		this.token16 = token16_ !== undefined ? token16_ : new token16.Ptr();
		this.depths = depths_ !== undefined ? depths_ : ($sliceType($Int16)).nil;
		this.leaf = leaf_ !== undefined ? leaf_ : false;
	});
	token32 = $pkg.token32 = $newType(0, "Struct", "autocompletion.token32", "token32", "github.com/scampi/gosparqled/autocompletion", function(pegRule_, begin_, end_, next_) {
		this.$val = this;
		this.pegRule = pegRule_ !== undefined ? pegRule_ : 0;
		this.begin = begin_ !== undefined ? begin_ : 0;
		this.end = end_ !== undefined ? end_ : 0;
		this.next = next_ !== undefined ? next_ : 0;
	});
	tokens32 = $pkg.tokens32 = $newType(0, "Struct", "autocompletion.tokens32", "tokens32", "github.com/scampi/gosparqled/autocompletion", function(tree_, ordered_) {
		this.$val = this;
		this.tree = tree_ !== undefined ? tree_ : ($sliceType(token32)).nil;
		this.ordered = ordered_ !== undefined ? ordered_ : ($sliceType(($sliceType(token32)))).nil;
	});
	state32 = $pkg.state32 = $newType(0, "Struct", "autocompletion.state32", "state32", "github.com/scampi/gosparqled/autocompletion", function(token32_, depths_, leaf_) {
		this.$val = this;
		this.token32 = token32_ !== undefined ? token32_ : new token32.Ptr();
		this.depths = depths_ !== undefined ? depths_ : ($sliceType($Int32)).nil;
		this.leaf = leaf_ !== undefined ? leaf_ : false;
	});
	Sparql = $pkg.Sparql = $newType(0, "Struct", "autocompletion.Sparql", "Sparql", "github.com/scampi/gosparqled/autocompletion", function(Bgp_, Buffer_, buffer_, rules_, Parse_, Reset_, tokenTree_) {
		this.$val = this;
		this.Bgp = Bgp_ !== undefined ? Bgp_ : ($ptrType(Bgp)).nil;
		this.Buffer = Buffer_ !== undefined ? Buffer_ : "";
		this.buffer = buffer_ !== undefined ? buffer_ : ($sliceType($Int32)).nil;
		this.rules = rules_ !== undefined ? rules_ : ($arrayType(($funcType([], [$Bool], false)), 99)).zero();
		this.Parse = Parse_ !== undefined ? Parse_ : $throwNilPointerError;
		this.Reset = Reset_ !== undefined ? Reset_ : $throwNilPointerError;
		this.tokenTree = tokenTree_ !== undefined ? tokenTree_ : null;
	});
	textPosition = $pkg.textPosition = $newType(0, "Struct", "autocompletion.textPosition", "textPosition", "github.com/scampi/gosparqled/autocompletion", function(line_, symbol_) {
		this.$val = this;
		this.line = line_ !== undefined ? line_ : 0;
		this.symbol = symbol_ !== undefined ? symbol_ : 0;
	});
	parseError = $pkg.parseError = $newType(0, "Struct", "autocompletion.parseError", "parseError", "github.com/scampi/gosparqled/autocompletion", function(p_) {
		this.$val = this;
		this.p = p_ !== undefined ? p_ : ($ptrType(Sparql)).nil;
	});
	Bgp.Ptr.prototype.setSubject = function(s) {
		var b;
		b = this;
		s = strings.TrimSpace(s);
		if ((s.length === 0)) {
			return;
		}
		b.triplePattern.S = s;
	};
	Bgp.prototype.setSubject = function(s) { return this.$val.setSubject(s); };
	Bgp.Ptr.prototype.setPredicate = function(p) {
		var b;
		b = this;
		p = strings.TrimSpace(p);
		if ((p.length === 0)) {
			return;
		}
		b.triplePattern.P = p;
	};
	Bgp.prototype.setPredicate = function(p) { return this.$val.setPredicate(p); };
	Bgp.Ptr.prototype.setObject = function(o) {
		var b;
		b = this;
		o = strings.TrimSpace(o);
		if ((o.length === 0)) {
			return;
		}
		b.triplePattern.O = o;
	};
	Bgp.prototype.setObject = function(o) { return this.$val.setObject(o); };
	Bgp.Ptr.prototype.addTriplePattern = function() {
		var b, tp;
		b = this;
		tp = new triplePattern.Ptr(b.triplePattern.S, b.triplePattern.P, b.triplePattern.O);
		b.Tps = $append(b.Tps, tp);
	};
	Bgp.prototype.addTriplePattern = function() { return this.$val.addTriplePattern(); };
	Bgp.Ptr.prototype.trimToScope = function() {
		var b, _map, _key, size, _ref, _i, tp, scoped, _ref$1, _i$1, tp$1;
		b = this;
		b.scope = (_map = new $Map(), _key = "?POF", _map[_key] = { k: _key, v: true }, _map);
		size = 0;
		while (!((size === $keys(b.scope).length))) {
			size = $keys(b.scope).length;
			_ref = b.Tps;
			_i = 0;
			while (_i < _ref.$length) {
				tp = new triplePattern.Ptr(); $copy(tp, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), triplePattern);
				if (tp.in$(b.scope)) {
					tp.addToScope(b.scope);
				}
				_i++;
			}
		}
		scoped = ($sliceType(triplePattern)).nil;
		_ref$1 = b.Tps;
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			tp$1 = new triplePattern.Ptr(); $copy(tp$1, ((_i$1 < 0 || _i$1 >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i$1]), triplePattern);
			if (tp$1.in$(b.scope)) {
				scoped = $append(scoped, tp$1);
			}
			_i$1++;
		}
		b.Tps = scoped;
	};
	Bgp.prototype.trimToScope = function() { return this.$val.trimToScope(); };
	triplePattern.Ptr.prototype.in$ = function(scope) {
		var tp, _entry, _entry$1, _entry$2;
		tp = this;
		if ((_entry = scope[tp.S], _entry !== undefined ? _entry.v : false) || (_entry$1 = scope[tp.P], _entry$1 !== undefined ? _entry$1.v : false) || (_entry$2 = scope[tp.O], _entry$2 !== undefined ? _entry$2.v : false)) {
			return true;
		}
		return false;
	};
	triplePattern.prototype.in$ = function(scope) { return this.$val.in$(scope); };
	triplePattern.Ptr.prototype.addToScope = function(scope) {
		var tp, _key, _key$1, _key$2;
		tp = this;
		_key = tp.S; (scope || $throwRuntimeError("assignment to entry in nil map"))[_key] = { k: _key, v: true };
		_key$1 = tp.P; (scope || $throwRuntimeError("assignment to entry in nil map"))[_key$1] = { k: _key$1, v: true };
		_key$2 = tp.O; (scope || $throwRuntimeError("assignment to entry in nil map"))[_key$2] = { k: _key$2, v: true };
	};
	triplePattern.prototype.addToScope = function(scope) { return this.$val.addToScope(scope); };
	Bgp.Ptr.prototype.RecommendationQuery = function() {
		var b;
		b = this;
		b.trimToScope();
		return b.Template.Render(new ($sliceType($emptyInterface))([b]));
	};
	Bgp.prototype.RecommendationQuery = function() { return this.$val.RecommendationQuery(); };
	node32.Ptr.prototype.print = function(depth, buffer) {
		var node, c, x;
		node = this;
		while (!(node === ($ptrType(node32)).nil)) {
			c = 0;
			while (c < depth) {
				fmt.Printf(" ", new ($sliceType($emptyInterface))([]));
				c = c + (1) >> 0;
			}
			fmt.Printf("\x1B[34m%v\x1B[m %v\n", new ($sliceType($emptyInterface))([new $String((x = node.token32.pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x]))), new $String(strconv.Quote(buffer.substring(node.token32.begin, node.token32.end)))]));
			if (!(node.up === ($ptrType(node32)).nil)) {
				node.up.print(depth + 1 >> 0, buffer);
			}
			node = node.next;
		}
	};
	node32.prototype.print = function(depth, buffer) { return this.$val.print(depth, buffer); };
	node32.Ptr.prototype.Print = function(buffer) {
		var ast;
		ast = this;
		ast.print(0, buffer);
	};
	node32.prototype.Print = function(buffer) { return this.$val.Print(buffer); };
	token16.Ptr.prototype.isParentOf = function(u) {
		var t;
		t = this;
		return t.begin <= u.begin && t.end >= u.end && t.next > u.next;
	};
	token16.prototype.isParentOf = function(u) { return this.$val.isParentOf(u); };
	token16.Ptr.prototype.getToken32 = function() {
		var t;
		t = this;
		return new token32.Ptr(t.pegRule, (t.begin >> 0), (t.end >> 0), (t.next >> 0));
	};
	token16.prototype.getToken32 = function() { return this.$val.getToken32(); };
	token16.Ptr.prototype.String = function() {
		var t, x;
		t = this;
		return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", new ($sliceType($emptyInterface))([new $String((x = t.pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x]))), new $Int16(t.begin), new $Int16(t.end), new $Int16(t.next)]));
	};
	token16.prototype.String = function() { return this.$val.String(); };
	tokens16.Ptr.prototype.trim = function(length) {
		var t;
		t = this;
		t.tree = $subslice(t.tree, 0, length);
	};
	tokens16.prototype.trim = function(length) { return this.$val.trim(length); };
	tokens16.Ptr.prototype.Print = function() {
		var t, _ref, _i, token;
		t = this;
		_ref = t.tree;
		_i = 0;
		while (_i < _ref.$length) {
			token = new token16.Ptr(); $copy(token, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token16);
			fmt.Println(new ($sliceType($emptyInterface))([new $String(token.String())]));
			_i++;
		}
	};
	tokens16.prototype.Print = function() { return this.$val.Print(); };
	tokens16.Ptr.prototype.Order = function() {
		var t, depths, _ref, _i, i, token, depth, length, _lhs, _index, _tmp, _tmp$1, ordered, pool, _ref$1, _i$1, i$1, depth$1, _tmp$2, _tmp$3, _tmp$4, _ref$2, _i$2, i$2, token$1, depth$2, x, x$1, _lhs$1, _index$1;
		t = this;
		if (!(t.ordered === ($sliceType(($sliceType(token16)))).nil)) {
			return t.ordered;
		}
		depths = ($sliceType($Int16)).make(1, 32767);
		_ref = t.tree;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			token = new token16.Ptr(); $copy(token, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token16);
			if (token.pegRule === 0) {
				t.tree = $subslice(t.tree, 0, i);
				break;
			}
			depth = (token.next >> 0);
			length = depths.$length;
			if (depth >= length) {
				depths = $subslice(depths, 0, (depth + 1 >> 0));
			}
			_lhs = depths; _index = depth; (_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index] = ((_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index]) + (1) << 16 >> 16;
			_i++;
		}
		depths = $append(depths, 0);
		_tmp = ($sliceType(($sliceType(token16)))).make(depths.$length); _tmp$1 = ($sliceType(token16)).make((t.tree.$length + depths.$length >> 0)); ordered = _tmp; pool = _tmp$1;
		_ref$1 = depths;
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			i$1 = _i$1;
			depth$1 = ((_i$1 < 0 || _i$1 >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i$1]);
			depth$1 = depth$1 + (1) << 16 >> 16;
			_tmp$2 = $subslice(pool, 0, depth$1); _tmp$3 = $subslice(pool, depth$1); _tmp$4 = 0; (i$1 < 0 || i$1 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i$1] = _tmp$2; pool = _tmp$3; (i$1 < 0 || i$1 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + i$1] = _tmp$4;
			_i$1++;
		}
		_ref$2 = t.tree;
		_i$2 = 0;
		while (_i$2 < _ref$2.$length) {
			i$2 = _i$2;
			token$1 = new token16.Ptr(); $copy(token$1, ((_i$2 < 0 || _i$2 >= _ref$2.$length) ? $throwRuntimeError("index out of range") : _ref$2.$array[_ref$2.$offset + _i$2]), token16);
			depth$2 = token$1.next;
			token$1.next = (i$2 << 16 >> 16);
			$copy((x = ((depth$2 < 0 || depth$2 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth$2]), x$1 = ((depth$2 < 0 || depth$2 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth$2]), ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1])), token$1, token16);
			_lhs$1 = depths; _index$1 = depth$2; (_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1] = ((_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1]) + (1) << 16 >> 16;
			_i$2++;
		}
		t.ordered = ordered;
		return ordered;
	};
	tokens16.prototype.Order = function() { return this.$val.Order(); };
	tokens16.Ptr.prototype.AST = function($b) {
		var $this = this, $args = arguments, $r, $s = 0, t, tokens, _r, stack, _ref, _ok, _tuple, _r$1, token, node;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		t = $this;
		tokens = t.Tokens();
		_r = $recv(tokens, true); /* */ $s = 1; case 1: if (_r && _r.constructor === Function) { _r = _r(); }
		stack = new element.Ptr(new node32.Ptr(_r[0], ($ptrType(node32)).nil, ($ptrType(node32)).nil), ($ptrType(element)).nil);
		_ref = tokens;
		/* while (true) { */ case 2: if(!(true)) { $s = 3; continue; }
			_r$1 = $recv(_ref, true); /* */ $s = 4; case 4: if (_r$1 && _r$1.constructor === Function) { _r$1 = _r$1(); }
			_tuple = _r$1; token = new token32.Ptr(); $copy(token, _tuple[0], token32); _ok = _tuple[1];
			if (!_ok) {
				/* break; */ $s = 3; continue;
			}
			if (token.begin === token.end) {
				/* continue; */ $s = 2; continue;
			}
			node = new node32.Ptr(token, ($ptrType(node32)).nil, ($ptrType(node32)).nil);
			while (!(stack === ($ptrType(element)).nil) && stack.node.token32.begin >= token.begin && stack.node.token32.end <= token.end) {
				stack.node.next = node.up;
				node.up = stack.node;
				stack = stack.down;
			}
			stack = new element.Ptr(node, stack);
		/* } */ $s = 2; continue; case 3:
		return stack.node;
		/* */ case -1: } return; } };
	};
	tokens16.prototype.AST = function($b) { return this.$val.AST($b); };
	tokens16.Ptr.prototype.PreOrder = function() {
		var t, _tmp, _tmp$1, s, ordered;
		t = this;
		_tmp = new ($chanType(state16, false, false))(6); _tmp$1 = t.Order(); s = _tmp; ordered = _tmp$1;
		$go((function() {
			var states, _ref, _i, i, _tmp$2, _tmp$3, _tmp$4, depths, state, depth, write, x, _lhs, _index, _tmp$5, x$1, x$2, x$3, x$4, _tmp$6, x$5, x$6, a, b, i$1, _tmp$7, x$7, x$8, _tmp$8, x$9, c, j, x$10, x$11, x$12, next, x$13, x$14, c$1, _lhs$1, _index$1, _tmp$9, _tmp$10, _tmp$11, _lhs$2, _index$2, _tmp$12, x$15, x$16, _tmp$13, c$2, parent, _tmp$14, x$17, x$18, x$19, x$20, _tmp$15, _tmp$16, x$21, x$22;
			states = ($arrayType(state16, 8)).zero(); $copy(states, ($arrayType(state16, 8)).zero(), ($arrayType(state16, 8)));
			_ref = states;
			_i = 0;
			while (_i < 8) {
				i = _i;
				((i < 0 || i >= states.length) ? $throwRuntimeError("index out of range") : states[i]).depths = ($sliceType($Int16)).make(ordered.$length);
				_i++;
			}
			_tmp$2 = ($sliceType($Int16)).make(ordered.$length); _tmp$3 = 0; _tmp$4 = 1; depths = _tmp$2; state = _tmp$3; depth = _tmp$4;
			write = (function(t$1, leaf, $b) {
				var $this = this, $args = arguments, $r, $s = 0, S, _tmp$5, _r, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10;
				/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
				S = new state16.Ptr(); $copy(S, ((state < 0 || state >= states.length) ? $throwRuntimeError("index out of range") : states[state]), state16);
				_tmp$5 = (_r = ((state + 1 >> 0)) % 8, _r === _r ? _r : $throwRuntimeError("integer divide by zero")); _tmp$6 = t$1.pegRule; _tmp$7 = t$1.begin; _tmp$8 = t$1.end; _tmp$9 = (depth << 16 >> 16); _tmp$10 = leaf; state = _tmp$5; S.token16.pegRule = _tmp$6; S.token16.begin = _tmp$7; S.token16.end = _tmp$8; S.token16.next = _tmp$9; S.leaf = _tmp$10;
				$copySlice(S.depths, depths);
				$r = $send(s, $clone(S, state16), true); /* */ $s = 1; case 1: if ($r && $r.constructor === Function) { $r = $r(); }
				/* */ case -1: } return; } };
			});
			$copy(((state < 0 || state >= states.length) ? $throwRuntimeError("index out of range") : states[state]).token16, (x = ((0 < 0 || 0 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + 0]), ((0 < 0 || 0 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + 0])), token16);
			_lhs = depths; _index = 0; (_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index] = ((_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index]) + (1) << 16 >> 16;
			state = state + (1) >> 0;
			_tmp$5 = new token16.Ptr(); $copy(_tmp$5, (x$1 = (x$2 = depth - 1 >> 0, ((x$2 < 0 || x$2 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x$2])), x$3 = (x$4 = depth - 1 >> 0, ((x$4 < 0 || x$4 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + x$4])) - 1 << 16 >> 16, ((x$3 < 0 || x$3 >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + x$3])), token16); _tmp$6 = new token16.Ptr(); $copy(_tmp$6, (x$5 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$6 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]), ((x$6 < 0 || x$6 >= x$5.$length) ? $throwRuntimeError("index out of range") : x$5.$array[x$5.$offset + x$6])), token16); a = new token16.Ptr(); $copy(a, _tmp$5, token16); b = new token16.Ptr(); $copy(b, _tmp$6, token16);
			depthFirstSearch:
			while (true) {
				while (true) {
					i$1 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]);
					if (i$1 > 0) {
						_tmp$7 = new token16.Ptr(); $copy(_tmp$7, (x$7 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$8 = i$1 - 1 << 16 >> 16, ((x$8 < 0 || x$8 >= x$7.$length) ? $throwRuntimeError("index out of range") : x$7.$array[x$7.$offset + x$8])), token16); _tmp$8 = (x$9 = depth - 1 >> 0, ((x$9 < 0 || x$9 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + x$9])); c = new token16.Ptr(); $copy(c, _tmp$7, token16); j = _tmp$8;
						if (a.isParentOf($clone(c, token16)) && (j < 2 || !(x$10 = (x$11 = depth - 1 >> 0, ((x$11 < 0 || x$11 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x$11])), x$12 = j - 2 << 16 >> 16, ((x$12 < 0 || x$12 >= x$10.$length) ? $throwRuntimeError("index out of range") : x$10.$array[x$10.$offset + x$12])).isParentOf($clone(c, token16)))) {
							if (!((c.end === b.begin))) {
								write(new token16.Ptr(100, c.end, b.begin, 0), true);
							}
							break;
						}
					}
					if (a.begin < b.begin) {
						write(new token16.Ptr(99, a.begin, b.begin, 0), true);
					}
					break;
				}
				next = depth + 1 >> 0;
				c$1 = new token16.Ptr(); $copy(c$1, (x$13 = ((next < 0 || next >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + next]), x$14 = ((next < 0 || next >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + next]), ((x$14 < 0 || x$14 >= x$13.$length) ? $throwRuntimeError("index out of range") : x$13.$array[x$13.$offset + x$14])), token16);
				if (!((c$1.pegRule === 0)) && b.isParentOf($clone(c$1, token16))) {
					write($clone(b, token16), false);
					_lhs$1 = depths; _index$1 = depth; (_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1] = ((_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1]) + (1) << 16 >> 16;
					_tmp$9 = next; _tmp$10 = new token16.Ptr(); $copy(_tmp$10, b, token16); _tmp$11 = new token16.Ptr(); $copy(_tmp$11, c$1, token16); depth = _tmp$9; $copy(a, _tmp$10, token16); $copy(b, _tmp$11, token16);
					continue;
				}
				write($clone(b, token16), true);
				_lhs$2 = depths; _index$2 = depth; (_index$2 < 0 || _index$2 >= _lhs$2.$length) ? $throwRuntimeError("index out of range") : _lhs$2.$array[_lhs$2.$offset + _index$2] = ((_index$2 < 0 || _index$2 >= _lhs$2.$length) ? $throwRuntimeError("index out of range") : _lhs$2.$array[_lhs$2.$offset + _index$2]) + (1) << 16 >> 16;
				_tmp$12 = new token16.Ptr(); $copy(_tmp$12, (x$15 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$16 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]), ((x$16 < 0 || x$16 >= x$15.$length) ? $throwRuntimeError("index out of range") : x$15.$array[x$15.$offset + x$16])), token16); _tmp$13 = true; c$2 = new token16.Ptr(); $copy(c$2, _tmp$12, token16); parent = _tmp$13;
				while (true) {
					if (!((c$2.pegRule === 0)) && a.isParentOf($clone(c$2, token16))) {
						$copy(b, c$2, token16);
						continue depthFirstSearch;
					} else if (parent && !((b.end === a.end))) {
						write(new token16.Ptr(101, b.end, a.end, 0), true);
					}
					depth = depth - (1) >> 0;
					if (depth > 0) {
						_tmp$14 = new token16.Ptr(); $copy(_tmp$14, (x$17 = (x$18 = depth - 1 >> 0, ((x$18 < 0 || x$18 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x$18])), x$19 = (x$20 = depth - 1 >> 0, ((x$20 < 0 || x$20 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + x$20])) - 1 << 16 >> 16, ((x$19 < 0 || x$19 >= x$17.$length) ? $throwRuntimeError("index out of range") : x$17.$array[x$17.$offset + x$19])), token16); _tmp$15 = new token16.Ptr(); $copy(_tmp$15, a, token16); _tmp$16 = new token16.Ptr(); $copy(_tmp$16, (x$21 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$22 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]), ((x$22 < 0 || x$22 >= x$21.$length) ? $throwRuntimeError("index out of range") : x$21.$array[x$21.$offset + x$22])), token16); $copy(a, _tmp$14, token16); $copy(b, _tmp$15, token16); $copy(c$2, _tmp$16, token16);
						parent = a.isParentOf($clone(b, token16));
						continue;
					}
					break depthFirstSearch;
				}
			}
			$close(s);
		}), []);
		return [s, ordered];
	};
	tokens16.prototype.PreOrder = function() { return this.$val.PreOrder(); };
	tokens16.Ptr.prototype.PrintSyntax = function($b) {
		var $this = this, $args = arguments, $r, $s = 0, t, _tuple, tokens, ordered, max, _ref, _ok, _tuple$1, _r, token, _tmp, _tmp$1, _tmp$2, i, leaf, depths, x, x$1, x$2, x$3, _tmp$3, _tmp$4, _tmp$5, i$1, leaf$1, depths$1, x$4, x$5, x$6, x$7, _tmp$6, _tmp$7, c, end, i$2, j, i$3, j$1, _tmp$8, _tmp$9, _tmp$10, i$4, leaf$2, depths$2, x$8, x$9, x$10, x$11;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		t = $this;
		_tuple = t.PreOrder(); tokens = _tuple[0]; ordered = _tuple[1];
		max = -1;
		_ref = tokens;
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			_r = $recv(_ref, true); /* */ $s = 3; case 3: if (_r && _r.constructor === Function) { _r = _r(); }
			_tuple$1 = _r; token = new state16.Ptr(); $copy(token, _tuple$1[0], state16); _ok = _tuple$1[1];
			if (!_ok) {
				/* break; */ $s = 2; continue;
			}
			if (!token.leaf) {
				fmt.Printf("%v", new ($sliceType($emptyInterface))([new $Int16(token.token16.begin)]));
				_tmp = 0; _tmp$1 = (token.token16.next >> 0); _tmp$2 = token.depths; i = _tmp; leaf = _tmp$1; depths = _tmp$2;
				while (i < leaf) {
					fmt.Printf(" \x1B[36m%v\x1B[m", new ($sliceType($emptyInterface))([new $String((x = (x$1 = ((i < 0 || i >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i]), x$2 = ((i < 0 || i >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + i]) - 1 << 16 >> 16, ((x$2 < 0 || x$2 >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + x$2])).pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x])))]));
					i = i + (1) >> 0;
				}
				fmt.Printf(" \x1B[36m%v\x1B[m\n", new ($sliceType($emptyInterface))([new $String((x$3 = token.token16.pegRule, ((x$3 < 0 || x$3 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$3])))]));
			} else if (token.token16.begin === token.token16.end) {
				fmt.Printf("%v", new ($sliceType($emptyInterface))([new $Int16(token.token16.begin)]));
				_tmp$3 = 0; _tmp$4 = (token.token16.next >> 0); _tmp$5 = token.depths; i$1 = _tmp$3; leaf$1 = _tmp$4; depths$1 = _tmp$5;
				while (i$1 < leaf$1) {
					fmt.Printf(" \x1B[31m%v\x1B[m", new ($sliceType($emptyInterface))([new $String((x$4 = (x$5 = ((i$1 < 0 || i$1 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i$1]), x$6 = ((i$1 < 0 || i$1 >= depths$1.$length) ? $throwRuntimeError("index out of range") : depths$1.$array[depths$1.$offset + i$1]) - 1 << 16 >> 16, ((x$6 < 0 || x$6 >= x$5.$length) ? $throwRuntimeError("index out of range") : x$5.$array[x$5.$offset + x$6])).pegRule, ((x$4 < 0 || x$4 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$4])))]));
					i$1 = i$1 + (1) >> 0;
				}
				fmt.Printf(" \x1B[31m%v\x1B[m\n", new ($sliceType($emptyInterface))([new $String((x$7 = token.token16.pegRule, ((x$7 < 0 || x$7 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$7])))]));
			} else {
				_tmp$6 = token.token16.begin; _tmp$7 = token.token16.end; c = _tmp$6; end = _tmp$7;
				while (c < end) {
					i$2 = (c >> 0);
					if ((max + 1 >> 0) < i$2) {
						j = max;
						while (j < i$2) {
							fmt.Printf("skip %v %v\n", new ($sliceType($emptyInterface))([new $Int(j), new $String(token.token16.String())]));
							j = j + (1) >> 0;
						}
						max = i$2;
					} else {
						i$3 = (c >> 0);
						if (i$3 <= max) {
							j$1 = i$3;
							while (j$1 <= max) {
								fmt.Printf("dupe %v %v\n", new ($sliceType($emptyInterface))([new $Int(j$1), new $String(token.token16.String())]));
								j$1 = j$1 + (1) >> 0;
							}
						} else {
							max = (c >> 0);
						}
					}
					fmt.Printf("%v", new ($sliceType($emptyInterface))([new $Int16(c)]));
					_tmp$8 = 0; _tmp$9 = (token.token16.next >> 0); _tmp$10 = token.depths; i$4 = _tmp$8; leaf$2 = _tmp$9; depths$2 = _tmp$10;
					while (i$4 < leaf$2) {
						fmt.Printf(" \x1B[34m%v\x1B[m", new ($sliceType($emptyInterface))([new $String((x$8 = (x$9 = ((i$4 < 0 || i$4 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i$4]), x$10 = ((i$4 < 0 || i$4 >= depths$2.$length) ? $throwRuntimeError("index out of range") : depths$2.$array[depths$2.$offset + i$4]) - 1 << 16 >> 16, ((x$10 < 0 || x$10 >= x$9.$length) ? $throwRuntimeError("index out of range") : x$9.$array[x$9.$offset + x$10])).pegRule, ((x$8 < 0 || x$8 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$8])))]));
						i$4 = i$4 + (1) >> 0;
					}
					fmt.Printf(" \x1B[34m%v\x1B[m\n", new ($sliceType($emptyInterface))([new $String((x$11 = token.token16.pegRule, ((x$11 < 0 || x$11 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$11])))]));
					c = c + (1) << 16 >> 16;
				}
				fmt.Printf("\n", new ($sliceType($emptyInterface))([]));
			}
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } };
	};
	tokens16.prototype.PrintSyntax = function($b) { return this.$val.PrintSyntax($b); };
	tokens16.Ptr.prototype.PrintSyntaxTree = function(buffer, $b) {
		var $this = this, $args = arguments, $r, $s = 0, t, _tuple, tokens, _ref, _ok, _tuple$1, _r, token, c, x;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		t = $this;
		_tuple = t.PreOrder(); tokens = _tuple[0];
		_ref = tokens;
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			_r = $recv(_ref, true); /* */ $s = 3; case 3: if (_r && _r.constructor === Function) { _r = _r(); }
			_tuple$1 = _r; token = new state16.Ptr(); $copy(token, _tuple$1[0], state16); _ok = _tuple$1[1];
			if (!_ok) {
				/* break; */ $s = 2; continue;
			}
			c = 0;
			while (c < (token.token16.next >> 0)) {
				fmt.Printf(" ", new ($sliceType($emptyInterface))([]));
				c = c + (1) >> 0;
			}
			fmt.Printf("\x1B[34m%v\x1B[m %v\n", new ($sliceType($emptyInterface))([new $String((x = token.token16.pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x]))), new $String(strconv.Quote(buffer.substring(token.token16.begin, token.token16.end)))]));
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } };
	};
	tokens16.prototype.PrintSyntaxTree = function(buffer, $b) { return this.$val.PrintSyntaxTree(buffer, $b); };
	tokens16.Ptr.prototype.Add = function(rule, begin, end, depth, index) {
		var t, x;
		t = this;
		$copy((x = t.tree, ((index < 0 || index >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + index])), new token16.Ptr(rule, (begin << 16 >> 16), (end << 16 >> 16), (depth << 16 >> 16)), token16);
	};
	tokens16.prototype.Add = function(rule, begin, end, depth, index) { return this.$val.Add(rule, begin, end, depth, index); };
	tokens16.Ptr.prototype.Tokens = function() {
		var t, s;
		t = this;
		s = new ($chanType(token32, false, false))(16);
		$go((function($b) {
			var $this = this, $args = arguments, $r, $s = 0, _ref, _i, v;
			/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
			_ref = t.tree;
			_i = 0;
			/* while (_i < _ref.$length) { */ case 1: if(!(_i < _ref.$length)) { $s = 2; continue; }
				v = new token16.Ptr(); $copy(v, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token16);
				$r = $send(s, $clone(v.getToken32(), token32), true); /* */ $s = 3; case 3: if ($r && $r.constructor === Function) { $r = $r(); }
				_i++;
			/* } */ $s = 1; continue; case 2:
			$close(s);
			/* */ case -1: } return; } };
		}), []);
		return s;
	};
	tokens16.prototype.Tokens = function() { return this.$val.Tokens(); };
	tokens16.Ptr.prototype.Error = function() {
		var t, ordered, length, _tmp, _tmp$1, tokens, _ref, _i, i, x, o, x$1;
		t = this;
		ordered = t.Order();
		length = ordered.$length;
		_tmp = ($sliceType(token32)).make(length); _tmp$1 = length - 1 >> 0; tokens = _tmp; length = _tmp$1;
		_ref = tokens;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			o = (x = length - i >> 0, ((x < 0 || x >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x]));
			if (o.$length > 1) {
				$copy(((i < 0 || i >= tokens.$length) ? $throwRuntimeError("index out of range") : tokens.$array[tokens.$offset + i]), (x$1 = o.$length - 2 >> 0, ((x$1 < 0 || x$1 >= o.$length) ? $throwRuntimeError("index out of range") : o.$array[o.$offset + x$1])).getToken32(), token32);
			}
			_i++;
		}
		return tokens;
	};
	tokens16.prototype.Error = function() { return this.$val.Error(); };
	token32.Ptr.prototype.isParentOf = function(u) {
		var t;
		t = this;
		return t.begin <= u.begin && t.end >= u.end && t.next > u.next;
	};
	token32.prototype.isParentOf = function(u) { return this.$val.isParentOf(u); };
	token32.Ptr.prototype.getToken32 = function() {
		var t;
		t = this;
		return new token32.Ptr(t.pegRule, t.begin, t.end, t.next);
	};
	token32.prototype.getToken32 = function() { return this.$val.getToken32(); };
	token32.Ptr.prototype.String = function() {
		var t, x;
		t = this;
		return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", new ($sliceType($emptyInterface))([new $String((x = t.pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x]))), new $Int32(t.begin), new $Int32(t.end), new $Int32(t.next)]));
	};
	token32.prototype.String = function() { return this.$val.String(); };
	tokens32.Ptr.prototype.trim = function(length) {
		var t;
		t = this;
		t.tree = $subslice(t.tree, 0, length);
	};
	tokens32.prototype.trim = function(length) { return this.$val.trim(length); };
	tokens32.Ptr.prototype.Print = function() {
		var t, _ref, _i, token;
		t = this;
		_ref = t.tree;
		_i = 0;
		while (_i < _ref.$length) {
			token = new token32.Ptr(); $copy(token, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token32);
			fmt.Println(new ($sliceType($emptyInterface))([new $String(token.String())]));
			_i++;
		}
	};
	tokens32.prototype.Print = function() { return this.$val.Print(); };
	tokens32.Ptr.prototype.Order = function() {
		var t, depths, _ref, _i, i, token, depth, length, _lhs, _index, _tmp, _tmp$1, ordered, pool, _ref$1, _i$1, i$1, depth$1, _tmp$2, _tmp$3, _tmp$4, _ref$2, _i$2, i$2, token$1, depth$2, x, x$1, _lhs$1, _index$1;
		t = this;
		if (!(t.ordered === ($sliceType(($sliceType(token32)))).nil)) {
			return t.ordered;
		}
		depths = ($sliceType($Int32)).make(1, 32767);
		_ref = t.tree;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			token = new token32.Ptr(); $copy(token, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token32);
			if (token.pegRule === 0) {
				t.tree = $subslice(t.tree, 0, i);
				break;
			}
			depth = (token.next >> 0);
			length = depths.$length;
			if (depth >= length) {
				depths = $subslice(depths, 0, (depth + 1 >> 0));
			}
			_lhs = depths; _index = depth; (_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index] = ((_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index]) + (1) >> 0;
			_i++;
		}
		depths = $append(depths, 0);
		_tmp = ($sliceType(($sliceType(token32)))).make(depths.$length); _tmp$1 = ($sliceType(token32)).make((t.tree.$length + depths.$length >> 0)); ordered = _tmp; pool = _tmp$1;
		_ref$1 = depths;
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			i$1 = _i$1;
			depth$1 = ((_i$1 < 0 || _i$1 >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i$1]);
			depth$1 = depth$1 + (1) >> 0;
			_tmp$2 = $subslice(pool, 0, depth$1); _tmp$3 = $subslice(pool, depth$1); _tmp$4 = 0; (i$1 < 0 || i$1 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i$1] = _tmp$2; pool = _tmp$3; (i$1 < 0 || i$1 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + i$1] = _tmp$4;
			_i$1++;
		}
		_ref$2 = t.tree;
		_i$2 = 0;
		while (_i$2 < _ref$2.$length) {
			i$2 = _i$2;
			token$1 = new token32.Ptr(); $copy(token$1, ((_i$2 < 0 || _i$2 >= _ref$2.$length) ? $throwRuntimeError("index out of range") : _ref$2.$array[_ref$2.$offset + _i$2]), token32);
			depth$2 = token$1.next;
			token$1.next = (i$2 >> 0);
			$copy((x = ((depth$2 < 0 || depth$2 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth$2]), x$1 = ((depth$2 < 0 || depth$2 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth$2]), ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1])), token$1, token32);
			_lhs$1 = depths; _index$1 = depth$2; (_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1] = ((_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1]) + (1) >> 0;
			_i$2++;
		}
		t.ordered = ordered;
		return ordered;
	};
	tokens32.prototype.Order = function() { return this.$val.Order(); };
	tokens32.Ptr.prototype.AST = function($b) {
		var $this = this, $args = arguments, $r, $s = 0, t, tokens, _r, stack, _ref, _ok, _tuple, _r$1, token, node;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		t = $this;
		tokens = t.Tokens();
		_r = $recv(tokens, true); /* */ $s = 1; case 1: if (_r && _r.constructor === Function) { _r = _r(); }
		stack = new element.Ptr(new node32.Ptr(_r[0], ($ptrType(node32)).nil, ($ptrType(node32)).nil), ($ptrType(element)).nil);
		_ref = tokens;
		/* while (true) { */ case 2: if(!(true)) { $s = 3; continue; }
			_r$1 = $recv(_ref, true); /* */ $s = 4; case 4: if (_r$1 && _r$1.constructor === Function) { _r$1 = _r$1(); }
			_tuple = _r$1; token = new token32.Ptr(); $copy(token, _tuple[0], token32); _ok = _tuple[1];
			if (!_ok) {
				/* break; */ $s = 3; continue;
			}
			if (token.begin === token.end) {
				/* continue; */ $s = 2; continue;
			}
			node = new node32.Ptr(token, ($ptrType(node32)).nil, ($ptrType(node32)).nil);
			while (!(stack === ($ptrType(element)).nil) && stack.node.token32.begin >= token.begin && stack.node.token32.end <= token.end) {
				stack.node.next = node.up;
				node.up = stack.node;
				stack = stack.down;
			}
			stack = new element.Ptr(node, stack);
		/* } */ $s = 2; continue; case 3:
		return stack.node;
		/* */ case -1: } return; } };
	};
	tokens32.prototype.AST = function($b) { return this.$val.AST($b); };
	tokens32.Ptr.prototype.PreOrder = function() {
		var t, _tmp, _tmp$1, s, ordered;
		t = this;
		_tmp = new ($chanType(state32, false, false))(6); _tmp$1 = t.Order(); s = _tmp; ordered = _tmp$1;
		$go((function() {
			var states, _ref, _i, i, _tmp$2, _tmp$3, _tmp$4, depths, state, depth, write, x, _lhs, _index, _tmp$5, x$1, x$2, x$3, x$4, _tmp$6, x$5, x$6, a, b, i$1, _tmp$7, x$7, x$8, _tmp$8, x$9, c, j, x$10, x$11, x$12, next, x$13, x$14, c$1, _lhs$1, _index$1, _tmp$9, _tmp$10, _tmp$11, _lhs$2, _index$2, _tmp$12, x$15, x$16, _tmp$13, c$2, parent, _tmp$14, x$17, x$18, x$19, x$20, _tmp$15, _tmp$16, x$21, x$22;
			states = ($arrayType(state32, 8)).zero(); $copy(states, ($arrayType(state32, 8)).zero(), ($arrayType(state32, 8)));
			_ref = states;
			_i = 0;
			while (_i < 8) {
				i = _i;
				((i < 0 || i >= states.length) ? $throwRuntimeError("index out of range") : states[i]).depths = ($sliceType($Int32)).make(ordered.$length);
				_i++;
			}
			_tmp$2 = ($sliceType($Int32)).make(ordered.$length); _tmp$3 = 0; _tmp$4 = 1; depths = _tmp$2; state = _tmp$3; depth = _tmp$4;
			write = (function(t$1, leaf, $b) {
				var $this = this, $args = arguments, $r, $s = 0, S, _tmp$5, _r, _tmp$6, _tmp$7, _tmp$8, _tmp$9, _tmp$10;
				/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
				S = new state32.Ptr(); $copy(S, ((state < 0 || state >= states.length) ? $throwRuntimeError("index out of range") : states[state]), state32);
				_tmp$5 = (_r = ((state + 1 >> 0)) % 8, _r === _r ? _r : $throwRuntimeError("integer divide by zero")); _tmp$6 = t$1.pegRule; _tmp$7 = t$1.begin; _tmp$8 = t$1.end; _tmp$9 = (depth >> 0); _tmp$10 = leaf; state = _tmp$5; S.token32.pegRule = _tmp$6; S.token32.begin = _tmp$7; S.token32.end = _tmp$8; S.token32.next = _tmp$9; S.leaf = _tmp$10;
				$copySlice(S.depths, depths);
				$r = $send(s, $clone(S, state32), true); /* */ $s = 1; case 1: if ($r && $r.constructor === Function) { $r = $r(); }
				/* */ case -1: } return; } };
			});
			$copy(((state < 0 || state >= states.length) ? $throwRuntimeError("index out of range") : states[state]).token32, (x = ((0 < 0 || 0 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + 0]), ((0 < 0 || 0 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + 0])), token32);
			_lhs = depths; _index = 0; (_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index] = ((_index < 0 || _index >= _lhs.$length) ? $throwRuntimeError("index out of range") : _lhs.$array[_lhs.$offset + _index]) + (1) >> 0;
			state = state + (1) >> 0;
			_tmp$5 = new token32.Ptr(); $copy(_tmp$5, (x$1 = (x$2 = depth - 1 >> 0, ((x$2 < 0 || x$2 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x$2])), x$3 = (x$4 = depth - 1 >> 0, ((x$4 < 0 || x$4 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + x$4])) - 1 >> 0, ((x$3 < 0 || x$3 >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + x$3])), token32); _tmp$6 = new token32.Ptr(); $copy(_tmp$6, (x$5 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$6 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]), ((x$6 < 0 || x$6 >= x$5.$length) ? $throwRuntimeError("index out of range") : x$5.$array[x$5.$offset + x$6])), token32); a = new token32.Ptr(); $copy(a, _tmp$5, token32); b = new token32.Ptr(); $copy(b, _tmp$6, token32);
			depthFirstSearch:
			while (true) {
				while (true) {
					i$1 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]);
					if (i$1 > 0) {
						_tmp$7 = new token32.Ptr(); $copy(_tmp$7, (x$7 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$8 = i$1 - 1 >> 0, ((x$8 < 0 || x$8 >= x$7.$length) ? $throwRuntimeError("index out of range") : x$7.$array[x$7.$offset + x$8])), token32); _tmp$8 = (x$9 = depth - 1 >> 0, ((x$9 < 0 || x$9 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + x$9])); c = new token32.Ptr(); $copy(c, _tmp$7, token32); j = _tmp$8;
						if (a.isParentOf($clone(c, token32)) && (j < 2 || !(x$10 = (x$11 = depth - 1 >> 0, ((x$11 < 0 || x$11 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x$11])), x$12 = j - 2 >> 0, ((x$12 < 0 || x$12 >= x$10.$length) ? $throwRuntimeError("index out of range") : x$10.$array[x$10.$offset + x$12])).isParentOf($clone(c, token32)))) {
							if (!((c.end === b.begin))) {
								write(new token32.Ptr(100, c.end, b.begin, 0), true);
							}
							break;
						}
					}
					if (a.begin < b.begin) {
						write(new token32.Ptr(99, a.begin, b.begin, 0), true);
					}
					break;
				}
				next = depth + 1 >> 0;
				c$1 = new token32.Ptr(); $copy(c$1, (x$13 = ((next < 0 || next >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + next]), x$14 = ((next < 0 || next >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + next]), ((x$14 < 0 || x$14 >= x$13.$length) ? $throwRuntimeError("index out of range") : x$13.$array[x$13.$offset + x$14])), token32);
				if (!((c$1.pegRule === 0)) && b.isParentOf($clone(c$1, token32))) {
					write($clone(b, token32), false);
					_lhs$1 = depths; _index$1 = depth; (_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1] = ((_index$1 < 0 || _index$1 >= _lhs$1.$length) ? $throwRuntimeError("index out of range") : _lhs$1.$array[_lhs$1.$offset + _index$1]) + (1) >> 0;
					_tmp$9 = next; _tmp$10 = new token32.Ptr(); $copy(_tmp$10, b, token32); _tmp$11 = new token32.Ptr(); $copy(_tmp$11, c$1, token32); depth = _tmp$9; $copy(a, _tmp$10, token32); $copy(b, _tmp$11, token32);
					continue;
				}
				write($clone(b, token32), true);
				_lhs$2 = depths; _index$2 = depth; (_index$2 < 0 || _index$2 >= _lhs$2.$length) ? $throwRuntimeError("index out of range") : _lhs$2.$array[_lhs$2.$offset + _index$2] = ((_index$2 < 0 || _index$2 >= _lhs$2.$length) ? $throwRuntimeError("index out of range") : _lhs$2.$array[_lhs$2.$offset + _index$2]) + (1) >> 0;
				_tmp$12 = new token32.Ptr(); $copy(_tmp$12, (x$15 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$16 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]), ((x$16 < 0 || x$16 >= x$15.$length) ? $throwRuntimeError("index out of range") : x$15.$array[x$15.$offset + x$16])), token32); _tmp$13 = true; c$2 = new token32.Ptr(); $copy(c$2, _tmp$12, token32); parent = _tmp$13;
				while (true) {
					if (!((c$2.pegRule === 0)) && a.isParentOf($clone(c$2, token32))) {
						$copy(b, c$2, token32);
						continue depthFirstSearch;
					} else if (parent && !((b.end === a.end))) {
						write(new token32.Ptr(101, b.end, a.end, 0), true);
					}
					depth = depth - (1) >> 0;
					if (depth > 0) {
						_tmp$14 = new token32.Ptr(); $copy(_tmp$14, (x$17 = (x$18 = depth - 1 >> 0, ((x$18 < 0 || x$18 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x$18])), x$19 = (x$20 = depth - 1 >> 0, ((x$20 < 0 || x$20 >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + x$20])) - 1 >> 0, ((x$19 < 0 || x$19 >= x$17.$length) ? $throwRuntimeError("index out of range") : x$17.$array[x$17.$offset + x$19])), token32); _tmp$15 = new token32.Ptr(); $copy(_tmp$15, a, token32); _tmp$16 = new token32.Ptr(); $copy(_tmp$16, (x$21 = ((depth < 0 || depth >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + depth]), x$22 = ((depth < 0 || depth >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + depth]), ((x$22 < 0 || x$22 >= x$21.$length) ? $throwRuntimeError("index out of range") : x$21.$array[x$21.$offset + x$22])), token32); $copy(a, _tmp$14, token32); $copy(b, _tmp$15, token32); $copy(c$2, _tmp$16, token32);
						parent = a.isParentOf($clone(b, token32));
						continue;
					}
					break depthFirstSearch;
				}
			}
			$close(s);
		}), []);
		return [s, ordered];
	};
	tokens32.prototype.PreOrder = function() { return this.$val.PreOrder(); };
	tokens32.Ptr.prototype.PrintSyntax = function($b) {
		var $this = this, $args = arguments, $r, $s = 0, t, _tuple, tokens, ordered, max, _ref, _ok, _tuple$1, _r, token, _tmp, _tmp$1, _tmp$2, i, leaf, depths, x, x$1, x$2, x$3, _tmp$3, _tmp$4, _tmp$5, i$1, leaf$1, depths$1, x$4, x$5, x$6, x$7, _tmp$6, _tmp$7, c, end, i$2, j, i$3, j$1, _tmp$8, _tmp$9, _tmp$10, i$4, leaf$2, depths$2, x$8, x$9, x$10, x$11;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		t = $this;
		_tuple = t.PreOrder(); tokens = _tuple[0]; ordered = _tuple[1];
		max = -1;
		_ref = tokens;
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			_r = $recv(_ref, true); /* */ $s = 3; case 3: if (_r && _r.constructor === Function) { _r = _r(); }
			_tuple$1 = _r; token = new state32.Ptr(); $copy(token, _tuple$1[0], state32); _ok = _tuple$1[1];
			if (!_ok) {
				/* break; */ $s = 2; continue;
			}
			if (!token.leaf) {
				fmt.Printf("%v", new ($sliceType($emptyInterface))([new $Int32(token.token32.begin)]));
				_tmp = 0; _tmp$1 = (token.token32.next >> 0); _tmp$2 = token.depths; i = _tmp; leaf = _tmp$1; depths = _tmp$2;
				while (i < leaf) {
					fmt.Printf(" \x1B[36m%v\x1B[m", new ($sliceType($emptyInterface))([new $String((x = (x$1 = ((i < 0 || i >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i]), x$2 = ((i < 0 || i >= depths.$length) ? $throwRuntimeError("index out of range") : depths.$array[depths.$offset + i]) - 1 >> 0, ((x$2 < 0 || x$2 >= x$1.$length) ? $throwRuntimeError("index out of range") : x$1.$array[x$1.$offset + x$2])).pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x])))]));
					i = i + (1) >> 0;
				}
				fmt.Printf(" \x1B[36m%v\x1B[m\n", new ($sliceType($emptyInterface))([new $String((x$3 = token.token32.pegRule, ((x$3 < 0 || x$3 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$3])))]));
			} else if (token.token32.begin === token.token32.end) {
				fmt.Printf("%v", new ($sliceType($emptyInterface))([new $Int32(token.token32.begin)]));
				_tmp$3 = 0; _tmp$4 = (token.token32.next >> 0); _tmp$5 = token.depths; i$1 = _tmp$3; leaf$1 = _tmp$4; depths$1 = _tmp$5;
				while (i$1 < leaf$1) {
					fmt.Printf(" \x1B[31m%v\x1B[m", new ($sliceType($emptyInterface))([new $String((x$4 = (x$5 = ((i$1 < 0 || i$1 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i$1]), x$6 = ((i$1 < 0 || i$1 >= depths$1.$length) ? $throwRuntimeError("index out of range") : depths$1.$array[depths$1.$offset + i$1]) - 1 >> 0, ((x$6 < 0 || x$6 >= x$5.$length) ? $throwRuntimeError("index out of range") : x$5.$array[x$5.$offset + x$6])).pegRule, ((x$4 < 0 || x$4 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$4])))]));
					i$1 = i$1 + (1) >> 0;
				}
				fmt.Printf(" \x1B[31m%v\x1B[m\n", new ($sliceType($emptyInterface))([new $String((x$7 = token.token32.pegRule, ((x$7 < 0 || x$7 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$7])))]));
			} else {
				_tmp$6 = token.token32.begin; _tmp$7 = token.token32.end; c = _tmp$6; end = _tmp$7;
				while (c < end) {
					i$2 = (c >> 0);
					if ((max + 1 >> 0) < i$2) {
						j = max;
						while (j < i$2) {
							fmt.Printf("skip %v %v\n", new ($sliceType($emptyInterface))([new $Int(j), new $String(token.token32.String())]));
							j = j + (1) >> 0;
						}
						max = i$2;
					} else {
						i$3 = (c >> 0);
						if (i$3 <= max) {
							j$1 = i$3;
							while (j$1 <= max) {
								fmt.Printf("dupe %v %v\n", new ($sliceType($emptyInterface))([new $Int(j$1), new $String(token.token32.String())]));
								j$1 = j$1 + (1) >> 0;
							}
						} else {
							max = (c >> 0);
						}
					}
					fmt.Printf("%v", new ($sliceType($emptyInterface))([new $Int32(c)]));
					_tmp$8 = 0; _tmp$9 = (token.token32.next >> 0); _tmp$10 = token.depths; i$4 = _tmp$8; leaf$2 = _tmp$9; depths$2 = _tmp$10;
					while (i$4 < leaf$2) {
						fmt.Printf(" \x1B[34m%v\x1B[m", new ($sliceType($emptyInterface))([new $String((x$8 = (x$9 = ((i$4 < 0 || i$4 >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + i$4]), x$10 = ((i$4 < 0 || i$4 >= depths$2.$length) ? $throwRuntimeError("index out of range") : depths$2.$array[depths$2.$offset + i$4]) - 1 >> 0, ((x$10 < 0 || x$10 >= x$9.$length) ? $throwRuntimeError("index out of range") : x$9.$array[x$9.$offset + x$10])).pegRule, ((x$8 < 0 || x$8 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$8])))]));
						i$4 = i$4 + (1) >> 0;
					}
					fmt.Printf(" \x1B[34m%v\x1B[m\n", new ($sliceType($emptyInterface))([new $String((x$11 = token.token32.pegRule, ((x$11 < 0 || x$11 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$11])))]));
					c = c + (1) >> 0;
				}
				fmt.Printf("\n", new ($sliceType($emptyInterface))([]));
			}
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } };
	};
	tokens32.prototype.PrintSyntax = function($b) { return this.$val.PrintSyntax($b); };
	tokens32.Ptr.prototype.PrintSyntaxTree = function(buffer, $b) {
		var $this = this, $args = arguments, $r, $s = 0, t, _tuple, tokens, _ref, _ok, _tuple$1, _r, token, c, x;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		t = $this;
		_tuple = t.PreOrder(); tokens = _tuple[0];
		_ref = tokens;
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			_r = $recv(_ref, true); /* */ $s = 3; case 3: if (_r && _r.constructor === Function) { _r = _r(); }
			_tuple$1 = _r; token = new state32.Ptr(); $copy(token, _tuple$1[0], state32); _ok = _tuple$1[1];
			if (!_ok) {
				/* break; */ $s = 2; continue;
			}
			c = 0;
			while (c < (token.token32.next >> 0)) {
				fmt.Printf(" ", new ($sliceType($emptyInterface))([]));
				c = c + (1) >> 0;
			}
			fmt.Printf("\x1B[34m%v\x1B[m %v\n", new ($sliceType($emptyInterface))([new $String((x = token.token32.pegRule, ((x < 0 || x >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x]))), new $String(strconv.Quote(buffer.substring(token.token32.begin, token.token32.end)))]));
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } };
	};
	tokens32.prototype.PrintSyntaxTree = function(buffer, $b) { return this.$val.PrintSyntaxTree(buffer, $b); };
	tokens32.Ptr.prototype.Add = function(rule, begin, end, depth, index) {
		var t, x;
		t = this;
		$copy((x = t.tree, ((index < 0 || index >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + index])), new token32.Ptr(rule, (begin >> 0), (end >> 0), (depth >> 0)), token32);
	};
	tokens32.prototype.Add = function(rule, begin, end, depth, index) { return this.$val.Add(rule, begin, end, depth, index); };
	tokens32.Ptr.prototype.Tokens = function() {
		var t, s;
		t = this;
		s = new ($chanType(token32, false, false))(16);
		$go((function($b) {
			var $this = this, $args = arguments, $r, $s = 0, _ref, _i, v;
			/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
			_ref = t.tree;
			_i = 0;
			/* while (_i < _ref.$length) { */ case 1: if(!(_i < _ref.$length)) { $s = 2; continue; }
				v = new token32.Ptr(); $copy(v, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token32);
				$r = $send(s, $clone(v.getToken32(), token32), true); /* */ $s = 3; case 3: if ($r && $r.constructor === Function) { $r = $r(); }
				_i++;
			/* } */ $s = 1; continue; case 2:
			$close(s);
			/* */ case -1: } return; } };
		}), []);
		return s;
	};
	tokens32.prototype.Tokens = function() { return this.$val.Tokens(); };
	tokens32.Ptr.prototype.Error = function() {
		var t, ordered, length, _tmp, _tmp$1, tokens, _ref, _i, i, x, o, x$1;
		t = this;
		ordered = t.Order();
		length = ordered.$length;
		_tmp = ($sliceType(token32)).make(length); _tmp$1 = length - 1 >> 0; tokens = _tmp; length = _tmp$1;
		_ref = tokens;
		_i = 0;
		while (_i < _ref.$length) {
			i = _i;
			o = (x = length - i >> 0, ((x < 0 || x >= ordered.$length) ? $throwRuntimeError("index out of range") : ordered.$array[ordered.$offset + x]));
			if (o.$length > 1) {
				$copy(((i < 0 || i >= tokens.$length) ? $throwRuntimeError("index out of range") : tokens.$array[tokens.$offset + i]), (x$1 = o.$length - 2 >> 0, ((x$1 < 0 || x$1 >= o.$length) ? $throwRuntimeError("index out of range") : o.$array[o.$offset + x$1])).getToken32(), token32);
			}
			_i++;
		}
		return tokens;
	};
	tokens32.prototype.Error = function() { return this.$val.Error(); };
	tokens16.Ptr.prototype.Expand = function(index) {
		var t, tree, x, expanded, _ref, _i, i, v;
		t = this;
		tree = t.tree;
		if (index >= tree.$length) {
			expanded = ($sliceType(token32)).make((x = tree.$length, (((2 >>> 16 << 16) * x >> 0) + (2 << 16 >>> 16) * x) >> 0));
			_ref = tree;
			_i = 0;
			while (_i < _ref.$length) {
				i = _i;
				v = new token16.Ptr(); $copy(v, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token16);
				$copy(((i < 0 || i >= expanded.$length) ? $throwRuntimeError("index out of range") : expanded.$array[expanded.$offset + i]), v.getToken32(), token32);
				_i++;
			}
			return new tokens32.Ptr(expanded, ($sliceType(($sliceType(token32)))).nil);
		}
		return null;
	};
	tokens16.prototype.Expand = function(index) { return this.$val.Expand(index); };
	tokens32.Ptr.prototype.Expand = function(index) {
		var t, tree, x, expanded;
		t = this;
		tree = t.tree;
		if (index >= tree.$length) {
			expanded = ($sliceType(token32)).make((x = tree.$length, (((2 >>> 16 << 16) * x >> 0) + (2 << 16 >>> 16) * x) >> 0));
			$copySlice(expanded, tree);
			t.tree = expanded;
		}
		return null;
	};
	tokens32.prototype.Expand = function(index) { return this.$val.Expand(index); };
	translatePositions = function(buffer, positions) {
		var _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, length, translations, j, line, symbol, _ref, _i, _rune, i, c, _tmp$5, _tmp$6, _key;
		_tmp = positions.$length; _tmp$1 = new $Map(); _tmp$2 = 0; _tmp$3 = 1; _tmp$4 = 0; length = _tmp; translations = _tmp$1; j = _tmp$2; line = _tmp$3; symbol = _tmp$4;
		sort.Ints(positions);
		_ref = buffer.substring(0);
		_i = 0;
		search:
		while (_i < _ref.length) {
			_rune = $decodeRune(_ref, _i);
			i = _i;
			c = _rune[0];
			if (c === 10) {
				_tmp$5 = line + 1 >> 0; _tmp$6 = 0; line = _tmp$5; symbol = _tmp$6;
			} else {
				symbol = symbol + (1) >> 0;
			}
			if (i === ((j < 0 || j >= positions.$length) ? $throwRuntimeError("index out of range") : positions.$array[positions.$offset + j])) {
				_key = ((j < 0 || j >= positions.$length) ? $throwRuntimeError("index out of range") : positions.$array[positions.$offset + j]); (translations || $throwRuntimeError("assignment to entry in nil map"))[_key] = { k: _key, v: new textPosition.Ptr(line, symbol) };
				j = j + (1) >> 0;
				while (j < length) {
					if (!((i === ((j < 0 || j >= positions.$length) ? $throwRuntimeError("index out of range") : positions.$array[positions.$offset + j])))) {
						_i += _rune[1];
						continue search;
					}
					j = j + (1) >> 0;
				}
				break search;
			}
			_i += _rune[1];
		}
		return translations;
	};
	parseError.Ptr.prototype.Error = function() {
		var e, _tmp, _tmp$1, tokens, error, _tmp$2, x, _tmp$3, positions, p, _ref, _i, token, _tmp$4, _tmp$5, _tmp$6, _tmp$7, translations, _ref$1, _i$1, token$1, _tmp$8, _tmp$9, begin, end, x$1, _entry, _entry$1, _entry$2, _entry$3;
		e = this;
		_tmp = e.p.tokenTree.Error(); _tmp$1 = "\n"; tokens = _tmp; error = _tmp$1;
		_tmp$2 = ($sliceType($Int)).make((x = tokens.$length, (((2 >>> 16 << 16) * x >> 0) + (2 << 16 >>> 16) * x) >> 0)); _tmp$3 = 0; positions = _tmp$2; p = _tmp$3;
		_ref = tokens;
		_i = 0;
		while (_i < _ref.$length) {
			token = new token32.Ptr(); $copy(token, ((_i < 0 || _i >= _ref.$length) ? $throwRuntimeError("index out of range") : _ref.$array[_ref.$offset + _i]), token32);
			_tmp$4 = (token.begin >> 0); _tmp$5 = p + 1 >> 0; (p < 0 || p >= positions.$length) ? $throwRuntimeError("index out of range") : positions.$array[positions.$offset + p] = _tmp$4; p = _tmp$5;
			_tmp$6 = (token.end >> 0); _tmp$7 = p + 1 >> 0; (p < 0 || p >= positions.$length) ? $throwRuntimeError("index out of range") : positions.$array[positions.$offset + p] = _tmp$6; p = _tmp$7;
			_i++;
		}
		translations = translatePositions(e.p.Buffer, positions);
		_ref$1 = tokens;
		_i$1 = 0;
		while (_i$1 < _ref$1.$length) {
			token$1 = new token32.Ptr(); $copy(token$1, ((_i$1 < 0 || _i$1 >= _ref$1.$length) ? $throwRuntimeError("index out of range") : _ref$1.$array[_ref$1.$offset + _i$1]), token32);
			_tmp$8 = (token$1.begin >> 0); _tmp$9 = (token$1.end >> 0); begin = _tmp$8; end = _tmp$9;
			error = error + (fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n", new ($sliceType($emptyInterface))([new $String((x$1 = token$1.pegRule, ((x$1 < 0 || x$1 >= rul3s.length) ? $throwRuntimeError("index out of range") : rul3s[x$1]))), new $Int((_entry = translations[begin], _entry !== undefined ? _entry.v : new textPosition.Ptr()).line), new $Int((_entry$1 = translations[begin], _entry$1 !== undefined ? _entry$1.v : new textPosition.Ptr()).symbol), new $Int((_entry$2 = translations[end], _entry$2 !== undefined ? _entry$2.v : new textPosition.Ptr()).line), new $Int((_entry$3 = translations[end], _entry$3 !== undefined ? _entry$3.v : new textPosition.Ptr()).symbol), new $String(e.p.Buffer.substring(begin, end))])));
			_i$1++;
		}
		return error;
	};
	parseError.prototype.Error = function() { return this.$val.Error(); };
	Sparql.Ptr.prototype.PrintSyntaxTree = function() {
		var p;
		p = this;
		p.tokenTree.PrintSyntaxTree(p.Buffer);
	};
	Sparql.prototype.PrintSyntaxTree = function() { return this.$val.PrintSyntaxTree(); };
	Sparql.Ptr.prototype.Highlighter = function() {
		var p;
		p = this;
		p.tokenTree.PrintSyntax();
	};
	Sparql.prototype.Highlighter = function() { return this.$val.Highlighter(); };
	Sparql.Ptr.prototype.Execute = function($b) {
		var $this = this, $args = arguments, $r, $s = 0, p, _tmp, _tmp$1, _tmp$2, buffer, begin, end, _ref, _ok, _tuple, _r, token, _ref$1, _tmp$3, _tmp$4;
		/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
		p = $this;
		_tmp = p.Buffer; _tmp$1 = 0; _tmp$2 = 0; buffer = _tmp; begin = _tmp$1; end = _tmp$2;
		_ref = p.tokenTree.Tokens();
		/* while (true) { */ case 1: if(!(true)) { $s = 2; continue; }
			_r = $recv(_ref, true); /* */ $s = 3; case 3: if (_r && _r.constructor === Function) { _r = _r(); }
			_tuple = _r; token = new token32.Ptr(); $copy(token, _tuple[0], token32); _ok = _tuple[1];
			if (!_ok) {
				/* break; */ $s = 2; continue;
			}
			_ref$1 = token.pegRule;
			if (_ref$1 === 89) {
				_tmp$3 = (token.begin >> 0); _tmp$4 = (token.end >> 0); begin = _tmp$3; end = _tmp$4;
			} else if (_ref$1 === 90) {
				p.Bgp.setSubject(buffer.substring(begin, end));
			} else if (_ref$1 === 91) {
				p.Bgp.setSubject(buffer.substring(begin, end));
			} else if (_ref$1 === 92) {
				p.Bgp.setSubject("?POF");
			} else if (_ref$1 === 93) {
				p.Bgp.setPredicate("?POF");
			} else if (_ref$1 === 94) {
				p.Bgp.setPredicate(buffer.substring(begin, end));
			} else if (_ref$1 === 95) {
				p.Bgp.setPredicate(buffer.substring(begin, end));
			} else if (_ref$1 === 96) {
				p.Bgp.setObject(buffer.substring(begin, end));
				p.Bgp.addTriplePattern();
			} else if (_ref$1 === 97) {
				p.Bgp.setObject("?POF");
				p.Bgp.addTriplePattern();
			} else if (_ref$1 === 98) {
				p.Bgp.setObject("?FillVar");
				p.Bgp.addTriplePattern();
			}
		/* } */ $s = 1; continue; case 2:
		/* */ case -1: } return; } };
	};
	Sparql.prototype.Execute = function($b) { return this.$val.Execute($b); };
	Sparql.Ptr.prototype.Init = function() {
		var p, x, x$1, tree, _tmp, _tmp$1, _tmp$2, _tmp$3, _tmp$4, position, depth, tokenIndex, buffer, rules, add, matchDot;
		p = this;
		p.buffer = new ($sliceType($Int32))($stringToRunes(p.Buffer));
		if ((p.buffer.$length === 0) || !(((x = p.buffer, x$1 = p.buffer.$length - 1 >> 0, ((x$1 < 0 || x$1 >= x.$length) ? $throwRuntimeError("index out of range") : x.$array[x.$offset + x$1])) === 4))) {
			p.buffer = $append(p.buffer, 4);
		}
		tree = new tokens16.Ptr(($sliceType(token16)).make(32767), ($sliceType(($sliceType(token16)))).nil);
		_tmp = 0; _tmp$1 = 0; _tmp$2 = 0; _tmp$3 = p.buffer; _tmp$4 = ($arrayType(($funcType([], [$Bool], false)), 99)).zero(); $copy(_tmp$4, p.rules, ($arrayType(($funcType([], [$Bool], false)), 99))); position = _tmp; depth = _tmp$1; tokenIndex = _tmp$2; buffer = _tmp$3; rules = ($arrayType(($funcType([], [$Bool], false)), 99)).zero(); $copy(rules, _tmp$4, ($arrayType(($funcType([], [$Bool], false)), 99)));
		p.Parse = (function(rule) {
			var r, x$2, matches;
			r = 1;
			if (rule.$length > 0) {
				r = ((0 < 0 || 0 >= rule.$length) ? $throwRuntimeError("index out of range") : rule.$array[rule.$offset + 0]);
			}
			matches = (x$2 = p.rules, ((r < 0 || r >= x$2.length) ? $throwRuntimeError("index out of range") : x$2[r]))();
			p.tokenTree = tree;
			if (matches) {
				p.tokenTree.trim(tokenIndex);
				return null;
			}
			return new parseError.Ptr(p);
		});
		p.Reset = (function() {
			var _tmp$5, _tmp$6, _tmp$7;
			_tmp$5 = 0; _tmp$6 = 0; _tmp$7 = 0; position = _tmp$5; tokenIndex = _tmp$6; depth = _tmp$7;
		});
		add = (function(rule, begin) {
			var t;
			t = tree.Expand(tokenIndex);
			if (!($interfaceIsEqual(t, null))) {
				tree = t;
			}
			tree.Add(rule, begin, position, depth, tokenIndex);
			tokenIndex = tokenIndex + (1) >> 0;
		});
		matchDot = (function() {
			if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 4))) {
				position = position + (1) >> 0;
				return true;
			}
			return false;
		});
		$copy(rules, $toNativeArray("Func", [$throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position0, tokenIndex0, depth0, position1, position2, _tmp$8, _tmp$9, _tmp$10, position4, tokenIndex4, depth4, _tmp$11, _tmp$12, _tmp$13, position5, tokenIndex5, depth5, position7, position8, _tmp$14, _tmp$15, _tmp$16, position9, tokenIndex9, depth9, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, position11, tokenIndex11, depth11, _tmp$23, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, position13, tokenIndex13, depth13, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, position15, tokenIndex15, depth15, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40, position17, tokenIndex17, depth17, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46, position19, tokenIndex19, depth19, _tmp$47, _tmp$48, _tmp$49, _tmp$50, _tmp$51, _tmp$52, position23, tokenIndex23, depth23, _tmp$53, _tmp$54, _tmp$55, position24, tokenIndex24, depth24, _tmp$56, _tmp$57, _tmp$58, _tmp$59, _tmp$60, _tmp$61, _tmp$62, _tmp$63, _tmp$64, position22, tokenIndex22, depth22, _tmp$65, _tmp$66, _tmp$67, position26, tokenIndex26, depth26, _tmp$68, _tmp$69, _tmp$70, position27, tokenIndex27, depth27, _tmp$71, _tmp$72, _tmp$73, _tmp$74, _tmp$75, _tmp$76, _tmp$77, _tmp$78, _tmp$79, position29, _tmp$80, _tmp$81, _tmp$82, position30, position31, _tmp$83, _tmp$84, _tmp$85, position32, tokenIndex32, depth32, _tmp$86, _tmp$87, _tmp$88, _tmp$89, _tmp$90, _tmp$91, position34, tokenIndex34, depth34, _tmp$92, _tmp$93, _tmp$94, _tmp$95, _tmp$96, _tmp$97, position36, tokenIndex36, depth36, _tmp$98, _tmp$99, _tmp$100, _tmp$101, _tmp$102, _tmp$103, position38, tokenIndex38, depth38, _tmp$104, _tmp$105, _tmp$106, _tmp$107, _tmp$108, _tmp$109, position40, position41, _tmp$110, _tmp$111, _tmp$112, position42, tokenIndex42, depth42, position44, position45, _tmp$113, _tmp$114, _tmp$115, position46, tokenIndex46, depth46, _tmp$116, _tmp$117, _tmp$118, _tmp$119, _tmp$120, _tmp$121, position48, tokenIndex48, depth48, _tmp$122, _tmp$123, _tmp$124, _tmp$125, _tmp$126, _tmp$127, position50, tokenIndex50, depth50, _tmp$128, _tmp$129, _tmp$130, _tmp$131, _tmp$132, _tmp$133, position52, tokenIndex52, depth52, _tmp$134, _tmp$135, _tmp$136, _tmp$137, _tmp$138, _tmp$139, position54, tokenIndex54, depth54, position56, _tmp$140, _tmp$141, _tmp$142, position57, tokenIndex57, depth57, _tmp$143, _tmp$144, _tmp$145, _tmp$146, _tmp$147, _tmp$148, position59, tokenIndex59, depth59, _tmp$149, _tmp$150, _tmp$151, _tmp$152, _tmp$153, _tmp$154, position61, tokenIndex61, depth61, _tmp$155, _tmp$156, _tmp$157, _tmp$158, _tmp$159, _tmp$160, position63, tokenIndex63, depth63, _tmp$161, _tmp$162, _tmp$163, _tmp$164, _tmp$165, _tmp$166, position65, tokenIndex65, depth65, _tmp$167, _tmp$168, _tmp$169, _tmp$170, _tmp$171, _tmp$172, _tmp$173, _tmp$174, _tmp$175, position67, _tmp$176, _tmp$177, _tmp$178, position68, tokenIndex68, depth68, position70, _tmp$179, _tmp$180, _tmp$181, position71, tokenIndex71, depth71, _tmp$182, _tmp$183, _tmp$184, position73, tokenIndex73, depth73, _tmp$185, _tmp$186, _tmp$187, _tmp$188, _tmp$189, _tmp$190, _tmp$191, _tmp$192, _tmp$193, position75, tokenIndex75, depth75, _tmp$194, _tmp$195, _tmp$196, _tmp$197, _tmp$198, _tmp$199, _tmp$200, _tmp$201, _tmp$202, position77, tokenIndex77, depth77, _tmp$203, _tmp$204, _tmp$205, _tmp$206, _tmp$207, _tmp$208;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position0 = _tmp$5; tokenIndex0 = _tmp$6; depth0 = _tmp$7;
			position1 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 65; continue; }
				/* goto l0 */ $s = 1; continue;
			/* } */ case 65:
			position2 = position;
			depth = depth + (1) >> 0;
			/* l3: */ case 33:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position4 = _tmp$8; tokenIndex4 = _tmp$9; depth4 = _tmp$10;
			_tmp$11 = position; _tmp$12 = tokenIndex; _tmp$13 = depth; position5 = _tmp$11; tokenIndex5 = _tmp$12; depth5 = _tmp$13;
			position7 = position;
			depth = depth + (1) >> 0;
			position8 = position;
			depth = depth + (1) >> 0;
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position9 = _tmp$14; tokenIndex9 = _tmp$15; depth9 = _tmp$16;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 112))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 112))) {} else { $s = 66; continue; }
				/* goto l10 */ $s = 2; continue;
			/* } */ case 66:
			position = position + (1) >> 0;
			/* goto l9 */ $s = 3; continue;
			/* l10: */ case 2:
			_tmp$17 = position9; _tmp$18 = tokenIndex9; _tmp$19 = depth9; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 80))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 80))) {} else { $s = 67; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 67:
			position = position + (1) >> 0;
			/* l9: */ case 3:
			_tmp$20 = position; _tmp$21 = tokenIndex; _tmp$22 = depth; position11 = _tmp$20; tokenIndex11 = _tmp$21; depth11 = _tmp$22;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) {} else { $s = 68; continue; }
				/* goto l12 */ $s = 5; continue;
			/* } */ case 68:
			position = position + (1) >> 0;
			/* goto l11 */ $s = 6; continue;
			/* l12: */ case 5:
			_tmp$23 = position11; _tmp$24 = tokenIndex11; _tmp$25 = depth11; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) {} else { $s = 69; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 69:
			position = position + (1) >> 0;
			/* l11: */ case 6:
			_tmp$26 = position; _tmp$27 = tokenIndex; _tmp$28 = depth; position13 = _tmp$26; tokenIndex13 = _tmp$27; depth13 = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 70; continue; }
				/* goto l14 */ $s = 7; continue;
			/* } */ case 70:
			position = position + (1) >> 0;
			/* goto l13 */ $s = 8; continue;
			/* l14: */ case 7:
			_tmp$29 = position13; _tmp$30 = tokenIndex13; _tmp$31 = depth13; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 71; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 71:
			position = position + (1) >> 0;
			/* l13: */ case 8:
			_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position15 = _tmp$32; tokenIndex15 = _tmp$33; depth15 = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) {} else { $s = 72; continue; }
				/* goto l16 */ $s = 9; continue;
			/* } */ case 72:
			position = position + (1) >> 0;
			/* goto l15 */ $s = 10; continue;
			/* l16: */ case 9:
			_tmp$35 = position15; _tmp$36 = tokenIndex15; _tmp$37 = depth15; position = _tmp$35; tokenIndex = _tmp$36; depth = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) {} else { $s = 73; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 73:
			position = position + (1) >> 0;
			/* l15: */ case 10:
			_tmp$38 = position; _tmp$39 = tokenIndex; _tmp$40 = depth; position17 = _tmp$38; tokenIndex17 = _tmp$39; depth17 = _tmp$40;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 74; continue; }
				/* goto l18 */ $s = 11; continue;
			/* } */ case 74:
			position = position + (1) >> 0;
			/* goto l17 */ $s = 12; continue;
			/* l18: */ case 11:
			_tmp$41 = position17; _tmp$42 = tokenIndex17; _tmp$43 = depth17; position = _tmp$41; tokenIndex = _tmp$42; depth = _tmp$43;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 75; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 75:
			position = position + (1) >> 0;
			/* l17: */ case 12:
			_tmp$44 = position; _tmp$45 = tokenIndex; _tmp$46 = depth; position19 = _tmp$44; tokenIndex19 = _tmp$45; depth19 = _tmp$46;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 120))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 120))) {} else { $s = 76; continue; }
				/* goto l20 */ $s = 13; continue;
			/* } */ case 76:
			position = position + (1) >> 0;
			/* goto l19 */ $s = 14; continue;
			/* l20: */ case 13:
			_tmp$47 = position19; _tmp$48 = tokenIndex19; _tmp$49 = depth19; position = _tmp$47; tokenIndex = _tmp$48; depth = _tmp$49;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 88))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 88))) {} else { $s = 77; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 77:
			position = position + (1) >> 0;
			/* l19: */ case 14:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 78; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 78:
			depth = depth - (1) >> 0;
			add(55, position8);
			_tmp$50 = position; _tmp$51 = tokenIndex; _tmp$52 = depth; position23 = _tmp$50; tokenIndex23 = _tmp$51; depth23 = _tmp$52;
			_tmp$53 = position; _tmp$54 = tokenIndex; _tmp$55 = depth; position24 = _tmp$53; tokenIndex24 = _tmp$54; depth24 = _tmp$55;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) {} else { $s = 79; continue; }
				/* goto l25 */ $s = 15; continue;
			/* } */ case 79:
			position = position + (1) >> 0;
			/* goto l24 */ $s = 16; continue;
			/* l25: */ case 15:
			_tmp$56 = position24; _tmp$57 = tokenIndex24; _tmp$58 = depth24; position = _tmp$56; tokenIndex = _tmp$57; depth = _tmp$58;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) {} else { $s = 80; continue; }
				/* goto l23 */ $s = 17; continue;
			/* } */ case 80:
			position = position + (1) >> 0;
			/* l24: */ case 16:
			/* goto l6 */ $s = 4; continue;
			/* l23: */ case 17:
			_tmp$59 = position23; _tmp$60 = tokenIndex23; _tmp$61 = depth23; position = _tmp$59; tokenIndex = _tmp$60; depth = _tmp$61;
			/* if (!matchDot()) { */ if (!matchDot()) {} else { $s = 81; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 81:
			/* l21: */ case 22:
			_tmp$62 = position; _tmp$63 = tokenIndex; _tmp$64 = depth; position22 = _tmp$62; tokenIndex22 = _tmp$63; depth22 = _tmp$64;
			_tmp$65 = position; _tmp$66 = tokenIndex; _tmp$67 = depth; position26 = _tmp$65; tokenIndex26 = _tmp$66; depth26 = _tmp$67;
			_tmp$68 = position; _tmp$69 = tokenIndex; _tmp$70 = depth; position27 = _tmp$68; tokenIndex27 = _tmp$69; depth27 = _tmp$70;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) {} else { $s = 82; continue; }
				/* goto l28 */ $s = 18; continue;
			/* } */ case 82:
			position = position + (1) >> 0;
			/* goto l27 */ $s = 19; continue;
			/* l28: */ case 18:
			_tmp$71 = position27; _tmp$72 = tokenIndex27; _tmp$73 = depth27; position = _tmp$71; tokenIndex = _tmp$72; depth = _tmp$73;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) {} else { $s = 83; continue; }
				/* goto l26 */ $s = 20; continue;
			/* } */ case 83:
			position = position + (1) >> 0;
			/* l27: */ case 19:
			/* goto l22 */ $s = 21; continue;
			/* l26: */ case 20:
			_tmp$74 = position26; _tmp$75 = tokenIndex26; _tmp$76 = depth26; position = _tmp$74; tokenIndex = _tmp$75; depth = _tmp$76;
			/* if (!matchDot()) { */ if (!matchDot()) {} else { $s = 84; continue; }
				/* goto l22 */ $s = 21; continue;
			/* } */ case 84:
			/* goto l21 */ $s = 22; continue;
			/* l22: */ case 21:
			_tmp$77 = position22; _tmp$78 = tokenIndex22; _tmp$79 = depth22; position = _tmp$77; tokenIndex = _tmp$78; depth = _tmp$79;
			position29 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) {} else { $s = 85; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 85:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 86; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 86:
			depth = depth - (1) >> 0;
			add(72, position29);
			/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 87; continue; }
				/* goto l6 */ $s = 4; continue;
			/* } */ case 87:
			depth = depth - (1) >> 0;
			add(3, position7);
			/* goto l5 */ $s = 23; continue;
			/* l6: */ case 4:
			_tmp$80 = position5; _tmp$81 = tokenIndex5; _tmp$82 = depth5; position = _tmp$80; tokenIndex = _tmp$81; depth = _tmp$82;
			position30 = position;
			depth = depth + (1) >> 0;
			position31 = position;
			depth = depth + (1) >> 0;
			_tmp$83 = position; _tmp$84 = tokenIndex; _tmp$85 = depth; position32 = _tmp$83; tokenIndex32 = _tmp$84; depth32 = _tmp$85;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 98))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 98))) {} else { $s = 88; continue; }
				/* goto l33 */ $s = 24; continue;
			/* } */ case 88:
			position = position + (1) >> 0;
			/* goto l32 */ $s = 25; continue;
			/* l33: */ case 24:
			_tmp$86 = position32; _tmp$87 = tokenIndex32; _tmp$88 = depth32; position = _tmp$86; tokenIndex = _tmp$87; depth = _tmp$88;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 66))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 66))) {} else { $s = 89; continue; }
				/* goto l4 */ $s = 26; continue;
			/* } */ case 89:
			position = position + (1) >> 0;
			/* l32: */ case 25:
			_tmp$89 = position; _tmp$90 = tokenIndex; _tmp$91 = depth; position34 = _tmp$89; tokenIndex34 = _tmp$90; depth34 = _tmp$91;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) {} else { $s = 90; continue; }
				/* goto l35 */ $s = 27; continue;
			/* } */ case 90:
			position = position + (1) >> 0;
			/* goto l34 */ $s = 28; continue;
			/* l35: */ case 27:
			_tmp$92 = position34; _tmp$93 = tokenIndex34; _tmp$94 = depth34; position = _tmp$92; tokenIndex = _tmp$93; depth = _tmp$94;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) {} else { $s = 91; continue; }
				/* goto l4 */ $s = 26; continue;
			/* } */ case 91:
			position = position + (1) >> 0;
			/* l34: */ case 28:
			_tmp$95 = position; _tmp$96 = tokenIndex; _tmp$97 = depth; position36 = _tmp$95; tokenIndex36 = _tmp$96; depth36 = _tmp$97;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) {} else { $s = 92; continue; }
				/* goto l37 */ $s = 29; continue;
			/* } */ case 92:
			position = position + (1) >> 0;
			/* goto l36 */ $s = 30; continue;
			/* l37: */ case 29:
			_tmp$98 = position36; _tmp$99 = tokenIndex36; _tmp$100 = depth36; position = _tmp$98; tokenIndex = _tmp$99; depth = _tmp$100;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) {} else { $s = 93; continue; }
				/* goto l4 */ $s = 26; continue;
			/* } */ case 93:
			position = position + (1) >> 0;
			/* l36: */ case 30:
			_tmp$101 = position; _tmp$102 = tokenIndex; _tmp$103 = depth; position38 = _tmp$101; tokenIndex38 = _tmp$102; depth38 = _tmp$103;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 94; continue; }
				/* goto l39 */ $s = 31; continue;
			/* } */ case 94:
			position = position + (1) >> 0;
			/* goto l38 */ $s = 32; continue;
			/* l39: */ case 31:
			_tmp$104 = position38; _tmp$105 = tokenIndex38; _tmp$106 = depth38; position = _tmp$104; tokenIndex = _tmp$105; depth = _tmp$106;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 95; continue; }
				/* goto l4 */ $s = 26; continue;
			/* } */ case 95:
			position = position + (1) >> 0;
			/* l38: */ case 32:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 96; continue; }
				/* goto l4 */ $s = 26; continue;
			/* } */ case 96:
			depth = depth - (1) >> 0;
			add(58, position31);
			/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 97; continue; }
				/* goto l4 */ $s = 26; continue;
			/* } */ case 97:
			depth = depth - (1) >> 0;
			add(4, position30);
			/* l5: */ case 23:
			/* goto l3 */ $s = 33; continue;
			/* l4: */ case 26:
			_tmp$107 = position4; _tmp$108 = tokenIndex4; _tmp$109 = depth4; position = _tmp$107; tokenIndex = _tmp$108; depth = _tmp$109;
			depth = depth - (1) >> 0;
			add(2, position2);
			position40 = position;
			depth = depth + (1) >> 0;
			position41 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[7]()) { */ if (!rules[7]()) {} else { $s = 98; continue; }
				/* goto l0 */ $s = 1; continue;
			/* } */ case 98:
			_tmp$110 = position; _tmp$111 = tokenIndex; _tmp$112 = depth; position42 = _tmp$110; tokenIndex42 = _tmp$111; depth42 = _tmp$112;
			position44 = position;
			depth = depth + (1) >> 0;
			position45 = position;
			depth = depth + (1) >> 0;
			_tmp$113 = position; _tmp$114 = tokenIndex; _tmp$115 = depth; position46 = _tmp$113; tokenIndex46 = _tmp$114; depth46 = _tmp$115;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) {} else { $s = 99; continue; }
				/* goto l47 */ $s = 34; continue;
			/* } */ case 99:
			position = position + (1) >> 0;
			/* goto l46 */ $s = 35; continue;
			/* l47: */ case 34:
			_tmp$116 = position46; _tmp$117 = tokenIndex46; _tmp$118 = depth46; position = _tmp$116; tokenIndex = _tmp$117; depth = _tmp$118;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) {} else { $s = 100; continue; }
				/* goto l42 */ $s = 36; continue;
			/* } */ case 100:
			position = position + (1) >> 0;
			/* l46: */ case 35:
			_tmp$119 = position; _tmp$120 = tokenIndex; _tmp$121 = depth; position48 = _tmp$119; tokenIndex48 = _tmp$120; depth48 = _tmp$121;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) {} else { $s = 101; continue; }
				/* goto l49 */ $s = 37; continue;
			/* } */ case 101:
			position = position + (1) >> 0;
			/* goto l48 */ $s = 38; continue;
			/* l49: */ case 37:
			_tmp$122 = position48; _tmp$123 = tokenIndex48; _tmp$124 = depth48; position = _tmp$122; tokenIndex = _tmp$123; depth = _tmp$124;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) {} else { $s = 102; continue; }
				/* goto l42 */ $s = 36; continue;
			/* } */ case 102:
			position = position + (1) >> 0;
			/* l48: */ case 38:
			_tmp$125 = position; _tmp$126 = tokenIndex; _tmp$127 = depth; position50 = _tmp$125; tokenIndex50 = _tmp$126; depth50 = _tmp$127;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) {} else { $s = 103; continue; }
				/* goto l51 */ $s = 39; continue;
			/* } */ case 103:
			position = position + (1) >> 0;
			/* goto l50 */ $s = 40; continue;
			/* l51: */ case 39:
			_tmp$128 = position50; _tmp$129 = tokenIndex50; _tmp$130 = depth50; position = _tmp$128; tokenIndex = _tmp$129; depth = _tmp$130;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) {} else { $s = 104; continue; }
				/* goto l42 */ $s = 36; continue;
			/* } */ case 104:
			position = position + (1) >> 0;
			/* l50: */ case 40:
			_tmp$131 = position; _tmp$132 = tokenIndex; _tmp$133 = depth; position52 = _tmp$131; tokenIndex52 = _tmp$132; depth52 = _tmp$133;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 109))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 109))) {} else { $s = 105; continue; }
				/* goto l53 */ $s = 41; continue;
			/* } */ case 105:
			position = position + (1) >> 0;
			/* goto l52 */ $s = 42; continue;
			/* l53: */ case 41:
			_tmp$134 = position52; _tmp$135 = tokenIndex52; _tmp$136 = depth52; position = _tmp$134; tokenIndex = _tmp$135; depth = _tmp$136;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 77))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 77))) {} else { $s = 106; continue; }
				/* goto l42 */ $s = 36; continue;
			/* } */ case 106:
			position = position + (1) >> 0;
			/* l52: */ case 42:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 107; continue; }
				/* goto l42 */ $s = 36; continue;
			/* } */ case 107:
			depth = depth - (1) >> 0;
			add(62, position45);
			_tmp$137 = position; _tmp$138 = tokenIndex; _tmp$139 = depth; position54 = _tmp$137; tokenIndex54 = _tmp$138; depth54 = _tmp$139;
			position56 = position;
			depth = depth + (1) >> 0;
			_tmp$140 = position; _tmp$141 = tokenIndex; _tmp$142 = depth; position57 = _tmp$140; tokenIndex57 = _tmp$141; depth57 = _tmp$142;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) {} else { $s = 108; continue; }
				/* goto l58 */ $s = 43; continue;
			/* } */ case 108:
			position = position + (1) >> 0;
			/* goto l57 */ $s = 44; continue;
			/* l58: */ case 43:
			_tmp$143 = position57; _tmp$144 = tokenIndex57; _tmp$145 = depth57; position = _tmp$143; tokenIndex = _tmp$144; depth = _tmp$145;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) {} else { $s = 109; continue; }
				/* goto l54 */ $s = 45; continue;
			/* } */ case 109:
			position = position + (1) >> 0;
			/* l57: */ case 44:
			_tmp$146 = position; _tmp$147 = tokenIndex; _tmp$148 = depth; position59 = _tmp$146; tokenIndex59 = _tmp$147; depth59 = _tmp$148;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) {} else { $s = 110; continue; }
				/* goto l60 */ $s = 46; continue;
			/* } */ case 110:
			position = position + (1) >> 0;
			/* goto l59 */ $s = 47; continue;
			/* l60: */ case 46:
			_tmp$149 = position59; _tmp$150 = tokenIndex59; _tmp$151 = depth59; position = _tmp$149; tokenIndex = _tmp$150; depth = _tmp$151;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) {} else { $s = 111; continue; }
				/* goto l54 */ $s = 45; continue;
			/* } */ case 111:
			position = position + (1) >> 0;
			/* l59: */ case 47:
			_tmp$152 = position; _tmp$153 = tokenIndex; _tmp$154 = depth; position61 = _tmp$152; tokenIndex61 = _tmp$153; depth61 = _tmp$154;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 109))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 109))) {} else { $s = 112; continue; }
				/* goto l62 */ $s = 48; continue;
			/* } */ case 112:
			position = position + (1) >> 0;
			/* goto l61 */ $s = 49; continue;
			/* l62: */ case 48:
			_tmp$155 = position61; _tmp$156 = tokenIndex61; _tmp$157 = depth61; position = _tmp$155; tokenIndex = _tmp$156; depth = _tmp$157;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 77))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 77))) {} else { $s = 113; continue; }
				/* goto l54 */ $s = 45; continue;
			/* } */ case 113:
			position = position + (1) >> 0;
			/* l61: */ case 49:
			_tmp$158 = position; _tmp$159 = tokenIndex; _tmp$160 = depth; position63 = _tmp$158; tokenIndex63 = _tmp$159; depth63 = _tmp$160;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 114; continue; }
				/* goto l64 */ $s = 50; continue;
			/* } */ case 114:
			position = position + (1) >> 0;
			/* goto l63 */ $s = 51; continue;
			/* l64: */ case 50:
			_tmp$161 = position63; _tmp$162 = tokenIndex63; _tmp$163 = depth63; position = _tmp$161; tokenIndex = _tmp$162; depth = _tmp$163;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 115; continue; }
				/* goto l54 */ $s = 45; continue;
			/* } */ case 115:
			position = position + (1) >> 0;
			/* l63: */ case 51:
			_tmp$164 = position; _tmp$165 = tokenIndex; _tmp$166 = depth; position65 = _tmp$164; tokenIndex65 = _tmp$165; depth65 = _tmp$166;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) {} else { $s = 116; continue; }
				/* goto l66 */ $s = 52; continue;
			/* } */ case 116:
			position = position + (1) >> 0;
			/* goto l65 */ $s = 53; continue;
			/* l66: */ case 52:
			_tmp$167 = position65; _tmp$168 = tokenIndex65; _tmp$169 = depth65; position = _tmp$167; tokenIndex = _tmp$168; depth = _tmp$169;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) {} else { $s = 117; continue; }
				/* goto l54 */ $s = 45; continue;
			/* } */ case 117:
			position = position + (1) >> 0;
			/* l65: */ case 53:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 118; continue; }
				/* goto l54 */ $s = 45; continue;
			/* } */ case 118:
			depth = depth - (1) >> 0;
			add(63, position56);
			/* goto l55 */ $s = 54; continue;
			/* l54: */ case 45:
			_tmp$170 = position54; _tmp$171 = tokenIndex54; _tmp$172 = depth54; position = _tmp$170; tokenIndex = _tmp$171; depth = _tmp$172;
			/* l55: */ case 54:
			/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 119; continue; }
				/* goto l42 */ $s = 36; continue;
			/* } */ case 119:
			depth = depth - (1) >> 0;
			add(10, position44);
			/* goto l43 */ $s = 55; continue;
			/* l42: */ case 36:
			_tmp$173 = position42; _tmp$174 = tokenIndex42; _tmp$175 = depth42; position = _tmp$173; tokenIndex = _tmp$174; depth = _tmp$175;
			/* l43: */ case 55:
			/* if (!rules[11]()) { */ if (!rules[11]()) {} else { $s = 120; continue; }
				/* goto l0 */ $s = 1; continue;
			/* } */ case 120:
			position67 = position;
			depth = depth + (1) >> 0;
			_tmp$176 = position; _tmp$177 = tokenIndex; _tmp$178 = depth; position68 = _tmp$176; tokenIndex68 = _tmp$177; depth68 = _tmp$178;
			position70 = position;
			depth = depth + (1) >> 0;
			_tmp$179 = position; _tmp$180 = tokenIndex; _tmp$181 = depth; position71 = _tmp$179; tokenIndex71 = _tmp$180; depth71 = _tmp$181;
			/* if (!rules[39]()) { */ if (!rules[39]()) {} else { $s = 121; continue; }
				/* goto l72 */ $s = 56; continue;
			/* } */ case 121:
			_tmp$182 = position; _tmp$183 = tokenIndex; _tmp$184 = depth; position73 = _tmp$182; tokenIndex73 = _tmp$183; depth73 = _tmp$184;
			/* if (!rules[40]()) { */ if (!rules[40]()) {} else { $s = 122; continue; }
				/* goto l73 */ $s = 57; continue;
			/* } */ case 122:
			/* goto l74 */ $s = 58; continue;
			/* l73: */ case 57:
			_tmp$185 = position73; _tmp$186 = tokenIndex73; _tmp$187 = depth73; position = _tmp$185; tokenIndex = _tmp$186; depth = _tmp$187;
			/* l74: */ case 58:
			/* goto l71 */ $s = 59; continue;
			/* l72: */ case 56:
			_tmp$188 = position71; _tmp$189 = tokenIndex71; _tmp$190 = depth71; position = _tmp$188; tokenIndex = _tmp$189; depth = _tmp$190;
			/* if (!rules[40]()) { */ if (!rules[40]()) {} else { $s = 123; continue; }
				/* goto l68 */ $s = 60; continue;
			/* } */ case 123:
			_tmp$191 = position; _tmp$192 = tokenIndex; _tmp$193 = depth; position75 = _tmp$191; tokenIndex75 = _tmp$192; depth75 = _tmp$193;
			/* if (!rules[39]()) { */ if (!rules[39]()) {} else { $s = 124; continue; }
				/* goto l75 */ $s = 61; continue;
			/* } */ case 124:
			/* goto l76 */ $s = 62; continue;
			/* l75: */ case 61:
			_tmp$194 = position75; _tmp$195 = tokenIndex75; _tmp$196 = depth75; position = _tmp$194; tokenIndex = _tmp$195; depth = _tmp$196;
			/* l71: */ case 59:
			depth = depth - (1) >> 0;
			add(38, position70);
			/* goto l69 */ $s = 63; continue;
			/* l68: */ case 60:
			_tmp$197 = position68; _tmp$198 = tokenIndex68; _tmp$199 = depth68; position = _tmp$197; tokenIndex = _tmp$198; depth = _tmp$199;
			/* l69: */ case 63:
			depth = depth - (1) >> 0;
			add(37, position67);
			depth = depth - (1) >> 0;
			add(6, position41);
			depth = depth - (1) >> 0;
			add(5, position40);
			_tmp$200 = position; _tmp$201 = tokenIndex; _tmp$202 = depth; position77 = _tmp$200; tokenIndex77 = _tmp$201; depth77 = _tmp$202;
			/* if (!matchDot()) { */ if (!matchDot()) {} else { $s = 125; continue; }
				/* goto l77 */ $s = 64; continue;
			/* } */ case 125:
			/* goto l0 */ $s = 1; continue;
			/* l77: */ case 64:
			_tmp$203 = position77; _tmp$204 = tokenIndex77; _tmp$205 = depth77; position = _tmp$203; tokenIndex = _tmp$204; depth = _tmp$205;
			depth = depth - (1) >> 0;
			add(1, position1);
			return true;
			/* l0: */ case 1:
			_tmp$206 = position0; _tmp$207 = tokenIndex0; _tmp$208 = depth0; position = _tmp$206; tokenIndex = _tmp$207; depth = _tmp$208;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position83, tokenIndex83, depth83, position84, position85, _tmp$8, _tmp$9, _tmp$10, position86, tokenIndex86, depth86, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, position88, tokenIndex88, depth88, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, position90, tokenIndex90, depth90, _tmp$23, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, position92, tokenIndex92, depth92, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, position94, tokenIndex94, depth94, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40, position96, tokenIndex96, depth96, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46, position98, tokenIndex98, depth98, _tmp$47, _tmp$48, _tmp$49, position100, tokenIndex100, depth100, position102, _tmp$50, _tmp$51, _tmp$52, position103, tokenIndex103, depth103, _tmp$53, _tmp$54, _tmp$55, _tmp$56, _tmp$57, _tmp$58, position105, tokenIndex105, depth105, _tmp$59, _tmp$60, _tmp$61, _tmp$62, _tmp$63, _tmp$64, position107, tokenIndex107, depth107, _tmp$65, _tmp$66, _tmp$67, _tmp$68, _tmp$69, _tmp$70, position109, tokenIndex109, depth109, _tmp$71, _tmp$72, _tmp$73, _tmp$74, _tmp$75, _tmp$76, position111, tokenIndex111, depth111, _tmp$77, _tmp$78, _tmp$79, _tmp$80, _tmp$81, _tmp$82, position113, tokenIndex113, depth113, _tmp$83, _tmp$84, _tmp$85, _tmp$86, _tmp$87, _tmp$88, position115, tokenIndex115, depth115, _tmp$89, _tmp$90, _tmp$91, _tmp$92, _tmp$93, _tmp$94, position117, tokenIndex117, depth117, _tmp$95, _tmp$96, _tmp$97, _tmp$98, _tmp$99, _tmp$100, position119, _tmp$101, _tmp$102, _tmp$103, position120, tokenIndex120, depth120, _tmp$104, _tmp$105, _tmp$106, _tmp$107, _tmp$108, _tmp$109, position122, tokenIndex122, depth122, _tmp$110, _tmp$111, _tmp$112, _tmp$113, _tmp$114, _tmp$115, position124, tokenIndex124, depth124, _tmp$116, _tmp$117, _tmp$118, _tmp$119, _tmp$120, _tmp$121, position126, tokenIndex126, depth126, _tmp$122, _tmp$123, _tmp$124, _tmp$125, _tmp$126, _tmp$127, position128, tokenIndex128, depth128, _tmp$128, _tmp$129, _tmp$130, _tmp$131, _tmp$132, _tmp$133, position130, tokenIndex130, depth130, _tmp$134, _tmp$135, _tmp$136, _tmp$137, _tmp$138, _tmp$139, position132, tokenIndex132, depth132, _tmp$140, _tmp$141, _tmp$142, _tmp$143, _tmp$144, _tmp$145, _tmp$146, _tmp$147, _tmp$148, position134, tokenIndex134, depth134, position136, _tmp$149, _tmp$150, _tmp$151, position139, _tmp$152, _tmp$153, _tmp$154, position138, tokenIndex138, depth138, position140, _tmp$155, _tmp$156, _tmp$157, _tmp$158, _tmp$159, _tmp$160;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position83 = _tmp$5; tokenIndex83 = _tmp$6; depth83 = _tmp$7;
			position84 = position;
			depth = depth + (1) >> 0;
			position85 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position86 = _tmp$8; tokenIndex86 = _tmp$9; depth86 = _tmp$10;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) {} else { $s = 52; continue; }
				/* goto l87 */ $s = 1; continue;
			/* } */ case 52:
			position = position + (1) >> 0;
			/* goto l86 */ $s = 2; continue;
			/* l87: */ case 1:
			_tmp$11 = position86; _tmp$12 = tokenIndex86; _tmp$13 = depth86; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) {} else { $s = 53; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 53:
			position = position + (1) >> 0;
			/* l86: */ case 2:
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position88 = _tmp$14; tokenIndex88 = _tmp$15; depth88 = _tmp$16;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 54; continue; }
				/* goto l89 */ $s = 4; continue;
			/* } */ case 54:
			position = position + (1) >> 0;
			/* goto l88 */ $s = 5; continue;
			/* l89: */ case 4:
			_tmp$17 = position88; _tmp$18 = tokenIndex88; _tmp$19 = depth88; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 55; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 55:
			position = position + (1) >> 0;
			/* l88: */ case 5:
			_tmp$20 = position; _tmp$21 = tokenIndex; _tmp$22 = depth; position90 = _tmp$20; tokenIndex90 = _tmp$21; depth90 = _tmp$22;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) {} else { $s = 56; continue; }
				/* goto l91 */ $s = 6; continue;
			/* } */ case 56:
			position = position + (1) >> 0;
			/* goto l90 */ $s = 7; continue;
			/* l91: */ case 6:
			_tmp$23 = position90; _tmp$24 = tokenIndex90; _tmp$25 = depth90; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) {} else { $s = 57; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 57:
			position = position + (1) >> 0;
			/* l90: */ case 7:
			_tmp$26 = position; _tmp$27 = tokenIndex; _tmp$28 = depth; position92 = _tmp$26; tokenIndex92 = _tmp$27; depth92 = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 58; continue; }
				/* goto l93 */ $s = 8; continue;
			/* } */ case 58:
			position = position + (1) >> 0;
			/* goto l92 */ $s = 9; continue;
			/* l93: */ case 8:
			_tmp$29 = position92; _tmp$30 = tokenIndex92; _tmp$31 = depth92; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 59; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 59:
			position = position + (1) >> 0;
			/* l92: */ case 9:
			_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position94 = _tmp$32; tokenIndex94 = _tmp$33; depth94 = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 99))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 99))) {} else { $s = 60; continue; }
				/* goto l95 */ $s = 10; continue;
			/* } */ case 60:
			position = position + (1) >> 0;
			/* goto l94 */ $s = 11; continue;
			/* l95: */ case 10:
			_tmp$35 = position94; _tmp$36 = tokenIndex94; _tmp$37 = depth94; position = _tmp$35; tokenIndex = _tmp$36; depth = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 67))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 67))) {} else { $s = 61; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 61:
			position = position + (1) >> 0;
			/* l94: */ case 11:
			_tmp$38 = position; _tmp$39 = tokenIndex; _tmp$40 = depth; position96 = _tmp$38; tokenIndex96 = _tmp$39; depth96 = _tmp$40;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 62; continue; }
				/* goto l97 */ $s = 12; continue;
			/* } */ case 62:
			position = position + (1) >> 0;
			/* goto l96 */ $s = 13; continue;
			/* l97: */ case 12:
			_tmp$41 = position96; _tmp$42 = tokenIndex96; _tmp$43 = depth96; position = _tmp$41; tokenIndex = _tmp$42; depth = _tmp$43;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 63; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 63:
			position = position + (1) >> 0;
			/* l96: */ case 13:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 64; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 64:
			depth = depth - (1) >> 0;
			add(59, position85);
			_tmp$44 = position; _tmp$45 = tokenIndex; _tmp$46 = depth; position98 = _tmp$44; tokenIndex98 = _tmp$45; depth98 = _tmp$46;
			_tmp$47 = position; _tmp$48 = tokenIndex; _tmp$49 = depth; position100 = _tmp$47; tokenIndex100 = _tmp$48; depth100 = _tmp$49;
			position102 = position;
			depth = depth + (1) >> 0;
			_tmp$50 = position; _tmp$51 = tokenIndex; _tmp$52 = depth; position103 = _tmp$50; tokenIndex103 = _tmp$51; depth103 = _tmp$52;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) {} else { $s = 65; continue; }
				/* goto l104 */ $s = 14; continue;
			/* } */ case 65:
			position = position + (1) >> 0;
			/* goto l103 */ $s = 15; continue;
			/* l104: */ case 14:
			_tmp$53 = position103; _tmp$54 = tokenIndex103; _tmp$55 = depth103; position = _tmp$53; tokenIndex = _tmp$54; depth = _tmp$55;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) {} else { $s = 66; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 66:
			position = position + (1) >> 0;
			/* l103: */ case 15:
			_tmp$56 = position; _tmp$57 = tokenIndex; _tmp$58 = depth; position105 = _tmp$56; tokenIndex105 = _tmp$57; depth105 = _tmp$58;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 67; continue; }
				/* goto l106 */ $s = 17; continue;
			/* } */ case 67:
			position = position + (1) >> 0;
			/* goto l105 */ $s = 18; continue;
			/* l106: */ case 17:
			_tmp$59 = position105; _tmp$60 = tokenIndex105; _tmp$61 = depth105; position = _tmp$59; tokenIndex = _tmp$60; depth = _tmp$61;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 68; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 68:
			position = position + (1) >> 0;
			/* l105: */ case 18:
			_tmp$62 = position; _tmp$63 = tokenIndex; _tmp$64 = depth; position107 = _tmp$62; tokenIndex107 = _tmp$63; depth107 = _tmp$64;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) {} else { $s = 69; continue; }
				/* goto l108 */ $s = 19; continue;
			/* } */ case 69:
			position = position + (1) >> 0;
			/* goto l107 */ $s = 20; continue;
			/* l108: */ case 19:
			_tmp$65 = position107; _tmp$66 = tokenIndex107; _tmp$67 = depth107; position = _tmp$65; tokenIndex = _tmp$66; depth = _tmp$67;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) {} else { $s = 70; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 70:
			position = position + (1) >> 0;
			/* l107: */ case 20:
			_tmp$68 = position; _tmp$69 = tokenIndex; _tmp$70 = depth; position109 = _tmp$68; tokenIndex109 = _tmp$69; depth109 = _tmp$70;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 71; continue; }
				/* goto l110 */ $s = 21; continue;
			/* } */ case 71:
			position = position + (1) >> 0;
			/* goto l109 */ $s = 22; continue;
			/* l110: */ case 21:
			_tmp$71 = position109; _tmp$72 = tokenIndex109; _tmp$73 = depth109; position = _tmp$71; tokenIndex = _tmp$72; depth = _tmp$73;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 72; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 72:
			position = position + (1) >> 0;
			/* l109: */ case 22:
			_tmp$74 = position; _tmp$75 = tokenIndex; _tmp$76 = depth; position111 = _tmp$74; tokenIndex111 = _tmp$75; depth111 = _tmp$76;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 73; continue; }
				/* goto l112 */ $s = 23; continue;
			/* } */ case 73:
			position = position + (1) >> 0;
			/* goto l111 */ $s = 24; continue;
			/* l112: */ case 23:
			_tmp$77 = position111; _tmp$78 = tokenIndex111; _tmp$79 = depth111; position = _tmp$77; tokenIndex = _tmp$78; depth = _tmp$79;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 74; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 74:
			position = position + (1) >> 0;
			/* l111: */ case 24:
			_tmp$80 = position; _tmp$81 = tokenIndex; _tmp$82 = depth; position113 = _tmp$80; tokenIndex113 = _tmp$81; depth113 = _tmp$82;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) {} else { $s = 75; continue; }
				/* goto l114 */ $s = 25; continue;
			/* } */ case 75:
			position = position + (1) >> 0;
			/* goto l113 */ $s = 26; continue;
			/* l114: */ case 25:
			_tmp$83 = position113; _tmp$84 = tokenIndex113; _tmp$85 = depth113; position = _tmp$83; tokenIndex = _tmp$84; depth = _tmp$85;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) {} else { $s = 76; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 76:
			position = position + (1) >> 0;
			/* l113: */ case 26:
			_tmp$86 = position; _tmp$87 = tokenIndex; _tmp$88 = depth; position115 = _tmp$86; tokenIndex115 = _tmp$87; depth115 = _tmp$88;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 99))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 99))) {} else { $s = 77; continue; }
				/* goto l116 */ $s = 27; continue;
			/* } */ case 77:
			position = position + (1) >> 0;
			/* goto l115 */ $s = 28; continue;
			/* l116: */ case 27:
			_tmp$89 = position115; _tmp$90 = tokenIndex115; _tmp$91 = depth115; position = _tmp$89; tokenIndex = _tmp$90; depth = _tmp$91;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 67))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 67))) {} else { $s = 78; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 78:
			position = position + (1) >> 0;
			/* l115: */ case 28:
			_tmp$92 = position; _tmp$93 = tokenIndex; _tmp$94 = depth; position117 = _tmp$92; tokenIndex117 = _tmp$93; depth117 = _tmp$94;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 79; continue; }
				/* goto l118 */ $s = 29; continue;
			/* } */ case 79:
			position = position + (1) >> 0;
			/* goto l117 */ $s = 30; continue;
			/* l118: */ case 29:
			_tmp$95 = position117; _tmp$96 = tokenIndex117; _tmp$97 = depth117; position = _tmp$95; tokenIndex = _tmp$96; depth = _tmp$97;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 80; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 80:
			position = position + (1) >> 0;
			/* l117: */ case 30:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 81; continue; }
				/* goto l101 */ $s = 16; continue;
			/* } */ case 81:
			depth = depth - (1) >> 0;
			add(61, position102);
			/* goto l100 */ $s = 31; continue;
			/* l101: */ case 16:
			_tmp$98 = position100; _tmp$99 = tokenIndex100; _tmp$100 = depth100; position = _tmp$98; tokenIndex = _tmp$99; depth = _tmp$100;
			position119 = position;
			depth = depth + (1) >> 0;
			_tmp$101 = position; _tmp$102 = tokenIndex; _tmp$103 = depth; position120 = _tmp$101; tokenIndex120 = _tmp$102; depth120 = _tmp$103;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) {} else { $s = 82; continue; }
				/* goto l121 */ $s = 32; continue;
			/* } */ case 82:
			position = position + (1) >> 0;
			/* goto l120 */ $s = 33; continue;
			/* l121: */ case 32:
			_tmp$104 = position120; _tmp$105 = tokenIndex120; _tmp$106 = depth120; position = _tmp$104; tokenIndex = _tmp$105; depth = _tmp$106;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) {} else { $s = 83; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 83:
			position = position + (1) >> 0;
			/* l120: */ case 33:
			_tmp$107 = position; _tmp$108 = tokenIndex; _tmp$109 = depth; position122 = _tmp$107; tokenIndex122 = _tmp$108; depth122 = _tmp$109;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 84; continue; }
				/* goto l123 */ $s = 35; continue;
			/* } */ case 84:
			position = position + (1) >> 0;
			/* goto l122 */ $s = 36; continue;
			/* l123: */ case 35:
			_tmp$110 = position122; _tmp$111 = tokenIndex122; _tmp$112 = depth122; position = _tmp$110; tokenIndex = _tmp$111; depth = _tmp$112;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 85; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 85:
			position = position + (1) >> 0;
			/* l122: */ case 36:
			_tmp$113 = position; _tmp$114 = tokenIndex; _tmp$115 = depth; position124 = _tmp$113; tokenIndex124 = _tmp$114; depth124 = _tmp$115;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) {} else { $s = 86; continue; }
				/* goto l125 */ $s = 37; continue;
			/* } */ case 86:
			position = position + (1) >> 0;
			/* goto l124 */ $s = 38; continue;
			/* l125: */ case 37:
			_tmp$116 = position124; _tmp$117 = tokenIndex124; _tmp$118 = depth124; position = _tmp$116; tokenIndex = _tmp$117; depth = _tmp$118;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) {} else { $s = 87; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 87:
			position = position + (1) >> 0;
			/* l124: */ case 38:
			_tmp$119 = position; _tmp$120 = tokenIndex; _tmp$121 = depth; position126 = _tmp$119; tokenIndex126 = _tmp$120; depth126 = _tmp$121;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 117))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 117))) {} else { $s = 88; continue; }
				/* goto l127 */ $s = 39; continue;
			/* } */ case 88:
			position = position + (1) >> 0;
			/* goto l126 */ $s = 40; continue;
			/* l127: */ case 39:
			_tmp$122 = position126; _tmp$123 = tokenIndex126; _tmp$124 = depth126; position = _tmp$122; tokenIndex = _tmp$123; depth = _tmp$124;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 85))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 85))) {} else { $s = 89; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 89:
			position = position + (1) >> 0;
			/* l126: */ case 40:
			_tmp$125 = position; _tmp$126 = tokenIndex; _tmp$127 = depth; position128 = _tmp$125; tokenIndex128 = _tmp$126; depth128 = _tmp$127;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 99))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 99))) {} else { $s = 90; continue; }
				/* goto l129 */ $s = 41; continue;
			/* } */ case 90:
			position = position + (1) >> 0;
			/* goto l128 */ $s = 42; continue;
			/* l129: */ case 41:
			_tmp$128 = position128; _tmp$129 = tokenIndex128; _tmp$130 = depth128; position = _tmp$128; tokenIndex = _tmp$129; depth = _tmp$130;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 67))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 67))) {} else { $s = 91; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 91:
			position = position + (1) >> 0;
			/* l128: */ case 42:
			_tmp$131 = position; _tmp$132 = tokenIndex; _tmp$133 = depth; position130 = _tmp$131; tokenIndex130 = _tmp$132; depth130 = _tmp$133;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 92; continue; }
				/* goto l131 */ $s = 43; continue;
			/* } */ case 92:
			position = position + (1) >> 0;
			/* goto l130 */ $s = 44; continue;
			/* l131: */ case 43:
			_tmp$134 = position130; _tmp$135 = tokenIndex130; _tmp$136 = depth130; position = _tmp$134; tokenIndex = _tmp$135; depth = _tmp$136;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 93; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 93:
			position = position + (1) >> 0;
			/* l130: */ case 44:
			_tmp$137 = position; _tmp$138 = tokenIndex; _tmp$139 = depth; position132 = _tmp$137; tokenIndex132 = _tmp$138; depth132 = _tmp$139;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 100))) {} else { $s = 94; continue; }
				/* goto l133 */ $s = 45; continue;
			/* } */ case 94:
			position = position + (1) >> 0;
			/* goto l132 */ $s = 46; continue;
			/* l133: */ case 45:
			_tmp$140 = position132; _tmp$141 = tokenIndex132; _tmp$142 = depth132; position = _tmp$140; tokenIndex = _tmp$141; depth = _tmp$142;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 68))) {} else { $s = 95; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 95:
			position = position + (1) >> 0;
			/* l132: */ case 46:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 96; continue; }
				/* goto l98 */ $s = 34; continue;
			/* } */ case 96:
			depth = depth - (1) >> 0;
			add(60, position119);
			/* l100: */ case 31:
			/* goto l99 */ $s = 47; continue;
			/* l98: */ case 34:
			_tmp$143 = position98; _tmp$144 = tokenIndex98; _tmp$145 = depth98; position = _tmp$143; tokenIndex = _tmp$144; depth = _tmp$145;
			/* l99: */ case 47:
			_tmp$146 = position; _tmp$147 = tokenIndex; _tmp$148 = depth; position134 = _tmp$146; tokenIndex134 = _tmp$147; depth134 = _tmp$148;
			position136 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 42))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 42))) {} else { $s = 97; continue; }
				/* goto l135 */ $s = 48; continue;
			/* } */ case 97:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 98; continue; }
				/* goto l135 */ $s = 48; continue;
			/* } */ case 98:
			depth = depth - (1) >> 0;
			add(80, position136);
			/* goto l134 */ $s = 49; continue;
			/* l135: */ case 48:
			_tmp$149 = position134; _tmp$150 = tokenIndex134; _tmp$151 = depth134; position = _tmp$149; tokenIndex = _tmp$150; depth = _tmp$151;
			position139 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[42]()) { */ if (!rules[42]()) {} else { $s = 99; continue; }
				/* goto l83 */ $s = 3; continue;
			/* } */ case 99:
			depth = depth - (1) >> 0;
			add(9, position139);
			/* l137: */ case 51:
			_tmp$152 = position; _tmp$153 = tokenIndex; _tmp$154 = depth; position138 = _tmp$152; tokenIndex138 = _tmp$153; depth138 = _tmp$154;
			position140 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[42]()) { */ if (!rules[42]()) {} else { $s = 100; continue; }
				/* goto l138 */ $s = 50; continue;
			/* } */ case 100:
			depth = depth - (1) >> 0;
			add(9, position140);
			/* goto l137 */ $s = 51; continue;
			/* l138: */ case 50:
			_tmp$155 = position138; _tmp$156 = tokenIndex138; _tmp$157 = depth138; position = _tmp$155; tokenIndex = _tmp$156; depth = _tmp$157;
			/* l134: */ case 49:
			depth = depth - (1) >> 0;
			add(7, position84);
			return true;
			/* l83: */ case 3:
			_tmp$158 = position83; _tmp$159 = tokenIndex83; _tmp$160 = depth83; position = _tmp$158; tokenIndex = _tmp$159; depth = _tmp$160;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position141, tokenIndex141, depth141, position142, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position141 = _tmp$5; tokenIndex141 = _tmp$6; depth141 = _tmp$7;
			position142 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[7]()) { */ if (!rules[7]()) {} else { $s = 2; continue; }
				/* goto l141 */ $s = 1; continue;
			/* } */ case 2:
			/* if (!rules[11]()) { */ if (!rules[11]()) {} else { $s = 3; continue; }
				/* goto l141 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(8, position142);
			return true;
			/* l141: */ case 1:
			_tmp$8 = position141; _tmp$9 = tokenIndex141; _tmp$10 = depth141; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position145, tokenIndex145, depth145, position146, _tmp$8, _tmp$9, _tmp$10, position147, tokenIndex147, depth147, position149, _tmp$11, _tmp$12, _tmp$13, position150, tokenIndex150, depth150, _tmp$14, _tmp$15, _tmp$16, _tmp$17, _tmp$18, _tmp$19, position152, tokenIndex152, depth152, _tmp$20, _tmp$21, _tmp$22, _tmp$23, _tmp$24, _tmp$25, position154, tokenIndex154, depth154, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, position156, tokenIndex156, depth156, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, position158, tokenIndex158, depth158, _tmp$38, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position145 = _tmp$5; tokenIndex145 = _tmp$6; depth145 = _tmp$7;
			position146 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position147 = _tmp$8; tokenIndex147 = _tmp$9; depth147 = _tmp$10;
			position149 = position;
			depth = depth + (1) >> 0;
			_tmp$11 = position; _tmp$12 = tokenIndex; _tmp$13 = depth; position150 = _tmp$11; tokenIndex150 = _tmp$12; depth150 = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 119))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 119))) {} else { $s = 14; continue; }
				/* goto l151 */ $s = 1; continue;
			/* } */ case 14:
			position = position + (1) >> 0;
			/* goto l150 */ $s = 2; continue;
			/* l151: */ case 1:
			_tmp$14 = position150; _tmp$15 = tokenIndex150; _tmp$16 = depth150; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 87))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 87))) {} else { $s = 15; continue; }
				/* goto l147 */ $s = 3; continue;
			/* } */ case 15:
			position = position + (1) >> 0;
			/* l150: */ case 2:
			_tmp$17 = position; _tmp$18 = tokenIndex; _tmp$19 = depth; position152 = _tmp$17; tokenIndex152 = _tmp$18; depth152 = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 104))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 104))) {} else { $s = 16; continue; }
				/* goto l153 */ $s = 4; continue;
			/* } */ case 16:
			position = position + (1) >> 0;
			/* goto l152 */ $s = 5; continue;
			/* l153: */ case 4:
			_tmp$20 = position152; _tmp$21 = tokenIndex152; _tmp$22 = depth152; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 72))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 72))) {} else { $s = 17; continue; }
				/* goto l147 */ $s = 3; continue;
			/* } */ case 17:
			position = position + (1) >> 0;
			/* l152: */ case 5:
			_tmp$23 = position; _tmp$24 = tokenIndex; _tmp$25 = depth; position154 = _tmp$23; tokenIndex154 = _tmp$24; depth154 = _tmp$25;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 18; continue; }
				/* goto l155 */ $s = 6; continue;
			/* } */ case 18:
			position = position + (1) >> 0;
			/* goto l154 */ $s = 7; continue;
			/* l155: */ case 6:
			_tmp$26 = position154; _tmp$27 = tokenIndex154; _tmp$28 = depth154; position = _tmp$26; tokenIndex = _tmp$27; depth = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 19; continue; }
				/* goto l147 */ $s = 3; continue;
			/* } */ case 19:
			position = position + (1) >> 0;
			/* l154: */ case 7:
			_tmp$29 = position; _tmp$30 = tokenIndex; _tmp$31 = depth; position156 = _tmp$29; tokenIndex156 = _tmp$30; depth156 = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) {} else { $s = 20; continue; }
				/* goto l157 */ $s = 8; continue;
			/* } */ case 20:
			position = position + (1) >> 0;
			/* goto l156 */ $s = 9; continue;
			/* l157: */ case 8:
			_tmp$32 = position156; _tmp$33 = tokenIndex156; _tmp$34 = depth156; position = _tmp$32; tokenIndex = _tmp$33; depth = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) {} else { $s = 21; continue; }
				/* goto l147 */ $s = 3; continue;
			/* } */ case 21:
			position = position + (1) >> 0;
			/* l156: */ case 9:
			_tmp$35 = position; _tmp$36 = tokenIndex; _tmp$37 = depth; position158 = _tmp$35; tokenIndex158 = _tmp$36; depth158 = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 22; continue; }
				/* goto l159 */ $s = 10; continue;
			/* } */ case 22:
			position = position + (1) >> 0;
			/* goto l158 */ $s = 11; continue;
			/* l159: */ case 10:
			_tmp$38 = position158; _tmp$39 = tokenIndex158; _tmp$40 = depth158; position = _tmp$38; tokenIndex = _tmp$39; depth = _tmp$40;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 23; continue; }
				/* goto l147 */ $s = 3; continue;
			/* } */ case 23:
			position = position + (1) >> 0;
			/* l158: */ case 11:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 24; continue; }
				/* goto l147 */ $s = 3; continue;
			/* } */ case 24:
			depth = depth - (1) >> 0;
			add(64, position149);
			/* goto l148 */ $s = 12; continue;
			/* l147: */ case 3:
			_tmp$41 = position147; _tmp$42 = tokenIndex147; _tmp$43 = depth147; position = _tmp$41; tokenIndex = _tmp$42; depth = _tmp$43;
			/* l148: */ case 12:
			/* if (!rules[12]()) { */ if (!rules[12]()) {} else { $s = 25; continue; }
				/* goto l145 */ $s = 13; continue;
			/* } */ case 25:
			depth = depth - (1) >> 0;
			add(11, position146);
			return true;
			/* l145: */ case 13:
			_tmp$44 = position145; _tmp$45 = tokenIndex145; _tmp$46 = depth145; position = _tmp$44; tokenIndex = _tmp$45; depth = _tmp$46;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position160, tokenIndex160, depth160, position161, _tmp$8, _tmp$9, _tmp$10, position162, tokenIndex162, depth162, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position160 = _tmp$5; tokenIndex160 = _tmp$6; depth160 = _tmp$7;
			position161 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[65]()) { */ if (!rules[65]()) {} else { $s = 4; continue; }
				/* goto l160 */ $s = 1; continue;
			/* } */ case 4:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position162 = _tmp$8; tokenIndex162 = _tmp$9; depth162 = _tmp$10;
			/* if (!rules[8]()) { */ if (!rules[8]()) {} else { $s = 5; continue; }
				/* goto l163 */ $s = 2; continue;
			/* } */ case 5:
			/* goto l162 */ $s = 3; continue;
			/* l163: */ case 2:
			_tmp$11 = position162; _tmp$12 = tokenIndex162; _tmp$13 = depth162; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!rules[13]()) { */ if (!rules[13]()) {} else { $s = 6; continue; }
				/* goto l160 */ $s = 1; continue;
			/* } */ case 6:
			/* l162: */ case 3:
			/* if (!rules[66]()) { */ if (!rules[66]()) {} else { $s = 7; continue; }
				/* goto l160 */ $s = 1; continue;
			/* } */ case 7:
			depth = depth - (1) >> 0;
			add(12, position161);
			return true;
			/* l160: */ case 1:
			_tmp$14 = position160; _tmp$15 = tokenIndex160; _tmp$16 = depth160; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, position165, _tmp$5, _tmp$6, _tmp$7, position166, tokenIndex166, depth166, position168, position169, _tmp$8, _tmp$9, _tmp$10, position171, tokenIndex171, depth171, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, position172, tokenIndex172, depth172, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, _tmp$23, _tmp$24, _tmp$25, position174, tokenIndex174, depth174, position176, _tmp$26, _tmp$27, _tmp$28, position177, tokenIndex177, depth177, position179, position180, _tmp$29, _tmp$30, _tmp$31, position181, tokenIndex181, depth181, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, position183, tokenIndex183, depth183, _tmp$38, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, position185, tokenIndex185, depth185, _tmp$44, _tmp$45, _tmp$46, _tmp$47, _tmp$48, _tmp$49, position187, tokenIndex187, depth187, _tmp$50, _tmp$51, _tmp$52, _tmp$53, _tmp$54, _tmp$55, position189, tokenIndex189, depth189, _tmp$56, _tmp$57, _tmp$58, _tmp$59, _tmp$60, _tmp$61, position191, tokenIndex191, depth191, _tmp$62, _tmp$63, _tmp$64, _tmp$65, _tmp$66, _tmp$67, position193, tokenIndex193, depth193, _tmp$68, _tmp$69, _tmp$70, _tmp$71, _tmp$72, _tmp$73, position195, tokenIndex195, depth195, _tmp$74, _tmp$75, _tmp$76, _tmp$77, _tmp$78, _tmp$79, position197, tokenIndex197, depth197, _tmp$80, _tmp$81, _tmp$82, _tmp$83, _tmp$84, _tmp$85, _tmp$86, _tmp$87, _tmp$88, position199, tokenIndex199, depth199, _tmp$89, _tmp$90, _tmp$91, _tmp$92, _tmp$93, _tmp$94;
			/* */ while (true) { switch ($s) { case 0:
			position165 = position;
			depth = depth + (1) >> 0;
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position166 = _tmp$5; tokenIndex166 = _tmp$6; depth166 = _tmp$7;
			position168 = position;
			depth = depth + (1) >> 0;
			position169 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[19]()) { */ if (!rules[19]()) {} else { $s = 31; continue; }
				/* goto l166 */ $s = 1; continue;
			/* } */ case 31:
			/* l170: */ case 3:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position171 = _tmp$8; tokenIndex171 = _tmp$9; depth171 = _tmp$10;
			/* if (!rules[71]()) { */ if (!rules[71]()) {} else { $s = 32; continue; }
				/* goto l171 */ $s = 2; continue;
			/* } */ case 32:
			/* if (!rules[19]()) { */ if (!rules[19]()) {} else { $s = 33; continue; }
				/* goto l171 */ $s = 2; continue;
			/* } */ case 33:
			/* goto l170 */ $s = 3; continue;
			/* l171: */ case 2:
			_tmp$11 = position171; _tmp$12 = tokenIndex171; _tmp$13 = depth171; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position172 = _tmp$14; tokenIndex172 = _tmp$15; depth172 = _tmp$16;
			/* if (!rules[71]()) { */ if (!rules[71]()) {} else { $s = 34; continue; }
				/* goto l172 */ $s = 4; continue;
			/* } */ case 34:
			/* goto l173 */ $s = 5; continue;
			/* l172: */ case 4:
			_tmp$17 = position172; _tmp$18 = tokenIndex172; _tmp$19 = depth172; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* l173: */ case 5:
			depth = depth - (1) >> 0;
			add(18, position169);
			depth = depth - (1) >> 0;
			add(17, position168);
			/* goto l167 */ $s = 6; continue;
			/* l166: */ case 1:
			_tmp$20 = position166; _tmp$21 = tokenIndex166; _tmp$22 = depth166; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
			/* l167: */ case 6:
			_tmp$23 = position; _tmp$24 = tokenIndex; _tmp$25 = depth; position174 = _tmp$23; tokenIndex174 = _tmp$24; depth174 = _tmp$25;
			position176 = position;
			depth = depth + (1) >> 0;
			_tmp$26 = position; _tmp$27 = tokenIndex; _tmp$28 = depth; position177 = _tmp$26; tokenIndex177 = _tmp$27; depth177 = _tmp$28;
			position179 = position;
			depth = depth + (1) >> 0;
			position180 = position;
			depth = depth + (1) >> 0;
			_tmp$29 = position; _tmp$30 = tokenIndex; _tmp$31 = depth; position181 = _tmp$29; tokenIndex181 = _tmp$30; depth181 = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) {} else { $s = 35; continue; }
				/* goto l182 */ $s = 7; continue;
			/* } */ case 35:
			position = position + (1) >> 0;
			/* goto l181 */ $s = 8; continue;
			/* l182: */ case 7:
			_tmp$32 = position181; _tmp$33 = tokenIndex181; _tmp$34 = depth181; position = _tmp$32; tokenIndex = _tmp$33; depth = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) {} else { $s = 36; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 36:
			position = position + (1) >> 0;
			/* l181: */ case 8:
			_tmp$35 = position; _tmp$36 = tokenIndex; _tmp$37 = depth; position183 = _tmp$35; tokenIndex183 = _tmp$36; depth183 = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 112))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 112))) {} else { $s = 37; continue; }
				/* goto l184 */ $s = 10; continue;
			/* } */ case 37:
			position = position + (1) >> 0;
			/* goto l183 */ $s = 11; continue;
			/* l184: */ case 10:
			_tmp$38 = position183; _tmp$39 = tokenIndex183; _tmp$40 = depth183; position = _tmp$38; tokenIndex = _tmp$39; depth = _tmp$40;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 80))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 80))) {} else { $s = 38; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 38:
			position = position + (1) >> 0;
			/* l183: */ case 11:
			_tmp$41 = position; _tmp$42 = tokenIndex; _tmp$43 = depth; position185 = _tmp$41; tokenIndex185 = _tmp$42; depth185 = _tmp$43;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 39; continue; }
				/* goto l186 */ $s = 12; continue;
			/* } */ case 39:
			position = position + (1) >> 0;
			/* goto l185 */ $s = 13; continue;
			/* l186: */ case 12:
			_tmp$44 = position185; _tmp$45 = tokenIndex185; _tmp$46 = depth185; position = _tmp$44; tokenIndex = _tmp$45; depth = _tmp$46;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 40; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 40:
			position = position + (1) >> 0;
			/* l185: */ case 13:
			_tmp$47 = position; _tmp$48 = tokenIndex; _tmp$49 = depth; position187 = _tmp$47; tokenIndex187 = _tmp$48; depth187 = _tmp$49;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 41; continue; }
				/* goto l188 */ $s = 14; continue;
			/* } */ case 41:
			position = position + (1) >> 0;
			/* goto l187 */ $s = 15; continue;
			/* l188: */ case 14:
			_tmp$50 = position187; _tmp$51 = tokenIndex187; _tmp$52 = depth187; position = _tmp$50; tokenIndex = _tmp$51; depth = _tmp$52;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 42; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 42:
			position = position + (1) >> 0;
			/* l187: */ case 15:
			_tmp$53 = position; _tmp$54 = tokenIndex; _tmp$55 = depth; position189 = _tmp$53; tokenIndex189 = _tmp$54; depth189 = _tmp$55;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) {} else { $s = 43; continue; }
				/* goto l190 */ $s = 16; continue;
			/* } */ case 43:
			position = position + (1) >> 0;
			/* goto l189 */ $s = 17; continue;
			/* l190: */ case 16:
			_tmp$56 = position189; _tmp$57 = tokenIndex189; _tmp$58 = depth189; position = _tmp$56; tokenIndex = _tmp$57; depth = _tmp$58;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) {} else { $s = 44; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 44:
			position = position + (1) >> 0;
			/* l189: */ case 17:
			_tmp$59 = position; _tmp$60 = tokenIndex; _tmp$61 = depth; position191 = _tmp$59; tokenIndex191 = _tmp$60; depth191 = _tmp$61;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) {} else { $s = 45; continue; }
				/* goto l192 */ $s = 18; continue;
			/* } */ case 45:
			position = position + (1) >> 0;
			/* goto l191 */ $s = 19; continue;
			/* l192: */ case 18:
			_tmp$62 = position191; _tmp$63 = tokenIndex191; _tmp$64 = depth191; position = _tmp$62; tokenIndex = _tmp$63; depth = _tmp$64;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) {} else { $s = 46; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 46:
			position = position + (1) >> 0;
			/* l191: */ case 19:
			_tmp$65 = position; _tmp$66 = tokenIndex; _tmp$67 = depth; position193 = _tmp$65; tokenIndex193 = _tmp$66; depth193 = _tmp$67;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) {} else { $s = 47; continue; }
				/* goto l194 */ $s = 20; continue;
			/* } */ case 47:
			position = position + (1) >> 0;
			/* goto l193 */ $s = 21; continue;
			/* l194: */ case 20:
			_tmp$68 = position193; _tmp$69 = tokenIndex193; _tmp$70 = depth193; position = _tmp$68; tokenIndex = _tmp$69; depth = _tmp$70;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) {} else { $s = 48; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 48:
			position = position + (1) >> 0;
			/* l193: */ case 21:
			_tmp$71 = position; _tmp$72 = tokenIndex; _tmp$73 = depth; position195 = _tmp$71; tokenIndex195 = _tmp$72; depth195 = _tmp$73;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) {} else { $s = 49; continue; }
				/* goto l196 */ $s = 22; continue;
			/* } */ case 49:
			position = position + (1) >> 0;
			/* goto l195 */ $s = 23; continue;
			/* l196: */ case 22:
			_tmp$74 = position195; _tmp$75 = tokenIndex195; _tmp$76 = depth195; position = _tmp$74; tokenIndex = _tmp$75; depth = _tmp$76;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) {} else { $s = 50; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 50:
			position = position + (1) >> 0;
			/* l195: */ case 23:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 51; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 51:
			depth = depth - (1) >> 0;
			add(83, position180);
			/* if (!rules[65]()) { */ if (!rules[65]()) {} else { $s = 52; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 52:
			_tmp$77 = position; _tmp$78 = tokenIndex; _tmp$79 = depth; position197 = _tmp$77; tokenIndex197 = _tmp$78; depth197 = _tmp$79;
			/* if (!rules[8]()) { */ if (!rules[8]()) {} else { $s = 53; continue; }
				/* goto l198 */ $s = 24; continue;
			/* } */ case 53:
			/* goto l197 */ $s = 25; continue;
			/* l198: */ case 24:
			_tmp$80 = position197; _tmp$81 = tokenIndex197; _tmp$82 = depth197; position = _tmp$80; tokenIndex = _tmp$81; depth = _tmp$82;
			/* if (!rules[13]()) { */ if (!rules[13]()) {} else { $s = 54; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 54:
			/* l197: */ case 25:
			/* if (!rules[66]()) { */ if (!rules[66]()) {} else { $s = 55; continue; }
				/* goto l178 */ $s = 9; continue;
			/* } */ case 55:
			depth = depth - (1) >> 0;
			add(15, position179);
			/* goto l177 */ $s = 26; continue;
			/* l178: */ case 9:
			_tmp$83 = position177; _tmp$84 = tokenIndex177; _tmp$85 = depth177; position = _tmp$83; tokenIndex = _tmp$84; depth = _tmp$85;
			/* if (!rules[16]()) { */ if (!rules[16]()) {} else { $s = 56; continue; }
				/* goto l174 */ $s = 27; continue;
			/* } */ case 56:
			/* l177: */ case 26:
			depth = depth - (1) >> 0;
			add(14, position176);
			_tmp$86 = position; _tmp$87 = tokenIndex; _tmp$88 = depth; position199 = _tmp$86; tokenIndex199 = _tmp$87; depth199 = _tmp$88;
			/* if (!rules[71]()) { */ if (!rules[71]()) {} else { $s = 57; continue; }
				/* goto l199 */ $s = 28; continue;
			/* } */ case 57:
			/* goto l200 */ $s = 29; continue;
			/* l199: */ case 28:
			_tmp$89 = position199; _tmp$90 = tokenIndex199; _tmp$91 = depth199; position = _tmp$89; tokenIndex = _tmp$90; depth = _tmp$91;
			/* l200: */ case 29:
			/* if (!rules[13]()) { */ if (!rules[13]()) {} else { $s = 58; continue; }
				/* goto l174 */ $s = 27; continue;
			/* } */ case 58:
			/* goto l175 */ $s = 30; continue;
			/* l174: */ case 27:
			_tmp$92 = position174; _tmp$93 = tokenIndex174; _tmp$94 = depth174; position = _tmp$92; tokenIndex = _tmp$93; depth = _tmp$94;
			/* l175: */ case 30:
			depth = depth - (1) >> 0;
			add(13, position165);
			return true;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position203, tokenIndex203, depth203, position204, _tmp$8, _tmp$9, _tmp$10, position205, tokenIndex205, depth205, position207, _tmp$11, _tmp$12, _tmp$13, position208, tokenIndex208, depth208, _tmp$14, _tmp$15, _tmp$16, _tmp$17, _tmp$18, _tmp$19, position210, tokenIndex210, depth210, _tmp$20, _tmp$21, _tmp$22, _tmp$23, _tmp$24, _tmp$25, position212, tokenIndex212, depth212, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, position214, tokenIndex214, depth214, _tmp$32, _tmp$33, _tmp$34, _tmp$35, _tmp$36, _tmp$37, position216, tokenIndex216, depth216, _tmp$38, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position203 = _tmp$5; tokenIndex203 = _tmp$6; depth203 = _tmp$7;
			position204 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[12]()) { */ if (!rules[12]()) {} else { $s = 14; continue; }
				/* goto l203 */ $s = 1; continue;
			/* } */ case 14:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position205 = _tmp$8; tokenIndex205 = _tmp$9; depth205 = _tmp$10;
			position207 = position;
			depth = depth + (1) >> 0;
			_tmp$11 = position; _tmp$12 = tokenIndex; _tmp$13 = depth; position208 = _tmp$11; tokenIndex208 = _tmp$12; depth208 = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 117))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 117))) {} else { $s = 15; continue; }
				/* goto l209 */ $s = 2; continue;
			/* } */ case 15:
			position = position + (1) >> 0;
			/* goto l208 */ $s = 3; continue;
			/* l209: */ case 2:
			_tmp$14 = position208; _tmp$15 = tokenIndex208; _tmp$16 = depth208; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 85))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 85))) {} else { $s = 16; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 16:
			position = position + (1) >> 0;
			/* l208: */ case 3:
			_tmp$17 = position; _tmp$18 = tokenIndex; _tmp$19 = depth; position210 = _tmp$17; tokenIndex210 = _tmp$18; depth210 = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) {} else { $s = 17; continue; }
				/* goto l211 */ $s = 5; continue;
			/* } */ case 17:
			position = position + (1) >> 0;
			/* goto l210 */ $s = 6; continue;
			/* l211: */ case 5:
			_tmp$20 = position210; _tmp$21 = tokenIndex210; _tmp$22 = depth210; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) {} else { $s = 18; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 18:
			position = position + (1) >> 0;
			/* l210: */ case 6:
			_tmp$23 = position; _tmp$24 = tokenIndex; _tmp$25 = depth; position212 = _tmp$23; tokenIndex212 = _tmp$24; depth212 = _tmp$25;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 19; continue; }
				/* goto l213 */ $s = 7; continue;
			/* } */ case 19:
			position = position + (1) >> 0;
			/* goto l212 */ $s = 8; continue;
			/* l213: */ case 7:
			_tmp$26 = position212; _tmp$27 = tokenIndex212; _tmp$28 = depth212; position = _tmp$26; tokenIndex = _tmp$27; depth = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 20; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 20:
			position = position + (1) >> 0;
			/* l212: */ case 8:
			_tmp$29 = position; _tmp$30 = tokenIndex; _tmp$31 = depth; position214 = _tmp$29; tokenIndex214 = _tmp$30; depth214 = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) {} else { $s = 21; continue; }
				/* goto l215 */ $s = 9; continue;
			/* } */ case 21:
			position = position + (1) >> 0;
			/* goto l214 */ $s = 10; continue;
			/* l215: */ case 9:
			_tmp$32 = position214; _tmp$33 = tokenIndex214; _tmp$34 = depth214; position = _tmp$32; tokenIndex = _tmp$33; depth = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) {} else { $s = 22; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 22:
			position = position + (1) >> 0;
			/* l214: */ case 10:
			_tmp$35 = position; _tmp$36 = tokenIndex; _tmp$37 = depth; position216 = _tmp$35; tokenIndex216 = _tmp$36; depth216 = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 110))) {} else { $s = 23; continue; }
				/* goto l217 */ $s = 11; continue;
			/* } */ case 23:
			position = position + (1) >> 0;
			/* goto l216 */ $s = 12; continue;
			/* l217: */ case 11:
			_tmp$38 = position216; _tmp$39 = tokenIndex216; _tmp$40 = depth216; position = _tmp$38; tokenIndex = _tmp$39; depth = _tmp$40;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 78))) {} else { $s = 24; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 24:
			position = position + (1) >> 0;
			/* l216: */ case 12:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 25; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 25:
			depth = depth - (1) >> 0;
			add(84, position207);
			/* if (!rules[16]()) { */ if (!rules[16]()) {} else { $s = 26; continue; }
				/* goto l205 */ $s = 4; continue;
			/* } */ case 26:
			/* goto l206 */ $s = 13; continue;
			/* l205: */ case 4:
			_tmp$41 = position205; _tmp$42 = tokenIndex205; _tmp$43 = depth205; position = _tmp$41; tokenIndex = _tmp$42; depth = _tmp$43;
			/* l206: */ case 13:
			depth = depth - (1) >> 0;
			add(16, position204);
			return true;
			/* l203: */ case 1:
			_tmp$44 = position203; _tmp$45 = tokenIndex203; _tmp$46 = depth203; position = _tmp$44; tokenIndex = _tmp$45; depth = _tmp$46;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position220, tokenIndex220, depth220, position221, _tmp$8, _tmp$9, _tmp$10, position222, tokenIndex222, depth222, position224, _tmp$11, _tmp$12, _tmp$13, position225, tokenIndex225, depth225, position227, _tmp$14, _tmp$15, _tmp$16, position230, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, position233, _tmp$23, _tmp$24, _tmp$25, position234, tokenIndex234, depth234, position236, _tmp$26, _tmp$27, _tmp$28, position238, tokenIndex238, depth238, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, position239, position240, position241, _tmp$35, _tmp$36, _tmp$37;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position220 = _tmp$5; tokenIndex220 = _tmp$6; depth220 = _tmp$7;
			position221 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position222 = _tmp$8; tokenIndex222 = _tmp$9; depth222 = _tmp$10;
			position224 = position;
			depth = depth + (1) >> 0;
			_tmp$11 = position; _tmp$12 = tokenIndex; _tmp$13 = depth; position225 = _tmp$11; tokenIndex225 = _tmp$12; depth225 = _tmp$13;
			position227 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[42]()) { */ if (!rules[42]()) {} else { $s = 11; continue; }
				/* goto l226 */ $s = 1; continue;
			/* } */ case 11:
			depth = depth - (1) >> 0;
			add(89, position227);
			add(90, position);
			/* goto l225 */ $s = 2; continue;
			/* l226: */ case 1:
			_tmp$14 = position225; _tmp$15 = tokenIndex225; _tmp$16 = depth225; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			position230 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[21]()) { */ if (!rules[21]()) {} else { $s = 12; continue; }
				/* goto l229 */ $s = 3; continue;
			/* } */ case 12:
			depth = depth - (1) >> 0;
			add(89, position230);
			add(91, position);
			/* goto l225 */ $s = 2; continue;
			/* l229: */ case 3:
			_tmp$17 = position225; _tmp$18 = tokenIndex225; _tmp$19 = depth225; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* if (!rules[41]()) { */ if (!rules[41]()) {} else { $s = 13; continue; }
				/* goto l223 */ $s = 4; continue;
			/* } */ case 13:
			add(92, position);
			/* l225: */ case 2:
			depth = depth - (1) >> 0;
			add(20, position224);
			/* if (!rules[25]()) { */ if (!rules[25]()) {} else { $s = 14; continue; }
				/* goto l223 */ $s = 4; continue;
			/* } */ case 14:
			/* goto l222 */ $s = 5; continue;
			/* l223: */ case 4:
			_tmp$20 = position222; _tmp$21 = tokenIndex222; _tmp$22 = depth222; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
			position233 = position;
			depth = depth + (1) >> 0;
			_tmp$23 = position; _tmp$24 = tokenIndex; _tmp$25 = depth; position234 = _tmp$23; tokenIndex234 = _tmp$24; depth234 = _tmp$25;
			position236 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[76]()) { */ if (!rules[76]()) {} else { $s = 15; continue; }
				/* goto l235 */ $s = 6; continue;
			/* } */ case 15:
			/* if (!rules[36]()) { */ if (!rules[36]()) {} else { $s = 16; continue; }
				/* goto l235 */ $s = 6; continue;
			/* } */ case 16:
			/* l237: */ case 8:
			_tmp$26 = position; _tmp$27 = tokenIndex; _tmp$28 = depth; position238 = _tmp$26; tokenIndex238 = _tmp$27; depth238 = _tmp$28;
			/* if (!rules[36]()) { */ if (!rules[36]()) {} else { $s = 17; continue; }
				/* goto l238 */ $s = 7; continue;
			/* } */ case 17:
			/* goto l237 */ $s = 8; continue;
			/* l238: */ case 7:
			_tmp$29 = position238; _tmp$30 = tokenIndex238; _tmp$31 = depth238; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
			/* if (!rules[77]()) { */ if (!rules[77]()) {} else { $s = 18; continue; }
				/* goto l235 */ $s = 6; continue;
			/* } */ case 18:
			depth = depth - (1) >> 0;
			add(23, position236);
			/* goto l234 */ $s = 9; continue;
			/* l235: */ case 6:
			_tmp$32 = position234; _tmp$33 = tokenIndex234; _tmp$34 = depth234; position = _tmp$32; tokenIndex = _tmp$33; depth = _tmp$34;
			position239 = position;
			depth = depth + (1) >> 0;
			position240 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 91))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 91))) {} else { $s = 19; continue; }
				/* goto l220 */ $s = 10; continue;
			/* } */ case 19:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 20; continue; }
				/* goto l220 */ $s = 10; continue;
			/* } */ case 20:
			depth = depth - (1) >> 0;
			add(67, position240);
			/* if (!rules[25]()) { */ if (!rules[25]()) {} else { $s = 21; continue; }
				/* goto l220 */ $s = 10; continue;
			/* } */ case 21:
			position241 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 93))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 93))) {} else { $s = 22; continue; }
				/* goto l220 */ $s = 10; continue;
			/* } */ case 22:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 23; continue; }
				/* goto l220 */ $s = 10; continue;
			/* } */ case 23:
			depth = depth - (1) >> 0;
			add(68, position241);
			depth = depth - (1) >> 0;
			add(24, position239);
			/* l234: */ case 9:
			depth = depth - (1) >> 0;
			add(22, position233);
			/* if (!rules[25]()) { */ if (!rules[25]()) {} else { $s = 24; continue; }
				/* goto l220 */ $s = 10; continue;
			/* } */ case 24:
			/* l222: */ case 5:
			depth = depth - (1) >> 0;
			add(19, position221);
			return true;
			/* l220: */ case 10:
			_tmp$35 = position220; _tmp$36 = tokenIndex220; _tmp$37 = depth220; position = _tmp$35; tokenIndex = _tmp$36; depth = _tmp$37;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position243, tokenIndex243, depth243, position244, _ref, position246, position247, _tmp$8, _tmp$9, _tmp$10, position248, tokenIndex248, depth248, position250, _ref$1, c, c$1, c$2, _tmp$11, _tmp$12, _tmp$13, position252, tokenIndex252, depth252, _tmp$14, _tmp$15, _tmp$16, position254, tokenIndex254, depth254, c$3, _tmp$17, _tmp$18, _tmp$19, c$4, _tmp$20, _tmp$21, _tmp$22, c$5, _tmp$23, _tmp$24, _tmp$25, c$6, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, position258, position259, _tmp$32, _tmp$33, _tmp$34, position260, tokenIndex260, depth260, position262, _tmp$35, _tmp$36, _tmp$37, position263, tokenIndex263, depth263, _tmp$38, _tmp$39, _tmp$40, _tmp$41, _tmp$42, _tmp$43, position265, tokenIndex265, depth265, _tmp$44, _tmp$45, _tmp$46, _tmp$47, _tmp$48, _tmp$49, position267, tokenIndex267, depth267, _tmp$50, _tmp$51, _tmp$52, _tmp$53, _tmp$54, _tmp$55, position269, tokenIndex269, depth269, _tmp$56, _tmp$57, _tmp$58, _tmp$59, _tmp$60, _tmp$61, position271, _tmp$62, _tmp$63, _tmp$64, position272, tokenIndex272, depth272, _tmp$65, _tmp$66, _tmp$67, _tmp$68, _tmp$69, _tmp$70, position274, tokenIndex274, depth274, _tmp$71, _tmp$72, _tmp$73, _tmp$74, _tmp$75, _tmp$76, position276, tokenIndex276, depth276, _tmp$77, _tmp$78, _tmp$79, _tmp$80, _tmp$81, _tmp$82, position278, tokenIndex278, depth278, _tmp$83, _tmp$84, _tmp$85, _tmp$86, _tmp$87, _tmp$88, position280, tokenIndex280, depth280, _tmp$89, _tmp$90, _tmp$91, position282, position283, _tmp$92, _tmp$93, _tmp$94, position285, tokenIndex285, depth285, _tmp$95, _tmp$96, _tmp$97, position286, tokenIndex286, depth286, _tmp$98, _tmp$99, _tmp$100, _tmp$101, _tmp$102, _tmp$103, _tmp$104, _tmp$105, _tmp$106, position287, tokenIndex287, depth287, _tmp$107, _tmp$108, _tmp$109, position289, tokenIndex289, depth289, _tmp$110, _tmp$111, _tmp$112, position293, tokenIndex293, depth293, c$7, _tmp$113, _tmp$114, _tmp$115, c$8, _tmp$116, _tmp$117, _tmp$118, position292, tokenIndex292, depth292, _tmp$119, _tmp$120, _tmp$121, position295, tokenIndex295, depth295, c$9, _tmp$122, _tmp$123, _tmp$124, c$10, _tmp$125, _tmp$126, _tmp$127, _tmp$128, _tmp$129, _tmp$130, position298, tokenIndex298, depth298, _ref$2, c$11, c$12, c$13, _tmp$131, _tmp$132, _tmp$133, position300, tokenIndex300, depth300, _ref$3, c$14, c$15, c$16, _tmp$134, _tmp$135, _tmp$136, _tmp$137, _tmp$138, _tmp$139, _tmp$140, _tmp$141, _tmp$142, _tmp$143, _tmp$144, _tmp$145, position303, _tmp$146, _tmp$147, _tmp$148, position304, tokenIndex304, depth304, _tmp$149, _tmp$150, _tmp$151, position306, tokenIndex306, depth306, _tmp$152, _tmp$153, _tmp$154, _tmp$155, _tmp$156, _tmp$157, c$17, _tmp$158, _tmp$159, _tmp$160, position309, tokenIndex309, depth309, c$18, _tmp$161, _tmp$162, _tmp$163, _tmp$164, _tmp$165, _tmp$166, position310, tokenIndex310, depth310, _tmp$167, _tmp$168, _tmp$169, position313, tokenIndex313, depth313, c$19, _tmp$170, _tmp$171, _tmp$172, _tmp$173, _tmp$174, _tmp$175, _tmp$176, _tmp$177, _tmp$178;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position243 = _tmp$5; tokenIndex243 = _tmp$6; depth243 = _tmp$7;
			position244 = position;
			depth = depth + (1) >> 0;
			_ref = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* switch (0) { default: if (_ref === 40) { */ if (_ref === 40) {} else if (_ref === 91 || _ref === 95) { $s = 57; continue; } else if (_ref === 70 || _ref === 84 || _ref === 102 || _ref === 116) { $s = 58; continue; } else if (_ref === 34) { $s = 59; continue; } else if (_ref === 60) { $s = 60; continue; } else { $s = 61; continue; }
				position246 = position;
				depth = depth + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 40))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 40))) {} else { $s = 63; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 63:
				position = position + (1) >> 0;
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 64; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 64:
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 41))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 41))) {} else { $s = 65; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 65:
				position = position + (1) >> 0;
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 66; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 66:
				depth = depth - (1) >> 0;
				add(51, position246);
				/* break; */ $s = 62; continue;
			/* } else if (_ref === 91 || _ref === 95) { */ $s = 62; continue; case 57: 
				position247 = position;
				depth = depth + (1) >> 0;
				_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position248 = _tmp$8; tokenIndex248 = _tmp$9; depth248 = _tmp$10;
				position250 = position;
				depth = depth + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 95))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 95))) {} else { $s = 67; continue; }
					/* goto l249 */ $s = 2; continue;
				/* } */ case 67:
				position = position + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 58))) {} else { $s = 68; continue; }
					/* goto l249 */ $s = 2; continue;
				/* } */ case 68:
				position = position + (1) >> 0;
				_ref$1 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* switch (0) { default: if (_ref$1 === 48 || _ref$1 === 49 || _ref$1 === 50 || _ref$1 === 51 || _ref$1 === 52 || _ref$1 === 53 || _ref$1 === 54 || _ref$1 === 55 || _ref$1 === 56 || _ref$1 === 57) { */ if (_ref$1 === 48 || _ref$1 === 49 || _ref$1 === 50 || _ref$1 === 51 || _ref$1 === 52 || _ref$1 === 53 || _ref$1 === 54 || _ref$1 === 55 || _ref$1 === 56 || _ref$1 === 57) {} else if (_ref$1 === 65 || _ref$1 === 66 || _ref$1 === 67 || _ref$1 === 68 || _ref$1 === 69 || _ref$1 === 70 || _ref$1 === 71 || _ref$1 === 72 || _ref$1 === 73 || _ref$1 === 74 || _ref$1 === 75 || _ref$1 === 76 || _ref$1 === 77 || _ref$1 === 78 || _ref$1 === 79 || _ref$1 === 80 || _ref$1 === 81 || _ref$1 === 82 || _ref$1 === 83 || _ref$1 === 84 || _ref$1 === 85 || _ref$1 === 86 || _ref$1 === 87 || _ref$1 === 88 || _ref$1 === 89 || _ref$1 === 90) { $s = 69; continue; } else { $s = 70; continue; }
					c = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c < 48 || c > 57) { */ if (c < 48 || c > 57) {} else { $s = 72; continue; }
						/* goto l249 */ $s = 2; continue;
					/* } */ case 72:
					position = position + (1) >> 0;
					/* break; */ $s = 71; continue;
				/* } else if (_ref$1 === 65 || _ref$1 === 66 || _ref$1 === 67 || _ref$1 === 68 || _ref$1 === 69 || _ref$1 === 70 || _ref$1 === 71 || _ref$1 === 72 || _ref$1 === 73 || _ref$1 === 74 || _ref$1 === 75 || _ref$1 === 76 || _ref$1 === 77 || _ref$1 === 78 || _ref$1 === 79 || _ref$1 === 80 || _ref$1 === 81 || _ref$1 === 82 || _ref$1 === 83 || _ref$1 === 84 || _ref$1 === 85 || _ref$1 === 86 || _ref$1 === 87 || _ref$1 === 88 || _ref$1 === 89 || _ref$1 === 90) { */ $s = 71; continue; case 69: 
					c$1 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$1 < 65 || c$1 > 90) { */ if (c$1 < 65 || c$1 > 90) {} else { $s = 73; continue; }
						/* goto l249 */ $s = 2; continue;
					/* } */ case 73:
					position = position + (1) >> 0;
					/* break; */ $s = 71; continue;
				/* } else { */ $s = 71; continue; case 70: 
					c$2 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$2 < 97 || c$2 > 122) { */ if (c$2 < 97 || c$2 > 122) {} else { $s = 74; continue; }
						/* goto l249 */ $s = 2; continue;
					/* } */ case 74:
					position = position + (1) >> 0;
					/* break; */ $s = 71; continue;
				/* } } */ case 71:
				_tmp$11 = position; _tmp$12 = tokenIndex; _tmp$13 = depth; position252 = _tmp$11; tokenIndex252 = _tmp$12; depth252 = _tmp$13;
				_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position254 = _tmp$14; tokenIndex254 = _tmp$15; depth254 = _tmp$16;
				c$3 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$3 < 97 || c$3 > 122) { */ if (c$3 < 97 || c$3 > 122) {} else { $s = 75; continue; }
					/* goto l255 */ $s = 3; continue;
				/* } */ case 75:
				position = position + (1) >> 0;
				/* goto l254 */ $s = 4; continue;
				/* l255: */ case 3:
				_tmp$17 = position254; _tmp$18 = tokenIndex254; _tmp$19 = depth254; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
				c$4 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$4 < 65 || c$4 > 90) { */ if (c$4 < 65 || c$4 > 90) {} else { $s = 76; continue; }
					/* goto l256 */ $s = 5; continue;
				/* } */ case 76:
				position = position + (1) >> 0;
				/* goto l254 */ $s = 4; continue;
				/* l256: */ case 5:
				_tmp$20 = position254; _tmp$21 = tokenIndex254; _tmp$22 = depth254; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
				c$5 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$5 < 48 || c$5 > 57) { */ if (c$5 < 48 || c$5 > 57) {} else { $s = 77; continue; }
					/* goto l257 */ $s = 6; continue;
				/* } */ case 77:
				position = position + (1) >> 0;
				/* goto l254 */ $s = 4; continue;
				/* l257: */ case 6:
				_tmp$23 = position254; _tmp$24 = tokenIndex254; _tmp$25 = depth254; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
				c$6 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$6 < 46 || c$6 > 95) { */ if (c$6 < 46 || c$6 > 95) {} else { $s = 78; continue; }
					/* goto l252 */ $s = 7; continue;
				/* } */ case 78:
				position = position + (1) >> 0;
				/* l254: */ case 4:
				/* goto l253 */ $s = 8; continue;
				/* l252: */ case 7:
				_tmp$26 = position252; _tmp$27 = tokenIndex252; _tmp$28 = depth252; position = _tmp$26; tokenIndex = _tmp$27; depth = _tmp$28;
				/* l253: */ case 8:
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 79; continue; }
					/* goto l249 */ $s = 2; continue;
				/* } */ case 79:
				depth = depth - (1) >> 0;
				add(49, position250);
				/* goto l248 */ $s = 9; continue;
				/* l249: */ case 2:
				_tmp$29 = position248; _tmp$30 = tokenIndex248; _tmp$31 = depth248; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
				position258 = position;
				depth = depth + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 91))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 91))) {} else { $s = 80; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 80:
				position = position + (1) >> 0;
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 81; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 81:
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 93))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 93))) {} else { $s = 82; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 82:
				position = position + (1) >> 0;
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 83; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 83:
				depth = depth - (1) >> 0;
				add(50, position258);
				/* l248: */ case 9:
				depth = depth - (1) >> 0;
				add(48, position247);
				/* break; */ $s = 62; continue;
			/* } else if (_ref === 70 || _ref === 84 || _ref === 102 || _ref === 116) { */ $s = 62; continue; case 58: 
				position259 = position;
				depth = depth + (1) >> 0;
				_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position260 = _tmp$32; tokenIndex260 = _tmp$33; depth260 = _tmp$34;
				position262 = position;
				depth = depth + (1) >> 0;
				_tmp$35 = position; _tmp$36 = tokenIndex; _tmp$37 = depth; position263 = _tmp$35; tokenIndex263 = _tmp$36; depth263 = _tmp$37;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 84; continue; }
					/* goto l264 */ $s = 10; continue;
				/* } */ case 84:
				position = position + (1) >> 0;
				/* goto l263 */ $s = 11; continue;
				/* l264: */ case 10:
				_tmp$38 = position263; _tmp$39 = tokenIndex263; _tmp$40 = depth263; position = _tmp$38; tokenIndex = _tmp$39; depth = _tmp$40;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 85; continue; }
					/* goto l261 */ $s = 12; continue;
				/* } */ case 85:
				position = position + (1) >> 0;
				/* l263: */ case 11:
				_tmp$41 = position; _tmp$42 = tokenIndex; _tmp$43 = depth; position265 = _tmp$41; tokenIndex265 = _tmp$42; depth265 = _tmp$43;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 114))) {} else { $s = 86; continue; }
					/* goto l266 */ $s = 13; continue;
				/* } */ case 86:
				position = position + (1) >> 0;
				/* goto l265 */ $s = 14; continue;
				/* l266: */ case 13:
				_tmp$44 = position265; _tmp$45 = tokenIndex265; _tmp$46 = depth265; position = _tmp$44; tokenIndex = _tmp$45; depth = _tmp$46;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 82))) {} else { $s = 87; continue; }
					/* goto l261 */ $s = 12; continue;
				/* } */ case 87:
				position = position + (1) >> 0;
				/* l265: */ case 14:
				_tmp$47 = position; _tmp$48 = tokenIndex; _tmp$49 = depth; position267 = _tmp$47; tokenIndex267 = _tmp$48; depth267 = _tmp$49;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 117))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 117))) {} else { $s = 88; continue; }
					/* goto l268 */ $s = 15; continue;
				/* } */ case 88:
				position = position + (1) >> 0;
				/* goto l267 */ $s = 16; continue;
				/* l268: */ case 15:
				_tmp$50 = position267; _tmp$51 = tokenIndex267; _tmp$52 = depth267; position = _tmp$50; tokenIndex = _tmp$51; depth = _tmp$52;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 85))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 85))) {} else { $s = 89; continue; }
					/* goto l261 */ $s = 12; continue;
				/* } */ case 89:
				position = position + (1) >> 0;
				/* l267: */ case 16:
				_tmp$53 = position; _tmp$54 = tokenIndex; _tmp$55 = depth; position269 = _tmp$53; tokenIndex269 = _tmp$54; depth269 = _tmp$55;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 90; continue; }
					/* goto l270 */ $s = 17; continue;
				/* } */ case 90:
				position = position + (1) >> 0;
				/* goto l269 */ $s = 18; continue;
				/* l270: */ case 17:
				_tmp$56 = position269; _tmp$57 = tokenIndex269; _tmp$58 = depth269; position = _tmp$56; tokenIndex = _tmp$57; depth = _tmp$58;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 91; continue; }
					/* goto l261 */ $s = 12; continue;
				/* } */ case 91:
				position = position + (1) >> 0;
				/* l269: */ case 18:
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 92; continue; }
					/* goto l261 */ $s = 12; continue;
				/* } */ case 92:
				depth = depth - (1) >> 0;
				add(56, position262);
				/* goto l260 */ $s = 19; continue;
				/* l261: */ case 12:
				_tmp$59 = position260; _tmp$60 = tokenIndex260; _tmp$61 = depth260; position = _tmp$59; tokenIndex = _tmp$60; depth = _tmp$61;
				position271 = position;
				depth = depth + (1) >> 0;
				_tmp$62 = position; _tmp$63 = tokenIndex; _tmp$64 = depth; position272 = _tmp$62; tokenIndex272 = _tmp$63; depth272 = _tmp$64;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) {} else { $s = 93; continue; }
					/* goto l273 */ $s = 20; continue;
				/* } */ case 93:
				position = position + (1) >> 0;
				/* goto l272 */ $s = 21; continue;
				/* l273: */ case 20:
				_tmp$65 = position272; _tmp$66 = tokenIndex272; _tmp$67 = depth272; position = _tmp$65; tokenIndex = _tmp$66; depth = _tmp$67;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) {} else { $s = 94; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 94:
				position = position + (1) >> 0;
				/* l272: */ case 21:
				_tmp$68 = position; _tmp$69 = tokenIndex; _tmp$70 = depth; position274 = _tmp$68; tokenIndex274 = _tmp$69; depth274 = _tmp$70;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) {} else { $s = 95; continue; }
					/* goto l275 */ $s = 22; continue;
				/* } */ case 95:
				position = position + (1) >> 0;
				/* goto l274 */ $s = 23; continue;
				/* l275: */ case 22:
				_tmp$71 = position274; _tmp$72 = tokenIndex274; _tmp$73 = depth274; position = _tmp$71; tokenIndex = _tmp$72; depth = _tmp$73;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 65))) {} else { $s = 96; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 96:
				position = position + (1) >> 0;
				/* l274: */ case 23:
				_tmp$74 = position; _tmp$75 = tokenIndex; _tmp$76 = depth; position276 = _tmp$74; tokenIndex276 = _tmp$75; depth276 = _tmp$76;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) {} else { $s = 97; continue; }
					/* goto l277 */ $s = 24; continue;
				/* } */ case 97:
				position = position + (1) >> 0;
				/* goto l276 */ $s = 25; continue;
				/* l277: */ case 24:
				_tmp$77 = position276; _tmp$78 = tokenIndex276; _tmp$79 = depth276; position = _tmp$77; tokenIndex = _tmp$78; depth = _tmp$79;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) {} else { $s = 98; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 98:
				position = position + (1) >> 0;
				/* l276: */ case 25:
				_tmp$80 = position; _tmp$81 = tokenIndex; _tmp$82 = depth; position278 = _tmp$80; tokenIndex278 = _tmp$81; depth278 = _tmp$82;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) {} else { $s = 99; continue; }
					/* goto l279 */ $s = 26; continue;
				/* } */ case 99:
				position = position + (1) >> 0;
				/* goto l278 */ $s = 27; continue;
				/* l279: */ case 26:
				_tmp$83 = position278; _tmp$84 = tokenIndex278; _tmp$85 = depth278; position = _tmp$83; tokenIndex = _tmp$84; depth = _tmp$85;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) {} else { $s = 100; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 100:
				position = position + (1) >> 0;
				/* l278: */ case 27:
				_tmp$86 = position; _tmp$87 = tokenIndex; _tmp$88 = depth; position280 = _tmp$86; tokenIndex280 = _tmp$87; depth280 = _tmp$88;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 101; continue; }
					/* goto l281 */ $s = 28; continue;
				/* } */ case 101:
				position = position + (1) >> 0;
				/* goto l280 */ $s = 29; continue;
				/* l281: */ case 28:
				_tmp$89 = position280; _tmp$90 = tokenIndex280; _tmp$91 = depth280; position = _tmp$89; tokenIndex = _tmp$90; depth = _tmp$91;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 102; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 102:
				position = position + (1) >> 0;
				/* l280: */ case 29:
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 103; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 103:
				depth = depth - (1) >> 0;
				add(57, position271);
				/* l260: */ case 19:
				depth = depth - (1) >> 0;
				add(47, position259);
				/* break; */ $s = 62; continue;
			/* } else if (_ref === 34) { */ $s = 62; continue; case 59: 
				position282 = position;
				depth = depth + (1) >> 0;
				position283 = position;
				depth = depth + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 34))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 34))) {} else { $s = 104; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 104:
				position = position + (1) >> 0;
				/* l284: */ case 32:
				_tmp$92 = position; _tmp$93 = tokenIndex; _tmp$94 = depth; position285 = _tmp$92; tokenIndex285 = _tmp$93; depth285 = _tmp$94;
				_tmp$95 = position; _tmp$96 = tokenIndex; _tmp$97 = depth; position286 = _tmp$95; tokenIndex286 = _tmp$96; depth286 = _tmp$97;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 34))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 34))) {} else { $s = 105; continue; }
					/* goto l286 */ $s = 30; continue;
				/* } */ case 105:
				position = position + (1) >> 0;
				/* goto l285 */ $s = 31; continue;
				/* l286: */ case 30:
				_tmp$98 = position286; _tmp$99 = tokenIndex286; _tmp$100 = depth286; position = _tmp$98; tokenIndex = _tmp$99; depth = _tmp$100;
				/* if (!matchDot()) { */ if (!matchDot()) {} else { $s = 106; continue; }
					/* goto l285 */ $s = 31; continue;
				/* } */ case 106:
				/* goto l284 */ $s = 32; continue;
				/* l285: */ case 31:
				_tmp$101 = position285; _tmp$102 = tokenIndex285; _tmp$103 = depth285; position = _tmp$101; tokenIndex = _tmp$102; depth = _tmp$103;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 34))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 34))) {} else { $s = 107; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 107:
				position = position + (1) >> 0;
				depth = depth - (1) >> 0;
				add(45, position283);
				_tmp$104 = position; _tmp$105 = tokenIndex; _tmp$106 = depth; position287 = _tmp$104; tokenIndex287 = _tmp$105; depth287 = _tmp$106;
				_tmp$107 = position; _tmp$108 = tokenIndex; _tmp$109 = depth; position289 = _tmp$107; tokenIndex289 = _tmp$108; depth289 = _tmp$109;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 64))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 64))) {} else { $s = 108; continue; }
					/* goto l290 */ $s = 33; continue;
				/* } */ case 108:
				position = position + (1) >> 0;
				_tmp$110 = position; _tmp$111 = tokenIndex; _tmp$112 = depth; position293 = _tmp$110; tokenIndex293 = _tmp$111; depth293 = _tmp$112;
				c$7 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$7 < 97 || c$7 > 122) { */ if (c$7 < 97 || c$7 > 122) {} else { $s = 109; continue; }
					/* goto l294 */ $s = 34; continue;
				/* } */ case 109:
				position = position + (1) >> 0;
				/* goto l293 */ $s = 35; continue;
				/* l294: */ case 34:
				_tmp$113 = position293; _tmp$114 = tokenIndex293; _tmp$115 = depth293; position = _tmp$113; tokenIndex = _tmp$114; depth = _tmp$115;
				c$8 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$8 < 65 || c$8 > 90) { */ if (c$8 < 65 || c$8 > 90) {} else { $s = 110; continue; }
					/* goto l290 */ $s = 33; continue;
				/* } */ case 110:
				position = position + (1) >> 0;
				/* l293: */ case 35:
				/* l291: */ case 39:
				_tmp$116 = position; _tmp$117 = tokenIndex; _tmp$118 = depth; position292 = _tmp$116; tokenIndex292 = _tmp$117; depth292 = _tmp$118;
				_tmp$119 = position; _tmp$120 = tokenIndex; _tmp$121 = depth; position295 = _tmp$119; tokenIndex295 = _tmp$120; depth295 = _tmp$121;
				c$9 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$9 < 97 || c$9 > 122) { */ if (c$9 < 97 || c$9 > 122) {} else { $s = 111; continue; }
					/* goto l296 */ $s = 36; continue;
				/* } */ case 111:
				position = position + (1) >> 0;
				/* goto l295 */ $s = 37; continue;
				/* l296: */ case 36:
				_tmp$122 = position295; _tmp$123 = tokenIndex295; _tmp$124 = depth295; position = _tmp$122; tokenIndex = _tmp$123; depth = _tmp$124;
				c$10 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$10 < 65 || c$10 > 90) { */ if (c$10 < 65 || c$10 > 90) {} else { $s = 112; continue; }
					/* goto l292 */ $s = 38; continue;
				/* } */ case 112:
				position = position + (1) >> 0;
				/* l295: */ case 37:
				/* goto l291 */ $s = 39; continue;
				/* l292: */ case 38:
				_tmp$125 = position292; _tmp$126 = tokenIndex292; _tmp$127 = depth292; position = _tmp$125; tokenIndex = _tmp$126; depth = _tmp$127;
				/* l297: */ case 43:
				_tmp$128 = position; _tmp$129 = tokenIndex; _tmp$130 = depth; position298 = _tmp$128; tokenIndex298 = _tmp$129; depth298 = _tmp$130;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 45))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 45))) {} else { $s = 113; continue; }
					/* goto l298 */ $s = 40; continue;
				/* } */ case 113:
				position = position + (1) >> 0;
				_ref$2 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* switch (0) { default: if (_ref$2 === 48 || _ref$2 === 49 || _ref$2 === 50 || _ref$2 === 51 || _ref$2 === 52 || _ref$2 === 53 || _ref$2 === 54 || _ref$2 === 55 || _ref$2 === 56 || _ref$2 === 57) { */ if (_ref$2 === 48 || _ref$2 === 49 || _ref$2 === 50 || _ref$2 === 51 || _ref$2 === 52 || _ref$2 === 53 || _ref$2 === 54 || _ref$2 === 55 || _ref$2 === 56 || _ref$2 === 57) {} else if (_ref$2 === 65 || _ref$2 === 66 || _ref$2 === 67 || _ref$2 === 68 || _ref$2 === 69 || _ref$2 === 70 || _ref$2 === 71 || _ref$2 === 72 || _ref$2 === 73 || _ref$2 === 74 || _ref$2 === 75 || _ref$2 === 76 || _ref$2 === 77 || _ref$2 === 78 || _ref$2 === 79 || _ref$2 === 80 || _ref$2 === 81 || _ref$2 === 82 || _ref$2 === 83 || _ref$2 === 84 || _ref$2 === 85 || _ref$2 === 86 || _ref$2 === 87 || _ref$2 === 88 || _ref$2 === 89 || _ref$2 === 90) { $s = 114; continue; } else { $s = 115; continue; }
					c$11 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$11 < 48 || c$11 > 57) { */ if (c$11 < 48 || c$11 > 57) {} else { $s = 117; continue; }
						/* goto l298 */ $s = 40; continue;
					/* } */ case 117:
					position = position + (1) >> 0;
					/* break; */ $s = 116; continue;
				/* } else if (_ref$2 === 65 || _ref$2 === 66 || _ref$2 === 67 || _ref$2 === 68 || _ref$2 === 69 || _ref$2 === 70 || _ref$2 === 71 || _ref$2 === 72 || _ref$2 === 73 || _ref$2 === 74 || _ref$2 === 75 || _ref$2 === 76 || _ref$2 === 77 || _ref$2 === 78 || _ref$2 === 79 || _ref$2 === 80 || _ref$2 === 81 || _ref$2 === 82 || _ref$2 === 83 || _ref$2 === 84 || _ref$2 === 85 || _ref$2 === 86 || _ref$2 === 87 || _ref$2 === 88 || _ref$2 === 89 || _ref$2 === 90) { */ $s = 116; continue; case 114: 
					c$12 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$12 < 65 || c$12 > 90) { */ if (c$12 < 65 || c$12 > 90) {} else { $s = 118; continue; }
						/* goto l298 */ $s = 40; continue;
					/* } */ case 118:
					position = position + (1) >> 0;
					/* break; */ $s = 116; continue;
				/* } else { */ $s = 116; continue; case 115: 
					c$13 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$13 < 97 || c$13 > 122) { */ if (c$13 < 97 || c$13 > 122) {} else { $s = 119; continue; }
						/* goto l298 */ $s = 40; continue;
					/* } */ case 119:
					position = position + (1) >> 0;
					/* break; */ $s = 116; continue;
				/* } } */ case 116:
				/* l299: */ case 42:
				_tmp$131 = position; _tmp$132 = tokenIndex; _tmp$133 = depth; position300 = _tmp$131; tokenIndex300 = _tmp$132; depth300 = _tmp$133;
				_ref$3 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* switch (0) { default: if (_ref$3 === 48 || _ref$3 === 49 || _ref$3 === 50 || _ref$3 === 51 || _ref$3 === 52 || _ref$3 === 53 || _ref$3 === 54 || _ref$3 === 55 || _ref$3 === 56 || _ref$3 === 57) { */ if (_ref$3 === 48 || _ref$3 === 49 || _ref$3 === 50 || _ref$3 === 51 || _ref$3 === 52 || _ref$3 === 53 || _ref$3 === 54 || _ref$3 === 55 || _ref$3 === 56 || _ref$3 === 57) {} else if (_ref$3 === 65 || _ref$3 === 66 || _ref$3 === 67 || _ref$3 === 68 || _ref$3 === 69 || _ref$3 === 70 || _ref$3 === 71 || _ref$3 === 72 || _ref$3 === 73 || _ref$3 === 74 || _ref$3 === 75 || _ref$3 === 76 || _ref$3 === 77 || _ref$3 === 78 || _ref$3 === 79 || _ref$3 === 80 || _ref$3 === 81 || _ref$3 === 82 || _ref$3 === 83 || _ref$3 === 84 || _ref$3 === 85 || _ref$3 === 86 || _ref$3 === 87 || _ref$3 === 88 || _ref$3 === 89 || _ref$3 === 90) { $s = 120; continue; } else { $s = 121; continue; }
					c$14 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$14 < 48 || c$14 > 57) { */ if (c$14 < 48 || c$14 > 57) {} else { $s = 123; continue; }
						/* goto l300 */ $s = 41; continue;
					/* } */ case 123:
					position = position + (1) >> 0;
					/* break; */ $s = 122; continue;
				/* } else if (_ref$3 === 65 || _ref$3 === 66 || _ref$3 === 67 || _ref$3 === 68 || _ref$3 === 69 || _ref$3 === 70 || _ref$3 === 71 || _ref$3 === 72 || _ref$3 === 73 || _ref$3 === 74 || _ref$3 === 75 || _ref$3 === 76 || _ref$3 === 77 || _ref$3 === 78 || _ref$3 === 79 || _ref$3 === 80 || _ref$3 === 81 || _ref$3 === 82 || _ref$3 === 83 || _ref$3 === 84 || _ref$3 === 85 || _ref$3 === 86 || _ref$3 === 87 || _ref$3 === 88 || _ref$3 === 89 || _ref$3 === 90) { */ $s = 122; continue; case 120: 
					c$15 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$15 < 65 || c$15 > 90) { */ if (c$15 < 65 || c$15 > 90) {} else { $s = 124; continue; }
						/* goto l300 */ $s = 41; continue;
					/* } */ case 124:
					position = position + (1) >> 0;
					/* break; */ $s = 122; continue;
				/* } else { */ $s = 122; continue; case 121: 
					c$16 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
					/* if (c$16 < 97 || c$16 > 122) { */ if (c$16 < 97 || c$16 > 122) {} else { $s = 125; continue; }
						/* goto l300 */ $s = 41; continue;
					/* } */ case 125:
					position = position + (1) >> 0;
					/* break; */ $s = 122; continue;
				/* } } */ case 122:
				/* goto l299 */ $s = 42; continue;
				/* l300: */ case 41:
				_tmp$134 = position300; _tmp$135 = tokenIndex300; _tmp$136 = depth300; position = _tmp$134; tokenIndex = _tmp$135; depth = _tmp$136;
				/* goto l297 */ $s = 43; continue;
				/* l298: */ case 40:
				_tmp$137 = position298; _tmp$138 = tokenIndex298; _tmp$139 = depth298; position = _tmp$137; tokenIndex = _tmp$138; depth = _tmp$139;
				/* goto l289 */ $s = 44; continue;
				/* l290: */ case 33:
				_tmp$140 = position289; _tmp$141 = tokenIndex289; _tmp$142 = depth289; position = _tmp$140; tokenIndex = _tmp$141; depth = _tmp$142;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 94))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 94))) {} else { $s = 126; continue; }
					/* goto l287 */ $s = 45; continue;
				/* } */ case 126:
				position = position + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 94))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 94))) {} else { $s = 127; continue; }
					/* goto l287 */ $s = 45; continue;
				/* } */ case 127:
				position = position + (1) >> 0;
				/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 128; continue; }
					/* goto l287 */ $s = 45; continue;
				/* } */ case 128:
				/* l289: */ case 44:
				/* goto l288 */ $s = 46; continue;
				/* l287: */ case 45:
				_tmp$143 = position287; _tmp$144 = tokenIndex287; _tmp$145 = depth287; position = _tmp$143; tokenIndex = _tmp$144; depth = _tmp$145;
				/* l288: */ case 46:
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 129; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 129:
				depth = depth - (1) >> 0;
				add(44, position282);
				/* break; */ $s = 62; continue;
			/* } else if (_ref === 60) { */ $s = 62; continue; case 60: 
				/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 130; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 130:
				/* break; */ $s = 62; continue;
			/* } else { */ $s = 62; continue; case 61: 
				position303 = position;
				depth = depth + (1) >> 0;
				_tmp$146 = position; _tmp$147 = tokenIndex; _tmp$148 = depth; position304 = _tmp$146; tokenIndex304 = _tmp$147; depth304 = _tmp$148;
				_tmp$149 = position; _tmp$150 = tokenIndex; _tmp$151 = depth; position306 = _tmp$149; tokenIndex306 = _tmp$150; depth306 = _tmp$151;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 43))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 43))) {} else { $s = 131; continue; }
					/* goto l307 */ $s = 47; continue;
				/* } */ case 131:
				position = position + (1) >> 0;
				/* goto l306 */ $s = 48; continue;
				/* l307: */ case 47:
				_tmp$152 = position306; _tmp$153 = tokenIndex306; _tmp$154 = depth306; position = _tmp$152; tokenIndex = _tmp$153; depth = _tmp$154;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 45))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 45))) {} else { $s = 132; continue; }
					/* goto l304 */ $s = 49; continue;
				/* } */ case 132:
				position = position + (1) >> 0;
				/* l306: */ case 48:
				/* goto l305 */ $s = 50; continue;
				/* l304: */ case 49:
				_tmp$155 = position304; _tmp$156 = tokenIndex304; _tmp$157 = depth304; position = _tmp$155; tokenIndex = _tmp$156; depth = _tmp$157;
				/* l305: */ case 50:
				c$17 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$17 < 48 || c$17 > 57) { */ if (c$17 < 48 || c$17 > 57) {} else { $s = 133; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 133:
				position = position + (1) >> 0;
				/* l308: */ case 52:
				_tmp$158 = position; _tmp$159 = tokenIndex; _tmp$160 = depth; position309 = _tmp$158; tokenIndex309 = _tmp$159; depth309 = _tmp$160;
				c$18 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$18 < 48 || c$18 > 57) { */ if (c$18 < 48 || c$18 > 57) {} else { $s = 134; continue; }
					/* goto l309 */ $s = 51; continue;
				/* } */ case 134:
				position = position + (1) >> 0;
				/* goto l308 */ $s = 52; continue;
				/* l309: */ case 51:
				_tmp$161 = position309; _tmp$162 = tokenIndex309; _tmp$163 = depth309; position = _tmp$161; tokenIndex = _tmp$162; depth = _tmp$163;
				_tmp$164 = position; _tmp$165 = tokenIndex; _tmp$166 = depth; position310 = _tmp$164; tokenIndex310 = _tmp$165; depth310 = _tmp$166;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 46))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 46))) {} else { $s = 135; continue; }
					/* goto l310 */ $s = 53; continue;
				/* } */ case 135:
				position = position + (1) >> 0;
				/* l312: */ case 55:
				_tmp$167 = position; _tmp$168 = tokenIndex; _tmp$169 = depth; position313 = _tmp$167; tokenIndex313 = _tmp$168; depth313 = _tmp$169;
				c$19 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
				/* if (c$19 < 48 || c$19 > 57) { */ if (c$19 < 48 || c$19 > 57) {} else { $s = 136; continue; }
					/* goto l313 */ $s = 54; continue;
				/* } */ case 136:
				position = position + (1) >> 0;
				/* goto l312 */ $s = 55; continue;
				/* l313: */ case 54:
				_tmp$170 = position313; _tmp$171 = tokenIndex313; _tmp$172 = depth313; position = _tmp$170; tokenIndex = _tmp$171; depth = _tmp$172;
				/* goto l311 */ $s = 56; continue;
				/* l310: */ case 53:
				_tmp$173 = position310; _tmp$174 = tokenIndex310; _tmp$175 = depth310; position = _tmp$173; tokenIndex = _tmp$174; depth = _tmp$175;
				/* l311: */ case 56:
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 137; continue; }
					/* goto l243 */ $s = 1; continue;
				/* } */ case 137:
				depth = depth - (1) >> 0;
				add(46, position303);
				/* break; */ $s = 62; continue;
			/* } } */ case 62:
			depth = depth - (1) >> 0;
			add(21, position244);
			return true;
			/* l243: */ case 1:
			_tmp$176 = position243; _tmp$177 = tokenIndex243; _tmp$178 = depth243; position = _tmp$176; tokenIndex = _tmp$177; depth = _tmp$178;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position317, tokenIndex317, depth317, position318, _tmp$8, _tmp$9, _tmp$10, position319, tokenIndex319, depth319, _tmp$11, _tmp$12, _tmp$13, position323, _tmp$14, _tmp$15, _tmp$16, position325, _tmp$17, _tmp$18, _tmp$19, position326, tokenIndex326, depth326, position328, _tmp$20, _tmp$21, _tmp$22, _tmp$23, _tmp$24, _tmp$25;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position317 = _tmp$5; tokenIndex317 = _tmp$6; depth317 = _tmp$7;
			position318 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position319 = _tmp$8; tokenIndex319 = _tmp$9; depth319 = _tmp$10;
			/* if (!rules[41]()) { */ if (!rules[41]()) {} else { $s = 7; continue; }
				/* goto l320 */ $s = 1; continue;
			/* } */ case 7:
			add(93, position);
			/* goto l319 */ $s = 2; continue;
			/* l320: */ case 1:
			_tmp$11 = position319; _tmp$12 = tokenIndex319; _tmp$13 = depth319; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			position323 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[42]()) { */ if (!rules[42]()) {} else { $s = 8; continue; }
				/* goto l322 */ $s = 3; continue;
			/* } */ case 8:
			depth = depth - (1) >> 0;
			add(89, position323);
			add(94, position);
			/* goto l319 */ $s = 2; continue;
			/* l322: */ case 3:
			_tmp$14 = position319; _tmp$15 = tokenIndex319; _tmp$16 = depth319; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			position325 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[27]()) { */ if (!rules[27]()) {} else { $s = 9; continue; }
				/* goto l317 */ $s = 4; continue;
			/* } */ case 9:
			depth = depth - (1) >> 0;
			add(26, position325);
			/* l319: */ case 2:
			/* if (!rules[34]()) { */ if (!rules[34]()) {} else { $s = 10; continue; }
				/* goto l317 */ $s = 4; continue;
			/* } */ case 10:
			_tmp$17 = position; _tmp$18 = tokenIndex; _tmp$19 = depth; position326 = _tmp$17; tokenIndex326 = _tmp$18; depth326 = _tmp$19;
			position328 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 59))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 59))) {} else { $s = 11; continue; }
				/* goto l326 */ $s = 5; continue;
			/* } */ case 11:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 12; continue; }
				/* goto l326 */ $s = 5; continue;
			/* } */ case 12:
			depth = depth - (1) >> 0;
			add(69, position328);
			/* if (!rules[25]()) { */ if (!rules[25]()) {} else { $s = 13; continue; }
				/* goto l326 */ $s = 5; continue;
			/* } */ case 13:
			/* goto l327 */ $s = 6; continue;
			/* l326: */ case 5:
			_tmp$20 = position326; _tmp$21 = tokenIndex326; _tmp$22 = depth326; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
			/* l327: */ case 6:
			depth = depth - (1) >> 0;
			add(25, position318);
			return true;
			/* l317: */ case 4:
			_tmp$23 = position317; _tmp$24 = tokenIndex317; _tmp$25 = depth317; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position330, tokenIndex330, depth330, position331, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position330 = _tmp$5; tokenIndex330 = _tmp$6; depth330 = _tmp$7;
			position331 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[28]()) { */ if (!rules[28]()) {} else { $s = 2; continue; }
				/* goto l330 */ $s = 1; continue;
			/* } */ case 2:
			depth = depth - (1) >> 0;
			add(27, position331);
			return true;
			/* l330: */ case 1:
			_tmp$8 = position330; _tmp$9 = tokenIndex330; _tmp$10 = depth330; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position332, tokenIndex332, depth332, position333, _tmp$8, _tmp$9, _tmp$10, position335, tokenIndex335, depth335, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position332 = _tmp$5; tokenIndex332 = _tmp$6; depth332 = _tmp$7;
			position333 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[29]()) { */ if (!rules[29]()) {} else { $s = 4; continue; }
				/* goto l332 */ $s = 1; continue;
			/* } */ case 4:
			/* l334: */ case 3:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position335 = _tmp$8; tokenIndex335 = _tmp$9; depth335 = _tmp$10;
			/* if (!rules[73]()) { */ if (!rules[73]()) {} else { $s = 5; continue; }
				/* goto l335 */ $s = 2; continue;
			/* } */ case 5:
			/* if (!rules[28]()) { */ if (!rules[28]()) {} else { $s = 6; continue; }
				/* goto l335 */ $s = 2; continue;
			/* } */ case 6:
			/* goto l334 */ $s = 3; continue;
			/* l335: */ case 2:
			_tmp$11 = position335; _tmp$12 = tokenIndex335; _tmp$13 = depth335; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			depth = depth - (1) >> 0;
			add(28, position333);
			return true;
			/* l332: */ case 1:
			_tmp$14 = position332; _tmp$15 = tokenIndex332; _tmp$16 = depth332; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position336, tokenIndex336, depth336, position337, position338, position339, _tmp$8, _tmp$9, _tmp$10, position340, tokenIndex340, depth340, _tmp$11, _tmp$12, _tmp$13, position342, _ref, position344, position345, _tmp$14, _tmp$15, _tmp$16, position346, tokenIndex346, depth346, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, position348, tokenIndex348, depth348, _tmp$23, _tmp$24, _tmp$25, position351, tokenIndex351, depth351, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, position354, tokenIndex354, depth354, position355, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position336 = _tmp$5; tokenIndex336 = _tmp$6; depth336 = _tmp$7;
			position337 = position;
			depth = depth + (1) >> 0;
			position338 = position;
			depth = depth + (1) >> 0;
			position339 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position340 = _tmp$8; tokenIndex340 = _tmp$9; depth340 = _tmp$10;
			/* if (!rules[75]()) { */ if (!rules[75]()) {} else { $s = 12; continue; }
				/* goto l340 */ $s = 1; continue;
			/* } */ case 12:
			/* goto l341 */ $s = 2; continue;
			/* l340: */ case 1:
			_tmp$11 = position340; _tmp$12 = tokenIndex340; _tmp$13 = depth340; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* l341: */ case 2:
			position342 = position;
			depth = depth + (1) >> 0;
			_ref = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* switch (0) { default: if (_ref === 40) { */ if (_ref === 40) {} else if (_ref === 33) { $s = 13; continue; } else if (_ref === 97) { $s = 14; continue; } else { $s = 15; continue; }
				/* if (!rules[76]()) { */ if (!rules[76]()) {} else { $s = 17; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 17:
				/* if (!rules[27]()) { */ if (!rules[27]()) {} else { $s = 18; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 18:
				/* if (!rules[77]()) { */ if (!rules[77]()) {} else { $s = 19; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 19:
				/* break; */ $s = 16; continue;
			/* } else if (_ref === 33) { */ $s = 16; continue; case 13: 
				position344 = position;
				depth = depth + (1) >> 0;
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 33))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 33))) {} else { $s = 20; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 20:
				position = position + (1) >> 0;
				/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 21; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 21:
				depth = depth - (1) >> 0;
				add(79, position344);
				position345 = position;
				depth = depth + (1) >> 0;
				_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position346 = _tmp$14; tokenIndex346 = _tmp$15; depth346 = _tmp$16;
				/* if (!rules[33]()) { */ if (!rules[33]()) {} else { $s = 22; continue; }
					/* goto l347 */ $s = 4; continue;
				/* } */ case 22:
				/* goto l346 */ $s = 5; continue;
				/* l347: */ case 4:
				_tmp$17 = position346; _tmp$18 = tokenIndex346; _tmp$19 = depth346; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
				/* if (!rules[76]()) { */ if (!rules[76]()) {} else { $s = 23; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 23:
				_tmp$20 = position; _tmp$21 = tokenIndex; _tmp$22 = depth; position348 = _tmp$20; tokenIndex348 = _tmp$21; depth348 = _tmp$22;
				/* if (!rules[33]()) { */ if (!rules[33]()) {} else { $s = 24; continue; }
					/* goto l348 */ $s = 6; continue;
				/* } */ case 24:
				/* l350: */ case 8:
				_tmp$23 = position; _tmp$24 = tokenIndex; _tmp$25 = depth; position351 = _tmp$23; tokenIndex351 = _tmp$24; depth351 = _tmp$25;
				/* if (!rules[73]()) { */ if (!rules[73]()) {} else { $s = 25; continue; }
					/* goto l351 */ $s = 7; continue;
				/* } */ case 25:
				/* if (!rules[33]()) { */ if (!rules[33]()) {} else { $s = 26; continue; }
					/* goto l351 */ $s = 7; continue;
				/* } */ case 26:
				/* goto l350 */ $s = 8; continue;
				/* l351: */ case 7:
				_tmp$26 = position351; _tmp$27 = tokenIndex351; _tmp$28 = depth351; position = _tmp$26; tokenIndex = _tmp$27; depth = _tmp$28;
				/* goto l349 */ $s = 9; continue;
				/* l348: */ case 6:
				_tmp$29 = position348; _tmp$30 = tokenIndex348; _tmp$31 = depth348; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
				/* l349: */ case 9:
				/* if (!rules[77]()) { */ if (!rules[77]()) {} else { $s = 27; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 27:
				/* l346: */ case 5:
				depth = depth - (1) >> 0;
				add(32, position345);
				/* break; */ $s = 16; continue;
			/* } else if (_ref === 97) { */ $s = 16; continue; case 14: 
				/* if (!rules[78]()) { */ if (!rules[78]()) {} else { $s = 28; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 28:
				/* break; */ $s = 16; continue;
			/* } else { */ $s = 16; continue; case 15: 
				/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 29; continue; }
					/* goto l336 */ $s = 3; continue;
				/* } */ case 29:
				/* break; */ $s = 16; continue;
			/* } } */ case 16:
			depth = depth - (1) >> 0;
			add(31, position342);
			depth = depth - (1) >> 0;
			add(30, position339);
			depth = depth - (1) >> 0;
			add(89, position338);
			add(95, position);
			/* l353: */ case 11:
			_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position354 = _tmp$32; tokenIndex354 = _tmp$33; depth354 = _tmp$34;
			position355 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 47))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 47))) {} else { $s = 30; continue; }
				/* goto l354 */ $s = 10; continue;
			/* } */ case 30:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 31; continue; }
				/* goto l354 */ $s = 10; continue;
			/* } */ case 31:
			depth = depth - (1) >> 0;
			add(74, position355);
			/* if (!rules[29]()) { */ if (!rules[29]()) {} else { $s = 32; continue; }
				/* goto l354 */ $s = 10; continue;
			/* } */ case 32:
			/* goto l353 */ $s = 11; continue;
			/* l354: */ case 10:
			_tmp$35 = position354; _tmp$36 = tokenIndex354; _tmp$37 = depth354; position = _tmp$35; tokenIndex = _tmp$36; depth = _tmp$37;
			depth = depth - (1) >> 0;
			add(29, position337);
			return true;
			/* l336: */ case 3:
			_tmp$38 = position336; _tmp$39 = tokenIndex336; _tmp$40 = depth336; position = _tmp$38; tokenIndex = _tmp$39; depth = _tmp$40;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position359, tokenIndex359, depth359, position360, _ref, _tmp$8, _tmp$9, _tmp$10, position362, tokenIndex362, depth362, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position359 = _tmp$5; tokenIndex359 = _tmp$6; depth359 = _tmp$7;
			position360 = position;
			depth = depth + (1) >> 0;
			_ref = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* switch (0) { default: if (_ref === 94) { */ if (_ref === 94) {} else if (_ref === 97) { $s = 4; continue; } else { $s = 5; continue; }
				/* if (!rules[75]()) { */ if (!rules[75]()) {} else { $s = 7; continue; }
					/* goto l359 */ $s = 1; continue;
				/* } */ case 7:
				_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position362 = _tmp$8; tokenIndex362 = _tmp$9; depth362 = _tmp$10;
				/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 8; continue; }
					/* goto l363 */ $s = 2; continue;
				/* } */ case 8:
				/* goto l362 */ $s = 3; continue;
				/* l363: */ case 2:
				_tmp$11 = position362; _tmp$12 = tokenIndex362; _tmp$13 = depth362; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
				/* if (!rules[78]()) { */ if (!rules[78]()) {} else { $s = 9; continue; }
					/* goto l359 */ $s = 1; continue;
				/* } */ case 9:
				/* l362: */ case 3:
				/* break; */ $s = 6; continue;
			/* } else if (_ref === 97) { */ $s = 6; continue; case 4: 
				/* if (!rules[78]()) { */ if (!rules[78]()) {} else { $s = 10; continue; }
					/* goto l359 */ $s = 1; continue;
				/* } */ case 10:
				/* break; */ $s = 6; continue;
			/* } else { */ $s = 6; continue; case 5: 
				/* if (!rules[43]()) { */ if (!rules[43]()) {} else { $s = 11; continue; }
					/* goto l359 */ $s = 1; continue;
				/* } */ case 11:
				/* break; */ $s = 6; continue;
			/* } } */ case 6:
			depth = depth - (1) >> 0;
			add(33, position360);
			return true;
			/* l359: */ case 1:
			_tmp$14 = position359; _tmp$15 = tokenIndex359; _tmp$16 = depth359; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, position365, position366, _tmp$5, _tmp$6, _tmp$7, position367, tokenIndex367, depth367, position369, _tmp$8, _tmp$9, _tmp$10, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, position375, tokenIndex375, depth375, position376, _tmp$17, _tmp$18, _tmp$19;
			/* */ while (true) { switch ($s) { case 0:
			position365 = position;
			depth = depth + (1) >> 0;
			position366 = position;
			depth = depth + (1) >> 0;
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position367 = _tmp$5; tokenIndex367 = _tmp$6; depth367 = _tmp$7;
			position369 = position;
			depth = depth + (1) >> 0;
			/* if (!rules[36]()) { */ if (!rules[36]()) {} else { $s = 6; continue; }
				/* goto l368 */ $s = 1; continue;
			/* } */ case 6:
			depth = depth - (1) >> 0;
			add(89, position369);
			add(96, position);
			/* goto l367 */ $s = 2; continue;
			/* l368: */ case 1:
			_tmp$8 = position367; _tmp$9 = tokenIndex367; _tmp$10 = depth367; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			/* if (!rules[41]()) { */ if (!rules[41]()) {} else { $s = 7; continue; }
				/* goto l371 */ $s = 3; continue;
			/* } */ case 7:
			add(97, position);
			/* goto l367 */ $s = 2; continue;
			/* l371: */ case 3:
			_tmp$11 = position367; _tmp$12 = tokenIndex367; _tmp$13 = depth367; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			add(98, position);
			/* l367: */ case 2:
			depth = depth - (1) >> 0;
			add(35, position366);
			/* l374: */ case 5:
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position375 = _tmp$14; tokenIndex375 = _tmp$15; depth375 = _tmp$16;
			position376 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 44))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 44))) {} else { $s = 8; continue; }
				/* goto l375 */ $s = 4; continue;
			/* } */ case 8:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 9; continue; }
				/* goto l375 */ $s = 4; continue;
			/* } */ case 9:
			depth = depth - (1) >> 0;
			add(70, position376);
			/* if (!rules[34]()) { */ if (!rules[34]()) {} else { $s = 10; continue; }
				/* goto l375 */ $s = 4; continue;
			/* } */ case 10:
			/* goto l374 */ $s = 5; continue;
			/* l375: */ case 4:
			_tmp$17 = position375; _tmp$18 = tokenIndex375; _tmp$19 = depth375; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			depth = depth - (1) >> 0;
			add(34, position365);
			return true;
			/* */ case -1: } return; }
		}), $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position378, tokenIndex378, depth378, position379, _tmp$8, _tmp$9, _tmp$10, position380, tokenIndex380, depth380, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position378 = _tmp$5; tokenIndex378 = _tmp$6; depth378 = _tmp$7;
			position379 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position380 = _tmp$8; tokenIndex380 = _tmp$9; depth380 = _tmp$10;
			/* if (!rules[42]()) { */ if (!rules[42]()) {} else { $s = 4; continue; }
				/* goto l381 */ $s = 1; continue;
			/* } */ case 4:
			/* goto l380 */ $s = 2; continue;
			/* l381: */ case 1:
			_tmp$11 = position380; _tmp$12 = tokenIndex380; _tmp$13 = depth380; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!rules[21]()) { */ if (!rules[21]()) {} else { $s = 5; continue; }
				/* goto l378 */ $s = 3; continue;
			/* } */ case 5:
			/* l380: */ case 2:
			depth = depth - (1) >> 0;
			add(36, position379);
			return true;
			/* l378: */ case 3:
			_tmp$14 = position378; _tmp$15 = tokenIndex378; _tmp$16 = depth378; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position384, tokenIndex384, depth384, position385, position386, _tmp$8, _tmp$9, _tmp$10, position387, tokenIndex387, depth387, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, position389, tokenIndex389, depth389, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, position391, tokenIndex391, depth391, _tmp$23, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, position393, tokenIndex393, depth393, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, position395, tokenIndex395, depth395, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position384 = _tmp$5; tokenIndex384 = _tmp$6; depth384 = _tmp$7;
			position385 = position;
			depth = depth + (1) >> 0;
			position386 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position387 = _tmp$8; tokenIndex387 = _tmp$9; depth387 = _tmp$10;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 108))) {} else { $s = 12; continue; }
				/* goto l388 */ $s = 1; continue;
			/* } */ case 12:
			position = position + (1) >> 0;
			/* goto l387 */ $s = 2; continue;
			/* l388: */ case 1:
			_tmp$11 = position387; _tmp$12 = tokenIndex387; _tmp$13 = depth387; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 76))) {} else { $s = 13; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 13:
			position = position + (1) >> 0;
			/* l387: */ case 2:
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position389 = _tmp$14; tokenIndex389 = _tmp$15; depth389 = _tmp$16;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 14; continue; }
				/* goto l390 */ $s = 4; continue;
			/* } */ case 14:
			position = position + (1) >> 0;
			/* goto l389 */ $s = 5; continue;
			/* l390: */ case 4:
			_tmp$17 = position389; _tmp$18 = tokenIndex389; _tmp$19 = depth389; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 15; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 15:
			position = position + (1) >> 0;
			/* l389: */ case 5:
			_tmp$20 = position; _tmp$21 = tokenIndex; _tmp$22 = depth; position391 = _tmp$20; tokenIndex391 = _tmp$21; depth391 = _tmp$22;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 109))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 109))) {} else { $s = 16; continue; }
				/* goto l392 */ $s = 6; continue;
			/* } */ case 16:
			position = position + (1) >> 0;
			/* goto l391 */ $s = 7; continue;
			/* l392: */ case 6:
			_tmp$23 = position391; _tmp$24 = tokenIndex391; _tmp$25 = depth391; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 77))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 77))) {} else { $s = 17; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 17:
			position = position + (1) >> 0;
			/* l391: */ case 7:
			_tmp$26 = position; _tmp$27 = tokenIndex; _tmp$28 = depth; position393 = _tmp$26; tokenIndex393 = _tmp$27; depth393 = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 105))) {} else { $s = 18; continue; }
				/* goto l394 */ $s = 8; continue;
			/* } */ case 18:
			position = position + (1) >> 0;
			/* goto l393 */ $s = 9; continue;
			/* l394: */ case 8:
			_tmp$29 = position393; _tmp$30 = tokenIndex393; _tmp$31 = depth393; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 73))) {} else { $s = 19; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 19:
			position = position + (1) >> 0;
			/* l393: */ case 9:
			_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position395 = _tmp$32; tokenIndex395 = _tmp$33; depth395 = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 20; continue; }
				/* goto l396 */ $s = 10; continue;
			/* } */ case 20:
			position = position + (1) >> 0;
			/* goto l395 */ $s = 11; continue;
			/* l396: */ case 10:
			_tmp$35 = position395; _tmp$36 = tokenIndex395; _tmp$37 = depth395; position = _tmp$35; tokenIndex = _tmp$36; depth = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 21; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 21:
			position = position + (1) >> 0;
			/* l395: */ case 11:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 22; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 22:
			depth = depth - (1) >> 0;
			add(85, position386);
			/* if (!rules[87]()) { */ if (!rules[87]()) {} else { $s = 23; continue; }
				/* goto l384 */ $s = 3; continue;
			/* } */ case 23:
			depth = depth - (1) >> 0;
			add(39, position385);
			return true;
			/* l384: */ case 3:
			_tmp$38 = position384; _tmp$39 = tokenIndex384; _tmp$40 = depth384; position = _tmp$38; tokenIndex = _tmp$39; depth = _tmp$40;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position397, tokenIndex397, depth397, position398, position399, _tmp$8, _tmp$9, _tmp$10, position400, tokenIndex400, depth400, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16, position402, tokenIndex402, depth402, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22, position404, tokenIndex404, depth404, _tmp$23, _tmp$24, _tmp$25, _tmp$26, _tmp$27, _tmp$28, position406, tokenIndex406, depth406, _tmp$29, _tmp$30, _tmp$31, _tmp$32, _tmp$33, _tmp$34, position408, tokenIndex408, depth408, _tmp$35, _tmp$36, _tmp$37, _tmp$38, _tmp$39, _tmp$40, position410, tokenIndex410, depth410, _tmp$41, _tmp$42, _tmp$43, _tmp$44, _tmp$45, _tmp$46;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position397 = _tmp$5; tokenIndex397 = _tmp$6; depth397 = _tmp$7;
			position398 = position;
			depth = depth + (1) >> 0;
			position399 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position400 = _tmp$8; tokenIndex400 = _tmp$9; depth400 = _tmp$10;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 111))) {} else { $s = 14; continue; }
				/* goto l401 */ $s = 1; continue;
			/* } */ case 14:
			position = position + (1) >> 0;
			/* goto l400 */ $s = 2; continue;
			/* l401: */ case 1:
			_tmp$11 = position400; _tmp$12 = tokenIndex400; _tmp$13 = depth400; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 79))) {} else { $s = 15; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 15:
			position = position + (1) >> 0;
			/* l400: */ case 2:
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position402 = _tmp$14; tokenIndex402 = _tmp$15; depth402 = _tmp$16;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) {} else { $s = 16; continue; }
				/* goto l403 */ $s = 4; continue;
			/* } */ case 16:
			position = position + (1) >> 0;
			/* goto l402 */ $s = 5; continue;
			/* l403: */ case 4:
			_tmp$17 = position402; _tmp$18 = tokenIndex402; _tmp$19 = depth402; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) {} else { $s = 17; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 17:
			position = position + (1) >> 0;
			/* l402: */ case 5:
			_tmp$20 = position; _tmp$21 = tokenIndex; _tmp$22 = depth; position404 = _tmp$20; tokenIndex404 = _tmp$21; depth404 = _tmp$22;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 102))) {} else { $s = 18; continue; }
				/* goto l405 */ $s = 6; continue;
			/* } */ case 18:
			position = position + (1) >> 0;
			/* goto l404 */ $s = 7; continue;
			/* l405: */ case 6:
			_tmp$23 = position404; _tmp$24 = tokenIndex404; _tmp$25 = depth404; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 70))) {} else { $s = 19; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 19:
			position = position + (1) >> 0;
			/* l404: */ case 7:
			_tmp$26 = position; _tmp$27 = tokenIndex; _tmp$28 = depth; position406 = _tmp$26; tokenIndex406 = _tmp$27; depth406 = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 115))) {} else { $s = 20; continue; }
				/* goto l407 */ $s = 8; continue;
			/* } */ case 20:
			position = position + (1) >> 0;
			/* goto l406 */ $s = 9; continue;
			/* l407: */ case 8:
			_tmp$29 = position406; _tmp$30 = tokenIndex406; _tmp$31 = depth406; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 83))) {} else { $s = 21; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 21:
			position = position + (1) >> 0;
			/* l406: */ case 9:
			_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position408 = _tmp$32; tokenIndex408 = _tmp$33; depth408 = _tmp$34;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 101))) {} else { $s = 22; continue; }
				/* goto l409 */ $s = 10; continue;
			/* } */ case 22:
			position = position + (1) >> 0;
			/* goto l408 */ $s = 11; continue;
			/* l409: */ case 10:
			_tmp$35 = position408; _tmp$36 = tokenIndex408; _tmp$37 = depth408; position = _tmp$35; tokenIndex = _tmp$36; depth = _tmp$37;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 69))) {} else { $s = 23; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 23:
			position = position + (1) >> 0;
			/* l408: */ case 11:
			_tmp$38 = position; _tmp$39 = tokenIndex; _tmp$40 = depth; position410 = _tmp$38; tokenIndex410 = _tmp$39; depth410 = _tmp$40;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 116))) {} else { $s = 24; continue; }
				/* goto l411 */ $s = 12; continue;
			/* } */ case 24:
			position = position + (1) >> 0;
			/* goto l410 */ $s = 13; continue;
			/* l411: */ case 12:
			_tmp$41 = position410; _tmp$42 = tokenIndex410; _tmp$43 = depth410; position = _tmp$41; tokenIndex = _tmp$42; depth = _tmp$43;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 84))) {} else { $s = 25; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 25:
			position = position + (1) >> 0;
			/* l410: */ case 13:
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 26; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 26:
			depth = depth - (1) >> 0;
			add(86, position399);
			/* if (!rules[87]()) { */ if (!rules[87]()) {} else { $s = 27; continue; }
				/* goto l397 */ $s = 3; continue;
			/* } */ case 27:
			depth = depth - (1) >> 0;
			add(40, position398);
			return true;
			/* l397: */ case 3:
			_tmp$44 = position397; _tmp$45 = tokenIndex397; _tmp$46 = depth397; position = _tmp$44; tokenIndex = _tmp$45; depth = _tmp$46;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position412, tokenIndex412, depth412, position413, _ref, _tmp$8, _tmp$9, _tmp$10, position415, tokenIndex415, depth415, _ref$1, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position412 = _tmp$5; tokenIndex412 = _tmp$6; depth412 = _tmp$7;
			position413 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 60))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 60))) {} else { $s = 4; continue; }
				/* goto l412 */ $s = 1; continue;
			/* } */ case 4:
			position = position + (1) >> 0;
			_ref = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* switch (0) { default: if (_ref === 12) { */ if (_ref === 12) {} else if (_ref === 13) { $s = 5; continue; } else if (_ref === 10) { $s = 6; continue; } else if (_ref === 9) { $s = 7; continue; } else { $s = 8; continue; }
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 12))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 12))) {} else { $s = 10; continue; }
					/* goto l412 */ $s = 1; continue;
				/* } */ case 10:
				position = position + (1) >> 0;
				/* break; */ $s = 9; continue;
			/* } else if (_ref === 13) { */ $s = 9; continue; case 5: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 13))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 13))) {} else { $s = 11; continue; }
					/* goto l412 */ $s = 1; continue;
				/* } */ case 11:
				position = position + (1) >> 0;
				/* break; */ $s = 9; continue;
			/* } else if (_ref === 10) { */ $s = 9; continue; case 6: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 10))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 10))) {} else { $s = 12; continue; }
					/* goto l412 */ $s = 1; continue;
				/* } */ case 12:
				position = position + (1) >> 0;
				/* break; */ $s = 9; continue;
			/* } else if (_ref === 9) { */ $s = 9; continue; case 7: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 9))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 9))) {} else { $s = 13; continue; }
					/* goto l412 */ $s = 1; continue;
				/* } */ case 13:
				position = position + (1) >> 0;
				/* break; */ $s = 9; continue;
			/* } else { */ $s = 9; continue; case 8: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) {} else { $s = 14; continue; }
					/* goto l412 */ $s = 1; continue;
				/* } */ case 14:
				position = position + (1) >> 0;
				/* break; */ $s = 9; continue;
			/* } } */ case 9:
			/* l414: */ case 3:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position415 = _tmp$8; tokenIndex415 = _tmp$9; depth415 = _tmp$10;
			_ref$1 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* switch (0) { default: if (_ref$1 === 12) { */ if (_ref$1 === 12) {} else if (_ref$1 === 13) { $s = 15; continue; } else if (_ref$1 === 10) { $s = 16; continue; } else if (_ref$1 === 9) { $s = 17; continue; } else { $s = 18; continue; }
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 12))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 12))) {} else { $s = 20; continue; }
					/* goto l415 */ $s = 2; continue;
				/* } */ case 20:
				position = position + (1) >> 0;
				/* break; */ $s = 19; continue;
			/* } else if (_ref$1 === 13) { */ $s = 19; continue; case 15: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 13))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 13))) {} else { $s = 21; continue; }
					/* goto l415 */ $s = 2; continue;
				/* } */ case 21:
				position = position + (1) >> 0;
				/* break; */ $s = 19; continue;
			/* } else if (_ref$1 === 10) { */ $s = 19; continue; case 16: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 10))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 10))) {} else { $s = 22; continue; }
					/* goto l415 */ $s = 2; continue;
				/* } */ case 22:
				position = position + (1) >> 0;
				/* break; */ $s = 19; continue;
			/* } else if (_ref$1 === 9) { */ $s = 19; continue; case 17: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 9))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 9))) {} else { $s = 23; continue; }
					/* goto l415 */ $s = 2; continue;
				/* } */ case 23:
				position = position + (1) >> 0;
				/* break; */ $s = 19; continue;
			/* } else { */ $s = 19; continue; case 18: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) {} else { $s = 24; continue; }
					/* goto l415 */ $s = 2; continue;
				/* } */ case 24:
				position = position + (1) >> 0;
				/* break; */ $s = 19; continue;
			/* } } */ case 19:
			/* goto l414 */ $s = 3; continue;
			/* l415: */ case 2:
			_tmp$11 = position415; _tmp$12 = tokenIndex415; _tmp$13 = depth415; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			depth = depth - (1) >> 0;
			add(41, position413);
			return true;
			/* l412: */ case 1:
			_tmp$14 = position412; _tmp$15 = tokenIndex412; _tmp$16 = depth412; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position418, tokenIndex418, depth418, position419, _tmp$8, _tmp$9, _tmp$10, position420, tokenIndex420, depth420, _tmp$11, _tmp$12, _tmp$13, position422, _tmp$14, _tmp$15, _tmp$16, position425, tokenIndex425, depth425, position427, _tmp$17, _tmp$18, _tmp$19, position428, tokenIndex428, depth428, position430, _tmp$20, _tmp$21, _tmp$22, position431, tokenIndex431, depth431, c, _tmp$23, _tmp$24, _tmp$25, c$1, _tmp$26, _tmp$27, _tmp$28, _tmp$29, _tmp$30, _tmp$31, c$2, _tmp$32, _tmp$33, _tmp$34, position424, tokenIndex424, depth424, _tmp$35, _tmp$36, _tmp$37, position433, tokenIndex433, depth433, position435, _tmp$38, _tmp$39, _tmp$40, position436, tokenIndex436, depth436, position438, _tmp$41, _tmp$42, _tmp$43, position439, tokenIndex439, depth439, c$3, _tmp$44, _tmp$45, _tmp$46, c$4, _tmp$47, _tmp$48, _tmp$49, _tmp$50, _tmp$51, _tmp$52, c$5, _tmp$53, _tmp$54, _tmp$55, _tmp$56, _tmp$57, _tmp$58;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position418 = _tmp$5; tokenIndex418 = _tmp$6; depth418 = _tmp$7;
			position419 = position;
			depth = depth + (1) >> 0;
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position420 = _tmp$8; tokenIndex420 = _tmp$9; depth420 = _tmp$10;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 63))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 63))) {} else { $s = 18; continue; }
				/* goto l421 */ $s = 1; continue;
			/* } */ case 18:
			position = position + (1) >> 0;
			/* goto l420 */ $s = 2; continue;
			/* l421: */ case 1:
			_tmp$11 = position420; _tmp$12 = tokenIndex420; _tmp$13 = depth420; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 36))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 36))) {} else { $s = 19; continue; }
				/* goto l418 */ $s = 3; continue;
			/* } */ case 19:
			position = position + (1) >> 0;
			/* l420: */ case 2:
			position422 = position;
			depth = depth + (1) >> 0;
			_tmp$14 = position; _tmp$15 = tokenIndex; _tmp$16 = depth; position425 = _tmp$14; tokenIndex425 = _tmp$15; depth425 = _tmp$16;
			position427 = position;
			depth = depth + (1) >> 0;
			_tmp$17 = position; _tmp$18 = tokenIndex; _tmp$19 = depth; position428 = _tmp$17; tokenIndex428 = _tmp$18; depth428 = _tmp$19;
			position430 = position;
			depth = depth + (1) >> 0;
			_tmp$20 = position; _tmp$21 = tokenIndex; _tmp$22 = depth; position431 = _tmp$20; tokenIndex431 = _tmp$21; depth431 = _tmp$22;
			c = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c < 97 || c > 122) { */ if (c < 97 || c > 122) {} else { $s = 20; continue; }
				/* goto l432 */ $s = 4; continue;
			/* } */ case 20:
			position = position + (1) >> 0;
			/* goto l431 */ $s = 5; continue;
			/* l432: */ case 4:
			_tmp$23 = position431; _tmp$24 = tokenIndex431; _tmp$25 = depth431; position = _tmp$23; tokenIndex = _tmp$24; depth = _tmp$25;
			c$1 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c$1 < 65 || c$1 > 90) { */ if (c$1 < 65 || c$1 > 90) {} else { $s = 21; continue; }
				/* goto l429 */ $s = 6; continue;
			/* } */ case 21:
			position = position + (1) >> 0;
			/* l431: */ case 5:
			depth = depth - (1) >> 0;
			add(54, position430);
			/* goto l428 */ $s = 7; continue;
			/* l429: */ case 6:
			_tmp$26 = position428; _tmp$27 = tokenIndex428; _tmp$28 = depth428; position = _tmp$26; tokenIndex = _tmp$27; depth = _tmp$28;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 95))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 95))) {} else { $s = 22; continue; }
				/* goto l426 */ $s = 8; continue;
			/* } */ case 22:
			position = position + (1) >> 0;
			/* l428: */ case 7:
			depth = depth - (1) >> 0;
			add(53, position427);
			/* goto l425 */ $s = 9; continue;
			/* l426: */ case 8:
			_tmp$29 = position425; _tmp$30 = tokenIndex425; _tmp$31 = depth425; position = _tmp$29; tokenIndex = _tmp$30; depth = _tmp$31;
			c$2 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c$2 < 48 || c$2 > 57) { */ if (c$2 < 48 || c$2 > 57) {} else { $s = 23; continue; }
				/* goto l418 */ $s = 3; continue;
			/* } */ case 23:
			position = position + (1) >> 0;
			/* l425: */ case 9:
			/* l423: */ case 17:
			_tmp$32 = position; _tmp$33 = tokenIndex; _tmp$34 = depth; position424 = _tmp$32; tokenIndex424 = _tmp$33; depth424 = _tmp$34;
			_tmp$35 = position; _tmp$36 = tokenIndex; _tmp$37 = depth; position433 = _tmp$35; tokenIndex433 = _tmp$36; depth433 = _tmp$37;
			position435 = position;
			depth = depth + (1) >> 0;
			_tmp$38 = position; _tmp$39 = tokenIndex; _tmp$40 = depth; position436 = _tmp$38; tokenIndex436 = _tmp$39; depth436 = _tmp$40;
			position438 = position;
			depth = depth + (1) >> 0;
			_tmp$41 = position; _tmp$42 = tokenIndex; _tmp$43 = depth; position439 = _tmp$41; tokenIndex439 = _tmp$42; depth439 = _tmp$43;
			c$3 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c$3 < 97 || c$3 > 122) { */ if (c$3 < 97 || c$3 > 122) {} else { $s = 24; continue; }
				/* goto l440 */ $s = 10; continue;
			/* } */ case 24:
			position = position + (1) >> 0;
			/* goto l439 */ $s = 11; continue;
			/* l440: */ case 10:
			_tmp$44 = position439; _tmp$45 = tokenIndex439; _tmp$46 = depth439; position = _tmp$44; tokenIndex = _tmp$45; depth = _tmp$46;
			c$4 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c$4 < 65 || c$4 > 90) { */ if (c$4 < 65 || c$4 > 90) {} else { $s = 25; continue; }
				/* goto l437 */ $s = 12; continue;
			/* } */ case 25:
			position = position + (1) >> 0;
			/* l439: */ case 11:
			depth = depth - (1) >> 0;
			add(54, position438);
			/* goto l436 */ $s = 13; continue;
			/* l437: */ case 12:
			_tmp$47 = position436; _tmp$48 = tokenIndex436; _tmp$49 = depth436; position = _tmp$47; tokenIndex = _tmp$48; depth = _tmp$49;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 95))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 95))) {} else { $s = 26; continue; }
				/* goto l434 */ $s = 14; continue;
			/* } */ case 26:
			position = position + (1) >> 0;
			/* l436: */ case 13:
			depth = depth - (1) >> 0;
			add(53, position435);
			/* goto l433 */ $s = 15; continue;
			/* l434: */ case 14:
			_tmp$50 = position433; _tmp$51 = tokenIndex433; _tmp$52 = depth433; position = _tmp$50; tokenIndex = _tmp$51; depth = _tmp$52;
			c$5 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c$5 < 48 || c$5 > 57) { */ if (c$5 < 48 || c$5 > 57) {} else { $s = 27; continue; }
				/* goto l424 */ $s = 16; continue;
			/* } */ case 27:
			position = position + (1) >> 0;
			/* l433: */ case 15:
			/* goto l423 */ $s = 17; continue;
			/* l424: */ case 16:
			_tmp$53 = position424; _tmp$54 = tokenIndex424; _tmp$55 = depth424; position = _tmp$53; tokenIndex = _tmp$54; depth = _tmp$55;
			depth = depth - (1) >> 0;
			add(52, position422);
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 28; continue; }
				/* goto l418 */ $s = 3; continue;
			/* } */ case 28:
			depth = depth - (1) >> 0;
			add(42, position419);
			return true;
			/* l418: */ case 3:
			_tmp$56 = position418; _tmp$57 = tokenIndex418; _tmp$58 = depth418; position = _tmp$56; tokenIndex = _tmp$57; depth = _tmp$58;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position441, tokenIndex441, depth441, position442, _tmp$8, _tmp$9, _tmp$10, position444, tokenIndex444, depth444, _tmp$11, _tmp$12, _tmp$13, position445, tokenIndex445, depth445, _tmp$14, _tmp$15, _tmp$16, _tmp$17, _tmp$18, _tmp$19, _tmp$20, _tmp$21, _tmp$22;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position441 = _tmp$5; tokenIndex441 = _tmp$6; depth441 = _tmp$7;
			position442 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 60))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 60))) {} else { $s = 5; continue; }
				/* goto l441 */ $s = 1; continue;
			/* } */ case 5:
			position = position + (1) >> 0;
			/* l443: */ case 4:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position444 = _tmp$8; tokenIndex444 = _tmp$9; depth444 = _tmp$10;
			_tmp$11 = position; _tmp$12 = tokenIndex; _tmp$13 = depth; position445 = _tmp$11; tokenIndex445 = _tmp$12; depth445 = _tmp$13;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 62))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 62))) {} else { $s = 6; continue; }
				/* goto l445 */ $s = 2; continue;
			/* } */ case 6:
			position = position + (1) >> 0;
			/* goto l444 */ $s = 3; continue;
			/* l445: */ case 2:
			_tmp$14 = position445; _tmp$15 = tokenIndex445; _tmp$16 = depth445; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			/* if (!matchDot()) { */ if (!matchDot()) {} else { $s = 7; continue; }
				/* goto l444 */ $s = 3; continue;
			/* } */ case 7:
			/* goto l443 */ $s = 4; continue;
			/* l444: */ case 3:
			_tmp$17 = position444; _tmp$18 = tokenIndex444; _tmp$19 = depth444; position = _tmp$17; tokenIndex = _tmp$18; depth = _tmp$19;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 62))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 62))) {} else { $s = 8; continue; }
				/* goto l441 */ $s = 1; continue;
			/* } */ case 8:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 9; continue; }
				/* goto l441 */ $s = 1; continue;
			/* } */ case 9:
			depth = depth - (1) >> 0;
			add(43, position442);
			return true;
			/* l441: */ case 1:
			_tmp$20 = position441; _tmp$21 = tokenIndex441; _tmp$22 = depth441; position = _tmp$20; tokenIndex = _tmp$21; depth = _tmp$22;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position467, tokenIndex467, depth467, position468, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position467 = _tmp$5; tokenIndex467 = _tmp$6; depth467 = _tmp$7;
			position468 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 123))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 123))) {} else { $s = 2; continue; }
				/* goto l467 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l467 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(65, position468);
			return true;
			/* l467: */ case 1:
			_tmp$8 = position467; _tmp$9 = tokenIndex467; _tmp$10 = depth467; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position469, tokenIndex469, depth469, position470, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position469 = _tmp$5; tokenIndex469 = _tmp$6; depth469 = _tmp$7;
			position470 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 125))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 125))) {} else { $s = 2; continue; }
				/* goto l469 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l469 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(66, position470);
			return true;
			/* l469: */ case 1:
			_tmp$8 = position469; _tmp$9 = tokenIndex469; _tmp$10 = depth469; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position475, tokenIndex475, depth475, position476, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position475 = _tmp$5; tokenIndex475 = _tmp$6; depth475 = _tmp$7;
			position476 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 46))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 46))) {} else { $s = 2; continue; }
				/* goto l475 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l475 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(71, position476);
			return true;
			/* l475: */ case 1:
			_tmp$8 = position475; _tmp$9 = tokenIndex475; _tmp$10 = depth475; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position478, tokenIndex478, depth478, position479, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position478 = _tmp$5; tokenIndex478 = _tmp$6; depth478 = _tmp$7;
			position479 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 124))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 124))) {} else { $s = 2; continue; }
				/* goto l478 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l478 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(73, position479);
			return true;
			/* l478: */ case 1:
			_tmp$8 = position478; _tmp$9 = tokenIndex478; _tmp$10 = depth478; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position481, tokenIndex481, depth481, position482, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position481 = _tmp$5; tokenIndex481 = _tmp$6; depth481 = _tmp$7;
			position482 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 94))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 94))) {} else { $s = 2; continue; }
				/* goto l481 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l481 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(75, position482);
			return true;
			/* l481: */ case 1:
			_tmp$8 = position481; _tmp$9 = tokenIndex481; _tmp$10 = depth481; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position483, tokenIndex483, depth483, position484, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position483 = _tmp$5; tokenIndex483 = _tmp$6; depth483 = _tmp$7;
			position484 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 40))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 40))) {} else { $s = 2; continue; }
				/* goto l483 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l483 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(76, position484);
			return true;
			/* l483: */ case 1:
			_tmp$8 = position483; _tmp$9 = tokenIndex483; _tmp$10 = depth483; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position485, tokenIndex485, depth485, position486, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position485 = _tmp$5; tokenIndex485 = _tmp$6; depth485 = _tmp$7;
			position486 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 41))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 41))) {} else { $s = 2; continue; }
				/* goto l485 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l485 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(77, position486);
			return true;
			/* l485: */ case 1:
			_tmp$8 = position485; _tmp$9 = tokenIndex485; _tmp$10 = depth485; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position487, tokenIndex487, depth487, position488, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position487 = _tmp$5; tokenIndex487 = _tmp$6; depth487 = _tmp$7;
			position488 = position;
			depth = depth + (1) >> 0;
			/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 97))) {} else { $s = 2; continue; }
				/* goto l487 */ $s = 1; continue;
			/* } */ case 2:
			position = position + (1) >> 0;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 3; continue; }
				/* goto l487 */ $s = 1; continue;
			/* } */ case 3:
			depth = depth - (1) >> 0;
			add(78, position488);
			return true;
			/* l487: */ case 1:
			_tmp$8 = position487; _tmp$9 = tokenIndex487; _tmp$10 = depth487; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			return false;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, (function() {
			var $this = this, $args = arguments, $s = 0, _tmp$5, _tmp$6, _tmp$7, position497, tokenIndex497, depth497, position498, c, _tmp$8, _tmp$9, _tmp$10, position500, tokenIndex500, depth500, c$1, _tmp$11, _tmp$12, _tmp$13, _tmp$14, _tmp$15, _tmp$16;
			/* */ while (true) { switch ($s) { case 0:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position497 = _tmp$5; tokenIndex497 = _tmp$6; depth497 = _tmp$7;
			position498 = position;
			depth = depth + (1) >> 0;
			c = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c < 48 || c > 57) { */ if (c < 48 || c > 57) {} else { $s = 4; continue; }
				/* goto l497 */ $s = 1; continue;
			/* } */ case 4:
			position = position + (1) >> 0;
			/* l499: */ case 3:
			_tmp$8 = position; _tmp$9 = tokenIndex; _tmp$10 = depth; position500 = _tmp$8; tokenIndex500 = _tmp$9; depth500 = _tmp$10;
			c$1 = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* if (c$1 < 48 || c$1 > 57) { */ if (c$1 < 48 || c$1 > 57) {} else { $s = 5; continue; }
				/* goto l500 */ $s = 2; continue;
			/* } */ case 5:
			position = position + (1) >> 0;
			/* goto l499 */ $s = 3; continue;
			/* l500: */ case 2:
			_tmp$11 = position500; _tmp$12 = tokenIndex500; _tmp$13 = depth500; position = _tmp$11; tokenIndex = _tmp$12; depth = _tmp$13;
			/* if (!rules[88]()) { */ if (!rules[88]()) {} else { $s = 6; continue; }
				/* goto l497 */ $s = 1; continue;
			/* } */ case 6:
			depth = depth - (1) >> 0;
			add(87, position498);
			return true;
			/* l497: */ case 1:
			_tmp$14 = position497; _tmp$15 = tokenIndex497; _tmp$16 = depth497; position = _tmp$14; tokenIndex = _tmp$15; depth = _tmp$16;
			return false;
			/* */ case -1: } return; }
		}), (function() {
			var $this = this, $args = arguments, $s = 0, position502, _tmp$5, _tmp$6, _tmp$7, position504, tokenIndex504, depth504, _ref, _tmp$8, _tmp$9, _tmp$10;
			/* */ while (true) { switch ($s) { case 0:
			position502 = position;
			depth = depth + (1) >> 0;
			/* l503: */ case 2:
			_tmp$5 = position; _tmp$6 = tokenIndex; _tmp$7 = depth; position504 = _tmp$5; tokenIndex504 = _tmp$6; depth504 = _tmp$7;
			_ref = ((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]);
			/* switch (0) { default: if (_ref === 11) { */ if (_ref === 11) {} else if (_ref === 12) { $s = 3; continue; } else if (_ref === 10) { $s = 4; continue; } else if (_ref === 13) { $s = 5; continue; } else if (_ref === 9) { $s = 6; continue; } else { $s = 7; continue; }
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 11))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 11))) {} else { $s = 9; continue; }
					/* goto l504 */ $s = 1; continue;
				/* } */ case 9:
				position = position + (1) >> 0;
				/* break; */ $s = 8; continue;
			/* } else if (_ref === 12) { */ $s = 8; continue; case 3: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 12))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 12))) {} else { $s = 10; continue; }
					/* goto l504 */ $s = 1; continue;
				/* } */ case 10:
				position = position + (1) >> 0;
				/* break; */ $s = 8; continue;
			/* } else if (_ref === 10) { */ $s = 8; continue; case 4: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 10))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 10))) {} else { $s = 11; continue; }
					/* goto l504 */ $s = 1; continue;
				/* } */ case 11:
				position = position + (1) >> 0;
				/* break; */ $s = 8; continue;
			/* } else if (_ref === 13) { */ $s = 8; continue; case 5: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 13))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 13))) {} else { $s = 12; continue; }
					/* goto l504 */ $s = 1; continue;
				/* } */ case 12:
				position = position + (1) >> 0;
				/* break; */ $s = 8; continue;
			/* } else if (_ref === 9) { */ $s = 8; continue; case 6: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 9))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 9))) {} else { $s = 13; continue; }
					/* goto l504 */ $s = 1; continue;
				/* } */ case 13:
				position = position + (1) >> 0;
				/* break; */ $s = 8; continue;
			/* } else { */ $s = 8; continue; case 7: 
				/* if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) { */ if (!((((position < 0 || position >= buffer.$length) ? $throwRuntimeError("index out of range") : buffer.$array[buffer.$offset + position]) === 32))) {} else { $s = 14; continue; }
					/* goto l504 */ $s = 1; continue;
				/* } */ case 14:
				position = position + (1) >> 0;
				/* break; */ $s = 8; continue;
			/* } } */ case 8:
			/* goto l503 */ $s = 2; continue;
			/* l504: */ case 1:
			_tmp$8 = position504; _tmp$9 = tokenIndex504; _tmp$10 = depth504; position = _tmp$8; tokenIndex = _tmp$9; depth = _tmp$10;
			depth = depth - (1) >> 0;
			add(88, position502);
			return true;
			/* */ case -1: } return; }
		}), $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError, $throwNilPointerError]), ($arrayType(($funcType([], [$Bool], false)), 99)));
		$copy(p.rules, rules, ($arrayType(($funcType([], [$Bool], false)), 99)));
	};
	Sparql.prototype.Init = function() { return this.$val.Init(); };
	$pkg.$init = function() {
		($ptrType(triplePattern)).methods = [["addToScope", "addToScope", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [], false, -1], ["in$", "in", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [$Bool], false, -1]];
		triplePattern.init([["S", "S", "", $String, ""], ["P", "P", "", $String, ""], ["O", "O", "", $String, ""]]);
		($ptrType(Bgp)).methods = [["RecommendationQuery", "RecommendationQuery", "", [], [$String], false, -1], ["addToScope", "addToScope", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [], false, 0], ["addTriplePattern", "addTriplePattern", "github.com/scampi/gosparqled/autocompletion", [], [], false, -1], ["in$", "in", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [$Bool], false, 0], ["setObject", "setObject", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, -1], ["setPredicate", "setPredicate", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, -1], ["setSubject", "setSubject", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, -1], ["trimToScope", "trimToScope", "github.com/scampi/gosparqled/autocompletion", [], [], false, -1]];
		Bgp.init([["triplePattern", "", "github.com/scampi/gosparqled/autocompletion", triplePattern, ""], ["Tps", "Tps", "", ($sliceType(triplePattern)), ""], ["scope", "scope", "github.com/scampi/gosparqled/autocompletion", ($mapType($String, $Bool)), ""], ["Template", "Template", "", ($ptrType(mustache.Template)), ""]]);
		tokenTree.init([["AST", "AST", "", [], [($ptrType(node32))], false], ["Add", "Add", "", [pegRule, $Int, $Int, $Int, $Int], [], false], ["Error", "Error", "", [], [($sliceType(token32))], false], ["Expand", "Expand", "", [$Int], [tokenTree], false], ["Print", "Print", "", [], [], false], ["PrintSyntax", "PrintSyntax", "", [], [], false], ["PrintSyntaxTree", "PrintSyntaxTree", "", [$String], [], false], ["Tokens", "Tokens", "", [], [($chanType(token32, false, true))], false], ["trim", "trim", "github.com/scampi/gosparqled/autocompletion", [$Int], [], false]]);
		($ptrType(node32)).methods = [["Print", "Print", "", [$String], [], false, -1], ["String", "String", "", [], [$String], false, 0], ["getToken32", "getToken32", "github.com/scampi/gosparqled/autocompletion", [], [token32], false, 0], ["isParentOf", "isParentOf", "github.com/scampi/gosparqled/autocompletion", [token32], [$Bool], false, 0], ["isZero", "isZero", "github.com/scampi/gosparqled/autocompletion", [], [$Bool], false, 0], ["print", "print", "github.com/scampi/gosparqled/autocompletion", [$Int, $String], [], false, -1]];
		node32.init([["token32", "", "github.com/scampi/gosparqled/autocompletion", token32, ""], ["up", "up", "github.com/scampi/gosparqled/autocompletion", ($ptrType(node32)), ""], ["next", "next", "github.com/scampi/gosparqled/autocompletion", ($ptrType(node32)), ""]]);
		element.init([["node", "node", "github.com/scampi/gosparqled/autocompletion", ($ptrType(node32)), ""], ["down", "down", "github.com/scampi/gosparqled/autocompletion", ($ptrType(element)), ""]]);
		($ptrType(token16)).methods = [["String", "String", "", [], [$String], false, -1], ["getToken32", "getToken32", "github.com/scampi/gosparqled/autocompletion", [], [token32], false, -1], ["isParentOf", "isParentOf", "github.com/scampi/gosparqled/autocompletion", [token16], [$Bool], false, -1], ["isZero", "isZero", "github.com/scampi/gosparqled/autocompletion", [], [$Bool], false, -1]];
		token16.init([["pegRule", "", "github.com/scampi/gosparqled/autocompletion", pegRule, ""], ["begin", "begin", "github.com/scampi/gosparqled/autocompletion", $Int16, ""], ["end", "end", "github.com/scampi/gosparqled/autocompletion", $Int16, ""], ["next", "next", "github.com/scampi/gosparqled/autocompletion", $Int16, ""]]);
		($ptrType(tokens16)).methods = [["AST", "AST", "", [], [($ptrType(node32))], false, -1], ["Add", "Add", "", [pegRule, $Int, $Int, $Int, $Int], [], false, -1], ["Error", "Error", "", [], [($sliceType(token32))], false, -1], ["Expand", "Expand", "", [$Int], [tokenTree], false, -1], ["Order", "Order", "", [], [($sliceType(($sliceType(token16))))], false, -1], ["PreOrder", "PreOrder", "", [], [($chanType(state16, false, true)), ($sliceType(($sliceType(token16))))], false, -1], ["Print", "Print", "", [], [], false, -1], ["PrintSyntax", "PrintSyntax", "", [], [], false, -1], ["PrintSyntaxTree", "PrintSyntaxTree", "", [$String], [], false, -1], ["Tokens", "Tokens", "", [], [($chanType(token32, false, true))], false, -1], ["trim", "trim", "github.com/scampi/gosparqled/autocompletion", [$Int], [], false, -1]];
		tokens16.init([["tree", "tree", "github.com/scampi/gosparqled/autocompletion", ($sliceType(token16)), ""], ["ordered", "ordered", "github.com/scampi/gosparqled/autocompletion", ($sliceType(($sliceType(token16)))), ""]]);
		($ptrType(state16)).methods = [["String", "String", "", [], [$String], false, 0], ["getToken32", "getToken32", "github.com/scampi/gosparqled/autocompletion", [], [token32], false, 0], ["isParentOf", "isParentOf", "github.com/scampi/gosparqled/autocompletion", [token16], [$Bool], false, 0], ["isZero", "isZero", "github.com/scampi/gosparqled/autocompletion", [], [$Bool], false, 0]];
		state16.init([["token16", "", "github.com/scampi/gosparqled/autocompletion", token16, ""], ["depths", "depths", "github.com/scampi/gosparqled/autocompletion", ($sliceType($Int16)), ""], ["leaf", "leaf", "github.com/scampi/gosparqled/autocompletion", $Bool, ""]]);
		($ptrType(token32)).methods = [["String", "String", "", [], [$String], false, -1], ["getToken32", "getToken32", "github.com/scampi/gosparqled/autocompletion", [], [token32], false, -1], ["isParentOf", "isParentOf", "github.com/scampi/gosparqled/autocompletion", [token32], [$Bool], false, -1], ["isZero", "isZero", "github.com/scampi/gosparqled/autocompletion", [], [$Bool], false, -1]];
		token32.init([["pegRule", "", "github.com/scampi/gosparqled/autocompletion", pegRule, ""], ["begin", "begin", "github.com/scampi/gosparqled/autocompletion", $Int32, ""], ["end", "end", "github.com/scampi/gosparqled/autocompletion", $Int32, ""], ["next", "next", "github.com/scampi/gosparqled/autocompletion", $Int32, ""]]);
		($ptrType(tokens32)).methods = [["AST", "AST", "", [], [($ptrType(node32))], false, -1], ["Add", "Add", "", [pegRule, $Int, $Int, $Int, $Int], [], false, -1], ["Error", "Error", "", [], [($sliceType(token32))], false, -1], ["Expand", "Expand", "", [$Int], [tokenTree], false, -1], ["Order", "Order", "", [], [($sliceType(($sliceType(token32))))], false, -1], ["PreOrder", "PreOrder", "", [], [($chanType(state32, false, true)), ($sliceType(($sliceType(token32))))], false, -1], ["Print", "Print", "", [], [], false, -1], ["PrintSyntax", "PrintSyntax", "", [], [], false, -1], ["PrintSyntaxTree", "PrintSyntaxTree", "", [$String], [], false, -1], ["Tokens", "Tokens", "", [], [($chanType(token32, false, true))], false, -1], ["trim", "trim", "github.com/scampi/gosparqled/autocompletion", [$Int], [], false, -1]];
		tokens32.init([["tree", "tree", "github.com/scampi/gosparqled/autocompletion", ($sliceType(token32)), ""], ["ordered", "ordered", "github.com/scampi/gosparqled/autocompletion", ($sliceType(($sliceType(token32)))), ""]]);
		($ptrType(state32)).methods = [["String", "String", "", [], [$String], false, 0], ["getToken32", "getToken32", "github.com/scampi/gosparqled/autocompletion", [], [token32], false, 0], ["isParentOf", "isParentOf", "github.com/scampi/gosparqled/autocompletion", [token32], [$Bool], false, 0], ["isZero", "isZero", "github.com/scampi/gosparqled/autocompletion", [], [$Bool], false, 0]];
		state32.init([["token32", "", "github.com/scampi/gosparqled/autocompletion", token32, ""], ["depths", "depths", "github.com/scampi/gosparqled/autocompletion", ($sliceType($Int32)), ""], ["leaf", "leaf", "github.com/scampi/gosparqled/autocompletion", $Bool, ""]]);
		Sparql.methods = [["AST", "AST", "", [], [($ptrType(node32))], false, 6], ["Add", "Add", "", [pegRule, $Int, $Int, $Int, $Int], [], false, 6], ["Error", "Error", "", [], [($sliceType(token32))], false, 6], ["Expand", "Expand", "", [$Int], [tokenTree], false, 6], ["Print", "Print", "", [], [], false, 6], ["PrintSyntax", "PrintSyntax", "", [], [], false, 6], ["RecommendationQuery", "RecommendationQuery", "", [], [$String], false, 0], ["Tokens", "Tokens", "", [], [($chanType(token32, false, true))], false, 6], ["addToScope", "addToScope", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [], false, 0], ["addTriplePattern", "addTriplePattern", "github.com/scampi/gosparqled/autocompletion", [], [], false, 0], ["in$", "in", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [$Bool], false, 0], ["setObject", "setObject", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, 0], ["setPredicate", "setPredicate", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, 0], ["setSubject", "setSubject", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, 0], ["trim", "trim", "github.com/scampi/gosparqled/autocompletion", [$Int], [], false, 6], ["trimToScope", "trimToScope", "github.com/scampi/gosparqled/autocompletion", [], [], false, 0]];
		($ptrType(Sparql)).methods = [["AST", "AST", "", [], [($ptrType(node32))], false, 6], ["Add", "Add", "", [pegRule, $Int, $Int, $Int, $Int], [], false, 6], ["Error", "Error", "", [], [($sliceType(token32))], false, 6], ["Execute", "Execute", "", [], [], false, -1], ["Expand", "Expand", "", [$Int], [tokenTree], false, 6], ["Highlighter", "Highlighter", "", [], [], false, -1], ["Init", "Init", "", [], [], false, -1], ["Print", "Print", "", [], [], false, 6], ["PrintSyntax", "PrintSyntax", "", [], [], false, 6], ["PrintSyntaxTree", "PrintSyntaxTree", "", [], [], false, -1], ["RecommendationQuery", "RecommendationQuery", "", [], [$String], false, 0], ["Tokens", "Tokens", "", [], [($chanType(token32, false, true))], false, 6], ["addToScope", "addToScope", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [], false, 0], ["addTriplePattern", "addTriplePattern", "github.com/scampi/gosparqled/autocompletion", [], [], false, 0], ["in$", "in", "github.com/scampi/gosparqled/autocompletion", [($mapType($String, $Bool))], [$Bool], false, 0], ["setObject", "setObject", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, 0], ["setPredicate", "setPredicate", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, 0], ["setSubject", "setSubject", "github.com/scampi/gosparqled/autocompletion", [$String], [], false, 0], ["trim", "trim", "github.com/scampi/gosparqled/autocompletion", [$Int], [], false, 6], ["trimToScope", "trimToScope", "github.com/scampi/gosparqled/autocompletion", [], [], false, 0]];
		Sparql.init([["Bgp", "", "", ($ptrType(Bgp)), ""], ["Buffer", "Buffer", "", $String, ""], ["buffer", "buffer", "github.com/scampi/gosparqled/autocompletion", ($sliceType($Int32)), ""], ["rules", "rules", "github.com/scampi/gosparqled/autocompletion", ($arrayType(($funcType([], [$Bool], false)), 99)), ""], ["Parse", "Parse", "", ($funcType([($sliceType($Int))], [$error], true)), ""], ["Reset", "Reset", "", ($funcType([], [], false)), ""], ["tokenTree", "", "github.com/scampi/gosparqled/autocompletion", tokenTree, ""]]);
		textPosition.init([["line", "line", "github.com/scampi/gosparqled/autocompletion", $Int, ""], ["symbol", "symbol", "github.com/scampi/gosparqled/autocompletion", $Int, ""]]);
		($ptrType(parseError)).methods = [["Error", "Error", "", [], [$String], false, -1]];
		parseError.init([["p", "p", "github.com/scampi/gosparqled/autocompletion", ($ptrType(Sparql)), ""]]);
		rul3s = $toNativeArray("String", ["Unknown", "queryContainer", "prolog", "prefixDecl", "baseDecl", "query", "selectQuery", "select", "subSelect", "projectionElem", "datasetClause", "whereClause", "groupGraphPattern", "graphPattern", "graphPatternNotTriples", "optionalGraphPattern", "groupOrUnionGraphPattern", "basicGraphPattern", "triplesBlock", "triplesSameSubjectPath", "varOrTerm", "graphTerm", "triplesNodePath", "collectionPath", "blankNodePropertyListPath", "propertyListPath", "verbPath", "path", "pathAlternative", "pathSequence", "pathElt", "pathPrimary", "pathNegatedPropertySet", "pathOneInPropertySet", "objectListPath", "objectPath", "graphNodePath", "solutionModifier", "limitOffsetClauses", "limit", "offset", "pof", "var", "iri", "literal", "string", "numericLiteral", "booleanLiteral", "blankNode", "blankNodeLabel", "anon", "nil", "VARNAME", "PN_CHARS_U", "PN_CHARS_BASE", "PREFIX", "TRUE", "FALSE", "BASE", "SELECT", "REDUCED", "DISTINCT", "FROM", "NAMED", "WHERE", "LBRACE", "RBRACE", "LBRACK", "RBRACK", "SEMICOLON", "COMMA", "DOT", "COLON", "PIPE", "SLASH", "INVERSE", "LPAREN", "RPAREN", "ISA", "NOT", "STAR", "QUESTION", "PLUS", "OPTIONAL", "UNION", "LIMIT", "OFFSET", "INTEGER", "ws", "PegText", "Action0", "Action1", "Action2", "Action3", "Action4", "Action5", "Action6", "Action7", "Action8", "Pre_", "_In_", "_Suf"]);
	};
	return $pkg;
})();
$packages["/home/stecam/documents/prog/go/src/github.com/scampi/gosparqled"] = (function() {
	var $pkg = {}, js = $packages["github.com/gopherjs/gopherjs/js"], mustache = $packages["github.com/hoisie/mustache"], autocompletion = $packages["github.com/scampi/gosparqled/autocompletion"], tmpl, tp, _tuple, RecommendationQuery, main;
	RecommendationQuery = $pkg.RecommendationQuery = function(query, callback) {
		$go((function(query$1, $b) {
			var $this = this, $args = arguments, $r, $s = 0, s, err;
			/* */ if(!$b) { $nonblockingCall(); }; return function() { while (true) { switch ($s) { case 0:
			s = new autocompletion.Sparql.Ptr(new autocompletion.Bgp.Ptr(new autocompletion.triplePattern.Ptr(), ($sliceType(autocompletion.triplePattern)).nil, false, tp), query$1, ($sliceType($Int32)).nil, ($arrayType(($funcType([], [$Bool], false)), 99)).zero(), $throwNilPointerError, $throwNilPointerError, null);
			s.Init();
			err = s.Parse(new ($sliceType($Int))([]));
			/* if ($interfaceIsEqual(err, null)) { */ if ($interfaceIsEqual(err, null)) {} else { $s = 1; continue; }
				$r = s.Execute(true); /* */ $s = 3; case 3: if ($r && $r.constructor === Function) { $r = $r(); }
				callback(s.Bgp.RecommendationQuery());
			/* } else { */ $s = 2; continue; case 1: 
				callback(query$1 + "\n" + err.Error());
			/* } */ case 2:
			/* */ case -1: } return; } };
		}), [query]);
	};
	main = function() {
		var _map, _key;
		$global.autocompletion = $externalize((_map = new $Map(), _key = "RecommendationQuery", _map[_key] = { k: _key, v: new ($funcType([$String, ($funcType([$String], [], false))], [], false))(RecommendationQuery) }, _map), ($mapType($String, $emptyInterface)));
	};
	$pkg.$run = function($b) {
		$packages["github.com/gopherjs/gopherjs/js"].$init();
		$packages["runtime"].$init();
		$packages["errors"].$init();
		$packages["sync/atomic"].$init();
		$packages["sync"].$init();
		$packages["io"].$init();
		$packages["unicode"].$init();
		$packages["unicode/utf8"].$init();
		$packages["bytes"].$init();
		$packages["math"].$init();
		$packages["syscall"].$init();
		$packages["strings"].$init();
		$packages["time"].$init();
		$packages["os"].$init();
		$packages["strconv"].$init();
		$packages["reflect"].$init();
		$packages["fmt"].$init();
		$packages["sort"].$init();
		$packages["path/filepath"].$init();
		$packages["io/ioutil"].$init();
		$packages["path"].$init();
		$packages["github.com/hoisie/mustache"].$init();
		$packages["github.com/scampi/gosparqled/autocompletion"].$init();
		$pkg.$init();
		main();
	};
	$pkg.$init = function() {
		tmpl = "\n    SELECT DISTINCT ?POF\n    WHERE {\n    {{#Tps}}\n        {{{S}}} {{{P}}} {{{O}}} .\n    {{/Tps}}\n    }\n    LIMIT 10\n";
		_tuple = mustache.ParseString(tmpl); tp = _tuple[0];
	};
	return $pkg;
})();
$error.implementedBy = [$packages["errors"].errorString.Ptr, $packages["github.com/gopherjs/gopherjs/js"].Error.Ptr, $packages["github.com/hoisie/mustache"].parseError, $packages["github.com/hoisie/mustache"].parseError.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].parseError.Ptr, $packages["os"].LinkError.Ptr, $packages["os"].PathError.Ptr, $packages["os"].SyscallError.Ptr, $packages["reflect"].ValueError.Ptr, $packages["runtime"].NotSupportedError.Ptr, $packages["runtime"].TypeAssertionError.Ptr, $packages["runtime"].errorString, $packages["syscall"].Errno, $packages["time"].ParseError.Ptr, $ptrType($packages["runtime"].errorString), $ptrType($packages["syscall"].Errno)];
$packages["github.com/gopherjs/gopherjs/js"].Object.implementedBy = [$packages["github.com/gopherjs/gopherjs/js"].Error, $packages["github.com/gopherjs/gopherjs/js"].Error.Ptr];
$packages["sync"].Locker.implementedBy = [$packages["sync"].Mutex.Ptr, $packages["sync"].RWMutex.Ptr, $packages["sync"].poolLocal.Ptr, $packages["sync"].rlocker.Ptr, $packages["syscall"].mmapper.Ptr];
$packages["io"].Reader.implementedBy = [$packages["bytes"].Buffer.Ptr, $packages["fmt"].ss.Ptr, $packages["os"].File.Ptr];
$packages["io"].RuneReader.implementedBy = [$packages["bytes"].Buffer.Ptr, $packages["fmt"].ss.Ptr];
$packages["io"].Writer.implementedBy = [$packages["bytes"].Buffer.Ptr, $packages["fmt"].pp.Ptr, $packages["os"].File.Ptr, $ptrType($packages["fmt"].buffer)];
$packages["os"].FileInfo.implementedBy = [$packages["os"].fileStat.Ptr];
$packages["reflect"].Type.implementedBy = [$packages["reflect"].arrayType.Ptr, $packages["reflect"].chanType.Ptr, $packages["reflect"].funcType.Ptr, $packages["reflect"].interfaceType.Ptr, $packages["reflect"].mapType.Ptr, $packages["reflect"].ptrType.Ptr, $packages["reflect"].rtype.Ptr, $packages["reflect"].sliceType.Ptr, $packages["reflect"].structType.Ptr];
$packages["fmt"].Formatter.implementedBy = [];
$packages["fmt"].GoStringer.implementedBy = [];
$packages["fmt"].State.implementedBy = [$packages["fmt"].pp.Ptr];
$packages["fmt"].Stringer.implementedBy = [$packages["bytes"].Buffer.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].node32.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].state16.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].state32.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].token16.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].token32.Ptr, $packages["os"].FileMode, $packages["reflect"].ChanDir, $packages["reflect"].Kind, $packages["reflect"].Value, $packages["reflect"].Value.Ptr, $packages["reflect"].arrayType.Ptr, $packages["reflect"].chanType.Ptr, $packages["reflect"].funcType.Ptr, $packages["reflect"].interfaceType.Ptr, $packages["reflect"].mapType.Ptr, $packages["reflect"].ptrType.Ptr, $packages["reflect"].rtype.Ptr, $packages["reflect"].sliceType.Ptr, $packages["reflect"].structType.Ptr, $packages["strconv"].decimal.Ptr, $packages["time"].Duration, $packages["time"].Location.Ptr, $packages["time"].Month, $packages["time"].Time, $packages["time"].Time.Ptr, $packages["time"].Weekday, $ptrType($packages["os"].FileMode), $ptrType($packages["reflect"].ChanDir), $ptrType($packages["reflect"].Kind), $ptrType($packages["time"].Duration), $ptrType($packages["time"].Month), $ptrType($packages["time"].Weekday)];
$packages["fmt"].runeUnreader.implementedBy = [$packages["bytes"].Buffer.Ptr, $packages["fmt"].ss.Ptr];
$packages["github.com/scampi/gosparqled/autocompletion"].tokenTree.implementedBy = [$packages["github.com/scampi/gosparqled/autocompletion"].tokens16.Ptr, $packages["github.com/scampi/gosparqled/autocompletion"].tokens32.Ptr];
$go($packages["/home/stecam/documents/prog/go/src/github.com/scampi/gosparqled"].$run, [], true);

})();
//# sourceMappingURL=gosparqled.js.map
