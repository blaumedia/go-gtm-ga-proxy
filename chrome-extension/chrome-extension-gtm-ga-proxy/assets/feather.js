!(function (e, n) {
    'object' == typeof exports && 'object' == typeof module
      ? (module.exports = n())
      : 'function' == typeof define && define.amd
      ? define([], n)
      : 'object' == typeof exports
      ? (exports.feather = n())
      : (e.feather = n());
  })('undefined' != typeof self ? self : this, function () {
    return (function (e) {
      var n = {};
      function i(t) {
        if (n[t]) return n[t].exports;
        var l = (n[t] = { i: t, l: !1, exports: {} });
        return e[t].call(l.exports, l, l.exports, i), (l.l = !0), l.exports;
      }
      return (
        (i.m = e),
        (i.c = n),
        (i.d = function (e, n, t) {
          i.o(e, n) ||
            Object.defineProperty(e, n, {
              configurable: !1,
              enumerable: !0,
              get: t,
            });
        }),
        (i.r = function (e) {
          Object.defineProperty(e, '__esModule', { value: !0 });
        }),
        (i.n = function (e) {
          var n =
            e && e.__esModule
              ? function () {
                  return e.default;
                }
              : function () {
                  return e;
                };
          return i.d(n, 'a', n), n;
        }),
        (i.o = function (e, n) {
          return Object.prototype.hasOwnProperty.call(e, n);
        }),
        (i.p = ''),
        i((i.s = 80))
      );
    })([
      function (e, n, i) {
        (function (n) {
          var i = 'object',
            t = function (e) {
              return e && e.Math == Math && e;
            };
          e.exports =
            t(typeof globalThis == i && globalThis) ||
            t(typeof window == i && window) ||
            t(typeof self == i && self) ||
            t(typeof n == i && n) ||
            Function('return this')();
        }.call(this, i(75)));
      },
      function (e, n) {
        var i = {}.hasOwnProperty;
        e.exports = function (e, n) {
          return i.call(e, n);
        };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(11),
          r = i(33),
          o = i(62),
          a = t.Symbol,
          c = l('wks');
        e.exports = function (e) {
          return c[e] || (c[e] = (o && a[e]) || (o ? a : r)('Symbol.' + e));
        };
      },
      function (e, n, i) {
        var t = i(6);
        e.exports = function (e) {
          if (!t(e)) throw TypeError(String(e) + ' is not an object');
          return e;
        };
      },
      function (e, n) {
        e.exports = function (e) {
          try {
            return !!e();
          } catch (e) {
            return !0;
          }
        };
      },
      function (e, n, i) {
        var t = i(8),
          l = i(7),
          r = i(10);
        e.exports = t
          ? function (e, n, i) {
              return l.f(e, n, r(1, i));
            }
          : function (e, n, i) {
              return (e[n] = i), e;
            };
      },
      function (e, n) {
        e.exports = function (e) {
          return 'object' == typeof e ? null !== e : 'function' == typeof e;
        };
      },
      function (e, n, i) {
        var t = i(8),
          l = i(35),
          r = i(3),
          o = i(18),
          a = Object.defineProperty;
        n.f = t
          ? a
          : function (e, n, i) {
              if ((r(e), (n = o(n, !0)), r(i), l))
                try {
                  return a(e, n, i);
                } catch (e) {}
              if ('get' in i || 'set' in i)
                throw TypeError('Accessors not supported');
              return 'value' in i && (e[n] = i.value), e;
            };
      },
      function (e, n, i) {
        var t = i(4);
        e.exports = !t(function () {
          return (
            7 !=
            Object.defineProperty({}, 'a', {
              get: function () {
                return 7;
              },
            }).a
          );
        });
      },
      function (e, n) {
        e.exports = {};
      },
      function (e, n) {
        e.exports = function (e, n) {
          return {
            enumerable: !(1 & e),
            configurable: !(2 & e),
            writable: !(4 & e),
            value: n,
          };
        };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(19),
          r = i(17),
          o = t['__core-js_shared__'] || l('__core-js_shared__', {});
        (e.exports = function (e, n) {
          return o[e] || (o[e] = void 0 !== n ? n : {});
        })('versions', []).push({
          version: '3.1.3',
          mode: r ? 'pure' : 'global',
          copyright: '© 2019 Denis Pushkarev (zloirock.ru)',
        });
      },
      function (e, n, i) {
        'use strict';
        Object.defineProperty(n, '__esModule', { value: !0 });
        var t = o(i(43)),
          l = o(i(41)),
          r = o(i(40));
        function o(e) {
          return e && e.__esModule ? e : { default: e };
        }
        n.default = Object.keys(l.default)
          .map(function (e) {
            return new t.default(e, l.default[e], r.default[e]);
          })
          .reduce(function (e, n) {
            return (e[n.name] = n), e;
          }, {});
      },
      function (e, n) {
        e.exports = [
          'constructor',
          'hasOwnProperty',
          'isPrototypeOf',
          'propertyIsEnumerable',
          'toLocaleString',
          'toString',
          'valueOf',
        ];
      },
      function (e, n, i) {
        var t = i(72),
          l = i(20);
        e.exports = function (e) {
          return t(l(e));
        };
      },
      function (e, n) {
        e.exports = {};
      },
      function (e, n, i) {
        var t = i(11),
          l = i(33),
          r = t('keys');
        e.exports = function (e) {
          return r[e] || (r[e] = l(e));
        };
      },
      function (e, n) {
        e.exports = !1;
      },
      function (e, n, i) {
        var t = i(6);
        e.exports = function (e, n) {
          if (!t(e)) return e;
          var i, l;
          if (n && 'function' == typeof (i = e.toString) && !t((l = i.call(e))))
            return l;
          if ('function' == typeof (i = e.valueOf) && !t((l = i.call(e))))
            return l;
          if (!n && 'function' == typeof (i = e.toString) && !t((l = i.call(e))))
            return l;
          throw TypeError("Can't convert object to primitive value");
        };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(5);
        e.exports = function (e, n) {
          try {
            l(t, e, n);
          } catch (i) {
            t[e] = n;
          }
          return n;
        };
      },
      function (e, n) {
        e.exports = function (e) {
          if (void 0 == e) throw TypeError("Can't call method on " + e);
          return e;
        };
      },
      function (e, n) {
        var i = Math.ceil,
          t = Math.floor;
        e.exports = function (e) {
          return isNaN((e = +e)) ? 0 : (e > 0 ? t : i)(e);
        };
      },
      function (e, n, i) {
        var t;
        /*!
    Copyright (c) 2016 Jed Watson.
    Licensed under the MIT License (MIT), see
    http://jedwatson.github.io/classnames
  */
        /*!
    Copyright (c) 2016 Jed Watson.
    Licensed under the MIT License (MIT), see
    http://jedwatson.github.io/classnames
  */
        !(function () {
          'use strict';
          var i = (function () {
            function e() {}
            function n(e, n) {
              for (var i = n.length, t = 0; t < i; ++t) l(e, n[t]);
            }
            e.prototype = Object.create(null);
            var i = {}.hasOwnProperty;
            var t = /\s+/;
            function l(e, l) {
              if (l) {
                var r = typeof l;
                'string' === r
                  ? (function (e, n) {
                      for (var i = n.split(t), l = i.length, r = 0; r < l; ++r)
                        e[i[r]] = !0;
                    })(e, l)
                  : Array.isArray(l)
                  ? n(e, l)
                  : 'object' === r
                  ? (function (e, n) {
                      for (var t in n) i.call(n, t) && (e[t] = !!n[t]);
                    })(e, l)
                  : 'number' === r &&
                    (function (e, n) {
                      e[n] = !0;
                    })(e, l);
              }
            }
            return function () {
              for (var i = arguments.length, t = Array(i), l = 0; l < i; l++)
                t[l] = arguments[l];
              var r = new e();
              n(r, t);
              var o = [];
              for (var a in r) r[a] && o.push(a);
              return o.join(' ');
            };
          })();
          void 0 !== e && e.exports
            ? (e.exports = i)
            : void 0 ===
                (t = function () {
                  return i;
                }.apply(n, [])) || (e.exports = t);
        })();
      },
      function (e, n, i) {
        var t = i(7).f,
          l = i(1),
          r = i(2)('toStringTag');
        e.exports = function (e, n, i) {
          e &&
            !l((e = i ? e : e.prototype), r) &&
            t(e, r, { configurable: !0, value: n });
        };
      },
      function (e, n, i) {
        var t = i(20);
        e.exports = function (e) {
          return Object(t(e));
        };
      },
      function (e, n, i) {
        var t = i(1),
          l = i(24),
          r = i(16),
          o = i(63),
          a = r('IE_PROTO'),
          c = Object.prototype;
        e.exports = o
          ? Object.getPrototypeOf
          : function (e) {
              return (
                (e = l(e)),
                t(e, a)
                  ? e[a]
                  : 'function' == typeof e.constructor &&
                    e instanceof e.constructor
                  ? e.constructor.prototype
                  : e instanceof Object
                  ? c
                  : null
              );
            };
      },
      function (e, n, i) {
        'use strict';
        var t,
          l,
          r,
          o = i(25),
          a = i(5),
          c = i(1),
          p = i(2),
          y = i(17),
          h = p('iterator'),
          x = !1;
        [].keys &&
          ('next' in (r = [].keys())
            ? (l = o(o(r))) !== Object.prototype && (t = l)
            : (x = !0)),
          void 0 == t && (t = {}),
          y ||
            c(t, h) ||
            a(t, h, function () {
              return this;
            }),
          (e.exports = { IteratorPrototype: t, BUGGY_SAFARI_ITERATORS: x });
      },
      function (e, n, i) {
        var t = i(21),
          l = Math.min;
        e.exports = function (e) {
          return e > 0 ? l(t(e), 9007199254740991) : 0;
        };
      },
      function (e, n, i) {
        var t = i(1),
          l = i(14),
          r = i(68),
          o = i(15),
          a = r(!1);
        e.exports = function (e, n) {
          var i,
            r = l(e),
            c = 0,
            p = [];
          for (i in r) !t(o, i) && t(r, i) && p.push(i);
          for (; n.length > c; ) t(r, (i = n[c++])) && (~a(p, i) || p.push(i));
          return p;
        };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(11),
          r = i(5),
          o = i(1),
          a = i(19),
          c = i(36),
          p = i(37),
          y = p.get,
          h = p.enforce,
          x = String(c).split('toString');
        l('inspectSource', function (e) {
          return c.call(e);
        }),
          (e.exports = function (e, n, i, l) {
            var c = !!l && !!l.unsafe,
              p = !!l && !!l.enumerable,
              y = !!l && !!l.noTargetGet;
            'function' == typeof i &&
              ('string' != typeof n || o(i, 'name') || r(i, 'name', n),
              (h(i).source = x.join('string' == typeof n ? n : ''))),
              e !== t
                ? (c ? !y && e[n] && (p = !0) : delete e[n],
                  p ? (e[n] = i) : r(e, n, i))
                : p
                ? (e[n] = i)
                : a(n, i);
          })(Function.prototype, 'toString', function () {
            return ('function' == typeof this && y(this).source) || c.call(this);
          });
      },
      function (e, n) {
        var i = {}.toString;
        e.exports = function (e) {
          return i.call(e).slice(8, -1);
        };
      },
      function (e, n, i) {
        var t = i(8),
          l = i(73),
          r = i(10),
          o = i(14),
          a = i(18),
          c = i(1),
          p = i(35),
          y = Object.getOwnPropertyDescriptor;
        n.f = t
          ? y
          : function (e, n) {
              if (((e = o(e)), (n = a(n, !0)), p))
                try {
                  return y(e, n);
                } catch (e) {}
              if (c(e, n)) return r(!l.f.call(e, n), e[n]);
            };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(31).f,
          r = i(5),
          o = i(29),
          a = i(19),
          c = i(71),
          p = i(65);
        e.exports = function (e, n) {
          var i,
            y,
            h,
            x,
            s,
            u = e.target,
            d = e.global,
            f = e.stat;
          if ((i = d ? t : f ? t[u] || a(u, {}) : (t[u] || {}).prototype))
            for (y in n) {
              if (
                ((x = n[y]),
                (h = e.noTargetGet ? (s = l(i, y)) && s.value : i[y]),
                !p(d ? y : u + (f ? '.' : '#') + y, e.forced) && void 0 !== h)
              ) {
                if (typeof x == typeof h) continue;
                c(x, h);
              }
              (e.sham || (h && h.sham)) && r(x, 'sham', !0), o(i, y, x, e);
            }
        };
      },
      function (e, n) {
        var i = 0,
          t = Math.random();
        e.exports = function (e) {
          return 'Symbol('.concat(
            void 0 === e ? '' : e,
            ')_',
            (++i + t).toString(36)
          );
        };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(6),
          r = t.document,
          o = l(r) && l(r.createElement);
        e.exports = function (e) {
          return o ? r.createElement(e) : {};
        };
      },
      function (e, n, i) {
        var t = i(8),
          l = i(4),
          r = i(34);
        e.exports =
          !t &&
          !l(function () {
            return (
              7 !=
              Object.defineProperty(r('div'), 'a', {
                get: function () {
                  return 7;
                },
              }).a
            );
          });
      },
      function (e, n, i) {
        var t = i(11);
        e.exports = t('native-function-to-string', Function.toString);
      },
      function (e, n, i) {
        var t,
          l,
          r,
          o = i(76),
          a = i(0),
          c = i(6),
          p = i(5),
          y = i(1),
          h = i(16),
          x = i(15),
          s = a.WeakMap;
        if (o) {
          var u = new s(),
            d = u.get,
            f = u.has,
            g = u.set;
          (t = function (e, n) {
            return g.call(u, e, n), n;
          }),
            (l = function (e) {
              return d.call(u, e) || {};
            }),
            (r = function (e) {
              return f.call(u, e);
            });
        } else {
          var v = h('state');
          (x[v] = !0),
            (t = function (e, n) {
              return p(e, v, n), n;
            }),
            (l = function (e) {
              return y(e, v) ? e[v] : {};
            }),
            (r = function (e) {
              return y(e, v);
            });
        }
        e.exports = {
          set: t,
          get: l,
          has: r,
          enforce: function (e) {
            return r(e) ? l(e) : t(e, {});
          },
          getterFor: function (e) {
            return function (n) {
              var i;
              if (!c(n) || (i = l(n)).type !== e)
                throw TypeError('Incompatible receiver, ' + e + ' required');
              return i;
            };
          },
        };
      },
      function (e, n, i) {
        'use strict';
        Object.defineProperty(n, '__esModule', { value: !0 });
        var t =
            Object.assign ||
            function (e) {
              for (var n = 1; n < arguments.length; n++) {
                var i = arguments[n];
                for (var t in i)
                  Object.prototype.hasOwnProperty.call(i, t) && (e[t] = i[t]);
              }
              return e;
            },
          l = o(i(22)),
          r = o(i(12));
        function o(e) {
          return e && e.__esModule ? e : { default: e };
        }
        n.default = function () {
          var e =
            arguments.length > 0 && void 0 !== arguments[0] ? arguments[0] : {};
          if ('undefined' == typeof document)
            throw new Error(
              '`feather.replace()` only works in a browser environment.'
            );
          var n = document.querySelectorAll('[data-feather]');
          Array.from(n).forEach(function (n) {
            return (function (e) {
              var n =
                  arguments.length > 1 && void 0 !== arguments[1]
                    ? arguments[1]
                    : {},
                i = (function (e) {
                  return Array.from(e.attributes).reduce(function (e, n) {
                    return (e[n.name] = n.value), e;
                  }, {});
                })(e),
                o = i['data-feather'];
              delete i['data-feather'];
              var a = r.default[o].toSvg(
                  t({}, n, i, { class: (0, l.default)(n.class, i.class) })
                ),
                c = new DOMParser()
                  .parseFromString(a, 'image/svg+xml')
                  .querySelector('svg');
              e.parentNode.replaceChild(c, e);
            })(n, e);
          });
        };
      },
      function (e, n, i) {
        'use strict';
        Object.defineProperty(n, '__esModule', { value: !0 });
        var t,
          l = i(12),
          r = (t = l) && t.__esModule ? t : { default: t };
        n.default = function (e) {
          var n =
            arguments.length > 1 && void 0 !== arguments[1] ? arguments[1] : {};
          if (
            (console.warn(
              'feather.toSvg() is deprecated. Please use feather.icons[name].toSvg() instead.'
            ),
            !e)
          )
            throw new Error(
              'The required `key` (icon name) parameter is missing.'
            );
          if (!r.default[e])
            throw new Error(
              "No icon matching '" +
                e +
                "'. See the complete list of icons at https://feathericons.com"
            );
          return r.default[e].toSvg(n);
        };
      },
      function (e) {
        e.exports = {
          plus: ['add', 'new'],
          'plus-circle': ['add', 'new'],
          'plus-square': ['add', 'new'],
          x: ['cancel', 'close', 'delete', 'remove', 'times', 'clear'],
        };
      },
      function (e) {
        e.exports = {
          'chevrons-down':
            '<polyline points="7 13 12 18 17 13"></polyline><polyline points="7 6 12 11 17 6"></polyline>',
          'chevrons-up':
            '<polyline points="17 11 12 6 7 11"></polyline><polyline points="17 18 12 13 7 18"></polyline>',
          plus:
            '<line x1="12" y1="5" x2="12" y2="19"></line><line x1="5" y1="12" x2="19" y2="12"></line>',
          x:
            '<line x1="18" y1="6" x2="6" y2="18"></line><line x1="6" y1="6" x2="18" y2="18"></line>',
        };
      },
      function (e) {
        e.exports = {
          xmlns: 'http://www.w3.org/2000/svg',
          width: 24,
          height: 24,
          viewBox: '0 0 24 24',
          fill: 'none',
          stroke: 'currentColor',
          'stroke-width': 2,
          'stroke-linecap': 'round',
          'stroke-linejoin': 'round',
        };
      },
      function (e, n, i) {
        'use strict';
        Object.defineProperty(n, '__esModule', { value: !0 });
        var t =
            Object.assign ||
            function (e) {
              for (var n = 1; n < arguments.length; n++) {
                var i = arguments[n];
                for (var t in i)
                  Object.prototype.hasOwnProperty.call(i, t) && (e[t] = i[t]);
              }
              return e;
            },
          l = (function () {
            function e(e, n) {
              for (var i = 0; i < n.length; i++) {
                var t = n[i];
                (t.enumerable = t.enumerable || !1),
                  (t.configurable = !0),
                  'value' in t && (t.writable = !0),
                  Object.defineProperty(e, t.key, t);
              }
            }
            return function (n, i, t) {
              return i && e(n.prototype, i), t && e(n, t), n;
            };
          })(),
          r = a(i(22)),
          o = a(i(42));
        function a(e) {
          return e && e.__esModule ? e : { default: e };
        }
        var c = (function () {
          function e(n, i) {
            var l =
              arguments.length > 2 && void 0 !== arguments[2] ? arguments[2] : [];
            !(function (e, n) {
              if (!(e instanceof n))
                throw new TypeError('Cannot call a class as a function');
            })(this, e),
              (this.name = n),
              (this.contents = i),
              (this.tags = l),
              (this.attrs = t({}, o.default, { class: 'feather feather-' + n }));
          }
          return (
            l(e, [
              {
                key: 'toSvg',
                value: function () {
                  var e =
                    arguments.length > 0 && void 0 !== arguments[0]
                      ? arguments[0]
                      : {};
                  return (
                    '<svg ' +
                    (function (e) {
                      return Object.keys(e)
                        .map(function (n) {
                          return n + '="' + e[n] + '"';
                        })
                        .join(' ');
                    })(
                      t({}, this.attrs, e, {
                        class: (0, r.default)(this.attrs.class, e.class),
                      })
                    ) +
                    '>' +
                    this.contents +
                    '</svg>'
                  );
                },
              },
              {
                key: 'toString',
                value: function () {
                  return this.contents;
                },
              },
            ]),
            e
          );
        })();
        n.default = c;
      },
      function (e, n, i) {
        'use strict';
        var t = o(i(12)),
          l = o(i(39)),
          r = o(i(38));
        function o(e) {
          return e && e.__esModule ? e : { default: e };
        }
        e.exports = { icons: t.default, toSvg: l.default, replace: r.default };
      },
      function (e, n, i) {
        e.exports = i(0);
      },
      function (e, n, i) {
        var t = i(2)('iterator'),
          l = !1;
        try {
          var r = 0,
            o = {
              next: function () {
                return { done: !!r++ };
              },
              return: function () {
                l = !0;
              },
            };
          (o[t] = function () {
            return this;
          }),
            Array.from(o, function () {
              throw 2;
            });
        } catch (e) {}
        e.exports = function (e, n) {
          if (!n && !l) return !1;
          var i = !1;
          try {
            var r = {};
            (r[t] = function () {
              return {
                next: function () {
                  return { done: (i = !0) };
                },
              };
            }),
              e(r);
          } catch (e) {}
          return i;
        };
      },
      function (e, n, i) {
        var t = i(30),
          l = i(2)('toStringTag'),
          r =
            'Arguments' ==
            t(
              (function () {
                return arguments;
              })()
            );
        e.exports = function (e) {
          var n, i, o;
          return void 0 === e
            ? 'Undefined'
            : null === e
            ? 'Null'
            : 'string' ==
              typeof (i = (function (e, n) {
                try {
                  return e[n];
                } catch (e) {}
              })((n = Object(e)), l))
            ? i
            : r
            ? t(n)
            : 'Object' == (o = t(n)) && 'function' == typeof n.callee
            ? 'Arguments'
            : o;
        };
      },
      function (e, n, i) {
        var t = i(47),
          l = i(9),
          r = i(2)('iterator');
        e.exports = function (e) {
          if (void 0 != e) return e[r] || e['@@iterator'] || l[t(e)];
        };
      },
      function (e, n, i) {
        'use strict';
        var t = i(18),
          l = i(7),
          r = i(10);
        e.exports = function (e, n, i) {
          var o = t(n);
          o in e ? l.f(e, o, r(0, i)) : (e[o] = i);
        };
      },
      function (e, n, i) {
        var t = i(2),
          l = i(9),
          r = t('iterator'),
          o = Array.prototype;
        e.exports = function (e) {
          return void 0 !== e && (l.Array === e || o[r] === e);
        };
      },
      function (e, n, i) {
        var t = i(3);
        e.exports = function (e, n, i, l) {
          try {
            return l ? n(t(i)[0], i[1]) : n(i);
          } catch (n) {
            var r = e.return;
            throw (void 0 !== r && t(r.call(e)), n);
          }
        };
      },
      function (e, n) {
        e.exports = function (e) {
          if ('function' != typeof e)
            throw TypeError(String(e) + ' is not a function');
          return e;
        };
      },
      function (e, n, i) {
        var t = i(52);
        e.exports = function (e, n, i) {
          if ((t(e), void 0 === n)) return e;
          switch (i) {
            case 0:
              return function () {
                return e.call(n);
              };
            case 1:
              return function (i) {
                return e.call(n, i);
              };
            case 2:
              return function (i, t) {
                return e.call(n, i, t);
              };
            case 3:
              return function (i, t, l) {
                return e.call(n, i, t, l);
              };
          }
          return function () {
            return e.apply(n, arguments);
          };
        };
      },
      function (e, n, i) {
        'use strict';
        var t = i(53),
          l = i(24),
          r = i(51),
          o = i(50),
          a = i(27),
          c = i(49),
          p = i(48);
        e.exports = function (e) {
          var n,
            i,
            y,
            h,
            x = l(e),
            s = 'function' == typeof this ? this : Array,
            u = arguments.length,
            d = u > 1 ? arguments[1] : void 0,
            f = void 0 !== d,
            g = 0,
            v = p(x);
          if (
            (f && (d = t(d, u > 2 ? arguments[2] : void 0, 2)),
            void 0 == v || (s == Array && o(v)))
          )
            for (i = new s((n = a(x.length))); n > g; g++)
              c(i, g, f ? d(x[g], g) : x[g]);
          else
            for (h = v.call(x), i = new s(); !(y = h.next()).done; g++)
              c(i, g, f ? r(h, d, [y.value, g], !0) : y.value);
          return (i.length = g), i;
        };
      },
      function (e, n, i) {
        var t = i(32),
          l = i(54);
        t(
          {
            target: 'Array',
            stat: !0,
            forced: !i(46)(function (e) {
              Array.from(e);
            }),
          },
          { from: l }
        );
      },
      function (e, n, i) {
        var t = i(6),
          l = i(3);
        e.exports = function (e, n) {
          if ((l(e), !t(n) && null !== n))
            throw TypeError("Can't set " + String(n) + ' as a prototype');
        };
      },
      function (e, n, i) {
        var t = i(56);
        e.exports =
          Object.setPrototypeOf ||
          ('__proto__' in {}
            ? (function () {
                var e,
                  n = !1,
                  i = {};
                try {
                  (e = Object.getOwnPropertyDescriptor(
                    Object.prototype,
                    '__proto__'
                  ).set).call(i, []),
                    (n = i instanceof Array);
                } catch (e) {}
                return function (i, l) {
                  return t(i, l), n ? e.call(i, l) : (i.__proto__ = l), i;
                };
              })()
            : void 0);
      },
      function (e, n, i) {
        var t = i(0).document;
        e.exports = t && t.documentElement;
      },
      function (e, n, i) {
        var t = i(28),
          l = i(13);
        e.exports =
          Object.keys ||
          function (e) {
            return t(e, l);
          };
      },
      function (e, n, i) {
        var t = i(8),
          l = i(7),
          r = i(3),
          o = i(59);
        e.exports = t
          ? Object.defineProperties
          : function (e, n) {
              r(e);
              for (var i, t = o(n), a = t.length, c = 0; a > c; )
                l.f(e, (i = t[c++]), n[i]);
              return e;
            };
      },
      function (e, n, i) {
        var t = i(3),
          l = i(60),
          r = i(13),
          o = i(15),
          a = i(58),
          c = i(34),
          p = i(16)('IE_PROTO'),
          y = function () {},
          h = function () {
            var e,
              n = c('iframe'),
              i = r.length;
            for (
              n.style.display = 'none',
                a.appendChild(n),
                n.src = String('javascript:'),
                (e = n.contentWindow.document).open(),
                e.write('<script>document.F=Object</script>'),
                e.close(),
                h = e.F;
              i--;
  
            )
              delete h.prototype[r[i]];
            return h();
          };
        (e.exports =
          Object.create ||
          function (e, n) {
            var i;
            return (
              null !== e
                ? ((y.prototype = t(e)),
                  (i = new y()),
                  (y.prototype = null),
                  (i[p] = e))
                : (i = h()),
              void 0 === n ? i : l(i, n)
            );
          }),
          (o[p] = !0);
      },
      function (e, n, i) {
        var t = i(4);
        e.exports =
          !!Object.getOwnPropertySymbols &&
          !t(function () {
            return !String(Symbol());
          });
      },
      function (e, n, i) {
        var t = i(4);
        e.exports = !t(function () {
          function e() {}
          return (
            (e.prototype.constructor = null),
            Object.getPrototypeOf(new e()) !== e.prototype
          );
        });
      },
      function (e, n, i) {
        'use strict';
        var t = i(26).IteratorPrototype,
          l = i(61),
          r = i(10),
          o = i(23),
          a = i(9),
          c = function () {
            return this;
          };
        e.exports = function (e, n, i) {
          var p = n + ' Iterator';
          return (
            (e.prototype = l(t, { next: r(1, i) })),
            o(e, p, !1, !0),
            (a[p] = c),
            e
          );
        };
      },
      function (e, n, i) {
        var t = i(4),
          l = /#|\.prototype\./,
          r = function (e, n) {
            var i = a[o(e)];
            return i == p || (i != c && ('function' == typeof n ? t(n) : !!n));
          },
          o = (r.normalize = function (e) {
            return String(e).replace(l, '.').toLowerCase();
          }),
          a = (r.data = {}),
          c = (r.NATIVE = 'N'),
          p = (r.POLYFILL = 'P');
        e.exports = r;
      },
      function (e, n) {
        n.f = Object.getOwnPropertySymbols;
      },
      function (e, n, i) {
        var t = i(21),
          l = Math.max,
          r = Math.min;
        e.exports = function (e, n) {
          var i = t(e);
          return i < 0 ? l(i + n, 0) : r(i, n);
        };
      },
      function (e, n, i) {
        var t = i(14),
          l = i(27),
          r = i(67);
        e.exports = function (e) {
          return function (n, i, o) {
            var a,
              c = t(n),
              p = l(c.length),
              y = r(o, p);
            if (e && i != i) {
              for (; p > y; ) if ((a = c[y++]) != a) return !0;
            } else
              for (; p > y; y++)
                if ((e || y in c) && c[y] === i) return e || y || 0;
            return !e && -1;
          };
        };
      },
      function (e, n, i) {
        var t = i(28),
          l = i(13).concat('length', 'prototype');
        n.f =
          Object.getOwnPropertyNames ||
          function (e) {
            return t(e, l);
          };
      },
      function (e, n, i) {
        var t = i(0),
          l = i(69),
          r = i(66),
          o = i(3),
          a = t.Reflect;
        e.exports =
          (a && a.ownKeys) ||
          function (e) {
            var n = l.f(o(e)),
              i = r.f;
            return i ? n.concat(i(e)) : n;
          };
      },
      function (e, n, i) {
        var t = i(1),
          l = i(70),
          r = i(31),
          o = i(7);
        e.exports = function (e, n) {
          for (var i = l(n), a = o.f, c = r.f, p = 0; p < i.length; p++) {
            var y = i[p];
            t(e, y) || a(e, y, c(n, y));
          }
        };
      },
      function (e, n, i) {
        var t = i(4),
          l = i(30),
          r = ''.split;
        e.exports = t(function () {
          return !Object('z').propertyIsEnumerable(0);
        })
          ? function (e) {
              return 'String' == l(e) ? r.call(e, '') : Object(e);
            }
          : Object;
      },
      function (e, n, i) {
        'use strict';
        var t = {}.propertyIsEnumerable,
          l = Object.getOwnPropertyDescriptor,
          r = l && !t.call({ 1: 2 }, 1);
        n.f = r
          ? function (e) {
              var n = l(this, e);
              return !!n && n.enumerable;
            }
          : t;
      },
      function (e, n, i) {
        'use strict';
        var t = i(32),
          l = i(64),
          r = i(25),
          o = i(57),
          a = i(23),
          c = i(5),
          p = i(29),
          y = i(2),
          h = i(17),
          x = i(9),
          s = i(26),
          u = s.IteratorPrototype,
          d = s.BUGGY_SAFARI_ITERATORS,
          f = y('iterator'),
          g = function () {
            return this;
          };
        e.exports = function (e, n, i, y, s, v, m) {
          l(i, n, y);
          var w,
            M,
            b,
            z = function (e) {
              if (e === s && O) return O;
              if (!d && e in H) return H[e];
              switch (e) {
                case 'keys':
                case 'values':
                case 'entries':
                  return function () {
                    return new i(this, e);
                  };
              }
              return function () {
                return new i(this);
              };
            },
            A = n + ' Iterator',
            k = !1,
            H = e.prototype,
            V = H[f] || H['@@iterator'] || (s && H[s]),
            O = (!d && V) || z(s),
            j = ('Array' == n && H.entries) || V;
          if (
            (j &&
              ((w = r(j.call(new e()))),
              u !== Object.prototype &&
                w.next &&
                (h ||
                  r(w) === u ||
                  (o ? o(w, u) : 'function' != typeof w[f] && c(w, f, g)),
                a(w, A, !0, !0),
                h && (x[A] = g))),
            'values' == s &&
              V &&
              'values' !== V.name &&
              ((k = !0),
              (O = function () {
                return V.call(this);
              })),
            (h && !m) || H[f] === O || c(H, f, O),
            (x[n] = O),
            s)
          )
            if (
              ((M = {
                values: z('values'),
                keys: v ? O : z('keys'),
                entries: z('entries'),
              }),
              m)
            )
              for (b in M) (!d && !k && b in H) || p(H, b, M[b]);
            else t({ target: n, proto: !0, forced: d || k }, M);
          return M;
        };
      },
      function (e, n) {
        var i;
        i = (function () {
          return this;
        })();
        try {
          i = i || Function('return this')() || (0, eval)('this');
        } catch (e) {
          'object' == typeof window && (i = window);
        }
        e.exports = i;
      },
      function (e, n, i) {
        var t = i(0),
          l = i(36),
          r = t.WeakMap;
        e.exports = 'function' == typeof r && /native code/.test(l.call(r));
      },
      function (e, n, i) {
        var t = i(21),
          l = i(20);
        e.exports = function (e, n, i) {
          var r,
            o,
            a = String(l(e)),
            c = t(n),
            p = a.length;
          return c < 0 || c >= p
            ? i
              ? ''
              : void 0
            : (r = a.charCodeAt(c)) < 55296 ||
              r > 56319 ||
              c + 1 === p ||
              (o = a.charCodeAt(c + 1)) < 56320 ||
              o > 57343
            ? i
              ? a.charAt(c)
              : r
            : i
            ? a.slice(c, c + 2)
            : o - 56320 + ((r - 55296) << 10) + 65536;
        };
      },
      function (e, n, i) {
        'use strict';
        var t = i(77),
          l = i(37),
          r = i(74),
          o = l.set,
          a = l.getterFor('String Iterator');
        r(
          String,
          'String',
          function (e) {
            o(this, { type: 'String Iterator', string: String(e), index: 0 });
          },
          function () {
            var e,
              n = a(this),
              i = n.string,
              l = n.index;
            return l >= i.length
              ? { value: void 0, done: !0 }
              : ((e = t(i, l, !0)),
                (n.index += e.length),
                { value: e, done: !1 });
          }
        );
      },
      function (e, n, i) {
        i(78), i(55);
        var t = i(45);
        e.exports = t.Array.from;
      },
      function (e, n, i) {
        i(79), (e.exports = i(44));
      },
    ]);
  });