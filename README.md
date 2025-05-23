<div align="center">
 <a href="https://www.speakeasy.com/" target="_blank">
   <picture>
       <source media="(prefers-color-scheme: light)" srcset="https://github.com/user-attachments/assets/21dd5d3a-aefc-4cd3-abee-5e17ef1d4dad">
       <source media="(prefers-color-scheme: dark)" srcset="https://github.com/user-attachments/assets/0a747f98-d228-462d-9964-fd87bf93adc5">
       <img width="100px" src="https://github.com/user-attachments/assets/21dd5d3a-aefc-4cd3-abee-5e17ef1d4dad#gh-light-mode-only" alt="Speakeasy">
   </picture>
 </a>
  <h1>Speakeasy</h1>
  <p>Build APIs your users love ❤️ with Speakeasy</p>
  <div>
   <a href="https://speakeasy.com/docs/create-client-sdks/" target="_blank"><b>Docs Quickstart</b></a>&nbsp;&nbsp;//&nbsp;&nbsp;<a href="https://join.slack.com/t/speakeasy-dev/shared_invite/zt-1cwb3flxz-lS5SyZxAsF_3NOq5xc8Cjw" target="_blank"><b>Join us on Slack</b></a>
  </div>
 <br />

</div>

# jsonpath

<a href="https://pkg.go.dev/github.com/speakeasy-api/jsonpath?tab=doc"><img alt="Go Doc" src="https://img.shields.io/badge/godoc-reference-blue.svg?style=for-the-badge"></a>

This is a full implementation of [RFC 9535](https://datatracker.ietf.org/doc/rfc9535/)

It is build to be wasm compatible. A playground application is available at [overlay.speakeasy.com](https://overlay.speakeasy.com/)

Everything within RFC9535 is in scope. Grammars outside RFC 9535 are not in scope.

## Installation

This application is included in the [speakeasy](https://github.com/speakeasy-api/speakeasy) CLI, but is also available as a standalone library.

## ABNF grammar

```
   jsonpath-query      = root-identifier segments
   segments            = *(S segment)

   B                   = %x20 /    ; Space
                         %x09 /    ; Horizontal tab
                         %x0A /    ; Line feed or New line
                         %x0D      ; Carriage return
   S                   = *B        ; optional blank space
   root-identifier     = "$"
   selector            = name-selector /
                         wildcard-selector /
                         slice-selector /
                         index-selector /
                         filter-selector
   name-selector       = string-literal

   string-literal      = %x22 *double-quoted %x22 /     ; "string"
                         %x27 *single-quoted %x27       ; 'string'

   double-quoted       = unescaped /
                         %x27      /                    ; '
                         ESC %x22  /                    ; \"
                         ESC escapable

   single-quoted       = unescaped /
                         %x22      /                    ; "
                         ESC %x27  /                    ; \'
                         ESC escapable

   ESC                 = %x5C                           ; \ backslash

   unescaped           = %x20-21 /                      ; see RFC 8259
                            ; omit 0x22 "
                         %x23-26 /
                            ; omit 0x27 '
                         %x28-5B /
                            ; omit 0x5C \
                         %x5D-D7FF /
                            ; skip surrogate code points
                         %xE000-10FFFF

   escapable           = %x62 / ; b BS backspace U+0008
                         %x66 / ; f FF form feed U+000C
                         %x6E / ; n LF line feed U+000A
                         %x72 / ; r CR carriage return U+000D
                         %x74 / ; t HT horizontal tab U+0009
                         "/"  / ; / slash (solidus) U+002F
                         "\"  / ; \ backslash (reverse solidus) U+005C
                         (%x75 hexchar) ;  uXXXX U+XXXX

   hexchar             = non-surrogate /
                         (high-surrogate "\" %x75 low-surrogate)
   non-surrogate       = ((DIGIT / "A"/"B"/"C" / "E"/"F") 3HEXDIG) /
                         ("D" %x30-37 2HEXDIG )
   high-surrogate      = "D" ("8"/"9"/"A"/"B") 2HEXDIG
   low-surrogate       = "D" ("C"/"D"/"E"/"F") 2HEXDIG

   HEXDIG              = DIGIT / "A" / "B" / "C" / "D" / "E" / "F"
   wildcard-selector   = "*"
   index-selector      = int                        ; decimal integer

   int                 = "0" /
                         (["-"] DIGIT1 *DIGIT)      ; - optional
   DIGIT1              = %x31-39                    ; 1-9 non-zero digit
   slice-selector      = [start S] ":" S [end S] [":" [S step ]]

   start               = int       ; included in selection
   end                 = int       ; not included in selection
   step                = int       ; default: 1
   filter-selector     = "?" S logical-expr
   logical-expr        = logical-or-expr
   logical-or-expr     = logical-and-expr *(S "||" S logical-and-expr)
                           ; disjunction
                           ; binds less tightly than conjunction
   logical-and-expr    = basic-expr *(S "&&" S basic-expr)
                           ; conjunction
                           ; binds more tightly than disjunction

   basic-expr          = paren-expr /
                         comparison-expr /
                         test-expr

   paren-expr          = [logical-not-op S] "(" S logical-expr S ")"
                                           ; parenthesized expression
   logical-not-op      = "!"               ; logical NOT operator
   test-expr           = [logical-not-op S]
                         (filter-query / ; existence/non-existence
                          function-expr) ; LogicalType or NodesType
   filter-query        = rel-query / jsonpath-query
   rel-query           = current-node-identifier segments
   current-node-identifier = "@"
   comparison-expr     = comparable S comparison-op S comparable
   literal             = number / string-literal /
                         true / false / null
   comparable          = literal /
                         singular-query / ; singular query value
                         function-expr    ; ValueType
   comparison-op       = "==" / "!=" /
                         "<=" / ">=" /
                         "<"  / ">"

   singular-query      = rel-singular-query / abs-singular-query
   rel-singular-query  = current-node-identifier singular-query-segments
   abs-singular-query  = root-identifier singular-query-segments
   singular-query-segments = *(S (name-segment / index-segment))
   name-segment        = ("[" name-selector "]") /
                         ("." member-name-shorthand)
   index-segment       = "[" index-selector "]"
   number              = (int / "-0") [ frac ] [ exp ] ; decimal number
   frac                = "." 1*DIGIT                  ; decimal fraction
   exp                 = "e" [ "-" / "+" ] 1*DIGIT    ; decimal exponent
   true                = %x74.72.75.65                ; true
   false               = %x66.61.6c.73.65             ; false
   null                = %x6e.75.6c.6c                ; null
   function-name       = function-name-first *function-name-char
   function-name-first = LCALPHA
   function-name-char  = function-name-first / "_" / DIGIT
   LCALPHA             = %x61-7A  ; "a".."z"

   function-expr       = function-name "(" S [function-argument
                            *(S "," S function-argument)] S ")"
   function-argument   = literal /
                         filter-query / ; (includes singular-query)
                         logical-expr /
                         function-expr
   segment             = child-segment / descendant-segment
   child-segment       = bracketed-selection /
                         ("."
                          (wildcard-selector /
                           member-name-shorthand))

   bracketed-selection = "[" S selector *(S "," S selector) S "]"

   member-name-shorthand = name-first *name-char
   name-first          = ALPHA /
                         "_"   /
                         %x80-D7FF /
                            ; skip surrogate code points
                         %xE000-10FFFF
   name-char           = name-first / DIGIT

   DIGIT               = %x30-39              ; 0-9
   ALPHA               = %x41-5A / %x61-7A    ; A-Z / a-z
   descendant-segment  = ".." (bracketed-selection /
                               wildcard-selector /
                               member-name-shorthand)

                Figure 2: Collected ABNF of JSONPath Queries

   Figure 3 contains the collected ABNF grammar that defines the syntax
   of a JSONPath Normalized Path while also using the rules root-
   identifier, ESC, DIGIT, and DIGIT1 from Figure 2.

   normalized-path      = root-identifier *(normal-index-segment)
   normal-index-segment = "[" normal-selector "]"
   normal-selector      = normal-name-selector / normal-index-selector
   normal-name-selector = %x27 *normal-single-quoted %x27 ; 'string'
   normal-single-quoted = normal-unescaped /
                          ESC normal-escapable
   normal-unescaped     =    ; omit %x0-1F control codes
                          %x20-26 /
                             ; omit 0x27 '
                          %x28-5B /
                             ; omit 0x5C \
                          %x5D-D7FF /
                             ; skip surrogate code points
                          %xE000-10FFFF

   normal-escapable     = %x62 / ; b BS backspace U+0008
                          %x66 / ; f FF form feed U+000C
                          %x6E / ; n LF line feed U+000A
                          %x72 / ; r CR carriage return U+000D
                          %x74 / ; t HT horizontal tab U+0009
                          "'" /  ; ' apostrophe U+0027
                          "\" /  ; \ backslash (reverse solidus) U+005C
                          (%x75 normal-hexchar)
                                          ; certain values u00xx U+00XX
   normal-hexchar       = "0" "0"
                          (
                             ("0" %x30-37) / ; "00"-"07"
                                ; omit U+0008-U+000A BS HT LF
                             ("0" %x62) /    ; "0b"
                                ; omit U+000C-U+000D FF CR
                             ("0" %x65-66) / ; "0e"-"0f"
                             ("1" normal-HEXDIG)
                          )
   normal-HEXDIG        = DIGIT / %x61-66    ; "0"-"9", "a"-"f"
   normal-index-selector = "0" / (DIGIT1 *DIGIT)
                           ; non-negative decimal integer
```

## Contributing

We welcome contributions to this repository! Please open a Github issue or a Pull Request if you have an implementation for a bug fix or feature. This repository is compliant with the [jsonpath standard compliance test suite](https://github.com/jsonpath-standard/jsonpath-compliance-test-suite/tree/9277705cda4489c3d0d984831e7656e48145399b) 
